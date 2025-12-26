package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config wraps viper with framework-specific methods
type Config struct {
	*viper.Viper
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	v := viper.New()

	// Set defaults
	v.SetDefault("app.name", "forge")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.debug", true)
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", "8000")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)

	// Set config file
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Environment variables
	v.SetEnvPrefix("FORGE")
	v.AutomaticEnv()

	// Read config file (optional - don't error if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file found but another error occurred
			fmt.Printf("Error reading config file: %v\n", err)
		}
	}

	return &Config{Viper: v}
}

// GetString gets a string value with a default
func (c *Config) GetString(key, defaultValue string) string {
	if c.Viper.IsSet(key) {
		return c.Viper.GetString(key)
	}
	return defaultValue
}

// GetInt gets an int value with a default
func (c *Config) GetInt(key string, defaultValue int) int {
	if c.Viper.IsSet(key) {
		return c.Viper.GetInt(key)
	}
	return defaultValue
}

// GetBool gets a bool value with a default
func (c *Config) GetBool(key string, defaultValue bool) bool {
	if c.Viper.IsSet(key) {
		return c.Viper.GetBool(key)
	}
	return defaultValue
}

// IsDevelopment returns true if app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.GetString("app.env", "development") == "development"
}

// IsProduction returns true if app is in production mode
func (c *Config) IsProduction() bool {
	return c.GetString("app.env", "development") == "production"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	// Try environment variable first
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	driver := c.GetString("database.driver", "postgres")
	
	// SQLite uses a file path, not connection string
	if driver == "sqlite" || driver == "sqlite3" {
		dbname := c.GetString("database.name", "db.sqlite3")
		return dbname
	}

	// Build PostgreSQL connection string
	host := c.GetString("database.host", "localhost")
	port := c.GetInt("database.port", 5432)
	user := c.GetString("database.user", "postgres")
	password := c.GetString("database.password", "")
	dbname := c.GetString("database.name", "forge")
	sslmode := c.GetString("database.sslmode", "disable")

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

// GetDriver returns the database driver name
func (c *Config) GetDriver() string {
	return c.GetString("database.driver", "postgres")
}
