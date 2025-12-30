package components

import (
	"fmt"
	"html/template"
)

// Form represents a form component
type Form struct {
	Action      string
	Method      string
	Enctype     string
	Fields      []FormField
	Fieldsets   []FormFieldset
	Actions     []FormAction
	Errors      []string
	CSRFToken   string
	CSSClass    string
}

// FormField represents a form field
type FormField struct {
	Name        string
	Label       string
	Type        string
	Value       interface{}
	Widget      template.HTML
	HelpText    string
	Required    bool
	ReadOnly    bool
	Errors      []string
	CSSClass    string
	Placeholder string
	Attrs       map[string]string
}

// FormFieldset represents a form fieldset
type FormFieldset struct {
	Name        string
	Fields      []FormField
	Collapsed   bool
	Description string
	CSSClass    string
}

// FormAction represents a form action button
type FormAction struct {
	Label     string
	Name      string
	Value     string
	Type      string // "submit", "button", "link"
	URL       string // For link type
	Class     string
	Icon      string
	Primary   bool
}

// Render renders the form as HTML
func (f *Form) Render() template.HTML {
	html := `<form`
	if f.Action != "" {
		html += ` action="` + template.HTMLEscapeString(f.Action) + `"`
	}
	if f.Method != "" {
		html += ` method="` + template.HTMLEscapeString(f.Method) + `"`
	} else {
		html += ` method="post"`
	}
	if f.Enctype != "" {
		html += ` enctype="` + template.HTMLEscapeString(f.Enctype) + `"`
	}
	if f.CSSClass != "" {
		html += ` class="` + template.HTMLEscapeString(f.CSSClass) + `"`
	}
	html += `>`

	// CSRF Token
	if f.CSRFToken != "" {
		html += `<input type="hidden" name="csrf_token" value="` + template.HTMLEscapeString(f.CSRFToken) + `">`
	}

	// Errors
	if len(f.Errors) > 0 {
		html += `<div class="form-errors">`
		html += `<ul class="error-list">`
		for _, err := range f.Errors {
			html += `<li class="error-item">` + template.HTMLEscapeString(err) + `</li>`
		}
		html += `</ul>`
		html += `</div>`
	}

	// Fieldsets or Fields
	if len(f.Fieldsets) > 0 {
		for _, fieldset := range f.Fieldsets {
			html += string(f.renderFieldset(fieldset))
		}
	} else {
		for _, field := range f.Fields {
			html += string(f.renderField(field))
		}
	}

	// Actions
	if len(f.Actions) > 0 {
		html += `<div class="form-actions">`
		for _, action := range f.Actions {
			html += string(f.renderAction(action))
		}
		html += `</div>`
	}

	html += `</form>`
	return template.HTML(html)
}

// renderFieldset renders a fieldset
func (f *Form) renderFieldset(fieldset FormFieldset) template.HTML {
	html := `<fieldset class="form-fieldset`
	if fieldset.Collapsed {
		html += ` collapsed`
	}
	if fieldset.CSSClass != "" {
		html += ` ` + template.HTMLEscapeString(fieldset.CSSClass)
	}
	html += `">`

	if fieldset.Name != "" {
		html += `<legend class="fieldset-legend">` + template.HTMLEscapeString(fieldset.Name)
		if fieldset.Collapsed {
			html += `<button type="button" class="toggle-fieldset">Expand</button>`
		}
		html += `</legend>`
	}

	if fieldset.Description != "" {
		html += `<p class="fieldset-description">` + template.HTMLEscapeString(fieldset.Description) + `</p>`
	}

	html += `<div class="fieldset-fields">`
	for _, field := range fieldset.Fields {
		html += string(f.renderField(field))
	}
	html += `</div>`
	html += `</fieldset>`

	return template.HTML(html)
}

