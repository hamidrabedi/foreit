# Quick Start Guide - Gogo Framework

## Installation

```bash
go get github.com/gogo/pkg/orm
go get github.com/gogo/pkg/settings
go get github.com/gogo/pkg/endpoints
# ... other modules as needed
```

## Basic Application

```go
package main

import (
    "github.com/gogo/pkg/settings"
    "github.com/gogo/pkg/orm"
    "github.com/gogo/pkg/endpoints"
    "github.com/gogo/pkg/pipeline"
    "github.com/gofiber/fiber/v2"
)

type Config struct {
    DatabaseURL string `env:"DATABASE_URL"`
    Port        int    `env:"PORT" default:"8080"`
}

func main() {
    // Load settings
    cfg := settings.Load[Config]()
    
    // Setup database
    client, _ := orm.NewClient("postgres", cfg.DatabaseURL)
    
    // Setup Fiber
    app := fiber.New()
    
    // Middleware
    app.Use(pipeline.Logging())
    app.Use(pipeline.Recovery())
    
    // API
    router := endpoints.NewRouter(app, "/api")
    // router.RegisterResource("users", userResource)
    
    app.Listen(fmt.Sprintf(":%d", cfg.Port))
}
```

## Using Modules

### ORM - Database Access

```go
// Define repository
type UserRepo struct {
    orm.Repository[models.User, *ent.UserQuery]
}

// Use
users := repo.Query().Where(user.NameEQ("John")).All(ctx)
```

### Settings - Configuration

```go
type Config struct {
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    Port        int    `env:"PORT" default:"8080"`
}

cfg := settings.Load[Config]()
```

### Endpoints - API

```go
type UserResource struct {
    endpoints.Resource[models.User, *ent.UserQuery]
}

func (r *UserResource) Index(ctx *endpoints.Context) ([]*models.User, error) {
    return r.Repo.Query().All(ctx.Request.Context())
}

router := endpoints.NewRouter(app, "/api")
router.RegisterResource("users", &UserResource{})
```

### Auth - Authentication

```go
jwt := auth.NewJWT("secret")
app.Use(auth.Middleware(jwt))

// Check permissions
if err := auth.Require[models.Post](ctx, "edit", post); err != nil {
    return err
}
```

### Workers - Background Jobs

```go
job := &SendEmailJob{To: "user@example.com"}
workers.Enqueue(ctx, job)

workers.Start(ctx, 5)
```

### Cache - Caching

```go
cache.Set(ctx, "user:123", user, 10*time.Minute)
user, _ := cache.Get(ctx, "user:123")
```

## Next Steps

- See `examples/modular/` for a complete example
- Read module-specific documentation in each `pkg/*/README.md`
- Check `docs/MODULAR_ARCHITECTURE.md` for architecture details

