package models

import (
	"github.com/forgego/forge/pkg/schema"
)

// OrderItem represents an item in an order
type OrderItem struct {
	schema.BaseSchema
}

// Fields returns all field definitions for OrderItem
func (OrderItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("order_id").Required().VerboseName("Order ID").Build(),
		schema.Int64("product_id").Required().VerboseName("Product ID").Build(),
		schema.Int64("variant_id").Optional().VerboseName("Variant ID").Build(),
		schema.String("product_name").Required().MaxLength(500).VerboseName("Product Name").Build(),
		schema.String("product_sku").Required().MaxLength(100).VerboseName("Product SKU").Build(),
		schema.Int32("quantity").Required().VerboseName("Quantity").Build(),
		schema.Decimal("unit_price").MaxDigits(12).DecimalPlaces(2).Required().VerboseName("Unit Price").Build(),
		schema.Decimal("total_price").MaxDigits(12).DecimalPlaces(2).Required().VerboseName("Total Price").Build(),
		schema.Decimal("discount_amount").MaxDigits(12).DecimalPlaces(2).Default(0.0).VerboseName("Discount Amount").Build(),
		schema.Decimal("tax_amount").MaxDigits(12).DecimalPlaces(2).Default(0.0).VerboseName("Tax Amount").Build(),
		schema.JSON("product_data").Optional().VerboseName("Product Snapshot Data").Build(),
		schema.JSON("variant_data").Optional().VerboseName("Variant Snapshot Data").Build(),
		schema.DateTime("created_at").AutoNowAdd().Build(),
		schema.DateTime("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (OrderItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "order_items",
		VerboseName:       "Order Item",
		VerboseNamePlural: "Order Items",
		OrderBy:           []string{"id"},
		Indexes: []schema.Index{
			{Name: "idx_order_item_order_id", Fields: []string{"order_id"}, Unique: false},
			{Name: "idx_order_item_product_id", Fields: []string{"product_id"}, Unique: false},
			{Name: "idx_order_item_variant_id", Fields: []string{"variant_id"}, Unique: false},
		},
	}
}

// Relations returns all relationship definitions
func (OrderItem) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").Required().OnDelete(schema.CascadeCASCADE).RelatedName("items").Build(),
		schema.ForeignKey("product_id", "Product").Required().OnDelete(schema.CascadePROTECT).RelatedName("order_items").Build(),
		schema.ForeignKey("variant_id", "ProductVariant").Optional().OnDelete(schema.CascadePROTECT).RelatedName("order_items").Build(),
	}
}

// Hooks returns model lifecycle hooks
func (OrderItem) Hooks() *schema.ModelHooks {
	return nil
}

