package main

import (
	"fmt"
	"log"
	"net/http"

	"ecommerce/admin"
	"ecommerce/api"
	"ecommerce/models"

	"github.com/forgego/forge"
	"github.com/forgego/forge/pkg/admin"
	forgeapi "github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/api/authentication"
	"github.com/forgego/forge/pkg/users"
)

func main() {
	// Load configuration
	cfg := forge.NewConfig()
	settings := forge.LoadSettings(cfg)

	// Create logger
	logger, err := forge.NewLogger(cfg.IsDevelopment())
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Connect to database
	database, err := forge.NewDBFromConfig(cfg)
	if err != nil {
		logger.Warn("Failed to connect to database", forge.Error(err))
		logger.Info("Server will start without database connection")
	} else {
		defer database.Close()
		logger.Info("Database connection established")

		// Set database on all managers (after code generation)
		setManagersDB(database)
	}

	// Initialize user system
	var userSystem *users.UserSystem
	if database != nil {
		userSystem, err = users.SetupUserSystem(database, nil)
		if err != nil {
			logger.Warn("Failed to initialize user system", forge.Error(err))
		} else {
			logger.Info("User system initialized successfully")
		}
	}

	// Create server
	server, err := forge.NewServer(cfg, settings, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Register routes
	server.RegisterRoutes(func(router *forge.Router) {
		// Health check
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome to Ecommerce API!")
		})

		// Register admin routes
		if settings.Admin.Enabled {
			var sessionManager interface{} = nil
			var db interface{} = nil
			if database != nil {
				db = database
			}
			forge.RegisterAdminRoutes(router, settings.Admin.Path, sessionManager, db)
		}

		// Register admin models (old system - kept for compatibility)
		registerAdminModels()

		// Register admin models with new v2 system (full features)
		if settings.Admin.Enabled {
			admin.SetupAdmin(router, settings.Admin.Path)
		}

		// Register user system routes
		if userSystem != nil {
			userSystem.RegisterRoutes(router)
			logger.Info("User system routes registered")
		}

		// Register REST API routes
		if database != nil {
			// Setup complete API with all features
			forgeapi.SetupCompleteAPI()
			
			// Configure default authentication
			forgeapi.SetDefaultAuthentication(
				authentication.NewTokenAuthentication(func(token string) (interface{}, error) {
					// Token lookup logic
					return nil, nil
				}),
			)
			
			// Simple API routes (basic CRUD)
			simpleRouter := forgeapi.NewRouter("/api/v1")
			api.RegisterSimpleAPIViewsets(simpleRouter)
			simpleRouter.RegisterRoutes(router)
			
			// Complex API routes (full features)
			complexRouter := forgeapi.NewEnhancedRouter("/api/v1")
			api.RegisterComplexAPIViewsets(complexRouter)
			complexRouter.RegisterRoutesEnhanced(router)
			
			// Legacy routes (backward compatible)
			legacyRouter := forgeapi.NewRouter("/api/v1/legacy")
			api.RegisterAPIViewsets(legacyRouter)
			legacyRouter.RegisterRoutes(router)
		}
	})

	// Start server
	fmt.Printf("Starting server on %s:%s\n", settings.Server.Host, settings.Server.Port)
	fmt.Printf("Admin interface available at http://%s:%s%s\n", settings.Server.Host, settings.Server.Port, settings.Admin.Path)
	fmt.Printf("REST API available at http://%s:%s/api/v1/\n", settings.Server.Host, settings.Server.Port)
	if userSystem != nil {
		fmt.Printf("User system API available at http://%s:%s/api/auth/ and http://%s:%s/api/users/\n", 
			settings.Server.Host, settings.Server.Port, settings.Server.Host, settings.Server.Port)
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// setManagersDB sets the database connection on all model managers
// This function will work after code generation creates the managers
func setManagersDB(database *forge.DB) {
	// Note: These will be available after running `forge generate`
	// The generated code creates ModelName.Objects for each model

	// Core models
	// models.Customer.Objects.SetDB(database)
	// models.CustomerProfile.Objects.SetDB(database)
	// models.Address.Objects.SetDB(database)

	// Product models
	// models.Brand.Objects.SetDB(database)
	// models.Supplier.Objects.SetDB(database)
	// models.Category.Objects.SetDB(database)
	// models.Product.Objects.SetDB(database)
	// models.ProductVariant.Objects.SetDB(database)
	// models.Inventory.Objects.SetDB(database)
	// models.Warehouse.Objects.SetDB(database)

	// Order models
	// models.Order.Objects.SetDB(database)
	// models.OrderItem.Objects.SetDB(database)
	// models.Payment.Objects.SetDB(database)
	// models.Shipping.Objects.SetDB(database)

	// Review model
	// models.Review.Objects.SetDB(database)
}

// registerAdminModels registers all models with the admin interface
// Note: After code generation, use RegisterModelWithManager to include managers
func registerAdminModels() {
	// Core models with admin configuration
	admin.RegisterModelWithOptions(
		&models.Customer{},
		admin.WithListDisplay("id", "email", "first_name", "last_name", "is_active", "total_orders", "created_at"),
		admin.WithSearchFields("email", "first_name", "last_name"),
		admin.WithListFilter("is_active", "is_verified", "is_premium", "created_at"),
	)
	admin.RegisterModel(&models.CustomerProfile{})
	admin.RegisterModelWithOptions(
		&models.Address{},
		admin.WithListDisplay("id", "customer_id", "type", "city", "state", "country", "is_default"),
		admin.WithSearchFields("street_address", "city", "state", "postal_code"),
	)

	// Product models
	admin.RegisterModelWithOptions(
		&models.Brand{},
		admin.WithListDisplay("id", "name", "slug", "is_active", "created_at"),
		admin.WithSearchFields("name", "slug"),
	)
	admin.RegisterModel(&models.Supplier{})
	admin.RegisterModelWithOptions(
		&models.Category{},
		admin.WithListDisplay("id", "name", "slug", "parent_id", "is_active"),
		admin.WithSearchFields("name", "slug"),
	)
	admin.RegisterModelWithOptions(
		&models.Product{},
		admin.WithListDisplay("id", "name", "sku", "price", "stock_quantity", "is_active", "is_featured", "status"),
		admin.WithSearchFields("name", "sku", "slug", "description"),
		admin.WithListFilter("is_active", "is_featured", "status", "brand_id", "created_at"),
	)
	admin.RegisterModel(&models.ProductVariant{})
	admin.RegisterModelWithOptions(
		&models.Inventory{},
		admin.WithListDisplay("id", "product_id", "warehouse_id", "quantity", "available_quantity"),
		admin.WithListFilter("warehouse_id", "product_id"),
	)
	admin.RegisterModel(&models.Warehouse{})

	// Order models
	admin.RegisterModelWithOptions(
		&models.Order{},
		admin.WithListDisplay("id", "order_number", "customer_id", "status", "total_amount", "payment_status", "placed_at"),
		admin.WithSearchFields("order_number"),
		admin.WithListFilter("status", "payment_status", "shipping_status", "placed_at"),
	)
	admin.RegisterModel(&models.OrderItem{})
	admin.RegisterModel(&models.Payment{})
	admin.RegisterModel(&models.Shipping{})

	// Review model
	admin.RegisterModelWithOptions(
		&models.Review{},
		admin.WithListDisplay("id", "product_id", "customer_id", "rating", "is_approved", "created_at"),
		admin.WithListFilter("rating", "is_approved", "is_verified_purchase"),
	)
}
