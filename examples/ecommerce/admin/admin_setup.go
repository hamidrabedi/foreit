package admin

import (
	"context"
	"reflect"

	"ecommerce/models"

adminv2 "github.com/forgego/forge/pkg/admin"
adminhttp "github.com/forgego/forge/pkg/admin/http"
	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/schema"
)

// SetupAdmin configures the full admin system for all ecommerce models
// This uses the new type-safe v2 admin system with all features
// NOTE: After code generation, uncomment manager assignments and use actual managers
func SetupAdmin(router *httplib.Router, adminPath string) {
	// Register all models with comprehensive admin configuration
	registerAllModels()

	// Register admin routes
	adminRouter := adminhttp.NewRouter(adminv2.GetGlobalRegistry())
	adminRouter.RegisterRoutes(router, adminPath)
}

// registerAllModels registers all ecommerce models with full admin features
func registerAllModels() {
	// ============================================
	// CUSTOMER MODELS
	// ============================================
	registerCustomer()
	registerCustomerProfile()
	registerAddress()

	// ============================================
	// PRODUCT MODELS
	// ============================================
	registerBrand()
	registerSupplier()
	registerCategory()
	registerProduct()
	registerProductVariant()
	registerInventory()
	registerWarehouse()

	// ============================================
	// ORDER MODELS
	// ============================================
	registerOrder()
	registerOrderItem()
	registerPayment()
	registerShipping()

	// ============================================
	// REVIEW MODEL
	// ============================================
	registerReview()
}

// ============================================
// CUSTOMER MODELS
// ============================================

// registerCustomer registers Customer model with comprehensive admin features
func registerCustomer() {
	// After code generation, uncomment and use:
	// customerAdmin := adminv2.Register(
	// 	&models.Customer{},
	// 	models.Customer.Objects,
	// 	&adminv2.Config[*models.Customer]{...}
	// )

	// For now, create admin using reflection-based field access
	// This will work before code generation
	customerAdmin := adminv2.Register(
		&models.Customer{},
		nil, // Will be models.Customer.Objects after generation
		&adminv2.Config[*models.Customer]{
			// List Display - fields shown in list view
			ListDisplay: []adminv2.FieldExpr[*models.Customer, interface{}]{
				createStringField[*models.Customer]("email"),
				createStringField[*models.Customer]("first_name"),
				createStringField[*models.Customer]("last_name"),
				createBoolField[*models.Customer]("is_active"),
				createBoolField[*models.Customer]("is_verified"),
				createBoolField[*models.Customer]("is_premium"),
				createInt64Field[*models.Customer]("total_orders"),
				createDecimalField[*models.Customer]("lifetime_value"),
				createTimeField[*models.Customer]("created_at"),
			},

			// Search Fields - fields to search in
			SearchFields: []adminv2.FieldExpr[*models.Customer, interface{}]{
				createStringField[*models.Customer]("email"),
				createStringField[*models.Customer]("first_name"),
				createStringField[*models.Customer]("last_name"),
				createStringField[*models.Customer]("phone"),
			},

			// Filters - filters for list view
			ListFilter: []adminv2.Filter[*models.Customer]{
				adminv2.NewBooleanFilter(createBoolField[*models.Customer]("is_active")),
				adminv2.NewBooleanFilter(createBoolField[*models.Customer]("is_verified")),
				adminv2.NewBooleanFilter(createBoolField[*models.Customer]("is_premium")),
				// Date filter - will be implemented after code generation
				// adminv2.NewDateFilter(createTimeField[*models.Customer]("created_at")),
			},

			// Ordering - default ordering
			Ordering: []adminv2.Ordering[*models.Customer]{
				adminv2.OrderBy(createTimeField[*models.Customer]("created_at")).Desc(),
				adminv2.OrderBy(createStringField[*models.Customer]("email")).Asc(),
			},

			// Bulk Actions
			Actions: []adminv2.Action[*models.Customer]{
				adminv2.NewAction(
					"activate",
					"Activate selected customers",
					func(ctx context.Context, customers []*models.Customer) error {
						for _, customer := range customers {
							setFieldValue(customer, "is_active", true)
							// After generation: customer.IsActive = true
							// if err := models.Customer.Objects.Update(ctx, customer); err != nil {
							// 	return err
							// }
						}
						return nil
					},
				),
				adminv2.NewAction(
					"deactivate",
					"Deactivate selected customers",
					func(ctx context.Context, customers []*models.Customer) error {
						for _, customer := range customers {
							setFieldValue(customer, "is_active", false)
						}
						return nil
					},
				),
				adminv2.NewAction(
					"mark_premium",
					"Mark as premium members",
					func(ctx context.Context, customers []*models.Customer) error {
						for _, customer := range customers {
							setFieldValue(customer, "is_premium", true)
						}
						return nil
					},
				),
			},

			// Fieldsets - group form fields
			Fieldsets: []adminv2.Fieldset[*models.Customer]{
				adminv2.NewFieldset(
					"Account Information",
					createStringField[*models.Customer]("email"),
					createStringField[*models.Customer]("phone"),
					createStringField[*models.Customer]("password_hash"),
				),
				adminv2.NewFieldset(
					"Personal Information",
					createStringField[*models.Customer]("first_name"),
					createStringField[*models.Customer]("last_name"),
				),
				adminv2.NewFieldset(
					"Status & Permissions",
					createBoolField[*models.Customer]("is_active"),
					createBoolField[*models.Customer]("is_verified"),
					createBoolField[*models.Customer]("is_premium"),
				),
				adminv2.NewFieldset(
					"Statistics",
					createInt64Field[*models.Customer]("total_orders"),
					createDecimalField[*models.Customer]("lifetime_value"),
				),
			},

			// Inlines - related models edited inline
			// Note: Inlines require proper field expressions for parent relationship
			// This is a placeholder - full implementation after code generation

			ListPerPage:       25,
			VerboseName:       "Customer",
			VerboseNamePlural: "Customers",
		},
	)

	adminhttp.RegisterAdminForHTTP(customerAdmin)
}

