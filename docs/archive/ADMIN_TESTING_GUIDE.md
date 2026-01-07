# Forge Admin - Comprehensive Testing Guide

## Overview

This guide covers testing strategies for the Forge admin framework, including unit tests, integration tests, and end-to-end tests.

## Test Coverage

### Completed Test Suites

1. **History System** (`forge/admin/advanced/history_test.go`)
   - DefaultHistoryStore functionality
   - History entry creation and retrieval
   - User history tracking
   - Object history tracking
   - Benchmarks for performance

2. **Security Features** (`forge/server/security_test.go`)
   - XSS sanitization
   - HTML escaping
   - SQL injection prevention
   - Input validation
   - Query parameterization
   - Content Security Policy
   - Benchmarks for security operations

3. **Widgets** (`forge/admin/widgets/widgets_test.go`)
   - Rich text editor widget
   - File upload widget
   - Select search widget
   - Widget registry
   - Widget rendering
   - Widget parsing
   - XSS prevention in widgets
   - Benchmarks for widget operations

4. **Migration Recovery** (`forge/db/migrate/execute/recover_test.go`)
   - Dirty state detection
   - Migration integrity validation
   - Partial rollback
   - Checksum comparison
   - SQL statement splitting
   - Recovery steps generation

## Running Tests

### Run All Tests
```bash
go test ./forge/...
```

### Run Specific Package Tests
```bash
go test ./forge/admin/advanced/...
go test ./forge/server/...
go test ./forge/admin/widgets/...
go test ./forge/db/migrate/execute/...
```

### Run With Coverage
```bash
go test -cover ./forge/...
go test -coverprofile=coverage.out ./forge/...
go tool cover -html=coverage.out
```

### Run Benchmarks
```bash
go test -bench=. ./forge/admin/advanced/...
go test -bench=. ./forge/server/...
go test -bench=. ./forge/admin/widgets/...
```

## Test Structure

### Unit Test Example

```go
func TestFeature(t *testing.T) {
    t.Run("SubTest1", func(t *testing.T) {
        // Arrange
        input := setupTestData()
        
        // Act
        result := functionUnderTest(input)
        
        // Assert
        assert.Equal(t, expected, result)
    })
    
    t.Run("SubTest2", func(t *testing.T) {
        // Test another aspect
    })
}
```

### Integration Test Example

```go
func TestIntegration(t *testing.T) {
    // Setup database
    db := setupTestDB(t)
    defer db.Close()
    
    // Setup admin
    site := admin.NewSite("Test", "/admin/")
    admin := setupTestAdmin(site, db)
    
    // Test workflow
    instance := createInstance()
    err := admin.Save(ctx, instance)
    require.NoError(t, err)
    
    // Verify
    retrieved, err := admin.Get(ctx, instance.ID)
    require.NoError(t, err)
    assert.Equal(t, instance.Title, retrieved.Title)
}
```

## Testing Best Practices

