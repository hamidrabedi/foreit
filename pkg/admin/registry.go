package admin

import (
	"fmt"
	"reflect"
	"sync"
)

// Registry maintains a registry of admin instances
type Registry struct {
	admins map[string]AdminInterface
	mu     sync.RWMutex
}

var globalRegistry = &Registry{
	admins: make(map[string]AdminInterface),
}

// AdminInterface is the interface for dynamic access to admin instances
type AdminInterface interface {
	ModelName() string
	ModelType() reflect.Type
	ManagerInterface() interface{}
	ConfigInterface() interface{}
}

// register registers an admin instance
func (r *Registry) register(admin AdminInterface) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := admin.ModelName()
	r.admins[name] = admin
}

// Get retrieves an admin instance by name
func (r *Registry) Get(name string) (AdminInterface, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	admin, ok := r.admins[name]
	if !ok {
		return nil, fmt.Errorf("admin %s is not registered", name)
	}

	return admin, nil
}

// GetAll returns all registered admin instances
func (r *Registry) GetAll() map[string]AdminInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]AdminInterface)
	for k, v := range r.admins {
		result[k] = v
	}

	return result
}

// ModelType returns the model type
func (a *Admin[T]) ModelType() reflect.Type {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

// ManagerInterface returns the manager as interface{} for dynamic access
func (a *Admin[T]) ManagerInterface() interface{} {
	return a.manager
}

// ConfigInterface returns the config as interface{} for dynamic access
func (a *Admin[T]) ConfigInterface() interface{} {
	return a.config
}
