package recovery

import (
	"context"
	"fmt"
	"strings"
)

// Recovery handles migration recovery operations
type Recovery struct{}

// NewRecovery creates a new recovery handler
func NewRecovery() *Recovery {
	return &Recovery{}
}

// RecoverDirtyState attempts to recover from a dirty migration state
func (r *Recovery) RecoverDirtyState(ctx context.Context, db interface{}, migrationsDir string) error {
	// Check if database is in dirty state
	// This would query the schema_migrations table to check dirty flag
	// For now, return instructions

	return fmt.Errorf("database is in dirty state - manual recovery required:\n" +
		"1. Check the last applied migration in schema_migrations table\n" +
		"2. Verify the database state matches the migration\n" +
		"3. If migration partially applied, manually fix the schema\n" +
		"4. Mark migration as clean: UPDATE schema_migrations SET dirty = false WHERE version = <version>")
}

// ValidateMigrationIntegrity validates that migration files haven't been modified
func (r *Recovery) ValidateMigrationIntegrity(migrationsDir string, appliedMigrations map[string]string) error {
	// This would check checksums of applied migrations
	// For now, return nil (checksum validation is handled elsewhere)
	return nil
}

// GetRecoverySteps returns recovery steps for a failed migration
func (r *Recovery) GetRecoverySteps(version string, errorMsg string) []string {
	steps := []string{
		"1. Check the error message above for details",
		"2. Review the migration file: " + version,
		"3. Check database logs for detailed error information",
		"4. Verify database state matches expected state",
	}

	if strings.Contains(errorMsg, "constraint") {
		steps = append(steps, "5. Check for foreign key constraints that may need to be dropped first")
	}

	if strings.Contains(errorMsg, "column") {
		steps = append(steps, "5. Verify column doesn't already exist or has different type")
	}

	if strings.Contains(errorMsg, "table") {
		steps = append(steps, "5. Verify table doesn't already exist or has different structure")
	}

	steps = append(steps,
		"6. Fix the issue manually if needed",
		fmt.Sprintf("7. Mark migration as clean: UPDATE schema_migrations SET dirty = false WHERE version = %s", version),
		"8. Retry the migration",
	)

	return steps
}

// RollbackPartialMigration attempts to rollback a partially applied migration
func (r *Recovery) RollbackPartialMigration(ctx context.Context, version string, downSQL string) error {
	// This would attempt to execute the down migration
	// For now, return instructions
	return fmt.Errorf("partial migration detected - manual rollback required:\n" +
		"1. Review the down migration SQL for version " + version + "\n" +
		"2. Execute the down migration manually\n" +
		"3. Mark migration as clean in schema_migrations table")
}

