# Admin Framework - Integration with Existing Forge Systems

> **How the admin framework integrates with and extends existing Forge ORM, Schema, Filter, API, and Server systems**

## Overview

The admin framework is **not a replacement** for existing Forge systems - it's a **layer on top** that provides:
- Auto-generated UI from existing schemas
- REST/GraphQL API endpoints using existing API system
- Advanced filtering using existing filter system
- CRUD operations using existing ORM
- Server integration using existing server package

## I. Schema System Integration

### Using Existing Schema

The admin system reads from the **existing** `schema.Schema` interface:

```go
// forge/schema/schema.go (EXISTING)
type Schema interface {
    Fields() []Field
    Relations() []Relation
    Meta() Meta
    Hooks() *ModelHooks
}
```

### Admin Enhancement (NEW)

Admin adds **metadata extraction** without modifying core schema:

```go
// forge/admin/schema/discovery.go (NEW)
package adminschema

import (
    "github.com/forgego/forge/schema"
)

// DiscoverFields extracts field metadata from existing schema
func DiscoverFields[T any](s schema.Schema) ([]FieldInfo, error) {
    fields := s.Fields()
    result := make([]FieldInfo, len(fields))

    for i, field := range fields {
        result[i] = FieldInfo{
            Name:         field.Name(),
            Type:         field.Type(),
            Label:        field.VerboseName(),
            HelpText:     field.HelpText(),
            Required:     field.Required(),
            ReadOnly:     field.ReadOnly(),
            Choices:      extractChoices(field),
            Widget:       determineWidget(field), // NEW: infer widget from field type
            Validators:   extractValidators(field),
            DefaultValue: field.Default(),
        }
    }

    return result, nil
}

// NEW: Widget inference from field type
func determineWidget(field schema.Field) string {
    switch field.Type() {
    case schema.TypeString:
        if field.MaxLength() > 200 {
            return "textarea"
        }
        return "text"
    case schema.TypeText:
        return "rich_text"
    case schema.TypeEmail:
        return "email"
    case schema.TypePassword:
        return "password"
    case schema.TypeInt64, schema.TypeInt32:
        return "number"
    case schema.TypeBool:
        return "checkbox"
    case schema.TypeDate:
        return "date"
    case schema.TypeDateTime:
        return "datetime"
    case schema.TypeFile:
        return "file_upload"
    case schema.TypeImage:
        return "image_upload"
    default:
        return "text"
    }
}
```

### Example Usage

```go
// User model with EXISTING schema
type User struct {
    schema.BaseModel
}

func (User) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("email").Required().MaxLength(255).Build(),
        schema.String("name").Required().MaxLength(100).Build(),
        schema.Text("bio").Build(),
        schema.Image("avatar").Build(),
        schema.DateTime("created_at").AutoNow().Build(),
    }
}

// Admin automatically discovers all fields and their widgets!
// No need to duplicate field definitions
```

---

## II. ORM Integration

### Using Existing Manager & QuerySet

Admin uses the **existing** `orm.Manager[T]` and `orm.QuerySet[T]`:

```go
// forge/admin/orm/manager.go (NEW - wraps existing ORM)
package adminorm

import (
    "github.com/forgego/forge/orm"
)

// AdminManager wraps existing orm.Manager with admin-specific features
type AdminManager[T any] struct {
    manager *orm.Manager[T]  // REUSE existing manager
}

func NewAdminManager[T any](manager *orm.Manager[T]) (*AdminManager[T], error) {
    return &AdminManager[T]{
        manager: manager,
    }, nil
}

// GetQueryset returns the base QuerySet from existing ORM
func (m *AdminManager[T]) GetQueryset(ctx context.Context) (orm.QuerySet[T], error) {
    // REUSE existing Manager.All() to get base queryset
    qs, err := m.manager.Filter(orm.Q{}) // Empty filter = all
    if err != nil {
        return nil, err
    }
    return qs, nil
}

// Create, Update, Delete delegate to existing Manager
func (m *AdminManager[T]) Create(ctx context.Context, instance *T) error {
    return m.manager.Create(ctx, instance) // REUSE
}

func (m *AdminManager[T]) Update(ctx context.Context, instance *T) error {
    return m.manager.Update(ctx, instance) // REUSE
}

func (m *AdminManager[T]) Delete(ctx context.Context, instance *T) error {
    return m.manager.Delete(ctx, instance) // REUSE
}

// Manager returns the underlying Manager for direct access
func (m *AdminManager[T]) Manager() *orm.Manager[T] {
    return m.manager
}
```

### Admin List Handler (REUSES ORM)

