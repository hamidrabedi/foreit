# forge Framework - Architecture Documentation

## Overview

forge is a Django-like Go framework that prioritizes type safety while offering dynamic capabilities. It provides a full-stack web framework with ORM, admin interface, code generation, and extensibility.

## Table of Contents

- [Core Philosophy](#core-philosophy)
- [Architecture Layers](#architecture-layers)
- [Component Architecture](#component-architecture)
  - [Schema Definition System](#1-schema-definition-system)
  - [Code Generation System](#2-code-generation-system)
  - [Query System](#3-query-system)
  - [Database Layer](#4-database-layer)
  - [HTTP & Routing](#5-http--routing)
  - [Admin System](#6-admin-system)
  - [Security](#7-security)
  - [Validation](#8-validation)
  - [Authentication](#9-authentication)
  - [Configuration](#10-configuration)
  - [Logging](#11-logging)
  - [Utilities](#12-utilities)
- [Data Flow](#data-flow)
- [Extension Points](#extension-points)
- [Technology Stack](#technology-stack)
- [Design Principles](#design-principles)
- [Detailed Architecture](#detailed-architecture)
  - [QuerySet Architecture](#queryset-architecture)
  - [Schema Architecture Details](#schema-architecture-details)
  - [API Design Principles](#api-design-principles)

## Core Philosophy

1. **Type-Safe First**: Primary API uses generics for compile-time safety
2. **Dynamic When Needed**: Secondary API for runtime flexibility
3. **Convention over Configuration**: Sensible defaults everywhere
4. **Fully Extensible**: Everything can be extended/overridden
5. **Security by Default**: Built-in protections
6. **Code Generation**: AST-based generation for type-safe code

## Architecture Layers

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

## Component Architecture

### 1. Schema Definition System

**Location:** `pkg/schema/`

**Purpose:** Define models declaratively in Go code

**Components:**
- `schema.go` - Core Schema interface
- `field.go` - Field definitions with builders
- `relation.go` - Relationship definitions
- `meta.go` - Model metadata
- `hooks.go` - Lifecycle hooks
- `constraint_builder.go` - Database constraint builders

**Example:**
```go
type User struct {
    schema.Schema
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        fields.Int64("id").Primary().AutoIncrement(),
        fields.String("username").Unique().Required(),
        fields.String("email").Unique().Required(),
    }
}
```

### 2. Code Generation System

**Location:** `pkg/generator/`

**Purpose:** Generate type-safe Go code from schema definitions

**Components:**
- `generator.go` - Main generation orchestrator
- `ast_parser.go` - Go AST parser for schema extraction
- `writer.go` - Code file writer
- `templates.go` - Code generation templates
- `templates/` - Template files for Manager and QuerySet generation

**Generated Files:**
- `models/*.gen.go` - Model structs
- `models/*_fields.gen.go` - FieldExpr definitions
- `models/*_manager.gen.go` - Manager with CRUD operations ✅
- `models/*_queryset.gen.go` - Type-safe QuerySet wrapper ✅

**Generation Pattern:**
```go
// Generated UserQuerySet embeds BaseQuerySet
type UserQuerySet struct {
    *query.BaseQuerySet[User]
}

// Generated Manager provides CRUD operations
type UserManagerType struct {
    db *db.DB
}
```

### 3. Query System

**Location:** `pkg/query/`

**Purpose:** Type-safe and dynamic query building

**Components:**
- `field_expr.go` - Type-safe field accessors (`FieldExpr[T]`)
- `query_expr.go` - Query conditions (`QueryExpr`)
- `queryset.go` - QuerySet implementation with `BaseQuerySet[T]`
- `manager.go` - Manager with CRUD operations
- `dynamic.go` - Dynamic query API
- `aggregates.go` - Aggregate functions (structure ready)
- `annotations.go` - Annotation support (structure ready)

**Architecture:**
- `BaseQuerySet[T]` - Exported base implementation that generated QuerySets embed
- Generated QuerySets (e.g., `UserQuerySet`) embed `*query.BaseQuerySet[User]`
- All QuerySet methods delegate to embedded `BaseQuerySet` for consistency
- Type safety maintained through generics at compile time

**Usage:**
```go
// Type-safe
users, err := User.Objects.Filter(
    User.Fields.IsActive.Equals(true).And(
        User.Fields.DateJoined.Greater(lastMonth),
    ),
).All(ctx)

// Dynamic
users, err := User.Objects.FilterDynamic(
    query.Q("is_active", true),
).All(ctx)
```

### 4. Database Layer

**Location:** `pkg/db/`

**Purpose:** Database connection, transactions, migrations

**Components:**
- `db.go` - Database connection wrapper
- `transaction.go` - Transaction management
- `migrations.go` - Migration system integration

**Migration System:** `pkg/migrations/` and `pkg/migrate/`
- Complete migration system with state management
- Schema detection and diff generation
- SQL builder for migrations

**Features:**
- Connection pooling
- Transaction support with savepoints
- Migration management
- SQL builder with proper escaping and parameter binding

### 5. HTTP & Routing

**Location:** `pkg/http/`

**Purpose:** HTTP server and routing (chi wrapper)

**Components:**
- `router.go` - Chi router wrapper
- `middleware.go` - Middleware layer
- `context.go` - Request context utilities
- `server.go` - Server wrapper

**Middleware Stack:**
1. Request ID
2. Real IP
3. Recoverer
4. Logger (zap)
5. Session (scs)
6. CSRF (gorilla/csrf)
7. Authentication

### 6. Admin System

**Location:** `pkg/admin/`

**Purpose:** Type-safe admin interface with full CRUD operations

**Components:**
- `admin.go` - Type-safe Admin[T] and Config[T]
- `registry.go` - Admin model registry
- `list_view.go` - List view with pagination, search, filtering
- `detail_view.go` - Detail view
- `form_view.go` - Create/update forms
- `fields.go` - Type-safe field expressions
- `filters.go` - Filter system
- `actions.go` - Bulk actions
- `widgets.go` - Form widgets
- `export.go` - CSV/JSON export
- `http/` - HTTP handlers and routing

**Features:**
- ✅ Type-safe with generics
- ✅ Complete HTTP handlers (List, Detail, Create, Update, Delete)
- ✅ Rich form widgets
- ✅ Filters, search, pagination
- ✅ Bulk actions
- ✅ Export functionality
- ✅ Inlines and fieldsets

### 7. Security

**Location:** `pkg/security/`

**Purpose:** Security features

**Components:**
- `csrf.go` - CSRF protection (gorilla/csrf)
- `sessions.go` - Session management (scs)
- `xss.go` - XSS protection
- `sql_injection.go` - SQL injection prevention

### 8. Validation

**Location:** `pkg/validation/`

**Purpose:** Data validation

**Components:**
- `validator.go` - Validator wrapper (go-playground/validator)
- `tags.go` - Validation tag helpers
- `integration.go` - Schema integration

### 9. Authentication

**Location:** `pkg/auth/` and `pkg/users/`

**Purpose:** Authentication and authorization

**Components:**
- `pkg/auth/password.go` - Password hashing (bcrypt)
- `pkg/users/` - Complete user system with:
  - User management (CRUD)
  - Authentication services
  - Session management
  - Permission system (RBAC)
  - Password management
  - Authentication backends (password, token)

### 10. Configuration

**Location:** `pkg/config/`

**Purpose:** Application configuration

**Components:**
- `config.go` - Viper wrapper
- `settings.go` - Framework settings structure

### 11. Logging

**Location:** `pkg/logging/`

**Purpose:** Structured logging

**Components:**
- `logger.go` - Zap logger wrapper
- `middleware.go` - Request logging middleware

### 12. Utilities

**Location:** `pkg/utils/`

**Purpose:** Helper utilities

**Components:**
- `strcase.go` - String case conversion (strcase)
- `uuid.go` - UUID utilities (google/uuid)

### 13. API Framework

**Location:** `pkg/api/`

**Purpose:** DRF-like REST API framework

**Components:**
- `serializers/` - Complete serializer system with field types
- `viewset.go` - BaseViewSet with CRUD operations
- `authentication/` - Token, JWT, Basic, Session, API Key auth
- `permissions/` - Permission system (AllowAny, IsAuthenticated, IsAdminUser, etc.)
- `throttling/` - Rate limiting (AnonRateThrottle, UserRateThrottle, ScopedRateThrottle)
- `renderers/` - JSON, XML, YAML, HTML, CSV renderers
- `parsers/` - JSON, XML, Form, MultiPart parsers
- `filters/` - Field filtering and search
- `pagination.go` - PageNumber and LimitOffset pagination
- `exceptions/` - Complete exception hierarchy
- `versioning/` - API versioning support
- `caching/` - Cache backends
- `docs/` - OpenAPI documentation generation

**Status:** ✅ Complete - Production ready

### 14. User System

**Location:** `pkg/users/`

**Purpose:** Complete user management and authentication

**Components:**
- `models/` - User, Session, Permission, Group, Token models
- `repository/` - Data access layer (Repository pattern)
- `service/` - Business logic layer (User, Auth, Password, Permission services)
- `backends/` - Authentication backends (password, token)
- `serializers/` - API serializers
- `handlers/` - HTTP handlers/viewsets
- `middleware/` - Authentication middleware

**Status:** ✅ Complete - Production ready

### 15. CLI Tools

**Location:** `pkg/cli/`

**Purpose:** Command-line interface for framework operations

**Components:**
- `commands/` - CLI commands (new, generate, migrate, runserver, etc.)
- `templates/` - Project and code templates
- `root.go` - CLI root command

**Status:** ✅ Complete

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
  - SQL Builder generates SQL with parameter binding
  - Proper identifier escaping
  - Type-safe query execution
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

### Utilities
- **stretchr/testify** - Testing
- **iancoleman/strcase** - String utilities
- **google/uuid** - UUID generation

## Design Principles

1. **Wrap, Don't Expose**: All third-party libraries are wrapped
2. **Type-Safe First**: Primary API uses generics
3. **Convention over Configuration**: Sensible defaults
4. **Extensibility**: Everything can be extended
5. **Security by Default**: Built-in protections
6. **Code Generation**: Reduce boilerplate

## Detailed Architecture

### QuerySet Architecture

The QuerySet system provides a type-safe, Django-like query interface using Go generics.

**BaseQuerySet Structure:**
```go
type BaseQuerySet[T any] struct {
    table           string
    conditions      []QueryExpr
    excludes        []QueryExpr
    orderBy         []string
    limitVal        *int
    offsetVal       *int
    distinct        bool
    selectFields    []string
    selectRelated   []string
    prefetchRelated []string
    onlyFields      []string
    deferFields     []string
    aggregates      []Aggregate
    annotations     []AnnotationExpr
    db              interface{}  // Database connection
}
```

**Key Design Decisions:**
1. **Exported Type**: `BaseQuerySet` is exported so generated code can embed it
2. **Generic Type Parameter**: `[T any]` ensures type safety for model instances
3. **Composition**: Generated QuerySets embed `*BaseQuerySet[T]` for all functionality
4. **Database Connection**: QuerySets can hold DB connection or get from context

**Generated QuerySet Pattern:**
```go
type UserQuerySet struct {
    *query.BaseQuerySet[User]
}

func NewUserQuerySet() *UserQuerySet {
    base := query.NewBaseQuerySet[User]("users")
    return &UserQuerySet{BaseQuerySet: base}
}
```

**Query Execution Flow:**
1. Query construction via method chaining
2. `buildSQL()` constructs SQL from conditions using SQL builder
3. SQL executed via `database/sql` with proper parameter binding
4. Rows scanned into model instances using reflection
5. Relations loaded if `SelectRelated`/`PrefetchRelated` used

**Implementation Status:**
- ✅ Query building (Filter, Exclude, OrderBy, Limit, Offset, Distinct)
- ✅ Query execution (All, Get, First, Last, Count, Exists)
- ✅ SQL generation from QueryExpr
- ✅ Row scanning into model instances
- ✅ SQL builder with proper escaping and parameter binding
- ✅ Manager CRUD operations (Create, Update, Delete with hooks)
- 🚧 SelectRelated/PrefetchRelated (structure ready)
- 🚧 Aggregates execution (structure ready)
- 🚧 Annotations execution (structure ready)

### Schema Architecture Details

**Field Structure:**
```go
type Field struct {
    Name          string
    Type          FieldType
    Required      bool
    Blank         bool
    Default       interface{}
    HelpText      string
    VerboseName   string
    DBColumn      string
    DBIndex       bool
    Unique        bool
    PrimaryKey    bool
    AutoIncrement bool
    Validators    []Validator
    ValidationTag string
    MaxLength     *int
    MinLength     *int
    Choices       []Choice
    Editable      bool
    AutoNow       bool
    AutoNowAdd    bool
}
```

**Relation Structure:**
```go
type Relation struct {
    Name        string
    Type        RelationType  // ForeignKey, OneToOne, OneToMany, ManyToMany
    Target      string
    Required    bool
    OnDelete    CascadeType
    OnUpdate    CascadeType
    RelatedName string
    Through     string  // For ManyToMany
}
```

**Meta Structure:**
```go
type Meta struct {
    TableName          string
    OrderBy            []string
    VerboseName        string
    VerboseNamePlural  string
    Indexes            []Index
    UniqueTogether     [][]string
    Constraints        []UniqueConstraint
    Permissions        []Permission
    DefaultPermissions bool
}
```

**Hook Execution Order:**

Create: `BeforeSave` → `BeforeCreate` → Database Insert → `AfterCreate` → `AfterSave`

Update: `BeforeSave` → `BeforeUpdate` → Database Update → `AfterUpdate` → `AfterSave`

Delete: `BeforeDelete` → Database Delete → `AfterDelete`

### API Design Principles

1. **Type-Safe First**: Primary API uses generics for compile-time safety
2. **Django-Inspired**: Familiar patterns for Django developers
3. **Composable**: APIs can be combined
4. **Extensible**: Everything can be extended
5. **Convention over Configuration**: Sensible defaults

**API Patterns:**
- **Fluent Interface**: Method chaining for queries
- **Builder Pattern**: Chainable field builders
- **Factory Pattern**: Router, server creation
- **Registry Pattern**: Model and extension registration

**Database Schema Mapping:**

| Go Type | PostgreSQL Type | Nullable |
|---------|----------------|----------|
| `int64` | `BIGINT` | No |
| `*int64` | `BIGINT` | Yes |
| `int32` | `INTEGER` | No |
| `string` | `TEXT` or `VARCHAR(n)` | No |
| `*string` | `TEXT` or `VARCHAR(n)` | Yes |
| `bool` | `BOOLEAN` | No |
| `time.Time` | `TIMESTAMP` | No |
| `*time.Time` | `TIMESTAMP` | Yes |
| `float64` | `DOUBLE PRECISION` | No |
| `[]byte` | `BYTEA` or `JSONB` | No |

## Package Structure Summary

The framework is organized into the following packages under `pkg/`:

```
pkg/
├── admin/          # Type-safe admin interface ✅
├── api/            # REST API framework ✅
├── auth/           # Basic authentication utilities
├── cli/            # Command-line interface ✅
├── config/         # Configuration management
├── db/             # Database connection and transactions
├── errors/         # Error handling
├── generator/      # Code generation system ✅
├── http/           # HTTP routing and middleware
├── logging/        # Structured logging
├── migrate/        # Migration system ✅
├── migrations/     # Migration utilities
├── models/         # Base models
├── query/          # Type-safe ORM ✅
├── registry/       # Extension registry
├── schema/         # Schema definition system ✅
├── security/       # Security features
├── users/          # User system ✅
├── utils/          # Utility functions
└── validation/     # Validation system
```

**Status Legend:**
- ✅ Complete and production-ready
- 🚧 In progress / structure ready
- 📋 Planned

For detailed architecture of specific packages, see:
- [API Architecture](API_ARCHITECTURE.md) - API framework details
- [User System Architecture](USER_SYSTEM_ARCHITECTURE.md) - User system details
