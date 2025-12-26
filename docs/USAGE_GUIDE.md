# forge Usage Guide

## Getting Started

### Installation

```bash
# Clone the repository
git clone https://github.com/forgego/forge.git
cd forge

# Install dependencies
go mod download

# Build the CLI
go build ./cmd/forge
```

### Create a New Project

```bash
# Create a new Forge project (interactive prompts)
forge new myproject

# Create with specific template
forge new myproject --template advanced

# Create with options
forge new myproject --database postgres --docker
```

The command will prompt you for:
- Project template (Simple or Advanced)
- Database type (postgres, mysql, sqlite)
- Docker setup (optional)

## Project Structure

Forge supports two project templates:

### Simple Template (Default)
```
myproject/
├── cmd/server/main.go   # Application entry point
├── app/                 # Django-style apps
│   ├── users/
│   │   ├── models.go
│   │   ├── admin.go
│   │   └── api.go
│   └── blog/
│       ├── models.go
│       ├── admin.go
│       └── api.go
├── migrations/          # Database migrations
├── config/              # Configuration
├── static/              # Static files
├── templates/           # HTML templates
└── go.mod
```

### Advanced Template
Includes additional directories:
- `domain/` - Pure business logic
- `infra/` - Infrastructure layer
- `pkg/` - Shared utilities

