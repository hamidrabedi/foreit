package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewStartProjectCommand creates a new startproject command
func NewStartProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "startproject [project-name]",
		Short: "Create a new Go Admin project",
		Long:  "Creates a new Go Admin project with the specified name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			return createProject(projectName)
		},
	}

	return cmd
}

func createProject(projectName string) error {
	// Validate project name
	if !isValidProjectName(projectName) {
		return fmt.Errorf("invalid project name: %s (must be a valid Go package name)", projectName)
	}

	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectPath := filepath.Join(wd, projectName)

	// Check if directory exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists", projectName)
	}

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create directory structure
	dirs := []string{
		"cmd/web",
		"internal/models/ent/schema",
		"internal/admin",
		"migrations",
		"docs",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectPath, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.18

require (
	entgo.io/ent v0.13.0
	github.com/gogo v0.1.0
	github.com/gofiber/fiber/v2 v2.40.1
)
`, projectName)

	if err := os.WriteFile(filepath.Join(projectPath, "go.mod"), []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	// Create main.go
	mainContent := `package main

import (
	"log"

	"github.com/gogo/pkg/gogo"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app, err := gogo.New(&gogo.AppConfig{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost/dbname?sslmode=disable"),
		Port:        getEnvInt("PORT", 8080),
		SecretKey:   getEnv("SECRET_KEY", "change-me-in-production"),
		Debug:       getEnvBool("DEBUG", false),
		EnableConsole: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Gogo Framework",
			"console": "/console",
			"api":     "/api/v1",
		})
	})

	log.Println("Starting server on :8080")
	log.Println("Console: http://localhost:8080/console")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}
`

	if err := os.WriteFile(filepath.Join(projectPath, "cmd/web/main.go"), []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}

	// Create README
	readmeContent := fmt.Sprintf(`# %s

A Gogo Framework project.

## Getting Started

1. Install dependencies:
   `+"```bash\n   go mod tidy\n   ```\n\n"+`2. Define your Ent schemas in `+"`internal/models/ent/schema/`\n\n"+`3. Generate Ent code:
   `+"```bash\n   go run -mod=mod entgo.io/ent/cmd/ent generate ./internal/models/ent/schema\n   ```\n\n"+`4. Run migrations:
   `+"```bash\n   gogo migrate\n   ```\n\n"+`5. Start the server:
   `+"```bash\n   go run main.go\n   ```\n\n"+`## Console Interface

- Console: http://localhost:8080/console
- API: http://localhost:8080/api/v1
`, strings.Title(projectName))

	if err := os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	fmt.Printf("✅ Created project '%s' at %s\n", projectName, projectPath)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  # Define your Ent schemas in internal/models/ent/schema/")
	fmt.Println("  # Generate Ent code: go run -mod=mod entgo.io/ent/cmd/ent generate ./internal/models/ent/schema")
	fmt.Println("  # Start server: go run cmd/web/main.go")

	return nil
}

func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}
	// Check if it's a valid Go package name
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

