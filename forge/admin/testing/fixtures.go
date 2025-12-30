package testing

import (
	"context"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
)

// FixtureManager manages test fixtures
type FixtureManager[T any] struct {
	admin   *admin.Admin[T]
	manager *orm.Manager[T]
}

// NewFixtureManager creates a new fixture manager
func NewFixtureManager[T any](admin *admin.Admin[T], manager *orm.Manager[T]) *FixtureManager[T] {
	return &FixtureManager[T]{
		admin:   admin,
		manager: manager,
	}
}

// CreateFixture creates a test fixture
func (fm *FixtureManager[T]) CreateFixture(ctx context.Context, instance *T) error {
	return fm.manager.Create(ctx, instance)
}

// CreateFixtures creates multiple test fixtures
func (fm *FixtureManager[T]) CreateFixtures(ctx context.Context, instances []*T) error {
	for _, instance := range instances {
		if err := fm.CreateFixture(ctx, instance); err != nil {
			return err
		}
	}
	return nil
}

// CleanupFixtures cleans up test fixtures
func (fm *FixtureManager[T]) CleanupFixtures(ctx context.Context, instances []*T) error {
	for _, instance := range instances {
		if err := fm.manager.Delete(ctx, instance); err != nil {
			return err
		}
	}
	return nil
}

// GetFixture gets a fixture by ID
func (fm *FixtureManager[T]) GetFixture(ctx context.Context, id int64) (*T, error) {
	return fm.manager.Get(ctx, id)
}
