package inventory

import (
	"context"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Warehouse represents a physical warehouse or storage location
type Warehouse struct {
	schema.BaseSchema
}

func (Warehouse) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		
		// Identification
		schema.String("name").Required().MaxLength(200).Unique().Build(),
		schema.String("code").Required().MaxLength(50).Unique().
			HelpText("Warehouse code for reference").Build(),
		
		// Contact
		schema.String("contact_name").MaxLength(200).Optional().Build(),
		schema.String("contact_email").MaxLength(255).Optional().Build(),
		schema.String("contact_phone").MaxLength(20).Optional().Build(),
		
		// Address
		schema.String("address_line1").Required().MaxLength(255).Build(),
		schema.String("address_line2").MaxLength(255).Optional().Build(),
		schema.String("city").Required().MaxLength(100).Build(),
		schema.String("state").MaxLength(100).Optional().Build(),
		schema.String("postal_code").Required().MaxLength(20).Build(),
		schema.String("country_code").Required().MaxLength(2).Build(),
		schema.String("country_name").Required().MaxLength(100).Build(),
		
		// Geolocation
		schema.Float64("latitude").Optional().Build(),
		schema.Float64("longitude").Optional().Build(),
		
		// Settings
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("is_primary").Default(false).
			HelpText("Primary warehouse for fulfillment").Build(),
		schema.Int32("priority").Default(0).
			HelpText("Fulfillment priority (higher = preferred)").Build(),
		
		// Capacity
		schema.Int32("total_capacity").Optional().
			HelpText("Total storage capacity (units)").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Warehouse) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "warehouses",
		VerboseName:      "Warehouse",
		VerboseNamePlural: "Warehouses",
		OrderBy:          []string{"-is_primary", "-priority", "name"},
		Indexes: []schema.Index{
			{Name: "idx_warehouse_code", Fields: []string{"code"}, Unique: true},
			{Name: "idx_warehouse_active", Fields: []string{"is_active"}},
			{Name: "idx_warehouse_primary", Fields: []string{"is_primary"}},
		},
	}
}

func (Warehouse) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Warehouse) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Ensure only one primary warehouse
			return nil
		},
	}
}

// Stock represents inventory levels for product variants
type Stock struct {
	schema.BaseSchema
}

func (Stock) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("product_variant_id").Required().Build(),
		schema.Int64("warehouse_id").Required().Build(),
		
		// Quantities
		schema.Int32("quantity").Required().Default(0).
			HelpText("Available quantity").Build(),
		schema.Int32("reserved_quantity").Default(0).
			HelpText("Quantity reserved in pending orders").Build(),
		schema.Int32("available_quantity").Default(0).
			HelpText("quantity - reserved_quantity").Build(),
		
		// Reordering
		schema.Int32("reorder_point").Default(10).
			HelpText("Minimum quantity before reorder").Build(),
		schema.Int32("reorder_quantity").Default(50).
			HelpText("Quantity to reorder").Build(),
		
		// Location
		schema.String("bin_location").MaxLength(100).Optional().
			HelpText("Physical location in warehouse (e.g., A-12-3)").Build(),
		
		// Status
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("allow_backorder").Default(false).Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("last_counted_at").Optional().
			HelpText("Last physical count").Build(),
	}
}

func (Stock) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "stock",
		VerboseName:      "Stock",
		VerboseNamePlural: "Stock",
		OrderBy:          []string{"warehouse_id", "product_variant_id"},
		Indexes: []schema.Index{
			{Name: "idx_stock_variant", Fields: []string{"product_variant_id"}},
			{Name: "idx_stock_warehouse", Fields: []string{"warehouse_id"}},
			{Name: "idx_stock_quantity", Fields: []string{"quantity"}},
		},
		UniqueTogether: [][]string{
			{"product_variant_id", "warehouse_id"},
		},
	}
}

func (Stock) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_variant_id", "ProductVariant").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("stock_records").Build(),
		schema.ForeignKey("warehouse_id", "Warehouse").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("stock_records").Build(),
	}
}

func (Stock) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Calculate available_quantity
			// Check for low stock alert
			return nil
		},
		AfterSave: func(ctx context.Context, instance interface{}) error {
			// Update product variant total stock
			// Trigger low stock alert if needed
			return nil
		},
	}
}

// StockMovement represents inventory transactions
type StockMovement struct {
	schema.BaseSchema
}

