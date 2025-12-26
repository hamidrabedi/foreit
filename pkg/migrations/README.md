# ForgeGo Migration System

A production-ready migration system for ForgeGo with comprehensive features for managing database schema changes.

## Features

### Core Features
- ✅ **Automatic Migration Generation** - Generate migrations from model definitions
- ✅ **Change Detection** - Automatically detect schema changes (tables, columns, indexes, constraints, foreign keys)
- ✅ **Multi-Database Support** - PostgreSQL and SQLite support with driver-specific SQL generation
- ✅ **State Management** - Track schema state and load from migration files
- ✅ **Reversible Migrations** - Generate both up and down migrations
- ✅ **Version Control** - Sequential versioning with checksum validation

### Production Features
- ✅ **Migration Validation** - Validate migrations before execution
- ✅ **Migration Linting** - Check migrations for common issues and best practices
- ✅ **Dry-Run Mode** - Preview migrations without applying them
- ✅ **Detailed Status** - View applied, pending, and out-of-order migrations
- ✅ **Error Recovery** - Tools and guidance for recovering from failed migrations
- ✅ **Checksum Validation** - Detect if migration files have been modified
- ✅ **Migration Squashing** - Combine multiple migrations into one

### CLI Commands

```bash
# Apply migrations
forge migrate up [--dry-run] [--path ./migrations]

# Show migration status
forge migrate status [--path ./migrations]

# Show specific migration or plan
forge migrate show [version] [--path ./migrations] [--models ./models]

# Rollback last migration
forge migrate rollback [--path ./migrations]

# Lint migrations
forge migrate lint [--path ./migrations]
```

## Architecture

### Package Structure

```
migrations/
├── core/           # Core types, interfaces, errors
├── detection/      # Change detection logic
├── generation/     # Migration generation
├── sql/            # SQL generation (PostgreSQL, SQLite)
├── state/          # Schema state management
├── execution/       # Migration execution and status
├── validation/     # Migration validation
├── linting/        # Migration linting
├── recovery/       # Error recovery
└── squashing/      # Migration squashing
```

### Design Patterns

- **Strategy Pattern** - Database-specific SQL generation
- **Builder Pattern** - Migration plan construction
- **Factory Pattern** - Generator creation with dependencies
- **Repository Pattern** - State management abstraction
- **Observer Pattern** - Migration event handling (ready for extension)

## Usage Examples

### Generate Migrations

```go
gen, err := migrations.NewGenerator("./models", "./migrations")
if err != nil {
    log.Fatal(err)
}

err = gen.GenerateMigrations("add_user_table")
if err != nil {
    log.Fatal(err)
}
```

### Validate Migrations

```go
validator := validation.NewValidator()
plan := &core.MigrationPlan{...}
if err := validator.ValidateMigration(plan); err != nil {
    log.Fatal(err)
}
```

### Lint Migrations

```go
linter := linting.NewLinter()
results, err := linter.LintMigrationsDir("./migrations")
for _, result := range results {
    fmt.Printf("%s: %s\n", result.Level, result.Message)
}
```

## Best Practices

1. **Always Review Generated Migrations** - Use `migrate show` to review before applying
2. **Use Dry-Run Mode** - Test migrations with `--dry-run` flag
3. **Lint Before Committing** - Run `migrate lint` before committing migrations
4. **Keep Migrations Small** - Split large changes into multiple migrations
5. **Test Down Migrations** - Ensure rollback works correctly
6. **Use Checksums** - Enable checksum validation in production
7. **Document Complex Changes** - Add comments for complex migrations

## Error Handling

The migration system provides comprehensive error handling:

- **Domain-Specific Errors** - `MigrationError` with error codes
- **Recovery Guidance** - Step-by-step recovery instructions
- **Validation Errors** - Clear validation failure messages
- **State Mismatch Detection** - Detect and report schema drift

## Production Considerations

- ✅ **Transaction Safety** - Migrations should be wrapped in transactions
- ✅ **Backup Before Migration** - Always backup before applying migrations
- ✅ **Test in Staging** - Test migrations in staging environment first
- ✅ **Monitor Migration Status** - Use `migrate status` to monitor state
- ✅ **Checksum Validation** - Validate migration integrity in production
- ✅ **Rollback Plan** - Always have a rollback plan ready

## Extension Points

The system is designed for extensibility:

- **Custom Change Types** - Implement `Change` interface for custom changes
- **Custom SQL Builders** - Implement `SQLBuilder` for new databases
- **Custom Validators** - Add custom validation rules
- **Custom Linters** - Add custom linting rules
- **Custom State Managers** - Implement `StateManager` for custom state storage

## Migration File Format

Migrations follow the standard format:
- `{version}_{name}.up.sql` - Up migration
- `{version}_{name}.down.sql` - Down migration

Example:
- `000001_create_users.up.sql`
- `000001_create_users.down.sql`

## State Management

The system maintains schema state in memory and can load from:
- Migration files (parsing SQL)
- Model definitions
- Database introspection (future)

## Performance

- Efficient change detection using maps for O(1) lookups
- Minimal SQL parsing (only when needed)
- Lazy state loading
- Optimized SQL generation

## Security

- SQL injection prevention through parameterized queries
- Dangerous operation detection
- Checksum validation to prevent tampering
- Validation of all changes before execution

