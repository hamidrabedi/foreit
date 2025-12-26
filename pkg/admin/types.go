package admin

import (
	"github.com/forgego/forge/pkg/query"
	"github.com/forgego/forge/pkg/schema"
)

// ManagerInterface defines the interface that all managers must implement
type ManagerInterface[T any] interface {
	Get(ctx interface{}, id int64) (*T, error)
	All(ctx interface{}) ([]*T, error)
	Create(ctx interface{}, instance *T) error
	Update(ctx interface{}, instance *T) error
	Delete(ctx interface{}, instance *T) error
	Filter(expr ...query.QueryExpr) query.QuerySet[T]
}

// AdminModelGeneric is a generic admin model that uses type-safe operations
type AdminModelGeneric[T any] struct {
	Name          string
	Model         T
	Manager       ManagerInterface[T]
	
	// Auto-generated from Meta and field definitions
	ListDisplay   []interface{}
	ListFilter    []interface{}
	SearchFields  []interface{}
	ReadOnlyFields []interface{}
	
	// Extensibility
	CustomAdmin CustomAdmin
}

// GetManager returns the manager as ManagerInterface
func (m *AdminModelGeneric[T]) GetManager() ManagerInterface[T] {
	return m.Manager
}

// AdminModel is defined in registry.go to avoid circular dependencies

// SchemaModel is a constraint for models that implement schema.Schema
type SchemaModel interface {
	schema.Schema
}

// ModelWithID is a constraint for models with ID field
type ModelWithID interface {
	GetID() int64
	SetID(int64)
}

