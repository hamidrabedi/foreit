package models

import (
	"context"
)

// FieldDefinition represents a field in a model definition.
// This is a type-erased interface to avoid import cycles with the field package.
type FieldDefinition interface {
	GetName() string
	GetColumn() string
}

// FieldDefinitionWithValidators is an optional interface for fields that have validators.
// This allows API serializers to access validators without breaking the type-erased design.
type FieldDefinitionWithValidators interface {
	FieldDefinition
	GetValidators() []interface{} // Returns validators as []interface{} to avoid import cycles
}

// ModelDefinition represents a Django-style model definition with typed fields.
// Use RegisterSchema or RegisterSchemas to create ModelDefinitions from schemas.
type ModelDefinition[T any] struct {
	tableName     string
	fields        []FieldDefinition
	relationships []RelationDefinition
	validators    []ModelValidator
	meta          ModelDefinitionMeta
	hooks         *ModelHooks[T]
}

// ModelDefinitionMeta contains model-level configuration
type ModelDefinitionMeta struct {
	TableName     string
	Timestamps    bool
	SoftDelete    bool
	Ordering      []string
	Indexes       []Index
	UniqueIndexes []Index
}

// Index represents a database index
type Index struct {
	Name   string
	Fields []string
	Unique bool
}

// RelationDefinition represents a model relationship
type RelationDefinition struct {
	Name         string
	Type         RelationType
	RelatedModel string
	ForeignKey   string
	BackRef      string
	Through      string
	OnDelete     string
	OnUpdate     string
}

// ModelHooks allows users to override model behavior
type ModelHooks[T any] struct {
	BeforeCreate func(ctx context.Context, instance *T) error
	AfterCreate  func(ctx context.Context, instance *T) error
	BeforeUpdate func(ctx context.Context, instance *T) error
	AfterUpdate  func(ctx context.Context, instance *T) error
	BeforeDelete func(ctx context.Context, instance *T) error
	AfterDelete  func(ctx context.Context, instance *T) error
	BeforeSave   func(ctx context.Context, instance *T) error
	AfterSave    func(ctx context.Context, instance *T) error
	Clean        func(instance *T) error
	Save         func(ctx context.Context, instance *T) error
	Delete       func(ctx context.Context, instance *T) error
}


// GetTableName returns the table name
func (m *ModelDefinition[T]) GetTableName() string {
	return m.tableName
}

// GetFields returns all fields
func (m *ModelDefinition[T]) GetFields() []FieldDefinition {
	return m.fields
}

// GetMeta returns model metadata
func (m *ModelDefinition[T]) GetMeta() ModelDefinitionMeta {
	return m.meta
}

// GetHooks returns model hooks
func (m *ModelDefinition[T]) GetHooks() *ModelHooks[T] {
	return m.hooks
}

// GetRelationships returns all relationships
func (m *ModelDefinition[T]) GetRelationships() []RelationDefinition {
	return m.relationships
}

// GetValidators returns all validators
func (m *ModelDefinition[T]) GetValidators() []ModelValidator {
	return m.validators
}

