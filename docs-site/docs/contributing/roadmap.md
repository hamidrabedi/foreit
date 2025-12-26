---
sidebar_position: 3
---

# Roadmap

Development roadmap and future plans for forge framework.

## Vision

forge aims to be the **Django of Go** - a full-featured, type-safe web framework that combines the best of Django's developer experience with Go's performance and type safety.

## Current Status: Foundation Complete ✅

- ✅ Schema definition system
- ✅ Code generation infrastructure
- ✅ Manager & QuerySet generation templates
- ✅ Type-safe ORM API
- ✅ Database layer with SQL builder
- ✅ HTTP routing and middleware
- ✅ Security features
- ✅ Admin system structure
- ✅ All recommended libraries integrated

## Roadmap Timeline

> **Note:** Timeline is estimated and subject to change based on development progress and community feedback.

### Q1 2025: Core Completion

**Goal:** Complete core ORM and admin functionality

**Deliverables:**
- [x] Basic QuerySet execution (All, Get, First, Last, Count, Exists) ✅
- [x] Complete QuerySet execution with SQL builder ✅
- [ ] Full admin interface (list, detail, create, edit, delete)
- [x] Complete code generation (Manager, QuerySet templates) ✅
- [ ] AST parser enhancements (extract all field options, relations, meta, hooks)
- [ ] Manager CRUD implementation (Create, Update, Delete with hooks)
- [ ] Basic testing suite

**Success Metrics:**
- Can build a complete CRUD app
- Admin interface fully functional
- All QuerySet methods working

### Q2 2025: API & Advanced Features

**Goal:** REST API generation and advanced ORM features

**Deliverables:**
- [ ] Auto-generated REST APIs (DRF-like)
- [ ] OpenAPI/Swagger documentation
- [ ] Advanced QuerySet features (F() expressions, subqueries)
- [ ] Caching layer
- [ ] Background tasks

**Success Metrics:**
- Can build a complete API
- API documentation auto-generated
- Performance benchmarks meet targets

### Q3 2025: Developer Experience

**Goal:** Improve developer experience and tooling

**Deliverables:**
- [ ] Complete CLI toolset
- [ ] Development tools (hot reload, debug toolbar)
- [ ] Comprehensive documentation
- [ ] Example applications
- [ ] Testing utilities

**Success Metrics:**
- Developer onboarding < 1 hour
- Complete documentation
- 5+ example applications

### Q4 2025: Production & Ecosystem

**Goal:** Production-ready and plugin ecosystem

**Deliverables:**
- [ ] Production optimizations
- [ ] Monitoring and observability
- [ ] Deployment guides
- [ ] Plugin system completion
- [ ] Official plugins (auth, storage, email)

**Success Metrics:**
- Production deployments
- 10+ community plugins
- Performance benchmarks

## Feature Priorities

### P0 - Critical (Must Have)

1. **Query Execution** - ✅ Basic QuerySet methods working
   - ✅ SQL builder with proper escaping and parameter binding ✅
   - 🚧 Implement SelectRelated/PrefetchRelated
   - 🚧 Implement aggregates/annotations
2. **Manager CRUD** - Implement Create, Update, Delete with hooks
3. **Admin Interface** - Complete CRUD operations
4. **Code Generation** - ✅ Templates complete
   - 🚧 Enhance AST parser to extract all options
5. **Testing** - Comprehensive test suite

### P1 - High Priority

1. **REST API Generation** - Auto-generate APIs
2. **Advanced ORM** - F() expressions, subqueries
3. **CLI Tools** - Complete command set
4. **Documentation** - Comprehensive docs

### P2 - Medium Priority

1. **Caching** - Query and model caching
2. **Background Tasks** - Async task execution
3. **Monitoring** - Metrics and observability
4. **Plugin System** - Complete plugin architecture

### P3 - Nice to Have

1. **GraphQL Support** - Optional GraphQL API
2. **WebSocket Support** - Real-time features
3. **Multi-tenancy** - Built-in multi-tenancy
4. **Internationalization** - i18n support

## Recent Achievements ✅

### Code Generation Complete (December 2024)
- ✅ Manager generation templates - Complete
- ✅ QuerySet generation templates - Complete
- ✅ BaseQuerySet exported and working correctly
- ✅ Type-safe embedding pattern established
- ✅ All templates compile and generate valid code

## Success Criteria

### MVP (Minimum Viable Product)

- [ ] Can define models with schema
- [ ] Can generate code from schemas
- [ ] Can query database with type-safe API
- [ ] Can use admin interface for CRUD
- [ ] Can build a simple blog application

**Target:** End of Q1 2025

### Beta Release

- [ ] All core features working
- [ ] REST API generation
- [ ] Comprehensive documentation
- [ ] Example applications
- [ ] Basic plugin system

**Target:** End of Q2 2025

### 1.0 Release

- [ ] Production-ready
- [ ] Performance benchmarks
- [ ] Security audit
- [ ] Plugin ecosystem
- [ ] Community adoption

**Target:** End of Q4 2025

## Community Goals

### Year 1

- **GitHub Stars:** 1,000+
- **Contributors:** 20+
- **Plugins:** 10+
- **Production Users:** 50+

### Year 2

- **GitHub Stars:** 5,000+
- **Contributors:** 50+
- **Plugins:** 50+
- **Production Users:** 500+

## See Also

- [Development Guide](development) - Contributing guide
- [Architecture](architecture) - Framework architecture

