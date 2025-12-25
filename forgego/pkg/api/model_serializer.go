package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/forgego/forge/pkg/models"
)

type ModelSerializer[T any] struct {
	*BaseSerializer[T]
	Fields        []string
	ExcludeFields []string
	ReadOnlyFields []string
	modelDef      *models.ModelDefinition[T]
}

func NewModelSerializer[T any]() *ModelSerializer[T] {
	return &ModelSerializer[T]{
		BaseSerializer: NewBaseSerializer[T](),
		Fields:        []string{},
		ExcludeFields: []string{},
		ReadOnlyFields: []string{},
		modelDef:      nil,
	}
}

func NewModelSerializerWithDefinition[T any](modelDef *models.ModelDefinition[T]) *ModelSerializer[T] {
	serializer := &ModelSerializer[T]{
		BaseSerializer: NewBaseSerializer[T](),
		Fields:        []string{},
		ExcludeFields: []string{},
		ReadOnlyFields: []string{},
		modelDef:      modelDef,
	}
	serializer.buildFieldAdapters()
	return serializer
}

func (s *ModelSerializer[T]) WithModelDefinition(modelDef *models.ModelDefinition[T]) *ModelSerializer[T] {
	s.modelDef = modelDef
	s.buildFieldAdapters()
	return s
}

func (s *ModelSerializer[T]) buildFieldAdapters() {
	if s.modelDef == nil {
		return
	}
	
	modelFields := s.modelDef.GetFields()
	adapters := AdaptModelFields(modelFields)
	
	for fieldName, adapter := range adapters {
		// Apply field selection/exclusion filters
		if s.shouldExclude(fieldName) {
			continue
		}
		if len(s.Fields) > 0 && !s.containsField(fieldName) {
			continue
		}
		if s.isReadOnly(fieldName) {
			if fieldAdapter, ok := adapter.(*FieldAdapter); ok {
				fieldAdapter.ReadOnly(true)
			}
		}
		s.BaseSerializer.AddField(fieldName, adapter)
	}
}

func (s *ModelSerializer[T]) WithFields(fields ...string) *ModelSerializer[T] {
	s.Fields = fields
	// Rebuild field adapters if model definition is set
	if s.modelDef != nil {
		s.BaseSerializer.fields = make(map[string]Field)
		s.buildFieldAdapters()
	}
	return s
}

func (s *ModelSerializer[T]) Exclude(fields ...string) *ModelSerializer[T] {
	s.ExcludeFields = fields
	// Rebuild field adapters if model definition is set
	if s.modelDef != nil {
		s.BaseSerializer.fields = make(map[string]Field)
		s.buildFieldAdapters()
	}
	return s
}

func (s *ModelSerializer[T]) ReadOnly(fields ...string) *ModelSerializer[T] {
	s.ReadOnlyFields = fields
	// Rebuild field adapters if model definition is set
	if s.modelDef != nil {
		s.BaseSerializer.fields = make(map[string]Field)
		s.buildFieldAdapters()
	}
	return s
}

func (s *ModelSerializer[T]) ToRepresentation(obj *T) (interface{}, error) {
	if len(s.BaseSerializer.fields) > 0 {
		return s.BaseSerializer.ToRepresentation(obj)
	}
	return s.toRepresentationReflection(obj)
}

func (s *ModelSerializer[T]) toRepresentationReflection(obj *T) (interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", val.Kind())
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanInterface() {
			continue
		}

		if field.Name == "Schema" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		jsonFieldName := fieldName
		if jsonTag != "" && jsonTag != "-" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonFieldName = jsonTag[:idx]
			} else {
				jsonFieldName = jsonTag
			}
		} else {
			jsonFieldName = toSnakeCase(fieldName)
		}

		// Check exclusion and field selection using both field name and JSON field name
		if s.shouldExclude(fieldName) || s.shouldExclude(jsonFieldName) {
			continue
		}

		if len(s.Fields) > 0 && !s.containsField(fieldName) && !s.containsField(jsonFieldName) {
			continue
		}

		if s.isReadOnly(fieldName) || s.isReadOnly(jsonFieldName) {
			continue
		}

		// Use JSON field name for output
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				result[jsonFieldName] = nil
			} else {
				value := fieldVal.Elem().Interface()
				if t, ok := value.(time.Time); ok {
					result[jsonFieldName] = t.Format(time.RFC3339)
				} else {
					result[jsonFieldName] = value
				}
			}
		} else {
			value := fieldVal.Interface()
			if t, ok := value.(time.Time); ok {
				result[jsonFieldName] = t.Format(time.RFC3339)
			} else {
				result[jsonFieldName] = value
			}
		}
	}

	return result, nil
}


