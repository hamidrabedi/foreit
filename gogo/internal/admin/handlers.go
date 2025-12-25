package admin

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CRUDHandler handles CRUD operations for a registered model
type CRUDHandler struct {
	modelMeta         *ModelMeta
	client            interface{} // *ent.Client
	entHelper         *EntClientHelper
	queryBuilder      *QueryBuilder
	permissionChecker *PermissionChecker
	hookRegistry      *HookRegistry
}

// NewCRUDHandler creates a new CRUD handler for a model
func NewCRUDHandler(meta *ModelMeta, client interface{}, registry *Registry) *CRUDHandler {
	return &CRUDHandler{
		modelMeta:         meta,
		client:            client,
		entHelper:         NewEntClientHelper(client),
		queryBuilder:      NewQueryBuilder(meta),
		permissionChecker: NewPermissionChecker(registry),
		hookRegistry:      registry.GetHookRegistry(),
	}
}

// List handles GET /admin/api/{model} - List all records with pagination and filters
func (h *CRUDHandler) List(c *fiber.Ctx) error {
	// Check permissions using permission checker
	permCtx := &Context{
		User:     GetUserFromContext(c),
		Request:  c,
		Model:    h.modelMeta,
		Action:   "list",
	}
	
	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to list this resource",
		})
	}

	// Execute before hook
	hookCtx := &HookContext{
		Model:   h.modelMeta,
		Action:  "list",
		User:    GetUserFromContext(c),
		Request: c,
	}
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeList, hookCtx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	// Parse all query parameters
	queryParams := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	// Parse query parameters using query builder
	params := h.queryBuilder.ParseQueryParams(queryParams)
	hookCtx.Query = params

	// Build and execute query using Ent client
	ctx := context.Background()
	results, total, err := h.entHelper.List(ctx, h.modelMeta.Name, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list records",
			"details": err.Error(),
		})
	}

	// Convert results to map format for JSON response
	data := make([]map[string]interface{}, len(results))
	for i, result := range results {
		data[i] = h.modelToMap(result)
	}

	response := fiber.Map{
		"data": data,
		"pagination": CalculatePagination(params.Page, params.PageSize, total),
		"filters": params.Filters,
		"sort": fiber.Map{
			"by":    params.SortBy,
			"order": params.SortOrder,
		},
	}
	
	// Execute after hook
	hookCtx.Result = response
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterList, hookCtx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}

// Get handles GET /admin/api/{model}/{id} - Get a single record
func (h *CRUDHandler) Get(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "view",
	}
	
	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to view this resource",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	// Get record by ID using Ent client
	ctx := context.Background()
	result, err := h.entHelper.Get(ctx, h.modelMeta.Name, id)
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get record",
			"details": err.Error(),
		})
	}

	// Check rule-based permissions with the actual resource
	permCtx.Resource = result
	allowed, err = h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to view this resource",
		})
	}

	response := fiber.Map{
		"data": h.modelToMap(result),
	}

	return c.JSON(response)
}

// Create handles POST /admin/api/{model} - Create a new record
func (h *CRUDHandler) Create(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "create",
	}
	
	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to create this resource",
		})
	}

	// Parse request body
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Execute before hook
	hookCtx := &HookContext{
		Model:   h.modelMeta,
		Action:  "create",
		User:    GetUserFromContext(c),
		Request: c,
		Data:    data,
	}
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeCreate, hookCtx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	// Validate data
	if err := h.validateData(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Create record using Ent client
	ctx := context.Background()
	result, err := h.entHelper.Create(ctx, h.modelMeta.Name, data)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create record",
			"details": err.Error(),
		})
	}

	response := fiber.Map{
		"data":    h.modelToMap(result),
		"message": "Record created successfully",
	}
	
	// Execute after hook
	hookCtx.Result = response
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterCreate, hookCtx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Update handles PUT /admin/api/{model}/{id} - Update a record
func (h *CRUDHandler) Update(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "update",
	}
	
	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this resource",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	// Parse request body
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Execute before hook
	hookCtx := &HookContext{
		Model:    h.modelMeta,
		Action:   "update",
		User:     GetUserFromContext(c),
		Request:  c,
		Data:     data,
		Resource: id,
	}
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeUpdate, hookCtx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	// Validate data
	if err := h.validateData(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Get existing record for permission check
	ctx := context.Background()
	existing, err := h.entHelper.Get(ctx, h.modelMeta.Name, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get record",
			"details": err.Error(),
		})
	}

	// Check rule-based permissions with the existing resource
	permCtx.Resource = existing
	allowed, err = h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to update this resource",
		})
	}

	// Update record using Ent client
	result, err := h.entHelper.Update(ctx, h.modelMeta.Name, id, data)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to update record",
			"details": err.Error(),
		})
	}

	response := fiber.Map{
		"data":    h.modelToMap(result),
		"message": "Record updated successfully",
	}
	
	// Execute after hook
	hookCtx.Result = response
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterUpdate, hookCtx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}