func (StockMovement) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("stock_id").Required().Build(),
		schema.Int64("product_variant_id").Required().Build(),
		schema.Int64("warehouse_id").Required().Build(),
		
		// Movement details
		schema.String("type").Required().MaxLength(50).
			HelpText("Type: purchase, sale, return, adjustment, transfer, damage, loss, count").Build(),
		schema.Int32("quantity").Required().
			HelpText("Quantity change (positive or negative)").Build(),
		schema.Int32("quantity_before").Required().
			HelpText("Quantity before movement").Build(),
		schema.Int32("quantity_after").Required().
			HelpText("Quantity after movement").Build(),
		
		// Reference
		schema.String("reference_type").MaxLength(50).Optional().
			HelpText("Reference type: order, return, transfer, etc.").Build(),
		schema.Int64("reference_id").Optional().
			HelpText("ID of related record (order_id, transfer_id, etc.)").Build(),
		schema.String("reference_number").MaxLength(100).Optional().
			HelpText("Human-readable reference number").Build(),
		
		// Transfer (if applicable)
		schema.Int64("from_warehouse_id").Optional().
			HelpText("Source warehouse for transfers").Build(),
		schema.Int64("to_warehouse_id").Optional().
			HelpText("Destination warehouse for transfers").Build(),
		
		// Additional information
		schema.String("reason").MaxLength(255).Optional().
			HelpText("Reason for movement").Build(),
		schema.Text("notes").Optional().Build(),
		
		// User tracking
		schema.Int64("user_id").Optional().
			HelpText("User who performed the movement").Build(),
		schema.String("user_name").MaxLength(200).Optional().Build(),
		
		// Cost (for accounting)
		schema.Float64("unit_cost").Optional().
			HelpText("Cost per unit at time of movement").Build(),
		schema.Float64("total_cost").Optional().
			HelpText("Total cost of movement").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("movement_date").Required().
			HelpText("Date of movement (can be backdated)").Build(),
	}
}

func (StockMovement) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "stock_movements",
		VerboseName:      "Stock Movement",
		VerboseNamePlural: "Stock Movements",
		OrderBy:          []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_movement_stock", Fields: []string{"stock_id"}},
			{Name: "idx_movement_variant", Fields: []string{"product_variant_id"}},
			{Name: "idx_movement_warehouse", Fields: []string{"warehouse_id"}},
			{Name: "idx_movement_type", Fields: []string{"type"}},
			{Name: "idx_movement_reference", Fields: []string{"reference_type", "reference_id"}},
			{Name: "idx_movement_date", Fields: []string{"movement_date"}},
		},
	}
}

func (StockMovement) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("stock_id", "Stock").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("movements").Build(),
		schema.ForeignKey("product_variant_id", "ProductVariant").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("stock_movements").Build(),
		schema.ForeignKey("warehouse_id", "Warehouse").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("stock_movements").Build(),
		schema.ForeignKey("from_warehouse_id", "Warehouse").
			OnDelete(schema.CascadeSET_NULL).
			Optional().
			RelatedName("outbound_transfers").Build(),
		schema.ForeignKey("to_warehouse_id", "Warehouse").
			OnDelete(schema.CascadeSET_NULL).
			Optional().
			RelatedName("inbound_transfers").Build(),
	}
}

func (StockMovement) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeCreate: func(ctx context.Context, instance interface{}) error {
			// Capture quantity_before
			// Calculate quantity_after
			// Calculate total_cost
			return nil
		},
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Update stock quantity
			// Create alert if low stock
			return nil
		},
	}
}

// StockAlert represents low stock alerts
type StockAlert struct {
	schema.BaseSchema
}

func (StockAlert) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("stock_id").Required().Build(),
		schema.Int64("product_variant_id").Required().Build(),
		schema.Int64("warehouse_id").Required().Build(),
		
		// Alert details
		schema.String("alert_type").Required().MaxLength(50).
			HelpText("Type: low_stock, out_of_stock, overstock").Build(),
		schema.Int32("current_quantity").Required().Build(),
		schema.Int32("threshold").Required().
			HelpText("Threshold that triggered alert").Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("active").
			HelpText("Status: active, acknowledged, resolved, dismissed").Build(),
		
		// Resolution
		schema.Int64("resolved_by_user_id").Optional().Build(),
		schema.String("resolved_by_user_name").MaxLength(200).Optional().Build(),
		schema.Text("resolution_notes").Optional().Build(),
		
		// Notification
		schema.Bool("notification_sent").Default(false).Build(),
		schema.Time("notification_sent_at").Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("acknowledged_at").Optional().Build(),
		schema.Time("resolved_at").Optional().Build(),
	}
}

