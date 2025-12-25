# Gogo Package - High-Level Application Builder

The `gogo` package provides a high-level API for building applications with all Gogo modules.

## Usage

### Simple Application

```go
package main

import (
    "github.com/gogo/pkg/gogo"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app, err := gogo.New(&gogo.AppConfig{
        DatabaseURL: "postgres://...",
        Port: 8080,
        SecretKey: "your-secret-key",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Register resources
    // app.RegisterResource("users", userResource)
    
    // Add routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Hello Gogo!"})
    })
    
    app.Listen(":8080")
}
```

### With All Features

```go
app, err := gogo.New(&gogo.AppConfig{
    DatabaseURL: "postgres://...",
    Port: 8080,
    SecretKey: "secret",
    EnableConsole: true,
    EnableWorkers: true,
    EnableCache: true,
    EnableSessions: true,
    EnableI18n: true,
    EnableStatic: true,
    ConsolePath: "/admin",
    APIPath: "/api",
    StaticPath: "/static",
    StaticRoot: "./public",
    LocalesPath: "./locales",
    DefaultLocale: "en",
    RedisAddr: "localhost:6379",
    RedisPassword: "",
    RedisDB: 0,
    WorkersConcurrency: 10,
})
```

## Features

- Automatic module setup
- Configuration management
- Resource registration
- Route management
- Graceful shutdown

