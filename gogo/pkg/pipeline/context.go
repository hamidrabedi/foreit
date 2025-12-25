package pipeline

import (
	"github.com/gofiber/fiber/v2"
)

// Context provides request context through the middleware pipeline
type Context struct {
	*fiber.Ctx
	values map[string]interface{}
}

// NewContext creates a new pipeline context
func NewContext(c *fiber.Ctx) *Context {
	return &Context{
		Ctx:    c,
		values: make(map[string]interface{}),
	}
}

// Set stores a value in the context
func (c *Context) Set(key string, value interface{}) {
	if c.values == nil {
		c.values = make(map[string]interface{})
	}
	c.values[key] = value
}

// Get retrieves a value from the context
func (c *Context) Get(key string) (interface{}, bool) {
	if c.values == nil {
		return nil, false
	}
	value, ok := c.values[key]
	return value, ok
}

// GetString retrieves a string value
func (c *Context) GetString(key string, defaultValue string) string {
	if value, ok := c.Get(key); ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt retrieves an int value
func (c *Context) GetInt(key string, defaultValue int) int {
	if value, ok := c.Get(key); ok {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return defaultValue
}

// Handler is a middleware handler function
type Handler func(*Context) error

// Middleware is a function that wraps a handler
type Middleware func(Handler) Handler

