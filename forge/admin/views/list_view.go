package views

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

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
	Page                int
	PageSize            int
	TotalCount          int64
	TotalPages          int
	Search              string
	DisplayFields       []string // Field names to display
	EditableFields      []string // Field names editable in list
	DisplayLinks        []string // Field names that link to detail
	DateHierarchy       string   // Field name for date hierarchy
	EmptyValueDisplay   string
	HasAddPermission    bool
	HasChangePermission bool
	HasDeletePermission bool
}

// Render renders the list view and returns the data
func (lv *ListView[T]) Render(ctx context.Context, r *http.Request, user interface{}) (*ListData[T], error) {
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

	// Apply search if provided
	search := GetSearch(r)
	if search != "" {
		filteredQs, err = lv.applySearch(filteredQs, search)
		if err != nil {
			return nil, fmt.Errorf("failed to apply search: %w", err)
		}
	}

	// Apply ordering
	filteredQs = lv.applyOrdering(filteredQs)

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
	dateHierarchy := lv.getDateHierarchy(config)

	return &ListData[T]{
		Objects:             objects,
		Page:                page,
		PageSize:            pageSize,
		TotalCount:          totalCount,
		TotalPages:          totalPages,
		Search:              search,
		DisplayFields:       displayFields,
		EditableFields:      editableFields,
		DisplayLinks:        displayLinks,
		DateHierarchy:       dateHierarchy,
		EmptyValueDisplay:   lv.getEmptyValueDisplay(config),
		HasAddPermission:    lv.admin.HasAddPermission(ctx, user),
		HasChangePermission: true, // Would check per object
		HasDeletePermission: true, // Would check per object
	}, nil
}

// applySearch applies search to queryset using ORM field expressions
func (lv *ListView[T]) applySearch(qs orm.QuerySet[T], searchTerm string) (orm.QuerySet[T], error) {
	config := lv.admin.Config()
	if config == nil || len(config.SearchFields) == 0 || searchTerm == "" {
		return qs, nil
	}

	// Get field accessor for type-safe field expressions
	fa, err := lv.admin.Manager().GetFieldAccessor()
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
		// Build OR chain: expr1 OR expr2 OR expr3 ...
		q := orm.NewQ(searchExpressions[0])
		for i := 1; i < len(searchExpressions); i++ {
			q = q.Or(orm.NewQ(searchExpressions[i]))
		}
		combinedExpr = q
	}

	// Apply filter
	return qs.Filter(combinedExpr), nil
}

// applyOrdering applies ordering to queryset
func (lv *ListView[T]) applyOrdering(qs orm.QuerySet[T]) orm.QuerySet[T] {
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
	var orderFields []orm.OrderField
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
func (lv *ListView[T]) parseOrderBy(orderBy []string) []orm.OrderField {
	orderFields := make([]orm.OrderField, 0, len(orderBy))
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
