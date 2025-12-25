package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrator handles database migrations
type Migrator struct {
	migrate *migrate.Migrate
}

// NewMigrator creates a new migrator
func NewMigrator(db *sql.DB, driver string, migrationsPath string) (*Migrator, error) {
	var driverInstance migrate.Driver
	var err error
	
	switch driver {
	case "postgres":
		driverInstance, err = postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres driver: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
	
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		driver,
		driverInstance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}
	
	return &Migrator{
		migrate: m,
	}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Steps runs n migrations (positive for up, negative for down)
func (m *Migrator) Steps(n int) error {
	if err := m.migrate.Steps(n); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Version returns the current migration version
func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}

// Force sets the migration version without running migrations
func (m *Migrator) Force(version int) error {
	return m.migrate.Force(version)
}

// Close closes the migrator
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}

