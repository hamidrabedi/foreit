package inventory

import (
	"context"
	
	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/db"
)

// RegisterAdmin registers inventory models with the admin interface
func RegisterAdmin(ctx context.Context, registry *admin.Registry, database *db.DB) {
	// Warehouse admin
	warehouseConfig := &admin.ModelConfig{
		Name:             "Warehouse",
		PluralName:       "Warehouses",
		Icon:             "🏭",
		ListDisplay:      []string{"id", "name", "code", "city", "country_code", "is_active", "is_primary", "priority"},
		ListFilter:       []string{"is_active", "is_primary", "country_code"},
		SearchFields:     []string{"name", "code", "city", "address_line1"},
		OrderBy:          []string{"-is_primary", "-priority", "name"},
		PerPage:          20,
		Actions:          []string{"delete", "activate", "deactivate"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("Warehouse", &Warehouse{}, warehouseConfig)
	
	// Stock admin
	stockConfig := &admin.ModelConfig{
		Name:             "Stock",
		PluralName:       "Stock",
		Icon:             "📊",
		ListDisplay:      []string{"id", "product_variant_id", "warehouse_id", "quantity", "reserved_quantity", "available_quantity", "reorder_point", "is_active"},
		ListFilter:       []string{"warehouse_id", "is_active"},
		SearchFields:     []string{"bin_location"},
		OrderBy:          []string{"warehouse_id", "quantity"},
		PerPage:          20,
		Actions:          []string{"delete", "adjust", "count", "export"},
		ExportFormats:    []string{"csv", "json", "xlsx"},
		BulkActions:      true,
	}
	registry.Register("Stock", &Stock{}, stockConfig)
	
	// StockMovement admin
	movementConfig := &admin.ModelConfig{
		Name:             "Stock Movement",
		PluralName:       "Stock Movements",
		Icon:             "📈",
		ListDisplay:      []string{"id", "product_variant_id", "warehouse_id", "type", "quantity", "quantity_before", "quantity_after", "reference_number", "movement_date"},
		ListFilter:       []string{"type", "warehouse_id", "movement_date"},
		SearchFields:     []string{"reference_number", "reference_type", "reason"},
		OrderBy:          []string{"-movement_date", "-created_at"},
		PerPage:          20,
		Actions:          []string{"export"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("StockMovement", &StockMovement{}, movementConfig)
	
	// StockAlert admin
	alertConfig := &admin.ModelConfig{
		Name:             "Stock Alert",
		PluralName:       "Stock Alerts",
		Icon:             "⚠️",
		ListDisplay:      []string{"id", "product_variant_id", "warehouse_id", "alert_type", "current_quantity", "threshold", "status", "created_at"},
		ListFilter:       []string{"alert_type", "status", "warehouse_id"},
		SearchFields:     []string{"product_variant_id"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"acknowledge", "resolve", "dismiss"},
		BulkActions:      true,
	}
	registry.Register("StockAlert", &StockAlert{}, alertConfig)
	
	// StockTransfer admin
	transferConfig := &admin.ModelConfig{
		Name:             "Stock Transfer",
		PluralName:       "Stock Transfers",
		Icon:             "↔️",
		ListDisplay:      []string{"id", "transfer_number", "from_warehouse_id", "to_warehouse_id", "product_variant_id", "quantity", "status", "created_at"},
		ListFilter:       []string{"status", "from_warehouse_id", "to_warehouse_id"},
		SearchFields:     []string{"transfer_number", "tracking_number"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"approve", "ship", "complete", "cancel"},
	}
	registry.Register("StockTransfer", &StockTransfer{}, transferConfig)
}
