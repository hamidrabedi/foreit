package catalog

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers catalog models with the admin interface
func RegisterAdmin(ctx context.Context) {
	// Category admin
	admin.Register(&admin.Config[Category]{
		Icon: "FolderTree",
		ListDisplay: []admin.Field{
			CategoryFields.Name,
			CategoryFields.ParentID,
			CategoryFields.SortOrder,
			CategoryFields.IsActive,
			CategoryFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			CategoryFields.IsActive,
			CategoryFields.ParentID,
			CategoryFields.Level,
		},
		SearchFields: []admin.Field{
			CategoryFields.Name,
			CategoryFields.Slug,
			CategoryFields.Description,
		},
		Ordering: []admin.Field{
			CategoryFields.SortOrder,
			CategoryFields.Name,
		},
	})

	// Brand admin
	admin.Register(&admin.Config[Brand]{
		Icon: "Tag",
		ListDisplay: []admin.Field{
			BrandFields.Name,
			BrandFields.Slug,
			BrandFields.WebsiteURL,
			BrandFields.IsActive,
			BrandFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			BrandFields.IsActive,
		},
		SearchFields: []admin.Field{
			BrandFields.Name,
			BrandFields.Slug,
			BrandFields.Description,
		},
		Ordering: []admin.Field{
			BrandFields.Name,
		},
	})

	// Product admin
	admin.Register(&admin.Config[Product]{
		Icon: "Package",
		ListDisplay: []admin.Field{
			ProductFields.Name,
			ProductFields.SKU,
			ProductFields.Price,
			ProductFields.StockQuantity,
			ProductFields.IsActive,
			ProductFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			ProductFields.IsActive,
			ProductFields.IsFeatured,
			ProductFields.CategoryID,
			ProductFields.BrandID,
		},
		SearchFields: []admin.Field{
			ProductFields.Name,
			ProductFields.SKU,
			ProductFields.Description,
		},
		Ordering: []admin.Field{
			ProductFields.CreatedAt.Desc(),
		},
		Actions: []admin.Action[Product]{
			{
				Name:  "activate",
				Label: "Activate Products",
				Handler: func(ctx context.Context, instances []*Product) error {
					for _, p := range instances {
						p.IsActive = true
						// Save would be handled by a manager call or similar depending on ORM context
					}
					return nil
				},
			},
			{
				Name:  "deactivate",
				Label: "Deactivate Products",
				Handler: func(ctx context.Context, instances []*Product) error {
					for _, p := range instances {
						p.IsActive = false
					}
					return nil
				},
			},
		},
	})

	// ProductVariant admin
	admin.Register(&admin.Config[ProductVariant]{
		Icon: "Layers",
		ListDisplay: []admin.Field{
			ProductVariantFields.Name,
			ProductVariantFields.SKU,
			ProductVariantFields.Price,
			ProductVariantFields.StockQuantity,
			ProductVariantFields.IsActive,
		},
		ListFilter: []admin.Field{
			ProductVariantFields.IsActive,
			ProductVariantFields.ProductID,
		},
		SearchFields: []admin.Field{
			ProductVariantFields.Name,
			ProductVariantFields.SKU,
		},
	})

	// ProductImage admin
	admin.Register(&admin.Config[ProductImage]{
		Icon: "Image",
		ListDisplay: []admin.Field{
			ProductImageFields.ProductID,
			ProductImageFields.AltText,
			ProductImageFields.IsPrimary,
			ProductImageFields.SortOrder,
		},
		ListFilter: []admin.Field{
			ProductImageFields.IsPrimary,
			ProductImageFields.ProductID,
		},
	})
}
