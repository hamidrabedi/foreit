---
sidebar_position: 1
---

# Models Guide

Models are the foundation of your forge application. They define your data structure and provide the interface for interacting with your database.

## What is a Model?

A model is a Go struct that implements the `Schema` interface. It defines:

- **Fields** - The data columns in your database table
- **Relations** - Relationships to other models
- **Meta** - Metadata like table name, ordering, indexes
- **Hooks** - Lifecycle callbacks for create, update, delete operations

## Defining a Model

Here's a basic model definition:

```go
package models

import (
    "github.com/forgego/forge/pkg/schema"
)

type Post struct {
    schema.BaseSchema
}

func (Post) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("title").Required().MaxLength(200).Build(),
        schema.Text("content").Required().Build(),
        schema.Bool("published").Default(false).Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
    }
}

func (Post) Meta() schema.Meta {
    return schema.Meta{
        TableName:        "posts",
        VerboseName:      "Post",
        VerboseNamePlural: "Posts",
    }
}

func (Post) Relations() []schema.Relation {
    return []schema.Relation{}
}

func (Post) Hooks() *schema.ModelHooks {
    return nil
}
```

## Field Types

forge supports many field types:

### Numeric Fields

```go
schema.Int64("id").Primary().AutoIncrement().Build()
schema.Int32("age").Build()
schema.Float64("price").Build()
schema.Decimal("amount").MaxDigits(10).DecimalPlaces(2).Build()
```

### String Fields

```go
schema.String("username").MaxLength(150).Required().Build()
schema.Text("description").Build()
schema.Email("email").Unique().Required().Build()
schema.URL("website").Build()
schema.Slug("slug").MaxLength(200).Build()
```

### Boolean Fields

```go
schema.Bool("is_active").Default(true).Build()
schema.Bool("is_staff").Default(false).Build()
```

### Date and Time Fields

```go
schema.Time("created_at").AutoNowAdd().Build()
schema.Time("updated_at").AutoNow().Build()
schema.Date("birth_date").Build()
schema.DateTime("last_login").Build()
```

### Special Fields

```go
schema.UUID("id").Primary().Build()
schema.JSON("metadata").Build()
schema.Bytes("avatar").Build()
```

## Field Options

All field types support chainable options:

```go
schema.String("username")
    .Required()              // NOT NULL constraint
    .Unique()               // UNIQUE constraint
    .Primary()              // PRIMARY KEY
    .Index()                // Create database index
    .DBColumn("user_name")  // Custom column name
    .Default("guest")       // Default value
    .MaxLength(150)         // Maximum length (strings)
    .MinLength(3)           // Minimum length (strings)
    .HelpText("Username")   // Help text for admin
    .VerboseName("Username") // Human-readable name
    .Choices(...)           // Predefined choices
    .Null()                 // Allow NULL
    .Blank()                // Allow blank (strings)
    .Editable(false)        // Read-only in admin
    .AutoNow()              // Set on save
    .AutoNowAdd()           // Set on create only
    .Build()
```

## Relations

### ForeignKey (Many-to-One)

```go
import "github.com/forgego/forge/pkg/schema/relations"

func (Post) Relations() []schema.Relation {
    return []schema.Relation{
        relations.ForeignKey("author", "User").
            Required().
            OnDelete(schema.Cascade).
            RelatedName("posts"),
    }
}
```

### OneToOne

```go
func (User) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToOne("profile", "UserProfile").
            OnDelete(schema.Cascade),
    }
}
```

### ManyToMany

```go
func (Post) Relations() []schema.Relation {
    return []schema.Relation{
        relations.ManyToMany("tags", "Tag").
            Through("post_tags").
            RelatedName("posts"),
    }
}
```

## Meta Options

The `Meta()` method returns metadata about your model:

```go
func (User) Meta() schema.Meta {
    return schema.Meta{
        TableName: "users",              // Custom table name
        OrderBy: []string{"-date_joined"}, // Default ordering
        VerboseName: "User",              // Singular name
        VerboseNamePlural: "Users",       // Plural name
        Indexes: []schema.Index{          // Custom indexes
            {
                Name: "idx_user_email",
                Fields: []string{"email"},
                Unique: false,
            },
        },
        UniqueTogether: [][]string{       // Unique constraints
            {"username", "email"},
        },
    }
}
```

