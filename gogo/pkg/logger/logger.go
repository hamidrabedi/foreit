package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog logger
type Logger struct {
	logger zerolog.Logger
}

// New creates a new logger
func New() *Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	
	// Use console writer for development
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	
	logger := zerolog.New(output).With().Timestamp().Logger()
	
	return &Logger{
		logger: logger,
	}
}

// NewProduction creates a production logger (JSON output)
func NewProduction() *Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &Logger{logger: logger}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}

// WithFields adds multiple fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	logger := l.logger.With()
	for k, v := range fields {
		logger = logger.Interface(k, v)
	}
	return &Logger{logger: logger.Logger()}
}

// WithError adds an error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With().Err(err).Logger(),
	}
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level zerolog.Level) {
	l.logger = l.logger.Level(level)
}

// Global logger instance
var globalLogger *Logger

// Init initializes the global logger
func Init(production bool) {
	if production {
		globalLogger = NewProduction()
	} else {
		globalLogger = New()
	}
	log.Logger = globalLogger.logger
}

// Get returns the global logger
func Get() *Logger {
	if globalLogger == nil {
		globalLogger = New()
	}
	return globalLogger
}

// Debug logs using global logger
func Debug(msg string) {
	Get().Debug(msg)
}

// Info logs using global logger
func Info(msg string) {
	Get().Info(msg)
}

// Warn logs using global logger
func Warn(msg string) {
	Get().Warn(msg)
}

// Error logs using global logger
func Error(msg string) {
	Get().Error(msg)
}

// Fatal logs using global logger
func Fatal(msg string) {
	Get().Fatal(msg)
}

