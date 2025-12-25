package gogo

import (
	"github.com/gogo/pkg/routing"
	"github.com/gofiber/fiber/v2"
)

// Router returns a helper to access router methods
// This is a convenience method
func (a *App) RouterHelper() *RouterHelper {
	return &RouterHelper{app: a}
}

// RouterHelper provides helper methods for routing
type RouterHelper struct {
	app *App
}

// Name creates a Name route option
func (h *RouterHelper) Name(name string) routing.RouteOption {
	return routing.Name(name)
}

// Middleware creates a Middleware route option
func (h *RouterHelper) Middleware(handlers ...fiber.Handler) routing.RouteOption {
	return routing.Middleware(handlers...)
}

