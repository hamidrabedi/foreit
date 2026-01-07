# Migration Guide: v1.x to v2.0

This guide helps you migrate your code from Forge v1.x to v2.0, focusing on the naming and API changes introduced in v1.5.0.

## Overview

The main changes in v2.0 are:
1. **Query Expression API**: `NewQ()` → Use `And()`, `Or()`, `Not()` directly
2. **Field References**: `NewFieldQueryExpr()` → Use `Where()` or field methods
3. **Field Types**: `FieldExpression[T]` → `Field[T]` (alias provided for compatibility)
4. **Field References**: New `F()` and `FieldRef()` functions for runtime field references

## Breaking Changes

### 1. NewQ() Removed

**Before (v1.x)**:
```go
q := orm.NewQ(nameField.Eq("John"))
qs.Filter(q)
```

**After (v2.0)**:
```go
// Option 1: Direct (simplest)
qs.Filter(nameField.Eq("John"))

// Option 2: Complex queries with And/Or
qs.Filter(And(
    nameField.Eq("John"),
    ageField.Gt(18),
))
```

**Migration Steps**:
1. Replace `orm.NewQ(expr)` with just `expr` when filtering directly
2. Replace `q.Or(orm.NewQ(other))` with `Or(expr, other)`
3. Replace `q.And(orm.NewQ(other))` with `And(expr, other)`

### 2. NewFieldQueryExpr() Removed

**Before (v1.x)**:
```go
expr := orm.NewFieldQueryExpr("age", orm.OpGreater, 18)
qs.Filter(orm.NewQ(expr))
```

**After (v2.0)**:
```go
// Option 1: Where() function (explicit, SQL-like)
qs.Filter(Where("age", OpGreater, 18))

// Option 2: Type-safe field methods (best)
qs.Filter(User.Age.Gt(18))
```

**Migration Steps**:
1. Replace `NewFieldQueryExpr(field, op, value)` with `Where(field, op, value)`
2. Prefer type-safe field methods when available: `User.Age.Gt(18)`

### 3. FieldExpression → Field

**Before (v1.x)**:
```go
field := orm.NewField[string]("name", "users")
// Returns FieldExpression[string]
```

**After (v2.0)**:
```go
field := orm.NewField[string]("name", "users")
// Returns Field[string] (FieldExpression is an alias for backward compatibility)
```

**Migration Steps**:
- No code changes needed - `FieldExpression[T]` is an alias for `Field[T]`
- Update type annotations to use `Field[T]` in new code
- `FieldExpression[T]` will be removed in v3.0

### 4. FieldExpr Deprecated

**Before (v1.x)**:
```go
field := orm.NewFieldExpr[string]("name", "users")
expr := field.Equals("John")
```

**After (v2.0)**:
```go
field := orm.NewField[string]("name", "users")
expr := field.Eq("John")
```

**Migration Steps**:
1. Replace `NewFieldExpr` with `NewField`
2. Replace method names:
   - `Equals()` → `Eq()`
   - `NotEquals()` → `Ne()`
   - `Greater()` → `Gt()`
   - `GreaterOrEqual()` → `Gte()`
   - `Less()` → `Lt()`
   - `LessOrEqual()` → `Lte()`

## New Features

### 1. F() and FieldRef() for Runtime Field References

**New in v1.5.0**:
```go
// Django-style (short)
qs.Filter(F("age").Gt(18))

// Explicit alternative
qs.Filter(FieldRef("age").Gt(18))

// Type-safe (best, when available)
qs.Filter(User.Age.Gt(18))
```

### 2. Where() Function

**New in v1.5.0**:
```go
// SQL-like, explicit
qs.Filter(Where("age", OpGreater, 18))
qs.Filter(Where("name", OpEquals, "John"))
```

### 3. And(), Or(), Not() Functions

