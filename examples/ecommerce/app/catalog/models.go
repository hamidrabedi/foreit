package catalog

import (
	"context"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Category represents a product category with hierarchical support
type Category struct {
	schema.BaseSchema
}

func (Category) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(200).Unique().
			HelpText("Category name").Build(),
		schema.String("slug").Required().MaxLength(200).Unique().
			HelpText("URL-friendly identifier").Build(),
		schema.Text("description").Null().
			HelpText("Category description").Build(),
		schema.Int64("parent_id").Null().
			HelpText("Parent category for hierarchy").Build(),
		schema.String("image_url").MaxLength(500).Null().
			HelpText("Category image").Build(),
		schema.Int32("sort_order").Default(0).
			HelpText("Display order").Build(),
		schema.Bool("is_active").Default(true).
			HelpText("Is category active").Build(),
		schema.Int32("level").Default(0).
			HelpText("Hierarchy level (0=root)").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Category) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "categories",
		VerboseName:       "Category",
		VerboseNamePlural: "Categories",
		OrderBy:           []string{"sort_order", "name"},
		Indexes: []schema.Index{
			{Name: "idx_category_slug", Fields: []string{"slug"}, Unique: true},
			{Name: "idx_category_parent", Fields: []string{"parent_id"}},
			{Name: "idx_category_active", Fields: []string{"is_active"}},
		},
	}
}

func (Category) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("parent_id", "Category", "parent").
			OnDelete(schema.SetNull).
			Null().
			RelatedName("children").
			HelpText("Parent category"),
	}
}

func (Category) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Auto-generate slug from name if not provided
			// Calculate hierarchy level
			return nil
		},
	}
}

// Brand represents a product brand
type Brand struct {
	schema.BaseSchema
}

func (Brand) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(200).Unique().Build(),
		schema.String("slug").Required().MaxLength(200).Unique().Build(),
		schema.Text("description").Null().Build(),
		schema.String("logo_url").MaxLength(500).Null().Build(),
		schema.String("website_url").MaxLength(500).Null().Build(),
		schema.Bool("is_active").Default(true).Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Brand) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "brands",
		VerboseName:       "Brand",
		VerboseNamePlural: "Brands",
		OrderBy:           []string{"name"},
		Indexes: []schema.Index{
			{Name: "idx_brand_slug", Fields: []string{"slug"}, Unique: true},
		},
	}
}

func (Brand) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Brand) Hooks() *schema.ModelHooks {
	return nil
}

// Product represents a product in the catalog
type Product struct {
	schema.BaseSchema
}

func (Product) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(255).
			HelpText("Product name").Build(),
		schema.String("slug").Required().MaxLength(255).Unique().
			HelpText("URL-friendly identifier").Build(),
		schema.String("sku").Required().MaxLength(100).Unique().
			HelpText("Stock Keeping Unit").Build(),
		schema.Text("description").Required().
			HelpText("Full product description").Build(),
		schema.Text("short_description").Null().
			HelpText("Brief description for listings").Build(),

		// Relationships
		schema.Int64("category_id").Required().
			HelpText("Product category").Build(),
		schema.Int64("brand_id").Null().
			HelpText("Product brand").Build(),

		// Pricing
		schema.Float64("price").Required().
			HelpText("Base price").Build(),
		schema.Float64("cost_price").Null().
			HelpText("Cost to seller").Build(),
		schema.Float64("compare_at_price").Null().
			HelpText("Original price for discounts").Build(),

		// Inventory (base level)
		schema.Int32("stock_quantity").Default(0).
			HelpText("Total stock across all warehouses").Build(),
		schema.Bool("track_inventory").Default(true).
			HelpText("Whether to track inventory").Build(),
		schema.Bool("allow_backorder").Default(false).
			HelpText("Allow orders when out of stock").Build(),

		// Physical attributes
		schema.Float64("weight").Null().
			HelpText("Weight in kg").Build(),
		schema.Float64("length").Null().
			HelpText("Length in cm").Build(),
		schema.Float64("width").Null().
			HelpText("Width in cm").Build(),
		schema.Float64("height").Null().
			HelpText("Height in cm").Build(),

		// Status
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("is_featured").Default(false).Build(),
		schema.Bool("is_digital").Default(false).
			HelpText("Digital product (no shipping)").Build(),

		// SEO
		schema.String("meta_title").MaxLength(255).Null().Build(),
		schema.Text("meta_description").Null().Build(),
		schema.String("meta_keywords").MaxLength(500).Null().Build(),

		// Stats
		schema.Int32("view_count").Default(0).Build(),
		schema.Int32("order_count").Default(0).Build(),
		schema.Float64("rating_average").Default(0.0).Build(),
		schema.Int32("rating_count").Default(0).Build(),

		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("published_at").Null().Build(),
	}
}