```go
// forge/admin/http/handlers/list.go (NEW)
package handlers

import (
    "github.com/forgego/forge/admin"
    "github.com/forgego/forge/orm"
)

func HandleList[T any](adm *admin.Admin[T]) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Get base queryset from existing ORM
        qs, err := adm.Manager().GetQueryset(ctx)
        if err != nil {
            respondError(w, err)
            return
        }

        // Apply filters using EXISTING filter system (see next section)
        qs = applyFilters(qs, adm.FilterSet(), r.URL.Query())

        // Apply search using EXISTING ORM search capabilities
        if search := r.URL.Query().Get("search"); search != "" {
            qs = qs.Search(search, adm.Config().SearchFields...) // REUSE existing ORM search
        }

        // Apply ordering using EXISTING ORM
        if ordering := parseOrdering(r.URL.Query().Get("ordering")); len(ordering) > 0 {
            qs = qs.OrderBy(ordering...) // REUSE existing ORM ordering
        }

        // Apply pagination using EXISTING ORM
        page, pageSize := parsePagination(r.URL.Query())
        offset := (page - 1) * pageSize
        qs = qs.Limit(pageSize).Offset(offset) // REUSE existing ORM pagination

        // Execute query using EXISTING ORM
        results, err := qs.All(ctx)
        if err != nil {
            respondError(w, err)
            return
        }

        // Get total count using EXISTING ORM
        total, err := qs.Count(ctx)
        if err != nil {
            respondError(w, err)
            return
        }

        // Respond
        respondJSON(w, PaginatedResponse{
            Count:      total,
            PageSize:   pageSize,
            Page:       page,
            TotalPages: (total + pageSize - 1) / pageSize,
            Results:    results,
        })
    }
}
```

---

## III. Filter System Integration

### Extending Existing Filter System

Admin **extends** the existing filter system:

```go
// forge/filter/filter.go (EXISTING)
package filter

type Filter interface {
    Apply(qs QuerySet, value interface{}) QuerySet
    Name() string
    Type() FilterType
}

type FilterSet[T any] struct {
    filters map[string]Filter
}

func (fs *FilterSet[T]) Apply(qs QuerySet, params map[string]interface{}) QuerySet {
    for name, value := range params {
        if f, ok := fs.filters[name]; ok {
            qs = f.Apply(qs, value)
        }
    }
    return qs
}
```

### Admin Filter Builder (NEW - uses existing filters)

```go
// forge/admin/filter/builder.go (NEW)
package adminfilter

import (
    "github.com/forgego/forge/filter"
    "github.com/forgego/forge/schema"
    "github.com/forgego/forge/orm"
)

// BuildFiltersFromSchema auto-generates filters from schema
func BuildFiltersFromSchema[T any](s schema.Schema, fs *filter.FilterSet[T]) error {
    fields := s.Fields()

    for _, field := range fields {
        // Create appropriate filter based on field type
        switch field.Type() {
        case schema.TypeString, schema.TypeEmail:
            // REUSE existing TextFilter
            fs.Add(filter.NewTextFilter(field.Name()))
            // Add contains filter
            fs.Add(filter.NewTextFilter(field.Name() + "__contains"))

        case schema.TypeInt64, schema.TypeInt32:
            // REUSE existing NumberFilter
            fs.Add(filter.NewNumberFilter(field.Name()))
            fs.Add(filter.NewNumberFilter(field.Name() + "__gt"))
            fs.Add(filter.NewNumberFilter(field.Name() + "__lt"))
            fs.Add(filter.NewNumberFilter(field.Name() + "__gte"))
            fs.Add(filter.NewNumberFilter(field.Name() + "__lte"))

        case schema.TypeBool:
            // REUSE existing BooleanFilter
            fs.Add(filter.NewBooleanFilter(field.Name()))

        case schema.TypeDate, schema.TypeDateTime:
            // REUSE existing DateRangeFilter
            fs.Add(filter.NewDateRangeFilter(field.Name()))
            fs.Add(filter.NewDateFilter(field.Name() + "__gte"))
            fs.Add(filter.NewDateFilter(field.Name() + "__lte"))

        case schema.TypeForeignKey:
            // REUSE existing RelationFilter
            relatedModel := field.RelatedModel()
            fs.Add(filter.NewRelationFilter(field.Name(), relatedModel))

        case schema.TypeManyToMany:
            // REUSE existing ManyToManyFilter
            relatedModel := field.RelatedModel()
            fs.Add(filter.NewManyToManyFilter(field.Name(), relatedModel))
        }

        // Add choices filter if field has choices
        if choices := field.Choices(); len(choices) > 0 {
            fs.Add(filter.NewChoicesFilter(field.Name(), choices))
        }
    }

    return nil
}

// AdminFilterSet wraps existing FilterSet with metadata
type AdminFilterSet[T any] struct {
    filterSet *filter.FilterSet[T] // REUSE existing
    metadata  []FilterMetadata      // NEW: for API/UI
}

func NewAdminFilterSet[T any]() (*AdminFilterSet[T], error) {
    return &AdminFilterSet[T]{
        filterSet: filter.NewFilterSet[T](), // REUSE existing
        metadata:  []FilterMetadata{},
    }, nil
}

func (afs *AdminFilterSet[T]) FilterSet() *filter.FilterSet[T] {
    return afs.filterSet
}

func (afs *AdminFilterSet[T]) Apply(qs orm.QuerySet[T], params map[string]interface{}) orm.QuerySet[T] {
    // Delegate to existing FilterSet
    return afs.filterSet.Apply(qs, params)
}
```

