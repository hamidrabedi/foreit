# API Framework Architecture & Design Patterns

## Overview

This document describes the architecture, design patterns, and structure of the API framework. The design follows DRF principles while leveraging Go's strengths for performance and type safety.

## Architecture Principles

### 1. **Separation of Concerns**
- Each component has a single, well-defined responsibility
- Clear boundaries between layers
- Easy to test and maintain

### 2. **Strategy Pattern First**
- Pluggable components (auth, permissions, throttling, filters)
- Easy to extend and customize
- Default implementations provided

### 3. **Type Safety**
- Leverage Go's type system
- Compile-time validation where possible
- Clear interfaces and contracts

### 4. **Performance by Default**
- Zero-allocation paths where possible
- Efficient serialization
- Connection pooling
- Caching support

### 5. **Developer Experience**
- Simple APIs for common cases
- Powerful APIs for advanced cases
- Clear error messages
- Excellent documentation

## Directory Structure

```
pkg/api/
├── core/                    # Core framework components
│   ├── request.go          # Request wrapper with context
│   ├── response.go         # Response wrapper
│   ├── context.go          # Request context utilities
│   └── middleware.go       # Core middleware interface
│
├── serializers/             # Serialization system
│   ├── serializer.go        # Base serializer interface
│   ├── fields/              # Field types
│   │   ├── base.go         # Base field
│   │   ├── string.go       # StringField
│   │   ├── integer.go      # IntegerField
│   │   ├── float.go        # FloatField
│   │   ├── boolean.go      # BooleanField
│   │   ├── datetime.go     # DateTimeField
│   │   ├── email.go        # EmailField
│   │   ├── url.go          # URLField
│   │   ├── uuid.go         # UUIDField
│   │   ├── decimal.go      # DecimalField
│   │   ├── json.go         # JSONField
│   │   ├── file.go         # FileField
│   │   ├── image.go        # ImageField
│   │   ├── relation.go     # Relation fields
│   │   └── method.go       # SerializerMethodField
│   ├── validation.go       # Validation logic
│   ├── model_serializer.go # ModelSerializer
│   └── nested.go           # Nested serializers
│
├── viewsets/               # ViewSet system
│   ├── viewset.go          # Base ViewSet interface
│   ├── base.go             # BaseViewSet implementation
│   ├── mixins.go           # ViewSet mixins
│   ├── actions.go          # Custom action decorator
│   └── hooks.go            # ViewSet hooks
│
├── authentication/         # Authentication system
│   ├── auth.go             # Base authentication interface
│   ├── token.go            # TokenAuthentication
│   ├── jwt.go             # JWTAuthentication
│   ├── session.go          # SessionAuthentication
│   ├── basic.go           # BasicAuthentication
│   ├── apikey.go          # APIKeyAuthentication
│   └── result.go          # Authentication result
│
├── permissions/            # Permission system
│   ├── permission.go      # Base permission interface
│   ├── allow_any.go      # AllowAny
│   ├── is_authenticated.go # IsAuthenticated
│   ├── is_admin.go       # IsAdminUser
│   ├── is_owner.go       # IsOwnerOrReadOnly
│   ├── django_model.go   # DjangoModelPermissions
│   └── composite.go      # Permission composition
│
├── throttling/            # Throttling system
│   ├── throttle.go       # Base throttle interface
│   ├── anon_rate.go      # AnonRateThrottle
│   ├── user_rate.go      # UserRateThrottle
│   ├── scoped_rate.go    # ScopedRateThrottle
│   └── cache.go          # Throttle cache backend
│
├── renderers/             # Response renderers
│   ├── renderer.go        # Base renderer interface
│   ├── json.go            # JSONRenderer
│   ├── xml.go             # XMLRenderer
│   ├── yaml.go            # YAMLRenderer
│   ├── html.go            # HTMLRenderer (browsable API)
│   ├── csv.go             # CSVRenderer
│   └── negotiation.go     # Content negotiation
│
├── parsers/               # Request parsers
│   ├── parser.go          # Base parser interface
│   ├── json.go            # JSONParser
│   ├── xml.go             # XMLParser
│   ├── form.go            # FormParser
│   ├── multipart.go       # MultiPartParser
│   └── yaml.go            # YAMLParser
│
├── filters/               # Filtering system
│   ├── backend.go         # Filter backend interface
│   ├── django_filter.go   # DjangoFilterBackend
│   ├── search.go          # SearchFilter
│   ├── ordering.go        # OrderingFilter
│   └── filterset.go       # FilterSet (existing)
│
├── exceptions/            # Exception handling
│   ├── exception.go       # Base exception
│   ├── validation.go      # ValidationError
│   ├── auth.go            # Authentication exceptions
│   ├── permission.go      # Permission exceptions
│   ├── not_found.go       # NotFound
│   ├── throttled.go       # Throttled
│   └── handler.go         # Exception handler
│
├── versioning/            # API versioning
│   ├── version.go          # Version interface
│   ├── url_path.go        # URLPathVersioning
│   ├── query_param.go     # QueryParameterVersioning
│   ├── header.go          # HeaderVersioning
│   └── namespace.go       # NamespaceVersioning
│
├── caching/               # Caching system
│   ├── backend.go         # Cache backend interface
│   ├── memory.go          # MemoryCache
│   ├── redis.go           # RedisCache
│   └── key.go             # Cache key generation
│
├── pagination/            # Pagination (existing)
│   ├── pagination.go      # Pagination types
│   └── cursor.go          # Cursor pagination (future)
│
├── router/                # Routing system
│   ├── router.go          # API Router
│   ├── route.go           # Route definition
│   └── registry.go        # Route registry
│
├── metadata/              # API metadata
│   ├── metadata.go        # Metadata interface
│   └── options.go         # OPTIONS handler
│
├── testing/                # Testing utilities
│   ├── client.go          # APIClient
│   ├── factory.go         # Test factories
│   └── assertions.go     # Test assertions
│
├── docs/                  # Documentation generation
│   ├── openapi.go         # OpenAPI generator
│   ├── schema.go          # Schema generator
│   └── swagger.go         # Swagger UI
│
└── config/                # Configuration
    ├── settings.go        # API settings
    └── defaults.go        # Default configurations
```

