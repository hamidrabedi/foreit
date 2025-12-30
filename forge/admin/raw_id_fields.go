package admin

import (
	"fmt"
	"html/template"
	"strconv"
)

// RawIDWidget is a widget for raw ID input (for foreign keys)
type RawIDWidget struct{}

// NewRawIDWidget creates a new raw ID widget
func NewRawIDWidget() Widget {
	return &RawIDWidget{}
}

// Type returns the widget type
func (w *RawIDWidget) Type() string {
	return "raw_id"
}

// Render renders the raw ID input widget
func (w *RawIDWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	html := fmt.Sprintf(`<input type="text" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control raw-id-field">`

	// Add lookup button (would trigger popup in real implementation)
	html += fmt.Sprintf(` <a href="#" class="raw-id-lookup" data-field="%s">Lookup</a>`, name)

	return template.HTML(html)
}

// Parse parses the raw ID value
func (w *RawIDWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	// Parse as int64 (typical for IDs)
	// In a real implementation, would handle different ID types
	return strconv.ParseInt(value, 10, 64)
}
