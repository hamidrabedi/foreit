package views

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/forgego/forge/admin"
	adminorm "github.com/forgego/forge/admin/orm"
	adminschema "github.com/forgego/forge/admin/schema"
	"github.com/forgego/forge/orm"
)

// ListView represents a type-safe list view that integrates with ORM and filters
type ListView[T any] struct {
	*BaseView[T]
}

// NewListView creates a new list view
func NewListView[T any](admin *admin.Admin[T]) *ListView[T] {
	return &ListView[T]{
		BaseView: NewBaseView(admin),
	}
}

// ListData contains data for rendering the list view
type ListData[T any] struct {
	Objects             []*T
	Rows                []ListRowData
	Page                int
	PageSize            int
	TotalCount          int64
	TotalPages          int
	Search              string
	SearchFields        []string
	DisplayFields       []ListFieldHeader
	EditableFields      []string
	DisplayLinks        []string
	DateHierarchy       string
	EmptyValueDisplay   string
	Filters             []FilterData
	DateHierarchyData   *DateHierarchyData
	Actions             []ActionData
	HasAddPermission    bool
	HasViewPermission   bool
	HasChangePermission bool
	HasDeletePermission bool
	ModelName           string
	AddURL              string
	GroupBy             string              // Field to group by
	GroupedRows         []GroupedRowData    // Grouped rows if grouping is enabled
	Aggregates          []AggregateData     // Aggregate summaries
}

type ListFieldHeader struct {
	Name       string
	Label      string
	Sortable   bool
	Sorted     bool   // Whether this field is currently sorted
	SortOrder  string // "asc" or "desc"
	SortURL    string // URL to sort by this field
}

// GroupedRowData represents a grouped row
type GroupedRowData struct {
	GroupValue   interface{}
	GroupLabel   string
	Rows         []ListRowData
	Subtotal     map[string]AggregateValue
	IsGroup      bool
	IsSubtotal   bool
}

// AggregateData represents aggregate information
type AggregateData struct {
	Field     string
	Label     string
	Function  string // "sum", "count", "avg", "min", "max"
	Value     AggregateValue
}

// AggregateValue represents an aggregate value
type AggregateValue struct {
	Type  string      // "number", "string", "count"
	Value interface{}
}

type ListRowData struct {
	ID        interface{}
	Fields    []ListFieldVal
	ChangeURL string
	DeleteURL string
}

type ListFieldVal struct {
	Value  interface{}
	IsLink bool
	URL    string
}

type ActionData struct {
	Name  string
	Label string
}

// FilterData contains data for a filter
type FilterData struct {
	Name    string
	Label   string
	Type    string // "choices", "boolean", "date"
	Value   interface{}
	Options []FilterOption
}

// FilterOption contains data for a filter option
type FilterOption struct {
	Label    string
	Value    interface{}
	Selected bool
}

// DateHierarchyData contains data for date hierarchy navigation
type DateHierarchyData struct {
	FieldName string
	Years     []string
	Months    []string
	Days      []string
	Current   map[string]string // {year, month, day}
}

