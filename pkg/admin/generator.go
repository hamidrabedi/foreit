package admin

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/forgego/forge/pkg/admin/templates"
)

// AdminGenerator generates admin views and templates
type AdminGenerator struct {
	templates *template.Template
}

// NewAdminGenerator creates a new admin generator
func NewAdminGenerator() *AdminGenerator {
	tmpl := template.New("admin").Funcs(templates.FuncMap())
	return &AdminGenerator{
		templates: tmpl,
	}
}

// GenerateViews generates admin views for all registered models
func (g *AdminGenerator) GenerateViews(outputDir string) error {
	models := GetAllModels()
	
	for name, model := range models {
		if err := g.generateModelViews(name, model, outputDir); err != nil {
			return fmt.Errorf("failed to generate views for %s: %w", name, err)
		}
	}
	
	return nil
}

// generateModelViews generates views for a single model
func (g *AdminGenerator) generateModelViews(name string, model *AdminModel, outputDir string) error {
	// TODO: Implement view generation
	// - List view
	// - Detail view
	// - Create view
	// - Update view
	// - Delete view
	// - Forms
	// - Templates
	
	return nil
}

// GenerateTemplates generates HTML templates for admin interface
func (g *AdminGenerator) GenerateTemplates(outputDir string) error {
	// Parse base template
	basePath := filepath.Join(outputDir, "base.html")
	_, err := template.New("base.html").Funcs(templates.FuncMap()).ParseFiles(basePath)
	if err != nil {
		// Base template might not exist yet - that's okay
		// Will be created by admin system
	}
	
	// Generate list template (example - full implementation would be more complex)
	// TODO: Implement full template generation with proper field mapping
	
	// TODO: Generate detail and form templates
	
	return nil
}

