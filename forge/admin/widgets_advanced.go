package admin

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

// FileInputWidget is a file upload widget
type FileInputWidget struct{}

// NewFileInput creates a new file input widget
func NewFileInput() Widget {
	return &FileInputWidget{}
}

// Type returns the widget type
func (w *FileInputWidget) Type() string {
	return "file"
}

// Render renders the file input widget
func (w *FileInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	html := fmt.Sprintf(`<input type="file" name="%s" id="id_%s"`, name, name)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control file-input">`
	
	// Show current file if exists
	if value != nil {
		if filePath, ok := value.(string); ok && filePath != "" {
			html += fmt.Sprintf(`<div class="current-file">Current: %s</div>`, template.HTMLEscapeString(filePath))
		}
	}
	
	return template.HTML(html)
}

// Parse parses the file input value
func (w *FileInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// ImageInputWidget is an image upload widget with preview
type ImageInputWidget struct{}

// NewImageInput creates a new image input widget
func NewImageInput() Widget {
	return &ImageInputWidget{}
}

// Type returns the widget type
func (w *ImageInputWidget) Type() string {
	return "image"
}

// Render renders the image input widget
func (w *ImageInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	html := fmt.Sprintf(`<input type="file" name="%s" id="id_%s" accept="image/*"`, name, name)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control image-input">`
	
	// Show current image preview if exists
	if value != nil {
		if imagePath, ok := value.(string); ok && imagePath != "" {
			html += fmt.Sprintf(`<div class="current-image"><img src="%s" alt="Current image" style="max-width: 200px; max-height: 200px;"></div>`, template.HTMLEscapeString(imagePath))
		}
	}
	
	return template.HTML(html)
}

// Parse parses the image input value
func (w *ImageInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// RichTextWidget is a rich text editor widget
type RichTextWidget struct {
	rows int
	cols int
}

// NewRichText creates a new rich text widget
func NewRichText(rows, cols int) Widget {
	return &RichTextWidget{
		rows: rows,
		cols: cols,
	}
}

// Type returns the widget type
func (w *RichTextWidget) Type() string {
	return "richtext"
}

// Render renders the rich text widget
func (w *RichTextWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
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
	html += ` class="form-control richtext-editor">`
	html += template.HTMLEscapeString(valueStr)
	html += `</textarea>`
	
	// Add rich text editor initialization script
	html += fmt.Sprintf(`<script>initRichTextEditor('id_%s');</script>`, name)
	
	return template.HTML(html)
}

// Parse parses the rich text value
func (w *RichTextWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// ColorInputWidget is a color picker widget
type ColorInputWidget struct{}

// NewColorInput creates a new color input widget
func NewColorInput() Widget {
	return &ColorInputWidget{}
}

// Type returns the widget type
func (w *ColorInputWidget) Type() string {
	return "color"
}

// Render renders the color input widget
func (w *ColorInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := "#000000"
	if value != nil {
		if str, ok := value.(string); ok {
			valueStr = str
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}
	}

	html := fmt.Sprintf(`<input type="color" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control color-input">`
	return template.HTML(html)
}

// Parse parses the color value
func (w *ColorInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// URLInputWidget is a URL input widget
type URLInputWidget struct{}

// NewURLInput creates a new URL input widget
func NewURLInput() Widget {
	return &URLInputWidget{}
}

// Type returns the widget type
func (w *URLInputWidget) Type() string {
	return "url"
}

// Render renders the URL input widget
func (w *URLInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	html := fmt.Sprintf(`<input type="url" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control url-input">`
	return template.HTML(html)
}

// Parse parses the URL value
func (w *URLInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// PasswordInputWidget is a password input widget
type PasswordInputWidget struct{}

// NewPasswordInput creates a new password input widget
func NewPasswordInput() Widget {
	return &PasswordInputWidget{}
}

// Type returns the widget type
func (w *PasswordInputWidget) Type() string {
	return "password"
}

// Render renders the password input widget
func (w *PasswordInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	html := fmt.Sprintf(`<input type="password" name="%s" id="id_%s"`, name, name)
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control password-input">`
	return template.HTML(html)
}

// Parse parses the password value
func (w *PasswordInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// HiddenInputWidget is a hidden input widget
type HiddenInputWidget struct{}

// NewHiddenInput creates a new hidden input widget
func NewHiddenInput() Widget {
	return &HiddenInputWidget{}
}

// Type returns the widget type
func (w *HiddenInputWidget) Type() string {
	return "hidden"
}

// Render renders the hidden input widget
func (w *HiddenInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	html := fmt.Sprintf(`<input type="hidden" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += `>`
	return template.HTML(html)
}

// Parse parses the hidden input value
func (w *HiddenInputWidget) Parse(value string) (interface{}, error) {
	return value, nil
}

// TimeInputWidget is a time input widget
type TimeInputWidget struct{}

// NewTimeInput creates a new time input widget
func NewTimeInput() Widget {
	return &TimeInputWidget{}
}

// Type returns the widget type
func (w *TimeInputWidget) Type() string {
	return "time"
}

// Render renders the time input widget
func (w *TimeInputWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	valueStr := ""
	if value != nil {
		if t, ok := value.(time.Time); ok {
			if !t.IsZero() {
				valueStr = t.Format("15:04")
			}
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}
	}

	html := fmt.Sprintf(`<input type="time" name="%s" id="id_%s" value="%s"`, name, name, template.HTMLEscapeString(valueStr))
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	html += ` class="form-control time-input">`
	return template.HTML(html)
}

// Parse parses the time value
func (w *TimeInputWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}
	// Parse time in HH:MM format
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid time format")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	return time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC), nil
}
