# Cache Module

Caching layer with support for tags and TTL.

## Usage

### Basic Caching

```go
// Set with TTL
cache.Set(ctx, "user:123", user, 10*time.Minute)

// Get
user, err := cache.Get(ctx, "user:123")

// Check existence
exists, _ := cache.Has(ctx, "user:123")

// Delete
cache.Delete(ctx, "user:123")
```

### Remember Pattern

```go
user, err := cache.Remember(ctx, "user:123", 10*time.Minute, func() (interface{}, error) {
    return getUserFromDB(123)
})
```

### Tagged Caching

```go
// Set with tags
cache.TagSet(ctx, "post:456", post, 1*time.Hour, "posts", "user:123")

// Invalidate by tag
cache.TagInvalidate(ctx, "posts") // Invalidates all posts
```

## Features

- In-memory store
- TTL support
- Tag-based invalidation
- Remember pattern
- Automatic cleanup

