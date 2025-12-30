package fields

import (
	"github.com/forgego/forge/schema"
)

// FieldRegistry manages field type mappings
type FieldRegistry struct {
	typeMappings map[schema.FieldType]FieldTypeInfo
}

// FieldTypeInfo contains information about a field type
type FieldTypeInfo struct {
	TypeName      string
	GoType        string
	DefaultWidget string
}

// NewFieldRegistry creates a new field registry
func NewFieldRegistry() *FieldRegistry {
	registry := &FieldRegistry{
		typeMappings: make(map[schema.FieldType]FieldTypeInfo),
	}

	// Register default field types
	registry.registerDefaults()

	return registry
}

// registerDefaults registers default field type mappings
func (fr *FieldRegistry) registerDefaults() {
	fr.typeMappings[schema.TypeString] = FieldTypeInfo{
		TypeName:      "string",
		GoType:        "string",
		DefaultWidget: "text",
	}
	fr.typeMappings[schema.TypeText] = FieldTypeInfo{
		TypeName:      "text",
		GoType:        "string",
		DefaultWidget: "textarea",
	}
	fr.typeMappings[schema.TypeInt64] = FieldTypeInfo{
		TypeName:      "int64",
		GoType:        "int64",
		DefaultWidget: "number",
	}
	fr.typeMappings[schema.TypeInt32] = FieldTypeInfo{
		TypeName:      "int32",
		GoType:        "int32",
		DefaultWidget: "number",
	}
	fr.typeMappings[schema.TypeBool] = FieldTypeInfo{
		TypeName:      "bool",
		GoType:        "bool",
		DefaultWidget: "checkbox",
	}
	fr.typeMappings[schema.TypeDateTime] = FieldTypeInfo{
		TypeName:      "datetime",
		GoType:        "time.Time",
		DefaultWidget: "datetime",
	}
	fr.typeMappings[schema.TypeDate] = FieldTypeInfo{
		TypeName:      "date",
		GoType:        "time.Time",
		DefaultWidget: "date",
	}
	fr.typeMappings[schema.TypeEmail] = FieldTypeInfo{
		TypeName:      "email",
		GoType:        "string",
		DefaultWidget: "email",
	}
	fr.typeMappings[schema.TypeURL] = FieldTypeInfo{
		TypeName:      "url",
		GoType:        "string",
		DefaultWidget: "url",
	}
	fr.typeMappings[schema.TypeForeignKey] = FieldTypeInfo{
		TypeName:      "foreignkey",
		GoType:        "int64",
		DefaultWidget: "select",
	}
	fr.typeMappings[schema.TypeOneToOne] = FieldTypeInfo{
		TypeName:      "onetoone",
		GoType:        "int64",
		DefaultWidget: "select",
	}
	fr.typeMappings[schema.TypeManyToMany] = FieldTypeInfo{
		TypeName:      "manytomany",
		GoType:        "[]int64",
		DefaultWidget: "multiselect",
	}
}

// GetFieldTypeInfo gets field type information
func (fr *FieldRegistry) GetFieldTypeInfo(fieldType schema.FieldType) (FieldTypeInfo, bool) {
	info, ok := fr.typeMappings[fieldType]
	return info, ok
}

// RegisterFieldType registers a custom field type
func (fr *FieldRegistry) RegisterFieldType(fieldType schema.FieldType, info FieldTypeInfo) {
	fr.typeMappings[fieldType] = info
}
