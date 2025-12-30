// Package db provides database functionality
package db

import (
	"database/sql"
	"fmt"

	"github.com/forgego/forge/config"
)

// DB wraps database/sql.DB with additional functionality
type DB struct {
	*sql.DB
}

// NewDBFromConfig creates a new database connection from config
func NewDBFromConfig(cfg *config.Config) (*DB, error) {
	driver := cfg.GetDriver()
	host := cfg.GetString("database.host", "localhost")
	port := cfg.GetInt("database.port", 5432)
	user := cfg.GetString("database.user", "postgres")
	password := cfg.GetString("database.password", "")
	name := cfg.GetString("database.name", "forge")
	sslmode := cfg.GetString("database.sslmode", "disable")

	var dsn string
	if driver == "postgres" || driver == "postgresql" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, name, sslmode)
	} else if driver == "sqlite" || driver == "sqlite3" {
		dsn = name // For SQLite, name is the file path
	} else {
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	return NewDB(dsn)
}

// NewDB creates a new database connection
// For PostgreSQL, dsn should be a connection string
// For SQLite, dsn should be a file path
func NewDB(dsn string) (*DB, error) {
	// Try PostgreSQL first
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		// Try SQLite as fallback
		sqlDB, err = sql.Open("sqlite3", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: sqlDB}, nil
}
