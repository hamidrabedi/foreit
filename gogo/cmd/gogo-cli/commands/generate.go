package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate code (resource, console, policy, etc.)",
	Long: `Generate code for resources, consoles, policies, and more.

Examples:
  gogo generate resource User
  gogo generate console User
  gogo generate policy User
  gogo generate serializer User`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		genType := args[0]
		name := args[1]
		
		switch genType {
		case "resource":
			generateResource(name)
		case "console":
			generateConsole(name)
		case "policy":
			generatePolicy(name)
		case "serializer":
			generateSerializer(name)
		default:
			fmt.Printf("Unknown type: %s\n", genType)
			fmt.Println("Available types: resource, console, policy, serializer")
			os.Exit(1)
		}
	},
}

func generateResource(name string) {
	tmpl := `package resources

import (
	"context"
	
	"github.com/gogo/pkg/endpoints"
	"github.com/gogo/pkg/orm"
	"yourproject/internal/models"
	"yourproject/internal/models/ent"
)

type {{.Name}}Resource struct {
	*endpoints.BaseResource[models.{{.Name}}, *ent.{{.Name}}Query]
}

func New{{.Name}}Resource(client *orm.Client) *{{.Name}}Resource {
	return &{{.Name}}Resource{
		BaseResource: endpoints.NewResource[models.{{.Name}}, *ent.{{.Name}}Query](client),
	}
}

func (r *{{.Name}}Resource) Index(ctx *endpoints.Context) ([]*models.{{.Name}}, error) {
	return r.Query().All(ctx.Request.Context())
}

func (r *{{.Name}}Resource) Show(ctx *endpoints.Context, id interface{}) (*models.{{.Name}}, error) {
	return r.Query().Where(ent.{{.Name}}IDEQ(id.(int))).Only(ctx.Request.Context())
}

func (r *{{.Name}}Resource) Create(ctx *endpoints.Context) (*models.{{.Name}}, error) {
	// Implement create logic
	return nil, nil
}

func (r *{{.Name}}Resource) Update(ctx *endpoints.Context, id interface{}) (*models.{{.Name}}, error) {
	// Implement update logic
	return nil, nil
}

func (r *{{.Name}}Resource) Destroy(ctx *endpoints.Context, id interface{}) error {
	// Implement delete logic
	return nil
}
`
	
	generateFile("resources", name, tmpl)
}

func generateConsole(name string) {
	tmpl := `package consoles

import (
	"github.com/gogo/pkg/console"
	"yourproject/internal/models"
)

type {{.Name}}Console struct {
	*console.BaseConsole[models.{{.Name}}]
}

func New{{.Name}}Console() *{{.Name}}Console {
	return &{{.Name}}Console{
		BaseConsole: console.New[models.{{.Name}}](&console.Options{
			ListDisplay:  []string{"id", "name", "created_at"},
			SearchFields: []string{"name"},
		}),
	}
}
`
	
	generateFile("consoles", name, tmpl)
}

func generatePolicy(name string) {
	tmpl := `package policies

import (
	"context"
	
	"github.com/gogo/pkg/auth"
	"yourproject/internal/models"
)

type {{.Name}}Policy struct {
	auth.BasePolicy[models.{{.Name}}]
}

func (p *{{.Name}}Policy) CanView(ctx context.Context, user interface{}, obj *models.{{.Name}}) bool {
	// Implement view permission
	return true
}

func (p *{{.Name}}Policy) CanCreate(ctx context.Context, user interface{}) bool {
	// Implement create permission
	return true
}

func (p *{{.Name}}Policy) CanUpdate(ctx context.Context, user interface{}, obj *models.{{.Name}}) bool {
	// Implement update permission
	return true
}

func (p *{{.Name}}Policy) CanDelete(ctx context.Context, user interface{}, obj *models.{{.Name}}) bool {
	// Implement delete permission
	return true
}
`
	
	generateFile("policies", name, tmpl)
}

func generateSerializer(name string) {
	tmpl := `package serializers

import (
	"github.com/gogo/pkg/endpoints"
	"yourproject/internal/models"
)

type {{.Name}}Serializer struct {
	*endpoints.BaseSerializer[models.{{.Name}}]
}

func New{{.Name}}Serializer() *{{.Name}}Serializer {
	return &{{.Name}}Serializer{
		BaseSerializer: endpoints.NewSerializer[models.{{.Name}}](),
	}
}

func (s *{{.Name}}Serializer) ToRepresentation(obj *models.{{.Name}}) map[string]interface{} {
	return map[string]interface{}{
		"id":   obj.ID,
		// Add fields
	}
}

func (s *{{.Name}}Serializer) ToInternal(data map[string]interface{}) *models.{{.Name}} {
	// Convert data to model
	return nil
}
`
	
	generateFile("serializers", name, tmpl)
}

func generateFile(dir, name, tmpl string) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}
	
	// Create file
	filename := filepath.Join(dir, fmt.Sprintf("%s.go", toSnakeCase(name)))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	
	// Parse and execute template
	t, err := template.New("code").Parse(tmpl)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}
	
	data := struct {
		Name string
	}{
		Name: name,
	}
	
	if err := t.Execute(file, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Generated %s\n", filename)
}

// RegisterGenerateCommand registers the generate command
func RegisterGenerateCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(generateCmd)
}

func toSnakeCase(s string) string {
	// Simple conversion - in production use a proper library
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return result
}

