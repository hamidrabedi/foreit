package admin

import (
	"fmt"
	"sync"
)

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

type HookFunc func(ctx *HookContext) error

type HookContext struct {
	Model     *ModelMeta
	Action    string
	User      interface{}
	Request   interface{}
	Data      interface{}
	Resource  interface{}
	Query     interface{}
	Result    interface{}
	Error     error
}

type HookRegistry struct {
	hooks map[string]map[HookType][]HookFunc
	mu    sync.RWMutex
}

func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[string]map[HookType][]HookFunc),
	}
}

func (hr *HookRegistry) Register(modelName string, hookType HookType, hook HookFunc) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	
	if hr.hooks[modelName] == nil {
		hr.hooks[modelName] = make(map[HookType][]HookFunc)
	}
	
	hr.hooks[modelName][hookType] = append(hr.hooks[modelName][hookType], hook)
}

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
	
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return fmt.Errorf("hook execution failed: %w", err)
		}
	}
	
	return nil
}

func (hr *HookRegistry) GetHooks(modelName string, hookType HookType) []HookFunc {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	
	modelHooks, exists := hr.hooks[modelName]
	if !exists {
		return nil
	}
	
	return modelHooks[hookType]
}

func (hr *HookRegistry) Clear(modelName string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	
	delete(hr.hooks, modelName)
}

func (hr *HookRegistry) ClearAll() {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	
	hr.hooks = make(map[string]map[HookType][]HookFunc)
}

