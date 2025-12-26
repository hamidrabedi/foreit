package operations

import (
	"fmt"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/state"
)

// AddField adds a field to a model
type AddField struct {
	*BaseOperation
	ModelName       string
	FieldName       string
	Field           generator.FieldDefinition
	PreserveDefault bool // Whether the default value is permanent
}

// NewAddField creates a new AddField operation
func NewAddField(modelName, fieldName string, field generator.FieldDefinition, preserveDefault bool) *AddField {
	return &AddField{
		BaseOperation:   NewBaseOperation(true, true, CategoryAddition),
		ModelName:       modelName,
		FieldName:       fieldName,
		Field:           field,
		PreserveDefault: preserveDefault,
	}
}

// StateForwards adds the field to the model state
func (op *AddField) StateForwards(appLabel string, s *state.SchemaState) error {
	table, exists := s.Tables[op.ModelName]
	if !exists {
		return fmt.Errorf("table %s does not exist", op.ModelName)
	}
	table.Columns[op.FieldName] = &state.ColumnState{
		Name:     op.FieldName,
		Type:     op.Field.Type,
		Required: op.Field.Required,
		Default:  op.Field.Default,
	}
	return nil
}

// DatabaseForwards adds the column to the database
func (op *AddField) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.AddColumn(op.ModelName, op.Field)
}

// DatabaseBackwards removes the column from the database
func (op *AddField) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.RemoveColumn(op.ModelName, op.FieldName)
}

// Describe returns a human-readable description
func (op *AddField) Describe() string {
	return fmt.Sprintf("Add field %s to %s", op.FieldName, op.ModelName)
}

// RemoveField removes a field from a model
type RemoveField struct {
	*BaseOperation
	ModelName string
	FieldName string
}

// NewRemoveField creates a new RemoveField operation
func NewRemoveField(modelName, fieldName string) *RemoveField {
	return &RemoveField{
		BaseOperation: NewBaseOperation(true, true, CategoryRemoval),
		ModelName:     modelName,
		FieldName:     fieldName,
	}
}

// StateForwards removes the field from the model state
func (op *RemoveField) StateForwards(appLabel string, s *state.SchemaState) error {
	table, exists := s.Tables[op.ModelName]
	if !exists {
		return fmt.Errorf("table %s does not exist", op.ModelName)
	}
	delete(table.Columns, op.FieldName)
	return nil
}

// DatabaseForwards removes the column from the database
func (op *RemoveField) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.RemoveColumn(op.ModelName, op.FieldName)
}

// DatabaseBackwards adds the column back to the database
func (op *RemoveField) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	// We need the old field definition from fromState
	table, exists := fromState.Tables[op.ModelName]
	if !exists {
		return fmt.Errorf("table %s does not exist in fromState", op.ModelName)
	}
	oldColumn, exists := table.Columns[op.FieldName]
	if !exists {
		return fmt.Errorf("column %s does not exist in fromState", op.FieldName)
	}
	// Convert ColumnState to FieldDefinition
	field := generator.FieldDefinition{
		Name:     op.FieldName,
		Type:     oldColumn.Type,
		Required: oldColumn.Required,
		Default:  oldColumn.Default,
	}
	return schemaEditor.AddColumn(op.ModelName, field)
}

// Describe returns a human-readable description
func (op *RemoveField) Describe() string {
	return fmt.Sprintf("Remove field %s from %s", op.FieldName, op.ModelName)
}

// AlterField alters a field's properties
type AlterField struct {
	*BaseOperation
	ModelName string
	FieldName string
	OldField  generator.FieldDefinition
	NewField  generator.FieldDefinition
}

// NewAlterField creates a new AlterField operation
func NewAlterField(modelName, fieldName string, oldField, newField generator.FieldDefinition) *AlterField {
	return &AlterField{
		BaseOperation: NewBaseOperation(true, true, CategoryAlteration),
		ModelName:     modelName,
		FieldName:     fieldName,
		OldField:      oldField,
		NewField:      newField,
	}
}

// StateForwards updates the field in the model state
func (op *AlterField) StateForwards(appLabel string, s *state.SchemaState) error {
	table, exists := s.Tables[op.ModelName]
	if !exists {
		return fmt.Errorf("table %s does not exist", op.ModelName)
	}
	column, exists := table.Columns[op.FieldName]
	if !exists {
		return fmt.Errorf("column %s does not exist", op.FieldName)
	}
	column.Type = op.NewField.Type
	column.Required = op.NewField.Required
	column.Default = op.NewField.Default
	return nil
}

// DatabaseForwards alters the column in the database
func (op *AlterField) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.AlterColumn(op.ModelName, op.OldField, op.NewField)
}

// DatabaseBackwards reverts the column alteration
func (op *AlterField) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.AlterColumn(op.ModelName, op.NewField, op.OldField)
}

// Describe returns a human-readable description
func (op *AlterField) Describe() string {
	return fmt.Sprintf("Alter field %s on %s", op.FieldName, op.ModelName)
}

// RenameField renames a field
type RenameField struct {
	*BaseOperation
	ModelName string
	OldName   string
	NewName   string
}

// NewRenameField creates a new RenameField operation
func NewRenameField(modelName, oldName, newName string) *RenameField {
	return &RenameField{
		BaseOperation: NewBaseOperation(true, true, CategoryAlteration),
		ModelName:     modelName,
		OldName:       oldName,
		NewName:       newName,
	}
}

// StateForwards renames the field in the model state
func (op *RenameField) StateForwards(appLabel string, s *state.SchemaState) error {
	table, exists := s.Tables[op.ModelName]
	if !exists {
		return fmt.Errorf("table %s does not exist", op.ModelName)
	}
	column, exists := table.Columns[op.OldName]
	if !exists {
		return fmt.Errorf("column %s does not exist", op.OldName)
	}
	table.Columns[op.NewName] = column
	delete(table.Columns, op.OldName)
	return nil
}

// DatabaseForwards renames the column in the database
func (op *RenameField) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.RenameColumn(op.ModelName, op.OldName, op.NewName)
}

// DatabaseBackwards reverts the column rename
func (op *RenameField) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *state.SchemaState) error {
	return schemaEditor.RenameColumn(op.ModelName, op.NewName, op.OldName)
}

// Describe returns a human-readable description
func (op *RenameField) Describe() string {
	return fmt.Sprintf("Rename field %s on %s to %s", op.OldName, op.ModelName, op.NewName)
}

