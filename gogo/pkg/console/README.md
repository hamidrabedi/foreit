# Console Module

Admin interface for managing data (like Django admin).

## Usage

### Auto-Register with Defaults

```go
console.Register[models.User](nil)
```

### Register with Custom Console

```go
type UserConsole struct {
    console.ModelConsole[models.User]
}

func (c *UserConsole) ListDisplay() []string {
    return []string{"name", "email", "created_at"}
}

func (c *UserConsole) ListFilters() []console.Filter {
    return []console.Filter{
        console.Filter{
            Field: "status",
            Type:  console.FilterTypeChoice,
            Choices: []console.Choice{
                {Value: "active", Label: "Active"},
                {Value: "inactive", Label: "Inactive"},
            },
        },
    }
}

console.Register[models.User](&UserConsole{})
```

### Install Routes

```go
console.InstallRoutes(app, "/console")
```

## Features

- Auto-generated forms from Ent schemas
- List views with filtering and search
- Detail views
- Create/Edit forms
- Custom actions
- Export/import (planned)

