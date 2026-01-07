package models

import "github.com/forgego/forge/schema"

// ProductWithJSONB demonstrates JSONB columns
type ProductWithJSONB struct {
	schema.BaseSchema
}

func (ProductWithJSONB) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(255),
		schema.JSON("attributes"), // Maps to JSONB in Postgres
		schema.JSON("metadata"),
	}
}

func (ProductWithJSONB) Meta() schema.Meta {
	return schema.Meta{
		TableName: "products_jsonb",
	}
}

func (ProductWithJSONB) Relations() []schema.Relation {
	return []schema.Relation{}
}

// UserWithUUID demonstrates UUID fields
type UserWithUUID struct {
	schema.BaseSchema
}

func (UserWithUUID) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.UUID("uuid").WithRequired().WithDBDefault("gen_random_uuid()"),
		schema.String("email").WithRequired().WithMaxLength(255).WithUnique(),
	}
}

func (UserWithUUID) Meta() schema.Meta {
	return schema.Meta{
		TableName: "users_uuid",
	}
}

func (UserWithUUID) Relations() []schema.Relation {
	return []schema.Relation{}
}

// DocumentWithTimestamps demonstrates timestamp with time zone
type DocumentWithTimestamps struct {
	schema.BaseSchema
}

func (DocumentWithTimestamps) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("title").WithRequired().WithMaxLength(255),
		schema.DateTime("created_at"),
		schema.DateTime("updated_at"),
	}
}

func (DocumentWithTimestamps) Meta() schema.Meta {
	return schema.Meta{
		TableName: "documents",
	}
}

func (DocumentWithTimestamps) Relations() []schema.Relation {
	return []schema.Relation{}
}
