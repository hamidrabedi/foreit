# Naming Architecture & Design

> **Philosophy**: Every name must be meaningful, intentional, and justified. We are a framework - our API is our contract with developers.

## Table of Contents

1. [Core Principles](#core-principles)
2. [Package Naming](#package-naming)
3. [Type Naming](#type-naming)
4. [Function Naming](#function-naming)
5. [Variable Naming](#variable-naming)
6. [Interface Naming](#interface-naming)
7. [Constant Naming](#constant-naming)
8. [Method Naming](#method-naming)
9. [Query Expression API](#query-expression-api)
10. [Field Expression API](#field-expression-api)
11. [Framework-Specific Conventions](#framework-specific-conventions)
12. [Anti-Patterns](#anti-patterns)
13. [Migration Guide](#migration-guide)

---

## Core Principles

### 1. **Clarity Over Brevity**
- Bad: `NewQ()` - What is Q? Why "New"?
- Good: `Q()` - Direct, mimics Django/Laravel conventions
- Good: `Where()` - Clear intent

### 2. **Framework Consistency**
- When inspired by Django/Laravel/Rails, use their naming unless there's a strong Go-specific reason
- Example: Django uses `Q()` and `F()`, not `NewQ()` or `NewF()`

### 3. **Go Idioms First**
- Exported names start with capital letters
- Interface names end with `-er` when appropriate (Reader, Writer)
- Constructor functions should be meaningful

### 4. **Self-Documenting Code**
- Names should make code readable as prose
- Avoid unnecessary abbreviations

### 5. **Type Safety Without Verbosity**
- Leverage Go generics for type safety
- Don't encode type information in names when generics handle it

---

## Package Naming

### Standard Packages
```
forge/
├── orm/          # Object-Relational Mapping
├── schema/       # Schema Definition
├── db/           # Database Layer
├── migrate/      # Migration System
├── filter/       # Filtering System
├── api/          # REST API Framework
├── admin/        # Admin Interface
├── cli/          # Command Line Interface
├── validate/     # Validation
├── identity/     # Authentication & Authorization
├── server/       # HTTP Server
├── log/          # Logging
├── config/       # Configuration
└── registry/     # Plugin/Component Registry
```

### Rules
- **Singular nouns**: `orm` not `orms`
- **Clear purpose**: Package name indicates what it does
- **No underscores**: Use camelCase for multi-word (rare)
- **Short but meaningful**: `orm` is short but universally understood

### Examples
```go
// Good
import "github.com/forgego/forge/orm"
import "github.com/forgego/forge/filter"

// Bad
import "github.com/forgego/forge/utils"     // Too generic
import "github.com/forgego/forge/helpers"   // What kind of helpers?
```

---

## Type Naming

### Core ORM Types

#### QuerySet
```go
// Primary query interface
type QuerySet[T any] interface {
    Filter(expr Expression) QuerySet[T]
    Exclude(expr Expression) QuerySet[T]
    // ...
}

// Implementation
type BaseQuerySet[T any] struct { }

// Constructor
func NewQuerySet[T any](table string) (QuerySet[T], error)
```

**Justification**:
- `QuerySet` matches Django terminology
- Developers coming from Django/Python understand this immediately
- Generic `[T any]` provides type safety

#### Manager
```go
// Model manager - Django-inspired
type Manager[T any] struct {
    tableName string
    db        *db.DB
}

func NewManager[T any](tableName string, database *db.DB) *Manager[T]
```

**Justification**:
- `Manager` is Django's model manager concept
- Manages model lifecycle and queries
- Clear responsibility

### Expression Types

#### Current (BAD)
```go
// ❌ Problematic naming
type QueryExpr struct { }                              // Vague
func NewFieldQueryExpr(field, op, value) QueryExpr     // Too verbose
func NewQ(expr Expression) *Q                          // Meaningless abbreviation
```

#### Proposed (GOOD)
```go
// ✅ Clear and purposeful
type Expression interface {
    ToSQL(builder *SQLBuilder) (string, []interface{}, error)
    Resolve(schema *ModelSchema) error
}

// Query expressions - Django Q() equivalent
func Q(conditions ...Condition) Expression

// Field expressions - Django F() equivalent  
func F(fieldPath string) FieldRef

// Conditions
func Where(field string, op Operator, value interface{}) Condition
func And(conditions ...Expression) Expression
func Or(conditions ...Expression) Expression
func Not(expr Expression) Expression
```

**Justification**:
- `Q()` matches Django's query object
- `F()` matches Django's field reference
- `Where()` is SQL-like and clear
- No "New" prefix - these are factory functions, not constructors

### Field Types

```go
// Field reference for type-safe queries
type Field[T any] struct {
    path  string
    table string
}

// Constructor - justified "New" because it's a struct initialization
func NewField[T any](path, table string) Field[T]

// But for common use, provide shorthand
func F(path string) FieldRef  // Runtime field reference
```

### Schema Types

```go
// Schema definition
type Schema interface {
    Fields() []Field
    Relations() []Relation
    Meta() Meta
}

// Field definition
type Field struct {
    Name       string
    Type       FieldType
    Required   bool
    // ...
}

// Field types enum
type FieldType int

const (
    TypeInt64 FieldType = iota
    TypeString
    TypeBool
    // ...
)
```

**Justification**:
- Clear hierarchy: Schema → Field → FieldType
- Enums prefixed with `Type` to show they're type variants

---

## Function Naming

### Constructors vs Factories

#### Use "New" Prefix When:
1. Creating concrete struct instances
2. Initialization requires validation/setup
3. Return type is the concrete struct

```go
// Good uses of "New"
func NewQuerySet[T any](table string) (QuerySet[T], error)
func NewManager[T any](table string, db *db.DB) *Manager[T]
func NewField[T any](path, table string) Field[T]
func NewSQLBuilder() *SQLBuilder
```

#### Use Direct Name When:
1. Creating expressions/builders (factory pattern)
2. Simple value construction
3. Mimicking well-known API (Django Q, F)

```go
// Good factory functions (no "New")
func Q(conditions ...Condition) Expression
func F(fieldPath string) FieldRef
func Where(field string, op Operator, value interface{}) Condition
func And(expressions ...Expression) Expression
func Or(expressions ...Expression) Expression

// Field operations (no "New")
func CharField(name string) Field
func IntegerField(name string) Field
func ForeignKey(name, relatedModel string) Field
```

### Query Building Functions

```go
// Current (BAD)
qs.Filter(orm.NewQ(nameField.Eq("John")))  // ❌ Too verbose

// Proposed (GOOD)  
qs.Filter(nameField.Eq("John"))             // ✅ Direct, clean
qs.Filter(Q(Where("name", Equals, "John"))) // ✅ Alternative SQL-like
qs.Filter(And(
    name.Eq("John"),
    age.Gt(18),
))                                          // ✅ Readable
```

### Utility Functions

```go
// Good naming for utilities
func EscapeIdentifier(name string) string
func Paginate[T any](qs QuerySet[T], page, pageSize int) QuerySet[T]
func ParseFilters(queryParams url.Values) ([]Expression, error)
```

**Avoid**:
- Generic names like `Process()`, `Handle()`, `Do()`
- Redundant prefixes: `FilterFilter()`, `QueryQuerySet()`

---

## Variable Naming

### Local Variables

```go
// Good - clear and concise
func GetUser(ctx context.Context, id int64) (*User, error) {
    qs, err := NewQuerySet[User]("users")
    if err != nil {
        return nil, err
    }
    
    user, err := qs.Filter(F("id").Eq(id)).Get(ctx)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// Bad - too abbreviated
func GetUser(ctx context.Context, id int64) (*User, error) {
    q, e := NewQuerySet[User]("users")  // ❌ What is 'q'? What is 'e'?
    if e != nil {
        return nil, e
    }
    
    u, e := q.Filter(F("id").Eq(id)).Get(ctx)  // ❌ 'u' is ambiguous
    return u, e
}
```

### Struct Fields

```go
// Good - full names
type Config struct {
    DatabaseURL     string
    MaxConnections  int
    EnableLogging   bool
    SecretKey       []byte
}

// Bad - unnecessary abbreviations
type Config struct {
    DbURL    string  // ❌ Just spell it out
    MaxConns int     // ❌ Connections is fine
    LoggingOn bool   // ❌ Enable is clearer
    Secret   []byte  // ❌ SecretKey is more specific
}
```

### Receiver Names

```go
// Good - consistent, short but clear
func (qs *BaseQuerySet[T]) Filter(expr Expression) QuerySet[T]
func (m *Manager[T]) Create(ctx context.Context, obj *T) error
func (f *Field[T]) Eq(value T) Expression

// Bad - inconsistent or unclear
func (q *BaseQuerySet[T]) Filter(expr Expression) QuerySet[T]  // ❌ 'q' is ambiguous
func (this *Manager[T]) Create(ctx context.Context, obj *T) error  // ❌ Don't use 'this'
func (fld *Field[T]) Eq(value T) Expression  // ❌ Unnecessary abbreviation
```

**Rules**:
- Use 1-2 letter abbreviation of type name
- Be consistent across all methods of a type
- Common: `qs` (QuerySet), `m` (Manager), `f` (Field), `s` (Schema)

---

## Interface Naming

### Standard Interfaces

```go
// Good - follows Go conventions
type Expression interface {
    ToSQL(*SQLBuilder) (string, []interface{}, error)
}

type Serializer interface {
    Serialize(data interface{}) ([]byte, error)
    Deserialize([]byte) (interface{}, error)
}

type Validator interface {
    Validate(value interface{}) error
}

type Renderer interface {
    Render(data interface{}) ([]byte, error)
}
```

### Naming Rules

1. **Use `-er` suffix when describing capability**:
   - Reader, Writer, Builder, Validator, Renderer

2. **Use noun when describing contract**:
   - Expression, Schema, Field, Relation

3. **Avoid redundant prefixes**:
   - Bad: `IExpression`, `ExpressionInterface`
   - Good: `Expression`

---

## Constant Naming

### Enums and Types

```go
// Field types
type FieldType int

const (
    TypeInt64 FieldType = iota
    TypeString
    TypeBool
    TypeTime
)

// Operators
type Operator string

const (
    OpEquals    Operator = "="
    OpNotEquals Operator = "!="
    OpGreater   Operator = ">"
    OpLess      Operator = "<"
    OpIn        Operator = "IN"
    OpContains  Operator = "LIKE"
)

// Combiners
type Combiner string

const (
    CombineAnd Combiner = "AND"
    CombineOr  Combiner = "OR"
)
```

**Rules**:
- Prefix with type name: `Op`, `Type`, `Combine`
- Use PascalCase
- Make intent obvious from name

### Configuration Constants

```go
const (
    DefaultPageSize     = 20
    MaxPageSize         = 100
    DefaultTimeout      = 30 * time.Second
    MaxUploadSize       = 10 * 1024 * 1024  // 10MB
)
```

---

## Method Naming

### Query Methods

```go
// QuerySet methods - match Django API
Filter(expr Expression) QuerySet[T]
Exclude(expr Expression) QuerySet[T]
OrderBy(fields ...OrderField) QuerySet[T]
Limit(n int) QuerySet[T]
Offset(n int) QuerySet[T]
Distinct(fields ...string) QuerySet[T]

// Execution methods
All(ctx context.Context) ([]*T, error)
Get(ctx context.Context) (*T, error)
First(ctx context.Context) (*T, error)
Count(ctx context.Context) (int64, error)
Exists(ctx context.Context) (bool, error)

// Modification methods
Create(ctx context.Context, obj *T) error
Update(ctx context.Context, updates map[string]interface{}) (int64, error)
Delete(ctx context.Context) (int64, error)
```

### Field Expression Methods

```go
// Comparison operations
Eq(value T) Expression      // Equal
Ne(value T) Expression      // Not equal
Gt(value T) Expression      // Greater than
Gte(value T) Expression     // Greater than or equal
Lt(value T) Expression      // Less than
Lte(value T) Expression     // Less than or equal

// String operations
Contains(value string) Expression
StartsWith(value string) Expression
EndsWith(value string) Expression
IContains(value string) Expression  // Case-insensitive

// Collection operations
In(values ...T) Expression
NotIn(values ...T) Expression

// Null checks
IsNull() Expression
IsNotNull() Expression
```

**Justification**:
- Short but clear: `Eq` vs `Equal` (Django uses `__exact`)
- Follows Django lookup conventions
- Go convention: shortened but readable

---

## Query Expression API

### Current Problems

```go
// ❌ Current implementation - confusing
q := orm.NewQ(nameField.Eq("John"))
combined := q.And(orm.NewQ(ageField.Gt(18)))
qs.Filter(combined)

// ❌ Issues:
// 1. NewQ wrapping is redundant
// 2. Not idiomatic (Django uses Q() directly)
// 3. Too verbose for common operations
```

### Proposed Solution

```go
// ✅ Direct expression (most common)
qs.Filter(name.Eq("John"))
qs.Filter(age.Gt(18))

// ✅ Complex queries with Q() (when needed)
qs.Filter(Q(
    name.Eq("John"),
    age.Gt(18),
))

// ✅ Boolean operations
qs.Filter(And(
    name.Eq("John"),
    Or(
        age.Gt(18),
        age.Lt(65),
    ),
))

// ✅ Negation
qs.Exclude(name.Eq("John"))
qs.Filter(Not(age.Gt(65)))
```

### Implementation

```go
// Q creates a query expression from conditions
// Mimics Django's Q object
func Q(conditions ...Condition) Expression {
    if len(conditions) == 0 {
        return &EmptyExpression{}
    }
    if len(conditions) == 1 {
        return conditions[0]
    }
    return And(conditions...)
}

// F creates a field reference
// Mimics Django's F object for field references
func F(fieldPath string) FieldRef {
    return FieldRef{path: fieldPath}
}

// Boolean combiners
func And(expressions ...Expression) Expression {
    return &BoolExpression{
        operator: CombineAnd,
        children: expressions,
    }
}

func Or(expressions ...Expression) Expression {
    return &BoolExpression{
        operator: CombineOr,
        children: expressions,
    }
}

func Not(expr Expression) Expression {
    return &NotExpression{inner: expr}
}

// Where creates a simple condition (alternative to field methods)
func Where(field string, op Operator, value interface{}) Condition {
    return &SimpleCondition{
        field: field,
        op:    op,
        value: value,
    }
}
```

---

## Field Expression API

### Type-Safe Fields

```go
// Generated or manually created type-safe fields
type UserFields struct {
    ID        Field[int64]
    Name      Field[string]
    Email     Field[string]
    Age       Field[int]
    CreatedAt Field[time.Time]
}

var User = UserFields{
    ID:        NewField[int64]("id", "users"),
    Name:      NewField[string]("name", "users"),
    Email:     NewField[string]("email", "users"),
    Age:       NewField[int]("age", "users"),
    CreatedAt: NewField[time.Time]("created_at", "users"),
}

// Usage
qs.Filter(User.Name.Eq("John"))
qs.Filter(User.Age.Gt(18))
qs.OrderBy(User.CreatedAt.Desc())
```

### String-Based Fields (Fallback)

```go
// When type-safe fields aren't available
qs.Filter(F("name").Eq("John"))
qs.Filter(Where("age", OpGreater, 18))
```

---

## Framework-Specific Conventions

### Admin Package

```go
// Admin configuration
type Admin[T any] struct { }

func Register[T any](schema Schema, manager *Manager[T], config *Config[T]) (*Admin[T], error)

// Field configuration
type FieldConfig struct {
    Name       string
    Label      string
    Widget     Widget
    Readonly   bool
    Searchable bool
}

// List display
type ListConfig struct {
    DisplayFields []string
    SearchFields  []string
    FilterFields  []string
    OrderingField string
}
```

### API Package

```go
// ViewSet - Django REST Framework naming
type ViewSet[T any] interface {
    List(ctx *Context) error
    Create(ctx *Context) error
    Retrieve(ctx *Context) error
    Update(ctx *Context) error
    Delete(ctx *Context) error
}

// Serializer - DRF naming
type Serializer interface {
    Serialize(obj interface{}) (map[string]interface{}, error)
    Deserialize(data map[string]interface{}) (interface{}, error)
    Validate(data map[string]interface{}) error
}

// ViewSet actions
func (vs *BaseViewSet[T]) List(ctx *Context) error
func (vs *BaseViewSet[T]) Create(ctx *Context) error
```

### Filter Package

```go
// FilterSet - Django-filter naming
type FilterSet[T any] struct {
    filters map[string]Filter[T]
}

func NewFilterSet[T any]() *FilterSet[T]

// Filter types
type CharFilter[T any] struct { }
type IntegerFilter[T any] struct { }
type BooleanFilter[T any] struct { }
type DateFilter[T any] struct { }
type ChoiceFilter[T any] struct { }
```

---

## Anti-Patterns

### ❌ Avoid These

```go
// 1. Meaningless abbreviations with "New"
func NewQ(expr Expression) *Q  // ❌ What is Q?
func NewF(field string) *F     // ❌ What is F?

// 2. Redundant type prefixes
type QuerySetQuerySet[T any] interface { }  // ❌ Redundant
type FilterFilter[T any] struct { }         // ❌ Redundant

// 3. Generic utility packages
package utils   // ❌ Too vague
package helpers // ❌ What kind of helpers?
package common  // ❌ Everything is common

// 4. Hungarian notation
func GetStrName() string      // ❌ Type in name
var intCount int              // ❌ Redundant
var boolIsActive bool         // ❌ Redundant

// 5. Unnecessary abbreviations in exported API
func GetUsr() *User           // ❌ Just spell it out
func ProcReq() error          // ❌ Unclear
func ValidFld() bool          // ❌ Ambiguous

// 6. Inconsistent naming
func (qs *QuerySet) Filter()  // ✅
func (q *QuerySet) Exclude()  // ❌ Inconsistent receiver name

// 7. Redundant context
package orm
type ORMQuerySet interface { }  // ❌ Package already says "orm"
```

### ✅ Use These Instead

```go
// 1. Meaningful factory functions
func Q(conditions ...Condition) Expression  // ✅ Clear, matches Django
func F(fieldPath string) FieldRef          // ✅ Clear, matches Django

// 2. Clean type names
type QuerySet[T any] interface { }  // ✅
type Filter[T any] interface { }    // ✅

// 3. Specific packages
package stringutil    // ✅ Specific purpose
package timeutil      // ✅ Specific purpose

// 4. Clear names
func GetName() string       // ✅
var count int               // ✅
var isActive bool           // ✅

// 5. Full words in API
func GetUser() *User        // ✅
func ProcessRequest() error // ✅
func ValidateField() bool   // ✅

// 6. Consistent receivers
func (qs *QuerySet) Filter()  // ✅
func (qs *QuerySet) Exclude() // ✅

// 7. Avoid package redundancy
package orm
type QuerySet interface { }  // ✅ Clean in context
```

---

## Migration Guide

### Phase 1: Query Expressions

**Changes**:
```go
// OLD
q := orm.NewQ(expr)
combined := q.Or(orm.NewQ(other))

// NEW
q := orm.Q(expr)
combined := q.Or(other)

// OR simply
qs.Filter(expr)
qs.Filter(Or(expr1, expr2))
```

**Deprecation**:
```go
// Mark as deprecated
// Deprecated: Use Q() instead. NewQ will be removed in v2.0.
func NewQ(expr Expression) *QObject { return Q(expr) }
```

### Phase 2: Field References

**Changes**:
```go
// OLD
f := orm.NewFieldQueryExpr("name", orm.OpEquals, "John")

// NEW
f := User.Name.Eq("John")
// OR
f := F("name").Eq("John")
// OR
f := Where("name", OpEquals, "John")
```

### Phase 3: Constructor Cleanup

**Review all `New*` functions**:
- Keep `New*` for actual constructors (NewQuerySet, NewManager)
- Replace factory `New*` with direct functions (Q, F, Where)

---

## Naming Checklist

Before adding any new name to the framework, ask:

1. **Is it meaningful?**
   - ✅ Can someone unfamiliar with the codebase understand it?
   - ❌ Does it require documentation to understand?

2. **Is it consistent?**
   - ✅ Does it follow patterns established in the framework?
   - ✅ Does it match Django/Laravel conventions where applicable?

3. **Is it justified?**
   - ✅ If it's an abbreviation, is it widely recognized?
   - ✅ If it's verbose, is the verbosity necessary?

4. **Is it idiomatic?**
   - ✅ Does it follow Go conventions?
   - ✅ Does it leverage Go's type system well?

5. **Is it future-proof?**
   - ✅ Will it make sense as the framework grows?
   - ✅ Is it specific enough to avoid conflicts?

---

## Summary

### Key Principles

1. **Clarity over brevity** - spell it out if there's any ambiguity
2. **Framework consistency** - match Django/Laravel/Rails when it makes sense
3. **Go idioms** - follow Go conventions for exported names
4. **Self-documenting** - code should read like prose
5. **Type safety** - leverage generics instead of encoding types in names

### Quick Reference

| Pattern | Bad | Good | Why |
|---------|-----|------|-----|
| Query Objects | `NewQ()` | `Q()` | Matches Django, clearer |
| Field References | `NewF()` | `F()` | Matches Django, clearer |
| Conditions | `NewFieldQueryExpr()` | `Where()` | Self-documenting |
| Variables | `u`, `e`, `q` | `user`, `err`, `qs` | Clear intent |
| Constants | `EQUALS` | `OpEquals` | Namespaced |
| Interfaces | `IValidator` | `Validator` | Go convention |
| Receivers | `this`, `self` | `v`, `qs`, `m` | Go convention |

### Implementation Priority

1. **High Priority** (Breaking changes):
   - Replace `NewQ()` with `Q()`
   - Replace `NewF()` with `F()`
   - Simplify `NewFieldQueryExpr()`

2. **Medium Priority** (Deprecations):
   - Standardize constructor naming
   - Clean up internal naming
   - Document naming conventions

3. **Low Priority** (Polish):
   - Variable name consistency
   - Comment improvements
   - Example code updates

---

**Last Updated**: 2026-01-01
**Version**: 1.0
**Status**: Draft - Pending Review
