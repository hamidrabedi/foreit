# forge Framework

A Django-like Go framework with full type safety, code generation, and extensibility.

## Features

- **Type-Safe ORM**: Full Django ORM features with compile-time type checking
- **Code Generation**: AST-based code generation for models, managers, and querysets
- **Auto-Generated Admin**: Django-style admin interface auto-generated from model registry
- **Dual API**: Type-safe primary API + dynamic secondary API for runtime flexibility
- **Migration System**: Built-in migration system with golang-migrate
- **Security**: Built-in CSRF, XSS, and SQL injection protection
- **Extensible**: Everything is extendable/overridable via plugins
- **SQL Builder**: Type-safe SQL generation with proper escaping and parameter binding

## Quick Start

### 1. Create a New Project

```bash
forge new myapp
cd myapp
```

### 2. Configure Database

Edit `config/config.yaml` with your PostgreSQL credentials.

### 3. Define Models

Edit `models/example.go` or create new model files:

```go
package models

import "github.com/forgego/forge/internal/schema"

type Post struct {
	schema.BaseSchema
}

func (Post) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("title").Required().MaxLength(200).Build(),
		schema.String("content").Required().Build(),
	}
}

func (Post) Meta() schema.Meta {
	return schema.Meta{
		TableName: "posts",
		VerboseName: "Post",
	}
}

func (Post) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Post) Hooks() *schema.ModelHooks {
	return nil
}
```

### 4. Generate Code

```bash
forge generate
```

### 5. Run Migrations

```bash
forge migrate
```

### 6. Start Server

```bash
forge runserver
```

Visit `http://localhost:8000/admin/` for the auto-generated admin interface!

For detailed instructions, see the [Getting Started Guide](docs/GETTING_STARTED.md).

## Documentation

- **[Getting Started](docs/GETTING_STARTED.md)** - 10-minute quick start guide ⭐
- **[REST API Guide](docs/REST_API.md)** - Build APIs for React/Vue frontends
- **[HTMX Patterns](docs/HTMX_PATTERNS.md)** - Admin interface patterns and best practices
- **[UI Strategy](docs/UI_STRATEGY.md)** - UI technology choices and approach
- **[Architecture](docs/ARCHITECTURE.md)** - Framework architecture
- **[Schema Reference](docs/SCHEMA_REFERENCE.md)** - Schema definition guide
- **[API Reference](docs/API_REFERENCE.md)** - API documentation
- **[Usage Guide](docs/USAGE_GUIDE.md)** - Step-by-step tutorials
- **[Features](docs/FEATURES.md)** - Current and planned features
- **[Roadmap](docs/ROADMAP.md)** - Development roadmap
- **[Development](docs/DEVELOPMENT.md)** - Contributing guide

## Example

See `examples/blog/` for a complete example application.

## Technology Stack

- **Go 1.24+**
- **chi/v5** - HTTP router
- **database/sql** - Standard SQL interface
- **golang-migrate** - Migrations
- **zap** - Logging
- **viper** - Configuration
- **testify** - Testing

## Status

**Current:** MVP Complete - Core features working! 🎉

**Implementation Status:**
- ✅ Schema system - Complete
- ✅ Code generation - Complete
- ✅ Type-safe ORM - Complete (All, Get, First, Last, Count, Exists)
- ✅ Q Objects - Complete (And, Or, Not)
- ✅ Manager CRUD - Complete (Create, Update, Delete with hooks)
- ✅ Admin Interface - Complete (List view, Create/Edit forms)
- ✅ REST API - Complete (Full CRUD with pagination, filtering, ordering)
- ✅ Plugin System - Complete
- ✅ CLI - Complete (new, generate, migrate, runserver)

See [Features](docs/FEATURES.md) and [Roadmap](docs/ROADMAP.md) for details.

## Contributing

See [Development Guide](docs/DEVELOPMENT.md) for contribution guidelines.

## License

MIT
