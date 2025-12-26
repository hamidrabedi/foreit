---
sidebar_position: 6
---

# Security Guide

forge includes built-in security features to protect your application from common web vulnerabilities.

## CSRF Protection

Cross-Site Request Forgery (CSRF) protection is enabled by default for all state-changing operations.

### Using CSRF in Forms

Include the CSRF token in your forms:

```html
<form method="POST" action="/posts/create/">
    <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
    <!-- form fields -->
</form>
```

### CSRF in API Requests

For API requests, include the token in headers:

```javascript
fetch('/api/posts/', {
    method: 'POST',
    headers: {
        'X-CSRF-Token': getCSRFToken(),
        'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
});
```

### Exempting Views

Exempt specific views from CSRF protection:

```go
import "github.com/forgego/forge/pkg/security"

router.Post("/api/webhook", csrf.Exempt(webhookHandler))
```

## XSS Protection

forge automatically escapes output in templates to prevent XSS attacks.

### Safe HTML

If you need to output HTML, mark it as safe:

```go
// In your template
{{.Content | safeHTML}}
```

But be careful - only use this for trusted content!

### Content Security Policy

Set Content Security Policy headers:

```go
import "github.com/forgego/forge/pkg/security"

router.Use(security.ContentSecurityPolicy(
    "default-src 'self'; script-src 'self' 'unsafe-inline'",
))
```

## SQL Injection Prevention

forge uses parameterized queries to prevent SQL injection.

### Always Use QuerySet

Never use raw SQL with user input:

```go
// Bad: SQL injection risk
query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)

// Good: Parameterized query
users, err := User.Objects.
    Filter(User.Fields.Username.Equals(username)).
    All(ctx)
```

### Raw Queries

If you must use raw SQL, use parameterized queries:

```go
rows, err := db.Query(
    "SELECT * FROM users WHERE username = $1 AND email = $2",
    username, email,
)
```

## Authentication

### Password Hashing

Always hash passwords:

```go
import "github.com/forgego/forge/pkg/auth"

hashed, err := auth.HashPassword(password)
if err != nil {
    return err
}
```

### Password Verification

Verify passwords:

```go
err := auth.VerifyPassword(hashed, password)
if err != nil {
    // Invalid password
}
```

### Session Management

Use secure sessions:

```go
import "github.com/forgego/forge/pkg/security"

session := security.NewSession()
session.Set("user_id", userID)
session.Save(w)
```

## Input Validation

### Validate User Input

Always validate user input:

```go
import "github.com/forgego/forge/pkg/validation"

func CreateUser(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    
    if err := validation.ValidateUsername(username); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Proceed with creation
}
```

### Model Validation

Use model validation:

```go
user := &User{Username: username}
if err := user.Validate(); err != nil {
    // Handle validation error
}
```

## Rate Limiting

Protect your API from abuse:

```go
import "github.com/forgego/forge/pkg/security"

router.Use(security.RateLimit(100, time.Minute)) // 100 requests per minute
```

## HTTPS

Always use HTTPS in production:

```go
// Redirect HTTP to HTTPS
router.Use(security.RequireHTTPS())
```

## Security Headers

Add security headers:

```go
import "github.com/forgego/forge/pkg/security"

router.Use(security.SecurityHeaders(
    security.XSSProtection(),
    security.ContentTypeOptions(),
    security.FrameOptions("DENY"),
    security.ReferrerPolicy("strict-origin-when-cross-origin"),
))
```

## File Upload Security

### Validate File Types

```go
func handleUpload(w http.ResponseWriter, r *http.Request) {
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Validate file type
    allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
    if !contains(allowedTypes, header.Header.Get("Content-Type")) {
        http.Error(w, "Invalid file type", http.StatusBadRequest)
        return
    }
    
    // Validate file size
    if header.Size > 5*1024*1024 { // 5MB
        http.Error(w, "File too large", http.StatusBadRequest)
        return
    }
    
    // Process file
}
```

### Store Files Securely

- Store uploads outside web root
- Use random filenames
- Validate file content, not just extension
- Scan for malware (in production)

## Best Practices

1. **Always Validate Input** - Never trust user input
2. **Use Parameterized Queries** - Prevent SQL injection
3. **Hash Passwords** - Never store plain text passwords
4. **Use HTTPS** - Encrypt data in transit
5. **Set Security Headers** - Add defense in depth
6. **Rate Limit APIs** - Prevent abuse
7. **Keep Dependencies Updated** - Patch vulnerabilities
8. **Log Security Events** - Monitor for attacks

## Security Checklist

Before deploying to production:

- [ ] CSRF protection enabled
- [ ] XSS protection enabled
- [ ] SQL injection prevention (parameterized queries)
- [ ] Passwords hashed
- [ ] HTTPS enabled
- [ ] Security headers set
- [ ] Rate limiting configured
- [ ] Input validation in place
- [ ] File uploads secured
- [ ] Error messages don't leak information
- [ ] Dependencies up to date
- [ ] Logging configured

## Next Steps

- [Deployment Guide](deployment) - Secure deployment practices
- [API Reference](../reference/schema) - Security-related APIs

