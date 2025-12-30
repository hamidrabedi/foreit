package fields

import (
	"fmt"
	"reflect"

	"github.com/forgego/forge/admin/widgets"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// FieldExpr is a concrete implementation of Field that uses ORM field expressions
type FieldExpr[T any, V any] struct {
	*BaseField[T, V]
}

// NewFieldExpr creates a new field expression from schema and ORM
func NewFieldExpr[T any, V any](
	name string,
	schemaField schema.Field,
	accessor *orm.FieldAccessor[T],
	widget widgets.Widget,
) (*FieldExpr[T, V], error) {
	// Get type-safe field expression from ORM
	fieldExpr := orm.FieldFor[T, V](accessor, name)

	baseField := NewBaseField[T, V](
		name,
		schemaField,
		fieldExpr,
		accessor,
		widget,
	)

	return &FieldExpr[T, V]{
		BaseField: baseField,
	}, nil
}

// Get gets the field value from an instance (type-safe)
func (f *FieldExpr[T, V]) Get(instance *T) V {
	// Use reflection to get field value
	// In a full implementation, this would use ORM's field accessor methods
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fieldValue := val.FieldByName(f.name)
	if !fieldValue.IsValid() {
		var zero V
		return zero
	}

	if fieldValue.CanInterface() {
		if v, ok := fieldValue.Interface().(V); ok {
			return v
		}
	}

	var zero V
	return zero
}

// Set sets the field value on an instance (type-safe)
func (f *FieldExpr[T, V]) Set(instance *T, value V) error {
	// Use reflection to set field value
	// In a full implementation, this would use ORM's field accessor methods
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	fieldValue := val.FieldByName(f.name)
	if !fieldValue.IsValid() {
		return fmt.Errorf("field %s not found", f.name)
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("field %s cannot be set", f.name)
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Type().AssignableTo(fieldValue.Type()) {
		fieldValue.Set(valueReflect)
		return nil
	}

	if valueReflect.Type().ConvertibleTo(fieldValue.Type()) {
		fieldValue.Set(valueReflect.Convert(fieldValue.Type()))
		return nil
	}

	return fmt.Errorf("cannot assign value of type %v to field %s of type %v",
		valueReflect.Type(), f.name, fieldValue.Type())
}

// CreateFieldFromSchema creates a field from schema information
func CreateFieldFromSchema[T any, V any](
	schemaField schema.Field,
	accessor *orm.FieldAccessor[T],
	widgetRegistry *widgets.WidgetRegistry,
) (Field[T, V], error) {
	// Get widget for field type
	widget := widgetRegistry.GetWidgetForFieldType(schemaField.Type)
	if widget == nil {
		// Use default text widget
		widget = widgets.NewTextInput()
	}

	return NewFieldExpr[T, V](
		schemaField.Name,
		schemaField,
		accessor,
		widget,
	)
}
