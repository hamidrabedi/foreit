package admin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	registry            *Registry
	app                 *fiber.App
	basePath            string
	serviceIntegration  *ServiceIntegration
	viewSetIntegration  *ViewSetIntegration
}

func NewRouter(registry *Registry, app *fiber.App, basePath string) *Router {
	return &Router{
		registry: registry,
		app:      app,
		basePath: strings.TrimSuffix(basePath, "/"),
	}
}

func (r *Router) WithServiceIntegration(integration *ServiceIntegration) *Router {
	r.serviceIntegration = integration
	return r
}

func (r *Router) WithViewSetIntegration(integration *ViewSetIntegration) *Router {
	r.viewSetIntegration = integration
	return r
}

func (r *Router) Install(client interface{}) error {
	if err := r.installAPIRoutes(client); err != nil {
		return err
	}

	if err := r.installUIRoutes(); err != nil {
		return err
	}

	r.installOpenAPIRoute()
	return nil
}

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

func (r *Router) installAPIRoutes(client interface{}) error {
	if r.serviceIntegration == nil && r.viewSetIntegration == nil {
		return nil
	}
	
	models := r.registry.GetAllModels()
	apiGroup := r.app.Group(r.basePath + "/api")
	
	for modelName, meta := range models {
		resourceName := r.modelNameToResource(modelName)
		basePath := fmt.Sprintf("/%s", resourceName)
		
		// Try ViewSet integration first
		if r.viewSetIntegration != nil {
			if err := r.installViewSetAPIRoutes(apiGroup, basePath, modelName, meta); err == nil {
				continue
			}
		}
		
		// Fall back to Service integration
		if r.serviceIntegration != nil {
			if err := r.installServiceAPIRoutes(apiGroup, basePath, modelName, meta); err != nil {
				return fmt.Errorf("failed to install routes for %s: %w", modelName, err)
			}
		}
	}
	
	return nil
}

func (r *Router) installViewSetAPIRoutes(apiGroup fiber.Router, basePath, modelName string, meta *ModelMeta) error {
	// This will be implemented with type-safe handlers
	// For now, we'll use the UI routes
	return nil
}

func (r *Router) installServiceAPIRoutes(apiGroup fiber.Router, basePath, modelName string, meta *ModelMeta) error {
	// This will be implemented with service handlers
	return nil
}

func (r *Router) modelNameToResource(modelName string) string {
	var result strings.Builder
	for i, r := range modelName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString("_")
		}
		result.WriteRune(r)
	}

	resourceName := strings.ToLower(result.String())

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

func (r *Router) GetModelNameFromResource(resourceName string) (string, error) {
	models := r.registry.GetAllModels()

	for modelName := range models {
		if r.modelNameToResource(modelName) == resourceName {
			return modelName, nil
		}
	}

	return "", fmt.Errorf("model not found for resource: %s", resourceName)
}

func (r *Router) installUIRoutes() error {
	if r.viewSetIntegration == nil && r.serviceIntegration == nil {
		return nil
	}
	
	models := r.registry.GetAllModels()
	uiGroup := r.app.Group(r.basePath)
	
	// Install index page
	uiGroup.Get("/", r.handleIndex)
	
	for modelName, meta := range models {
		resourceName := r.modelNameToResource(modelName)
		basePath := fmt.Sprintf("/%s", resourceName)
		
		// Try ViewSet integration first
		if r.viewSetIntegration != nil {
			if err := r.installViewSetUIRoutes(uiGroup, basePath, modelName, meta); err == nil {
				continue
			}
		}
		
		// Fall back to Service integration
		if r.serviceIntegration != nil {
			if err := r.installServiceUIRoutes(uiGroup, basePath, modelName, meta); err != nil {
				return fmt.Errorf("failed to install UI routes for %s: %w", modelName, err)
			}
		}
	}
	
	return nil
}

func (r *Router) handleIndex(c *fiber.Ctx) error {
	models := r.registry.GetAllModels()
	modelList := make([]map[string]interface{}, 0, len(models))
	
	for name := range models {
		resourceName := r.modelNameToResource(name)
		modelList = append(modelList, map[string]interface{}{
			"name":   name,
			"label":  name,
			"url":    fmt.Sprintf("%s/%s/", r.basePath, resourceName),
		})
	}
	
	return c.JSON(fiber.Map{
		"models": modelList,
	})
}

