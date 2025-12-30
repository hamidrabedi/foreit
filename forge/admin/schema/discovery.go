package schema

import (
	"fmt"
	"reflect"

	"github.com/forgego/forge/schema"
)

// DiscoverFields discovers all fields from a schema instance and returns admin field information
func DiscoverFields[T any](schemaInstance schema.Schema) ([]FieldInfo, error) {
	fields := schemaInstance.Fields()
	adminFields := make([]FieldInfo, 0, len(fields))

	for _, schemaField := range fields {
		fieldInfo, err := mapSchemaFieldToAdminField[T](schemaField)
		if err != nil {
			return nil, fmt.Errorf("failed to map field %s: %w", schemaField.Name, err)
		}
		adminFields = append(adminFields, fieldInfo)
	}

	return adminFields, nil
}

// DiscoverRelations discovers all relations from a schema instance
func DiscoverRelations[T any](schemaInstance schema.Schema) ([]RelationInfo, error) {
	relations := schemaInstance.Relations()
	adminRelations := make([]RelationInfo, 0, len(relations))

	for _, schemaRelation := range relations {
		relationInfo := mapSchemaRelationToAdminRelation(schemaRelation)
		adminRelations = append(adminRelations, relationInfo)
	}

	return adminRelations, nil
}

// DiscoverMeta discovers model metadata from schema
func DiscoverMeta(schemaInstance schema.Schema) MetaInfo {
	meta := schemaInstance.Meta()
	return MetaInfo{
		TableName:         meta.TableName,
		VerboseName:       meta.VerboseName,
		VerboseNamePlural: meta.VerboseNamePlural,
		AppLabel:          meta.AppLabel,
		OrderBy:           meta.OrderBy,
	}
}

// FieldInfo contains admin field information derived from schema
type FieldInfo struct {
	Name          string
	Type          schema.FieldType
	SchemaField   schema.Field
	Required      bool
	ReadOnly      bool
	PrimaryKey    bool
	Unique        bool
	HasChoices    bool
	Choices       []schema.Choice
	VerboseName   string
	HelpText      string
	DefaultValue   interface{}
	MaxLength     *int
	MinLength     *int
	MaxValue      *float64
	MinValue      *float64
	DBColumn      string
	AutoIncrement bool
	AutoNow       bool
	AutoNowAdd    bool
}

// RelationInfo contains admin relation information derived from schema
type RelationInfo struct {
	Name        string
	Type        schema.RelationType
	TargetModel string
	RelatedName string
	OnDelete    schema.CascadeType
	OnUpdate    schema.CascadeType
	Optional    bool
}

// MetaInfo contains model metadata
type MetaInfo struct {
	TableName         string
	VerboseName       string
	VerboseNamePlural string
	AppLabel          string
	OrderBy           []string
}

// mapSchemaFieldToAdminField maps a schema field to admin field info
func mapSchemaFieldToAdminField[T any](schemaField schema.Field) (FieldInfo, error) {
	fieldInfo := FieldInfo{
		Name:          schemaField.Name,
		Type:          schemaField.Type,
		SchemaField:   schemaField,
		Required:      schemaField.Required,
		ReadOnly:      !schemaField.Editable, // Read-only when not editable
		PrimaryKey:    schemaField.PrimaryKey,
		Unique:        schemaField.Unique,
		HasChoices:    len(schemaField.Choices) > 0,
		Choices:       schemaField.Choices,
		VerboseName:   schemaField.VerboseName,
		HelpText:      schemaField.HelpText,
		DefaultValue:  schemaField.Default,
		MaxLength:     schemaField.MaxLength,
		MinLength:     schemaField.MinLength,
		MaxValue:      schemaField.MaxValue,
		MinValue:      schemaField.MinValue,
		DBColumn:      schemaField.DBColumn,
		AutoIncrement: schemaField.AutoIncrement,
		AutoNow:       schemaField.AutoNow,
		AutoNowAdd:    schemaField.AutoNowAdd,
	}

	if fieldInfo.DBColumn == "" {
		fieldInfo.DBColumn = fieldInfo.Name
	}

	if fieldInfo.VerboseName == "" {
		fieldInfo.VerboseName = fieldInfo.Name
	}

	return fieldInfo, nil
}

// mapSchemaRelationToAdminRelation maps a schema relation to admin relation info
func mapSchemaRelationToAdminRelation(schemaRelation schema.Relation) RelationInfo {
	return RelationInfo{
		Name:        schemaRelation.Name,
		Type:        schemaRelation.Type,
		TargetModel: schemaRelation.To,
		RelatedName: schemaRelation.RelatedName,
		OnDelete:    schemaRelation.OnDelete,
		OnUpdate:    schemaRelation.OnUpdate,
		Optional:    schemaRelation.LimitChoicesTo == nil, // Simplified check
	}
}

// GetModelType gets the Go type for a model from schema
func GetModelType(schemaInstance schema.Schema) (reflect.Type, error) {
	instanceValue := reflect.ValueOf(schemaInstance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}
	return instanceValue.Type(), nil
}
