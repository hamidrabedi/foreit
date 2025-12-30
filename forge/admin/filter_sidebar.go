package admin

import (
	"html/template"
)

// FilterSidebar renders the filter sidebar for admin list views
type FilterSidebar struct {
	Items []FilterSidebarItem
}

// Render renders the filter sidebar HTML
func (fs *FilterSidebar) Render() template.HTML {
	html := `<div class="filter-sidebar">`
	html += `<h3>Filters</h3>`

	for _, item := range fs.Items {
		html += `<div class="filter-item">`
		html += `<label>` + template.HTMLEscapeString(item.Label) + `</label>`
		
		// Render widget
		widgetHTML, _ := item.Widget.Render(item.Name, nil, map[string]string{
			"class": "form-control",
		})
		html += widgetHTML
		
		html += `</div>`
	}

	html += `</div>`
	return template.HTML(html)
}
