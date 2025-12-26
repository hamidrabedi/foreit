---
sidebar_position: 5
---

# Migrations Guide

forge uses golang-migrate for database migrations. Migrations allow you to version control your database schema and apply changes incrementally.

## Creating Migrations

### Automatic Migration Creation

Generate migrations from your models:

```bash
forge makemigrations
```

This creates migration files in the `migrations/` directory:

```
migrations/
├── 000001_create_users.up.sql
├── 000001_create_users.down.sql
├── 000002_create_posts.up.sql
└── 000002_create_posts.down.sql
```

### Manual Migration Creation

Create migration files manually:

```bash
forge makemigrations add_user_email_index
```

This creates:
- `migrations/XXXXX_add_user_email_index.up.sql`
- `migrations/XXXXX_add_user_email_index.down.sql`

Edit these files to add your SQL:

**up.sql:**
```sql
CREATE INDEX idx_user_email ON users(email);
```

**down.sql:**
```sql
DROP INDEX IF EXISTS idx_user_email;
```

## Applying Migrations

### Apply All Migrations

```bash
forge migrate
```

This applies all pending migrations in order.

### Apply Specific Migration

```bash
forge migrate up 2
```

Applies the next 2 migrations.

### Rollback Migrations

```bash
forge migrate down 1
```

Rolls back the last migration.

### Check Migration Status

```bash
forge migrate version
```

Shows the current migration version.

## Migration Files

### Up Migration

The `up.sql` file contains SQL to apply the migration:

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(150) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_email ON users(email);
```

### Down Migration

The `down.sql` file contains SQL to reverse the migration:

```sql
DROP INDEX IF EXISTS idx_user_email;
DROP TABLE IF EXISTS users;
```

## Best Practices

### 1. Always Create Down Migrations

Every migration should be reversible:

```sql
-- up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- down.sql
ALTER TABLE users DROP COLUMN phone;
```

### 2. Test Migrations

Test both up and down migrations:

```bash
# Apply
forge migrate up

# Rollback
forge migrate down

# Apply again
forge migrate up
```

### 3. Use Transactions

Wrap migrations in transactions when possible:

```sql
BEGIN;

CREATE TABLE posts (...);
CREATE INDEX idx_posts_author ON posts(author_id);

COMMIT;
```

### 4. Avoid Data Loss

Be careful with destructive operations:

```sql
-- Bad: Data loss
ALTER TABLE users DROP COLUMN email;

-- Good: Migrate data first
ALTER TABLE users ADD COLUMN new_email VARCHAR(255);
UPDATE users SET new_email = email;
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users RENAME COLUMN new_email TO email;
```

### 5. Use IF EXISTS / IF NOT EXISTS

Make migrations idempotent:

```sql
-- Good
CREATE INDEX IF NOT EXISTS idx_user_email ON users(email);
DROP INDEX IF EXISTS idx_user_email;

-- Bad
CREATE INDEX idx_user_email ON users(email);  -- Fails if exists
```

## Common Migration Patterns

### Adding a Column

```sql
-- up.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- down.sql
ALTER TABLE users DROP COLUMN phone;
```

### Removing a Column

```sql
-- up.sql
ALTER TABLE users DROP COLUMN phone;

-- down.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

### Changing Column Type

```sql
-- up.sql
ALTER TABLE users ALTER COLUMN age TYPE INTEGER USING age::INTEGER;

-- down.sql
ALTER TABLE users ALTER COLUMN age TYPE VARCHAR(10) USING age::VARCHAR;
```

### Adding an Index

```sql
-- up.sql
CREATE INDEX idx_user_email ON users(email);

-- down.sql
DROP INDEX IF EXISTS idx_user_email;
```

### Adding a Foreign Key

```sql
-- up.sql
ALTER TABLE posts ADD COLUMN author_id BIGINT;
ALTER TABLE posts ADD CONSTRAINT fk_posts_author 
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;

-- down.sql
ALTER TABLE posts DROP CONSTRAINT IF EXISTS fk_posts_author;
ALTER TABLE posts DROP COLUMN author_id;
```

### Adding a Unique Constraint

```sql
-- up.sql
ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);

-- down.sql
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_username;
```

## Migration from Models

When you change your models, regenerate migrations:

```bash
# 1. Update your model
# Edit models/user.go

# 2. Generate migrations
forge makemigrations

# 3. Review generated SQL
# Check migrations/XXXXX_*.sql files

# 4. Apply migrations
forge migrate
```

## Production Migrations

### Backup First

Always backup before applying migrations in production:

```bash
pg_dump myapp_db > backup.sql
```

### Test in Staging

Test migrations in a staging environment first.

### Apply During Maintenance Window

Schedule migrations during low-traffic periods.

### Monitor

Watch for errors and performance issues after applying migrations.

## Troubleshooting

### Migration Already Applied

If a migration was partially applied:

```bash
# Check current version
forge migrate version

# Force to specific version
forge migrate force 5
```

### Migration Failed

If a migration fails:

1. Check the error message
2. Fix the SQL in the migration file
3. Rollback if needed: `forge migrate down`
4. Fix and reapply: `forge migrate up`

### Database Out of Sync

If your database is out of sync:

```bash
# Check migration status
forge migrate version

# Apply missing migrations
forge migrate up
```

## Next Steps

- [Models Guide](/docs/guides/models) - Learn about model definitions
- [Database Guide](/docs/reference/schema) - Database schema reference
- [Deployment Guide](/docs/guides/deployment) - Deploy with migrations

