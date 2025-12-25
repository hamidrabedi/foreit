package admin

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/service"
)

type ServiceAdapter[T any] struct {
	service   service.ResourceService[T]
	modelMeta *ModelMeta
}

func NewServiceAdapter[T any](svc service.ResourceService[T], modelMeta *ModelMeta) *ServiceAdapter[T] {
	return &ServiceAdapter[T]{
		service:   svc,
		modelMeta: modelMeta,
	}
}

func (a *ServiceAdapter[T]) GetService() service.ResourceService[T] {
	return a.service
}

func (a *ServiceAdapter[T]) GetModelMeta() *ModelMeta {
	return a.modelMeta
}

type ServiceIntegration struct {
	registry    *Registry
	rendererRegistry *RendererRegistry
	templates   *TemplateEngine
	db          interface{}
	adapters    map[string]interface{}
}

func NewServiceIntegration(registry *Registry, rendererRegistry *RendererRegistry,
	templates *TemplateEngine, db interface{}) *ServiceIntegration {
	return &ServiceIntegration{
		registry:    registry,
		rendererRegistry: rendererRegistry,
		templates:   templates,
		db:          db,
		adapters:    make(map[string]interface{}),
	}
}

func RegisterService[T any](integration *ServiceIntegration,
	model interface{},
	svc service.ResourceService[T],
	options ...Option) error {

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	modelName := modelType.Name()

	err := integration.registry.Register(model, options...)
	if err != nil {
		return fmt.Errorf("failed to register model: %w", err)
	}

	meta, err := integration.registry.GetModel(modelName)
	if err != nil {
		return fmt.Errorf("failed to get model metadata: %w", err)
	}

	adapter := NewServiceAdapter(svc, meta)
	integration.adapters[modelName] = adapter

	return nil
}

func (i *ServiceIntegration) GetAdapter(modelName string) (interface{}, bool) {
	adapter, ok := i.adapters[modelName]
	return adapter, ok
}

func (i *ServiceIntegration) GetRegistry() *Registry {
	return i.registry
}

func (i *ServiceIntegration) GetRendererRegistry() *RendererRegistry {
	return i.rendererRegistry
}

func (i *ServiceIntegration) GetTemplates() *TemplateEngine {
	return i.templates
}

func (i *ServiceIntegration) GetDB() interface{} {
	return i.db
}

func (i *ServiceIntegration) CreateViewSetIntegration() *ViewSetIntegration {
	db, _ := i.db.(*models.DB)
	return NewViewSetIntegration(
		i.registry,
		i.rendererRegistry,
		i.templates,
		db,
	)
}

func CreateServiceUIHandler[T any](integration *ServiceIntegration, modelName string) (*ServiceUIHandler[T], error) {
	meta, err := integration.registry.GetModel(modelName)
	if err != nil {
		return nil, err
	}

	adapterInterface, ok := integration.GetAdapter(modelName)
	if !ok {
		return nil, fmt.Errorf("no service registered for model %s", modelName)
	}

	adapter, ok := adapterInterface.(*ServiceAdapter[T])
	if !ok {
		return nil, fmt.Errorf("invalid adapter type for model %s", modelName)
	}

	svc := adapter.GetService()
	return NewServiceUIHandler(meta, svc, integration.registry,
		integration.rendererRegistry, integration.templates), nil
}

type ServiceUIHandler[T any] struct {
	modelMeta   *ModelMeta
	service     service.ResourceService[T]
	registry    *Registry
	rendererRegistry *RendererRegistry
	templates   *TemplateEngine
}

func NewServiceUIHandler[T any](meta *ModelMeta, svc service.ResourceService[T], registry *Registry,
	rendererRegistry *RendererRegistry, templates *TemplateEngine) *ServiceUIHandler[T] {
	return &ServiceUIHandler[T]{
		modelMeta:   meta,
		service:     svc,
		registry:    registry,
		rendererRegistry: rendererRegistry,
		templates:   templates,
	}
}

