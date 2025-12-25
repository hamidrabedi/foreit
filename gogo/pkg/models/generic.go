package models

import (
	"context"
)

// Model is the base interface all models must implement
type Model interface {
	GetID() interface{}
	SetID(id interface{})
	IsNew() bool
	String() string
}

// Manager is a generic manager for type-safe operations
type Manager[T Model] interface {
	// Save saves a model (create or update)
	Save(ctx context.Context, model *T) error
	
	// Delete deletes a model
	Delete(ctx context.Context, model *T) error
	
	// Get retrieves a model by ID
	Get(ctx context.Context, id interface{}) (*T, error)
	
	// All returns all models
	All(ctx context.Context) ([]*T, error)
	
	// Filter returns a type-safe queryset
	Filter(ctx context.Context) QuerySet[T]
	
	// Create creates a new model
	Create(ctx context.Context, model *T) error
	
	// Update updates an existing model
	Update(ctx context.Context, model *T) error
}

// QuerySet is a generic, immutable queryset
type QuerySet[T Model] interface {
	// Filter adds a filter condition (returns new QuerySet)
	Filter(condition Q) QuerySet[T]
	
	// Exclude adds an exclude condition (returns new QuerySet)
	Exclude(condition Q) QuerySet[T]
	
	// OrderBy adds ordering (returns new QuerySet)
	OrderBy(fields ...string) QuerySet[T]
	
	// Limit limits results (returns new QuerySet)
	Limit(n int) QuerySet[T]
	
	// Offset sets offset (returns new QuerySet)
	Offset(n int) QuerySet[T]
	
	// Get retrieves a single result
	Get(ctx context.Context) (*T, error)
	
	// All retrieves all results
	All(ctx context.Context) ([]*T, error)
	
	// Count returns the count
	Count(ctx context.Context) (int, error)
	
	// Exists checks if any results exist
	Exists(ctx context.Context) (bool, error)
	
	// First returns the first result
	First(ctx context.Context) (*T, error)
	
	// Last returns the last result
	Last(ctx context.Context) (*T, error)
}

// BaseManager provides a generic base implementation
type BaseManager[T Model] struct {
	repo Repository[T]
	meta *ModelMeta
}

// NewBaseManager creates a new generic manager
func NewBaseManager[T Model](repo Repository[T], meta *ModelMeta) *BaseManager[T] {
	return &BaseManager[T]{
		repo: repo,
		meta: meta,
	}
}

// Save saves a model
func (m *BaseManager[T]) Save(ctx context.Context, model *T) error {
	if (*model).IsNew() {
		return m.Create(ctx, model)
	}
	return m.Update(ctx, model)
}

// Create creates a new model
func (m *BaseManager[T]) Create(ctx context.Context, model *T) error {
	return m.repo.Create(ctx, model)
}

// Update updates an existing model
func (m *BaseManager[T]) Update(ctx context.Context, model *T) error {
	return m.repo.Update(ctx, model)
}

// Delete deletes a model
func (m *BaseManager[T]) Delete(ctx context.Context, model *T) error {
	return m.repo.Delete(ctx, model)
}

// Get retrieves a model by ID
func (m *BaseManager[T]) Get(ctx context.Context, id interface{}) (*T, error) {
	return m.repo.GetByID(ctx, id)
}

// All returns all models
func (m *BaseManager[T]) All(ctx context.Context) ([]*T, error) {
	return m.repo.All(ctx)
}

// Filter returns a queryset
func (m *BaseManager[T]) Filter(ctx context.Context) QuerySet[T] {
	return NewQuerySetImpl[T](m.repo, m.meta)
}

// Repository is a generic repository interface
type Repository[T Model] interface {
	Create(ctx context.Context, model *T) error
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, model *T) error
	GetByID(ctx context.Context, id interface{}) (*T, error)
	All(ctx context.Context) ([]*T, error)
	Query(ctx context.Context, qs *QuerySetImpl[T]) ([]*T, error)
	Count(ctx context.Context, qs *QuerySetImpl[T]) (int, error)
}

// Ensure QuerySetImpl implements QuerySet
var _ QuerySet[Model] = (*QuerySetImpl[Model])(nil)

