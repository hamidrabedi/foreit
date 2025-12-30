package orm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/forgego/forge/orm"
)

// AdminQuerySet wraps orm.QuerySet for admin use with additional admin-specific methods
type AdminQuerySet[T any] struct {
	queryset orm.QuerySet[T]
	schema   *orm.ModelSchema
}

// NewAdminQuerySet creates a new admin queryset wrapper
func NewAdminQuerySet[T any](queryset orm.QuerySet[T]) (*AdminQuerySet[T], error) {
	schema, err := orm.GetModelSchema[T]()
	if err != nil {
		return nil, err
	}

	return &AdminQuerySet[T]{
		queryset: queryset,
		schema:   schema,
	}, nil
}

// QuerySet returns the underlying ORM queryset
func (aqs *AdminQuerySet[T]) QuerySet() orm.QuerySet[T] {
	return aqs.queryset
}

// ApplySearch applies search to the queryset using search fields
func (aqs *AdminQuerySet[T]) ApplySearch(ctx context.Context, searchTerm string, searchFields []string) (orm.QuerySet[T], error) {
	if searchTerm == "" || len(searchFields) == 0 {
		return aqs.queryset, nil
	}

	// Get field accessor
	fa, err := orm.NewFieldAccessor[T]()
	if err != nil {
		return nil, fmt.Errorf("failed to get field accessor: %w", err)
	}

	// Build OR conditions for search using Q expressions
	var searchExpressions []orm.Expression

	for _, fieldName := range searchFields {
		// Get field info to determine type
		fieldInfo := aqs.schema.GetField(fieldName)
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
		return aqs.queryset, nil
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
	result := aqs.queryset.Filter(combinedExpr)
	return result, nil
}

// ApplyOrdering applies ordering to the queryset
func (aqs *AdminQuerySet[T]) ApplyOrdering(orderFields []orm.OrderField) orm.QuerySet[T] {
	if len(orderFields) == 0 {
		return aqs.queryset
	}
	// Convert []OrderField to []any
	anyFields := make([]any, len(orderFields))
	for i, f := range orderFields {
		anyFields[i] = f
	}
	return aqs.queryset.OrderBy(anyFields...)
}

// Paginate applies pagination to the queryset
func (aqs *AdminQuerySet[T]) Paginate(page, pageSize int) orm.QuerySet[T] {
	offset := (page - 1) * pageSize
	return aqs.queryset.Offset(offset).Limit(pageSize)
}

// Count returns the total count
func (aqs *AdminQuerySet[T]) Count(ctx context.Context) (int64, error) {
	return aqs.queryset.Count(ctx)
}

// All returns all results
func (aqs *AdminQuerySet[T]) All(ctx context.Context) ([]*T, error) {
	return aqs.queryset.All(ctx)
}

// Get returns a single result
func (aqs *AdminQuerySet[T]) Get(ctx context.Context) (*T, error) {
	return aqs.queryset.Get(ctx)
}

// First returns the first result
func (aqs *AdminQuerySet[T]) First(ctx context.Context) (*T, error) {
	return aqs.queryset.First(ctx)
}

// Last returns the last result
func (aqs *AdminQuerySet[T]) Last(ctx context.Context) (*T, error) {
	return aqs.queryset.Last(ctx)
}

// SelectRelated applies select_related optimization
func (aqs *AdminQuerySet[T]) SelectRelated(relations ...string) *AdminQuerySet[T] {
	// Convert []string to []any
	anyRelations := make([]any, len(relations))
	for i, r := range relations {
		anyRelations[i] = r
	}
	return &AdminQuerySet[T]{
		queryset: aqs.queryset.SelectRelated(anyRelations...),
		schema:   aqs.schema,
	}
}

// PrefetchRelated applies prefetch_related optimization
func (aqs *AdminQuerySet[T]) PrefetchRelated(relations ...string) *AdminQuerySet[T] {
	// Convert []string to []any
	anyRelations := make([]any, len(relations))
	for i, r := range relations {
		anyRelations[i] = r
	}
	return &AdminQuerySet[T]{
		queryset: aqs.queryset.PrefetchRelated(anyRelations...),
		schema:   aqs.schema,
	}
}