func (Product) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "products",
		VerboseName:       "Product",
		VerboseNamePlural: "Products",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_product_slug", Fields: []string{"slug"}, Unique: true},
			{Name: "idx_product_sku", Fields: []string{"sku"}, Unique: true},
			{Name: "idx_product_category", Fields: []string{"category_id"}},
			{Name: "idx_product_brand", Fields: []string{"brand_id"}},
			{Name: "idx_product_active", Fields: []string{"is_active"}},
			{Name: "idx_product_featured", Fields: []string{"is_featured"}},
			{Name: "idx_product_price", Fields: []string{"price"}},
		},
	}
}

func (Product) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("category_id", "Category", "category").
			OnDelete(schema.Protect).
			Required().
			RelatedName("products"),
		schema.ForeignKey("brand_id", "Brand", "brand").
			OnDelete(schema.SetNull).
			Null().
			RelatedName("products"),
	}
}

func (Product) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Auto-generate slug if not provided
			// Validate price > 0
			// Set published_at if becoming active
			return nil
		},
		AfterSave: func(ctx context.Context, instance interface{}) error {
			// Update search index
			// Clear cache
			return nil
		},
	}
}

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
	schema.BaseSchema
}

func (ProductVariant) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("product_id").Required().Build(),

		// Variant identification
		schema.String("sku").Required().MaxLength(100).Unique().
			HelpText("Unique SKU for this variant").Build(),
		schema.String("name").Required().MaxLength(255).
			HelpText("Variant name (e.g., 'Large Red')").Build(),

		// Variant options
		schema.String("option1_name").MaxLength(100).Null().
			HelpText("First option name (e.g., 'Size')").Build(),
		schema.String("option1_value").MaxLength(100).Null().
			HelpText("First option value (e.g., 'Large')").Build(),
		schema.String("option2_name").MaxLength(100).Null().
			HelpText("Second option name (e.g., 'Color')").Build(),
		schema.String("option2_value").MaxLength(100).Null().
			HelpText("Second option value (e.g., 'Red')").Build(),
		schema.String("option3_name").MaxLength(100).Null().
			HelpText("Third option name").Build(),
		schema.String("option3_value").MaxLength(100).Null().
			HelpText("Third option value").Build(),

		// Pricing (can override product price)
		schema.Float64("price").Null().
			HelpText("Override price (uses product price if null)").Build(),
		schema.Float64("compare_at_price").Null().Build(),
		schema.Float64("cost_price").Null().Build(),

		// Inventory
		schema.Int32("stock_quantity").Default(0).Build(),
		schema.Int32("reserved_quantity").Default(0).
			HelpText("Quantity in pending orders").Build(),
		schema.Bool("track_inventory").Default(true).Build(),

		// Physical attributes (can override product)
		schema.Float64("weight").Null().Build(),
		schema.Float64("length").Null().Build(),
		schema.Float64("width").Null().Build(),
		schema.Float64("height").Null().Build(),

		// Status
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("is_default").Default(false).
			HelpText("Default variant for product").Build(),

		// Display
		schema.String("image_url").MaxLength(500).Null().Build(),
		schema.Int32("sort_order").Default(0).Build(),

		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (ProductVariant) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "product_variants",
		VerboseName:       "Product Variant",
		VerboseNamePlural: "Product Variants",
		OrderBy:           []string{"product_id", "sort_order"},
		Indexes: []schema.Index{
			{Name: "idx_variant_sku", Fields: []string{"sku"}, Unique: true},
			{Name: "idx_variant_product", Fields: []string{"product_id"}},
			{Name: "idx_variant_active", Fields: []string{"is_active"}},
		},
		UniqueTogether: [][]string{
			{"product_id", "option1_value", "option2_value", "option3_value"},
		},
	}
}

