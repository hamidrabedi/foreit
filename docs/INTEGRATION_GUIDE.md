# Logging & Error System Integration Guide

## Quick Start

### 1. Basic Server Setup with New Systems

```go
package main

import (
	"github.com/forgego/forge/config"
	"github.com/forgego/forge/log"
	"github.com/forgego/forge/server"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()
	settings := config.LoadSettings(cfg)

	// Create logger
	logger, err := log.NewLogger(settings.App.Debug)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Create server (automatically includes error handler and request ID)
	server, err := server.NewServer(cfg, settings, logger)
	if err != nil {
		logger.Fatal("Failed to create server", log.Error(err))
	}

	// Register routes
	server.RegisterRoutes(func(router *server.Router) {
		router.Get("/api/users", handleUsers)
	})

	// Start server
	if err := server.Start(); err != nil {
		logger.Fatal("Server failed", log.Error(err))
	}
}
```

### 2. Using Error System

```go
import "github.com/forgego/forge/api/errors"

func handleUsers(w http.ResponseWriter, r *http.Request) {
	// Create validation error
	valErr := errors.NewValidationError("Validation failed")
	valErr.AddFieldError("email", "Invalid email format", errors.CodeInvalidEmailFormatField)
	
	// Return error (will be handled by middleware)
	// The error handler will automatically convert to RFC 7807 Problem Details
	return valErr
}

// Or use legacy exceptions (automatically mapped)
import "github.com/forgego/forge/api/exceptions"

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	err := exceptions.NewNotFound("User not found")
	// Error handler will convert to Problem Details automatically
	return err
}
```

### 3. Advanced Logger Configuration

```go
// Using builder pattern
logger := log.NewBuilder().
	Production().
	Level(log.LevelInfo).
	AddConsoleOutput(log.LevelInfo).
	AddFileOutput("logs/app.log", log.LevelInfo, 100, 30, 10, true).
	Caller(true).
	Sampling(100, 100).
	Build()

// Or from config file
config := log.DefaultLoggingConfig(false)
// Modify config as needed
logger, err := log.NewLoggerFromConfig(config)
```

### 4. Custom Error Handler Configuration

```go
import "github.com/forgego/forge/api/errors"

// Create custom error handler
handlerConfig := errors.DefaultHandlerConfig()
handlerConfig.Logger = logger.Logger
handlerConfig.TypeBaseURL = "https://api.myapp.com/problems"
handlerConfig.LinkHeaderURL = "https://api.myapp.com/docs/errors"
handlerConfig.IncludeLinkHeader = true

handler := errors.NewHandler(handlerConfig)

// Use in router
router.Use(handler.Middleware())
```

## Middleware Order

The recommended middleware order is:

1. **RequestIDMiddleware** - Generate request IDs (first)
2. **ErrorHandler** - Handle errors and panics (early)
3. **Logging Middleware** - Log requests
4. Other middleware (CORS, auth, etc.)

This is automatically set up in `server.NewServer()`.

## Error Response Format

All errors are automatically converted to RFC 7807 Problem Details:

```json
{
  "type": "https://api.example.com/problems/validation-error",
  "title": "Validation Failed",
  "status": 400,
  "detail": "The request contains invalid data",
  "instance": "/api/v1/users",
  "code": "VALIDATION_ERROR",
  "meta": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  "errors": {
    "email": [
      {
        "message": "Invalid email format",
        "code": "INVALID_EMAIL_FORMAT"
      }
    ]
  }
}
```

## Log Output Examples

### Development Mode (Colored)
```
15:04:05.123 INFO  server.go:45 Starting server | address=:8080 environment=development
15:04:05.456 INFO  middleware.go:23 HTTP request | method=GET path=/api/users status=200 duration=45ms request_id=abc-123
```

### Production Mode (JSON)
```json
{"level":"info","ts":1705314245.123,"caller":"server.go:45","msg":"Starting server","address":":8080","environment":"production"}
{"level":"info","ts":1705314245.456,"caller":"middleware.go:23","msg":"HTTP request","method":"GET","path":"/api/users","status":200,"duration":45000000,"request_id":"abc-123"}
```

## Configuration File

Add to your `config.yaml`:

```yaml
logging:
  level: "info"
  format: "json"
  outputs:
    - type: "console"
      enabled: true
      level: "debug"
    - type: "file"
      enabled: true
      path: "logs/app.log"
      level: "info"
      rotation:
        max_size: 100
        max_age: 30
        max_backups: 10
        compress: true

errors:
  problem_details:
    type_base_url: "https://api.example.com/problems"
  request_id:
    header_name: "X-Request-ID"
    generate_if_missing: true
  sanitization:
    hide_database_errors: true
    hide_stack_traces: true
    redact_pii: true
```

## Migration from Old System

### Old Way
```go
logger, err := log.NewLogger(true)
// Errors returned as simple JSON
```

### New Way
```go
logger, err := log.NewLogger(true) // Still works!
// Or use builder:
logger := log.NewBuilder().Development().Build()

// Errors automatically use RFC 7807 format
// Request IDs automatically added
// Sanitization automatically enabled
```

The new system is backward compatible - existing code continues to work, but you get all the new features automatically!
