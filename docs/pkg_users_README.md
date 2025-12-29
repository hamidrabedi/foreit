# User System Package

Complete user management and authentication system for forge framework.

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

### 2. Initialize User System

```go
import "github.com/forgego/forge/pkg/users"

userSystem, err := users.SetupUserSystem(database, nil)
```

### 3. Register Routes

```go
import forgehttp "github.com/forgego/forge/pkg/http"

router := forgehttp.NewRouter()
userSystem.RegisterRoutes(router)
```

### 4. Use Services

```go
// Register a user
user, err := userSystem.UserService.Register(ctx, &service.RegisterRequest{
    Username: "johndoe",
    Email: "john@example.com",
    Password: "SecurePassword123!",
})

// Authenticate
authUser, err := userSystem.BackendRegistry.Authenticate(ctx, map[string]string{
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
UserSystem
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
import "github.com/forgego/forge/pkg/users/config"

userConfig := config.DefaultUserConfig()
userConfig.PasswordPolicy.MinLength = 12
userConfig.LockoutConfig.MaxFailedAttempts = 5

userSystem, err := users.SetupUserSystem(database, userConfig)
```

## Database Migrations

Run migrations to create all required tables:

```bash
forge migrate up
```

Migrations are located in `pkg/users/migrations/`

## Examples

See `docs/USER_SYSTEM_QUICK_START.md` for more examples.
