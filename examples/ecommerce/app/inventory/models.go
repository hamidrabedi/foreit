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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),

		// Identification
		schema.String("name").WithRequired().WithMaxLength(200).WithUnique(),
		schema.String("code").WithRequired().WithMaxLength(50).WithUnique().
			WithHelpText("Warehouse code for reference"),

		// Contact
		schema.String("contact_name").WithMaxLength(200).WithOptional(),
		schema.String("contact_email").WithMaxLength(255).WithOptional(),
		schema.String("contact_phone").WithMaxLength(20).WithOptional(),

		// Address
		schema.String("address_line1").WithRequired().WithMaxLength(255),
		schema.String("address_line2").WithMaxLength(255).WithOptional(),
		schema.String("city").WithRequired().WithMaxLength(100),
		schema.String("state").WithMaxLength(100).WithOptional(),
		schema.String("postal_code").WithRequired().WithMaxLength(20),
		schema.String("country_code").WithRequired().WithMaxLength(2),
		schema.String("country_name").WithRequired().WithMaxLength(100),

		// Geolocation
		schema.Float64("latitude").WithOptional(),
		schema.Float64("longitude").WithOptional(),

		// Settings
		schema.Bool("is_active").WithDefault(true),
		schema.Bool("is_primary").WithDefault(false).
			WithHelpText("Primary warehouse for fulfillment"),
		schema.Int32("priority").WithDefault(0).
			WithHelpText("Fulfillment priority (higher = preferred)"),

		// Capacity
		schema.Int32("total_capacity").WithOptional().
			WithHelpText("Total storage capacity (units)"),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Warehouse) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "warehouses",
		VerboseName:       "Warehouse",
		VerboseNamePlural: "Warehouses",
		OrderBy:           []string{"-is_primary", "-priority", "name"},
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_variant_id").WithRequired(),
		schema.Int64("warehouse_id").WithRequired(),

		// Quantities
		schema.Int32("quantity").WithRequired().WithDefault(0).
			WithHelpText("Available quantity"),
		schema.Int32("reserved_quantity").WithDefault(0).
			WithHelpText("Quantity reserved in pending orders"),
		schema.Int32("available_quantity").WithDefault(0).
			WithHelpText("quantity - reserved_quantity"),

		// Reordering
		schema.Int32("reorder_point").WithDefault(10).
			WithHelpText("Minimum quantity before reorder"),
		schema.Int32("reorder_quantity").WithDefault(50).
			WithHelpText("Quantity to reorder"),

		// Location
		schema.String("bin_location").WithMaxLength(100).WithOptional().
			WithHelpText("Physical location in warehouse (e.g., A-12-3)"),

		// Status
		schema.Bool("is_active").WithDefault(true),
		schema.Bool("allow_backorder").WithDefault(false),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("last_counted_at").WithOptional().
			WithHelpText("Last physical count"),
	}
}

func (Stock) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "stock",
		VerboseName:       "Stock",
		VerboseNamePlural: "Stock",
		OrderBy:           []string{"warehouse_id", "product_variant_id"},
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
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("stock_records"),
		schema.ForeignKey("warehouse_id", "Warehouse").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("stock_records"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("stock_id").WithRequired(),
		schema.Int64("product_variant_id").WithRequired(),
		schema.Int64("warehouse_id").WithRequired(),

		// Movement details
		schema.String("type").WithRequired().WithMaxLength(50).
			WithHelpText("Type: purchase, sale, return, adjustment, transfer, damage, loss, count"),
		schema.Int32("quantity").WithRequired().
			WithHelpText("Quantity change (positive or negative)"),
		schema.Int32("quantity_before").WithRequired().
			WithHelpText("Quantity before movement"),
		schema.Int32("quantity_after").WithRequired().
			WithHelpText("Quantity after movement"),

		// Reference
		schema.String("reference_type").WithMaxLength(50).WithOptional().
			WithHelpText("Reference type: order, return, transfer, etc."),
		schema.Int64("reference_id").WithOptional().
			WithHelpText("ID of related record (order_id, transfer_id, etc.)"),
		schema.String("reference_number").WithMaxLength(100).WithOptional().
			WithHelpText("Human-readable reference number"),

		// Transfer (if applicable)
		schema.Int64("from_warehouse_id").WithOptional().
			WithHelpText("Source warehouse for transfers"),
		schema.Int64("to_warehouse_id").WithOptional().
			WithHelpText("Destination warehouse for transfers"),

		// Additional information
		schema.String("reason").WithMaxLength(255).WithOptional().
			WithHelpText("Reason for movement"),
		schema.Text("notes").WithOptional(),

		// User tracking
		schema.Int64("user_id").WithOptional().
			WithHelpText("User who performed the movement"),
		schema.String("user_name").WithMaxLength(200).WithOptional(),

		// Cost (for accounting)
		schema.Float64("unit_cost").WithOptional().
			WithHelpText("Cost per unit at time of movement"),
		schema.Float64("total_cost").WithOptional().
			WithHelpText("Total cost of movement"),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("movement_date").WithRequired().
			WithHelpText("Date of movement (can be backdated)"),
	}
}

