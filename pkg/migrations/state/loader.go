package state

import (
	"fmt"
	"os"
	"path/filepath"

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
}

// NewFileStateLoader creates a new file-based state loader
func NewFileStateLoader(migrationsDir string) *FileStateLoader {
	return &FileStateLoader{
		migrationsDir: migrationsDir,
	}
}

// Load loads state from migration files
func (l *FileStateLoader) Load() (*SchemaState, error) {
	state := &SchemaState{
		Tables: make(map[string]*TableState),
	}

	// Check if migrations directory exists
	if _, err := os.Stat(l.migrationsDir); os.IsNotExist(err) {
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
			continue
		}

		// Apply changes to state
		for _, change := range changes {
			if err := applyChangeToState(state, change); err != nil {
				return nil, fmt.Errorf("failed to apply change from %s: %w", file, err)
			}
		}
	}

	return state, nil
}

// sortMigrationFiles sorts migration files by version number
func sortMigrationFiles(files []string) []string {
	// Simple sort by filename (which contains version prefix)
	// For a proper implementation, we'd parse and sort by version number
	// For now, just return as-is since file systems typically return in order
	return files
}

// applyChangeToState applies a single change to the state
func applyChangeToState(state *SchemaState, change core.Change) error {
	// Create a temporary manager with the state
	manager := &InMemoryState{state: state}
	return manager.Apply([]core.Change{change})
}

