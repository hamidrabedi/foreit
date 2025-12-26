// Package forge provides database functionality
package forge

import (
	pkgDB "github.com/forgego/forge/pkg/db"
)

// DB wraps the pkg DB type
type DB = pkgDB.DB

// NewDBFromConfig creates a new database connection from config
func NewDBFromConfig(cfg *Config) (*DB, error) {
	return pkgDB.NewDBFromConfig(cfg)
}

// NewDB creates a new database connection
func NewDB(dsn string) (*DB, error) {
	return pkgDB.NewDB(dsn)
}
