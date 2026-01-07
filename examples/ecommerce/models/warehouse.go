package models

import "github.com/forgego/forge/schema"

// Warehouse represents a storage location.
type Warehouse struct {
	schema.BaseSchema
}

func (Warehouse) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200),
		schema.String("code").WithRequired().WithMaxLength(50).WithUnique(),
		schema.String("city").WithMaxLength(100).WithOptional(),
		schema.String("country").WithMaxLength(2).WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Warehouse) Meta() schema.Meta {
	return schema.Meta{
		TableName: "warehouses",
	}
}

func (Warehouse) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Warehouse) Hooks() *schema.ModelHooks {
	return nil
}
