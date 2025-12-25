package field

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func AutoStringField[T any](fieldName string, opts ...Option) Field[T, string] {
	getter, setter := makeAccessors[T, string](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, string](fieldName, columnName, getter, setter, opts...)
}

func AutoIntField[T any](fieldName string, opts ...Option) Field[T, int] {
	getter, setter := makeAccessors[T, int](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, int](fieldName, columnName, getter, setter, opts...)
}

func AutoInt64Field[T any](fieldName string, opts ...Option) Field[T, int64] {
	getter, setter := makeAccessors[T, int64](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, int64](fieldName, columnName, getter, setter, opts...)
}

func AutoBoolField[T any](fieldName string, opts ...Option) Field[T, bool] {
	getter, setter := makeAccessors[T, bool](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, bool](fieldName, columnName, getter, setter, opts...)
}

func AutoTimeField[T any](fieldName string, opts ...Option) Field[T, time.Time] {
	getter, setter := makeAccessors[T, time.Time](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, time.Time](fieldName, columnName, getter, setter, opts...)
}

func AutoTimePtrField[T any](fieldName string, opts ...Option) Field[T, *time.Time] {
	getter, setter := makeAccessors[T, *time.Time](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, *time.Time](fieldName, columnName, getter, setter, opts...)
}

func AutoFloat64Field[T any](fieldName string, opts ...Option) Field[T, float64] {
	getter, setter := makeAccessors[T, float64](fieldName)
	columnName := toSnakeCase(fieldName)
	opts = append([]Option{Column(columnName)}, opts...)
	return NewField[T, float64](fieldName, columnName, getter, setter, opts...)
}

func makeAccessors[T any, V any](fieldName string) (func(*T) *V, func(*T, V)) {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	field, found := typ.FieldByName(fieldName)
	if !found {
		panic(fmt.Sprintf("field %s not found in type %s", fieldName, typ.Name()))
	}
	
	fieldType := field.Type
	var valueType reflect.Type
	if fieldType.Kind() == reflect.Ptr {
		valueType = fieldType.Elem()
	} else {
		valueType = fieldType
	}
	
	var expectedType reflect.Type
	var zeroV V
	expectedType = reflect.TypeOf(zeroV)
	if expectedType.Kind() == reflect.Ptr {
		expectedType = expectedType.Elem()
	}
	
	if valueType != expectedType {
		panic(fmt.Sprintf("field %s type %s does not match expected type %s", fieldName, valueType, expectedType))
	}
	
	getter := func(t *T) *V {
		val := reflect.ValueOf(t)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		fieldVal := val.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			return nil
		}
		
		if fieldType.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				return nil
			}
			return fieldVal.Interface().(*V)
		}
		
		ptr := reflect.New(fieldType)
		ptr.Elem().Set(fieldVal)
		return ptr.Interface().(*V)
	}
	
		setter := func(t *T, v V) {
		val := reflect.ValueOf(t)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		fieldVal := val.FieldByName(fieldName)
		if !fieldVal.CanSet() {
			return
		}
		
		if fieldType.Kind() == reflect.Ptr {
			vVal := reflect.ValueOf(v)
			if !vVal.IsValid() || vVal.IsNil() {
				fieldVal.Set(reflect.Zero(fieldType))
			} else {
				ptr := reflect.New(valueType)
				ptr.Elem().Set(vVal)
				fieldVal.Set(ptr)
			}
		} else {
			fieldVal.Set(reflect.ValueOf(v))
		}
	}
	
	return getter, setter
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