See [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for details.

## Quick Start

### 1. Create Project

```bash
forge new myproject
cd myproject
```

### 2. Add an App

```bash
forge add app blog
```

This creates `app/blog/` with `models.go`, `admin.go`, and `api.go`.

### 3. Add a Model

```bash
forge add model --app blog --name Post
```

This interactively prompts for fields and adds the model to `app/blog/models.go`.

### 4. Generate Code

```bash
forge generate
```

### 5. Create Migrations

```bash
forge makemigrations
forge migrate
```

### 6. Start Server

```bash
forge runserver
```

## CLI Commands

### Project Management

- `forge new [name]` - Create a new project
- `forge add app [name]` - Add a new app
- `forge add model [name]` - Add a model to an app
- `forge add handler [name]` - Add an HTTP handler
- `forge add api [name]` - Add a REST API endpoint
- `forge add service [name]` - Add a service
- `forge auth` - Scaffold authentication app

### Development

- `forge generate` - Generate code from models
- `forge makemigrations` - Create migration files
- `forge migrate` - Run migrations
- `forge runserver` - Start development server
- `forge test` - Run tests
- `forge shell` - Open interactive shell

## Step-by-Step Tutorial

### Step 1: Define Models

You can either manually create models or use the CLI:

```bash
forge add model --app blog --name Post
```

Or manually create `app/blog/models.go`:

```go
package models

import (
    "context"
    
    "github.com/forgego/forge/internal/schema"
    "github.com/forgego/forge/internal/schema/fields"
    "github.com/forgego/forge/internal/schema/relations"
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
            MaxLength(150),
        fields.String("email").
            Unique().
            Required(),
        fields.String("password").
            Required().
            MaxLength(128),
        fields.Bool("is_active").
            Default(true),
        fields.Time("date_joined").
            AutoNowAdd(),
    }
}

func (User) Relations() []schema.Relation {
    return []schema.Relation{
        relations.OneToMany("posts", "Post").
            RelatedName("author"),
    }
}

func (User) Meta() schema.Meta {
    return schema.Meta{
        TableName: "users",
        OrderBy: []string{"-date_joined"},
        VerboseName: "User",
        VerboseNamePlural: "Users",
    }
}

func (User) Hooks() schema.ModelHooks {
    return schema.ModelHooks{
        BeforeCreate: func(ctx context.Context, instance interface{}) error {
            user := instance.(*User)
            // Hash password, etc.
            return nil
        },
    }
}
```

### Step 2: Generate Code

```bash
# Generate code from schema definitions
forge generate

# This will generate:
# - app/blog/post.gen.go (Post struct)
# - app/blog/post_fields.gen.go (FieldExpr definitions)
# - app/blog/post_manager.gen.go (Manager and QuerySet)
```

This generates:
- `models/user.gen.go` - User struct
- `models/user_fields.gen.go` - FieldExpr definitions
- `models/user_manager.gen.go` - Manager and QuerySet

The framework uses a SQL builder that generates SQL queries with proper escaping and parameter binding automatically.

### Step 4: Use in Your Code

```bash
# Create migration
forge makemigrations

# Apply migrations
forge migrate
```


```go
package main

import (
    "context"
    
    "myproject/models"
    "github.com/forgego/forge/internal/db"
    "github.com/forgego/forge/internal/config"
)

func main() {
    // Load config
    cfg := config.NewConfig()
    
    // Connect to database
    database, err := db.NewDBFromConfig(cfg)
    if err != nil {
        panic(err)
    }
    defer database.Close()
    
    ctx := context.Background()
    
    // Create user
    user := &models.User{
        Username: "john",
        Email: "john@example.com",
        Password: "secret",
    }
    err = models.User.Objects.Create(ctx, user)
    
    // Query users
    users, err := models.User.Objects.Filter(
        models.User.Fields.IsActive.Equals(true),
    ).All(ctx)
    
    // Type-safe query
    user, err := models.User.Objects.Filter(
        models.User.Fields.Username.Equals("john"),
    ).Get(ctx)
}
```

### Step 6: Register Admin

```go
import "github.com/forgego/forge/internal/admin"

func init() {
    admin.RegisterModel(&models.User{})
}
```

### Step 7: Start Server

```bash
forge runserver
```

Visit `http://localhost:8000/admin` to access the admin interface.

## Common Patterns

### Pattern 1: Custom Manager Methods

```go
// In generated manager
func (m *UserManager) GetByEmail(ctx context.Context, email string) (*User, error) {
    return m.Filter(
        User.Fields.Email.Equals(email),
    ).Get(ctx)
}
```

### Pattern 2: Model Methods

```go
// In generated model
func (u *User) GetFullName() string {
    return u.FirstName + " " + u.LastName
}

func (u *User) IsAdmin() bool {
    return u.IsStaff || u.IsSuperuser
}
```

### Pattern 3: Custom Validation

```go
func (User) Hooks() schema.ModelHooks {
    return schema.ModelHooks{
        Clean: func(ctx context.Context, instance interface{}) error {
            user := instance.(*User)
            if user.Username == "" {
                return errors.New("username is required")
            }
            if len(user.Username) < 3 {
                return errors.New("username must be at least 3 characters")
            }
            return nil
        },
    }
}
```

### Pattern 4: Transactions

```go
err := db.WithTx(ctx, func(tx *db.Tx) error {
    user := &User{Username: "john"}
    if err := User.Objects.Create(ctx, user); err != nil {
        return err
    }
    
    post := &Post{Author: user, Title: "Hello"}
    return Post.Objects.Create(ctx, post)
})
```

### Pattern 5: Complex Queries

```go
// Active staff users joined in last month
users, err := User.Objects.Filter(
    User.Fields.IsActive.Equals(true).And(
        User.Fields.IsStaff.Equals(true),
    ).And(
        User.Fields.DateJoined.Greater(
            time.Now().AddDate(0, -1, 0),
        ),
    ),
).OrderBy("-date_joined").Limit(10).All(ctx)
```

## Configuration

### config.yaml

```yaml
app:
  name: "My Application"
  env: "development"
  debug: true

server:
  host: "localhost"
  port: "8000"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: ""
  name: "mydb"
  sslmode: "disable"

security:
  secret_key: "change-me-in-production"
  csrf_secret_key: "change-me-in-production"
  session_secret: "change-me-in-production"

admin:
  enabled: true
  path: "/admin"
  title: "My Admin"
```

### Environment Variables

```bash
export FORGE_DATABASE_HOST=localhost
export FORGE_DATABASE_NAME=mydb
export FORGE_APP_ENV=production
```

## Best Practices

### 1. Schema Definitions

- Keep schemas in separate files
- Use descriptive field names
- Add help text for clarity
- Use appropriate field types

### 2. Code Generation

- Run `forge generate` after schema changes
- Commit generated files to version control
- Review generated code

### 3. Queries

- Prefer type-safe API
- Use SelectRelated for JOINs
- Use PrefetchRelated for separate queries
- Add indexes for frequently queried fields

### 4. Security

- Always hash passwords in BeforeCreate hook
- Use CSRF protection for forms
- Validate all user input
- Use parameterized queries

### 5. Performance

- Use Select() to limit fields
- Use PrefetchRelated() efficiently
- Add database indexes
- Use transactions for multiple operations

## Troubleshooting

### Issue: Code generation fails

**Solution:**
- Check schema syntax
- Ensure all imports are correct
- Run `go mod tidy`

### Issue: Migrations fail

**Solution:**
- Check database connection
- Verify migration files
- Check for conflicting migrations

### Issue: QuerySet methods not working

**Solution:**
- Check database connection
- Verify model registration
- Ensure code generation was run (`forge generate`)

## Examples

See `examples/blog/` for a complete example application.