## Design Patterns

### 1. Strategy Pattern

**Purpose:** Allow interchangeable algorithms for authentication, permissions, throttling, filters, renderers, parsers.

**Implementation:**

```go
// Authentication Strategy
type Authentication interface {
    Authenticate(r *http.Request) (*AuthResult, error)
    AuthenticateHeader(r *http.Request) string
}

// Permission Strategy
type Permission interface {
    HasPermission(r *http.Request, view ViewSet) bool
    HasObjectPermission(r *http.Request, view ViewSet, obj interface{}) bool
}

// Throttle Strategy
type Throttle interface {
    AllowRequest(r *http.Request, view ViewSet) (bool, time.Duration, error)
    GetScope(r *http.Request, view ViewSet) string
}
```

**Benefits:**
- Easy to add new authentication methods
- Easy to customize permissions
- Pluggable components

### 2. Chain of Responsibility Pattern

**Purpose:** Process requests through a chain of middleware, authentication, permissions, throttling.

**Implementation:**

```go
type MiddlewareChain struct {
    handlers []Middleware
}

func (c *MiddlewareChain) Process(r *http.Request, next http.Handler) http.Handler {
    // Process in reverse order (last added, first executed)
    for i := len(c.handlers) - 1; i >= 0; i-- {
        next = c.handlers[i](next)
    }
    return next
}

// Request processing flow:
// 1. Authentication
// 2. Permissions
// 3. Throttling
// 4. Parsing
// 5. ViewSet action
// 6. Serialization
// 7. Rendering
```

**Benefits:**
- Flexible request processing
- Easy to add/remove steps
- Clear execution order

### 3. Template Method Pattern

**Purpose:** Define skeleton of ViewSet operations, allow customization.

**Implementation:**

```go
type BaseViewSet struct {
    // Template methods
    GetQueryset() interface{}
    GetSerializer() Serializer
    GetPermissions() []Permission
    GetThrottles() []Throttle
    
    // Hooks for customization
    PerformCreate(obj interface{}) error
    PerformUpdate(obj interface{}) error
    PerformDestroy(obj interface{}) error
}

// Default implementation
func (vs *BaseViewSet) Create(w http.ResponseWriter, r *http.Request) {
    // 1. Authenticate
    // 2. Check permissions
    // 3. Check throttles
    // 4. Parse request
    // 5. Validate
    // 6. PerformCreate (hook)
    // 7. Serialize
    // 8. Render
}
```

**Benefits:**
- Consistent behavior
- Easy to override specific steps
- Reduces code duplication

### 4. Factory Pattern

**Purpose:** Create serializers, renderers, parsers dynamically.

**Implementation:**