func (StockMovement) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "stock_movements",
		VerboseName:       "Stock Movement",
		VerboseNamePlural: "Stock Movements",
		OrderBy:           []string{"-created_at"},
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
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("movements"),
		schema.ForeignKey("product_variant_id", "ProductVariant").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("stock_movements"),
		schema.ForeignKey("warehouse_id", "Warehouse").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("stock_movements"),
		schema.ForeignKey("from_warehouse_id", "Warehouse").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("outbound_transfers"),
		schema.ForeignKey("to_warehouse_id", "Warehouse").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("inbound_transfers"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("stock_id").WithRequired(),
		schema.Int64("product_variant_id").WithRequired(),
		schema.Int64("warehouse_id").WithRequired(),

		// Alert details
		schema.String("alert_type").WithRequired().WithMaxLength(50).
			WithHelpText("Type: low_stock, out_of_stock, overstock"),
		schema.Int32("current_quantity").WithRequired(),
		schema.Int32("threshold").WithRequired().
			WithHelpText("Threshold that triggered alert"),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("active").
			WithHelpText("Status: active, acknowledged, resolved, dismissed"),

		// Resolution
		schema.Int64("resolved_by_user_id").WithOptional(),
		schema.String("resolved_by_user_name").WithMaxLength(200).WithOptional(),
		schema.Text("resolution_notes").WithOptional(),

		// Notification
		schema.Bool("notification_sent").WithDefault(false),
		schema.Time("notification_sent_at").WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("acknowledged_at").WithOptional(),
		schema.Time("resolved_at").WithOptional(),
	}
}

func (StockAlert) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "stock_alerts",
		VerboseName:       "Stock Alert",
		VerboseNamePlural: "Stock Alerts",
		OrderBy:           []string{"-created_at"},
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
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("alerts"),
		schema.ForeignKey("product_variant_id", "ProductVariant").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("stock_alerts"),
		schema.ForeignKey("warehouse_id", "Warehouse").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("stock_alerts"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),

		// Transfer identification
		schema.String("transfer_number").WithRequired().WithMaxLength(50).WithUnique(),

		// Warehouses
		schema.Int64("from_warehouse_id").WithRequired(),
		schema.Int64("to_warehouse_id").WithRequired(),

		// Product
		schema.Int64("product_variant_id").WithRequired(),
		schema.Int32("quantity").WithRequired(),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Status: pending, in_transit, completed, cancelled"),

		// Tracking
		schema.String("tracking_number").WithMaxLength(255).WithOptional(),
		schema.String("carrier").WithMaxLength(100).WithOptional(),

		// Notes
		schema.Text("notes").WithOptional(),
		schema.Text("reason").WithOptional(),

		// User tracking
		schema.Int64("requested_by_user_id").WithOptional(),
		schema.String("requested_by_user_name").WithMaxLength(200).WithOptional(),
		schema.Int64("approved_by_user_id").WithOptional(),
		schema.String("approved_by_user_name").WithMaxLength(200).WithOptional(),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("shipped_at").WithOptional(),
		schema.Time("completed_at").WithOptional(),
		schema.Time("cancelled_at").WithOptional(),
	}
}

func (StockTransfer) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "stock_transfers",
		VerboseName:       "Stock Transfer",
		VerboseNamePlural: "Stock Transfers",
		OrderBy:           []string{"-created_at"},
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
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("outbound_stock_transfers"),
		schema.ForeignKey("to_warehouse_id", "Warehouse").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("inbound_stock_transfers"),
		schema.ForeignKey("product_variant_id", "ProductVariant").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("transfers"),
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
