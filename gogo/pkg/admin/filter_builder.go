package admin

import (
	"github.com/gogo/pkg/models"
)

// FilterBuilder helps build type-safe filters from admin FilterSpec
// This bridges admin FieldRef (for display) with models.FieldRef (for queries)
type FilterBuilder[T any] struct {
	fieldMap map[string]models.FieldRef[interface{}, T]
}

// NewFilterBuilder creates a new filter builder
func NewFilterBuilder[T any]() *FilterBuilder[T] {
	return &FilterBuilder[T]{
		fieldMap: make(map[string]models.FieldRef[interface{}, T]),
	}
}

// RegisterField maps an admin field name to a models FieldRef
func (b *FilterBuilder[T]) RegisterField(adminFieldName string, fieldRef models.FieldRef[interface{}, T]) *FilterBuilder[T] {
	b.fieldMap[adminFieldName] = fieldRef
	return b
}

// BuildCondition builds a type-safe condition from a FilterSpec and value
func (b *FilterBuilder[T]) BuildCondition(spec FilterSpec[T], value interface{}) *models.Condition[T] {
	fieldRef, ok := b.fieldMap[spec.Field.Name()]
	if !ok {
		return nil
	}
	
	// Build condition based on filter type
	switch spec.Type {
	case FilterTypeBoolean:
		if boolVal, ok := value.(bool); ok {
			return models.NewCondition[T](fieldRef.ApplyEq(boolVal))
		}
	case FilterTypeChoice:
		return models.NewCondition[T](fieldRef.ApplyEq(value))
	case FilterTypeList:
		if values, ok := value.([]interface{}); ok {
			return models.NewCondition[T](fieldRef.ApplyIn(values))
		}
	}
	
	return nil
}

// BuildSearchCondition builds a search condition for string fields
func (b *FilterBuilder[T]) BuildSearchCondition(fieldName, searchTerm string) *models.Condition[T] {
	fieldRef, ok := b.fieldMap[fieldName]
	if !ok {
		return nil
	}
	
	// For string fields, use IContains
	// This requires the field to be a StringFieldRef
	// In a full implementation, we'd check the field type
	return models.NewCondition[T](func(query interface{}) interface{} {
		// This is a placeholder - would use StringFieldRef.ApplyIContains
		return query
	})
}

