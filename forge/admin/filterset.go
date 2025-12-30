package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/forgego/forge/filter"
	query "github.com/forgego/forge/orm"
)

// FilterSetIntegration integrates FilterSet with admin ListView
type FilterSetIntegration[T any] struct {
	filterSet *filter.FilterSet[T]
	parser    *filter.Parser
}

// NewFilterSetIntegration creates a new filter set integration for admin
func NewFilterSetIntegration[T any](fs *filter.FilterSet[T]) *FilterSetIntegration[T] {
	return &FilterSetIntegration[T]{
		filterSet: fs,
		parser:    filter.NewParser(fs.GetSecurityConfig()),
	}
}

// ApplyToListView applies filters to a ListView
func (fsi *FilterSetIntegration[T]) ApplyToListView(ctx context.Context, r *http.Request, lv *ListView[T]) error {
	// Parse query parameters
	ast, err := fsi.parser.ParseFilterNode(r, nil)
	if err != nil {
		return fmt.Errorf("failed to parse filters: %w", err)
	}

	if ast != nil {
		fsi.filterSet.SetAST(ast)
		qs, err := fsi.filterSet.ApplyAST(ctx, ast)
		if err != nil {
			return fmt.Errorf("failed to apply filters: %w", err)
		}
		lv.queryset = qs
	}

	return nil
}

// GetFilterSidebarData returns data for rendering the filter sidebar
func (fsi *FilterSetIntegration[T]) GetFilterSidebarData(ctx context.Context, qs query.QuerySet[T]) ([]FilterSidebarItem, error) {
	items := make([]FilterSidebarItem, 0)

	// Iterate through filters and build sidebar items
	filters := fsi.filterSet.GetFilters()
	for name, f := range filters {
		options, err := f.GetOptions(ctx, qs)
		if err != nil {
			continue
		}

		item := FilterSidebarItem{
			Name:    name,
			Label:   name, // Would get from filter
			Widget:  f.GetWidget(),
			Options: options,
		}

		items = append(items, item)
	}

	return items, nil
}

// FilterSidebarItem represents an item in the filter sidebar
type FilterSidebarItem struct {
	Name    string
	Label   string
	Widget  filter.Widget
	Options []filter.FilterOption
}

