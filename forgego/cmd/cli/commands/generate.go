package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func RegisterGenerateCommand(rootCmd *cobra.Command) {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate static typed models from templates",
		Long: `Generate static typed models and structs from model template files.

The command looks for model template files (e.g., user.go.tmpl) in the models/ directory
and generates the corresponding Go model files with proper struct definitions.

Example:
  forge generate
  forge generate --input models/
  forge generate --output models/generated/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			inputDir := cmd.Flag("input").Value.String()
			outputDir := cmd.Flag("output").Value.String()
			
			if inputDir == "" {
				inputDir = "models"
			}
			if outputDir == "" {
				outputDir = inputDir
			}
			
			return generateModels(inputDir, outputDir)
		},
	}

	generateCmd.Flags().String("input", "models", "Input directory containing model templates")
	generateCmd.Flags().String("output", "", "Output directory for generated models (default: same as input)")

	rootCmd.AddCommand(generateCmd)
}

func generateModels(inputDir, outputDir string) error {
	// Check if input directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory does not exist: %s", inputDir)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Find all .go.tmpl files
	templateFiles := []string{}
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go.tmpl") {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(templateFiles) == 0 {
		fmt.Printf("No template files (.go.tmpl) found in %s\n", inputDir)
		fmt.Println("\nCreate model template files like this:")
		fmt.Println("  models/user.go.tmpl")
		fmt.Println("\nExample template:")
		fmt.Println(`  package models`)
		fmt.Println(``)
		fmt.Println(`  type {{.ModelName}} struct {`)
		fmt.Println(`      models.Schema`)
		fmt.Println(`      ID        int64`)
		fmt.Println(`      Name      string`)
		fmt.Println(`      Email     string`)
		fmt.Println(`      CreatedAt time.Time`)
		fmt.Println(`  }`)
		return nil
	}

	// Process each template file
	for _, templateFile := range templateFiles {
		relPath, err := filepath.Rel(inputDir, templateFile)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Remove .tmpl extension
		outputFile := strings.TrimSuffix(relPath, ".tmpl")
		outputPath := filepath.Join(outputDir, outputFile)

		// Read template
		content, err := os.ReadFile(templateFile)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", templateFile, err)
		}

		// For now, just copy the template content (in the future, this could process templates)
		// Users should fill in the templates with actual model definitions
		if err := os.WriteFile(outputPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", outputPath, err)
		}

		fmt.Printf("✓ Generated: %s\n", outputPath)
	}

	fmt.Printf("\n✓ Generated %d model(s) from templates\n", len(templateFiles))
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the generated model files")
	fmt.Println("  2. Add Fields() methods to define your schema")
	fmt.Println("  3. Register models in your application")

	return nil
}
