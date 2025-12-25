package models

// Schema is embedded in user models to enable builder methods.
type Schema struct{}

// SchemaInterface defines methods models should implement.
type SchemaInterface interface {
	Fields() []FieldDescriptor
	Relations() []RelationDescriptor
}

// AdvancedSchemaInterface defines optional advanced methods.
type AdvancedSchemaInterface interface {
	SchemaInterface
	TableName() string
	Indexes() []IndexDescriptor
	Hooks() *ModelHooks[interface{}] // Use interface{} as a placeholder
}

// IndexDescriptor represents a database index.
type IndexDescriptor struct {
	Name    string
	Fields  []string
	Unique  bool
	Partial string // Optional WHERE clause
}

