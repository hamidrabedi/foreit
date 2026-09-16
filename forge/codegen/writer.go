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

	// Generated REST APIs address rows by integer IDs; reject models with
	// non-integer primary keys instead of emitting broken CRUD.
	if err := ValidateAPIModels(definitions); err != nil {
		return err
	}

	packageName := definitions[0].Package

	t := template.New("api").Funcs(template.FuncMap{
		"ToLower":       strings.ToLower,
		"ToSnake":       utils.ToSnake,
		"Pluralize":     utils.Pluralize,
		"ToKebab":       utils.ToKebab,
		"ToKebabPlural": utils.ToKebabPlural,
		"ToCamel":       utils.ToCamel,
		"ToPascal":      utils.ToPascal,
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
	writeFilename := filename
	if fi, err := os.Lstat(filename); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		resolved, err := filepath.EvalSymlinks(filename)
		if err != nil {
			// A dangling symlink (its target does not exist yet) cannot be
			// fully evaluated. Resolve the existing parent portion of the
			// link and keep the final target name, so generation writes
			// through the link and creates the target instead of aborting.
			resolved, err = resolveDanglingSymlink(filename)
			if err != nil {
				return fmt.Errorf("failed to resolve destination symlink %s: %w", filename, err)
			}
		}
		writeFilename = resolved
	}
	if err := ensureDestinationWritable(writeFilename); err != nil {
		return err
	}
	dir := filepath.Dir(writeFilename)
	base := filepath.Base(writeFilename)

	var perm os.FileMode
	if fi, err := os.Stat(writeFilename); err == nil {
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

	if err := os.Rename(tmpPath, writeFilename); err != nil {
		return fmt.Errorf("failed to rename temp file to %s: %w", writeFilename, err)
	}

	keepTemp = true
	return nil
}

// ensureDestinationWritable fails when an existing destination file is not
// writable by the current process, so a file deliberately marked read-only is
// never silently replaced through the temp-file-and-rename path (which only
// requires a writable parent directory). A destination that does not exist is
// left alone. Probing with O_WRONLY (no truncation, no creation) reflects the
// effective permissions of the current process, including ownership.
func ensureDestinationWritable(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("refusing to overwrite read-only file %s: %w", path, err)
	}
	return f.Close()
}

// resolveDanglingSymlink resolves filename, which must be a symlink whose
// target does not exist yet, by evaluating the existing parent portion of the
// link target while keeping the final target name. The result still passes
// through the link, so writing to it creates the target.
func resolveDanglingSymlink(filename string) (string, error) {
	target, err := os.Readlink(filename)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(filename), target)
	}
	parent := filepath.Dir(target)
	resolvedParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return "", err
	}
	return filepath.Join(resolvedParent, filepath.Base(target)), nil
}
