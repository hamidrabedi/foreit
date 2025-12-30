package admin

import (
	"context"
	"net/http"
)

// ViewHooks provides hooks for customizing admin views in the new admincore system
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
// Returns interface{} to avoid import cycle - will be cast to *views.ListView[T] in handlers
type ChangelistViewHook[T any] func(ctx context.Context, admin *Admin[T], request *http.Request) (interface{}, error)

// ChangeViewHook allows customizing the change (edit) view
// Returns interface{} to avoid import cycle - will be cast to *views.FormView[T] in handlers
type ChangeViewHook[T any] func(ctx context.Context, admin *Admin[T], obj *T, request *http.Request) (interface{}, error)

// AddViewHook allows customizing the add (create) view
// Returns interface{} to avoid import cycle - will be cast to *views.FormView[T] in handlers
type AddViewHook[T any] func(ctx context.Context, admin *Admin[T], request *http.Request) (interface{}, error)

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

// GetURLsHook allows adding custom URLs to admin
type GetURLsHook[T any] func(admin *Admin[T]) []URLPattern

// URLPattern represents a URL pattern for admin
type URLPattern struct {
	Path    string
	Handler http.HandlerFunc
	Name    string
}

// LogEntry represents a change history entry
type LogEntry struct {
	ID        int64
	ObjectID  int64
	Action    string
	UserID    int64
	Timestamp int64
	Changes   map[string]interface{}
}
