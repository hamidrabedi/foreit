---
sidebar_position: 3
description: Build REST APIs with forge's API framework. Django REST Framework style but type-safe and fast.
keywords:
  - forge api
  - rest api go
  - django rest framework go
  - api framework
  - viewsets
image: /img/forge-social-card.jpg
---

# REST API Framework

The API framework gives you Django REST Framework functionality in Go - type-safe, fast, and production-ready. Build complete CRUD APIs with minimal code.

## Quick Start

Create a full CRUD API for a model:

```go
package api

import (
    "github.com/forgego/forge/api"
    "github.com/forgego/forge/api/serializers"
    "myapp/models"
)

type PostSerializer struct {
    serializers.BaseSerializer
}

func (s *PostSerializer) Fields() []serializers.Field {
    return []serializers.Field{
        serializers.Int64Field("id").ReadOnly().Build(),
        serializers.StringField("title").Required().Build(),
        serializers.TextField("content").Required().Build(),
        serializers.BooleanField("published").Build(),
        serializers.DateTimeField("created_at").ReadOnly().Build(),
        serializers.DateTimeField("updated_at").ReadOnly().Build(),
    }
}

func PostViewSet() *api.BaseViewSet {
    return api.NewBaseViewSet(
        &PostSerializer{},
        models.Post.Objects.All(),
        &models.Post{},
    )
}

// Register in your router
func RegisterRoutes(router *chi.Mux) {
    router.Route("/api/posts", func(r chi.Router) {
        r.Get("/", PostViewSet().List)
        r.Post("/", PostViewSet().Create)
        r.Get("/{id}", PostViewSet().Retrieve)
        r.Put("/{id}", PostViewSet().Update)
        r.Delete("/{id}", PostViewSet().Destroy)
    })
}
```

That's it - you now have a full REST API at `/api/posts/` with:
- GET `/api/posts/` - List all posts
- POST `/api/posts/` - Create a new post
- GET `/api/posts/123/` - Get a specific post
- PUT `/api/posts/123/` - Update a post
- DELETE `/api/posts/123/` - Delete a post

## Serializers

### Basic Serializer

```go
type UserSerializer struct {
    serializers.BaseSerializer
}

func (s *UserSerializer) Fields() []serializers.Field {
    return []serializers.Field{
        serializers.Int64Field("id").ReadOnly().Build(),
        serializers.StringField("username").Required().MaxLength(150).Build(),
        serializers.EmailField("email").Required().Build(),
        serializers.CharField("first_name").MaxLength(30).Build(),
        serializers.CharField("last_name").MaxLength(30).Build(),
        serializers.BooleanField("is_active").Default(true).Build(),
        serializers.DateTimeField("date_joined").ReadOnly().Build(),
    }
}
```

### Nested Serializers

```go
type PostSerializer struct {
    serializers.BaseSerializer
}

func (s *PostSerializer) Fields() []serializers.Field {
    return []serializers.Field{
        serializers.Int64Field("id").ReadOnly().Build(),
        serializers.StringField("title").Required().Build(),
        serializers.TextField("content").Required().Build(),
        serializers.ForeignKeyField("author", &UserSerializer{}).Build(),
        serializers.DateTimeField("created_at").ReadOnly().Build(),
    }
}
```

### Custom Methods

Add custom fields to your serializer:

```go
type PostSerializer struct {
    serializers.BaseSerializer
}

func (s *PostSerializer) Fields() []serializers.Field {
    return []serializers.Field{
        serializers.Int64Field("id").ReadOnly().Build(),
        serializers.StringField("title").Required().Build(),
        serializers.TextField("content").Required().Build(),
        serializers.MethodField("word_count", s.getWordCount).Build(),
        serializers.DateTimeField("created_at").ReadOnly().Build(),
    }
}

func (s *PostSerializer) getWordCount(obj interface{}) (interface{}, error) {
    post := obj.(*models.Post)
    words := strings.Split(post.Content, " ")
    return len(words), nil
}
```

### Validation

Add custom validation:

```go
func (s *UserSerializer) Validate(data map[string]interface{}) error {
    password := data["password"].(string)
    confirm := data["password_confirm"].(string)
    
    if password != confirm {
        return errors.New("Passwords don't match")
    }
    
    if len(password) < 8 {
        return errors.New("Password must be at least 8 characters")
    }
    
    return nil
}

func (s *UserSerializer) ValidateEmail(value string) error {
    if !strings.Contains(value, "@") {
        return errors.New("Invalid email format")
    }
    return nil
}
```

## ViewSets

### BaseViewSet

The `BaseViewSet` gives you CRUD operations out of the box:

```go
type PostViewSet struct {
    *api.BaseViewSet
}

func NewPostViewSet() *PostViewSet {
    return &PostViewSet{
        BaseViewSet: api.NewBaseViewSet(
            &PostSerializer{},
            models.Post.Objects.All(),
            &models.Post{},
        ),
    }
}
```

