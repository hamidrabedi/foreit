package operations

import (
	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrate/schema"
)

// OperationCategory represents the category of an operation
type OperationCategory string

const (
	CategoryAddition   OperationCategory = "+"
	CategoryRemoval    OperationCategory = "-"
	CategoryAlteration OperationCategory = "~"
	CategoryPython     OperationCategory = "p"
	CategorySQL        OperationCategory = "s"
	CategoryMixed      OperationCategory = "?"
)

// Operation is the base interface for all migration operations
// Similar to Django's Operation class with state_forwards and database_forwards
type Operation interface {
	// StateForwards mutates the in-memory state to represent what this operation performs
	// This is used to track model history without touching the database
	StateForwards(appLabel string, s *schema.SchemaState) error

	// DatabaseForwards performs the mutation on the database schema in the forwards direction
	DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error

	// DatabaseBackwards performs the mutation in the reverse direction
	DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error

	// Reversible returns whether this operation can be reversed
	Reversible() bool

	// ReducesToSQL returns whether this operation can be represented as SQL
	ReducesToSQL() bool

	// Describe returns a human-readable description of the operation
	Describe() string

	// Category returns the operation category
	Category() OperationCategory
}

// SchemaEditor is an interface for performing database schema operations
type SchemaEditor interface {
	// Table operations
	CreateTable(table *generator.ModelDefinition) error
	DeleteTable(tableName string) error
	RenameTable(oldName, newName string) error
	AlterTable(tableName string, changes []interface{}) error

	// Column operations
	AddColumn(tableName string, column generator.FieldDefinition) error
	RemoveColumn(tableName, columnName string) error
	AlterColumn(tableName string, oldColumn, newColumn generator.FieldDefinition) error
	RenameColumn(tableName, oldName, newName string) error

	// Index operations
	AddIndex(tableName string, index *generator.IndexDefinition) error
	RemoveIndex(tableName, indexName string) error
	RenameIndex(tableName, oldName, newName string) error

	// Constraint operations
	AddConstraint(tableName string, constraint *generator.ConstraintDefinition) error
	RemoveConstraint(tableName, constraintName string) error

	// Execute raw SQL
	ExecuteSQL(sql string) error
}

// BaseOperation provides default implementations for common operation methods
type BaseOperation struct {
	reversible     bool
	reducesToSQL   bool
	category       OperationCategory
}

// NewBaseOperation creates a new base operation
func NewBaseOperation(reversible, reducesToSQL bool, category OperationCategory) *BaseOperation {
	return &BaseOperation{
		reversible:   reversible,
		reducesToSQL: reducesToSQL,
		category:     category,
	}
}

// Reversible returns whether this operation can be reversed
func (b *BaseOperation) Reversible() bool {
	return b.reversible
}

// ReducesToSQL returns whether this operation can be represented as SQL
func (b *BaseOperation) ReducesToSQL() bool {
	return b.reducesToSQL
}

// Category returns the operation category
func (b *BaseOperation) Category() OperationCategory {
	return b.category
}
