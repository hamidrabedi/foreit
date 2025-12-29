package admin

// Complete Production Example
//
// This demonstrates a complete production setup for the admin system.

/*
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/forgego/forge/pkg/admin"
	"github.com/forgego/forge/pkg/admin/http"
	httplib "github.com/forgego/forge/pkg/http"
	"your-app/models"
	"your-app/db"
)

func main() {
	// Setup database
	database := db.NewDBFromConfig(cfg)
	
	// Set database on managers
	models.User.Objects.SetDB(database)
	models.Post.Objects.SetDB(database)

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
				admin.BoolField(
					"is_staff",
					func(u *models.User) bool { return u.IsStaff },
					func(u *models.User, v bool) { u.IsStaff = v },
				),
			},
			SearchFields: []admin.FieldExpr[*models.User, any]{
				admin.StringField("username",
					func(u *models.User) string { return u.Username },
					func(u *models.User, v string) { u.Username = v },
				),
				admin.StringField("email",
					func(u *models.User) string { return u.Email },
					func(u *models.User, v string) { u.Email = v },
				),
			},
			ListFilter: []admin.Filter[*models.User]{
				admin.NewBooleanFilter(
					admin.BoolField("is_active",
						func(u *models.User) bool { return u.IsActive },
						func(u *models.User, v bool) { u.IsActive = v },
					),
				),
				admin.NewBooleanFilter(
					admin.BoolField("is_staff",
						func(u *models.User) bool { return u.IsStaff },
						func(u *models.User, v bool) { u.IsStaff = v },
					),
				),
			},
			Ordering: []admin.Ordering[*models.User]{
				admin.OrderBy(
					admin.StringField("username",
						func(u *models.User) string { return u.Username },
						func(u *models.User, v string) { u.Username = v },
					),
				).Desc(),
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
				admin.NewAction(
					"deactivate",
					"Deactivate selected users",
					func(ctx context.Context, users []*models.User) error {
						for _, user := range users {
							user.IsActive = false
							if err := models.User.Objects.Update(ctx, user); err != nil {
								return err
							}
						}
						return nil
					},
				),
			},
			Fieldsets: []admin.Fieldset[*models.User]{
				admin.NewFieldset(
					"Account Information",
					admin.StringField("username", getter, setter),
					admin.StringField("email", getter, setter),
				),
				admin.NewFieldset(
					"Permissions",
					admin.BoolField("is_active", getter, setter),
					admin.BoolField("is_staff", getter, setter),
				),
			},
		},
	)

	// Register Post model with Comment inline
	postAdmin := admin.Register(
		&models.Post{},
		models.Post.Objects,
		&admin.Config[*models.Post]{
			ListDisplay: []admin.FieldExpr[*models.Post, any]{
				admin.StringField("title", getter, setter),
				admin.StringField("author", getter, setter),
			},
			Inlines: []admin.Inline[*models.Post, any]{
				admin.TabularInline(
					&models.Comment{},
					models.Comment.Objects,
					admin.FieldExpr[*models.Comment, *models.Post]{
						// Parent field (foreign key)
					},
					[]admin.FieldExpr[*models.Comment, any]{
						admin.StringField("author", getter, setter),
						admin.StringField("content", getter, setter),
					},
				),
			},
		},
	)

	// Register for HTTP handlers
	http.RegisterAdminForHTTP(userAdmin)
	http.RegisterAdminForHTTP(postAdmin)

	// Setup router
	router := httplib.NewRouter()

	// Register admin routes
	adminRouter := http.NewRouter(admin.GetGlobalRegistry())
	adminRouter.RegisterRoutes(router, "/admin")

	// Add authentication middleware if needed
	// router.Use(authMiddleware)

	// Start server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
*/
