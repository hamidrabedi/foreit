package customers

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers customer models with the admin interface
func RegisterAdmin(ctx context.Context) {
	// CustomerGroup admin
	admin.Register(&admin.Config[CustomerGroup]{
		Icon: "Users",
		ListDisplay: []admin.Field{
			CustomerGroupFields.Name,
			CustomerGroupFields.Code,
			CustomerGroupFields.DiscountPercentage,
			CustomerGroupFields.IsActive,
			CustomerGroupFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			CustomerGroupFields.IsActive,
		},
	})

	// Customer admin
	admin.Register(&admin.Config[Customer]{
		Icon: "User",
		ListDisplay: []admin.Field{
			CustomerFields.Email,
			CustomerFields.FirstName,
			CustomerFields.LastName,
			CustomerFields.IsActive,
			CustomerFields.IsVerified,
			CustomerFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			CustomerFields.IsActive,
			CustomerFields.IsVerified,
		},
		SearchFields: []admin.Field{
			CustomerFields.Email,
			CustomerFields.FirstName,
			CustomerFields.LastName,
		},
	})

	// Address admin
	admin.Register(&admin.Config[Address]{
		Icon: "MapPin",
		ListDisplay: []admin.Field{
			AddressFields.CustomerID,
			AddressFields.AddressType,
			AddressFields.City,
			AddressFields.StateProvince,
			AddressFields.CountryCode,
		},
		ListFilter: []admin.Field{
			AddressFields.AddressType,
			AddressFields.CountryCode,
		},
	})

	// WishList admin
	admin.Register(&admin.Config[WishList]{
		Icon: "Heart",
		ListDisplay: []admin.Field{
			WishListFields.CustomerID,
			WishListFields.Name,
			WishListFields.IsPublic,
			WishListFields.IsDefault,
			WishListFields.CreatedAt,
		},
	})
}
