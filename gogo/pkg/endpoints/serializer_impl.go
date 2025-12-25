package endpoints

import (
	"encoding/json"
	"reflect"
	"strings"
)

type ModelSerializer[T any] struct{}

func NewModelSerializer[T any]() *ModelSerializer[T] {
	return &ModelSerializer[T]{}
}

func (s *ModelSerializer[T]) ToRepresentation(obj *T) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

func (s *ModelSerializer[T]) ToInternal(data map[string]interface{}) (*T, error) {
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

func (s *ModelSerializer[T]) ToRepresentationList(objs []*T) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, len(objs))
	for i, obj := range objs {
		repr, err := s.ToRepresentation(obj)
		if err != nil {
			return nil, err
		}
		results[i] = repr
	}
	return results, nil
}

func (s *ModelSerializer[T]) Validate(data map[string]interface{}) error {
	return nil
}

func (s *ModelSerializer[T]) GetFields() []string {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	fields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonTag = jsonTag[:idx]
			}
			fields = append(fields, jsonTag)
		} else {
			fields = append(fields, field.Name)
		}
	}
	
	return fields
}

