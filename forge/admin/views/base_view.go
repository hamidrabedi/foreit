package views

import (
	"context"
	"net/http"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
)

// BaseView provides base functionality for all admin views
type BaseView[T any] struct {
	admin *admin.Admin[T]
}

// NewBaseView creates a new base view
func NewBaseView[T any](admin *admin.Admin[T]) *BaseView[T] {
	return &BaseView[T]{
		admin: admin,
	}
}

// Admin returns the admin instance
func (bv *BaseView[T]) Admin() *admin.Admin[T] {
	return bv.admin
}

// GetQueryset gets the base queryset with optional custom filtering
func (bv *BaseView[T]) GetQueryset(ctx context.Context) (orm.QuerySet[T], error) {
	return bv.admin.GetQueryset(ctx)
}

// ApplyFilters applies filters from request to queryset
func (bv *BaseView[T]) ApplyFilters(ctx context.Context, r *http.Request, qs orm.QuerySet[T]) (orm.QuerySet[T], error) {
	// Set base queryset on filterset
	filterset := bv.admin.FilterSet()
	filterset = filterset.WithQueryset(qs)

	// Apply filters from request
	return filterset.ApplyRequest(ctx, r)
}

// GetPage gets the page number from request
func GetPage(r *http.Request) int {
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		// Parse page number
		// For now, default to 1
		_ = p
	}
	return page
}

// GetPageSize gets the page size from request or config
func GetPageSize(r *http.Request, defaultSize int) int {
	pageSize := defaultSize
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		// Parse page size
		// For now, use default
		_ = ps
	}
	return pageSize
}

// GetSearch gets the search term from request
func GetSearch(r *http.Request) string {
	return r.URL.Query().Get("search")
}
