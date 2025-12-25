package admin

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/api"
	"github.com/forgego/forge/pkg/models"
)

type ViewSetAdapter[T any] struct {
	viewSet  api.ViewSet[T]
	manager  models.Manager[T]
	modelDef *models.ModelDefinition[T]
	modelMeta *ModelMeta
}

func NewViewSetAdapter[T any](
	viewSet api.ViewSet[T],
	manager models.Manager[T],
	modelDef *models.ModelDefinition[T],
	modelMeta *ModelMeta,
) *ViewSetAdapter[T] {
	return &ViewSetAdapter[T]{
		viewSet:   viewSet,
		manager:  manager,
		modelDef: modelDef,
		modelMeta: modelMeta,
	}
}

func (a *ViewSetAdapter[T]) GetViewSet() api.ViewSet[T] {
	return a.viewSet
}

func (a *ViewSetAdapter[T]) GetManager() models.Manager[T] {
	return a.manager
}

func (a *ViewSetAdapter[T]) GetModelDefinition() *models.ModelDefinition[T] {
	return a.modelDef
}

func (a *ViewSetAdapter[T]) GetModelMeta() *ModelMeta {
	return a.modelMeta
}

type ViewSetIntegration struct {
	registry    *Registry
	rendererRegistry *RendererRegistry
	templates   *TemplateEngine
	db          *models.DB
	adapters    map[string]interface{}
}

func NewViewSetIntegration(
	registry *Registry,
	rendererRegistry *RendererRegistry,
	templates *TemplateEngine,
	db *models.DB,
) *ViewSetIntegration {
	return &ViewSetIntegration{
		registry:    registry,
		rendererRegistry: rendererRegistry,
		templates:   templates,
		db:          db,
		adapters:    make(map[string]interface{}),
	}
}

func RegisterViewSet[T any](
	integration *ViewSetIntegration,
	model interface{},
	viewSet api.ViewSet[T],
	manager models.Manager[T],
	modelDef *models.ModelDefinition[T],
	options ...Option,
) error {
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

	if modelDef != nil {
		if meta.TableName == "" {
			meta.TableName = modelDef.GetTableName()
		}
	}

	adapter := NewViewSetAdapter(viewSet, manager, modelDef, meta)
	integration.adapters[modelName] = adapter

	return nil
}

func (i *ViewSetIntegration) GetAdapter(modelName string) (interface{}, bool) {
	adapter, ok := i.adapters[modelName]
	return adapter, ok
}

func CreateViewSetUIHandler[T any](
	integration *ViewSetIntegration,
	modelName string,
) (*ViewSetUIHandler[T], error) {
	meta, err := integration.registry.GetModel(modelName)
	if err != nil {
		return nil, err
	}

	adapterInterface, ok := integration.GetAdapter(modelName)
	if !ok {
		return nil, fmt.Errorf("no viewset registered for model %s", modelName)
	}

	adapter, ok := adapterInterface.(*ViewSetAdapter[T])
	if !ok {
		return nil, fmt.Errorf("invalid adapter type for model %s", modelName)
	}

	return NewViewSetUIHandler(
		meta,
		adapter.GetViewSet(),
		adapter.GetManager(),
		adapter.GetModelDefinition(),
		integration.registry,
		integration.rendererRegistry,
		integration.templates,
	), nil
}

type ViewSetUIHandler[T any] struct {
	modelMeta   *ModelMeta
	viewSet     api.ViewSet[T]
	manager     models.Manager[T]
	modelDef    *models.ModelDefinition[T]
	serializer  api.Serializer[T]
	registry    *Registry
	rendererRegistry *RendererRegistry
	templates   *TemplateEngine
}

func NewViewSetUIHandler[T any](
	meta *ModelMeta,
	viewSet api.ViewSet[T],
	manager models.Manager[T],
	modelDef *models.ModelDefinition[T],
	registry *Registry,
	rendererRegistry *RendererRegistry,
	templates *TemplateEngine,
) *ViewSetUIHandler[T] {
	var serializer api.Serializer[T]
	if baseViewSet, ok := viewSet.(*api.BaseViewSet[T]); ok {
		serializer = baseViewSet.GetSerializer()
	}
	
	return &ViewSetUIHandler[T]{
		modelMeta:   meta,
		viewSet:     viewSet,
		manager:     manager,
		modelDef:    modelDef,
		serializer:  serializer,
		registry:    registry,
		rendererRegistry: rendererRegistry,
		templates:   templates,
	}
}

