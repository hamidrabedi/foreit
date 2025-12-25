package serializer

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Serializer provides serialization utilities
type Serializer struct{}

// ToJSON converts any value to JSON
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON converts JSON to a value
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ToMap converts a struct to a map
func ToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	
	return result, nil
}

// FromMap converts a map to a struct
func FromMap(data map[string]interface{}, v interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(jsonData, v)
}

// GetJSONFields returns JSON field names from a struct
func GetJSONFields(v interface{}) []string {
	t := reflect.TypeOf(v)
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
		}
	}
	
	return fields
}

// OmitFields removes specified fields from a map
func OmitFields(data map[string]interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{})
	omitMap := make(map[string]bool)
	for _, field := range fields {
		omitMap[field] = true
	}
	
	for k, v := range data {
		if !omitMap[k] {
			result[k] = v
		}
	}
	
	return result
}

// PickFields keeps only specified fields in a map
func PickFields(data map[string]interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{})
	pickMap := make(map[string]bool)
	for _, field := range fields {
		pickMap[field] = true
	}
	
	for k, v := range data {
		if pickMap[k] {
			result[k] = v
		}
	}
	
	return result
}

