package api

import (
	"net/http"
	
	forgeapi "github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/api/authentication"
	"github.com/forgego/forge/pkg/api/filters"
	"github.com/forgego/forge/pkg/api/permissions"
	"github.com/forgego/forge/pkg/api/renderers"
	"github.com/forgego/forge/pkg/api/throttling"
	"ecommerce/models"
)

// ============================================================================
// SIMPLE SCENARIOS - Basic CRUD with minimal configuration
// ============================================================================

// RegisterSimpleAPIViewsets registers simple viewsets with basic functionality
// These demonstrate the simplest way to use the API framework
func RegisterSimpleAPIViewsets(apiRouter *forgeapi.Router) {
	// Simple public read-only endpoints (no auth required)
	registerSimpleBrandAPI(apiRouter)
	registerSimpleCategoryAPI(apiRouter)
	registerSimpleProductAPI(apiRouter)
}

// Simple Brand API - Public read-only
func registerSimpleBrandAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewBrandSerializer,
		nil, // Will be set to models.Brand.Objects.Filter() after generation
		&models.Brand{},
	)
	
	// Simple configuration: Allow anyone to read, no auth required
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}
	
	apiRouter.Register("brands", viewset)
}

// Simple Category API - Public read-only
func registerSimpleCategoryAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewCategorySerializer,
		nil, // Will be set to models.Category.Objects.Filter() after generation
		&models.Category{},
	)
	
	// Simple: Allow anyone to read
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}
	
	apiRouter.Register("categories", viewset)
}

// Simple Product API - Public read-only with basic search
func registerSimpleProductAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewProductSerializer,
		nil, // Will be set to models.Product.Objects.Filter() after generation
		&models.Product{},
	)
	
	// Simple: Allow anyone to read
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}
	
	apiRouter.Register("products", viewset)
}

// ============================================================================
// COMPLEX SCENARIOS - Full-featured APIs with all framework capabilities
// ============================================================================

// RegisterComplexAPIViewsets registers complex viewsets with full features
// These demonstrate all capabilities: auth, permissions, throttling, filters, etc.
func RegisterComplexAPIViewsets(apiRouter *forgeapi.EnhancedRouter) {
	// Complex authenticated endpoints
	registerComplexCustomerAPI(apiRouter)
	registerComplexOrderAPI(apiRouter)
	registerComplexPaymentAPI(apiRouter)
	registerComplexReviewAPI(apiRouter)
	registerComplexInventoryAPI(apiRouter)
}

// Complex Customer API - Full authentication, permissions, throttling
func registerComplexCustomerAPI(apiRouter *forgeapi.EnhancedRouter) {
	// Create fully integrated viewset
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewCustomerSerializer,
		nil, // Will be set to models.Customer.Objects.Filter() after generation
		&models.Customer{},
	)
	
	// Authentication: Token + JWT
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			// Lookup customer by token
			// In real app: return models.Customer.Objects.GetByToken(ctx, token)
			return nil, nil
		}),
		// JWT authentication would go here
	}
	
	// Permissions: Authenticated users can read, only themselves can update
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
		// Add IsOwnerOrReadOnly for update operations
	}
	
	// Throttling: Different rates for authenticated vs anonymous
	cache := throttling.NewMemoryCache()
	viewset.ThrottleClasses = []throttling.Throttle{
		throttling.NewAnonRateThrottle("50/hour", cache),
		throttling.NewUserRateThrottle("1000/day", cache),
	}
	
	// Content negotiation: JSON, XML, HTML
	viewset.RendererClasses = []renderers.Renderer{
		renderers.NewJSONRenderer(),
		renderers.NewXMLRenderer(),
		renderers.NewHTMLRenderer(),
	}
	
	// Filters: Search and ordering
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"email", "first_name", "last_name", "phone"}),
		filters.NewOrderingFilter([]string{"created_at", "email", "last_name", "total_orders"}),
	}
	
	apiRouter.RegisterEnhanced("customers", viewset)
}

// Complex Order API - Full business logic with permissions
func registerComplexOrderAPI(apiRouter *forgeapi.EnhancedRouter) {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewOrderSerializer,
		nil, // Will be set to models.Order.Objects.Filter() after generation
		&models.Order{},
	)
	
	// Authentication required
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}
	
	// Permissions: Customers can only see their own orders
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
		// IsOwnerOrReadOnly would filter by customer_id
	}
	
	// Throttling: Prevent order spam
	cache := throttling.NewMemoryCache()
	viewset.ThrottleClasses = []throttling.Throttle{
		throttling.NewUserRateThrottle("10/hour", cache), // Max 10 orders per hour
	}
	
	// Filters: Search by order number, filter by status
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"order_number", "customer_id"}),
		filters.NewOrderingFilter([]string{"placed_at", "total_amount", "status"}),
	}
	
	apiRouter.RegisterEnhanced("orders", viewset)
	
	// Register custom action: Cancel order
	apiRouter.RegisterAction("cancel", &forgeapi.ActionConfig{
		Methods: []string{"POST"},
		Detail:  true,
	}, func(w http.ResponseWriter, r *http.Request) {
		// Custom cancel order logic
		// This would be implemented in a custom viewset method
		http.Error(w, "Cancel order endpoint - implement business logic here", http.StatusNotImplemented)
	})
}

