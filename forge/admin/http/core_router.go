package http

import (
	"github.com/forgego/forge/admin"
	admintemplates "github.com/forgego/forge/admin/templates"
	"github.com/forgego/forge/server"
	httplib "github.com/forgego/forge/server"
	"github.com/gorilla/csrf"
)

// CoreRouter provides routing for the new core admin system
type CoreRouter struct {
	handler  *CoreHandler
	registry *admin.Registry
	csrfKey  []byte
}

// NewCoreRouter creates a new core admin router
func NewCoreRouter(registry *admin.Registry, templateDir string, sessionManager *server.SessionManager) *CoreRouter {
	// Create template engine and renderer
	engine := admintemplates.NewEngine(templateDir)
	renderer := admintemplates.NewRenderer(engine)

	// Load templates
	if err := renderer.LoadTemplates(templateDir); err != nil {
		// Log error but continue - templates can be loaded later
		_ = err
	}

	return &CoreRouter{
		handler:  NewCoreHandler(registry, renderer, sessionManager),
		registry: registry,
		csrfKey:  []byte("32-byte-long-auth-key-for-admin-"), // Default for dev
	}
}

// WithCSRFKey sets the CSRF key
func (r *CoreRouter) WithCSRFKey(key []byte) *CoreRouter {
	r.csrfKey = key
	return r
}

// RegisterRoutes registers all admin routes on the HTTP router
func (r *CoreRouter) RegisterRoutes(router *httplib.Router, path string) {
	// Create CSRF middleware
	CSRF := csrf.Protect(r.csrfKey, csrf.Path("/"))

	router.Route(path, func(subRouter *httplib.Router) {
		// Apply CSRF to admin only
		subRouter.Use(CSRF)

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
func (r *CoreRouter) registerModelRoutes(router *httplib.Router, modelName string) {
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

		// List editable (inline editing in list view)
		subRouter.Post("/list-editable/", r.handler.HandleListEditable(modelName))
	})
}
