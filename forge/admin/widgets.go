package admin

import (
	"fmt"
	"html/template"
	"strconv"
	"time"
)

// Widget represents a form widget
type Widget interface {
	Render(name string, value interface{}, attrs map[string]string) template.HTML
	Parse(value string) (interface{}, error)
	Type() string
}

// TextInputWidget is a text input widget
type TextInputWidget struct{}

// NewTextInput creates a new text input widget
func NewTextInput() Widget {
	return &TextInputWidget{}
}

// Type returns the widget type
func (w *TextInputWidget) Type() string {
	return "text"
}

// Render renders the widget
func (w *TextInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	html := fmt.Sprintf(`<input type="text" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control">`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *TextInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// NumberInputWidget is a number input widget
type NumberInputWidget struct{}

// NewNumberInput creates a new number input widget
func NewNumberInput() Widget {
	return &NumberInputWidget{}
}

// Type returns the widget type
func (w *NumberInputWidget) Type() string {
	return "number"
}

// Render renders the widget
func (w *NumberInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	html := fmt.Sprintf(`<input type="number" name="%s" id="id_%s" value="%s"`, name, name, valueStr)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control">`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *NumberInputWidget) Parse(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}

// EmailInputWidget is an email input widget
type EmailInputWidget struct{}

// NewEmailInput creates a new email input widget
func NewEmailInput() Widget {
	return &EmailInputWidget{}
}

// Type returns the widget type
func (w *EmailInputWidget) Type() string {
	return "email"
}

// Render renders the widget
func (w *EmailInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	html := fmt.Sprintf(`<input type="email" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control">`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *EmailInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// TextareaWidget is a textarea widget
type TextareaWidget struct {
	rows int
	cols int
}

// NewTextarea creates a new textarea widget
func NewTextarea(rows, cols int) Widget {
	return &TextareaWidget{
		rows: rows,
		cols: cols,
	}
}

// Type returns the widget type
func (w *TextareaWidget) Type() string {
	return "textarea"
}

// Render renders the widget
func (w *TextareaWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	rows := w.rows
	if rows == 0 {
		rows = 10
	}
	cols := w.cols
	if cols == 0 {
		cols = 40
	}

	html := fmt.Sprintf(`<textarea name="%s" id="id_%s" rows="%d" cols="%d"`, name, name, rows, cols)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control">`
	html += template.HTMLEscapeString(valueStr)
	html += `</textarea>`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *TextareaWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// CheckboxWidget is a checkbox widget
type CheckboxWidget struct{}

// NewCheckbox creates a new checkbox widget
func NewCheckbox() Widget {
	return &CheckboxWidget{}
}

// Type returns the widget type
func (w *CheckboxWidget) Type() string {
	return "checkbox"
}

// Render renders the widget
func (w *CheckboxWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	checked := ""
	if value != nil {
		if boolVal, ok := value.(bool); ok && boolVal {
			checked = " checked"
		} else if strVal, ok := value.(string); ok && (strVal == "true" || strVal == "on" || strVal == "1") {
			checked = " checked"
		}
	}

	html := fmt.Sprintf(`<input type="checkbox" name="%s" id="id_%s"%s`, name, name, checked)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-check-input">`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *CheckboxWidget) Parse(value string) (interface{}, error) {
	return value == "on" || value == "true" || value == "1", nil
}

// SelectWidget is a select widget
type SelectWidget struct {
	choices []Choice[interface{}]
}

// NewSelect creates a new select widget
func NewSelect(choices []Choice[interface{}]) Widget {
	return &SelectWidget{
		choices: choices,
	}
}

// Type returns the widget type
func (w *SelectWidget) Type() string {
	return "select"
}

// Render renders the widget
func (w *SelectWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	html := fmt.Sprintf(`<select name="%s" id="id_%s"`, name, name)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-select">`
	html += `<option value="">---------</option>`

	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	for _, choice := range w.choices {
		choiceValue := fmt.Sprintf("%v", choice.Value)
		selected := ""
		if choiceValue == valueStr {
			selected = " selected"
		}
		html += fmt.Sprintf(`<option value="%s"%s>%s</option>`,
			template.HTMLEscapeString(choiceValue),
			selected,
			template.HTMLEscapeString(choice.Label),
		)
	}

	html += `</select>`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *SelectWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// DateInputWidget is a date input widget
type DateInputWidget struct{}

// NewDateInput creates a new date input widget
func NewDateInput() Widget {
	return &DateInputWidget{}
}

// Type returns the widget type
func (w *DateInputWidget) Type() string {
	return "date"
}

// Render renders the widget
func (w *DateInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		if t, ok := value.(time.Time); ok {
			if !t.IsZero() {
				valueStr = t.Format("2006-01-02")
			}
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}
	}

	html := fmt.Sprintf(`<input type="date" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control">`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *DateInputWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", value)
}

// DateTimeInputWidget is a datetime input widget
type DateTimeInputWidget struct{}

// NewDateTimeInput creates a new datetime input widget
func NewDateTimeInput() Widget {
	return &DateTimeInputWidget{}
}

// Type returns the widget type
func (w *DateTimeInputWidget) Type() string {
	return "datetime-local"
}

// Render renders the widget
func (w *DateTimeInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		if t, ok := value.(time.Time); ok {
			if !t.IsZero() {
				valueStr = t.Format("2006-01-02T15:04")
			}
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}
	}

	html := fmt.Sprintf(`<input type="datetime-local" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control">`
	return template.HTML(html)
}

// Parse parses the widget value
func (w *DateTimeInputWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02T15:04", value)
}
