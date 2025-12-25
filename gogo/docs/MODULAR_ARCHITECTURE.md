# Gogo Framework Architecture

## Vision

A comprehensive, modular web framework for Go that provides everything you need to build modern applications - from database models to APIs, admin interfaces, background jobs, and more. Inspired by Django's "batteries-included" philosophy but designed for Go's strengths.

## Core Philosophy

1. **Composable Modules**: Use only what you need, everything works together
2. **Type Safety First**: Leverage Ent's code generation for compile-time safety
3. **Go Idioms**: Patterns that feel natural in Go, not Python translations
4. **Developer Experience**: Fast feedback, clear errors, great tooling
5. **Production Ready**: Built-in features for real-world applications

## Module Structure

```
gogo/
├── pkg/
│   ├── orm/            # Database layer (Ent integration & helpers)
│   ├── settings/       # Application configuration
│   ├── console/        # Admin interface (like Django admin)
│   ├── endpoints/      # API framework (REST + more)
│   ├── routing/        # URL routing & reverse lookup
│   ├── auth/           # Authentication & authorization
│   ├── pipeline/       # Request/response middleware pipeline
│   ├── workers/        # Background job processing
│   ├── cache/          # Caching layer
│   ├── sessions/       # Session management
│   ├── i18n/           # Internationalization
│   ├── static/         # Static file serving
│   └── utils/          # Shared utilities
```

## Module 1: ORM (`pkg/orm`)

**Purpose**: Type-safe database operations using Ent

