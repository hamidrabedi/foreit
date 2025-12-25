# Testing Module

Testing utilities for Gogo applications.

## Usage

### Test Application

```go
func TestMyFeature(t *testing.T) {
    app, client, err := testing.TestApp()
    if err != nil {
        t.Fatalf("Failed to create test app: %v", err)
    }
    defer client.Close()
    
    // Test your feature
}
```

### Test Client

```go
func TestDatabase(t *testing.T) {
    client, err := testing.TestClient()
    if err != nil {
        t.Fatalf("Failed to create test client: %v", err)
    }
    defer client.Close()
    
    // Test database operations
}
```

## Features

- Test application creation
- Test database client
- Test request helpers
- Database cleanup utilities

