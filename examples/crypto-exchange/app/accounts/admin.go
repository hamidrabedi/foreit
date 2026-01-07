package accounts

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers account models with the admin interface.
func RegisterAdmin(ctx context.Context) {
	admin.Register(&admin.Config[UserProfile]{
		Icon: "User",
		ListDisplay: []admin.Field{
			UserProfileFieldsInstance.Id,
			UserProfileFieldsInstance.UserId,
			UserProfileFieldsInstance.DisplayName,
			UserProfileFieldsInstance.Email,
			UserProfileFieldsInstance.Tier,
			UserProfileFieldsInstance.IsActive,
			UserProfileFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			UserProfileFieldsInstance.Tier,
			UserProfileFieldsInstance.IsActive,
		},
		SearchFields: []admin.Field{
			UserProfileFieldsInstance.DisplayName,
			UserProfileFieldsInstance.Email,
		},
		Ordering: []admin.Field{
			UserProfileFieldsInstance.CreatedAt.Desc(),
		},
		Actions: []admin.Action[UserProfile]{
			admin.NewAction("deactivate", "Deactivate", func(ctx context.Context, instances []*UserProfile) error {
				for _, profile := range instances {
					profile.IsActive = false
				}
				return nil
			}).WithDescription("Disable selected accounts").WithDangerous(true),
		},
	})

	admin.Register(&admin.Config[APIKey]{
		Icon: "Key",
		ListDisplay: []admin.Field{
			APIKeyFieldsInstance.Id,
			APIKeyFieldsInstance.UserId,
			APIKeyFieldsInstance.Name,
			APIKeyFieldsInstance.IsActive,
			APIKeyFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			APIKeyFieldsInstance.IsActive,
		},
		SearchFields: []admin.Field{
			APIKeyFieldsInstance.Name,
			APIKeyFieldsInstance.Key,
		},
		Ordering: []admin.Field{
			APIKeyFieldsInstance.CreatedAt.Desc(),
		},
	})
}