func (ProductVariant) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product", "product").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("variants"),
	}
}

func (ProductVariant) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Ensure only one default variant per product
			// Validate SKU uniqueness
			return nil
		},
		AfterSave: func(ctx context.Context, instance interface{}) error {
			// Update product stock_quantity
			return nil
		},
	}
}

// ProductImage represents product images
type ProductImage struct {
	schema.BaseSchema
}

func (ProductImage) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("variant_id").Null().
			HelpText("Optional: Associate image with specific variant").Build(),

		schema.String("image_url").Required().MaxLength(500).Build(),
		schema.String("thumbnail_url").MaxLength(500).Null().Build(),
		schema.String("alt_text").MaxLength(255).Null().
			HelpText("Alternative text for accessibility").Build(),

		schema.Int32("sort_order").Default(0).Build(),
		schema.Bool("is_primary").Default(false).
			HelpText("Primary product image").Build(),

		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (ProductImage) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "product_images",
		VerboseName:       "Product Image",
		VerboseNamePlural: "Product Images",
		OrderBy:           []string{"product_id", "sort_order"},
		Indexes: []schema.Index{
			{Name: "idx_image_product", Fields: []string{"product_id"}},
			{Name: "idx_image_variant", Fields: []string{"variant_id"}},
		},
	}
}

func (ProductImage) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product", "product").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("images"),
		schema.ForeignKey("variant_id", "ProductVariant", "variant").
			OnDelete(schema.Cascade).
			Null().
			RelatedName("images"),
	}
}

func (ProductImage) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Ensure only one primary image per product
			return nil
		},
	}
}

// ProductAttribute represents dynamic product attributes
type ProductAttribute struct {
	schema.BaseSchema
}

func (ProductAttribute) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(200).Unique().
			HelpText("Attribute name (e.g., 'Color', 'Size')").Build(),
		schema.String("code").Required().MaxLength(100).Unique().
			HelpText("Code for programmatic access").Build(),
		schema.String("type").Required().MaxLength(50).
			HelpText("Type: text, number, select, multiselect").Build(),
		schema.Bool("is_filterable").Default(true).
			HelpText("Show in filters").Build(),
		schema.Bool("is_visible").Default(true).
			HelpText("Show on product page").Build(),
		schema.Int32("sort_order").Default(0).Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
	}
}

func (ProductAttribute) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "product_attributes",
		VerboseName:       "Product Attribute",
		VerboseNamePlural: "Product Attributes",
		OrderBy:           []string{"sort_order", "name"},
	}
}

func (ProductAttribute) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (ProductAttribute) Hooks() *schema.ModelHooks {
	return nil
}

// ProductAttributeValue represents attribute values for products
type ProductAttributeValue struct {
	schema.BaseSchema
}

func (ProductAttributeValue) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("attribute_id").Required().Build(),
		schema.String("value").Required().MaxLength(500).Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
	}
}

func (ProductAttributeValue) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "product_attribute_values",
		VerboseName:       "Product Attribute Value",
		VerboseNamePlural: "Product Attribute Values",
		Indexes: []schema.Index{
			{Name: "idx_attr_value_product", Fields: []string{"product_id"}},
			{Name: "idx_attr_value_attribute", Fields: []string{"attribute_id"}},
		},
		UniqueTogether: [][]string{
			{"product_id", "attribute_id"},
		},
	}
}

func (ProductAttributeValue) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product", "product").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("attribute_values"),
		schema.ForeignKey("attribute_id", "ProductAttribute", "attribute").
			OnDelete(schema.Cascade).
			Required().
			RelatedName("values"),
	}
}

func (ProductAttributeValue) Hooks() *schema.ModelHooks {
	return nil
}

// RegisterModels registers catalog models with the framework
func RegisterModels() {
	registry.RegisterModel(&Category{})
	registry.RegisterModel(&Brand{})
	registry.RegisterModel(&Product{})
	registry.RegisterModel(&ProductVariant{})
	registry.RegisterModel(&ProductImage{})
	registry.RegisterModel(&ProductAttribute{})
	registry.RegisterModel(&ProductAttributeValue{})
}