func (h *ViewSetUIHandler[T]) ListUI(c *fiber.Ctx) error {
	ctx := c.UserContext()
	qs := h.manager.All()

	// Apply search
	searchQuery := c.Query("search")
	if searchQuery != "" && len(h.modelMeta.Options.SearchFields) > 0 {
		searchConditions := h.buildSearchConditions(searchQuery)
		if len(searchConditions) > 0 {
			qs = qs.Filter(searchConditions...)
		}
	}

	// Apply filters
	h.applyFilters(qs, c)

	// Get total count before pagination
	total, err := qs.Count(ctx)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error counting: %v", err))
	}

	// Apply pagination
	page := h.getPage(c)
	pageSize := h.getPageSize(c)
	offset := (page - 1) * pageSize
	qs = qs.Limit(pageSize).Offset(offset)

	// Apply sorting
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order", "asc")
	if sortBy != "" {
		direction := models.OrderAsc
		if sortOrder == "desc" {
			direction = models.OrderDesc
		}
		qs = qs.OrderBy(models.OrderBy{Field: sortBy, Direction: direction})
	} else {
		// Default sort by ID descending
		qs = qs.OrderBy(models.OrderBy{Field: "id", Direction: models.OrderDesc})
	}

	results, err := qs.All(ctx)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error: %v", err))
	}

	objects := make([]map[string]interface{}, 0, len(results))
	for _, obj := range results {
		objects = append(objects, h.objectToMap(obj))
	}

	data := fiber.Map{
		"model":    h.modelMeta,
		"objects":  objects,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"search":   searchQuery,
	}

	html, err := h.templates.RenderTemplate(h.modelMeta.Name, "changelist", data)
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Template error: %v", err))
	}

	return c.Type("html").SendString(html)
}

func (h *ViewSetUIHandler[T]) buildSearchConditions(query string) []models.Condition {
	if query == "" || len(h.modelMeta.Options.SearchFields) == 0 {
		return nil
	}

	conditions := make([]models.Condition, 0)
	for _, fieldName := range h.modelMeta.Options.SearchFields {
		condition := models.NewStringCondition(fieldName, "LIKE", "%"+query+"%")
		conditions = append(conditions, condition)
	}

	return conditions
}

func (h *ViewSetUIHandler[T]) applyFilters(qs models.QuerySet[T], c *fiber.Ctx) {
	// Apply filters from query parameters
	// Format: filter_<fieldname>=<value>
	// Note: This is a simplified version - in production, you'd parse all query params
	// For now, we'll apply filters manually based on known filterable fields
	for _, fieldName := range h.modelMeta.Options.FilterableFields {
		filterValue := c.Query("filter_" + fieldName)
		if filterValue != "" {
			condition := models.NewStringCondition(fieldName, "=", filterValue)
			qs = qs.Filter(condition)
		}
	}
}

func (h *ViewSetUIHandler[T]) getPage(c *fiber.Ctx) int {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func (h *ViewSetUIHandler[T]) getPageSize(c *fiber.Ctx) int {
	pageSize, err := strconv.Atoi(c.Query("page_size", "20"))
	if err != nil || pageSize < 1 {
		return 20
	}
	if pageSize > 100 {
		return 100
	}
	return pageSize
}

func (h *ViewSetUIHandler[T]) FormUI(c *fiber.Ctx) error {
	id := c.Params("id")
	if id != "" {
		ctx := c.UserContext()
		condition := parseIDCondition(id)
		obj, err := h.manager.Get(ctx, condition)
		if err != nil {
			return c.Status(404).SendString("Record not found")
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

func (h *ViewSetUIHandler[T]) DetailUI(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).SendString("ID is required")
	}

	ctx := c.UserContext()
	condition := parseIDCondition(id)
	obj, err := h.manager.Get(ctx, condition)
	if err != nil {
		return c.Status(404).SendString("Record not found")
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

func (h *ViewSetUIHandler[T]) objectToMap(obj *T) map[string]interface{} {
	if h.serializer != nil {
		rep, err := h.serializer.ToRepresentation(obj)
		if err == nil {
			if result, ok := rep.(map[string]interface{}); ok {
				return result
			}
		}
	}
	
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		if fieldVal.CanInterface() {
			result[field.Name] = fieldVal.Interface()
		}
	}
	return result
}

func (h *ViewSetUIHandler[T]) GetFieldAdapters() map[string]api.Field {
	if h.modelDef == nil {
		return nil
	}
	
	modelFields := h.modelDef.GetFields()
	return api.AdaptModelFields(modelFields)
}

func (h *ViewSetUIHandler[T]) GetFieldValidators(fieldName string) []interface{} {
	adapters := h.GetFieldAdapters()
	if adapters == nil {
		return nil
	}
	
	adapter, ok := adapters[fieldName]
	if !ok {
		return nil
	}
	
	if fieldAdapter, ok := adapter.(*api.FieldAdapter); ok {
		validators := fieldAdapter.GetModelFieldDefinition()
		if withValidators, ok := validators.(models.FieldDefinitionWithValidators); ok {
			return withValidators.GetValidators()
		}
	}
	
	return nil
}

// parseIDCondition parses an ID parameter and returns the appropriate condition.
// Tries to parse as int64 first, falls back to string if parsing fails.
func parseIDCondition(id string) models.Condition {
	if intID, err := strconv.ParseInt(id, 10, 64); err == nil {
		return models.NewIntCondition("id", "=", intID)
	}
	return models.NewStringCondition("id", "=", id)
}

