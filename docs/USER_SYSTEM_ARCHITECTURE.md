# User System Architecture & Design Patterns

## Overview

This document describes the refactored architecture and design patterns for the new user system implementation.

## Directory Structure

```
pkg/users/
├── models/              # User models and related entities
│   ├── user.go         # User model
│   ├── session.go      # UserSession model
│   ├── permission.go   # Permission model
│   ├── group.go        # Group model
│   └── token.go        # Various token models
├── repository/         # Data access layer (Repository pattern)
│   ├── interface.go    # Repository interfaces
│   ├── user.go         # User repository implementation
│   ├── session.go      # Session repository
│   └── permission.go   # Permission repository
├── service/            # Business logic layer (Service pattern)
│   ├── interface.go    # Service interfaces
│   ├── user.go         # User service
│   ├── auth.go         # Authentication service
│   ├── password.go     # Password service
│   └── permission.go   # Permission service
├── backends/           # Authentication backends (Strategy pattern)
│   ├── interface.go    # Backend interface
│   ├── password.go     # Password backend
│   ├── token.go        # Token backend
│   ├── oauth.go        # OAuth backend
│   └── registry.go     # Backend registry
├── serializers/        # API serializers
│   ├── user.go         # User serializers
│   ├── auth.go         # Auth serializers
│   └── permission.go   # Permission serializers
├── handlers/           # HTTP handlers/viewsets
│   ├── user.go         # User viewset
│   ├── auth.go         # Auth endpoints
│   └── permission.go   # Permission endpoints
├── middleware/         # Authentication middleware
│   ├── auth.go         # Auth middleware
│   └── permission.go   # Permission middleware
├── validators/         # Validation logic
│   ├── password.go     # Password validators
│   └── user.go         # User validators
└── config/            # Configuration
    └── settings.go     # User system settings
```

## Design Patterns

### 1. Repository Pattern

**Purpose:** Abstract data access layer, making it easy to swap implementations and test business logic.

**Interface:**
```go
type UserRepository interface {
    // CRUD operations
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetByUsername(ctx context.Context, username string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
    
    // Query operations
    List(ctx context.Context, filters *UserFilters) ([]*User, error)
    Count(ctx context.Context, filters *UserFilters) (int64, error)
    Exists(ctx context.Context, email string) (bool, error)
}
```

**Benefits:**
- Separation of concerns
- Easy to mock for testing
- Can swap database implementations
- Centralized data access logic

### 2. Service Pattern

**Purpose:** Encapsulate business logic, orchestrate multiple repositories, and provide a clean API.

**Interface:**
```go
type UserService interface {
    // User management
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
    DeleteUser(ctx context.Context, id int64) error
    GetUser(ctx context.Context, id int64) (*User, error)
    
    // Authentication
    Register(ctx context.Context, req *RegisterRequest) (*User, error)
    VerifyEmail(ctx context.Context, token string) error
    
    // Password management
    ChangePassword(ctx context.Context, userID int64, req *ChangePasswordRequest) error
    RequestPasswordReset(ctx context.Context, email string) error
    ResetPassword(ctx context.Context, token string, newPassword string) error
}
```

**Benefits:**
- Business logic in one place
- Can orchestrate multiple repositories
- Transaction management
- Validation and error handling

### 3. Strategy Pattern (Authentication Backends)

**Purpose:** Support multiple authentication methods (password, token, OAuth) with a unified interface.

**Interface:**
```go
type AuthenticationBackend interface {
    // Authenticate attempts to authenticate using this backend
    Authenticate(ctx context.Context, credentials map[string]string) (*User, error)
    
    // GetUser retrieves a user by identifier (for token-based auth)
    GetUser(ctx context.Context, identifier string) (*User, error)
    
    // Supports returns true if this backend can handle the given credential type
    Supports(credentialType string) bool
    
    // Name returns the backend name
    Name() string
}
```

**Benefits:**
- Easy to add new authentication methods
- Backends are independent
- Can chain multiple backends
- Testable in isolation

### 4. Factory Pattern

**Purpose:** Create instances of services, repositories, and backends with proper dependencies.

```go
type UserSystemFactory interface {
    NewUserRepository() UserRepository
    NewUserService() UserService
    NewAuthService() AuthService
    NewBackendRegistry() BackendRegistry
}
```

### 5. Builder Pattern

**Purpose:** Construct complex objects (like User) step by step.

```go
type UserBuilder struct {
    user *User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{user: &User{}}
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
    b.user.Email = email
    return b
}

func (b *UserBuilder) WithUsername(username string) *UserBuilder {
    b.user.Username = username
    return b
}

func (b *UserBuilder) Build() (*User, error) {
    // Validate and return
    return b.user, nil
}
```

## Layer Responsibilities

### Models Layer
- **Responsibility:** Data structures, validation rules, relationships
- **No business logic** - just data representation
- **No database access** - pure structs