**Key Concepts**:
- **Repository Pattern**: Type-safe repositories for each model
- **Query Sets**: Fluent query builders (using Ent's generated code)
- **Migrations**: Schema versioning and migrations
- **Relationships**: Type-safe relationship handling

**Usage**:
```go
// Repository for a model
type UserRepo struct {
    orm.Repository[models.User, *ent.UserQuery]
}

// Type-safe queries
users := repo.Query().
    Where(user.NameEQ("John")).
    Limit(10).
    All(ctx)

// Create with validation
user, err := repo.Create(ctx, &models.User{
    Name: "John",
    Email: "john@example.com",
})
```

**Features**:
- Bulk operations
- Transactions
- Query optimization hints
- Connection pooling
- Database-specific features

## Module 2: Settings (`pkg/settings`)

**Purpose**: Application configuration management

**Key Concepts**:
- **Settings Registry**: Centralized configuration
- **Environment Loading**: Automatic env var mapping
- **Validation**: Type-safe config validation
- **Hot Reload**: Optional runtime config updates

**Usage**:
```go
type AppSettings struct {
    DatabaseURL string `setting:"database.url" env:"DB_URL" required:"true"`
    Port        int    `setting:"server.port" env:"PORT" default:"8080"`
    Debug       bool   `setting:"app.debug" env:"DEBUG" default:"false"`
}

settings := settings.Load[AppSettings]()
// Access: settings.DatabaseURL, settings.Port, etc.
```

**Features**:
- Multiple sources (env, files, CLI args)
- Nested configuration
- Type conversion
- Validation rules
- Default values
- Secret management

## Module 3: Console (`pkg/console`)

**Purpose**: Admin interface for managing data (like Django admin)

**Key Concepts**:
- **Model Console**: Auto-generated admin for each model
- **Console Actions**: Custom bulk actions
- **Console Filters**: Advanced filtering UI
- **Console Forms**: Custom form layouts

**Usage**:
```go
// Auto-register with defaults
console.Register[models.User]()

// Or with custom console
type UserConsole struct {
    console.ModelConsole[models.User]
}

func (c *UserConsole) DisplayFields() []string {
    return []string{"name", "email", "created_at"}
}

func (c *UserConsole) ListFilters() []console.Filter {
    return []console.Filter{
        console.DateFilter("created_at"),
        console.ChoiceFilter("status", []string{"active", "inactive"}),
    }
}

console.Register[models.User](&UserConsole{})
```

**Features**:
- Auto-generated forms from Ent schemas
- Inline editing for related models
- Custom actions (bulk operations)
- Export/import (CSV, JSON, Excel)
- Search and filtering
- Permission-based field visibility

## Module 4: Endpoints (`pkg/endpoints`)

**Purpose**: API framework for building RESTful and other APIs

**Key Concepts**:
- **Resource Handlers**: Type-safe request handlers
- **Serializers**: Request/response transformation
- **Request Processors**: Query parsing, filtering, pagination
- **Response Formatters**: JSON, XML, MessagePack, etc.

**Usage**:
```go
// Define a resource handler
type UserResource struct {
    endpoints.Resource[models.User, *ent.UserQuery]
}

func (r *UserResource) Index(ctx *endpoints.Context) ([]models.User, error) {
    // List users with filtering, pagination, etc.
    return r.Query().
        Where(user.ActiveEQ(true)).
        All(ctx.Request.Context())
}

func (r *UserResource) Show(ctx *endpoints.Context) (*models.User, error) {
    id := ctx.Param("id")
    return r.Repo().GetByID(ctx.Request.Context(), id)
}

func (r *UserResource) Create(ctx *endpoints.Context) (*models.User, error) {
    var data CreateUserRequest
    if err := ctx.Bind(&data); err != nil {
        return nil, err
    }
    
    // Validation happens automatically via serializer
    return r.Repo().Create(ctx.Request.Context(), &data)
}

// Register
endpoints.RegisterResource("users", &UserResource{
    Resource: endpoints.NewResource[models.User](client),
})
```

**Features**:
- Automatic request parsing
- Built-in filtering (query params → Ent predicates)
- Pagination helpers
- Sorting
- Field selection (sparse fieldsets)
- Content negotiation
- Rate limiting per endpoint
- Request/response logging

## Module 5: Routing (`pkg/routing`)

**Purpose**: URL routing with reverse lookup

**Key Concepts**:
- **Route Groups**: Organize routes by feature
- **Named Routes**: Reverse URL generation
- **Route Middleware**: Per-route middleware
- **Route Constraints**: Type-safe path parameters

**Usage**:
```go
router := routing.NewRouter()

// Simple route
router.Get("/", homeHandler, routing.Name("home"))

// Route group
api := router.Group("/api/v1", middleware.Auth())
api.Get("/users", listUsers, routing.Name("api.users.list"))
api.Get("/users/:id", showUser, routing.Name("api.users.show"))

// Reverse URL
url := router.URL("api.users.show", routing.Param("id", "123"))
// Returns: "/api/v1/users/123"
```

**Features**:
- Named routes
- Reverse URL generation
- Route groups with middleware
- Type-safe path parameters
- Route caching
- Route listing (for debugging)

## Module 6: Auth (`pkg/auth`)

**Purpose**: Authentication and authorization

**Key Concepts**:
- **Authenticators**: Multiple auth methods (JWT, Session, API Key, etc.)
- **Policies**: Flexible permission system
- **Guards**: Route protection
- **User Context**: Request-scoped user access

**Usage**:
```go
// Define policies
type PostPolicy struct {
    auth.Policy[models.Post]
}

func (p *PostPolicy) CanView(user *models.User, post *models.Post) bool {
    return post.Published || post.AuthorID == user.ID
}

func (p *PostPolicy) CanEdit(user *models.User, post *models.Post) bool {
    return post.AuthorID == user.ID || user.IsAdmin
}

// Use in endpoints
func (r *PostResource) Show(ctx *endpoints.Context) (*models.Post, error) {
    post, err := r.Repo().GetByID(ctx.Request.Context(), ctx.Param("id"))
    if err != nil {
        return nil, err
    }
    
    // Check policy
    if !auth.Can(ctx.User(), "view", post) {
        return nil, endpoints.ErrForbidden
    }
    
    return post, nil
}
```

**Features**:
- Multiple authentication methods
- Policy-based authorization
- Field-level permissions
- Role-based access control (RBAC)
- Permission caching
- Audit logging

## Module 7: Pipeline (`pkg/pipeline`)

**Purpose**: Request/response middleware pipeline

**Key Concepts**:
- **Middleware Chain**: Composable middleware
- **Pipeline Stages**: Request → Handler → Response
- **Context Propagation**: Shared context through pipeline
- **Error Handling**: Centralized error processing

**Usage**:
```go
// Built-in middleware
app.Use(pipeline.Logging())
app.Use(pipeline.Recovery())
app.Use(pipeline.CORS())
app.Use(pipeline.RateLimit(100, time.Minute))

// Custom middleware
app.Use(func(next pipeline.Handler) pipeline.Handler {
    return func(ctx *pipeline.Context) error {
        // Before handler
        ctx.Set("start_time", time.Now())
        
        err := next(ctx)
        
        // After handler
        duration := time.Since(ctx.Get("start_time").(time.Time))
        log.Printf("Request took %v", duration)
        
        return err
    }
})
```

**Features**:
- Request logging
- Error recovery
- CORS handling
- Rate limiting
- Request ID generation
- Compression
- Security headers

## Module 8: Workers (`pkg/workers`)

**Purpose**: Background job processing

**Key Concepts**:
- **Job Queue**: Async job processing
- **Job Types**: Type-safe job definitions
- **Schedulers**: Cron-like scheduling
- **Retries**: Automatic retry with backoff

**Usage**:
```go
// Define a job
type SendEmailJob struct {
    To      string
    Subject string
    Body    string
}

func (j *SendEmailJob) Execute(ctx context.Context) error {
    return sendEmail(j.To, j.Subject, j.Body)
}

// Enqueue job
workers.Enqueue(&SendEmailJob{
    To: "user@example.com",
    Subject: "Welcome",
    Body: "Welcome to our app!",
})

// Schedule recurring job
workers.Schedule("0 0 * * *", &DailyReportJob{})
```

**Features**:
- Redis-based distributed queue (Asynq)
- Automatic retries with exponential backoff
- Cron-based scheduling
- Job priorities (critical, default, low)
- Delayed jobs
- Built-in Web UI for monitoring
- Rate limiting support

## Module 9: Cache (`pkg/cache`)

**Purpose**: Caching layer

**Key Concepts**:
- **Cache Stores**: Multiple backends (memory, Redis, etc.)
- **Cache Tags**: Tag-based invalidation
- **Cache Patterns**: Cache-aside, write-through, etc.

**Usage**:
```go
// Cache with TTL
cache.Set("user:123", user, 5*time.Minute)

// Cache with tags
cache.Set("post:456", post, cache.Tags("posts", "user:123"))

// Invalidate by tag
cache.Invalidate("posts") // Invalidates all posts

// Cache helper for endpoints
func (r *UserResource) Show(ctx *endpoints.Context) (*models.User, error) {
    id := ctx.Param("id")
    
    var user models.User
    if err := cache.Get("user:"+id, &user); err == nil {
        return &user, nil
    }
    
    user, err := r.Repo().GetByID(ctx.Request.Context(), id)
    if err != nil {
        return nil, err
    }
    
    cache.Set("user:"+id, user, 10*time.Minute)
    return &user, nil
}
```

## Module 10: Sessions (`pkg/sessions`)

**Purpose**: Session management

**Usage**:
```go
// Configure session store
app.Use(sessions.New(sessions.Config{
    Store: sessions.NewCookieStore("secret-key"),
    Lifetime: 24 * time.Hour,
}))

// Use in handlers
func loginHandler(ctx *fiber.Ctx) error {
    session := sessions.Get(ctx)
    session.Set("user_id", user.ID)
    return session.Save()
}
```

## Module 11: i18n (`pkg/i18n`)

**Purpose**: Internationalization

**Usage**:
```go
// Load translations
i18n.Load("locales")

// Use in handlers
func handler(ctx *fiber.Ctx) error {
    t := i18n.T(ctx, "welcome.message")
    return ctx.JSON(fiber.Map{"message": t})
}
```

## Module 12: Static (`pkg/static`)

**Purpose**: Static file serving

**Usage**:
```go
app.Use(static.New(static.Config{
    Root: "./public",
    Prefix: "/static",
    Compress: true,
}))
```

## Complete Example

```go
package main

import (
    "github.com/gogo/pkg/orm"
    "github.com/gogo/pkg/settings"
    "github.com/gogo/pkg/console"
    "github.com/gogo/pkg/endpoints"
    "github.com/gogo/pkg/routing"
    "github.com/gogo/pkg/auth"
    "github.com/gogo/pkg/pipeline"
    "github.com/gogo/pkg/workers"
    "github.com/gofiber/fiber/v2"
)

type AppSettings struct {
    DatabaseURL string `env:"DB_URL"`
    Port        int    `env:"PORT" default:"8080"`
}

func main() {
    // Load settings
    cfg := settings.Load[AppSettings]()
    
    // Setup database
    client := orm.NewClient(cfg.DatabaseURL)
    
    // Setup Fiber
    app := fiber.New()
    
    // Pipeline
    app.Use(pipeline.Logging())
    app.Use(pipeline.Recovery())
    app.Use(pipeline.CORS())
    
    // Auth
    app.Use(auth.Middleware(auth.JWT()))
    
    // Routing
    router := routing.NewRouter()
    
    // Endpoints
    endpoints.RegisterResource("users", &UserResource{
        Resource: endpoints.NewResource[models.User](client),
    })
    
    // Console (admin)
    console.Register[models.User]()
    console.Register[models.Post]()
    console.Install(app, "/console")
    
    // Workers
    workers.Start()
    
    // Install routes
    router.Install(app)
    
    app.Listen(fmt.Sprintf(":%d", cfg.Port))
}
```

## Design Principles

1. **Go-First**: Patterns that make sense in Go, not Python translations
2. **Type Safety**: Leverage Ent's code generation everywhere
3. **Composability**: Every module works standalone or together
4. **Developer Experience**: Fast feedback, clear errors
5. **Production Ready**: Built-in features for real apps

## Beyond Django

Additional features Django doesn't have:
- Type-safe everything (thanks to Ent)
- More flexible auth system
- Native WebSocket support (via Fiber)
- gRPC support (planned)
