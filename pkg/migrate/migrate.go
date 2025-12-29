// Package migrate provides a public API for migration generation, state management,
// and SQL generation. This is the main entry point for the migration system.
//
// Core types are re-exported from the types sub-package for convenience.
package migrate

import (
	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrate/types"
)

// ChangeDetector is the interface for detecting changes
type ChangeDetector interface {
	DetectChanges(current, previous []*generator.ModelDefinition) ([]Change, error)
}

// Re-export types for convenience
type (
	Change     = types.Change
	ChangeType = types.ChangeType
	Driver     = types.Driver
)

// Re-export change types
type (
	CreateTable      = types.CreateTable
	DropTable        = types.DropTable
	RenameTable      = types.RenameTable
	AddColumn        = types.AddColumn
	DropColumn       = types.DropColumn
	ModifyColumn     = types.ModifyColumn
	RenameColumn     = types.RenameColumn
	AddIndex         = types.AddIndex
	DropIndex        = types.DropIndex
	ModifyIndex      = types.ModifyIndex
	AddForeignKey    = types.AddForeignKey
	DropForeignKey   = types.DropForeignKey
	ModifyForeignKey = types.ModifyForeignKey
	AddConstraint    = types.AddConstraint
	DropConstraint   = types.DropConstraint
	RunSQL           = types.RunSQL
	RunGo            = types.RunGo
	UnknownChange    = types.UnknownChange
)

// Re-export constants
const (
	ChangeTypeCreateTable      = types.ChangeTypeCreateTable
	ChangeTypeDropTable        = types.ChangeTypeDropTable
	ChangeTypeRenameTable      = types.ChangeTypeRenameTable
	ChangeTypeAddColumn        = types.ChangeTypeAddColumn
	ChangeTypeDropColumn       = types.ChangeTypeDropColumn
	ChangeTypeModifyColumn     = types.ChangeTypeModifyColumn
	ChangeTypeRenameColumn     = types.ChangeTypeRenameColumn
	ChangeTypeAddIndex         = types.ChangeTypeAddIndex
	ChangeTypeDropIndex        = types.ChangeTypeDropIndex
	ChangeTypeModifyIndex      = types.ChangeTypeModifyIndex
	ChangeTypeAddForeignKey    = types.ChangeTypeAddForeignKey
	ChangeTypeDropForeignKey   = types.ChangeTypeDropForeignKey
	ChangeTypeModifyForeignKey = types.ChangeTypeModifyForeignKey
	ChangeTypeAddConstraint    = types.ChangeTypeAddConstraint
	ChangeTypeDropConstraint   = types.ChangeTypeDropConstraint
	ChangeTypeRunSQL           = types.ChangeTypeRunSQL
	ChangeTypeRunGo            = types.ChangeTypeRunGo
	ChangeTypeUnknown          = types.ChangeTypeUnknown
	DriverPostgreSQL           = types.DriverPostgreSQL
	DriverSQLite               = types.DriverSQLite
	DriverSQLite3              = types.DriverSQLite3
)

// Re-export migration types
type (
	Migration     = types.Migration
	MigrationPlan = types.MigrationPlan
	Dependency    = types.Dependency
)

// NewMigrationPlan creates a new migration plan
func NewMigrationPlan(version, name string, changes []Change) *MigrationPlan {
	return types.NewMigrationPlan(version, name, changes)
}
