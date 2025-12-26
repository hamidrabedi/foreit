# Admin Interface Usage Guide

The forge admin interface provides a Django-like auto-generated admin panel for your models.

## Quick Start

### 1. Register Models with Managers

To use the admin interface, you need to register your models along with their managers:

```go
package main

import (
    "github.com/forgego/forge/internal/admin"
    "github.com/forgego/forge/internal/db"
    "your-app/models"
)

func main() {
    // Set up database connection
    database := db.NewDBFromConfig(cfg)

    // Set database on managers
    models.User.Objects.SetDB(database)
    models.Post.Objects.SetDB(database)

    // Register models with admin (including managers)
    admin.RegisterModelWithManager(&models.User{}, models.User.Objects)
    admin.RegisterModelWithManager(&models.Post{}, models.Post.Objects)

    // Or with custom options
    admin.RegisterModelWithOptions(
        &models.User{},
        admin.WithManager(models.User.Objects),
        admin.WithListDisplay("username", "email", "is_active"),
        admin.WithSearchFields("username", "email"),
    )
}
```

### 2. Register Admin Routes

```go
import (
    httplib "github.com/forgego/forge/internal/http"
    "github.com/forgego/forge/internal/admin"
)

func main() {
    router := httplib.NewRouter()

    // Register admin routes
    admin.RegisterAdminRoutes(router, "/admin")

    // ... rest of your setup
}
```

## Features

### List View

- **Pagination**: Automatic pagination with configurable page size
- **Search**: Full-text search on specified fields
- **Sorting**: Click column headers to sort
- **Inline Editing**: Double-click cells to edit inline
- **Bulk Actions**: Select multiple items for bulk operations (coming soon)

### Detail View

- **Read-only Display**: View all fields of a model instance
- **Related Objects**: View related objects (coming soon)

### Create/Edit Forms

- **Auto-generated Forms**: Forms automatically generated from model schema
- **Field Types**: Proper input types for different field types (text, number, date, etc.)
- **Validation**: Client and server-side validation
- **Rich Text Editing**: TinyMCE for textarea fields
- **Date Pickers**: Flatpickr for date/time fields
- **Enhanced Dropdowns**: Select2 for select fields

### Delete

- **Confirmation**: Browser confirmation dialog
- **HTMX Integration**: Smooth row removal with HTMX

## Admin Options

### List Display

Control which fields appear in the list view:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithManager(models.User.Objects),
    admin.WithListDisplay("username", "email", "is_active", "date_joined"),
)
```

### Search Fields

Specify which fields are searchable:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithManager(models.User.Objects),
    admin.WithSearchFields("username", "email", "first_name", "last_name"),
)
```

### List Filters

Add filters to the sidebar:

```go
admin.RegisterModelWithOptions(
    &models.Post{},
    admin.WithManager(models.Post.Objects),
    admin.WithListFilter("status", "category", "author"),
)
```

### Read-only Fields

Mark fields as read-only in forms:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithManager(models.User.Objects),
    admin.WithReadOnlyFields("id", "date_joined", "last_login"),
)
```

## Custom Admin Classes

For more control, implement a custom admin class:

```go
type UserAdmin struct{}

func (a *UserAdmin) GetListDisplay() []interface{} {
    return []interface{}{"username", "email", "is_active"}
}

func (a *UserAdmin) GetListFilter() []interface{} {
    return []interface{}{"is_active", "is_staff"}
}

func (a *UserAdmin) GetSearchFields() []interface{} {
    return []interface{}{"username", "email"}
}

func (a *UserAdmin) GetReadOnlyFields() []interface{} {
    return []interface{}{"id", "date_joined"}
}

// Register with custom admin
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithManager(models.User.Objects),
    admin.WithCustomAdmin(&UserAdmin{}),
)
```

## HTMX Features

The admin interface uses HTMX for dynamic interactions:

- **Dynamic Search**: Search as you type (500ms debounce)
- **Inline Editing**: Double-click to edit, blur or Enter to save
- **Pagination**: Click page numbers to load new pages without full reload
- **Delete**: Smooth row removal with confirmation

## Styling

The admin uses Bootstrap 5 for styling. All templates can be customized by:

1. Overriding template files
2. Adding custom CSS
3. Using Bootstrap utility classes

## Security

- **CSRF Protection**: All forms include CSRF tokens
- **XSS Protection**: All field values are HTML-escaped
- **Authentication**: Add authentication middleware (coming soon)

## Examples

See `examples/` directory for complete examples of admin usage.

## Troubleshooting

### "Manager not set" Error

Make sure you register models with their managers:

```go
admin.RegisterModelWithManager(&models.User{}, models.User.Objects)
```

### "Model not registered" Error

Register your models before starting the server:

```go
admin.RegisterModelWithManager(&models.User{}, models.User.Objects)
admin.RegisterAdminRoutes(router, "/admin")
```

### Database Connection Not Set

Set the database connection on managers:

```go
database := db.NewDBFromConfig(cfg)
models.User.Objects.SetDB(database)
```

## Next Steps

- See [HTMX Patterns](HTMX_PATTERNS.md) for advanced patterns
- See [REST API Guide](REST_API.md) for API integration
- See [Architecture](ARCHITECTURE.md) for framework internals
