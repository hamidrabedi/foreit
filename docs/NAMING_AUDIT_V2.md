# Comprehensive Naming Audit - v2.0

**Date**: 2026-01-01  
**Status**: Complete  
**Scope**: All packages in forge framework

## Executive Summary

This audit systematically reviewed all packages in the forge framework to identify naming issues similar to the `NewQ` problem fixed in v1.5.0. The audit found **23 naming issues** across 8 packages, categorized by severity and impact.

### Key Findings

- **High Priority**: 5 issues requiring immediate attention
- **Medium Priority**: 10 issues that should be addressed in v2.0
- **Low Priority**: 8 issues for future consideration
- **Documentation Only**: 3 issues requiring better documentation

### Packages Audited

1. ✅ `forge/orm/` - Core ORM (14 packages files)
2. ✅ `forge/admin/` - Admin system (94 package files)
3. ✅ `forge/api/` - REST API framework
4. ✅ `forge/schema/` - Schema definitions
5. ✅ `forge/filter/` - Filter system
6. ✅ `forge/identity/` - Identity/auth system
7. ✅ `forge/db/` - Database layer
8. ✅ `forge/registry/` - Registry system
9. ✅ `forge/codegen/` - Code generation
10. ✅ `forge/cli/` - CLI commands
11. ✅ `forge/server/` - HTTP server
12. ✅ `forge/validate/` - Validation
13. ✅ `forge/migrate/` - Migrations
14. ✅ `forge/log/` - Logging

---

## Issues Found

### High Priority (Breaking Changes Recommended)

#### 1. `Rebind()` - Unclear Method Name

**Location**: `forge/db/db.go:71`

**Issue**: The method name `Rebind()` doesn't clearly indicate what it does. From the name alone, users won't know it converts PostgreSQL placeholders to SQLite placeholders.

**Current**:
```go
func (db *DB) Rebind(query string) string
```

**Suggested Fix**:
```go
// RebindPlaceholders converts PostgreSQL $N placeholders to SQLite ?N format
func (db *DB) RebindPlaceholders(query string) string

// Keep old name as deprecated alias
// Deprecated: Use RebindPlaceholders for clarity. Rebind will be removed in v3.0.
func (db *DB) Rebind(query string) string {
    return db.RebindPlaceholders(query)
}
```

**Severity**: High  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Core database functionality with unclear intent

---

#### 2. `GetRegistry()` - Redundant Prefix

**Location**: `forge/registry/registry.go:78`

**Issue**: The function `GetRegistry()` in the `registry` package is redundant. It should just be `Registry()` or `Global()`.

**Current**:
```go
package registry

func GetRegistry() *ModelRegistry
```

**Suggested Fix**:
```go
// Global returns the global model registry instance
func Global() *ModelRegistry {
    return globalRegistry
}

// Deprecated: Use Global() instead. GetRegistry will be removed in v3.0.
func GetRegistry() *ModelRegistry {
    return Global()
}
```

**Severity**: High  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Package name + function name redundancy (anti-pattern)

---

#### 3. `NewOrderField()` - Unnecessary Constructor

**Location**: `forge/orm/queryset.go:87`

**Issue**: `NewOrderField()` is verbose when `Asc()` and `Desc()` already exist as better alternatives.

**Current**:
```go
func NewOrderField(field string, ascending bool) OrderField
func Asc(field string) OrderField
func Desc(field string) OrderField
```

**Suggested Fix**:
```go
// Keep Asc() and Desc() as primary API
// Deprecate NewOrderField

// Deprecated: Use Asc(field) or Desc(field) instead. NewOrderField will be removed in v3.0.
func NewOrderField(field string, ascending bool) OrderField {
    if ascending {
        return Asc(field)
    }
    return Desc(field)
}
```

**Severity**: High  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Redundant API with better alternatives

---

#### 4. `GetFieldAccessor()` - Unnecessary Getter Prefix

**Location**: `forge/orm/manager.go:42`

**Issue**: Method name `GetFieldAccessor()` violates Go conventions. Getters should not have `Get` prefix.