### 1. Use Table-Driven Tests

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test", "TEST", false},
        {"invalid input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Function(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

### 2. Use Test Fixtures

```go
// fixtures.go
type TestFixtures struct {
    DB    *sql.DB
    Site  *admin.Site
    Posts []*BlogPost
}

func SetupFixtures(t *testing.T) *TestFixtures {
    db := setupTestDB(t)
    site := admin.NewSite("Test", "/admin/")
    
    posts := []*BlogPost{
        {ID: 1, Title: "Post 1"},
        {ID: 2, Title: "Post 2"},
    }
    
    return &TestFixtures{
        DB:    db,
        Site:  site,
        Posts: posts,
    }
}

func (f *TestFixtures) Teardown() {
    f.DB.Close()
}
```

### 3. Mock External Dependencies

```go
type MockHistoryStore struct {
    LogActionFunc      func(ctx context.Context, entry *HistoryEntry) error
    GetObjectHistoryFunc func(ctx context.Context, objType string, objID int64) ([]*HistoryEntry, error)
}

func (m *MockHistoryStore) LogAction(ctx context.Context, entry *HistoryEntry) error {
    if m.LogActionFunc != nil {
        return m.LogActionFunc(ctx, entry)
    }
    return nil
}
```

### 4. Test Error Paths

```go
func TestErrorHandling(t *testing.T) {
    t.Run("nil input", func(t *testing.T) {
        err := Function(nil)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "nil input")
    })
    
    t.Run("invalid state", func(t *testing.T) {
        err := Function(invalidInput)
        assert.Error(t, err)
    })
}
```

### 5. Use Subtests for Organization

```go
func TestAdmin(t *testing.T) {
    t.Run("Create", func(t *testing.T) {
        // Test creation
    })
    
    t.Run("Update", func(t *testing.T) {
        // Test updates
    })
    
    t.Run("Delete", func(t *testing.T) {
        // Test deletion
    })
}
```

## Performance Testing

### Benchmark Example

```go
func BenchmarkOperation(b *testing.B) {
    setup := prepareData()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Operation(setup)
    }
}
```

### Benchmark Results Interpretation

```
BenchmarkHistoryStore/LogAction-8         500000      3000 ns/op      500 B/op      10 allocs/op
BenchmarkHistoryStore/GetObjectHistory-8  100000     15000 ns/op     2000 B/op      50 allocs/op
```

- First number: iterations
- Second number: ns/op (nanoseconds per operation)
- Third number: B/op (bytes allocated per operation)
- Fourth number: allocs/op (allocations per operation)

## Testing Checklist

### Core Features
- [x] History logging
- [x] Security features (XSS, SQL injection)
- [x] Widget rendering
- [x] Widget parsing
- [x] Migration recovery
- [ ] Permission system
- [ ] Form validation
- [ ] HTTP handlers
- [ ] Template rendering

### Integration Tests
- [ ] Complete admin workflow (create, read, update, delete)
- [ ] Authentication and authorization
- [ ] File uploads
- [ ] Bulk actions
- [ ] Export functionality
- [ ] Search and filtering

### End-to-End Tests
- [ ] Browser automation tests
- [ ] Full user workflows
- [ ] Performance under load

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      - name: Run tests
        run: go test -v -cover ./...
      - name: Run benchmarks
        run: go test -bench=. -benchmem ./...
```

## Coverage Goals

| Component | Current Coverage | Goal |
|-----------|-----------------|------|
| History System | 80% | 90% |
| Security | 85% | 95% |
| Widgets | 75% | 85% |
| Migration Recovery | 70% | 85% |
| HTTP Handlers | 30% | 80% |
| Views | 20% | 75% |
| Overall | 50% | 80% |

## Testing Tools

### Required Dependencies

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/stretchr/testify/mock
go get github.com/mattn/go-sqlite3  # For test database
```

### Recommended Tools

- **testify** - Assertions and mocks
- **gomock** - Mock generation
- **go-sqlmock** - Database mocking
- **httptest** - HTTP testing
- **golangci-lint** - Linting

## Common Test Patterns

### 1. HTTP Handler Testing

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/admin/posts/", nil)
    rec := httptest.NewRecorder()
    
    handler.ServeHTTP(rec, req)
    
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Contains(t, rec.Body.String(), "Posts")
}
```

### 2. Database Testing

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)
    
    // Run migrations
    _, err = db.Exec(schema)
    require.NoError(t, err)
    
    return db
}
```

### 3. Context Testing

```go
func TestWithContext(t *testing.T) {
    ctx := context.WithValue(context.Background(), "user", mockUser)
    
    result, err := Function(ctx)
    
    assert.NoError(t, err)
}
```

## Debugging Tests

### Verbose Output
```bash
go test -v ./...
```

### Run Single Test
```bash
go test -run TestSpecificTest ./package
```

### Debug With Delve
```bash
dlv test ./package -- -test.run TestSpecificTest
```

## Next Steps

1. Add integration tests for HTTP handlers
2. Add E2E tests for complete workflows
3. Increase coverage for views and templates
4. Add performance regression tests
5. Add security vulnerability tests
6. Set up automated testing in CI/CD

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Coverage](https://blog.golang.org/cover)
