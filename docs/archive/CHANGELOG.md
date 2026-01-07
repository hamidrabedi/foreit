# Changelog

All notable changes to the forge framework will be documented in this file.

## [Unreleased]

## [2.0.0] - Naming Consistency Improvements

### Added
- **New Naming Functions**: Clearer, more idiomatic Go naming
  - `RebindPlaceholders()` - Clearer name for database placeholder conversion
  - `Global()` - Access global registry without redundant package name
  - `FieldAccessor()` - Go-conventional getter (no Get prefix)
  - `ParsePaginationParams()` - Clearer intent (parsing vs getting)
- **Migration Script**: Automated migration script for v1.x → v2.0
- **Migration Checklist**: Step-by-step migration checklist

### Changed
- **Database Operations**: `Rebind()` → `RebindPlaceholders()` (deprecated)
- **Registry Access**: `GetRegistry()` → `Global()` (deprecated)
- **Ordering**: `NewOrderField()` → `Asc()` / `Desc()` (deprecated)
- **Field Access**: `GetFieldAccessor()` → `FieldAccessor()` (deprecated)
- **Pagination**: `GetPaginationParams()` → `ParsePaginationParams()` (deprecated)

### Deprecated
- `Rebind()` - Use `RebindPlaceholders()` for clarity
- `GetRegistry()` - Use `Global()` to avoid package name redundancy
- `NewOrderField()` - Use `Asc()` or `Desc()` for clarity
- `GetFieldAccessor()` - Use `FieldAccessor()` (Go convention: no Get prefix)
- `GetPaginationParams()` - Use `ParsePaginationParams()` for clarity

### Documentation
- Updated `MIGRATION_V1_TO_V2.md` with v2.0 naming changes
- Added `MIGRATION_CHECKLIST_V2.md` - Step-by-step migration checklist
- Added `NAMING_AUDIT_V2.md` - Comprehensive naming audit report
- Added `NAMING_ISSUES_CATEGORIZED.md` - Categorized issues by priority
- Updated all examples to use new naming

### Fixed
- Improved naming consistency across all packages
- Better Go idiom compliance
- Clearer API intent through better naming

## [1.5.0] - Naming Refactor & API Improvements

### Added
- **New Query Expression API**: `And()`, `Or()`, `Not()` functions for explicit boolean combinations
- **Runtime Field References**: `F()` and `FieldRef()` functions for dynamic field access
- **SQL-like Conditions**: `Where()` function for explicit field conditions
- **Comprehensive Naming Documentation**: Complete naming architecture and conventions
- **Migration Guide**: Detailed guide for migrating from v1.x to v2.0

### Changed
- **Field Type Consolidation**: `FieldExpression[T]` renamed to `Field[T]` (backward-compatible alias provided)
- **Method Names**: Shorter, clearer method names:
  - `Equals()` → `Eq()`
  - `NotEquals()` → `Ne()`
  - `Greater()` → `Gt()`
  - `GreaterOrEqual()` → `Gte()`
  - `Less()` → `Lt()`
  - `LessOrEqual()` → `Lte()`
- **Get() Semantics**: Enhanced documentation clarifying `Manager.Get()` vs `QuerySet.Get()`
- **Internal Migration**: Updated admin package to use new API

### Deprecated
- `NewQ()` - Use `And()`, `Or()`, `Not()` directly instead
- `NewFieldQueryExpr()` - Use `Where()` or field methods instead
- `FieldExpr[T]` - Use `Field[T]` (via `NewField`) instead
- `NewFieldExpr()` - Use `NewField()` instead

### Documentation
- Added `NAMING_ARCHITECTURE.md` - Complete naming conventions and design philosophy
- Added `NAMING_AUDIT.md` - Current issues and migration plan
- Added `NAMING_QUICK_REFERENCE.md` - Developer cheat sheet
- Added `NAMING_INDEX.md` - Navigation and status dashboard
- Added `MIGRATION_V1_TO_V2.md` - User migration guide
- Updated `API_REFERENCE.md` with new API examples

### Fixed
- Receiver naming consistency across all types
- Clarified Get() method semantics in godoc
- Improved code readability with meaningful names

## [Unreleased] (Previous)

### Added
- Comprehensive documentation system
- Error handling package with proper error types
- UUID field builder
- String case conversion utilities (strcase)
- Testing framework (testify)

### Changed
- Improved error messages throughout codebase
- Replaced generic "not implemented" errors with typed errors
- Enhanced code organization and structure
- **Documentation consolidation**: Merged detailed architecture docs into ARCHITECTURE.md
- **Documentation cleanup**: Removed 50+ outdated progress/status/completion reports
- **Documentation refactoring**: Merged admin documentation into single comprehensive guide
- Updated INDEX.md with organized structure and all current documentation
- Updated FEATURES.md and ROADMAP.md with accurate implementation status
- Fixed broken documentation references

### Fixed
- Documentation inconsistencies
- Error handling patterns
- Code organization issues
- Broken links to deleted documentation files

## [0.1.0] - Foundation Release

### Added
- Schema definition system
- Code generation infrastructure
- Type-safe ORM API structure
- Database layer with SQL builder
- HTTP routing and middleware
- Security features
- Admin system structure
- Extension/plugin architecture
- Migration system
- Configuration system
- Logging system
- Validation system
- Authentication utilities

### Infrastructure
- Chi router integration
- golang-migrate integration
- gorilla/csrf integration
- scs session management
- zap logging
- viper configuration
- validator integration
- bcrypt password hashing
- sprig template functions
- HTMX support

## Notes

- Version numbers follow [Semantic Versioning](https://semver.org/)
- Unreleased changes are features in development
- Breaking changes will be clearly marked