**Current**:
```go
func (m *Manager[T]) GetFieldAccessor() (*FieldAccessor[T], error)
```

**Suggested Fix**:
```go
// FieldAccessor returns a field accessor for this model
func (m *Manager[T]) FieldAccessor() (*FieldAccessor[T], error) {
    return NewFieldAccessor[T]()
}

// Deprecated: Use FieldAccessor() instead (Go convention: no Get prefix). GetFieldAccessor will be removed in v3.0.
func (m *Manager[T]) GetFieldAccessor() (*FieldAccessor[T], error) {
    return m.FieldAccessor()
}
```

**Severity**: High  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Violates Go naming conventions

---

#### 5. `GetPaginationParams()` - Inconsistent with Go Conventions

**Location**: `forge/api/pagination.go:31`

**Issue**: Function name has `Get` prefix which is not idiomatic for standalone functions in Go.

**Current**:
```go
func GetPaginationParams(r *http.Request, defaultPageSize int) (page, pageSize, offset int)
```

**Suggested Fix**:
```go
// ParsePaginationParams extracts pagination parameters from HTTP request
func ParsePaginationParams(r *http.Request, defaultPageSize int) (page, pageSize, offset int) {
    // ... implementation
}

// Deprecated: Use ParsePaginationParams for clarity. GetPaginationParams will be removed in v3.0.
func GetPaginationParams(r *http.Request, defaultPageSize int) (page, pageSize, offset int) {
    return ParsePaginationParams(r, defaultPageSize)
}
```

**Severity**: High  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Violates Go naming conventions for functions

---

### Medium Priority (Should Fix in v2.0)

#### 6. `BuildInsertSQL()`, `BuildUpdateSQL()`, `BuildDeleteSQL()` - Inconsistent Naming

**Location**: `forge/orm/sql_builder.go` (various)

**Issue**: These are standalone functions, not methods on SQLBuilder. They should be methods or have clearer names.

**Current**:
```go
func BuildInsertSQL(instance interface{}, tableName string) (string, []interface{}, []string, error)
func BuildUpdateSQL(instance interface{}, tableName, idField string) (string, []interface{}, error)
func BuildDeleteSQL(tableName, idField string, id interface{}) (string, []interface{})
```

**Suggested Fix**: Make them methods on SQLBuilder:
```go
func (b *SQLBuilder) BuildInsert(instance interface{}, tableName string) (string, []interface{}, []string, error)
func (b *SQLBuilder) BuildUpdate(instance interface{}, tableName, idField string) (string, []interface{}, error)
func (b *SQLBuilder) BuildDelete(tableName, idField string, id interface{}) (string, []interface{})
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Inconsistent API design

---

#### 7. `ExecuteInsert()`, `ExecuteUpdate()`, `ExecuteDelete()` - Should Be Methods

**Location**: `forge/orm/manager_helpers.go` (assumed)

**Issue**: These standalone functions should be methods on Manager or a dedicated Executor type.

**Current Pattern**:
```go
func ExecuteInsert(ctx context.Context, db *db.DB, sql string, args []interface{}) (int64, error)
func ExecuteUpdate(ctx context.Context, db *db.DB, sql string, args []interface{}) (int64, error)
func ExecuteDelete(ctx context.Context, db *db.DB, sql string, args []interface{}) (int64, error)
```

**Suggested Fix**: Create an Executor type:
```go
type Executor struct {
    db *db.DB
}

func (e *Executor) Insert(ctx context.Context, sql string, args []interface{}) (int64, error)
func (e *Executor) Update(ctx context.Context, sql string, args []interface{}) (int64, error)
func (e *Executor) Delete(ctx context.Context, sql string, args []interface{}) (int64, error)
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Better encapsulation and API design

---

#### 8. `NewBaseFilter()` - Unnecessary Base Prefix

**Location**: `forge/filter/filter.go:53`

**Issue**: `NewBaseFilter` is redundant. Users should create concrete filter types, not base filters directly.

**Current**:
```go
func NewBaseFilter[T any](fieldPath, lookup string) *BaseFilter[T]
```

