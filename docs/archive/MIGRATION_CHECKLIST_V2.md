# Migration Checklist for v2.0 Naming Changes

**Version**: v2.0  
**Date**: 2026-01-01

This checklist helps you migrate your code from v1.x to v2.0 naming conventions.

## Automated Migration

Run the automated migration script:

```bash
./scripts/migrate_naming_v2.sh
```

This will automatically update most of your code. Then manually review the changes.

## Manual Checklist

### ✅ Fix 1: Database Rebind

- [ ] Search for `.Rebind(` in your code
- [ ] Replace with `.RebindPlaceholders(`
- [ ] Update any `database.Rebind(` calls
- [ ] Update any `db.Rebind(` calls
- [ ] Test database operations

**Example**:
```go
// Old
sql := db.Rebind(query)

// New
sql := db.RebindPlaceholders(query)
```

### ✅ Fix 2: Registry Access

- [ ] Search for `registry.GetRegistry()` in your code
- [ ] Replace with `registry.Global()`
- [ ] Test registry operations

**Example**:
```go
// Old
reg := registry.GetRegistry()

// New
reg := registry.Global()
```

### ✅ Fix 3: Ordering Fields

- [ ] Search for `NewOrderField(` in your code
- [ ] Replace `NewOrderField(field, true)` with `Asc(field)`
- [ ] Replace `NewOrderField(field, false)` with `Desc(field)`
- [ ] Test ordering operations

**Example**:
```go
// Old
qs.OrderBy(orm.NewOrderField("created_at", true))
qs.OrderBy(orm.NewOrderField("price", false))

// New
qs.OrderBy(orm.Asc("created_at"))
qs.OrderBy(orm.Desc("price"))
```

### ✅ Fix 4: Field Accessor

- [ ] Search for `.GetFieldAccessor()` in your code
- [ ] Replace with `.FieldAccessor()`
- [ ] Test field access operations

**Example**:
```go
// Old
fa, err := manager.GetFieldAccessor()

// New
fa, err := manager.FieldAccessor()
```

### ✅ Fix 5: Pagination Parameters

- [ ] Search for `GetPaginationParams(` in your code
- [ ] Replace with `ParsePaginationParams(`
- [ ] Update any `api.GetPaginationParams(` calls
- [ ] Test pagination functionality

**Example**:
```go
// Old
page, size, offset := api.GetPaginationParams(r, 20)

// New
page, size, offset := api.ParsePaginationParams(r, 20)
```

## Testing

After migration:

- [ ] Run all tests: `go test ./...`
- [ ] Build all packages: `go build ./...`
- [ ] Test database operations
- [ ] Test pagination
- [ ] Test ordering
- [ ] Test field accessors

## Verification

Check that:

- [ ] No deprecated function warnings in your code
- [ ] All tests pass
- [ ] All builds succeed
- [ ] No linter errors
- [ ] Code review completed

## Rollback Plan

If you need to rollback:

1. Use git to revert changes: `git revert <commit>`
2. Or restore from backup branch created by migration script
3. Old functions still work (they're deprecated, not removed)

## Support

If you encounter issues:

1. Check [MIGRATION_V1_TO_V2.md](MIGRATION_V1_TO_V2.md) for detailed guide
2. Review [NAMING_AUDIT_V2.md](NAMING_AUDIT_V2.md) for context
3. Open an issue on GitHub
4. Ask for help on Discord

## Timeline

- **v2.0 Release**: Deprecations added (old functions still work)
- **v2.x Releases**: Deprecation warnings continue
- **v3.0 Release**: Deprecated functions removed (6+ months from v2.0)

You have plenty of time to migrate!
