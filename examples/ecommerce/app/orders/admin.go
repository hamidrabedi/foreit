package orders

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers order models with the admin interface
func RegisterAdmin(ctx context.Context) {
	// Cart admin
	admin.Register(&admin.Config[Cart]{
		Icon: "ShoppingCart",
		ListDisplay: []admin.Field{
			CartFields.ID,
			CartFields.GuestEmail,
			CartFields.Total,
			CartFields.Status,
			CartFields.IsAbandoned,
			CartFields.LastActivityAt,
			CartFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			CartFields.Status,
			CartFields.IsAbandoned,
		},
	})

	// CartItem admin
	admin.Register(&admin.Config[CartItem]{
		Icon: "ShoppingBag",
		ListDisplay: []admin.Field{
			CartItemFields.CartID,
			CartItemFields.ProductName,
			CartItemFields.VariantName,
			CartItemFields.Quantity,
			CartItemFields.UnitPrice,
			CartItemFields.Total,
		},
		ListFilter: []admin.Field{
			CartItemFields.CartID,
		},
	})

	// Order admin
	admin.Register(&admin.Config[Order]{
		Icon: "ClipboardList",
		ListDisplay: []admin.Field{
			OrderFields.OrderNumber,
			OrderFields.CustomerEmail,
			OrderFields.Total,
			OrderFields.Status,
			OrderFields.PaymentStatus,
			OrderFields.FulfillmentStatus,
			OrderFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			OrderFields.Status,
			OrderFields.PaymentStatus,
			OrderFields.FulfillmentStatus,
		},
		SearchFields: []admin.Field{
			OrderFields.OrderNumber,
			OrderFields.CustomerEmail,
		},
	})

	// OrderItem admin
	admin.Register(&admin.Config[OrderItem]{
		Icon: "Package",
		ListDisplay: []admin.Field{
			OrderItemFields.OrderID,
			OrderItemFields.ProductName,
			OrderItemFields.Quantity,
			OrderItemFields.UnitPrice,
			OrderItemFields.Total,
			OrderItemFields.FulfillmentStatus,
		},
	})

	// Payment admin
	admin.Register(&admin.Config[Payment]{
		Icon: "CreditCard",
		ListDisplay: []admin.Field{
			PaymentFields.OrderID,
			PaymentFields.Amount,
			PaymentFields.Currency,
			PaymentFields.PaymentMethod,
			PaymentFields.Status,
			PaymentFields.TransactionID,
		},
		ListFilter: []admin.Field{
			PaymentFields.Status,
			PaymentFields.PaymentMethod,
		},
	})

	// Shipment admin
	admin.Register(&admin.Config[Shipment]{
		Icon: "Truck",
		ListDisplay: []admin.Field{
			ShipmentFields.OrderID,
			ShipmentFields.TrackingNumber,
			ShipmentFields.Carrier,
			ShipmentFields.Status,
			ShipmentFields.ShippedAt,
		},
		ListFilter: []admin.Field{
			ShipmentFields.Status,
			ShipmentFields.Carrier,
		},
	})
}
