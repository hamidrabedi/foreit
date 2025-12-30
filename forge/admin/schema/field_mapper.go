package schema

import (
	"github.com/forgego/forge/schema"
)

// FieldMapper maps schema fields to admin field types
type FieldMapper struct{}

// NewFieldMapper creates a new field mapper
func NewFieldMapper() *FieldMapper {
	return &FieldMapper{}
}

// MapFieldType maps a schema field type to an admin field type identifier
func (fm *FieldMapper) MapFieldType(fieldType schema.FieldType) string {
	switch fieldType {
	case schema.TypeString:
		return "string"
	case schema.TypeInt64:
		return "int64"
	case schema.TypeInt32:
		return "int32"
	case schema.TypeBool:
		return "bool"
	case schema.TypeTime, schema.TypeDate, schema.TypeDateTime:
		return "datetime"
	case schema.TypeFloat32, schema.TypeFloat64:
		return "float"
	case schema.TypeDecimal:
		return "decimal"
	case schema.TypeText:
		return "text"
	case schema.TypeEmail:
		return "email"
	case schema.TypeURL:
		return "url"
	case schema.TypeUUID:
		return "uuid"
	case schema.TypeJSON:
		return "json"
	case schema.TypeBytes:
		return "bytes"
	case schema.TypeForeignKey:
		return "foreignkey"
	case schema.TypeManyToMany:
		return "manytomany"
	case schema.TypeOneToOne:
		return "onetoone"
	default:
		return "unknown"
	}
}

// ShouldDisplayInList determines if a field should be displayed in list view by default
func (fm *FieldMapper) ShouldDisplayInList(field schema.Field) bool {
	// Don't show primary keys, passwords, or very long text fields by default
	if field.PrimaryKey {
		return false
	}
	if field.Type == schema.TypeText && field.MaxLength != nil && *field.MaxLength > 500 {
		return false
	}
	if !field.Serialize {
		// Write-only fields (not serialized) shouldn't be in list by default
		return false
	}
	return true
}

// ShouldDisplayInForm determines if a field should be displayed in form by default
func (fm *FieldMapper) ShouldDisplayInForm(field schema.Field) bool {
	// Don't show auto-generated fields in forms
	if field.AutoIncrement {
		return false
	}
	if field.AutoNow || field.AutoNowAdd {
		return false
	}
	return true
}

// GetDefaultWidgetType returns the default widget type for a field
func (fm *FieldMapper) GetDefaultWidgetType(field schema.Field) string {
	switch field.Type {
	case schema.TypeBool:
		return "checkbox"
	case schema.TypeText:
		return "textarea"
	case schema.TypeDate:
		return "date"
	case schema.TypeDateTime, schema.TypeTime:
		return "datetime"
	case schema.TypeEmail:
		return "email"
	case schema.TypeURL:
		return "url"
	case schema.TypeForeignKey, schema.TypeOneToOne:
		return "select"
	case schema.TypeManyToMany:
		return "multiselect"
	case schema.TypeJSON:
		return "json"
	default:
		return "text"
	}
}
