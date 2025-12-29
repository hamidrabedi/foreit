package admin

import (
	"context"

	"github.com/forgego/forge/pkg/query"
)

// Filter represents a filter for list views
type Filter[T any] interface {
	Name() string
	Label() string
	GetOptions(ctx context.Context, qs query.QuerySet[T]) ([]FilterOption, error)
	Apply(ctx context.Context, qs query.QuerySet[T], value interface{}) (query.QuerySet[T], error)
}

// FilterOption represents a filter option
type FilterOption struct {
	Label string
	Value interface{}
	Count int64
}

// BooleanFilter is a filter for boolean fields
type BooleanFilter[T any] struct {
	field FieldExpr[T, bool]
	name  string
	label string
}

// NewBooleanFilter creates a boolean filter
func NewBooleanFilter[T any](field FieldExpr[T, bool]) Filter[T] {
	return &BooleanFilter[T]{
		field: field,
		name:  field.Name(),
		label: field.Name(),
	}
}

// Name returns the filter name
func (f *BooleanFilter[T]) Name() string {
	return f.name
}

// Label returns the filter label
func (f *BooleanFilter[T]) Label() string {
	return f.label
}

// GetOptions returns filter options
func (f *BooleanFilter[T]) GetOptions(ctx context.Context, qs query.QuerySet[T]) ([]FilterOption, error) {
	return []FilterOption{
		{Label: "Yes", Value: true},
		{Label: "No", Value: false},
	}, nil
}

// Apply applies the filter to a queryset
func (f *BooleanFilter[T]) Apply(ctx context.Context, qs query.QuerySet[T], value interface{}) (query.QuerySet[T], error) {
	boolValue, ok := value.(bool)
	if !ok {
		return qs, nil
	}

	// Create query expression
	expr := query.NewFieldQueryExpr(f.name, query.OpEquals, boolValue)
	return qs.Filter(expr), nil
}

// ChoiceFilter is a filter for choice fields
type ChoiceFilter[T any, F comparable] struct {
	field   FieldExpr[T, F]
	name    string
	label   string
	choices []Choice[F]
}

// Choice represents a choice option
type Choice[F any] struct {
	Label string
	Value F
}

// NewChoiceFilter creates a choice filter
func NewChoiceFilter[T any, F comparable](field FieldExpr[T, F], choices []Choice[F]) Filter[T] {
	return &ChoiceFilter[T, F]{
		field:   field,
		name:    field.Name(),
		label:   field.Name(),
		choices: choices,
	}
}

// Name returns the filter name
func (f *ChoiceFilter[T, F]) Name() string {
	return f.name
}

// Label returns the filter label
func (f *ChoiceFilter[T, F]) Label() string {
	return f.label
}

// GetOptions returns filter options
func (f *ChoiceFilter[T, F]) GetOptions(ctx context.Context, qs query.QuerySet[T]) ([]FilterOption, error) {
	options := make([]FilterOption, len(f.choices))
	for i, choice := range f.choices {
		options[i] = FilterOption{
			Label: choice.Label,
			Value: choice.Value,
		}
	}
	return options, nil
}

// Apply applies the filter to a queryset
func (f *ChoiceFilter[T, F]) Apply(ctx context.Context, qs query.QuerySet[T], value interface{}) (query.QuerySet[T], error) {
	// Type assert value
	val, ok := value.(F)
	if !ok {
		return qs, nil
	}

	// Create query expression
	expr := query.NewFieldQueryExpr(f.name, query.OpEquals, val)
	return qs.Filter(expr), nil
}
