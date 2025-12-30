package filter

import (
	"context"
	"net/http"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/orm"
)

// AdminFilterSet wraps filter.FilterSet for admin use
type AdminFilterSet[T any] struct {
	filterset *filter.FilterSet[T]
	config    *FilterConfig[T]
}

// FilterConfig contains admin-specific filter configuration
type FilterConfig[T any] struct {
	EnabledFilters []string // List of enabled filter names
	FilterOptions  map[string]interface{}
}

// NewAdminFilterSet creates a new admin filterset wrapper
func NewAdminFilterSet[T any]() (*AdminFilterSet[T], error) {
	filterset, err := filter.NewFilterSet[T]()
	if err != nil {
		return nil, err
	}

	return &AdminFilterSet[T]{
		filterset: filterset,
		config: &FilterConfig[T]{
			EnabledFilters: []string{},
			FilterOptions:  make(map[string]interface{}),
		},
	}, nil
}

// WithQueryset sets the base queryset for filtering
func (afs *AdminFilterSet[T]) WithQueryset(qs orm.QuerySet[T]) *AdminFilterSet[T] {
	afs.filterset = afs.filterset.WithQueryset(qs)
	return afs
}

// ApplyRequest applies filters from an HTTP request
func (afs *AdminFilterSet[T]) ApplyRequest(ctx context.Context, r *http.Request) (orm.QuerySet[T], error) {
	return afs.filterset.ApplyQueryParams(ctx, r)
}

// AddFilter adds a filter to the filterset
func (afs *AdminFilterSet[T]) AddFilter(name string, f filter.Filter[T]) *AdminFilterSet[T] {
	afs.filterset.AddFilter(name, f)
	return afs
}

// GetFilter gets a filter by name
func (afs *AdminFilterSet[T]) GetFilter(name string) (filter.Filter[T], bool) {
	return afs.filterset.GetFilter(name)
}

// GetFilters returns all filters
func (afs *AdminFilterSet[T]) GetFilters() map[string]filter.Filter[T] {
	return afs.filterset.GetFilters()
}

// FilterSet returns the underlying filter filterset
func (afs *AdminFilterSet[T]) FilterSet() *filter.FilterSet[T] {
	return afs.filterset
}

// SetConfig sets the filter configuration
func (afs *AdminFilterSet[T]) SetConfig(config *FilterConfig[T]) *AdminFilterSet[T] {
	afs.config = config
	return afs
}

// GetConfig returns the filter configuration
func (afs *AdminFilterSet[T]) GetConfig() *FilterConfig[T] {
	return afs.config
}
