# Static Design Patterns for Gogo

## Philosophy

Replace reflection with Go generics and Ent's code generation for type safety and performance.

## Pattern 1: Repository Pattern

**Instead of**: Reflection-based helpers
**Use**: Type-safe repositories

```go
// Generic repository interface
type Repository[T any, Q any] interface {
    Query() Q
    GetByID(ctx context.Context, id interface{}) (*T, error)
    Create(ctx context.Context, data *T) (*T, error)
    Update(ctx context.Context, id interface{}, data *T) (*T, error)
    Delete(ctx context.Context, id interface{}) error
}

// Concrete implementation using Ent types
type UserRepository struct {
    client *ent.Client
}

func (r *UserRepository) Query() *ent.UserQuery {
    return r.client.User.Query()
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
    return r.client.User.Get(ctx, id)
}

// No reflection needed - Ent generates everything
```

## Pattern 2: Resource Handlers

**Instead of**: Generic handlers with reflection
**Use**: Type-safe resource handlers

```go
// Resource base with generics
type Resource[T any, Q any] struct {
    Repo Repository[T, Q]
}

// Concrete resource
type UserResource struct {
    Resource[models.User, *ent.UserQuery]
}

func (r *UserResource) Index(ctx *endpoints.Context) ([]models.User, error) {
    // Type-safe query
    return r.Repo.Query().
        Where(user.ActiveEQ(true)).
        Limit(10).
        All(ctx.Request.Context())
}
```

## Pattern 3: Serializers with Ent Descriptors

**Instead of**: Reflection-based field setters
**Use**: Ent's generated field descriptors

```go
// Serializer uses Ent's generated types
type Serializer[T any] struct {
    fields map[string]FieldDescriptor
}

// Ent generates field descriptors at compile time
// No reflection needed to access them
func NewModelSerializer[T any]() *Serializer[T] {
    // Extract from Ent's generated code
    // UserFields.Name, UserFields.Email, etc.
    return &Serializer[T]{
        fields: extractFieldsFromEnt[T](),
    }
}

func (s *Serializer[T]) ToInternalValue(data map[string]interface{}) (*T, error) {
    // Use Ent's Create builder (generated)
    // No MethodByName calls needed
    var zero T
    builder := getCreateBuilder[T](data)
    return builder.Save(ctx)
}
```

## Pattern 4: Query Processors

**Instead of**: Reflection-based query building
**Use**: Ent's predicate system

```go
// Query processor converts HTTP params to Ent predicates
type QueryProcessor[Q any] struct {
    Query Q
}

func (p *QueryProcessor[Q]) ApplyFilter(field, op, value string) error {
    // Use Ent's generated predicate functions
    // No reflection - Ent generates WhereX, WhereXEQ, etc.
    switch op {
    case "eq":
        p.Query = applyEQ(p.Query, field, value)
    case "gt":
        p.Query = applyGT(p.Query, field, value)
    // ...
    }
    return nil
}

// Type-safe predicate application
func applyEQ[Q any](q Q, field, value string) Q {
    // Ent generates type-safe methods
    // For User: q.Where(user.NameEQ(value))
    // Compile-time checked!
}
```

## Pattern 5: Console Registration

**Instead of**: Reflection-based model registration
**Use**: Generic registration

```go
// Console registry with generics
type ConsoleRegistry struct {
    consoles map[string]ConsoleInfo
}

type ConsoleInfo[T any] struct {
    Model  T
    Console Console[T]
}

func Register[T any](console Console[T]) {
    var zero T
    name := getTypeName(zero)
    registry.consoles[name] = ConsoleInfo[T]{
        Model: zero,
        Console: console,
    }
}

// Usage
console.Register[models.User](&UserConsole{})
// Type-safe, no reflection
```

## Pattern 6: Endpoint Registration

**Instead of**: String-based model names
**Use**: Type-safe resource registration

```go
// Resource registry
type ResourceRegistry struct {
    resources map[string]ResourceHandler
}

func RegisterResource[T any, Q any](
    name string,
    resource *Resource[T, Q],
) {
    registry.resources[name] = resource
}

// Usage
endpoints.RegisterResource("users", &UserResource{
    Resource: endpoints.NewResource[models.User](client),
})
```

## Pattern 7: Policy System

**Instead of**: String-based permission checks
**Use**: Type-safe policies

```go
// Policy interface
type Policy[T any] interface {
    CanView(user *models.User, obj T) bool
    CanEdit(user *models.User, obj T) bool
    CanDelete(user *models.User, obj T) bool
}

// Concrete policy
type PostPolicy struct {
    Policy[models.Post]
}

func (p *PostPolicy) CanView(user *models.User, post *models.Post) bool {
    return post.Published || post.AuthorID == user.ID
}

// Type-safe check
func Can[T any](user *models.User, action string, obj T) bool {
    policy := getPolicy[T]()
    switch action {
    case "view":
        return policy.CanView(user, obj)
    // ...
    }
}
```

## Benefits

1. **Compile-Time Safety**: All errors caught at compile time
2. **Performance**: Zero reflection overhead
3. **IDE Support**: Full autocomplete and type checking
4. **Maintainability**: Clear types, easier to understand
5. **Go Idioms**: Proper use of generics and interfaces

## Migration Strategy

1. Start with Repository pattern (easiest, biggest win)
2. Move to Resource handlers
3. Update Serializers
4. Refactor Console
5. Remove all reflection code
