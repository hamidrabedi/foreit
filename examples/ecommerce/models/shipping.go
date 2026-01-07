package models

import "github.com/forgego/forge/schema"

// Shipping represents shipment tracking for an order.
type Shipping struct {
	schema.BaseSchema
}

func (Shipping) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("order_id").WithRequired(),
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending"),
		schema.String("tracking_number").WithMaxLength(255).WithOptional(),
		schema.String("carrier").WithMaxLength(100).WithOptional(),
		schema.DateTime("shipped_at").WithOptional(),
		schema.DateTime("delivered_at").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Shipping) Meta() schema.Meta {
	return schema.Meta{
		TableName: "shipping",
	}
}

func (Shipping) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("shipping"),
	}
}

func (Shipping) Hooks() *schema.ModelHooks {
	return nil
}
