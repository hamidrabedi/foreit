// Package migrations provides migration generation and management.
package migrations

import (
	"github.com/forgego/forge/pkg/migrations/generation"
)

// Generator generates migrations from model definitions
// This is a compatibility wrapper around the new generation package.
type Generator struct {
	generator *generation.MigrationGenerator
}

// NewGenerator creates a new migration generator
func NewGenerator(modelsDir, migrationsDir string) (*Generator, error) {
	gen, err := generation.NewMigrationGeneratorWithDefaults(modelsDir, migrationsDir)
	if err != nil {
		return nil, err
	}
	return &Generator{
		generator: gen,
	}, nil
}

// GenerateMigrations generates migration files from model definitions
func (g *Generator) GenerateMigrations(name string) error {
	return g.generator.GenerateMigrations(name)
}
