package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var startprojectCmd = &cobra.Command{
	Use:   "startproject [name]",
	Short: "Create a new Gogo project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		createProject(projectName)
	},
}

func createProject(name string) {
	// Create project directory
	if err := os.MkdirAll(name, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		os.Exit(1)
	}
	
	// Create main.go
	mainContent := `package main

import (
	"log"
	
	"github.com/forgego/forge/pkg/app"
)

func main() {
	appInstance, err := app.New(&app.Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Port:        getEnvInt("PORT", 8080),
		SecretKey:   getEnv("SECRET_KEY", "change-me"),
		Debug:       getEnvBool("DEBUG", false),
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Add your routes here
	
	if err := appInstance.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
`
	
	writeFile(filepath.Join(name, "main.go"), mainContent)
	
	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/forgego/forge/pkg/app latest
	github.com/gofiber/fiber/v2 latest
)
`, name)
	
	writeFile(filepath.Join(name, "go.mod"), goModContent)
	
	// Create .env.example
	envContent := `DATABASE_URL=postgres://user:pass@localhost/dbname
PORT=8080
SECRET_KEY=change-me-in-production
DEBUG=false
`
	
	writeFile(filepath.Join(name, ".env.example"), envContent)
	
	// Create README
	readmeContent := fmt.Sprintf(`# %s

A Forge framework application.

## Setup

1. Copy .env.example to .env
2. Update DATABASE_URL
3. Run: go mod tidy
4. Run: go run main.go

## Structure

- main.go - Application entry point
- models/ - Model templates (.go.tmpl files) and generated models
- resources/ - API resources
- policies/ - Authorization policies

## Usage

1. Create model templates in models/ directory (e.g., user.go.tmpl)
2. Run: forge generate
3. Fill in the generated model files with Fields() methods
4. Register models in your application
`, name)
	
	writeFile(filepath.Join(name, "README.md"), readmeContent)
	
	fmt.Printf("Project '%s' created successfully!\n", name)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}

// RegisterStartProjectCommand registers the startproject command
func RegisterStartProjectCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(startprojectCmd)
}

