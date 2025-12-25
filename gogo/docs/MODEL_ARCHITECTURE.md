# Model Architecture - Schema-First Codegen Approach

## Philosophy

**Core Principle:** Schema → Codegen → Type-Safe Runtime

This architecture combines the best of:
- **Ent** (Go-native codegen, type safety)
- **Prisma** (Great DX, schema-first)
- **Django** (Batteries-included, admin integration)
- **Pydantic** (Validation, serialization)
- **SQLAlchemy** (Sessions, relationships)

## Architecture Layers

```
┌─────────────────────────────────────────┐
│  1. Schema Definition (Go DSL)          │
│     - Type-safe, in Go code             │
│     - Version controlled                │
│     - Ent-compatible                    │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│  2. Code Generator                       │
│     - Models (structs)                  │
│     - Type-safe query builders          │
│     - Graph/edge traversal              │
│     - Migration helpers                 │
│     - Admin metadata                    │
│     - Validation rules                 │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│  3. Generated Code                       │
│     - Zero reflection                   │
│     - Full IDE support                  │
│     - Compile-time safety               │
└──────────────┬──────────────────────────┘
               ↓
┌─────────────────────────────────────────┐
│  4. Runtime Layer                        │
│     - Models package (validation)       │
│     - Sessions (transactions)           │
│     - Signals (lifecycle events)        │
│     - Admin (ModelAdmin)                │
└─────────────────────────────────────────┘
```

## Why This Approach Wins

### ✅ Compile-Time Safety
- No runtime reflection for core operations
- Type errors caught at compile time
- IDE autocomplete works perfectly

### ✅ Performance
- Zero reflection overhead
- Direct struct field access
- Optimized query builders

### ✅ Go Idiomatic
- Codegen is standard in Go (protobuf, OpenAPI)
- Explicit, not magical
- Easy to understand and debug

### ✅ Extensible
- Everything is an interface
- Everything can be overridden
- Composition over inheritance

## Schema Definition

### Using Ent (Recommended)

```go
// schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id"),
        field.String("email").Unique(),
        field.String("name").MaxLen(255),
        field.Enum("role").Values("user", "admin").Default("user"),
        field.Time("created_at").Default(time.Now),
        field.Time("updated_at").UpdateDefault(time.Now),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("posts", Post.Type),
        edge.To("comments", Comment.Type),
    }
}

func (User) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("email", "role"),
    }
}
```

### Generated Code Usage

```go
// Fully typed, autocomplete works
user, err := client.User.
    Create().
    SetEmail("john@example.com").
    SetName("John").
    SetRole(user.RoleAdmin).
    Save(ctx)

// Type-safe queries
users, err := client.User.
    Query().
    Where(user.EmailContains("@example.com")).
    Where(user.RoleEQ(user.RoleAdmin)).
    Order(user.ByCreatedAt(ent.Desc)).
    Limit(10).
    All(ctx)

// Graph queries with edges
users, err := client.User.
    Query().
    WithPosts(func(q *ent.PostQuery) {
        q.Where(post.StatusEQ("published"))
    }).
    All(ctx)
```

## Model Layer Integration

### Embedding BaseModel

```go
// Your generated Ent model
type User struct {
    ent.User
    models.BaseModel
}

// Now you get Django-like methods
user := &User{
    User: *entUser,
    BaseModel: *models.NewBaseModel(),
}

// Set manager (connects to Ent client)
user.SetManager(userManager)

// Save with validation and signals
if err := user.Save(ctx); err != nil {
    // Handle error
}

// Delete with hooks
if err := user.Delete(ctx); err != nil {
    // Handle error
}
```

### Validation Integration

```go
// Add validation tags to Ent schema
type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("email").
            Unique().
            Validate(func(s string) error {
                return models.DefaultValidator.ValidateField(&struct{
                    Email string `validate:"required,email"`
                }{Email: s}, "Email")
            }),
    }
}
```

## Admin Integration

### Auto-Generated Admin

```go
// Admin is generated from schema
type UserAdmin struct {
    *admin.BaseModelAdmin
}

func NewUserAdmin() *UserAdmin {
    admin := admin.NewBaseModelAdmin("User")
    
    // Auto-configured from schema
    admin.SetListDisplay("id", "email", "name", "role", "created_at")
    admin.SetSearchFields("email", "name")
    admin.SetListFilter(
        admin.FilterSpec{Field: "role", Type: admin.FilterTypeChoice},
        admin.FilterSpec{Field: "created_at", Type: admin.FilterTypeDate},
    )
    
    return &UserAdmin{BaseModelAdmin: admin}
}

// Override anything
func (a *UserAdmin) SaveModel(ctx context.Context, obj interface{}, form interface{}, change bool) error {
    // Custom save logic
    return a.BaseModelAdmin.SaveModel(ctx, obj, form, change)
}
```

## QuerySet Integration

### Using with Ent

```go
// Ent query builder
entQuery := client.User.Query().
    Where(user.EmailContains("@example.com"))

// Convert to our QuerySet for additional features
qs := models.NewQuerySet("User", userManager)
// Apply Ent predicates to QuerySet
users, err := qs.All(ctx)
```

## Migration System

### Auto-Detection (Future)

```bash
# Detect schema changes
gogo migration create add_user_bio

# Generated migration
# migrations/20241221_add_user_bio.go
```

```go
func (m *Migration20241221) Up(ctx context.Context) error {
    return m.Exec(`
        ALTER TABLE users ADD COLUMN bio TEXT;
        CREATE INDEX idx_users_bio ON users(bio);
    `)
}

func (m *Migration20241221) Down(ctx context.Context) error {
    return m.Exec(`ALTER TABLE users DROP COLUMN bio;`)
}
```

## Key Design Decisions

### 1. Schema-First, Not Reflection-First
- Define schema explicitly
- Generate code from schema
- No runtime introspection

### 2. Type Safety Over Magic
- Compile-time errors, not runtime
- Explicit APIs, not hidden behavior
- IDE support, not guesswork

### 3. Composition Over Inheritance
- Interfaces everywhere
- Embed structs for behavior
- Override methods as needed

### 4. Go Idioms First
- Codegen is standard
- Explicit is better than implicit
- Fast compile times

## Integration Points

### Ent → Models
- Use Ent for schema and codegen
- Embed BaseModel for Django-like methods
- Use Manager to bridge Ent and Models

### Models → Admin
- Auto-generate admin from schema
- Use ModelAdmin for customization
- Everything overrideable

### Models → Validation
- Struct tags for validation
- Pydantic-inspired validators
- Integrated with Save()

### Models → Serialization
- JSON serialization with field control
- Pydantic-inspired serializers
- Flexible inclusion/exclusion

## Future Enhancements

1. **Auto-Migration Generation**
   - Detect schema changes
   - Generate migration files
   - Apply/rollback support

2. **Better Admin Integration**
   - Auto-generate from Ent schema
   - Rich form widgets
   - Inline editing

3. **Enhanced Validation DSL**
   - Built into Ent fields
   - Custom validators
   - Cross-field validation

4. **Single CLI Tool**
   - `gogo generate` - codegen
   - `gogo migrate` - migrations
   - `gogo admin` - admin UI
   - `gogo validate` - validation

5. **Django-Quality Docs**
   - Comprehensive examples
   - Best practices
   - Migration guides

