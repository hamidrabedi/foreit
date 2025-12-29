package sqlgen

import (
	"github.com/forgego/forge/pkg/migrate/errors"
	"github.com/forgego/forge/pkg/migrate/types"
)

// SQLBuilder generates SQL DDL statements from changes
type SQLBuilder interface {
	// BuildUpSQL generates the up migration SQL for a list of changes
	BuildUpSQL(changes []types.Change) (string, error)
	
	// BuildDownSQL generates the down migration SQL for a list of changes
	BuildDownSQL(changes []types.Change) (string, error)
	
	// BuildCreateTable generates CREATE TABLE SQL
	BuildCreateTable(change *types.CreateTable) (string, error)
	
	// BuildAddColumn generates ALTER TABLE ADD COLUMN SQL
	BuildAddColumn(change *types.AddColumn) (string, error)
	
	// BuildModifyColumn generates ALTER TABLE ALTER COLUMN SQL
	BuildModifyColumn(change *types.ModifyColumn) (string, error)
	
	// BuildAddIndex generates CREATE INDEX SQL
	BuildAddIndex(change *types.AddIndex) (string, error)
	
	// BuildModifyIndex generates DROP and CREATE INDEX SQL
	BuildModifyIndex(change *types.ModifyIndex) (string, error)
	
	// BuildAddForeignKey generates ALTER TABLE ADD CONSTRAINT FOREIGN KEY SQL
	BuildAddForeignKey(change *types.AddForeignKey) (string, error)
	
	// BuildModifyForeignKey generates DROP and ADD FOREIGN KEY SQL
	BuildModifyForeignKey(change *types.ModifyForeignKey) (string, error)
	
	// BuildAddConstraint generates ALTER TABLE ADD CONSTRAINT SQL
	BuildAddConstraint(change *types.AddConstraint) (string, error)
}

// NewSQLBuilder creates a new SQL builder for the given driver
func NewSQLBuilder(driver types.Driver) (SQLBuilder, error) {
	switch {
	case driver.IsPostgreSQL():
		return NewPostgreSQLBuilder(), nil
	case driver.IsSQLite():
		return NewSQLiteBuilder(), nil
	default:
		return nil, errors.NewMigrationError(
			errors.ErrInvalidChange,
			"unsupported database driver: "+driver.String(),
			nil,
		)
	}
}
