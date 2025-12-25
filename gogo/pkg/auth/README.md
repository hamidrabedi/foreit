# Auth Module

Authentication and authorization system with policies and roles.

## Concepts

### Authenticators

Authenticators handle authentication:

```go
// JWT authentication
jwt := auth.NewJWT("secret-key")
app.Use(auth.Middleware(jwt))

// Session authentication
session := auth.NewSession("user")
app.Use(auth.Middleware(session))

// API Key authentication
apiKey := auth.NewAPIKey(func(key string) (interface{}, error) {
    // Validate API key and return user
    return getUserByAPIKey(key), nil
})
app.Use(auth.Middleware(apiKey))
```

### Policies

Policies define authorization rules:

```go
type PostPolicy struct {
    auth.Policy[models.Post]
}

func (p *PostPolicy) CanView(user interface{}, post *models.Post) bool {
    return post.Published || post.AuthorID == getUserID(user)
}

func (p *PostPolicy) CanEdit(user interface{}, post *models.Post) bool {
    return post.AuthorID == getUserID(user) || isAdmin(user)
}

// Register policy
auth.Register[models.Post](&PostPolicy{})

// Use in handlers
if err := auth.Require[models.Post](ctx, "edit", post); err != nil {
    return err
}
```

### Roles

Role-based access control:

```go
// Check roles
if auth.HasRole(user, auth.RoleAdmin) {
    // Admin only
}

// Require role middleware
app.Use(auth.RequireRole(auth.RoleAdmin))
app.Use(auth.RequireAnyRole(auth.RoleAdmin, auth.RoleModerator))
```

## Features

- Multiple authentication methods (JWT, Session, API Key)
- Policy-based authorization
- Role-based access control (RBAC)
- Type-safe policy checking
- Middleware integration

