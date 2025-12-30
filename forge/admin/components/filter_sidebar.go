package components

import (
	"fmt"
	"html/template"
)

// FilterSidebar represents a filter sidebar component
type FilterSidebar struct {
	Filters    []Filter
	ActiveURL  string
	ClearURL   string
	Title      string
	Collapsed  bool
}

// Filter represents a filter option
type Filter struct {
	Name        string
	Label       string
	Type        string // "boolean", "choice", "date", "number", "related"
	Options     []FilterOption
	Selected    []string
	Multiple    bool
	Lookups     []Lookup // For date/number filters
	SelectedLookup string
}

// FilterOption represents a filter option value
type FilterOption struct {
	Label string
	Value string
	Count int64 // Number of results with this filter
	Selected bool
}

// Lookup represents a filter lookup (e.g., "exact", "gt", "lt" for dates)
type Lookup struct {
	Label string
	Value string
	Selected bool
}

// Render renders the filter sidebar as HTML
func (fs *FilterSidebar) Render() template.HTML {
	html := `<div class="admin-filter-sidebar`
	if fs.Collapsed {
		html += ` collapsed`
	}
	html += `">`

	if fs.Title != "" {
		html += `<h3 class="filter-title">` + template.HTMLEscapeString(fs.Title) + `</h3>`
	}

	if fs.ClearURL != "" {
		html += `<a href="` + template.HTMLEscapeString(fs.ClearURL) + `" class="filter-clear">Clear all filters</a>`
	}

	html += `<ul class="filter-list">`
	for _, filter := range fs.Filters {
		html += string(fs.renderFilter(filter))
	}
	html += `</ul>`
	html += `</div>`

	return template.HTML(html)
}

// renderFilter renders a single filter
func (fs *FilterSidebar) renderFilter(filter Filter) template.HTML {
	html := `<li class="filter-item">`
	html += `<div class="filter-header">`
	html += `<h4 class="filter-name">` + template.HTMLEscapeString(filter.Label) + `</h4>`
	html += `</div>`

	html += `<div class="filter-content">`

	switch filter.Type {
	case "boolean":
		html += string(fs.renderBooleanFilter(filter))
	case "choice":
		html += string(fs.renderChoiceFilter(filter))
	case "date", "datetime":
		html += string(fs.renderDateFilter(filter))
	case "number":
		html += string(fs.renderNumberFilter(filter))
	case "related":
		html += string(fs.renderRelatedFilter(filter))
	default:
		html += string(fs.renderChoiceFilter(filter))
	}

	html += `</div>`
	html += `</li>`

	return template.HTML(html)
}

// renderBooleanFilter renders a boolean filter
func (fs *FilterSidebar) renderBooleanFilter(filter Filter) template.HTML {
	html := `<div class="filter-boolean">`
	for _, option := range filter.Options {
		checked := ""
		if option.Selected {
			checked = ` checked`
		}
		html += `<label class="filter-checkbox">`
		html += `<input type="checkbox" name="` + template.HTMLEscapeString(filter.Name) + `" value="` + template.HTMLEscapeString(option.Value) + `"` + checked + `>`
		html += `<span>` + template.HTMLEscapeString(option.Label) + `</span>`
		if option.Count > 0 {
			html += `<span class="filter-count">(` + template.HTMLEscapeString(fmt.Sprintf("%d", option.Count)) + `)</span>`
		}
		html += `</label>`
	}
	html += `</div>`
	return template.HTML(html)
}

// renderChoiceFilter renders a choice filter
func (fs *FilterSidebar) renderChoiceFilter(filter Filter) template.HTML {
	html := `<div class="filter-choice">`
	if filter.Multiple {
		html += `<select name="` + template.HTMLEscapeString(filter.Name) + `" multiple class="filter-select">`
	} else {
		html += `<select name="` + template.HTMLEscapeString(filter.Name) + `" class="filter-select">`
		html += `<option value="">All</option>`
	}

	for _, option := range filter.Options {
		selected := ""
		if option.Selected {
			selected = ` selected`
		}
		html += `<option value="` + template.HTMLEscapeString(option.Value) + `"` + selected + `>`
		html += template.HTMLEscapeString(option.Label)
		if option.Count > 0 {
			html += ` (` + template.HTMLEscapeString(fmt.Sprintf("%d", option.Count)) + `)`
		}
		html += `</option>`
	}
	html += `</select>`
	html += `</div>`
	return template.HTML(html)
}

// renderDateFilter renders a date filter
func (fs *FilterSidebar) renderDateFilter(filter Filter) template.HTML {
	html := `<div class="filter-date">`
	if len(filter.Lookups) > 0 {
		html += `<select name="` + template.HTMLEscapeString(filter.Name) + `_lookup" class="filter-lookup">`
		for _, lookup := range filter.Lookups {
			selected := ""
			if lookup.Selected {
				selected = ` selected`
			}
			html += `<option value="` + template.HTMLEscapeString(lookup.Value) + `"` + selected + `>`
			html += template.HTMLEscapeString(lookup.Label)
			html += `</option>`
		}
		html += `</select>`
	}
	html += `<input type="date" name="` + template.HTMLEscapeString(filter.Name) + `" class="filter-date-input">`
	html += `</div>`
	return template.HTML(html)
}

// renderNumberFilter renders a number filter
func (fs *FilterSidebar) renderNumberFilter(filter Filter) template.HTML {
	html := `<div class="filter-number">`
	if len(filter.Lookups) > 0 {
		html += `<select name="` + template.HTMLEscapeString(filter.Name) + `_lookup" class="filter-lookup">`
		for _, lookup := range filter.Lookups {
			selected := ""
			if lookup.Selected {
				selected = ` selected`
			}
			html += `<option value="` + template.HTMLEscapeString(lookup.Value) + `"` + selected + `>`
			html += template.HTMLEscapeString(lookup.Label)
			html += `</option>`
		}
		html += `</select>`
	}
	html += `<input type="number" name="` + template.HTMLEscapeString(filter.Name) + `" class="filter-number-input">`
	html += `</div>`
	return template.HTML(html)
}

// renderRelatedFilter renders a related filter
func (fs *FilterSidebar) renderRelatedFilter(filter Filter) template.HTML {
	// Similar to choice filter but with autocomplete support
	return fs.renderChoiceFilter(filter)
}
