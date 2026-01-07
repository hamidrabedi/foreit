package models

import "github.com/forgego/forge/schema"

// Author represents an author
type Author struct {
	schema.BaseSchema
}

func (Author) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(255),
		schema.String("email").WithRequired().WithMaxLength(255).WithUnique(),
	}
}

func (Author) Meta() schema.Meta {
	return schema.Meta{
		TableName: "authors",
	}
}

func (Author) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Article with foreign key to Author
type Article struct {
	schema.BaseSchema
}

func (Article) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("title").WithRequired().WithMaxLength(200),
		schema.Text("body"),
		schema.Int64("author_id").WithRequired(),
	}
}

func (Article) Meta() schema.Meta {
	return schema.Meta{
		TableName: "articles",
	}
}

func (Article) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("author_id", "Author").WithOnDelete("CASCADE"),
	}
}

// Profile with OneToOne relationship to Author
type Profile struct {
	schema.BaseSchema
}

func (Profile) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Text("bio"),
		schema.String("website").WithMaxLength(255),
		schema.Int64("author_id").WithRequired().WithUnique(),
	}
}

func (Profile) Meta() schema.Meta {
	return schema.Meta{
		TableName: "profiles",
	}
}

func (Profile) Relations() []schema.Relation {
	return []schema.Relation{
		schema.OneToOne("author_id", "Author").WithOnDelete("CASCADE"),
	}
}

// Tag for many-to-many relationship
type Tag struct {
	schema.BaseSchema
}

func (Tag) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(100).WithUnique(),
		schema.String("slug").WithRequired().WithMaxLength(100).WithUnique(),
	}
}

func (Tag) Meta() schema.Meta {
	return schema.Meta{
		TableName: "tags",
	}
}

func (Tag) Relations() []schema.Relation {
	return []schema.Relation{}
}

// ArticleTag for many-to-many relationship between Article and Tag
type ArticleTag struct {
	schema.BaseSchema
}

func (ArticleTag) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("article_id").WithRequired(),
		schema.Int64("tag_id").WithRequired(),
	}
}

func (ArticleTag) Meta() schema.Meta {
	return schema.Meta{
		TableName: "article_tags",
		Indexes: []schema.Index{
			schema.Index("article_id", "tag_id").WithUnique(),
		},
	}
}

func (ArticleTag) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("article_id", "Article").WithOnDelete("CASCADE"),
		schema.ForeignKey("tag_id", "Tag").WithOnDelete("CASCADE"),
	}
}
