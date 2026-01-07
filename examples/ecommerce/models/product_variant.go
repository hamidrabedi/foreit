package models

import "github.com/forgego/forge/schema"

// ProductVariant represents a product variant.
type ProductVariant struct {
	schema.BaseSchema
}

func (ProductVariant) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.String("sku").WithMaxLength(100).WithOptional(),
		schema.String("name").WithMaxLength(200).WithOptional(),
		schema.Decimal("price").WithOptional().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.JSON("attributes").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (ProductVariant) Meta() schema.Meta {
	return schema.Meta{
		TableName: "product_variants",
	}
}

func (ProductVariant) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("variants"),
	}
}

func (ProductVariant) Hooks() *schema.ModelHooks {
	return nil
}
