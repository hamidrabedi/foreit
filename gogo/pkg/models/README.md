# Models Package

A comprehensive model system inspired by Prisma, Pydantic, SQLAlchemy, and Ent, designed to be Go-idiomatic and fully composable.

## Architecture: Schema-First Codegen

This package follows the **schema-first codegen** approach (like Ent/Prisma):

1. **Define schema** (using Ent or custom DSL)
2. **Generate code** (type-safe structs, queries, admin metadata)
3. **Use runtime** (validation, serialization, sessions, signals)

**Why this approach?**
- ✅ Compile-time safety (no runtime reflection)
- ✅ Full IDE support (autocomplete, type checking)
- ✅ Go-idiomatic (codegen is standard in Go)
- ✅ High performance (zero reflection overhead)

## Design Philosophy

- **Type Safety**: Leverage Go generics and Ent's code generation
- **Composability**: Everything is an interface, everything can be overridden
- **Validation**: Pydantic-inspired validation with struct tags
- **Serialization**: Flexible JSON serialization with field inclusion/exclusion
- **Sessions**: SQLAlchemy-inspired session and transaction management
- **Signals**: Django-like signals for lifecycle events

## Core Components

### Model Interface

```go
type Model interface {
    Save(ctx context.Context) error
    Delete(ctx context.Context) error
    String() string
    IsNew() bool
    GetID() interface{}
    SetID(id interface{})
    GetCreatedAt() *time.Time
    GetUpdatedAt() *time.Time
}
```

### BaseModel

Embed `BaseModel` in your structs to get Django-like model behavior:

```go
type User struct {
    models.BaseModel
    Name  string `json:"name" validate:"required,min_length=3"`
    Email string `json:"email" validate:"required,email"`
}

// Usage
user := &User{
    BaseModel: *models.NewBaseModel(),
    Name: "John",
    Email: "john@example.com",
}

// Set manager (connects to database)
user.SetManager(userManager)

// Save with validation and signals
if err := user.Save(ctx); err != nil {
    // Handle error
}
```

### Validation (Pydantic-inspired)

```go
// Struct tags for validation
type User struct {
    Name  string `validate:"required,min_length=3,max_length=100"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=120"`
}

// Validate manually
if err := models.DefaultValidator.Validate(user); err != nil {
    // Handle validation errors
}

// Get detailed errors
errors := models.DefaultValidator.GetValidationErrors(user)
```

### Serialization (Pydantic-inspired)

```go
serializer := models.NewJSONSerializer()

// Serialize with field exclusion
data, _ := serializer.ExcludeFields("password", "secret").Serialize(user)

// Serialize with field inclusion only
data, _ := serializer.IncludeFields("id", "name", "email").Serialize(user)

// Deserialize
var user User
serializer.Deserialize(data, &user)
```

### Sessions (SQLAlchemy-inspired)

```go
// Get session from context
session, _ := models.ContextSession(ctx)

// Use transaction
err := models.WithTransaction(ctx, session, func(txCtx context.Context) error {
    user1.Save(txCtx)
    user2.Save(txCtx)
    return nil // Auto-commits, or rollback on error
})
```

### QuerySets (Prisma/Django-inspired)

```go
// Chainable queries
users, _ := manager.Filter(ctx).
    Filter(models.QFilter("age", ">=", 18)).
    Filter(models.QFilter("active", "=", true)).
    OrderBy("name", "-created_at").
    Limit(10).
    All(ctx)

// Using F expressions
users, _ := manager.Filter(ctx).
    Filter(models.FField("age").Gte(18)).
    Filter(models.FField("email").Contains("@example.com")).
    All(ctx)

// Complex Q objects
q := models.NewQ().
    And("age", ">=", 18).
    Or("role", "=", "admin")

users, _ := manager.Filter(ctx).Filter(q).All(ctx)
```

### Signals (Django-inspired)

```go
// Connect to signals
models.PreSave.Connect(func(ctx context.Context, sender models.Model, signalType models.SignalType) error {
    // Called before any model is saved
    return nil
})

models.PostCreate.Connect(func(ctx context.Context, sender models.Model, signalType models.SignalType) error {
    // Called after a new model is created
    return nil
})
```

### Lifecycle Hooks

```go
type User struct {
    models.BaseModel
    Name string
}

// Implement ModelWithHooks for custom logic
func (u *User) BeforeSave(ctx context.Context) error {
    // Custom logic before save
    return nil
}

func (u *User) AfterCreate(ctx context.Context) error {
    // Send welcome email, etc.
    return nil
}
```

## Integration with Ent

The models package is designed to work seamlessly with Ent:

1. Use Ent for schema definition and code generation
2. Embed `BaseModel` in your Ent-generated types (or wrap them)
3. Create managers that use Ent clients
4. Leverage Ent's type-safe queries

## Best Practices

1. **Always validate**: Use struct tags and call `Validate()` before saving
2. **Use transactions**: Wrap related operations in `WithTransaction`
3. **Leverage signals**: Use signals for cross-cutting concerns (logging, notifications)
4. **Compose, don't inherit**: Use interfaces and composition
5. **Type safety first**: Prefer generics and code generation over reflection

