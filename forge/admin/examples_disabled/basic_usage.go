//go:build ignore_examples
// +build ignore_examples

package examples

import (
	"context"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// ExampleModel represents a basic model
type ExampleModel struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	IsActive  bool   `db:"is_active"`
	CreatedAt string `db:"created_at"`
}

// ExampleSchema implements schema.Schema for ExampleModel
type ExampleSchema struct{}

func (s *ExampleSchema) Fields() []schema.Field {
	return []schema.Field{
		{
			Name:          "ID",
			Type:          schema.TypeInt64,
			PrimaryKey:    true,
			AutoIncrement: true,
			Required:      true,
			VerboseName:   "ID",
		},
		{
			Name:        "Name",
			Type:        schema.TypeString,
			Required:    true,
			MaxLength:   intPtr(255),
			VerboseName: "Name",
			HelpText:    "The name of the item",
		},
		{
			Name:        "Email",
			Type:        schema.TypeEmail,
			Required:    true,
			VerboseName: "Email",
			HelpText:    "Email address",
		},
		{
			Name:        "IsActive",
			Type:        schema.TypeBool,
			Default:     true,
			VerboseName: "Is Active",
		},
		{
			Name:        "CreatedAt",
			Type:        schema.TypeDateTime,
			AutoNowAdd:  true,
			Editable:    false, // Read-only
			VerboseName: "Created At",
		},
	}
}

func (s *ExampleSchema) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (s *ExampleSchema) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "example_models",
		VerboseName:       "Example Model",
		VerboseNamePlural: "Example Models",
		AppLabel:          "examples",
	}
}

func (s *ExampleSchema) Hooks() *schema.ModelHooks {
	return nil
}

// ExampleBasicUsage demonstrates basic admin registration
func ExampleBasicUsage() {
	// Create admin configuration
	config := &admin.Config[ExampleModel]{
		VerboseName:       "Example Model",
		VerboseNamePlural: "Example Models",
		ListPerPage:       20,
		// Fields will be auto-discovered from schema
	}

	// Register admin
	admin, err := admin.Register(config)
	if err != nil {
		panic(err)
	}

	_ = admin
}

// ExampleWithCustomConfig demonstrates admin with custom configuration
func ExampleWithCustomConfig() {
	schemaInstance := &ExampleSchema{}
	var manager *orm.Manager[ExampleModel]

	config := &admin.Config[ExampleModel]{
		VerboseName:       "Example Model",
		VerboseNamePlural: "Example Models",
		ListPerPage:       50,
		// Custom list display using field names
		ListDisplay: []interface{}{
			"Name",
			"Email",
			"IsActive",
		},
		// Search fields
		SearchFields: []interface{}{
			"Name",
			"Email",
		},
		// Ordering
		Ordering: []admin.Ordering[ExampleModel]{
			*admin.OrderBy[ExampleModel]("CreatedAt").Desc(),
		},
	}

	admin, err := admin.Register(schemaInstance, manager, config)
	if err != nil {
		panic(err)
	}

	_ = admin
}

// ExampleWithActions demonstrates admin with custom actions
func ExampleWithActions() {
	schemaInstance := &ExampleSchema{}
	var manager *orm.Manager[ExampleModel]

	config := &admin.Config[ExampleModel]{
		VerboseName:       "Example Model",
		VerboseNamePlural: "Example Models",
		Actions: []admin.Action[ExampleModel]{
			admin.NewAction[ExampleModel](
				"activate_selected",
				"Activate selected",
				func(ctx context.Context, instances []*ExampleModel) error {
					for _, instance := range instances {
						instance.IsActive = true
						if err := manager.Update(ctx, instance); err != nil {
							return err
						}
					}
					return nil
				},
			).WithDescription("Activate the selected items"),
		},
	}

	admin, err := admin.Register(schemaInstance, manager, config)
	if err != nil {
		panic(err)
	}

	_ = admin
}

// Helper function
func intPtr(i int) *int {
	return &i
}
