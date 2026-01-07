package models

import "github.com/forgego/forge/schema"

// Category represents a product category.
type Category struct {
	schema.BaseSchema
}

func (Category) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200),
		schema.String("slug").WithRequired().WithMaxLength(200).WithUnique(),
		schema.Text("description").WithOptional(),
		schema.Int64("parent_id").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Category) Meta() schema.Meta {
	return schema.Meta{
		TableName: "categories",
	}
}

func (Category) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("parent_id", "Category").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("children"),
	}
}

func (Category) Hooks() *schema.ModelHooks {
	return nil
}
