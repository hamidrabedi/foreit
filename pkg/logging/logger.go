package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with framework-specific methods
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger with framework defaults
func NewLogger(development bool) (*Logger, error) {
	var config zap.Config
	if development {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

// NewNopLogger creates a no-op logger for testing
func NewNopLogger() *Logger {
	return &Logger{Logger: zap.NewNop()}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// WithFields adds structured fields to the logger
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// String creates a zap.String field
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int creates a zap.Int field
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Error creates a zap.Error field
func Error(err error) zap.Field {
	return zap.Error(err)
}