func (h *ServiceUIHandler[T]) ListUI(c *fiber.Ctx) error {
	ctx := c.UserContext()
	params := service.ParseListParams(c)
	
	result, err := h.service.List(ctx, params)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(403).SendString("You don't have permission to perform this action")
		}
		return c.Status(500).SendString(fmt.Sprintf("Error: %v", err))
	}
	
	objects := make([]map[string]interface{}, 0, len(result.Items))
	for _, obj := range result.Items {
		objects = append(objects, h.objectToMap(obj))
	}
	
	page := params.Page
	if page == 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize == 0 {
		pageSize = 20
	}
	
	data := fiber.Map{
		"model":    h.modelMeta,
		"objects":  objects,
		"page":     page,
		"pageSize": pageSize,
		"total":    result.Total,
		"search":   params.Search,
	}
	
	html, err := h.templates.RenderTemplate(h.modelMeta.Name, "changelist", data)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
	}
	
	return c.Type("html").SendString(html)
}

func (h *ServiceUIHandler[T]) FormUI(c *fiber.Ctx) error {
	id := c.Params("id")
	if id != "" {
		ctx := c.UserContext()
		obj, err := h.service.Retrieve(ctx, id)
		if err != nil {
			if err == service.ErrNotFound {
				return c.Status(404).SendString("Record not found")
			}
			return c.Status(500).SendString(fmt.Sprintf("Error: %v", err))
		}
		
		data := fiber.Map{
			"model":  h.modelMeta,
			"object": h.objectToMap(obj),
		}
		
		html, err := h.templates.RenderTemplate(h.modelMeta.Name, "form", data)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
		}
		return c.Type("html").SendString(html)
	}
	
	data := fiber.Map{
		"model":  h.modelMeta,
		"object": nil,
	}
	
	html, err := h.templates.RenderTemplate(h.modelMeta.Name, "form", data)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
	}
	return c.Type("html").SendString(html)
}

func (h *ServiceUIHandler[T]) DetailUI(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("ID is required")
	}
	
	ctx := c.UserContext()
	obj, err := h.service.Retrieve(ctx, id)
	if err != nil {
		if err == service.ErrPermissionDenied {
			return c.Status(403).SendString("You don't have permission to perform this action")
		}
		if err == service.ErrNotFound {
			return c.Status(404).SendString("Record not found")
		}
		return c.Status(500).SendString(fmt.Sprintf("Error: %v", err))
	}
	
	data := fiber.Map{
		"model":  h.modelMeta,
		"object": h.objectToMap(obj),
	}
	
	html, err := h.templates.RenderTemplate(h.modelMeta.Name, "detail", data)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
	}
	
	return c.Type("html").SendString(html)
}

func (h *ServiceUIHandler[T]) objectToMap(obj T) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		// Skip Schema embedded field
		if field.Name == "Schema" {
			continue
		}
		
		if fieldVal.CanInterface() {
			fieldName := field.Name
			// Convert to snake_case for consistency
			fieldName = toSnakeCase(fieldName)
			
			if fieldVal.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					result[fieldName] = nil
				} else {
					result[fieldName] = fieldVal.Elem().Interface()
				}
			} else {
				result[fieldName] = fieldVal.Interface()
			}
		}
	}
	return result
}

// toSnakeCase converts a camelCase or PascalCase string to snake_case.
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return string(result)
}

// Interface-based registration functions (for runtime type erasure)

func RegisterServiceForInterface(integration *ServiceIntegration,
	model interface{},
	svc interface{},
	options ...Option) error {

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	modelName := modelType.Name()

	err := integration.registry.Register(model, options...)
	if err != nil {
		return fmt.Errorf("failed to register model: %w", err)
	}

	meta, err := integration.registry.GetModel(modelName)
	if err != nil {
		return fmt.Errorf("failed to get model metadata: %w", err)
	}

	adapter := NewInterfaceServiceAdapter(svc, meta)
	integration.adapters[modelName] = adapter

	return nil
}

