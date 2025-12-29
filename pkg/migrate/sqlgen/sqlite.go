package sqlgen

import (
	"fmt"
	"strings"

	"github.com/forgego/forge/pkg/migrate/errors"
	"github.com/forgego/forge/pkg/migrate/types"
)

// SQLiteBuilder implements SQLBuilder for SQLite
type SQLiteBuilder struct {
	*baseBuilder
}

// NewSQLiteBuilder creates a new SQLite SQL builder
func NewSQLiteBuilder() SQLBuilder {
	return &SQLiteBuilder{
		baseBuilder: &baseBuilder{
			isSQLite:  true,
			isPostgres: false,
		},
	}
}

// BuildUpSQL generates the up migration SQL for a list of changes
func (b *SQLiteBuilder) BuildUpSQL(changes []types.Change) (string, error) {
	var statements []string

	// Order changes by dependency
	orderedChanges := orderChanges(changes)

	for _, change := range orderedChanges {
		sql, err := b.buildChangeUpSQL(change)
		if err != nil {
			return "", errors.NewMigrationError(
				errors.ErrInvalidChange,
				fmt.Sprintf("failed to generate SQL for %s", change.Type()),
				err,
			)
		}
		if sql != "" {
			statements = append(statements, sql)
		}
	}

	return strings.Join(statements, "\n\n"), nil
}

// BuildDownSQL generates the down migration SQL for a list of changes
func (b *SQLiteBuilder) BuildDownSQL(changes []types.Change) (string, error) {
	var statements []string

	// Reverse order for down migrations
	orderedChanges := orderChanges(changes)
	for i := len(orderedChanges) - 1; i >= 0; i-- {
		change := orderedChanges[i]
		sql, err := b.buildChangeDownSQL(change)
		if err != nil {
			return "", errors.NewMigrationError(
				errors.ErrInvalidChange,
				fmt.Sprintf("failed to generate down SQL for %s", change.Type()),
				err,
			)
		}
		if sql != "" {
			statements = append(statements, sql)
		}
	}

	return strings.Join(statements, "\n\n"), nil
}

// buildChangeUpSQL generates up SQL for a single change
func (b *SQLiteBuilder) buildChangeUpSQL(change types.Change) (string, error) {
	switch c := change.(type) {
	case *types.CreateTable:
		return b.BuildCreateTable(c)
	case *types.DropTable:
		return fmt.Sprintf("DROP TABLE IF EXISTS %s;", c.Table), nil
	case *types.RenameTable:
		return fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", c.OldName, c.NewName), nil
	case *types.AddColumn:
		return b.BuildAddColumn(c)
	case *types.DropColumn:
		// SQLite doesn't support DROP COLUMN directly
		return fmt.Sprintf("-- SQLite does not support DROP COLUMN directly\n-- Column %s.%s should be dropped manually", c.Table, c.ColumnName), nil
	case *types.ModifyColumn:
		return b.BuildModifyColumn(c)
	case *types.RenameColumn:
		return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", c.Table, c.OldName, c.NewName), nil
	case *types.AddIndex:
		return b.BuildAddIndex(c)
	case *types.DropIndex:
		return fmt.Sprintf("DROP INDEX IF EXISTS %s;", c.IndexName), nil
	case *types.ModifyIndex:
		return b.BuildModifyIndex(c)
	case *types.AddForeignKey:
		return b.BuildAddForeignKey(c)
	case *types.DropForeignKey:
		return fmt.Sprintf("-- SQLite does not support DROP FOREIGN KEY directly\n-- Foreign key %s should be dropped manually", c.FKName), nil
	case *types.ModifyForeignKey:
		return b.BuildModifyForeignKey(c)
	case *types.AddConstraint:
		return b.BuildAddConstraint(c)
	case *types.DropConstraint:
		return fmt.Sprintf("-- SQLite has limited constraint support\n-- Constraint %s should be dropped manually", c.ConstraintName), nil
	case *types.RunSQL:
		// For RunSQL, return the SQL directly
		return c.SQL, nil
	case *types.RunGo:
		// For RunGo, embed as a comment with special marker for execution
		// The actual execution will be handled by the runner
		return fmt.Sprintf("-- RUNGO: %s\n-- This migration requires Go code execution", c.UpFunc), nil
	default:
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			fmt.Sprintf("unknown change type: %T", change),
			nil,
		)
	}
}

