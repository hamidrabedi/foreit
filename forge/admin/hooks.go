package admin

import (
	"context"
	"net/http"

	query "github.com/forgego/forge/orm"
)

// ViewHooks provides hooks for customizing admin views
type ViewHooks[T any] struct {
	admin *Admin[T]
}

// NewViewHooks creates a new view hooks instance
func NewViewHooks[T any](admin *Admin[T]) *ViewHooks[T] {
	return &ViewHooks[T]{
		admin: admin,
	}
}

// ChangelistViewHook allows customizing the changelist (list) view
type ChangelistViewHook[T any] func(ctx context.Context, admin *Admin[T], request *http.Request) (*ListView[T], error)

// ChangeViewHook allows customizing the change (edit) view
type ChangeViewHook[T any] func(ctx context.Context, admin *Admin[T], obj *T, request *http.Request) (*FormView[T], error)

// AddViewHook allows customizing the add (create) view
type AddViewHook[T any] func(ctx context.Context, admin *Admin[T], request *http.Request) (*FormView[T], error)

// DeleteViewHook allows customizing the delete view
type DeleteViewHook[T any] func(ctx context.Context, admin *Admin[T], obj *T, request *http.Request) error

// HistoryViewHook allows customizing the history view
type HistoryViewHook[T any] func(ctx context.Context, admin *Admin[T], objID int64, request *http.Request) ([]*LogEntry, error)

// ResponseAddHook allows customizing response after add
type ResponseAddHook[T any] func(ctx context.Context, admin *Admin[T], obj *T, request *http.Request, response http.ResponseWriter) error

// ResponseChangeHook allows customizing response after change
type ResponseChangeHook[T any] func(ctx context.Context, admin *Admin[T], obj *T, request *http.Request, response http.ResponseWriter) error

// ResponseDeleteHook allows customizing response after delete
type ResponseDeleteHook[T any] func(ctx context.Context, admin *Admin[T], obj *T, request *http.Request, response http.ResponseWriter) error

// MessageUser sends a message to the user
func MessageUser(ctx context.Context, message string, level MessageLevel) {
	// In a real implementation, this would store messages in the session/flash
	// For now, it's a placeholder
	_ = ctx
	_ = message
	_ = level
}

// MessageLevel represents the level of a message
type MessageLevel string

const (
	MessageSuccess MessageLevel = "success"
	MessageInfo    MessageLevel = "info"
	MessageWarning MessageLevel = "warning"
	MessageError   MessageLevel = "error"
)

// GetURLsHook allows adding custom URLs to admin
type GetURLsHook[T any] func(admin *Admin[T]) []URLPattern

// URLPattern represents a URL pattern for admin
type URLPattern struct {
	Path    string
	Handler http.HandlerFunc
	Name    string
}

// GetQuerysetHook allows customizing the base queryset
type GetQuerysetHook[T any] func(ctx context.Context, admin *Admin[T], qs query.QuerySet[T]) (query.QuerySet[T], error)
