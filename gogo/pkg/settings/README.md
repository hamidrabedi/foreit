# Settings Module

Application configuration management with environment variable loading, validation, and defaults.

## Usage

### Basic Usage

```go
type AppSettings struct {
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    Port        int    `env:"PORT" default:"8080"`
    Debug       bool   `env:"DEBUG" default:"false"`
}

settings := settings.Load[AppSettings]()
// Access: settings.DatabaseURL, settings.Port, etc.
```

### With Custom Sources

```go
loader := settings.NewLoader(
    &settings.EnvSource{},
    &settings.FileSource{Path: "config.yaml"},
)

settings, err := loader.Load[AppSettings]()
```

### Global Registry

```go
// Register a setting
settings.Register("app.name", "My App")

// Retrieve
name := settings.GetString("app.name", "Default Name")
```

## Features

- Environment variable loading
- Default values
- Required field validation
- Type conversion (string, int, bool, float, slices)
- Multiple configuration sources
- Global registry for runtime settings

