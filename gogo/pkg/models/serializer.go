package models

import (
	"encoding/json"
	"reflect"
	"time"
)

// Serializer provides Pydantic-inspired serialization
type Serializer interface {
	// Serialize converts a model to JSON
	Serialize(model interface{}) ([]byte, error)
	
	// Deserialize converts JSON to a model
	Deserialize(data []byte, model interface{}) error
	
	// ToDict converts a model to a map
	ToDict(model interface{}) (map[string]interface{}, error)
	
	// FromDict creates a model from a map
	FromDict(data map[string]interface{}, model interface{}) error
}

// JSONSerializer implements JSON serialization
type JSONSerializer struct {
	excludeFields map[string]bool
	includeFields map[string]bool
}

// NewJSONSerializer creates a new JSON serializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{
		excludeFields: make(map[string]bool),
		includeFields: make(map[string]bool),
	}
}

// ExcludeFields excludes fields from serialization
func (s *JSONSerializer) ExcludeFields(fields ...string) *JSONSerializer {
	for _, field := range fields {
		s.excludeFields[field] = true
	}
	return s
}

// IncludeFields includes only specified fields
func (s *JSONSerializer) IncludeFields(fields ...string) *JSONSerializer {
	for _, field := range fields {
		s.includeFields[field] = true
	}
	return s
}

// Serialize converts a model to JSON
func (s *JSONSerializer) Serialize(model interface{}) ([]byte, error) {
	dict, err := s.ToDict(model)
	if err != nil {
		return nil, err
	}
	return json.Marshal(dict)
}

// Deserialize converts JSON to a model
func (s *JSONSerializer) Deserialize(data []byte, model interface{}) error {
	var dict map[string]interface{}
	if err := json.Unmarshal(data, &dict); err != nil {
		return err
	}
	return s.FromDict(dict, model)
}

// ToDict converts a model to a map
func (s *JSONSerializer) ToDict(model interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	typ := val.Type()
	result := make(map[string]interface{})
	
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}
		
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		
		fieldName := getJSONFieldName(jsonTag)
		
		// Apply include/exclude filters
		if len(s.includeFields) > 0 && !s.includeFields[fieldName] {
			continue
		}
		if s.excludeFields[fieldName] {
			continue
		}
		
		// Handle special types
		value := fieldVal.Interface()
		if timeVal, ok := value.(time.Time); ok {
			result[fieldName] = timeVal.Format(time.RFC3339)
		} else if timePtr, ok := value.(*time.Time); ok && timePtr != nil {
			result[fieldName] = timePtr.Format(time.RFC3339)
		} else {
			result[fieldName] = value
		}
	}
	
	return result, nil
}

// FromDict creates a model from a map
func (s *JSONSerializer) FromDict(data map[string]interface{}, model interface{}) error {
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr {
		return ErrNotAPointer
	}
	
	val = val.Elem()
	typ := val.Type()
	
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		if !fieldVal.CanSet() {
			continue
		}
		
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		
		fieldName := getJSONFieldName(jsonTag)
		value, ok := data[fieldName]
		if !ok {
			continue
		}
		
		if err := setFieldValue(fieldVal, value); err != nil {
			return err
		}
	}
	
	return nil
}

func getJSONFieldName(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

func setFieldValue(fieldVal reflect.Value, value interface{}) error {
	if !fieldVal.CanSet() {
		return ErrCannotSetField
	}
	
	fieldType := fieldVal.Type()
	
	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		if value == nil {
			fieldVal.Set(reflect.Zero(fieldType))
			return nil
		}
		// Create new pointer and set value
		elemType := fieldType.Elem()
		newVal := reflect.New(elemType)
		if err := setFieldValue(newVal.Elem(), value); err != nil {
			return err
		}
		fieldVal.Set(newVal)
		return nil
	}
	
	// Type conversion
	val := reflect.ValueOf(value)
	if val.Type().AssignableTo(fieldType) {
		fieldVal.Set(val)
		return nil
	}
	
	// Try conversion
	if val.Type().ConvertibleTo(fieldType) {
		fieldVal.Set(val.Convert(fieldType))
		return nil
	}
	
	return ErrTypeMismatch
}

var (
	ErrNotAPointer   = errors.New("model must be a pointer")
	ErrCannotSetField = errors.New("cannot set field")
	ErrTypeMismatch  = errors.New("type mismatch")
)

