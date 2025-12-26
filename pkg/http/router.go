package http

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

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}
