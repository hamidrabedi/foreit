# Identity Package

Complete identity and access management (IAM) system for forge framework.

## Features

- ✅ User management (CRUD operations)
- ✅ Authentication (password, token-based)
- ✅ Session management
- ✅ Password management (reset, change, validation)
- ✅ Permission system (RBAC)
- ✅ Group management
- ✅ Email verification
- ✅ Account lockout
- ✅ Rate limiting configuration
- ✅ Security features

## Quick Start

### 1. Setup Database

```go
import "github.com/forgego/forge/pkg/db"

database, err := db.NewDBFromConfig(cfg)
```

### 2. Initialize Identity System

```go
import "github.com/forgego/forge/pkg/identity"

identitySystem, err := identity.SetupIdentitySystem(database, nil)
```

### 3. Register Routes

```go
import forgehttp "github.com/forgego/forge/pkg/http"

router := forgehttp.NewRouter()
identitySystem.RegisterRoutes(router)
```

### 4. Use Services

```go
// Register a user
user, err := identitySystem.UserService.Register(ctx, &service.RegisterRequest{
    Username: "johndoe",
    Email: "john@example.com",
    Password: "SecurePassword123!",
})

// Authenticate
authUser, err := identitySystem.BackendRegistry.Authenticate(ctx, map[string]string{
    "username": "johndoe",
    "password": "SecurePassword123!",
})
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Users
- `GET /api/users/` - List users (authenticated)
- `POST /api/users/` - Create user (superuser)
- `GET /api/users/{id}` - Get user (authenticated)
- `PUT /api/users/{id}` - Update user (authenticated)
- `DELETE /api/users/{id}` - Delete user (superuser)

## Architecture

```
IdentitySystem
├── Repositories (Data Access)
│   ├── UserRepository
│   ├── SessionRepository
│   ├── TokenRepository
│   ├── PermissionRepository
│   └── GroupRepository
├── Services (Business Logic)
│   ├── UserService
│   ├── AuthService
│   ├── PasswordService
│   └── PermissionService
├── Backends (Authentication)
│   ├── PasswordBackend
│   ├── TokenBackend
│   └── BackendRegistry
├── Handlers (HTTP)
│   ├── UserHandler
│   └── AuthHandler
├── Middleware
│   └── AuthenticationMiddleware
└── Serializers
    ├── UserSerializer
    └── AuthSerializer
```

## Configuration

```go
import "github.com/forgego/forge/pkg/identity/config"

identityConfig := config.DefaultIdentityConfig()
identityConfig.PasswordPolicy.MinLength = 12
identityConfig.LockoutConfig.MaxFailedAttempts = 5

identitySystem, err := identity.SetupIdentitySystem(database, identityConfig)
```

## Database Migrations

Run migrations to create all required tables:

```bash
forge migrate up
```

Migrations are located in `pkg/identity/migrations/`

## Examples

See `docs/USER_SYSTEM_QUICK_START.md` for more examples.
