package models

import (
	"context"
	"fmt"
)

// ModelManager[T] provides Django-style model management
// This is the NEW manager system that works with ModelDefinition[T]
// Implements the Manager[T] interface
type ModelManager[T any] struct {
	model *ModelDefinition[T]
	db    *DB
}

// NewModelManager creates a new ModelManager for a model
func NewModelManager[T any](model *ModelDefinition[T], db *DB) *ModelManager[T] {
	return &ModelManager[T]{
		model: model,
		db:    db,
	}
}

// Objects returns a QuerySet for this model
// This is the Django-style API: Users.Objects().Filter(...)
func (m *ModelManager[T]) Objects() QuerySet[T] {
	return NewQuerySet[T](m.db, m.model)
}

// All returns a QuerySet for all instances
func (m *ModelManager[T]) All() QuerySet[T] {
	return m.Objects()
}

// Filter returns a QuerySet filtered by the given conditions
func (m *ModelManager[T]) Filter(conditions ...Condition) QuerySet[T] {
	return m.Objects().Filter(conditions...)
}

// FilterQ returns a QuerySet filtered by the given QueryExpr
func (m *ModelManager[T]) FilterQ(q *QueryExpr) QuerySet[T] {
	return m.Objects().FilterQ(q)
}

// First returns the first instance matching the conditions
func (m *ModelManager[T]) First(ctx context.Context, conditions ...Condition) (*T, error) {
	return m.Objects().Filter(conditions...).First(ctx)
}

// Count returns the count of instances matching the conditions
func (m *ModelManager[T]) Count(ctx context.Context, conditions ...Condition) (int, error) {
	return m.Objects().Filter(conditions...).Count(ctx)
}

// Exists checks if any instance matches the conditions
func (m *ModelManager[T]) Exists(ctx context.Context, conditions ...Condition) (bool, error) {
	return m.Objects().Filter(conditions...).Exists(ctx)
}

// Create creates a new instance
func (m *ModelManager[T]) Create(ctx context.Context, instance *T) error {
	// Run validation
	if err := ValidateModel(m.model, instance); err != nil {
		return err
	}

	// Run hooks
	if err := RunHooks(m.model, ctx, instance, "BeforeCreate"); err != nil {
		return err
	}

	// Insert into database
	_, err := m.db.NewInsert().Model(instance).Exec(ctx)
	if err != nil {
		return err
	}

	// Run after hooks
	if err := RunHooks(m.model, ctx, instance, "AfterCreate"); err != nil {
		return err
	}

	return nil
}

// Update updates an existing instance
func (m *ModelManager[T]) Update(ctx context.Context, instance *T) error {
	// Run validation
	if err := ValidateModel(m.model, instance); err != nil {
		return err
	}

	// Run hooks
	if err := RunHooks(m.model, ctx, instance, "BeforeUpdate"); err != nil {
		return err
	}

	// Update in database
	_, err := m.db.NewUpdate().Model(instance).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	// Run after hooks
	if err := RunHooks(m.model, ctx, instance, "AfterUpdate"); err != nil {
		return err
	}

	return nil
}

// Delete deletes an instance
func (m *ModelManager[T]) Delete(ctx context.Context, instance *T) error {
	// Run hooks
	if err := RunHooks(m.model, ctx, instance, "BeforeDelete"); err != nil {
		return err
	}

	// Delete from database
	_, err := m.db.NewDelete().Model(instance).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	// Run after hooks
	if err := RunHooks(m.model, ctx, instance, "AfterDelete"); err != nil {
		return err
	}

	return nil
}

// Get retrieves a single instance by conditions
func (m *ModelManager[T]) Get(ctx context.Context, conditions ...Condition) (*T, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("Get requires at least one condition")
	}
	return m.Objects().Filter(conditions...).Get(ctx)
}

// GetModel returns the model definition
func (m *ModelManager[T]) GetModel() *ModelDefinition[T] {
	return m.model
}

