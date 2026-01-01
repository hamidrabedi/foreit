package filter

import (
	"context"
	"fmt"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/orm"
)

// FacetCount represents a count for a filter option
type FacetCount struct {
	Value interface{}
	Count int64
	Label string
}

// FacetResult contains facet counts for a filter
type FacetResult struct {
	FilterName string
	Facets     []FacetCount
	TotalCount int64
}

// FacetCalculator calculates facet counts for filters
type FacetCalculator[T any] struct {
	filterset *filter.FilterSet[T]
}

// NewFacetCalculator creates a new facet calculator
func NewFacetCalculator[T any](filterset *filter.FilterSet[T]) *FacetCalculator[T] {
	return &FacetCalculator[T]{
		filterset: filterset,
	}
}

// CalculateFacets calculates facet counts for all filters
func (fc *FacetCalculator[T]) CalculateFacets(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	excludeFilter string, // Filter to exclude when calculating facets
) (map[string]*FacetResult, error) {
	results := make(map[string]*FacetResult)
	filters := fc.filterset.GetFilters()

	for filterName, f := range filters {
		if filterName == excludeFilter {
			continue
		}

		facetResult, err := fc.CalculateFilterFacets(ctx, baseQs, filterName, f)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate facets for filter %s: %w", filterName, err)
		}

		if facetResult != nil {
			results[filterName] = facetResult
		}
	}

	return results, nil
}

// CalculateFilterFacets calculates facet counts for a specific filter
func (fc *FacetCalculator[T]) CalculateFilterFacets(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	filterName string,
	f filter.Filter[T],
) (*FacetResult, error) {
	// Get filter options
	options, err := f.GetOptions(ctx, baseQs)
	if err != nil {
		return nil, fmt.Errorf("failed to get filter options: %w", err)
	}

	if len(options) == 0 {
		return nil, nil // No options to count
	}

	// Calculate count for each option
	facets := make([]FacetCount, 0, len(options))
	totalCount := int64(0)

	for _, option := range options {
		// Apply filter with this option value
		filteredQs, err := f.Apply(ctx, baseQs, option.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to apply filter for option %v: %w", option.Value, err)
		}

		// Count results
		count, err := filteredQs.Count(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count results: %w", err)
		}

		facets = append(facets, FacetCount{
			Value: option.Value,
			Count: count,
			Label: option.Label,
		})

		totalCount += count
	}

	return &FacetResult{
		FilterName: filterName,
		Facets:     facets,
		TotalCount: totalCount,
	}, nil
}

// GetFacetCountsForFilter gets facet counts for a specific filter
func (fc *FacetCalculator[T]) GetFacetCountsForFilter(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	filterName string,
) ([]FacetCount, error) {
	f, ok := fc.filterset.GetFilter(filterName)
	if !ok {
		return nil, fmt.Errorf("filter %s not found", filterName)
	}

	result, err := fc.CalculateFilterFacets(ctx, baseQs, filterName, f)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return []FacetCount{}, nil
	}

	return result.Facets, nil
}

// AddFacetCountsToOptions adds facet counts to filter options
func (fc *FacetCalculator[T]) AddFacetCountsToOptions(
	ctx context.Context,
	baseQs orm.QuerySet[T],
	filterName string,
) ([]filter.FilterOption, error) {
	// Get base options
	f, ok := fc.filterset.GetFilter(filterName)
	if !ok {
		return nil, fmt.Errorf("filter %s not found", filterName)
	}

	options, err := f.GetOptions(ctx, baseQs)
	if err != nil {
		return nil, fmt.Errorf("failed to get filter options: %w", err)
	}

	// Get facet counts
	facets, err := fc.GetFacetCountsForFilter(ctx, baseQs, filterName)
	if err != nil {
		return nil, err
	}

	// Create a map of value -> count
	facetMap := make(map[interface{}]int64)
	for _, facet := range facets {
		facetMap[facet.Value] = facet.Count
	}

	// Add counts to options
	result := make([]filter.FilterOption, len(options))
	for i, option := range options {
		count := facetMap[option.Value]
		result[i] = filter.FilterOption{
			Label: option.Label,
			Value: option.Value,
			Count: count,
		}
	}

	return result, nil
}
