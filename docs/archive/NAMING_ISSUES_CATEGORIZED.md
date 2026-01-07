# Naming Issues - Categorized Summary

**Date**: 2026-01-01  
**Source**: NAMING_AUDIT_V2.md

## Quick Reference

This document provides a categorized view of all naming issues for quick decision-making.

---

## By Severity

### 🔴 High Priority (5 issues)

Must be addressed in v2.0 with deprecation paths.

| # | Issue | Location | Fix Complexity |
|---|-------|----------|----------------|
| 1 | `Rebind()` unclear name | `forge/db/db.go:71` | Easy |
| 2 | `GetRegistry()` redundant | `forge/registry/registry.go:78` | Easy |
| 3 | `NewOrderField()` unnecessary | `forge/orm/queryset.go:87` | Easy |
| 4 | `GetFieldAccessor()` violates Go convention | `forge/orm/manager.go:42` | Easy |
| 5 | `GetPaginationParams()` violates Go convention | `forge/api/pagination.go:31` | Easy |

**Estimated Effort**: 2-3 hours  
**Breaking Changes**: Yes (all need deprecation)

---

### 🟡 Medium Priority (10 issues)

Should be addressed in v2.x releases.

| # | Issue | Location | Fix Complexity |
|---|-------|----------|----------------|
| 6 | `BuildInsertSQL()` inconsistent | `forge/orm/sql_builder.go` | Medium |
| 7 | `ExecuteInsert()` should be methods | `forge/orm/manager_helpers.go` | Medium |
| 8 | `NewBaseFilter()` unnecessary base | `forge/filter/filter.go:53` | Easy |
| 9 | `NewAnnotation()` inconsistent | `forge/orm/annotations.go:12` | Easy |
| 10 | `DefaultSettings()` verbose | `forge/api/api.go:50` | Easy |
| 11 | `GetFieldPath()` getter overuse | `forge/filter/filter.go` | Medium |
| 12 | `SetLabel()` setter pattern | `forge/filter/filter.go:71` | Easy |
| 13 | `NewTextInput()` verbose | `forge/admin/widgets.go` | Easy |
| 14 | `BuildPaginatedResponse()` inconsistent | `forge/api/pagination.go:82` | Easy |
| 15 | `RegisterAggregate()` incomplete | `forge/orm/aggregates.go:87` | Easy |

**Estimated Effort**: 8-10 hours  
**Breaking Changes**: Yes (most need deprecation)

---

### 🟢 Low Priority (8 issues)

Can be addressed in v3.0 or later.

| # | Issue | Location | Fix Complexity |
|---|-------|----------|----------------|
| 16 | `OrderField` vs `Ordering[T]` | Multiple | Medium |
| 17 | `Choice[V any]` over-engineered | `forge/admin/widgets.go:11` | Medium |
| 18 | `AggregateFunc` redundant | `forge/orm/aggregates.go:11` | Easy |
| 19 | `ValuesQuerySet[T]` verbose | `forge/orm/queryset.go:42` | Easy |
| 20 | `UpdateMap` too generic | `forge/orm/update_builder.go` | Easy |
| 21 | `ModelWithID` naming | `forge/orm/manager.go:146` | Easy |
| 22 | `Base` prefix pattern | Multiple | Hard |
| 23 | `rebindPostgresToSQLite()` specific | `forge/db/db.go:85` | None |

**Estimated Effort**: 5-6 hours  
**Breaking Changes**: Yes (major refactoring)

---

### 📝 Documentation Only (3 issues)

No code changes needed.

| # | Issue | Location | Fix Complexity |
|---|-------|----------|----------------|
| 24 | `Aggregate` struct docs | `forge/orm/aggregates.go:4` | Easy |
| 25 | `Expression` interface examples | `forge/orm/expression.go` | Easy |
| 26 | `Filter[T]` implementation guide | `forge/filter/filter.go:11` | Medium |

**Estimated Effort**: 2-3 hours  
**Breaking Changes**: None

---

## By Package

### forge/orm/ (7 issues)
- 🔴 High: #3, #4
- 🟡 Medium: #6, #7, #9, #15
- 🟢 Low: #16, #18, #19, #20, #21, #22
- 📝 Docs: #24, #25

### forge/admin/ (4 issues)
- 🟡 Medium: #13
- 🟢 Low: #16, #17, #22

### forge/api/ (3 issues)
- 🔴 High: #5
- 🟡 Medium: #10, #14

### forge/filter/ (4 issues)
- 🟡 Medium: #8, #11, #12
- 📝 Docs: #26

### forge/db/ (1 issue)
- 🔴 High: #1

### forge/registry/ (1 issue)
- 🔴 High: #2

### forge/server/ (2 issues)
- 🟡 Medium: #10

---

## By Migration Impact

