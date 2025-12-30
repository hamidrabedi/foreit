package examples

import (
	"context"
	"fmt"

	"github.com/forgego/forge/admin"
	adminvalidation "github.com/forgego/forge/admin/validation"
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
			// Validate before saving
			validator := adminvalidation.NewFormValidator(schemaInstance)
			errors, err := validator.Validate(formData)
			if err != nil {
				return err
			}
			if len(errors) > 0 {
				// Return validation errors
				return fmt.Errorf("validation failed: %v", errors)
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
