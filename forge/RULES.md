# Forge Framework Development Rules

This document outlines the mandatory rules and standards that ALL agents must follow when developing, testing, or contributing to the Forge framework. These rules ensure consistency, maintainability, and alignment with Forge's core principles of type safety, Django-like patterns, and extensibility.

## Table of Contents

- [Core Principles](#core-principles)
- [Coding Standards](#coding-standards)
- [Testing Standards](#testing-standards)
- [Schema and Model Rules](#schema-and-model-rules)
- [ORM and QuerySet Rules](#orm-and-queryset-rules)
- [API Development Rules](#api-development-rules)
- [Migration Rules](#migration-rules)
- [Error Handling](#error-handling)
- [Documentation Standards](#documentation-standards)

## Core Principles

1. **Type Safety First**: Forge prioritizes compile-time type checking over runtime flexibility
2. **Django Patterns**: Follow Django's design patterns where they make sense in Go
3. **Interface-Based Design**: Use interfaces for abstraction and testability
4. **Extensibility**: Design components to be pluggable and overridable
5. **Code Generation**: Leverage AST-based code generation for type safety
6. **Test-Driven**: Write comprehensive tests before/after implementation

## Coding Standards

### General Go Conventions

1. **Naming**:
   - Exported identifiers: `PascalCase`
   - Unexported identifiers: `camelCase`
   - Acronyms: `UUID`, `API`, `HTTP`, `JSON`, `SQL` (all caps)
   - Methods: Verb-based names (`GetUser()`, `CreateRecord()`)

2. **Imports**:
   - Standard library first
   - Third-party packages second
   - Internal packages last
   - Group with blank lines
   - Use aliases for conflicting names

```go
import (
    "context"
    "fmt"

    "github.com/stretchr/testify/assert"

    "github.com/forgego/forge/schema"
)
```

3. **Error Handling**:
   - Always handle errors explicitly
   - Wrap errors with context using `fmt.Errorf`
   - Use `errors.Is()` and `errors.As()` for error checking
   - Prefer early returns over nested if statements

4. **Panic Usage**:
   - Only panic for programmer errors or unrecoverable states
   - Never panic in user-facing code paths
   - Use panic in development/testing, but recover in production

### Framework-Specific Patterns

1. **Interface Implementation**:
   - Implement interfaces explicitly (no implicit satisfaction)
   - Use embedding for default behavior (`BaseSchema`, `BaseQuerySet`)
   - Define small, focused interfaces

2. **Generics Usage**:
   - Use generics for type-safe collections (`QuerySet[T]`)
   - Avoid over-genericizing - prefer concrete types when possible
   - Use type constraints appropriately (`any`, `comparable`)

3. **Context Usage**:
   - Always accept `context.Context` in public APIs
   - Use context for cancellation and timeouts
   - Pass context through call chains

4. **Builder Pattern**:
   - Use builders for complex object construction
   - Return builder for chaining (`field.WithRequired().WithMaxLength(100).Build()`)
   - Validate in `Build()` method

## Testing Standards

### Test Organization

1. **Test File Naming**:
   - Unit tests: `*_test.go` in same package
   - Integration tests: `tests/` directory with descriptive names
   - E2E tests: `tests/e2e/` directory

2. **Test Naming**:
   - `TestFunctionName` for unit tests
   - `TestFeature_Scenario` for integration tests
   - Descriptive names that explain what is being tested

### Testing Frameworks

1. **Assertion Library**: Use `testify/assert` exclusively
2. **Mocking**: Use `testify/mock` for interfaces
3. **HTTP Testing**: Use `net/http/httptest` for API endpoints

### Test Categories

1. **Unit Tests** (`*_test.go`):
   - Test pure functions and methods
   - Mock external dependencies
   - Fast execution (< 100ms per test)
   - No I/O operations

2. **Integration Tests** (`tests/integration/`):
   - Real database connections (PostgreSQL)
   - Full component interaction
   - Test actual SQL execution
   - Use timestamps in database names for isolation

3. **E2E Tests** (`tests/e2e/`):
   - Full CLI command execution
   - Real filesystem operations
   - Test complete workflows

### Testing Patterns

1. **Setup/Teardown**:
   - Use `t.Cleanup()` for resource cleanup
   - Setup test databases with unique names
   - Use test helpers for common setup

```go
func TestMyFeature(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // Setup database
    opts := testhelpers.PostgresOpts{...}
    db, dsn, cleanup, err := testhelpers.StartPostgresContainer(ctx, opts)
    require.NoError(t, err)
    defer cleanup()
}
```

2. **Table-Driven Tests**:
   - Use for testing multiple inputs/outputs
   - Clear test case structs with descriptive names

3. **Test Helpers**:
   - Pure assertion helpers in `tests/helpers/`
   - Infrastructure setup in `tests/infra/`
   - Reuse across test packages

4. **Coverage Goals**:
   - Unit tests: 80%+ coverage
   - Integration tests: All critical paths
   - E2E tests: Main workflows

## Schema and Model Rules

1. **Schema Interface**:
   - Always embed `schema.BaseSchema`
   - Implement `Fields()`, `Relations()`, `Meta()`, `Hooks()` methods
   - Return concrete field builders from `Fields()`

2. **Field Definitions**:
   - Use fluent builders: `schema.String("name").WithRequired().WithMaxLength(100).Build()`
   - Set appropriate validations and constraints
   - Use `VerboseName()` and `HelpText()` for user-facing fields

3. **Model Structure**:
   - Define models as structs with JSON tags
   - Use pointer receivers for methods where appropriate
   - Implement `schema.Schema` interface

```go
type User struct {
    ID    int64  `json:"id"`
    Email string `json:"email"`
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").WithPrimary().WithAutoIncrement().Build(),
        schema.String("email").WithRequired().WithUnique().Build(),
    }
}
```

4. **Relationships**:
   - Use `ForeignKey()`, `OneToOne()`, `ManyToMany()`
   - Set `OnDelete()` behavior explicitly
   - Define reverse relations with `RelatedName()`

## ORM and QuerySet Rules

1. **QuerySet Usage**:
   - Always chain methods fluently
   - Use type-safe field expressions when available
   - Set database connection with `SetDB()`

```go
users, err := qs.Filter(User.Email.Eq("test@example.com")).
    OrderBy(User.CreatedAt.Desc()).
    Limit(10).
    All(ctx)
```

2. **Field Expressions**:
   - Use generated type-safe expressions: `User.Email.Eq(value)`
   - Prefer expressions over raw strings
   - Combine with `And()`, `Or()`, `Not()`

3. **Error Handling**:
   - Check errors from all QuerySet operations
   - Use `Get()` for single results, `First()` for ordering-dependent
   - Handle SQL-specific errors appropriately

4. **Performance**:
   - Use `SelectRelated()` and `PrefetchRelated()` to avoid N+1 queries
   - Use `Only()` and `Defer()` for field selection
   - Use `Values()` and `ValuesList()` for bulk operations

## API Development Rules

1. **ViewSet Pattern**:
   - Extend `BaseViewSet` or `EnhancedBaseViewSet`
   - Implement required methods: `GetQueryset()`, `GetSerializer()`
   - Set `PermissionClasses` appropriately

2. **Serializer Pattern**:
   - Implement `Serializer` interface
   - Handle validation in `Validate()` method
   - Use field-level and object-level validation

3. **Router Configuration**:
   - Use `NewRouter()` or `NewEnhancedRouter()`
   - Register viewsets with descriptive names
   - Configure custom actions with `RegisterAction()`

4. **Permissions**:
   - Use built-in permissions: `IsAuthenticated()`, `IsAdmin()`, `AllowAny()`
   - Create custom permissions extending `BasePermission`
   - Check permissions in viewset methods

## Migration Rules

1. **Migration Generation**:
   - Run `forge makemigrations` after model changes
   - Review generated SQL before applying
   - Test migrations on development data first

2. **Migration Execution**:
   - Use `forge migrate` for applying migrations
   - Use `forge rollback` for reverting (when safe)
   - Never modify existing migrations manually

3. **Migration Testing**:
   - Test all change types: create, drop, modify, rename
   - Test reversibility where possible
   - Test edge cases (nullable columns, foreign keys, indexes)

4. **Version Control**:
   - Commit migration files with model changes
   - Never commit migration files without corresponding model changes
   - Use descriptive migration names

## Error Handling

1. **Error Types**:
   - Use custom error types for framework-specific errors
   - Implement `error` interface
   - Use error wrapping for context

2. **HTTP Errors**:
   - Use appropriate HTTP status codes
   - Return structured error responses
   - Include error details for debugging

3. **Validation Errors**:
   - Collect multiple validation errors
   - Return field-specific error messages
   - Use consistent error format

## Documentation Standards

1. **Code Comments**:
   - Comment exported functions, types, and methods
   - Explain complex logic and edge cases
   - Use complete sentences

2. **Package Documentation**:
   - Every package must have a `doc.go` file
   - Explain package purpose and usage
   - Include examples where helpful

3. **Example Code**:
   - Provide runnable examples in `*_example_test.go` files
   - Cover common use cases
   - Test examples with `go test -run=Example`

4. **API Documentation**:
   - Document REST API endpoints
   - Include request/response formats
   - Document authentication and permissions

---

## Enforcement

These rules are MANDATORY for all Forge development. Violations will be caught by:

- **Code Review**: All changes require review
- **CI/CD Pipeline**: Automated checks for style and coverage
- **Testing**: Comprehensive test suites ensure compliance

## Contributing

When adding new features:

1. Follow these rules from the start
2. Write tests first (TDD approach)
3. Update this document if rules need clarification
4. Ensure all tests pass: `go test ./...`
5. Check coverage meets minimum thresholds

Remember: Consistency across the codebase is crucial for maintainability and user experience.</content>
</xai:function_call António</xai:function_call name="update_todo_list">
<parameter name="todos">[x] Analyze forge framework structure and patterns\n[x] Review testing practices and guidelines\n[x] Review coding style and conventions\n[x] Create comprehensive rules document\n[x] Place rules file in forge directory
