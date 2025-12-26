// Package migrations provides migration generation and management.
package migrations

import (
	"github.com/forgego/forge/pkg/migrations/core"
)

// MigrationPlan represents a complete migration plan
type MigrationPlan = core.MigrationPlan

// NewMigrationPlan creates a new migration plan
func NewMigrationPlan(version, name string, changes []Change) *MigrationPlan {
	// Change is a type alias for core.Change, so we can pass directly
	return core.NewMigrationPlan(version, name, changes)
}