**New in v1.5.0**:
```go
// Explicit boolean combiners
qs.Filter(And(
    User.Name.Eq("John"),
    Or(
        User.Age.Gt(18),
        User.Role.Eq("admin"),
    ),
))

// Negation
qs.Exclude(Not(User.Age.Gt(65)))
```

## Detailed Migration Examples

### Example 1: Simple Filter

**Before**:
```go
q := orm.NewQ(nameField.Eq("John"))
users, err := userManager.Filter(q).All(ctx)
```

**After**:
```go
// Direct (simplest)
users, err := userManager.Filter(nameField.Eq("John")).All(ctx)
```

### Example 2: Boolean Logic

**Before**:
```go
q1 := orm.NewQ(nameField.Eq("John"))
q2 := orm.NewQ(ageField.Gt(18))
combined := q1.Or(q2)
users, err := userManager.Filter(combined).All(ctx)
```

**After**:
```go
users, err := userManager.Filter(Or(
    nameField.Eq("John"),
    ageField.Gt(18),
)).All(ctx)
```

### Example 3: Complex Queries

**Before**:
```go
q1 := orm.NewQ(nameField.Contains("John"))
q2 := orm.NewQ(ageField.Gte(18)).And(orm.NewQ(ageField.Lte(65)))
q3 := orm.NewQ(statusField.Eq("active"))
final := q1.And(q2).And(q3)
users, err := userManager.Filter(final).All(ctx)
```

**After**:
```go
users, err := userManager.Filter(And(
    nameField.Contains("John"),
    ageField.Gte(18),
    ageField.Lte(65),
    statusField.Eq("active"),
)).All(ctx)
```

### Example 4: String-Based (Without Type-Safe Fields)

**Before**:
```go
expr := orm.NewFieldQueryExpr("email", orm.OpContains, "@example.com")
users, err := userManager.Filter(orm.NewQ(expr)).All(ctx)
```

**After**:
```go
// Option 1: Where() function
users, err := userManager.Filter(
    Where("email", OpContains, "@example.com"),
).All(ctx)

// Option 2: F() function
users, err := userManager.Filter(
    F("email").Contains("@example.com"),
).All(ctx)
```

### Example 5: FieldExpr Migration

**Before**:
```go
field := orm.NewFieldExpr[string]("name", "users")
expr := field.Equals("John")
users, err := qs.Filter(orm.NewQ(expr)).All(ctx)
```

**After**:
```go
field := orm.NewField[string]("name", "users")
users, err := qs.Filter(field.Eq("John")).All(ctx)
```

## Compatibility Layer

During the v1.5.0 → v2.0 transition period, deprecated functions remain available with warnings:

- `NewQ()` - Still works, but logs deprecation warning
- `NewFieldQueryExpr()` - Still works, but logs deprecation warning
- `FieldExpression[T]` - Type alias, fully compatible
- `FieldExpr[T]` - Still works, but deprecated

## Automated Migration

### Using grep/sed (Simple Cases)

```bash
# Replace NewQ wrapping (simple cases)
find . -name "*.go" -exec sed -i 's/orm\.NewQ(\([^)]*\))/\1/g' {} \;

# Replace NewFieldQueryExpr with Where
find . -name "*.go" -exec sed -i 's/orm\.NewFieldQueryExpr(\([^,]*\), \(.*\), \(.*\))/orm.Where(\1, \2, \3)/g' {} \;
```

**Note**: Manual review required for complex cases.

### Using gorename (For Type Renames)

```bash
# Rename FieldExpression to Field (if needed)
gorename -from '"github.com/forgego/forge/orm".FieldExpression' -to 'Field'
```

## Testing Your Migration

After migrating:

1. **Run tests**: `go test ./...`
2. **Check for deprecation warnings**: Look for log messages about deprecated functions
3. **Verify functionality**: Ensure queries return expected results
4. **Update documentation**: Update any internal docs or examples

## Timeline

