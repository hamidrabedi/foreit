// Package forge provides logging functionality
package forge

import (
	pkgLogging "github.com/forgego/forge/pkg/logging"
	"go.uber.org/zap"
)

// Logger wraps the pkg logger
type Logger = pkgLogging.Logger

// NewLogger creates a new logger with framework defaults
func NewLogger(development bool) (*Logger, error) {
	return pkgLogging.NewLogger(development)
}

// NewNopLogger creates a no-op logger for testing
func NewNopLogger() *Logger {
	return pkgLogging.NewNopLogger()
}

// String creates a zap.String field
func String(key, value string) zap.Field {
	return pkgLogging.String(key, value)
}

// Int creates a zap.Int field
func Int(key string, value int) zap.Field {
	return pkgLogging.Int(key, value)
}

// Error creates a zap.Error field
func Error(err error) zap.Field {
	return pkgLogging.Error(err)
}
