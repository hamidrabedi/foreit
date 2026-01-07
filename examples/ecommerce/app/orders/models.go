package orders

import (
	"context"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Cart represents a shopping cart
type Cart struct {
	schema.BaseSchema
}

func (Cart) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("customer_id").WithOptional().
			WithHelpText("Customer (null for guest carts)"),

		// Guest cart identification
		schema.String("session_id").WithMaxLength(255).WithOptional().
			WithHelpText("Session ID for guest carts"),
		schema.String("guest_email").WithMaxLength(255).WithOptional(),

		// Pricing
		schema.Float64("subtotal").WithDefault(0.0).
			WithHelpText("Cart subtotal before discounts"),
		schema.Float64("discount_amount").WithDefault(0.0),
		schema.Float64("tax_amount").WithDefault(0.0),
		schema.Float64("shipping_amount").WithDefault(0.0),
		schema.Float64("total").WithDefault(0.0).
			WithHelpText("Cart total"),

		// Coupon
		schema.Int64("coupon_id").WithOptional(),
		schema.String("coupon_code").WithMaxLength(50).WithOptional(),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("active").
			WithHelpText("Status: active, abandoned, converted"),
		schema.Bool("is_abandoned").WithDefault(false),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("last_activity_at").WithOptional(),
		schema.Time("converted_at").WithOptional().
			WithHelpText("When cart was converted to order"),
	}
}

func (Cart) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "carts",
		VerboseName:       "Shopping Cart",
		VerboseNamePlural: "Shopping Carts",
		OrderBy:           []string{"-updated_at"},
		Indexes: []schema.Index{
			{Name: "idx_cart_customer", Fields: []string{"customer_id"}},
			{Name: "idx_cart_session", Fields: []string{"session_id"}},
			{Name: "idx_cart_status", Fields: []string{"status"}},
			{Name: "idx_cart_abandoned", Fields: []string{"is_abandoned"}},
		},
	}
}

func (Cart) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithOptional().
			WithRelatedName("carts"),
		schema.ForeignKey("coupon_id", "Coupon").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("carts"),
	}
}

func (Cart) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterSave: func(ctx context.Context, instance interface{}) error {
			// Recalculate cart totals
			// Check for abandoned cart
			return nil
		},
	}
}

// CartItem represents items in a shopping cart
type CartItem struct {
	schema.BaseSchema
}

func (CartItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("cart_id").WithRequired(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("variant_id").WithOptional(),

		// Quantity and pricing
		schema.Int32("quantity").WithRequired().WithDefault(1),
		schema.Float64("unit_price").WithRequired().
			WithHelpText("Price per unit at time of adding"),
		schema.Float64("discount_amount").WithDefault(0.0),
		schema.Float64("tax_amount").WithDefault(0.0),
		schema.Float64("total").WithRequired().
			WithHelpText("Line total (quantity * unit_price - discount + tax)"),

		// Product snapshot (for price changes)
		schema.String("product_name").WithMaxLength(255).WithOptional(),
		schema.String("variant_name").WithMaxLength(255).WithOptional(),
		schema.String("image_url").WithMaxLength(500).WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (CartItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "cart_items",
		VerboseName:       "Cart Item",
		VerboseNamePlural: "Cart Items",
		OrderBy:           []string{"created_at"},
		Indexes: []schema.Index{
			{Name: "idx_cart_item_cart", Fields: []string{"cart_id"}},
			{Name: "idx_cart_item_product", Fields: []string{"product_id"}},
			{Name: "idx_cart_item_variant", Fields: []string{"variant_id"}},
		},
		UniqueTogether: [][]string{
			{"cart_id", "product_id", "variant_id"},
		},
	}
}

func (CartItem) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("cart_id", "Cart").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("items"),
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("cart_items"),
		schema.ForeignKey("variant_id", "ProductVariant").
			WithOnDelete(schema.CascadeCASCADE).
			WithOptional().
			WithRelatedName("cart_items"),
	}
}

func (CartItem) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Calculate line total
			// Capture product snapshot
			return nil
		},
		AfterSave: func(ctx context.Context, instance interface{}) error {
			// Update cart totals
			return nil
		},
	}
}

// Order represents a customer order
type Order struct {
	schema.BaseSchema
}