### Usage in Admin

```go
// User admin with auto-generated filters
userAdmin, err := admin.Register[User](
    userSchema,
    userManager,
    &admin.Config[User]{
        // Filters automatically generated from schema!
        // But you can override:
        ListFilter: []interface{}{
            "email",           // Auto: TextFilter
            "is_active",       // Auto: BooleanFilter
            "created_at",      // Auto: DateRangeFilter
            CustomStatusFilter(), // Custom filter
        },
    },
)

// Behind the scenes, admin uses EXISTING filter.FilterSet
```

---

## IV. API System Integration

### Extending Existing API Package

```go
// forge/api/api.go (EXISTING)
package api

type ViewSet[T any] interface {
    List(ctx context.Context, r Request) Response
    Retrieve(ctx context.Context, id int64) Response
    Create(ctx context.Context, data T) Response
    Update(ctx context.Context, id int64, data T) Response
    Destroy(ctx context.Context, id int64) Response
}
```

### Admin API Router (NEW - uses existing API system)

```go
// forge/admin/api/rest/router.go (NEW)
package rest

import (
    "github.com/forgego/forge/admin"
    "github.com/forgego/forge/api"
    "github.com/forgego/forge/server"
    "github.com/go-chi/chi/v5"
)

type AdminAPIRouter struct {
    registry *admin.Registry
    server   *server.Server // REUSE existing server
}

func NewAdminAPIRouter(registry *admin.Registry, srv *server.Server) *AdminAPIRouter {
    return &AdminAPIRouter{
        registry: registry,
        server:   srv,
    }
}

func (r *AdminAPIRouter) RegisterRoutes(router chi.Router) {
    // Use existing server's middleware
    router.Use(r.server.Middleware.CORS())         // REUSE
    router.Use(r.server.Middleware.RateLimit())    // REUSE
    router.Use(r.server.Middleware.Authenticate()) // REUSE

    // Admin metadata endpoints
    router.Get("/api/admin/meta", r.handleMetaList)
    router.Get("/api/admin/meta/{model}", r.handleMetaDetail)

    // Model CRUD endpoints (auto-registered for each model)
    for name, adm := range r.registry.GetAll() {
        r.registerModelRoutes(router, name, adm)
    }

    // Global search
    router.Get("/api/admin/search", r.handleGlobalSearch)

    // File upload (REUSE existing api.FileHandler)
    router.Post("/api/admin/upload", api.FileUploadHandler(r.server.Config.UploadPath))
}

func (r *AdminAPIRouter) registerModelRoutes(router chi.Router, name string, adm admin.AdminInterface) {
    basePath := fmt.Sprintf("/api/admin/%s", name)

    router.Route(basePath, func(r chi.Router) {
        // REUSE existing API patterns
        r.Get("/", r.handleList(adm))           // List with pagination/filtering
        r.Post("/", r.handleCreate(adm))        // Create
        r.Get("/{id}", r.handleDetail(adm))     // Retrieve
        r.Patch("/{id}", r.handleUpdate(adm))   // Partial update
        r.Put("/{id}", r.handleReplace(adm))    // Full update
        r.Delete("/{id}", r.handleDelete(adm))  // Delete

        // Bulk operations
        r.Post("/bulk-create", r.handleBulkCreate(adm))
        r.Post("/bulk-update", r.handleBulkUpdate(adm))
        r.Delete("/bulk-delete", r.handleBulkDelete(adm))

        // Actions
        r.Post("/action/{action}", r.handleAction(adm))

        // Autocomplete for relation fields
        r.Get("/autocomplete", r.handleAutocomplete(adm))
    })
}
```

