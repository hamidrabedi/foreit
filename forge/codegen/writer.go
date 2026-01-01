package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/forgego/forge/utils"
)

// Writer writes generated code to files
type Writer struct {
	templates map[string]*template.Template
}

// NewWriter creates a new writer
func NewWriter() *Writer {
	return &Writer{
		templates: make(map[string]*template.Template),
	}
}

// WriteModel writes the generated model struct
func (w *Writer) WriteModel(def *ModelDefinition, outputDir string) error {
	// Create gen subdirectory
	genDir := filepath.Join(outputDir, "gen")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return fmt.Errorf("failed to create gen directory: %w", err)
	}

	t := template.New("model").Funcs(template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})
	t, err := t.Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse model template: %w", err)
	}

	// Prepare template data with Meta accessible
	data := map[string]interface{}{
		"Package": def.Package, // Keep original package for imports if needed
		"Name":    def.Name,
		"Fields":  def.Fields,
		"Meta":    def.Meta,
	}

	filename := filepath.Join(genDir, utils.ToSnake(def.Name)+".gen.go")
	return w.writeTemplate(t, data, filename)
}

// WriteFields writes the generated FieldExpr definitions
func (w *Writer) WriteFields(def *ModelDefinition, outputDir string) error {
	// Create gen subdirectory
	genDir := filepath.Join(outputDir, "gen")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return fmt.Errorf("failed to create gen directory: %w", err)
	}

	// Add custom template function for ToLower
	t := template.New("fields").Funcs(template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})

	t, err := t.Parse(fieldsTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse fields template: %w", err)
	}

	filename := filepath.Join(genDir, utils.ToSnake(def.Name)+"_fields.gen.go")
	return w.writeTemplate(t, def, filename)
}

// WriteRelations writes the generated RelationExpr definitions
func (w *Writer) WriteRelations(def *ModelDefinition, outputDir string) error {
	// Create gen subdirectory
	genDir := filepath.Join(outputDir, "gen")
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return fmt.Errorf("failed to create gen directory: %w", err)
	}

	// Generate RelationExpr file
	filename := filepath.Join(genDir, utils.ToSnake(def.Name)+"_relations.gen.go")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create relations file: %w", err)
	}
	defer f.Close()

	t := template.New("relations").Funcs(template.FuncMap{
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})

	t, err = t.Parse(relationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse relation template: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Package":   def.Package,
		"ModelName": def.Name,
		"Relations": def.Relations,
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute relation template: %w", err)
	}

	return nil
}

// WriteManager writes the generated Manager
func (w *Writer) WriteManager(def *ModelDefinition, outputDir string) error {
	t := template.New("manager").Funcs(template.FuncMap{
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})

	t, err := t.Parse(managerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse manager template: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Package":   def.Package,
		"Name":      def.Name,
		"TableName": def.Meta.TableName,
	}

	if data["TableName"] == "" {
		data["TableName"] = utils.ToSnake(def.Name) + "s"
	}

	filename := filepath.Join(outputDir, utils.ToSnake(def.Name)+"_manager.gen.go")
	return w.writeTemplate(t, data, filename)
}

// WriteQuerySet writes the generated QuerySet
func (w *Writer) WriteQuerySet(def *ModelDefinition, outputDir string) error {
	t := template.New("queryset").Funcs(template.FuncMap{
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})

	t, err := t.Parse(querysetTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse queryset template: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Package":   def.Package,
		"Name":      def.Name,
		"TableName": def.Meta.TableName,
	}

	if data["TableName"] == "" {
		data["TableName"] = utils.ToSnake(def.Name) + "s"
	}

	filename := filepath.Join(outputDir, utils.ToSnake(def.Name)+"_queryset.gen.go")
	return w.writeTemplate(t, data, filename)
}

// writeTemplate writes a template to a file
func (w *Writer) writeTemplate(t *template.Template, data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	if err := t.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