// Complex Payment API - Secure with strict permissions
func registerComplexPaymentAPI(apiRouter *forgeapi.EnhancedRouter) {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewPaymentSerializer,
		nil, // Will be set to models.Payment.Objects.Filter() after generation
		&models.Payment{},
	)
	
	// Strong authentication
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}
	
	// Strict permissions: Only authenticated users, read-only for most
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
		// Payments are typically read-only for customers
	}
	
	// Aggressive throttling for payment endpoints
	cache := throttling.NewMemoryCache()
	viewset.ThrottleClasses = []throttling.Throttle{
		throttling.NewUserRateThrottle("5/hour", cache), // Very strict
	}
	
	apiRouter.RegisterEnhanced("payments", viewset)
}

// Complex Review API - Public read, authenticated write
func registerComplexReviewAPI(apiRouter *forgeapi.EnhancedRouter) {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewReviewSerializer,
		nil, // Will be set to models.Review.Objects.Filter() after generation
		&models.Review{},
	)
	
	// Authentication for writes only
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}
	
	// Permissions: Read for all, write for authenticated
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticatedOrReadOnly(),
	}
	
	// Moderate throttling
	cache := throttling.NewMemoryCache()
	viewset.ThrottleClasses = []throttling.Throttle{
		throttling.NewAnonRateThrottle("100/hour", cache),
		throttling.NewUserRateThrottle("20/hour", cache), // Prevent review spam
	}
	
	// Filters: Search reviews, filter by rating
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"title", "comment"}),
		filters.NewOrderingFilter([]string{"created_at", "rating"}),
	}
	
	apiRouter.RegisterEnhanced("reviews", viewset)
}

// Complex Inventory API - Admin-only with advanced features
func registerComplexInventoryAPI(apiRouter *forgeapi.EnhancedRouter) {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewInventorySerializer,
		nil, // Will be set to models.Inventory.Objects.Filter() after generation
		&models.Inventory{},
	)
	
	// Authentication required
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}
	
	// Admin-only permissions
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
		permissions.NewIsAdminUser(),
	}
	
	// Filters: Search by product, filter by warehouse
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"product_id", "warehouse_id"}),
		filters.NewOrderingFilter([]string{"quantity", "available_quantity"}),
	}
	
	apiRouter.RegisterEnhanced("inventory", viewset)
}

// ============================================================================
// LEGACY - Original simple registration (backward compatible)
// ============================================================================

// RegisterAPIViewsets registers all REST API viewsets (legacy, uses simple viewsets)
func RegisterAPIViewsets(apiRouter *forgeapi.Router) {
	// Use simple registration for backward compatibility
	RegisterSimpleAPIViewsets(apiRouter)
	
	// Also register remaining endpoints with simple configuration
	registerAddressAPI(apiRouter)
	registerSupplierAPI(apiRouter)
	registerProductVariantAPI(apiRouter)
	registerWarehouseAPI(apiRouter)
	registerOrderItemAPI(apiRouter)
	registerShippingAPI(apiRouter)
}

func registerAddressAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewAddressSerializer,
		nil,
		&models.Address{},
	)
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
	}
	apiRouter.Register("addresses", viewset)
}

func registerSupplierAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewSupplierSerializer,
		nil,
		&models.Supplier{},
	)
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}
	apiRouter.Register("suppliers", viewset)
}

func registerProductVariantAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewProductVariantSerializer,
		nil,
		&models.ProductVariant{},
	)
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}
	apiRouter.Register("product-variants", viewset)
}

func registerWarehouseAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewWarehouseSerializer,
		nil,
		&models.Warehouse{},
	)
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
	}
	apiRouter.Register("warehouses", viewset)
}

func registerOrderItemAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewOrderItemSerializer,
		nil,
		&models.OrderItem{},
	)
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
	}
	apiRouter.Register("order-items", viewset)
}

func registerShippingAPI(apiRouter *forgeapi.Router) {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewShippingSerializer,
		nil,
		&models.Shipping{},
	)
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
	}
	apiRouter.Register("shipping", viewset)
}
