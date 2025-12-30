package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Engine manages template loading and rendering
type Engine struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
	baseDir   string
	funcMap   template.FuncMap
}

// NewEngine creates a new template engine
func NewEngine(baseDir string) *Engine {
	engine := &Engine{
		templates: make(map[string]*template.Template),
		baseDir:   baseDir,
		funcMap:   make(template.FuncMap),
	}

	// Register default template functions
	engine.registerDefaultFuncs()

	return engine
}

// registerDefaultFuncs registers default template functions
func (e *Engine) registerDefaultFuncs() {
	e.funcMap["eq"] = func(a, b interface{}) bool {
		return a == b
	}
	e.funcMap["ne"] = func(a, b interface{}) bool {
		return a != b
	}
	e.funcMap["gt"] = func(a, b int) bool {
		return a > b
	}
	e.funcMap["lt"] = func(a, b int) bool {
		return a < b
	}
	e.funcMap["add"] = func(a, b int) int {
		return a + b
	}
	e.funcMap["sub"] = func(a, b int) int {
		return a - b
	}
	e.funcMap["mul"] = func(a, b int) int {
		return a * b
	}
	e.funcMap["div"] = func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a / b
	}
	e.funcMap["mod"] = func(a, b int) int {
		if b == 0 {
			return 0
		}
		return a % b
	}
	e.funcMap["default"] = func(def, val interface{}) interface{} {
		if val == nil || val == "" {
			return def
		}
		return val
	}
	e.funcMap["empty"] = func(val interface{}) bool {
		if val == nil {
			return true
		}
		switch v := val.(type) {
		case string:
			return v == ""
		case []interface{}:
			return len(v) == 0
		case map[string]interface{}:
			return len(v) == 0
		}
		return false
	}
	e.funcMap["escape"] = template.HTMLEscapeString
	e.funcMap["safe"] = func(s string) template.HTML {
		return template.HTML(s)
	}
}

// RegisterFunc registers a custom template function
func (e *Engine) RegisterFunc(name string, fn interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.funcMap[name] = fn
}

// LoadTemplate loads a template from file
func (e *Engine) LoadTemplate(name, filePath string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	fullPath := filepath.Join(e.baseDir, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", fullPath, err)
	}

	tmpl, err := template.New(name).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	e.templates[name] = tmpl
	return nil
}

// LoadTemplateString loads a template from string
func (e *Engine) LoadTemplateString(name, content string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	tmpl, err := template.New(name).Funcs(e.funcMap).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	e.templates[name] = tmpl
	return nil
}

// GetTemplate gets a template by name
func (e *Engine) GetTemplate(name string) (*template.Template, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tmpl, ok := e.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}

	return tmpl, nil
}

// Render renders a template with data
func (e *Engine) Render(w io.Writer, name string, data interface{}) error {
	tmpl, err := e.GetTemplate(name)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

// RenderString renders a template to a string
func (e *Engine) RenderString(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := e.Render(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MustRender renders a template and panics on error
func (e *Engine) MustRender(w io.Writer, name string, data interface{}) {
	if err := e.Render(w, name, data); err != nil {
		panic(fmt.Sprintf("template render error: %v", err))
	}
}