// renderField renders a form field
func (f *Form) renderField(field FormField) template.HTML {
	html := `<div class="form-field`
	if len(field.Errors) > 0 {
		html += ` has-error`
	}
	if field.Required {
		html += ` required`
	}
	if field.ReadOnly {
		html += ` readonly`
	}
	if field.CSSClass != "" {
		html += ` ` + template.HTMLEscapeString(field.CSSClass)
	}
	html += `">`

	// Label
	if field.Label != "" {
		html += `<label for="id_` + template.HTMLEscapeString(field.Name) + `" class="form-label">`
		html += template.HTMLEscapeString(field.Label)
		if field.Required {
			html += `<span class="required">*</span>`
		}
		html += `</label>`
	}

	// Widget
	html += `<div class="form-widget">`
	if field.Widget != "" {
		html += string(field.Widget)
	} else {
		// Default widget based on type
		html += f.renderDefaultWidget(field)
	}
	html += `</div>`

	// Help text
	if field.HelpText != "" {
		html += `<p class="form-help">` + template.HTMLEscapeString(field.HelpText) + `</p>`
	}

	// Errors
	if len(field.Errors) > 0 {
		html += `<ul class="field-errors">`
		for _, err := range field.Errors {
			html += `<li class="field-error">` + template.HTMLEscapeString(err) + `</li>`
		}
		html += `</ul>`
	}

	html += `</div>`
	return template.HTML(html)
}

// renderDefaultWidget renders a default widget based on field type
func (f *Form) renderDefaultWidget(field FormField) string {
	attrs := `name="` + template.HTMLEscapeString(field.Name) + `" id="id_` + template.HTMLEscapeString(field.Name) + `"`
	if field.Required {
		attrs += ` required`
	}
	if field.ReadOnly {
		attrs += ` readonly`
	}
	if field.Placeholder != "" {
		attrs += ` placeholder="` + template.HTMLEscapeString(field.Placeholder) + `"`
	}
	for k, v := range field.Attrs {
		attrs += ` ` + template.HTMLEscapeString(k) + `="` + template.HTMLEscapeString(v) + `"`
	}

	value := ""
	if field.Value != nil {
		value = fmt.Sprintf("%v", field.Value)
	}

	switch field.Type {
	case "textarea":
		return `<textarea ` + attrs + ` class="form-control">` + template.HTMLEscapeString(value) + `</textarea>`
	case "checkbox":
		checked := ""
		if value == "true" || value == "1" || value == "on" {
			checked = ` checked`
		}
		return `<input type="checkbox" ` + attrs + checked + ` class="form-check-input">`
	case "select":
		return `<select ` + attrs + ` class="form-select"></select>`
	default:
		return `<input type="` + template.HTMLEscapeString(field.Type) + `" ` + attrs + ` value="` + template.HTMLEscapeString(value) + `" class="form-control">`
	}
}

// renderAction renders a form action
func (f *Form) renderAction(action FormAction) template.HTML {
	if action.Type == "link" {
		html := `<a href="` + template.HTMLEscapeString(action.URL) + `" class="btn`
		if action.Primary {
			html += ` btn-primary`
		}
		if action.Class != "" {
			html += ` ` + template.HTMLEscapeString(action.Class)
		}
		html += `">`
		if action.Icon != "" {
			html += `<i class="icon-` + template.HTMLEscapeString(action.Icon) + `"></i> `
		}
		html += template.HTMLEscapeString(action.Label) + `</a>`
		return template.HTML(html)
	}

	html := `<button type="` + template.HTMLEscapeString(action.Type) + `" class="btn`
	if action.Primary {
		html += ` btn-primary`
	}
	if action.Class != "" {
		html += ` ` + template.HTMLEscapeString(action.Class)
	}
	html += `"`
	if action.Name != "" {
		html += ` name="` + template.HTMLEscapeString(action.Name) + `"`
	}
	if action.Value != "" {
		html += ` value="` + template.HTMLEscapeString(action.Value) + `"`
	}
	html += `>`
	if action.Icon != "" {
		html += `<i class="icon-` + template.HTMLEscapeString(action.Icon) + `"></i> `
	}
	html += template.HTMLEscapeString(action.Label) + `</button>`
	return template.HTML(html)
}
