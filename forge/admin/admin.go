package admin

import (
	"context"
	"fmt"
	"reflect"

	adminfilter "github.com/forgego/forge/admin/filter"
	adminorm "github.com/forgego/forge/admin/orm"
	adminschema "github.com/forgego/forge/admin/schema"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// Admin represents a type-safe admin configuration for a model
// It fully integrates with schema, ORM, and filter systems
type Admin[T any] struct {
	// Core dependencies
	schema      schema.Schema
	manager     *adminorm.AdminManager[T]
	filterset   *adminfilter.AdminFilterSet[T]
	config      *Config[T]
	modelSchema *orm.ModelSchema

	// Discovered information
	fields    []adminschema.FieldInfo
	relations []adminschema.RelationInfo
	meta      adminschema.MetaInfo

	// Metadata
	name string
}

// Register registers a model with the admin system
// It automatically discovers fields, relations, and metadata from the schema
func Register[T any](
	schemaInstance schema.Schema,
	manager *orm.Manager[T],
	config *Config[T],
) (*Admin[T], error) {
	// Get model schema from ORM
	modelSchema, err := orm.GetModelSchema[T]()
	if err != nil {
		return nil, fmt.Errorf("failed to get model schema: %w", err)
	}

	// Create admin manager wrapper
	adminManager, err := adminorm.NewAdminManager(manager)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin manager: %w", err)
	}

	// Create filterset
	adminFilterSet, err := adminfilter.NewAdminFilterSet[T]()
	if err != nil {
		return nil, fmt.Errorf("failed to create filterset: %w", err)
	}

	// Discover fields from schema
	fields, err := adminschema.DiscoverFields[T](schemaInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to discover fields: %w", err)
	}

	// Discover relations from schema
	relations, err := adminschema.DiscoverRelations[T](schemaInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to discover relations: %w", err)
	}

	// Discover metadata
	meta := adminschema.DiscoverMeta(schemaInstance)

	// Build filters from schema
	if err := adminfilter.BuildFiltersFromSchema[T](schemaInstance, adminFilterSet.FilterSet()); err != nil {
		return nil, fmt.Errorf("failed to build filters: %w", err)
	}

	// Get model name
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	name := typ.Name()
	if meta.TableName != "" {
		name = meta.TableName
	}

	admin := &Admin[T]{
		schema:      schemaInstance,
		manager:     adminManager,
		filterset:   adminFilterSet,
		config:      config,
		modelSchema: modelSchema,
		fields:      fields,
		relations:   relations,
		meta:        meta,
		name:        name,
	}

	// Auto-configure from schema if config is minimal
	if config != nil {
		admin.autoConfigureFromSchema()
	}

	// Register with global registry
	globalRegistry.register(admin)

	return admin, nil
}

// autoConfigureFromSchema auto-configures admin from discovered schema information
func (a *Admin[T]) autoConfigureFromSchema() {
	if a.config == nil {
		return
	}

	// Auto-configure list display if not set
	if len(a.config.ListDisplay) == 0 {
		fieldMapper := adminschema.NewFieldMapper()
		for _, field := range a.fields {
			if fieldMapper.ShouldDisplayInList(field.SchemaField) {
				a.config.ListDisplay = append(a.config.ListDisplay, field.Name)
			}
		}
	}

	// Auto-configure search fields if not set
	if len(a.config.SearchFields) == 0 {
		// Use string fields for search
		for _, field := range a.fields {
			if field.Type == schema.TypeString || field.Type == schema.TypeText || field.Type == schema.TypeEmail {
				a.config.SearchFields = append(a.config.SearchFields, field.Name)
			}
		}
	}
}

// ModelName returns the name of the model
func (a *Admin[T]) ModelName() string {
	return a.name
}

// Manager returns the admin manager
func (a *Admin[T]) Manager() *adminorm.AdminManager[T] {
	return a.manager
}

// FilterSet returns the admin filterset
func (a *Admin[T]) FilterSet() *adminfilter.AdminFilterSet[T] {
	return a.filterset
}

// Config returns the configuration
func (a *Admin[T]) Config() *Config[T] {
	return a.config
}

// Schema returns the schema instance
func (a *Admin[T]) Schema() schema.Schema {
	return a.schema
}

// ModelSchema returns the ORM model schema
func (a *Admin[T]) ModelSchema() *orm.ModelSchema {
	return a.modelSchema
}

// Fields returns discovered field information
func (a *Admin[T]) Fields() []adminschema.FieldInfo {
	return a.fields
}

// Relations returns discovered relation information
func (a *Admin[T]) Relations() []adminschema.RelationInfo {
	return a.relations
}

// Meta returns model metadata
func (a *Admin[T]) Meta() adminschema.MetaInfo {
	return a.meta
}

// GetQueryset returns the base queryset for this admin
func (a *Admin[T]) GetQueryset(ctx context.Context) (orm.QuerySet[T], error) {
	// Start with manager's base queryset
	qs, err := a.manager.GetQueryset(ctx)
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

// ModelType returns the model type for Admin[T]
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
	return a.manager.Manager()
}

// ConfigInterface returns the config as interface{} for dynamic access
func (a *Admin[T]) ConfigInterface() interface{} {
	return a.config
}
