package models

import "github.com/forgego/forge/schema"

// Supplier represents a product supplier.
type Supplier struct {
	schema.BaseSchema
}

func (Supplier) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200),
		schema.String("slug").WithRequired().WithMaxLength(200).WithUnique(),
		schema.String("email").WithMaxLength(255).WithOptional(),
		schema.String("phone").WithMaxLength(20).WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Supplier) Meta() schema.Meta {
	return schema.Meta{
		TableName: "suppliers",
	}
}

func (Supplier) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Supplier) Hooks() *schema.ModelHooks {
	return nil
}
