# Ecommerce API Examples - Simple & Complex Scenarios

This document demonstrates both **simple** and **complex** API scenarios using the Forge API framework.

## 📚 Table of Contents

1. [Simple Scenarios](#simple-scenarios)
2. [Complex Scenarios](#complex-scenarios)
3. [Usage Examples](#usage-examples)
4. [API Endpoints](#api-endpoints)

## 🎯 Simple Scenarios

Simple scenarios demonstrate basic CRUD operations with minimal configuration. Perfect for:
- Learning the framework
- Quick prototypes
- Public read-only APIs
- Internal tools

### Example 1: Public Read-Only API

```go
// Simplest possible API - no authentication, no permissions
viewset := forgeapi.NewEnhancedBaseViewSet(
    NewBrandSerializer,
    models.Brand.Objects.Filter(),
    &models.Brand{},
)

// That's it! Default allows all access
apiRouter.Register("brands", viewset)
```

**Features:**
- ✅ Basic CRUD operations
- ✅ Pagination
- ✅ No authentication required
- ✅ No permissions
- ✅ No throttling

**Endpoints:**
- `GET /api/v1/brands/` - List all brands
- `GET /api/v1/brands/{id}/` - Get brand details

### Example 2: Authenticated API

```go
viewset := forgeapi.NewEnhancedBaseViewSet(
    NewCustomerSerializer,
    models.Customer.Objects.Filter(),
    &models.Customer{},
)

// Add authentication
viewset.AuthenticationClasses = []authentication.Authentication{
    authentication.NewTokenAuthentication(lookupUserByToken),
}

// Require authentication
viewset.PermissionClasses = []permissions.Permission{
    permissions.NewIsAuthenticated(),
}
```

**Features:**
- ✅ Authentication required
- ✅ Token-based auth
- ✅ All CRUD operations

### Example 3: Read-Only Public, Write Requires Auth

```go
viewset := forgeapi.NewEnhancedBaseViewSet(
    NewProductSerializer,
    models.Product.Objects.Filter(),
    &models.Product{},
)

viewset.AuthenticationClasses = []authentication.Authentication{
    authentication.NewTokenAuthentication(lookupUserByToken),
}

// Read for all, write for authenticated
viewset.PermissionClasses = []permissions.Permission{
    permissions.NewIsAuthenticatedOrReadOnly(),
}
```

**Features:**
- ✅ Public read access
- ✅ Authenticated write access
- ✅ Common pattern for content APIs

## 🚀 Complex Scenarios

Complex scenarios demonstrate the full power of the framework with:
- Multiple authentication methods
- Advanced permissions
- Rate limiting
- Content negotiation
- Advanced filtering
- Custom actions

### Example 1: Full-Featured Customer API

```go
viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
    NewCustomerSerializer,
    models.Customer.Objects.Filter(),
    &models.Customer{},
)

// Multiple authentication methods
viewset.AuthenticationClasses = []authentication.Authentication{
    authentication.NewTokenAuthentication(lookupUserByToken),
    authentication.NewJWTAuthentication(secretKey, lookupUserFromJWT),
}

// Permissions: Authenticated users only
viewset.PermissionClasses = []permissions.Permission{
    permissions.NewIsAuthenticated(),
}

// Throttling: Different rates for different user types
cache := throttling.NewMemoryCache()
viewset.ThrottleClasses = []throttling.Throttle{
    throttling.NewAnonRateThrottle("50/hour", cache),
    throttling.NewUserRateThrottle("1000/day", cache),
}

// Content negotiation: Multiple formats
viewset.RendererClasses = []renderers.Renderer{
    renderers.NewJSONRenderer(),
    renderers.NewXMLRenderer(),
    renderers.NewHTMLRenderer(),
}

// Advanced filtering
viewset.FilterBackends = []filters.FilterBackend{
    filters.NewSearchFilter([]string{"email", "first_name", "last_name"}),
    filters.NewOrderingFilter([]string{"created_at", "email"}),
}
```

**Features:**
- ✅ Multiple authentication methods
- ✅ Rate limiting
- ✅ Content negotiation (JSON/XML/HTML)
- ✅ Search and ordering
- ✅ Full CRUD operations

**Usage:**
```bash
# Get JSON response
curl -H "Accept: application/json" /api/v1/customers/

# Get XML response
curl -H "Accept: application/xml" /api/v1/customers/

# Search customers
curl /api/v1/customers/?search=john

# Order by email
curl /api/v1/customers/?ordering=email
```

### Example 2: Secure Order API

```go
viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
    NewOrderSerializer,
    models.Order.Objects.Filter(),
    &models.Order{},
)

// Strong authentication
viewset.AuthenticationClasses = []authentication.Authentication{
    authentication.NewTokenAuthentication(lookupUserByToken),
}

// Permissions: Users can only see their own orders
viewset.PermissionClasses = []permissions.Permission{
    permissions.NewIsAuthenticated(),
    // IsOwnerOrReadOnly would filter by customer_id
}

// Aggressive throttling to prevent abuse
cache := throttling.NewMemoryCache()
viewset.ThrottleClasses = []throttling.Throttle{
    throttling.NewUserRateThrottle("10/hour", cache), // Max 10 orders/hour
}

// Filters
viewset.FilterBackends = []filters.FilterBackend{
    filters.NewSearchFilter([]string{"order_number"}),
    filters.NewOrderingFilter([]string{"placed_at", "total_amount"}),
}
```

**Features:**
- ✅ Strict authentication
- ✅ Owner-based permissions
- ✅ Rate limiting (10 orders/hour)
- ✅ Search and filtering

### Example 3: Admin-Only Inventory API

```go
viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
    NewInventorySerializer,
    models.Inventory.Objects.Filter(),
    &models.Inventory{},
)

viewset.AuthenticationClasses = []authentication.Authentication{
    authentication.NewTokenAuthentication(lookupUserByToken),
}

// Admin-only access
viewset.PermissionClasses = []permissions.Permission{
    permissions.NewIsAuthenticated(),
    permissions.NewIsAdminUser(),
}

// Advanced filtering for inventory management
viewset.FilterBackends = []filters.FilterBackend{
    filters.NewSearchFilter([]string{"product_id", "warehouse_id"}),
    filters.NewOrderingFilter([]string{"quantity", "available_quantity"}),
}
```

**Features:**
- ✅ Admin-only access
- ✅ Advanced filtering
- ✅ Inventory management

### Example 4: Review API with Custom Actions

```go
viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
    NewReviewSerializer,
    models.Review.Objects.Filter(),
    &models.Review{},
)

// Public read, authenticated write
viewset.PermissionClasses = []permissions.Permission{
    permissions.NewIsAuthenticatedOrReadOnly(),
}

// Moderate throttling to prevent spam
cache := throttling.NewMemoryCache()
viewset.ThrottleClasses = []throttling.Throttle{
    throttling.NewAnonRateThrottle("100/hour", cache),
    throttling.NewUserRateThrottle("20/hour", cache), // Prevent review spam
}

// Register custom actions
apiRouter.RegisterAction("approve", &forgeapi.ActionConfig{
    Methods: []string{"POST"},
    Detail:  true,
}, approveReviewHandler)

apiRouter.RegisterAction("report", &forgeapi.ActionConfig{
    Methods: []string{"POST"},
    Detail:  true,
}, reportReviewHandler)
```

**Features:**
- ✅ Public read access
- ✅ Authenticated write
- ✅ Rate limiting
- ✅ Custom actions (approve, report)

**Endpoints:**
- `GET /api/v1/reviews/` - List reviews
- `POST /api/v1/reviews/{id}/approve/` - Approve review (custom action)
- `POST /api/v1/reviews/{id}/report/` - Report review (custom action)

## 📊 Comparison Table

| Feature | Simple | Complex |
|---------|--------|---------|
| Authentication | ❌ | ✅ Multiple methods |
| Permissions | Basic | Advanced (owner-based, admin) |
| Throttling | ❌ | ✅ Rate limiting |
| Content Negotiation | ❌ | ✅ JSON/XML/HTML |
| Filtering | Basic | Advanced (search, ordering) |
| Custom Actions | ❌ | ✅ |
| Error Handling | Basic | Comprehensive |
| Security | Basic | Production-ready |

## 🎓 Learning Path

1. **Start Simple**: Use `Example1_SimplePublicAPI` to understand basics
2. **Add Auth**: Move to `Example2_AuthenticatedAPI` 
3. **Add Permissions**: Try `Example3_ReadOnlyPublic`
4. **Add Throttling**: Use `Example4_WithThrottling`
5. **Add Filters**: Try `Example5_WithFilters`
6. **Go Full**: Use `Example7_Complete` for production

## 🔗 API Endpoints

### Simple APIs

- `GET /api/v1/brands/` - List brands (public)
- `GET /api/v1/categories/` - List categories (public)
- `GET /api/v1/products/` - List products (public)

### Complex APIs

- `GET /api/v1/customers/` - List customers (authenticated)
- `GET /api/v1/orders/` - List orders (authenticated, owner-only)
- `GET /api/v1/payments/` - List payments (authenticated, read-only)
- `GET /api/v1/reviews/` - List reviews (public read, auth write)
- `GET /api/v1/inventory/` - List inventory (admin-only)

### Custom Actions

- `POST /api/v1/orders/{id}/cancel/` - Cancel order
- `POST /api/v1/reviews/{id}/approve/` - Approve review
- `POST /api/v1/reviews/{id}/report/` - Report review

## 💡 Best Practices

### Simple Scenarios
- Use for public data
- Use for internal tools
- Use for learning
- Keep it minimal

### Complex Scenarios
- Use for production APIs
- Use for sensitive data
- Use for high-traffic endpoints
- Use for business-critical operations

## 🚀 Quick Start

```go
// Simple
simpleRouter := forgeapi.NewRouter("/api/v1")
api.RegisterSimpleAPIViewsets(simpleRouter)

// Complex
complexRouter := forgeapi.NewEnhancedRouter("/api/v1")
api.RegisterComplexAPIViewsets(complexRouter)
```

---

For more examples, see `api/examples.go` in the ecommerce directory.
