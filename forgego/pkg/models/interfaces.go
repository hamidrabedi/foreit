package models

import "context"

// QuerySet defines the interface for type-safe database queries.
type QuerySet[T any] interface {
	Filter(conditions ...Condition) QuerySet[T]
	FilterQ(q *QueryExpr) QuerySet[T]
	Exclude(conditions ...Condition) QuerySet[T]
	OrderBy(orders ...OrderBy) QuerySet[T]
	Limit(limit int) QuerySet[T]
	Offset(offset int) QuerySet[T]
	Distinct() QuerySet[T]
	Select(fields ...string) QuerySet[T]
	SelectRelated(relations ...string) QuerySet[T]
	PrefetchRelated(relations ...string) QuerySet[T]
	All(ctx context.Context) ([]*T, error)
	Get(ctx context.Context) (*T, error)
	First(ctx context.Context) (*T, error)
	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context) (bool, error)
	Delete(ctx context.Context) (int, error)
	Update(ctx context.Context, values map[string]interface{}) (int, error)
}

// Manager defines the interface for model management operations.
type Manager[T any] interface {
	All() QuerySet[T]
	Filter(conditions ...Condition) QuerySet[T]
	FilterQ(q *QueryExpr) QuerySet[T]
	Get(ctx context.Context, conditions ...Condition) (*T, error)
	First(ctx context.Context, conditions ...Condition) (*T, error)
	Count(ctx context.Context, conditions ...Condition) (int, error)
	Exists(ctx context.Context, conditions ...Condition) (bool, error)
	Create(ctx context.Context, obj *T) error
	Update(ctx context.Context, obj *T) error
	Delete(ctx context.Context, obj *T) error
}

// QuerySetFactory creates a QuerySet for a given type.
type QuerySetFactory[T any] func(db *DB) QuerySet[T]

// ManagerFactory creates a Manager for a given type.
type ManagerFactory[T any] func(db *DB) Manager[T]

// Expr represents a SQL expression that can be converted to SQL.
type Expr interface {
	ToSQL() (string, []interface{})
}
