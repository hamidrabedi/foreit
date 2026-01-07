# Naming Audit - Current Issues & Migration Plan

> **Critical Issue**: Framework naming is inconsistent and some names lack clear justification (like `NewQ`, `NewF`).

## Executive Summary

This document audits all current naming in the Forge framework and provides a migration plan to fix issues identified in [NAMING_ARCHITECTURE.md](./NAMING_ARCHITECTURE.md).

---

## Current Issues by Severity

### 🔴 CRITICAL - Breaking Changes Required

#### 1. Query Expression API - `NewQ()`

**Current**:
```go
q := orm.NewQ(nameField.Eq("John"))
qs.Filter(q)
```

**Issues**:
- ❌ `NewQ` is meaningless - Q is an abbreviation with no context
- ❌ Wrapping is redundant - expressions already work in Filter()
- ❌ Not idiomatic - Django uses `Q()` directly, not `NewQ()`
- ❌ Inconsistent - why `New` here but not elsewhere?

**Fixed**:
```go
// Option 1: Direct (most common)
qs.Filter(nameField.Eq("John"))

// Option 2: Q for complex queries
qs.Filter(Q(
    nameField.Eq("John"),
    ageField.Gt(18),
))

// Option 3: Boolean combiners
qs.Filter(And(
    nameField.Eq("John"),
    Or(
        ageField.Gt(18),
        statusField.Eq("active"),
    ),
))
```

**Impact**: High - used throughout codebase
**Files affected**: ~50+

#### 2. Field Reference API - `NewF()`

**Current**:
```go
// Not currently used but planned
f := orm.NewF("name")
```

**Issues**:
- ❌ `NewF` is even worse than `NewQ`
- ❌ F means nothing without Django context
- ❌ Should match Django's F() for field references

**Fixed**:
```go
// F creates a field reference (Django-style)
f := F("name")

// Or use type-safe fields
User.Name.Eq("John")
```

**Impact**: Medium - not widely used yet
**Files affected**: ~10

#### 3. `NewFieldQueryExpr()` - Too Verbose

**Current**:
```go
expr := orm.NewFieldQueryExpr("price", orm.OpGreater, 100.0)
```

**Issues**:
- ❌ Name is way too long
- ❌ Not user-friendly for a common operation
- ❌ Should be replaced with field methods or `Where()`

**Fixed**:
```go
// Option 1: Type-safe
expr := Product.Price.Gt(100.0)

// Option 2: Where clause
expr := Where("price", OpGreater, 100.0)

// Option 3: Keep NewFieldQueryExpr as internal helper only
```

**Impact**: High - used in many places
**Files affected**: ~30

---

### 🟡 MEDIUM - Deprecations & Improvements

#### 4. Inconsistent Constructor Naming

**Current Issues**:
```go
// Some use New prefix appropriately
NewQuerySet[T]()    // ✅ Good - returns concrete struct
NewManager[T]()     // ✅ Good - returns concrete struct
NewSQLBuilder()     // ✅ Good - returns concrete struct

// Others use New unnecessarily
NewBaseFilter[T]()  // ⚠️ Could be just a struct literal
NewFilterError()    // ⚠️ Could be inline

// Inconsistent patterns
NewFieldAccessor[T]()  // Returns *FieldAccessor
NewAdminManager[T]()   // Returns *AdminManager
```

**Fixed Pattern**:
```go
// Use "New" for:
// 1. Complex initialization
// 2. Error handling needed
// 3. Interface returns

func NewQuerySet[T any](table string) (QuerySet[T], error)
func NewManager[T any](table string, db *DB) *Manager[T]

// Don't use "New" for:
// 1. Simple structs (use literals)
// 2. Factory functions (use direct names)

// Instead of NewFilterError:
return &FilterError{Field: field, Message: msg}

// Instead of NewBaseFilter:
filter := &BaseFilter[T]{fieldPath: path}
```

**Impact**: Medium
**Files affected**: ~40

#### 5. Type Name Redundancy

