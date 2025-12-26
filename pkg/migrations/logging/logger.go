package logging

import (
	"log"
	"os"
)

// Logger provides structured logging for migrations
type Logger struct {
	*log.Logger
	verbose bool
}

// NewLogger creates a new migration logger
func NewLogger(verbose bool) *Logger {
	return &Logger{
		Logger:  log.New(os.Stdout, "[migrations] ", log.LstdFlags),
		verbose: verbose,
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.verbose {
		l.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.Printf("[WARN] "+format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.Printf("[ERROR] "+format, v...)
}

// Debug logs a debug message (only if verbose)
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.verbose {
		l.Printf("[DEBUG] "+format, v...)
	}
}

// MigrationStart logs the start of a migration
func (l *Logger) MigrationStart(version, name string) {
	l.Info("Starting migration: %s (%s)", version, name)
}

// MigrationComplete logs the completion of a migration
func (l *Logger) MigrationComplete(version, name string) {
	l.Info("Completed migration: %s (%s)", version, name)
}

// MigrationFailed logs a failed migration
func (l *Logger) MigrationFailed(version, name string, err error) {
	l.Error("Failed migration: %s (%s): %v", version, name, err)
}

// ChangeDetected logs a detected change
func (l *Logger) ChangeDetected(changeType, table string) {
	l.Debug("Detected change: %s on table %s", changeType, table)
}

