---
sidebar_position: 2
---

# Quick Start

Get up and running with forge in 10 minutes. This tutorial will walk you through creating a simple blog application.

## Step 1: Create a New Project

Use the forge CLI to create a new project:

```bash
forge new myblog
cd myblog
```

This creates a new project with the following structure:

```
myblog/
├── main.go              # Application entry point
├── go.mod               # Go module file
├── config/
│   └── config.yaml      # Configuration file
├── models/              # Your model definitions
│   └── example.go       # Example model
└── migrations/          # Database migrations
```

## Step 2: Configure Database

Edit `config/config.yaml` and set your database connection:

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: myblog_db
  sslmode: disable

server:
  host: localhost
  port: "8000"

admin:
  enabled: true
  path: "/admin"
```

Create the database:

```bash
psql -U postgres -c "CREATE DATABASE myblog_db;"
```

## Step 3: Define Your Models

Edit `models/example.go` or create `models/post.go`:

```go
package models

import (
    "github.com/forgego/forge/pkg/schema"
)

// Post represents a blog post
type Post struct {
    schema.BaseSchema
}

// Fields returns all field definitions
func (Post) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("title").Required().MaxLength(200).Build(),
        schema.String("slug").Unique().MaxLength(200).Build(),
        schema.Text("content").Required().Build(),
        schema.Bool("published").Default(false).Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
        schema.Time("updated_at").AutoNow().Build(),
    }
}

// Meta returns model metadata
func (Post) Meta() schema.Meta {
    return schema.Meta{
        TableName:        "posts",
        VerboseName:      "Post",
        VerboseNamePlural: "Posts",
        OrderBy:          []string{"-created_at"},
    }
}

// Relations returns relationship definitions
func (Post) Relations() []schema.Relation {
    return []schema.Relation{}
}

// Hooks returns model lifecycle hooks
func (Post) Hooks() *schema.ModelHooks {
    return nil
}
```

## Step 4: Generate Code

Generate type-safe code from your model definitions:

```bash
forge generate
```

This creates:
- `models/post.gen.go` - Generated model struct
- `models/post_fields.gen.go` - Type-safe field expressions
- `models/post_manager.gen.go` - Manager with CRUD operations
- `models/post_queryset.gen.go` - QuerySet for filtering

## Step 5: Register Models

Update `main.go` to register your models:

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/forgego/forge/pkg/admin"
    "github.com/forgego/forge/pkg/config"
    "github.com/forgego/forge/pkg/db"
    "github.com/forgego/forge/pkg/logging"
    httplib "github.com/forgego/forge/pkg/http"
    "github.com/forgego/forge/pkg/registry"
    "myblog/models"
)

func main() {
    // Load configuration
    cfg := config.NewConfig()
    settings := config.LoadSettings(cfg)

    // Create logger
    logger, err := logging.NewLogger(cfg.IsDevelopment())
    if err != nil {
        log.Fatal(err)
    }
    defer logger.Sync()

    // Connect to database
    database, err := db.NewDBFromConfig(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    // Register models
    registry.RegisterModel(&models.Post{})

    // Register admin models
    admin.RegisterModel(&models.Post{})

    // Create server
    server, err := httplib.NewServer(cfg, settings, logger)
    if err != nil {
        log.Fatal(err)
    }

    // Register routes
    server.RegisterRoutes(func(router *httplib.Router) {
        router.Get("/", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "Welcome to MyBlog!")
        })

        // Register admin routes
        if settings.Admin.Enabled {
            admin.RegisterAdminRoutes(router, settings.Admin.Path)
        }
    })

    // Start server
    fmt.Printf("Starting server on %s:%s\n", settings.Server.Host, settings.Server.Port)
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Step 6: Run Migrations

Create and apply database migrations:

```bash
forge makemigrations
forge migrate
```

This creates the database tables for your models.

## Step 7: Start the Server

```bash
forge runserver
```

Or:

```bash
go run main.go
```

## Step 8: Access Admin Interface

Visit `http://localhost:8000/admin/` to access the auto-generated admin interface.

You can:
- View all posts in a table
- Create new posts
- Edit existing posts
- Delete posts
- Search and filter posts

## Using the ORM

Now you can use the ORM in your code:

```go
import (
    "context"
    "myblog/models"
)

ctx := context.Background()

// Get all published posts
posts, err := models.Post.Objects.
    Filter(models.Post.Fields.Published.Equals(true)).
    OrderBy("-created_at").
    All(ctx)

// Get a single post
post, err := models.Post.Objects.Get(ctx, 1)

// Create a new post
newPost := &models.Post{
    Title:     "My First Post",
    Content:   "This is the content...",
    Slug:      "my-first-post",
    Published: true,
}
err := models.Post.Objects.Create(ctx, newPost)

// Update a post
post.Title = "Updated Title"
err := models.Post.Objects.Update(ctx, post)

// Delete a post
err := models.Post.Objects.Delete(ctx, post)
```

## What's Next?

Congratulations! You've created your first forge application. Now you can:

- [Learn about Models](guides/models) - Deep dive into model definitions
- [Explore Queries](guides/queries) - Learn about QuerySet and filtering
- [Customize Admin](guides/admin) - Customize the admin interface
- [Build REST APIs](guides/rest-api) - Create APIs for your frontend
- [Check out Examples](examples/blog) - See complete example applications

