package models

import (
	"github.com/forgego/forge/pkg/schema"
)

// Category represents a book category
type Category struct {
	schema.BaseSchema
}

// Fields returns all field definitions for Category
func (Category) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(100).Unique().VerboseName("Name").Build(),
		schema.String("slug").Required().MaxLength(100).Unique().VerboseName("Slug").Build(),
		schema.String("description").VerboseName("Description").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (Category) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "categories",
		VerboseName:       "Category",
		VerboseNamePlural: "Categories",
		OrderBy:           []string{"name"},
		Indexes: []schema.Index{
			{Name: "idx_category_slug", Fields: []string{"slug"}, Unique: true},
		},
	}
}

// Relations returns all relationship definitions
func (Category) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Hooks returns model lifecycle hooks
func (Category) Hooks() *schema.ModelHooks {
	return nil
}