### Breaking Changes Required (18 issues)
Issues that need deprecation paths:
- High: #1, #2, #3, #4, #5
- Medium: #6, #7, #8, #9, #10, #11, #12, #13, #14
- Low: #16, #17, #19, #20, #21

### Non-Breaking (5 issues)
- Medium: #15 (already broken/incomplete)
- Low: #18, #22, #23
- Docs: All documentation issues

---

## By Fix Complexity

### Easy (16 issues)
Can be fixed with simple rename + deprecation:
- #1, #2, #3, #4, #5, #8, #9, #10, #12, #13, #14, #15, #18, #19, #20, #21

### Medium (7 issues)
Require some refactoring:
- #6, #7, #11, #16, #17, #26

### Hard (1 issue)
Major refactoring required:
- #22 (Base prefix pattern across multiple packages)

### None (2 issues)
Documentation only:
- #23, #24, #25

---

## Recommended Action Plan

### Immediate (Next Sprint)
Focus on high-priority, easy fixes:
1. ✅ Fix #1: `Rebind()` → `RebindPlaceholders()`
2. ✅ Fix #2: `GetRegistry()` → `Global()`
3. ✅ Fix #3: Deprecate `NewOrderField()`
4. ✅ Fix #4: `GetFieldAccessor()` → `FieldAccessor()`
5. ✅ Fix #5: `GetPaginationParams()` → `ParsePaginationParams()`

**Deliverable**: v2.0-alpha with 5 deprecations

---

### Short Term (v2.0)
Address medium-priority easy fixes:
1. ✅ Fix #8: Make `NewBaseFilter()` internal
2. ✅ Fix #9: `NewAnnotation()` → `Annotate()`
3. ✅ Fix #10: `DefaultSettings()` → `NewSettings()`
4. ✅ Fix #12: `SetLabel()` → `WithLabel()`
5. ✅ Fix #13: `NewTextInput()` → `TextInput()`
6. ✅ Fix #14: `BuildPaginatedResponse()` → `NewPaginatedResponse()`
7. ✅ Fix #15: Document or remove `RegisterAggregate()`

**Deliverable**: v2.0 stable with 12 deprecations

---

### Medium Term (v2.1-v2.5)
Address medium-priority complex fixes:
1. ✅ Fix #6: Refactor `BuildInsertSQL()` to methods
2. ✅ Fix #7: Create `Executor` type
3. ✅ Fix #11: Remove `Get` prefix from interface methods

**Deliverable**: v2.x releases with improved APIs

---

### Long Term (v3.0)
Major refactoring:
1. ✅ Fix #22: Reconsider `Base` prefix pattern
2. ✅ Fix #16: Standardize `OrderField` vs `Ordering`
3. ✅ Remove all deprecated functions

**Deliverable**: v3.0 with clean, consistent API

---

## Impact Analysis

### User-Facing Impact

**High Priority Issues**: Affect most users
- `GetFieldAccessor()` - Used in custom field logic
- `GetPaginationParams()` - Used in all paginated APIs
- `Rebind()` - Used in multi-database apps

**Medium Priority Issues**: Affect some users
- Widget constructors - Used in custom admin
- Filter builders - Used in advanced filtering
- SQL builders - Used in custom queries

**Low Priority Issues**: Affect few users
- Internal types and helpers
- Advanced customization

### Internal Impact

**Code Changes Required**:
- High: ~50 files
- Medium: ~30 files
- Low: ~20 files

**Test Updates Required**:
- High: ~20 test files
- Medium: ~15 test files
- Low: ~10 test files

**Documentation Updates Required**:
- All priority levels: ~15 doc files

---

## Success Metrics

### v2.0 Goals
- ✅ All high-priority issues addressed
- ✅ 50% of medium-priority issues addressed
- ✅ Zero new naming issues introduced
- ✅ All deprecations documented

### v2.x Goals
- ✅ 100% of medium-priority issues addressed
- ✅ 50% of low-priority issues addressed
- ✅ Migration guide complete

### v3.0 Goals
- ✅ All deprecated functions removed
- ✅ 100% naming consistency
- ✅ Zero naming-related issues in backlog

---

## Conclusion

The naming audit identified 26 issues across the framework. Most are easy to fix with proper deprecation paths. The recommended approach is:

1. **v2.0**: Fix 5 high-priority issues (2-3 hours)
2. **v2.x**: Fix 10 medium-priority issues (8-10 hours)
3. **v3.0**: Fix 8 low-priority issues (5-6 hours)

**Total Estimated Effort**: 15-19 hours of development work

**User Impact**: Minimal with proper deprecation periods (6 months recommended)

**Benefits**: 
- Clearer, more intuitive API
- Better Go idiom compliance
- Easier onboarding for new developers
- Reduced cognitive load
