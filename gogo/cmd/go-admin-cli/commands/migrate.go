package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewMigrateCommand creates a new migrate command
func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  "Runs Ent migrations to update the database schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createMigration(args[0])
		},
	})

	return cmd
}

func runMigrations() error {
	// Check if we're in a project directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check for go.mod
	if _, err := os.Stat(filepath.Join(wd, "go.mod")); os.IsNotExist(err) {
		return fmt.Errorf("not in a Go project directory (go.mod not found)")
	}

	// Check for Ent schemas
	schemaPath := filepath.Join(wd, "internal/models/ent/schema")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		return fmt.Errorf("Ent schemas not found at %s", schemaPath)
	}

	fmt.Println("Running migrations...")
	fmt.Println("Note: This is a placeholder. In production, this would:")
	fmt.Println("  1. Generate Ent code: go run -mod=mod entgo.io/ent/cmd/ent generate ./internal/models/ent/schema")
	fmt.Println("  2. Create migration files")
	fmt.Println("  3. Apply migrations to database")

	// TODO: Implement actual migration logic
	// - Generate Ent code
	// - Create migration files
	// - Apply to database

	return nil
}

func createMigration(name string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	migrationsPath := filepath.Join(wd, "migrations")
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create migration file
	migrationFile := filepath.Join(migrationsPath, fmt.Sprintf("%s.sql", name))
	migrationContent := fmt.Sprintf(`-- Migration: %s
-- Created: %s

-- Add your migration SQL here

`, name, "now()")

	if err := os.WriteFile(migrationFile, []byte(migrationContent), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	fmt.Printf("✅ Created migration file: %s\n", migrationFile)
	return nil
}