```go
type SerializerFactory interface {
    CreateSerializer(data interface{}) Serializer
}

type RendererFactory interface {
    CreateRenderer(format string) Renderer
}

// Default factories
func NewSerializerFactory() SerializerFactory {
    return &defaultSerializerFactory{}
}
```

**Benefits:**
- Centralized creation logic
- Easy to swap implementations
- Supports dependency injection

### 5. Builder Pattern

**Purpose:** Build complex objects (queries, serializers, responses) step by step.

**Implementation:**

```go
type QueryBuilder struct {
    queryset interface{}
    filters  []Filter
    ordering []string
    limit    *int
    offset   *int
}

func (b *QueryBuilder) Filter(f Filter) *QueryBuilder {
    b.filters = append(b.filters, f)
    return b
}

func (b *QueryBuilder) OrderBy(fields ...string) *QueryBuilder {
    b.ordering = append(b.ordering, fields...)
    return b
}

func (b *QueryBuilder) Build() interface{} {
    // Build final query
}
```

**Benefits:**
- Fluent API
- Easy to construct complex queries
- Immutable operations

### 6. Decorator Pattern

**Purpose:** Add functionality to ViewSets (actions, middleware, caching).

**Implementation:**

```go
// Action decorator
func Action(methods []string, detail bool) func(ViewSet) ViewSet {
    return func(vs ViewSet) ViewSet {
        // Register action
        return vs
    }
}

// Usage
type PostViewSet struct {
    *BaseViewSet
}

@Action(methods=["POST"], detail=true)
func (vs *PostViewSet) Publish(w http.ResponseWriter, r *http.Request) {
    // Custom action
}
```

**Benefits:**
- Add functionality without modifying base
- Composable features
- Clean separation

### 7. Observer Pattern

**Purpose:** Event system for hooks (before_create, after_create, etc.).

**Implementation:**

```go
type EventType string

const (
    EventBeforeCreate EventType = "before_create"
    EventAfterCreate  EventType = "after_create"
    EventBeforeUpdate EventType = "before_update"
    EventAfterUpdate  EventType = "after_update"
)

type EventObserver interface {
    OnEvent(event EventType, data interface{})
}

type EventEmitter struct {
    observers []EventObserver
}

func (e *EventEmitter) Emit(event EventType, data interface{}) {
    for _, obs := range e.observers {
        obs.OnEvent(event, data)
    }
}
```

**Benefits:**
- Loose coupling
- Easy to add hooks
- Extensible

### 8. Repository Pattern

**Purpose:** Abstract data access (already implemented via QuerySet/Manager).

**Implementation:**

```go
// Already exists in query package
type QuerySet interface {
    Filter(expr QueryExpr) QuerySet
    All(ctx context.Context) ([]interface{}, error)
    Get(ctx context.Context, id int64) (interface{}, error)
    // ...
}
```

**Benefits:**
- Testable
- Swappable backends
- Clean separation

## Request Flow

```
HTTP Request
    ↓
Router (matches route)
    ↓
Middleware Chain
    ├─ Logging
    ├─ CORS
    ├─ Request ID
    └─ ...
    ↓
ViewSet Handler
    ↓
Authentication
    ├─ Try TokenAuth
    ├─ Try JWTAuth
    └─ Try SessionAuth
    ↓
Permissions
    ├─ Check view-level permissions
    └─ Check object-level permissions (if applicable)
    ↓
Throttling
    ├─ Check rate limits
    └─ Return 429 if exceeded
    ↓
Content Negotiation
    ├─ Determine parser (from Content-Type)
    └─ Determine renderer (from Accept)
    ↓
Parse Request
    └─ Parse body using selected parser
    ↓
ViewSet Action
    ├─ Get queryset
    ├─ Apply filters
    ├─ Apply ordering
    ├─ Apply pagination
    ├─ Execute query
    └─ Get serializer
    ↓
Serialize
    ├─ Validate data (if write operation)
    ├─ Transform to internal representation
    └─ Save to database (if write operation)
    ↓
Render Response
    └─ Render using selected renderer
    ↓
HTTP Response
```

## Component Interactions

### ViewSet → Serializer
- ViewSet gets serializer instance
- Passes data to serializer
- Serializer validates and transforms

### ViewSet → Authentication
- ViewSet checks authentication
- Gets user from request context
- Uses user for permissions/queries

### ViewSet → Permissions
- ViewSet checks permissions before action
- Checks object permissions after retrieving object
- Returns 401/403 if denied

