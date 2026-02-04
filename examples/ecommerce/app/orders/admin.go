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
			CartFieldsInstance.Id,
			CartFieldsInstance.GuestEmail,
			CartFieldsInstance.Total,
			CartFieldsInstance.Status,
			CartFieldsInstance.IsAbandoned,
			CartFieldsInstance.LastActivityAt,
			CartFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			CartFieldsInstance.Status,
			CartFieldsInstance.IsAbandoned,
		},
	})

	// CartItem admin
	admin.Register(&admin.Config[CartItem]{
		Icon: "ShoppingBag",
		ListDisplay: []admin.Field{
			CartItemFieldsInstance.CartId,
			CartItemFieldsInstance.ProductName,
			CartItemFieldsInstance.VariantName,
			CartItemFieldsInstance.Quantity,
			CartItemFieldsInstance.UnitPrice,
			CartItemFieldsInstance.Total,
		},
		ListFilter: []admin.Field{
			CartItemFieldsInstance.CartId,
		},
	})

	// Order admin
	admin.Register(&admin.Config[Order]{
		Icon: "ClipboardList",
		ListDisplay: []admin.Field{
			OrderFieldsInstance.OrderNumber,
			OrderFieldsInstance.CustomerEmail,
			OrderFieldsInstance.Total,
			OrderFieldsInstance.Status,
			OrderFieldsInstance.PaymentStatus,
			OrderFieldsInstance.FulfillmentStatus,
			OrderFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			OrderFieldsInstance.Status,
			OrderFieldsInstance.PaymentStatus,
			OrderFieldsInstance.FulfillmentStatus,
		},
		SearchFields: []admin.Field{
			OrderFieldsInstance.OrderNumber,
			OrderFieldsInstance.CustomerEmail,
		},
	})

	// OrderItem admin
	admin.Register(&admin.Config[OrderItem]{
		Icon: "Package",
		ListDisplay: []admin.Field{
			OrderItemFieldsInstance.OrderId,
			OrderItemFieldsInstance.ProductName,
			OrderItemFieldsInstance.Quantity,
			OrderItemFieldsInstance.UnitPrice,
			OrderItemFieldsInstance.Total,
			OrderItemFieldsInstance.FulfillmentStatus,
		},
	})

	// Payment admin
	admin.Register(&admin.Config[Payment]{
		Icon: "CreditCard",
		ListDisplay: []admin.Field{
			PaymentFieldsInstance.OrderId,
			PaymentFieldsInstance.Amount,
			PaymentFieldsInstance.Currency,
			PaymentFieldsInstance.PaymentMethod,
			PaymentFieldsInstance.Status,
			PaymentFieldsInstance.TransactionId,
		},
		ListFilter: []admin.Field{
			PaymentFieldsInstance.Status,
			PaymentFieldsInstance.PaymentMethod,
		},
	})

	// Shipment admin
	admin.Register(&admin.Config[Shipment]{
		Icon: "Truck",
		ListDisplay: []admin.Field{
			ShipmentFieldsInstance.OrderId,
			ShipmentFieldsInstance.TrackingNumber,
			ShipmentFieldsInstance.Carrier,
			ShipmentFieldsInstance.Status,
			ShipmentFieldsInstance.ShippedAt,
		},
		ListFilter: []admin.Field{
			ShipmentFieldsInstance.Status,
			ShipmentFieldsInstance.Carrier,
		},
	})
}
