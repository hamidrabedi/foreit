# Gogo Framework

> **A comprehensive, modular web framework for Go inspired by Django**

[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Gogo is a production-ready web framework for Go that provides Django-like features with Go's type safety and performance.

## ✨ Features

- 🚀 **Type-Safe** - Uses Go generics and Ent types (no reflection in core paths)
- 🧩 **Modular** - 14 independent, composable modules
- 📡 **RESTful API** - Complete CRUD operations with query processing
- 🎛️ **Admin Interface** - Django-like admin console
- 🔐 **Authentication** - JWT, Session, API Key support
- ⚡ **Fast** - Built on Fiber for high performance
- 🔄 **Background Jobs** - Asynq-powered distributed task queue with scheduling
- 💾 **Caching** - Tag-based cache invalidation
- 🌍 **i18n** - Internationalization support
- 🛠️ **CLI Tools** - Developer-friendly command-line interface

## 📦 Modules

### Core
- `pkg/orm` - Type-safe database operations
- `pkg/settings` - Configuration management
- `pkg/endpoints` - RESTful API framework
- `pkg/routing` - URL routing
- `pkg/pipeline` - Middleware pipeline
- `pkg/auth` - Authentication & authorization

### Advanced
- `pkg/console` - Admin interface
- `pkg/workers` - Background jobs
- `pkg/cache` - Caching
- `pkg/sessions` - Session management

### Utilities
- `pkg/i18n` - Internationalization
- `pkg/static` - Static file serving
- `pkg/utils` - Shared utilities
- `pkg/gogo` - High-level application builder

## 🚀 Quick Start

### Installation

```bash
go get github.com/gogo/pkg/gogo
```

### Basic Example

```go
package main

import (
    "github.com/gogo/pkg/gogo"
    "github.com/gogo/pkg/endpoints"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    app, _ := gogo.New(&gogo.AppConfig{
        Port: 8080,
        SecretKey: "your-secret-key",
    })
    
    // Register resources
    // app.RegisterResource("users", userResource)
    
    app.Listen(":8080")
}
```

## 📚 Documentation

See [docs/README.md](docs/README.md) for complete documentation.

## 🎯 Features in Detail

### RESTful API

```go
// Automatic CRUD endpoints
GET    /api/v1/users          // List users
GET    /api/v1/users/:id      // Get user
POST   /api/v1/users          // Create user
PUT    /api/v1/users/:id      // Update user
DELETE /api/v1/users/:id      // Delete user
```

### Query Processing

```
GET /api/v1/users?page=1&page_size=20
GET /api/v1/users?name__contains=john
GET /api/v1/users?age__gte=18&sort_by=created_at&sort_order=desc
```

### Admin Console

```
GET    /admin                 // List models
GET    /admin/users           // List users
GET    /admin/users/:id       // Show user
POST   /admin/users           // Create user
PUT    /admin/users/:id       // Update user
DELETE /admin/users/:id       // Delete user
```

## 🛠️ CLI Tools

```bash
# Create a new project
gogo startproject myapp

# Create a new app
gogo startapp blog

# Generate code
gogo generate resource User
gogo generate console User
gogo generate policy User

# Run migrations
gogo migrate up
```

## 📖 Examples

- [Complete App](examples/complete-app/)
- [Working Example](examples/working/)
- [Modular Example](examples/modular/)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Django and Django REST Framework
- Built on [Fiber](https://github.com/gofiber/fiber)
- Uses [Ent](https://entgo.io/) for database operations

---

**Made with ❤️ for the Go community**
