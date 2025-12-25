package admin

import (
	"github.com/gofiber/fiber/v2"
)

type Engine struct {
	registry *Registry
	client   interface{}
}

func New(client interface{}) *Engine {
	registry := NewRegistry()
	registry.SetClient(client)
	
	return &Engine{
		registry: registry,
		client:   client,
	}
}

func (e *Engine) Register(model interface{}, options ...Option) error {
	return e.registry.Register(model, options...)
}

func (e *Engine) Install(app *fiber.App, basePath string) error {
	router := NewRouter(e.registry, app, basePath)
	return router.Install(e.client)
}

func (e *Engine) GetRegistry() *Registry {
	return e.registry
}

func (e *Engine) GetModel(name string) (*ModelMeta, error) {
	return e.registry.GetModel(name)
}

func (e *Engine) GetAllModels() map[string]*ModelMeta {
	return e.registry.GetAllModels()
}

func (e *Engine) GenerateOpenAPI(baseURL string) (*OpenAPISpec, error) {
	generator := NewOpenAPIGenerator(e.registry, baseURL)
	return generator.Generate()
}

func (e *Engine) GenerateOpenAPIJSON(baseURL string) ([]byte, error) {
	generator := NewOpenAPIGenerator(e.registry, baseURL)
	return generator.GenerateJSON()
}

func (e *Engine) RegisterHook(modelName string, hookType HookType, hook HookFunc) {
	registry := e.registry.GetHookRegistry()
	registry.Register(modelName, hookType, hook)
}

