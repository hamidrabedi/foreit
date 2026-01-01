package customers

import (
	"context"
	
	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/db"
)

// RegisterAdmin registers customer models with the admin interface
func RegisterAdmin(ctx context.Context, registry *admin.Registry, database *db.DB) {
	// CustomerGroup admin
	customerGroupConfig := &admin.ModelConfig{
		Name:             "Customer Group",
		PluralName:       "Customer Groups",
		Icon:             "👥",
		ListDisplay:      []string{"id", "name", "code", "discount_percentage", "is_active", "created_at"},
		ListFilter:       []string{"is_active"},
		SearchFields:     []string{"name", "code", "description"},
		OrderBy:          []string{"name"},
		PerPage:          20,
		Actions:          []string{"delete", "activate", "deactivate"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("CustomerGroup", &CustomerGroup{}, customerGroupConfig)
	
	// Customer admin
	customerConfig := &admin.ModelConfig{
		Name:             "Customer",
		PluralName:       "Customers",
		Icon:             "👤",
		ListDisplay:      []string{"id", "email", "first_name", "last_name", "customer_group_id", "total_orders", "total_spent", "is_active", "is_verified", "created_at"},
		ListFilter:       []string{"is_active", "is_verified", "customer_group_id", "accepts_marketing"},
		SearchFields:     []string{"email", "first_name", "last_name", "company_name"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"delete", "activate", "deactivate", "verify", "export"},
		ExportFormats:    []string{"csv", "json"},
		BulkActions:      true,
	}
	registry.Register("Customer", &Customer{}, customerConfig)
	
	// Address admin
	addressConfig := &admin.ModelConfig{
		Name:             "Address",
		PluralName:       "Addresses",
		Icon:             "📍",
		ListDisplay:      []string{"id", "customer_id", "address_type", "city", "state_province", "country_code", "is_default_shipping", "is_default_billing"},
		ListFilter:       []string{"address_type", "country_code", "is_default_shipping", "is_default_billing"},
		SearchFields:     []string{"first_name", "last_name", "address_line1", "city", "postal_code"},
		OrderBy:          []string{"customer_id", "-is_default_shipping"},
		PerPage:          20,
		Actions:          []string{"delete"},
	}
	registry.Register("Address", &Address{}, addressConfig)
	
	// WishList admin
	wishListConfig := &admin.ModelConfig{
		Name:             "Wish List",
		PluralName:       "Wish Lists",
		Icon:             "❤️",
		ListDisplay:      []string{"id", "customer_id", "name", "is_public", "is_default", "created_at"},
		ListFilter:       []string{"is_public", "is_default"},
		SearchFields:     []string{"name", "description"},
		OrderBy:          []string{"customer_id", "-is_default"},
		PerPage:          20,
		Actions:          []string{"delete"},
	}
	registry.Register("WishList", &WishList{}, wishListConfig)
	
	// WishListItem admin
	wishListItemConfig := &admin.ModelConfig{
		Name:             "Wish List Item",
		PluralName:       "Wish List Items",
		Icon:             "💝",
		ListDisplay:      []string{"id", "wish_list_id", "product_id", "variant_id", "desired_quantity", "priority", "created_at"},
		ListFilter:       []string{"wish_list_id", "priority"},
		SearchFields:     []string{"notes"},
		OrderBy:          []string{"wish_list_id", "-priority"},
		PerPage:          20,
		Actions:          []string{"delete"},
	}
	registry.Register("WishListItem", &WishListItem{}, wishListItemConfig)
}