### Serializer Integration (REUSES existing API serializers)

```go
// forge/admin/api/serializers/model.go (NEW - extends existing)
package serializers

import (
    "github.com/forgego/forge/api/serializers"
)

// ModelSerializer extends existing API serializer
type ModelSerializer[T any] struct {
    serializers.BaseSerializer[T] // REUSE existing base
    admin                          *admin.Admin[T]
}

func NewModelSerializer[T any](adm *admin.Admin[T]) *ModelSerializer[T] {
    return &ModelSerializer[T]{
        BaseSerializer: serializers.NewBaseSerializer[T](), // REUSE existing
        admin:          adm,
    }
}

func (s *ModelSerializer[T]) Serialize(instance *T) map[string]interface{} {
    // REUSE existing serialization
    data := s.BaseSerializer.Serialize(instance)

    // Add admin-specific metadata
    data["_meta"] = map[string]interface{}{
        "can_edit":   s.admin.HasChangePermission(context.TODO(), nil, instance),
        "can_delete": s.admin.HasDeletePermission(context.TODO(), nil, instance),
    }

    return data
}
```

---

## V. Server Integration

### Using Existing Server Package

```go
// forge/server/server.go (EXISTING)
package server

type Server struct {
    Router     chi.Router
    Config     *Config
    Middleware MiddlewareRegistry
    DB         *db.DB
}
```

### Admin Server Integration (NEW - extends existing)

```go
// forge/admin/server.go (NEW)
package admin

import (
    "github.com/forgego/forge/server"
    "github.com/forgego/forge/admin/api/rest"
)

// RegisterAdmin adds admin routes to existing server
func RegisterAdmin(srv *server.Server, registry *Registry) error {
    // Create admin API router
    adminAPI := rest.NewAdminAPIRouter(registry, srv)

    // Register routes on existing server
    adminAPI.RegisterRoutes(srv.Router)

    // Optionally serve static admin UI
    if srv.Config.AdminUI.Enabled {
        srv.Router.Handle("/admin/*", http.FileServer(http.Dir(srv.Config.AdminUI.StaticPath)))
    }

    return nil
}

// Example usage
func main() {
    // Create existing server
    srv := server.New(&server.Config{
        Host: "localhost",
        Port: 8000,
        DB:   dbConfig,
    })

    // Register models with ORM (EXISTING)
    userManager, _ := orm.NewManager[User]("users")
    userManager.SetDB(srv.DB)

    // Create admin registry
    adminRegistry := admin.NewRegistry()

    // Register models with admin
    admin.Register[User](
        &UserSchema{},
        userManager,
        &admin.Config[User]{},
    )

    // Integrate admin with existing server
    admin.RegisterAdmin(srv, adminRegistry)

    // Start server (EXISTING)
    srv.Start()
}
```

---

## VI. Extending Existing Systems (When Needed)

### 6.1 Adding Features to Filter System

If admin needs a filter type that doesn't exist:

```go
// forge/filter/facet.go (ADD to existing filter package)
package filter

// FacetFilter for faceted search (needed by admin)
type FacetFilter struct {
    field string
}

func NewFacetFilter(field string) *FacetFilter {
    return &FacetFilter{field: field}
}

func (f *FacetFilter) Apply(qs QuerySet, value interface{}) QuerySet {
    // Implementation
    return qs
}

func (f *FacetFilter) GetFacets(qs QuerySet) map[string]int {
    // Return facet counts
    return nil
}
```

Then admin can use it:

```go
// forge/admin/filter/facets.go (NEW - uses the filter we added)
package adminfilter

import "github.com/forgego/forge/filter"

func BuildFacetedSearch[T any](field string) filter.Filter {
    return filter.NewFacetFilter(field) // REUSE
}
```

### 6.2 Adding Features to ORM

If admin needs ORM features that don't exist:

```go
// forge/orm/search.go (ADD to existing ORM package)
package orm

// Search performs full-text search (needed by admin)
func (qs *QuerySet[T]) Search(query string, fields ...string) QuerySet[T] {
    // Build search expression
    expr := buildSearchExpression(query, fields)

    // Apply using existing Filter
    return qs.Filter(expr)
}

// Annotate adds computed fields (needed by admin for aggregations)
func (qs *QuerySet[T]) Annotate(name string, expr Expression) QuerySet[T] {
    qs.annotations[name] = expr
    return qs
}
```

Then admin uses it:

```go
// forge/admin/orm/queryset.go (NEW - uses ORM features)
package adminorm

func (m *AdminManager[T]) SearchQueryset(ctx context.Context, query string, fields []string) (orm.QuerySet[T], error) {
    qs, err := m.GetQueryset(ctx)
    if err != nil {
        return nil, err
    }

    return qs.Search(query, fields...), nil // REUSE ORM search
}
```

