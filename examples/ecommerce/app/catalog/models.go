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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200).WithUnique().
			WithHelpText("Category name"),
		schema.String("slug").WithRequired().WithMaxLength(200).WithUnique().
			WithHelpText("URL-friendly identifier"),
		schema.Text("description").WithOptional().
			WithHelpText("Category description"),
		schema.Int64("parent_id").WithOptional().
			WithHelpText("Parent category for hierarchy"),
		schema.String("image_url").WithMaxLength(500).WithOptional().
			WithHelpText("Category image"),
		schema.Int32("sort_order").WithDefault(0).
			WithHelpText("Display order"),
		schema.Bool("is_active").WithDefault(true).
			WithHelpText("Is category active"),
		schema.Int32("level").WithDefault(0).
			WithHelpText("Hierarchy level (0=root)"),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
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
		schema.ForeignKey("parent_id", "Category").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("children"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200).WithUnique(),
		schema.String("slug").WithRequired().WithMaxLength(200).WithUnique(),
		schema.Text("description").WithOptional(),
		schema.String("logo_url").WithMaxLength(500).WithOptional(),
		schema.String("website_url").WithMaxLength(500).WithOptional(),
		schema.Bool("is_active").WithDefault(true),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
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
	IsActive bool `json:"is_active" db:"is_active"`
}

func (Product) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(255).
			WithHelpText("Product name"),
		schema.String("slug").WithRequired().WithMaxLength(255).WithUnique().
			WithHelpText("URL-friendly identifier"),
		schema.String("sku").WithRequired().WithMaxLength(100).WithUnique().
			WithHelpText("Stock Keeping Unit"),
		schema.Text("description").WithRequired().
			WithHelpText("Full product description"),
		schema.Text("short_description").WithOptional().
			WithHelpText("Brief description for listings"),

		// Relationships
		schema.Int64("category_id").WithRequired().
			WithHelpText("Product category"),
		schema.Int64("brand_id").WithOptional().
			WithHelpText("Product brand"),

		// Pricing
		schema.Float64("price").WithRequired().
			WithHelpText("Base price"),
		schema.Float64("cost_price").WithOptional().
			WithHelpText("Cost to seller"),
		schema.Float64("compare_at_price").WithOptional().
			WithHelpText("Original price for discounts"),

		// Inventory (base level)
		schema.Int32("stock_quantity").WithDefault(0).
			WithHelpText("Total stock across all warehouses"),
		schema.Bool("track_inventory").WithDefault(true).
			WithHelpText("Whether to track inventory"),
		schema.Bool("allow_backorder").WithDefault(false).
			WithHelpText("Allow orders when out of stock"),

		// Physical attributes
		schema.Float64("weight").WithOptional().
			WithHelpText("Weight in kg"),
		schema.Float64("length").WithOptional().
			WithHelpText("Length in cm"),
		schema.Float64("width").WithOptional().
			WithHelpText("Width in cm"),
		schema.Float64("height").WithOptional().
			WithHelpText("Height in cm"),

		// Status
		schema.Bool("is_active").WithDefault(true),
		schema.Bool("is_featured").WithDefault(false),
		schema.Bool("is_digital").WithDefault(false).
			WithHelpText("Digital product (no shipping)"),

		// SEO
		schema.String("meta_title").WithMaxLength(255).WithOptional(),
		schema.Text("meta_description").WithOptional(),
		schema.String("meta_keywords").WithMaxLength(500).WithOptional(),

		// Stats
		schema.Int32("view_count").WithDefault(0),
		schema.Int32("order_count").WithDefault(0),
		schema.Float64("rating_average").WithDefault(0.0),
		schema.Int32("rating_count").WithDefault(0),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("published_at").WithOptional(),
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
		schema.ForeignKey("category_id", "Category").
			WithOnDelete(schema.CascadePROTECT).
			WithRequired().
			WithRelatedName("products"),
		schema.ForeignKey("brand_id", "Brand").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("products"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),

		// Variant identification
		schema.String("sku").WithRequired().WithMaxLength(100).WithUnique().
			WithHelpText("Unique SKU for this variant"),
		schema.String("name").WithRequired().WithMaxLength(255).
			WithHelpText("Variant name (e.g., 'Large Red')"),

		// Variant options
		schema.String("option1_name").WithMaxLength(100).WithOptional().
			WithHelpText("First option name (e.g., 'Size')"),
		schema.String("option1_value").WithMaxLength(100).WithOptional().
			WithHelpText("First option value (e.g., 'Large')"),
		schema.String("option2_name").WithMaxLength(100).WithOptional().
			WithHelpText("Second option name (e.g., 'Color')"),
		schema.String("option2_value").WithMaxLength(100).WithOptional().
			WithHelpText("Second option value (e.g., 'Red')"),
		schema.String("option3_name").WithMaxLength(100).WithOptional().
			WithHelpText("Third option name"),
		schema.String("option3_value").WithMaxLength(100).WithOptional().
			WithHelpText("Third option value"),

		// Pricing (can override product price)
		schema.Float64("price").WithOptional().
			WithHelpText("Override price (uses product price if null)"),
		schema.Float64("compare_at_price").WithOptional(),
		schema.Float64("cost_price").WithOptional(),

		// Inventory
		schema.Int32("stock_quantity").WithDefault(0),
		schema.Int32("reserved_quantity").WithDefault(0).
			WithHelpText("Quantity in pending orders"),
		schema.Bool("track_inventory").WithDefault(true),

		// Physical attributes (can override product)
		schema.Float64("weight").WithOptional(),
		schema.Float64("length").WithOptional(),
		schema.Float64("width").WithOptional(),
		schema.Float64("height").WithOptional(),

		// Status
		schema.Bool("is_active").WithDefault(true),
		schema.Bool("is_default").WithDefault(false).
			WithHelpText("Default variant for product"),

		// Display
		schema.String("image_url").WithMaxLength(500).WithOptional(),
		schema.Int32("sort_order").WithDefault(0),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
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
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("variants"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("variant_id").WithOptional().
			WithHelpText("Optional: Associate image with specific variant"),

		schema.String("image_url").WithRequired().WithMaxLength(500),
		schema.String("thumbnail_url").WithMaxLength(500).WithOptional(),
		schema.String("alt_text").WithMaxLength(255).WithOptional().
			WithHelpText("Alternative text for accessibility"),

		schema.Int32("sort_order").WithDefault(0),
		schema.Bool("is_primary").WithDefault(false).
			WithHelpText("Primary product image"),

		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
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
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("images"),
		schema.ForeignKey("product_image_id", "ProductImage").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("images"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200).WithUnique().
			WithHelpText("Attribute name (e.g., 'Color', 'Size')"),
		schema.String("code").WithRequired().WithMaxLength(100).WithUnique().
			WithHelpText("Code for programmatic access"),
		schema.String("type").WithRequired().WithMaxLength(50).
			WithHelpText("Type: text, number, select, multiselect"),
		schema.Bool("is_filterable").WithDefault(true).
			WithHelpText("Show in filters"),
		schema.Bool("is_visible").WithDefault(true).
			WithHelpText("Show on product page"),
		schema.Int32("sort_order").WithDefault(0),
		schema.Time("created_at").WithAutoNowAdd(),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("attribute_id").WithRequired(),
		schema.String("value").WithRequired().WithMaxLength(500),
		schema.Time("created_at").WithAutoNowAdd(),
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
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("attribute_values"),
		schema.ForeignKey("attribute_id", "ProductAttribute").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("values"),
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
