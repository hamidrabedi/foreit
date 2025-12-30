package admin

// Example usage of the admin system
//
// This file demonstrates how to use the new type-safe admin system.
// Copy this pattern to register your models.

/*
import (
	"context"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/admin/http"
	httplib "github.com/forgego/forge/server"
	"your-app/models"
)

func SetupAdmin(router *httplib.Router) {
	// Register User model
	userAdmin := admin.Register(
		&models.User{},
		models.User.Objects,
		&admin.Config[*models.User]{
			ListDisplay: []admin.FieldExpr[*models.User, any]{
				admin.StringField(
					"username",
					func(u *models.User) string { return u.Username },
					func(u *models.User, v string) { u.Username = v },
				),
				admin.StringField(
					"email",
					func(u *models.User) string { return u.Email },
					func(u *models.User, v string) { u.Email = v },
				),
				admin.BoolField(
					"is_active",
					func(u *models.User) bool { return u.IsActive },
					func(u *models.User, v bool) { u.IsActive = v },
				),
			},
			SearchFields: []admin.FieldExpr[*models.User, any]{
				admin.StringField("username", getter, setter),
				admin.StringField("email", getter, setter),
			},
			ListFilter: []admin.Filter[*models.User]{
				admin.NewBooleanFilter(
					admin.BoolField("is_active", getter, setter),
				),
			},
			Ordering: []admin.Ordering[*models.User]{
				admin.OrderBy(usernameField).Desc(),
			},
			ListPerPage: 25,
			Actions: []admin.Action[*models.User]{
				admin.NewAction(
					"activate",
					"Activate selected users",
					func(ctx context.Context, users []*models.User) error {
						for _, user := range users {
							user.IsActive = true
							if err := models.User.Objects.Update(ctx, user); err != nil {
								return err
							}
						}
						return nil
					},
				),
			},
		},
	)

	// Register for HTTP handlers
	http.RegisterAdminForHTTP(userAdmin)

	// Register routes
	adminRouter := http.NewRouter(admin.GetGlobalRegistry())
	adminRouter.RegisterRoutes(router, "/admin")
}

// Get global registry
func GetGlobalRegistry() *admin.Registry {
	return globalRegistry
}
*/