func (Order) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),

		// Order identification
		schema.String("order_number").WithRequired().WithMaxLength(50).WithUnique().
			WithHelpText("Human-readable order number"),

		// Customer
		schema.Int64("customer_id").WithRequired(),
		schema.String("customer_email").WithRequired().WithMaxLength(255),
		schema.String("customer_first_name").WithRequired().WithMaxLength(100),
		schema.String("customer_last_name").WithRequired().WithMaxLength(100),
		schema.String("customer_phone").WithMaxLength(20).WithOptional(),

		// Pricing
		schema.Float64("subtotal").WithRequired().
			WithHelpText("Order subtotal before discounts"),
		schema.Float64("discount_amount").WithDefault(0.0),
		schema.Float64("tax_amount").WithRequired(),
		schema.Float64("shipping_amount").WithRequired(),
		schema.Float64("total").WithRequired().
			WithHelpText("Order total"),

		// Coupon
		schema.Int64("coupon_id").WithOptional(),
		schema.String("coupon_code").WithMaxLength(50).WithOptional(),
		schema.Float64("coupon_discount").WithDefault(0.0),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Status: pending, processing, shipped, delivered, cancelled, refunded"),
		schema.String("payment_status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Payment status: pending, paid, failed, refunded"),
		schema.String("fulfillment_status").WithRequired().WithMaxLength(20).WithDefault("unfulfilled").
			WithHelpText("Fulfillment: unfulfilled, partial, fulfilled"),

		// Addresses
		schema.Int64("shipping_address_id").WithOptional(),
		schema.Int64("billing_address_id").WithOptional(),

		// Shipping address snapshot
		schema.String("shipping_first_name").WithMaxLength(100).WithOptional(),
		schema.String("shipping_last_name").WithMaxLength(100).WithOptional(),
		schema.String("shipping_company").WithMaxLength(200).WithOptional(),
		schema.String("shipping_address_line1").WithMaxLength(255).WithOptional(),
		schema.String("shipping_address_line2").WithMaxLength(255).WithOptional(),
		schema.String("shipping_city").WithMaxLength(100).WithOptional(),
		schema.String("shipping_state").WithMaxLength(100).WithOptional(),
		schema.String("shipping_postal_code").WithMaxLength(20).WithOptional(),
		schema.String("shipping_country_code").WithMaxLength(2).WithOptional(),
		schema.String("shipping_country_name").WithMaxLength(100).WithOptional(),
		schema.String("shipping_phone").WithMaxLength(20).WithOptional(),

		// Billing address snapshot
		schema.String("billing_first_name").WithMaxLength(100).WithOptional(),
		schema.String("billing_last_name").WithMaxLength(100).WithOptional(),
		schema.String("billing_company").WithMaxLength(200).WithOptional(),
		schema.String("billing_address_line1").WithMaxLength(255).WithOptional(),
		schema.String("billing_address_line2").WithMaxLength(255).WithOptional(),
		schema.String("billing_city").WithMaxLength(100).WithOptional(),
		schema.String("billing_state").WithMaxLength(100).WithOptional(),
		schema.String("billing_postal_code").WithMaxLength(20).WithOptional(),
		schema.String("billing_country_code").WithMaxLength(2).WithOptional(),
		schema.String("billing_country_name").WithMaxLength(100).WithOptional(),

		// Payment
		schema.String("payment_method").WithMaxLength(50).WithOptional().
			WithHelpText("Payment method: credit_card, paypal, stripe, etc."),
		schema.String("payment_transaction_id").WithMaxLength(255).WithOptional(),

		// Shipping
		schema.String("shipping_method").WithMaxLength(100).WithOptional(),
		schema.String("tracking_number").WithMaxLength(255).WithOptional(),
		schema.String("carrier").WithMaxLength(100).WithOptional(),

		// Additional information
		schema.Text("customer_notes").WithOptional().
			WithHelpText("Notes from customer"),
		schema.Text("admin_notes").WithOptional().
			WithHelpText("Internal notes"),

		// IP and user agent (for fraud detection)
		schema.String("ip_address").WithMaxLength(45).WithOptional(),
		schema.String("user_agent").WithMaxLength(500).WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("paid_at").WithOptional(),
		schema.Time("shipped_at").WithOptional(),
		schema.Time("delivered_at").WithOptional(),
		schema.Time("cancelled_at").WithOptional(),
		schema.Time("expires_at").WithOptional().
			WithHelpText("Order expiry for pending orders"),
	}
}

