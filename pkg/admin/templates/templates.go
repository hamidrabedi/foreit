package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html *.tmpl
var templateFS embed.FS

// LoadTemplates loads all templates from the embedded filesystem
func LoadTemplates() (*template.Template, error) {
	// Parse base template first
	baseContent, err := templateFS.ReadFile("base.html")
	if err != nil {
		return nil, err
	}

	// Create template with functions and parse base
	tmpl := template.New("base").Funcs(FuncMap())
	tmpl, err = tmpl.Parse(string(baseContent))
	if err != nil {
		return nil, err
	}

	// Parse child templates
	files, err := templateFS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || file.Name() == "base.html" {
			continue
		}

		content, err := templateFS.ReadFile(file.Name())
		if err != nil {
			continue // Skip files that can't be read
		}

		// Get template name without extension
		templateName := file.Name()
		if len(templateName) > 5 && templateName[len(templateName)-5:] == ".tmpl" {
			templateName = templateName[:len(templateName)-5]
		} else if len(templateName) > 5 && templateName[len(templateName)-5:] == ".html" {
			templateName = templateName[:len(templateName)-5]
		}

		// Parse as associated template
		_, err = tmpl.New(templateName).Parse(string(content))
		if err != nil {
			continue // Skip templates that can't be parsed
		}
	}

	return tmpl, nil
}
