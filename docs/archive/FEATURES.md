# forge Framework - Features Documentation

## Current Features

### Code Generation ✅
- ✅ AST-based schema parsing
- ✅ Model struct generation
- ✅ FieldExpr generation for type-safe field access
- ✅ **Manager generation** - Complete templates ✅
- ✅ **QuerySet generation** - Complete templates ✅
- ✅ SQL builder with proper escaping and parameter binding

### ✅ Core Framework

1. **Schema Definition System**
   - Declarative model definitions in Go
   - Full Django field options
   - Relationship definitions (ForeignKey, OneToOne, ManyToMany)
   - Model metadata (Meta options)
   - Lifecycle hooks

2. **Code Generation**
   - AST-based schema parsing
   - Model struct generation
   - FieldExpr generation
   - **Manager generation** ✅
   - **QuerySet generation** ✅
   - SQL builder with proper escaping

3. **Type-Safe ORM**
   - `FieldExpr[T]` for type-safe field access ✅
   - `QueryExpr` for complex queries ✅
   - QuerySet API (Django-like) ✅
     - Filter, Exclude, OrderBy, Limit, Offset, Distinct ✅
     - All, Get, First, Last, Count, Exists ✅
     - Select, Only, Defer (structure ready)
     - SelectRelated, PrefetchRelated (structure ready)
     - Aggregates, Annotations (structure ready)
   - Manager API ✅ (Complete CRUD with hooks)
   - Dynamic query API ✅

4. **Database Layer**
   - PostgreSQL support
   - Connection pooling
   - Transaction management
   - Migration system (golang-migrate)
   - SQL builder with proper escaping and parameter binding

5. **HTTP & Routing**
   - Chi router wrapper
   - Middleware stack
   - Request context utilities
   - Server wrapper

6. **Security**
   - CSRF protection (gorilla/csrf)
   - Session management (scs)
   - XSS protection
   - SQL injection prevention

7. **Validation**
   - go-playground/validator integration
   - Auto-generated validation tags
   - Custom validators

8. **Logging**
   - Structured logging (zap)
   - Request logging middleware

9. **Configuration**
   - Viper integration
   - YAML, JSON, env var support
   - Framework settings

10. **Admin System** ✅
    - Type-safe admin interface with generics
    - Django-style registration
    - Complete HTTP handlers (List, Detail, Create, Update, Delete)
    - Rich form widgets
    - Filters, search, pagination
    - Bulk actions
    - CSV/JSON export
    - Inlines and fieldsets
    - Permission system

11. **Authentication** ✅
    - Password hashing (bcrypt)
    - Authentication middleware
    - Complete user system with authentication backends
    - Session management
    - Permission system (RBAC)

12. **Utilities**
    - String case conversion (strcase)
    - UUID generation (google/uuid)
    - Testing toolkit (testify)

## Feature Status Summary

### ✅ Complete Features

| Feature | Status | Details |
|---------|--------|---------|
| Schema Definition | ✅ Complete | Full Django field options, relations, meta, hooks |
| Code Generation | ✅ Complete | Models, Managers, QuerySets, FieldExpr |
| Type-Safe ORM | ✅ Complete | QuerySet API, Manager CRUD, SQL builder |
| Admin System | ✅ Complete | Type-safe admin with HTTP handlers, widgets, export |
| REST API Framework | ✅ Complete | ViewSets, Serializers, Auth, Permissions, Throttling |
| User System | ✅ Complete | User management, auth, sessions, permissions |
| Migration System | ✅ Complete | Built-in migrations with golang-migrate |
| CLI Tools | ✅ Complete | new, generate, migrate, runserver |
| Security | ✅ Complete | CSRF, XSS, SQL injection protection |
| Plugin System | ✅ Complete | Extensible plugin architecture |

### 🚧 In Progress / Structure Ready

| Feature | Status | Next Steps |
|---------|--------|------------|
| Advanced ORM | 🚧 Structure Ready | SelectRelated, PrefetchRelated, Aggregates, Annotations |
| AST Parser | 🚧 Partial | Extract all field/relation/meta/hook options |
| Admin Templates | 🚧 Core Complete | HTML template rendering, rich text editor |
| Testing Infrastructure | 🚧 Partial | Comprehensive test suite, fixtures |

### 📋 Planned Features

See [Roadmap](ROADMAP.md) for detailed planned features including:
- REST API auto-generation
- Advanced query features (F() expressions, subqueries)
- Caching layer
- Background tasks
- Development tools
- Monitoring & observability
- GraphQL support
- WebSocket support
- Multi-tenancy
- Internationalization

**For detailed roadmap and timeline, see [ROADMAP.md](ROADMAP.md)**

## Architecture Decisions

### Why SQL Builder?
- Type-safe SQL generation from Go code
- Proper identifier escaping to prevent SQL injection
- Parameter binding for all values
- Works with standard database/sql
- Go-first approach (not SQL-first)

### Why chi?
- Lightweight, composable
- Standard library compatible
- Framework-friendly
- Widely adopted

### Why code generation?
- Type safety
- Performance (no reflection)
- IDE support
- Compile-time checks

### Why dual API?
- Type-safe for common cases
- Dynamic for edge cases
- Best of both worlds