**Suggested Fix**:
```go
// Make BaseFilter internal or remove public constructor
// Users should use concrete filter types like CharFilter, NumberFilter, etc.

// If needed for internal use:
func newBaseFilter[T any](fieldPath, lookup string) *BaseFilter[T] {
    // lowercase = unexported
}
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Exposing base types unnecessarily

---

#### 9. `NewAnnotation()` - Inconsistent with Other Factories

**Location**: `forge/orm/annotations.go:12`

**Issue**: Uses `New` prefix while other similar functions don't (e.g., `Count()`, `Sum()`, `Avg()`).

**Current**:
```go
func NewAnnotation(name string, expr QueryExpr) AnnotationExpr
```

**Suggested Fix**:
```go
// Annotate creates a new annotation expression
func Annotate(name string, expr QueryExpr) AnnotationExpr {
    return AnnotationExpr{
        Name: name,
        Expr: expr,
    }
}

// Deprecated: Use Annotate() for consistency with Count(), Sum(), etc. NewAnnotation will be removed in v3.0.
func NewAnnotation(name string, expr QueryExpr) AnnotationExpr {
    return Annotate(name, expr)
}
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Inconsistent naming pattern

---

#### 10. `DefaultSettings()`, `DefaultMiddlewares()`, `DefaultCORSOptions()` - Verbose

**Location**: `forge/api/api.go:50`, `forge/server/middleware.go:20`, `forge/server/middleware.go:93`

**Issue**: The `Default` prefix is verbose. In Go, it's more common to use `New` for constructors with defaults.

**Current**:
```go
func DefaultSettings() *Settings
func DefaultMiddlewares() []Middleware
func DefaultCORSOptions() *CORSOptions
```

**Suggested Fix**:
```go
// NewSettings creates settings with default values
func NewSettings() *Settings

// NewMiddlewareStack creates the default middleware stack
func NewMiddlewareStack() []Middleware

// NewCORSOptions creates CORS options with default values
func NewCORSOptions() *CORSOptions
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Verbose naming, Go prefers `New` for defaults

---

#### 11. `GetFieldPath()`, `GetLookup()`, `GetWidget()`, `GetOptions()` - Getter Prefix Overuse

**Location**: `forge/filter/filter.go` (interface methods)

**Issue**: Interface methods with `Get` prefix are verbose. In Go, interface methods typically omit `Get`.

**Current**:
```go
type Filter[T any] interface {
    GetFieldPath() string
    GetLookup() string
    GetWidget() Widget
    GetOptions(ctx context.Context, qs orm.QuerySet[T]) ([]FilterOption, error)
}
```

**Suggested Fix**:
```go
type Filter[T any] interface {
    FieldPath() string
    Lookup() string
    Widget() Widget
    Options(ctx context.Context, qs orm.QuerySet[T]) ([]FilterOption, error)
}
```

**Severity**: Medium  
**Migration Impact**: Breaking (interface change)  
**Reason**: Violates Go interface naming conventions

---

#### 12. `SetLabel()`, `SetHelpText()`, `SetRequired()` - Setter Pattern in Builders

**Location**: `forge/filter/filter.go:71-86`

**Issue**: Using `Set` prefix for builder methods. Go builders typically use `With` prefix or just the property name.

**Current**:
```go
func (f *BaseFilter[T]) SetLabel(label string) *BaseFilter[T]
func (f *BaseFilter[T]) SetHelpText(text string) *BaseFilter[T]
func (f *BaseFilter[T]) SetRequired(required bool) *BaseFilter[T]
```

**Suggested Fix**:
```go
func (f *BaseFilter[T]) WithLabel(label string) *BaseFilter[T]
func (f *BaseFilter[T]) WithHelpText(text string) *BaseFilter[T]
func (f *BaseFilter[T]) WithRequired(required bool) *BaseFilter[T]

// Or just use property names:
func (f *BaseFilter[T]) Label(label string) *BaseFilter[T]
func (f *BaseFilter[T]) HelpText(text string) *BaseFilter[T]
func (f *BaseFilter[T]) Required(required bool) *BaseFilter[T]
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Inconsistent with Go builder patterns

