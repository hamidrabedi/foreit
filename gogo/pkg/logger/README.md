# Logger Module

Structured logging using [zerolog](https://github.com/rs/zerolog), a fast and structured logger.

## Features

- Structured logging (JSON in production, console in development)
- Multiple log levels
- Context fields
- Error tracking
- Fast performance

## Usage

### Basic Usage

```go
import "github.com/gogo/pkg/logger"

// Initialize
logger.Init(false) // false = development mode

// Use global logger
logger.Info("Server starting")
logger.Error("Failed to connect")
logger.WithField("user_id", 123).Info("User logged in")
```

### Instance Logger

```go
log := logger.New()
log.Info("Message")
log.WithField("key", "value").Debug("Debug message")
log.WithError(err).Error("Error occurred")
```

### Production Logger

```go
log := logger.NewProduction()
log.Info("Production log") // JSON output
```

### Log Levels

```go
log.Debug("Debug message")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message")
log.Fatal("Fatal message") // Exits program
```

### With Fields

```go
log.WithFields(map[string]interface{}{
    "user_id": 123,
    "action": "login",
}).Info("User action")
```

