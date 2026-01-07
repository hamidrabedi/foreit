package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// URLParam gets a URL parameter from chi router
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// Router wraps chi.Router with framework-specific extensions
type Router struct {
	chi.Router
	// middleware is reserved for future use
	// nolint:unused // Reserved for future middleware management
	middleware []func(http.Handler) http.Handler
	notFound   http.HandlerFunc
	methodNotAllowed http.HandlerFunc
}

// NewRouter creates a new framework router with default middleware
func NewRouter() *Router {
	r := chi.NewRouter()

	// Add default chi middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)

	return &Router{
		Router: r,
		middleware: []func(http.Handler) http.Handler{},
	}
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...func(http.Handler) http.Handler) {
	r.Router.Use(middleware...)
}

// Route creates a new sub-router with a pattern
func (r *Router) Route(pattern string, fn func(*Router)) {
	r.Router.Route(pattern, func(c chi.Router) {
		subRouter := &Router{Router: c}
		fn(subRouter)
	})
}

// Mount mounts another router at a pattern
func (r *Router) Mount(pattern string, handler http.Handler) {
	r.Router.Mount(pattern, handler)
}

// Get adds a GET route
func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.Router.Get(pattern, handler)
}

// Post adds a POST route
func (r *Router) Post(pattern string, handler http.HandlerFunc) {
	r.Router.Post(pattern, handler)
}

// Put adds a PUT route
func (r *Router) Put(pattern string, handler http.HandlerFunc) {
	r.Router.Put(pattern, handler)
}

// Patch adds a PATCH route
func (r *Router) Patch(pattern string, handler http.HandlerFunc) {
	r.Router.Patch(pattern, handler)
}

// Delete adds a DELETE route
func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
	r.Router.Delete(pattern, handler)
}

// Head adds a HEAD route
func (r *Router) Head(pattern string, handler http.HandlerFunc) {
	r.Router.Head(pattern, handler)
}

// Options adds an OPTIONS route
func (r *Router) Options(pattern string, handler http.HandlerFunc) {
	r.Router.Options(pattern, handler)
}

// Connect adds a CONNECT route
func (r *Router) Connect(pattern string, handler http.HandlerFunc) {
	r.Router.Connect(pattern, handler)
}

// Trace adds a TRACE route
func (r *Router) Trace(pattern string, handler http.HandlerFunc) {
	r.Router.Trace(pattern, handler)
}

// With applies middleware to a route group
func (r *Router) With(middleware ...func(http.Handler) http.Handler) *Router {
	subRouter := &Router{
		Router:     r.Router.With(middleware...).(chi.Router),
		middleware: append(r.middleware, middleware...),
		notFound:   r.notFound,
		methodNotAllowed: r.methodNotAllowed,
	}
	return subRouter
}

// Group creates a new route group with middleware
func (r *Router) Group(fn func(*Router)) {
	r.Route("", fn)
}

// NotFound sets a custom 404 handler
func (r *Router) NotFound(handler http.HandlerFunc) {
	r.notFound = handler
	r.Router.NotFound(handler)
}

// MethodNotAllowed sets a custom 405 handler
func (r *Router) MethodNotAllowed(handler http.HandlerFunc) {
	r.methodNotAllowed = handler
	r.Router.MethodNotAllowed(handler)
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

