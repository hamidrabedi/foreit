# Project Structure Guide

Forge provides two project structure templates to suit different needs.

## Template 1: Simple App-Based (Default)

The default template is simple and familiar, perfect for most projects:

```
myproject/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── app/                          # USER CODE - Django-style apps
│   ├── users/
│   │   ├── models.go            # ORM models
│   │   ├── admin.go             # Admin configuration (per-app)
│   │   ├── api.go               # REST/GraphQL endpoints
│   │   └── handlers.go          # HTTP handlers (optional)
│   └── blog/
│       ├── models.go
│       ├── admin.go
│       └── api.go
├── migrations/                   # Database migrations
├── config/
│   └── config.yaml              # Configuration
├── static/                       # Static files
├── templates/                    # HTML templates
├── go.mod
└── README.md
```

**Key points:**
- Admin config lives in each app's `admin.go` (no top-level admin directory)
- Simple, Django-like structure
- Everything essential, nothing extra
- Perfect for most projects

## Template 2: Advanced Hybrid Structure (Optional)

For larger projects or teams wanting clean architecture:

```
myproject/
├── cmd/
│   └── server/
│       └── main.go
├── app/                          # USER CODE - Django-style apps
│   ├── users/
│   │   ├── models.go
│   │   ├── admin.go
│   │   └── api.go
│   └── blog/
│       ├── models.go
│       ├── admin.go
│       └── api.go
├── domain/                       # PURE BUSINESS LOGIC (optional)
│   ├── user/
│   │   ├── entity.go
│   │   ├── service.go
│   │   └── policy.go
│   └── billing/
│       └── service.go
├── infra/                        # INFRASTRUCTURE (external services)
│   ├── db/
│   │   └── connection.go
│   ├── redis/
│   │   └── client.go
│   └── email/
│       └── sender.go
├── pkg/                          # SHARED UTILITIES (cross-app)
│   ├── validators/
│   ├── utils/
│   └── types/
├── migrations/
├── config/
├── static/
├── templates/
├── go.mod
└── README.md
```

**Key points:**
- Includes `domain/` for clean architecture
- Includes `infra/` for infrastructure concerns
- Includes `pkg/` for shared utilities
- Admin still lives in each app's `admin.go` (no top-level admin directory)

## App Structure

Each app follows a flat file structure:

- `models.go` - ORM model definitions
- `admin.go` - Admin configuration (auto-registers via `init()`)
- `api.go` - REST/GraphQL API endpoints (auto-registers via `init()`)
- `handlers.go` - HTTP handlers (optional)
- `services.go` - App-level services (optional)
- `tests/` - App-specific tests (optional)

## Creating Projects

```bash
# Create simple project (default)
forge new myproject

# Create advanced project
forge new myproject --template advanced
```

## Adding Apps

```bash
# Add a new app
forge add app blog

# This creates:
# app/blog/models.go
# app/blog/admin.go
# app/blog/api.go
```

## Adding Components

```bash
# Add a model
forge add model --app blog --name Post

# Add a handler
forge add handler --app blog --name list_posts --method GET --path /posts

# Add an API endpoint
forge add api --app blog --model Post --name posts

# Add a service
forge add service --app blog --name PostService
```

## Best Practices

1. **Keep apps focused**: Each app should represent a single domain concept
2. **Use flat files**: Keep models, admin, and API in single files per app
3. **Auto-registration**: Use `init()` functions for automatic discovery
4. **Shared code**: Use `pkg/` (advanced template) for cross-app utilities
5. **Domain logic**: Use `domain/` (advanced template) for pure business logic

