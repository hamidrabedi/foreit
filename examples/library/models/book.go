package models

import (
	"github.com/forgego/forge/schema"
)

// Book represents a book in the library
type Book struct {
	schema.BaseSchema
}

// Fields returns all field definitions for Book
func (Book) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("title").Required().MaxLength(500).VerboseName("Title").Build(),
		schema.String("isbn").MaxLength(20).Unique().VerboseName("ISBN").Build(),
		schema.Int64("author_id").Required().VerboseName("Author ID").Build(),
		schema.Int64("category_id").VerboseName("Category ID").Build(),
		schema.String("description").VerboseName("Description").Build(),
		schema.Int64("pages").Default(0).VerboseName("Pages").Build(),
		schema.Bool("available").Default(true).VerboseName("Available").Build(),
		schema.Float64("price").Default(0.0).VerboseName("Price").Build(),
		schema.Time("published_at").VerboseName("Published Date").Build(),
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (Book) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "books",
		VerboseName:       "Book",
		VerboseNamePlural: "Books",
		OrderBy:           []string{"-created_at", "title"},
		Indexes: []schema.Index{
			{Name: "idx_book_author", Fields: []string{"author_id"}, Unique: false},
			{Name: "idx_book_category", Fields: []string{"category_id"}, Unique: false},
			{Name: "idx_book_isbn", Fields: []string{"isbn"}, Unique: true},
			{Name: "idx_book_available", Fields: []string{"available"}, Unique: false},
		},
	}
}

// Relations returns all relationship definitions
func (Book) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Hooks returns model lifecycle hooks
func (Book) Hooks() *schema.ModelHooks {
	return nil
}
