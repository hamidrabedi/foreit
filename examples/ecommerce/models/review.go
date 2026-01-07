package models

import "github.com/forgego/forge/schema"

// Review represents a product review.
type Review struct {
	schema.BaseSchema
}

func (Review) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("customer_id").WithRequired(),
		schema.Int32("rating").WithRequired(),
		schema.String("title").WithMaxLength(200).WithOptional(),
		schema.Text("body").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Review) Meta() schema.Meta {
	return schema.Meta{
		TableName: "reviews",
	}
}

func (Review) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("reviews"),
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("reviews"),
	}
}

func (Review) Hooks() *schema.ModelHooks {
	return nil
}
