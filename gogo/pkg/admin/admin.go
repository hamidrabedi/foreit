package admin

import (
	"entgo.io/ent"

	"github.com/gogo/internal/admin"
	"github.com/gofiber/fiber/v2"
)

// Engine is the main admin engine that manages the registry and provides the public API
type Engine struct {
	registry *admin.Registry
	client   interface{} // *ent.Client
}

// New creates a new admin engine
func New(client interface{}) *Engine {
	registry := admin.NewRegistry()
	registry.SetClient(client)
	
	return &Engine{
		registry: registry,
		client:   client,
	}
}

// Register registers a model with the admin engine
// This is the main public API for registering models
func (e *Engine) Register(model interface{}, options ...admin.Option) error {
	// Register the model
	if err := e.registry.Register(model, options...); err != nil {
		return err
	}

	// If the model implements ent.Schema, introspect it
	if schema, ok := model.(ent.Schema); ok {
		introspector := admin.NewSchemaIntrospector(e.registry)
		_, err := introspector.Introspect(schema)
		return err
	}

	return nil
}

// Install installs the admin engine on a Fiber app
func (e *Engine) Install(app *fiber.App, basePath string) error {
	router := admin.NewRouter(e.registry, app, basePath)
	return router.Install(e.client)
}

// GetRegistry returns the internal registry (for advanced usage)
func (e *Engine) GetRegistry() *admin.Registry {
	return e.registry
}

// GetModel retrieves model metadata by name
func (e *Engine) GetModel(name string) (*admin.ModelMeta, error) {
	return e.registry.GetModel(name)
}

// GetAllModels returns all registered models
func (e *Engine) GetAllModels() map[string]*admin.ModelMeta {
	return e.registry.GetAllModels()
}

// GenerateOpenAPI generates an OpenAPI 3.0 specification
func (e *Engine) GenerateOpenAPI(baseURL string) (*admin.OpenAPISpec, error) {
	generator := admin.NewOpenAPIGenerator(e.registry, baseURL)
	return generator.Generate()
}

// GenerateOpenAPIJSON generates OpenAPI spec as JSON
func (e *Engine) GenerateOpenAPIJSON(baseURL string) ([]byte, error) {
	generator := admin.NewOpenAPIGenerator(e.registry, baseURL)
	return generator.GenerateJSON()
}

// RegisterHook registers a hook for a model
func (e *Engine) RegisterHook(modelName string, hookType admin.HookType, hook admin.HookFunc) {
	registry := e.registry.GetHookRegistry()
	registry.Register(modelName, hookType, hook)
}

