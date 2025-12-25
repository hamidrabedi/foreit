# ORM Module

Type-safe database operations using Ent's code generation.

## Concepts

### Repository Pattern

The Repository pattern provides a clean interface for data access:

```go
type UserRepository struct {
    orm.Repository[models.User, *ent.UserQuery]
}

// Type-safe queries
users := repo.Query().
    Where(user.NameEQ("John")).
    Limit(10).
    All(ctx)
```

### Client

The ORM client wraps Ent's client with additional utilities:

```go
client, err := orm.NewClient("postgres", "postgres://...")
// Use with Ent client
client.Open(entClient)
```

### Migrations

Migrations are handled through Ent's migration system. This module provides helpers:

```go
migrator := orm.NewMigrator(client)
migrator.Migrate(ctx)
```

### Bulk Operations

Bulk operations for efficient batch processing:

```go
bulk := orm.NewBulkOperations[models.User](repo)
results, err := bulk.BulkCreate(ctx, users)
```

## Features

- Type-safe repositories
- Query builders (using Ent's generated code)
- Transactions
- Bulk operations
- Migration helpers
- Connection management

