# Forge Ecommerce - Complete CLI Usage Guide

This guide demonstrates all Forge CLI commands using the ecommerce example project.

## Table of Contents

1. [Project Creation](#project-creation)
2. [Code Generation](#code-generation)
3. [Migration Management](#migration-management)
4. [Development Server](#development-server)
5. [Admin Management](#admin-management)
6. [Development Tools](#development-tools)
7. [Adding Features](#adding-features)
8. [Testing](#testing)
9. [Production Deployment](#production-deployment)

---

## Project Creation

### Create a New Project

```bash
# Basic project creation
forge new ecommerce

# With options
forge new ecommerce --template=simple --database=postgres --path=./myshop

# With Docker setup
forge new ecommerce --docker

# Non-interactive mode
forge new ecommerce --template=simple --database=postgres --docker
```

**What this creates:**
- Project structure with proper Go module
- Configuration files
- Example models
- Database configuration
- Static and template directories

**Options:**
- `--template`: `simple` (default) or `advanced` (clean architecture)
- `--database`: `postgres` (default), `mysql`, or `sqlite`
- `--path`: Custom project path
- `--docker`: Include Docker and compose setup

---

## Code Generation

The `generate` command is the heart of Forge - it creates type-safe code from your model definitions.

### Basic Generation

```bash
# Generate code from all models
cd ecommerce
forge generate

# Output:
# Scanning apps in app/...
#   Generating for catalog...
#   Generating for customers...
#   Generating for orders...
#   Generating for inventory...
#   Generating for marketing...
# ✓ Generated code for 5 apps
```

**What this generates for each model:**
- `*_gen.go` - Model struct with all fields
- `*_fields_gen.go` - Type-safe field expressions
- `*_manager_gen.go` - Manager with CRUD operations
- `*_queryset_gen.go` - QuerySet for filtering

### Targeted Generation

```bash
# Generate for specific app
forge generate --models=./app/catalog

# Generate for specific directory
forge generate --models=./app/catalog --output=./app/catalog

# Dry run (preview without writing)
forge generate --dry-run
```

### Generation Examples

**Before generation** (`app/catalog/models.go`):
```go
type Product struct {
    schema.BaseSchema
}

func (Product) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").WithPrimary().WithAutoIncrement(),
        schema.String("name").WithRequired().WithMaxLength(255),
        schema.Float64("price").WithRequired(),
    }
}
```

**After generation** (`app/catalog/product_gen.go`):
```go
type ProductGenerated struct {
    ID    int64   `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

var ProductFields = struct{
    ID    *FieldExpr[int64]
    Name  *FieldExpr[string]
    Price *FieldExpr[float64]
}{
    ID:    NewFieldExpr[int64]("id"),
    Name:  NewFieldExpr[string]("name"),
    Price: NewFieldExpr[float64]("price"),
}

var ProductObjects = NewManager[Product]()
```

---

## Migration Management

Forge uses golang-migrate for database migrations.

### Create Migrations

```bash
# Auto-generate migrations from models
forge makemigrations

# Output:
# Created new migration: 001_initial_schema.up.sql
# Created new migration: 001_initial_schema.down.sql
```

**Generated SQL** (`migrations/001_initial_schema.up.sql`):
```sql
-- Categories table
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL UNIQUE,
    slug VARCHAR(200) NOT NULL UNIQUE,
    description TEXT,
    parent_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    image_url VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    level INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_category_slug ON categories(slug);
CREATE INDEX idx_category_parent ON categories(parent_id);
CREATE INDEX idx_category_active ON categories(is_active);

-- Products table
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    sku VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    short_description TEXT,
    price DOUBLE PRECISION NOT NULL,
    cost_price DOUBLE PRECISION,
    compare_at_price DOUBLE PRECISION,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    brand_id BIGINT REFERENCES brands(id) ON DELETE SET NULL,
    stock_quantity INTEGER DEFAULT 0,
    -- ... more fields ...
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_slug ON products(slug);
CREATE INDEX idx_product_sku ON products(sku);
CREATE INDEX idx_product_category ON products(category_id);
-- ... more indexes ...
```

### Apply Migrations

```bash
# Apply all pending migrations
forge migrate

# Apply migrations with verbose output
forge migrate --verbose

# Apply specific number of migrations
forge migrate --steps=1

# Output:
# Running migrations...
# ✓ Applied 001_initial_schema.up.sql
# ✓ Applied 002_add_product_variants.up.sql
# Database schema is up to date
```

### Migration Status

```bash
# Check migration status
forge migrate status

# Output:
# Migration Status:
# ✓ 001_initial_schema              Applied at 2024-01-15 10:30:00
# ✓ 002_add_product_variants         Applied at 2024-01-15 10:30:01
# ✓ 003_add_inventory_tables         Applied at 2024-01-15 10:30:02
# ○ 004_add_reviews_tables           Pending
```

### Rollback Migrations

```bash
# Rollback last migration
forge migrate rollback

# Rollback specific number of migrations
forge migrate rollback --steps=2

# Rollback to specific version
forge migrate rollback --version=001

# Output:
# Rolling back migrations...
# ✓ Rolled back 004_add_reviews_tables.down.sql
# ✓ Rolled back 003_add_inventory_tables.down.sql
```

### Show Migration SQL

```bash
# Show SQL for a migration
forge migrate show 001_initial_schema

# Output:
# Migration: 001_initial_schema
# Type: up
# ---
# CREATE TABLE categories (
#   id BIGSERIAL PRIMARY KEY,
#   ...
# );
```

### Lint Migrations

```bash
# Check migrations for issues
forge migrate lint

# Output:
# Linting migrations...
# ✓ 001_initial_schema              No issues
# ✓ 002_add_product_variants         No issues
# ⚠ 003_add_inventory_tables         Warning: Large data migration, consider batching
# ✗ 004_add_reviews_tables           Error: Missing down migration
```

### Advanced Migration Commands

```bash
# Force migration version (use with caution)
forge migrate force 003

# Mark migration as applied without running
forge migrate fake 004

# Squash migrations (combine multiple into one)
forge migrate squash --from=001 --to=005 --name=initial_combined

# Create empty migration (for data migrations)
forge migrate create add_sample_products

# Group migrations (for team collaboration)
forge migrate group --name=sprint-1 --from=001 --to=010
```

---

## Development Server

### Run Server

```bash
# Start development server
forge runserver

# Output:
# Starting Forge Ecommerce System...
# Database connected successfully
# Models registered
# Admin models registered
# API routes registered at /api/v1
# Admin interface registered at /admin/
# 
# 🚀 Server running at http://localhost:8000
# 📊 Admin: http://localhost:8000/admin/
# 🔌 API: http://localhost:8000/api/v1/
```

### Custom Host and Port

```bash
# Custom port
forge runserver --port=9000

# Custom host (for Docker, etc.)
forge runserver --host=0.0.0.0 --port=8080

# With auto-reload (watches for file changes)
forge runserver --reload

# Output:
# 🚀 Server running at http://0.0.0.0:8080
# 📊 Admin: http://0.0.0.0:8080/admin/
# 🔌 API: http://0.0.0.0:8080/api/v1/
# 
# Watching for changes...
```

### Direct Execution

```bash
# Alternative: run directly with go
go run main.go

# With hot reload using air
air
```

---

## Admin Management

### Create Superuser

```bash
# Interactive mode
forge createsuperuser

# Output:
# Email: admin@example.com
# Password: ********
# Confirm Password: ********
# First Name: Admin
# Last Name: User
# ✓ Superuser created successfully

# Non-interactive mode
forge createsuperuser \
    --email=admin@example.com \
    --password=SecurePass123! \
    --first-name=Admin \
    --last-name=User

# Create regular user
forge createsuperuser --no-superuser
```

### Admin Access

Once created, access the admin at:
- http://localhost:8000/admin/

**Admin features available:**
- 📦 Products - Full CRUD, bulk actions, CSV export
- 👥 Customers - Search, filter, export
- 📋 Orders - View, process, print invoices
- 📊 Inventory - Stock management, alerts
- ⭐ Reviews - Moderate, approve/reject
- 🎟️ Coupons - Create, manage promotions

---

## Development Tools

### Check Project

```bash
# Run comprehensive project checks
forge check

# Output:
# Running project checks...
# 
# ✓ Database connection
# ✓ Model definitions
# ✓ Migration status
# ✓ Generated code up to date
# ✗ Linter errors found (3)
# ⚠ Security warnings (1)
# 
# Linter Errors:
#   app/catalog/models.go:45: unused variable 'ctx'
#   app/orders/services.go:120: error return value not checked
# 
# Security Warnings:
#   config/config.yaml: Using default session secret
# 
# Run 'forge check --fix' to auto-fix issues
```

### Interactive Shell

```bash
# Start interactive Go shell with models loaded
forge shell

# Output:
# Forge Interactive Shell
# Models loaded: Category, Product, Order, Customer, Stock, Review
# Database connected
# 
# >>> 

# Example usage in shell:
>>> products, _ := catalog.Product.Objects.All(ctx)
>>> fmt.Printf("Total products: %d\n", len(products))
Total products: 150

>>> customer, _ := customers.Customer.Objects.Get(ctx, 1)
>>> fmt.Printf("Customer: %s %s\n", customer.FirstName, customer.LastName)
Customer: John Doe
```

### Run Tests

```bash
# Run all tests
forge test

# Run specific package tests
forge test ./app/catalog

# With coverage
forge test --cover

# With verbose output
forge test -v

# Output:
# Running tests...
# 
# app/catalog
#   ✓ TestProductModel (0.01s)
#   ✓ TestCategoryHierarchy (0.02s)
#   ✓ TestProductVariants (0.01s)
# 
# app/orders
#   ✓ TestOrderCreation (0.03s)
#   ✓ TestOrderProcessing (0.05s)
#   ✓ TestPaymentFlow (0.04s)
# 
# Coverage: 87.5%
# 
# PASS
# Total: 45 tests, 45 passed, 0 failed
# Time: 2.5s
```

---

## Adding Features

### Add New App

```bash
# Create new app module
forge add app loyalty

# Output:
# Created app structure:
#   app/loyalty/
#   app/loyalty/models.go
#   app/loyalty/admin.go
#   app/loyalty/api.go
#   app/loyalty/filters.go
# ✓ App 'loyalty' created successfully
```

### Add New Model

```bash
# Add model interactively
forge add model LoyaltyPoints --app=loyalty

# Output:
# Model name: LoyaltyPoints
# Table name (loyalty_points): [Enter]
# 
# Add fields:
# Field name: customer_id
# Field type: (String/Int64/Bool/Time/Float64): Int64
# Required? (y/n): y
# Add another field? (y/n): y
# 
# Field name: points
# Field type: Int64
# Required? (y/n): y
# Add another field? (y/n): n
# 
# ✓ Added model LoyaltyPoints to app/loyalty/models.go
```

### Add API Endpoint

```bash
# Add custom API endpoint
forge add handler calculate_loyalty \
    --app=loyalty \
    --method=POST \
    --path=/api/v1/loyalty/calculate

# Output:
# Created handler: app/loyalty/handlers.go
# ✓ Handler 'calculate_loyalty' added
```

### Add Service

```bash
# Add business logic service
forge add service LoyaltyService --app=loyalty

# Output:
# Created service: app/loyalty/services.go
# ✓ Service 'LoyaltyService' added
```

### Add Filter

```bash
# Add custom filter
forge add filter ActiveProductsFilter --app=catalog

# Output:
# Created filter: app/catalog/filters.go
# ✓ Filter 'ActiveProductsFilter' added
```

---

## Testing Workflow

### Complete Testing Flow

```bash
# 1. Check code
forge check

# 2. Run tests
forge test

# 3. Generate code
forge generate

# 4. Create migrations
forge makemigrations

# 5. Apply migrations (to test DB)
DATABASE_URL=postgres://localhost/ecommerce_test forge migrate

# 6. Run integration tests
forge test --integration

# 7. Check coverage
forge test --cover --html

# 8. Run linters
forge lint

# 9. Security check
forge check --security
```

---

## Production Deployment

### Build for Production

```bash
# Build optimized binary
go build -o ecommerce -ldflags="-s -w" main.go

# Build with version info
go build -o ecommerce \
    -ldflags="-X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y%m%d%H%M%S)" \
    main.go
```

### Migration Commands for Production

```bash
# Dry run migrations (shows SQL without applying)
forge migrate --dry-run

# Apply migrations with transaction safety
forge migrate --safe

# Backup before migration
forge migrate --backup

# Apply with timeout
forge migrate --timeout=300s
```

### Docker Deployment

```bash
# Build Docker image
docker build -t ecommerce:latest .

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f app

# Run migrations in container
docker-compose exec app forge migrate

# Create superuser in container
docker-compose exec app forge createsuperuser
```

---

## Advanced CLI Usage

### Configuration

```bash
# Use custom config file
forge runserver --config=./config/production.yaml

# Override config values
forge runserver --db-host=postgres.example.com --db-name=prod_db

# Environment variables
export FORGE_DB_HOST=localhost
export FORGE_DB_PORT=5432
forge runserver
```

### Debugging

```bash
# Enable debug mode
forge runserver --debug

# Verbose output
forge migrate --verbose

# Show SQL queries
forge runserver --log-sql

# Profile performance
forge runserver --profile --profile-port=6060
```

### Batch Operations

```bash
# Export data
forge export products --format=csv --output=products.csv

# Import data
forge import products --from=products.csv

# Seed database
forge seed --file=seeds/sample_data.yaml

# Clear data
forge flush --model=Product --confirm
```

---

## Complete Example Workflow

Here's a complete workflow from project creation to deployment:

```bash
# 1. Create project
forge new ecommerce --template=simple --database=postgres --docker
cd ecommerce

# 2. Configure database
# Edit config/config.yaml with your database credentials

# 3. Create database
createdb ecommerce_db

# 4. Generate code from models
forge generate

# 5. Create and apply migrations
forge makemigrations
forge migrate

# 6. Check everything is working
forge check

# 7. Create admin user
forge createsuperuser

# 8. Start development server
forge runserver

# 9. Access the application
# - Homepage: http://localhost:8000
# - Admin: http://localhost:8000/admin/
# - API: http://localhost:8000/api/v1/products/

# 10. Make changes to models
# Edit app/catalog/models.go

# 11. Regenerate code
forge generate

# 12. Create new migrations
forge makemigrations

# 13. Apply new migrations
forge migrate

# 14. Run tests
forge test

# 15. Build for production
go build -o ecommerce main.go

# 16. Deploy
./ecommerce
```

---

## CLI Tips and Tricks

### Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
alias fr='forge runserver'
alias fg='forge generate'
alias fm='forge migrate'
alias fmm='forge makemigrations'
alias ft='forge test'
alias fc='forge check'
```

### Auto-completion

```bash
# Enable bash completion
forge completion bash > /etc/bash_completion.d/forge

# Enable zsh completion
forge completion zsh > ~/.zsh/completion/_forge
```

### Watch Mode

```bash
# Auto-regenerate on file changes
forge watch generate

# Auto-migrate on schema changes
forge watch makemigrations
```

### Scripts

Create `scripts/dev.sh`:

```bash
#!/bin/bash
forge generate && \
forge makemigrations && \
forge migrate && \
forge check && \
forge runserver --reload
```

---

## Troubleshooting

### Common Issues

**Issue: "Model not found"**
```bash
# Solution: Regenerate code
forge generate
```

**Issue: "Migration already exists"**
```bash
# Solution: Check status and force if needed
forge migrate status
forge migrate force 005
```

**Issue: "Database connection failed"**
```bash
# Solution: Check configuration
forge check
# Or test connection
forge shell
```

**Issue: "Port already in use"**
```bash
# Solution: Use different port
forge runserver --port=9000
# Or kill existing process
lsof -ti:8000 | xargs kill
```

---

## Help and Documentation

```bash
# General help
forge --help

# Command-specific help
forge migrate --help
forge generate --help
forge runserver --help

# Version info
forge version

# List all commands
forge commands
```

---

## Summary

This ecommerce example demonstrates all Forge CLI commands:

✅ **Project Management**: new, add app/model/handler
✅ **Code Generation**: generate (with all options)
✅ **Migrations**: makemigrations, migrate, rollback, status, show, lint, force, fake, squash
✅ **Development**: runserver, check, shell, test
✅ **Admin**: createsuperuser
✅ **Deployment**: build, docker integration

Every command has been designed to make Django developers feel at home while leveraging Go's type safety and performance.

