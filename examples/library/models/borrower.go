package models

import (
	"github.com/forgego/forge/pkg/schema"
)

// Borrower represents a library member who can borrow books
type Borrower struct {
	schema.BaseSchema
}

// Fields returns all field definitions for Borrower
func (Borrower) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("first_name").Required().MaxLength(100).VerboseName("First Name").Build(),
		schema.String("last_name").Required().MaxLength(100).VerboseName("Last Name").Build(),
		schema.String("email").Required().MaxLength(255).Unique().VerboseName("Email").Build(),
		schema.String("phone").MaxLength(20).VerboseName("Phone").Build(),
		schema.String("address").VerboseName("Address").Build(),
		schema.Bool("is_active").Default(true).VerboseName("Active").Build(),
		schema.Int64("max_books").Default(5).VerboseName("Max Books").Build(),
		schema.Time("joined_at").AutoNowAdd().VerboseName("Joined Date").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (Borrower) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "borrowers",
		VerboseName:       "Borrower",
		VerboseNamePlural: "Borrowers",
		OrderBy:           []string{"last_name", "first_name"},
		Indexes: []schema.Index{
			{Name: "idx_borrower_email", Fields: []string{"email"}, Unique: true},
		},
	}
}

// Relations returns all relationship definitions
func (Borrower) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Hooks returns model lifecycle hooks
func (Borrower) Hooks() *schema.ModelHooks {
	return nil
}

