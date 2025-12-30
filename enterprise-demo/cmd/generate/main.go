package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	codegen "github.com/forgego/forge/codegen"
)

func main() {
	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Find app directory
	modelsDir := filepath.Join(wd, "app", "enterprise")
	outputDir := modelsDir

	// Check if models directory exists
	if _, err := os.Stat(modelsDir); os.IsNotExist(err) {
		log.Fatalf("Models directory does not exist: %s", modelsDir)
	}

	fmt.Printf("Generating code from %s to %s...\n", modelsDir, outputDir)

	// Create generator
	gen := codegen.NewGenerator(modelsDir, outputDir)

	// Generate code
	if err := gen.Generate(); err != nil {
		log.Fatalf("Code generation failed: %v", err)
	}

	fmt.Println("✓ Code generation completed successfully!")
}
