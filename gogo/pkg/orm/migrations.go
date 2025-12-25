package orm

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// Migrator handles versioned migrations using golang-migrate
type Migrator struct {
	db       *sql.DB
	driver   string
	migrate  *migrate.Migrate
	migrationsPath string
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB, driver, migrationsPath string) (*Migrator, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	m := &Migrator{
		db:             sqlDB,
		driver:         driver,
		migrationsPath: migrationsPath,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return m, nil
}

// init initializes the migrator
func (m *Migrator) init() error {
	var driver migrate.Driver
	var err error

	switch m.driver {
	case "postgres":
		driver, err = postgres.WithInstance(m.db, &postgres.Config{})
		if err != nil {
			return fmt.Errorf("failed to create postgres driver: %w", err)
		}
	case "sqlite3", "sqlite":
		driver, err = sqlite3.WithInstance(m.db, &sqlite3.Config{})
		if err != nil {
			return fmt.Errorf("failed to create sqlite driver: %w", err)
		}
	default:
		return fmt.Errorf("unsupported driver: %s", m.driver)
	}

	absPath, err := filepath.Abs(m.migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	migrateInstance, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		m.driver,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	m.migrate = migrateInstance
	return nil
}

// Up applies all pending migrations
func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	if err := m.migrate.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}
	return nil
}

// Steps applies or rolls back N migrations
func (m *Migrator) Steps(n int) error {
	if err := m.migrate.Steps(n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to step migrations: %w", err)
	}
	return nil
}

// Version returns the current migration version
func (m *Migrator) Version() (uint, bool, error) {
	version, dirty, err := m.migrate.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	return version, dirty, err
}

// Force sets the migration version (use with caution)
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
