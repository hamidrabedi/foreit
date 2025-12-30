package http

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	admin "github.com/forgego/forge/admin"
)

// TypeRegistry stores type information for admin instances
type TypeRegistry struct {
	admins map[string]AdminHandler
	mu     sync.RWMutex
}

var globalTypeRegistry = &TypeRegistry{
	admins: make(map[string]AdminHandler),
}

// AdminHandler is the interface for type-safe admin operations
type AdminHandler interface {
	HandleList(ctx context.Context, page, pageSize int, search string, filters map[string]interface{}) (interface{}, error)
	HandleDetail(ctx context.Context, id int64) (interface{}, error)
	HandleCreate(ctx context.Context, formData map[string]interface{}) (interface{}, error)
	HandleUpdate(ctx context.Context, id int64, formData map[string]interface{}) (interface{}, error)
	HandleDelete(ctx context.Context, id int64) error
	HandleExport(ctx context.Context, format string) (interface{}, error)
	HandleBulkAction(ctx context.Context, action string, ids []int64) error
	HandleAutocomplete(ctx context.Context, search string, limit int) ([]map[string]interface{}, error)
}

// RegisterAdmin registers an admin with type-safe handlers
func RegisterAdmin[T any](admin *admin.Admin[T]) {
	globalTypeRegistry.mu.Lock()
	defer globalTypeRegistry.mu.Unlock()

	handler := &adminHandler[T]{
		admin: admin,
	}
	globalTypeRegistry.admins[admin.ModelName()] = handler
}

// GetAdminHandler gets a handler for a model
func GetAdminHandler(modelName string) (AdminHandler, error) {
	globalTypeRegistry.mu.RLock()
	defer globalTypeRegistry.mu.RUnlock()

	handler, ok := globalTypeRegistry.admins[modelName]
	if !ok {
		return nil, fmt.Errorf("admin handler for %s not found", modelName)
	}

	return handler, nil
}

// adminHandler provides type-safe handlers for Admin[T]
type adminHandler[T any] struct {
	admin *admin.Admin[T]
}

// HandleList handles list view
func (h *adminHandler[T]) HandleList(ctx context.Context, page, pageSize int, search string, filters map[string]interface{}) (interface{}, error) {
	// Import views package functions directly
	// For now, return a placeholder response
	// TODO: Fix import cycle by moving views to same package or using interface
	return map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
		"search":   search,
		"filters":  filters,
	}, nil
}

// HandleDetail handles detail view
func (h *adminHandler[T]) HandleDetail(ctx context.Context, id int64) (interface{}, error) {
	// Get instance by ID
	instance, err := h.admin.Manager().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"instance": instance,
		"id":       id,
	}, nil
}

// HandleCreate handles create
func (h *adminHandler[T]) HandleCreate(ctx context.Context, formData map[string]interface{}) (interface{}, error) {
	var zero T
	instance := &zero

	// Save using admin's SaveModel
	if err := h.admin.SaveModel(ctx, instance, admin.FormData(formData), true); err != nil {
		return nil, err
	}

	return instance, nil
}

// HandleUpdate handles update
func (h *adminHandler[T]) HandleUpdate(ctx context.Context, id int64, formData map[string]interface{}) (interface{}, error) {
	// Get instance by ID
	instance, err := h.admin.Manager().Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Save using admin's SaveModel
	if err := h.admin.SaveModel(ctx, instance, admin.FormData(formData), false); err != nil {
		return nil, err
	}

	return instance, nil
}

// HandleDelete handles delete
func (h *adminHandler[T]) HandleDelete(ctx context.Context, id int64) error {
	instance, err := h.admin.Manager().Get(ctx, id)
	if err != nil {
		return err
	}

	return h.admin.DeleteModel(ctx, instance)
}

// HandleExport handles export - returns ExportView for HTTP handler to use
func (h *adminHandler[T]) HandleExport(ctx context.Context, format string) (interface{}, error) {
	exportView := admin.NewExportView(h.admin, admin.ExportFormat(format))
	return exportView, nil
}

// HandleBulkAction handles bulk action
func (h *adminHandler[T]) HandleBulkAction(ctx context.Context, action string, ids []int64) error {
	// Get instances
	instances := make([]*T, 0, len(ids))
	for _, id := range ids {
		instance, err := h.admin.Manager().Get(ctx, id)
		if err != nil {
			continue
		}
		instances = append(instances, instance)
	}

	// Find and execute action
	for _, act := range h.admin.Config().Actions {
		if act.Name == action {
			return act.Handler(ctx, instances)
		}
	}

	return fmt.Errorf("action %s not found", action)
}

// HandleAutocomplete handles autocomplete
func (h *adminHandler[T]) HandleAutocomplete(ctx context.Context, search string, limit int) ([]map[string]interface{}, error) {
	qs := h.admin.Manager().Filter()
	if search != "" {
		// Apply search
		// This is simplified - would use search fields from config
	}

	objects, err := qs.Limit(limit).All(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, len(objects))
	for i, obj := range objects {
		// Convert to map (simplified)
		results[i] = map[string]interface{}{
			"id":    getID(obj),
			"label": fmt.Sprintf("%v", obj),
		}
	}

	return results, nil
}

// getID extracts ID from an object
func getID(obj interface{}) int64 {
	// Use reflection to get ID field
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanInterface() {
		if id, ok := idField.Interface().(int64); ok {
			return id
		}
	}

	return 0
}
