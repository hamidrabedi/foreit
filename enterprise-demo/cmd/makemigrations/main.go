package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/forgego/forge/migrate"
)

func main() {
	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Find directories
	modelsDir := filepath.Join(wd, "app", "enterprise")
	migrationsDir := filepath.Join(wd, "migrations")

	// Check if models directory exists
	if _, err := os.Stat(modelsDir); os.IsNotExist(err) {
		log.Fatalf("Models directory does not exist: %s", modelsDir)
	}

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		log.Fatalf("Failed to create migrations directory: %v", err)
	}

	fmt.Printf("Generating migrations from %s to %s...\n", modelsDir, migrationsDir)

	// Create migration generator
	gen, err := migrate.NewGenerator(modelsDir, migrationsDir)
	if err != nil {
		log.Fatalf("Failed to create migration generator: %v", err)
	}

	// Generate migrations
	if err := gen.GenerateMigrations("initial"); err != nil {
		log.Fatalf("Migration generation failed: %v", err)
	}

	fmt.Println("✓ Migration generation completed successfully!")
	fmt.Printf("  Check migrations in: %s\n", migrationsDir)
}