func (s *ModelSerializer[T]) FromCreate(body []byte) (*T, error) {
	if len(s.BaseSerializer.fields) > 0 {
		return s.BaseSerializer.FromCreate(body)
	}
	return s.fromCreateReflection(body)
}

func (s *ModelSerializer[T]) fromCreateReflection(body []byte) (*T, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	var zero T
	obj := &zero
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		if field.Name == "Schema" {
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
		} else {
			fieldName = toSnakeCase(fieldName)
		}

		if s.isReadOnly(fieldName) {
			continue
		}

		value, exists := data[fieldName]
		if !exists {
			camelName := field.Name
			if value, exists = data[camelName]; !exists {
				continue
			}
		}

		if err := s.setFieldValue(fieldVal, value); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return obj, nil
}

func (s *ModelSerializer[T]) FromUpdate(obj *T, body []byte) error {
	if len(s.BaseSerializer.fields) > 0 {
		return s.BaseSerializer.FromUpdate(obj, body)
	}
	return s.fromUpdateReflection(obj, body)
}

func (s *ModelSerializer[T]) fromUpdateReflection(obj *T, body []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		if field.Name == "Schema" {
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
		} else {
			fieldName = toSnakeCase(fieldName)
		}

		if s.isReadOnly(fieldName) {
			continue
		}

		value, exists := data[fieldName]
		if !exists {
			camelName := field.Name
			if value, exists = data[camelName]; !exists {
				continue
			}
		}

		if err := s.setFieldValue(fieldVal, value); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldName, err)
		}
	}

	return nil
}

func (s *ModelSerializer[T]) shouldExclude(fieldName string) bool {
	for _, excluded := range s.ExcludeFields {
		if excluded == fieldName {
			return true
		}
	}
	return false
}

func (s *ModelSerializer[T]) containsField(fieldName string) bool {
	for _, field := range s.Fields {
		if field == fieldName {
			return true
		}
	}
	return false
}

func (s *ModelSerializer[T]) isReadOnly(fieldName string) bool {
	for _, readonly := range s.ReadOnlyFields {
		if readonly == fieldName {
			return true
		}
	}
	return false
}

func (s *ModelSerializer[T]) setFieldValue(fieldVal reflect.Value, value interface{}) error {
	if !fieldVal.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	fieldType := fieldVal.Type()

	if value == nil {
		if fieldType.Kind() == reflect.Ptr {
			fieldVal.Set(reflect.Zero(fieldType))
			return nil
		}
		return fmt.Errorf("cannot set nil to non-pointer field")
	}

	valueType := reflect.TypeOf(value)

	if valueType.AssignableTo(fieldType) {
		fieldVal.Set(reflect.ValueOf(value))
		return nil
	}

	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		elemType := fieldType.Elem()
		if valueType.AssignableTo(elemType) {
			newVal := reflect.New(elemType)
			newVal.Elem().Set(reflect.ValueOf(value))
			fieldVal.Set(newVal)
			return nil
		}
		// Handle time.Time pointer from string
		if elemType == reflect.TypeOf(time.Time{}) {
			if str, ok := value.(string); ok {
				t, err := time.Parse(time.RFC3339, str)
				if err != nil {
					return fmt.Errorf("invalid time format: %w", err)
				}
				newVal := reflect.New(elemType)
				newVal.Elem().Set(reflect.ValueOf(t))
				fieldVal.Set(newVal)
				return nil
			}
		}
		// Try conversion for pointer types
		if valueType.ConvertibleTo(elemType) {
			newVal := reflect.New(elemType)
			newVal.Elem().Set(reflect.ValueOf(value).Convert(elemType))
			fieldVal.Set(newVal)
			return nil
		}
	}

	// Handle time.Time from string
	if fieldType == reflect.TypeOf(time.Time{}) {
		if str, ok := value.(string); ok {
			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return fmt.Errorf("invalid time format: %w", err)
			}
			fieldVal.Set(reflect.ValueOf(t))
			return nil
		}
	}

	// Handle numeric conversions (JSON unmarshals numbers as float64)
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

	// Handle string conversions
	if fieldType.Kind() == reflect.String {
		if str, ok := value.(string); ok {
			fieldVal.SetString(str)
			return nil
		}
		if valueType.ConvertibleTo(fieldType) {
			fieldVal.Set(reflect.ValueOf(value).Convert(fieldType))
			return nil
		}
	}

	// Handle bool conversions
	if fieldType.Kind() == reflect.Bool {
		if boolVal, ok := value.(bool); ok {
			fieldVal.SetBool(boolVal)
			return nil
		}
	}

	// Try type conversion as last resort
	if valueType.ConvertibleTo(fieldType) {
		fieldVal.Set(reflect.ValueOf(value).Convert(fieldType))
		return nil
	}

	return fmt.Errorf("cannot convert %v to %s", value, fieldType)
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}


