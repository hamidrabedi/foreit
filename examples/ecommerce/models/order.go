package models

import "github.com/forgego/forge/schema"

// Order represents a customer order.
type Order struct {
	schema.BaseSchema
}

func (Order) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.UUID("order_number").WithRequired(),
		schema.Int64("customer_id").WithRequired(),
		schema.Decimal("subtotal").WithRequired().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.Decimal("total_amount").WithRequired().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.String("currency").WithRequired().WithMaxLength(3).WithDefault("USD"),
		schema.Int64("billing_address_id").WithRequired(),
		schema.Int64("shipping_address_id").WithRequired(),
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending"),
		schema.String("payment_status").WithRequired().WithMaxLength(20).WithDefault("pending"),
		schema.String("shipping_status").WithRequired().WithMaxLength(20).WithDefault("pending"),
		schema.JSON("metadata").WithOptional(),
		schema.DateTime("placed_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Order) Meta() schema.Meta {
	return schema.Meta{
		TableName: "orders",
		Indexes: []schema.Index{
			{Name: "idx_order_order_number", Fields: []string{"order_number"}, Unique: true},
			{Name: "idx_order_customer_id", Fields: []string{"customer_id"}},
			{Name: "idx_order_status", Fields: []string{"status"}},
		},
	}
}

func (Order) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("orders"),
		schema.ForeignKey("billing_address_id", "Address").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("billing_orders"),
		schema.ForeignKey("shipping_address_id", "Address").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("shipping_orders"),
	}
}

func (Order) Hooks() *schema.ModelHooks {
	return nil
}