**Current**:
```go
// In orm package
type QueryExpr struct {}      // ⚠️ Vague name
type FieldQueryExpr struct {} // ❌ Never defined, just function name
```

**Issues**:
- `QueryExpr` should be renamed to clarify it's an expression
- Confusion between `QueryExpr`, `Expression`, and `FieldExpression`

**Fixed**:
```go
// Clear hierarchy
type Expression interface {
    ToSQL(*SQLBuilder) (string, []interface{}, error)
}

type FieldExpression[T any] struct {
    path  string
    table string
}

type BoolExpression struct {
    operator Combiner
    children []Expression
}

type ComparisonExpression struct {
    field    string
    operator Operator
    value    interface{}
}
```

**Impact**: Medium
**Files affected**: ~25

#### 6. Unclear Variable Names in Generated Code

**Current**:
```go
// Generated field accessors
q := orm.NewQ(...)           // ❌ What is q?
f := orm.NewField[T](...)    // ⚠️ Too generic
```

**Fixed**:
```go
// Better names in examples and generated code
query := Q(...)
condition := Where(...)
nameField := NewField[string]("name", "users")
```

---

### 🟢 LOW - Polish & Consistency

#### 7. Receiver Name Consistency

**Current Issues**:
```go
// Inconsistent receivers across types
func (qs *BaseQuerySet) Filter() {}
func (q *BaseQuerySet) Exclude() {}   // ❌ Different receiver
```

**Fixed**:
```go
// Always use same receiver for a type
func (qs *BaseQuerySet) Filter() {}
func (qs *BaseQuerySet) Exclude() {}  // ✅ Consistent
```

**Impact**: Low - internal only
**Files affected**: All Go files (lint check)

#### 8. Package Import Aliases

**Current**:
```go
import (
    adminfilter "github.com/forgego/forge/admin/filter"
    adminorm "github.com/forgego/forge/admin/orm"
    adminschema "github.com/forgego/forge/admin/schema"
)
```

**Issues**:
- ⚠️ Necessary but verbose
- Could indicate packages should be organized differently

**Possible Future Fix**:
```go
// Reorganize packages to avoid conflicts
forge/
├── orm/          # Core ORM
├── filter/       # Core filtering
└── admin/
    ├── ui/       # UI components (no conflict with orm.UI)
    └── config/   # Admin config (uses core orm/filter)
```

**Impact**: Low - works but could be cleaner
**Decision**: Keep for now, revisit in v2.0

---

## File-by-File Audit

### forge/orm/

| File | Issues | Priority | Status |
|------|--------|----------|--------|
| `query_expr.go` | `NewQ()`, `NewFieldQueryExpr()` | 🔴 Critical | Pending |
| `expression.go` | `NewField()` OK, need `F()` shorthand | 🟡 Medium | Pending |
| `queryset.go` | `NewQuerySet()` OK | ✅ Good | - |
| `manager.go` | `NewManager()` OK | ✅ Good | - |
| `field_expr.go` | Needs review | 🟡 Medium | Pending |
| `annotations.go` | `NewAnnotation()` OK | ✅ Good | - |

### forge/filter/

| File | Issues | Priority | Status |
|------|--------|----------|--------|
| `filter.go` | `NewFilterError()` unnecessary | 🟡 Medium | Pending |
| `filterset.go` | `NewFilterSet()` OK | ✅ Good | - |
| `filters/*.go` | Various `New*Filter()` - review | 🟡 Medium | Pending |

### forge/admin/

| File | Issues | Priority | Status |
|------|--------|----------|--------|
| `admin.go` | `Register()` good, import aliases verbose | 🟢 Low | OK |
| `config.go` | `NewFieldset()` OK | ✅ Good | - |
| `orm/manager.go` | Uses `NewQ()` - will break | 🔴 Critical | Pending |
| `views/*.go` | Uses `NewQ()` heavily | 🔴 Critical | Pending |

### forge/api/

| File | Issues | Priority | Status |
|------|--------|----------|--------|
| `viewset.go` | Good naming | ✅ Good | - |
| `serializer.go` | Good naming | ✅ Good | - |
| `serializers/fields/*.go` | `NewBaseField()` OK | ✅ Good | - |

