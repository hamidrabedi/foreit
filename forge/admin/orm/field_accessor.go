package orm

import (
	"fmt"
	"reflect"

	"github.com/forgego/forge/orm"
)

// FieldAccessor provides type-safe field access for admin operations
type FieldAccessor[T any] struct {
	accessor *orm.FieldAccessor[T]
	schema   *orm.ModelSchema
}

// NewFieldAccessor creates a new admin field accessor
func NewFieldAccessor[T any]() (*FieldAccessor[T], error) {
	accessor, err := orm.NewFieldAccessor[T]()
	if err != nil {
		return nil, err
	}

	schema, err := orm.GetModelSchema[T]()
	if err != nil {
		return nil, err
	}

	return &FieldAccessor[T]{
		accessor: accessor,
		schema:   schema,
	}, nil
}

// GetFieldValue gets a field value from an instance (type-safe)
// V is the field type
func GetFieldValue[T any, V any](fa *FieldAccessor[T], instance *T, fieldName string) (V, error) {
	var zero V

	// Validate field exists
	fieldInfo := fa.schema.GetField(fieldName)
	if fieldInfo == nil {
		return zero, fmt.Errorf("field %s not found", fieldName)
	}

	// Use reflection to get the value
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	fieldValue := instanceValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return zero, fmt.Errorf("field %s is not valid", fieldName)
	}

	if !fieldValue.CanInterface() {
		return zero, fmt.Errorf("field %s cannot be accessed", fieldName)
	}

	// Convert to type V
	val := fieldValue.Interface()
	if v, ok := val.(V); ok {
		return v, nil
	}

	// Try type conversion if types don't match exactly
	valType := reflect.TypeOf(val)
	expectedType := reflect.TypeOf(zero)
	if valType.ConvertibleTo(expectedType) {
		converted := reflect.ValueOf(val).Convert(expectedType)
		return converted.Interface().(V), nil
	}

	return zero, fmt.Errorf("field %s has type %v, cannot convert to %v", fieldName, valType, expectedType)
}

// SetFieldValue sets a field value on an instance (type-safe)
func SetFieldValue[T any, V any](fa *FieldAccessor[T], instance *T, fieldName string, value V) error {
	// Validate field exists
	fieldInfo := fa.schema.GetField(fieldName)
	if fieldInfo == nil {
		return fmt.Errorf("field %s not found", fieldName)
	}

	// Use reflection to set the value
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	fieldValue := instanceValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return fmt.Errorf("field %s is not valid", fieldName)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	// Convert value to field type if needed
	valType := reflect.TypeOf(value)
	fieldType := fieldValue.Type()

	if valType.AssignableTo(fieldType) {
		fieldValue.Set(reflect.ValueOf(value))
		return nil
	}

	// Try type conversion
	if valType.ConvertibleTo(fieldType) {
		converted := reflect.ValueOf(value).Convert(fieldType)
		fieldValue.Set(converted)
		return nil
	}

	return fmt.Errorf("cannot assign value of type %v to field %s of type %v", valType, fieldName, fieldType)
}

// GetFieldExpression gets a type-safe field expression for queries
func GetFieldExpression[T any, V any](fa *FieldAccessor[T], fieldName string) orm.FieldExpression[V] {
	return orm.FieldFor[T, V](fa.accessor, fieldName)
}

// Accessor returns the underlying ORM field accessor
func (fa *FieldAccessor[T]) Accessor() *orm.FieldAccessor[T] {
	return fa.accessor
}

// Schema returns the model schema
func (fa *FieldAccessor[T]) Schema() *orm.ModelSchema {
	return fa.schema
}
