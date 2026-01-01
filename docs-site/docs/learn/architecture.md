---
sidebar_position: 2
description: Understand forge's architecture, design principles, and how it works. Learn about code generation, layers, components, and extension points.
keywords:
  - forge architecture
  - forge design
  - code generation
  - forge internals
  - framework architecture
image: /img/forge-social-card.jpg
---

# Architecture

Understanding forge's architecture helps you build better applications and extend the framework.

## How forge works

forge uses a layered architecture that separates concerns and enables code generation:

1. **You define schemas** - Declarative model definitions
2. **forge generates code** - Type-safe structs, managers, querysets
3. **You use the ORM** - Type-safe database queries
4. **forge handles the rest** - Migrations, admin, APIs

This architecture gives you Django-like productivity with Go's type safety and performance.

## Architecture layers

forge is organized into distinct layers, each with a specific responsibility:

```
┌─────────────────────────────────────────────────────────┐
│                    User Application                      │
│  (Schema Definitions, Models, Views, Controllers)        │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│              Code Generation Layer                        │
│  (AST Parser → SQL Generation → Go Code Generation)      │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│              Framework API Layer                          │
│  (QuerySet, Manager, FieldExpr, QueryExpr)               │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│              Database Layer                               │
│  (SQL Builder, Parameter Binding, Transactions)          │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│              Infrastructure Layer                         │
│  (HTTP Router, Middleware, Security, Config, Logging)     │
└─────────────────────────────────────────────────────────┘
```

## Core principles

forge follows these design principles:

1. **Type-Safe First** - Primary API uses generics for compile-time safety
2. **Dynamic When Needed** - Secondary API for runtime flexibility
3. **Convention over Configuration** - Sensible defaults everywhere
4. **Fully Extensible** - Everything can be extended/overridden
5. **Security by Default** - Built-in protections
6. **Code Generation** - AST-based generation for type-safe code

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                          │
│  User Models, Views, Controllers, Routes                      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  Code Generation Layer                        │
│  AST Parser → Schema Analysis → Code Generator                │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    Framework API Layer                        │
│  Admin, API, ORM, Identity, Filter                            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    Database Layer                             │
│  SQL Builder, Query Execution, Transactions, Migrations       │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                         │
│  HTTP Server, Security, Logging, Config, Server               │
└─────────────────────────────────────────────────────────────┘
```

## Component Architecture

### 1. Schema Definition System

**Location:** `pkg/schema/`

Defines models declaratively in Go code:

```go
type User struct {
    schema.BaseSchema
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("username").Required().Unique().Build(),
    }
}
```

### 2. Code Generation System

**Location:** `pkg/generator/`

Generates type-safe Go code from schema definitions:
- Model structs
- FieldExpr definitions
- Manager with CRUD operations
- QuerySet wrappers

### 3. Query System

**Location:** `pkg/query/`

Type-safe and dynamic query building:
- `FieldExpr[T]` for type-safe field access
- `QueryExpr` for query conditions
- `BaseQuerySet[T]` for query execution
- SQL builder with proper escaping

### 4. Database Layer

**Location:** `pkg/db/`

Database connection, transactions, migrations:
- Connection pooling
- Transaction management
- Migration system integration
- SQL builder with parameter binding

### 5. HTTP & Routing

**Location:** `pkg/http/`

HTTP server and routing (chi wrapper):
- Router wrapper
- Middleware stack
- Request context utilities
- Server wrapper

### 6. Admin System

**Location:** `pkg/admin/`

Type-safe admin interface:
- Type-safe Admin[T] and Config[T]
- Complete HTTP handlers
- Rich form widgets
- Filters, search, pagination

### 7. REST API Framework

**Location:** `pkg/api/`

DRF-like API framework:
- Serializers
- ViewSets
- Authentication
- Permissions
- Throttling

## Data Flow

### Request Flow

```
HTTP Request
  ↓
Chi Router
  ↓
Middleware Stack
  ↓
Framework Handler
  ↓
QuerySet / Manager
  ↓
SQL Builder (with parameter binding)
  ↓
PostgreSQL
  ↓
Response (JSON/HTML)
```

### Code Generation Flow

```
Schema Definition (Go)
  ↓
AST Parser
  ↓
Model Definition
  ↓
Code Generator
  ↓
Generated Files:
  - Model structs
  - FieldExpr
  - Manager/QuerySet
```

## Design Patterns

### 1. Builder Pattern

Used extensively for field definitions:

```go
schema.String("username").
    Required().
    Unique().
    MaxLength(150).
    Build()
```

### 2. Strategy Pattern

Used for authentication, permissions, throttling:

```go
type Authentication interface {
    Authenticate(r *http.Request) (*AuthResult, error)
}
```

### 3. Repository Pattern

QuerySet/Manager abstracts data access:

```go
type QuerySet[T any] interface {
    Filter(expr QueryExpr) QuerySet[T]
    All(ctx context.Context) ([]*T, error)
}
```

### 4. Template Method Pattern

BaseViewSet defines skeleton, allows customization:

```go
type BaseViewSet struct {
    // Template methods
    GetQueryset() interface{}
    GetSerializer() Serializer
}
```

## Extension Points

### 1. Model Extensions
- Add fields to existing models
- Add relations
- Add hooks
- Override methods

### 2. Admin Extensions
- Customize list display
- Add filters
- Add actions
- Custom widgets

### 3. Query Extensions
- Custom QuerySet methods
- Custom aggregates
- Custom annotations

### 4. Middleware
- Custom middleware injection
- Request/response modification

### 5. Plugin System
- Register plugins
- Plugin lifecycle hooks
- Plugin dependencies

## Technology Stack

### Core
- **Go 1.21+** - Programming language
- **database/sql** - Database interface
- **go/ast** - Code generation

### HTTP & Routing
- **chi/v5** - HTTP router
- **chi/middleware** - Middleware

### Database
- **database/sql** - Standard library SQL interface
- **lib/pq** - PostgreSQL driver
- **golang-migrate/v4** - Migrations

### Security
- **gorilla/csrf** - CSRF protection
- **alexedwards/scs/v2** - Sessions
- **golang.org/x/crypto/bcrypt** - Password hashing

### Validation & Configuration
- **go-playground/validator/v10** - Validation
- **spf13/viper** - Configuration

### Logging & Templates
- **zap** - Structured logging
- **Masterminds/sprig/v3** - Template functions

## Design Principles

1. **Wrap, Don't Expose** - All third-party libraries are wrapped
2. **Type-Safe First** - Primary API uses generics
3. **Convention over Configuration** - Sensible defaults
4. **Extensibility** - Everything can be extended
5. **Security by Default** - Built-in protections
6. **Code Generation** - Reduce boilerplate

## Next Steps

For complete architecture documentation, see the [Architecture Deep Dive](/docs/deep-dives/architecture).

You may also want to explore:
- [Design Principles](/docs/deep-dives/design-principles) - Framework design principles
- [Features Overview](/docs/deep-dives/features-overview) - Complete feature list
- [API Architecture](/docs/learn/api-architecture) - REST API framework architecture
- [User System Architecture](/docs/learn/user-system-architecture) - User system design
