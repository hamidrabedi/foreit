# Endpoints Module

API framework for building RESTful endpoints with type-safe resource handlers.

## Concepts

### Resource Handlers

Resource handlers define CRUD operations for a model:

```go
type UserResource struct {
    endpoints.Resource[models.User, *ent.UserQuery]
}

func (r *UserResource) Index(ctx *endpoints.Context) ([]*models.User, error) {
    return r.Repo.Query().All(ctx.Request.Context())
}
```

### Serializers

Serializers handle request/response transformation:

```go
serializer := endpoints.NewModelSerializer[models.User]()
if err := serializer.Validate(data); err != nil {
    return err
}
user, _ := serializer.ToInternalValue(data)
```

### Router

Router registers resources and generates routes:

```go
router := endpoints.NewRouter(app, "/api")
router.RegisterResource("users", &UserResource{})
```

## Features

- Type-safe resource handlers
- Automatic route generation
- Request/response serialization
- Query parameter processing (filters, pagination, sorting)
- Error handling
- Content negotiation

