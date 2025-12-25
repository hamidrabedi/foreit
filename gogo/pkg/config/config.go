package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Loader loads configuration from various sources
type Loader struct {
	viper *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	v := viper.New()
	
	// Set defaults
	v.SetEnvPrefix("GOGO")
	v.AutomaticEnv()
	
	return &Loader{
		viper: v,
	}
}

// LoadFromFile loads configuration from a file
func (l *Loader) LoadFromFile(path string) error {
	l.viper.SetConfigFile(path)
	return l.viper.ReadInConfig()
}

// LoadFromEnv loads configuration from .env file
func (l *Loader) LoadFromEnv(path string) error {
	if path == "" {
		path = ".env"
	}
	return godotenv.Load(path)
}

// LoadFromEnvFile loads from .env file and sets environment variables
func (l *Loader) LoadFromEnvFile(path string) error {
	if path == "" {
		path = ".env"
	}
	
	envMap, err := godotenv.Read(path)
	if err != nil {
		return err
	}
	
	for key, value := range envMap {
		os.Setenv(key, value)
	}
	
	return nil
}

// Get gets a configuration value
func (l *Loader) Get(key string) interface{} {
	return l.viper.Get(key)
}

// GetString gets a string value
func (l *Loader) GetString(key string) string {
	return l.viper.GetString(key)
}

// GetInt gets an int value
func (l *Loader) GetInt(key string) int {
	return l.viper.GetInt(key)
}

// GetBool gets a bool value
func (l *Loader) GetBool(key string) bool {
	return l.viper.GetBool(key)
}

// Set sets a configuration value
func (l *Loader) Set(key string, value interface{}) {
	l.viper.Set(key, value)
}

// SetDefault sets a default value
func (l *Loader) SetDefault(key string, value interface{}) {
	l.viper.SetDefault(key, value)
}

// Load loads configuration with automatic detection
func Load() (*Loader, error) {
	loader := NewLoader()
	
	// Try to load .env file
	if err := loader.LoadFromEnvFile(".env"); err != nil {
		// .env file is optional
	}
	
	// Try to load config file
	if err := loader.LoadFromFile("config.yaml"); err != nil {
		// Try config.json
		if err := loader.LoadFromFile("config.json"); err != nil {
			// Config file is optional
		}
	}
	
	return loader, nil
}

// LoadInto loads configuration into a struct
func LoadInto[T any](prefix string) (*T, error) {
	loader, err := Load()
	if err != nil {
		return nil, err
	}
	
	var config T
	if err := loader.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}

