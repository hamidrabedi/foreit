package orders

import (
	"context"
	
	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/db"
)

// RegisterAdmin registers order models with the admin interface
func RegisterAdmin(ctx context.Context, registry *admin.Registry, database *db.DB) {
	// Cart admin
	cartConfig := &admin.ModelConfig{
		Name:             "Shopping Cart",
		PluralName:       "Shopping Carts",
		Icon:             "🛒",
		ListDisplay:      []string{"id", "customer_id", "guest_email", "total", "status", "is_abandoned", "last_activity_at", "created_at"},
		ListFilter:       []string{"status", "is_abandoned"},
		SearchFields:     []string{"customer_id", "guest_email", "session_id"},
		OrderBy:          []string{"-updated_at"},
		PerPage:          20,
		Actions:          []string{"delete", "export"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("Cart", &Cart{}, cartConfig)
	
	// CartItem admin
	cartItemConfig := &admin.ModelConfig{
		Name:             "Cart Item",
		PluralName:       "Cart Items",
		Icon:             "🛍️",
		ListDisplay:      []string{"id", "cart_id", "product_name", "variant_name", "quantity", "unit_price", "total"},
		ListFilter:       []string{"cart_id"},
		SearchFields:     []string{"product_name", "variant_name"},
		OrderBy:          []string{"cart_id", "created_at"},
		PerPage:          20,
		Actions:          []string{"delete"},
	}
	registry.Register("CartItem", &CartItem{}, cartItemConfig)
	
	// Order admin
	orderConfig := &admin.ModelConfig{
		Name:             "Order",
		PluralName:       "Orders",
		Icon:             "📋",
		ListDisplay:      []string{"id", "order_number", "customer_email", "total", "status", "payment_status", "fulfillment_status", "created_at"},
		ListFilter:       []string{"status", "payment_status", "fulfillment_status"},
		SearchFields:     []string{"order_number", "customer_email", "customer_first_name", "customer_last_name"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"delete", "export", "print_invoice", "mark_paid", "mark_shipped"},
		ExportFormats:    []string{"csv", "json", "pdf"},
		BulkActions:      true,
	}
	registry.Register("Order", &Order{}, orderConfig)
	
	// OrderItem admin
	orderItemConfig := &admin.ModelConfig{
		Name:             "Order Item",
		PluralName:       "Order Items",
		Icon:             "📦",
		ListDisplay:      []string{"id", "order_id", "product_name", "variant_name", "quantity", "unit_price", "total", "fulfillment_status"},
		ListFilter:       []string{"fulfillment_status", "order_id"},
		SearchFields:     []string{"product_name", "product_sku", "variant_sku"},
		OrderBy:          []string{"order_id", "created_at"},
		PerPage:          20,
		Actions:          []string{"fulfill", "refund"},
	}
	registry.Register("OrderItem", &OrderItem{}, orderItemConfig)
	
	// Payment admin
	paymentConfig := &admin.ModelConfig{
		Name:             "Payment",
		PluralName:       "Payments",
		Icon:             "💳",
		ListDisplay:      []string{"id", "order_id", "amount", "currency", "payment_method", "status", "transaction_id", "created_at"},
		ListFilter:       []string{"status", "payment_method", "payment_gateway"},
		SearchFields:     []string{"transaction_id", "order_id"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"refund", "export"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("Payment", &Payment{}, paymentConfig)
	
	// Shipment admin
	shipmentConfig := &admin.ModelConfig{
		Name:             "Shipment",
		PluralName:       "Shipments",
		Icon:             "🚚",
		ListDisplay:      []string{"id", "order_id", "tracking_number", "carrier", "status", "shipped_at", "estimated_delivery_at", "delivered_at"},
		ListFilter:       []string{"status", "carrier"},
		SearchFields:     []string{"tracking_number", "order_id", "recipient_name"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"delete", "print_label", "track"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("Shipment", &Shipment{}, shipmentConfig)
}
