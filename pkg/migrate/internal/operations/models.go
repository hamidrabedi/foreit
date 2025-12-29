package operations

import (
	"fmt"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrate/schema"
)

// CreateModel creates a new model/table
type CreateModel struct {
	*BaseOperation
	ModelName string
	Fields    []generator.FieldDefinition
	Options   map[string]interface{} // Meta options like db_table, ordering, etc.
}

// NewCreateModel creates a new CreateModel operation
func NewCreateModel(modelName string, fields []generator.FieldDefinition, options map[string]interface{}) *CreateModel {
	return &CreateModel{
		BaseOperation: NewBaseOperation(true, true, CategoryAddition),
		ModelName:     modelName,
		Fields:        fields,
		Options:       options,
	}
}

// StateForwards adds the model to the state
func (op *CreateModel) StateForwards(appLabel string, s *schema.SchemaState) error {
	tableName := op.ModelName
	if dbTable, ok := op.Options["db_table"].(string); ok {
		tableName = dbTable
	}

	table := &schema.TableState{
		Columns: make(map[string]*schema.ColumnState),
	}
	for _, field := range op.Fields {
		table.Columns[field.Name] = &schema.ColumnState{
			Name:     field.Name,
			Type:     field.Type,
			Required: field.Required,
			Default:  field.Default,
		}
	}
	s.Tables[tableName] = table
	return nil
}

// DatabaseForwards creates the table in the database
func (op *CreateModel) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	tableName := op.ModelName
	if dbTable, ok := op.Options["db_table"].(string); ok {
		tableName = dbTable
	}

	model := &generator.ModelDefinition{
		Name:   op.ModelName,
		Fields: op.Fields,
		Meta: generator.MetaDefinition{
			TableName: tableName,
		},
	}
	return schemaEditor.CreateTable(model)
}

// DatabaseBackwards drops the table from the database
func (op *CreateModel) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	tableName := op.ModelName
	if dbTable, ok := op.Options["db_table"].(string); ok {
		tableName = dbTable
	}
	return schemaEditor.DeleteTable(tableName)
}

// Describe returns a human-readable description
func (op *CreateModel) Describe() string {
	return fmt.Sprintf("Create model %s", op.ModelName)
}

// DeleteModel deletes a model/table
type DeleteModel struct {
	*BaseOperation
	ModelName string
}

// NewDeleteModel creates a new DeleteModel operation
func NewDeleteModel(modelName string) *DeleteModel {
	return &DeleteModel{
		BaseOperation: NewBaseOperation(true, true, CategoryRemoval),
		ModelName:     modelName,
	}
}

// StateForwards removes the model from the state
func (op *DeleteModel) StateForwards(appLabel string, s *schema.SchemaState) error {
	delete(s.Tables, op.ModelName)
	return nil
}

// DatabaseForwards drops the table from the database
func (op *DeleteModel) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return schemaEditor.DeleteTable(op.ModelName)
}

// DatabaseBackwards recreates the table
func (op *DeleteModel) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	table, exists := fromState.Tables[op.ModelName]
	if !exists {
		return fmt.Errorf("table %s does not exist in fromState", op.ModelName)
	}
	// Convert TableState to ModelDefinition
	fields := make([]generator.FieldDefinition, 0, len(table.Columns))
	for name, column := range table.Columns {
		fields = append(fields, generator.FieldDefinition{
			Name:     name,
			Type:     column.Type,
			Required: column.Required,
			Default:  column.Default,
		})
	}
	model := &generator.ModelDefinition{
		Name:   op.ModelName,
		Fields: fields,
	}
	return schemaEditor.CreateTable(model)
}

// Describe returns a human-readable description
func (op *DeleteModel) Describe() string {
	return fmt.Sprintf("Delete model %s", op.ModelName)
}

// RenameModel renames a model/table
type RenameModel struct {
	*BaseOperation
	OldName string
	NewName string
}

// NewRenameModel creates a new RenameModel operation
func NewRenameModel(oldName, newName string) *RenameModel {
	return &RenameModel{
		BaseOperation: NewBaseOperation(true, true, CategoryAlteration),
		OldName:       oldName,
		NewName:       newName,
	}
}

// StateForwards renames the model in the state
func (op *RenameModel) StateForwards(appLabel string, s *schema.SchemaState) error {
	table, exists := s.Tables[op.OldName]
	if !exists {
		return fmt.Errorf("table %s does not exist", op.OldName)
	}
	s.Tables[op.NewName] = table
	delete(s.Tables, op.OldName)
	return nil
}

// DatabaseForwards renames the table in the database
func (op *RenameModel) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return schemaEditor.RenameTable(op.OldName, op.NewName)
}

// DatabaseBackwards reverts the table rename
func (op *RenameModel) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return schemaEditor.RenameTable(op.NewName, op.OldName)
}

// Describe returns a human-readable description
func (op *RenameModel) Describe() string {
	return fmt.Sprintf("Rename model %s to %s", op.OldName, op.NewName)
}

// AlterModelOptions alters model Meta options (doesn't affect database)
type AlterModelOptions struct {
	*BaseOperation
	ModelName string
	Options   map[string]interface{}
}

// NewAlterModelOptions creates a new AlterModelOptions operation
func NewAlterModelOptions(modelName string, options map[string]interface{}) *AlterModelOptions {
	return &AlterModelOptions{
		BaseOperation: NewBaseOperation(true, false, CategoryAlteration),
		ModelName:     modelName,
		Options:       options,
	}
}

// StateForwards updates model options in state (for historical model tracking)
func (op *AlterModelOptions) StateForwards(appLabel string, s *schema.SchemaState) error {
	// This doesn't affect the database, only the in-memory state
	// Used for tracking model options for historical models
	return nil
}

// DatabaseForwards is a no-op (options don't affect database)
func (op *AlterModelOptions) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return nil
}

// DatabaseBackwards is a no-op
func (op *AlterModelOptions) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return nil
}

// Describe returns a human-readable description
func (op *AlterModelOptions) Describe() string {
	return fmt.Sprintf("Alter model options for %s", op.ModelName)
}

// AlterModelTable alters the database table name for a model
type AlterModelTable struct {
	*BaseOperation
	ModelName string
	OldTable  string
	NewTable  string
}

// NewAlterModelTable creates a new AlterModelTable operation
func NewAlterModelTable(modelName, oldTable, newTable string) *AlterModelTable {
	return &AlterModelTable{
		BaseOperation: NewBaseOperation(true, true, CategoryAlteration),
		ModelName:     modelName,
		OldTable:      oldTable,
		NewTable:      newTable,
	}
}

// StateForwards updates the table name in state
func (op *AlterModelTable) StateForwards(appLabel string, s *schema.SchemaState) error {
	table, exists := s.Tables[op.OldTable]
	if !exists {
		return fmt.Errorf("table %s does not exist", op.OldTable)
	}
	s.Tables[op.NewTable] = table
	delete(s.Tables, op.OldTable)
	return nil
}

// DatabaseForwards renames the table in the database
func (op *AlterModelTable) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return schemaEditor.RenameTable(op.OldTable, op.NewTable)
}

// DatabaseBackwards reverts the table rename
func (op *AlterModelTable) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return schemaEditor.RenameTable(op.NewTable, op.OldTable)
}

// Describe returns a human-readable description
func (op *AlterModelTable) Describe() string {
	return fmt.Sprintf("Alter db_table for %s from %s to %s", op.ModelName, op.OldTable, op.NewTable)
}