// registerCustomerProfile registers CustomerProfile model
func registerCustomerProfile() {
	customerProfileAdmin := adminv2.Register(
		&models.CustomerProfile{},
		nil, // Will be models.CustomerProfile.Objects after generation
		&adminv2.Config[*models.CustomerProfile]{
			ListDisplay: []adminv2.FieldExpr[*models.CustomerProfile, interface{}]{
				createInt64Field[*models.CustomerProfile]("customer_id"),
				createStringField[*models.CustomerProfile]("gender"),
				createStringField[*models.CustomerProfile]("preferred_language"),
				createTimeField[*models.CustomerProfile]("created_at"),
			},
			SearchFields: []adminv2.FieldExpr[*models.CustomerProfile, interface{}]{
				createStringField[*models.CustomerProfile]("preferred_language"),
			},
			ListFilter: []adminv2.Filter[*models.CustomerProfile]{
				adminv2.NewChoiceFilter(
					createStringField[*models.CustomerProfile]("gender"),
					[]adminv2.Choice[string]{
						{Label: "Male", Value: "male"},
						{Label: "Female", Value: "female"},
						{Label: "Other", Value: "other"},
					},
				),
			},
			Fieldsets: []adminv2.Fieldset[*models.CustomerProfile]{
				adminv2.NewFieldset(
					"Personal Details",
					createDateField[*models.CustomerProfile]("date_of_birth"),
					createStringField[*models.CustomerProfile]("gender"),
					createStringField[*models.CustomerProfile]("avatar_url"),
				),
				adminv2.NewFieldset(
					"Preferences",
					createStringField[*models.CustomerProfile]("preferred_language"),
					createStringField[*models.CustomerProfile]("timezone"),
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(customerProfileAdmin)
}

// registerAddress registers Address model
func registerAddress() {
	addressAdmin := adminv2.Register(
		&models.Address{},
		nil, // Will be models.Address.Objects after generation
		&adminv2.Config[*models.Address]{
			ListDisplay: []adminv2.FieldExpr[*models.Address, interface{}]{
				createInt64Field[*models.Address]("customer_id"),
				createStringField[*models.Address]("type"),
				createStringField[*models.Address]("city"),
				createStringField[*models.Address]("state"),
				createStringField[*models.Address]("country"),
				createBoolField[*models.Address]("is_default"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Address, interface{}]{
				createStringField[*models.Address]("street_address"),
				createStringField[*models.Address]("city"),
				createStringField[*models.Address]("state"),
				createStringField[*models.Address]("postal_code"),
			},
			ListFilter: []adminv2.Filter[*models.Address]{
				adminv2.NewChoiceFilter(
					createStringField[*models.Address]("type"),
					[]adminv2.Choice[string]{
						{Label: "Billing", Value: "billing"},
						{Label: "Shipping", Value: "shipping"},
						{Label: "Both", Value: "both"},
					},
				),
				adminv2.NewBooleanFilter(createBoolField[*models.Address]("is_default")),
				adminv2.NewBooleanFilter(createBoolField[*models.Address]("is_active")),
			},
			Fieldsets: []adminv2.Fieldset[*models.Address]{
				adminv2.NewFieldset(
					"Contact Information",
					createStringField[*models.Address]("first_name"),
					createStringField[*models.Address]("last_name"),
					createStringField[*models.Address]("company"),
					createStringField[*models.Address]("phone"),
				),
				adminv2.NewFieldset(
					"Address",
					createStringField[*models.Address]("address_line1"),
					createStringField[*models.Address]("address_line2"),
					createStringField[*models.Address]("city"),
					createStringField[*models.Address]("state"),
					createStringField[*models.Address]("postal_code"),
					createStringField[*models.Address]("country"),
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(addressAdmin)
}

// ============================================
// PRODUCT MODELS
// ============================================

// registerBrand registers Brand model
func registerBrand() {
	brandAdmin := adminv2.Register(
		&models.Brand{},
		nil, // Will be models.Brand.Objects after generation
		&adminv2.Config[*models.Brand]{
			ListDisplay: []adminv2.FieldExpr[*models.Brand, interface{}]{
				createStringField[*models.Brand]("name"),
				createStringField[*models.Brand]("slug"),
				createBoolField[*models.Brand]("is_active"),
				createBoolField[*models.Brand]("is_featured"),
				createTimeField[*models.Brand]("created_at"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Brand, interface{}]{
				createStringField[*models.Brand]("name"),
				createStringField[*models.Brand]("slug"),
			},
			ListFilter: []adminv2.Filter[*models.Brand]{
				adminv2.NewBooleanFilter(createBoolField[*models.Brand]("is_active")),
				adminv2.NewBooleanFilter(createBoolField[*models.Brand]("is_featured")),
			},
			Actions: []adminv2.Action[*models.Brand]{
				adminv2.NewAction(
					"feature",
					"Feature selected brands",
					func(ctx context.Context, brands []*models.Brand) error {
						for _, brand := range brands {
							setFieldValue(brand, "is_featured", true)
						}
						return nil
					},
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(brandAdmin)
}

// registerSupplier registers Supplier model
func registerSupplier() {
	supplierAdmin := adminv2.Register(
		&models.Supplier{},
		nil, // Will be models.Supplier.Objects after generation
		&adminv2.Config[*models.Supplier]{
			ListDisplay: []adminv2.FieldExpr[*models.Supplier, interface{}]{
				createStringField[*models.Supplier]("name"),
				createStringField[*models.Supplier]("code"),
				createStringField[*models.Supplier]("contact_name"),
				createStringField[*models.Supplier]("email"),
				createStringField[*models.Supplier]("city"),
				createStringField[*models.Supplier]("country"),
				createBoolField[*models.Supplier]("is_active"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Supplier, interface{}]{
				createStringField[*models.Supplier]("name"),
				createStringField[*models.Supplier]("code"),
				createStringField[*models.Supplier]("email"),
			},
			ListFilter: []adminv2.Filter[*models.Supplier]{
				adminv2.NewBooleanFilter(createBoolField[*models.Supplier]("is_active")),
			},
			Fieldsets: []adminv2.Fieldset[*models.Supplier]{
				adminv2.NewFieldset(
					"Company Information",
					createStringField[*models.Supplier]("name"),
					createStringField[*models.Supplier]("code"),
					createStringField[*models.Supplier]("website_url"),
				),
				adminv2.NewFieldset(
					"Contact Information",
					createStringField[*models.Supplier]("contact_name"),
					createStringField[*models.Supplier]("email"),
					createStringField[*models.Supplier]("phone"),
				),
				adminv2.NewFieldset(
					"Address",
					createStringField[*models.Supplier]("address_line1"),
					createStringField[*models.Supplier]("address_line2"),
					createStringField[*models.Supplier]("city"),
					createStringField[*models.Supplier]("state"),
					createStringField[*models.Supplier]("postal_code"),
					createStringField[*models.Supplier]("country"),
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(supplierAdmin)
}

// registerCategory registers Category model with hierarchical support
func registerCategory() {
	categoryAdmin := adminv2.Register(
		&models.Category{},
		nil, // Will be models.Category.Objects after generation
		&adminv2.Config[*models.Category]{
			ListDisplay: []adminv2.FieldExpr[*models.Category, interface{}]{
				createStringField[*models.Category]("name"),
				createStringField[*models.Category]("slug"),
				createInt64Field[*models.Category]("parent_id"),
				createInt32Field[*models.Category]("level"),
				createInt32Field[*models.Category]("sort_order"),
				createBoolField[*models.Category]("is_active"),
				createBoolField[*models.Category]("is_featured"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Category, interface{}]{
				createStringField[*models.Category]("name"),
				createStringField[*models.Category]("slug"),
			},
			ListFilter: []adminv2.Filter[*models.Category]{
				adminv2.NewBooleanFilter(createBoolField[*models.Category]("is_active")),
				adminv2.NewBooleanFilter(createBoolField[*models.Category]("is_featured")),
			},
			DateHierarchy: createTimeField[*models.Category]("created_at"),
			Fieldsets: []adminv2.Fieldset[*models.Category]{
				adminv2.NewFieldset(
					"Basic Information",
					createStringField[*models.Category]("name"),
					createStringField[*models.Category]("slug"),
					createInt64Field[*models.Category]("parent_id"),
					createInt32Field[*models.Category]("level"),
					createInt32Field[*models.Category]("sort_order"),
				),
				adminv2.NewFieldset(
					"Content",
					createStringField[*models.Category]("description"),
					createStringField[*models.Category]("image_url"),
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(categoryAdmin)
}

// registerProduct registers Product model with full ecommerce features
func registerProduct() {
	productAdmin := adminv2.Register(
		&models.Product{},
		nil, // Will be models.Product.Objects after generation
		&adminv2.Config[*models.Product]{
			ListDisplay: []adminv2.FieldExpr[*models.Product, interface{}]{
				createStringField[*models.Product]("name"),
				createStringField[*models.Product]("sku"),
				createDecimalField[*models.Product]("price"),
				createInt32Field[*models.Product]("stock_quantity"),
				createStringField[*models.Product]("status"),
				createBoolField[*models.Product]("is_active"),
				createBoolField[*models.Product]("is_featured"),
				createInt64Field[*models.Product]("brand_id"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Product, interface{}]{
				createStringField[*models.Product]("name"),
				createStringField[*models.Product]("sku"),
				createStringField[*models.Product]("slug"),
				createStringField[*models.Product]("description"),
			},
			ListFilter: []adminv2.Filter[*models.Product]{
				adminv2.NewBooleanFilter(createBoolField[*models.Product]("is_active")),
				adminv2.NewBooleanFilter(createBoolField[*models.Product]("is_featured")),
				adminv2.NewBooleanFilter(createBoolField[*models.Product]("is_digital")),
				adminv2.NewChoiceFilter(
					createStringField[*models.Product]("status"),
					[]adminv2.Choice[string]{
						{Label: "Draft", Value: "draft"},
						{Label: "Active", Value: "active"},
						{Label: "Archived", Value: "archived"},
						{Label: "Discontinued", Value: "discontinued"},
					},
				),
			},
			Ordering: []adminv2.Ordering[*models.Product]{
				adminv2.OrderBy(createTimeField[*models.Product]("created_at")).Desc(),
				adminv2.OrderBy(createStringField[*models.Product]("name")).Asc(),
			},
			Actions: []adminv2.Action[*models.Product]{
				adminv2.NewAction(
					"activate",
					"Activate selected products",
					func(ctx context.Context, products []*models.Product) error {
						for _, product := range products {
							setFieldValue(product, "is_active", true)
							setFieldValue(product, "status", "active")
						}
						return nil
					},
				),
				adminv2.NewAction(
					"feature",
					"Feature selected products",
					func(ctx context.Context, products []*models.Product) error {
						for _, product := range products {
							setFieldValue(product, "is_featured", true)
						}
						return nil
					},
				),
				adminv2.NewAction(
					"archive",
					"Archive selected products",
					func(ctx context.Context, products []*models.Product) error {
						for _, product := range products {
							setFieldValue(product, "status", "archived")
						}
						return nil
					},
				),
			},
			Fieldsets: []adminv2.Fieldset[*models.Product]{
				adminv2.NewFieldset(
					"Basic Information",
					createStringField[*models.Product]("name"),
					createStringField[*models.Product]("slug"),
					createStringField[*models.Product]("sku"),
					createStringField[*models.Product]("description"),
					createStringField[*models.Product]("short_description"),
				),
				adminv2.NewFieldset(
					"Pricing",
					createDecimalField[*models.Product]("price"),
					createDecimalField[*models.Product]("compare_at_price"),
					createDecimalField[*models.Product]("cost_price"),
					createStringField[*models.Product]("currency"),
				),
				adminv2.NewFieldset(
					"Inventory",
					createInt32Field[*models.Product]("stock_quantity"),
					createInt32Field[*models.Product]("low_stock_threshold"),
					createBoolField[*models.Product]("track_inventory"),
				),
				adminv2.NewFieldset(
					"Shipping",
					createBoolField[*models.Product]("requires_shipping"),
					createFloat64Field[*models.Product]("weight"),
					createFloat64Field[*models.Product]("length"),
					createFloat64Field[*models.Product]("width"),
					createFloat64Field[*models.Product]("height"),
				),
				adminv2.NewFieldset(
					"Status & Options",
					createStringField[*models.Product]("status"),
					createBoolField[*models.Product]("is_active"),
					createBoolField[*models.Product]("is_featured"),
					createBoolField[*models.Product]("is_digital"),
				),
			},
			// Inlines would be added here after code generation
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(productAdmin)
}

// registerProductVariant registers ProductVariant model
func registerProductVariant() {
	variantAdmin := adminv2.Register(
		&models.ProductVariant{},
		nil, // Will be models.ProductVariant.Objects after generation
		&adminv2.Config[*models.ProductVariant]{
			ListDisplay: []adminv2.FieldExpr[*models.ProductVariant, interface{}]{
				createInt64Field[*models.ProductVariant]("product_id"),
				createStringField[*models.ProductVariant]("name"),
				createStringField[*models.ProductVariant]("sku"),
				createStringField[*models.ProductVariant]("option1"),
				createStringField[*models.ProductVariant]("option2"),
				createDecimalField[*models.ProductVariant]("price"),
				createInt32Field[*models.ProductVariant]("stock_quantity"),
				createBoolField[*models.ProductVariant]("is_active"),
			},
			SearchFields: []adminv2.FieldExpr[*models.ProductVariant, interface{}]{
				createStringField[*models.ProductVariant]("name"),
				createStringField[*models.ProductVariant]("sku"),
			},
			ListFilter: []adminv2.Filter[*models.ProductVariant]{
				adminv2.NewBooleanFilter(createBoolField[*models.ProductVariant]("is_active")),
				adminv2.NewBooleanFilter(createBoolField[*models.ProductVariant]("is_default")),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(variantAdmin)
}

// registerInventory registers Inventory model
func registerInventory() {
	inventoryAdmin := adminv2.Register(
		&models.Inventory{},
		nil, // Will be models.Inventory.Objects after generation
		&adminv2.Config[*models.Inventory]{
			ListDisplay: []adminv2.FieldExpr[*models.Inventory, interface{}]{
				createInt64Field[*models.Inventory]("product_id"),
				createInt64Field[*models.Inventory]("variant_id"),
				createInt64Field[*models.Inventory]("warehouse_id"),
				createInt32Field[*models.Inventory]("quantity"),
				createInt32Field[*models.Inventory]("available_quantity"),
				createInt32Field[*models.Inventory]("reserved_quantity"),
			},
			ListFilter: []adminv2.Filter[*models.Inventory]{
				adminv2.NewBooleanFilter(createBoolField[*models.Inventory]("is_active")),
			},
			ListPerPage: 50,
		},
	)

	adminhttp.RegisterAdminForHTTP(inventoryAdmin)
}

// registerWarehouse registers Warehouse model
func registerWarehouse() {
	warehouseAdmin := adminv2.Register(
		&models.Warehouse{},
		nil, // Will be models.Warehouse.Objects after generation
		&adminv2.Config[*models.Warehouse]{
			ListDisplay: []adminv2.FieldExpr[*models.Warehouse, interface{}]{
				createStringField[*models.Warehouse]("name"),
				createStringField[*models.Warehouse]("code"),
				createStringField[*models.Warehouse]("city"),
				createStringField[*models.Warehouse]("state"),
				createStringField[*models.Warehouse]("country"),
				createBoolField[*models.Warehouse]("is_active"),
				createBoolField[*models.Warehouse]("is_primary"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Warehouse, interface{}]{
				createStringField[*models.Warehouse]("name"),
				createStringField[*models.Warehouse]("code"),
				createStringField[*models.Warehouse]("city"),
			},
			ListFilter: []adminv2.Filter[*models.Warehouse]{
				adminv2.NewBooleanFilter(createBoolField[*models.Warehouse]("is_active")),
				adminv2.NewBooleanFilter(createBoolField[*models.Warehouse]("is_primary")),
			},
			Fieldsets: []adminv2.Fieldset[*models.Warehouse]{
				adminv2.NewFieldset(
					"Basic Information",
					createStringField[*models.Warehouse]("name"),
					createStringField[*models.Warehouse]("code"),
				),
				adminv2.NewFieldset(
					"Address",
					createStringField[*models.Warehouse]("address_line1"),
					createStringField[*models.Warehouse]("address_line2"),
					createStringField[*models.Warehouse]("city"),
					createStringField[*models.Warehouse]("state"),
					createStringField[*models.Warehouse]("postal_code"),
					createStringField[*models.Warehouse]("country"),
				),
				adminv2.NewFieldset(
					"Contact",
					createStringField[*models.Warehouse]("phone"),
					createStringField[*models.Warehouse]("email"),
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(warehouseAdmin)
}

// ============================================
// ORDER MODELS
// ============================================

// registerOrder registers Order model with full order management
func registerOrder() {
	orderAdmin := adminv2.Register(
		&models.Order{},
		nil, // Will be models.Order.Objects after generation
		&adminv2.Config[*models.Order]{
			ListDisplay: []adminv2.FieldExpr[*models.Order, interface{}]{
				createStringField[*models.Order]("order_number"),
				createInt64Field[*models.Order]("customer_id"),
				createStringField[*models.Order]("status"),
				createDecimalField[*models.Order]("total_amount"),
				createStringField[*models.Order]("payment_status"),
				createStringField[*models.Order]("shipping_status"),
				createTimeField[*models.Order]("placed_at"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Order, interface{}]{
				createStringField[*models.Order]("order_number"),
			},
			ListFilter: []adminv2.Filter[*models.Order]{
				adminv2.NewChoiceFilter(
					createStringField[*models.Order]("status"),
					[]adminv2.Choice[string]{
						{Label: "Pending", Value: "pending"},
						{Label: "Confirmed", Value: "confirmed"},
						{Label: "Processing", Value: "processing"},
						{Label: "Shipped", Value: "shipped"},
						{Label: "Delivered", Value: "delivered"},
						{Label: "Cancelled", Value: "cancelled"},
						{Label: "Refunded", Value: "refunded"},
					},
				),
				adminv2.NewChoiceFilter(
					createStringField[*models.Order]("payment_status"),
					[]adminv2.Choice[string]{
						{Label: "Pending", Value: "pending"},
						{Label: "Paid", Value: "paid"},
						{Label: "Failed", Value: "failed"},
						{Label: "Refunded", Value: "refunded"},
					},
				),
				adminv2.NewChoiceFilter(
					createStringField[*models.Order]("shipping_status"),
					[]adminv2.Choice[string]{
						{Label: "Pending", Value: "pending"},
						{Label: "Processing", Value: "processing"},
						{Label: "Shipped", Value: "shipped"},
						{Label: "Delivered", Value: "delivered"},
						{Label: "Returned", Value: "returned"},
					},
				),
				// Date filter - placeholder (will work after code generation)
				// adminv2.NewDateFilter(createTimeField[*models.Order]("placed_at")),
			},
			DateHierarchy: createTimeField[*models.Order]("placed_at"),
			Ordering: []adminv2.Ordering[*models.Order]{
				adminv2.OrderBy(createTimeField[*models.Order]("placed_at")).Desc(),
			},
			Actions: []adminv2.Action[*models.Order]{
				adminv2.NewAction(
					"confirm",
					"Confirm selected orders",
					func(ctx context.Context, orders []*models.Order) error {
						for _, order := range orders {
							setFieldValue(order, "status", "confirmed")
						}
						return nil
					},
				),
				adminv2.NewAction(
					"ship",
					"Mark as shipped",
					func(ctx context.Context, orders []*models.Order) error {
						for _, order := range orders {
							setFieldValue(order, "status", "shipped")
							setFieldValue(order, "shipping_status", "shipped")
						}
						return nil
					},
				),
				adminv2.NewAction(
					"cancel",
					"Cancel selected orders",
					func(ctx context.Context, orders []*models.Order) error {
						for _, order := range orders {
							setFieldValue(order, "status", "cancelled")
						}
						return nil
					},
				),
			},
			Fieldsets: []adminv2.Fieldset[*models.Order]{
				adminv2.NewFieldset(
					"Order Information",
					createStringField[*models.Order]("order_number"),
					createInt64Field[*models.Order]("customer_id"),
					createStringField[*models.Order]("status"),
				),
				adminv2.NewFieldset(
					"Totals",
					createDecimalField[*models.Order]("subtotal"),
					createDecimalField[*models.Order]("tax_amount"),
					createDecimalField[*models.Order]("shipping_amount"),
					createDecimalField[*models.Order]("discount_amount"),
					createDecimalField[*models.Order]("total_amount"),
					createStringField[*models.Order]("currency"),
				),
				adminv2.NewFieldset(
					"Status",
					createStringField[*models.Order]("payment_status"),
					createStringField[*models.Order]("shipping_status"),
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(orderAdmin)
}

// registerOrderItem registers OrderItem model
func registerOrderItem() {
	orderItemAdmin := adminv2.Register(
		&models.OrderItem{},
		nil, // Will be models.OrderItem.Objects after generation
		&adminv2.Config[*models.OrderItem]{
			ListDisplay: []adminv2.FieldExpr[*models.OrderItem, interface{}]{
				createInt64Field[*models.OrderItem]("order_id"),
				createStringField[*models.OrderItem]("product_name"),
				createStringField[*models.OrderItem]("product_sku"),
				createInt32Field[*models.OrderItem]("quantity"),
				createDecimalField[*models.OrderItem]("unit_price"),
				createDecimalField[*models.OrderItem]("total_price"),
			},
			ListFilter: []adminv2.Filter[*models.OrderItem]{
				// Filter by order_id would be useful
			},
			ListPerPage: 50,
		},
	)

	adminhttp.RegisterAdminForHTTP(orderItemAdmin)
}

// registerPayment registers Payment model
func registerPayment() {
	paymentAdmin := adminv2.Register(
		&models.Payment{},
		nil, // Will be models.Payment.Objects after generation
		&adminv2.Config[*models.Payment]{
			ListDisplay: []adminv2.FieldExpr[*models.Payment, interface{}]{
				createInt64Field[*models.Payment]("order_id"),
				createStringField[*models.Payment]("transaction_id"),
				createStringField[*models.Payment]("payment_method"),
				createStringField[*models.Payment]("status"),
				createDecimalField[*models.Payment]("amount"),
				createTimeField[*models.Payment]("created_at"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Payment, interface{}]{
				createStringField[*models.Payment]("transaction_id"),
				createStringField[*models.Payment]("gateway_transaction_id"),
			},
			ListFilter: []adminv2.Filter[*models.Payment]{
				adminv2.NewChoiceFilter(
					createStringField[*models.Payment]("status"),
					[]adminv2.Choice[string]{
						{Label: "Pending", Value: "pending"},
						{Label: "Processing", Value: "processing"},
						{Label: "Completed", Value: "completed"},
						{Label: "Failed", Value: "failed"},
						{Label: "Refunded", Value: "refunded"},
						{Label: "Cancelled", Value: "cancelled"},
					},
				),
				adminv2.NewChoiceFilter(
					createStringField[*models.Payment]("payment_method"),
					[]adminv2.Choice[string]{
						{Label: "Credit Card", Value: "credit_card"},
						{Label: "Debit Card", Value: "debit_card"},
						{Label: "PayPal", Value: "paypal"},
						{Label: "Stripe", Value: "stripe"},
						{Label: "Bank Transfer", Value: "bank_transfer"},
						{Label: "Cash on Delivery", Value: "cash_on_delivery"},
					},
				),
			},
			Actions: []adminv2.Action[*models.Payment]{
				adminv2.NewAction(
					"refund",
					"Refund selected payments",
					func(ctx context.Context, payments []*models.Payment) error {
						for _, payment := range payments {
							setFieldValue(payment, "status", "refunded")
							setFieldValue(payment, "is_refunded", true)
						}
						return nil
					},
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(paymentAdmin)
}

// registerShipping registers Shipping model
func registerShipping() {
	shippingAdmin := adminv2.Register(
		&models.Shipping{},
		nil, // Will be models.Shipping.Objects after generation
		&adminv2.Config[*models.Shipping]{
			ListDisplay: []adminv2.FieldExpr[*models.Shipping, interface{}]{
				createInt64Field[*models.Shipping]("order_id"),
				createStringField[*models.Shipping]("carrier"),
				createStringField[*models.Shipping]("service"),
				createStringField[*models.Shipping]("tracking_number"),
				createStringField[*models.Shipping]("status"),
				createTimeField[*models.Shipping]("created_at"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Shipping, interface{}]{
				createStringField[*models.Shipping]("tracking_number"),
			},
			ListFilter: []adminv2.Filter[*models.Shipping]{
				adminv2.NewChoiceFilter(
					createStringField[*models.Shipping]("status"),
					[]adminv2.Choice[string]{
						{Label: "Pending", Value: "pending"},
						{Label: "Label Created", Value: "label_created"},
						{Label: "In Transit", Value: "in_transit"},
						{Label: "Out for Delivery", Value: "out_for_delivery"},
						{Label: "Delivered", Value: "delivered"},
						{Label: "Exception", Value: "exception"},
						{Label: "Returned", Value: "returned"},
					},
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(shippingAdmin)
}

// ============================================
// REVIEW MODEL
// ============================================

// registerReview registers Review model with moderation features
func registerReview() {
	reviewAdmin := adminv2.Register(
		&models.Review{},
		nil, // Will be models.Review.Objects after generation
		&adminv2.Config[*models.Review]{
			ListDisplay: []adminv2.FieldExpr[*models.Review, interface{}]{
				createInt64Field[*models.Review]("product_id"),
				createInt64Field[*models.Review]("customer_id"),
				createInt32Field[*models.Review]("rating"),
				createStringField[*models.Review]("title"),
				createBoolField[*models.Review]("is_approved"),
				createBoolField[*models.Review]("is_verified_purchase"),
				createTimeField[*models.Review]("created_at"),
			},
			SearchFields: []adminv2.FieldExpr[*models.Review, interface{}]{
				createStringField[*models.Review]("title"),
				createStringField[*models.Review]("comment"),
			},
			ListFilter: []adminv2.Filter[*models.Review]{
				adminv2.NewChoiceFilter(
					createInt32Field[*models.Review]("rating"),
					[]adminv2.Choice[int32]{
						{Label: "1 Star", Value: 1},
						{Label: "2 Stars", Value: 2},
						{Label: "3 Stars", Value: 3},
						{Label: "4 Stars", Value: 4},
						{Label: "5 Stars", Value: 5},
					},
				),
				adminv2.NewBooleanFilter(createBoolField[*models.Review]("is_approved")),
				adminv2.NewBooleanFilter(createBoolField[*models.Review]("is_verified_purchase")),
			},
			Actions: []adminv2.Action[*models.Review]{
				adminv2.NewAction(
					"approve",
					"Approve selected reviews",
					func(ctx context.Context, reviews []*models.Review) error {
						for _, review := range reviews {
							setFieldValue(review, "is_approved", true)
						}
						return nil
					},
				),
				adminv2.NewAction(
					"reject",
					"Reject selected reviews",
					func(ctx context.Context, reviews []*models.Review) error {
						for _, review := range reviews {
							setFieldValue(review, "is_approved", false)
						}
						return nil
					},
				),
			},
			ListPerPage: 25,
		},
	)

	adminhttp.RegisterAdminForHTTP(reviewAdmin)
}

// ============================================
// HELPER FUNCTIONS - Reflection-based field access
// ============================================
// These functions use reflection to access fields before code generation
// After code generation, replace with direct field access using getters/setters

// createStringField creates a string field expression using reflection
func createStringField[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return adminv2.NewFieldExpr(
		fieldName,
		func(instance *T) interface{} {
			return getFieldValue(instance, fieldName)
		},
		func(instance *T, value interface{}) {
			setFieldValue(instance, fieldName, value)
		},
		schema.Field{
			Name: fieldName,
			Type: schema.TypeString,
		},
	)
}

// createInt64Field creates an int64 field expression
func createInt64Field[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName) // Uses same reflection approach
}

// createInt32Field creates an int32 field expression
func createInt32Field[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName)
}

// createBoolField creates a bool field expression
func createBoolField[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName)
}

// createDecimalField creates a decimal field expression
func createDecimalField[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName)
}

// createTimeField creates a time.Time field expression
func createTimeField[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName)
}

// createDateField creates a date field expression
func createDateField[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName)
}

// createFloat64Field creates a float64 field expression
func createFloat64Field[T any](fieldName string) adminv2.FieldExpr[T, interface{}] {
	return createStringField[T](fieldName)
}

// getFieldValue gets a field value using reflection
func getFieldValue(instance interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	if field.CanInterface() {
		return field.Interface()
	}

	return nil
}

// setFieldValue sets a field value using reflection
func setFieldValue(instance interface{}, fieldName string, value interface{}) {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return
	}

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Type().AssignableTo(field.Type()) {
		field.Set(valueReflect)
	} else if valueReflect.Type().ConvertibleTo(field.Type()) {
		field.Set(valueReflect.Convert(field.Type()))
	}
}
