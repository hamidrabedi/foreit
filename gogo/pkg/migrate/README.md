# Migrate Module

Database migrations using [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

## Features

- SQL migration files
- Up/Down migrations
- Version tracking
- Multiple database drivers

## Usage

```go
import (
    "database/sql"
    "github.com/gogo/pkg/migrate"
    _ "github.com/lib/pq"
)

db, _ := sql.Open("postgres", "postgres://...")
migrator, _ := migrate.NewMigrator(db, "postgres", "./migrations")

// Run all pending migrations
migrator.Up()

// Rollback last migration
migrator.Down()

// Run specific number of migrations
migrator.Steps(2) // Up 2
migrator.Steps(-1) // Down 1

// Get current version
version, dirty, _ := migrator.Version()
```

## Migration Files

Create migration files in `migrations/` directory:

```
migrations/
  000001_create_users.up.sql
  000001_create_users.down.sql
  000002_add_email_index.up.sql
  000002_add_email_index.down.sql
```

