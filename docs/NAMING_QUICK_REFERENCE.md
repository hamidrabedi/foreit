# Naming Quick Reference

> Fast lookup for Forge framework naming conventions. For detailed explanations, see [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md).

## Quick Decision Tree

```
Need to name something?
│
├─ Is it a package?
│  └─ Use singular noun: orm, filter, admin
│
├─ Is it a type/struct?
│  ├─ Interface? → Noun or -er suffix: Expression, Validator
│  └─ Struct? → Clear noun: QuerySet, Manager, Field
│
├─ Is it a function?
│  ├─ Constructor (returns concrete struct)?
│  │  └─ Use New prefix: NewQuerySet(), NewManager()
│  │
│  ├─ Factory (returns interface/builds objects)?
│  │  └─ Direct name: Q(), F(), Where(), CharField()
│  │
│  └─ Regular function?
│     └─ Verb + noun: CreateUser(), BuildSQL(), ValidateField()
│
├─ Is it a constant?
│  └─ UPPER or prefixed PascalCase: OpEquals, TypeString
│
├─ Is it a variable?
│  ├─ Receiver? → 1-2 letter type abbreviation: qs, m, f
│  ├─ Local? → Full word, clear intent: user, count, isActive
│  └─ Struct field? → PascalCase, no abbreviations: DatabaseURL
│
└─ Is it a method?
   ├─ Query operation? → Django-style: Filter(), All(), Get()
   └─ Regular method? → Go-style: String(), Error(), Close()
```

---

## Common Patterns

### Constructors (Use "New")

```go
// ✅ Good - Complex initialization, returns concrete type
func NewQuerySet[T any](table string) (QuerySet[T], error)
func NewManager[T any](table string, db *DB) *Manager[T]
func NewSQLBuilder() *SQLBuilder
func NewFilterSet[T any]() *FilterSet[T]
```

### Factories (No "New")

```go
// ✅ Good - Simple construction, returns interface/wrapper
func Q(conditions ...Condition) Expression
func F(fieldPath string) FieldRef
func Where(field string, op Operator, value interface{}) Condition
func And(expressions ...Expression) Expression
func Or(expressions ...Expression) Expression

// Field factories
func CharField(name string) Field
func IntegerField(name string) Field
func ForeignKey(name, related string) Field
```

### Query Operations

```go
// ✅ Django-style query methods
Filter(expr Expression) QuerySet[T]
Exclude(expr Expression) QuerySet[T]
OrderBy(fields ...OrderField) QuerySet[T]
Limit(n int) QuerySet[T]
All(ctx context.Context) ([]*T, error)
Get(ctx context.Context) (*T, error)
Count(ctx context.Context) (int64, error)
```

### Field Expressions

```go
// ✅ Short but clear
field.Eq(value)       // Equal
field.Ne(value)       // Not equal
field.Gt(value)       // Greater than
field.Gte(value)      // Greater than or equal
field.Lt(value)       // Less than
field.Lte(value)      // Less than or equal
field.Contains(str)   // String contains
field.In(values...)   // In list
field.IsNull()        // Is NULL
```

---

## Anti-Patterns (DON'T)

```go
// ❌ Bad - Meaningless abbreviations
func NewQ(expr Expression) *Q
func NewF(field string) *F
func GetUsr() *User
func ProcReq() error

// ❌ Bad - Redundant prefixes
type QuerySetQuerySet interface {}
package utils  // Too generic

// ❌ Bad - Hungarian notation
var strName string
var intCount int

// ❌ Bad - Inconsistent receivers
func (qs *QuerySet) Filter() {}
func (q *QuerySet) Exclude() {}  // Should be 'qs'

// ❌ Bad - Generic names
func Do() error
func Process() error
func Handle() error
```

---

## Package-Specific Conventions

### forge/orm

```go
// Types
type QuerySet[T any] interface
type Manager[T any] struct
type Expression interface
type Field[T any] struct

// Constructors
func NewQuerySet[T any](table string) (QuerySet[T], error)
func NewManager[T any](table, db) *Manager[T]
func NewField[T any](path, table string) Field[T]

// Factories
func Q(conditions ...Condition) Expression
func F(path string) FieldRef
func Where(field, op, value) Condition
func And(...Expression) Expression
func Or(...Expression) Expression

// Operators
const (
    OpEquals    Operator = "="
    OpGreater   Operator = ">"
    OpContains  Operator = "LIKE"
)
```

### forge/filter

```go
// Types
type FilterSet[T any] struct
type Filter[T any] interface

// Constructors
func NewFilterSet[T any]() *FilterSet[T]

// Filters
func CharFilter[T any](field string) Filter[T]
func IntegerFilter[T any](field string) Filter[T]
func BooleanFilter[T any](field string) Filter[T]
func DateFilter[T any](field string) Filter[T]
```

### forge/admin

