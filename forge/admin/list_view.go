package admin

import (
	"context"
	"fmt"

	query "github.com/forgego/forge/orm"
)

// ListView represents a type-safe list view
type ListView[T any] struct {
	admin    *Admin[T]
	queryset query.QuerySet[T]
	page     int
	pageSize int
	search   string
	filters  map[string]interface{}
}

// ListData contains data for rendering the list view
type ListData[T any] struct {
	Objects           []*T
	Page              int
	PageSize          int
	TotalCount        int64
	TotalPages        int
	Search            string
	DisplayFields     []FieldExpr[T, interface{}]
	EditableFields    []FieldExpr[T, interface{}]
	DisplayLinks      []FieldExpr[T, interface{}]
	DateHierarchy     FieldExpr[T, interface{}]
	EmptyValueDisplay string
}

// NewListView creates a new list view
func NewListView[T any](admin *Admin[T]) *ListView[T] {
	return &ListView[T]{
		admin:    admin,
		page:     1,
		pageSize: admin.Config().ListPerPage,
		filters:  make(map[string]interface{}),
	}
}

// WithPage sets the page number
func (v *ListView[T]) WithPage(page int) *ListView[T] {
	v.page = page
	return v
}

// WithPageSize sets the page size
func (v *ListView[T]) WithPageSize(size int) *ListView[T] {
	v.pageSize = size
	return v
}

// WithSearch sets the search query
func (v *ListView[T]) WithSearch(search string) *ListView[T] {
	v.search = search
	return v
}

// WithFilter sets a filter value
func (v *ListView[T]) WithFilter(name string, value interface{}) *ListView[T] {
	v.filters[name] = value
	return v
}

// Render renders the list view and returns the data
func (v *ListView[T]) Render(ctx context.Context) (*ListData[T], error) {
	// Get queryset from manager
	qs, err := v.admin.GetQueryset(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queryset: %w", err)
	}

	// Apply search
	if v.search != "" && len(v.admin.Config().SearchFields) > 0 {
		qs = v.applySearch(qs, v.search)
	}

	// Apply filters (both old-style and new FilterSet if available)
	for name, value := range v.filters {
		qs = v.applyFilter(qs, name, value)
	}

	// Apply ordering
	if len(v.admin.Config().Ordering) > 0 {
		qs = v.applyOrdering(qs)
	}

	// Get total count
	totalCount, err := qs.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	// Apply pagination
	offset := (v.page - 1) * v.pageSize
	qs = qs.Offset(offset).Limit(v.pageSize)

	// Get objects
	objects, err := qs.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects: %w", err)
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(v.pageSize) - 1) / int64(v.pageSize))

	config := v.admin.Config()

	// Get display fields
	displayFields := config.ListDisplay
	if len(displayFields) == 0 {
		// Default: use all fields if none specified
		displayFields = []FieldExpr[T, interface{}]{}
	}

	return &ListData[T]{
		Objects:           objects,
		Page:              v.page,
		PageSize:          v.pageSize,
		TotalCount:        totalCount,
		TotalPages:        totalPages,
		Search:            v.search,
		DisplayFields:     displayFields,
		EditableFields:    config.ListEditable,
		DisplayLinks:      config.ListDisplayLinks,
		DateHierarchy:     config.DateHierarchy,
		EmptyValueDisplay: config.EmptyValueDisplay,
	}, nil
}

// applySearch applies search to queryset
func (v *ListView[T]) applySearch(qs query.QuerySet[T], search string) query.QuerySet[T] {
	// Build OR expression for all search fields
	var exprs []query.Expression
	for _, field := range v.admin.Config().SearchFields {
		// Create contains expression for each search field
		ormField := query.NewField[string](field.Name(), "")
		expr := ormField.Contains(search)
		exprs = append(exprs, expr)
	}

	// Combine with OR
	if len(exprs) > 0 {
		combined := query.NewQ(exprs[0])
		for i := 1; i < len(exprs); i++ {
			combined = combined.Or(query.NewQ(exprs[i]))
		}
		qs = qs.Filter(combined)
	}

	return qs
}

// applyFilter applies a filter to queryset
func (v *ListView[T]) applyFilter(qs query.QuerySet[T], name string, value interface{}) query.QuerySet[T] {
	// Find filter by name
	for _, filter := range v.admin.Config().ListFilter {
		if filter.Name() == name {
			// Apply filter (this will use context from queryset)
			ctx := context.Background() // TODO: Get from queryset
			filtered, err := filter.Apply(ctx, qs, value)
			if err == nil {
				return filtered
			}
		}
	}
	return qs
}

// applyOrdering applies ordering to queryset
func (v *ListView[T]) applyOrdering(qs query.QuerySet[T]) query.QuerySet[T] {
	var orderFields []query.OrderField
	for _, ordering := range v.admin.Config().Ordering {
		fieldName := ordering.Field().Name()
		if ordering.IsDescending() {
			orderFields = append(orderFields, query.Desc(fieldName))
		} else {
			orderFields = append(orderFields, query.Asc(fieldName))
		}
	}
	if len(orderFields) > 0 {
		qs = qs.OrderBy(orderFields...)
	}
	return qs
}
