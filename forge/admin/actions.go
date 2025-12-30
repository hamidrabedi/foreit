package admin

import (
	"context"
)

// Action represents a bulk action
type Action[T any] struct {
	Name        string
	Label       string
	Description string
	Handler     func(ctx context.Context, instances []*T) error
	Permissions []string
}

// NewAction creates a new action
func NewAction[T any](
	name string,
	label string,
	handler func(ctx context.Context, instances []*T) error,
) Action[T] {
	return Action[T]{
		Name:    name,
		Label:   label,
		Handler: handler,
	}
}

// WithDescription sets the action description
func (a Action[T]) WithDescription(desc string) Action[T] {
	a.Description = desc
	return a
}

// WithPermissions sets required permissions
func (a Action[T]) WithPermissions(perms ...string) Action[T] {
	a.Permissions = perms
	return a
}
