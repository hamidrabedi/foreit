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
   - Manager API (templates complete, CRUD methods use SQL builder)
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

10. **Admin System**
    - Auto-generation from models
    - Django-style registration
    - HTMX integration
    - Sprig template functions

11. **Authentication**
    - Password hashing (bcrypt)
    - Authentication middleware

12. **Utilities**
    - String case conversion (strcase)
    - UUID generation (google/uuid)
    - Testing toolkit (testify)

## Future Features (Planned)

### Phase 1: Core Completion (Weeks 1-4)

1. **Complete Code Generation**
   - [x] Manager generation templates ✅
   - [x] QuerySet generation templates ✅
   - [ ] Complete AST parser (extract all field options)
   - [ ] Relation extraction and generation
   - [ ] Meta extraction and generation
   - [ ] Hook extraction and generation

2. **Query Execution**
   - [x] Implement QuerySet.All() ✅ (uses database/sql, reflection-based scanning)
   - [x] Implement QuerySet.Get() ✅
   - [x] Implement QuerySet.First(), Last(), Count(), Exists() ✅
   - [ ] Implement SelectRelated (JOIN queries) - structure ready
   - [ ] Implement PrefetchRelated (separate queries) - structure ready
   - [ ] Implement aggregates execution - structure ready
   - [ ] Implement annotations execution - structure ready
   - [ ] Implement Values/ValuesList - structure ready
   - [ ] Implement BulkUpdate/BulkCreate - structure ready
   - [x] SQL builder for type-safe query execution ✅

3. **Admin Interface**
   - [ ] List view implementation
   - [ ] Detail view implementation
   - [ ] Create/Edit form implementation
   - [ ] Delete confirmation
   - [ ] Search functionality
   - [ ] Filtering
   - [ ] Pagination
   - [ ] Bulk actions

### Phase 2: API Generation (Weeks 5-8)

4. **REST API Auto-Generation**
   - [ ] ViewSet generation
   - [ ] Serializer generation
   - [ ] Router registration
   - [ ] Pagination
   - [ ] Filtering
   - [ ] Search
   - [ ] Ordering
   - [ ] Permissions

5. **API Features**
   - [ ] OpenAPI/Swagger documentation
   - [ ] API versioning
   - [ ] Rate limiting
   - [ ] CORS support
   - [ ] API authentication (JWT, API keys)

### Phase 3: Advanced Features (Weeks 9-12)

6. **Advanced ORM Features**
   - [ ] F() expressions (database functions)
   - [ ] Q() object improvements
   - [ ] Subqueries
   - [ ] Raw SQL support
   - [ ] Database functions
   - [ ] Window functions
   - [ ] Full-text search

7. **Caching**
   - [ ] Query result caching
   - [ ] Model instance caching
   - [ ] Cache invalidation
   - [ ] Redis support
   - [ ] In-memory cache

8. **Background Tasks**
   - [ ] Task queue integration
   - [ ] Scheduled tasks
   - [ ] Async task execution
   - [ ] Task monitoring

### Phase 4: Developer Experience (Weeks 13-16)

9. **CLI Enhancements**
   - [ ] `forge startapp` - Create new app
   - [ ] `forge shell` - Interactive shell
   - [ ] `forge test` - Test runner
   - [ ] `forge collectstatic` - Static file collection
   - [ ] `forge createsuperuser` - Admin user creation
   - [ ] `forge dbshell` - Database shell

10. **Development Tools**
    - [ ] Hot reload (file watching)
    - [ ] Debug toolbar
    - [ ] Query logging
    - [ ] Performance profiling
    - [ ] Code coverage

11. **Documentation**
    - [ ] Auto-generated API docs
    - [ ] Schema documentation
    - [ ] Tutorials
    - [ ] Examples

### Phase 5: Production Features (Weeks 17-20)

12. **Performance**
    - [ ] Query optimization
    - [ ] Connection pooling tuning
    - [ ] Query result caching
    - [ ] Database query analysis

13. **Monitoring**
    - [ ] Metrics collection
    - [ ] Health checks
    - [ ] Error tracking
    - [ ] Performance monitoring

14. **Deployment**
    - [ ] Docker support
    - [ ] Kubernetes manifests
    - [ ] Deployment guides
    - [ ] Production checklist

### Phase 6: Ecosystem (Weeks 21-24)

15. **Plugin System**
    - [ ] Plugin registry
    - [ ] Plugin lifecycle
    - [ ] Plugin dependencies
    - [ ] Plugin marketplace

16. **Official Plugins**
    - [ ] Authentication plugin (OAuth, JWT)
    - [ ] File storage plugin (S3, local)
    - [ ] Email plugin
    - [ ] Search plugin (Elasticsearch)
    - [ ] Cache plugin (Redis)

17. **Community**
    - [ ] Plugin examples
    - [ ] Best practices guide
    - [ ] Community plugins
    - [ ] Contributing guide

## Feature Matrix

| Feature | Status | Priority | Phase |
|---------|--------|----------|-------|
| Schema Definition | ✅ Complete | P0 | Done |
| Code Generation (Basic) | ✅ Complete | P0 | Done |
| Type-Safe ORM API | ✅ Complete | P0 | Done |
| Query Execution | 🚧 In Progress | P0 | Phase 1 |
| Admin Interface | 🚧 In Progress | P0 | Phase 1 |
| REST API Generation | 📋 Planned | P1 | Phase 2 |
| Advanced ORM | 📋 Planned | P1 | Phase 3 |
| Caching | 📋 Planned | P2 | Phase 3 |
| Background Tasks | 📋 Planned | P2 | Phase 3 |
| CLI Tools | 🚧 Partial | P1 | Phase 4 |
| Development Tools | 📋 Planned | P2 | Phase 4 |
| Performance | 📋 Planned | P1 | Phase 5 |
| Monitoring | 📋 Planned | P2 | Phase 5 |
| Plugin System | ✅ Structure | P1 | Phase 6 |

**Legend:**
- ✅ Complete
- 🚧 In Progress
- 📋 Planned
- P0 = Critical, P1 = High, P2 = Medium

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

