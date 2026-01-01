package admin

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/forgego/forge/schema"
	validation "github.com/forgego/forge/validate"
)

// ValidateForm validates form data against admin config
func ValidateForm[T any](admin *Admin[T], instance *T, formData FormData, isNew bool) map[string]string {
	errors := make(map[string]string)

	validator := validation.NewValidator()
	fieldValidator := validation.NewFieldValidator(validator)

	// Get all fields
	discoveredFields := admin.Fields()

	for _, fieldInfo := range discoveredFields {
		// Skip if AutoNow, AutoNowAdd, AutoIncrement
		if fieldInfo.AutoNow || fieldInfo.AutoNowAdd || fieldInfo.AutoIncrement {
			continue
		}

		// Skip if explicitly ReadOnly in config or schema
		if fieldInfo.ReadOnly {
			continue
		}

		value, exists := formData[fieldInfo.Name]

		if !exists {
			value = nil
		} else {
			// Convert string values to appropriate types for validation
			value = convertValue(value, fieldInfo.Type)
		}

		if err := fieldValidator.ValidateField(fieldInfo.SchemaField, value); err != nil {
			msg := err.Error()
			// Strip field name prefix if present
			prefix := fmt.Sprintf("%s: ", fieldInfo.Name)
			if len(msg) > len(prefix) && msg[:len(prefix)] == prefix {
				msg = msg[len(prefix):]
			}
			errors[fieldInfo.Name] = msg
		}
	}

	return errors
}

// convertValue converts string values to appropriate types for validation
func convertValue(value interface{}, fieldType schema.FieldType) interface{} {
	if value == nil {
		return nil
	}

	strVal, isString := value.(string)
	if !isString {
		return value
	}

	switch fieldType {
	case schema.TypeInt64, schema.TypeInt32:
		if i, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			if fieldType == schema.TypeInt32 {
				return int32(i)
			}
			return i
		}
	case schema.TypeFloat64, schema.TypeFloat32, schema.TypeDecimal:
		if f, err := strconv.ParseFloat(strVal, 64); err == nil {
			if fieldType == schema.TypeFloat32 {
				return float32(f)
			}
			return f
		}
	case schema.TypeBool:
		if b, err := strconv.ParseBool(strVal); err == nil {
			return b
		}
	}

	return value
}

// getFieldNameFromFormField extracts field name from FormField
func getFieldNameFromFormField[T any](field FormField[T]) string {
	return field.Name()
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
