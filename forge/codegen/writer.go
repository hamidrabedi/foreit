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

// WriteCombined writes all generated code to a single gen.go file
func (w *Writer) WriteCombined(definitions []*ModelDefinition, outputDir string) error {
	if len(definitions) == 0 {
		return nil
	}

	// Get package name from first model
	packageName := definitions[0].Package

	// Create template functions
	t := template.New("combined").Funcs(template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})

	t, err := t.Parse(combinedTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse combined template: %w", err)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Package": packageName,
		"Models":  definitions,
	}

	// Create gen.go file in the output directory
	filename := filepath.Join(outputDir, "gen.go")
	return w.writeTemplate(t, data, filename)
}

// WriteAPI writes all generated REST API code to an api_gen.go file
func (w *Writer) WriteAPI(definitions []*ModelDefinition, outputDir string) error {
	if len(definitions) == 0 {
		return nil
	}

	packageName := definitions[0].Package

	t := template.New("api").Funcs(template.FuncMap{
		"ToLower":  strings.ToLower,
		"ToSnake":  utils.ToSnake,
		"ToCamel":  utils.ToCamel,
		"ToPascal": utils.ToPascal,
	})

	t, err := t.Parse(apiTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse api template: %w", err)
	}

	data := map[string]interface{}{
		"Package": packageName,
		"Models":  definitions,
	}

	filename := filepath.Join(outputDir, "api_gen.go")
	return w.writeTemplate(t, data, filename)
}

// writeTemplate writes a template to a file. Replacement is atomic on Unix.
// On Windows, os.Rename is not guaranteed atomic and may fail while the
// destination is open elsewhere.
func (w *Writer) writeTemplate(t *template.Template, data interface{}, filename string) error {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)

	var perm os.FileMode
	if fi, err := os.Stat(filename); err == nil {
		perm = fi.Mode().Perm()
	} else {
		perm = os.FileMode(0o666 &^ processUmask())
	}

	tmpFile, err := os.CreateTemp(dir, "."+base+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	keepTemp := false
	defer func() {
		if !keepTemp {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	if err := tmpFile.Chmod(perm); err != nil {
		return fmt.Errorf("failed to set permissions on temp file: %w", err)
	}

	if err := t.Execute(tmpFile, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tmpPath, filename); err != nil {
		return fmt.Errorf("failed to rename temp file to %s: %w", filename, err)
	}

	keepTemp = true
	return nil
}
