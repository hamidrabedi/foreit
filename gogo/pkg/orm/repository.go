package orm

import (
	"context"
)

// Repository defines the interface for type-safe data access
type Repository[T any, Q any] interface {
	// Query returns a query builder for the model
	Query() Q
	
	// GetByID retrieves a single record by ID
	GetByID(ctx context.Context, id interface{}) (*T, error)
	
	// Get retrieves a single record matching the query
	Get(ctx context.Context, query Q) (*T, error)
	
	// All retrieves all records matching the query
	All(ctx context.Context, query Q) ([]*T, error)
	
	// Count returns the count of records matching the query
	Count(ctx context.Context, query Q) (int, error)
	
	// Create creates a new record
	Create(ctx context.Context, data *T) (*T, error)
	
	// Update updates an existing record
	Update(ctx context.Context, id interface{}, data *T) (*T, error)
	
	// Delete deletes a record by ID
	Delete(ctx context.Context, id interface{}) error
	
	// Exists checks if a record exists
	Exists(ctx context.Context, id interface{}) (bool, error)
}

// BaseRepository provides a base implementation that can be embedded
type BaseRepository[T any, Q any] struct {
	client interface{} // *ent.Client
	query  func(interface{}) Q
	create func(interface{}) interface{} // Returns create builder
	update func(interface{}, interface{}) interface{} // Returns update builder
}

// NewBaseRepository creates a new base repository
// This is a helper - concrete repositories should be generated or created per model
func NewBaseRepository[T any, Q any](
	client interface{},
	query func(interface{}) Q,
	create func(interface{}) interface{},
	update func(interface{}, interface{}) interface{},
) *BaseRepository[T, Q] {
	return &BaseRepository[T, Q]{
		client: client,
		query:  query,
		create: create,
		update: update,
	}
}

// Query returns a query builder
func (r *BaseRepository[T, Q]) Query() Q {
	return r.query(r.client)
}

// Transaction represents a database transaction
type Transaction interface {
	Commit() error
	Rollback() error
}

// TransactionalRepository extends Repository with transaction support
type TransactionalRepository[T any, Q any] interface {
	Repository[T, Q]
	
	// WithTransaction executes a function within a transaction
	WithTransaction(ctx context.Context, fn func(context.Context, Repository[T, Q]) error) error
}

