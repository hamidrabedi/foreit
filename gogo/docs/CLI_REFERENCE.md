# CLI Reference

Gogo provides a command-line interface for common tasks.

## Commands

### startproject

Create a new Gogo project.

```bash
gogo startproject myproject
```

Creates:
- Project structure
- Main application file
- Configuration files
- Example resources

### startapp

Create a new app within a project.

```bash
gogo startapp myapp
```

Creates:
- App directory
- Basic structure
- Example models

### generate

Generate code for common patterns.

#### Generate Resource

```bash
gogo generate resource User
```

Creates a resource handler with CRUD operations.

#### Generate Console

```bash
gogo generate console User
```

Creates a console (admin) interface.

#### Generate Policy

```bash
gogo generate policy User
```

Creates an authorization policy.

#### Generate Serializer

```bash
gogo generate serializer User
```

Creates a serializer for data transformation.

### migrate

Run database migrations.

```bash
# Run all pending migrations
gogo migrate up

# Rollback last migration
gogo migrate down

# Show migration status
gogo migrate status

# Create new migration
gogo migrate create add_users_table
```

## Examples

### Complete Workflow

```bash
# Create project
gogo startproject myapp
cd myapp

# Create app
gogo startapp blog

# Generate resources
gogo generate resource Post
gogo generate resource Comment

# Generate consoles
gogo generate console Post
gogo generate console Comment

# Run migrations
gogo migrate up

# Start server
go run main.go
```