### Repository Layer
- **Responsibility:** Database operations, query building
- **No business logic** - just data access
- **Returns models** - not DTOs

### Service Layer
- **Responsibility:** Business logic, validation, orchestration
- **Uses repositories** - not direct database access
- **Transaction management** - coordinates multiple operations
- **Error handling** - converts errors to domain errors

### Handler Layer
- **Responsibility:** HTTP request/response handling
- **Uses services** - not repositories directly
- **Serialization** - converts models to JSON
- **Error responses** - formats errors for API

### Backend Layer
- **Responsibility:** Authentication strategy implementation
- **Independent** - each backend is self-contained
- **Stateless** - can be used across requests

## Data Flow

```
HTTP Request
    ↓
Middleware (Auth, Permission, Rate Limit)
    ↓
Handler (Parse request, validate format)
    ↓
Serializer (Validate input, deserialize)
    ↓
Service (Business logic, validation)
    ↓
Repository (Database operations)
    ↓
Database
    ↓
Repository (Return model)
    ↓
Service (Transform, business logic)
    ↓
Serializer (Serialize to JSON)
    ↓
Handler (Format response)
    ↓
HTTP Response
```

## Dependency Injection

All dependencies are injected through constructors:

```go
// Repository depends on database
type userRepository struct {
    db *db.DB
}

func NewUserRepository(db *db.DB) UserRepository {
    return &userRepository{db: db}
}

// Service depends on repositories
type userService struct {
    userRepo UserRepository
    emailService EmailService
}

func NewUserService(userRepo UserRepository, emailService EmailService) UserService {
    return &userService{
        userRepo: userRepo,
        emailService: emailService,
    }
}

// Handler depends on services
type userHandler struct {
    userService UserService
    serializer UserSerializer
}

func NewUserHandler(userService UserService, serializer UserSerializer) *userHandler {
    return &userHandler{
        userService: userService,
        serializer: serializer,
    }
}
```

## Error Handling

### Domain Errors
```go
type UserError struct {
    Code    string
    Message string
    Field   string // For validation errors
}

// Predefined errors
var (
    ErrUserNotFound = &UserError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrEmailExists = &UserError{Code: "EMAIL_EXISTS", Message: "Email already exists"}
    ErrInvalidCredentials = &UserError{Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
)
```

### Error Conversion
Services convert domain errors to HTTP status codes:
- `ErrUserNotFound` → 404
- `ErrEmailExists` → 409
- `ErrInvalidCredentials` → 401
- Validation errors → 400

## Testing Strategy

### Unit Tests
- **Repositories:** Mock database
- **Services:** Mock repositories
- **Handlers:** Mock services
- **Backends:** Test in isolation

### Integration Tests
- Test full flow: Handler → Service → Repository → Database
- Use test database
- Clean up after tests

### Test Doubles
```go
// Mock repository
type mockUserRepository struct {
    users map[int64]*User
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, ErrUserNotFound
    }
    return user, nil
}
```

## Configuration

Centralized configuration for user system:

```go
type UserConfig struct {
    // Password policy
    PasswordPolicy PasswordPolicy
    
    // Account lockout
    LockoutConfig LockoutConfig
    
    // Rate limiting
    RateLimitConfig RateLimitConfig
    
    // Email settings
    EmailVerificationRequired bool
    EmailFromAddress string
    
    // Session settings
    SessionTimeout time.Duration
    RememberMeDuration time.Duration
}
```

## Security Considerations

1. **Password Hashing:** Always use bcrypt/argon2
2. **Token Security:** Hash tokens before storing
3. **SQL Injection:** Use parameterized queries (handled by repository)
4. **XSS:** Sanitize user input (handled by serializer)
5. **CSRF:** Use CSRF tokens (handled by middleware)
6. **Rate Limiting:** Limit auth attempts (handled by middleware)

## Current Status

✅ **Complete and Production Ready**

All components of the user system have been implemented:
- ✅ User models (User, Session, Permission, Group, Token)
- ✅ Repositories (User, Session, Permission, Group)
- ✅ Services (User, Auth, Password, Permission)
- ✅ Authentication backends (Password, Token)
- ✅ Serializers for API
- ✅ HTTP handlers/viewsets
- ✅ Authentication middleware
- ✅ Password validators

See [User System Package](pkg_users_README.md) for usage examples.

## Benefits of This Architecture

1. **Separation of Concerns:** Each layer has a single responsibility
2. **Testability:** Easy to mock dependencies
3. **Maintainability:** Changes are localized
4. **Extensibility:** Easy to add new features
5. **Reusability:** Services can be used by different handlers
6. **Type Safety:** Go's type system ensures correctness

## Usage

The user system is ready for production use. See [User System Package Documentation](pkg_users_README.md) for complete usage guide and examples.
