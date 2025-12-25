package models

import (
	"fmt"
	"reflect"
	"sync"
)

// Registry manages model registration (Ent-inspired)
type Registry struct {
	models    map[string]*ModelInfo
	managers  map[string]Manager
	mu        sync.RWMutex
}

// ModelInfo contains metadata about a registered model
type ModelInfo struct {
	Type    reflect.Type
	Manager Manager
	Meta    *Meta
}

// NewRegistry creates a new model registry
func NewRegistry() *Registry {
	return &Registry{
		models:   make(map[string]*ModelInfo),
		managers: make(map[string]Manager),
	}
}

// Register registers a model type with its manager
func (r *Registry) Register(modelType interface{}, manager Manager, meta *Meta) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	typ := reflect.TypeOf(modelType)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	name := typ.Name()
	if _, exists := r.models[name]; exists {
		return fmt.Errorf("model %s is already registered", name)
	}
	
	if meta == nil {
		meta = DefaultMeta()
	}
	
	r.models[name] = &ModelInfo{
		Type:    typ,
		Manager: manager,
		Meta:    meta,
	}
	r.managers[name] = manager
	
	return nil
}

// GetModelInfo retrieves model info by name
func (r *Registry) GetModelInfo(name string) (*ModelInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	info, exists := r.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s is not registered", name)
	}
	
	return info, nil
}

// GetManager retrieves a manager by model name
func (r *Registry) GetManager(name string) (Manager, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	manager, exists := r.managers[name]
	if !exists {
		return nil, fmt.Errorf("no manager registered for model %s", name)
	}
	
	return manager, nil
}

// GetAllModels returns all registered models
func (r *Registry) GetAllModels() map[string]*ModelInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make(map[string]*ModelInfo)
	for k, v := range r.models {
		result[k] = v
	}
	return result
}

// CreateInstance creates a new instance of a registered model
func (r *Registry) CreateInstance(name string) (Model, error) {
	info, err := r.GetModelInfo(name)
	if err != nil {
		return nil, err
	}
	
	instance := reflect.New(info.Type).Interface()
	if model, ok := instance.(Model); ok {
		if baseModel, ok := model.(interface{ SetManager(Manager) }); ok {
			baseModel.SetManager(info.Manager)
		}
		return model, nil
	}
	
	return nil, fmt.Errorf("model %s does not implement Model interface", name)
}

// Global registry instance
var globalRegistry = NewRegistry()

// RegisterModel registers a model in the global registry
func RegisterModel(modelType interface{}, manager Manager, meta *Meta) error {
	return globalRegistry.Register(modelType, manager, meta)
}

// GetModelInfo retrieves model info from global registry
func GetModelInfo(name string) (*ModelInfo, error) {
	return globalRegistry.GetModelInfo(name)
}

// GetModelManager retrieves a manager from global registry
func GetModelManager(name string) (Manager, error) {
	return globalRegistry.GetManager(name)
}

