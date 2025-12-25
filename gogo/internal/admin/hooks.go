package admin

import (
	"fmt"
	"sync"
)

// HookType represents the type of hook
type HookType string

const (
	HookBeforeCreate HookType = "before_create"
	HookAfterCreate  HookType = "after_create"
	HookBeforeUpdate HookType = "before_update"
	HookAfterUpdate  HookType = "after_update"
	HookBeforeDelete HookType = "before_delete"
	HookAfterDelete  HookType = "after_delete"
	HookBeforeList   HookType = "before_list"
	HookAfterList    HookType = "after_list"
	HookBeforeGet    HookType = "before_get"
	HookAfterGet     HookType = "after_get"
)

// HookFunc is a function that can be registered as a hook
type HookFunc func(ctx *HookContext) error

// HookContext provides context for hook execution
type HookContext struct {
	Model     *ModelMeta
	Action    string
	User      interface{}
	Request   interface{}
	Data      interface{} // Input data for create/update
	Resource  interface{} // Resource being accessed
	Query     interface{} // Query object (for list/get)
	Result    interface{} // Result data (for after hooks)
	Error     error       // Error from operation
}

// HookRegistry manages hooks for models
type HookRegistry struct {
	hooks map[string]map[HookType][]HookFunc
	mu    sync.RWMutex
}

// NewHookRegistry creates a new hook registry
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[string]map[HookType][]HookFunc),
	}
}

// Register registers a hook for a model
func (hr *HookRegistry) Register(modelName string, hookType HookType, hook HookFunc) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	
	if hr.hooks[modelName] == nil {
		hr.hooks[modelName] = make(map[HookType][]HookFunc)
	}
	
	hr.hooks[modelName][hookType] = append(hr.hooks[modelName][hookType], hook)
}

// Execute executes all hooks of a given type for a model
func (hr *HookRegistry) Execute(modelName string, hookType HookType, ctx *HookContext) error {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	
	modelHooks, exists := hr.hooks[modelName]
	if !exists {
		return nil
	}
	
	hooks, exists := modelHooks[hookType]
	if !exists {
		return nil
	}
	
	// Execute all hooks in order
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return fmt.Errorf("hook execution failed: %w", err)
		}
	}
	
	return nil
}

// GetHooks returns all hooks for a model and hook type
func (hr *HookRegistry) GetHooks(modelName string, hookType HookType) []HookFunc {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	
	modelHooks, exists := hr.hooks[modelName]
	if !exists {
		return nil
	}
	
	return modelHooks[hookType]
}

// Clear clears all hooks for a model (useful for testing)
func (hr *HookRegistry) Clear(modelName string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	
	delete(hr.hooks, modelName)
}

// ClearAll clears all hooks
func (hr *HookRegistry) ClearAll() {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	
	hr.hooks = make(map[string]map[HookType][]HookFunc)
}

