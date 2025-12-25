# Pipeline Module

Request/response middleware pipeline for Fiber applications.

## Usage

### Basic Middleware

```go
app.Use(pipeline.Logging())
app.Use(pipeline.Recovery())
app.Use(pipeline.CORS())
app.Use(pipeline.RequestID())
```

### Rate Limiting

```go
app.Use(pipeline.RateLimit(100, time.Minute))
```

### Security Headers

```go
app.Use(pipeline.SecurityHeaders())
```

### Custom Middleware

```go
app.Use(func(c *fiber.Ctx) error {
    // Before handler
    start := time.Now()
    
    err := c.Next()
    
    // After handler
    duration := time.Since(start)
    log.Printf("Request took %v", duration)
    
    return err
})
```

### Chaining

```go
app.Use(pipeline.Chain(
    pipeline.Logging(),
    pipeline.Recovery(),
    pipeline.CORS(),
))
```

## Features

- Request logging
- Error recovery
- CORS handling
- Rate limiting
- Security headers
- Request ID generation
- Compression
- User context extraction

