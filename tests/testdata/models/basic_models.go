package models

import "github.com/forgego/forge/schema"

// User is a basic user model
type User struct {
	schema.BaseSchema
}

func (User) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("email").WithRequired().WithMaxLength(255).WithUnique(),
		schema.String("username").WithRequired().WithMaxLength(150),
		schema.Bool("is_active").WithDefault(true),
		schema.DateTime("created_at"),
		schema.DateTime("updated_at"),
	}
}

func (User) Meta() schema.Meta {
	return schema.Meta{
		TableName: "users",
	}
}

func (User) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Post is a basic post model
type Post struct {
	schema.BaseSchema
}

func (Post) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("title").WithRequired().WithMaxLength(200),
		schema.Text("content"),
		schema.String("status").WithDefault("draft").WithMaxLength(50),
		schema.DateTime("published_at"),
		schema.DateTime("created_at"),
		schema.DateTime("updated_at"),
	}
}

func (Post) Meta() schema.Meta {
	return schema.Meta{
		TableName: "posts",
	}
}

func (Post) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Comment is a basic comment model
type Comment struct {
	schema.BaseSchema
}

func (Comment) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Text("content").WithRequired(),
		schema.Bool("is_approved").WithDefault(false),
		schema.DateTime("created_at"),
		schema.DateTime("updated_at"),
	}
}

func (Comment) Meta() schema.Meta {
	return schema.Meta{
		TableName: "comments",
	}
}

func (Comment) Relations() []schema.Relation {
	return []schema.Relation{}
}
