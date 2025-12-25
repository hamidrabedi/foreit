package admin

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TemplateEngine manages template rendering with override support
type TemplateEngine struct {
	basePath  string
	overrides map[string]string // model -> override path
	templates map[string]*template.Template
	registry  *RendererRegistry
	mu        sync.RWMutex
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(basePath string, registry *RendererRegistry) *TemplateEngine {
	return &TemplateEngine{
		basePath:  basePath,
		overrides: make(map[string]string),
		templates: make(map[string]*template.Template),
		registry:  registry,
	}
}

// LoadTemplates loads templates from the base path
func (e *TemplateEngine) LoadTemplates() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Load default templates
	defaultPath := filepath.Join(e.basePath, "defaults")
	if err := e.loadTemplatesFromPath(defaultPath, "defaults"); err != nil {
		return err
	}

	// Load override templates
	overridePath := filepath.Join(e.basePath, "overrides")
	if err := e.loadTemplatesFromPath(overridePath, "overrides"); err != nil {
		// Overrides are optional, so log but don't fail
		fmt.Printf("Warning: Could not load override templates: %v\n", err)
	}

	return nil
}

// loadTemplatesFromPath loads templates from a directory
func (e *TemplateEngine) loadTemplatesFromPath(path, prefix string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Path doesn't exist, skip
	}

	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(filePath, ".html") {
			return nil
		}

		relPath, err := filepath.Rel(path, filePath)
		if err != nil {
			return err
		}

		tmpl, err := template.New(info.Name()).Funcs(e.getFuncMap()).ParseFiles(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", filePath, err)
		}

		key := fmt.Sprintf("%s/%s", prefix, relPath)
		e.templates[key] = tmpl

		return nil
	})
}

// RenderTemplate renders a template with override support
func (e *TemplateEngine) RenderTemplate(model, view string, data interface{}) (string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// 1. Check for override
	overrideKey := fmt.Sprintf("overrides/%s/%s.html", model, view)
	if tmpl, ok := e.templates[overrideKey]; ok {
		return e.executeTemplate(tmpl, data)
	}

	// 2. Use default
	defaultKey := fmt.Sprintf("defaults/%s.html", view)
	if tmpl, ok := e.templates[defaultKey]; ok {
		return e.executeTemplate(tmpl, data)
	}

	return "", fmt.Errorf("template not found: %s/%s", model, view)
}

// executeTemplate executes a template with data
func (e *TemplateEngine) executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}
	return buf.String(), nil
}

// getFuncMap returns template functions
func (e *TemplateEngine) getFuncMap() template.FuncMap {
	return template.FuncMap{
		"renderField": e.renderField,
		"htmxAttrs":   e.buildHTMXAttrs,
		"safeHTML":    func(s string) template.HTML { return template.HTML(s) },
		"add":         func(a, b int) int { return a + b },
		"sub":         func(a, b int) int { return a - b },
		"mul":         func(a, b int) int { return a * b },
		"div":         func(a, b int) int { if b == 0 { return 0 }; return a / b },
		"int":         func(v interface{}) int {
			switch val := v.(type) {
			case int:
				return val
			case int64:
				return int(val)
			case int32:
				return int(val)
			default:
				return 0
			}
		},
		"index": func(m map[string]interface{}, key string) interface{} { return m[key] },
	}
}

// renderField is a template helper for rendering fields
func (e *TemplateEngine) renderField(field *FieldMeta, value interface{}, htmxAttrs map[string]string) template.HTML {
	renderer, err := e.registry.GetRenderer(string(field.Type))
	if err != nil {
		return template.HTML(fmt.Sprintf("<!-- Error: %v -->", err))
	}

	ctx := RenderContext{
		Field:     field,
		Value:     value,
		HTMXAttrs: htmxAttrs,
	}

	html, err := renderer.RenderHTML(ctx)
	if err != nil {
		return template.HTML(fmt.Sprintf("<!-- Error: %v -->", err))
	}

	return html
}

// buildHTMXAttrs is a template helper for building HTMX attributes
func (e *TemplateEngine) buildHTMXAttrs(attrs map[string]string) string {
	return buildHTMXAttrs(attrs)
}

// RegisterOverride registers a template override path for a model
func (e *TemplateEngine) RegisterOverride(model, overridePath string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.overrides[model] = overridePath
}