// buildChangeDownSQL generates down SQL for a single change
func (b *SQLiteBuilder) buildChangeDownSQL(change types.Change) (string, error) {
	switch c := change.(type) {
	case *types.CreateTable:
		return fmt.Sprintf("DROP TABLE IF EXISTS %s;", c.TableName()), nil
	case *types.DropTable:
		// Cannot generate down SQL for DropTable without original table definition
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"cannot generate down SQL for DropTable without table definition",
			nil,
		)
	case *types.RenameTable:
		return fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", c.NewName, c.OldName), nil
	case *types.AddColumn:
		return fmt.Sprintf("-- SQLite does not support DROP COLUMN directly\n-- Column %s.%s should be dropped manually", c.Table, c.Column.Name), nil
	case *types.DropColumn:
		// Cannot generate down SQL for DropColumn without original column definition
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"cannot generate down SQL for DropColumn without column definition",
			nil,
		)
	case *types.ModifyColumn:
		// Reverse the modification
		return b.buildModifyColumnDown(c)
	case *types.RenameColumn:
		return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", c.Table, c.NewName, c.OldName), nil
	case *types.AddIndex:
		indexName := c.Index.Name
		if indexName == "" {
			indexName = fmt.Sprintf("idx_%s_%s", c.Table, strings.Join(c.Index.Fields, "_"))
		}
		return fmt.Sprintf("DROP INDEX IF EXISTS %s;", indexName), nil
	case *types.DropIndex:
		// Cannot generate down SQL for DropIndex without original index definition
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"cannot generate down SQL for DropIndex without index definition",
			nil,
		)
	case *types.ModifyIndex:
		// Reverse the modification
		reversed := &types.ModifyIndex{
			Table:   c.Table,
			OldIndex: c.NewIndex,
			NewIndex: c.OldIndex,
		}
		return b.BuildModifyIndex(reversed)
	case *types.AddForeignKey:
		fkName := fmt.Sprintf("fk_%s_%s", c.Table, c.Relation.Name)
		return fmt.Sprintf("-- SQLite does not support DROP FOREIGN KEY directly\n-- Foreign key %s should be dropped manually", fkName), nil
	case *types.DropForeignKey:
		// Cannot generate down SQL for DropForeignKey without original FK definition
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"cannot generate down SQL for DropForeignKey without FK definition",
			nil,
		)
	case *types.ModifyForeignKey:
		// Reverse the modification
		reversed := &types.ModifyForeignKey{
			Table:       c.Table,
			OldFK:       c.NewFK,
			NewFK:       c.OldFK,
			TargetTable: c.TargetTable,
		}
		return b.BuildModifyForeignKey(reversed)
	case *types.AddConstraint:
		return fmt.Sprintf("-- SQLite has limited constraint support\n-- Constraint %s should be dropped manually", c.Constraint.Name), nil
	case *types.DropConstraint:
		// Cannot generate down SQL for DropConstraint without original constraint definition
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"cannot generate down SQL for DropConstraint without constraint definition",
			nil,
		)
	case *types.RunSQL:
		// For RunSQL, return the reverse SQL if provided
		if c.ReverseSQL != "" {
			return c.ReverseSQL, nil
		}
		if !c.CanReverse {
			return "-- This data migration is not reversible", nil
		}
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"RunSQL migration is marked as reversible but no ReverseSQL provided",
			nil,
		)
	case *types.RunGo:
		// For RunGo, embed as a comment with special marker for execution
		if c.DownFunc != "" {
			return fmt.Sprintf("-- RUNGO: %s\n-- This migration requires Go code execution for rollback", c.DownFunc), nil
		}
		if !c.CanReverse {
			return "-- This data migration is not reversible", nil
		}
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			"RunGo migration is marked as reversible but no DownFunc provided",
			nil,
		)
	default:
		return "", errors.NewMigrationError(
			errors.ErrInvalidChange,
			fmt.Sprintf("unknown change type: %T", change),
			nil,
		)
	}
}

// BuildModifyColumn generates ALTER TABLE ALTER COLUMN statement for SQLite
// SQLite has very limited ALTER TABLE support, so this generates a comment
func (b *SQLiteBuilder) BuildModifyColumn(c *types.ModifyColumn) (string, error) {
	newType := mapFieldTypeToSQL(c.NewColumn, true, false)
	oldType := mapFieldTypeToSQL(c.OldColumn, true, false)
	
	// SQLite has limited ALTER TABLE support
	// This would require table recreation in practice
	return fmt.Sprintf("-- SQLite does not support ALTER COLUMN directly\n-- Column %s.%s type changed from %s to %s\n-- This requires table recreation in practice",
		c.Table, c.NewColumn.Name, oldType, newType), nil
}

// buildModifyColumnDown generates the reverse of ModifyColumn
func (b *SQLiteBuilder) buildModifyColumnDown(c *types.ModifyColumn) (string, error) {
	reversed := &types.ModifyColumn{
		Table:     c.Table,
		OldColumn: c.NewColumn,
		NewColumn: c.OldColumn,
	}
	return b.BuildModifyColumn(reversed)
}

// BuildCreateTable delegates to baseBuilder
func (b *SQLiteBuilder) BuildCreateTable(c *types.CreateTable) (string, error) {
	return b.baseBuilder.BuildCreateTable(c)
}

// BuildAddColumn delegates to baseBuilder
func (b *SQLiteBuilder) BuildAddColumn(c *types.AddColumn) (string, error) {
	return b.baseBuilder.BuildAddColumn(c)
}

// BuildAddIndex delegates to baseBuilder
func (b *SQLiteBuilder) BuildAddIndex(c *types.AddIndex) (string, error) {
	return b.baseBuilder.BuildAddIndex(c)
}

// BuildModifyIndex delegates to baseBuilder
func (b *SQLiteBuilder) BuildModifyIndex(c *types.ModifyIndex) (string, error) {
	return b.baseBuilder.BuildModifyIndex(c)
}

// BuildAddForeignKey delegates to baseBuilder
func (b *SQLiteBuilder) BuildAddForeignKey(c *types.AddForeignKey) (string, error) {
	return b.baseBuilder.BuildAddForeignKey(c)
}

// BuildModifyForeignKey delegates to baseBuilder
func (b *SQLiteBuilder) BuildModifyForeignKey(c *types.ModifyForeignKey) (string, error) {
	return b.baseBuilder.BuildModifyForeignKey(c)
}

// BuildAddConstraint delegates to baseBuilder
func (b *SQLiteBuilder) BuildAddConstraint(c *types.AddConstraint) (string, error) {
	return b.baseBuilder.BuildAddConstraint(c)
}
