# Admin System - Type-Safe Admin Interface

The forge admin system provides a fully type-safe, Django-inspired admin interface using Go generics. It's production-ready with HTTP handlers, widgets, export, and error handling.

## Features

- ✅ **Type-Safe**: Full compile-time type checking with generics
- ✅ **Easy to Use**: Simple, intuitive API
- ✅ **Fully Featured**: Comparable to Django Admin
- ✅ **Production Ready**: All features needed for real-world use
- ✅ **HTTP Handlers**: Complete REST API for admin operations
- ✅ **Widgets**: Rich form widgets for all field types
- ✅ **Export**: CSV and JSON export functionality

## Quick Start

### 1. Register a Model

```go
import (
	"github.com/forgego/forge/pkg/admin"
	"github.com/forgego/forge/pkg/admin/http"
	httplib "github.com/forgego/forge/pkg/http"
)

// Register User model
userAdmin := admin.Register(
	&models.User{},
	models.User.Objects,
	&admin.Config[*models.User]{
		ListDisplay: []admin.FieldExpr[*models.User, any]{
			admin.StringField(
				"username",
				func(u *models.User) string { return u.Username },
				func(u *models.User, v string) { u.Username = v },
			),
			admin.StringField(
				"email",
				func(u *models.User) string { return u.Email },
				func(u *models.User, v string) { u.Email = v },
			),
		},
		SearchFields: []admin.FieldExpr[*models.User, any]{
			admin.StringField("username",
				func(u *models.User) string { return u.Username },
				func(u *models.User, v string) { u.Username = v },
			),
			admin.StringField("email",
				func(u *models.User) string { return u.Email },
				func(u *models.User, v string) { u.Email = v },
			),
		},
		ListFilter: []admin.Filter[*models.User]{
			admin.NewBooleanFilter(
				admin.BoolField("is_active",
					func(u *models.User) bool { return u.IsActive },
					func(u *models.User, v bool) { u.IsActive = v },
				),
			),
		},
		ListPerPage: 25,
	},
)

// Register for HTTP handlers
http.RegisterAdminForHTTP(userAdmin)
```

### 2. Register Routes

```go
router := httplib.NewRouter()

adminRouter := http.NewRouter(admin.GetGlobalRegistry())
adminRouter.RegisterRoutes(router, "/admin")

http.ListenAndServe(":8080", router)
```

### 3. Access Admin

Visit `http://localhost:8080/admin/` to see the admin dashboard.

## Field Expressions

### Creating Field Expressions

```go
// String field
usernameField := admin.StringField(
	"username",
	func(u *models.User) string { return u.Username },
	func(u *models.User, v string) { u.Username = v },
)

// Int64 field
idField := admin.Int64Field(
	"id",
	func(u *models.User) int64 { return u.ID },
	func(u *models.User, v int64) { u.ID = v },
)

// Bool field
isActiveField := admin.BoolField(
	"is_active",
	func(u *models.User) bool { return u.IsActive },
	func(u *models.User, v bool) { u.IsActive = v },
)

// Time field
createdAtField := admin.TimeField(
	"created_at",
	func(u *models.User) time.Time { return u.CreatedAt },
	func(u *models.User, v time.Time) { u.CreatedAt = v },
)
```

## Configuration Options

### ListDisplay

Fields to show in list view:

```go
ListDisplay: []admin.FieldExpr[*models.User, any]{
	usernameField,
	emailField,
	isActiveField,
}
```

### SearchFields

Fields to search in:

```go
SearchFields: []admin.FieldExpr[*models.User, any]{
	usernameField,
	emailField,
}
```

### ListFilter

Filters for list view:

```go
ListFilter: []admin.Filter[*models.User]{
	admin.NewBooleanFilter(isActiveField),
	admin.NewChoiceFilter(
		statusField,
		[]admin.Choice[string]{
			{Label: "Active", Value: "active"},
			{Label: "Inactive", Value: "inactive"},
		},
	),
}
```

### Ordering

Default ordering:

```go
Ordering: []admin.Ordering[*models.User]{
	admin.OrderBy(createdAtField).Desc(),
	admin.OrderBy(usernameField).Asc(),
}
```

### Actions

Bulk actions:

```go
Actions: []admin.Action[*models.User]{
	admin.NewAction(
		"activate",
		"Activate selected users",
		func(ctx context.Context, users []*models.User) error {
			for _, user := range users {
				user.IsActive = true
				if err := models.User.Objects.Update(ctx, user); err != nil {
					return err
				}
			}
			return nil
		},
	),
}
```

### Fieldsets

Group form fields:

```go
Fieldsets: []admin.Fieldset[*models.User]{
	admin.NewFieldset(
		"Account Information",
		usernameField,
		emailField,
	),
	admin.NewFieldset(
		"Permissions",
		isActiveField,
		isStaffField,
	).WithCollapsed(true),
}
```

## HTTP Routes

All routes are registered under `/admin/`:

- `GET /admin/` - Admin index/dashboard
- `GET /admin/{model}/` - List view
- `GET /admin/{model}/new/` - Create form
- `POST /admin/{model}/new/` - Create POST
- `GET /admin/{model}/{id}/` - Detail view
- `GET /admin/{model}/{id}/change/` - Update form
- `POST /admin/{model}/{id}/change/` - Update POST
- `DELETE /admin/{model}/{id}/delete/` - Delete
- `GET /admin/{model}/export/` - Export (CSV/JSON)
- `POST /admin/{model}/bulk-action/` - Bulk actions
- `GET /admin/{model}/autocomplete/` - Autocomplete

## Production Considerations

### Security
- Add authentication middleware
- Implement permission checking
- Add CSRF protection
- Validate all inputs

### Performance
- Use pagination (default: 20 per page)
- Optimize queries
- Add caching if needed
- Use database indexes

### Error Handling
- Use AdminError for admin-specific errors
- Log errors appropriately
- Return user-friendly messages
- Handle edge cases

## Architecture

### Core Components

- **Admin[T]**: Type-safe admin instance
- **Config[T]**: Type-safe configuration
- **FieldExpr[T, F]**: Type-safe field expressions
- **Filter[T]**: Type-safe filters
- **Action[T]**: Type-safe bulk actions
- **Inline[T, R]**: Type-safe inline editing
- **ListView[T]**: Type-safe list view
- **DetailView[T]**: Type-safe detail view
- **FormView[T]**: Type-safe form view

## Complete Example

See `pkg/admin/complete_example.go` for a full production example.
