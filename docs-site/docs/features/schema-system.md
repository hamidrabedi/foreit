---
sidebar_position: 1
description: Define your models with forge's schema system. Type-safe field definitions, relationships, and metadata that just work.
keywords:
  - forge schema
  - model definitions
  - go orm schema
  - django models go
  - type-safe models
image: /img/forge-social-card.jpg
---

# Schema System

The schema system is where you define your models. If you've used Django, you'll feel right at home - but with Go's type safety keeping you honest.

## Defining a Model

Here's how simple it is to define a model:

```go
type User struct {
    schema.BaseSchema
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("username").Required().Unique().MaxLength(150).Build(),
        schema.String("email").Required().Unique().MaxLength(254).Build(),
        schema.Text("bio").Blank().Build(),
        schema.Bool("is_active").Default(true).Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
        schema.Time("updated_at").AutoNow().Build(),
    }
}

func (User) Meta() schema.Meta {
    return schema.Meta{
        TableName: "users",
        VerboseName: "User",
        VerboseNamePlural: "Users",
        Ordering: []string{"-created_at"},
    }
}

func (User) Relations() []schema.Relation {
    return []schema.Relation{
        schema.HasMany("posts", "Post", "user_id"),
    }
}
```

That's it. No tags, no reflection magic, just clean Go code that your IDE can understand.

## Field Types

forge gives you all the field types you'd expect:

### Basic Types
```go
schema.Int64("id").Primary().AutoIncrement().Build()
schema.String("name").Required().MaxLength(100).Build()
schema.Text("description").Blank().Build()
schema.Bool("active").Default(true).Build()
```

### Numbers
```go
schema.Float32("price").Default(0.0).Build()
schema.Float64("rating").Min(0).Max(5).Build()
schema.Decimal("amount").Precision(10).Scale(2).Build()
```

### Dates and Times
```go
schema.Date("birth_date").Build()
schema.Time("created_at").AutoNowAdd().Build()
schema.DateTime("published_at").Build()
```

### Special Types
```go
schema.Email("email").Required().Build()
schema.URL("website").Blank().Build()
schema.UUID("token").Required().Build()
schema.JSON("metadata").Blank().Build()
schema.Bytes("file_data").Build()
```

## Field Options

Each field type has the options you need:

### Common Options
- `Required()` - Field cannot be null/empty
- `Blank()` - Field can be empty (for strings)
- `Default(value)` - Default value for new records
- `Unique()` - Field must be unique
- `Index()` - Create database index
- `HelpText(text)` - Help text for admin forms

### String Options
- `MaxLength(n)` - Maximum length
- `MinLength(n)` - Minimum length
- `Choices(choices)` - Limit to specific values

### Number Options
- `Min(value)` - Minimum value
- `Max(value)` - Maximum value

### Database Options
- `DBColumn(name)` - Custom column name
- `DBType(type)` - Custom database type
- `DBIndex(type)` - Custom index type

## Relationships

### ForeignKey
```go
schema.ForeignKey("user_id", "User", "id").
    OnDelete(schema.Cascade).
    Build()
```

### OneToOne
```go
schema.OneToOne("profile_id", "Profile", "user_id").
    Build()
```

### ManyToMany
```go
schema.ManyToMany("tags", "Tag", "post_tags").
    Build()
```

### HasMany
```go
schema.HasMany("posts", "Post", "user_id").
    Build()
```

## Model Metadata

Control how your model behaves:

```go
func (User) Meta() schema.Meta {
    return schema.Meta{
        TableName: "users",
        VerboseName: "User",
        VerboseNamePlural: "Users",
        Ordering: []string{"-created_at"},
        Indexes: []schema.Index{
            {Name: "idx_user_email", Fields: []string{"email"}},
        },
        UniqueTogether: []schema.UniqueTogether{
            {Fields: []string{"username", "email"}},
        },
        Permissions: []string{"add_user", "change_user", "delete_user"},
    }
}
```

## Lifecycle Hooks

Run code at specific points in your model's lifecycle:

```go
func (User) Hooks() *schema.ModelHooks {
    return &schema.ModelHooks{
        BeforeSave: func(ctx context.Context, model interface{}) error {
            user := model.(*User)
            // Validate something before saving
            return nil
        },
        AfterSave: func(ctx context.Context, model interface{}) error {
            user := model.(*User)
            // Send welcome email, update cache, etc.
            return nil
        },
        BeforeDelete: func(ctx context.Context, model interface{}) error {
            user := model.(*User)
            // Clean up related data
            return nil
        },
    }
}
```

## Constraints

Add database constraints:

```go
func (User) Constraints() []schema.Constraint {
    return []schema.Constraint{
        schema.CheckConstraint("chk_user_age", "age >= 18"),
        schema.UniqueConstraint("uq_user_email", "email"),
    }
}
```

## What Gets Generated?

When you run `forge generate`, forge creates:

1. **Model struct** - A Go struct with proper database tags
2. **FieldExpr** - Type-safe field accessors for queries
3. **Manager** - CRUD operations (Create, Update, Delete, Get)
4. **QuerySet** - Query building and execution

### Generated Model
```go
type User struct {
    ID        int64     `db:"id"`
    Username  string    `db:"username"`
    Email     string    `db:"email"`
    Bio       string    `db:"bio"`
    IsActive  bool      `db:"is_active"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

### Generated FieldExpr
```go
type UserFields struct {
    ID        FieldExpr[int64]
    Username  FieldExpr[string]
    Email     FieldExpr[string]
    Bio       FieldExpr[string]
    IsActive  FieldExpr[bool]
    CreatedAt FieldExpr[time.Time]
    UpdatedAt FieldExpr[time.Time]
}
```

## Using Your Model

Once generated, you can use your model like this:

```go
// Create a user
user := &User{
    Username: "johndoe",
    Email: "john@example.com",
}
err := User.Objects.Create(ctx, user)

// Query users
users, err := User.Objects.
    Filter(User.Fields.IsActive.Equals(true)).
    OrderBy("-created_at").
    Limit(10).
    All(ctx)

// Get a specific user
user, err := User.Objects.Get(ctx, User.Fields.ID.Equals(1))

// Update a user
err := User.Objects.Filter(User.Fields.ID.Equals(1)).
    Update(ctx, map[string]interface{}{
        "username": "newname",
    })

// Delete a user
err := User.Objects.Filter(User.Fields.ID.Equals(1)).Delete(ctx)
```

## Best Practices

1. **Keep it simple** - Don't over-engineer your models
2. **Use constraints** - Let the database enforce data integrity
3. **Add hooks sparingly** - Only use hooks when you really need them
4. **Document your models** - Use comments to explain complex relationships
5. **Test your models** - Write tests for custom validation and hooks

## Next Steps

- [Code Generation](/docs/features/code-generation) - See how forge generates code
- [ORM System](/docs/features/orm-system) - Learn about the QuerySet API
- [Admin Interface](/docs/features/admin-system) - Auto-generate admin from your models
