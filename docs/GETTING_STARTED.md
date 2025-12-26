# Getting Started with forge

This guide will help you create your first forge application in 10 minutes.

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Basic knowledge of Go

## Step 1: Create a New Project

```bash
forge new myapp
cd myapp
```

This creates a new project with the following structure:

```
myapp/
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
  dbname: myapp_db
  sslmode: disable
```

Create the database:

```sql
CREATE DATABASE myapp_db;
```

## Step 3: Define Your Models

Edit `models/example.go` or create new model files. Here's a complete example:

```go
package models

import (
	"github.com/forgego/forge/internal/schema"
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
		schema.String("content").Required().Build(),
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

	"github.com/forgego/forge/internal/admin"
	"github.com/forgego/forge/internal/config"
	"github.com/forgego/forge/internal/db"
	"github.com/forgego/forge/internal/logging"
	httplib "github.com/forgego/forge/internal/http"
	"github.com/forgego/forge/internal/registry"
	"your-module/models"
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
			fmt.Fprintf(w, "Welcome to MyApp!")
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

### Type-Safe Queries

```go
import (
	"context"
	"your-module/models"
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
	Published: true,
}
err := models.Post.Objects.Create(ctx, newPost)

// Update a post
post.Title = "Updated Title"
err := models.Post.Objects.Update(ctx, post)

// Delete a post
err := models.Post.Objects.Delete(ctx, post)
```

### Complex Queries

```go
// Filter with multiple conditions
posts, err := models.Post.Objects.
	Filter(
		models.Post.Fields.Published.Equals(true).
			And(models.Post.Fields.CreatedAt.GreaterThan(someDate)),
	).
	OrderBy("-created_at").
	Limit(10).
	All(ctx)

// Count posts
count, err := models.Post.Objects.
	Filter(models.Post.Fields.Published.Equals(true)).
	Count(ctx)

// Check if any posts exist
exists, err := models.Post.Objects.
	Filter(models.Post.Fields.Published.Equals(true)).
	Exists(ctx)
```

## Next Steps

- Read the [API Reference](API_REFERENCE.md) for detailed API documentation
- Check out [examples](../examples/) for complete example applications
- Learn about [Plugins](../docs/ARCHITECTURE.md#plugins) to extend the framework
- See [Schema Reference](SCHEMA_REFERENCE.md) for all field types and options

## Troubleshooting

### Database Connection Errors

Make sure:

- PostgreSQL is running
- Database credentials in `config.yaml` are correct
- Database exists (run `CREATE DATABASE myapp_db;`)

### Code Generation Errors

- Make sure your model files are in the `models/` directory
- Check that model structs embed `schema.BaseSchema`
- Verify `Fields()`, `Meta()`, `Relations()`, and `Hooks()` methods are defined

### Admin Not Showing

- Make sure models are registered with `admin.RegisterModel()`
- Check that `admin.enabled: true` in `config.yaml`
- Verify admin path is correct in configuration

## Example Applications

- **[Blog Example](../examples/blog/)** - Simple blog with User model
- **[Demo Example](../examples/demo/)** - Full-featured demo with multiple models
- **[Full Featured Example](../examples/fullfeatured/)** - Complete example with all features