### Custom ViewSet

Override methods for custom behavior:

```go
type PostViewSet struct {
    *api.BaseViewSet
}

func (v *PostViewSet) List(w http.ResponseWriter, r *http.Request) {
    // Custom list logic
    queryset := v.GetQueryset().Filter(models.Post.Fields.Published.Equals(true))
    
    // Add pagination
    page := r.URL.Query().Get("page")
    if page == "" {
        page = "1"
    }
    
    paginated := api.Paginate(queryset, page, 20)
    v.RenderResponse(w, paginated)
}

func (v *PostViewSet) Create(w http.ResponseWriter, r *http.Request) {
    // Custom create logic
    var data map[string]interface{}
    json.NewDecoder(r.Body).Decode(&data)
    
    // Set author from current user
    user := getCurrentUser(r)
    data["author_id"] = user.ID
    
    serializer := &PostSerializer{}
    if err := serializer.Validate(data); err != nil {
        v.RenderError(w, err, http.StatusBadRequest)
        return
    }
    
    post := serializer.Create(data)
    v.RenderResponse(w, post, http.StatusCreated)
}
```

### Custom Actions

Add custom actions to your ViewSet:

```go
func (v *PostViewSet) RegisterRoutes(router chi.Router) {
    // Standard CRUD routes
    router.Get("/", v.List)
    router.Post("/", v.Create)
    router.Get("/{id}", v.Retrieve)
    router.Put("/{id}", v.Update)
    router.Delete("/{id}", v.Destroy)
    
    // Custom actions
    router.Post("/{id}/publish", v.Publish)
    router.Get("/{id}/stats", v.Stats)
    router.Post("/{id}/like", v.Like)
}

func (v *PostViewSet) Publish(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    post, err := models.Post.Objects.Get(context.Background(), 
        models.Post.Fields.ID.Equals(parseID(id)))
    if err != nil {
        v.RenderError(w, err, http.StatusNotFound)
        return
    }
    
    post.Published = true
    models.Post.Objects.Save(post)
    
    v.RenderResponse(w, post)
}
```

## Authentication

### Token Authentication

```go
func setupAuth() {
    api.UseAuthentication(&api.TokenAuthentication{
        Keyword: "Bearer",
    })
}
```

### JWT Authentication

```go
func setupAuth() {
    api.UseAuthentication(&api.JWTAuthentication{
        SecretKey: "your-secret-key",
        Algorithm: "HS256",
        Expiration: time.Hour * 24,
    })
}
```

### Multiple Authentication

```go
func setupAuth() {
    api.UseAuthentication(
        &api.JWTAuthentication{...},
        &api.TokenAuthentication{...},
        &api.SessionAuthentication{...},
    )
}
```

## Permissions

### Built-in Permissions

```go
func PostViewSet() *api.BaseViewSet {
    return api.NewBaseViewSet(
        &PostSerializer{},
        models.Post.Objects.All(),
        &models.Post{},
    ).WithPermissions(
        &api.IsAuthenticated{},
        &api.IsAdminUser{},
    )
}
```

### Custom Permissions

```go
type IsAuthorOrReadOnly struct{}

func (p *IsAuthorOrReadOnly) HasPermission(request *http.Request, view interface{}, obj interface{}) bool {
    if request.Method == "GET" || request.Method == "HEAD" || request.Method == "OPTIONS" {
        return true
    }
    
    user := getCurrentUser(request)
    post := obj.(*models.Post)
    
    return post.AuthorID == user.ID || user.IsAdmin
}

func PostViewSet() *api.BaseViewSet {
    return api.NewBaseViewSet(
        &PostSerializer{},
        models.Post.Objects.All(),
        &models.Post{},
    ).WithPermissions(
        &api.IsAuthenticated{},
        &IsAuthorOrReadOnly{},
    )
}
```

## Throttling

### Rate Limiting

```go
func setupThrottling() {
    api.UseThrottling(
        &api.AnonRateThrottle{Rate: "100/hour"},
        &api.UserRateThrottle{Rate: "1000/hour"},
    )
}
```

### Scoped Throttling

```go
type UploadThrottle struct {
    api.ScopedRateThrottle
}

func (t *UploadThrottle) GetRate(request *http.Request) string {
    if request.URL.Path == "/api/upload" {
        return "10/hour"
    }
    return "1000/hour"
}
```

## Renderers

### Multiple Renderers

```go
func PostViewSet() *api.BaseViewSet {
    return api.NewBaseViewSet(
        &PostSerializer{},
        models.Post.Objects.All(),
        &models.Post{},
    ).WithRenderers(
        &api.JSONRenderer{},
        &api.XMLRenderer{},
        &api.CSVRenderer{},
    )
}
```