### forge/schema/

| File | Issues | Priority | Status |
|------|--------|----------|--------|
| `field.go` | Good naming | ✅ Good | - |
| `registry.go` | `NewFieldBuilder()` OK | ✅ Good | - |
| `field_config.go` | `NewFieldOptions()` OK | ✅ Good | - |

---

## Migration Plan

### Phase 1: Preparation (Week 1)

**Goals**:
- Document current usage
- Create compatibility layer
- Write migration guide

**Tasks**:
1. ✅ Create NAMING_ARCHITECTURE.md (Done)
2. ✅ Create NAMING_AUDIT.md (Done)
3. ⬜ Grep all usages of `NewQ`, `NewF`, `NewFieldQueryExpr`
4. ⬜ Create compatibility shims
5. ⬜ Update examples with new API

**Deliverables**:
- Migration guide for users
- Compatibility layer code
- Updated examples

### Phase 2: New API Implementation (Week 2)

**Goals**:
- Implement new Q() and F() functions
- Implement Where() helper
- Keep old API with deprecation warnings

**Tasks**:
1. ⬜ Implement `Q()` factory function
2. ⬜ Implement `F()` field reference function
3. ⬜ Implement `Where()` condition function
4. ⬜ Implement `And()`, `Or()`, `Not()` combiners
5. ⬜ Add deprecation warnings to `NewQ()`, `NewFieldQueryExpr()`
6. ⬜ Write comprehensive tests

**Code Changes**:
```go
// forge/orm/query.go

// Q creates a query expression from conditions
// This replaces NewQ() and provides a cleaner API
func Q(conditions ...Condition) Expression {
    if len(conditions) == 0 {
        return &EmptyExpression{}
    }
    if len(conditions) == 1 {
        return conditions[0]
    }
    return And(conditions...)
}

// F creates a field reference for use in queries
// Matches Django's F() for field references
func F(fieldPath string) FieldRef {
    return FieldRef{path: fieldPath}
}

// Where creates a simple field condition
// Alternative to field expression methods
func Where(field string, op Operator, value interface{}) Condition {
    return &SimpleCondition{
        field: field,
        op:    op,
        value: value,
    }
}

// Deprecated: Use Q() instead. NewQ will be removed in v2.0.
func NewQ(expr Expression) *QObject {
    log.Warn("NewQ is deprecated, use Q() instead")
    return &QObject{expr: expr}
}

// Deprecated: Use Where() or field methods instead. Removed in v2.0.
func NewFieldQueryExpr(field string, op Operator, value interface{}) QueryExpr {
    log.Warn("NewFieldQueryExpr is deprecated, use Where() or field methods instead")
    return QueryExpr{field: field, op: op, value: value}
}
```

### Phase 3: Internal Migration (Week 3-4)

**Goals**:
- Update all internal code to use new API
- Fix all deprecation warnings
- Update tests

**Tasks**:
1. ⬜ Update forge/admin to use new API
2. ⬜ Update forge/filter to use new API  
3. ⬜ Update all examples
4. ⬜ Update all tests
5. ⬜ Run full test suite
6. ⬜ Fix any broken functionality

**Priority Order**:
1. Core ORM (forge/orm/)
2. Admin (forge/admin/) - highest user impact
3. Filter (forge/filter/)
4. Examples (examples/)
5. Tests (tests/)

### Phase 4: Documentation Update (Week 5)

**Goals**:
- Update all documentation
- Create migration guide for users
- Update tutorials and examples

**Tasks**:
1. ⬜ Update API_REFERENCE.md
2. ⬜ Update GETTING_STARTED.md
3. ⬜ Update USAGE_GUIDE.md
4. ⬜ Create MIGRATION_V1_TO_V2.md
5. ⬜ Update docs-site/
6. ⬜ Update README examples

### Phase 5: Announcement & Deprecation (Week 6)

**Goals**:
- Announce changes to users
- Provide clear migration path
- Set timeline for removal

