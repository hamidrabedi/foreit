//go:build ignore_examples
// +build ignore_examples

package examples

import (
	"context"
	"fmt"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
)

// ExampleAdvancedUsage demonstrates advanced admin features
func ExampleAdvancedUsage() {
	schemaInstance := &ExampleSchema{}
	var manager *orm.Manager[ExampleModel]

	config := &admin.Config[ExampleModel]{
		VerboseName:       "Example Model",
		VerboseNamePlural: "Example Models",

		// Custom queryset filtering
		GetQueryset: func(ctx context.Context, admin *admin.Admin[ExampleModel], qs orm.QuerySet[ExampleModel]) (orm.QuerySet[ExampleModel], error) {
			// Only show active items
			// This would use ORM field expressions
			return qs, nil
		},

		// Custom save logic
		SaveModel: func(ctx context.Context, admin *admin.Admin[ExampleModel], instance *ExampleModel, formData admin.FormData, isNew bool) error {
			// Basic validation - in production, you'd use a proper validation library
			if instance.Name == "" {
				return fmt.Errorf("name is required")
			}
			if instance.Email == "" {
				return fmt.Errorf("email is required")
			}

			// Custom save logic
			if isNew {
				return manager.Create(ctx, instance)
			}
			return manager.Update(ctx, instance)
		},

		// Custom permissions
		HasAddPermission: func(ctx context.Context, admin *admin.Admin[ExampleModel], user interface{}) bool {
			// Check if user can add
			return true
		},

		HasChangePermission: func(ctx context.Context, admin *admin.Admin[ExampleModel], user interface{}, obj *ExampleModel) bool {
			// Check if user can change this specific object
			return true
		},
	}

	admin, err := admin.Register(schemaInstance, manager, config)
	if err != nil {
		panic(err)
	}

	_ = admin
}
