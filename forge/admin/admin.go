package admin

import (
	"context"
	"reflect"

	query "github.com/forgego/forge/orm"
	"github.com/forgego/forge/orm"
)

// Admin represents a type-safe admin configuration for a model
type Admin[T any] struct {
	model   T
	manager *query.Manager[T]
	config  *Config[T]
	name    string
}

// Register registers a model with the admin system
func Register[T any](model T, manager *query.Manager[T], config *Config[T]) *Admin[T] {
	// Get model name from type
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	name := typ.Name()

	admin := &Admin[T]{
		model:   model,
		manager: manager,
		config:  config,
		name:    name,
	}

	// Register with global registry
	globalRegistry.register(admin)

	return admin
}

// ModelName returns the name of the model
func (a *Admin[T]) ModelName() string {
	return a.name
}

// Manager returns the manager for this admin
func (a *Admin[T]) Manager() *query.Manager[T] {
	return a.manager
}

// Config returns the configuration for this admin
func (a *Admin[T]) Config() *Config[T] {
	return a.config
}

// GetQueryset returns the base queryset for this admin
func (a *Admin[T]) GetQueryset(ctx context.Context) (query.QuerySet[T], error) {
	// Start with manager's base queryset
	qs, err := a.manager.Filter(nil)
	if err != nil {
		return nil, err
	}

	// Apply custom queryset hook if provided
	if a.config != nil && a.config.GetQueryset != nil {
		return a.config.GetQueryset(ctx, a, qs)
	}

	return qs, nil
}

// SaveModel saves a model instance
func (a *Admin[T]) SaveModel(ctx context.Context, instance *T, formData FormData, isNew bool) error {
	// Apply save hooks if provided
	if a.config != nil {
		if a.config.SaveModel != nil {
			return a.config.SaveModel(ctx, a, instance, formData, isNew)
		}
	}

	// Default save behavior
	if isNew {
		return a.manager.Create(ctx, instance)
	}
	return a.manager.Update(ctx, instance)
}

// DeleteModel deletes a model instance
func (a *Admin[T]) DeleteModel(ctx context.Context, instance *T) error {
	// Apply delete hooks if provided
	if a.config != nil && a.config.DeleteModel != nil {
		return a.config.DeleteModel(ctx, a, instance)
	}

	// Default delete behavior
	return a.manager.Delete(ctx, instance)
}

// HasAddPermission checks if the user has permission to add objects
func (a *Admin[T]) HasAddPermission(ctx context.Context, user interface{}) bool {
	if a.config != nil && a.config.HasAddPermission != nil {
		return a.config.HasAddPermission(ctx, a, user)
	}
	// Default: check with permission checker
	if a.config != nil && a.config.PermissionChecker != nil {
		return a.config.PermissionChecker.HasPermission(ctx, user, GetPermissionName(a.name, PermAdd))
	}
	return true // Default allow
}

// HasChangePermission checks if the user has permission to change objects
func (a *Admin[T]) HasChangePermission(ctx context.Context, user interface{}, obj *T) bool {
	if a.config != nil && a.config.HasChangePermission != nil {
		return a.config.HasChangePermission(ctx, a, user, obj)
	}
	// Default: check with permission checker
	if a.config != nil && a.config.PermissionChecker != nil {
		return a.config.PermissionChecker.HasPermission(ctx, user, GetPermissionName(a.name, PermChange))
	}
	return true // Default allow
}

// HasDeletePermission checks if the user has permission to delete objects
func (a *Admin[T]) HasDeletePermission(ctx context.Context, user interface{}, obj *T) bool {
	if a.config != nil && a.config.HasDeletePermission != nil {
		return a.config.HasDeletePermission(ctx, a, user, obj)
	}
	// Default: check with permission checker
	if a.config != nil && a.config.PermissionChecker != nil {
		return a.config.PermissionChecker.HasPermission(ctx, user, GetPermissionName(a.name, PermDelete))
	}
	return true // Default allow
}

// HasViewPermission checks if the user has permission to view objects
func (a *Admin[T]) HasViewPermission(ctx context.Context, user interface{}, obj *T) bool {
	if a.config != nil && a.config.HasViewPermission != nil {
		return a.config.HasViewPermission(ctx, a, user, obj)
	}
	// Default: check with permission checker
	if a.config != nil && a.config.PermissionChecker != nil {
		return a.config.PermissionChecker.HasPermission(ctx, user, GetPermissionName(a.name, PermView))
	}
	return true // Default allow
}

// HasModulePermission checks if the user has permission to access this module
func (a *Admin[T]) HasModulePermission(ctx context.Context, user interface{}) bool {
	if a.config != nil && a.config.HasModulePermission != nil {
		return a.config.HasModulePermission(ctx, a, user)
	}
	// Default: check with permission checker
	if a.config != nil && a.config.PermissionChecker != nil {
		return a.config.PermissionChecker.HasPermission(ctx, user, GetPermissionName(a.name, PermView))
	}
	return true // Default allow
}
