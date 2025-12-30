# Forge Admin - Redesigned Architecture

## Overview

The Forge Admin system has been completely redesigned to be a first-class citizen that fully integrates with the existing schema, ORM, and filter systems. This new architecture provides:

- **Schema-First Design**: Auto-discovery of fields, relations, and metadata
- **Full ORM Integration**: All queries use `orm.QuerySet[T]` and `orm.Manager[T]`
- **Filter Integration**: All filtering uses `filter.FilterSet[T]`
- **Type Safety**: Full type safety with Go generics throughout
- **Modern UI**: Template-based rendering with reusable components

## Architecture

### Core Components

```
forge/admin/
├── core/              # Core admin types (Admin[T], Config[T], Registry, Site)
├── schema/            # Schema integration (discovery, field mapping, relation mapping)
├── orm/               # ORM integration (AdminManager, AdminQuerySet, FieldAccessor)
├── filter/            # Filter integration (AdminFilterSet, filter builder, filter UI)
├── fields/            # Field system (type-safe fields, field registry, field types)
├── widgets/           # Widget system (widget registry, widget selection)
├── views/             # View system (ListView, FormView, DetailView)
├── templates/         # Template engine and renderer
├── components/        # Reusable UI components (Table, Form, Pagination, FilterSidebar)
├── http/              # HTTP handlers and routing
├── advanced/          # Advanced features (Actions, History, Permissions)
├── security/          # Security features (CSRF, XSS prevention)
├── codegen/           # Code generation tools
├── testing/           # Testing helpers and fixtures
└── examples/          # Usage examples
```

## Quick Start

### 1. Define Your Schema

```go
type UserSchema struct{}

func (s *UserSchema) Fields() []schema.Field {
    return []schema.Field{
        {
            Name:        "ID",
            Type:        schema.TypeInt64,
            PrimaryKey:  true,
            AutoIncrement: true,
        },
        {
            Name:        "Name",
            Type:        schema.TypeString,
            Required:    true,
            MaxLength:   intPtr(255),
            VerboseName: "Name",
        },
        {
            Name:        "Email",
            Type:        schema.TypeEmail,
            Required:    true,
            VerboseName: "Email",
        },
    }
}

func (s *UserSchema) Relations() []schema.Relation {
    return []schema.Relation{}
}

func (s *UserSchema) Meta() schema.Meta {
    return schema.Meta{
        TableName:         "users",
        VerboseName:       "User",
        VerboseNamePlural: "Users",
    }
}

func (s *UserSchema) Hooks() *schema.ModelHooks {
    return nil
}
```

### 2. Register Admin

```go
import (
    admincore "github.com/forgego/forge/admin/core"
    "github.com/forgego/forge/orm"
    "github.com/forgego/forge/schema"
)

// Create schema instance
schemaInstance := &UserSchema{}

// Get ORM manager (from your ORM setup)
manager := orm.NewManager[User](db)

// Create admin configuration
config := &admincore.Config[User]{
    VerboseName:       "User",
    VerboseNamePlural: "Users",
    ListPerPage:       20,
    // Fields will be auto-discovered from schema
}

// Register admin
admin, err := admincore.Register(schemaInstance, manager, config)
if err != nil {
    log.Fatal(err)
}
```

### 3. Set Up HTTP Routes

```go
import (
    admincore "github.com/forgego/forge/admin/core"
    adminhttp "github.com/forgego/forge/admin/http"
    httplib "github.com/forgego/forge/server"
)

// Get the global registry
registry := admincore.GetGlobalRegistry()

// Create router with template directory
router := adminhttp.NewCoreRouter(registry, "./templates")

// Register routes
httpRouter := httplib.NewRouter()
router.RegisterRoutes(httpRouter, "/admin")
```

## Features

### Auto-Discovery

The admin system automatically discovers:
- **Fields**: From `schema.Schema.Fields()`
- **Relations**: From `schema.Schema.Relations()`
- **Metadata**: From `schema.Schema.Meta()`
- **Filters**: Auto-generated based on field types

