package models

import "github.com/forgego/forge/schema"

// Brand represents a product brand.
type Brand struct {
	schema.BaseSchema
}

func (Brand) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200),
		schema.String("slug").WithRequired().WithMaxLength(200).WithUnique(),
		schema.Text("description").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Brand) Meta() schema.Meta {
	return schema.Meta{
		TableName: "brands",
	}
}

func (Brand) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Brand) Hooks() *schema.ModelHooks {
	return nil
}
