package models

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	URL    string `mapstructure:"url" validate:"required"`
	Driver string `mapstructure:"driver" default:"postgres"` // "postgres", "sqlite3", "mysql"
}

