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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("customer_id").Optional().
			HelpText("Customer (null for guest carts)").Build(),
		
		// Guest cart identification
		schema.String("session_id").MaxLength(255).Optional().
			HelpText("Session ID for guest carts").Build(),
		schema.String("guest_email").MaxLength(255).Optional().Build(),
		
		// Pricing
		schema.Float64("subtotal").Default(0.0).
			HelpText("Cart subtotal before discounts").Build(),
		schema.Float64("discount_amount").Default(0.0).Build(),
		schema.Float64("tax_amount").Default(0.0).Build(),
		schema.Float64("shipping_amount").Default(0.0).Build(),
		schema.Float64("total").Default(0.0).
			HelpText("Cart total").Build(),
		
		// Coupon
		schema.Int64("coupon_id").Optional().Build(),
		schema.String("coupon_code").MaxLength(50).Optional().Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("active").
			HelpText("Status: active, abandoned, converted").Build(),
		schema.Bool("is_abandoned").Default(false).Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("last_activity_at").Optional().Build(),
		schema.Time("converted_at").Optional().
			HelpText("When cart was converted to order").Build(),
	}
}

func (Cart) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "carts",
		VerboseName:      "Shopping Cart",
		VerboseNamePlural: "Shopping Carts",
		OrderBy:          []string{"-updated_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Optional().
			RelatedName("carts").Build(),
		schema.ForeignKey("coupon_id", "Coupon").
			OnDelete(schema.CascadeSET_NULL).
			Optional().
			RelatedName("carts").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("cart_id").Required().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("variant_id").Optional().Build(),
		
		// Quantity and pricing
		schema.Int32("quantity").Required().Default(1).Build(),
		schema.Float64("unit_price").Required().
			HelpText("Price per unit at time of adding").Build(),
		schema.Float64("discount_amount").Default(0.0).Build(),
		schema.Float64("tax_amount").Default(0.0).Build(),
		schema.Float64("total").Required().
			HelpText("Line total (quantity * unit_price - discount + tax)").Build(),
		
		// Product snapshot (for price changes)
		schema.String("product_name").MaxLength(255).Optional().Build(),
		schema.String("variant_name").MaxLength(255).Optional().Build(),
		schema.String("image_url").MaxLength(500).Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (CartItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "cart_items",
		VerboseName:      "Cart Item",
		VerboseNamePlural: "Cart Items",
		OrderBy:          []string{"created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("items").Build(),
		schema.ForeignKey("product_id", "Product").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("cart_items").Build(),
		schema.ForeignKey("variant_id", "ProductVariant").
			OnDelete(schema.CascadeCASCADE).
			Optional().
			RelatedName("cart_items").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		
		// Order identification
		schema.String("order_number").Required().MaxLength(50).Unique().
			HelpText("Human-readable order number").Build(),
		
		// Customer
		schema.Int64("customer_id").Required().Build(),
		schema.String("customer_email").Required().MaxLength(255).Build(),
		schema.String("customer_first_name").Required().MaxLength(100).Build(),
		schema.String("customer_last_name").Required().MaxLength(100).Build(),
		schema.String("customer_phone").MaxLength(20).Optional().Build(),
		
		// Pricing
		schema.Float64("subtotal").Required().
			HelpText("Order subtotal before discounts").Build(),
		schema.Float64("discount_amount").Default(0.0).Build(),
		schema.Float64("tax_amount").Required().Build(),
		schema.Float64("shipping_amount").Required().Build(),
		schema.Float64("total").Required().
			HelpText("Order total").Build(),
		
		// Coupon
		schema.Int64("coupon_id").Optional().Build(),
		schema.String("coupon_code").MaxLength(50).Optional().Build(),
		schema.Float64("coupon_discount").Default(0.0).Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("pending").
			HelpText("Status: pending, processing, shipped, delivered, cancelled, refunded").Build(),
		schema.String("payment_status").Required().MaxLength(20).Default("pending").
			HelpText("Payment status: pending, paid, failed, refunded").Build(),
		schema.String("fulfillment_status").Required().MaxLength(20).Default("unfulfilled").
			HelpText("Fulfillment: unfulfilled, partial, fulfilled").Build(),
		
		// Addresses
		schema.Int64("shipping_address_id").Optional().Build(),
		schema.Int64("billing_address_id").Optional().Build(),
		
		// Shipping address snapshot
		schema.String("shipping_first_name").MaxLength(100).Optional().Build(),
		schema.String("shipping_last_name").MaxLength(100).Optional().Build(),
		schema.String("shipping_company").MaxLength(200).Optional().Build(),
		schema.String("shipping_address_line1").MaxLength(255).Optional().Build(),
		schema.String("shipping_address_line2").MaxLength(255).Optional().Build(),
		schema.String("shipping_city").MaxLength(100).Optional().Build(),
		schema.String("shipping_state").MaxLength(100).Optional().Build(),
		schema.String("shipping_postal_code").MaxLength(20).Optional().Build(),
		schema.String("shipping_country_code").MaxLength(2).Optional().Build(),
		schema.String("shipping_country_name").MaxLength(100).Optional().Build(),
		schema.String("shipping_phone").MaxLength(20).Optional().Build(),
		
		// Billing address snapshot
		schema.String("billing_first_name").MaxLength(100).Optional().Build(),
		schema.String("billing_last_name").MaxLength(100).Optional().Build(),
		schema.String("billing_company").MaxLength(200).Optional().Build(),
		schema.String("billing_address_line1").MaxLength(255).Optional().Build(),
		schema.String("billing_address_line2").MaxLength(255).Optional().Build(),
		schema.String("billing_city").MaxLength(100).Optional().Build(),
		schema.String("billing_state").MaxLength(100).Optional().Build(),
		schema.String("billing_postal_code").MaxLength(20).Optional().Build(),
		schema.String("billing_country_code").MaxLength(2).Optional().Build(),
		schema.String("billing_country_name").MaxLength(100).Optional().Build(),
		
		// Payment
		schema.String("payment_method").MaxLength(50).Optional().
			HelpText("Payment method: credit_card, paypal, stripe, etc.").Build(),
		schema.String("payment_transaction_id").MaxLength(255).Optional().Build(),
		
		// Shipping
		schema.String("shipping_method").MaxLength(100).Optional().Build(),
		schema.String("tracking_number").MaxLength(255).Optional().Build(),
		schema.String("carrier").MaxLength(100).Optional().Build(),
		
		// Additional information
		schema.Text("customer_notes").Optional().
			HelpText("Notes from customer").Build(),
		schema.Text("admin_notes").Optional().
			HelpText("Internal notes").Build(),
		
		// IP and user agent (for fraud detection)
		schema.String("ip_address").MaxLength(45).Optional().Build(),
		schema.String("user_agent").MaxLength(500).Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("paid_at").Optional().Build(),
		schema.Time("shipped_at").Optional().Build(),
		schema.Time("delivered_at").Optional().Build(),
		schema.Time("cancelled_at").Optional().Build(),
		schema.Time("expires_at").Optional().
			HelpText("Order expiry for pending orders").Build(),
	}
}

func (Order) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "orders",
		VerboseName:      "Order",
		VerboseNamePlural: "Orders",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadePROTECT).
			Required().
			RelatedName("orders").Build(),
		schema.ForeignKey("coupon_id", "Coupon").
			OnDelete(schema.CascadeSET_NULL).
			Optional().
			RelatedName("orders").Build(),
		schema.ForeignKey("shipping_address_id", "Address").
			OnDelete(schema.CascadePROTECT).
			Optional().
			RelatedName("shipping_orders").Build(),
		schema.ForeignKey("billing_address_id", "Address").
			OnDelete(schema.CascadePROTECT).
			Optional().
			RelatedName("billing_orders").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("order_id").Required().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("variant_id").Optional().Build(),
		
		// Product snapshot (prices can change)
		schema.String("product_name").Required().MaxLength(255).Build(),
		schema.String("product_sku").Required().MaxLength(100).Build(),
		schema.String("variant_name").MaxLength(255).Optional().Build(),
		schema.String("variant_sku").MaxLength(100).Optional().Build(),
		schema.String("image_url").MaxLength(500).Optional().Build(),
		
		// Quantity and pricing
		schema.Int32("quantity").Required().Default(1).Build(),
		schema.Float64("unit_price").Required().
			HelpText("Price per unit at time of order").Build(),
		schema.Float64("discount_amount").Default(0.0).Build(),
		schema.Float64("tax_amount").Default(0.0).Build(),
		schema.Float64("total").Required().
			HelpText("Line total").Build(),
		
		// Fulfillment
		schema.Int32("quantity_fulfilled").Default(0).Build(),
		schema.Int32("quantity_refunded").Default(0).Build(),
		schema.String("fulfillment_status").Required().MaxLength(20).Default("unfulfilled").
			HelpText("Status: unfulfilled, fulfilled, refunded").Build(),
		
		// Physical attributes (for shipping calculations)
		schema.Float64("weight").Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (OrderItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "order_items",
		VerboseName:      "Order Item",
		VerboseNamePlural: "Order Items",
		OrderBy:          []string{"created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("items").Build(),
		schema.ForeignKey("product_id", "Product").
			OnDelete(schema.CascadePROTECT).
			Required().
			RelatedName("order_items").Build(),
		schema.ForeignKey("variant_id", "ProductVariant").
			OnDelete(schema.CascadePROTECT).
			Optional().
			RelatedName("order_items").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("order_id").Required().Build(),
		
		// Payment identification
		schema.String("transaction_id").MaxLength(255).Optional().
			HelpText("External payment gateway transaction ID").Build(),
		
		// Amount
		schema.Float64("amount").Required().
			HelpText("Payment amount").Build(),
		schema.String("currency").Required().MaxLength(3).Default("USD").Build(),
		
		// Payment method
		schema.String("payment_method").Required().MaxLength(50).
			HelpText("Method: credit_card, paypal, stripe, bank_transfer, etc.").Build(),
		schema.String("payment_gateway").MaxLength(50).Optional().
			HelpText("Gateway used: stripe, paypal, square, etc.").Build(),
		
		// Card details (last 4 digits only for security)
		schema.String("card_last4").MaxLength(4).Optional().Build(),
		schema.String("card_brand").MaxLength(50).Optional().Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("pending").
			HelpText("Status: pending, completed, failed, refunded, cancelled").Build(),
		
		// Additional information
		schema.Text("gateway_response").Optional().
			HelpText("Full gateway response (JSON)").Build(),
		schema.String("failure_reason").MaxLength(500).Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("completed_at").Optional().Build(),
		schema.Time("failed_at").Optional().Build(),
		schema.Time("refunded_at").Optional().Build(),
	}
}