func (r *Router) installViewSetUIRoutes(uiGroup fiber.Router, basePath, modelName string, meta *ModelMeta) error {
	// Create handler using reflection
	handler := r.createViewSetUIHandler(modelName)
	if handler == nil {
		return fmt.Errorf("failed to create handler for %s", modelName)
	}
	
	// Use reflection to call methods
	handlerValue := reflect.ValueOf(handler)
	
	// List
	listMethod := handlerValue.MethodByName("ListUI")
	if listMethod.IsValid() {
		uiGroup.Get(basePath+"/", func(c *fiber.Ctx) error {
			results := listMethod.Call([]reflect.Value{reflect.ValueOf(c)})
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			return nil
		})
	}
	
	// Add
	formMethod := handlerValue.MethodByName("FormUI")
	if formMethod.IsValid() {
		uiGroup.Get(basePath+"/add/", func(c *fiber.Ctx) error {
			results := formMethod.Call([]reflect.Value{reflect.ValueOf(c)})
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			return nil
		})
	}
	
	// Change
	uiGroup.Get(basePath+"/:id/change/", func(c *fiber.Ctx) error {
		results := formMethod.Call([]reflect.Value{reflect.ValueOf(c)})
		if len(results) > 0 {
			if err, ok := results[0].Interface().(error); ok && err != nil {
				return err
			}
		}
		return nil
	})
	
	// Detail
	detailMethod := handlerValue.MethodByName("DetailUI")
	if detailMethod.IsValid() {
		uiGroup.Get(basePath+"/:id/", func(c *fiber.Ctx) error {
			results := detailMethod.Call([]reflect.Value{reflect.ValueOf(c)})
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			return nil
		})
	}
	
	return nil
}

func (r *Router) installServiceUIRoutes(uiGroup fiber.Router, basePath, modelName string, meta *ModelMeta) error {
	handler := r.createServiceUIHandler(modelName)
	if handler == nil {
		return fmt.Errorf("failed to create service handler for %s", modelName)
	}
	
	handlerValue := reflect.ValueOf(handler)
	
	// List
	listMethod := handlerValue.MethodByName("ListUI")
	if listMethod.IsValid() {
		uiGroup.Get(basePath+"/", func(c *fiber.Ctx) error {
			results := listMethod.Call([]reflect.Value{reflect.ValueOf(c)})
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			return nil
		})
	}
	
	// Add
	formMethod := handlerValue.MethodByName("FormUI")
	if formMethod.IsValid() {
		uiGroup.Get(basePath+"/add/", func(c *fiber.Ctx) error {
			results := formMethod.Call([]reflect.Value{reflect.ValueOf(c)})
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			return nil
		})
	}
	
	// Change
	uiGroup.Get(basePath+"/:id/change/", func(c *fiber.Ctx) error {
		results := formMethod.Call([]reflect.Value{reflect.ValueOf(c)})
		if len(results) > 0 {
			if err, ok := results[0].Interface().(error); ok && err != nil {
				return err
			}
		}
		return nil
	})
	
	// Detail
	detailMethod := handlerValue.MethodByName("DetailUI")
	if detailMethod.IsValid() {
		uiGroup.Get(basePath+"/:id/", func(c *fiber.Ctx) error {
			results := detailMethod.Call([]reflect.Value{reflect.ValueOf(c)})
			if len(results) > 0 {
				if err, ok := results[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			return nil
		})
	}
	
	return nil
}

func (r *Router) createViewSetUIHandler(modelName string) interface{} {
	if r.viewSetIntegration == nil {
		return nil
	}
	
	return &genericViewSetHandler{
		integration: r.viewSetIntegration,
		modelName:   modelName,
		registry:    r.registry,
	}
}

func (r *Router) createServiceUIHandler(modelName string) interface{} {
	if r.serviceIntegration == nil {
		return nil
	}
	
	return &genericServiceHandler{
		integration: r.serviceIntegration,
		modelName:   modelName,
		registry:    r.registry,
	}
}

type genericViewSetHandler struct {
	integration *ViewSetIntegration
	modelName   string
	registry    *Registry
}

func (h *genericViewSetHandler) ListUI(c *fiber.Ctx) error {
	_, err := h.registry.GetModel(h.modelName)
	if err != nil {
		return c.Status(404).SendString("Model not found")
	}
	
	// Use reflection to call the typed handler
	// For now, return a simple response
	return c.SendString(fmt.Sprintf("List UI for %s", h.modelName))
}

func (h *genericViewSetHandler) FormUI(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Form UI for %s", h.modelName))
}

func (h *genericViewSetHandler) DetailUI(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Detail UI for %s", h.modelName))
}

func (h *genericViewSetHandler) handleCreate(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Create handler for %s", h.modelName))
}

func (h *genericViewSetHandler) handleUpdate(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Update handler for %s", h.modelName))
}

func (h *genericViewSetHandler) handleDelete(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Delete handler for %s", h.modelName))
}

type genericServiceHandler struct {
	integration *ServiceIntegration
	modelName   string
	registry    *Registry
}

func (h *genericServiceHandler) ListUI(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Service List UI for %s", h.modelName))
}

func (h *genericServiceHandler) FormUI(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Service Form UI for %s", h.modelName))
}

func (h *genericServiceHandler) DetailUI(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Service Detail UI for %s", h.modelName))
}

