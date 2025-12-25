package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewStartAppCommand creates a new startapp command
func NewStartAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "startapp [app-name]",
		Short: "Create a new app in the current project",
		Long:  "Creates a new app (like Django's startapp) with models and admin registration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			return createApp(appName)
		},
	}

	return cmd
}

func createApp(appName string) error {
	// Validate app name
	if !isValidProjectName(appName) {
		return fmt.Errorf("invalid app name: %s (must be a valid Go package name)", appName)
	}

	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	appPath := filepath.Join(wd, "internal", "apps", appName)

	// Check if directory exists
	if _, err := os.Stat(appPath); err == nil {
		return fmt.Errorf("app %s already exists", appName)
	}

	// Create app directory structure
	dirs := []string{
		"models/ent/schema",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(appPath, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create example model schema
	schemaContent := fmt.Sprintf(`package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// %s represents a %s model
type %s struct {
	ent.Schema
}

func (%s) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.Time("created_at").Default(func() time.Time {
			return time.Now()
		}),
		field.Time("updated_at").Default(func() time.Time {
			return time.Now()
		}).UpdateDefault(func() time.Time {
			return time.Now()
		}),
	}
}
`, strings.Title(appName), appName, strings.Title(appName), strings.Title(appName))

	if err := os.WriteFile(filepath.Join(appPath, "models/ent/schema", strings.ToLower(appName)+".go"), []byte(schemaContent), 0644); err != nil {
		return fmt.Errorf("failed to create schema file: %w", err)
	}

	// Create admin registration file
	adminContent := fmt.Sprintf(`package %s

import (
	"github.com/gogo/pkg/admin"
)

// RegisterAdmin registers %s models with the admin engine
func RegisterAdmin(engine *admin.Engine) {
	// Register your models here
	// engine.Register(&models.%s{})
}
`, appName, appName, strings.Title(appName))

	if err := os.WriteFile(filepath.Join(appPath, "admin.go"), []byte(adminContent), 0644); err != nil {
		return fmt.Errorf("failed to create admin.go: %w", err)
	}

	fmt.Printf("✅ Created app '%s' at %s\n", appName, appPath)
	fmt.Println("\nNext steps:")
	fmt.Printf("  # Edit models in internal/apps/%s/models/ent/schema/\n", appName)
	fmt.Println("  # Generate Ent code")
	fmt.Printf("  # Register in admin: import and call %s.RegisterAdmin(adminEngine)\n", appName)

	return nil
}

