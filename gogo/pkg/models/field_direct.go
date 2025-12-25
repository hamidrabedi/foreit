package models

import (
	"reflect"
	"sync"
	"gorm.io/gorm"
)

// FieldAccessor provides direct access to a struct field
type FieldAccessor struct {
	fieldIndex int
	fieldType  reflect.Type
	columnName string
}

// FieldRegistry caches field accessors for models
type FieldRegistry struct {
	mu     sync.RWMutex
	fields map[reflect.Type]map[string]*FieldAccessor
}

var globalRegistry = &FieldRegistry{
	fields: make(map[reflect.Type]map[string]*FieldAccessor),
}

// RegisterModel registers a model type and extracts field metadata
func RegisterModel[T any](columnMap map[string]string) {
	var zero T
	modelType := reflect.TypeOf(zero)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	
	fields := make(map[string]*FieldAccessor)
	
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		columnName := getColumnName(field, columnMap)
		
		fields[field.Name] = &FieldAccessor{
			fieldIndex: i,
			fieldType:  field.Type,
			columnName: columnName,
		}
	}
	
	globalRegistry.fields[modelType] = fields
}

// getColumnName extracts column name from struct tag or map
func getColumnName(field reflect.StructField, columnMap map[string]string) string {
	if col, ok := columnMap[field.Name]; ok {
		return col
	}
	if tag := field.Tag.Get("gorm"); tag != "" {
		if column := extractGormColumn(tag); column != "" {
			return column
		}
	}
	if tag := field.Tag.Get("db"); tag != "" {
		return tag
	}
	return toSnakeCase(field.Name)
}

// FieldRef provides type-safe field references that map directly to struct fields
type FieldRef[T any, M any] struct {
	fieldName  string
	accessor   *FieldAccessor
	modelType  reflect.Type
}

// NewFieldRef creates a type-safe field reference for a model
func NewFieldRef[T any, M any](fieldName string) *FieldRef[T, M] {
	var zero M
	modelType := reflect.TypeOf(zero)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	
	globalRegistry.mu.RLock()
	fields, ok := globalRegistry.fields[modelType]
	globalRegistry.mu.RUnlock()
	
	if !ok {
		panic("model not registered: " + modelType.Name())
	}
	
	accessor, ok := fields[fieldName]
	if !ok {
		panic("field not found: " + fieldName)
	}
	
	return &FieldRef[T, M]{
		fieldName: fieldName,
		accessor:  accessor,
		modelType: modelType,
	}
}

// Column returns the database column name
func (f *FieldRef[T, M]) Column() string {
	return f.accessor.columnName
}

// ApplyEq applies an equality condition directly to GORM query
func (f *FieldRef[T, M]) ApplyEq(value T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" = ?", value)
	}
}

// ApplyNe applies a not-equal condition
func (f *FieldRef[T, M]) ApplyNe(value T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" != ?", value)
	}
}

// ApplyGt applies a greater-than condition
func (f *FieldRef[T, M]) ApplyGt(value T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" > ?", value)
	}
}

// ApplyGte applies a greater-than-or-equal condition
func (f *FieldRef[T, M]) ApplyGte(value T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" >= ?", value)
	}
}

// ApplyLt applies a less-than condition
func (f *FieldRef[T, M]) ApplyLt(value T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" < ?", value)
	}
}

// ApplyLte applies a less-than-or-equal condition
func (f *FieldRef[T, M]) ApplyLte(value T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" <= ?", value)
	}
}

// ApplyIn applies an IN condition
func (f *FieldRef[T, M]) ApplyIn(values []T) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" IN ?", values)
	}
}

// StringFieldRef provides string-specific operations
type StringFieldRef[M any] struct {
	*FieldRef[string, M]
}

// NewStringFieldRef creates a string field reference
func NewStringFieldRef[M any](fieldName string) *StringFieldRef[M] {
	return &StringFieldRef[M]{
		FieldRef: NewFieldRef[string, M](fieldName),
	}
}

// ApplyContains applies a LIKE condition
func (f *StringFieldRef[M]) ApplyContains(value string) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" LIKE ?", "%"+value+"%")
	}
}

// ApplyIContains applies a case-insensitive LIKE condition
func (f *StringFieldRef[M]) ApplyIContains(value string) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" ILIKE ?", "%"+value+"%")
	}
}

// ApplyStartsWith applies a starts-with condition
func (f *StringFieldRef[M]) ApplyStartsWith(value string) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" LIKE ?", value+"%")
	}
}

// ApplyEndsWith applies an ends-with condition
func (f *StringFieldRef[M]) ApplyEndsWith(value string) func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName+" LIKE ?", "%"+value)
	}
}

// ApplyIsNull applies an IS NULL condition
func (f *FieldRef[T, M]) ApplyIsNull() func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName + " IS NULL")
	}
}

// ApplyIsNotNull applies an IS NOT NULL condition
func (f *FieldRef[T, M]) ApplyIsNotNull() func(*gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(f.accessor.columnName + " IS NOT NULL")
	}
}

// Condition represents a type-safe query condition
type Condition[M any] struct {
	apply func(*gorm.DB) *gorm.DB
}

// NewCondition creates a new condition
func NewCondition[M any](apply func(*gorm.DB) *gorm.DB) *Condition[M] {
	return &Condition[M]{apply: apply}
}

// Apply applies the condition to a GORM query
func (c *Condition[M]) Apply(query *gorm.DB) *gorm.DB {
	return c.apply(query)
}

// And combines two conditions with AND
func (c *Condition[M]) And(other *Condition[M]) *Condition[M] {
	return NewCondition[M](func(query *gorm.DB) *gorm.DB {
		return other.apply(c.apply(query))
	})
}

// Or combines two conditions with OR (GORM's Or method)
func (c *Condition[M]) Or(other *Condition[M]) *Condition[M] {
	return NewCondition[M](func(query *gorm.DB) *gorm.DB {
		return query.Where(c.apply).Or(other.apply)
	})
}

// Helper functions
func extractGormColumn(tag string) string {
	// Simple extraction - can be enhanced
	return ""
}

func toSnakeCase(s string) string {
	// Simple conversion - can be enhanced
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return result
}

