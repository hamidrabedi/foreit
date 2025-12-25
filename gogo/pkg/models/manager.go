package models

import (
	"context"
	"fmt"
)

// DefaultManager provides a default manager implementation
type DefaultManager struct {
	modelType string
	repo      Repository
}

// Repository interface for database operations
type Repository interface {
	Save(ctx context.Context, model Model) error
	Delete(ctx context.Context, model Model) error
	GetByID(ctx context.Context, id interface{}) (Model, error)
	All(ctx context.Context) ([]Model, error)
	Query(ctx context.Context, qs *QuerySetImpl) ([]Model, error)
	Count(ctx context.Context, qs *QuerySetImpl) (int, error)
}

// NewDefaultManager creates a new default manager
func NewDefaultManager(modelType string, repo Repository) *DefaultManager {
	return &DefaultManager{
		modelType: modelType,
		repo:      repo,
	}
}

// Save saves a model
func (m *DefaultManager) Save(ctx context.Context, model Model) error {
	return m.repo.Save(ctx, model)
}

// Delete deletes a model
func (m *DefaultManager) Delete(ctx context.Context, model Model) error {
	return m.repo.Delete(ctx, model)
}

// Get retrieves a model by ID
func (m *DefaultManager) Get(ctx context.Context, id interface{}) (Model, error) {
	return m.repo.GetByID(ctx, id)
}

// All returns all models
func (m *DefaultManager) All(ctx context.Context) ([]Model, error) {
	return m.repo.All(ctx)
}

// Filter returns a queryset
func (m *DefaultManager) Filter(ctx context.Context) QuerySet {
	return NewQuerySet(m.modelType, m)
}

// QuerySetAll implements QuerySetManager
func (m *DefaultManager) QuerySetAll(ctx context.Context, qs *QuerySetImpl) ([]Model, error) {
	return m.repo.Query(ctx, qs)
}

// QuerySetCount implements QuerySetManager
func (m *DefaultManager) QuerySetCount(ctx context.Context, qs *QuerySetImpl) (int, error) {
	return m.repo.Count(ctx, qs)
}

// ManagerRegistry stores managers for different model types
type ManagerRegistry struct {
	managers map[string]Manager
}

var globalManagerRegistry = &ManagerRegistry{
	managers: make(map[string]Manager),
}

// RegisterManager registers a manager for a model type
func RegisterManager(modelType string, manager Manager) {
	globalManagerRegistry.managers[modelType] = manager
}

// GetManager retrieves a manager for a model type
func GetManager(modelType string) (Manager, error) {
	manager, ok := globalManagerRegistry.managers[modelType]
	if !ok {
		return nil, fmt.Errorf("no manager registered for model type: %s", modelType)
	}
	return manager, nil
}

