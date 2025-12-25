package models

import (
	"time"
)

// ModelDefinition defines a type-safe model schema
type ModelDefinition[T Model] struct {
	name   string
	fields map[string]FieldDef
	pk     string
}

// FieldDef is a type-erased field definition
type FieldDef interface {
	Name() string
	Type() string
	Required() bool
	Unique() bool
	Indexed() bool
}

// NewModelDefinition creates a new type-safe model definition
func NewModelDefinition[T Model](name string) *ModelDefinition[T] {
	return &ModelDefinition[T]{
		name:   name,
		fields: make(map[string]FieldDef),
		pk:     "id",
	}
}

// AddField adds a type-safe field to the definition
func (d *ModelDefinition[T]) AddField(field FieldDef) *ModelDefinition[T] {
	d.fields[field.Name()] = field
	return d
}

// SetPrimaryKey sets the primary key field
func (d *ModelDefinition[T]) SetPrimaryKey(fieldName string) *ModelDefinition[T] {
	d.pk = fieldName
	return d
}

// Name returns the model name
func (d *ModelDefinition[T]) Name() string {
	return d.name
}

// Fields returns all field definitions
func (d *ModelDefinition[T]) Fields() map[string]FieldDef {
	return d.fields
}

// PrimaryKey returns the primary key field name
func (d *ModelDefinition[T]) PrimaryKey() string {
	return d.pk
}

// ModelBuilder provides a fluent API for building model definitions
type ModelBuilder[T Model] struct {
	def *ModelDefinition[T]
}

// NewModelBuilder creates a new model builder
func NewModelBuilder[T Model](name string) *ModelBuilder[T] {
	return &ModelBuilder[T]{
		def: NewModelDefinition[T](name),
	}
}

// String adds a string field
func (b *ModelBuilder[T]) String(name string) *Field[string] {
	field := NewField[string](name)
	b.def.AddField(field)
	return field
}

// Int adds an int field
func (b *ModelBuilder[T]) Int(name string) *Field[int] {
	field := NewField[int](name)
	b.def.AddField(field)
	return field
}

// Float adds a float64 field
func (b *ModelBuilder[T]) Float(name string) *Field[float64] {
	field := NewField[float64](name)
	b.def.AddField(field)
	return field
}

// Bool adds a bool field
func (b *ModelBuilder[T]) Bool(name string) *Field[bool] {
	field := NewField[bool](name)
	b.def.AddField(field)
	return field
}

// Time adds a time.Time field
func (b *ModelBuilder[T]) Time(name string) *Field[time.Time] {
	field := NewField[time.Time](name)
	b.def.AddField(field)
	return field
}

// Build returns the model definition
func (b *ModelBuilder[T]) Build() *ModelDefinition[T] {
	return b.def
}

// Example usage:
// type User struct {
//     ID        int
//     Email     string
//     Name      string
//     CreatedAt time.Time
// }
//
// var UserModel = NewModelBuilder[User]("User").
//     Int("id").Required().Indexed().
//     String("email").Required().Unique().
//     String("name").Required().
//     Time("created_at").Default(time.Now()).
//     Build()

