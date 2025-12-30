# forge Framework - Complete Architecture Documentation

This document provides a comprehensive, in-depth architecture guide for the forge framework, covering all design patterns, implementation details, features, and architectural decisions.

**Document Version:** 1.0  
**Last Updated:** January 2025  
**Total Sections:** 50+  
**Target Length:** 3000+ lines

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Core Philosophy & Design Principles](#core-philosophy--design-principles)
3. [Architecture Overview](#architecture-overview)
4. [System Architecture Layers](#system-architecture-layers)
5. [Component Architecture](#component-architecture)
6. [Design Patterns](#design-patterns)
7. [Type System & Generics](#type-system--generics)
8. [Schema Definition System](#schema-definition-system)
9. [Code Generation System](#code-generation-system)
10. [Query System Architecture](#query-system-architecture)
11. [Database Layer Architecture](#database-layer-architecture)
12. [Migration System Architecture](#migration-system-architecture)
13. [Admin System Architecture](#admin-system-architecture)
14. [API Framework Architecture](#api-framework-architecture)
15. [User System Architecture](#user-system-architecture)
16. [HTTP & Routing Architecture](#http--routing-architecture)
17. [Security Architecture](#security-architecture)
18. [Validation Architecture](#validation-architecture)
19. [Configuration Architecture](#configuration-architecture)
20. [Logging Architecture](#logging-architecture)
21. [CLI Architecture](#cli-architecture)
22. [Extension & Plugin System](#extension--plugin-system)
23. [Data Flow & Request Processing](#data-flow--request-processing)
24. [Error Handling Architecture](#error-handling-architecture)
25. [Testing Architecture](#testing-architecture)
26. [Performance Considerations](#performance-considerations)
27. [Concurrency & Thread Safety](#concurrency--thread-safety)
28. [Memory Management](#memory-management)
29. [Database Query Optimization](#database-query-optimization)
30. [Caching Strategy](#caching-strategy)
31. [API Design Principles](#api-design-principles)
32. [Type Safety Mechanisms](#type-safety-mechanisms)
33. [Code Generation Details](#code-generation-details)
34. [SQL Builder Architecture](#sql-builder-architecture)
35. [Transaction Management](#transaction-management)
36. [Middleware System](#middleware-system)
37. [Authentication Flow](#authentication-flow)
38. [Permission System](#permission-system)
39. [Serialization System](#serialization-system)
40. [Content Negotiation](#content-negotiation)
41. [Pagination System](#pagination-system)
42. [Filtering System](#filtering-system)
43. [Throttling System](#throttling-system)
44. [Versioning System](#versioning-system)
45. [Exception Handling](#exception-handling)
46. [State Management](#state-management)
47. [Dependency Injection](#dependency-injection)
48. [Lifecycle Management](#lifecycle-management)
49. [Extension Points](#extension-points)
50. [Best Practices & Guidelines](#best-practices--guidelines)
51. [Future Architecture Considerations](#future-architecture-considerations)

---

## Executive Summary

forge is a comprehensive, Django-inspired Go web framework that combines the best of both worlds: Django's developer-friendly conventions and Go's performance and type safety. The framework provides:

- **Type-Safe ORM**: Full-featured ORM with compile-time type safety using Go generics
- **Code Generation**: AST-based code generation for zero-boilerplate model definitions
- **Admin Interface**: Complete admin system with type-safe configuration
- **REST API Framework**: DRF-like API framework with serializers, viewsets, authentication, and permissions
- **User System**: Complete user management with authentication, sessions, and permissions
- **Migration System**: Automatic migration generation and execution
- **Security**: Built-in CSRF, XSS, SQL injection protection
- **CLI Tools**: Comprehensive command-line interface for development

The framework is designed with extensibility, type safety, and developer experience as core principles.

---

## Core Philosophy & Design Principles

### 1. Type-Safe First

**Principle:** Primary APIs use Go generics for compile-time type safety.

**Implementation:**
- All QuerySet operations are type-safe: `QuerySet[User]` ensures only `User` instances are returned
- Field expressions are type-safe: `FieldExpr[User, string]` ensures correct field types
- Admin system uses generics: `Admin[User]` and `Config[User]` for type-safe configuration
- Manager operations are type-safe: `Manager[User]` ensures type-safe CRUD operations

**Benefits:**
- Compile-time error detection
- IDE autocomplete and refactoring support
- No runtime type errors
- Better documentation through types

**Example:**
```go
// Type-safe QuerySet
users, err := User.Objects.Filter(
    User.Fields.Email.Equals("user@example.com"),
).All(ctx)

// Compiler ensures 'users' is []*User, not []*Post
// Compiler ensures Email field is string, not int
```

### 2. Dynamic When Needed

**Principle:** Secondary APIs provide runtime flexibility for edge cases.

**Implementation:**
- Dynamic query API: `FilterDynamic()` for runtime query building
- Reflection-based serialization for unknown types
- Runtime schema inspection
- Dynamic admin registration

**Benefits:**
- Flexibility for complex scenarios
- Support for dynamic data structures
- Runtime introspection capabilities

**Example:**
```go
// Dynamic query when field names are determined at runtime
users, err := User.Objects.FilterDynamic(
    query.Q("email", "user@example.com"),
    query.Q("is_active", true),
).All(ctx)
```

### 3. Convention over Configuration

**Principle:** Sensible defaults everywhere, minimal configuration required.

**Implementation:**
- Default table names from model names
- Auto-generated field names
- Default admin configuration
- Convention-based routing
- Default middleware stack

**Benefits:**
- Faster development
- Less boilerplate
- Consistent patterns
- Easy onboarding

**Example:**
```go
// Minimal configuration - defaults work
admin.Register(&User{}, userManager, nil)

// vs explicit configuration
admin.Register(&User{}, userManager, &admin.Config[User]{
    ListDisplay: []admin.FieldExpr[User, interface{}]{...},
    // ... many options
})
```

### 4. Fully Extensible

**Principle:** Everything can be extended, overridden, or customized.

**Implementation:**
- Interface-based design
- Strategy pattern for pluggable components
- Hook system for lifecycle events
- Plugin system for extensions
- Customizable templates

**Benefits:**
- Adapt to any use case
- No framework limitations
- Easy to customize
- Support for complex requirements

**Example:**
```go
// Custom authentication backend
type CustomAuthBackend struct{}

func (b *CustomAuthBackend) Authenticate(ctx context.Context, creds map[string]string) (*User, error) {
    // Custom authentication logic
}

// Register custom backend
registry.RegisterBackend("custom", &CustomAuthBackend{})
```

### 5. Security by Default

**Principle:** Built-in security protections, no opt-in required.

**Implementation:**
- CSRF protection enabled by default
- SQL injection prevention (parameter binding)
- XSS protection in templates
- Secure session management
- Password hashing (bcrypt)
- Input validation

**Benefits:**
- Secure out of the box
- No security oversights
- Best practices enforced
- Reduced attack surface

### 6. Code Generation

**Principle:** AST-based code generation reduces boilerplate while maintaining type safety.

**Implementation:**
- Parse Go AST to extract schema definitions
- Generate type-safe Manager and QuerySet code
- Generate FieldExpr definitions
- Template-based code generation
- Incremental generation support

**Benefits:**
- Zero boilerplate for common operations
- Type safety maintained
- Performance (no reflection overhead)
- IDE support for generated code

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Application                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Models      │  │   Admin       │  │   API         │         │
│  │  (Schema)     │  │   Config      │  │   ViewSets    │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Code Generation Layer                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  AST Parser │→ │   Generator   │→ │  Templates    │         │
│  │  (go/ast)   │  │  (Manager/    │  │  (Code Gen)   │         │
│  │             │  │   QuerySet)   │  │               │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Framework API Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  QuerySet    │  │   Manager     │  │  FieldExpr    │         │
│  │  [T]         │  │   [T]         │  │  [T, F]       │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  QueryExpr   │  │   SQL Builder │  │  Migrations   │         │
│  │              │  │   (Safe)      │  │              │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Database Layer                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Connection  │  │  Transaction  │  │  SQL Exec    │         │
│  │  Pool        │  │  Management   │  │  (database/  │         │
│  │              │  │               │  │   sql)       │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  HTTP Router │  │  Middleware  │  │  Security     │         │
│  │  (chi)       │  │  Stack       │  │  (CSRF/XSS)   │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Config      │  │  Logging     │  │  Validation   │         │
│  │  (Viper)     │  │  (zap)       │  │  (validator)  │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

### Component Interaction Flow

```
1. User defines schema in Go code
   ↓
2. Code generator parses AST and generates Manager/QuerySet
   ↓
3. User creates Admin or API ViewSet
   ↓
4. HTTP request arrives
   ↓
5. Middleware stack processes request
   ↓
6. Handler extracts QuerySet from context
   ↓
7. QuerySet builds SQL query with parameter binding
   ↓
8. SQL executed against database
   ↓
9. Results serialized (JSON/HTML)
   ↓
10. Response sent to client
```

---

## System Architecture Layers

### Layer 1: User Application Layer

**Responsibility:** User-defined models, admin configurations, API viewsets, handlers.

**Components:**
- Schema definitions (models implementing `schema.Schema`)
- Admin configurations (`admin.Config[T]`)
- API viewsets (`api.ViewSet`)
- HTTP handlers
- Business logic

**Characteristics:**
- Framework-agnostic business logic
- Uses framework APIs
- No direct database access
- Type-safe operations

### Layer 2: Code Generation Layer

**Responsibility:** Generate type-safe code from schema definitions.

**Components:**
- AST parser (`generator.ASTParser`)
- Code generator (`generator.Generator`)
- Template engine
- Writer (`generator.Writer`)

**Process:**
1. Scan Go files for schema definitions
2. Parse AST to extract model information
3. Generate Manager code with CRUD operations
4. Generate QuerySet code with query methods
5. Generate FieldExpr definitions
6. Write generated files

### Layer 3: Framework API Layer

**Responsibility:** Type-safe ORM, query building, admin system, API framework.

**Components:**
- QuerySet system (`query.QuerySet[T]`)
- Manager system (`query.Manager[T]`)
- FieldExpr system (`query.FieldExpr[T, F]`)
- Admin system (`admin.Admin[T]`)
- API framework (`api.ViewSet`)
- User system (`users.*`)

**Characteristics:**
- Type-safe with generics
- Composable APIs
- Extensible design
- Performance-optimized

### Layer 4: Database Layer

**Responsibility:** Database connections, transactions, SQL execution.

**Components:**
- Connection pool (`db.DB`)
- Transaction management (`db.Tx`)
- SQL builder (`query.SQLBuilder`)
- Migration system (`migrations.*`)

**Characteristics:**
- Connection pooling
- Transaction support
- Parameter binding
- SQL injection prevention

### Layer 5: Infrastructure Layer

**Responsibility:** HTTP routing, middleware, security, configuration, logging.

**Components:**
- HTTP router (`http.Router`)
- Middleware stack
- Security (`security.*`)
- Configuration (`config.Config`)
- Logging (`logging.Logger`)

**Characteristics:**
- Framework-agnostic where possible
- Pluggable components
- Sensible defaults
- Production-ready

---

## Component Architecture

### Schema Definition System

**Location:** `pkg/schema/`

**Purpose:** Declarative model definitions with Django-like field options.

#### Core Components

**1. Schema Interface**
```go
type Schema interface {
    Fields() []Field
    Relations() []Relation
    Meta() Meta
    Hooks() *ModelHooks
}
```

**2. Field System**
- **Field Types:** String, Int64, Int32, Float64, Bool, Time, JSON, UUID, etc.
- **Field Options:** Required, Unique, PrimaryKey, AutoIncrement, Default, MaxLength, Choices, etc.
- **Field Builders:** Fluent API for field construction
- **Field Traits:** Mixins for common field behaviors

**3. Relation System**
- **ForeignKey:** Many-to-one relationships
- **OneToOne:** One-to-one relationships
- **ManyToMany:** Many-to-many relationships
- **Cascade Options:** CASCADE, SET_NULL, SET_DEFAULT, RESTRICT, NO_ACTION

**4. Meta System**
- Table name customization
- Ordering defaults
- Indexes
- Unique constraints
- Permissions
- Verbose names

**5. Hooks System**
- **BeforeSave:** Before any save operation
- **AfterSave:** After any save operation
- **BeforeCreate:** Before insert
- **AfterCreate:** After insert
- **BeforeUpdate:** Before update
- **AfterUpdate:** After update
- **BeforeDelete:** Before delete
- **AfterDelete:** After delete

#### Field Builder API

```go
// Type-specific builders
fields.String("username").Required().Unique().MaxLength(150)
fields.Int64("age").Default(0).Min(0).Max(150)
fields.Bool("is_active").Default(true)
fields.Time("created_at").AutoNowAdd()
fields.JSON("metadata").Default(map[string]interface{}{})

// Relations
fields.ForeignKey("author", "User").Required().OnDelete(schema.CASCADE)
fields.ManyToMany("tags", "Tag").Through("PostTag")
```

#### Schema Registry

**Purpose:** Discover and register all models in the application.

**Features:**
- Automatic model discovery
- Schema validation
- Dependency resolution
- Extension point registration

---

### Code Generation System

**Location:** `pkg/generator/`

**Purpose:** Generate type-safe Go code from schema definitions.

#### Architecture

**1. AST Parser (`ast_parser.go`)**

**Responsibilities:**
- Parse Go source files
- Extract schema definitions
- Identify model structs
- Extract field information
- Extract relation information
- Extract meta options
- Extract hooks

**Process:**
```go
1. Load Go source files
2. Parse with go/parser
3. Walk AST with go/ast
4. Identify Schema implementations
5. Extract Fields() method
6. Extract Relations() method
7. Extract Meta() method
8. Extract Hooks() method
9. Build ModelDefinition
```

**2. Code Generator (`generator.go`)**

**Responsibilities:**
- Orchestrate generation process
- Manage template rendering
- Write generated files
- Handle incremental generation

**3. Template System (`templates/`)**

**Templates:**
- `manager.tmpl` - Manager generation
- `queryset.tmpl` - QuerySet generation
- `fields.tmpl` - FieldExpr generation
- `model.tmpl` - Model struct generation

**Template Features:**
- Go template syntax
- Sprig functions
- Type-safe code generation
- Proper imports
- Documentation comments

**4. Writer (`writer.go`)**

**Responsibilities:**
- Write generated files
- Format code with gofmt
- Handle file conflicts
- Backup existing files

#### Generated Code Structure

**Manager Generation:**
```go
// Generated: models/user_manager.gen.go
type UserManagerType struct {
    db *db.DB
    table string
}

func (m *UserManagerType) Create(ctx context.Context, user *User) error {
    // Generated CRUD with hooks
}

func (m *UserManagerType) Update(ctx context.Context, user *User) error {
    // Generated update with hooks
}

func (m *UserManagerType) Delete(ctx context.Context, id int64) error {
    // Generated delete with hooks
}
```

**QuerySet Generation:**
```go
// Generated: models/user_queryset.gen.go
type UserQuerySet struct {
    *query.BaseQuerySet[User]
}

func NewUserQuerySet() *UserQuerySet {
    base := query.NewBaseQuerySet[User]("users")
    return &UserQuerySet{BaseQuerySet: base}
}

// Type-safe methods
func (qs *UserQuerySet) FilterByEmail(email string) *UserQuerySet {
    return qs.Filter(User.Fields.Email.Equals(email))
}
```

**FieldExpr Generation:**
```go
// Generated: models/user_fields.gen.go
type UserFields struct {
    ID       query.FieldExpr[User, int64]
    Username query.FieldExpr[User, string]
    Email    query.FieldExpr[User, string]
    // ...
}

var User = struct {
    Objects *UserManagerType
    Fields  UserFields
}{
    Objects: &UserManagerType{},
    Fields: UserFields{
        ID:       query.NewFieldExpr[User, int64]("id"),
        Username: query.NewFieldExpr[User, string]("username"),
        Email:    query.NewFieldExpr[User, string]("email"),
    },
}
```

#### Generation Process

```
1. User runs: forge generate
   ↓
2. Generator scans models/ directory
   ↓
3. AST Parser extracts schema definitions
   ↓
4. Generator processes each model:
   a. Generate Manager code
   b. Generate QuerySet code
   c. Generate FieldExpr code
   d. Generate model struct (if needed)
   ↓
5. Writer writes generated files
   ↓
6. Code formatted with gofmt
   ↓
7. User imports and uses generated code
```

---

### Query System Architecture

**Location:** `pkg/query/`

**Purpose:** Type-safe query building and execution with Django-like API.

#### Core Components

**1. QuerySet Interface**

```go
type QuerySet[T any] interface {
    // Filtering
    Filter(expr QueryExpr) QuerySet[T]
    Exclude(expr QueryExpr) QuerySet[T]
    
    // Ordering
    OrderBy(fields ...string) QuerySet[T]
    Reverse() QuerySet[T]
    
    // Limiting
    Limit(n int) QuerySet[T]
    Offset(n int) QuerySet[T]
    Distinct() QuerySet[T]
    
    // Selection
    Select(fields ...string) QuerySet[T]
    Only(fields ...string) QuerySet[T]
    Defer(fields ...string) QuerySet[T]
    
    // Relations
    SelectRelated(relations ...string) QuerySet[T]
    PrefetchRelated(relations ...string) QuerySet[T]
    
    // Aggregation
    Aggregate(aggregates ...Aggregate) QuerySet[T]
    Annotate(annotations ...AnnotationExpr) QuerySet[T]
    
    // Execution
    All(ctx context.Context) ([]*T, error)
    Get(ctx context.Context) (*T, error)
    First(ctx context.Context) (*T, error)
    Last(ctx context.Context) (*T, error)
    Count(ctx context.Context) (int64, error)
    Exists(ctx context.Context) (bool, error)
    
    // Updates
    Update(ctx context.Context, updates map[string]interface{}) (int64, error)
    Delete(ctx context.Context) (int64, error)
}
```

**2. BaseQuerySet Implementation**

**Structure:**
```go
type BaseQuerySet[T any] struct {
    table           string
    schema          *ModelSchema
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
    db              interface{}
    mu              sync.RWMutex
}
```

**Key Design Decisions:**

1. **Exported Type:** `BaseQuerySet` is exported so generated QuerySets can embed it
2. **Generic Type Parameter:** `[T any]` ensures type safety
3. **Immutable Operations:** Each method returns a new QuerySet (chainable)
4. **Lazy Evaluation:** Queries built but not executed until `All()`, `Get()`, etc.
5. **Thread Safety:** Mutex protection for concurrent access

**3. FieldExpr System**

**Purpose:** Type-safe field access and query building.

**Interface:**
```go
type FieldExpr[T any, F any] interface {
    // Get field value from instance
    Get(instance *T) F
    
    // Set field value on instance
    Set(instance *T, value F)
    
    // Query expressions
    Equals(value F) QueryExpr
    NotEquals(value F) QueryExpr
    GreaterThan(value F) QueryExpr
    GreaterThanOrEqual(value F) QueryExpr
    LessThan(value F) QueryExpr
    LessThanOrEqual(value F) QueryExpr
    In(values []F) QueryExpr
    NotIn(values []F) QueryExpr
    IsNull() QueryExpr
    IsNotNull() QueryExpr
    Like(pattern string) QueryExpr
    ILike(pattern string) QueryExpr
    Contains(value string) QueryExpr
    StartsWith(value string) QueryExpr
    EndsWith(value string) QueryExpr
    
    // Logical operations
    And(other QueryExpr) QueryExpr
    Or(other QueryExpr) QueryExpr
    Not() QueryExpr
}
```

**Implementation:**
```go
type fieldExpr[T any, F any] struct {
    name       string
    columnName string
    getter     func(*T) F
    setter     func(*T, F)
    fieldType  reflect.Type
}
```

**4. QueryExpr System**

**Purpose:** Represent query conditions in a composable way.

**Interface:**
```go
type QueryExpr interface {
    // Build SQL condition
    BuildSQL(builder *SQLBuilder) (string, []interface{}, error)
    
    // Logical operations
    And(other QueryExpr) QueryExpr
    Or(other QueryExpr) QueryExpr
    Not() QueryExpr
}
```

**Implementations:**
- `FieldExprQuery` - Field-based queries
- `ComparisonExpr` - Comparison operations
- `LogicalExpr` - AND/OR/NOT operations
- `RawExpr` - Raw SQL expressions
- `SubqueryExpr` - Subquery expressions

**5. SQL Builder**

**Purpose:** Generate safe SQL with proper escaping and parameter binding.

**Features:**
- Identifier escaping (table/column names)
- Parameter binding (prevents SQL injection)
- Dialect support (PostgreSQL, SQLite)
- Query optimization
- Query caching

**Example:**
```go
builder := NewSQLBuilder("postgres")
sql, args, err := builder.Select("users").
    Where("email = ?", "user@example.com").
    OrderBy("created_at DESC").
    Limit(10).
    Build()

// Generates:
// SELECT * FROM "users" WHERE email = $1 ORDER BY created_at DESC LIMIT $2
// Args: ["user@example.com", 10]
```

**6. Manager System**

**Purpose:** Type-safe CRUD operations with lifecycle hooks.

**Interface:**
```go
type Manager[T any] interface {
    // Create
    Create(ctx context.Context, instance *T) error
    
    // Read
    Get(ctx context.Context, id int64) (*T, error)
    GetBy(ctx context.Context, field string, value interface{}) (*T, error)
    All(ctx context.Context) ([]*T, error)
    
    // Update
    Update(ctx context.Context, instance *T) error
    UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error
    
    // Delete
    Delete(ctx context.Context, id int64) error
    DeleteBy(ctx context.Context, field string, value interface{}) error
    
    // QuerySet
    Objects() QuerySet[T]
}
```

**Hook Execution:**
```
Create:
  BeforeSave → BeforeCreate → SQL Insert → AfterCreate → AfterSave

Update:
  BeforeSave → BeforeUpdate → SQL Update → AfterUpdate → AfterSave

Delete:
  BeforeDelete → SQL Delete → AfterDelete
```

#### Query Execution Flow

```
1. User builds query:
   User.Objects.Filter(User.Fields.Email.Equals("user@example.com"))
   ↓
2. QuerySet builds QueryExpr tree
   ↓
3. User calls execution method: .All(ctx)
   ↓
4. QuerySet builds SQL using SQLBuilder
   ↓
5. SQL executed with database/sql
   ↓
6. Rows scanned into model instances
   ↓
7. Relations loaded if SelectRelated/PrefetchRelated used
   ↓
8. Results returned to user
```

#### Performance Optimizations

1. **Query Caching:** Cache compiled SQL queries
2. **Connection Pooling:** Reuse database connections
3. **Lazy Loading:** Load relations only when needed
4. **Batch Operations:** Bulk insert/update support
5. **Query Optimization:** Optimize query plans
6. **Field Selection:** Only select needed fields

---

### Database Layer Architecture

**Location:** `pkg/db/`

**Purpose:** Database connection management, transactions, and SQL execution.

#### Components

**1. Database Connection (`db.go`)**

**Structure:**
```go
type DB struct {
    *sql.DB
    dsn string
    config *Config
}

type Config struct {
    DSN             string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}
```

**Features:**
- Connection pooling
- Connection lifecycle management
- Health checks
- Query logging
- Metrics collection

**2. Transaction Management (`transaction.go`)**

**Features:**
- Savepoint support
- Nested transactions
- Transaction context
- Rollback on error
- Commit/rollback hooks

**Usage:**
```go
err := db.WithTransaction(ctx, func(tx *db.Tx) error {
    // Multiple operations in transaction
    user := &User{Email: "user@example.com"}
    if err := userManager.Create(ctx, user); err != nil {
        return err
    }
    // Transaction commits automatically
    return nil
})
```

**3. Migration Integration (`migrations.go`)**

**Features:**
- Migration table management
- Migration execution
- Migration rollback
- Migration status
- Migration validation

---

### Migration System Architecture

**Location:** `pkg/migrations/` and `pkg/migrate/`

**Purpose:** Automatic migration generation and execution.

#### Architecture Overview

```
Model Definitions (Go)
    ↓
AST Parser
    ↓
Model Definitions (Structured)
    ↓
State Manager (Load Previous State)
    ↓
Change Detector (Compare States)
    ↓
Changes[] (Diff)
    ↓
SQL Builder (Generate SQL)
    ↓
Migration Files (*.up.sql, *.down.sql)
    ↓
Migration Runner (Execute)
    ↓
Database
```

#### Core Components

**1. State Management**

**Purpose:** Track database schema state for incremental migrations.

**Components:**
- `StateManager` - Interface for state management
- `SchemaState` - In-memory schema representation
- `FileStateLoader` - Load state from migration files
- `InMemoryState` - In-memory state implementation

**State Structure:**
```go
type SchemaState struct {
    Tables map[string]*TableState
}

type TableState struct {
    Name        string
    Columns     map[string]*ColumnState
    Indexes     map[string]*IndexState
    ForeignKeys map[string]*ForeignKeyState
    Constraints map[string]*ConstraintState
}
```

**2. Change Detection**

**Purpose:** Detect differences between current and previous schema state.

**Components:**
- `Detector` - Main change detector
- `ColumnDetector` - Column change detection
- `IndexDetector` - Index change detection
- `ForeignKeyDetector` - Foreign key change detection
- `ConstraintDetector` - Constraint change detection

**Change Types:**
- CreateTable
- DropTable
- RenameTable
- AddColumn
- DropColumn
- ModifyColumn
- RenameColumn
- AddIndex
- DropIndex
- ModifyIndex
- AddForeignKey
- DropForeignKey
- AddConstraint
- DropConstraint

**3. SQL Generation**

**Purpose:** Generate database-agnostic SQL from changes.

**Components:**
- `SQLBuilder` - Base SQL builder
- `PostgreSQLBuilder` - PostgreSQL-specific SQL
- `SQLiteBuilder` - SQLite-specific SQL
- `Dialect` - Database dialect interface

**4. Migration Execution**

**Purpose:** Execute migrations safely with validation.

**Components:**
- `Runner` - Migration execution
- `Validator` - Migration validation
- `Status` - Migration status tracking
- `Recovery` - Migration recovery

**Features:**
- Checksum validation
- Transactional migrations
- Rollback support
- Dry-run mode
- Migration squashing

---

### Admin System Architecture

**Location:** `pkg/admin/`

**Purpose:** Type-safe admin interface with full CRUD operations.

#### Architecture

```
Admin[T] (Type-Safe Admin Instance)
    ↓
Config[T] (Type-Safe Configuration)
    ↓
Registry (Model Registration)
    ↓
HTTP Handlers (List, Detail, Create, Update, Delete)
    ↓
Views (ListView, DetailView, FormView)
    ↓
Manager[T] (Data Access)
    ↓
Database
```

#### Core Components

**1. Admin Type**

```go
type Admin[T any] struct {
    model    T
    manager  *query.Manager[T]
    config   *Config[T]
    registry *Registry
    name     string
}
```

**2. Config Type**

```go
type Config[T any] struct {
    // List view
    ListDisplay   []FieldExpr[T, interface{}]
    ListFilter    []Filter[T]
    SearchFields  []FieldExpr[T, interface{}]
    DateHierarchy FieldExpr[T, interface{}]
    Ordering      []Ordering[T]
    ListPerPage   int
    
    // Form view
    Fieldsets        []Fieldset[T]
    ReadOnlyFields   []FieldExpr[T, interface{}]
    AutocompleteFields []FieldExpr[T, interface{}]
    RawIDFields      []FieldExpr[T, interface{}]
    
    // Actions
    Actions []Action[T]
    
    // Inlines
    Inlines []Inline[T, interface{}]
    
    // Customization
    VerboseName       string
    VerboseNamePlural string
    
    // Custom methods
    GetQueryset func(ctx context.Context, manager *query.Manager[T]) (query.QuerySet[T], error)
    SaveModel   func(ctx context.Context, instance *T, form FormData, isNew bool) error
    DeleteModel func(ctx context.Context, instance *T) error
}
```

**3. Views**

**ListView:**
- Pagination
- Search
- Filtering
- Ordering
- Bulk actions
- Export

**DetailView:**
- Field display
- Inlines
- Actions
- History

**FormView:**
- Field rendering
- Validation
- Widgets
- Fieldsets
- Inlines
- Save options

**4. HTTP Handlers**

**Routes:**
- `GET /admin/` - Admin index
- `GET /admin/{model}/` - List view
- `GET /admin/{model}/new/` - Create form
- `POST /admin/{model}/new/` - Create handler
- `GET /admin/{model}/{id}/` - Detail view
- `GET /admin/{model}/{id}/edit/` - Update form
- `POST /admin/{model}/{id}/edit/` - Update handler
- `POST /admin/{model}/{id}/delete/` - Delete handler
- `GET /admin/{model}/export/` - Export handler
- `GET /admin/{model}/autocomplete/` - Autocomplete handler

**5. Widgets**

**Types:**
- TextInput
- TextArea
- Select
- Checkbox
- Radio
- DatePicker
- DateTimePicker
- FileUpload
- ImageUpload
- RichTextEditor

**6. Filters**

**Types:**
- BooleanFilter
- ChoiceFilter
- DateFilter
- DateTimeFilter
- NumberFilter
- TextFilter

**7. Actions**

**Types:**
- BulkDelete
- BulkUpdate
- Custom actions

---

### API Framework Architecture

**Location:** `pkg/api/`

**Purpose:** DRF-like REST API framework with serializers, viewsets, authentication, and permissions.

#### Architecture

```
HTTP Request
    ↓
Router (chi)
    ↓
Middleware Stack
    ↓
Authentication
    ↓
Permissions
    ↓
Throttling
    ↓
ViewSet
    ↓
Serializer
    ↓
QuerySet
    ↓
Response Renderer
    ↓
HTTP Response
```

#### Core Components

**1. ViewSet System**

**BaseViewSet:**
```go
type BaseViewSet struct {
    Serializer func() Serializer
    Queryset   interface{}
    Model      interface{}
}
```

**EnhancedBaseViewSet:**
```go
type EnhancedBaseViewSet struct {
    BaseViewSet
    AuthenticationClasses []authentication.Authentication
    PermissionClasses     []permissions.Permission
    ThrottleClasses       []throttling.Throttle
}
```

**Actions:**
- `List` - GET /resource/
- `Create` - POST /resource/
- `Retrieve` - GET /resource/{id}/
- `Update` - PUT /resource/{id}/
- `PartialUpdate` - PATCH /resource/{id}/
- `Destroy` - DELETE /resource/{id}/
- Custom actions

**2. Serializer System**

**Base Serializer:**
```go
type Serializer interface {
    Validate(data map[string]interface{}) error
    ToRepresentation(instance interface{}) (map[string]interface{}, error)
    ToInternalValue(data map[string]interface{}) (interface{}, error)
    Create(validatedData map[string]interface{}) (interface{}, error)
    Update(instance interface{}, validatedData map[string]interface{}) (interface{}, error)
}
```

**Field Types:**
- StringField
- IntegerField
- FloatField
- BooleanField
- DateTimeField
- EmailField
- URLField
- UUIDField
- ReadOnlyField
- WriteOnlyField
- HiddenField
- SerializerMethodField

**3. Authentication System**

**Interface:**
```go
type Authentication interface {
    Authenticate(r *http.Request) (*AuthResult, error)
    AuthenticateHeader(r *http.Request) string
}
```

**Implementations:**
- TokenAuthentication
- JWTAuthentication
- SessionAuthentication
- BasicAuthentication
- APIKeyAuthentication

**4. Permission System**

**Interface:**
```go
type Permission interface {
    HasPermission(r *http.Request, view ViewSet) bool
    HasObjectPermission(r *http.Request, view ViewSet, obj interface{}) bool
}
```

**Implementations:**
- AllowAny
- IsAuthenticated
- IsAuthenticatedOrReadOnly
- IsAdminUser
- IsOwnerOrReadOnly

**5. Throttling System**

**Interface:**
```go
type Throttle interface {
    AllowRequest(r *http.Request, view ViewSet) (bool, time.Duration, error)
    GetScope(r *http.Request, view ViewSet) string
}
```

**Implementations:**
- AnonRateThrottle
- UserRateThrottle
- ScopedRateThrottle

**6. Renderers**

**Types:**
- JSONRenderer
- XMLRenderer
- YAMLRenderer
- HTMLRenderer
- CSVRenderer

**7. Parsers**

**Types:**
- JSONParser
- XMLParser
- FormParser
- MultiPartParser
- YAMLParser

**8. Pagination**

**Types:**
- PageNumberPagination
- LimitOffsetPagination

**9. Filtering**

**Types:**
- SearchFilter
- OrderingFilter
- FieldFilter

**10. Versioning**

**Types:**
- URLPathVersioning
- QueryParameterVersioning
- HeaderVersioning

---

### User System Architecture

**Location:** `pkg/identity/`

**Purpose:** Complete user management with authentication, sessions, and permissions.

#### Architecture

```
HTTP Request
    ↓
Middleware (Auth)
    ↓
Handler (User/Auth)
    ↓
Service Layer
    ↓
Repository Layer
    ↓
Database
```

#### Components

**1. Models**

**User Model:**
- ID, Username, Email, Password (hashed)
- IsActive, IsStaff, IsSuperuser
- DateJoined, LastLogin
- FirstName, LastName

**Session Model:**
- UserID, Token, ExpiresAt
- IPAddress, UserAgent
- CreatedAt, LastActivity

**Permission Model:**
- Name, Codename
- ContentType

**Group Model:**
- Name
- Permissions (many-to-many)

**Token Models:**
- EmailVerificationToken
- PasswordResetToken

**2. Repository Layer**

**Interfaces:**
- UserRepository
- SessionRepository
- PermissionRepository
- TokenRepository

**Responsibilities:**
- Database operations
- Query building
- Data access only (no business logic)

**3. Service Layer**

**Interfaces:**
- UserService
- AuthService
- PasswordService
- PermissionService

**Responsibilities:**
- Business logic
- Validation
- Transaction management
- Error handling

**4. Authentication Backends**

**Interface:**
```go
type AuthenticationBackend interface {
    Authenticate(ctx context.Context, credentials map[string]string) (*User, error)
    GetUser(ctx context.Context, identifier string) (*User, error)
    Supports(credentialType string) bool
    Name() string
}
```

**Implementations:**
- PasswordBackend
- TokenBackend

**5. Handlers**

**User Handler:**
- List users
- Create user
- Update user
- Delete user
- Get user

**Auth Handler:**
- Login
- Logout
- Register
- Password reset
- Email verification

**6. Middleware**

**Auth Middleware:**
- Extract authentication token
- Authenticate user
- Set user in context
- Handle authentication errors

---

### HTTP & Routing Architecture

**Location:** `pkg/http/`

**Purpose:** HTTP server, routing, and middleware management.

#### Components

**1. Router**

**Wrapper around chi router:**
```go
type Router struct {
    chi.Router
    middleware []func(http.Handler) http.Handler
}
```

**Features:**
- RESTful routing
- Route groups
- Middleware support
- Sub-routers
- Route mounting

**2. Middleware Stack**

**Default Middleware:**
1. RequestID
2. RealIP
3. Recoverer
4. Logger
5. Session
6. CSRF
7. Authentication

**3. Server**

**Server wrapper:**
```go
type Server struct {
    Router *Router
    Config *ServerConfig
}
```

**Features:**
- Graceful shutdown
- Health checks
- Metrics
- Request logging

---

### Security Architecture

**Location:** `pkg/security/`

**Purpose:** Built-in security protections.

#### Components

**1. CSRF Protection**

**Implementation:**
- gorilla/csrf
- Token generation
- Token validation
- Token rotation

**2. XSS Protection**

**Implementation:**
- HTML escaping in templates
- Content Security Policy
- Input sanitization

**3. SQL Injection Prevention**

**Implementation:**
- Parameter binding
- Identifier escaping
- Query validation

**4. Session Management**

**Implementation:**
- alexedwards/scs/v2
- Secure cookies
- Session expiration
- Session rotation

---

### Validation Architecture

**Location:** `pkg/validation/`

**Purpose:** Data validation with go-playground/validator.

#### Components

**1. Validator**

**Wrapper:**
```go
type Validator struct {
    validator *validator.Validate
}
```

**2. Validation Tags**

**Auto-generated from schema:**
- `required`
- `min`, `max`
- `email`
- `url`
- `uuid`
- Custom validators

**3. Field Validation**

**Field-level validation:**
- Type validation
- Constraint validation
- Custom validation

---

### Configuration Architecture

**Location:** `pkg/config/`

**Purpose:** Application configuration management.

#### Components

**1. Config**

**Viper wrapper:**
```go
type Config struct {
    viper *viper.Viper
}
```

**2. Settings**

**Framework settings:**
- Database configuration
- Server configuration
- Security settings
- Feature flags

**3. Configuration Sources**

**Priority:**
1. Environment variables
2. Config files (YAML, JSON)
3. Defaults

---

### Logging Architecture

**Location:** `pkg/logging/`

**Purpose:** Structured logging with zap.

#### Components

**1. Logger**

**Zap wrapper:**
```go
type Logger struct {
    zap *zap.Logger
}
```

**2. Middleware**

**Request logging:**
- Request method, path
- Response status, duration
- User information
- Error logging

**3. Log Levels**

- Debug
- Info
- Warn
- Error
- Fatal

---

### CLI Architecture

**Location:** `pkg/cli/`

**Purpose:** Command-line interface for framework operations.

#### Commands

**Project Commands:**
- `forge new` - Create new project
- `forge add-app` - Add new app
- `forge add-model` - Add new model
- `forge add-api` - Add API endpoint

**Generation Commands:**
- `forge generate` - Generate code from schemas

**Migration Commands:**
- `forge makemigrations` - Create migrations
- `forge migrate` - Apply migrations
- `forge migrate rollback` - Rollback migrations
- `forge migrate status` - Show migration status

**Server Commands:**
- `forge runserver` - Run development server

**Admin Commands:**
- `forge createsuperuser` - Create admin user

---

## Design Patterns

### 1. Strategy Pattern

**Usage:** Authentication, permissions, throttling, filters, renderers, parsers.

**Example:**
```go
type Authentication interface {
    Authenticate(r *http.Request) (*AuthResult, error)
}

// Multiple implementations
type TokenAuth struct{}
type JWTAuth struct{}
type SessionAuth struct{}
```

### 2. Repository Pattern

**Usage:** Data access layer abstraction.

**Example:**
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
}

type userRepository struct {
    db *db.DB
}
```

### 3. Service Pattern

**Usage:** Business logic encapsulation.

**Example:**
```go
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
}

type userService struct {
    repo UserRepository
}
```

### 4. Factory Pattern

**Usage:** Object creation with dependencies.

**Example:**
```go
type IdentitySystemFactory interface {
    NewUserRepository() UserRepository
    NewUserService() UserService
}

type userSystemFactory struct {
    db *db.DB
}
```

### 5. Builder Pattern

**Usage:** Complex object construction.

**Example:**
```go
type FieldBuilder struct {
    field *Field
}

func (b *FieldBuilder) Required() *FieldBuilder {
    b.field.Required = true
    return b
}

func (b *FieldBuilder) Unique() *FieldBuilder {
    b.field.Unique = true
    return b
}
```

### 6. Observer Pattern

**Usage:** Lifecycle hooks.

**Example:**
```go
type ModelHooks struct {
    BeforeSave []func(*T) error
    AfterSave  []func(*T) error
}

func (m *Model) BeforeSave() error {
    // Hook execution
}
```

### 7. Template Method Pattern

**Usage:** Code generation templates.

**Example:**
```go
// Template defines structure
type ManagerTemplate struct {
    ModelName string
    Fields    []Field
}

// Generated code follows template
func (m *UserManager) Create(ctx context.Context, user *User) error {
    // Template method implementation
}
```

### 8. Chain of Responsibility

**Usage:** Middleware stack.

**Example:**
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic
        next.ServeHTTP(w, r)
    })
}
```

### 9. Decorator Pattern

**Usage:** ViewSet enhancements.

**Example:**
```go
type EnhancedViewSet struct {
    BaseViewSet
    AuthenticationClasses []Authentication
    PermissionClasses     []Permission
}
```

### 10. Registry Pattern

**Usage:** Model and extension registration.

**Example:**
```go
type Registry struct {
    models map[string]interface{}
    mutex  sync.RWMutex
}

func (r *Registry) Register(name string, model interface{}) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    r.models[name] = model
}
```

---

## Type System & Generics

### Generic Type Parameters

**QuerySet:**
```go
type QuerySet[T any] interface {
    All(ctx context.Context) ([]*T, error)
    Get(ctx context.Context) (*T, error)
}
```

**Manager:**
```go
type Manager[T any] interface {
    Create(ctx context.Context, instance *T) error
    Get(ctx context.Context, id int64) (*T, error)
}
```

**FieldExpr:**
```go
type FieldExpr[T any, F any] interface {
    Get(instance *T) F
    Set(instance *T, value F)
    Equals(value F) QueryExpr
}
```

**Admin:**
```go
type Admin[T any] struct {
    model   T
    manager *query.Manager[T]
    config  *Config[T]
}
```

### Type Constraints

**Model Constraint:**
```go
type Model interface {
    GetID() int64
}
```

**Schema Constraint:**
```go
type Schema interface {
    Fields() []Field
    Relations() []Relation
    Meta() Meta
}
```

### Type Safety Mechanisms

1. **Compile-Time Type Checking:** Generics ensure type safety at compile time
2. **Type Inference:** Go infers types from context
3. **Type Assertions:** Runtime type checking when needed
4. **Reflection:** Used sparingly for dynamic operations

---

## Data Flow & Request Processing

### HTTP Request Flow

```
1. HTTP Request arrives
   ↓
2. Chi Router matches route
   ↓
3. Middleware Stack:
   a. RequestID
   b. RealIP
   c. Recoverer
   d. Logger
   e. Session
   f. CSRF
   g. Authentication
   ↓
4. Handler extracts:
   - User from context
   - QuerySet from viewset
   - Request data
   ↓
5. Permission Check
   ↓
6. Throttle Check
   ↓
7. ViewSet Action:
   a. Parse request (Serializer)
   b. Validate data
   c. Execute query (QuerySet)
   d. Serialize response
   ↓
8. Response Renderer
   ↓
9. HTTP Response
```

### Query Execution Flow

```
1. User builds query:
   User.Objects.Filter(User.Fields.Email.Equals("user@example.com"))
   ↓
2. QuerySet builds QueryExpr:
   FieldExprQuery{Field: "email", Op: "=", Value: "user@example.com"}
   ↓
3. User calls execution: .All(ctx)
   ↓
4. QuerySet builds SQL:
   SELECT * FROM "users" WHERE "email" = $1
   ↓
5. SQL executed with database/sql
   ↓
6. Rows scanned into []*User
   ↓
7. Relations loaded if needed
   ↓
8. Results returned
```

### Code Generation Flow

```
1. User defines schema:
   type User struct {
       schema.BaseSchema
   }
   func (User) Fields() []schema.Field {
       return []schema.Field{
           fields.String("email").Required().Unique(),
       }
   }
   ↓
2. User runs: forge generate
   ↓
3. AST Parser:
   a. Loads Go files
   b. Parses AST
   c. Extracts schema definitions
   ↓
4. Generator:
   a. Generates Manager code
   b. Generates QuerySet code
   c. Generates FieldExpr code
   ↓
5. Writer:
   a. Writes generated files
   b. Formats with gofmt
   ↓
6. User imports generated code
```

---

## Performance Considerations

### Query Optimization

1. **Connection Pooling:** Reuse database connections
2. **Query Caching:** Cache compiled SQL queries
3. **Lazy Loading:** Load relations only when needed
4. **Field Selection:** Only select needed fields
5. **Batch Operations:** Bulk insert/update
6. **Index Usage:** Proper database indexes

### Memory Management

1. **Object Pooling:** Reuse objects where possible
2. **Streaming:** Stream large result sets
3. **Garbage Collection:** Minimize allocations
4. **Buffer Reuse:** Reuse buffers for serialization

### Concurrency

1. **Thread Safety:** Mutex protection for shared state
2. **Connection Pooling:** Thread-safe connection pool
3. **Query Isolation:** Each query is independent
4. **Transaction Isolation:** Proper transaction handling

---

## Extension Points

### 1. Custom Field Types

```go
type CustomField struct {
    schema.BaseField
}

func (f *CustomField) ToSQL() string {
    return "CUSTOM_TYPE"
}
```

### 2. Custom Authentication Backends

```go
type CustomAuth struct{}

func (a *CustomAuth) Authenticate(r *http.Request) (*AuthResult, error) {
    // Custom authentication logic
}
```

### 3. Custom Permissions

```go
type CustomPermission struct{}

func (p *CustomPermission) HasPermission(r *http.Request, view ViewSet) bool {
    // Custom permission logic
}
```

### 4. Custom Renderers

```go
type CustomRenderer struct{}

func (r *CustomRenderer) Render(data interface{}) ([]byte, error) {
    // Custom rendering logic
}
```

### 5. Custom Middleware

```go
func CustomMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom middleware logic
        next.ServeHTTP(w, r)
    })
}
```

---

## Best Practices & Guidelines

### 1. Schema Definition

- Use type-specific field builders
- Set appropriate constraints
- Use relations for relationships
- Define hooks for business logic
- Use Meta for table-level options

### 2. Query Building

- Use type-safe FieldExpr when possible
- Use dynamic queries for runtime flexibility
- Chain operations for readability
- Use Select() to limit fields
- Use PrefetchRelated() for relations

### 3. Admin Configuration

- Configure ListDisplay for readability
- Add filters for common queries
- Use fieldsets for form organization
- Add actions for bulk operations
- Customize verbose names

### 4. API Design

- Use serializers for request/response
- Set appropriate permissions
- Use throttling for rate limiting
- Use pagination for large datasets
- Version APIs appropriately

### 5. Error Handling

- Use framework error types
- Provide meaningful error messages
- Log errors appropriately
- Handle errors at appropriate levels

---

## Future Architecture Considerations

### Planned Features

1. **GraphQL Support:** GraphQL schema generation and resolvers
2. **WebSocket Support:** Real-time updates and channels
3. **Caching Layer:** Redis and in-memory caching
4. **Background Tasks:** Task queue integration
5. **Multi-tenancy:** Built-in tenant isolation
6. **Internationalization:** Translation system
7. **Advanced ORM Features:** SelectRelated, PrefetchRelated, Aggregates, Annotations
8. **Testing Infrastructure:** Comprehensive test suite and fixtures

### Architecture Evolution

1. **Plugin System:** Enhanced plugin architecture
2. **Microservices Support:** Service mesh integration
3. **Event System:** Event-driven architecture
4. **Observability:** Metrics, tracing, logging
5. **Performance:** Further optimizations

---

## Conclusion

The forge framework provides a comprehensive, type-safe, and extensible foundation for building web applications in Go. Its architecture is designed with:

- **Type Safety:** Generics ensure compile-time safety
- **Developer Experience:** Convention over configuration
- **Performance:** Optimized for speed and efficiency
- **Extensibility:** Everything can be extended
- **Security:** Built-in protections
- **Code Generation:** Zero boilerplate

The framework continues to evolve, with new features and improvements being added regularly. For the latest information, see the [Roadmap](ROADMAP.md) and [Changelog](CHANGELOG.md).

---

**Document End**

*This document is maintained as part of the forge framework documentation. For questions or contributions, please see the [Development Guide](DEVELOPMENT.md).*
