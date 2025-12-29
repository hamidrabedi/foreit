# Ecommerce API Examples - Updated

This ecommerce example now includes both **simple** and **complex** API scenarios demonstrating the full capabilities of the Forge API framework.

## 🎯 What's New

### Simple Scenarios
- ✅ Basic CRUD operations
- ✅ Minimal configuration
- ✅ Public read-only APIs
- ✅ Quick setup examples

### Complex Scenarios  
- ✅ Full authentication (Token + JWT)
- ✅ Advanced permissions (Owner-based, Admin-only)
- ✅ Rate limiting/throttling
- ✅ Content negotiation (JSON/XML/HTML)
- ✅ Advanced filtering (Search, Ordering)
- ✅ Custom actions
- ✅ Production-ready configurations

## 📁 File Structure

```
examples/ecommerce/api/
├── serializers.go      # Serializers for all models
├── viewsets.go         # Simple & Complex viewsets
└── examples.go         # 10 example scenarios
```

## 🚀 Quick Start

### Simple API Setup

```go
// Basic public read-only API
simpleRouter := forgeapi.NewRouter("/api/v1")
api.RegisterSimpleAPIViewsets(simpleRouter)
```

**Endpoints:**
- `GET /api/v1/brands/` - Public brands
- `GET /api/v1/categories/` - Public categories  
- `GET /api/v1/products/` - Public products

### Complex API Setup

```go
// Full-featured API with all capabilities
complexRouter := forgeapi.NewEnhancedRouter("/api/v1")
api.RegisterComplexAPIViewsets(complexRouter)
```

**Endpoints:**
- `GET /api/v1/customers/` - Authenticated customers
- `GET /api/v1/orders/` - Owner-only orders
- `GET /api/v1/payments/` - Secure payments
- `GET /api/v1/reviews/` - Public read, auth write
- `GET /api/v1/inventory/` - Admin-only inventory

## 📚 Example Scenarios

See `api/examples.go` for 10 complete examples:

1. **SimplePublicAPI** - Basic CRUD, no auth
2. **AuthenticatedAPI** - Token authentication
3. **ReadOnlyPublic** - Public read, auth write
4. **WithThrottling** - Rate limiting
5. **WithFilters** - Search and ordering
6. **WithContentNegotiation** - Multiple formats
7. **Complete** - All features together
8. **AdminOnly** - Admin-only access
9. **CustomAction** - Custom endpoints
10. **ProductionReady** - Best practices

## 📖 Documentation

See `API_EXAMPLES.md` for detailed documentation with:
- Code examples
- Feature comparisons
- Usage patterns
- Best practices
- Learning path

## 🔧 Configuration

The server (`cmd/server/main.go`) now registers:
- Simple APIs at `/api/v1/`
- Complex APIs at `/api/v1/` (enhanced router)
- Legacy APIs at `/api/v1/legacy/` (backward compatible)

## ✨ Features Demonstrated

### Simple Scenarios
- Basic ViewSet usage
- Public access
- Minimal configuration

### Complex Scenarios
- Multiple authentication methods
- Advanced permissions
- Rate limiting
- Content negotiation
- Advanced filtering
- Custom actions
- Error handling
- Security best practices

---

For detailed examples, see `API_EXAMPLES.md` and `api/examples.go`.