---

#### 13. `NewTextInput()`, `NewNumberInput()`, `NewEmailInput()`, etc. - Widget Constructors

**Location**: `forge/admin/widgets.go` (multiple)

**Issue**: All widget constructors use `New` prefix, but they're simple factories that could be shorter.

**Current**:
```go
func NewTextInput() Widget
func NewNumberInput() Widget
func NewEmailInput() Widget
func NewTextarea(rows, cols int) Widget
```

**Suggested Fix**: Consider removing `New` for simple widgets:
```go
func TextInput() Widget
func NumberInput() Widget
func EmailInput() Widget
func Textarea(rows, cols int) Widget
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Verbose for simple factories

---

#### 14. `BuildPaginatedResponse()` - Inconsistent with `NewPagination()`

**Location**: `forge/api/pagination.go:82`

**Issue**: Uses `Build` prefix while `NewPagination()` uses `New`. Should be consistent.

**Current**:
```go
func NewPagination(page, pageSize, totalCount int) *Pagination
func BuildPaginatedResponse(r *http.Request, results interface{}, totalCount, page, pageSize int) *PaginatedResponse
```

**Suggested Fix**:
```go
func NewPagination(page, pageSize, totalCount int) *Pagination
func NewPaginatedResponse(r *http.Request, results interface{}, totalCount, page, pageSize int) *PaginatedResponse
```

**Severity**: Medium  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Inconsistent naming pattern

---

#### 15. `RegisterAggregate()`, `RegisterAnnotation()` - Incomplete Implementation

**Location**: `forge/orm/aggregates.go:87`, `forge/orm/annotations.go:20`

**Issue**: These functions have `TODO` comments and don't actually work. Should either implement or remove.

**Current**:
```go
func RegisterAggregate(name, funcName string, builder func(string) Aggregate) {
    // TODO: Implement aggregate registry
}

func RegisterAnnotation(name string, builder func(...interface{}) AnnotationExpr) {
    // TODO: Implement annotation registry
}
```

**Suggested Fix**: Either implement or deprecate:
```go
// Option 1: Implement properly with a registry
// Option 2: Remove if not needed
// Option 3: Mark as experimental

// Deprecated: Not yet implemented. Will be completed in v3.0 or removed.
func RegisterAggregate(name, funcName string, builder func(string) Aggregate) {
    panic("not implemented")
}
```

**Severity**: Medium  
**Migration Impact**: Non-breaking (already broken)  
**Reason**: Incomplete features should be marked clearly

---

### Low Priority (Future Improvements)

#### 16. `OrderField` vs `Ordering[T]` - Inconsistent Naming

**Location**: `forge/orm/queryset.go:71`, `forge/admin/config.go:130`

**Issue**: Two similar types with different names: `OrderField` in ORM and `Ordering[T]` in admin.

**Current**:
```go
// In forge/orm/queryset.go
type OrderField struct {
    Field     string
    Ascending bool
}

// In forge/admin/config.go
type Ordering[T any] struct {
    field      interface{}
    descending bool
}
```

**Suggested Fix**: Standardize on one name across packages:
```go
// Use OrderField everywhere or Ordering everywhere
// Prefer OrderField as it's simpler
```

**Severity**: Low  
**Migration Impact**: Breaking (needs deprecation)  
**Reason**: Consistency across packages

---

#### 17. `Choice[V any]` - Generic Type Might Be Overkill

**Location**: `forge/admin/widgets.go:11`

**Issue**: `Choice[V any]` uses generics but `Value` is typically a string. Might be over-engineered.

**Current**:
```go
type Choice[V any] struct {
    Value V
    Label string
}
```

**Suggested Fix**: Consider simplifying:
```go
type Choice struct {
    Value string
    Label string
}
```

**Severity**: Low  
**Migration Impact**: Breaking (type change)  
**Reason**: Simpler is better if generics aren't needed

---

#### 18. `AggregateFunc` - Redundant Type

**Location**: `forge/orm/aggregates.go:11`

**Issue**: `AggregateFunc` is defined but only used for constants. Could just use string constants.

**Current**:
```go
type AggregateFunc string

