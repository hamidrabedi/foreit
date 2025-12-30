---
sidebar_position: 5
---

# Migrations Guide

forge includes a built-in migration system based on golang-migrate that automatically manages your database schema changes.

## Overview

Migrations are version-controlled database schema changes. They allow you to:

- Track schema changes over time
- Apply changes to different environments consistently
- Roll back changes if needed
- Collaborate with team members on schema changes

## Creating Migrations

### Automatic Migration Generation

Generate migrations from your model definitions:

```bash
forge makemigrations
```

This will:
1. Scan your models directory
2. Compare current models with database state
3. Generate migration files for any changes
4. Save migrations to `migrations/` directory

### Migration Files

Migrations are stored in `migrations/` directory:

```
migrations/
├── 000001_initial.up.sql
├── 000001_initial.down.sql
├── 000002_add_user_table.up.sql
├── 000002_add_user_table.down.sql
└── ...
```

Each migration has:
- **Up migration** (`.up.sql`) - Applies the change
- **Down migration** (`.down.sql`) - Reverts the change

## Applying Migrations

### Apply All Pending Migrations

```bash
forge migrate
```

This applies all migrations that haven't been run yet.

### Apply Specific Migration

```bash
forge migrate up 2
```

Applies migrations up to version 2.

### Rollback

```bash
forge migrate down 1
```

Rolls back the last migration.

### Check Migration Status

```bash
forge migrate status
```

Shows which migrations have been applied.

## Migration Examples

### Creating a Table

**Up Migration:**
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(150) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**Down Migration:**
```sql
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

### Adding a Column

**Up Migration:**
```sql
ALTER TABLE users ADD COLUMN last_login TIMESTAMP;
```

**Down Migration:**
```sql
ALTER TABLE users DROP COLUMN last_login;
```

### Adding a Foreign Key

**Up Migration:**
```sql
ALTER TABLE posts ADD COLUMN author_id BIGINT;
ALTER TABLE posts ADD CONSTRAINT fk_posts_author 
    FOREIGN KEY (author_id) REFERENCES users(id) 
    ON DELETE CASCADE;
```

**Down Migration:**
```sql
ALTER TABLE posts DROP CONSTRAINT fk_posts_author;
ALTER TABLE posts DROP COLUMN author_id;
```

### Creating an Index

**Up Migration:**
```sql
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

**Down Migration:**
```sql
DROP INDEX IF EXISTS idx_posts_created_at;
```

## Data Migrations

You can also include data migrations in your SQL files:

```sql
-- Update existing data
UPDATE users SET is_active = true WHERE is_active IS NULL;

-- Insert default data
INSERT INTO categories (name, slug) VALUES
    ('Technology', 'technology'),
    ('Science', 'science'),
    ('Arts', 'arts');
```

## Best Practices

### 1. Keep Migrations Small

Break large changes into multiple migrations:

```sql
-- Good: Separate migrations
-- 000001_add_users_table.up.sql
-- 000002_add_posts_table.up.sql
-- 000003_add_comments_table.up.sql

-- Bad: One large migration
-- 000001_add_all_tables.up.sql
```

### 2. Always Write Down Migrations

Every up migration should have a corresponding down migration:

```sql
-- Up
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- Down
ALTER TABLE users DROP COLUMN phone;
```

### 3. Test Migrations

Test both up and down migrations:

```bash
# Apply
forge migrate

# Rollback
forge migrate down 1

# Re-apply
forge migrate up 1
```

### 4. Don't Modify Existing Migrations

Once a migration is applied to production, don't modify it. Create a new migration instead.

### 5. Use Transactions

Wrap migrations in transactions when possible:

```sql
BEGIN;

ALTER TABLE users ADD COLUMN phone VARCHAR(20);
CREATE INDEX idx_users_phone ON users(phone);

COMMIT;
```

## Migration Workflow

### Development

1. Modify your models
2. Generate migrations: `forge makemigrations`
3. Review generated SQL
4. Apply migrations: `forge migrate`
5. Test your application

### Production

1. Review migrations before deploying
2. Backup database
3. Apply migrations: `forge migrate`
4. Verify application works
5. Keep backup until confident

## Troubleshooting

### Migration Conflicts

If you have conflicting migrations:

```bash
# Check status
forge migrate status

# Manually resolve conflicts
# Edit migration files as needed
```

### Failed Migrations

If a migration fails:

1. Check the error message
2. Fix the SQL in the migration file
3. Rollback if needed: `forge migrate down 1`
4. Fix and re-apply: `forge migrate up 1`

### Database State Mismatch

If your database state doesn't match migrations:

```bash
# Check current state
forge migrate status

# Force to specific version (use with caution)
forge migrate force 5
```

## Advanced Topics

### Custom Migration SQL

You can write custom SQL migrations:

```sql
-- Custom migration logic
DO $$
BEGIN
    -- Complex migration logic
    IF EXISTS (SELECT 1 FROM information_schema.tables 
               WHERE table_name = 'old_table') THEN
        -- Migration code
    END IF;
END $$;
```

### Migration Hooks

Add hooks to run code during migrations:

```go
// In your migration file or code
func RunMigration(ctx context.Context, db *sql.DB) error {
    // Custom migration logic
    return nil
}
```

## Next Steps

- [Schema Reference](/docs/api-reference/schema) - Learn about model definitions
- [Development Guide](/docs/contributing/development) - Contributing to forge
- [Deployment Guide](/docs/guides/deployment) - Deploying your application
