package inventory

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers inventory models with the admin interface
func RegisterAdmin(ctx context.Context) {
	// Warehouse admin
	admin.Register(&admin.Config[Warehouse]{
		Icon: "Warehouse",
		ListDisplay: []admin.Field{
			WarehouseFields.Name,
			WarehouseFields.Code,
			WarehouseFields.City,
			WarehouseFields.CountryCode,
			WarehouseFields.IsActive,
			WarehouseFields.IsPrimary,
		},
		ListFilter: []admin.Field{
			WarehouseFields.IsActive,
			WarehouseFields.IsPrimary,
			WarehouseFields.CountryCode,
		},
	})

	// Stock admin
	admin.Register(&admin.Config[Stock]{
		Icon: "Activity",
		ListDisplay: []admin.Field{
			StockFields.ProductVariantID,
			StockFields.WarehouseID,
			StockFields.Quantity,
			StockFields.ReservedQuantity,
			StockFields.IsActive,
		},
		ListFilter: []admin.Field{
			StockFields.WarehouseID,
			StockFields.IsActive,
		},
	})

	// StockMovement admin
	admin.Register(&admin.Config[StockMovement]{
		Icon: "TrendingUp",
		ListDisplay: []admin.Field{
			StockMovementFields.ProductVariantID,
			StockMovementFields.WarehouseID,
			StockMovementFields.Type,
			StockMovementFields.Quantity,
			StockMovementFields.MovementDate,
		},
		ListFilter: []admin.Field{
			StockMovementFields.Type,
			StockMovementFields.WarehouseID,
			StockMovementFields.MovementDate,
		},
	})

	// StockAlert admin
	admin.Register(&admin.Config[StockAlert]{
		Icon: "AlertTriangle",
		ListDisplay: []admin.Field{
			StockAlertFields.ProductVariantID,
			StockAlertFields.WarehouseID,
			StockAlertFields.AlertType,
			StockAlertFields.CurrentQuantity,
			StockAlertFields.Status,
		},
		ListFilter: []admin.Field{
			StockAlertFields.AlertType,
			StockAlertFields.Status,
			StockAlertFields.WarehouseID,
		},
	})

	// StockTransfer admin
	admin.Register(&admin.Config[StockTransfer]{
		Icon: "Repeat",
		ListDisplay: []admin.Field{
			StockTransferFields.TransferNumber,
			StockTransferFields.FromWarehouseID,
			StockTransferFields.ToWarehouseID,
			StockTransferFields.Quantity,
			StockTransferFields.Status,
		},
		ListFilter: []admin.Field{
			StockTransferFields.Status,
			StockTransferFields.FromWarehouseID,
			StockTransferFields.ToWarehouseID,
		},
	})
}
