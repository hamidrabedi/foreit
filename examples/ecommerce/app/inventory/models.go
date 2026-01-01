package inventory

import (
	"context"
	
	"github.com/forgego/forge/schema"
	"github.com/forgego/forge/registry"
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
		schema.String("contact_name").MaxLength(200).Null().Build(),
		schema.String("contact_email").MaxLength(255).Null().Build(),
		schema.String("contact_phone").MaxLength(20).Null().Build(),
		
		// Address
		schema.String("address_line1").Required().MaxLength(255).Build(),
		schema.String("address_line2").MaxLength(255).Null().Build(),
		schema.String("city").Required().MaxLength(100).Build(),
		schema.String("state").MaxLength(100).Null().Build(),
		schema.String("postal_code").Required().MaxLength(20).Build(),
		schema.String("country_code").Required().MaxLength(2).Build(),
		schema.String("country_name").Required().MaxLength(100).Build(),
		
		// Geolocation
		schema.Float64("latitude").Null().Build(),
		schema.Float64("longitude").Null().Build(),
		
		// Settings
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("is_primary").Default(false).
			HelpText("Primary warehouse for fulfillment").Build(),
		schema.Int32("priority").Default(0).
			HelpText("Fulfillment priority (higher = preferred)").Build(),
		
		// Capacity
		schema.Int32("total_capacity").Null().
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
		schema.String("bin_location").MaxLength(100).Null().
			HelpText("Physical location in warehouse (e.g., A-12-3)").Build(),
		
		// Status
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("allow_backorder").Default(false).Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("last_counted_at").Null().
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
		schema.ForeignKey("product_variant_id", "ProductVariant", "product_variant").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("stock_records"),
		schema.ForeignKey("warehouse_id", "Warehouse", "warehouse").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("stock_records"),
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
		schema.String("reference_type").MaxLength(50).Null().
			HelpText("Reference type: order, return, transfer, etc.").Build(),
		schema.Int64("reference_id").Null().
			HelpText("ID of related record (order_id, transfer_id, etc.)").Build(),
		schema.String("reference_number").MaxLength(100).Null().
			HelpText("Human-readable reference number").Build(),
		
		// Transfer (if applicable)
		schema.Int64("from_warehouse_id").Null().
			HelpText("Source warehouse for transfers").Build(),
		schema.Int64("to_warehouse_id").Null().
			HelpText("Destination warehouse for transfers").Build(),
		
		// Additional information
		schema.String("reason").MaxLength(255).Null().
			HelpText("Reason for movement").Build(),
		schema.Text("notes").Null().Build(),
		
		// User tracking
		schema.Int64("user_id").Null().
			HelpText("User who performed the movement").Build(),
		schema.String("user_name").MaxLength(200).Null().Build(),
		
		// Cost (for accounting)
		schema.Float64("unit_cost").Null().
			HelpText("Cost per unit at time of movement").Build(),
		schema.Float64("total_cost").Null().
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
		schema.ForeignKey("stock_id", "Stock", "stock").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("movements"),
		schema.ForeignKey("product_variant_id", "ProductVariant", "product_variant").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("stock_movements"),
		schema.ForeignKey("warehouse_id", "Warehouse", "warehouse").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("stock_movements"),
		schema.ForeignKey("from_warehouse_id", "Warehouse", "from_warehouse").
			OnDelete(schema.SetNull).
			Null().
			RelatedName("outbound_transfers"),
		schema.ForeignKey("to_warehouse_id", "Warehouse", "to_warehouse").
			OnDelete(schema.SetNull).
			Null().
			RelatedName("inbound_transfers"),
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
		schema.Int64("resolved_by_user_id").Null().Build(),
		schema.String("resolved_by_user_name").MaxLength(200).Null().Build(),
		schema.Text("resolution_notes").Null().Build(),
		
		// Notification
		schema.Bool("notification_sent").Default(false).Build(),
		schema.Time("notification_sent_at").Null().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("acknowledged_at").Null().Build(),
		schema.Time("resolved_at").Null().Build(),
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
		schema.ForeignKey("stock_id", "Stock", "stock").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("alerts"),
		schema.ForeignKey("product_variant_id", "ProductVariant", "product_variant").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("stock_alerts"),
		schema.ForeignKey("warehouse_id", "Warehouse", "warehouse").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("stock_alerts"),
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
		schema.String("tracking_number").MaxLength(255).Null().Build(),
		schema.String("carrier").MaxLength(100).Null().Build(),
		
		// Notes
		schema.Text("notes").Null().Build(),
		schema.Text("reason").Null().Build(),
		
		// User tracking
		schema.Int64("requested_by_user_id").Null().Build(),
		schema.String("requested_by_user_name").MaxLength(200).Null().Build(),
		schema.Int64("approved_by_user_id").Null().Build(),
		schema.String("approved_by_user_name").MaxLength(200).Null().Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("shipped_at").Null().Build(),
		schema.Time("completed_at").Null().Build(),
		schema.Time("cancelled_at").Null().Build(),
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
		schema.ForeignKey("from_warehouse_id", "Warehouse", "from_warehouse").
			OnDelete(schema.Protect).
			Required().
			RelatedName("outbound_stock_transfers"),
		schema.ForeignKey("to_warehouse_id", "Warehouse", "to_warehouse").
			OnDelete(schema.Protect).
			Required().
			RelatedName("inbound_stock_transfers"),
		schema.ForeignKey("product_variant_id", "ProductVariant", "product_variant").
			OnDelete(schema.Protect).
			Required().
			RelatedName("transfers"),
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
