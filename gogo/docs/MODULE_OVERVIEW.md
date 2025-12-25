# Gogo Framework - Module Overview

## Core Modules

### pkg/orm - Database Layer
Type-safe database operations using Ent's code generation.

**Key Features**:
- Repository pattern
- Type-safe queries
- Migrations
- Bulk operations
- Transactions

**Usage**:
```go
client, _ := orm.NewClient("postgres", dsn)
repo := &UserRepository{Client: client}
users := repo.Query().All(ctx)
```

### pkg/settings - Configuration
Application configuration management.

**Key Features**:
- Environment variable loading
- Configuration validation
- Default values
- Type conversion

**Usage**:
```go
type Config struct {
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    Port        int    `env:"PORT" default:"8080"`
}
cfg := settings.Load[Config]()
```

### pkg/endpoints - API Framework
RESTful API framework with resource handlers.

**Key Features**:
- Resource handlers
- Serializers
- Query processing
- Error handling
- Automatic routing

**Usage**:
```go
type UserResource struct {
    endpoints.Resource[models.User, *ent.UserQuery]
}
router.RegisterResource("users", &UserResource{})
```

### pkg/routing - URL Routing
URL routing with reverse lookup.

**Key Features**:
- Named routes
- Reverse URL generation
- Route groups
- Route middleware

**Usage**:
```go
router.Get("/users/:id", handler, routing.Name("users.show"))
url, _ := router.URL("users.show", routing.Param("id", "123"))
```

### pkg/pipeline - Middleware Pipeline
Request/response middleware pipeline.

**Key Features**:
- Built-in middleware (logging, recovery, CORS, etc.)
- Rate limiting
- Security headers
- Middleware chaining

**Usage**:
```go
app.Use(pipeline.Logging())
app.Use(pipeline.Recovery())
app.Use(pipeline.RateLimit(100, time.Minute))
```

### pkg/auth - Authentication & Authorization
Authentication and authorization system.

**Key Features**:
- Multiple authenticators (JWT, Session, API Key)
- Policy-based authorization
- Role-based access control
- Middleware integration

**Usage**:
```go
jwt := auth.NewJWT("secret")
app.Use(auth.Middleware(jwt))

auth.Register[models.Post](&PostPolicy{})
auth.Require[models.Post](ctx, "edit", post)
```

### pkg/console - Admin Interface
Django-like admin interface.

**Key Features**:
- Auto-generated admin from Ent schemas
- Custom console classes
- List/Detail views
- Filters and actions

**Usage**:
```go
console.Register[models.User](&console.Options{
    ListDisplay: []string{"name", "email"},
})
console.InstallRoutes(app, "/console")
```

### pkg/workers - Background Jobs
Background job processing powered by Asynq.

**Key Features**:
- Redis-based distributed queue
- Automatic retries with exponential backoff
- Cron-based scheduling
- Priority queues
- Built-in Web UI
- Delayed jobs

**Usage**:
```go
queue, _ := workers.NewAsynqQueue("localhost:6379", "", 0)
workers.SetDefaultQueue(queue)
workers.Register("send_email", &EmailHandler{})

job := &SendEmailJob{To: "user@example.com"}
workers.Enqueue(ctx, job)
workers.Start(ctx, 10)
workers.Schedule("0 6 * * *", &DailyReportJob{})
```

### pkg/cache - Caching
Caching layer with tags.

**Key Features**:
- In-memory store
- TTL support
- Tag-based invalidation
- Remember pattern

**Usage**:
```go
cache.Set(ctx, "user:123", user, 10*time.Minute)
cache.TagSet(ctx, "post:456", post, 1*time.Hour, "posts")
cache.TagInvalidate(ctx, "posts")
```

### pkg/sessions - Session Management
Session management for Fiber.

**Key Features**:
- Cookie-based sessions
- Automatic expiration
- Session regeneration
- Configurable lifetime

**Usage**:
```go
store := sessions.NewMemoryStore(sessions.DefaultConfig())
app.Use(sessions.Middleware(store))
sessions.Set(c, "user_id", user.ID)
```

## Module Dependencies

```
pkg/orm (foundation)
    ↓
pkg/endpoints (uses orm)
    ↓
pkg/console (uses orm, endpoints)
    ↓
pkg/auth (standalone)
pkg/pipeline (standalone)
pkg/routing (standalone)
pkg/workers (standalone)
pkg/cache (standalone)
pkg/sessions (standalone)
pkg/settings (standalone)
```

## Composing Modules

All modules are independent and can be used standalone or together:

```go
// Use only what you need
app.Use(pipeline.Logging())  // Just logging

// Or compose multiple modules
app.Use(pipeline.Logging())
app.Use(pipeline.Recovery())
app.Use(auth.Middleware(jwt))
app.Use(sessions.Middleware(store))
```

## Design Principles

1. **Type Safety**: Uses Go generics and Ent types
2. **Composability**: Modules work standalone or together
3. **Go Idioms**: Patterns designed for Go
4. **Original**: Unique naming and concepts
5. **Comprehensive**: More features than Django