// Render renders the list view and returns the data
func (lv *ListView[T]) Render(ctx context.Context, r *http.Request, user interface{}) (*ListData[T], error) {
	// Check view permission
	if !lv.admin.HasViewPermission(ctx, user, nil) {
		return nil, fmt.Errorf("permission denied")
	}

	// Check for ChangelistViewHook
	config := lv.admin.Config()
	if config != nil && config.ChangelistViewHook != nil {
		customView, err := config.ChangelistViewHook(ctx, lv.admin, r)
		if err != nil {
			return nil, fmt.Errorf("changelist view hook error: %w", err)
		}
		if customView != nil {
			// Type assert to ListView[T]
			if lv, ok := customView.(*ListView[T]); ok {
				return lv.Render(ctx, r, user)
			}
		}
	}

	// Get base queryset
	qs, err := lv.GetQueryset(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queryset: %w", err)
	}

	// Apply filters using FilterSet
	filteredQs, err := lv.ApplyFilters(ctx, r, qs)
	if err != nil {
		return nil, fmt.Errorf("failed to apply filters: %w", err)
	}

	// Apply date hierarchy filter
	dateHierarchy := lv.getDateHierarchy(config)
	if dateHierarchy != "" {
		filteredQs, err = lv.applyDateHierarchy(r, filteredQs, dateHierarchy)
		if err != nil {
			return nil, fmt.Errorf("failed to apply date hierarchy: %w", err)
		}
	}

	// Apply search if provided
	search := GetSearch(r)
	if search != "" {
		filteredQs, err = lv.applySearch(filteredQs, search)
		if err != nil {
			return nil, fmt.Errorf("failed to apply search: %w", err)
		}
	}

	// Apply ordering
	filteredQs = lv.applyOrdering(r, filteredQs)

	// Get total count
	totalCount, err := filteredQs.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	// Get pagination params
	page := GetPage(r)
	pageSize := GetPageSize(r, lv.getDefaultPageSize())

	// Apply pagination
	adminQs, err := adminorm.NewAdminQuerySet(filteredQs)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin queryset: %w", err)
	}
	paginatedQs := adminQs.Paginate(page, pageSize)

	// Execute query
	objects, err := paginatedQs.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects: %w", err)
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	// Get display fields from config
	displayFields := lv.getDisplayFields(config)
	editableFields := lv.getEditableFields(config)
	displayLinks := lv.getDisplayLinks(config)
	// Already defined: dateHierarchy := lv.getDateHierarchy(config)

	// Get current sort order from request
	currentSort := lv.getCurrentSort(r)
	
	// Prepare display headers
	displayFieldHeaders := make([]ListFieldHeader, 0, len(displayFields))
	for _, fieldName := range displayFields {
		label := fieldName
		// Try to find label from schema or field info
		for _, f := range lv.admin.Fields() {
			if f.Name == fieldName {
				label = f.VerboseName
				break
			}
		}
		
		// Check if field is sortable
		sortable := lv.isFieldSortable(fieldName, config)
		sorted := false
		sortOrder := ""
		sortURL := ""
		
		if sortable {
			// Check if currently sorted
			if currentSort.Field == fieldName {
				sorted = true
				sortOrder = currentSort.Order
				// Toggle order for next click
				nextOrder := "asc"
				if currentSort.Order == "asc" {
					nextOrder = "desc"
				}
				sortURL = lv.buildSortURL(r, fieldName, nextOrder)
			} else {
				// Default to ascending
				sortURL = lv.buildSortURL(r, fieldName, "asc")
			}
		}
		
		displayFieldHeaders = append(displayFieldHeaders, ListFieldHeader{
			Name:      fieldName,
			Label:     label,
			Sortable:  sortable,
			Sorted:    sorted,
			SortOrder: sortOrder,
			SortURL:   sortURL,
		})
	}

	// Prepare rows
	rows := make([]ListRowData, 0, len(objects))
	for _, obj := range objects {
		id := lv.getID(obj)
		rowFields := make([]ListFieldVal, 0, len(displayFields))
		for _, fieldHeader := range displayFieldHeaders {
			val := lv.getFieldValue(obj, fieldHeader.Name)
			isLink := false
			for _, linkField := range displayLinks {
				if linkField == fieldHeader.Name {
					isLink = true
					break
				}
			}
			rowFields = append(rowFields, ListFieldVal{
				Value:  val,
				IsLink: isLink,
				URL:    fmt.Sprintf("/admin/%s/%v/", lv.admin.ModelName(), id),
			})
		}
		rows = append(rows, ListRowData{
			ID:        id,
			Fields:    rowFields,
			ChangeURL: fmt.Sprintf("/admin/%s/%v/", lv.admin.ModelName(), id),
			DeleteURL: fmt.Sprintf("/admin/%s/%v/delete/", lv.admin.ModelName(), id),
		})
	}

	// Prepare actions
	actions := make([]ActionData, 0)
	configActions := lv.admin.Config().Actions
	for _, act := range configActions {
		actions = append(actions, ActionData{
			Name:  act.Name,
			Label: act.Label,
		})
	}

	// Prepare search fields
	searchFields := make([]string, 0)
	for _, f := range config.SearchFields {
		if sf, ok := f.(string); ok {
			searchFields = append(searchFields, sf)
		}
	}

	return &ListData[T]{
		Objects:             objects,
		Rows:                rows,
		Page:                page,
		PageSize:            pageSize,
		TotalCount:          totalCount,
		TotalPages:          totalPages,
		Search:              search,
		SearchFields:        searchFields,
		DisplayFields:       displayFieldHeaders,
		EditableFields:      editableFields,
		DisplayLinks:        displayLinks,
		DateHierarchy:       dateHierarchy,
		EmptyValueDisplay:   lv.getEmptyValueDisplay(config),
		Filters:             lv.buildFilterData(ctx, r),
		DateHierarchyData:   lv.buildDateHierarchyData(ctx, r, dateHierarchy),
		Actions:             actions,
		HasAddPermission:    lv.admin.HasAddPermission(ctx, user),
		HasViewPermission:   lv.admin.HasViewPermission(ctx, user, nil),
		HasChangePermission: true, // Would check per object
		HasDeletePermission: true, // Would check per object
		ModelName:           lv.admin.ModelName(),
		AddURL:              fmt.Sprintf("/admin/%s/add/", lv.admin.ModelName()),
	}, nil
}

