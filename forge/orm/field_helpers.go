package orm

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldFor creates a type-safe field expression
// T is the model type, V is the field type
func FieldFor[T any, V any](fa *FieldAccessor[T], name string) FieldExpression[V] {
	// Validate field exists and has correct type
	fieldInfo := fa.schema.GetField(name)
	fieldName := name
	if fieldInfo == nil {
		for i := range fa.schema.Fields {
			field := &fa.schema.Fields[i]
			if strings.EqualFold(field.Name, name) || strings.EqualFold(field.DBColumn, name) {
				fieldInfo = field
				fieldName = field.DBColumn
				if fieldName == "" {
					fieldName = field.Name
				}
				break
			}
		}
	}
	if fieldInfo == nil {
		panic(fmt.Sprintf("field %s not found on model", name))
	}

	expectedType := reflect.TypeOf((*V)(nil)).Elem()
	if fieldInfo.Type != expectedType {
		panic(fmt.Sprintf("field %s has type %v, not %v", name, fieldInfo.Type, expectedType))
	}

	return NewField[V](fieldName, fa.table)
}



