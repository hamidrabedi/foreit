# Admin System Testing Strategy

## Overview

This document outlines a comprehensive testing strategy for the admin system using a **testing pyramid** approach with unit, integration, and E2E tests.

## Testing Pyramid

```
        /\
       /E2E\          (10% - Full HTTP workflows)
      /----\
     /Integration\    (20% - Component interactions)
    /------------\
   /  Unit Tests  \   (70% - Isolated components)
  /----------------\
```

## Test Structure

### 1. Unit Tests (70%)

**Location**: `pkg/admin/*_test.go`

**Purpose**: Test individual components in isolation

**Files**:
- `admin_test.go` - Admin registration, config, registry
- `fields_test.go` - Field expressions (get/set)
- `filters_test.go` - Filter creation and options
- `actions_test.go` - Action creation and execution
- `views/*_test.go` - View rendering logic

**Example**:
```go
func TestFieldExpr_GetSet(t *testing.T) {
    user := &TestUser{Username: "testuser"}
    
    field := NewFieldExpr(
        "username",
        func(u *TestUser) interface{} { return u.Username },
        func(u *TestUser, v interface{}) { u.Username = v.(string) },
        schema.Field{Name: "username", Type: schema.TypeString},
    )
    
    if field.Get(user) != "testuser" {
        t.Errorf("Expected 'testuser', got '%v'", field.Get(user))
    }
}
```

### 2. Integration Tests (20%)

**Location**: `pkg/admin/integration_test.go`

**Purpose**: Test component interactions

**Tests**:
- CRUD workflows
- Filter application
- Bulk actions
- Search functionality

**Example**:
```go
func TestAdmin_Integration_CRUD(t *testing.T) {
    admin := Register(&TestUser{}, manager, config)
    
    // Test Create
    user := &TestUser{Username: "testuser"}
    err := admin.SaveModel(ctx, user, formData, true)
    // ...
}
```

### 3. E2E Tests (10%)

**Location**: `tests/admin/e2e_test.go` (separate package to avoid cycles)

**Purpose**: Test complete HTTP workflows

**Tests**:
- Admin index
- List view with pagination
- Create/Update/Delete forms
- Search and filters
- Bulk actions
- Export
- Autocomplete

**Example**:
```go
func TestE2E_AdminWorkflow(t *testing.T) {
    // Setup router
    router := setupTestRouter()
    
    // Test list view
    req := httptest.NewRequest("GET", "/admin/TestUser/", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Test Helpers

### TestAdminClient

Located in `pkg/admin/testing/helpers.go`:

```go
client := testing.NewTestAdminClient(registry)
response := client.Get("/TestUser/")
if response.Status() != 200 {
    t.Errorf("Expected 200, got %d", response.Status())
}
```

### TestManager

Mock manager for testing:

```go
manager := testing.NewTestManager[TestUser]()
user := &TestUser{Username: "test"}
manager.Create(ctx, user)
```

## Running Tests

### All Tests
```bash
go test ./pkg/admin/... -v
```

### Unit Tests Only
```bash
go test -short ./pkg/admin/...
```

### Integration Tests
```bash
go test -run Integration ./pkg/admin/...
```

### E2E Tests (separate package)
```bash
go test ./tests/admin/... -v
```

### With Coverage
```bash
go test ./pkg/admin/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Models

Use `TestUser` as a simple test model:

```go
type TestUser struct {
    ID        int64
    Username  string
    Email     string
    IsActive  bool
    CreatedAt time.Time
}
```

## What's Tested

### ✅ Core Components
- [x] Admin registration
- [x] Registry operations
- [x] Field expressions (get/set)
- [x] Filters (boolean, choice)
- [x] Actions (bulk operations)
- [x] Fieldsets
- [x] Ordering

### ✅ Views
- [x] ListView rendering
- [x] DetailView rendering
- [x] FormView rendering
- [x] Pagination
- [x] Search
- [x] Filtering

### ✅ HTTP Handlers
- [x] Index handler
- [x] List handler
- [x] Detail handler
- [x] Create handler
- [x] Update handler
- [x] Delete handler
- [x] Export handler
- [x] Autocomplete handler
- [x] Bulk action handler

### ✅ E2E Workflows
- [x] Admin index
- [x] List view
- [x] List with search
- [x] List with pagination
- [x] Create form
- [x] Detail view
- [x] Update form
- [x] Export
- [x] Autocomplete
- [x] 404 handling

## Test Patterns

### Table-Driven Tests
```go
func TestFieldExpr_Get(t *testing.T) {
    tests := []struct {
        name     string
        user     *TestUser
        field    FieldExpr[*TestUser, interface{}]
        expected interface{}
    }{
        {
            name:     "username",
            user:     &TestUser{Username: "test"},
            field:    usernameField,
            expected: "test",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.field.Get(tt.user)
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## Notes

- Tests use in-memory structures where possible for speed
- Integration tests may require database setup
- E2E tests use httptest for HTTP testing
- All tests are designed to be independent and parallelizable
- Import cycles are avoided by placing E2E tests in separate package

## Status

✅ **Testing Strategy Complete**

All core testing components have been implemented:
- ✅ Unit tests for admin components
- ✅ Integration tests for CRUD workflows
- ✅ E2E tests for HTTP handlers
- ✅ Test helpers and utilities

For ongoing testing improvements, see [Development Guide](DEVELOPMENT.md).