func (lv *ListView[T]) getID(instance interface{}) interface{} {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	f := val.FieldByName("ID")
	if f.IsValid() {
		return f.Interface()
	}
	return nil
}

func (lv *ListView[T]) getFieldValue(instance interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	f := val.FieldByName(fieldName)
	if f.IsValid() {
		return f.Interface()
	}
	return ""
}

// applySearch applies search to queryset using ORM field expressions
func (lv *ListView[T]) applySearch(qs orm.QuerySet[T], searchTerm string) (orm.QuerySet[T], error) {
	config := lv.admin.Config()
	if config == nil || len(config.SearchFields) == 0 || searchTerm == "" {
		return qs, nil
	}

	// Get field accessor for type-safe field expressions
	fa, err := lv.admin.Manager().FieldAccessor()
	if err != nil {
		return nil, fmt.Errorf("failed to get field accessor: %w", err)
	}

	// Get model schema to validate field types
	schema, err := orm.GetModelSchema[T]()
	if err != nil {
		return nil, fmt.Errorf("failed to get model schema: %w", err)
	}

	// Build OR conditions for search across all search fields
	var searchExpressions []orm.Expression

	for _, searchField := range config.SearchFields {
		fieldName := lv.getFieldName(searchField)
		if fieldName == "" {
			continue
		}

		// Get field info to check if it's a string field
		fieldInfo := schema.GetField(fieldName)
		if fieldInfo == nil {
			continue
		}

		// Only search string fields
		if fieldInfo.Type.Kind() == reflect.String {
			field := orm.FieldFor[T, string](fa, fieldName)
			containsExpr := field.Contains(searchTerm)
			searchExpressions = append(searchExpressions, containsExpr)
		}
	}

	if len(searchExpressions) == 0 {
		return qs, nil
	}

	// Combine expressions with OR
	var combinedExpr orm.Expression
	if len(searchExpressions) == 1 {
		combinedExpr = searchExpressions[0]
	} else {
		// Build OR chain: expr1 OR expr2 OR expr3 ... using new API
		combinedExpr = orm.Or(searchExpressions...)
	}

	// Apply filter
	return qs.Filter(combinedExpr), nil
}

// applyOrdering applies ordering to queryset
func (lv *ListView[T]) applyOrdering(r *http.Request, qs orm.QuerySet[T]) orm.QuerySet[T] {
	// 1. Check for URL parameter 'o' (ordering)
	if orderParam := r.URL.Query().Get("o"); orderParam != "" {
		orderFields := lv.parseOrderBy([]string{orderParam})
		if len(orderFields) > 0 {
			return qs.OrderBy(orderFields...)
		}
	}

	config := lv.admin.Config()
	if config == nil || len(config.Ordering) == 0 {
		// Use default ordering from schema meta
		meta := lv.admin.Meta()
		if len(meta.OrderBy) > 0 {
			orderFields := lv.parseOrderBy(meta.OrderBy)
			return qs.OrderBy(orderFields...)
		}
		return qs
	}

	// Convert config ordering to ORM order fields
	var orderFields []any
	for _, ordering := range config.Ordering {
		fieldName := lv.getFieldName(ordering.Field())
		if fieldName != "" {
			orderFields = append(orderFields, orm.OrderField{
				Field:     fieldName,
				Ascending: !ordering.IsDescending(),
			})
		}
	}

	if len(orderFields) > 0 {
		return qs.OrderBy(orderFields...)
	}

	return qs
}

// getFieldName extracts field name from interface{} (string or FieldExpr)
func (lv *ListView[T]) getFieldName(field interface{}) string {
	if name, ok := field.(string); ok {
		return name
	}
	// If it's a FieldExpr, we'd need to get the name from it
	// For now, return empty
	return ""
}

