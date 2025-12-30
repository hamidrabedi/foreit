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

// GetGlobalRegistry returns the global registry
func GetGlobalRegistry() *Registry {
	return globalRegistry
}

// ModelType, ManagerInterface, and ConfigInterface are defined in admin.go

// Permission names are defined in permissions.go
