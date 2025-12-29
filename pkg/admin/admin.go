package admin

import (
	"context"
	"reflect"

	"github.com/forgego/forge/pkg/query"
	"github.com/forgego/forge/pkg/schema"
)

// Admin is a type-safe admin instance for a model
type Admin[T any] struct {
	model    T
	manager  *query.Manager[T]
	config   *Config[T]
	registry *Registry
	name     string
}

// Config provides type-safe configuration for Admin
type Config[T any] struct {
	// List view configuration
	ListDisplay   []FieldExpr[T, interface{}]
	ListFilter    []Filter[T]
	SearchFields  []FieldExpr[T, interface{}]
	DateHierarchy FieldExpr[T, interface{}]
	Ordering      []Ordering[T]
	ListPerPage   int

	// Form configuration
	Fieldsets        []Fieldset[T]
	ReadOnlyFields   []FieldExpr[T, interface{}]
	AutocompleteFields []FieldExpr[T, interface{}]
	RawIDFields      []FieldExpr[T, interface{}]

	// Actions
	Actions []Action[T]

	// Inlines
	Inlines []Inline[T, interface{}]

	// Customization
	VerboseName       string
	VerboseNamePlural string
	SaveOnTop         bool
	SaveAs            bool
	SaveAsContinue    bool
	SaveAndAddAnother bool

	// Custom methods (type-safe)
	GetQueryset func(ctx context.Context, manager *query.Manager[T]) (query.QuerySet[T], error)
	SaveModel   func(ctx context.Context, instance *T, form FormData, isNew bool) error
	DeleteModel func(ctx context.Context, instance *T) error
	GetForm     func(ctx context.Context, instance *T, isNew bool) (Form[T], error)
}

// Register registers a model with the admin system
func Register[T any](
	model T,
	manager *query.Manager[T],
	config *Config[T],
) *Admin[T] {
	// Get model name
	var zero T
	modelType := reflect.TypeOf(zero)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	modelName := modelType.Name()

	// Create default config if nil
	if config == nil {
		config = &Config[T]{}
	}

	// Set defaults
	if config.ListPerPage == 0 {
		config.ListPerPage = 20
	}
	if config.VerboseName == "" {
		config.VerboseName = modelName
	}
	if config.VerboseNamePlural == "" {
		config.VerboseNamePlural = modelName + "s"
	}

	// Auto-generate config from schema if available
	if schemaModel, ok := any(model).(schema.Schema); ok {
		autoGenerateConfig(config, schemaModel)
	}

	admin := &Admin[T]{
		model:    model,
		manager:  manager,
		config:   config,
		registry: globalRegistry,
		name:     modelName,
	}

	// Register with global registry
	globalRegistry.register(admin)

	// Also register with type registry for HTTP handlers
	// This will be done when HTTP router is initialized

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

// Model returns the model instance
func (a *Admin[T]) Model() T {
	return a.model
}

// autoGenerateConfig auto-generates admin config from schema
func autoGenerateConfig[T any](config *Config[T], schemaModel schema.Schema) {
	fields := schemaModel.Fields()
	meta := schemaModel.Meta()

	// Auto-generate list display if not set
	if len(config.ListDisplay) == 0 {
		config.ListDisplay = autoGenerateListDisplay[T](fields, meta)
	}

	// Auto-generate list filter if not set
	if len(config.ListFilter) == 0 {
		config.ListFilter = autoGenerateListFilter[T](fields)
	}

	// Auto-generate search fields if not set
	if len(config.SearchFields) == 0 {
		config.SearchFields = autoGenerateSearchFields[T](fields)
	}

	// Auto-generate read-only fields if not set
	if len(config.ReadOnlyFields) == 0 {
		config.ReadOnlyFields = autoGenerateReadOnlyFields[T](fields)
	}
}

// autoGenerateListDisplay generates list display fields from schema
func autoGenerateListDisplay[T any](fields []schema.Field, meta schema.Meta) []FieldExpr[T, interface{}] {
	var listDisplay []FieldExpr[T, interface{}]

	// Use Meta.OrderBy if available
	if len(meta.OrderBy) > 0 {
		for _, fieldName := range meta.OrderBy {
			// Remove leading dash for descending order
			if len(fieldName) > 0 && fieldName[0] == '-' {
				fieldName = fieldName[1:]
			}
			// Create field expression (simplified - will be enhanced)
			// For now, we'll create a basic field expr
		}
	} else {
		// Default: use first 5 non-primary key fields
		count := 0
		for _, field := range fields {
			if !field.PrimaryKey && count < 5 {
				// Create field expression
				// This will be implemented when FieldExpr is created
				count++
			}
		}
	}

	return listDisplay
}

// autoGenerateListFilter generates list filter fields from schema
func autoGenerateListFilter[T any](fields []schema.Field) []Filter[T] {
	var filters []Filter[T]

	for _, field := range fields {
		// Add boolean fields
		if field.Type == schema.TypeBool {
			// Create boolean filter (will be implemented)
		}
	}

	return filters
}

// autoGenerateSearchFields generates search fields from schema
func autoGenerateSearchFields[T any](fields []schema.Field) []FieldExpr[T, interface{}] {
	var searchFields []FieldExpr[T, interface{}]

	for _, field := range fields {
		// Add string and text fields
		if field.Type == schema.TypeString || field.Type == schema.TypeText {
			// Create field expression (will be implemented)
		}
	}

	return searchFields
}

// autoGenerateReadOnlyFields generates read-only fields from schema
func autoGenerateReadOnlyFields[T any](fields []schema.Field) []FieldExpr[T, interface{}] {
	var readOnly []FieldExpr[T, interface{}]

	for _, field := range fields {
		// Primary keys are read-only
		if field.PrimaryKey {
			// Create field expression (will be implemented)
		}
		// Auto-now and auto-now-add fields are read-only
		if field.AutoNow || field.AutoNowAdd {
			// Create field expression (will be implemented)
		}
	}

	return readOnly
}

// GetQueryset returns the queryset for this admin
func (a *Admin[T]) GetQueryset(ctx context.Context) (query.QuerySet[T], error) {
	if a.config.GetQueryset != nil {
		return a.config.GetQueryset(ctx, a.manager)
	}
	return a.manager.Filter(), nil
}

// SaveModel saves a model instance
func (a *Admin[T]) SaveModel(ctx context.Context, instance *T, form FormData, isNew bool) error {
	if a.config.SaveModel != nil {
		return a.config.SaveModel(ctx, instance, form, isNew)
	}

	if isNew {
		return a.manager.Create(ctx, instance)
	}
	return a.manager.Update(ctx, instance)
}

// DeleteModel deletes a model instance
func (a *Admin[T]) DeleteModel(ctx context.Context, instance *T) error {
	if a.config.DeleteModel != nil {
		return a.config.DeleteModel(ctx, instance)
	}
	return a.manager.Delete(ctx, instance)
}

// FormData represents form data
type FormData map[string]interface{}