// getDisplayFields gets display field names from config
func (lv *ListView[T]) getDisplayFields(config *admin.Config[T]) []string {
	if config == nil || len(config.ListDisplay) == 0 {
		// Auto-discover from schema
		fields := lv.admin.Fields()
		fieldMapper := adminschema.NewFieldMapper()
		result := make([]string, 0)
		for _, field := range fields {
			if fieldMapper.ShouldDisplayInList(field.SchemaField) {
				result = append(result, field.Name)
			}
		}
		return result
	}

	result := make([]string, 0, len(config.ListDisplay))
	for _, field := range config.ListDisplay {
		if name := lv.getFieldName(field); name != "" {
			result = append(result, name)
		}
	}
	return result
}

// getEditableFields gets editable field names from config
func (lv *ListView[T]) getEditableFields(config *admin.Config[T]) []string {
	if config == nil {
		return []string{}
	}

	result := make([]string, 0, len(config.ListEditable))
	for _, field := range config.ListEditable {
		if name := lv.getFieldName(field); name != "" {
			result = append(result, name)
		}
	}
	return result
}

// getDisplayLinks gets display link field names from config
func (lv *ListView[T]) getDisplayLinks(config *admin.Config[T]) []string {
	if config == nil {
		return []string{}
	}

	result := make([]string, 0, len(config.ListDisplayLinks))
	for _, field := range config.ListDisplayLinks {
		if name := lv.getFieldName(field); name != "" {
			result = append(result, name)
		}
	}
	return result
}

// getDateHierarchy gets date hierarchy field name from config
func (lv *ListView[T]) getDateHierarchy(config *admin.Config[T]) string {
	if config == nil || config.DateHierarchy == nil {
		return ""
	}
	return lv.getFieldName(config.DateHierarchy)
}

// getEmptyValueDisplay gets empty value display from config
func (lv *ListView[T]) getEmptyValueDisplay(config *admin.Config[T]) string {
	if config == nil || config.EmptyValueDisplay == "" {
		return "-"
	}
	return config.EmptyValueDisplay
}

// getDefaultPageSize gets default page size from config
func (lv *ListView[T]) getDefaultPageSize() int {
	config := lv.admin.Config()
	if config != nil && config.ListPerPage > 0 {
		return config.ListPerPage
	}
	return 20 // Default
}

// parseOrderBy parses order by strings (e.g., "-created_at" for descending)
func (lv *ListView[T]) parseOrderBy(orderBy []string) []any {
	orderFields := make([]any, 0, len(orderBy))
	for _, field := range orderBy {
		ascending := true
		if len(field) > 0 && field[0] == '-' {
			ascending = false
			field = field[1:]
		}
		orderFields = append(orderFields, orm.OrderField{
			Field:     field,
			Ascending: ascending,
		})
	}
	return orderFields
}

// buildFilterData builds data for the filter sidebar
func (lv *ListView[T]) buildFilterData(ctx context.Context, r *http.Request) []FilterData {
	filterSet := lv.admin.FilterSet()
	if filterSet == nil {
		return nil
	}

	// Get base queryset to pass to GetOptions (as it might be used to count or filter options)
	qs, _ := lv.GetQueryset(ctx)

	filters := filterSet.FilterSet().GetFilters()
	result := make([]FilterData, 0, len(filters))

	for name, f := range filters {
		options, _ := f.GetOptions(ctx, qs)
		filterOptions := make([]FilterOption, 0, len(options))

		currentValue := r.URL.Query().Get(name)

		for _, opt := range options {
			filterOptions = append(filterOptions, FilterOption{
				Label:    opt.Label,
				Value:    opt.Value,
				Selected: fmt.Sprintf("%v", opt.Value) == currentValue,
			})
		}

		result = append(result, FilterData{
			Name:    name,
			Label:   f.Label(),
			Type:    "choices",
			Value:   currentValue,
			Options: filterOptions,
		})
	}

	return result
}

// buildDateHierarchyData builds data for date hierarchy navigation
func (lv *ListView[T]) buildDateHierarchyData(ctx context.Context, r *http.Request, fieldName string) *DateHierarchyData {
	if fieldName == "" {
		return nil
	}

	current := make(map[string]string)
	year := r.URL.Query().Get(fieldName + "__year")
	month := r.URL.Query().Get(fieldName + "__month")
	day := r.URL.Query().Get(fieldName + "__day")

	if year != "" {
		current["year"] = year
	}
	if month != "" {
		current["month"] = month
	}
	if day != "" {
		current["day"] = day
	}

	// Basic implementation - in a full implementation, this would query DB for unique dates
	// For now, return what's in the URL and some defaults
	years := []string{}
	// Logic to query unique years would go here
	// For now, if we have a year, we show it, otherwise maybe last 5 years
	if year != "" {
		years = append(years, year)
	} else {
		currYear := time.Now().Year()
		for i := 0; i < 5; i++ {
			years = append(years, strconv.Itoa(currYear-i))
		}
	}

	months := []string{}
	if year != "" {
		for m := 1; m <= 12; m++ {
			months = append(months, strconv.Itoa(m))
		}
	}

	return &DateHierarchyData{
		FieldName: fieldName,
		Years:     years,
		Months:    months,
		Days:      []string{},
		Current:   current,
	}
}

