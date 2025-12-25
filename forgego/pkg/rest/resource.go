package rest

import (
	"context"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/api"
)

// Resource defines the interface for API resource handlers
type Resource[T any, Q any] interface {
	// Index handles GET /resource - list all resources
	Index(ctx *Context) ([]*T, error)
	
	// Show handles GET /resource/:id - get a single resource
	Show(ctx *Context) (*T, error)
	
	// Create handles POST /resource - create a new resource
	Create(ctx *Context) (*T, error)
	
	// Update handles PUT /resource/:id - update a resource
	Update(ctx *Context) (*T, error)
	
	// Destroy handles DELETE /resource/:id - delete a resource
	Destroy(ctx *Context) error
}

// BaseResource provides a base implementation
type BaseResource[T any, Q any] struct {
	Repo       Repository[T, Q]
	Serializer api.Serializer[T]
}

// NewResource creates a new base resource
func NewResource[T any, Q any](repo Repository[T, Q]) *BaseResource[T, Q] {
	return &BaseResource[T, Q]{
		Repo:       repo,
		Serializer: api.NewModelSerializer[T](),
	}
}

// Repository interface for data access
type Repository[T any, Q any] interface {
	Query() Q
	GetByID(ctx context.Context, id interface{}) (*T, error)
	All(ctx context.Context, query Q) ([]*T, error)
	Count(ctx context.Context, query Q) (int, error)
	Create(ctx context.Context, data *T) (*T, error)
	Update(ctx context.Context, id interface{}, data *T) (*T, error)
	Delete(ctx context.Context, id interface{}) error
}

// Context provides request context for handlers
type Context struct {
	*fiber.Ctx
	User    interface{}
	Request *fiber.Request
}

// Param gets a path parameter
func (c *Context) Param(key string) string {
	return c.Ctx.Params(key)
}

// Query gets a query parameter
func (c *Context) Query(key string, defaultValue ...string) string {
	return c.Ctx.Query(key, defaultValue...)
}

// Bind binds request body to a struct
func (c *Context) Bind(dest interface{}) error {
	return c.Ctx.BodyParser(dest)
}

// JSON sends a JSON response
func (c *Context) JSON(data interface{}) error {
	return c.Ctx.JSON(data)
}

// Status sets the response status code
func (c *Context) Status(code int) *Context {
	c.Ctx.Status(code)
	return c
}

func (r *BaseResource[T, Q]) Index(ctx *Context) ([]*T, error) {
	query := r.Repo.Query()
	
	// Apply query parameters
	query = ApplyQuery(query, ctx.Ctx)
	
	// Parse pagination
	page, pageSize := ParsePagination(ctx.Ctx)
	
	// Apply pagination to query if supported
	if paginatedQuery, ok := any(query).(interface {
		Offset(int) interface{}
		Limit(int) interface{}
	}); ok {
		offset := (page - 1) * pageSize
		query = paginatedQuery.Offset(offset).(Q)
		query = paginatedQuery.Limit(pageSize).(Q)
	}
	
	results, err := r.Repo.All(ctx.Context(), query)
	if err != nil {
		return nil, err
	}
	
	return results, nil
}

func (r *BaseResource[T, Q]) Show(ctx *Context) (*T, error) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, NewError(400, "invalid_id", "Invalid ID format")
	}
	
	result, err := r.Repo.GetByID(ctx.Context(), idInt)
	if err != nil {
		return nil, &Error{Code: "not_found", Message: "Resource not found", Status: 404}
	}
	
	return result, nil
}

func (r *BaseResource[T, Q]) Create(ctx *Context) (*T, error) {
	var data T
	if err := ctx.Bind(&data); err != nil {
		return nil, NewError(400, "invalid_data", "Invalid request data")
	}
	
	result, err := r.Repo.Create(ctx.Context(), &data)
	if err != nil {
		return nil, &Error{Code: "create_failed", Message: err.Error(), Status: 500}
	}
	
	return result, nil
}

func (r *BaseResource[T, Q]) Update(ctx *Context) (*T, error) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, NewError(400, "invalid_id", "Invalid ID format")
	}
	
	var data T
	if err := ctx.Bind(&data); err != nil {
		return nil, NewError(400, "invalid_data", "Invalid request data")
	}
	
	result, err := r.Repo.Update(ctx.Context(), idInt, &data)
	if err != nil {
		return nil, &Error{Code: "update_failed", Message: err.Error(), Status: 500}
	}
	
	return result, nil
}

func (r *BaseResource[T, Q]) Destroy(ctx *Context) error {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &Error{Code: "invalid_id", Message: "Invalid ID format", Status: 400}
	}
	
	if err := r.Repo.Delete(ctx.Context(), idInt); err != nil {
		return &Error{Code: "delete_failed", Message: err.Error(), Status: 500}
	}
	
	return nil
}

