# Migration Guide - From Reflection to Static

This guide helps migrate from the old reflection-based code to the new modular, static architecture.

## Overview

The old architecture used reflection extensively. The new architecture uses Go generics and Ent's code generation for type safety.

## Key Changes

### 1. Database Access

**Before (Reflection)**:
```go
helper := NewEntClientHelper(client)
result, err := helper.Create(ctx, "User", data)
```

**After (Generics)**:
```go
repo := &UserRepository{Client: client}
result, err := repo.Create(ctx, user)
```

### 2. API Handlers

**Before (Reflection)**:
```go
handler := NewCRUDHandler(meta, client, registry)
// Uses reflection internally
```

**After (Generics)**:
```go
resource := &UserResource{
    Resource: endpoints.NewResource[models.User](client),
}
router.RegisterResource("users", resource)
```

### 3. Configuration

**Before**:
```go
// Manual env var reading
dbURL := os.Getenv("DATABASE_URL")
```

**After**:
```go
type Config struct {
    DatabaseURL string `env:"DATABASE_URL"`
}
cfg := settings.Load[Config]()
```

### 4. Middleware

**Before**:
```go
// Custom middleware
app.Use(func(c *fiber.Ctx) error {
    // ...
})
```

**After**:
```go
app.Use(pipeline.Logging())
app.Use(pipeline.Recovery())
app.Use(pipeline.CORS())
```

### 5. Authentication

**Before**:
```go
// Custom auth logic
```

**After**:
```go
jwt := auth.NewJWT("secret")
app.Use(auth.Middleware(jwt))
```

## Migration Steps

### Step 1: Update Imports

Replace old imports:
```go
// Old
import "github.com/gogo/internal/admin"

// New
import (
    "github.com/gogo/pkg/orm"
    "github.com/gogo/pkg/endpoints"
    "github.com/gogo/pkg/settings"
    // ... other modules
)
```

### Step 2: Replace Database Access

Create repositories for each model:
```go
type UserRepository struct {
    Client *ent.Client
}

func (r *UserRepository) Query() *ent.UserQuery {
    return r.Client.User.Query()
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
    return r.Client.User.Get(ctx, id)
}
```

### Step 3: Replace Handlers

Create resource handlers:
```go
type UserResource struct {
    endpoints.Resource[models.User, *ent.UserQuery]
}

func (r *UserResource) Index(ctx *endpoints.Context) ([]*models.User, error) {
    return r.Repo.Query().All(ctx.Request.Context())
}
```

### Step 4: Update Configuration

Use settings module:
```go
type Config struct {
    DatabaseURL string `env:"DATABASE_URL"`
    Port        int    `env:"PORT" default:"8080"`
}
cfg := settings.Load[Config]()
```

### Step 5: Update Middleware

Use pipeline module:
```go
app.Use(pipeline.Logging())
app.Use(pipeline.Recovery())
app.Use(pipeline.CORS())
```

## Benefits

1. **Type Safety**: Compile-time errors instead of runtime
2. **Performance**: No reflection overhead
3. **IDE Support**: Full autocomplete
4. **Maintainability**: Clear types
5. **Modularity**: Use only what you need

## Backward Compatibility

The old code still works. You can migrate gradually:
1. Start with new modules for new features
2. Migrate existing code module by module
3. Remove old code when migration is complete