```go
// Types
type Admin[T any] struct
type Config[T any] struct

// Registration
func Register[T any](schema, manager, config) (*Admin[T], error)

// Configuration
type FieldConfig struct {
    Name     string
    Label    string
    Widget   Widget
    Readonly bool
}
```

### forge/api

```go
// Types
type ViewSet[T any] interface
type Serializer interface

// ViewSet actions (Django REST Framework)
func (vs *ViewSet) List(ctx) error
func (vs *ViewSet) Create(ctx) error
func (vs *ViewSet) Retrieve(ctx) error
func (vs *ViewSet) Update(ctx) error
func (vs *ViewSet) Delete(ctx) error
```

### forge/schema

```go
// Types
type Schema interface
type Field struct
type FieldType int

// Field factories
func CharField(name string) Field
func IntegerField(name string) Field
func ForeignKey(name, related string) Field
func BooleanField(name string) Field

// Field types
const (
    TypeInt64  FieldType = iota
    TypeString
    TypeBool
)
```

---

## Variable Naming

### Receivers

```go
func (qs *QuerySet) Filter()     // ✅ qs = QuerySet
func (m *Manager) Create()       // ✅ m = Manager
func (f *Field) Validate()       // ✅ f = Field
func (s *Schema) Fields()        // ✅ s = Schema
func (v *ViewSet) List()         // ✅ v = ViewSet
func (ser *Serializer) Validate() // ✅ ser = Serializer
```

### Common Variables

```go
// ✅ Good
ctx context.Context
err error
qs QuerySet[T]
user *User
count int64
isActive bool

// ❌ Bad
c context.Context  // Use ctx
e error           // Use err  
q QuerySet[T]     // Use qs
u *User           // Use user
cnt int64         // Use count
```

### Loop Variables

```go
// ✅ Good
for i, user := range users {
    // i and user are clear
}

for _, field := range fields {
    // Underscore for unused index
}

for key, value := range mapping {
    // key/value are standard
}

// ❌ Bad
for i, v := range users {  // What is v?
    // Use descriptive name: user, item, etc.
}
```

---

## Type Naming

### Interfaces

```go
// ✅ Good - Describes capability (-er suffix)
type Validator interface
type Renderer interface
type Serializer interface
type Builder interface

// ✅ Good - Describes contract (noun)
type Expression interface
type Schema interface
type Field interface

// ❌ Bad
type IValidator interface      // No "I" prefix
type ValidatorInterface        // Redundant "Interface"
```

### Structs

```go
// ✅ Good - Clear, descriptive nouns
type QuerySet[T any] struct
type Manager[T any] struct
type Config struct
type FieldInfo struct

// ❌ Bad
type QS[T any] struct          // Too abbreviated
type Mgr[T any] struct         // Too abbreviated  
type Cfg struct                // Too abbreviated
```

### Enums

```go
// ✅ Good - Prefixed with type
type FieldType int

const (
    TypeInt64  FieldType = iota
    TypeString
    TypeBool
)

type Operator string

const (
    OpEquals   Operator = "="
    OpGreater  Operator = ">"
)

// ❌ Bad - No prefix (namespace pollution)
const (
    Int64  FieldType = iota  // Could conflict
    String                    // Conflicts with builtin
)
```

---

## Method Naming

### Query Methods (Django-style)

```go
Filter(Expression) QuerySet[T]
Exclude(Expression) QuerySet[T]
OrderBy(...OrderField) QuerySet[T]
Limit(int) QuerySet[T]
Offset(int) QuerySet[T]
Distinct(...string) QuerySet[T]

All(context.Context) ([]*T, error)
Get(context.Context) (*T, error)
First(context.Context) (*T, error)
Count(context.Context) (int64, error)
Exists(context.Context) (bool, error)

Create(context.Context, *T) error
Update(context.Context, map[string]interface{}) (int64, error)
Delete(context.Context) (int64, error)
```

### Field Operations

```go
Eq(value T) Expression
Ne(value T) Expression
Gt(value T) Expression
Gte(value T) Expression
Lt(value T) Expression
Lte(value T) Expression

Contains(string) Expression
StartsWith(string) Expression
EndsWith(string) Expression

In(...T) Expression
NotIn(...T) Expression
IsNull() Expression
IsNotNull() Expression
```

### Standard Go Methods

```go
String() string
Error() string
Close() error
Read([]byte) (int, error)
Write([]byte) (int, error)
```

---

## Constants

```go
// ✅ Good - Prefixed, clear
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
    MinPasswordLen  = 8
)

const (
    OpEquals   Operator = "="
    OpContains Operator = "LIKE"
)

// ❌ Bad - No prefix, unclear
const (
    PageSize = 20        // Default? Max? Min?
    Max      = 100       // Max what?
    EQUALS   = "="       // No type context
)
```

---

## File Naming

```go
// ✅ Good - Descriptive, snake_case
queryset.go
query_expr.go
field_expression.go
sql_builder.go

// ❌ Bad
qs.go              // Too abbreviated
querySet.go        // Use snake_case
query-expr.go      // Use underscore, not dash
```

