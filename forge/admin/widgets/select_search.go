package widgets

import (
	"fmt"
	"html/template"

	"github.com/forgego/forge/admin"
)

// SelectSearchWidget is a select widget with search/autocomplete
type SelectSearchWidget struct {
	choices     []admin.Choice[interface{}]
	searchable  bool
	placeholder string
	allowClear  bool
}

// NewSelectSearch creates a new searchable select widget
func NewSelectSearch(choices []admin.Choice[interface{}]) *SelectSearchWidget {
	return &SelectSearchWidget{
		choices:    choices,
		searchable: true,
		placeholder: "Search or select...",
		allowClear: true,
	}
}

// WithPlaceholder sets the placeholder text
func (w *SelectSearchWidget) WithPlaceholder(placeholder string) *SelectSearchWidget {
	w.placeholder = placeholder
	return w
}

// WithAllowClear allows clearing the selection
func (w *SelectSearchWidget) WithAllowClear(allowClear bool) *SelectSearchWidget {
	w.allowClear = allowClear
	return w
}

// Type returns the widget type
func (w *SelectSearchWidget) Type() string {
	return "select_search"
}

// Render renders the widget
func (w *SelectSearchWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}
	
	html := fmt.Sprintf(`<select name="%s" id="id_%s"`, name, name)
	html += ` class="form-select select-search"`
	html += fmt.Sprintf(` data-placeholder="%s"`, template.HTMLEscapeString(w.placeholder))
	
	if w.allowClear {
		html += ` data-allow-clear="true"`
	}
	
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	
	html += `>`
	
	if w.allowClear {
		html += `<option value="">---------</option>`
	}
	
	for _, choice := range w.choices {
		choiceValue := fmt.Sprintf("%v", choice.Value)
		selected := ""
		if choiceValue == valueStr {
			selected = ` selected`
		}
		html += fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			template.HTMLEscapeString(choiceValue),
			selected,
			template.HTMLEscapeString(choice.Label),
		)
	}
	
	html += `</select>`
	html += fmt.Sprintf(`<script>initSelectSearch('id_%s');</script>`, name)
	
	return template.HTML(html)
}

// Parse parses the widget value
func (w *SelectSearchWidget) Parse(value string) (interface{}, error) {
	return value, nil
}