func (Order) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "orders",
		VerboseName:       "Order",
		VerboseNamePlural: "Orders",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_order_number", Fields: []string{"order_number"}, Unique: true},
			{Name: "idx_order_customer", Fields: []string{"customer_id"}},
			{Name: "idx_order_status", Fields: []string{"status"}},
			{Name: "idx_order_payment_status", Fields: []string{"payment_status"}},
			{Name: "idx_order_created", Fields: []string{"created_at"}},
			{Name: "idx_order_email", Fields: []string{"customer_email"}},
		},
	}
}

func (Order) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("orders"),
		schema.ForeignKey("coupon_id", "Coupon").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("orders"),
		schema.ForeignKey("shipping_address_id", "Address").
			WithOnDelete(schema.CascadePROTECT).
			WithOptional().
			WithRelatedName("shipping_orders"),
		schema.ForeignKey("billing_address_id", "Address").
			WithOnDelete(schema.CascadePROTECT).
			WithOptional().
			WithRelatedName("billing_orders"),
	}
}

func (Order) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeCreate: func(ctx context.Context, instance interface{}) error {
			// Generate order number
			// Snapshot customer info
			// Snapshot addresses
			// Calculate totals
			// Set expiry date
			return nil
		},
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Send order confirmation email
			// Reserve inventory
			// Update customer stats
			return nil
		},
		BeforeUpdate: func(ctx context.Context, instance interface{}) error {
			// Handle status changes
			// Update timestamps based on status
			return nil
		},
	}
}

// OrderItem represents line items in an order
type OrderItem struct {
	schema.BaseSchema
}

func (OrderItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("order_id").WithRequired(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("variant_id").WithOptional(),

		// Product snapshot (prices can change)
		schema.String("product_name").WithRequired().WithMaxLength(255),
		schema.String("product_sku").WithRequired().WithMaxLength(100),
		schema.String("variant_name").WithMaxLength(255).WithOptional(),
		schema.String("variant_sku").WithMaxLength(100).WithOptional(),
		schema.String("image_url").WithMaxLength(500).WithOptional(),

		// Quantity and pricing
		schema.Int32("quantity").WithRequired().WithDefault(1),
		schema.Float64("unit_price").WithRequired().
			WithHelpText("Price per unit at time of order"),
		schema.Float64("discount_amount").WithDefault(0.0),
		schema.Float64("tax_amount").WithDefault(0.0),
		schema.Float64("total").WithRequired().
			WithHelpText("Line total"),

		// Fulfillment
		schema.Int32("quantity_fulfilled").WithDefault(0),
		schema.Int32("quantity_refunded").WithDefault(0),
		schema.String("fulfillment_status").WithRequired().WithMaxLength(20).WithDefault("unfulfilled").
			WithHelpText("Status: unfulfilled, fulfilled, refunded"),

		// Physical attributes (for shipping calculations)
		schema.Float64("weight").WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (OrderItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "order_items",
		VerboseName:       "Order Item",
		VerboseNamePlural: "Order Items",
		OrderBy:           []string{"created_at"},
		Indexes: []schema.Index{
			{Name: "idx_order_item_order", Fields: []string{"order_id"}},
			{Name: "idx_order_item_product", Fields: []string{"product_id"}},
			{Name: "idx_order_item_variant", Fields: []string{"variant_id"}},
		},
	}
}

func (OrderItem) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("items"),
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("order_items"),
		schema.ForeignKey("variant_id", "ProductVariant").
			WithOnDelete(schema.CascadePROTECT).
			WithOptional().
			WithRelatedName("order_items"),
	}
}

func (OrderItem) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Calculate line total
			return nil
		},
		AfterSave: func(ctx context.Context, instance interface{}) error {
			// Update order totals
			// Update product order_count
			return nil
		},
	}
}

// Payment represents payment transactions
type Payment struct {
	schema.BaseSchema
}

