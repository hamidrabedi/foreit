package models

import (
	"github.com/forgego/forge/schema"
)

// Author represents a book author
type Author struct {
	schema.BaseSchema
}

// Fields returns all field definitions for Author
func (Author) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(200).VerboseName("Name").Build(),
		schema.String("email").MaxLength(255).VerboseName("Email").Build(),
		schema.String("bio").VerboseName("Biography").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (Author) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "authors",
		VerboseName:       "Author",
		VerboseNamePlural: "Authors",
		OrderBy:           []string{"name"},
		Indexes: []schema.Index{
			{Name: "idx_author_email", Fields: []string{"email"}, Unique: true},
		},
	}
}

// Relations returns all relationship definitions
func (Author) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Hooks returns model lifecycle hooks
func (Author) Hooks() *schema.ModelHooks {
	return nil
}