// Delete handles DELETE /admin/api/{model}/{id} - Delete a record
func (h *CRUDHandler) Delete(c *fiber.Ctx) error {
	// Check permissions
	permCtx := &Context{
		User:    GetUserFromContext(c),
		Request: c,
		Model:   h.modelMeta,
		Action:  "delete",
	}
	
	allowed, err := h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this resource",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	// Execute before hook
	hookCtx := &HookContext{
		Model:    h.modelMeta,
		Action:   "delete",
		User:     GetUserFromContext(c),
		Request:  c,
		Resource: id,
	}
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookBeforeDelete, hookCtx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}

	// Get existing record for permission check
	ctx := context.Background()
	existing, err := h.entHelper.Get(ctx, h.modelMeta.Name, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get record",
			"details": err.Error(),
		})
	}

	// Check rule-based permissions with the existing resource
	permCtx.Resource = existing
	allowed, err = h.permissionChecker.CheckPermission(permCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Permission check failed",
			"details": err.Error(),
		})
	}
	if !allowed {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this resource",
		})
	}

	// Delete record using Ent client
	err = h.entHelper.Delete(ctx, h.modelMeta.Name, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete record",
			"details": err.Error(),
		})
	}
	
	// Execute after hook
	if err := h.hookRegistry.Execute(h.modelMeta.Name, HookAfterDelete, hookCtx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Hook execution failed",
			"details": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// parseFilters is now handled by QueryBuilder.ParseQueryParams
// This method is kept for backward compatibility but delegates to query builder
func (h *CRUDHandler) parseFilters(c *fiber.Ctx) map[string]interface{} {
	queryParams := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})
	
	params := h.queryBuilder.ParseQueryParams(queryParams)
	return params.Filters
}

// isFieldFilterable checks if a field is filterable
func (h *CRUDHandler) isFieldFilterable(fieldName string) bool {
	// Check if field is in filterable fields list
	if len(h.modelMeta.Options.FilterableFields) > 0 {
		for _, f := range h.modelMeta.Options.FilterableFields {
			if f == fieldName {
				return true
			}
		}
		return false
	}

	// Default: check field metadata
	for _, field := range h.modelMeta.Fields {
		if field.Name == fieldName && field.Filterable {
			return true
		}
	}

	return false
}

// validateData validates the data against the model's field definitions
func (h *CRUDHandler) validateData(data map[string]interface{}) error {
	var errors []string

	for _, field := range h.modelMeta.Fields {
		value, exists := data[field.Name]

		// Check required fields
		if field.Required && !exists {
			errors = append(errors, fmt.Sprintf("Field '%s' is required", field.Label))
			continue
		}

		// Skip validation if field doesn't exist and is optional
		if !exists {
			continue
		}

		// Check if field is read-only (should not be in update/create)
		if field.ReadOnly {
			errors = append(errors, fmt.Sprintf("Field '%s' is read-only", field.Label))
			continue
		}

		// Type validation
		if err := h.validateFieldType(field, value); err != nil {
			errors = append(errors, fmt.Sprintf("Field '%s': %s", field.Label, err.Error()))
		}

		// Enum validation
		if len(field.Choices) > 0 {
			if !h.isValidChoice(field, value) {
				errors = append(errors, fmt.Sprintf("Field '%s' must be one of: %v", field.Label, field.Choices))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}

// validateFieldType validates that a value matches the expected field type
func (h *CRUDHandler) validateFieldType(field FieldMeta, value interface{}) error {
	valueType := reflect.TypeOf(value)

	switch field.Type {
	case FieldTypeNumber:
		if !h.isNumericType(valueType) {
			return fmt.Errorf("expected number, got %s", valueType.Kind())
		}

	case FieldTypeBoolean:
		if valueType.Kind() != reflect.Bool {
			return fmt.Errorf("expected boolean, got %s", valueType.Kind())
		}

	case FieldTypeEmail, FieldTypeURL, FieldTypeText, FieldTypeTextarea:
		if valueType.Kind() != reflect.String {
			return fmt.Errorf("expected string, got %s", valueType.Kind())
		}

	case FieldTypeDate, FieldTypeDateTime, FieldTypeTime:
		// Accept string or time.Time
		if valueType.Kind() != reflect.String && valueType != reflect.TypeOf("") {
			return fmt.Errorf("expected date/time string, got %s", valueType.Kind())
		}
	}

	return nil
}

// isValidChoice checks if a value is a valid choice for an enum field
func (h *CRUDHandler) isValidChoice(field FieldMeta, value interface{}) bool {
	valueStr := fmt.Sprintf("%v", value)
	for _, choice := range field.Choices {
		if choice.Value == valueStr {
			return true
		}
	}
	return false
}

// isNumericType checks if a type is numeric
func (h *CRUDHandler) isNumericType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// modelToMap converts an Ent model to a map for JSON serialization
func (h *CRUDHandler) modelToMap(model interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(model)
	
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return result
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}
		
		// Get JSON tag or use field name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = field.Name
		} else {
			// Handle json tag options like "name,omitempty"
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonTag = jsonTag[:idx]
			}
		}
		
		// Convert field value to interface
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				result[jsonTag] = nil
			} else {
				result[jsonTag] = fieldVal.Elem().Interface()
			}
		} else {
			result[jsonTag] = fieldVal.Interface()
		}
	}
	
	return result
}

