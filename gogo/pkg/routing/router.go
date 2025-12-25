package routing

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// Router provides URL routing with reverse lookup
type Router struct {
	app     *fiber.App
	routes  map[string]*Route
	groups  []*Group
}

// Route represents a registered route
type Route struct {
	Name       string
	Method     string
	Path       string
	Handler    fiber.Handler
	Middleware []fiber.Handler
}

// Group represents a route group
type Group struct {
	Prefix     string
	Middleware []fiber.Handler
	Routes     []*Route
}

// NewRouter creates a new router
func NewRouter(app *fiber.App) *Router {
	return &Router{
		app:    app,
		routes: make(map[string]*Route),
		groups: make([]*Group, 0),
	}
}

// Get registers a GET route
func (r *Router) Get(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return r.register("GET", path, handler, options...)
}

// Post registers a POST route
func (r *Router) Post(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return r.register("POST", path, handler, options...)
}

// Put registers a PUT route
func (r *Router) Put(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return r.register("PUT", path, handler, options...)
}

// Delete registers a DELETE route
func (r *Router) Delete(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return r.register("DELETE", path, handler, options...)
}

// Patch registers a PATCH route
func (r *Router) Patch(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return r.register("PATCH", path, handler, options...)
}

// register registers a route
func (r *Router) register(method, path string, handler fiber.Handler, options ...RouteOption) *Route {
	route := &Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
	
	// Apply options
	for _, opt := range options {
		opt(route)
	}
	
	// Register with Fiber
	fiberRoute := r.app.Add(method, path, handler)
	if route.Name != "" {
		r.routes[route.Name] = route
	}
	
	// Apply middleware
	if len(route.Middleware) > 0 {
		for _, mw := range route.Middleware {
			fiberRoute.Use(mw)
		}
	}
	
	return route
}

// RouteOption configures a route
type RouteOption func(*Route)

// Name sets the route name for reverse lookup
func Name(name string) RouteOption {
	return func(r *Route) {
		r.Name = name
	}
}

// Middleware adds middleware to a route
func Middleware(handlers ...fiber.Handler) RouteOption {
	return func(r *Route) {
		r.Middleware = append(r.Middleware, handlers...)
	}
}

// Group creates a route group
func (r *Router) Group(prefix string, middleware ...fiber.Handler) *Group {
	group := &Group{
		Prefix:     prefix,
		Middleware: middleware,
		Routes:     make([]*Route, 0),
	}
	r.groups = append(r.groups, group)
	return group
}

// Get registers a GET route in a group
func (g *Group) Get(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return g.register("GET", path, handler, options...)
}

// Post registers a POST route in a group
func (g *Group) Post(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return g.register("POST", path, handler, options...)
}

// Put registers a PUT route in a group
func (g *Group) Put(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return g.register("PUT", path, handler, options...)
}

// Delete registers a DELETE route in a group
func (g *Group) Delete(path string, handler fiber.Handler, options ...RouteOption) *Route {
	return g.register("DELETE", path, handler, options...)
}

// register registers a route in a group
func (g *Group) register(method, path string, handler fiber.Handler, options ...RouteOption) *Route {
	fullPath := g.Prefix + path
	route := &Route{
		Method:  method,
		Path:    fullPath,
		Handler: handler,
	}
	
	// Apply options
	for _, opt := range options {
		opt(route)
	}
	
	g.Routes = append(g.Routes, route)
	return route
}

// URL generates a URL for a named route
func (r *Router) URL(name string, params ...Param) (string, error) {
	route, ok := r.routes[name]
	if !ok {
		return "", fmt.Errorf("route %s not found", name)
	}
	
	path := route.Path
	for _, param := range params {
		placeholder := ":" + param.Key
		path = strings.Replace(path, placeholder, param.Value, 1)
	}
	
	return path, nil
}

// Param represents a URL parameter
type Param struct {
	Key   string
	Value string
}

// Param creates a URL parameter
func Param(key, value string) Param {
	return Param{Key: key, Value: value}
}

// Reverse generates a URL (alias for URL)
func (r *Router) Reverse(name string, params ...Param) (string, error) {
	return r.URL(name, params...)
}

// ListRoutes returns all registered routes
func (r *Router) ListRoutes() []*Route {
	routes := make([]*Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}
	return routes
}