func NewInterfaceServiceAdapter(svc interface{}, modelMeta *ModelMeta) *InterfaceServiceAdapter {
	return &InterfaceServiceAdapter{
		service:   svc,
		modelMeta: modelMeta,
	}
}

type InterfaceServiceAdapter struct {
	service   interface{}
	modelMeta *ModelMeta
}

func (a *InterfaceServiceAdapter) GetService() interface{} {
	return a.service
}

func (a *InterfaceServiceAdapter) GetModelMeta() *ModelMeta {
	return a.modelMeta
}

func CreateInterfaceServiceUIHandler(integration *ServiceIntegration, modelName string, resourceService interface{}) interface{} {
	meta, err := integration.registry.GetModel(modelName)
	if err != nil {
		return nil
	}

	adapterInterface, ok := integration.GetAdapter(modelName)
	if !ok {
		return nil
	}

	adapter, ok := adapterInterface.(*InterfaceServiceAdapter)
	if !ok {
		return nil
	}

	svc := adapter.GetService()
	return NewInterfaceServiceUIHandler(meta, svc, integration.registry,
		integration.rendererRegistry, integration.templates)
}

func NewInterfaceServiceUIHandler(meta *ModelMeta, svc interface{}, registry *Registry,
	rendererRegistry *RendererRegistry, templates *TemplateEngine) interface{} {
	return &InterfaceServiceUIHandler{
		modelMeta:   meta,
		service:     svc,
		registry:    registry,
		rendererRegistry: rendererRegistry,
		templates:   templates,
	}
}

type InterfaceServiceUIHandler struct {
	modelMeta   *ModelMeta
	service     interface{}
	registry    *Registry
	rendererRegistry *RendererRegistry
	templates   *TemplateEngine
}

func (h *InterfaceServiceUIHandler) ListUI(c *fiber.Ctx) error {
	ctx := c.UserContext()
	params := service.ParseListParams(c)
	
	serviceValue := reflect.ValueOf(h.service)
	listMethod := serviceValue.MethodByName("List")
	if !listMethod.IsValid() {
		return c.Status(500).JSON(fiber.Map{"error": "List method not found"})
	}
	
	results := listMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(params),
	})
	
	if len(results) > 1 {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}
	
	return c.JSON(results[0].Interface())
}

func (h *InterfaceServiceUIHandler) FormUI(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	
	if id != "" {
		serviceValue := reflect.ValueOf(h.service)
		retrieveMethod := serviceValue.MethodByName("Retrieve")
		if !retrieveMethod.IsValid() {
			return c.Status(500).JSON(fiber.Map{"error": "Retrieve method not found"})
		}
		
		results := retrieveMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(id),
		})
		
		if len(results) > 1 {
			if err, ok := results[1].Interface().(error); ok && err != nil {
				return c.Status(404).JSON(fiber.Map{"error": err.Error()})
			}
		}
		
		obj := results[0].Interface()
		data := fiber.Map{
			"model":  h.modelMeta,
			"object": obj,
		}
		
		html, err := h.templates.RenderTemplate(h.modelMeta.Name, "form", data)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
		}
		
		return c.Type("html").SendString(html)
	}
	
	data := fiber.Map{
		"model":  h.modelMeta,
		"object": nil,
	}
	
	html, err := h.templates.RenderTemplate(h.modelMeta.Name, "form", data)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
	}
	
	return c.Type("html").SendString(html)
}

func (h *InterfaceServiceUIHandler) DetailUI(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := c.UserContext()
	
	serviceValue := reflect.ValueOf(h.service)
	retrieveMethod := serviceValue.MethodByName("Retrieve")
	if !retrieveMethod.IsValid() {
		return c.Status(500).JSON(fiber.Map{"error": "Retrieve method not found"})
	}
	
	results := retrieveMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})
	
	if len(results) > 1 {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
	}
	
	return c.JSON(results[0].Interface())
}

