package admin

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"

	"github.com/forgego/forge/pkg/admin/templates"
	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/query"
)

// ListViewData contains data for the admin list view
type ListViewData struct {
	Model       *AdminModel
	Objects     []interface{}
	Page        int
	PageSize    int
	TotalCount  int64
	TotalPages  int
	SearchQuery string
	BaseURL     string
}

// handleModelList handles the admin list view
func handleModelList(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canView(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to view this list.", http.StatusForbidden)
			return
		}

		ctx := r.Context()

		// Get pagination parameters
		page := httplib.GetQueryInt(r, "page", 1)
		listPerPage := GetListPerPage(model)
		pageSize := httplib.GetQueryInt(r, "page_size", listPerPage)
		searchQuery := httplib.GetQueryString(r, "search", "")
		
		// Get date hierarchy parameters
		dateHierarchy := GetDateHierarchy(model)
		year := httplib.GetQueryInt(r, "year", 0)
		month := httplib.GetQueryInt(r, "month", 0)
		day := httplib.GetQueryInt(r, "day", 0)

		// Get manager - try to use generic manager if available
		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		// Try to use generic manager operations
		// This is a type-safe path when manager is *query.Manager[T]
		if err := handleListWithGenericManager(w, r, ctx, modelName, model, manager, page, pageSize, searchQuery, dateHierarchy, year, month, day); err != nil {
			// Fallback to reflection-based approach
			if err := handleListWithReflection(w, r, ctx, modelName, model, manager, page, pageSize, searchQuery, dateHierarchy, year, month, day); err != nil {
				http.Error(w, fmt.Sprintf("Failed to process list: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

// handleListWithGenericManager attempts to use generic manager operations
func handleListWithGenericManager(w http.ResponseWriter, r *http.Request, ctx context.Context, modelName string, model *AdminModel, manager interface{}, page, pageSize int, searchQuery string, dateHierarchy string, year, month, day int) error {
	// Try to type assert to *query.Manager[T] for common types
	// This is a workaround since we can't know T at compile time
	// We'll use a type switch for known manager types
	
	// For now, return error to fallback to reflection
	// In a full implementation, we'd register managers with their types
	return fmt.Errorf("generic manager not available, using reflection fallback")
}

// handleListWithReflection is the fallback reflection-based implementation
func handleListWithReflection(w http.ResponseWriter, r *http.Request, ctx context.Context, modelName string, model *AdminModel, manager interface{}, page, pageSize int, searchQuery string, dateHierarchy string, year, month, day int) error {
	// Get manager's Filter method to get base queryset
	managerValue := reflect.ValueOf(manager)
	filterMethod := managerValue.MethodByName("Filter")
	if !filterMethod.IsValid() {
		return fmt.Errorf("manager does not have Filter method")
	}

	// Call Filter() to get base queryset (returns QuerySet)
	qsResults := filterMethod.Call([]reflect.Value{})
	if len(qsResults) == 0 {
		return fmt.Errorf("failed to create queryset")
	}
	qs := qsResults[0]
	
	// Apply custom queryset if available
	if model.ExtendedConfig != nil {
		if getQueryset, ok := model.ExtendedConfig["get_queryset"].(func(context.Context, interface{}) (interface{}, error)); ok {
			if customQs, err := getQueryset(ctx, manager); err == nil {
				qs = reflect.ValueOf(customQs)
			}
		}
	}
	
	// Apply date hierarchy
	if dateHierarchy != "" && (year > 0 || month > 0 || day > 0) {
		qs = ApplyDateHierarchy(qs, dateHierarchy, year, month, day)
	}
	
	// Apply ordering
	ordering := GetOrdering(model)
	if len(ordering) > 0 {
		qs = ApplyOrdering(qs, ordering)
	}
	
	// Apply filters
	qs = applyFilters(r, model, qs)

	// Apply search if provided
	if searchQuery != "" && len(model.SearchFields) > 0 {
		// Search across all search fields (OR condition)
		var searchExprs []query.QueryExpr
		for _, searchField := range model.SearchFields {
			if fieldName, ok := searchField.(string); ok {
				// Create a LIKE query for search
				searchExpr := query.NewFieldQueryExpr(fieldName, query.OpContains, searchQuery)
				searchExprs = append(searchExprs, searchExpr)
			}
		}
		
		// Combine with OR if multiple fields
		if len(searchExprs) > 0 {
			combinedExpr := searchExprs[0]
			for i := 1; i < len(searchExprs); i++ {
				combinedExpr = combinedExpr.Or(searchExprs[i])
			}
			
			filterMethod := qs.MethodByName("Filter")
			if filterMethod.IsValid() {
				qs = filterMethod.Call([]reflect.Value{reflect.ValueOf(combinedExpr)})[0]
			}
		}
	}

	// Get total count
	countMethod := qs.MethodByName("Count")
	var totalCount int64
	if countMethod.IsValid() {
		countResults := countMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
		if len(countResults) >= 2 {
			if errVal := countResults[1]; !errVal.IsNil() {
				if err, ok := errVal.Interface().(error); ok {
					return fmt.Errorf("failed to count objects: %w", err)
				}
			}
			if countVal, ok := countResults[0].Interface().(int64); ok {
				totalCount = countVal
			}
		}
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	offsetMethod := qs.MethodByName("Offset")
	if offsetMethod.IsValid() {
		qs = offsetMethod.Call([]reflect.Value{reflect.ValueOf(offset)})[0]
	}
	limitMethod := qs.MethodByName("Limit")
	if limitMethod.IsValid() {
		qs = limitMethod.Call([]reflect.Value{reflect.ValueOf(pageSize)})[0]
	}

	// Get objects
	allMethod := qs.MethodByName("All")
	var objectsInterface []interface{}
	if allMethod.IsValid() {
		allResults := allMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
		if len(allResults) >= 2 {
			if errVal := allResults[1]; !errVal.IsNil() {
				if err, ok := errVal.Interface().(error); ok {
					return fmt.Errorf("failed to get objects: %w", err)
				}
			}
			// Convert slice to []interface{}
			if allResults[0].IsValid() && allResults[0].Kind() == reflect.Slice {
				sliceValue := allResults[0]
				objectsInterface = make([]interface{}, sliceValue.Len())
				for i := 0; i < sliceValue.Len(); i++ {
					objectsInterface[i] = sliceValue.Index(i).Interface()
				}
			}
		}
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	// Prepare template data
	data := ListViewData{
		Model:       model,
		Objects:     objectsInterface,
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		SearchQuery: searchQuery,
		BaseURL:     fmt.Sprintf("/admin/%s/", modelName),
	}

	// Render template
	if err := renderListTemplate(w, r, data); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}

// getManagerForModel gets the manager for a model from the admin registry
func getManagerForModel(model *AdminModel) (interface{}, error) {
	// Check if manager is set
	if model.Manager == nil {
		return nil, fmt.Errorf("manager not set for model %s - use admin.RegisterModelWithManager()", model.Name)
	}
	
	return model.Manager, nil
}

// renderListTemplate renders the admin list template
func renderListTemplate(w http.ResponseWriter, r *http.Request, data ListViewData) error {
	// Load templates
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Prepare template data with field values and permissions
	type ObjectWithFields struct {
		ID      interface{}
		Fields  map[string]interface{}
		Actions map[string]string
		CanEdit bool
		CanDelete bool
	}

	var objectsWithFields []ObjectWithFields
	modelName := data.Model.Name
	for _, obj := range data.Objects {
		objFields := extractFieldValues(obj, data.Model)
		objID := getIDFromObject(obj)
		objectsWithFields = append(objectsWithFields, ObjectWithFields{
			ID:     objID,
			Fields: objFields,
			Actions: map[string]string{
				"view": fmt.Sprintf("%s%v/", data.BaseURL, objID),
				"edit": fmt.Sprintf("%s%v/change/", data.BaseURL, objID),
			},
			CanEdit:   canChange(r, modelName),
			CanDelete: canDelete(r, modelName),
		})
	}

	// Get list display fields
	listDisplay := data.Model.ListDisplay
	if len(listDisplay) == 0 {
		// Default: use all fields from first object
		if len(objectsWithFields) > 0 {
			for fieldName := range objectsWithFields[0].Fields {
				listDisplay = append(listDisplay, fieldName)
			}
		}
	}

		// Check if this is an HTMX request (only return table container)
		isHTMXRequest := r.Header.Get("HX-Request") == "true"
		
		if isHTMXRequest {
			// Return only the table container for HTMX swaps
			templateData := map[string]interface{}{
				"Model":      map[string]interface{}{
					"Name":       data.Model.Name,
					"ListDisplay": listDisplay,
				},
				"Objects":    objectsWithFields,
				"Page":       data.Page,
				"PageSize":   data.PageSize,
				"TotalCount": data.TotalCount,
				"TotalPages": data.TotalPages,
				"SearchQuery": data.SearchQuery,
				"BaseURL":    data.BaseURL,
				"CanAdd":     canAdd(r, modelName),
				"CanChange":  canChange(r, modelName),
				"CanDelete":  canDelete(r, modelName),
			}
			
			// Use embedded partial template
			return tmpl.ExecuteTemplate(w, "partial", templateData)
		}

		// Get filter fields (get manager from model)
		var filterFields []FilterField
		if data.Model.Manager != nil {
			filterFields = getFilterFields(r, data.Model, data.Model.Manager)
		}
		
		// Convert filters to template format
		filterTemplates := make([]map[string]interface{}, len(filterFields))
		for i, filter := range filterFields {
			choices := make([]map[string]interface{}, len(filter.Options))
			for j, opt := range filter.Options {
				// Build query string for filter choice
				queryString := fmt.Sprintf("%s?%s=%v", data.BaseURL, filter.Name, opt.Value)
				choices[j] = map[string]interface{}{
					"Display":     opt.Label,
					"QueryString": queryString,
					"Selected":    filter.Active == opt.Value,
				}
			}
			filterTemplates[i] = map[string]interface{}{
				"Title":   filter.Label,
				"Open":    filter.Active != nil,
				"Choices": choices,
			}
		}
		
		// Generate pagination data
		paginationData := generatePaginationData(data)
		
		// Generate actions data (if bulk actions are enabled)
		actionsData := generateActionsData(r, data, modelName)
		
		// Full page render
		templateData := map[string]interface{}{
			"Title":      fmt.Sprintf("%s List", data.Model.Name),
			"Model":      map[string]interface{}{
				"Name":       data.Model.Name,
				"ListDisplay": listDisplay,
			},
			"Objects":    objectsWithFields,
			"Page":       data.Page,
			"PageSize":   data.PageSize,
			"TotalCount": data.TotalCount,
			"TotalPages": data.TotalPages,
			"SearchQuery": data.SearchQuery,
			"BaseURL":    data.BaseURL,
			"Models":     GetAllModels(), // For navigation
			"CanAdd":     canAdd(r, modelName),
			"CanChange":  canChange(r, modelName),
			"CanDelete":  canDelete(r, modelName),
			"Filters":    filterTemplates,
			"Pagination": paginationData,
			"Actions":    actionsData,
		}

		// Execute list template
		return tmpl.ExecuteTemplate(w, "list", templateData)
}

// extractFieldValues extracts field values from an object using reflection
func extractFieldValues(obj interface{}, model *AdminModel) map[string]interface{} {
	fields := make(map[string]interface{})
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	// Get list display fields
	listDisplay := model.ListDisplay
	if len(listDisplay) == 0 {
		// Default: show all exported fields
		typ := objValue.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() && field.Name != "BaseSchema" {
				fieldValue := objValue.Field(i)
				if fieldValue.CanInterface() {
					fields[field.Name] = fieldValue.Interface()
				}
			}
		}
	} else {
		// Show only specified fields
		for _, fieldName := range listDisplay {
			if fieldStr, ok := fieldName.(string); ok {
				fieldValue := objValue.FieldByName(fieldStr)
				if fieldValue.IsValid() && fieldValue.CanInterface() {
					fields[fieldStr] = fieldValue.Interface()
				}
			}
		}
	}

	return fields
}

// getIDFromObject extracts ID from an object
func getIDFromObject(obj interface{}) interface{} {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	idField := objValue.FieldByName("ID")
	if idField.IsValid() && idField.CanInterface() {
		return idField.Interface()
	}

	return nil
}

// generatePaginationData generates pagination data for template
func generatePaginationData(data ListViewData) map[string]interface{} {
	// Generate page range (show up to 10 pages around current)
	pageRange := make([]int, 0)
	startPage := data.Page - 5
	if startPage < 1 {
		startPage = 1
	}
	endPage := startPage + 9
	if endPage > data.TotalPages {
		endPage = data.TotalPages
		startPage = endPage - 9
		if startPage < 1 {
			startPage = 1
		}
	}
	for i := startPage; i <= endPage; i++ {
		pageRange = append(pageRange, i)
	}
	
	// Build query string for search
	searchQuery := ""
	if data.SearchQuery != "" {
		searchQuery = data.SearchQuery
	}
	
	return map[string]interface{}{
		"PaginationRequired": data.TotalPages > 1,
		"PageRange":          pageRange,
		"CurrentPage":        data.Page,
		"TotalPages":         data.TotalPages,
		"ResultCount":        data.TotalCount,
		"ModelName":          data.Model.Name,
		"ModelNamePlural":    data.Model.Name + "s", // Simplified
		"BaseURL":            data.BaseURL,
		"SearchQuery":        searchQuery,
		"ShowAllURL":         fmt.Sprintf("%s?all=1", data.BaseURL),
	}
}

// generateActionsData generates bulk actions data for template
func generateActionsData(r *http.Request, data ListViewData, modelName string) map[string]interface{} {
	// Check if bulk actions are enabled
	actions := GetActions(data.Model)
	if len(actions) == 0 || !canChange(r, modelName) {
		return nil
	}
	
	// Build action form select options
	actionOptions := make([]map[string]interface{}, len(actions))
	for i, action := range actions {
		label := action.Label
		if label == "" {
			label = action.Description
		}
		if label == "" {
			label = action.Name
		}
		actionOptions[i] = map[string]interface{}{
			"Value": action.Name,
			"Label": label,
		}
	}
	
	resultCount := int64(len(data.Objects))
	
	// Build select options HTML
	var selectOptions strings.Builder
	selectOptions.WriteString(`<option value="">---------</option>`)
	for _, action := range actions {
		label := action.Label
		if label == "" {
			label = action.Description
		}
		if label == "" {
			label = action.Name
		}
		selectOptions.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, action.Name, label))
	}
	
	return map[string]interface{}{
		"ActionForm": map[string]interface{}{
			"Label":  "Action:",
			"Select": template.HTML(fmt.Sprintf(`<select name="action">%s</select>`, selectOptions.String())),
		},
		"ActionIndex":            0,
		"ActionsSelectionCounter": resultCount > 0,
		"SelectionNote":          fmt.Sprintf("0 of %d selected", resultCount),
		"SelectionNoteAll":       fmt.Sprintf("All %d selected", data.TotalCount),
		"TotalCount":             data.TotalCount,
		"ResultCount":            resultCount,
		"ModuleName":             modelName,
	}
}