### Type Safety

All operations are type-safe using Go generics:

```go
admin := &admincore.Admin[User]{}
manager := admin.Manager()  // Returns *AdminManager[User]
filterset := admin.FilterSet()  // Returns *AdminFilterSet[User]
```

### Customization

You can customize any aspect:

```go
config := &admincore.Config[User]{
    // List view
    ListDisplay: []interface{}{"Name", "Email"},
    SearchFields: []interface{}{"Name", "Email"},
    Ordering: []admincore.Ordering[User]{
        *admincore.OrderBy[User]("CreatedAt").Desc(),
    },
    
    // Form view
    Fieldsets: []admincore.Fieldset[User]{
        admincore.NewFieldset[User]("Basic Info", "Name", "Email"),
        admincore.NewFieldset[User]("Settings", "IsActive"),
    },
    
    // Actions
    Actions: []admincore.Action[User]{
        admincore.NewAction[User](
            "activate",
            "Activate selected",
            func(ctx context.Context, instances []*User) error {
                // Custom action logic
                return nil
            },
        ),
    },
}
```

## Integration Points

### Schema Integration

- Fields are discovered from `schema.Schema.Fields()`
- Field types are mapped to admin field types
- Field options (required, unique, choices) are respected

### ORM Integration

- All queries use `orm.QuerySet[T]`
- Field access uses `orm.FieldExpression[V]`
- Optimizations: `SelectRelated`, `PrefetchRelated`

### Filter Integration

- Filters are built from schema fields
- Uses `filter.FilterSet[T]` for all filtering
- Supports all filter types (boolean, choice, date, number, related)

## Advanced Features

### Actions

Bulk operations on selected records:

```go
config.Actions = []admincore.Action[User]{
    admincore.NewAction[User](
        "delete_selected",
        "Delete selected",
        func(ctx context.Context, instances []*User) error {
            for _, instance := range instances {
                admin.DeleteModel(ctx, instance)
            }
            return nil
        },
    ).WithPermissions("users.delete_user"),
}
```

### History

Change tracking:

```go
import "github.com/forgego/forge/admin/advanced"

historyManager := advanced.NewHistoryManager(admin, historyStore)
historyManager.LogChange(ctx, instance, advanced.ActionChange, user, changes)
```

### Permissions

Permission checking:

```go
config.PermissionChecker = customPermissionChecker
config.HasAddPermission = func(ctx context.Context, admin *admincore.Admin[User], user interface{}) bool {
    return checkPermission(user, "users.add_user")
}
```

## Template System

Templates are located in `forge/admin/templates/templates/`:

- `base.html`: Base layout with header, sidebar, footer
- `list.html`: List view template
- `form.html`: Create/update form template
- `detail.html`: Detail view template
- `index.html`: Admin dashboard template

## Components

Reusable UI components:

- **Table**: Data table with sorting, selection
- **Form**: Form with fieldsets, validation
- **Pagination**: Page navigation
- **FilterSidebar**: Filter sidebar with all filter types

## Security

- **CSRF Protection**: Built-in CSRF token management
- **XSS Prevention**: HTML sanitization utilities
- **Permission System**: Fine-grained access control

## Testing

Testing helpers:

```go
import "github.com/forgego/forge/admin/testing"

helper, err := testing.NewTestAdminHelper(schemaInstance, manager, config)
// Use helper for testing
```

## Code Generation

Generate admin code:

```go
import "github.com/forgego/forge/admin/codegen"

generator := codegen.NewAdminGenerator("./output")
generator.GenerateAdmin("User", "models", fields)
```

## Migration from Old System

The old admin system (`forge/admin/admin.go`) is still available for backward compatibility. To migrate:

1. Replace `admin.Register` with `admincore.Register`
2. Update imports to use `admincore` package
3. Update HTTP handlers to use `CoreHandler` and `CoreRouter`
4. Move templates to the new template directory structure

## Examples

See `forge/admin/examples/basic_usage.go` for complete examples.
