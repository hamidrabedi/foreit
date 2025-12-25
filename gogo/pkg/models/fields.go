package models

import (
	"time"
)

// Field represents a type-safe field definition
type Field[T any] struct {
	name     string
	required bool
	unique   bool
	indexed  bool
	default_ *T
	validate func(T) error
}

// NewField creates a new type-safe field
func NewField[T any](name string) *Field[T] {
	return &Field[T]{
		name: name,
	}
}

// Required marks the field as required
func (f *Field[T]) Required() *Field[T] {
	f.required = true
	return f
}

// Unique marks the field as unique
func (f *Field[T]) Unique() *Field[T] {
	f.unique = true
	return f
}

// Indexed marks the field as indexed
func (f *Field[T]) Indexed() *Field[T] {
	f.indexed = true
	return f
}

// Default sets a default value
func (f *Field[T]) Default(value T) *Field[T] {
	f.default_ = &value
	return f
}

// Validate adds a validation function
func (f *Field[T]) Validate(validator func(T) error) *Field[T] {
	f.validate = validator
	return f
}

// Name returns the field name
func (f *Field[T]) Name() string {
	return f.name
}

// StringField is a string field
type StringField = Field[string]

// IntField is an int field
type IntField = Field[int]

// FloatField is a float64 field
type FloatField = Field[float64]

// BoolField is a bool field
type BoolField = Field[bool]

// TimeField is a time.Time field
type TimeField = Field[time.Time]

// FieldRef provides type-safe field references for queries
type FieldRef[T any] struct {
	name string
}

// NewFieldRef creates a type-safe field reference
func NewFieldRef[T any](name string) *FieldRef[T] {
	return &FieldRef[T]{name: name}
}

// Name returns the field name
func (f *FieldRef[T]) Name() string {
	return f.name
}

// Eq creates an equality condition
func (f *FieldRef[T]) Eq(value T) Q {
	return Q{f.name: value}
}

// Ne creates a not-equal condition
func (f *FieldRef[T]) Ne(value T) Q {
	return Q{f.name + "__ne": value}
}

// Gt creates a greater-than condition
func (f *FieldRef[T]) Gt(value T) Q {
	return Q{f.name + "__gt": value}
}

// Gte creates a greater-than-or-equal condition
func (f *FieldRef[T]) Gte(value T) Q {
	return Q{f.name + "__gte": value}
}

// Lt creates a less-than condition
func (f *FieldRef[T]) Lt(value T) Q {
	return Q{f.name + "__lt": value}
}

// Lte creates a less-than-or-equal condition
func (f *FieldRef[T]) Lte(value T) Q {
	return Q{f.name + "__lte": value}
}

// In creates an IN condition
func (f *FieldRef[T]) In(values []T) Q {
	return Q{f.name + "__in": values}
}

// Contains creates a contains condition (for strings)
func (f *FieldRef[string]) Contains(value string) Q {
	return Q{f.name + "__contains": value}
}

// IContains creates a case-insensitive contains condition
func (f *FieldRef[string]) IContains(value string) Q {
	return Q{f.name + "__icontains": value}
}

// StartsWith creates a starts-with condition
func (f *FieldRef[string]) StartsWith(value string) Q {
	return Q{f.name + "__startswith": value}
}

// EndsWith creates an ends-with condition
func (f *FieldRef[string]) EndsWith(value string) Q {
	return Q{f.name + "__endswith": value}
}

// IsNull creates an IS NULL condition
func (f *FieldRef[T]) IsNull() Q {
	return Q{f.name + "__isnull": true}
}

// IsNotNull creates an IS NOT NULL condition
func (f *FieldRef[T]) IsNotNull() Q {
	return Q{f.name + "__isnull": false}
}