---

## Import Aliases

```go
// ✅ Good - Clear, necessary
import (
    "database/sql"
    
    "github.com/forgego/forge/orm"
    adminorm "github.com/forgego/forge/admin/orm"
)

// ✅ Good - Standard abbreviations
import (
    "context"
    
    "github.com/lib/pq"
    _ "github.com/mattn/go-sqlite3"  // Blank import for driver
)

// ❌ Bad - Unnecessary aliases
import (
    o "github.com/forgego/forge/orm"     // orm is fine
    f "github.com/forgego/forge/filter"  // filter is fine
)
```

---

## Comments and Documentation

```go
// ✅ Good - Starts with name, explains purpose
// QuerySet is a type-safe interface for building and executing database queries.
// It provides a fluent API similar to Django's QuerySet.
type QuerySet[T any] interface

// NewQuerySet creates a new QuerySet for the specified table.
// If tableName is empty, it will be derived from the model's schema.
func NewQuerySet[T any](tableName string) (QuerySet[T], error)

// ❌ Bad - Doesn't start with name, vague
// This is the main query interface
type QuerySet[T any] interface

// Creates a new query set
func NewQuerySet[T any](tableName string) (QuerySet[T], error)
```

---

## Real-World Examples

### Example 1: User Query

```go
// ✅ Good
users, err := userManager.
    Filter(And(
        User.Age.Gte(18),
        User.Email.Contains("@example.com"),
        User.IsActive.Eq(true),
    )).
    OrderBy(User.CreatedAt.Desc()).
    Limit(10).
    All(ctx)

// ❌ Bad
q1 := NewQ(NewFieldQueryExpr("age", OpGreaterOrEqual, 18))
q2 := NewQ(NewFieldQueryExpr("email", OpContains, "@example.com"))
q3 := NewQ(NewFieldQueryExpr("is_active", OpEquals, true))
combined := q1.And(q2).And(q3)
u, e := um.Filter(combined).OrderBy("-created_at").Limit(10).All(c)
```

### Example 2: Complex Filter

```go
// ✅ Good
products, err := productManager.
    Filter(Or(
        And(
            Product.Category.Eq("electronics"),
            Product.Price.Lte(1000),
        ),
        And(
            Product.Category.Eq("books"),
            Product.Rating.Gte(4.5),
        ),
    )).
    All(ctx)

// ❌ Bad  
q1 := NewQ(NewFieldQueryExpr("category", OpEquals, "electronics"))
q2 := NewQ(NewFieldQueryExpr("price", OpLessOrEqual, 1000))
q3 := q1.And(q2)
// ... too verbose
```

### Example 3: Admin Registration

```go
// ✅ Good
userAdmin, err := admin.Register(
    userSchema,
    userManager,
    &admin.Config[User]{
        ListDisplay:  []string{"id", "name", "email"},
        SearchFields: []string{"name", "email"},
        ListFilter:   []string{"is_active", "role"},
    },
)

// ❌ Bad
ua, e := admin.Reg(
    us,
    um,
    &admin.Cfg[User]{
        LD: []string{"id", "name", "email"},
        SF: []string{"name", "email"},
        LF: []string{"is_active", "role"},
    },
)
```

---

## Migration Checklist

When migrating code to new conventions:

- [ ] Replace `NewQ()` with direct expressions or `Q()`
- [ ] Replace `NewFieldQueryExpr()` with `Where()` or field methods
- [ ] Use type-safe fields where possible
- [ ] Ensure receiver names are consistent
- [ ] Use full words for variables (no abbreviations)
- [ ] Check that all constants have proper prefixes
- [ ] Verify import aliases are necessary
- [ ] Update comments to start with function/type name
- [ ] Run linter and fix any naming warnings

---

## Quick Lookup Table

| Situation | Good | Bad | Why |
|-----------|------|-----|-----|
| Query object | `Q()` | `NewQ()` | Matches Django, clearer |
| Field ref | `F()` | `NewF()` | Matches Django, clearer |
| Condition | `Where()` | `NewFieldQueryExpr()` | Self-documenting |
| Constructor | `NewQuerySet()` | `QuerySet()` | Go convention |
| Factory | `CharField()` | `NewCharField()` | Cleaner for factories |
| Receiver | `qs`, `m`, `f` | `q`, `mgr`, `fld` | Consistent, clear |
| Variable | `user`, `count` | `u`, `cnt` | Readable |
| Interface | `Validator` | `IValidator` | Go convention |
| Constant | `OpEquals` | `EQUALS` | Namespaced |
| Method | `Filter()` | `DoFilter()` | Concise, clear |

---

## Need More Detail?

- Full explanation: [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md)
- Current issues & migration: [NAMING_AUDIT.md](./NAMING_AUDIT.md)
- API Reference: [API_REFERENCE.md](./API_REFERENCE.md)

---

**Last Updated**: 2026-01-01
**Version**: 1.0
