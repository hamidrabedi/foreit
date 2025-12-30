---
sidebar_position: 3
description: Learn how to use and customize forge's auto-generated admin interface. Django-style admin panel with CRUD operations, search, filtering, and more.
keywords:
  - forge admin
  - admin interface
  - django admin go
  - auto-generated admin
  - forge admin panel
image: /img/forge-social-card.jpg
---

# Admin Interface

The admin interface is forge's auto-generated admin panel. Register your models and you get a full web UI for managing your data. No code required.

## Why use it?

Because building admin interfaces is boring:

- **Zero code** - Just register your models
- **Full CRUD** - Create, read, update, delete—all there
- **Search and filters** - Find what you need fast
- **Bulk actions** - Update or delete multiple records at once
- **Export** - Download your data as CSV or JSON

## Quick start

### 1. Register Models

In your `main.go`:

```go
import (
    "github.com/forgego/forge/pkg/admin"
    "myapp/models"
)

func main() {
    // ... your setup code ...
    
    admin.RegisterModel(&models.User{})
    admin.RegisterModel(&models.Post{})
    
    if settings.Admin.Enabled {
        admin.RegisterAdminRoutes(router, settings.Admin.Path)
    }
}
```

### 2. Access Admin

Start your server and visit `http://localhost:8000/admin/`

You'll see:
- A list of all registered models
- Links to view, add, and manage each model
- Search and filtering capabilities

## Admin Options

### List Display

Control which fields appear in the list view:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithListDisplay("username", "email", "is_active", "date_joined"),
)
```

### Search Fields

Enable search on specific fields:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithSearchFields("username", "email"),
)
```

### List Filter

Add filters to the sidebar:

```go
admin.RegisterModelWithOptions(
    &models.Post{},
    admin.WithListFilter("published", "created_at", "author"),
)
```

### Date Hierarchy

Add date-based navigation:

```go
admin.RegisterModelWithOptions(
    &models.Post{},
    admin.WithDateHierarchy("created_at"),
)
```

### Ordering

Set default ordering:

```go
admin.RegisterModelWithOptions(
    &models.Post{},
    admin.WithOrdering("-created_at", "title"),
)
```

### Read Only Fields

Make fields read-only in forms:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithReadOnlyFields("date_joined", "last_login"),
)
```

### Fieldsets

Organize form fields into groups:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithFieldsets(
        admin.Fieldset("Personal Information", "username", "email"),
        admin.Fieldset("Permissions", "is_active", "is_staff", "is_superuser"),
        admin.Fieldset("Important Dates", "date_joined", "last_login"),
    ),
)
```

## Customizing Admin

### Custom List View

You can customize the list view by creating a custom admin class:

```go
type UserAdmin struct {
    *admin.ModelAdmin
}

func (a *UserAdmin) ListDisplay() []string {
    return []string{"username", "email", "is_active", "date_joined"}
}

func (a *UserAdmin) SearchFields() []string {
    return []string{"username", "email"}
}

admin.RegisterModelWithAdmin(&models.User{}, &UserAdmin{})
```

### Custom Form

Customize the form for creating/editing:

```go
type UserAdmin struct {
    *admin.ModelAdmin
}

func (a *UserAdmin) GetFormFields() []string {
    return []string{"username", "email", "password", "is_active"}
}

func (a *UserAdmin) GetExcludeFields() []string {
    return []string{"date_joined", "last_login"}
}
```

### Custom Actions

Add custom bulk actions:

```go
type PostAdmin struct {
    *admin.ModelAdmin
}

func (a *PostAdmin) GetActions() []admin.Action {
    return []admin.Action{
        {
            Name: "publish",
            Label: "Publish selected posts",
            Handler: func(ids []int64) error {
                // Publish logic
                return nil
            },
        },
        {
            Name: "unpublish",
            Label: "Unpublish selected posts",
            Handler: func(ids []int64) error {
                // Unpublish logic
                return nil
            },
        },
    }
}
```

## Features

### List View

The list view provides:

- **Pagination** - Automatic pagination with configurable page size
- **Search** - Full-text search on specified fields
- **Sorting** - Click column headers to sort
- **Filtering** - Filter by field values
- **Inline Editing** - Double-click cells to edit inline (coming soon)
- **Bulk Actions** - Select multiple items for bulk operations

### Detail View

View all fields of a model instance:

- Read-only display of all fields
- Related objects display
- History/audit trail (coming soon)

### Create/Edit Forms

Auto-generated forms with:

- Proper input types for different field types
- Client and server-side validation
- Rich text editing for textarea fields
- Date pickers for date/time fields
- Enhanced dropdowns for select fields
- File upload support (coming soon)

### Delete

Delete with confirmation:

- Browser confirmation dialog
- Smooth row removal with HTMX
- Cascade delete handling

## Permissions

Control access to admin:

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithPermissions(
        admin.RequirePermission("can_view_user"),
        admin.RequirePermission("can_change_user"),
    ),
)
```

## Custom Templates

You can override admin templates by placing them in your project's `templates/admin/` directory:

```
templates/
└── admin/
    ├── base.html
    ├── index.html
    ├── list.html
    ├── form.html
    └── detail.html
```

## Styling

The admin uses Bootstrap 5 and can be customized with CSS:

```css
/* Custom admin styles */
.admin-header {
    background-color: #your-brand-color;
}
```

## Best Practices

1. **Use List Display** - Show only relevant fields in list view
2. **Enable Search** - Add search fields for frequently searched columns
3. **Use Filters** - Add filters for fields users commonly filter by
4. **Organize Fieldsets** - Group related fields together
5. **Set Read-Only** - Make auto-generated fields read-only

## Examples

### Blog Admin

```go
admin.RegisterModelWithOptions(
    &models.Post{},
    admin.WithListDisplay("title", "author", "published", "created_at"),
    admin.WithSearchFields("title", "content"),
    admin.WithListFilter("published", "author", "created_at"),
    admin.WithDateHierarchy("created_at"),
    admin.WithOrdering("-created_at"),
)
```

### User Admin

```go
admin.RegisterModelWithOptions(
    &models.User{},
    admin.WithListDisplay("username", "email", "is_active", "is_staff", "date_joined"),
    admin.WithSearchFields("username", "email"),
    admin.WithListFilter("is_active", "is_staff", "date_joined"),
    admin.WithFieldsets(
        admin.Fieldset("Personal Information", "username", "email"),
        admin.Fieldset("Permissions", "is_active", "is_staff", "is_superuser"),
    ),
    admin.WithReadOnlyFields("date_joined", "last_login"),
)
```

## Next Steps

- [REST API Guide](/docs/guides/rest-api) - Build APIs for your frontend
- [Security Guide](/docs/guides/security) - Secure your admin interface
- [Advanced Topics](/docs/advanced/plugins) - Extend the admin with plugins

