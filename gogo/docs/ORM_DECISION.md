# ORM & Migration Decision

## Final Architecture Decision

### Backend ORM: **GORM** ✅
- **Why**: Most mature, has generics API, handles SQL/scanning/relationships
- **What we use**: GORM for actual database operations
- **What we build**: Type-safe field descriptors + Django-style QuerySet wrapper

### Migrations: **golang-migrate/migrate** ✅
- **Why**: Industry standard, versioned migrations, SQL files, works with any ORM
- **What we use**: golang-migrate for versioned migrations
- **What we build**: CLI integration + auto-detection helpers

### Type Safety: **Codegen Field Descriptors** ✅
- **Why**: 100% compile-time safety, no strings, IDE autocomplete
- **What we build**: Codegen tool that generates typed field references from struct tags

## Architecture

```
User Code (Type-Safe)
    ↓
Generated Field Descriptors (codegen)
    ↓
Django-Style QuerySet Wrapper
    ↓
GORM Backend (handles SQL/scanning)
    ↓
Database

Migrations:
    ↓
golang-migrate (versioned SQL files)
    ↓
Database
```

## Why This Wins

1. **✅ Type-Safe**: Generated field descriptors = compile-time safety
2. **✅ Mature ORM**: GORM handles all the hard parts (SQL, scanning, relationships)
3. **✅ Versioned Migrations**: golang-migrate = industry standard
4. **✅ Django DX**: QuerySet API feels like Django
5. **✅ No Rebuilding**: Use proven tools, build thin wrapper
6. **✅ Production Ready**: GORM + golang-migrate = battle-tested

## What We Build

1. **Field Descriptor Generator** - Creates typed field references
2. **QuerySet Wrapper** - Django-style API over GORM
3. **Migration CLI** - Integrates golang-migrate with auto-detection
4. **Admin Integration** - Uses field descriptors for type-safe admin

## What We Don't Build

- ❌ SQL builder (GORM does this)
- ❌ Row scanning (GORM does this)
- ❌ Migration engine (golang-migrate does this)
- ❌ Connection pooling (GORM does this)

## Migration Workflow (Django-like)

```bash
# 1. Make model changes
# Edit models/user.go

# 2. Auto-detect changes
gogo makemigrations

# Generates: migrations/0001_initial.up.sql
#          migrations/0001_initial.down.sql

# 3. Apply migrations
gogo migrate

# 4. Rollback if needed
gogo migrate down 1
```