func (Payment) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("order_id").WithRequired(),

		// Payment identification
		schema.String("transaction_id").WithMaxLength(255).WithOptional().
			WithHelpText("External payment gateway transaction ID"),

		// Amount
		schema.Float64("amount").WithRequired().
			WithHelpText("Payment amount"),
		schema.String("currency").WithRequired().WithMaxLength(3).WithDefault("USD"),

		// Payment method
		schema.String("payment_method").WithRequired().WithMaxLength(50).
			WithHelpText("Method: credit_card, paypal, stripe, bank_transfer, etc."),
		schema.String("payment_gateway").WithMaxLength(50).WithOptional().
			WithHelpText("Gateway used: stripe, paypal, square, etc."),

		// Card details (last 4 digits only for security)
		schema.String("card_last4").WithMaxLength(4).WithOptional(),
		schema.String("card_brand").WithMaxLength(50).WithOptional(),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Status: pending, completed, failed, refunded, cancelled"),

		// Additional information
		schema.Text("gateway_response").WithOptional().
			WithHelpText("Full gateway response (JSON)"),
		schema.String("failure_reason").WithMaxLength(500).WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("completed_at").WithOptional(),
		schema.Time("failed_at").WithOptional(),
		schema.Time("refunded_at").WithOptional(),
	}
}

func (Payment) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "payments",
		VerboseName:       "Payment",
		VerboseNamePlural: "Payments",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_payment_order", Fields: []string{"order_id"}},
			{Name: "idx_payment_transaction", Fields: []string{"transaction_id"}},
			{Name: "idx_payment_status", Fields: []string{"status"}},
		},
	}
}

func (Payment) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("payments"),
	}
}

func (Payment) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterUpdate: func(ctx context.Context, instance interface{}) error {
			// Update order payment status
			// Send payment notification
			return nil
		},
	}
}

// Shipment represents order shipments
type Shipment struct {
	schema.BaseSchema
}

func (Shipment) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("order_id").WithRequired(),

		// Tracking
		schema.String("tracking_number").WithRequired().WithMaxLength(255),
		schema.String("carrier").WithRequired().WithMaxLength(100).
			WithHelpText("Carrier: USPS, UPS, FedEx, DHL, etc."),
		schema.String("service_level").WithMaxLength(100).WithOptional().
			WithHelpText("Service: Standard, Express, Overnight, etc."),

		// Shipping address (snapshot)
		schema.String("recipient_name").WithRequired().WithMaxLength(200),
		schema.String("address_line1").WithRequired().WithMaxLength(255),
		schema.String("address_line2").WithMaxLength(255).WithOptional(),
		schema.String("city").WithRequired().WithMaxLength(100),
		schema.String("state").WithMaxLength(100).WithOptional(),
		schema.String("postal_code").WithRequired().WithMaxLength(20),
		schema.String("country_code").WithRequired().WithMaxLength(2),
		schema.String("country_name").WithRequired().WithMaxLength(100),
		schema.String("phone").WithMaxLength(20).WithOptional(),

		// Shipping details
		schema.Float64("weight").WithOptional().
			WithHelpText("Total weight in kg"),
		schema.Float64("shipping_cost").WithRequired(),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Status: pending, in_transit, delivered, failed, returned"),

		// Tracking events
		schema.Text("tracking_events").WithOptional().
			WithHelpText("JSON array of tracking events"),

		// Notes
		schema.Text("notes").WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("shipped_at").WithOptional(),
		schema.Time("estimated_delivery_at").WithOptional(),
		schema.Time("delivered_at").WithOptional(),
	}
}

func (Shipment) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "shipments",
		VerboseName:       "Shipment",
		VerboseNamePlural: "Shipments",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_shipment_order", Fields: []string{"order_id"}},
			{Name: "idx_shipment_tracking", Fields: []string{"tracking_number"}},
			{Name: "idx_shipment_status", Fields: []string{"status"}},
		},
	}
}

func (Shipment) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("shipments"),
	}
}

func (Shipment) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Send shipping notification
			// Update order status
			return nil
		},
		AfterUpdate: func(ctx context.Context, instance interface{}) error {
			// Send tracking updates
			// Update order delivery status
			return nil
		},
	}
}

// RegisterModels registers order models with the framework
func RegisterModels() {
	registry.RegisterModel(&Cart{})
	registry.RegisterModel(&CartItem{})
	registry.RegisterModel(&Order{})
	registry.RegisterModel(&OrderItem{})
	registry.RegisterModel(&Payment{})
	registry.RegisterModel(&Shipment{})
}
