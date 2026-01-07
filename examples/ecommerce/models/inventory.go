package models

import "github.com/forgego/forge/schema"

// Inventory represents stock levels for products.
type Inventory struct {
	schema.BaseSchema
}

func (Inventory) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.Int32("quantity").WithRequired().WithDefault(0),
		schema.Int32("reserved").WithDefault(0),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Inventory) Meta() schema.Meta {
	return schema.Meta{
		TableName: "inventory",
	}
}

func (Inventory) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("inventory_records"),
	}
}

func (Inventory) Hooks() *schema.ModelHooks {
	return nil
}
