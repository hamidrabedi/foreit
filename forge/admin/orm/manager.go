package orm

import (
	"context"

	"github.com/forgego/forge/orm"
)

// AdminManager wraps orm.Manager for admin use
type AdminManager[T any] struct {
	manager *orm.Manager[T]
	schema  *orm.ModelSchema
}

// NewAdminManager creates a new admin manager wrapper
func NewAdminManager[T any](manager *orm.Manager[T]) (*AdminManager[T], error) {
	schema, err := orm.GetModelSchema[T]()
	if err != nil {
		return nil, err
	}

	return &AdminManager[T]{
		manager: manager,
		schema:  schema,
	}, nil
}

// GetQueryset returns a base queryset from the manager
func (am *AdminManager[T]) GetQueryset(ctx context.Context) (orm.QuerySet[T], error) {
	return am.manager.Filter(nil)
}

// GetFieldAccessor returns the field accessor for type-safe field operations
func (am *AdminManager[T]) GetFieldAccessor() (*orm.FieldAccessor[T], error) {
	return am.manager.GetFieldAccessor()
}

// GetSchema returns the model schema
func (am *AdminManager[T]) GetSchema() *orm.ModelSchema {
	return am.schema
}

// Get retrieves a model instance by ID
func (am *AdminManager[T]) Get(ctx context.Context, id int64) (*T, error) {
	return am.manager.Get(ctx, id)
}

// Create creates a new model instance
func (am *AdminManager[T]) Create(ctx context.Context, instance *T) error {
	return am.manager.Create(ctx, instance)
}

// Update updates an existing model instance
func (am *AdminManager[T]) Update(ctx context.Context, instance *T) error {
	return am.manager.Update(ctx, instance)
}

// Delete deletes a model instance
func (am *AdminManager[T]) Delete(ctx context.Context, instance *T) error {
	return am.manager.Delete(ctx, instance)
}

// All returns all model instances
func (am *AdminManager[T]) All(ctx context.Context) ([]*T, error) {
	return am.manager.All(ctx)
}

// Manager returns the underlying ORM manager
func (am *AdminManager[T]) Manager() *orm.Manager[T] {
	return am.manager
}
