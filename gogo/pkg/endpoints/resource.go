package endpoints

import (
	"context"
	"github.com/gofiber/fiber/v2"
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
	Serializer Serializer[T]
}

// NewResource creates a new base resource
func NewResource[T any, Q any](repo Repository[T, Q]) *BaseResource[T, Q] {
	return &BaseResource[T, Q]{
		Repo:       repo,
		Serializer: NewModelSerializer[T](),
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

