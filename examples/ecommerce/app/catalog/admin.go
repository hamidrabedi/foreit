package catalog

import (
	"context"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/db"
)

// RegisterAdmin registers catalog models with the admin interface
func RegisterAdmin(ctx context.Context, registry *admin.Registry, database *db.DB) {
	// Category admin
	categoryConfig := &admin.ModelConfig{
		Name:          "Category",
		PluralName:    "Categories",
		Icon:          "📁",
		ListDisplay:   []string{"id", "name", "parent_id", "sort_order", "is_active", "created_at"},
		ListFilter:    []string{"is_active", "parent_id", "level"},
		SearchFields:  []string{"name", "slug", "description"},
		OrderBy:       []string{"sort_order", "name"},
		PerPage:       20,
		Actions:       []string{"delete", "activate", "deactivate"},
		ExportFormats: []string{"csv", "json"},
	}
	registry.Register("Category", &Category{}, categoryConfig)

	// Brand admin
	brandConfig := &admin.ModelConfig{
		Name:          "Brand",
		PluralName:    "Brands",
		Icon:          "🏷️",
		ListDisplay:   []string{"id", "name", "slug", "website_url", "is_active", "created_at"},
		ListFilter:    []string{"is_active"},
		SearchFields:  []string{"name", "slug", "description"},
		OrderBy:       []string{"name"},
		PerPage:       20,
		Actions:       []string{"delete", "activate", "deactivate"},
		ExportFormats: []string{"csv", "json"},
	}
	registry.Register("Brand", &Brand{}, brandConfig)

	// Product admin
	productConfig := &admin.ModelConfig{
		Name:          "Product",
		PluralName:    "Products",
		Icon:          "📦",
		ListDisplay:   []string{"id", "name", "sku", "price", "category_id", "brand_id", "stock_quantity", "is_active", "created_at"},
		ListFilter:    []string{"is_active", "is_featured", "category_id", "brand_id"},
		SearchFields:  []string{"name", "sku", "description"},
		OrderBy:       []string{"-created_at"},
		PerPage:       20,
		Actions:       []string{"delete", "activate", "deactivate", "feature", "unfeature", "export"},
		ExportFormats: []string{"csv", "json", "xml"},
		BulkActions:   true,
	}
	registry.Register("Product", &Product{}, productConfig)

	// ProductVariant admin
	variantConfig := &admin.ModelConfig{
		Name:          "Product Variant",
		PluralName:    "Product Variants",
		Icon:          "🔖",
		ListDisplay:   []string{"id", "product_id", "name", "sku", "price", "stock_quantity", "is_active"},
		ListFilter:    []string{"is_active", "product_id"},
		SearchFields:  []string{"name", "sku"},
		OrderBy:       []string{"product_id", "sort_order"},
		PerPage:       20,
		Actions:       []string{"delete", "activate", "deactivate"},
		ExportFormats: []string{"csv", "json"},
	}
	registry.Register("ProductVariant", &ProductVariant{}, variantConfig)

	// ProductImage admin
	imageConfig := &admin.ModelConfig{
		Name:         "Product Image",
		PluralName:   "Product Images",
		Icon:         "🖼️",
		ListDisplay:  []string{"id", "product_id", "variant_id", "alt_text", "is_primary", "sort_order"},
		ListFilter:   []string{"is_primary", "product_id"},
		SearchFields: []string{"alt_text"},
		OrderBy:      []string{"product_id", "sort_order"},
		PerPage:      20,
		Actions:      []string{"delete"},
	}
	registry.Register("ProductImage", &ProductImage{}, imageConfig)

	// ProductAttribute admin
	attributeConfig := &admin.ModelConfig{
		Name:         "Product Attribute",
		PluralName:   "Product Attributes",
		Icon:         "⚙️",
		ListDisplay:  []string{"id", "name", "code", "type", "is_filterable", "is_visible", "sort_order"},
		ListFilter:   []string{"type", "is_filterable", "is_visible"},
		SearchFields: []string{"name", "code"},
		OrderBy:      []string{"sort_order", "name"},
		PerPage:      20,
		Actions:      []string{"delete"},
	}
	registry.Register("ProductAttribute", &ProductAttribute{}, attributeConfig)

	// ProductAttributeValue admin
	attrValueConfig := &admin.ModelConfig{
		Name:         "Product Attribute Value",
		PluralName:   "Product Attribute Values",
		Icon:         "📝",
		ListDisplay:  []string{"id", "product_id", "attribute_id", "value", "created_at"},
		ListFilter:   []string{"product_id", "attribute_id"},
		SearchFields: []string{"value"},
		OrderBy:      []string{"product_id", "attribute_id"},
		PerPage:      20,
		Actions:      []string{"delete"},
	}
	registry.Register("ProductAttributeValue", &ProductAttributeValue{}, attrValueConfig)
}
