# SQL Builder Documentation

## Overview

forge uses a custom SQL builder that generates safe SQL queries with proper escaping and parameter binding. This replaces the previous sqlc integration, providing a Go-first approach to query building.

## Architecture

The SQL builder is located in `pkg/query/sql_builder.go` and provides:

- **Identifier Escaping**: All table and column names are properly escaped to prevent SQL injection
- **Parameter Binding**: All values use PostgreSQL-style parameter placeholders (`$1`, `$2`, etc.)
- **Type Safety**: Works seamlessly with the type-safe QuerySet API

## Key Features

### 1. Identifier Escaping

```go
// All identifiers are escaped using double quotes (PostgreSQL standard)
EscapeIdentifier("users")        // Returns: "users"
EscapeIdentifier("user_name")   // Returns: "user_name"
EscapeIdentifier(`user"name`)    // Returns: "user""name" (escaped quotes)
```

### 2. Parameter Binding

```go
builder := NewSQLBuilder()
placeholder := builder.AddArg("john")  // Returns: "$1"
args := builder.Args()                  // Returns: []interface{}{"john"}
```

### 3. Query Building

The SQL builder provides methods for building all common SQL clauses:

- `BuildSelect()` - SELECT clause with field escaping
- `BuildWhere()` - WHERE clause with parameter binding
- `BuildOrderBy()` - ORDER BY clause with field escaping
- `BuildLimit()` - LIMIT clause
- `BuildOffset()` - OFFSET clause
- `BuildUpdate()` - UPDATE query with SET clause
- `BuildInsert()` - INSERT query with VALUES clause
- `BuildDelete()` - DELETE query

## Usage Example

```go
// QuerySet automatically uses SQL builder
users, err := User.Objects.
    Filter(User.Fields.Username.Equals("john")).
    Filter(User.Fields.IsActive.Equals(true)).
    OrderBy("-date_joined").
    Limit(10).
    All(ctx)

// Internally generates:
// SELECT * FROM "users"
// WHERE "username" = $1 AND "is_active" = $2
// ORDER BY "date_joined" DESC
// LIMIT 10
// Args: []interface{}{"john", true}
```

## Security

The SQL builder provides multiple layers of security:

1. **Identifier Escaping**: Prevents SQL injection through table/column names
2. **Parameter Binding**: All values are passed as parameters, never concatenated
3. **No String Concatenation**: SQL is built programmatically, not via string formatting

## Comparison with sqlc

| Feature           | sqlc               | SQL Builder        |
| ----------------- | ------------------ | ------------------ |
| Approach          | SQL-first          | Go-first           |
| Type Safety       | Generated from SQL | Built from Go code |
| Escaping          | Manual             | Automatic          |
| Parameter Binding | Manual             | Automatic          |
| Integration       | External tool      | Built-in           |

## Implementation Details

### QuerySet Integration

The `BaseQuerySet.buildSQL()` method uses the SQL builder:

```go
func (b *BaseQuerySet[T]) buildSQL() (string, []interface{}) {
    builder := NewSQLBuilder()

    // Build SELECT clause
    selectClause := builder.BuildSelect(b.table, selectFields, b.distinct)

    // Build WHERE clause
    whereClause, _ := builder.BuildWhere(b.conditions, b.excludes)

    // Build ORDER BY, LIMIT, OFFSET
    // ...

    return strings.Join(parts, " "), builder.Args()
}
```

### Executor

The `Executor` uses standard `database/sql`:

```go
type Executor struct {
    db *sql.DB
}

func (e *Executor) ExecuteRaw(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
    return e.db.QueryContext(ctx, sql, args...)
}
```

## Benefits

1. **No External Dependencies**: Uses standard `database/sql`
2. **Go-First**: SQL is generated from Go code, not the other way around
3. **Automatic Security**: Escaping and parameter binding are handled automatically
4. **Consistent API**: Works seamlessly with the type-safe QuerySet API
5. **Simple**: No code generation step required

## Future Enhancements

Potential improvements:

- Support for other database dialects (MySQL, SQLite)
- Query optimization hints
- Query caching
- Query logging/debugging utilities
