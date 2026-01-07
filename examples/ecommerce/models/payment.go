package models

import "github.com/forgego/forge/schema"

// Payment represents a payment transaction.
type Payment struct {
	schema.BaseSchema
}

func (Payment) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("order_id").WithRequired(),
		schema.Decimal("amount").WithRequired().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.String("currency").WithRequired().WithMaxLength(3).WithDefault("USD"),
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending"),
		schema.String("transaction_id").WithMaxLength(255).WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Payment) Meta() schema.Meta {
	return schema.Meta{
		TableName: "payments",
	}
}

func (Payment) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("payments"),
	}
}

func (Payment) Hooks() *schema.ModelHooks {
	return nil
}