**Tasks**:
1. ⬜ Write blog post / announcement
2. ⬜ Update CHANGELOG.md
3. ⬜ Tag v1.x with deprecation warnings
4. ⬜ Set timeline for v2.0 (old API removal)

**Timeline**:
- v1.5.0: New API available, old API deprecated
- v1.6.0-1.9.0: Migration period (6 months)
- v2.0.0: Old API removed

### Phase 6: Removal (v2.0)

**Goals**:
- Remove deprecated APIs
- Clean up compatibility layer
- Final naming polish

**Tasks**:
1. ⬜ Remove `NewQ()` function
2. ⬜ Remove `NewFieldQueryExpr()` function
3. ⬜ Remove compatibility shims
4. ⬜ Clean up any remaining naming issues
5. ⬜ Final documentation pass

---

## Code Search Results

### Usages of `NewQ()`

```bash
$ grep -r "NewQ(" forge/ --include="*.go" | wc -l
15
```

**Files affected**:
- `forge/admin/http/autocomplete.go` (2 usages)
- `forge/admin/orm/queryset.go` (2 usages)  
- `forge/admin/views/list_view.go` (2 usages)
- `forge/orm/expression.go` (1 definition + 4 usages in tests)
- `forge/orm/integration_test.go` (4 usages)

**Estimated effort**: 2-3 hours to refactor

### Usages of `NewFieldQueryExpr()`

```bash
$ grep -r "NewFieldQueryExpr(" forge/ --include="*.go" | wc -l
25
```

**Files affected**:
- `forge/orm/query_expr.go` (1 definition)
- `forge/orm/sql_builder_test.go` (5 usages)
- `forge/orm/queryset_test.go` (2 usages)
- `forge/orm/integration_test.go` (1 usage)
- Various test files

**Estimated effort**: 3-4 hours to refactor

### Usages of `NewField()`

```bash
$ grep -r "NewField\[" forge/ --include="*.go" | wc -l
80+
```

**Status**: ✅ This one is GOOD - it's an actual constructor
**Action**: Keep as-is, optionally add `F()` shorthand

---

## Compatibility Layer

To ease migration, we'll provide a compatibility layer:

```go
// forge/orm/compat.go

// Deprecated compatibility layer for v1.x -> v2.0 migration

// NewQ creates a query expression (deprecated)
// Use Q() instead
// 
// Example:
//   // Old
//   q := orm.NewQ(expr)
//   // New  
//   q := orm.Q(expr)
//
// This function will be removed in v2.0
func NewQ(expr Expression) Expression {
    if internal.WarnDeprecated {
        log.Warn("NewQ is deprecated, use Q() instead")
    }
    return Q(expr)
}

// NewFieldQueryExpr creates a field condition (deprecated)
// Use Where() or field expression methods instead
//
// Example:
//   // Old
//   expr := orm.NewFieldQueryExpr("name", orm.OpEquals, "John")
//   // New
//   expr := orm.Where("name", orm.OpEquals, "John")
//   // Or even better (type-safe)
//   expr := User.Name.Eq("John")
//
// This function will be removed in v2.0
func NewFieldQueryExpr(field string, op Operator, value interface{}) Condition {
    if internal.WarnDeprecated {
        log.Warn("NewFieldQueryExpr is deprecated, use Where() or field methods instead")
    }
    return Where(field, op, value)
}
```

---

## Testing Strategy

### 1. Compatibility Tests

Ensure old and new APIs produce identical results:

```go
func TestNewAPICompatibility(t *testing.T) {
    // Test that new API produces same results as old
    
    // Old API
    oldExpr := NewFieldQueryExpr("name", OpEquals, "John")
    oldSQL, oldArgs, _ := oldExpr.ToSQL(builder)
    
    // New API
    newExpr := Where("name", OpEquals, "John")
    newSQL, newArgs, _ := newExpr.ToSQL(builder)
    
    assert.Equal(t, oldSQL, newSQL)
    assert.Equal(t, oldArgs, newArgs)
}
```

### 2. Migration Tests

Test that common migration patterns work:

