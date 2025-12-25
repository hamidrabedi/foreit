package app

import (
	"github.com/forgego/forge/pkg/admin"
	"github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/i18n"
	"github.com/forgego/forge/pkg/middleware"
	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/rest"
	"github.com/forgego/forge/pkg/workers"
)

// AppSettings gathers all framework module configurations
// Users can embed this struct in their own Settings to extend it
type AppSettings struct {
	// Core framework configs
	Database models.DatabaseConfig `mapstructure:"database"`
	Server   rest.ServerConfig      `mapstructure:"server"`
	Security middleware.SecurityConfig `mapstructure:"security"`
	Logging  middleware.LoggingConfig  `mapstructure:"logging"`
	Static   rest.StaticConfig      `mapstructure:"static"`
	I18n     i18n.I18nConfig        `mapstructure:"i18n"`
	
	// Module configs
	Admin   admin.AdminConfig    `mapstructure:"admin"`
	Workers workers.WorkersConfig `mapstructure:"workers"`
	API     api.APIConfig        `mapstructure:"api"`
}

// DefaultAppSettings returns app settings with default values
func DefaultAppSettings() *AppSettings {
	return &AppSettings{
		Database: models.DatabaseConfig{
			Driver: "postgres",
		},
		Server: rest.ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Security: middleware.SecurityConfig{
			Debug: false,
		},
		Logging: middleware.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Static: rest.StaticConfig{
			Enable: true,
			Path:   "/static",
			Root:   "./public",
		},
		I18n: i18n.I18nConfig{
			Enable:       false,
			LocalesPath:  "./locales",
			DefaultLocale: "en",
		},
		Admin: admin.AdminConfig{
			Enable:       true,
			Path:         "/admin",
			TemplatePath: "./templates/admin",
		},
		Workers: workers.WorkersConfig{
			Enable:      true,
			RedisAddr:   "localhost:6379",
			RedisPassword: "",
			RedisDB:     0,
			Concurrency: 10,
		},
		API: api.APIConfig{
			Path: "/api/v1",
		},
	}
}

