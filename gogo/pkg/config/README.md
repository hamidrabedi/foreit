# Config Module

Configuration management using [viper](https://github.com/spf13/viper) and [godotenv](https://github.com/joho/godotenv).

## Features

- Load from .env files
- Load from YAML/JSON config files
- Environment variable support
- Automatic type conversion
- Default values
- Struct unmarshaling

## Usage

### Basic Usage

```go
import "github.com/gogo/pkg/config"

loader, _ := config.Load()
dbURL := loader.GetString("DATABASE_URL")
port := loader.GetInt("PORT")
```

### Load into Struct

```go
type AppConfig struct {
    DatabaseURL string `mapstructure:"database_url"`
    Port        int    `mapstructure:"port"`
    Debug       bool   `mapstructure:"debug"`
}

cfg, _ := config.LoadInto[AppConfig]("")
```

### Load from .env

```go
loader := config.NewLoader()
loader.LoadFromEnvFile(".env")
```

### Load from YAML

```go
loader := config.NewLoader()
loader.LoadFromFile("config.yaml")
```

