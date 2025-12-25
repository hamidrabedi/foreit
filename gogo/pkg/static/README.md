# Static Module

Static file serving for Fiber applications.

## Usage

### Basic Usage

```go
app.Use(static.New(static.Config{
    Root: "./public",
    Prefix: "/static",
    Compress: true,
}))
```

### With Cache Headers

```go
app.Use(static.New(static.Config{
    Root: "./public",
    Prefix: "/static",
    MaxAge: 31536000, // 1 year
    Immutable: true,
}))
```

### Serve Single File

```go
app.Get("/favicon.ico", static.ServeFile("./public/favicon.ico", true))
```

## Features

- Directory serving
- MIME type detection
- Cache headers
- ETag support
- Compression
- Security (path validation)