## Model Hooks

Hooks allow you to run code at specific points in a model's lifecycle:

```go
func (User) Hooks() *schema.ModelHooks {
    return &schema.ModelHooks{
        BeforeCreate: func(ctx context.Context, instance interface{}) error {
            user := instance.(*User)
            // Hash password, set defaults, etc.
            return nil
        },
        AfterCreate: func(ctx context.Context, instance interface{}) error {
            // Send notification, log, etc.
            return nil
        },
        BeforeUpdate: func(ctx context.Context, instance interface{}) error {
            // Validate changes, etc.
            return nil
        },
        AfterUpdate: func(ctx context.Context, instance interface{}) error {
            // Update cache, etc.
            return nil
        },
        BeforeSave: func(ctx context.Context, instance interface{}) error {
            // Common logic for create/update
            return nil
        },
        AfterSave: func(ctx context.Context, instance interface{}) error {
            // Common logic for create/update
            return nil
        },
        BeforeDelete: func(ctx context.Context, instance interface{}) error {
            // Check dependencies, etc.
            return nil
        },
        AfterDelete: func(ctx context.Context, instance interface{}) error {
            // Cleanup, etc.
            return nil
        },
    }
}
```

## Complete Example

Here's a complete model with all features:

```go
package models

import (
    "context"
    "github.com/forgego/forge/pkg/schema"
    "github.com/forgego/forge/pkg/schema/relations"
)

type User struct {
    schema.BaseSchema
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("username").
            Unique().
            Required().
            MaxLength(150).
            HelpText("Required. 150 characters or fewer.").
            Build(),
        schema.String("email").
            Unique().
            Required().
            MaxLength(255).
            HelpText("Required. Must be a valid email address.").
            Build(),
        schema.String("password").
            Required().
            MaxLength(128).
            HelpText("Required. Password for authentication.").
            Build(),
        schema.Bool("is_active").
            Default(true).
            HelpText("Designates whether this user should be treated as active.").
            Build(),
        schema.Bool("is_staff").
            Default(false).
            HelpText("Designates whether this user can access the admin site.").
            Build(),
        schema.Time("date_joined").
            AutoNowAdd().
            HelpText("Date when user was created.").
            Build(),
        schema.Time("last_login").
            Null().
            HelpText("Last login timestamp.").
            Build(),
    }
}

func (User) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToMany("posts", "Post").
            RelatedName("author"),
        relations.OneToOne("profile", "UserProfile").
            OnDelete(schema.Cascade),
    }
}

func (User) Meta() schema.Meta {
    return schema.Meta{
        TableName: "users",
        OrderBy: []string{"-date_joined"},
        VerboseName: "User",
        VerboseNamePlural: "Users",
        Indexes: []schema.Index{
            {
                Name: "idx_user_email",
                Fields: []string{"email"},
            },
        },
    }
}

func (User) Hooks() *schema.ModelHooks {
    return &schema.ModelHooks{
        BeforeCreate: func(ctx context.Context, instance interface{}) error {
            user := instance.(*User)
            // Hash password
            if user.Password != "" {
                hashed, err := auth.HashPassword(user.Password)
                if err != nil {
                    return err
                }
                user.Password = hashed
            }
            return nil
        },
    }
}
```

## Best Practices

1. **Use descriptive field names** - Make your code self-documenting
2. **Set appropriate constraints** - Use `Required()`, `Unique()`, `MaxLength()` where needed
3. **Use indexes wisely** - Index frequently queried fields
4. **Keep hooks simple** - Hooks should be fast and not block
5. **Use relations appropriately** - Choose the right relation type for your use case

## Next Steps

- [Learn about Queries](/docs/guides/queries) - Query your models
- [Explore Relations](/docs/reference/relations) - Deep dive into relationships
- [Check Field Reference](/docs/reference/fields) - Complete field type reference

