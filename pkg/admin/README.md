# Admin Package - Django-Style Admin Interface

This package provides a comprehensive Django-style admin interface for managing models in your application.

## Features

### Core Features (Implemented)
- ✅ **Model Registration** - Register models for admin management
- ✅ **List View** - Paginated list view with search and filtering
- ✅ **Detail View** - View individual model instances
- ✅ **Create/Update Forms** - Full CRUD operations
- ✅ **Delete** - Delete model instances
- ✅ **Permissions** - Role-based access control
- ✅ **Authentication** - Login/logout functionality
- ✅ **Bulk Actions** - Perform actions on multiple objects

### Advanced Features (Newly Added)

#### 1. **Date Hierarchy**
Group and filter by date fields (year, month, day):
```go
config := ModelAdminConfig{
    DateHierarchy: "created_at",
}
ApplyModelAdminConfig(model, config)
```

#### 2. **Ordering**
Custom default ordering for list views:
```go
config := ModelAdminConfig{
    Ordering: []string{"-created_at", "name"}, // Descending by created_at, then name
}
```

#### 3. **List Per Page**
Customize items per page:
```go
config := ModelAdminConfig{
    ListPerPage: 50, // Show 50 items per page
}
```

#### 4. **Save Buttons**
Enhanced save options:
- Save and add another
- Save and continue editing
- Save as new (for updates)

#### 5. **Fieldsets**
Group form fields into collapsible sections:
```go
config := ModelAdminConfig{
    Fieldsets: []Fieldset{
        {
            Name: "Basic Information",
            Fields: []string{"name", "email", "phone"},
        },
        {
            Name: "Advanced",
            Fields: []string{"settings", "preferences"},
            Collapsed: true,
        },
    },
}
```

#### 6. **Export Functionality**
Export data as CSV or JSON:
- Access via `/admin/{model}/export/?format=csv` or `format=json`
- Exports all instances with list_display fields

#### 7. **Custom Queryset**
Override the default queryset:
```go
config := ModelAdminConfig{
    GetQueryset: func(ctx context.Context, manager interface{}) (interface{}, error) {
        // Custom filtering logic
        return manager.Filter(...), nil
    },
}
```

#### 8. **Custom Save/Delete Hooks**
Intercept save and delete operations:
```go
config := ModelAdminConfig{
    SaveModel: func(ctx context.Context, instance interface{}, form map[string]interface{}, isNew bool) error {
        // Custom validation or processing
        return nil
    },
    DeleteModel: func(ctx context.Context, instance interface{}) error {
        // Custom delete logic
        return nil
    },
}
```

#### 9. **Verbose Names**
Human-readable names for models:
```go
config := ModelAdminConfig{
    VerboseName: "User Account",
    VerboseNamePlural: "User Accounts",
}
```

#### 10. **Autocomplete Fields** (Planned)
For foreign key relationships with search.

#### 11. **Raw ID Fields** (Planned)
Show foreign keys as raw ID inputs instead of dropdowns.

#### 12. **Custom Admin Actions** (Planned)
Define custom bulk actions beyond delete.

## Usage Example

```go
import (
    "github.com/forgego/forge/pkg/admin"
)

// Register a model with full configuration
userModel := &models.User{}
userManager := models.NewUserManager(db)

config := admin.ModelAdminConfig{
    ListDisplay: []interface{}{"id", "username", "email", "is_active", "date_joined"},
    ListFilter: []interface{}{"is_active", "is_staff"},
    SearchFields: []interface{}{"username", "email"},
    DateHierarchy: "date_joined",
    Ordering: []string{"-date_joined"},
    ListPerPage: 25,
    Fieldsets: []admin.Fieldset{
        {
            Name: "Account",
            Fields: []string{"username", "email", "password"},
        },
        {
            Name: "Permissions",
            Fields: []string{"is_active", "is_staff", "is_superuser"},
        },
    },
    VerboseName: "User",
    VerboseNamePlural: "Users",
    SaveAsContinue: true,
    SaveAndAddAnother: true,
}

admin.ApplyModelAdminConfig(userModel, config)
admin.RegisterModelWithManager(userModel, userManager)
```

## Configuration Options

### ModelAdminConfig Fields

| Field | Type | Description |
|-------|------|-------------|
| `ListDisplay` | `[]interface{}` | Fields to show in list view |
| `ListFilter` | `[]interface{}` | Fields to filter by |
| `SearchFields` | `[]interface{}` | Fields to search in |
| `DateHierarchy` | `string` | Date field for grouping |
| `Ordering` | `[]string` | Default ordering (use `-` for descending) |
| `ListPerPage` | `int` | Items per page (default: 20) |
| `Fieldsets` | `[]Fieldset` | Group form fields |
| `ReadOnlyFields` | `[]interface{}` | Fields that cannot be edited |
| `AutocompleteFields` | `[]string` | Foreign keys with autocomplete |
| `RawIDFields` | `[]string` | Foreign keys as raw IDs |
| `Actions` | `[]AdminAction` | Custom bulk actions |
| `VerboseName` | `string` | Human-readable name |
| `VerboseNamePlural` | `string` | Plural name |
| `SaveOnTop` | `bool` | Show save buttons on top |
| `SaveAs` | `bool` | Show "Save as new" button |
| `SaveAsContinue` | `bool` | Show "Save and continue editing" |
| `SaveAndAddAnother` | `bool` | Show "Save and add another" |
| `GetQueryset` | `func` | Custom queryset function |
| `SaveModel` | `func` | Custom save hook |
| `DeleteModel` | `func` | Custom delete hook |
| `ExportFormats` | `[]string` | Supported export formats |

## Production Considerations

### Security
- All views require authentication (if session manager provided)
- Permission checks on all operations
- CSRF protection recommended for forms

### Performance
- Pagination enabled by default
- Custom queryset support for optimized queries
- Export functionality for large datasets

### Extensibility
- Custom admin classes via `CustomAdmin` interface
- Custom save/delete hooks
- Custom queryset filtering
- Template customization support

## Roadmap

- [ ] Inline editing for related models
- [ ] Autocomplete for foreign keys
- [ ] History/audit logging
- [ ] Custom admin views
- [ ] Template customization per model
- [ ] Advanced filtering UI
- [ ] Date/time widgets
- [ ] Rich text editors
- [ ] Image/file uploads

