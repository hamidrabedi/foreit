package endpoints

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Serializer handles request/response transformation
type Serializer[T any] interface {
	// Validate validates input data
	Validate(data map[string]interface{}) error
	
	// ToRepresentation converts a model to API representation
	ToRepresentation(instance *T) (map[string]interface{}, error)
	
	// ToInternalValue converts API data to model
	ToInternalValue(data map[string]interface{}) (*T, error)
	
	// ManyToRepresentation converts multiple models
	ManyToRepresentation(instances []*T) ([]map[string]interface{}, error)
}

// ModelSerializer provides a default serializer implementation
type ModelSerializer[T any] struct {
	fields map[string]FieldInfo
}

// FieldInfo contains information about a field
type FieldInfo struct {
	Name     string
	Type     reflect.Type
	Required bool
	ReadOnly bool
}

// NewModelSerializer creates a new model serializer
func NewModelSerializer[T any]() Serializer[T] {
	fields := extractFields[T]()
	return &ModelSerializer[T]{
		fields: fields,
	}
}

// extractFields extracts field information from a type
func extractFields[T any]() map[string]FieldInfo {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	fields := make(map[string]FieldInfo)
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}
		
		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		
		// Parse JSON tag
		name := jsonTag
		if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
			name = jsonTag[:commaIdx]
		}
		
		// Check for required tag
		required := field.Tag.Get("required") == "true"
		readOnly := field.Tag.Get("readonly") == "true"
		
		fields[name] = FieldInfo{
			Name:     name,
			Type:     field.Type,
			Required: required,
			ReadOnly: readOnly,
		}
	}
	
	return fields
}

// Validate validates input data
func (s *ModelSerializer[T]) Validate(data map[string]interface{}) error {
	for name, info := range s.fields {
		if info.Required {
			if _, ok := data[name]; !ok {
				return fmt.Errorf("field %s is required", name)
			}
		}
	}
	return nil
}

// ToRepresentation converts a model to API representation
func (s *ModelSerializer[T]) ToRepresentation(instance *T) (map[string]interface{}, error) {
	// Use JSON marshaling/unmarshaling as a simple approach
	// In production, this would use Ent field descriptors
	data, err := json.Marshal(instance)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// ToInternalValue converts API data to model
func (s *ModelSerializer[T]) ToInternalValue(data map[string]interface{}) (*T, error) {
	// Validate first
	if err := s.Validate(data); err != nil {
		return nil, err
	}
	
	// Convert to JSON and back to model
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	var result T
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// ManyToRepresentation converts multiple models
func (s *ModelSerializer[T]) ManyToRepresentation(instances []*T) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, len(instances))
	for i, instance := range instances {
		repr, err := s.ToRepresentation(instance)
		if err != nil {
			return nil, err
		}
		results[i] = repr
	}
	return results, nil
}

