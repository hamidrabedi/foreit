package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type BaseSerializer[T any] struct {
	fields map[string]Field
	validators []Validator
}

func NewBaseSerializer[T any]() *BaseSerializer[T] {
	return &BaseSerializer[T]{
		fields: make(map[string]Field),
		validators: []Validator{},
	}
}

func (s *BaseSerializer[T]) AddField(name string, field Field) *BaseSerializer[T] {
	s.fields[name] = field
	return s
}

func (s *BaseSerializer[T]) AddValidator(validator Validator) *BaseSerializer[T] {
	s.validators = append(s.validators, validator)
	return s
}

func (s *BaseSerializer[T]) ToRepresentation(obj *T) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	if len(s.fields) == 0 {
		return s.toRepresentationAuto(obj, val)
	}

	for fieldName, field := range s.fields {
		if field.GetReadOnly() {
			continue
		}

		source := field.GetSource()
		if source == "" {
			source = fieldName
		}

		fieldValue, err := s.getFieldValue(val, source)
		if err != nil {
			continue
		}

		repr, err := field.ToRepresentation(fieldValue)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldName, err)
		}

		if methodField, ok := field.(*SerializerMethodField); ok {
			methodValue := reflect.ValueOf(obj)
			if methodValue.Kind() == reflect.Ptr {
				methodValue = methodValue.Elem()
			}
			method := reflect.ValueOf(obj).MethodByName(methodField.Method)
			if method.IsValid() {
				results := method.Call([]reflect.Value{})
				if len(results) > 0 {
					repr = results[0].Interface()
				}
			}
		}

		result[fieldName] = repr
	}

	return result, nil
}

func (s *BaseSerializer[T]) toRepresentationAuto(obj *T, val reflect.Value) (interface{}, error) {
	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanInterface() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName := field.Name
		if jsonTag != "" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldName = jsonTag[:idx]
			} else {
				fieldName = jsonTag
			}
		}

		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				result[fieldName] = nil
			} else {
				result[fieldName] = fieldVal.Elem().Interface()
			}
		} else {
			result[fieldName] = fieldVal.Interface()
		}
	}

	return result, nil
}

func (s *BaseSerializer[T]) FromCreate(body []byte) (*T, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if err := s.Validate(data); err != nil {
		return nil, err
	}

	var zero T
	obj := &zero
	val := reflect.ValueOf(obj).Elem()

	if len(s.fields) == 0 {
		return s.fromCreateAuto(data, obj, val)
	}

	for fieldName, field := range s.fields {
		if field.GetReadOnly() || field.GetWriteOnly() {
			continue
		}

		source := field.GetSource()
		if source == "" {
			source = fieldName
		}

		value, exists := data[fieldName]
		if !exists {
			if field.GetRequired() {
				return nil, fmt.Errorf("field %s is required", fieldName)
			}
			if field.GetDefault() != nil {
				value = field.GetDefault()
			} else {
				continue
			}
		}

		internalValue, err := field.ToInternalValue(value)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldName, err)
		}

		if err := s.setFieldValue(val, source, internalValue); err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldName, err)
		}
	}

	return obj, nil
}

func (s *BaseSerializer[T]) fromCreateAuto(data map[string]interface{}, obj *T, val reflect.Value) (*T, error) {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		if jsonTag != "" && jsonTag != "-" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldName = jsonTag[:idx]
			} else {
				fieldName = jsonTag
			}
		}

		value, exists := data[fieldName]
		if !exists {
			continue
		}

		if err := s.setFieldValueSimple(fieldVal, value); err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldName, err)
		}
	}

	return obj, nil
}

func (s *BaseSerializer[T]) FromUpdate(obj *T, body []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := s.Validate(data); err != nil {
		return err
	}

	val := reflect.ValueOf(obj).Elem()

	if len(s.fields) == 0 {
		return s.fromUpdateAuto(data, val)
	}

	for fieldName, field := range s.fields {
		if field.GetReadOnly() {
			continue
		}

		value, exists := data[fieldName]
		if !exists {
			continue
		}

		internalValue, err := field.ToInternalValue(value)
		if err != nil {
			return fmt.Errorf("field %s: %w", fieldName, err)
		}

		source := field.GetSource()
		if source == "" {
			source = fieldName
		}

		if err := s.setFieldValue(val, source, internalValue); err != nil {
			return fmt.Errorf("field %s: %w", fieldName, err)
		}
	}

	return nil
}

