package components

import (
	"fmt"
	"html/template"
)

// Table represents a data table component
type Table struct {
	Columns      []Column
	Rows         []Row
	HasSelect    bool
	HasActions   bool
	EmptyMessage string
	CSSClass     string
}

// Column represents a table column
type Column struct {
	Label     string
	Field     string
	Sortable  bool
	SortOrder string // "asc", "desc", ""
	Width     string
	Align     string // "left", "center", "right"
}

// Row represents a table row
type Row struct {
	ID      interface{}
	Data    map[string]interface{}
	Actions []Action
	URL     string
}

// Action represents a row action
type Action struct {
	Label string
	URL   string
	Class string
	Icon  string
}

// Render renders the table as HTML
func (t *Table) Render() template.HTML {
	html := `<table class="admin-table`
	if t.CSSClass != "" {
		html += " " + t.CSSClass
	}
	html += `">`

	// Header
	html += `<thead><tr>`
	if t.HasSelect {
		html += `<th class="select-column"><input type="checkbox" class="select-all"></th>`
	}
	for _, col := range t.Columns {
		html += `<th`
		if col.Width != "" {
			html += ` width="` + template.HTMLEscapeString(col.Width) + `"`
		}
		if col.Align != "" {
			html += ` align="` + template.HTMLEscapeString(col.Align) + `"`
		}
		html += `>`
		if col.Sortable {
			html += `<a href="?sort=` + template.HTMLEscapeString(col.Field)
			if col.SortOrder == "asc" {
				html += `&order=desc`
			} else {
				html += `&order=asc`
			}
			html += `">` + template.HTMLEscapeString(col.Label) + `</a>`
		} else {
			html += template.HTMLEscapeString(col.Label)
		}
		html += `</th>`
	}
	if t.HasActions {
		html += `<th class="actions-column">Actions</th>`
	}
	html += `</tr></thead>`

	// Body
	html += `<tbody>`
	if len(t.Rows) == 0 {
		colspan := len(t.Columns)
		if t.HasSelect {
			colspan++
		}
		if t.HasActions {
			colspan++
		}
		html += `<tr><td colspan="` + template.HTMLEscapeString(string(rune(colspan))) + `" class="empty-state">`
		if t.EmptyMessage != "" {
			html += template.HTMLEscapeString(t.EmptyMessage)
		} else {
			html += `No data available.`
		}
		html += `</td></tr>`
	} else {
		for _, row := range t.Rows {
			html += `<tr`
			if row.URL != "" {
				html += ` data-url="` + template.HTMLEscapeString(row.URL) + `"`
			}
			html += `>`
			if t.HasSelect {
				html += `<td class="select-column"><input type="checkbox" name="selected" value="`
				html += template.HTMLEscapeString(fmt.Sprintf("%v", row.ID)) + `"></td>`
			}
			for _, col := range t.Columns {
				html += `<td>`
				if val, ok := row.Data[col.Field]; ok {
					html += template.HTMLEscapeString(fmt.Sprintf("%v", val))
				}
				html += `</td>`
			}
			if t.HasActions {
				html += `<td class="actions-column">`
				for _, action := range row.Actions {
					html += `<a href="` + template.HTMLEscapeString(action.URL) + `" class="btn btn-sm`
					if action.Class != "" {
						html += " " + action.Class
					}
					html += `">`
					if action.Icon != "" {
						html += `<i class="icon-` + template.HTMLEscapeString(action.Icon) + `"></i> `
					}
					html += template.HTMLEscapeString(action.Label) + `</a> `
				}
				html += `</td>`
			}
			html += `</tr>`
		}
	}
	html += `</tbody></table>`

	return template.HTML(html)
}
