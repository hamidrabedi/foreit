package schema

import (
	"fmt"
	"strings"

	"github.com/forgego/forge/pkg/migrate/types"
)

// SchemaState represents the current database schema state
type SchemaState struct {
	Tables map[string]*TableState
}

// TableState represents the state of a single table
type TableState struct {
	Name        string
	Columns     map[string]*ColumnState
	Indexes     map[string]*IndexState
	ForeignKeys map[string]*ForeignKeyState
	Constraints map[string]*ConstraintState
}

// ColumnState represents the state of a column
type ColumnState struct {
	Name          string
	Type          string
	GoType        string
	Required      bool
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Default       interface{}
	Options       map[string]interface{}
}

// IndexState represents the state of an index
type IndexState struct {
	Name   string
	Fields []string
	Unique bool
}

// ForeignKeyState represents the state of a foreign key
type ForeignKeyState struct {
	Name         string
	Column       string
	TargetTable  string
	TargetColumn string
	OnDelete     string
	OnUpdate     string
}

// ConstraintState represents the state of a constraint
type ConstraintState struct {
	Name      string
	Type      string
	Condition string
	Fields    []string
}

// StateManager manages schema state
type StateManager interface {
	// Load loads the current schema state
	Load() (*SchemaState, error)
	
	// Apply applies changes to the state
	Apply(changes []types.Change) error
	
	// GetState returns the current state
	GetState() *SchemaState
}

// InMemoryState is an in-memory implementation of StateManager
type InMemoryState struct {
	state *SchemaState
}

// NewInMemoryState creates a new in-memory state manager
func NewInMemoryState() StateManager {
	return &InMemoryState{
		state: &SchemaState{
			Tables: make(map[string]*TableState),
		},
	}
}

// Load returns the current state (no-op for in-memory)
func (s *InMemoryState) Load() (*SchemaState, error) {
	return s.state, nil
}

// Apply applies changes to the state
func (s *InMemoryState) Apply(changes []types.Change) error {
	for _, change := range changes {
		switch c := change.(type) {
		case *types.CreateTable:
			if err := s.applyCreateTable(c); err != nil {
				return err
			}
		case *types.DropTable:
			delete(s.state.Tables, c.Table)
		case *types.RenameTable:
			if table, exists := s.state.Tables[c.OldName]; exists {
				// If new name already exists, delete it first (rename overwrites)
				if c.OldName != c.NewName {
					delete(s.state.Tables, c.NewName)
				}
				table.Name = c.NewName
				s.state.Tables[c.NewName] = table
				if c.OldName != c.NewName {
					delete(s.state.Tables, c.OldName)
				}
			}
		case *types.AddColumn:
			if err := s.applyAddColumn(c); err != nil {
				return err
			}
		case *types.DropColumn:
			if table, exists := s.state.Tables[c.Table]; exists {
				delete(table.Columns, c.ColumnName)
			}
		case *types.ModifyColumn:
			if err := s.applyModifyColumn(c); err != nil {
				return err
			}
		case *types.RenameColumn:
			if table, exists := s.state.Tables[c.Table]; exists {
				if col, exists := table.Columns[c.OldName]; exists {
					col.Name = c.NewName
					table.Columns[c.NewName] = col
					delete(table.Columns, c.OldName)
				}
			}
		case *types.AddIndex:
			if err := s.applyAddIndex(c); err != nil {
				return err
			}
		case *types.DropIndex:
			if table, exists := s.state.Tables[c.Table]; exists {
				delete(table.Indexes, c.IndexName)
			}
		case *types.AddForeignKey:
			if err := s.applyAddForeignKey(c); err != nil {
				return err
			}
		case *types.DropForeignKey:
			if table, exists := s.state.Tables[c.Table]; exists {
				delete(table.ForeignKeys, c.FKName)
			}
		case *types.AddConstraint:
			if err := s.applyAddConstraint(c); err != nil {
				return err
			}
		case *types.DropConstraint:
			if table, exists := s.state.Tables[c.Table]; exists {
				delete(table.Constraints, c.ConstraintName)
			}
		}
	}
	return nil
}

// GetState returns the current state
func (s *InMemoryState) GetState() *SchemaState {
	return s.state
}

func (s *InMemoryState) applyCreateTable(c *types.CreateTable) error {
	tableName := c.TableName()
	tableState := &TableState{
		Name:        tableName,
		Columns:     make(map[string]*ColumnState),
		Indexes:     make(map[string]*IndexState),
		ForeignKeys: make(map[string]*ForeignKeyState),
		Constraints: make(map[string]*ConstraintState),
	}

	// Add columns
	for _, field := range c.Table.Fields {
		colState := &ColumnState{
			Name:          field.Name,
			Type:          field.Type,
			GoType:        field.GoType,
			Required:      field.Required,
			PrimaryKey:    field.PrimaryKey,
			AutoIncrement: field.AutoIncrement,
			Unique:        false,
			Default:       field.Default,
			Options:       field.Options,
		}
		if unique, ok := field.Options["unique"].(bool); ok && unique {
			colState.Unique = unique
		}
		tableState.Columns[field.Name] = colState
	}

	// Add indexes
	for _, idx := range c.Table.Meta.Indexes {
		idxState := &IndexState{
			Name:   idx.Name,
			Fields: idx.Fields,
			Unique: idx.Unique,
		}
		if idxState.Name == "" {
			idxState.Name = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(idx.Fields, "_"))
		}
		tableState.Indexes[idxState.Name] = idxState
	}

	// Add constraints
	for _, constr := range c.Table.Meta.Constraints {
		constrState := &ConstraintState{
			Name:      constr.Name,
			Type:      constr.Type,
			Condition: constr.Condition,
			Fields:    constr.Fields,
		}
		tableState.Constraints[constr.Name] = constrState
	}

	s.state.Tables[tableName] = tableState
	return nil
}

