package admin

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/forgego/forge/pkg/registry"
	"github.com/forgego/forge/pkg/schema"
)

// AdminRegistry maintains a registry of models for admin auto-generation
type AdminRegistry struct {
	models map[string]*AdminModel
	mu     sync.RWMutex
}

// AdminModel contains admin configuration for a model
// Auto-generated from model metadata (Django-style)
type AdminModel struct {
	Name          string
	Model         interface{}
	Manager       interface{} // Manager instance for CRUD operations (typically *query.Manager[T])
	
	// Auto-generated from Meta and field definitions
	ListDisplay   []interface{} // Supports both string and FieldExpr
	ListFilter    []interface{}
	SearchFields  []interface{}
	ReadOnlyFields []interface{}
	
	// Extended configuration (Django ModelAdmin features)
	ExtendedConfig map[string]interface{} // Stores ModelAdminConfig values
	
	// Extensibility
	CustomAdmin CustomAdmin
}

// CustomAdmin is the interface for custom admin classes
type CustomAdmin interface {
	GetListDisplay() []interface{}
	GetListFilter() []interface{}
	GetSearchFields() []interface{}
	GetReadOnlyFields() []interface{}
}

var globalAdminRegistry = &AdminRegistry{
	models: make(map[string]*AdminModel),
}

// RegisterModel registers a model for admin auto-generation
// This is the simple Django-style registration
func RegisterModel(model interface{}) error {
	return RegisterModelWithOptions(model)
}

// RegisterModelWithOptions registers a model with custom admin options
func RegisterModelWithOptions(model interface{}, options ...AdminOption) error {
	globalAdminRegistry.mu.Lock()
	defer globalAdminRegistry.mu.Unlock()
	
	// Get model name using type assertion or reflection as fallback
	var modelName string
	if named, ok := model.(interface{ Name() string }); ok {
		modelName = named.Name()
	} else {
		// Fallback to reflection for model name
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		modelName = modelType.Name()
	}
	
	// Get model info from registry (optional - model might not be registered yet)
	_, err := registry.GetModel(modelName)
	if err != nil {
		// Model not in registry yet - that's okay, we'll register it
		if err := registry.RegisterModel(model); err != nil {
			// Ignore if already registered
		}
	}
	
	// Auto-generate admin config from model metadata
	adminModel := &AdminModel{
		Name:  modelName,
		Model: model,
	}
	
	// Auto-generate from schema if available
	// Check if model implements schema.Schema interface
	if schemaModel, ok := model.(schema.Schema); ok {
		adminModel = autoGenerateAdminFromSchema(adminModel, schemaModel)
	}
	
	// Apply custom options
	for _, opt := range options {
		opt(adminModel)
	}
	
	globalAdminRegistry.models[modelName] = adminModel
	return nil
}

// autoGenerateAdminFromSchema is now in auto_generate.go

// GetModel retrieves an admin model
func GetModel(name string) (*AdminModel, error) {
	globalAdminRegistry.mu.RLock()
	defer globalAdminRegistry.mu.RUnlock()
	
	model, exists := globalAdminRegistry.models[name]
	if !exists {
		return nil, fmt.Errorf("admin model %s is not registered", name)
	}
	
	return model, nil
}

// GetAllModels returns all registered admin models
func GetAllModels() map[string]*AdminModel {
	globalAdminRegistry.mu.RLock()
	defer globalAdminRegistry.mu.RUnlock()
	
	result := make(map[string]*AdminModel)
	for k, v := range globalAdminRegistry.models {
		result[k] = v
	}
	
	return result
}

// AdminOption is a function that configures an admin model
type AdminOption func(*AdminModel)

// WithListDisplay sets the list display fields
func WithListDisplay(fields ...interface{}) AdminOption {
	return func(m *AdminModel) {
		m.ListDisplay = fields
	}
}

// WithListFilter sets the list filter fields
func WithListFilter(fields ...interface{}) AdminOption {
	return func(m *AdminModel) {
		m.ListFilter = fields
	}
}

// WithSearchFields sets the search fields
func WithSearchFields(fields ...interface{}) AdminOption {
	return func(m *AdminModel) {
		m.SearchFields = fields
	}
}

// WithReadOnlyFields sets the read-only fields
func WithReadOnlyFields(fields ...interface{}) AdminOption {
	return func(m *AdminModel) {
		m.ReadOnlyFields = fields
	}
}

// WithCustomAdmin sets a custom admin class
func WithCustomAdmin(custom CustomAdmin) AdminOption {
	return func(m *AdminModel) {
		m.CustomAdmin = custom
	}
}

// WithManager sets the manager for the model
func WithManager(manager interface{}) AdminOption {
	return func(m *AdminModel) {
		m.Manager = manager
	}
}

// RegisterModelWithManager registers a model with its manager
func RegisterModelWithManager(model interface{}, manager interface{}) error {
	return RegisterModelWithOptions(model, WithManager(manager))
}