### ViewSet → Throttling
- ViewSet checks throttles before action
- Returns 429 if rate limit exceeded
- Includes retry-after header

### Serializer → Validation
- Serializer uses validation package
- Validates field-level rules
- Validates object-level rules

### Filter → QuerySet
- Filter backend applies filters to queryset
- Returns modified queryset
- QuerySet executes SQL

## Configuration

### Global Settings

```go
type APISettings struct {
    // Authentication
    DefaultAuthentication []Authentication
    
    // Permissions
    DefaultPermissions []Permission
    
    // Throttling
    DefaultThrottles []Throttle
    ThrottleRates    map[string]string
    
    // Renderers
    DefaultRenderers []Renderer
    
    // Parsers
    DefaultParsers []Parser
    
    // Pagination
    PageSize int
    
    // Format
    DefaultFormat string
    
    // Exception Handler
    ExceptionHandler ExceptionHandler
}
```

### Per-ViewSet Settings

```go
type ViewSetConfig struct {
    Authentication []Authentication
    Permissions    []Permission
    Throttles      []Throttle
    Renderers      []Renderer
    Parsers         []Parser
    Pagination      PaginationClass
    FilterBackends  []FilterBackend
}
```

## Error Handling

### Exception Hierarchy

```
APIException (base)
├─ ValidationError (400)
├─ ParseError (400)
├─ AuthenticationFailed (401)
├─ NotAuthenticated (401)
├─ PermissionDenied (403)
├─ NotFound (404)
├─ MethodNotAllowed (405)
├─ NotAcceptable (406)
├─ UnsupportedMediaType (415)
└─ Throttled (429)
```

### Exception Handler

```go
type ExceptionHandler interface {
    HandleException(err error, r *http.Request) *ErrorResponse
}

// Default handler converts exceptions to HTTP responses
func DefaultExceptionHandler(err error, r *http.Request) *ErrorResponse {
    switch e := err.(type) {
    case *ValidationError:
        return &ErrorResponse{
            Status:  400,
            Code:    "validation_error",
            Message: "Validation failed",
            Details: e.Errors,
        }
    // ... other cases
    }
}
```

## Performance Considerations

### 1. Serialization
- Use reflection efficiently
- Cache field metadata
- Pool serializer instances
- Zero-allocation paths where possible

### 2. Query Execution
- Connection pooling
- Query caching
- Prepared statements
- Efficient scanning

### 3. Caching
- Response caching
- Query result caching
- Serializer caching
- Metadata caching

### 4. Concurrency
- Goroutine-safe components
- No shared mutable state
- Context propagation
- Proper locking

## Testing Strategy

### Unit Tests
- Test each component in isolation
- Mock dependencies
- Test error cases

### Integration Tests
- Test full request flow
- Test with real database
- Test authentication/permissions

### Performance Tests
- Benchmark serialization
- Benchmark queries
- Load testing

## Extension Points

### 1. Custom Authentication
```go
type CustomAuth struct{}

func (a *CustomAuth) Authenticate(r *http.Request) (*AuthResult, error) {
    // Custom logic
}
```

### 2. Custom Permissions
```go
type CustomPermission struct{}

func (p *CustomPermission) HasPermission(r *http.Request, view ViewSet) bool {
    // Custom logic
}
```

### 3. Custom Renderers
```go
type CustomRenderer struct{}

func (r *CustomRenderer) Render(data interface{}) ([]byte, error) {
    // Custom rendering
}
```

### 4. Custom Actions
```go
func (vs *PostViewSet) CustomAction(w http.ResponseWriter, r *http.Request) {
    // Custom action logic
}
```

## Current Status

✅ **Complete and Production Ready**

All components of the API framework have been implemented:
- ✅ Serializers with all field types
- ✅ ViewSets with CRUD operations
- ✅ Authentication (Token, JWT, Basic, Session, API Key)
- ✅ Permissions (AllowAny, IsAuthenticated, IsAdminUser, IsOwnerOrReadOnly)
- ✅ Throttling (AnonRateThrottle, UserRateThrottle, ScopedRateThrottle)
- ✅ Content negotiation (JSON, XML, YAML, HTML, CSV)
- ✅ Pagination (PageNumber, LimitOffset)
- ✅ Filtering and search
- ✅ Exception handling
- ✅ API versioning
- ✅ Caching backends

See [API Reference](API_REFERENCE.md) for usage examples.

---

**Last Updated:** January 2025
**Status:** ✅ Production Ready
