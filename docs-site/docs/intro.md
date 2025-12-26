---
sidebar_position: 1
---

# Welcome to forge

**forge** is a Django-like Go framework that brings the best of Django's developer experience to Go, with full type safety and modern Go features.

## What is forge?

forge is a full-stack web framework for Go that provides:

- **Type-Safe ORM** - Full Django ORM features with compile-time type checking
- **Code Generation** - AST-based code generation for models, managers, and querysets
- **Auto-Generated Admin** - Django-style admin interface auto-generated from your models
- **REST API** - Built-in REST API system similar to Django REST Framework
- **Migration System** - Built-in migration system with golang-migrate
- **Security** - Built-in CSRF, XSS, and SQL injection protection
- **Extensible** - Everything is extendable/overridable via plugins

## Why forge?

If you've ever worked with Django and loved its developer experience, but needed the performance and type safety of Go, forge is for you. It combines:

- **Django's ease of use** - Declarative models, auto-generated admin, sensible defaults
- **Go's performance** - Fast, compiled, efficient
- **Type safety** - Compile-time checking, no runtime surprises
- **Modern Go** - Uses generics, interfaces, and Go best practices

## Quick Example

Here's what a forge model looks like:

```go
type Post struct {
    schema.BaseSchema
}

func (Post) Fields() []schema.Field {
    return []schema.Field{
        schema.Int64("id").Primary().AutoIncrement().Build(),
        schema.String("title").Required().MaxLength(200).Build(),
        schema.String("content").Required().Build(),
        schema.Bool("published").Default(false).Build(),
        schema.Time("created_at").AutoNowAdd().Build(),
    }
}
```

Then use it with type-safe queries:

```go
// Get all published posts
posts, err := Post.Objects.
    Filter(Post.Fields.Published.Equals(true)).
    OrderBy("-created_at").
    All(ctx)
```

## Getting Started

Ready to build your first forge application? Check out our [Installation Guide](/docs/getting-started/installation) and [Quick Start Tutorial](/docs/getting-started/quickstart).

## Key Features

### Type-Safe ORM

Write queries with full compile-time type checking:

```go
users, err := User.Objects.
    Filter(User.Fields.IsActive.Equals(true)).
    OrderBy("-date_joined").
    Limit(10).
    All(ctx)
```

### Auto-Generated Admin

Just register your models and get a full admin interface:

```go
admin.RegisterModel(&models.Post{})
// Visit /admin/ to see your models!
```

### Code Generation

Generate type-safe code from your schema definitions:

```bash
forge generate
```

This creates type-safe managers, querysets, and field expressions automatically.

### REST API

Build APIs for React, Vue, or any frontend:

```go
viewset := api.NewBaseViewSet(
    NewPostSerializer(),
    Post.Objects.Filter(),
    &Post{},
)
```

## What's Next?

- [Install forge](/docs/getting-started/installation)
- [Follow the Quick Start](/docs/getting-started/quickstart)
- [Read the Guides](/docs/guides/models)
- [Check out Examples](/docs/examples/blog)

