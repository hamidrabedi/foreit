package models

import "github.com/forgego/forge/schema"

// Product represents a catalog item.
type Product struct {
	schema.BaseSchema
}

func (Product) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.UUID("sku").WithRequired(),
		schema.String("name").WithRequired().WithMaxLength(255),
		schema.String("slug").WithRequired().WithMaxLength(255),
		schema.Text("description").WithOptional(),
		schema.Decimal("price").WithRequired().WithMaxDigits(12).WithDecimalPlaces(2),
		schema.String("currency").WithRequired().WithMaxLength(3).WithDefault("USD"),
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("active"),
		schema.Bool("is_active").WithDefault(true),
		schema.Int64("brand_id").WithOptional(),
		schema.Int64("supplier_id").WithOptional(),
		schema.JSON("attributes").WithOptional(),
		schema.JSON("images").WithOptional(),
		schema.JSON("seo_data").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Product) Meta() schema.Meta {
	return schema.Meta{
		TableName: "products",
		Indexes: []schema.Index{
			{Name: "idx_product_sku", Fields: []string{"sku"}, Unique: true},
			{Name: "idx_product_slug", Fields: []string{"slug"}, Unique: true},
			{Name: "idx_product_status", Fields: []string{"status"}},
			{Name: "idx_product_is_active", Fields: []string{"is_active"}},
		},
	}
}

func (Product) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("brand_id", "Brand").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("products"),
		schema.ForeignKey("supplier_id", "Supplier").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("products"),
	}
}

func (Product) Hooks() *schema.ModelHooks {
	return nil
}
