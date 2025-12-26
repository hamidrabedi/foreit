# Changelog

All notable changes to the forge framework will be documented in this file.

## [Unreleased]

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
- **Documentation cleanup**: Removed 24+ outdated progress/status reports
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

