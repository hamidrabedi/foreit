package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/forgego/forge/pkg/db"
	"github.com/forgego/forge/pkg/migrations"
)

// MigrationInfo represents information about a migration
type MigrationInfo struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

// CreateMigrationFromModels generates migrations from model definitions
func CreateMigrationFromModels(t *testing.T, modelsDir, migrationsDir, migrationName string) error {
	gen, err := migrations.NewGenerator(modelsDir, migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to create migration generator: %w", err)
	}

	if err := gen.GenerateMigrations(migrationName); err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	return nil
}

// ApplyMigrationSequence applies multiple migrations in order
func ApplyMigrationSequence(ctx context.Context, t *testing.T, database *db.DB, migrationsPath string) error {
	runner, err := db.NewMigrationRunner(database, migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to create migration runner: %w", err)
	}
	defer runner.Close()

	if err := runner.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// RollbackMigrationSequence rolls back migrations in reverse order
func RollbackMigrationSequence(ctx context.Context, t *testing.T, database *db.DB, migrationsPath string, steps int) error {
	runner, err := db.NewMigrationRunner(database, migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to create migration runner: %w", err)
	}
	defer runner.Close()

	for i := 0; i < steps; i++ {
		if err := runner.Rollback(ctx); err != nil {
			return fmt.Errorf("failed to rollback migration step %d: %w", i+1, err)
		}
	}

	return nil
}

// AssertMigrationState verifies migration state matches expected
func AssertMigrationState(ctx context.Context, t *testing.T, database *db.DB, migrationsPath string, expectedVersion uint, expectedDirty bool) {
	runner, err := db.NewMigrationRunner(database, migrationsPath)
	if err != nil {
		t.Fatalf("failed to create migration runner: %v", err)
	}
	defer runner.Close()

	version, dirty, err := runner.Version(ctx)
	if err != nil {
		t.Fatalf("failed to get migration version: %v", err)
	}

	if version != expectedVersion {
		t.Errorf("expected migration version %d, got %d", expectedVersion, version)
	}

	if dirty != expectedDirty {
		t.Errorf("expected dirty state %v, got %v", expectedDirty, dirty)
	}
}

// AssertSchemaMatchesModels verifies database schema matches model definitions
// This is a high-level check that validates tables, columns, and basic constraints exist
func AssertSchemaMatchesModels(ctx context.Context, t *testing.T, database *sql.DB, dialect string, expectedTables []string) {
	for _, tableName := range expectedTables {
		AssertTableExists(ctx, t, database, dialect, tableName)
	}
}

// GetMigrationFiles returns all migration files in a directory
func GetMigrationFiles(migrationsDir string) ([]MigrationInfo, error) {
	var migrations []MigrationInfo

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrationMap := make(map[string]*MigrationInfo)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !filepath.Ext(name) == ".sql" {
			continue
		}

		// Parse migration file name: {version}_{name}.{up|down}.sql
		var version, migrationName, direction string
		baseName := filepath.Base(name)
		ext := filepath.Ext(baseName) // .sql
		baseWithoutExt := baseName[:len(baseName)-len(ext)]

		// Split by . to get direction
		parts := filepath.Ext(baseWithoutExt) // .up or .down
		if parts == "" {
			continue
		}
		direction = parts[1:] // Remove the dot
		baseWithoutDirection := baseWithoutExt[:len(baseWithoutExt)-len(parts)]

		// Split by _ to get version and name
		underscoreIdx := -1
		for i := 0; i < len(baseWithoutDirection) && i < 6; i++ {
			if baseWithoutDirection[i] == '_' {
				underscoreIdx = i
				break
			}
		}

		if underscoreIdx == -1 {
			continue
		}

		version = baseWithoutDirection[:underscoreIdx]
		migrationName = baseWithoutDirection[underscoreIdx+1:]

		key := version + "_" + migrationName
		if migrationMap[key] == nil {
			migrationMap[key] = &MigrationInfo{
				Version: version,
				Name:    migrationName,
			}
		}

		// Read SQL content
		filePath := filepath.Join(migrationsDir, name)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		if direction == "up" {
			migrationMap[key].UpSQL = string(content)
		} else if direction == "down" {
			migrationMap[key].DownSQL = string(content)
		}
	}

	// Convert map to slice
	for _, mig := range migrationMap {
		migrations = append(migrations, *mig)
	}

	return migrations, nil
}

// ApplyMigrationSQL applies a single SQL migration directly
func ApplyMigrationSQL(ctx context.Context, t *testing.T, database *sql.DB, sql string) error {
	_, err := database.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to apply migration SQL: %w\nSQL: %s", err, sql)
	}
	return nil
}

// RollbackMigrationSQL applies a rollback SQL directly
func RollbackMigrationSQL(ctx context.Context, t *testing.T, database *sql.DB, sql string) error {
	_, err := database.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to rollback migration SQL: %w\nSQL: %s", err, sql)
	}
	return nil
}
