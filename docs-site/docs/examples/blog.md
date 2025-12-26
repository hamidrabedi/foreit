---
sidebar_position: 1
---

# Blog Example

A complete blog application demonstrating forge features.

## Overview

This example shows how to build a blog with:

- Post model with author relationship
- Category model with many-to-many relationship
- Admin interface for content management
- REST API for frontend consumption

## Models

### Post Model

```go
package models

import (
    "github.com/forgego/forge/pkg/schema"
    "github.com/forgego/forge/pkg/schema/relations"
)

type Post struct {
    schema.BaseSchema
}

func (Post) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("title").Required().MaxLength(200).Build(),
        schema.String("slug").Unique().MaxLength(200).Build(),
        schema.Text("content").Required().Build(),
        schema.Text("excerpt").MaxLength(500).Build(),
        schema.Bool("published").Default(false).Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
        schema.Time("updated_at").AutoNow().Build(),
        schema.Time("published_at").Build(),
    }
}

func (Post) Relations() []schema.Relation {
    return []schema.Relation{
        relations.ForeignKey("author", "User").
            Required().
            OnDelete(schema.Cascade).
            RelatedName("posts"),
        relations.ManyToMany("categories", "Category").
            RelatedName("posts"),
    }
}

func (Post) Meta() schema.Meta {
    return schema.Meta{
        TableName:        "posts",
        VerboseName:      "Post",
        VerboseNamePlural: "Posts",
        OrderBy:          []string{"-created_at"},
    }
}

func (Post) Hooks() *schema.ModelHooks {
    return nil
}
```

### Category Model

```go
type Category struct {
    schema.BaseSchema
}

func (Category) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("name").Required().Unique().MaxLength(100).Build(),
        schema.String("slug").Required().Unique().MaxLength(100).Build(),
        schema.Text("description").Build(),
    }
}

func (Category) Meta() schema.Meta {
    return schema.Meta{
        TableName:        "categories",
        VerboseName:      "Category",
        VerboseNamePlural: "Categories",
    }
}

func (Category) Relations() []schema.Relation {
    return []schema.Relation{}
}

func (Category) Hooks() *schema.ModelHooks {
    return nil
}
```

## Views

### List Posts

```go
package views

import (
    "context"
    "encoding/json"
    "net/http"
    "myblog/models"
)

func ListPosts(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    
    posts, err := models.Post.Objects.
        Filter(models.Post.Fields.Published.Equals(true)).
        OrderBy("-created_at").
        PrefetchRelated("author", "categories").
        All(ctx)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
```

### Get Post

```go
func GetPost(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    
    // Extract ID from URL
    id := extractID(r.URL.Path)
    
    post, err := models.Post.Objects.
        Filter(models.Post.Fields.ID.Equals(id)).
        Filter(models.Post.Fields.Published.Equals(true)).
        SelectRelated("author").
        PrefetchRelated("categories").
        Get(ctx)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}
```

## Admin Setup

```go
admin.RegisterModelWithOptions(
    &models.Post{},
    admin.WithListDisplay("title", "author", "published", "created_at"),
    admin.WithSearchFields("title", "content"),
    admin.WithListFilter("published", "author", "created_at"),
    admin.WithDateHierarchy("created_at"),
)

admin.RegisterModelWithOptions(
    &models.Category{},
    admin.WithListDisplay("name", "slug"),
    admin.WithSearchFields("name"),
)
```

## REST API

```go
func RegisterPostViewSet(router *api.Router) {
    viewset := api.NewBaseViewSet(
        func() api.Serializer {
            return NewPostSerializer()
        },
        models.Post.Objects.Filter(models.Post.Fields.Published.Equals(true)),
        &models.Post{},
    )
    
    router.Register("posts", viewset)
}
```

## Usage

### Create a Post

```go
post := &models.Post{
    Title:     "My First Post",
    Slug:      "my-first-post",
    Content:   "This is the content...",
    Published: true,
    Author:    author,
}
err := models.Post.Objects.Create(ctx, post)
```

### Get Published Posts

```go
posts, err := models.Post.Objects.
    Filter(models.Post.Fields.Published.Equals(true)).
    OrderBy("-created_at").
    All(ctx)
```

### Filter by Category

```go
posts, err := models.Post.Objects.
    Filter(models.Post.Fields.Categories.Contains(categoryID)).
    All(ctx)
```

## See Also

- [Models Guide](../guides/models) - Learn about models
- [Queries Guide](../guides/queries) - Query examples
- [Admin Guide](../guides/admin) - Admin customization

