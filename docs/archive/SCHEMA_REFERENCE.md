# forge Schema Reference

> **Note:** This is the complete reference for defining schemas in forge. For architecture details, see [Architecture](ARCHITECTURE.md).

## Schema Definition

A schema is defined by implementing the `Schema` interface:

```go
type Schema interface {
    Fields() []Field
    Relations() []Relation
    Meta() Meta
    Hooks() ModelHooks
}
```

## Field Types

### Basic Field Types

```go
// Numeric
fields.Int64("id").Primary().AutoIncrement()
fields.Int32("age")
fields.Float64("price")

// String
fields.String("username").MaxLength(150).Required()
fields.Text("description")
fields.Email("email").Unique()
fields.URL("website")
fields.UUID("id").Primary()

// Boolean
fields.Bool("is_active").Default(true)

// Time
fields.Time("created_at").AutoNow()
fields.Date("birth_date")
fields.DateTime("last_login")
```

### Field Options

All field builders support chainable methods:

```go
fields.String("username")
    .Required()              // NOT NULL
    .Unique()                // UNIQUE constraint
    .Primary()               // PRIMARY KEY
    .Index()                 // db_index
    .DBColumn("user_name")  // Custom column name
    .DBIndex()               // Create index
    .Default("guest")        // Default value
    .MaxLength(150)         // Maximum length
    .MinLength(3)           // Minimum length
    .HelpText("Username")   // Help text
    .VerboseName("Username") // Human-readable name
    .Choices(...)            // Predefined choices
    .Null()                  // Allow NULL
    .Blank()                 // Allow blank (for strings)
    .Editable(false)         // Read-only in admin
    .AutoNow()               // Set on save
    .AutoNowAdd()            // Set on create only
```

## Relations

### ForeignKey (ManyToOne)

```go
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

### OneToMany

```go
func (User) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToMany("posts", "Post").
            RelatedName("author"),
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

### Relation Options

```go
relations.ForeignKey("author", "User")
    .Required()                    // Required relation
    .OnDelete(schema.Cascade)     // CASCADE, SET_NULL, SET_DEFAULT, RESTRICT, DO_NOTHING
    .OnUpdate(schema.Cascade)     // Same options
    .RelatedName("posts")         // Reverse relation name
    .Through("post_tags")         // Through table (M2M only)
    .FromField("post_id")         // Custom from field
    .ToField("tag_id")            // Custom to field
```

## Meta Options

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
        Constraints: []schema.UniqueConstraint{ // Check constraints
            {
                Name: "check_age",
                Check: "age >= 0",
            },
        },
        Permissions: []schema.Permission{ // Model permissions
            {
                Codename: "can_view",
                Name: "Can view user",
            },
        },
        DefaultPermissions: true,         // Add default CRUD permissions
    }
}
```

## Model Hooks

```go
func (User) Hooks() schema.ModelHooks {
    return schema.ModelHooks{
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
        Clean: func(ctx context.Context, instance interface{}) error {
            // Custom validation
            return nil
        },
        Validate: func(ctx context.Context, instance interface{}) error {
            // Additional validation
            return nil
        },
    }
}
```

## Complete Example

```go
package models

import (
    "context"
    "time"
    
    "github.com/forgego/forge/pkg/schema"
    "github.com/forgego/forge/pkg/schema/fields"
    "github.com/forgego/forge/pkg/schema/relations"
)

type User struct {
    schema.Schema
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        fields.Int64("id").Primary().AutoIncrement(),
        fields.String("username").
            Unique().
            Required().
            MaxLength(150).
            HelpText("Required. 150 characters or fewer."),
        fields.String("email").
            Unique().
            Required().
            MaxLength(255).
            HelpText("Required. Must be a valid email address."),
        fields.String("password").
            Required().
            MaxLength(128).
            HelpText("Required. Password for authentication."),
        fields.Bool("is_active").
            Default(true).
            HelpText("Designates whether this user should be treated as active."),
        fields.Bool("is_staff").
            Default(false).
            HelpText("Designates whether this user can access the admin site."),
        fields.Bool("is_superuser").
            Default(false).
            HelpText("Designates that this user has all permissions."),
        fields.Time("date_joined").
            AutoNowAdd().
            HelpText("Date when user was created."),
        fields.Time("last_login").
            Null().
            HelpText("Last login timestamp."),
    }
}

func (User) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToMany("posts", "Post").
            RelatedName("author"),
        relations.OneToOne("profile", "UserProfile").
            OnDelete(schema.Cascade),
        relations.ManyToMany("groups", "Group").
            Through("user_groups").
            RelatedName("users"),
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
        UniqueTogether: [][]string{
            {"username", "email"},
        },
    }
}

func (User) Hooks() schema.ModelHooks {
    return schema.ModelHooks{
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
        Clean: func(ctx context.Context, instance interface{}) error {
            user := instance.(*User)
            // Custom validation
            if user.Username == "" {
                return errors.New("username is required")
            }
            return nil
        },
    }
}
```

## Field Type Reference

| Field Type | Go Type | SQL Type | Builder |
|------------|---------|----------|---------|
| Int64 | `int64` | `BIGINT` | `fields.Int64()` |
| Int32 | `int32` | `INTEGER` | `fields.Int32()` |
| String | `string` | `TEXT` or `VARCHAR(n)` | `fields.String()` |
| Text | `string` | `TEXT` | `fields.Text()` |
| Bool | `bool` | `BOOLEAN` | `fields.Bool()` |
| Float64 | `float64` | `DOUBLE PRECISION` | `fields.Float64()` |
| Time | `time.Time` | `TIMESTAMP` | `fields.Time()` |
| Date | `time.Time` | `DATE` | `fields.Date()` |
| DateTime | `time.Time` | `TIMESTAMP` | `fields.DateTime()` |
| Email | `string` | `VARCHAR(255)` | `fields.Email()` |
| URL | `string` | `VARCHAR(255)` | `fields.URL()` |
| UUID | `string` | `UUID` or `VARCHAR(36)` | `fields.UUID()` |
| JSON | `[]byte` | `JSONB` | `fields.JSON()` |
| Bytes | `[]byte` | `BYTEA` | `fields.Bytes()` |

## Relation Type Reference

| Relation Type | Description | Reverse Accessor |
|---------------|-------------|------------------|
| ForeignKey | Many-to-One | `user.posts` (OneToMany) |
| OneToOne | One-to-One | `user.profile` ↔ `profile.user` |
| OneToMany | One-to-Many | `user.posts` (collection) |
| ManyToMany | Many-to-Many | `post.tags` ↔ `tag.posts` |

## Cascade Options

| Option | SQL | Description |
|--------|-----|-------------|
| `Cascade` | `ON DELETE CASCADE` | Delete related objects |
| `SetNull` | `ON DELETE SET NULL` | Set foreign key to NULL |
| `SetDefault` | `ON DELETE SET DEFAULT` | Set foreign key to default |
| `Restrict` | `ON DELETE RESTRICT` | Prevent deletion if related objects exist |
| `DoNothing` | `ON DELETE NO ACTION` | No action (database default) |