### Custom Renderer

```go
type PDFRenderer struct{}

func (r *PDFRenderer) Render(data interface{}) ([]byte, error) {
    // Convert data to PDF
    return generatePDF(data), nil
}

func (r *PDFRenderer) GetContentType() string {
    return "application/pdf"
}
```

## Filtering

### Basic Filtering

```go
func PostViewSet() *api.BaseViewSet {
    return api.NewBaseViewSet(
        &PostSerializer{},
        models.Post.Objects.All(),
        &models.Post{},
    ).WithFilters(
        "title", "author", "published", "created_at",
    )
}
```

### Custom Filtering

```go
type PostFilter struct {
    api.BaseFilter
}

func (f *PostFilter) Filter(queryset interface{}, value interface{}) interface{} {
    qs := queryset.(models.PostQuerySet)
    
    if published, ok := value.(bool); ok {
        return qs.Filter(models.Post.Fields.Published.Equals(published))
    }
    
    return qs
}

func setupFilters() {
    api.RegisterFilter("published", &PostFilter{})
}
```

## Pagination

### Page Number Pagination

```go
func setupPagination() {
    api.UsePagination(&api.PageNumberPagination{
        PageSize: 20,
        MaxPageSize: 100,
    })
}
```

### Limit Offset Pagination

```go
func setupPagination() {
    api.UsePagination(&api.LimitOffsetPagination{
        DefaultLimit: 20,
        MaxLimit: 100,
    })
}
```

## Versioning

### URL Versioning

```go
func setupVersioning() {
    api.UseVersioning(&api.URLPathVersioning{
        DefaultVersion: "v1",
        AllowedVersions: []string{"v1", "v2"},
    })
}

// Routes become:
// /api/v1/posts/
// /api/v2/posts/
```

### Header Versioning

```go
func setupVersioning() {
    api.UseVersioning(&api.AcceptHeaderVersioning{
        DefaultVersion: "v1",
        AllowedVersions: []string{"v1", "v2"},
    })
}
```

## OpenAPI Documentation

### Auto-Generated Docs

```go
func setupDocs() {
    api.SetupOpenAPI(&api.OpenAPIConfig{
        Title: "My API",
        Version: "1.0.0",
        Description: "API documentation",
        Contact: api.Contact{
            Name: "API Support",
            Email: "support@example.com",
        },
    })
}
```

### Custom Documentation

```go
func (v *PostViewSet) GetDocs() *api.Operation {
    return &api.Operation{
        Summary: "List posts",
        Description: "Get a list of all posts",
        Parameters: []api.Parameter{
            {
                Name: "page",
                In: "query",
                Description: "Page number",
                Required: false,
                Type: "integer",
            },
        },
        Responses: map[int]api.Response{
            200: {
                Description: "List of posts",
                Schema: &api.Schema{Ref: "#/components/schemas/PostList"},
            },
        },
    }
}
```

## Error Handling

### Custom Error Responses

```go
func (v *PostViewSet) HandleError(w http.ResponseWriter, err error, status int) {
    if errors.Is(err, models.ErrNotFound) {
        v.RenderError(w, map[string]string{
            "error": "Post not found",
            "code": "POST_NOT_FOUND",
        }, http.StatusNotFound)
        return
    }
    
    // Default error handling
    v.BaseViewSet.HandleError(w, err, status)
}
```

## Testing

### Test Your APIs

```go
func TestPostViewSet(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer db.Close()
    
    // Create test data
    user := &models.User{Username: "testuser"}
    models.User.Objects.Create(user)
    
    post := &models.Post{
        Title: "Test Post",
        Content: "Test content",
        AuthorID: user.ID,
    }
    models.Post.Objects.Create(post)
    
    // Test list endpoint
    req := httptest.NewRequest("GET", "/api/posts/", nil)
    w := httptest.NewRecorder()
    
    viewset := NewPostViewSet()
    viewset.List(w, req)
    
    resp := w.Result()
    assert.Equal(t, 200, resp.StatusCode)
    
    var posts []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&posts)
    assert.Equal(t, 1, len(posts))
    assert.Equal(t, "Test Post", posts[0]["title"])
}
```

## Best Practices

1. **Keep ViewSets small** - Split complex logic into services
2. **Use proper HTTP methods** - GET for reading, POST for creating, etc.
3. **Validate everything** - Never trust user input
4. **Handle errors gracefully** - Return meaningful error messages
5. **Document your APIs** - Use OpenAPI for documentation
6. **Test thoroughly** - Write tests for all endpoints
7. **Use caching** - Cache expensive operations
8. **Monitor performance** - Track API response times

## Next Steps

- [Authentication System](/docs/features/identity-system) - Complete user management
- [Filter System](/docs/features/filter-system) - Advanced filtering
- [Examples](/docs/examples/api) - Complete API examples