```go
func TestMigrationPatterns(t *testing.T) {
    // Pattern 1: Simple filter
    // Old: qs.Filter(NewQ(field.Eq(value)))
    // New: qs.Filter(field.Eq(value))
    
    // Pattern 2: Complex boolean
    // Old: NewQ(expr1).Or(NewQ(expr2))
    // New: Or(expr1, expr2)
    
    // Pattern 3: Field conditions
    // Old: NewFieldQueryExpr("name", OpEquals, "John")
    // New: Where("name", OpEquals, "John")
}
```

### 3. Deprecation Warning Tests

Ensure warnings are shown:

```go
func TestDeprecationWarnings(t *testing.T) {
    // Enable deprecation warnings
    internal.WarnDeprecated = true
    
    // Capture log output
    var buf bytes.Buffer
    log.SetOutput(&buf)
    
    // Use deprecated function
    _ = NewQ(someExpr)
    
    // Check warning was logged
    assert.Contains(t, buf.String(), "deprecated")
}
```

---

## User Migration Examples

### Example 1: Simple Filter

**Before**:
```go
q := orm.NewQ(nameField.Eq("John"))
users, err := userManager.Filter(q).All(ctx)
```

**After**:
```go
// Option 1: Direct (simplest)
users, err := userManager.Filter(nameField.Eq("John")).All(ctx)

// Option 2: Q() for clarity (if you prefer)
users, err := userManager.Filter(Q(nameField.Eq("John"))).All(ctx)
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
users, err := userManager.Filter(
    Where("email", OpContains, "@example.com"),
).All(ctx)
```

---

## Success Metrics

### Code Quality
- [ ] All deprecation warnings resolved
- [ ] Test coverage maintained or improved
- [ ] No regressions in functionality
- [ ] Documentation coverage at 100%

### Developer Experience
- [ ] Examples are clearer and more concise
- [ ] API is more intuitive (user feedback)
- [ ] Fewer questions about "What is Q?"
- [ ] Migration guide is comprehensive

### Performance
- [ ] No performance regression
- [ ] Benchmarks show comparable or better performance
- [ ] Memory usage unchanged or improved

---

## Risk Assessment

### High Risk
- **Breaking changes**: Managed with deprecation period
- **Widespread usage**: Mitigated with compatibility layer
- **User disruption**: Minimized with clear migration guide

### Medium Risk
- **Documentation gaps**: Addressed in Phase 4
- **Example code outdated**: Fixed in Phase 3
- **Third-party code**: Cannot control, but deprecation warnings help

### Low Risk
- **Internal refactoring**: Well-tested
- **New API addition**: Backwards compatible
- **Performance**: No algorithmic changes

---

## Decision Log

### Decision 1: Deprecation Period Length
**Decision**: 6 months (v1.5.0 to v2.0.0)
**Rationale**: Gives users time to migrate without rushing
**Date**: 2026-01-01

### Decision 2: Keep NewField()
**Decision**: Keep `NewField[T]()` as-is, add `F()` as optional shorthand
**Rationale**: NewField is a proper constructor, name is justified
**Date**: 2026-01-01

### Decision 3: Q() Returns Expression, Not Wrapper
**Decision**: `Q()` returns `Expression` interface, not `*Q` struct
**Rationale**: Simpler, more flexible, matches actual usage
**Date**: 2026-01-01

---

## Next Steps

1. **Immediate** (Today):
   - ✅ Create NAMING_ARCHITECTURE.md
   - ✅ Create NAMING_AUDIT.md
   - ⬜ Review with team

2. **This Week**:
   - ⬜ Implement new Q(), F(), Where() functions
   - ⬜ Add deprecation warnings
   - ⬜ Create compatibility tests

3. **Next Week**:
   - ⬜ Start internal migration
   - ⬜ Update examples
   - ⬜ Begin documentation updates

4. **This Month**:
   - ⬜ Complete internal migration
   - ⬜ Release v1.5.0 with new API
   - ⬜ Announce changes

---

**Last Updated**: 2026-01-01
**Status**: Draft - Ready for Review
**Reviewers**: @team
