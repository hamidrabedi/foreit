# forge Development Guide

## Development Setup

### Prerequisites

- Go 1.24+
- PostgreSQL 12+

### Setup

```bash
# Clone repository
git clone https://github.com/forgego/forge.git
cd forge

# Install dependencies
go mod download

# Run tests
go test ./...

# Build CLI
go build ./cli/cmd
```

## Project Structure

```
forge/
├── cmd/
│   └── forge/              # CLI tool
├── pkg/
│   ├── schema/             # Schema definitions
│   ├── generator/          # Code generation
│   ├── query/              # Query builder
│   ├── db/                 # Database layer
│   ├── http/               # HTTP router
│   ├── admin/              # Admin interface
│   ├── security/           # Security features
│   ├── validation/         # Validation
│   ├── auth/               # Authentication
│   ├── config/             # Configuration
│   ├── logging/            # Logging
│   ├── registry/           # Model registry
│   └── utils/              # Utilities
├── examples/
│   └── blog/               # Example application
├── docs/                   # Documentation
└── go.mod
```

## Development Workflow

### 1. Schema Changes

1. Modify schema definition
2. Run `forge generate`
3. Review generated code
4. Run tests
5. Update documentation if needed

### 2. Adding Features

1. Create feature branch
2. Implement feature
3. Write tests
4. Update documentation
5. Submit PR

### 3. Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/query
```

## Code Style

### Naming Conventions

- **Packages:** lowercase, single word
- **Types:** PascalCase
- **Functions:** PascalCase (exported), camelCase (internal)
- **Variables:** camelCase
- **Constants:** PascalCase or UPPER_SNAKE_CASE

### File Organization

- One type per file (for large types)
- Related types in same package
- Keep files under 500 lines when possible

### Error Handling

```go
// Good
if err != nil {
    return fmt.Errorf("context: %w", err)
}

// Avoid
if err != nil {
    return err  // Loses context
}
```

### Documentation

- Export all public APIs
- Use godoc comments
- Include examples in comments
- Document complex logic

## Testing Guidelines

### Unit Tests

```go
func TestFieldExpr(t *testing.T) {
    field := query.NewFieldExpr[string]("username", "users")

    q := field.Equals("john")
    assert.NotNil(t, q)

    sql, args := q.ToSQL()
    assert.Equal(t, "username = ?", sql)
    assert.Equal(t, []interface{}{"john"}, args)
}
```

### Integration Tests

```go
func TestUserCRUD(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    ctx := context.Background()

    // Create
    user := &User{Username: "test"}
    err := User.Objects.Create(ctx, user)
    require.NoError(t, err)

    // Read
    found, err := User.Objects.Get(ctx, user.ID)
    require.NoError(t, err)
    assert.Equal(t, "test", found.Username)

    // Update
    found.Username = "updated"
    err = User.Objects.Update(ctx, found)
    require.NoError(t, err)

    // Delete
    err = User.Objects.Delete(ctx, found)
    require.NoError(t, err)
}
```

## Contributing

### Pull Request Process

1. Fork the repository
2. Create feature branch
3. Make changes
4. Write/update tests
5. Update documentation
6. Submit PR with description

### Commit Messages

```
feat: Add UUID field type support
fix: Fix QuerySet.Count() implementation
docs: Update API reference
refactor: Simplify AST parser
test: Add tests for FieldExpr
```

### Code Review

- All PRs require review
- Tests must pass
- Documentation must be updated
- No breaking changes without discussion

## Architecture Decisions

### Why SQL Builder?

- Type-safe SQL generation from Go code
- Proper identifier escaping to prevent SQL injection
- Parameter binding for all values
- Works with standard database/sql
- SQL-first approach
- Industry standard

### Why chi?

- Lightweight, composable
- Standard library compatible
- Framework-friendly

### Why code generation?

- Type safety
- Performance
- IDE support

## Performance Guidelines

### Database

- Use indexes for frequently queried fields
- Use SelectRelated for JOINs
- Use PrefetchRelated for separate queries
- Batch operations when possible

### Code Generation

- Keep templates simple
- Avoid complex logic in templates
- Use helper functions

### Memory

- Reuse QuerySet instances when possible
- Close rows properly
- Use connection pooling

## Security Guidelines

### Input Validation

- Always validate user input
- Use parameterized queries
- Sanitize output

### Authentication

- Hash passwords (bcrypt)
- Use secure sessions
- Implement CSRF protection

### Error Handling

- Don't expose internal errors
- Log errors securely
- Return user-friendly messages

## Debugging

### Enable Debug Logging

```go
logger, _ := logging.NewLogger(true)  // development mode
```

### Query Logging

```go
// Log all queries
db.EnableQueryLogging(true)
```

### Profiling

```go
import _ "net/http/pprof"

// In main.go
go http.ListenAndServe(":6060", nil)
```

## Release Process

### Versioning

- Follow semantic versioning
- Major: Breaking changes
- Minor: New features
- Patch: Bug fixes

### Release Checklist

Use this checklist when preparing a new release:

- [ ] All tests passing
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Version bumped
- [ ] Tagged in git
- [ ] Release notes written