func (StockAlert) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "stock_alerts",
		VerboseName:      "Stock Alert",
		VerboseNamePlural: "Stock Alerts",
		OrderBy:          []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_alert_stock", Fields: []string{"stock_id"}},
			{Name: "idx_alert_variant", Fields: []string{"product_variant_id"}},
			{Name: "idx_alert_warehouse", Fields: []string{"warehouse_id"}},
			{Name: "idx_alert_status", Fields: []string{"status"}},
			{Name: "idx_alert_type", Fields: []string{"alert_type"}},
		},
	}
}

func (StockAlert) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("stock_id", "Stock").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("alerts").Build(),
		schema.ForeignKey("product_variant_id", "ProductVariant").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("stock_alerts").Build(),
		schema.ForeignKey("warehouse_id", "Warehouse").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("stock_alerts").Build(),
	}
}

func (StockAlert) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Send notification email
			return nil
		},
	}
}

// StockTransfer represents transfers between warehouses
type StockTransfer struct {
	schema.BaseSchema
}

func (StockTransfer) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		
		// Transfer identification
		schema.String("transfer_number").Required().MaxLength(50).Unique().Build(),
		
		// Warehouses
		schema.Int64("from_warehouse_id").Required().Build(),
		schema.Int64("to_warehouse_id").Required().Build(),
		
		// Product
		schema.Int64("product_variant_id").Required().Build(),
		schema.Int32("quantity").Required().Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("pending").
			HelpText("Status: pending, in_transit, completed, cancelled").Build(),
		
		// Tracking
		schema.String("tracking_number").MaxLength(255).Optional().Build(),
		schema.String("carrier").MaxLength(100).Optional().Build(),
		
		// Notes
		schema.Text("notes").Optional().Build(),
		schema.Text("reason").Optional().Build(),
		
		// User tracking
		schema.Int64("requested_by_user_id").Optional().Build(),
		schema.String("requested_by_user_name").MaxLength(200).Optional().Build(),
		schema.Int64("approved_by_user_id").Optional().Build(),
		schema.String("approved_by_user_name").MaxLength(200).Optional().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("shipped_at").Optional().Build(),
		schema.Time("completed_at").Optional().Build(),
		schema.Time("cancelled_at").Optional().Build(),
	}
}

func (StockTransfer) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "stock_transfers",
		VerboseName:      "Stock Transfer",
		VerboseNamePlural: "Stock Transfers",
		OrderBy:          []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_transfer_number", Fields: []string{"transfer_number"}, Unique: true},
			{Name: "idx_transfer_from", Fields: []string{"from_warehouse_id"}},
			{Name: "idx_transfer_to", Fields: []string{"to_warehouse_id"}},
			{Name: "idx_transfer_variant", Fields: []string{"product_variant_id"}},
			{Name: "idx_transfer_status", Fields: []string{"status"}},
		},
	}
}

func (StockTransfer) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("from_warehouse_id", "Warehouse").
			OnDelete(schema.CascadePROTECT).
			Required().
			RelatedName("outbound_stock_transfers").Build(),
		schema.ForeignKey("to_warehouse_id", "Warehouse").
			OnDelete(schema.CascadePROTECT).
			Required().
			RelatedName("inbound_stock_transfers").Build(),
		schema.ForeignKey("product_variant_id", "ProductVariant").
			OnDelete(schema.CascadePROTECT).
			Required().
			RelatedName("transfers").Build(),
	}
}

func (StockTransfer) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeCreate: func(ctx context.Context, instance interface{}) error {
			// Generate transfer number
			// Validate sufficient stock in source warehouse
			return nil
		},
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Create stock movements for both warehouses
			return nil
		},
		AfterUpdate: func(ctx context.Context, instance interface{}) error {
			// Handle status changes
			// Update stock movements
			return nil
		},
	}
}

// RegisterModels registers inventory models with the framework
func RegisterModels() {
	registry.RegisterModel(&Warehouse{})
	registry.RegisterModel(&Stock{})
	registry.RegisterModel(&StockMovement{})
	registry.RegisterModel(&StockAlert{})
	registry.RegisterModel(&StockTransfer{})
}
