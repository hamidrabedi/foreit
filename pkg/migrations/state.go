// Package migrations provides migration generation and management.
package migrations

import (
	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/state"
)

// State represents the current database schema state
// This is a compatibility wrapper around the new state package.
type State struct {
	manager state.StateManager
}

// LoadState builds the current state from existing migration files
func LoadState(migrationsDir string) (*State, error) {
	// Try to load from files first
	schemaState, err := state.LoadStateFromFiles(migrationsDir)
	if err != nil {
		// If loading fails, return empty state
		// This maintains backward compatibility
		return &State{
			manager: state.NewInMemoryState(),
		}, nil
	}

	// Create a state manager with the loaded state
	manager := state.NewInMemoryState()
	managerState := manager.GetState()
	*managerState = *schemaState

	return &State{
		manager: manager,
	}, nil
}

// ApplyMigration updates the state after applying a migration
func (s *State) ApplyMigration(changes []Change) error {
	// Change is a type alias for core.Change, so we can pass directly
	return s.manager.Apply(changes)
}

// ToModelDefinitions converts state back to ModelDefinitions (for comparison)
func (s *State) ToModelDefinitions() []*generator.ModelDefinition {
	schemaState := s.manager.GetState()
	return schemaState.ToModelDefinitions()
}
