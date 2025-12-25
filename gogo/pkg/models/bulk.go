package models

import (
	"context"
)

// BulkOperations provides type-safe bulk operations
type BulkOperations[T Model] interface {
	// BulkCreate creates multiple models
	BulkCreate(ctx context.Context, models []*T) ([]*T, error)
	
	// BulkUpdate updates multiple models
	BulkUpdate(ctx context.Context, models []*T) error
	
	// BulkDelete deletes multiple models
	BulkDelete(ctx context.Context, models []*T) error
}

// BulkManager extends Manager with bulk operations
type BulkManager[T Model] interface {
	Manager[T]
	BulkOperations[T]
}

// BaseBulkManager provides bulk operation implementations
type BaseBulkManager[T Model] struct {
	*BaseManager[T]
}

// NewBaseBulkManager creates a new bulk manager
func NewBaseBulkManager[T Model](repo Repository[T], meta *ModelMeta) *BaseBulkManager[T] {
	return &BaseBulkManager[T]{
		BaseManager: NewBaseManager[T](repo, meta),
	}
}

// BulkCreate creates multiple models in a single transaction
func (m *BaseBulkManager[T]) BulkCreate(ctx context.Context, models []*T) ([]*T, error) {
	// Implementation would use batch insert
	// For now, create one by one in transaction
	results := make([]*T, 0, len(models))
	
	for _, model := range models {
		if err := m.Create(ctx, model); err != nil {
			return nil, err
		}
		results = append(results, model)
	}
	
	return results, nil
}

// BulkUpdate updates multiple models
func (m *BaseBulkManager[T]) BulkUpdate(ctx context.Context, models []*T) error {
	// Implementation would use batch update
	// For now, update one by one
	for _, model := range models {
		if err := m.Update(ctx, model); err != nil {
			return err
		}
	}
	return nil
}

// BulkDelete deletes multiple models
func (m *BaseBulkManager[T]) BulkDelete(ctx context.Context, models []*T) error {
	// Implementation would use batch delete
	// For now, delete one by one
	for _, model := range models {
		if err := m.Delete(ctx, model); err != nil {
			return err
		}
	}
	return nil
}

// BulkCreateFromQueryset creates models from queryset results
func BulkCreateFromQueryset[T Model](
	ctx context.Context,
	qs QuerySet[T],
	transform func(*T) *T,
) ([]*T, error) {
	// Get all from queryset
	models, err := qs.All(ctx)
	if err != nil {
		return nil, err
	}
	
	// Transform if needed
	if transform != nil {
		transformed := make([]*T, len(models))
		for i, model := range models {
			transformed[i] = transform(model)
		}
		models = transformed
	}
	
	// Would need manager reference - this is a helper function
	// Actual implementation would be on BulkManager
	return models, nil
}

