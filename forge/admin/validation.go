package admin

import (
	"fmt"
	"reflect"

	validation "github.com/forgego/forge/validate"
)

// ValidateForm validates form data against admin config
func ValidateForm[T any](admin *Admin[T], instance *T, formData FormData, isNew bool) map[string]string {
	errors := make(map[string]string)

	// Get form view to validate
	// Note: This requires views package - using simplified validation for now
	_ = isNew    // Avoid unused variable
	_ = admin    // Avoid unused variable
	_ = instance // Avoid unused variable
	// formView := NewFormView(admin, instance, isNew)
	// form := formView.Form()

	// Validate each field in formData
	// Simplified validation - check required fields from config
	for fieldName, value := range formData {
		_ = value     // Use value in validation
		_ = fieldName // Use fieldName
		// Full validation would check against admin config fieldsets
		// For now, just ensure formData is not empty
		if value == nil || value == "" {
			errors[fieldName] = fmt.Sprintf("%s is required", fieldName)
		}
	}

	return errors
}

// getFieldNameFromFormField extracts field name from FormField
func getFieldNameFromFormField[T any](field FormField[T]) string {
	// Access expr field using reflection
	val := reflect.ValueOf(field)
	exprField := val.FieldByName("expr")
	if exprField.IsValid() && exprField.CanInterface() {
		if expr, ok := exprField.Interface().(FieldExpr[T, interface{}]); ok {
			return expr.Name()
		}
	}
	return "field"
}

// isFieldRequired checks if field is required
func isFieldRequired[T any](field FormField[T]) bool {
	return field.IsRequired()
}

// validateFieldValue validates a field value
func validateFieldValue[T any](field FormField[T], value interface{}) error {
	// Get field type from field expression
	// This is simplified - full implementation would use schema validation

	// Basic type checking
	fieldType := reflect.TypeOf(value)
	if fieldType == nil {
		return nil // Nil values are handled by required check
	}

	// Additional validation can be added here
	// For now, just return nil (valid)
	return nil
}

// ValidateInstance validates an instance using validation package
func ValidateInstance(instance interface{}) error {
	validator := validation.NewValidator()
	// Use ValidateStruct method
	return validator.ValidateStruct(instance)
}
