# Test Structure

This directory contains all tests for the Forge framework, organized by test type and purpose.

## Directory Structure

```
tests/
├── integration/                # DB / filesystem integration tests
│   ├── migrate/                # Migration system tests
│   │   ├── flow_test.go
│   │   ├── reversibility_test.go
│   │   ├── edge_cases_test.go
│   │   ├── ecommerce_test.go
│   │   └── postgres_features_test.go
│   │
│   ├── orm/                    # ORM/Query integration tests
│   │   ├── query_integration_test.go
│   │   └── example_test.go
│   │
│   └── schema/                 # Schema builder integration tests
│       └── schema_integration_test.go
│
├── e2e/                        # Full system end-to-end tests
│   └── cli/                    # CLI command tests
│       ├── forge_init_test.go
│       ├── migrate_up_test.go
│       ├── migrate_rollback_test.go
│       └── help_command_test.go
│
├── testdata/                   # STATIC INPUTS ONLY
│   ├── models/                 # Test model definitions
│   │   ├── blog.go
│   │   ├── ecommerce.go
│   │   └── complex.go
│   │
│   └── migrations/             # Sample migration files
│       └── sample_app/
│
├── infra/                      # Heavy infrastructure helpers
│   ├── docker/                 # Docker container management
│   │   ├── postgres.go         # Postgres container lifecycle
│   │   └── lifecycle.go        # Container lifecycle utilities
│   │
│   ├── database/               # Database connection helpers
│   │   ├── local_postgres.go   # Local Postgres connection
│   │   └── testdb.go           # Test database setup
│   │
│   └── filesystem/             # Filesystem sandbox helpers
│       └── fs_helper.go        # Temp directories, file operations
│
├── helpers/                    # Pure test utilities (NO IO)
│   ├── assertions.go           # Complex assertion helpers
│   ├── sql_assertions.go       # SQL/database assertions
│   ├── migration_assertions.go  # Migration-specific assertions
│   ├── postgres_features.go   # Postgres-specific feature tests
│   ├── cli_helper.go          # CLI execution helpers
│   └── doc.go                 # Documentation
│
└── README.md                   # This file
```

## Test Categories

### 🔴 integration/

**Real DB, Real SQL, Temporary filesystem**

Tests that interact with actual databases and filesystems. These are slower but test real behavior.

- Uses helpers from `infra/` for setup
- Most migration tests belong here
- Tests actual SQL execution and database state

**Key packages:**
- `integration/migrate` - Migration flow, reversibility, edge cases
- `integration/orm` - Query builder, ORM operations
- `integration/schema` - Schema builder integration

### 🟡 e2e/

**Runs compiled binaries, executes forge CLI**

Full system tests that run the actual CLI commands. These are the slowest tests.

- Tests real workflows end-to-end
- Executes `forge` binary
- Tests CLI commands: init, migrate, rollback, help

### 📦 testdata/

**ONLY static data - NO helpers, NO logic, NO assertions**

Static test data files:
- Model definitions
- SQL migration files
- Example schemas

**Rules:**
- 🚫 No helper functions
- 🚫 No test logic
- 🚫 No assertions
- ✅ Only data/definitions

### 🛠 infra/

**Heavy, stateful infrastructure**

Infrastructure helpers that manage external resources:
- Docker containers
- Database connections
- File system sandboxes

**Packages:**
- `infra/docker` - PostgresOpts, StartPostgresContainer
- `infra/database` - Database connection setup
- `infra/filesystem` - TempDirInTests, file operations

### 🧰 helpers/

**Pure logic utilities**

Pure helper functions with no external dependencies:
- Assertions and comparisons
- Builders and matchers
- Test utilities

**Rules:**
- 🚫 No database connections
- 🚫 No filesystem operations
- 🚫 No Docker
- ✅ Only pure functions

## Naming Rules

Test file names describe **behavior**, not implementation:

✅ **Good:**
- `migration_reversibility_test.go` - Tests reversibility behavior
- `edge_cases_test.go` - Tests edge cases
- `flow_test.go` - Tests migration flow

❌ **Bad:**
- `migration_system_test.go` - Too vague
- `comprehensive_test.go` - Not descriptive

**Preferred patterns:**
- `*_flow_test.go` - Tests workflow/flow
- `*_edge_cases_test.go` - Tests edge cases
- `*_integration_test.go` - Integration tests
- `*_e2e_test.go` - End-to-end tests

## Running Tests

### Run all tests:
```bash
go test ./tests/...
```

### Run specific category:
```bash
# Integration tests only
go test ./tests/integration/...

# E2E tests only
go test ./tests/e2e/...

# Migration tests only
go test ./tests/integration/migrate/...
```

### Run with database:
```bash
# Set Postgres connection (optional, defaults to localhost)
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=123

go test ./tests/integration/migrate/...
```

## Import Guidelines

When writing tests, import from the appropriate packages:

```go
// For database setup
import "github.com/forgego/forge/tests/infra/docker"
opts := docker.PostgresOpts{...}
db, dsn, cleanup, err := docker.StartPostgresContainer(ctx, opts)

// For filesystem operations
import "github.com/forgego/forge/tests/infra/filesystem"
tempDir, cleanup := filesystem.TempDirInTests(t, "prefix_")

// For assertions
import "github.com/forgego/forge/tests/helpers"
helpers.AssertTableExists(ctx, t, db, "postgres", "users")
helpers.AssertMigrationState(ctx, t, database, migrationsDir, 1, false)
```

## Migration Notes

This structure was refactored from the old `pkg_*` structure. Key changes:

- `pkg_migrations/` → `integration/migrate/`
- `pkg_query/` → `integration/orm/`
- `pkg_schema/` → `integration/schema/`
- `cmd_forge/` → `e2e/cli/`
- `testhelpers/` → Split into `helpers/` and `infra/`

All imports need to be updated to use the new package paths.
