package models

import "github.com/forgego/forge/schema"

// OrderItem represents a line item in an order.
type OrderItem struct {
	schema.BaseSchema
}

func (OrderItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("order_id").WithRequired(),
		schema.Int64("product_id").WithRequired(),
		schema.String("product_name").WithRequired().WithMaxLength(255),
		schema.String("product_sku").WithRequired().WithMaxLength(100),
		schema.Int32("quantity").WithRequired().WithDefault(1),
		schema.Decimal("unit_price").WithRequired().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.Decimal("total_price").WithRequired().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (OrderItem) Meta() schema.Meta {
	return schema.Meta{
		TableName: "order_items",
	}
}

func (OrderItem) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("items"),
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("order_items"),
	}
}

func (OrderItem) Hooks() *schema.ModelHooks {
	return nil
}