### 6.3 Adding Features to API

If admin needs API features:

```go
// forge/api/pagination.go (ADD to existing API package)
package api

type PaginatedResponse[T any] struct {
    Count      int64 `json:"count"`
    Next       string `json:"next,omitempty"`
    Previous   string `json:"previous,omitempty"`
    PageSize   int   `json:"page_size"`
    Page       int   `json:"page"`
    TotalPages int   `json:"total_pages"`
    Results    []T   `json:"results"`
}

func Paginate[T any](results []T, total int64, page int, pageSize int, baseURL string) *PaginatedResponse[T] {
    // Implementation
    return &PaginatedResponse[T]{
        Count:      total,
        PageSize:   pageSize,
        Page:       page,
        TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
        Results:    results,
        Next:       buildNextURL(baseURL, page, pageSize, total),
        Previous:   buildPrevURL(baseURL, page),
    }
}
```

---

## VII. Integration Summary

### What Admin REUSES:
✅ `schema.Schema` - No duplication of field definitions
✅ `orm.Manager[T]` - All CRUD operations
✅ `orm.QuerySet[T]` - Filtering, ordering, pagination
✅ `filter.FilterSet[T]` - All existing filters
✅ `api.ViewSet[T]` - API patterns
✅ `api.Serializer[T]` - Serialization logic
✅ `server.Server` - HTTP server, middleware, routing
✅ `server.Middleware` - Auth, CORS, rate limiting
✅ `db.DB` - Database connection

### What Admin ADDS:
➕ Metadata extraction from schemas
➕ Auto-generated CRUD UI
➕ Widget system for dashboards
➕ Plugin system for extensibility
➕ Admin-specific API endpoints
➕ Permission system for admin access
➕ Action framework for bulk operations
➕ Advanced search UI
➕ File manager UI

### What Admin MAY EXTEND (in core packages):
🔧 `orm.QuerySet.Search()` - Full-text search
🔧 `orm.QuerySet.Annotate()` - Computed fields
🔧 `filter.FacetFilter` - Faceted search
🔧 `api.PaginatedResponse` - Standardized pagination
🔧 `api.BulkOperations` - Bulk create/update/delete

---

## VIII. Example: Complete Integration

```go
// main.go
package main

import (
    "github.com/forgego/forge/server"
    "github.com/forgego/forge/orm"
    "github.com/forgego/forge/admin"
    "myapp/models"
)

func main() {
    // 1. Create server (EXISTING)
    srv := server.New(&server.Config{
        Host: "localhost",
        Port: 8000,
    })

    // 2. Set up database (EXISTING)
    db, _ := srv.DB()

    // 3. Create managers (EXISTING)
    userManager, _ := orm.NewManager[models.User]("users")
    userManager.SetDB(db)

    postManager, _ := orm.NewManager[models.Post]("posts")
    postManager.SetDB(db)

    // 4. Create admin registry (NEW)
    adminRegistry := admin.NewRegistry()

    // 5. Register models with admin (NEW)
    admin.Register[models.User](
        &models.UserSchema{},
        userManager,
        &admin.Config[models.User]{
            ListDisplay:  []string{"id", "email", "name", "is_active"},
            SearchFields: []string{"email", "name"},
            ListFilter:   []string{"is_active", "created_at"},
        },
    )

    admin.Register[models.Post](
        &models.PostSchema{},
        postManager,
        &admin.Config[models.Post]{
            ListDisplay:  []string{"id", "title", "author", "status", "created_at"},
            SearchFields: []string{"title", "content"},
            ListFilter:   []string{"status", "author", "created_at"},
            Actions: []admin.Action[models.Post]{
                admin.NewAction("publish", "Publish", func(ctx context.Context, posts []*models.Post) error {
                    for _, post := range posts {
                        post.Status = "published"
                        postManager.Update(ctx, post)
                    }
                    return nil
                }),
            },
        },
    )

    // 6. Register admin routes on existing server (NEW)
    admin.RegisterAdmin(srv, adminRegistry)

    // 7. Start server (EXISTING)
    srv.Start()
}
```

**Result**: Admin panel available at `/admin/` with:
- Auto-generated UI from schemas
- CRUD operations using existing ORM
- Filtering using existing filter system
- API using existing server/routing
- Zero duplication of business logic

---

**Document Version**: 1.0
**Companion to**: ADMIN_REDESIGN_ARCHITECTURE.md
**Last Updated**: 2026-01-01
