package widgets

import (
	"fmt"
	"html/template"
)

// RichTextWidget is a rich text editor widget
type RichTextWidget struct {
	Height    int    // Height in pixels
	Toolbar   string // Toolbar configuration
	Plugins   []string // Editor plugins
	ReadOnly  bool
}

// NewRichText creates a new rich text widget
func NewRichText() *RichTextWidget {
	return &RichTextWidget{
		Height:  300,
		Toolbar: "full", // "full", "basic", "minimal"
		Plugins: []string{"basic"},
	}
}

// WithHeight sets the editor height
func (w *RichTextWidget) WithHeight(height int) *RichTextWidget {
	w.Height = height
	return w
}

// WithToolbar sets the toolbar configuration
func (w *RichTextWidget) WithToolbar(toolbar string) *RichTextWidget {
	w.Toolbar = toolbar
	return w
}

// WithPlugins sets editor plugins
func (w *RichTextWidget) WithPlugins(plugins ...string) *RichTextWidget {
	w.Plugins = plugins
	return w
}

// Type returns the widget type
func (w *RichTextWidget) Type() string {
	return "richtext"
}

// Render renders the widget
func (w *RichTextWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}
	
	html := fmt.Sprintf(`<textarea name="%s" id="id_%s"`, name, name)
	html += fmt.Sprintf(` class="form-control richtext-editor"`)
	html += fmt.Sprintf(` data-toolbar="%s"`, template.HTMLEscapeString(w.Toolbar))
	html += fmt.Sprintf(` data-height="%d"`, w.Height)
	
	if w.ReadOnly {
		html += ` readonly`
	}
	
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	
	html += fmt.Sprintf(` style="height: %dpx;">`, w.Height)
	html += template.HTMLEscapeString(valueStr)
	html += `</textarea>`
	
	// Add initialization script
	html += fmt.Sprintf(`<script>initRichTextEditor('id_%s');</script>`, name)
	
	return template.HTML(html)
}

// Parse parses the widget value
func (w *RichTextWidget) Parse(value string) (interface{}, error) {
	// Rich text is stored as HTML
	return value, nil
}
