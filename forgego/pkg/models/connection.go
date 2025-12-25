package models

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
	_ "github.com/uptrace/bun/driver/sqliteshim"
	_ "modernc.org/sqlite"
)

// DB wraps bun.DB with driver and DSN information.
type DB struct {
	*bun.DB
	driver string
	dsn    string
}

// NewDB creates a new database connection.
// Supported drivers: "postgres", "sqlite", "sqlite3"
func NewDB(driver, dsn string) (*DB, error) {
	sqldb, err := openSQLDB(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	var bunDB *bun.DB
	switch driver {
	case "postgres":
		bunDB = bun.NewDB(sqldb, pgdialect.New())
	case "sqlite", "sqlite3":
		bunDB = bun.NewDB(sqldb, sqlitedialect.New())
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	return &DB{
		DB:     bunDB,
		driver: driver,
		dsn:    dsn,
	}, nil
}

func openSQLDB(driver, dsn string) (*sql.DB, error) {
	switch driver {
	case "postgres":
		return sql.Open("pg", dsn)
	case "sqlite", "sqlite3":
		sqldb, err := sql.Open("sqlite", dsn)
		if err != nil {
			return nil, err
		}
		return sqldb, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) Ping() error {
	return db.DB.Ping()
}