- **v1.5.0** (Current): New API available, old API deprecated with warnings
- **v1.6.0 - v1.9.0**: Migration period (6 months)
- **v2.0.0**: Old API removed, breaking changes take effect

## Getting Help

If you encounter issues during migration:

1. Check this guide for common patterns
2. Review the [API Reference](./API_REFERENCE.md)
3. Check [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md) for design rationale
4. Open an issue on GitHub with your specific case

## v2.0 Additional Naming Changes

In v2.0, we've improved naming consistency across the framework. These changes follow Go conventions and make the API clearer.

### 5. Rebind() → RebindPlaceholders()

**Before (v1.x)**:
```go
sql := db.Rebind(query)
```

**After (v2.0)**:
```go
sql := db.RebindPlaceholders(query)
```

**Migration**: Use automated script or search/replace `.Rebind(` with `.RebindPlaceholders(`

### 6. GetRegistry() → Global()

**Before (v1.x)**:
```go
registry := registry.GetRegistry()
```

**After (v2.0)**:
```go
registry := registry.Global()
```

**Migration**: Replace `registry.GetRegistry()` with `registry.Global()`

### 7. NewOrderField() → Asc() / Desc()

**Before (v1.x)**:
```go
qs.OrderBy(orm.NewOrderField("created_at", true))
qs.OrderBy(orm.NewOrderField("price", false))
```

**After (v2.0)**:
```go
qs.OrderBy(orm.Asc("created_at"))
qs.OrderBy(orm.Desc("price"))
```

**Migration**: 
- `NewOrderField(field, true)` → `Asc(field)`
- `NewOrderField(field, false)` → `Desc(field)`

### 8. GetFieldAccessor() → FieldAccessor()

**Before (v1.x)**:
```go
fa, err := manager.GetFieldAccessor()
```

**After (v2.0)**:
```go
fa, err := manager.FieldAccessor()
```

**Migration**: Replace `.GetFieldAccessor()` with `.FieldAccessor()`

### 9. GetPaginationParams() → ParsePaginationParams()

**Before (v1.x)**:
```go
page, size, offset := api.GetPaginationParams(r, 20)
```

**After (v2.0)**:
```go
page, size, offset := api.ParsePaginationParams(r, 20)
```

**Migration**: Replace `GetPaginationParams(` with `ParsePaginationParams(`

## Automated Migration for v2.0

Use the provided migration script:

```bash
./scripts/migrate_naming_v2.sh
```

This will automatically update:
- `Rebind()` → `RebindPlaceholders()`
- `GetRegistry()` → `Global()`
- `GetFieldAccessor()` → `FieldAccessor()`
- `GetPaginationParams()` → `ParsePaginationParams()`

**Note**: `NewOrderField()` requires manual review (true/false → Asc/Desc)

See [MIGRATION_CHECKLIST_V2.md](./MIGRATION_CHECKLIST_V2.md) for detailed checklist.

## Summary Checklist

### v1.5.0 Changes
- [ ] Replace `orm.NewQ()` with direct expressions or `And()/Or()`
- [ ] Replace `orm.NewFieldQueryExpr()` with `Where()` or field methods
- [ ] Update `NewFieldExpr` to `NewField` (if used)
- [ ] Update method names: `Equals` → `Eq`, `Greater` → `Gt`, etc.

### v2.0 Changes
- [ ] Replace `Rebind()` with `RebindPlaceholders()`
- [ ] Replace `GetRegistry()` with `Global()`
- [ ] Replace `NewOrderField()` with `Asc()` or `Desc()`
- [ ] Replace `GetFieldAccessor()` with `FieldAccessor()`
- [ ] Replace `GetPaginationParams()` with `ParsePaginationParams()`

### Final Steps
- [ ] Run automated migration script
- [ ] Test all queries after migration
- [ ] Remove any deprecation warnings from logs
- [ ] Update documentation and examples

---

**Last Updated**: 2026-01-01  
**For**: Forge Framework v1.5.0 → v2.0.0 Migration
