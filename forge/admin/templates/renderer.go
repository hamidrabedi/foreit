package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Renderer provides high-level template rendering with layout support
type Renderer struct {
	engine *Engine
	layout string
}

// NewRenderer creates a new template renderer
func NewRenderer(engine *Engine) *Renderer {
	return &Renderer{
		engine: engine,
		layout: "base",
	}
}

// SetLayout sets the default layout template
func (r *Renderer) SetLayout(layout string) {
	r.layout = layout
}

// Render renders a template with layout
func (r *Renderer) Render(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	// Get the content template
	contentTmpl, err := r.engine.GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template %s: %w", templateName, err)
	}

	// Render content to buffer
	var contentBuf bytes.Buffer
	if err := contentTmpl.Execute(&contentBuf, data); err != nil {
		return fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	// If no layout, just render content
	if r.layout == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := io.Copy(w, &contentBuf)
		return err
	}

	// Get layout template
	layoutTmpl, err := r.engine.GetTemplate(r.layout)
	if err != nil {
		// If layout not found, just render content
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := io.Copy(w, &contentBuf)
		return err
	}

	// Add content to data
	data["Content"] = template.HTML(contentBuf.String())

	// Render layout with content
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return layoutTmpl.Execute(w, data)
}

// RenderPartial renders a template without layout
func (r *Renderer) RenderPartial(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.engine.Render(w, templateName, data)
}

// RenderString renders a template to a string
func (r *Renderer) RenderString(templateName string, data map[string]interface{}) (string, error) {
	return r.engine.RenderString(templateName, data)
}

// LoadTemplates loads all templates from a directory
func (r *Renderer) LoadTemplates(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Only load .html and .tmpl files
		ext := filepath.Ext(path)
		if ext != ".html" && ext != ".tmpl" {
			return nil
		}

		// Get relative path from base dir
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// Use relative path as template name (without extension)
		name := relPath[:len(relPath)-len(ext)]

		return r.engine.LoadTemplate(name, relPath)
	})
}
