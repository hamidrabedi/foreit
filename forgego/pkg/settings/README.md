# Settings Module

Type-safe configuration infrastructure for the Forge framework using Viper. This package provides the **loader infrastructure only** - each package in the framework holds its own config struct.

## Philosophy

- **Package ownership**: Each package (admin, workers, api, models, etc.) defines its own config struct
- **App aggregation**: The `app` package gathers all configs into `AppSettings`
- **User extensibility**: Users can embed `AppSettings` and add their own configs
- **Type safety**: Compile-time checking for all configuration
- **Multi-source**: Supports YAML, JSON, TOML, environment variables, and more

## Package Configs

Each framework package defines its own config:

- `models.DatabaseConfig` - Database connection settings
- `rest.ServerConfig` - HTTP server settings
- `rest.StaticConfig` - Static file serving settings
- `middleware.SecurityConfig` - Security settings
- `middleware.LoggingConfig` - Logging settings
- `i18n.I18nConfig` - Internationalization settings
- `admin.AdminConfig` - Admin module settings
- `workers.WorkersConfig` - Workers module settings
- `api.APIConfig` - API module settings

## AppSettings

The `app` package gathers all framework configs:

```go
// pkg/app/config.go
type AppSettings struct {
    Database models.DatabaseConfig
    Server   rest.ServerConfig
    Security middleware.SecurityConfig
    Logging  middleware.LoggingConfig
    Static   rest.StaticConfig
    I18n     i18n.I18nConfig
    Admin    admin.AdminConfig
    Workers  workers.WorkersConfig
    API      api.APIConfig
}
```

## Usage

### Loading AppSettings

```go
import (
    "github.com/forgego/forge/pkg/app"
    "github.com/forgego/forge/pkg/settings"
)

// Load app settings from config files or environment variables
appSettings, err := settings.Load[app.AppSettings]()
if err != nil {
    log.Fatal(err)
}

// Create app with settings
application, err := app.New(appSettings)
```

### Extending AppSettings

Users can create their own Settings struct that embeds AppSettings:

```go
// settings/settings.go in your application
package settings

import (
    "github.com/forgego/forge/pkg/app"
)

type Settings struct {
    app.AppSettings  // Embed all framework settings

    // Add your own settings
    MyApp MyAppConfig `mapstructure:"myapp"`
}

type MyAppConfig struct {
    APIKey    string `mapstructure:"api_key" validate:"required"`
    MaxUsers  int    `mapstructure:"max_users" default:"100"`
}
```

Then load your custom settings:

```go
customSettings, err := settings.Load[Settings]()
appSettings := &customSettings.AppSettings
application, err := app.New(appSettings)
```

## Configuration Sources

The loader supports multiple configuration sources with the following priority:

1. **Environment variables** (highest priority)
2. **Config files** (YAML, JSON, TOML)
3. **Default values** (from struct tags or code)

### Environment Variables

Environment variables are automatically mapped. Use underscores and uppercase:

```bash
DATABASE_URL=postgres://localhost/mydb
SERVER_PORT=8080
ADMIN_PATH=/admin
WORKERS_REDIS_ADDR=localhost:6379
```

### Config Files

Create a `config.yaml` (or `config.json`, `config.toml`) in your project root:

```yaml
database:
  url: postgres://localhost/mydb
  driver: postgres

server:
  host: 0.0.0.0
  port: 8080

security:
  secret_key: your-secret-key
  debug: false

admin:
  enable: true
  path: /admin

workers:
  enable: true
  redis_addr: localhost:6379

myapp: # Your custom settings
  api_key: your-api-key
  max_users: 100
```

## Advanced Usage

### Custom Loader Configuration

```go
loader := settings.NewLoader()
loader.SetConfigName("myconfig")
loader.AddConfigPath("./config")
loader.SetConfigType("yaml")

settings, err := loader.Load[YourSettings]()
```

### Loading from Specific File

```go
settings, err := settings.LoadFromFile[YourSettings]("./config/production.yaml")
```

## Adding Your Own Package Config

If you're creating a third-party package that integrates with Forge:

1. Define your config struct in your package:

```go
// mypackage/config.go
package mypackage

type MyPackageConfig struct {
    Enable bool   `mapstructure:"enable" default:"true"`
    Path   string `mapstructure:"path" default:"/mypackage"`
}
```

2. Users can add it to their Settings:

```go
type Settings struct {
    app.AppSettings
    MyPackage mypackage.MyPackageConfig `mapstructure:"mypackage"`
}
```

This way, each package owns its config, and users compose them as needed.
