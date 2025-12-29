package state

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/forgego/forge/pkg/migrations/core"
)

// LoadStateFromFiles loads state from migration files
func LoadStateFromFiles(migrationsDir string) (*SchemaState, error) {
	loader := NewFileStateLoader(migrationsDir)
	return loader.Load()
}

// FileStateLoader loads state from SQL migration files
type FileStateLoader struct {
	migrationsDir string
	state         *SchemaState
}

// NewFileStateLoader creates a new file-based state loader
func NewFileStateLoader(migrationsDir string) StateManager {
	return &FileStateLoader{
		migrationsDir: migrationsDir,
		state:         nil, // Will be loaded on first Load() call
	}
}

// Load loads state from migration files
func (l *FileStateLoader) Load() (*SchemaState, error) {
	// Load state if not already loaded
	if l.state == nil {
		state := &SchemaState{
			Tables: make(map[string]*TableState),
		}

		// Check if migrations directory exists
		if _, err := os.Stat(l.migrationsDir); os.IsNotExist(err) {
			l.state = state
			return state, nil
		}

		// Find all up migration files
		pattern := filepath.Join(l.migrationsDir, "*_*.up.sql")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migrations directory: %w", err)
		}

		// Sort files by version
		sortedFiles := sortMigrationFiles(matches)

		// Parse each migration file and apply to state
		parser := NewSQLParser()
		for _, file := range sortedFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read migration file %s: %w", file, err)
			}

			changes, err := parser.ParseUpSQL(string(content))
			if err != nil {
				// Skip files that can't be parsed (might be manual SQL)
				// TODO: Add optional logging/warning for unparseable migration files
				// TODO: Consider storing parse errors for user review
				continue
			}

			// Sort changes to ensure CREATE TABLE comes before ALTER TABLE
			changes = sortChangesByType(changes)

			// Apply changes to state
			for _, change := range changes {
				// Skip UnknownChange - these are unparseable statements
				// They don't affect schema state reconstruction
				if _, ok := change.(*core.UnknownChange); ok {
					// Log or skip silently - UnknownChange doesn't affect state
					// TODO: Add optional verbose mode to log UnknownChange statements for debugging
					continue
				}

				if err := applyChangeToState(state, change); err != nil {
					// If table doesn't exist error, it might be because the change is out of order
					// Try to continue - this is a limitation of the SQL parser
					if strings.Contains(err.Error(), "does not exist") {
						// Skip this change - it might be a constraint on a table that will be created later
						// or it's a constraint that was already added
						// TODO: Improve change ordering logic to handle out-of-order constraints better
						// TODO: Add retry mechanism after all CREATE TABLE statements are processed
						continue
					}
					return nil, fmt.Errorf("failed to apply change from %s: %w", file, err)
				}
			}
		}

		l.state = state
	}

	return l.state, nil
}

// Apply applies changes to the state (required for StateManager interface)
func (l *FileStateLoader) Apply(changes []core.Change) error {
	// Load state first if not loaded
	if l.state == nil {
		_, err := l.Load()
		if err != nil {
			return err
		}
	}

	// Use InMemoryState to apply changes
	manager := &InMemoryState{state: l.state}
	return manager.Apply(changes)
}

// GetState returns the current state (required for StateManager interface)
func (l *FileStateLoader) GetState() *SchemaState {
	if l.state == nil {
		// Load state if not loaded
		// If Load fails, return empty state to avoid nil pointer
		state, err := l.Load()
		if err != nil {
			// Return empty state on error rather than nil
			return &SchemaState{
				Tables: make(map[string]*TableState),
			}
		}
		return state
	}
	return l.state
}

// applyChangeToState applies a single change to the state
func applyChangeToState(state *SchemaState, change core.Change) error {
	// Create a temporary manager with the state
	manager := &InMemoryState{state: state}
	return manager.Apply([]core.Change{change})
}
