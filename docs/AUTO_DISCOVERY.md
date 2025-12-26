# Auto-Discovery System

Forge uses build-time discovery to automatically find and register components in your project.

## How It Works

The framework uses `go:generate` directives to scan your project at build time and generate a registry file. This avoids runtime overhead and provides better performance.

## Discovery Rules

The framework automatically discovers:

| File Pattern | Purpose | Registration |
|-------------|---------|--------------|
| `app/*/models.go` | ORM models | Via schema registry |
| `app/*/admin.go` | Admin config (per-app) | Via `init()` functions |
| `app/*/api.go` | API endpoints | Via `init()` functions |
| `domain/*/service.go` | Business logic | Manual wiring (optional) |
| `infra/*/client.go` | Infrastructure | Manual wiring |

## Using init() for Auto-Registration

Each app's `admin.go` and `api.go` files use `init()` functions to automatically register components:

### Admin Registration

```go
package users

import (
	"github.com/forgego/forge/pkg/admin"
)

func init() {
	// Auto-register models for admin interface
	admin.RegisterModel(&User{})
	admin.RegisterModel(&Profile{})
}
```

### API Registration

```go
package users

import (
	"github.com/forgego/forge/pkg/api"
	httplib "github.com/forgego/forge/pkg/http"
)

func init() {
	// Auto-register API routes
	RegisterUserAPI(router)
}

func RegisterUserAPI(router *httplib.Router) {
	// Register your API endpoints here
}
```

## Build-Time Discovery

The main.go file includes a `go:generate` directive:

```go
//go:generate forge discover

package main
```

When you run `go generate`, this scans your project and generates a `registry_generated.go` file that contains all discovered components.

## Third-Party Packages

Third-party packages can register themselves via `init()` functions:

```go
import _ "github.com/forgego/auth"  // Auto-registers via init()
```

This allows for a plugin-like ecosystem where packages can automatically integrate with your Forge application.

## Manual Registration

For components that can't be auto-discovered (like domain services), you can manually wire them in your `main.go`:

```go
func main() {
	// ... setup code ...
	
	// Manually wire domain services
	userService := domain.NewUserService()
	// Use service...
}
```

## Performance Considerations

- Discovery happens at build time, not runtime
- Results are cached to avoid repeated scanning
- Efficient file walking with early termination
- No impact on application startup time

## Custom Discovery

You can customize discovery by:

1. Overriding templates in `.forge/templates/`
2. Using custom `init()` functions
3. Manual registration in `main.go`

