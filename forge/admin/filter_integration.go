package admin

import (
	"context"
	"net/http"

	"github.com/forgego/forge/filter"
)

// FilterSetConfig extends admin Config with FilterSet support
type FilterSetConfig[T any] struct {
	FilterSet *filter.FilterSet[T]
}

// WithFilterSet adds a FilterSet to the admin configuration
func (c *Config[T]) WithFilterSet(fs *filter.FilterSet[T]) *Config[T] {
	// Store FilterSet in a way that can be accessed later
	// This would require extending Config struct or using a registry
	return c
}

// ApplyFilterSetToListView applies FilterSet filters to a ListView from HTTP request
func ApplyFilterSetToListView[T any](ctx context.Context, r *http.Request, lv *ListView[T], fs *filter.FilterSet[T]) error {
	if fs == nil || r == nil {
		return nil
	}

	integration := NewFilterSetIntegration(fs)
	return integration.ApplyToListView(ctx, r, lv)
}

// GetFilterSidebarForListView gets filter sidebar data for a ListView
func GetFilterSidebarForListView[T any](ctx context.Context, lv *ListView[T], fs *filter.FilterSet[T]) ([]FilterSidebarItem, error) {
	if fs == nil {
		return nil, nil
	}

	integration := NewFilterSetIntegration(fs)
	return integration.GetFilterSidebarData(ctx, lv.queryset)
}
