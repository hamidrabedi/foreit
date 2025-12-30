// Package migrate provides a public API for migration generation, state management,
// and SQL generation. This package re-exports types and functions from forge/db/migrate
// for convenience and backward compatibility.
package migrate

import (
	codegen "github.com/forgego/forge/codegen"
	dbmigrate "github.com/forgego/forge/db/migrate"
)

// Re-export all types and functions from forge/db/migrate
type (
	Change         = dbmigrate.Change
	ChangeType     = dbmigrate.ChangeType
	Driver         = dbmigrate.Driver
	ChangeDetector = dbmigrate.ChangeDetector
)

// Re-export change types
type (
	CreateTable      = dbmigrate.CreateTable
	DropTable        = dbmigrate.DropTable
	RenameTable      = dbmigrate.RenameTable
	AddColumn        = dbmigrate.AddColumn
	DropColumn       = dbmigrate.DropColumn
	ModifyColumn     = dbmigrate.ModifyColumn
	RenameColumn     = dbmigrate.RenameColumn
	AddIndex         = dbmigrate.AddIndex
	DropIndex        = dbmigrate.DropIndex
	ModifyIndex      = dbmigrate.ModifyIndex
	AddForeignKey    = dbmigrate.AddForeignKey
	DropForeignKey   = dbmigrate.DropForeignKey
	ModifyForeignKey = dbmigrate.ModifyForeignKey
	AddConstraint    = dbmigrate.AddConstraint
	DropConstraint   = dbmigrate.DropConstraint
	RunSQL           = dbmigrate.RunSQL
	RunGo            = dbmigrate.RunGo
	UnknownChange    = dbmigrate.UnknownChange
)

// Re-export constants
const (
	ChangeTypeCreateTable      = dbmigrate.ChangeTypeCreateTable
	ChangeTypeDropTable        = dbmigrate.ChangeTypeDropTable
	ChangeTypeRenameTable      = dbmigrate.ChangeTypeRenameTable
	ChangeTypeAddColumn        = dbmigrate.ChangeTypeAddColumn
	ChangeTypeDropColumn       = dbmigrate.ChangeTypeDropColumn
	ChangeTypeModifyColumn     = dbmigrate.ChangeTypeModifyColumn
	ChangeTypeRenameColumn     = dbmigrate.ChangeTypeRenameColumn
	ChangeTypeAddIndex         = dbmigrate.ChangeTypeAddIndex
	ChangeTypeDropIndex        = dbmigrate.ChangeTypeDropIndex
	ChangeTypeModifyIndex      = dbmigrate.ChangeTypeModifyIndex
	ChangeTypeAddForeignKey    = dbmigrate.ChangeTypeAddForeignKey
	ChangeTypeDropForeignKey   = dbmigrate.ChangeTypeDropForeignKey
	ChangeTypeModifyForeignKey = dbmigrate.ChangeTypeModifyForeignKey
	ChangeTypeAddConstraint    = dbmigrate.ChangeTypeAddConstraint
	ChangeTypeDropConstraint   = dbmigrate.ChangeTypeDropConstraint
	ChangeTypeRunSQL           = dbmigrate.ChangeTypeRunSQL
	ChangeTypeRunGo            = dbmigrate.ChangeTypeRunGo
	ChangeTypeUnknown          = dbmigrate.ChangeTypeUnknown
	DriverPostgreSQL           = dbmigrate.DriverPostgreSQL
	DriverSQLite               = dbmigrate.DriverSQLite
	DriverSQLite3              = dbmigrate.DriverSQLite3
)

// Re-export migration types
type (
	Migration     = dbmigrate.Migration
	MigrationPlan = dbmigrate.MigrationPlan
	Dependency    = dbmigrate.Dependency
	State         = dbmigrate.State
	Generator     = dbmigrate.Generator
	SQLGenerator  = dbmigrate.SQLGenerator
)

// Re-export functions
var (
	NewMigrationPlan = dbmigrate.NewMigrationPlan
	Generate         = dbmigrate.Generate
	LoadState        = dbmigrate.LoadState
	NewDetector      = dbmigrate.NewDetector
	NewSQLGenerator  = dbmigrate.NewSQLGenerator
	NewGenerator     = dbmigrate.NewGenerator
)

// DetectChanges compares current models to previous state and returns all changes.
// This is a convenience wrapper that matches the CLI usage.
func DetectChanges(current, previous []*codegen.ModelDefinition) ([]Change, error) {
	return dbmigrate.DetectChanges(current, previous)
}
