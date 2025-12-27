# Forge Framework Tests

Comprehensive test suite for the Forge migration and schema system.

## Structure

```
tests/
├── go.mod                          # Tests module manifest
├── testhelpers/                    # Shared test utilities
│   ├── docker_testhelper.go       # Postgres/SQLite container setup
│   ├── sql_assertions.go          # Schema assertions
│   ├── fs_helper.go               # Filesystem utilities
│   └── cli_helper.go              # CLI invocation
├── pkg_migrations/                 # Migration tests
│   ├── migration_unit_test.go      # Unit tests (fast)
│   └── migration_integration_test.go # Integration tests (needs DB)
├── cmd_forge/                      # CLI tests
│   └── cli_e2e_test.go            # End-to-end tests
└── testdata/                       # Sample models for testing
    └── sample_app_models.go
```

## Running Tests

### Run All Tests
```bash
cd tests
go test ./...
```

### Run Unit Tests Only (Fast)
```bash
go test ./pkg_migrations -short
```

### Run Integration Tests
```bash
# SQLite (no external dependencies)
go test -v -timeout 30m ./pkg_migrations

# Postgres (requires Postgres 15+ running)
export RUN_POSTGRES_TESTS=true
go test -v -timeout 30m -run TestMigrationApplyPostgres ./pkg_migrations
```

### Run E2E Tests
```bash
go test -v -timeout 60m ./cmd_forge
```

### Run Specific Test
```bash
go test -v -run TestMigrationApplySQLite ./pkg_migrations
```

## Test Organization

**Unit Tests** - Fast, isolated tests without database dependencies
- Location: `pkg_migrations/migration_unit_test.go`
- Run with: `go test -short`

**Integration Tests** - Tests against real SQLite/Postgres databases
- Location: `pkg_migrations/migration_integration_test.go`
- Coverage: FK cascades, indexes, constraints, schema introspection

**E2E Tests** - Full CLI workflow tests
- Location: `cmd_forge/cli_e2e_test.go`
- Coverage: makemigrations, apply, status, migrations

## Test Helpers

The `testhelpers` package provides utilities:

- **`docker_testhelper.go`**
  - `StartPostgresContainer()` - Launch ephemeral Postgres
  - `StartSQLiteMemory()` - Create in-memory SQLite DB
  - `WaitForDBReady()` - Poll for DB readiness

- **`sql_assertions.go`**
  - `AssertTableExists()`, `AssertColumnExists()`, `AssertIndexExists()`
  - `AssertForeignKeyExists()`, `AssertRowCount()`
  - `RunSQLExpectSuccess()`, `RunSQLExpectError()`

- **`fs_helper.go`**
  - `TempWorkdir()` - Create temporary working directory
  - `ReadFileString()`, `WriteFileString()` - File I/O

- **`cli_helper.go`**
  - `RunCLI()` - Execute forge CLI with timeout and env vars

## Prerequisites

### For SQLite Tests (included)
- No external dependencies needed

### For Postgres Tests
- Docker installed (for `dockertest`)
- Postgres 15+ running, or `RUN_POSTGRES_TESTS` env var unset to skip

### For CLI Tests
- Forge CLI built: `cd .. && go build ./cmd/forge`

## Example Test

```go
func TestMyFeature(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start in-memory SQLite
	db, err := testhelpers.StartSQLiteMemory("")
	require.NoError(t, err)
	defer db.Close()

	// Run SQL
	testhelpers.RunSQLExpectSuccess(ctx, t, db, 
		"CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(255))")

	// Assert schema
	testhelpers.AssertTableExists(ctx, t, db, "sqlite", "users")
	testhelpers.AssertColumnExists(ctx, t, db, "sqlite", "users", "name")
	
	// Test business logic
	testhelpers.AssertRowCount(ctx, t, db, "users", 0)
}
```

## CI/CD

The main module's CI/CD (`.github/workflows/`) runs tests from this module:
- Unit tests on Go 1.20, 1.21
- Integration tests with Postgres 15
- E2E tests against live CLI

## Notes

- Tests module has its own `go.mod` to isolate test dependencies from the main module
- Main module users are not affected by test dependencies
- All test files follow Go conventions: `_test.go` suffix
- Use `testhelpers` package for database operations to ensure consistency

## Contributing

When adding new tests:
1. Place unit tests in `pkg_migrations/migration_unit_test.go`
2. Place integration tests in the appropriate location under their package
3. Use `testhelpers` for common operations
4. Update imports if moving between packages
5. Run tests locally before committing: `go test ./...`
