# Gogo Framework - Complete Guide

## Quick Start

### 1. Create a Project

```bash
gogo startproject myapp
cd myapp
```

### 2. Create an App

```bash
gogo startapp blog
```

### 3. Define Models

Create `models/user.go`:

```go
package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name" required:"true"`
    Email     string    `json:"email" required:"true" email:"true"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 4. Create Repository

Create `repositories/user.go`:

```go
package repositories

import (
    "context"
    "yourproject/models"
    "github.com/gogo/pkg/orm"
)

type UserRepository struct {
    *orm.BaseRepository[*models.User, interface{}]
}

func NewUserRepository(client *orm.Client) *UserRepository {
    return &UserRepository{
        BaseRepository: orm.NewBaseRepository[*models.User, interface{}](
            client,
            func(c interface{}) interface{} {
                // Return Ent query
            },
            func(c interface{}) interface{} {
                // Return Ent create builder
            },
            func(c interface{}, id interface{}) interface{} {
                // Return Ent update builder
            },
        ),
    }
}
```

### 5. Create Resource

Create `resources/user.go`:

```go
package resources

import (
    "yourproject/models"
    "yourproject/repositories"
    "github.com/gogo/pkg/endpoints"
)

type UserResource struct {
    *endpoints.BaseResource[*models.User, interface{}]
}

func NewUserResource(repo *repositories.UserRepository) *UserResource {
    return &UserResource{
        BaseResource: endpoints.NewResource[*models.User, interface{}](repo),
    }
}
```

### 6. Register in main.go

```go
package main

import (
    "github.com/gogo/pkg/gogo"
    "yourproject/resources"
    "yourproject/repositories"
)

func main() {
    app, _ := gogo.New(&gogo.AppConfig{
        DatabaseURL: "postgres://...",
        Port: 8080,
        SecretKey: "your-secret-key",
    })
    
    // Create repository
    userRepo := repositories.NewUserRepository(app.Client())
    
    // Create and register resource
    userResource := resources.NewUserResource(userRepo)
    app.RegisterResource("users", userResource)
    
    app.Listen(":8080")
}
```

## API Endpoints

### List Resources

```
GET /api/v1/users?page=1&page_size=20&sort_by=name&sort_order=asc
```

### Get Resource

```
GET /api/v1/users/:id
```

### Create Resource

```
POST /api/v1/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

### Update Resource

```
PUT /api/v1/users/:id
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

### Delete Resource

```
DELETE /api/v1/users/:id
```

## Query Parameters

### Filters

```
GET /api/v1/users?name__contains=john
GET /api/v1/users?age__gte=18
GET /api/v1/users?status__in=active,pending
```

### Pagination

```
GET /api/v1/users?page=2&page_size=10
```

### Sorting

```
GET /api/v1/users?sort_by=created_at&sort_order=desc
```

## Console (Admin)

### List Models

```
GET /admin
```

### List Model Items

```
GET /admin/users?page=1&page_size=20
```

### Create Item

```
POST /admin/users
```

### Update Item

```
PUT /admin/users/:id
```

### Delete Item

```
DELETE /admin/users/:id
```

## Authentication

```go
jwtAuth := auth.NewJWT("secret-key")
app.SetAuth(jwtAuth, true) // Required auth
// or
app.SetAuth(jwtAuth, false) // Optional auth
```

## Background Jobs

```go
workers.Register("send_email", &EmailJobHandler{})
workers.Start(context.Background(), 5)
```

## Caching

```go
cache.SetDefaultStore(cache.NewTaggedMemoryStore())
cache.Get("key")
cache.Set("key", value, time.Hour)
```

## Sessions

```go
sessionStore := sessions.NewMemoryStore(sessions.DefaultConfig())
app.Use(sessions.Middleware(sessionStore))
```

## Internationalization

```go
i18n.Load("locales", "en")
app.Use(i18n.Middleware())
```

## Static Files

```go
app.Use(static.New(static.Config{
    Root: "./public",
    Prefix: "/static",
}))
```

## Best Practices

1. **Use Repositories** - Keep data access logic in repositories
2. **Use Resources** - Keep API logic in resources
3. **Use Serializers** - Transform data properly
4. **Validate Input** - Always validate user input
5. **Handle Errors** - Use structured error responses
6. **Use Middleware** - Add cross-cutting concerns via middleware
7. **Test Everything** - Write tests for your code

## Examples

See `examples/complete-app/` for a complete working example.

