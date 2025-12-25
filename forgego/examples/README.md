# Examples

This directory contains example projects demonstrating how to use the Forge framework's configuration system.

## Architecture Overview

The Forge framework uses a **decentralized configuration approach**:

- **Each package owns its config**: Every framework package (admin, workers, api, models, etc.) defines its own config struct
- **App package aggregates**: The `app` package gathers all framework configs into `AppSettings`
- **Settings package provides infrastructure**: The `settings` package only provides the loader (Viper-based)
- **Users extend freely**: Users can embed `AppSettings` and add their own configs or third-party package configs

## Package Config Locations

| Package | Config File | Config Type |
|---------|------------|-------------|
| `models` | `pkg/models/config.go` | `models.DatabaseConfig` |
| `rest` | `pkg/rest/config.go` | `rest.ServerConfig`, `rest.StaticConfig` |
| `middleware` | `pkg/middleware/config.go` | `middleware.SecurityConfig`, `middleware.LoggingConfig` |
| `i18n` | `pkg/i18n/config.go` | `i18n.I18nConfig` |
| `admin` | `pkg/admin/config.go` | `admin.AdminConfig` |
| `workers` | `pkg/workers/config.go` | `workers.WorkersConfig` |
| `api` | `pkg/api/config.go` | `api.APIConfig` |
| `app` | `pkg/app/config.go` | `app.AppSettings` (aggregates all above) |

## Examples

### 1. Basic Example (`examples/basic/`)

Shows how to:
- Create a custom Settings struct that embeds `app.AppSettings`
- Add your own application-specific configuration
- Load settings from config files or environment variables
- Use the settings to create and run the application

**Key Files:**
- `main.go` - Application with custom settings
- `config.yaml` - Configuration file

### 2. Third-Party Package Integration (`examples/third-party-package/`)

Shows how to:
- Create a third-party package with its own config
- Integrate the package config into user's Settings
- Load and use the package with its config

**Key Files:**
- `main.go` - Application using third-party package
- `mypackage/config.go` - Package config definition
- `mypackage/package.go` - Package implementation
- `config.yaml` - Configuration including package settings

## Benefits of This Architecture

1. **Package Independence**: Each package manages its own config, no central dependency
2. **Easy Integration**: Third-party packages can easily integrate by defining their config
3. **Type Safety**: All configs are type-safe with compile-time checking
4. **Flexibility**: Users compose only the configs they need
5. **Scalability**: As the framework grows, new packages add their configs without touching existing code

## Usage Pattern

```go
// 1. Define your Settings (embed AppSettings + add your own)
type Settings struct {
    app.AppSettings
    MyApp MyAppConfig `mapstructure:"myapp"`
}

// 2. Load settings
settings, err := settings.Load[Settings]()

// 3. Use framework settings for app
application, err := app.New(&settings.AppSettings)

// 4. Use your custom settings for your logic
apiKey := settings.MyApp.APIKey
```

## Configuration Sources

Settings can be loaded from:
1. **Config files**: YAML, JSON, TOML (default: `config.yaml`)
2. **Environment variables**: Automatic mapping (e.g., `DATABASE_URL`, `SERVER_PORT`)
3. **Defaults**: Struct tag defaults or programmatic defaults

Environment variables take highest priority, then config files, then defaults.