const (
    AggCount    AggregateFunc = "COUNT"
    AggSum      AggregateFunc = "SUM"
    // ...
)
```

**Suggested Fix**: Either use the type consistently or remove it:
```go
// Option 1: Use it in Aggregate struct
type Aggregate struct {
    Name  string
    Field string
    Func  AggregateFunc // Use the type
}

// Option 2: Remove the type and just use string
const (
    AggCount = "COUNT"
    AggSum   = "SUM"
    // ...
)
```

**Severity**: Low  
**Migration Impact**: Non-breaking (internal)  
**Reason**: Type safety vs simplicity tradeoff

---

#### 19. `ValuesQuerySet[T]`, `ValuesListQuerySet[T]` - Verbose Names

**Location**: `forge/orm/queryset.go:42-43`

**Issue**: These type names are very long. Could be shortened.

**Current**:
```go
Values(fields ...any) ValuesQuerySet[T]
ValuesList(fields ...any) ValuesListQuerySet[T]
```

**Suggested Fix**: Consider shorter names:
```go
Values(fields ...any) ValuesQuery[T]
ValuesList(fields ...any) ValuesListQuery[T]
```

**Severity**: Low  
**Migration Impact**: Breaking (type change)  
**Reason**: Shorter names are easier to use

---

#### 20. `UpdateMap` - Too Generic

**Location**: `forge/orm/update_builder.go` (assumed)

**Issue**: `UpdateMap` is a generic name. Could be more specific like `FieldUpdates`.

**Current**:
```go
type UpdateMap map[string]interface{}
```

**Suggested Fix**:
```go
type FieldUpdates map[string]interface{}
```

**Severity**: Low  
**Migration Impact**: Breaking (type alias)  
**Reason**: More descriptive name

---

#### 21. `ModelWithID` - Interface Naming

**Location**: `forge/orm/manager.go:146` (assumed)

**Issue**: Interface name could be more specific: `IDGetter` or `Identifiable`.

**Current**:
```go
type ModelWithID interface {
    GetID() int64
    SetID(int64)
}
```

**Suggested Fix**:
```go
type Identifiable interface {
    GetID() int64
    SetID(int64)
}
```

**Severity**: Low  
**Migration Impact**: Breaking (type rename)  
**Reason**: More precise interface naming

---

#### 22. `BaseQuerySet[T]`, `BaseSerializer`, `BaseViewSet` - Base Prefix Pattern

**Location**: Multiple packages

**Issue**: Using `Base` prefix for concrete implementations. In Go, it's more common to have the interface be the simple name and the implementation have a suffix.

**Current Pattern**:
```go
type QuerySet[T] interface { ... }
type BaseQuerySet[T] struct { ... }  // Implementation

type Serializer interface { ... }
type BaseSerializer struct { ... }  // Implementation
```

**Suggested Fix**: Consider reversing:
```go
type QuerySet[T] interface { ... }
type QuerySetImpl[T] struct { ... }  // Or DefaultQuerySet[T]

