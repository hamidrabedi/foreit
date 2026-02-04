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
			CustomerGroupFieldsInstance.Name,
			CustomerGroupFieldsInstance.Code,
			CustomerGroupFieldsInstance.DiscountPercentage,
			CustomerGroupFieldsInstance.IsActive,
			CustomerGroupFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			CustomerGroupFieldsInstance.IsActive,
		},
	})

	// Customer admin
	admin.Register(&admin.Config[Customer]{
		Icon: "User",
		ListDisplay: []admin.Field{
			CustomerFieldsInstance.Email,
			CustomerFieldsInstance.FirstName,
			CustomerFieldsInstance.LastName,
			CustomerFieldsInstance.IsActive,
			CustomerFieldsInstance.IsVerified,
			CustomerFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			CustomerFieldsInstance.IsActive,
			CustomerFieldsInstance.IsVerified,
		},
		SearchFields: []admin.Field{
			CustomerFieldsInstance.Email,
			CustomerFieldsInstance.FirstName,
			CustomerFieldsInstance.LastName,
		},
	})

	// Address admin
	admin.Register(&admin.Config[Address]{
		Icon: "MapPin",
		ListDisplay: []admin.Field{
			AddressFieldsInstance.CustomerId,
			AddressFieldsInstance.AddressType,
			AddressFieldsInstance.City,
			AddressFieldsInstance.StateProvince,
			AddressFieldsInstance.CountryCode,
		},
		ListFilter: []admin.Field{
			AddressFieldsInstance.AddressType,
			AddressFieldsInstance.CountryCode,
		},
	})

	// WishList admin
	admin.Register(&admin.Config[WishList]{
		Icon: "Heart",
		ListDisplay: []admin.Field{
			WishListFieldsInstance.CustomerId,
			WishListFieldsInstance.Name,
			WishListFieldsInstance.IsPublic,
			WishListFieldsInstance.IsDefault,
			WishListFieldsInstance.CreatedAt,
		},
	})
}
