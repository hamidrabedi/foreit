package admin

import (
	"context"
	"fmt"

	"github.com/forgego/forge/pkg/query"
)

// ManagerHelper provides type-safe access to managers
type ManagerHelper[T any] struct {
	manager *query.Manager[T]
}

// NewManagerHelper creates a new manager helper
func NewManagerHelper[T any](manager *query.Manager[T]) *ManagerHelper[T] {
	return &ManagerHelper[T]{
		manager: manager,
	}
}

// GetManager returns the underlying manager
func (h *ManagerHelper[T]) GetManager() *query.Manager[T] {
	return h.manager
}

// GetManagerInterface returns the manager as ManagerInterface
func (h *ManagerHelper[T]) GetManagerInterface() ManagerInterface[T] {
	return &managerAdapter[T]{manager: h.manager}
}

// managerAdapter adapts query.Manager to ManagerInterface
type managerAdapter[T any] struct {
	manager *query.Manager[T]
}

func (a *managerAdapter[T]) Get(ctx interface{}, id int64) (*T, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	// Type assert to context.Context
	if c, ok := ctx.(context.Context); ok {
		return a.manager.Get(c, id)
	}
	return nil, fmt.Errorf("invalid context type, expected context.Context")
}

func (a *managerAdapter[T]) All(ctx interface{}) ([]*T, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	if c, ok := ctx.(context.Context); ok {
		return a.manager.All(c)
	}
	return nil, fmt.Errorf("invalid context type, expected context.Context")
}

func (a *managerAdapter[T]) Create(ctx interface{}, instance *T) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}
	if c, ok := ctx.(context.Context); ok {
		return a.manager.Create(c, instance)
	}
	return fmt.Errorf("invalid context type, expected context.Context")
}

func (a *managerAdapter[T]) Update(ctx interface{}, instance *T) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}
	if c, ok := ctx.(context.Context); ok {
		return a.manager.Update(c, instance)
	}
	return fmt.Errorf("invalid context type, expected context.Context")
}

func (a *managerAdapter[T]) Delete(ctx interface{}, instance *T) error {
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}
	if c, ok := ctx.(context.Context); ok {
		return a.manager.Delete(c, instance)
	}
	return fmt.Errorf("invalid context type, expected context.Context")
}

func (a *managerAdapter[T]) Filter(expr ...query.QueryExpr) query.QuerySet[T] {
	return a.manager.Filter(expr...)
}

