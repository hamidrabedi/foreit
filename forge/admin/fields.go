package admin

import (
	"fmt"
	"reflect"
	"time"

	"github.com/forgego/forge/schema"
)

// FieldExpr represents a type-safe field expression
type FieldExpr[T any, F any] struct {
	name     string
	getter   func(*T) F
	setter   func(*T, F)
	fieldDef schema.Field
}

// NewFieldExpr creates a new field expression
func NewFieldExpr[T any, F any](
	name string,
	getter func(*T) F,
	setter func(*T, F),
	fieldDef schema.Field,
) FieldExpr[T, F] {
	return FieldExpr[T, F]{
		name:     name,
		getter:   getter,
		setter:   setter,
		fieldDef: fieldDef,
	}
}

// Name returns the field name
func (f FieldExpr[T, F]) Name() string {
	return f.name
}

// Get gets the field value from an instance
func (f FieldExpr[T, F]) Get(instance *T) F {
	if f.getter != nil {
		return f.getter(instance)
	}
	// Fallback to reflection
	return f.getByReflection(instance)
}

// Set sets the field value on an instance
func (f FieldExpr[T, F]) Set(instance *T, value F) {
	if f.setter != nil {
		f.setter(instance, value)
		return
	}
	// Fallback to reflection
	f.setByReflection(instance, value)
}

// getByReflection gets field value using reflection
func (f FieldExpr[T, F]) getByReflection(instance *T) F {
	var zero F
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	fieldValue := instanceValue.FieldByName(f.name)
	if !fieldValue.IsValid() {
		return zero
	}

	if fieldValue.CanInterface() {
		if val, ok := fieldValue.Interface().(F); ok {
			return val
		}
	}

	return zero
}

// setByReflection sets field value using reflection
func (f FieldExpr[T, F]) setByReflection(instance *T, value F) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	fieldValue := instanceValue.FieldByName(f.name)
	if fieldValue.IsValid() && fieldValue.CanSet() {
		valueReflect := reflect.ValueOf(value)
		if valueReflect.Type().AssignableTo(fieldValue.Type()) {
			fieldValue.Set(valueReflect)
		} else if valueReflect.Type().ConvertibleTo(fieldValue.Type()) {
			fieldValue.Set(valueReflect.Convert(fieldValue.Type()))
		}
	}
}

// Format formats the field value as a string
func (f FieldExpr[T, F]) Format(instance *T) string {
	value := f.Get(instance)
	return fmt.Sprintf("%v", value)
}

// StringFieldExpr is a convenience type for string fields
type StringFieldExpr[T any] = FieldExpr[T, string]

// Int64FieldExpr is a convenience type for int64 fields
type Int64FieldExpr[T any] = FieldExpr[T, int64]

// BoolFieldExpr is a convenience type for bool fields
type BoolFieldExpr[T any] = FieldExpr[T, bool]

// TimeFieldExpr is a convenience type for time.Time fields
type TimeFieldExpr[T any] = FieldExpr[T, time.Time]

// FieldExprBuilder helps build field expressions
type FieldExprBuilder[T any] struct {
	name     string
	fieldDef schema.Field
}

// NewFieldExprBuilder creates a new field expression builder
func NewFieldExprBuilder[T any](name string, fieldDef schema.Field) *FieldExprBuilder[T] {
	return &FieldExprBuilder[T]{
		name:     name,
		fieldDef: fieldDef,
	}
}

// Build builds a field expression with getter and setter
func (b *FieldExprBuilder[T]) Build(getter interface{}, setter interface{}) FieldExpr[T, interface{}] {
	// This is a simplified version - full implementation would use reflection
	// to create proper typed getters and setters
	return FieldExpr[T, interface{}]{
		name:     b.name,
		fieldDef: b.fieldDef,
	}
}

// Helper functions to create common field expressions

// StringField creates a string field expression
func StringField[T any](name string, getter func(*T) string, setter func(*T, string)) StringFieldExpr[T] {
	return NewFieldExpr(name, getter, setter, schema.Field{
		Name: name,
		Type: schema.TypeString,
	})
}

// Int64Field creates an int64 field expression
func Int64Field[T any](name string, getter func(*T) int64, setter func(*T, int64)) Int64FieldExpr[T] {
	return NewFieldExpr(name, getter, setter, schema.Field{
		Name: name,
		Type: schema.TypeInt64,
	})
}

// BoolField creates a bool field expression
func BoolField[T any](name string, getter func(*T) bool, setter func(*T, bool)) BoolFieldExpr[T] {
	return NewFieldExpr(name, getter, setter, schema.Field{
		Name: name,
		Type: schema.TypeBool,
	})
}

// TimeField creates a time.Time field expression
func TimeField[T any](name string, getter func(*T) time.Time, setter func(*T, time.Time)) TimeFieldExpr[T] {
	return NewFieldExpr(name, getter, setter, schema.Field{
		Name: name,
		Type: schema.TypeTime,
	})
}
