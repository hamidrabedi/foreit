// Package generate re-exports the db/migrate/generate package
package generate

import (
	"github.com/forgego/forge/db/migrate/generate"
)

// Re-export types
type (
	MigrationGenerator = generate.MigrationGenerator
	Squasher           = generate.Squasher
)

// Re-export functions
var (
	NewMigrationGenerator            = generate.NewMigrationGenerator
	NewMigrationGeneratorWithDefaults = generate.NewMigrationGeneratorWithDefaults
	NewDetector                      = generate.NewDetector
	NewSquasher                      = generate.NewSquasher
)

