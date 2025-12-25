package orm

import (
	"context"
)

// BulkOperations provides bulk database operations
type BulkOperations[T any] struct {
	repo Repository[T, interface{}]
}

// NewBulkOperations creates a new bulk operations helper
func NewBulkOperations[T any](repo Repository[T, interface{}]) *BulkOperations[T] {
	return &BulkOperations[T]{
		repo: repo,
	}
}

// BulkCreate creates multiple records
func (b *BulkOperations[T]) BulkCreate(ctx context.Context, items []*T) ([]*T, error) {
	// This should use Ent's bulk create when available
	// For now, create one by one
	results := make([]*T, 0, len(items))
	for _, item := range items {
		result, err := b.repo.Create(ctx, item)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// BulkUpdate updates multiple records
func (b *BulkOperations[T]) BulkUpdate(ctx context.Context, updates map[interface{}]*T) ([]*T, error) {
	results := make([]*T, 0, len(updates))
	for id, data := range updates {
		result, err := b.repo.Update(ctx, id, data)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// BulkDelete deletes multiple records
func (b *BulkOperations[T]) BulkDelete(ctx context.Context, ids []interface{}) error {
	for _, id := range ids {
		if err := b.repo.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// BulkUpsert creates or updates multiple records
func (b *BulkOperations[T]) BulkUpsert(ctx context.Context, items []*T, keyFunc func(*T) interface{}) ([]*T, error) {
	results := make([]*T, 0, len(items))
	for _, item := range items {
		key := keyFunc(item)
		exists, _ := b.repo.Exists(ctx, key)
		
		var result *T
		var err error
		if exists {
			result, err = b.repo.Update(ctx, key, item)
		} else {
			result, err = b.repo.Create(ctx, item)
		}
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