func (s *InMemoryState) applyAddColumn(c *types.AddColumn) error {
	table, exists := s.state.Tables[c.Table]
	if !exists {
		return fmt.Errorf("table %s does not exist", c.Table)
	}

	colState := &ColumnState{
		Name:          c.Column.Name,
		Type:          c.Column.Type,
		GoType:        c.Column.GoType,
		Required:      c.Column.Required,
		PrimaryKey:    c.Column.PrimaryKey,
		AutoIncrement: c.Column.AutoIncrement,
		Unique:        false,
		Default:       c.Column.Default,
		Options:       c.Column.Options,
	}
	if unique, ok := c.Column.Options["unique"].(bool); ok && unique {
		colState.Unique = unique
	}
	table.Columns[c.Column.Name] = colState
	return nil
}

func (s *InMemoryState) applyModifyColumn(c *types.ModifyColumn) error {
	table, exists := s.state.Tables[c.Table]
	if !exists {
		return fmt.Errorf("table %s does not exist", c.Table)
	}

	// Delete old column if name changed
	if c.OldColumn.Name != c.NewColumn.Name {
		delete(table.Columns, c.OldColumn.Name)
	}

	colState := &ColumnState{
		Name:          c.NewColumn.Name,
		Type:          c.NewColumn.Type,
		GoType:        c.NewColumn.GoType,
		Required:      c.NewColumn.Required,
		PrimaryKey:    c.NewColumn.PrimaryKey,
		AutoIncrement: c.NewColumn.AutoIncrement,
		Unique:        false,
		Default:       c.NewColumn.Default,
		Options:       c.NewColumn.Options,
	}
	if unique, ok := c.NewColumn.Options["unique"].(bool); ok && unique {
		colState.Unique = unique
	}
	table.Columns[c.NewColumn.Name] = colState
	return nil
}

func (s *InMemoryState) applyAddIndex(c *types.AddIndex) error {
	table, exists := s.state.Tables[c.Table]
	if !exists {
		return fmt.Errorf("table %s does not exist", c.Table)
	}

	idxState := &IndexState{
		Name:   c.Index.Name,
		Fields: c.Index.Fields,
		Unique: c.Index.Unique,
	}
	if idxState.Name == "" {
		idxState.Name = fmt.Sprintf("idx_%s_%s", c.Table, strings.Join(c.Index.Fields, "_"))
	}
	table.Indexes[idxState.Name] = idxState
	return nil
}

func (s *InMemoryState) applyAddForeignKey(c *types.AddForeignKey) error {
	table, exists := s.state.Tables[c.Table]
	if !exists {
		return fmt.Errorf("table %s does not exist", c.Table)
	}

	fkName := fmt.Sprintf("fk_%s_%s", c.Table, c.Relation.Name)
	fkState := &ForeignKeyState{
		Name:        fkName,
		Column:      c.Relation.Name,
		TargetTable: c.TargetTable,
		TargetColumn: "id", // Default to id
		OnDelete:    "NO ACTION",
		OnUpdate:    "NO ACTION",
	}

	if onDelete, ok := c.Relation.Options["on_delete"].(string); ok {
		fkState.OnDelete = mapCascadeType(onDelete)
	}
	if onUpdate, ok := c.Relation.Options["on_update"].(string); ok {
		fkState.OnUpdate = mapCascadeType(onUpdate)
	}

	table.ForeignKeys[fkName] = fkState
	return nil
}

func (s *InMemoryState) applyAddConstraint(c *types.AddConstraint) error {
	table, exists := s.state.Tables[c.Table]
	if !exists {
		return fmt.Errorf("table %s does not exist", c.Table)
	}

	constrState := &ConstraintState{
		Name:      c.Constraint.Name,
		Type:      c.Constraint.Type,
		Condition: c.Constraint.Condition,
		Fields:    c.Constraint.Fields,
	}
	table.Constraints[c.Constraint.Name] = constrState
	return nil
}

// mapCascadeType maps cascade type strings to SQL
func mapCascadeType(cascade string) string {
	switch cascade {
	case "CASCADE":
		return "CASCADE"
	case "SET_NULL", "SET NULL":
		return "SET NULL"
	case "PROTECT":
		return "RESTRICT"
	case "SET_DEFAULT", "SET DEFAULT":
		return "SET DEFAULT"
	case "DO_NOTHING", "NO ACTION":
		return "NO ACTION"
	default:
		return "NO ACTION"
	}
}
