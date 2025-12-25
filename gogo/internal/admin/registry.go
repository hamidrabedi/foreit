package admin

import (
	"fmt"
	"reflect"
	"sync"

	"entgo.io/ent"
)

// Registry is the central registry for all admin models (like Django's admin.site)
type Registry struct {
	models      map[string]*ModelMeta
	mu          sync.RWMutex
	client      interface{} // *ent.Client - using interface{} to avoid circular dependency
	hookRegistry *HookRegistry
}

// GetHookRegistry returns the hook registry
func (r *Registry) GetHookRegistry() *HookRegistry {
	return r.hookRegistry
}

// NewRegistry creates a new admin registry
func NewRegistry() *Registry {
	return &Registry{
		models:       make(map[string]*ModelMeta),
		hookRegistry: NewHookRegistry(),
	}
}

// SetClient sets the Ent client for the registry
func (r *Registry) SetClient(client interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.client = client
}

// Register registers a model with the admin registry
func (r *Registry) Register(model interface{}, options ...Option) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	modelName := modelType.Name()
	if _, exists := r.models[modelName]; exists {
		return fmt.Errorf("model %s is already registered", modelName)
	}

	// Create model metadata
	meta := &ModelMeta{
		Name:      modelName,
		ModelType: modelType,
		Options:   &ModelOptions{},
		Permissions: &Permissions{
			CanList:   true,
			CanView:   true,
			CanCreate: true,
			CanUpdate: true,
			CanDelete: true,
			Rules:     make(map[string]string),
		},
	}

	// Apply options
	for _, opt := range options {
		opt(meta)
	}

	r.models[modelName] = meta
	return nil
}

// GetModel retrieves model metadata by name
func (r *Registry) GetModel(name string) (*ModelMeta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, exists := r.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s is not registered", name)
	}

	return meta, nil
}

// GetAllModels returns all registered models
func (r *Registry) GetAllModels() map[string]*ModelMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]*ModelMeta)
	for k, v := range r.models {
		result[k] = v
	}
	return result
}

// ModelMeta contains metadata about a registered model
type ModelMeta struct {
	Name         string
	ModelType    reflect.Type
	Schema       ent.Schema
	TableName    string
	Fields       []FieldMeta
	Relationships []RelationshipMeta
	Permissions  *Permissions
	Options      *ModelOptions
}

// Permissions defines what actions are allowed on a model
type Permissions struct {
	CanList   bool
	CanView   bool
	CanCreate bool
	CanUpdate bool
	CanDelete bool
	Rules     map[string]string // Rule-based permissions (like PocketBase)
}

// ModelOptions contains configuration options for a model
type ModelOptions struct {
	ExcludeFields   []string
	ReadOnlyFields  []string
	FilterableFields []string
	SortableFields  []string
	SearchFields    []string
	CustomActions   []Action
}

// FieldMeta contains metadata about a model field
type FieldMeta struct {
	Name        string
	Label       string
	Type        FieldType
	DBType      string
	Required    bool
	Unique      bool
	Default     interface{}
	Choices     []Choice
	HelpText    string
	ReadOnly    bool
	Filterable  bool
	Sortable    bool
	Searchable  bool
	Relationship *RelationshipMeta // If this field is a foreign key
}

// RelationshipMeta contains metadata about a relationship (edge)
type RelationshipMeta struct {
	Name         string // Edge name
	Type         RelationshipType // OneToOne, OneToMany, ManyToOne, ManyToMany
	TargetModel  string // Target model name
	Inverse      bool   // Is this the inverse side of the relationship
	Required     bool   // Is the relationship required
	Unique       bool   // Is the relationship unique (OneToOne)
	ForeignKey   string // Foreign key field name (if applicable)
}

// FieldType represents the type of an admin field
type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeEmail    FieldType = "email"
	FieldTypeURL      FieldType = "url"
	FieldTypeBoolean  FieldType = "boolean"
	FieldTypeDate     FieldType = "date"
	FieldTypeDateTime FieldType = "datetime"
	FieldTypeTime     FieldType = "time"
	FieldTypeSelect   FieldType = "select"
	FieldTypeMultiSelect FieldType = "multiselect"
	FieldTypeTextarea FieldType = "textarea"
	FieldTypeJSON     FieldType = "json"
	FieldTypeFile     FieldType = "file"
	FieldTypeImage    FieldType = "image"
	FieldTypeRelation FieldType = "relation"
)

// RelationshipType represents the type of a relationship
type RelationshipType string

const (
	RelationshipTypeOneToOne   RelationshipType = "one_to_one"
	RelationshipTypeOneToMany  RelationshipType = "one_to_many"
	RelationshipTypeManyToOne  RelationshipType = "many_to_one"
	RelationshipTypeManyToMany RelationshipType = "many_to_many"
)

// Choice represents a select option
type Choice struct {
	Value string
	Label string
}

// Action represents a custom action that can be performed on a model
type Action struct {
	Name        string
	Label       string
	Handler     func(ids []string) error
	Description string
}

// Option is a function that configures a model during registration
type Option func(*ModelMeta)

// WithPermissions sets permissions for a model
func WithPermissions(perms Permissions) Option {
	return func(meta *ModelMeta) {
		meta.Permissions = &perms
	}
}

// WithFields configures field options
func WithFields(fields FieldsConfig) Option {
	return func(meta *ModelMeta) {
		if len(fields.Exclude) > 0 {
			meta.Options.ExcludeFields = fields.Exclude
		}
		if len(fields.ReadOnly) > 0 {
			meta.Options.ReadOnlyFields = fields.ReadOnly
		}
	}
}

// FieldsConfig contains field configuration options
type FieldsConfig struct {
	Exclude  []string
	ReadOnly []string
}

// WithFilters sets filterable fields
func WithFilters(fields []string) Option {
	return func(meta *ModelMeta) {
		meta.Options.FilterableFields = fields
	}
}

// WithSorting sets sortable fields
func WithSorting(fields []string) Option {
	return func(meta *ModelMeta) {
		meta.Options.SortableFields = fields
	}
}

// WithSearch sets searchable fields
func WithSearch(fields []string) Option {
	return func(meta *ModelMeta) {
		meta.Options.SearchFields = fields
	}
}

// WithTableName sets a custom table name
func WithTableName(name string) Option {
	return func(meta *ModelMeta) {
		meta.TableName = name
	}
}

// WithRule adds a rule-based permission (like PocketBase)
func WithRule(action string, rule string) Option {
	return func(meta *ModelMeta) {
		if meta.Permissions.Rules == nil {
			meta.Permissions.Rules = make(map[string]string)
		}
		meta.Permissions.Rules[action] = rule
	}
}

// WithActions adds custom actions
func WithActions(actions []Action) Option {
	return func(meta *ModelMeta) {
		meta.Options.CustomActions = actions
	}
}

