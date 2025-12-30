package widgets

import (
	"fmt"
	"html/template"
)

// FileUploadWidget is a file upload widget
type FileUploadWidget struct {
	Accept      string // MIME types (e.g., "image/*", ".pdf")
	Multiple    bool
	MaxSize     int64  // Maximum file size in bytes
	MaxFiles    int    // Maximum number of files
}

// NewFileUpload creates a new file upload widget
func NewFileUpload() *FileUploadWidget {
	return &FileUploadWidget{
		Accept:   "*/*",
		Multiple: false,
		MaxSize:  10 * 1024 * 1024, // 10MB default
		MaxFiles: 1,
	}
}

// WithAccept sets accepted MIME types
func (w *FileUploadWidget) WithAccept(accept string) *FileUploadWidget {
	w.Accept = accept
	return w
}

// WithMultiple allows multiple file selection
func (w *FileUploadWidget) WithMultiple(multiple bool) *FileUploadWidget {
	w.Multiple = multiple
	return w
}

// WithMaxSize sets maximum file size
func (w *FileUploadWidget) WithMaxSize(maxSize int64) *FileUploadWidget {
	w.MaxSize = maxSize
	return w
}

// Type returns the widget type
func (w *FileUploadWidget) Type() string {
	return "file"
}

// Render renders the widget
func (w *FileUploadWidget) Render(name string, value interface{}, attrs map[string]string) template.HTML {
	html := fmt.Sprintf(`<input type="file" name="%s" id="id_%s"`, name, name)
	
	if w.Accept != "" {
		html += fmt.Sprintf(` accept="%s"`, template.HTMLEscapeString(w.Accept))
	}
	
	if w.Multiple {
		html += ` multiple`
	}
	
	html += fmt.Sprintf(` data-max-size="%d"`, w.MaxSize)
	html += fmt.Sprintf(` data-max-files="%d"`, w.MaxFiles)
	
	for k, v := range attrs {
		html += fmt.Sprintf(` %s="%s"`, k, template.HTMLEscapeString(v))
	}
	
	html += ` class="form-control file-upload">`
	
	// Add file preview if value exists
	if value != nil {
		if filePath, ok := value.(string); ok && filePath != "" {
			html += fmt.Sprintf(`<div class="file-preview"><a href="%s" target="_blank">Current file</a></div>`, template.HTMLEscapeString(filePath))
		}
	}
	
	return template.HTML(html)
}

// Parse parses the widget value (would be handled by form processing)
func (w *FileUploadWidget) Parse(value string) (interface{}, error) {
	// File uploads are handled via multipart form data, not string parsing
	return value, nil
}
