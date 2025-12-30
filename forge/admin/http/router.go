package http

import (
	"github.com/forgego/forge/admin"
	httplib "github.com/forgego/forge/server"
)

// Router provides routing for admin views
type Router struct {
	handler  *Handler
	registry *admin.Registry
}

// NewRouter creates a new admin router
func NewRouter(registry *admin.Registry) *Router {
	return &Router{
		handler:  NewHandler(registry),
		registry: registry,
	}
}

// RegisterRoutes registers all admin routes on the HTTP router
func (r *Router) RegisterRoutes(router *httplib.Router, path string) {
	router.Route(path, func(subRouter *httplib.Router) {
		// Admin index
		subRouter.Get("/", r.handler.HandleIndex())

		// Register routes for each model
		allAdmins := r.registry.GetAll()
		for modelName := range allAdmins {
			r.registerModelRoutes(subRouter, modelName)
		}
	})
}

// registerModelRoutes registers routes for a specific model
func (r *Router) registerModelRoutes(router *httplib.Router, modelName string) {
	router.Route("/"+modelName, func(subRouter *httplib.Router) {
		// List view
		subRouter.Get("/", r.handler.HandleList(modelName))

		// Export
		subRouter.Get("/export/", r.handler.HandleExport(modelName))

		// Bulk actions
		subRouter.Post("/bulk-action/", r.handler.HandleBulkAction(modelName))

		// Create
		subRouter.Get("/new/", r.handler.HandleCreate(modelName))
		subRouter.Post("/new/", r.handler.HandleCreate(modelName))

		// Detail
		subRouter.Get("/{id}/", r.handler.HandleDetail(modelName))

		// Update
		subRouter.Get("/{id}/change/", r.handler.HandleUpdate(modelName))
		subRouter.Post("/{id}/change/", r.handler.HandleUpdate(modelName))

		// Delete
		subRouter.Delete("/{id}/delete/", r.handler.HandleDelete(modelName))
		subRouter.Post("/{id}/delete/", r.handler.HandleDelete(modelName))

		// Autocomplete
		subRouter.Get("/autocomplete/", r.handler.HandleAutocomplete(modelName))
	})
}
