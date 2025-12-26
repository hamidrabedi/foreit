package schema

import (
	"github.com/google/uuid"
)

// UUIDFieldBuilder builds UUID fields
type UUIDFieldBuilder struct {
	FieldBuilder
}

// UUID creates a new UUID field builder
func UUID(name string) *UUIDFieldBuilder {
	return &UUIDFieldBuilder{
		FieldBuilder: FieldBuilder{
			field: Field{
				Name: name,
				Type: TypeUUID,
			},
		},
	}
}

// Required marks the field as required
func (b *UUIDFieldBuilder) Required() *UUIDFieldBuilder {
	b.field.Required = true
	return b
}

// Unique marks the field as unique
func (b *UUIDFieldBuilder) Unique() *UUIDFieldBuilder {
	b.field.Unique = true
	return b
}

// Primary marks the field as primary key
func (b *UUIDFieldBuilder) Primary() *UUIDFieldBuilder {
	b.field.PrimaryKey = true
	return b
}

// DefaultUUID sets a default UUID value
func (b *UUIDFieldBuilder) DefaultUUID(value uuid.UUID) *UUIDFieldBuilder {
	b.field.Default = value.String()
	return b
}

// DefaultNewUUID sets default to generate new UUID
func (b *UUIDFieldBuilder) DefaultNewUUID() *UUIDFieldBuilder {
	b.field.Default = uuid.New().String()
	return b
}

// HelpText sets the help text
func (b *UUIDFieldBuilder) HelpText(text string) *UUIDFieldBuilder {
	b.field.HelpText = text
	return b
}

// VerboseName sets the verbose name
func (b *UUIDFieldBuilder) VerboseName(name string) *UUIDFieldBuilder {
	b.field.VerboseName = name
	return b
}

// Build returns the built field
func (b *UUIDFieldBuilder) Build() Field {
	return b.field
}
