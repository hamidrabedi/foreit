// Package migrations provides migration generation and management.
// This package re-exports types from sub-packages for backward compatibility.
package migrations

import (
	"github.com/forgego/forge/pkg/migrations/core"
)

// Re-export core types
type (
	Change     = core.Change
	ChangeType = core.ChangeType
)

// Re-export change types
type (
	CreateTable     = core.CreateTable
	DropTable       = core.DropTable
	RenameTable     = core.RenameTable
	AddColumn       = core.AddColumn
	DropColumn      = core.DropColumn
	ModifyColumn    = core.ModifyColumn
	RenameColumn    = core.RenameColumn
	AddIndex        = core.AddIndex
	DropIndex       = core.DropIndex
	ModifyIndex     = core.ModifyIndex
	AddForeignKey   = core.AddForeignKey
	DropForeignKey  = core.DropForeignKey
	ModifyForeignKey = core.ModifyForeignKey
	AddConstraint   = core.AddConstraint
	DropConstraint  = core.DropConstraint
)

// Re-export constants
const (
	ChangeTypeCreateTable     = core.ChangeTypeCreateTable
	ChangeTypeDropTable       = core.ChangeTypeDropTable
	ChangeTypeRenameTable     = core.ChangeTypeRenameTable
	ChangeTypeAddColumn       = core.ChangeTypeAddColumn
	ChangeTypeDropColumn      = core.ChangeTypeDropColumn
	ChangeTypeModifyColumn    = core.ChangeTypeModifyColumn
	ChangeTypeRenameColumn    = core.ChangeTypeRenameColumn
	ChangeTypeAddIndex        = core.ChangeTypeAddIndex
	ChangeTypeDropIndex       = core.ChangeTypeDropIndex
	ChangeTypeModifyIndex     = core.ChangeTypeModifyIndex
	ChangeTypeAddForeignKey   = core.ChangeTypeAddForeignKey
	ChangeTypeDropForeignKey  = core.ChangeTypeDropForeignKey
	ChangeTypeModifyForeignKey = core.ChangeTypeModifyForeignKey
	ChangeTypeAddConstraint   = core.ChangeTypeAddConstraint
	ChangeTypeDropConstraint  = core.ChangeTypeDropConstraint
)
