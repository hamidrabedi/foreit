package filter

import (
	"fmt"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/filter/filters"
)

// TextSearchFilter provides enhanced text search with multiple lookup types
type TextSearchFilter[T any] struct {
	*filters.CharFilter[T]
	lookupType string // "contains", "startswith", "endswith", "exact"
}

// NewTextSearchFilter creates a new text search filter
func NewTextSearchFilter[T any](fieldPath string) *TextSearchFilter[T] {
	baseFilter := filters.NewCharFilter[T](fieldPath)
	return &TextSearchFilter[T]{
		CharFilter: baseFilter,
		lookupType: "icontains", // Default to case-insensitive contains
	}
}

// WithLookup sets the lookup type for the filter
func (f *TextSearchFilter[T]) WithLookup(lookup string) *TextSearchFilter[T] {
	validLookups := map[string]bool{
		"contains":    true,
		"icontains":   true,
		"startswith":  true,
		"istartswith": true,
		"endswith":    true,
		"iendswith":   true,
		"exact":       true,
		"iexact":      true,
	}

	if !validLookups[lookup] {
		return f // Invalid lookup, keep current
	}

	f.lookupType = lookup
	// Note: BaseFilter doesn't have SetLookup, lookup is set via GetLookup()
	// The lookup is stored in the BaseFilter's lookup field
	return f
}

// Contains sets the filter to use case-sensitive contains
func (f *TextSearchFilter[T]) Contains() *TextSearchFilter[T] {
	return f.WithLookup("contains")
}

// IContains sets the filter to use case-insensitive contains (default)
func (f *TextSearchFilter[T]) IContains() *TextSearchFilter[T] {
	return f.WithLookup("icontains")
}

// StartsWith sets the filter to use case-sensitive starts with
func (f *TextSearchFilter[T]) StartsWith() *TextSearchFilter[T] {
	return f.WithLookup("startswith")
}

// IStartsWith sets the filter to use case-insensitive starts with
func (f *TextSearchFilter[T]) IStartsWith() *TextSearchFilter[T] {
	return f.WithLookup("istartswith")
}

// EndsWith sets the filter to use case-sensitive ends with
func (f *TextSearchFilter[T]) EndsWith() *TextSearchFilter[T] {
	return f.WithLookup("endswith")
}

// IEndsWith sets the filter to use case-insensitive ends with
func (f *TextSearchFilter[T]) IEndsWith() *TextSearchFilter[T] {
	return f.WithLookup("iendswith")
}

// Exact sets the filter to use exact match
func (f *TextSearchFilter[T]) Exact() *TextSearchFilter[T] {
	return f.WithLookup("exact")
}

// IExact sets the filter to use case-insensitive exact match
func (f *TextSearchFilter[T]) IExact() *TextSearchFilter[T] {
	return f.WithLookup("iexact")
}

// GetLookupType returns the current lookup type
func (f *TextSearchFilter[T]) GetLookupType() string {
	return f.lookupType
}

// GetWidget returns the widget for this filter
func (f *TextSearchFilter[T]) GetWidget() filter.Widget {
	return &TextSearchWidget{
		lookupType: f.lookupType,
		fieldName:  f.GetFieldPath(),
	}
}

// TextSearchWidget is a widget for text search filters
type TextSearchWidget struct {
	lookupType string
	fieldName  string
}

// Type returns the widget type
func (w *TextSearchWidget) Type() string {
	return "text_search"
}

// Render renders the widget HTML with search input and lookup selector
func (w *TextSearchWidget) Render(name string, value interface{}, attrs map[string]string) (string, error) {
	valueStr := ""
	if value != nil {
		if str, ok := value.(string); ok {
			valueStr = str
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}
	}

	// Build lookup options
	lookupOptions := []struct {
		value string
		label string
	}{
		{"icontains", "Contains (case-insensitive)"},
		{"contains", "Contains (case-sensitive)"},
		{"istartswith", "Starts with (case-insensitive)"},
		{"startswith", "Starts with (case-sensitive)"},
		{"iendswith", "Ends with (case-insensitive)"},
		{"endswith", "Ends with (case-sensitive)"},
		{"iexact", "Exact (case-insensitive)"},
		{"exact", "Exact (case-sensitive)"},
	}

	lookupSelect := fmt.Sprintf(`<select name="%s_lookup" class="form-control text-search-lookup">`, name)
	for _, opt := range lookupOptions {
		selected := ""
		if opt.value == w.lookupType {
			selected = " selected"
		}
		lookupSelect += fmt.Sprintf(`<option value="%s"%s>%s</option>`, opt.value, selected, opt.label)
	}
	lookupSelect += "</select>"

	html := fmt.Sprintf(`
		<div class="text-search-filter">
			<div class="text-search-input-group">
				<input type="text" name="%s" value="%s" placeholder="Search..." class="form-control text-search-input" />
				%s
			</div>
		</div>`,
		name, valueStr, lookupSelect,
	)

	return html, nil
}

// Parse parses the widget value
func (w *TextSearchWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	return value, nil
}
