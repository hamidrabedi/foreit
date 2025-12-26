package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/forgego/forge/pkg/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DB wraps *sql.DB with additional functionality
type DB struct {
	*sql.DB
	dsn string
}

// Config holds database configuration
type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns a default database configuration
func DefaultConfig() *Config {
	return &Config{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// NewDB creates a new database connection with PostgreSQL
func NewDB(dsn string) (*DB, error) {
	return NewDBWithConfig(dsn, DefaultConfig())
}

// NewDBFromConfig creates a new database connection from config
func NewDBFromConfig(cfg *config.Config) (*DB, error) {
	driver := cfg.GetDriver()
	dsn := cfg.GetDSN()
	dbConfig := &Config{
		DSN:             dsn,
		MaxOpenConns:    cfg.GetInt("database.max_open_conns", 25),
		MaxIdleConns:    cfg.GetInt("database.max_idle_conns", 5),
		ConnMaxLifetime: time.Duration(cfg.GetInt("database.conn_max_lifetime", 300)) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.GetInt("database.conn_max_idle_time", 600)) * time.Second,
	}
	return NewDBWithDriver(driver, dsn, dbConfig)
}

// NewDBWithConfig creates a new database connection with custom config (PostgreSQL)
func NewDBWithConfig(dsn string, config *Config) (*DB, error) {
	return NewDBWithDriver("postgres", dsn, config)
}

// NewDBWithDriver creates a new database connection with specified driver
func NewDBWithDriver(driver, dsn string, config *Config) (*DB, error) {
	// Normalize driver name for sql.Open (it uses "sqlite3")
	driverName := driver
	if driver == "sqlite" {
		driverName = "sqlite3"
	}
	
	sqldb, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool (SQLite has different defaults)
	if config != nil {
		if driver == "sqlite" || driver == "sqlite3" {
			// SQLite has different connection pool settings
			sqldb.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes well
			sqldb.SetMaxIdleConns(1)
		} else {
			sqldb.SetMaxOpenConns(config.MaxOpenConns)
			sqldb.SetMaxIdleConns(config.MaxIdleConns)
			sqldb.SetConnMaxLifetime(config.ConnMaxLifetime)
			sqldb.SetConnMaxIdleTime(config.ConnMaxIdleTime)
		}
	}

	// Test connection
	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		DB:  sqldb,
		dsn: dsn,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Ping checks the database connection
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// HealthCheck performs a health check on the database
func (db *DB) HealthCheck() error {
	return db.Ping()
}
