package models

import (
	"context"
	"gorm.io/gorm"
)

// ManagerImpl provides a type-safe manager implementation
type ManagerImpl[T any] struct {
	db    *gorm.DB
	model T
}

// NewManager creates a new manager
func NewManager[T any](db *gorm.DB, model T) *ManagerImpl[T] {
	return &ManagerImpl[T]{
		db:    db,
		model: model,
	}
}

// Save saves a model (create or update)
func (m *ManagerImpl[T]) Save(ctx context.Context, model *T) error {
	return m.db.WithContext(ctx).Save(model).Error
}

// Create creates a new model
func (m *ManagerImpl[T]) Create(ctx context.Context, model *T) error {
	return m.db.WithContext(ctx).Create(model).Error
}

// Update updates an existing model
func (m *ManagerImpl[T]) Update(ctx context.Context, model *T) error {
	return m.db.WithContext(ctx).Save(model).Error
}

// Delete deletes a model
func (m *ManagerImpl[T]) Delete(ctx context.Context, model *T) error {
	return m.db.WithContext(ctx).Delete(model).Error
}

// Get retrieves a model by ID
func (m *ManagerImpl[T]) Get(ctx context.Context, id interface{}) (*T, error) {
	var result T
	err := m.db.WithContext(ctx).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// All returns all models
func (m *ManagerImpl[T]) All(ctx context.Context) ([]*T, error) {
	var results []*T
	err := m.db.WithContext(ctx).Find(&results).Error
	return results, err
}

// Filter returns a type-safe queryset
func (m *ManagerImpl[T]) Filter(ctx context.Context) *QuerySetImpl[T] {
	return NewQuerySet(m.db.WithContext(ctx), m.model)
}

