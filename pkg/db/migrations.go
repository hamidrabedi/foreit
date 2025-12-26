package db

import (
	"context"
	"fmt"

	"github.com/forgego/forge/pkg/config"
	"github.com/forgego/forge/pkg/migrations/execution"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationRunner handles migration execution using golang-migrate
type MigrationRunner struct {
	db            *DB
	migrate       *migrate.Migrate
	migrationsPath string
}

// NewMigrationRunner creates a new migration runner (package-level function)
func NewMigrationRunner(db *DB, migrationsPath string) (*MigrationRunner, error) {
	// Detect driver from DSN or use config
	cfg := config.NewConfig()
	driverName := cfg.GetDriver()
	
	var driver database.Driver
	var err error
	
	if driverName == "sqlite" || driverName == "sqlite3" {
		driver, err = sqlite3.WithInstance(db.DB, &sqlite3.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create sqlite driver: %w", err)
		}
		driverName = "sqlite3" // golang-migrate uses "sqlite3" as the driver name
	} else {
		driver, err = postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres driver: %w", err)
		}
		driverName = "postgres"
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		driverName,
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &MigrationRunner{
		db:            db,
		migrate:       m,
		migrationsPath: migrationsPath,
	}, nil
}

// Migrate applies all pending migrations
func (mr *MigrationRunner) Migrate(ctx context.Context) error {
	if err := mr.migrate.Up(); err != nil {
		if err == migrate.ErrNoChange {
			// No pending migrations - this is fine
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

// MigrateTo applies migrations up to a specific version
func (mr *MigrationRunner) MigrateTo(ctx context.Context, version uint) error {
	if err := mr.migrate.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}
	return nil
}

// Rollback rolls back the last migration
func (mr *MigrationRunner) Rollback(ctx context.Context) error {
	if err := mr.migrate.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}
	return nil
}

// RollbackTo rolls back to a specific version
func (mr *MigrationRunner) RollbackTo(ctx context.Context, version uint) error {
	if err := mr.migrate.Migrate(version); err != nil {
		return fmt.Errorf("failed to rollback to version %d: %w", version, err)
	}
	return nil
}

// Version returns the current migration version
func (mr *MigrationRunner) Version(ctx context.Context) (version uint, dirty bool, err error) {
	version, dirty, err = mr.migrate.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}
	return version, dirty, nil
}

// Status returns migration status information
func (mr *MigrationRunner) Status(ctx context.Context) (*MigrationStatus, error) {
	version, dirty, err := mr.Version(ctx)
	if err != nil {
		return nil, err
	}

	return &MigrationStatus{
		Version: version,
		Dirty:   dirty,
	}, nil
}

// GetDetailedStatus returns detailed migration status
func (mr *MigrationRunner) GetDetailedStatus(ctx context.Context) (*DetailedMigrationStatus, error) {
	// Use StatusReporter for detailed status
	reporter := execution.NewStatusReporter(mr.migrationsPath, mr.migrate)
	detailed, err := reporter.GetDetailedStatus(ctx)
	if err != nil {
		// Fallback to basic status if detailed status fails
		status, statusErr := mr.Status(ctx)
		if statusErr != nil {
			return nil, fmt.Errorf("failed to get status: %w", statusErr)
		}
		return &DetailedMigrationStatus{
			Current: fmt.Sprintf("%d", status.Version),
			Status:  "OK",
			Dirty:   status.Dirty,
		}, nil
	}

	// Convert execution.DetailedStatus to db.DetailedMigrationStatus
	result := &DetailedMigrationStatus{
		Current:    detailed.Current,
		Next:       detailed.Next,
		Status:     detailed.Status,
		Dirty:      detailed.Status == "DIRTY",
		Error:      detailed.Error,
		Applied:    make([]string, len(detailed.Applied)),
		Pending:    make([]string, len(detailed.Pending)),
		OutOfOrder: make([]string, len(detailed.OutOfOrder)),
	}

	for i, mig := range detailed.Applied {
		result.Applied[i] = fmt.Sprintf("[%s] %s", mig.Version, mig.Name)
	}
	for i, mig := range detailed.Pending {
		result.Pending[i] = fmt.Sprintf("[%s] %s", mig.Version, mig.Name)
	}
	for i, mig := range detailed.OutOfOrder {
		result.OutOfOrder[i] = fmt.Sprintf("[%s] %s", mig.Version, mig.Name)
	}

	return result, nil
}

// MigrationStatus represents the status of migrations
type MigrationStatus struct {
	Version uint
	Dirty   bool
}

// DetailedMigrationStatus represents detailed migration status
type DetailedMigrationStatus struct {
	Current    string
	Next       string
	Applied    []string
	Pending    []string
	OutOfOrder []string
	Status     string
	Dirty      bool
	Error      string
}

// Force sets a migration version and marks it as clean (for dirty state recovery)
// WARNING: Use with caution - only after manually fixing a failed migration
func (mr *MigrationRunner) Force(ctx context.Context, version uint) error {
	if err := mr.migrate.Force(int(version)); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}
	return nil
}

// Close closes the migration runner
func (mr *MigrationRunner) Close() error {
	sourceErr, dbErr := mr.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close migration source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close migration database: %w", dbErr)
	}
	return nil
}
