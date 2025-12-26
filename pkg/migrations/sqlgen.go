// Package migrations provides migration generation and management.
package migrations

import (
	"github.com/forgego/forge/pkg/migrations/core"
	"github.com/forgego/forge/pkg/migrations/sql"
)

// SQLGenerator generates SQL DDL statements from changes
// This is a compatibility wrapper around the new sql package.
type SQLGenerator struct {
	builder sql.SQLBuilder
	driver  string
}

// NewSQLGenerator creates a new SQL generator for the given driver
func NewSQLGenerator(driver string) *SQLGenerator {
	driverType := core.Driver(driver)
	builder, err := sql.NewSQLBuilder(driverType)
	if err != nil {
		// For compatibility, we'll create a builder that will error on use
		// In practice, this should not happen with valid drivers
		panic(err)
	}
	return &SQLGenerator{
		builder: builder,
		driver:  driver,
	}
}

// GenerateUpSQL generates the up migration SQL for a list of changes
func (g *SQLGenerator) GenerateUpSQL(changes []Change) (string, error) {
	// Change is a type alias for core.Change, so we can pass directly
	return g.builder.BuildUpSQL(changes)
}

// GenerateDownSQL generates the down migration SQL for a list of changes
func (g *SQLGenerator) GenerateDownSQL(changes []Change) (string, error) {
	// Change is a type alias for core.Change, so we can pass directly
	return g.builder.BuildDownSQL(changes)
}
