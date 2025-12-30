package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/schema"
)

// FormValidator validates form data against schema
type FormValidator struct {
	schemaInstance schema.Schema
}

// NewFormValidator creates a new form validator
func NewFormValidator(schemaInstance schema.Schema) *FormValidator {
	return &FormValidator{
		schemaInstance: schemaInstance,
	}
}

// Validate validates form data
func (fv *FormValidator) Validate(formData admin.FormData) (map[string][]string, error) {
	errors := make(map[string][]string)
	fields := fv.schemaInstance.Fields()
	
	for _, field := range fields {
		fieldErrors := fv.validateField(field, formData)
		if len(fieldErrors) > 0 {
			errors[field.Name] = fieldErrors
		}
	}
	
	if len(errors) > 0 {
		return errors, fmt.Errorf("validation failed")
	}
	
	return nil, nil
}

// validateField validates a single field
func (fv *FormValidator) validateField(field schema.Field, formData admin.FormData) []string {
	var errors []string
	value, exists := formData[field.Name]
	
	// Required check
	if field.Required && (!exists || isEmpty(value)) {
		errors = append(errors, fmt.Sprintf("%s is required", field.VerboseName))
		return errors
	}
	
	if !exists || isEmpty(value) {
		return errors // Skip validation for empty optional fields
	}
	
	// Type-specific validation
	switch field.Type {
	case schema.TypeString, schema.TypeText, schema.TypeEmail, schema.TypeURL:
		errors = append(errors, fv.validateStringField(field, value)...)
	case schema.TypeInt64, schema.TypeInt32:
		errors = append(errors, fv.validateIntField(field, value)...)
	case schema.TypeFloat32, schema.TypeFloat64, schema.TypeDecimal:
		errors = append(errors, fv.validateFloatField(field, value)...)
	case schema.TypeBool:
		// Boolean validation is usually straightforward
	case schema.TypeDate, schema.TypeDateTime, schema.TypeTime:
		errors = append(errors, fv.validateDateField(field, value)...)
	}
	
	// Custom validators
	for _, validator := range field.Validators {
		if err := validator.Validate(value); err != nil {
			errors = append(errors, err.Error())
		}
	}
	
	return errors
}

// validateStringField validates string fields
func (fv *FormValidator) validateStringField(field schema.Field, value interface{}) []string {
	var errors []string
	valueStr := fmt.Sprintf("%v", value)
	
	// Length validation
	if field.MinLength != nil && len(valueStr) < *field.MinLength {
		errors = append(errors, fmt.Sprintf("%s must be at least %d characters", field.VerboseName, *field.MinLength))
	}
	if field.MaxLength != nil && len(valueStr) > *field.MaxLength {
		errors = append(errors, fmt.Sprintf("%s must be no more than %d characters", field.VerboseName, *field.MaxLength))
	}
	
	// Email validation
	if field.Type == schema.TypeEmail {
		if !strings.Contains(valueStr, "@") {
			errors = append(errors, fmt.Sprintf("%s must be a valid email address", field.VerboseName))
		}
	}
	
	// URL validation
	if field.Type == schema.TypeURL {
		if !strings.HasPrefix(valueStr, "http://") && !strings.HasPrefix(valueStr, "https://") {
			errors = append(errors, fmt.Sprintf("%s must be a valid URL", field.VerboseName))
		}
	}
	
	// Choices validation
	if len(field.Choices) > 0 {
		validChoice := false
		for _, choice := range field.Choices {
			if fmt.Sprintf("%v", choice.Value) == valueStr {
				validChoice = true
				break
			}
		}
		if !validChoice {
			errors = append(errors, fmt.Sprintf("%s must be one of the available choices", field.VerboseName))
		}
	}
	
	return errors
}

// validateIntField validates integer fields
func (fv *FormValidator) validateIntField(field schema.Field, value interface{}) []string {
	var errors []string
	
	// Try to convert to int64
	var intValue int64
	switch v := value.(type) {
	case int64:
		intValue = v
	case int:
		intValue = int64(v)
	case string:
		// Try parsing
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s must be a valid integer", field.VerboseName))
			return errors
		}
		intValue = parsed
	default:
		errors = append(errors, fmt.Sprintf("%s must be a valid integer", field.VerboseName))
		return errors
	}
	
	// Value range validation
	if field.MinValue != nil && float64(intValue) < *field.MinValue {
		errors = append(errors, fmt.Sprintf("%s must be at least %v", field.VerboseName, *field.MinValue))
	}
	if field.MaxValue != nil && float64(intValue) > *field.MaxValue {
		errors = append(errors, fmt.Sprintf("%s must be no more than %v", field.VerboseName, *field.MaxValue))
	}
	
	return errors
}

// validateFloatField validates float fields
func (fv *FormValidator) validateFloatField(field schema.Field, value interface{}) []string {
	var errors []string
	
	// Try to convert to float64
	var floatValue float64
	switch v := value.(type) {
	case float64:
		floatValue = v
	case float32:
		floatValue = float64(v)
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s must be a valid number", field.VerboseName))
			return errors
		}
		floatValue = parsed
	default:
		errors = append(errors, fmt.Sprintf("%s must be a valid number", field.VerboseName))
		return errors
	}
	
	// Value range validation
	if field.MinValue != nil && floatValue < *field.MinValue {
		errors = append(errors, fmt.Sprintf("%s must be at least %v", field.VerboseName, *field.MinValue))
	}
	if field.MaxValue != nil && floatValue > *field.MaxValue {
		errors = append(errors, fmt.Sprintf("%s must be no more than %v", field.VerboseName, *field.MaxValue))
	}
	
	return errors
}

// validateDateField validates date/time fields
func (fv *FormValidator) validateDateField(field schema.Field, value interface{}) []string {
	var errors []string
	
	// Date validation would parse and check format
	// For now, just check if it's not empty
	valueStr := fmt.Sprintf("%v", value)
	if valueStr == "" && field.Required {
		errors = append(errors, fmt.Sprintf("%s is required", field.VerboseName))
	}
	
	return errors
}

// isEmpty checks if a value is empty
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	}
	
	return false
}
