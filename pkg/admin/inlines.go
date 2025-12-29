package admin

import (
	"github.com/forgego/forge/pkg/query"
)

// Inline represents inline editing for related models
type Inline[T any, R any] struct {
	model       R
	manager     *query.Manager[R]
	parentField FieldExpr[R, *T]
	fields      []FieldExpr[R, interface{}]
	extra       int
	maxNum      int
	style       InlineStyle
}

// InlineStyle specifies the display style
type InlineStyle string

const (
	InlineTabular InlineStyle = "tabular"
	InlineStacked InlineStyle = "stacked"
)

// TabularInline creates a tabular inline
func TabularInline[T any, R any](
	model R,
	manager *query.Manager[R],
	parentField FieldExpr[R, *T],
	fields []FieldExpr[R, interface{}],
) Inline[T, R] {
	return Inline[T, R]{
		model:       model,
		manager:     manager,
		parentField: parentField,
		fields:      fields,
		extra:       1,
		style:       InlineTabular,
	}
}

// StackedInline creates a stacked inline
func StackedInline[T any, R any](
	model R,
	manager *query.Manager[R],
	parentField FieldExpr[R, *T],
	fields []FieldExpr[R, interface{}],
) Inline[T, R] {
	return Inline[T, R]{
		model:       model,
		manager:     manager,
		parentField: parentField,
		fields:      fields,
		extra:       1,
		style:       InlineStacked,
	}
}

// WithExtra sets the number of extra empty forms
func (i Inline[T, R]) WithExtra(extra int) Inline[T, R] {
	i.extra = extra
	return i
}

// WithMaxNum sets the maximum number of forms
func (i Inline[T, R]) WithMaxNum(maxNum int) Inline[T, R] {
	i.maxNum = maxNum
	return i
}

// Model returns the inline model
func (i Inline[T, R]) Model() R {
	return i.model
}

// Manager returns the inline manager
func (i Inline[T, R]) Manager() *query.Manager[R] {
	return i.manager
}

// ParentField returns the parent field expression
func (i Inline[T, R]) ParentField() FieldExpr[R, *T] {
	return i.parentField
}

// Fields returns the fields to display
func (i Inline[T, R]) Fields() []FieldExpr[R, interface{}] {
	return i.fields
}

// Style returns the inline style
func (i Inline[T, R]) Style() InlineStyle {
	return i.style
}
