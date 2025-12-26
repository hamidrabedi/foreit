package admin

import (
	"fmt"
	"reflect"
)

// FormFieldsetData contains fieldset data for form rendering
type FormFieldsetData struct {
	Name      string
	Fields    []FormField
	Collapsed bool
	Classes   []string
}

// GenerateFieldsets generates fieldsets from model configuration
func GenerateFieldsets(model *AdminModel, instance interface{}, isCreate bool) []FormFieldsetData {
	fieldsets := GetFieldsets(model)
	
	// If no fieldsets configured, return all fields in a single fieldset
	if len(fieldsets) == 0 {
		allFields := generateFormFields(model, instance, isCreate)
		return []FormFieldsetData{
			{
				Name:   "",
				Fields: allFields,
			},
		}
	}

	// Generate fields for all fieldsets
	allFields := generateFormFields(model, instance, isCreate)
	fieldMap := make(map[string]FormField)
	for _, field := range allFields {
		fieldMap[field.Name] = field
	}

	result := make([]FormFieldsetData, 0, len(fieldsets))
	for _, fieldset := range fieldsets {
		fsData := FormFieldsetData{
			Name:      fieldset.Name,
			Collapsed: fieldset.Collapsed,
			Classes:   fieldset.Classes,
			Fields:    make([]FormField, 0),
		}

		for _, fieldName := range fieldset.Fields {
			if field, ok := fieldMap[fieldName]; ok {
				fsData.Fields = append(fsData.Fields, field)
			}
		}

		result = append(result, fsData)
	}

	return result
}

// GetFieldValue extracts a field value from an instance
func GetFieldValue(instance interface{}, fieldName string) interface{} {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	fieldValue := instanceValue.FieldByName(fieldName)
	if fieldValue.IsValid() && fieldValue.CanInterface() {
		return fieldValue.Interface()
	}
	return nil
}

// SetFieldValue sets a field value on an instance
func SetFieldValue(instance interface{}, fieldName string, value interface{}) error {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	fieldValue := instanceValue.FieldByName(fieldName)
	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return fmt.Errorf("field %s is not valid or cannot be set", fieldName)
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Type().AssignableTo(fieldValue.Type()) {
		fieldValue.Set(valueReflect)
	} else if valueReflect.Type().ConvertibleTo(fieldValue.Type()) {
		fieldValue.Set(valueReflect.Convert(fieldValue.Type()))
	} else {
		return fmt.Errorf("cannot assign value of type %T to field %s of type %s", value, fieldName, fieldValue.Type())
	}

	return nil
}

