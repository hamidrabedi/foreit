package admin

import (
	"context"
	"fmt"
	"reflect"
)

// FormView represents a type-safe form view
type FormView[T any] struct {
	admin    *Admin[T]
	instance *T
	form     Form[T]
	isNew    bool
}

// FormViewData contains data for rendering the form view
type FormViewData[T any] struct {
	Instance *T
	Form     Form[T]
	IsNew    bool
}

// NewFormView creates a new form view
func NewFormView[T any](admin *Admin[T], instance *T, isNew bool) *FormView[T] {
	form := NewForm(instance, isNew)

	// Generate form fields from config
	config := admin.Config()
	if config.GetForm != nil {
		// Use custom form generator
		customForm, err := config.GetForm(context.Background(), instance, isNew)
		if err == nil {
			form = customForm
		}
	} else {
		// Auto-generate form fields
		form = generateFormFields(admin, instance, isNew)
	}

	return &FormView[T]{
		admin:    admin,
		instance: instance,
		form:     form,
		isNew:    isNew,
	}
}

// generateFormFields generates form fields from admin config
func generateFormFields[T any](admin *Admin[T], instance *T, isNew bool) Form[T] {
	form := NewForm(instance, isNew)

	// If fieldsets are configured, use them
	config := admin.Config()
	if len(config.Fieldsets) > 0 {
		for _, fieldset := range config.Fieldsets {
			for _, field := range fieldset.Fields {
				// Check if field is read-only
				readonly := isReadOnly(admin, field)

				formField := FormField[T]{
					expr:     field,
					value:    field.Get(instance),
					required: false, // TODO: Get from field definition
					readonly: readonly,
				}
				form.AddField(formField)
			}
		}
	} else {
		// Use list display fields or all fields
		fields := config.ListDisplay
		if len(fields) == 0 {
			// TODO: Get all fields from schema
			return form
		}

		for _, field := range fields {
			readonly := isReadOnly(admin, field)

			formField := FormField[T]{
				expr:     field,
				value:    field.Get(instance),
				required: false,
				readonly: readonly,
			}
			form.AddField(formField)
		}
	}

	return form
}

// isReadOnly checks if a field is read-only
func isReadOnly[T any](admin *Admin[T], field FieldExpr[T, interface{}]) bool {
	for _, readonlyField := range admin.Config().ReadOnlyFields {
		if readonlyField.Name() == field.Name() {
			return true
		}
	}
	return false
}

// Render renders the form view and returns the data
func (v *FormView[T]) Render(ctx context.Context) (*FormViewData[T], error) {
	return &FormViewData[T]{
		Instance: v.instance,
		Form:     v.form,
		IsNew:    v.isNew,
	}, nil
}

// Instance returns the form instance
func (v *FormView[T]) Instance() *T {
	return v.instance
}

// Form returns the form
func (v *FormView[T]) Form() Form[T] {
	return v.form
}

// Save saves the form data
func (v *FormView[T]) Save(ctx context.Context, formData FormData) error {
	// Populate instance from form data
	if err := v.populateFromForm(formData); err != nil {
		return fmt.Errorf("failed to populate form: %w", err)
	}

	// Save using admin's SaveModel
	return v.admin.SaveModel(ctx, v.instance, formData, v.isNew)
}

// populateFromForm populates instance from form data
func (v *FormView[T]) populateFromForm(formData FormData) error {
	for _, field := range v.form.Fields() {
		// Check if field is readonly using reflection or exported method
		if isFieldReadonly(field) {
			continue
		}

		// Get field name from field expression
		fieldName := getFieldExprName(field)
		value, ok := formData[fieldName]
		if !ok {
			continue
		}

		// Set field value using field expression
		setFieldValue(field, v.instance, value)
	}

	return nil
}

// isFieldReadonly checks if a field is readonly
func isFieldReadonly[T any](field FormField[T]) bool {
	// Access readonly field using reflection
	val := reflect.ValueOf(field)
	readonlyField := val.FieldByName("readonly")
	if readonlyField.IsValid() && readonlyField.CanInterface() {
		if readonly, ok := readonlyField.Interface().(bool); ok {
			return readonly
		}
	}
	return false
}

// getFieldExprName gets the field name from field expression
func getFieldExprName[T any](field FormField[T]) string {
	// Access expr field using reflection
	val := reflect.ValueOf(field)
	exprField := val.FieldByName("expr")
	if exprField.IsValid() && exprField.CanInterface() {
		if expr, ok := exprField.Interface().(FieldExpr[T, interface{}]); ok {
			return expr.Name()
		}
	}
	return ""
}

// setFieldValue sets a field value using field expression
func setFieldValue[T any](field FormField[T], instance *T, value interface{}) {
	// Access expr field using reflection
	val := reflect.ValueOf(field)
	exprField := val.FieldByName("expr")
	if exprField.IsValid() && exprField.CanInterface() {
		if expr, ok := exprField.Interface().(FieldExpr[T, interface{}]); ok {
			expr.Set(instance, value)
		}
	}
}

// GetFormViewByID gets an instance by ID and creates a form view
func GetFormViewByID[T any](admin *Admin[T], ctx context.Context, id int64, isNew bool) (*FormView[T], error) {
	var instance *T
	var err error

	if isNew {
		// Create new instance
		var zero T
		instance = &zero
	} else {
		// Get existing instance
		instance, err = admin.Manager().Get(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get instance: %w", err)
		}
	}

	return NewFormView(admin, instance, isNew), nil
}
