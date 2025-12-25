# Sessions Module

Session management for Fiber applications.

## Usage

### Setup

```go
store := sessions.NewMemoryStore(sessions.DefaultConfig())
app.Use(sessions.Middleware(store))
```

### Use in Handlers

```go
func handler(c *fiber.Ctx) error {
    // Set session value
    sessions.Set(c, "user_id", user.ID)
    
    // Get session value
    userID, _ := sessions.Get(c, "user_id")
    
    // Delete session value
    sessions.Delete(c, "user_id")
    
    // Clear all session data
    sessions.Clear(c)
    
    return nil
}
```

## Features

- In-memory store
- Cookie-based sessions
- Automatic expiration
- Session regeneration
- Configurable lifetime

