package admin

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Router handles route registration for the admin engine
type Router struct {
	registry *Registry
	app      *fiber.App
	basePath string
}

// NewRouter creates a new admin router
func NewRouter(registry *Registry, app *fiber.App, basePath string) *Router {
	return &Router{
		registry: registry,
		app:      app,
		basePath: strings.TrimSuffix(basePath, "/"),
	}
}

// Install installs all admin routes
func (r *Router) Install(client interface{}) error {
	// Install API routes
	if err := r.installAPIRoutes(client); err != nil {
		return err
	}

	// Install OpenAPI endpoint
	r.installOpenAPIRoute()

	// Install UI routes (will be implemented in Phase 3)
	// r.installUIRoutes()

	return nil
}

// installOpenAPIRoute installs the OpenAPI spec endpoint
func (r *Router) installOpenAPIRoute() {
	apiGroup := r.app.Group(r.basePath + "/api")
	
	apiGroup.Get("/openapi.json", func(c *fiber.Ctx) error {
		generator := NewOpenAPIGenerator(r.registry, r.basePath)
		spec, err := generator.Generate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate OpenAPI spec",
				"details": err.Error(),
			})
		}
		return c.JSON(spec)
	})
}

// installAPIRoutes installs REST API routes for all registered models
func (r *Router) installAPIRoutes(client interface{}) error {
	models := r.registry.GetAllModels()

	// Create API group
	apiGroup := r.app.Group(r.basePath + "/api")

	for modelName, meta := range models {
		// Convert model name to URL-friendly format (e.g., "User" -> "users")
		resourceName := r.modelNameToResource(modelName)

		// Create handler for this model
		handler := NewCRUDHandler(meta, client, r.registry)

		// Register routes
		apiGroup.Get("/"+resourceName, handler.List)                    // GET /admin/api/{model}
		apiGroup.Get("/"+resourceName+"/:id", handler.Get)             // GET /admin/api/{model}/{id}
		apiGroup.Post("/"+resourceName, handler.Create)                // POST /admin/api/{model}
		apiGroup.Put("/"+resourceName+"/:id", handler.Update)         // PUT /admin/api/{model}/{id}
		apiGroup.Delete("/"+resourceName+"/:id", handler.Delete)       // DELETE /admin/api/{model}/{id}
		
		// Bulk operations
		apiGroup.Post("/"+resourceName+"/bulk", handler.BulkCreate)    // POST /admin/api/{model}/bulk
		apiGroup.Put("/"+resourceName+"/bulk", handler.BulkUpdate)     // PUT /admin/api/{model}/bulk
		apiGroup.Delete("/"+resourceName+"/bulk", handler.BulkDelete)  // DELETE /admin/api/{model}/bulk
	}

	return nil
}

// modelNameToResource converts a model name to a URL-friendly resource name
func (r *Router) modelNameToResource(modelName string) string {
	// Convert PascalCase to snake_case and pluralize
	var result strings.Builder
	for i, r := range modelName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString("_")
		}
		result.WriteRune(r)
	}

	resourceName := strings.ToLower(result.String())

	// Simple pluralization
	if strings.HasSuffix(resourceName, "y") {
		resourceName = strings.TrimSuffix(resourceName, "y") + "ies"
	} else if strings.HasSuffix(resourceName, "s") || strings.HasSuffix(resourceName, "x") ||
		strings.HasSuffix(resourceName, "z") || strings.HasSuffix(resourceName, "ch") ||
		strings.HasSuffix(resourceName, "sh") {
		resourceName += "es"
	} else {
		resourceName += "s"
	}

	return resourceName
}

// GetModelNameFromResource converts a resource name back to model name
func (r *Router) GetModelNameFromResource(resourceName string) (string, error) {
	models := r.registry.GetAllModels()

	// Try to find matching model
	for modelName := range models {
		if r.modelNameToResource(modelName) == resourceName {
			return modelName, nil
		}
	}

	return "", fmt.Errorf("model not found for resource: %s", resourceName)
}

