package operations

import (
	"github.com/forgego/forge/pkg/migrate/schema"
)

// RunPython executes Python code during migration
// In Go, this would be RunGo which executes Go code
type RunPython struct {
	*BaseOperation
	Code      func(apps interface{}) error // Forward function
	Reverse   func(apps interface{}) error // Reverse function (optional)
	Hints     map[string]interface{}        // Hints for historical models
	Elidable  bool                          // Can be elided if previous migration has same operation
}

// NewRunPython creates a new RunPython operation
func NewRunPython(code func(apps interface{}) error, reverse func(apps interface{}) error) *RunPython {
	reversible := reverse != nil
	return &RunPython{
		BaseOperation: NewBaseOperation(reversible, false, CategoryPython),
		Code:          code,
		Reverse:       reverse,
		Hints:         make(map[string]interface{}),
		Elidable:      false,
	}
}

// StateForwards is a no-op for RunPython (doesn't affect state)
func (op *RunPython) StateForwards(appLabel string, s *schema.SchemaState) error {
	return nil
}

// DatabaseForwards executes the Python/Go code
func (op *RunPython) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	// In Go implementation, this would execute Go code
	// For now, we'll need a way to execute Go functions
	// This requires a runtime or plugin system
	return op.Code(nil) // Pass historical models/apps here
}

// DatabaseBackwards executes the reverse code
func (op *RunPython) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	if op.Reverse == nil {
		return ErrIrreversible
	}
	return op.Reverse(nil) // Pass historical models/apps here
}

// Describe returns a human-readable description
func (op *RunPython) Describe() string {
	return "Run Python code"
}

// RunSQL executes raw SQL during migration
type RunSQL struct {
	*BaseOperation
	SQL         string
	ReverseSQL  string
	StateApps   map[string]interface{} // State changes (for SeparateDatabaseAndState)
	Hints       map[string]interface{}
}

// NewRunSQL creates a new RunSQL operation
func NewRunSQL(sql string, reverseSQL string) *RunSQL {
	reversible := reverseSQL != ""
	return &RunSQL{
		BaseOperation: NewBaseOperation(reversible, true, CategorySQL),
		SQL:           sql,
		ReverseSQL:   reverseSQL,
		StateApps:    make(map[string]interface{}),
		Hints:        make(map[string]interface{}),
	}
}

// StateForwards applies state changes if any
func (op *RunSQL) StateForwards(appLabel string, s *schema.SchemaState) error {
	// If StateApps is set, apply state changes
	// This is used with SeparateDatabaseAndState
	return nil
}

// DatabaseForwards executes the SQL
func (op *RunSQL) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	return schemaEditor.ExecuteSQL(op.SQL)
}

// DatabaseBackwards executes the reverse SQL
func (op *RunSQL) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	if op.ReverseSQL == "" {
		return ErrIrreversible
	}
	return schemaEditor.ExecuteSQL(op.ReverseSQL)
}

// Describe returns a human-readable description
func (op *RunSQL) Describe() string {
	return "Run SQL"
}

// SeparateDatabaseAndState allows you to separate database changes from state changes
// This is useful when you want to change the database schema without affecting
// the migration state, or vice versa
type SeparateDatabaseAndState struct {
	*BaseOperation
	DatabaseOperations []Operation // Operations that affect the database
	StateOperations    []Operation // Operations that only affect state
}

// NewSeparateDatabaseAndState creates a new SeparateDatabaseAndState operation
func NewSeparateDatabaseAndState(databaseOps, stateOps []Operation) *SeparateDatabaseAndState {
	return &SeparateDatabaseAndState{
		BaseOperation:      NewBaseOperation(true, false, CategoryMixed),
		DatabaseOperations: databaseOps,
		StateOperations:    stateOps,
	}
}

// StateForwards applies only state operations
func (op *SeparateDatabaseAndState) StateForwards(appLabel string, s *schema.SchemaState) error {
	for _, stateOp := range op.StateOperations {
		if err := stateOp.StateForwards(appLabel, s); err != nil {
			return err
		}
	}
	return nil
}

// DatabaseForwards applies only database operations
func (op *SeparateDatabaseAndState) DatabaseForwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	for _, dbOp := range op.DatabaseOperations {
		if err := dbOp.DatabaseForwards(appLabel, schemaEditor, fromState, toState); err != nil {
			return err
		}
	}
	return nil
}

// DatabaseBackwards reverses database operations
func (op *SeparateDatabaseAndState) DatabaseBackwards(appLabel string, schemaEditor SchemaEditor, fromState, toState *schema.SchemaState) error {
	// Reverse in opposite order
	for i := len(op.DatabaseOperations) - 1; i >= 0; i-- {
		dbOp := op.DatabaseOperations[i]
		if err := dbOp.DatabaseBackwards(appLabel, schemaEditor, fromState, toState); err != nil {
			return err
		}
	}
	return nil
}

// Describe returns a human-readable description
func (op *SeparateDatabaseAndState) Describe() string {
	return "Separate database and state operations"
}

// ErrIrreversible is returned when trying to reverse an irreversible operation
var ErrIrreversible = &IrreversibleError{}

// IrreversibleError represents an error when trying to reverse an irreversible operation
type IrreversibleError struct{}

func (e *IrreversibleError) Error() string {
	return "operation is not reversible"
}
