package migrations

import (
	"fmt"
	"path/filepath"

	"github.com/forgego/forge/pkg/cli/cmd"
	"github.com/forgego/forge/pkg/migrations"
	"github.com/spf13/cobra"
)

// LintCommand creates the migration lint command
type LintCommand struct{}

// NewLintCommand creates a new instance of LintCommand
func NewLintCommand() *LintCommand {
	return &LintCommand{}
}

// Definition returns the cobra command definition
func (c *LintCommand) Definition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint migration files",
		Long:  "Check migration files for common issues and best practices",
	}
	cmd.Flags().String("path", "./migrations", "Path to migrations directory")
	return cmd
}

// Execute runs the command logic
func (c *LintCommand) Execute(ctx *cmd.Context, args []string) error {
	// Get migrations path
	migrationsPath, err := ctx.Cmd.Flags().GetString("path")
	if err != nil {
		return fmt.Errorf("failed to get migrations path flag: %w", err)
	}
	if migrationsPath == "" {
		migrationsPath = "./migrations"
	}

	// Create linter
	linter := migrations.NewLinter()

	// Lint all migrations
	results, err := linter.LintMigrationsDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to lint migrations: %w", err)
	}

	// Group results by level
	errors := []migrations.LintResult{}
	warnings := []migrations.LintResult{}
	infos := []migrations.LintResult{}

	for _, result := range results {
		switch result.Level {
		case "error":
			errors = append(errors, result)
		case "warning":
			warnings = append(warnings, result)
		case "info":
			infos = append(infos, result)
		}
	}

	// Print results
	if len(errors) > 0 {
		fmt.Println("❌ Errors:")
		for _, result := range errors {
			fmt.Printf("  %s: %s\n", filepath.Base(result.File), result.Message)
		}
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Println("⚠️  Warnings:")
		for _, result := range warnings {
			fmt.Printf("  %s: %s\n", filepath.Base(result.File), result.Message)
		}
		fmt.Println()
	}

	if len(infos) > 0 {
		fmt.Println("ℹ️  Info:")
		for _, result := range infos {
			fmt.Printf("  %s: %s\n", filepath.Base(result.File), result.Message)
		}
		fmt.Println()
	}

	if len(results) == 0 {
		fmt.Println("✓ No issues found")
	} else {
		fmt.Printf("Summary: %d errors, %d warnings, %d info messages\n",
			len(errors), len(warnings), len(infos))

		if len(errors) > 0 {
			return fmt.Errorf("linting found %d error(s)", len(errors))
		}
	}

	return nil
}

