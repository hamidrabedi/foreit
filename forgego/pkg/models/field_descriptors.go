package models

import (
	"reflect"
)

// FieldDescriptor is the interface for field definitions.
type FieldDescriptor interface {
	GetName() string
	GetType() reflect.Type
	GetColumn() string
	GetOptions() FieldOptions
}

// FieldOptions contains field configuration.
type FieldOptions struct {
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Required      bool
	Default       interface{}
	MaxLength     *int
	MinLength     *int
	Min           *float64
	Max           *float64
	Choices       []Choice
}

// Choice represents a choice option for ChoiceField.
type Choice struct {
	Value string
	Label string
}

// Int64FieldDescriptor represents an int64 field.
type Int64FieldDescriptor struct {
	name    string
	column  string
	options FieldOptions
}

// StringFieldDescriptor represents a string field.
type StringFieldDescriptor struct {
	name    string
	column  string
	options FieldOptions
}

// BoolFieldDescriptor represents a bool field.
type BoolFieldDescriptor struct {
	name    string
	column  string
	options FieldOptions
}

// TimeFieldDescriptor represents a time.Time field.
type TimeFieldDescriptor struct {
	name    string
	column  string
	options FieldOptions
}

// Float64FieldDescriptor represents a float64 field.
type Float64FieldDescriptor struct {
	name    string
	column  string
	options FieldOptions
}

// Int64 creates a new int64 field descriptor.
func Int64(name string) *Int64FieldDescriptor {
	return &Int64FieldDescriptor{
		name:   name,
		column: toSnakeCase(name),
		options: FieldOptions{},
	}
}

// String creates a new string field descriptor.
func String(name string) *StringFieldDescriptor {
	return &StringFieldDescriptor{
		name:   name,
		column: toSnakeCase(name),
		options: FieldOptions{},
	}
}

// Bool creates a new bool field descriptor.
func Bool(name string) *BoolFieldDescriptor {
	return &BoolFieldDescriptor{
		name:   name,
		column: toSnakeCase(name),
		options: FieldOptions{},
	}
}

// Time creates a new time.Time field descriptor.
func Time(name string) *TimeFieldDescriptor {
	return &TimeFieldDescriptor{
		name:   name,
		column: toSnakeCase(name),
		options: FieldOptions{},
	}
}

// Float64 creates a new float64 field descriptor.
func Float64(name string) *Float64FieldDescriptor {
	return &Float64FieldDescriptor{
		name:   name,
		column: toSnakeCase(name),
		options: FieldOptions{},
	}
}

// GetName returns the field name.
func (f *Int64FieldDescriptor) GetName() string {
	return f.name
}

// GetType returns the field type.
func (f *Int64FieldDescriptor) GetType() reflect.Type {
	return reflect.TypeOf(int64(0))
}

// GetColumn returns the database column name.
func (f *Int64FieldDescriptor) GetColumn() string {
	return f.column
}

// GetOptions returns the field options.
func (f *Int64FieldDescriptor) GetOptions() FieldOptions {
	return f.options
}

// PrimaryKey marks the field as a primary key.
func (f *Int64FieldDescriptor) PrimaryKey() *Int64FieldDescriptor {
	f.options.PrimaryKey = true
	return f
}

// AutoIncrement marks the field as auto-incrementing.
func (f *Int64FieldDescriptor) AutoIncrement() *Int64FieldDescriptor {
	f.options.AutoIncrement = true
	return f
}

// Default sets the default value.
func (f *Int64FieldDescriptor) Default(value int64) *Int64FieldDescriptor {
	f.options.Default = value
	return f
}

// GetName returns the field name.
func (f *StringFieldDescriptor) GetName() string {
	return f.name
}

// GetType returns the field type.
func (f *StringFieldDescriptor) GetType() reflect.Type {
	return reflect.TypeOf("")
}

// GetColumn returns the database column name.
func (f *StringFieldDescriptor) GetColumn() string {
	return f.column
}

// GetOptions returns the field options.
func (f *StringFieldDescriptor) GetOptions() FieldOptions {
	return f.options
}

// PrimaryKey marks the field as a primary key.
func (f *StringFieldDescriptor) PrimaryKey() *StringFieldDescriptor {
	f.options.PrimaryKey = true
	return f
}

// Unique marks the field as unique.
func (f *StringFieldDescriptor) Unique() *StringFieldDescriptor {
	f.options.Unique = true
	return f
}

// Required marks the field as required.
func (f *StringFieldDescriptor) Required() *StringFieldDescriptor {
	f.options.Required = true
	return f
}

// Default sets the default value.
func (f *StringFieldDescriptor) Default(value string) *StringFieldDescriptor {
	f.options.Default = value
	return f
}

// MaxLength sets the maximum length.
func (f *StringFieldDescriptor) MaxLength(n int) *StringFieldDescriptor {
	f.options.MaxLength = &n
	return f
}

// MinLength sets the minimum length.
func (f *StringFieldDescriptor) MinLength(n int) *StringFieldDescriptor {
	f.options.MinLength = &n
	return f
}

// GetName returns the field name.
func (f *BoolFieldDescriptor) GetName() string {
	return f.name
}

// GetType returns the field type.
func (f *BoolFieldDescriptor) GetType() reflect.Type {
	return reflect.TypeOf(false)
}

// GetColumn returns the database column name.
func (f *BoolFieldDescriptor) GetColumn() string {
	return f.column
}

// GetOptions returns the field options.
func (f *BoolFieldDescriptor) GetOptions() FieldOptions {
	return f.options
}

// Default sets the default value.
func (f *BoolFieldDescriptor) Default(value bool) *BoolFieldDescriptor {
	f.options.Default = value
	return f
}

// GetName returns the field name.
func (f *TimeFieldDescriptor) GetName() string {
	return f.name
}

// GetType returns the field type.
func (f *TimeFieldDescriptor) GetType() reflect.Type {
	return reflect.TypeOf((*interface{})(nil)).Elem() // Will be time.Time in actual usage
}

// GetColumn returns the database column name.
func (f *TimeFieldDescriptor) GetColumn() string {
	return f.column
}

// GetOptions returns the field options.
func (f *TimeFieldDescriptor) GetOptions() FieldOptions {
	return f.options
}

// Default sets the default value (use time.Now or similar).
func (f *TimeFieldDescriptor) Default(value interface{}) *TimeFieldDescriptor {
	f.options.Default = value
	return f
}

// GetName returns the field name.
func (f *Float64FieldDescriptor) GetName() string {
	return f.name
}

// GetType returns the field type.
func (f *Float64FieldDescriptor) GetType() reflect.Type {
	return reflect.TypeOf(float64(0))
}

// GetColumn returns the database column name.
func (f *Float64FieldDescriptor) GetColumn() string {
	return f.column
}

// GetOptions returns the field options.
func (f *Float64FieldDescriptor) GetOptions() FieldOptions {
	return f.options
}

// Default sets the default value.
func (f *Float64FieldDescriptor) Default(value float64) *Float64FieldDescriptor {
	f.options.Default = value
	return f
}

// Min sets the minimum value.
func (f *Float64FieldDescriptor) Min(value float64) *Float64FieldDescriptor {
	f.options.Min = &value
	return f
}

// Max sets the maximum value.
func (f *Float64FieldDescriptor) Max(value float64) *Float64FieldDescriptor {
	f.options.Max = &value
	return f
}

