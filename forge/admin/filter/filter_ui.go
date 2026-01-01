package filter

import (
	"context"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/orm"
)

// FilterUI provides UI-related functionality for filters
type FilterUI[T any] struct {
	filterset      *filter.FilterSet[T]
	facetCalculator *FacetCalculator[T]
}

// NewFilterUI creates a new filter UI helper
func NewFilterUI[T any](filterset *filter.FilterSet[T]) *FilterUI[T] {
	return &FilterUI[T]{
		filterset:      filterset,
		facetCalculator: NewFacetCalculator(filterset),
	}
}

// GetFilterOptions gets filter options for UI rendering
func (fui *FilterUI[T]) GetFilterOptions(ctx context.Context, filterName string, qs orm.QuerySet[T]) ([]filter.FilterOption, error) {
	f, ok := fui.filterset.GetFilter(filterName)
	if !ok {
		return nil, nil
	}

	return f.GetOptions(ctx, qs)
}

// GetAllFilterOptions gets options for all filters
func (fui *FilterUI[T]) GetAllFilterOptions(ctx context.Context, qs orm.QuerySet[T]) (map[string][]filter.FilterOption, error) {
	result := make(map[string][]filter.FilterOption)
	filters := fui.filterset.GetFilters()

	for name, f := range filters {
		options, err := f.GetOptions(ctx, qs)
		if err != nil {
			return nil, err
		}
		result[name] = options
	}

	return result, nil
}

// GetFilterWidget gets the widget for a filter
func (fui *FilterUI[T]) GetFilterWidget(filterName string) filter.Widget {
	f, ok := fui.filterset.GetFilter(filterName)
	if !ok {
		return nil
	}

	return f.GetWidget()
}

// GetAllFilterWidgets gets widgets for all filters
func (fui *FilterUI[T]) GetAllFilterWidgets() map[string]filter.Widget {
	result := make(map[string]filter.Widget)
	filters := fui.filterset.GetFilters()

	for name, f := range filters {
		result[name] = f.GetWidget()
	}

	return result
}

// GetFilterOptionsWithFacets gets filter options with facet counts
func (fui *FilterUI[T]) GetFilterOptionsWithFacets(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	filterName string,
) ([]filter.FilterOption, error) {
	return fui.facetCalculator.AddFacetCountsToOptions(ctx, baseQs, filterName)
}

// GetAllFilterOptionsWithFacets gets all filter options with facet counts
func (fui *FilterUI[T]) GetAllFilterOptionsWithFacets(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	excludeFilter string,
) (map[string][]filter.FilterOption, error) {
	results := make(map[string][]filter.FilterOption)
	filters := fui.filterset.GetFilters()

	for name := range filters {
		if name == excludeFilter {
			continue
		}

		options, err := fui.GetFilterOptionsWithFacets(ctx, baseQs, name)
		if err != nil {
			return nil, err
		}

		results[name] = options
	}

	return results, nil
}

// GetFacetCounts gets facet counts for all filters
func (fui *FilterUI[T]) GetFacetCounts(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	excludeFilter string,
) (map[string]*FacetResult, error) {
	return fui.facetCalculator.CalculateFacets(ctx, baseQs, excludeFilter)
}
