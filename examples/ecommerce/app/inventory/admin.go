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
			WarehouseFieldsInstance.Name,
			WarehouseFieldsInstance.Code,
			WarehouseFieldsInstance.City,
			WarehouseFieldsInstance.CountryCode,
			WarehouseFieldsInstance.IsActive,
			WarehouseFieldsInstance.IsPrimary,
		},
		ListFilter: []admin.Field{
			WarehouseFieldsInstance.IsActive,
			WarehouseFieldsInstance.IsPrimary,
			WarehouseFieldsInstance.CountryCode,
		},
	})

	// Stock admin
	admin.Register(&admin.Config[Stock]{
		Icon: "Activity",
		ListDisplay: []admin.Field{
			StockFieldsInstance.ProductVariantId,
			StockFieldsInstance.WarehouseId,
			StockFieldsInstance.Quantity,
			StockFieldsInstance.ReservedQuantity,
			StockFieldsInstance.IsActive,
		},
		ListFilter: []admin.Field{
			StockFieldsInstance.WarehouseId,
			StockFieldsInstance.IsActive,
		},
	})

	// StockMovement admin
	admin.Register(&admin.Config[StockMovement]{
		Icon: "TrendingUp",
		ListDisplay: []admin.Field{
			StockMovementFieldsInstance.ProductVariantId,
			StockMovementFieldsInstance.WarehouseId,
			StockMovementFieldsInstance.Type,
			StockMovementFieldsInstance.Quantity,
			StockMovementFieldsInstance.MovementDate,
		},
		ListFilter: []admin.Field{
			StockMovementFieldsInstance.Type,
			StockMovementFieldsInstance.WarehouseId,
			StockMovementFieldsInstance.MovementDate,
		},
	})

	// StockAlert admin
	admin.Register(&admin.Config[StockAlert]{
		Icon: "AlertTriangle",
		ListDisplay: []admin.Field{
			StockAlertFieldsInstance.ProductVariantId,
			StockAlertFieldsInstance.WarehouseId,
			StockAlertFieldsInstance.AlertType,
			StockAlertFieldsInstance.CurrentQuantity,
			StockAlertFieldsInstance.Status,
		},
		ListFilter: []admin.Field{
			StockAlertFieldsInstance.AlertType,
			StockAlertFieldsInstance.Status,
			StockAlertFieldsInstance.WarehouseId,
		},
	})

	// StockTransfer admin
	admin.Register(&admin.Config[StockTransfer]{
		Icon: "Repeat",
		ListDisplay: []admin.Field{
			StockTransferFieldsInstance.TransferNumber,
			StockTransferFieldsInstance.FromWarehouseId,
			StockTransferFieldsInstance.ToWarehouseId,
			StockTransferFieldsInstance.Quantity,
			StockTransferFieldsInstance.Status,
		},
		ListFilter: []admin.Field{
			StockTransferFieldsInstance.Status,
			StockTransferFieldsInstance.FromWarehouseId,
			StockTransferFieldsInstance.ToWarehouseId,
		},
	})
}
