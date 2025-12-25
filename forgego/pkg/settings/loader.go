package settings

import (
	"fmt"
	"os"
	"strings"

	"github.com/forgego/forge/pkg/validation"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Loader wraps Viper for type-safe configuration loading
type Loader struct {
	v *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	v := viper.New()
	
	// Set default config file name (without extension)
	v.SetConfigName("config")
	
	// Add config paths (current directory, ./config, etc.)
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	
	// Enable environment variables
	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	
	// Set config type priority
	v.SetConfigType("yaml")
	
	// Set defaults from environment variables
	setDefaultsFromEnv(v)
	
	return &Loader{v: v}
}

// LoadFromFile loads configuration from a specific file
func (l *Loader) LoadFromFile(path string) error {
	l.v.SetConfigFile(path)
	return l.v.ReadInConfig()
}

// Load loads configuration into the provided type T
// Supports loading from config files, environment variables, and defaults
func Load[T any]() (*T, error) {
	loader := NewLoader()
	return loadWithLoader[T](loader)
}

// loadWithLoader loads configuration using the provided loader
func loadWithLoader[T any](loader *Loader) (*T, error) {
	var config T
	
	// Try to read config file (optional - won't error if not found)
	if err := loader.v.ReadInConfig(); err != nil {
		// Config file not found is OK, we'll use env vars and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// Unmarshal into the target type
	if err := loader.v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Validate the configuration
	validator := validation.NewStructValidator()
	if err := validator.Validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &config, nil
}

// LoadFromFile loads configuration from a specific file into type T
func LoadFromFile[T any](path string) (*T, error) {
	loader := NewLoader()
	if err := loader.LoadFromFile(path); err != nil {
		return nil, err
	}
	return loadWithLoader[T](loader)
}

// SetDefault sets a default value for a configuration key
func (l *Loader) SetDefault(key string, value interface{}) {
	l.v.SetDefault(key, value)
}

// Set sets a configuration value
func (l *Loader) Set(key string, value interface{}) {
	l.v.Set(key, value)
}

// GetViper returns the underlying Viper instance for advanced usage
func (l *Loader) GetViper() *viper.Viper {
	return l.v
}

// BindEnv binds an environment variable to a configuration key
func (l *Loader) BindEnv(key string, envVar ...string) {
	if len(envVar) > 0 {
		l.v.BindEnv(key, envVar[0])
	} else {
		l.v.BindEnv(key)
	}
}

// SetConfigName sets the name of the config file (without extension)
func (l *Loader) SetConfigName(name string) {
	l.v.SetConfigName(name)
}

// AddConfigPath adds a path to search for config files
func (l *Loader) AddConfigPath(path string) {
	l.v.AddConfigPath(path)
}

// SetConfigType sets the type of the config file (yaml, json, toml, etc.)
func (l *Loader) SetConfigType(configType string) {
	l.v.SetConfigType(configType)
}

// WatchConfig watches for config file changes (useful for hot reloading)
func (l *Loader) WatchConfig(onChange func()) {
	l.v.WatchConfig()
	l.v.OnConfigChange(func(e fsnotify.Event) {
		onChange()
	})
}

// Helper function to set defaults from environment variables
// This allows setting defaults from env vars before unmarshaling
func setDefaultsFromEnv(v *viper.Viper) {
	// Common environment variable patterns
	envVars := map[string]string{
		"DATABASE_URL":      "database.url",
		"DATABASE_DRIVER":    "database.driver",
		"PORT":               "server.port",
		"HOST":               "server.host",
		"SECRET_KEY":         "security.secret_key",
		"DEBUG":              "security.debug",
		"LOG_LEVEL":          "logging.level",
		"LOG_FORMAT":         "logging.format",
		"STATIC_ENABLE":      "static.enable",
		"STATIC_PATH":        "static.path",
		"STATIC_ROOT":        "static.root",
		"I18N_ENABLE":        "i18n.enable",
		"I18N_LOCALES_PATH":  "i18n.locales_path",
		"I18N_DEFAULT_LOCALE": "i18n.default_locale",
		"ADMIN_ENABLE":       "admin.enable",
		"ADMIN_PATH":         "admin.path",
		"ADMIN_TEMPLATE_PATH": "admin.template_path",
		"WORKERS_ENABLE":     "workers.enable",
		"WORKERS_REDIS_ADDR": "workers.redis_addr",
		"WORKERS_REDIS_PASSWORD": "workers.redis_password",
		"WORKERS_REDIS_DB":   "workers.redis_db",
		"WORKERS_CONCURRENCY": "workers.concurrency",
		"API_PATH":           "api.path",
	}
	
	for envVar, configKey := range envVars {
		if value := os.Getenv(envVar); value != "" {
			v.SetDefault(configKey, value)
		}
	}
}

