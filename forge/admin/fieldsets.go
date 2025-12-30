package admin

// Fieldset groups form fields together
type Fieldset[T any] struct {
	Name      string
	Fields    []FieldExpr[T, interface{}]
	Collapsed bool
	Classes   []string
}

// NewFieldset creates a new fieldset
func NewFieldset[T any](name string, fields ...FieldExpr[T, interface{}]) Fieldset[T] {
	return Fieldset[T]{
		Name:   name,
		Fields: fields,
	}
}

// WithCollapsed sets the fieldset to be collapsed by default
func (f Fieldset[T]) WithCollapsed(collapsed bool) Fieldset[T] {
	f.Collapsed = collapsed
	return f
}

// WithClasses adds CSS classes to the fieldset
func (f Fieldset[T]) WithClasses(classes ...string) Fieldset[T] {
	f.Classes = classes
	return f
}
