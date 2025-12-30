package api

import (
	"fmt"
	"net/http"

	"ecommerce/models"
	"github.com/forgego/forge/api"
	forgeapi "github.com/forgego/forge/api"
)

// ============================================================================
// EXAMPLE SCENARIOS - Demonstrating different use cases
// ============================================================================

// Example1_SimplePublicAPI demonstrates the simplest possible API setup
// No authentication, no permissions, just basic CRUD
func Example1_SimplePublicAPI() *forgeapi.EnhancedBaseViewSet {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewBrandSerializer,
		nil,
		&models.Brand{},
	)

	// That's it! Default permissions allow all access
	return viewset
}

// Example2_AuthenticatedAPI demonstrates basic authentication
// Users must be authenticated to access
func Example2_AuthenticatedAPI() *forgeapi.EnhancedBaseViewSet {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewCustomerSerializer,
		nil,
		&models.Customer{},
	)

	// Add authentication
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			// Lookup user by token
			return nil, nil
		}),
	}

	// Require authentication
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
	}

	return viewset
}

// Example3_ReadOnlyPublic demonstrates public read, authenticated write
// Common pattern for content APIs
func Example3_ReadOnlyPublic() *forgeapi.EnhancedBaseViewSet {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewProductSerializer,
		nil,
		&models.Product{},
	)

	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}

	// Read for all, write for authenticated
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticatedOrReadOnly(),
	}

	return viewset
}

// Example4_WithThrottling demonstrates rate limiting
// Prevents abuse and ensures fair usage
func Example4_WithThrottling() *forgeapi.EnhancedBaseViewSet {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewOrderSerializer,
		nil,
		&models.Order{},
	)

	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}

	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
	}

	// Add throttling
	cache := throttling.NewMemoryCache()
	viewset.ThrottleClasses = []throttling.Throttle{
		throttling.NewAnonRateThrottle("100/hour", cache),
		throttling.NewUserRateThrottle("1000/day", cache),
	}

	return viewset
}

// Example5_WithFilters demonstrates search and ordering
// Makes APIs more useful for clients
func Example5_WithFilters() *forgeapi.EnhancedBaseViewSetIntegrated {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewProductSerializer,
		nil,
		&models.Product{},
	)

	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}

	// Add filters
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"name", "description", "sku"}),
		filters.NewOrderingFilter([]string{"price", "created_at", "name"}),
	}

	return viewset
}

// Example6_WithContentNegotiation demonstrates multiple formats
// Supports JSON, XML, HTML responses
func Example6_WithContentNegotiation() *forgeapi.EnhancedBaseViewSetIntegrated {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewProductSerializer,
		nil,
		&models.Product{},
	)

	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewAllowAny(),
	}

	// Multiple renderers
	viewset.RendererClasses = []renderers.Renderer{
		renderers.NewJSONRenderer(),
		renderers.NewXMLRenderer(),
		renderers.NewHTMLRenderer(),
	}

	return viewset
}

// Example7_Complete demonstrates all features together
// Production-ready configuration
func Example7_Complete() *forgeapi.EnhancedBaseViewSetIntegrated {
	viewset := forgeapi.NewEnhancedBaseViewSetIntegrated(
		NewOrderSerializer,
		nil,
		&models.Order{},
	)

	// Authentication
	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}

	// Permissions
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
		// Add IsOwnerOrReadOnly for customer-specific filtering
	}

	// Throttling
	cache := throttling.NewMemoryCache()
	viewset.ThrottleClasses = []throttling.Throttle{
		throttling.NewUserRateThrottle("100/hour", cache),
	}

	// Filters
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"order_number"}),
		filters.NewOrderingFilter([]string{"placed_at", "total_amount"}),
	}

	// Content negotiation
	viewset.RendererClasses = []renderers.Renderer{
		renderers.NewJSONRenderer(),
		renderers.NewXMLRenderer(),
	}

	return viewset
}

// Example8_AdminOnly demonstrates admin-only access
// For sensitive operations
func Example8_AdminOnly() *forgeapi.EnhancedBaseViewSet {
	viewset := forgeapi.NewEnhancedBaseViewSet(
		NewInventorySerializer,
		nil,
		&models.Inventory{},
	)

	viewset.AuthenticationClasses = []authentication.Authentication{
		authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
			return nil, nil
		}),
	}

	// Admin only
	viewset.PermissionClasses = []permissions.Permission{
		permissions.NewIsAuthenticated(),
		permissions.NewIsAdminUser(),
	}

	return viewset
}

// Example9_CustomAction demonstrates custom actions
// Beyond standard CRUD operations
func Example9_CustomAction(apiRouter *forgeapi.EnhancedRouter) {
	// Register a custom action on orders
	apiRouter.RegisterAction("cancel", &forgeapi.ActionConfig{
		Methods: []string{"POST"},
		Detail:  true, // Requires order ID
	}, func(w http.ResponseWriter, r *http.Request) {
		// Custom cancel order logic
		orderID := r.URL.Query().Get("id")
		fmt.Fprintf(w, "Cancelling order %s", orderID)
	})

	// Register another custom action
	apiRouter.RegisterAction("refund", &forgeapi.ActionConfig{
		Methods: []string{"POST"},
		Detail:  true,
	}, func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")
		fmt.Fprintf(w, "Processing refund for order %s", orderID)
	})
}

// Example10_ProductionReady demonstrates production setup
// Uses CreateProductionViewSet helper
func Example10_ProductionReady() *forgeapi.EnhancedBaseViewSetIntegrated {
	// This helper applies all best practices
	viewset := forgeapi.CreateProductionViewSet(
		NewOrderSerializer,
		nil,
		&models.Order{},
	)

	// Additional customizations can be added
	viewset.FilterBackends = []filters.FilterBackend{
		filters.NewSearchFilter([]string{"order_number"}),
	}

	return viewset
}