func (s *BaseSerializer[T]) fromUpdateAuto(data map[string]interface{}, val reflect.Value) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		if jsonTag != "" && jsonTag != "-" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				fieldName = jsonTag[:idx]
			} else {
				fieldName = jsonTag
			}
		}

		value, exists := data[fieldName]
		if !exists {
			continue
		}

		if err := s.setFieldValueSimple(fieldVal, value); err != nil {
			return fmt.Errorf("field %s: %w", fieldName, err)
		}
	}

	return nil
}

func (s *BaseSerializer[T]) Validate(data map[string]interface{}) error {
	for _, validator := range s.validators {
		if err := validator.Validate(data); err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseSerializer[T]) getFieldValue(val reflect.Value, fieldName string) (interface{}, error) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			fieldVal := val.Field(i)
			if fieldVal.CanInterface() {
				if fieldVal.Kind() == reflect.Ptr {
					if fieldVal.IsNil() {
						return nil, nil
					}
					return fieldVal.Elem().Interface(), nil
				}
				return fieldVal.Interface(), nil
			}
		}
	}
	return nil, fmt.Errorf("field %s not found", fieldName)
}

func (s *BaseSerializer[T]) setFieldValue(val reflect.Value, fieldName string, value interface{}) error {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == fieldName {
			fieldVal := val.Field(i)
			if fieldVal.CanSet() {
				return s.setFieldValueSimple(fieldVal, value)
			}
		}
	}
	return fmt.Errorf("field %s not found or cannot be set", fieldName)
}

func (s *BaseSerializer[T]) setFieldValueSimple(fieldVal reflect.Value, value interface{}) error {
	if value == nil {
		if fieldVal.Kind() == reflect.Ptr {
			fieldVal.Set(reflect.Zero(fieldVal.Type()))
			return nil
		}
		return fmt.Errorf("cannot set nil to non-pointer field")
	}

	valueType := reflect.TypeOf(value)
	fieldType := fieldVal.Type()

	// Direct assignment
	if valueType.AssignableTo(fieldType) {
		fieldVal.Set(reflect.ValueOf(value))
		return nil
	}

	// Pointer types
	if fieldType.Kind() == reflect.Ptr {
		elemType := fieldType.Elem()
		if valueType.AssignableTo(elemType) {
			newVal := reflect.New(elemType)
			newVal.Elem().Set(reflect.ValueOf(value))
			fieldVal.Set(newVal)
			return nil
		}
		// Try conversion for pointer types
		if valueType.ConvertibleTo(elemType) {
			newVal := reflect.New(elemType)
			newVal.Elem().Set(reflect.ValueOf(value).Convert(elemType))
			fieldVal.Set(newVal)
			return nil
		}
	}

	// Numeric conversions (JSON unmarshals numbers as float64)
	if fieldType.Kind() == reflect.Int || fieldType.Kind() == reflect.Int8 ||
		fieldType.Kind() == reflect.Int16 || fieldType.Kind() == reflect.Int32 ||
		fieldType.Kind() == reflect.Int64 {
		if intVal, ok := value.(int); ok {
			fieldVal.SetInt(int64(intVal))
			return nil
		}
		if floatVal, ok := value.(float64); ok {
			fieldVal.SetInt(int64(floatVal))
			return nil
		}
		if int64Val, ok := value.(int64); ok {
			fieldVal.SetInt(int64Val)
			return nil
		}
	}

	if fieldType.Kind() == reflect.Uint || fieldType.Kind() == reflect.Uint8 ||
		fieldType.Kind() == reflect.Uint16 || fieldType.Kind() == reflect.Uint32 ||
		fieldType.Kind() == reflect.Uint64 {
		if intVal, ok := value.(int); ok {
			fieldVal.SetUint(uint64(intVal))
			return nil
		}
		if floatVal, ok := value.(float64); ok {
			fieldVal.SetUint(uint64(floatVal))
			return nil
		}
		if uint64Val, ok := value.(uint64); ok {
			fieldVal.SetUint(uint64Val)
			return nil
		}
	}

	// String conversions
	if fieldType.Kind() == reflect.String {
		if str, ok := value.(string); ok {
			fieldVal.SetString(str)
			return nil
		}
	}

	// Bool conversions
	if fieldType.Kind() == reflect.Bool {
		if boolVal, ok := value.(bool); ok {
			fieldVal.SetBool(boolVal)
			return nil
		}
	}

	// Type conversion as last resort
	if valueType.ConvertibleTo(fieldType) {
		fieldVal.Set(reflect.ValueOf(value).Convert(fieldType))
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", value, fieldType)
}

