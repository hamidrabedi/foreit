package generation

import (
	"fmt"

	"github.com/forgego/forge/pkg/cli/cmd"
	"github.com/forgego/forge/pkg/generator"
	"github.com/spf13/cobra"
)

// GenerateCommand creates the generate command
type GenerateCommand struct{}

// NewGenerateCommand creates a new instance of GenerateCommand
func NewGenerateCommand() *GenerateCommand {
	return &GenerateCommand{}
}

// Definition returns the cobra command definition
func (c *GenerateCommand) Definition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code from schema definitions",
		Long:  "Parse schema definitions and generate type-safe models, managers, and querysets",
	}
	cmd.Flags().String("models", "./models", "Directory containing schema definitions")
	cmd.Flags().String("output", "./models", "Output directory for generated code")
	return cmd
}

// Execute runs the command logic
func (c *GenerateCommand) Execute(ctx *cmd.Context, args []string) error {
	modelsDir, err := ctx.Cmd.Flags().GetString("models")
	if err != nil {
		return fmt.Errorf("failed to get models flag: %w", err)
	}
	if modelsDir == "" {
		modelsDir = "./models"
	}

	outputDir, err := ctx.Cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}
	if outputDir == "" {
		outputDir = "./models"
	}

	// Create generator and run
	gen := generator.NewGenerator(modelsDir, outputDir)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	fmt.Printf("✓ Generated code from %s to %s\n", modelsDir, outputDir)
	return nil
}

