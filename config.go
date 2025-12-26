// Package forge provides the public API for the Forge framework.
// This package re-exports commonly used types and functions from pkg packages.
package forge

// Re-export config types and functions
import (
	pkgConfig "github.com/forgego/forge/pkg/config"
)

// Config wraps the pkg config with a public API
type Config = pkgConfig.Config

// Settings wraps the pkg settings with a public API
type Settings = pkgConfig.Settings

// AppSettings wraps pkg AppSettings
type AppSettings = pkgConfig.AppSettings

// ServerSettings wraps pkg ServerSettings
type ServerSettings = pkgConfig.ServerSettings

// DatabaseSettings wraps pkg DatabaseSettings
type DatabaseSettings = pkgConfig.DatabaseSettings

// SecuritySettings wraps pkg SecuritySettings
type SecuritySettings = pkgConfig.SecuritySettings

// AdminSettings wraps pkg AdminSettings
type AdminSettings = pkgConfig.AdminSettings

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	return pkgConfig.NewConfig()
}

// LoadSettings loads settings from config
func LoadSettings(cfg *Config) *Settings {
	return pkgConfig.LoadSettings(cfg)
}