func (Payment) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "payments",
		VerboseName:      "Payment",
		VerboseNamePlural: "Payments",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("payments").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("order_id").Required().Build(),
		
		// Tracking
		schema.String("tracking_number").Required().MaxLength(255).Build(),
		schema.String("carrier").Required().MaxLength(100).
			HelpText("Carrier: USPS, UPS, FedEx, DHL, etc.").Build(),
		schema.String("service_level").MaxLength(100).Optional().
			HelpText("Service: Standard, Express, Overnight, etc.").Build(),
		
		// Shipping address (snapshot)
		schema.String("recipient_name").Required().MaxLength(200).Build(),
		schema.String("address_line1").Required().MaxLength(255).Build(),
		schema.String("address_line2").MaxLength(255).Optional().Build(),
		schema.String("city").Required().MaxLength(100).Build(),
		schema.String("state").MaxLength(100).Optional().Build(),
		schema.String("postal_code").Required().MaxLength(20).Build(),
		schema.String("country_code").Required().MaxLength(2).Build(),
		schema.String("country_name").Required().MaxLength(100).Build(),
		schema.String("phone").MaxLength(20).Optional().Build(),
		
		// Shipping details
		schema.Float64("weight").Optional().
			HelpText("Total weight in kg").Build(),
		schema.Float64("shipping_cost").Required().Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("pending").
			HelpText("Status: pending, in_transit, delivered, failed, returned").Build(),
		
		// Tracking events
		schema.Text("tracking_events").Optional().
			HelpText("JSON array of tracking events").Build(),
		
		// Notes
		schema.Text("notes").Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("shipped_at").Optional().Build(),
		schema.Time("estimated_delivery_at").Optional().Build(),
		schema.Time("delivered_at").Optional().Build(),
	}
}

func (Shipment) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "shipments",
		VerboseName:      "Shipment",
		VerboseNamePlural: "Shipments",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("shipments").Build(),
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