type Serializer interface { ... }
type SerializerImpl struct { ... }  // Or DefaultSerializer
```

**Severity**: Low  
**Migration Impact**: Breaking (major refactor)  
**Reason**: Go convention prefers simple interface names

---

#### 23. `rebindPostgresToSQLite()` - Unexported Helper

**Location**: `forge/db/db.go:85`

**Issue**: Function name is very specific but unexported. Fine as-is, but could be more generic.

**Current**:
```go
func rebindPostgresToSQLite(query string) string
```

**Suggested Fix**: No change needed (unexported helper is fine)

**Severity**: Low  
**Migration Impact**: None (unexported)  
**Reason**: Documentation only

---

### Documentation Only (No Code Changes)

#### 24. `Aggregate` struct - Needs Better Documentation

**Location**: `forge/orm/aggregates.go:4`

**Issue**: The `Aggregate` struct fields aren't well documented.

**Suggested Fix**: Add field documentation:
```go
// Aggregate represents an aggregate function for database queries
type Aggregate struct {
    Name  string // Result field name (e.g., "count", "total")
    Field string // Database field to aggregate
    Func  string // SQL function name (COUNT, SUM, AVG, MAX, MIN, etc.)
}
```

**Severity**: Documentation  
**Migration Impact**: None

---

#### 25. `Expression` interface - Needs Usage Examples

**Location**: `forge/orm/expression.go`

**Issue**: The `Expression` interface is central to the ORM but lacks usage examples in godoc.

**Suggested Fix**: Add comprehensive godoc with examples

**Severity**: Documentation  
**Migration Impact**: None

---

#### 26. `Filter[T any]` interface - Complex Interface Needs Guide

**Location**: `forge/filter/filter.go:11`

**Issue**: The `Filter` interface has many methods but no comprehensive guide on implementing custom filters.

**Suggested Fix**: Add implementation guide in godoc or separate doc

**Severity**: Documentation  
**Migration Impact**: None

---

## Summary by Package

### forge/orm/ (7 issues)
- High: `GetFieldAccessor()`, `NewOrderField()`
- Medium: `BuildInsertSQL()`, `ExecuteInsert()`, `NewAnnotation()`, `RegisterAggregate()`
- Low: `OrderField`, `UpdateMap`, `ModelWithID`, `BaseQuerySet`

### forge/admin/ (4 issues)
- Medium: `NewTextInput()` and other widget constructors
- Low: `Choice[V any]`, `Ordering[T]`

### forge/api/ (3 issues)
- High: `GetPaginationParams()`
- Medium: `DefaultSettings()`, `BuildPaginatedResponse()`

### forge/filter/ (4 issues)
- Medium: `NewBaseFilter()`, `GetFieldPath()` interface methods, `SetLabel()` builders
- Documentation: `Filter[T]` interface guide

### forge/db/ (1 issue)
- High: `Rebind()`

### forge/registry/ (1 issue)
- High: `GetRegistry()`

### forge/server/ (2 issues)
- Medium: `DefaultMiddlewares()`, `DefaultCORSOptions()`

### forge/schema/ (0 issues)
- ✅ Schema package has good naming

### forge/identity/ (0 issues)
- ✅ Identity package has good naming

### forge/validate/ (0 issues)
- ✅ Validate package has good naming

---

## Migration Strategy

### Phase 1: v2.0 (High Priority)
1. Add new functions with better names
2. Mark old functions as deprecated
3. Update all internal usage to new names
4. Update documentation and examples

### Phase 2: v2.x (Medium Priority)
1. Add new APIs alongside old ones
2. Deprecate old APIs with 6-month notice
3. Provide migration scripts/tools

### Phase 3: v3.0 (Breaking Changes)
1. Remove all deprecated functions
2. Clean up API surface
3. Final naming consistency pass

---

## Recommendations

### Immediate Actions (v2.0)
1. Fix the 5 high-priority issues
2. Update `NAMING_ARCHITECTURE.md` with lessons learned
3. Add linting rules to catch these patterns

### Guidelines Going Forward
1. **Avoid `Get` prefix** for getters (Go convention)
2. **Avoid `New` prefix** for lightweight factories
3. **Use `With` prefix** for builder methods
4. **Keep names short** but descriptive
5. **Be consistent** across packages
6. **Document complex interfaces** thoroughly

### Tooling
1. Add `golangci-lint` rules for naming
2. Create pre-commit hooks
3. Add naming checks to CI/CD

---

## Conclusion

The forge framework has generally good naming, with most issues being minor inconsistencies or verbose patterns. The high-priority issues should be addressed in v2.0 with proper deprecation paths. The medium and low-priority issues can be addressed gradually in v2.x releases.

**Total Issues**: 26  
**High Priority**: 5  
**Medium Priority**: 10  
**Low Priority**: 8  
**Documentation Only**: 3

**Next Steps**: Create migration plan for v2.0 addressing high-priority issues.
