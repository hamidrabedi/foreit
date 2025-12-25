package orm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Client provides database connection and migration support
type Client struct {
	db       *gorm.DB
	driver   string
	dsn      string
}

// NewClient creates a new database client
func NewClient(driver, dsn string) (*Client, error) {
	var db *gorm.DB
	var err error

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case "sqlite", "sqlite3":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Client{
		db:     db,
		driver: driver,
		dsn:    dsn,
	}, nil
}

// DB returns the underlying database connection
func (c *Client) DB() *gorm.DB {
	return c.db
}

// Migrator returns a migrator for versioned migrations
func (c *Client) Migrator(migrationsPath string) (*Migrator, error) {
	return NewMigrator(c.db, c.driver, migrationsPath)
}

// Close closes the database connection
func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping checks the database connection
func (c *Client) Ping() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
