package admin

// Form represents a type-safe form
type Form[T any] struct {
	instance *T
	fields   []FormField[T]
	errors   map[string]string
	isNew    bool
}

// FormField represents a form field
type FormField[T any] struct {
	expr     FieldExpr[T, interface{}]
	widget   Widget
	value    interface{}
	required bool
	readonly bool
	helpText string
	errors   []string
}

// NewForm creates a new form
func NewForm[T any](instance *T, isNew bool) Form[T] {
	return Form[T]{
		instance: instance,
		fields:   []FormField[T]{},
		errors:   make(map[string]string),
		isNew:    isNew,
	}
}

// AddField adds a field to the form
func (f *Form[T]) AddField(field FormField[T]) {
	f.fields = append(f.fields, field)
}

// Fields returns all form fields
func (f Form[T]) Fields() []FormField[T] {
	return f.fields
}

// Instance returns the form instance
func (f Form[T]) Instance() *T {
	return f.instance
}

// IsNew returns true if this is a new instance
func (f Form[T]) IsNew() bool {
	return f.isNew
}

// Errors returns form errors
func (f Form[T]) Errors() map[string]string {
	return f.errors
}

// AddError adds an error to the form
func (f *Form[T]) AddError(field string, message string) {
	f.errors[field] = message
}

// Widget interface is defined in widgets.go