// applyDateHierarchy applies date hierarchy filters to the queryset
func (lv *ListView[T]) applyDateHierarchy(r *http.Request, qs orm.QuerySet[T], fieldName string) (orm.QuerySet[T], error) {
	year := r.URL.Query().Get(fieldName + "__year")
	month := r.URL.Query().Get(fieldName + "__month")
	day := r.URL.Query().Get(fieldName + "__day")

	if year == "" {
		return qs, nil
	}

	iYear, err := strconv.Atoi(year)
	if err != nil {
		return qs, nil
	}

	// Use generic expressions for Year/Month/Day
	// In Forge ORM, we can use OpYear, OpMonth, OpDay

	// Year filter
	// Actually, a better way is to use a range for the year/month/day to be index-friendly
	start := time.Date(iYear, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	if month != "" {
		iMonth, _ := strconv.Atoi(month)
		start = time.Date(iYear, time.Month(iMonth), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)

		if day != "" {
			iDay, _ := strconv.Atoi(day)
			start = time.Date(iYear, time.Month(iMonth), iDay, 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 0, 1)
		}
	}

	// Use Gte and Lt for efficiency
	// We need to create expressions for these.
	// Since we don't have FieldExpr[T] here easily, we can use generic expressions or
	// assume there's a helper.

	// For simplicity, let's use the year/month/day operators if we assume ORM supports them
	// qs = qs.Filter(orm.NewQueryExpr(fieldName, orm.OpYear, iYear))
	// But let's use the date range as it's more robust.

	// We need something like orm.Field(name).Gte(val)
	// But we are in a generic context.

	// Let's use the genericExpression we defined in config.go if we can,
	// but it's not exported.
	// Wait, I can just use orm.NewField[time.Time](fieldName, "").GreaterOrEqual(start)

	// Create expressions using NewField and Gte/Lt
	expr1 := orm.NewField[time.Time](fieldName, "").Gte(start)
	expr2 := orm.NewField[time.Time](fieldName, "").Lt(end)

	return qs.Filter(expr1).Filter(expr2), nil
}

// SortInfo represents current sort state
type SortInfo struct {
	Field string
	Order string // "asc" or "desc"
}

// getCurrentSort extracts current sort from request
func (lv *ListView[T]) getCurrentSort(r *http.Request) SortInfo {
	orderParam := r.URL.Query().Get("o")
	if orderParam == "" {
		return SortInfo{}
	}

	// Parse order param (e.g., "created_at" or "-created_at")
	field := orderParam
	order := "asc"
	
	if strings.HasPrefix(field, "-") {
		field = field[1:]
		order = "desc"
	}

	return SortInfo{
		Field: field,
		Order: order,
	}
}

// isFieldSortable checks if a field is sortable
func (lv *ListView[T]) isFieldSortable(fieldName string, config *admin.Config[T]) bool {
	// Check if field is in SortableBy
	if config != nil && len(config.SortableBy) > 0 {
		for _, sortableField := range config.SortableBy {
			if sf, ok := sortableField.(string); ok && sf == fieldName {
				return true
			}
		}
		return false
	}

	// Default: all display fields are sortable unless explicitly disabled
	return true
}

// buildSortURL builds URL for sorting by a field
func (lv *ListView[T]) buildSortURL(r *http.Request, fieldName string, order string) string {
	// Get current URL
	u := r.URL
	
	// Create new query values
	query := u.Query()
	
	// Set order parameter
	orderValue := fieldName
	if order == "desc" {
		orderValue = "-" + fieldName
	}
	query.Set("o", orderValue)
	
	// Reset page to 1 when sorting
	query.Set("page", "1")
	
	// Build new URL
	newURL := url.URL{
		Path:     u.Path,
		RawQuery: query.Encode(),
	}
	
	return newURL.String()
}
