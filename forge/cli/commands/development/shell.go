package development

import (
	"fmt"

	"github.com/forgego/forge/cli/core"
	"github.com/spf13/cobra"
)

// ShellCommand creates the shell command
type ShellCommand struct{}

// NewShellCommand creates a new instance of ShellCommand
func NewShellCommand() *ShellCommand {
	return &ShellCommand{}
}

// Definition returns the cobra command definition
func (c *ShellCommand) Definition() *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Interactive shell",
		Long:  "Start an interactive shell with framework context",
	}
}

// Execute runs the command logic
func (c *ShellCommand) Execute(ctx *core.Context, args []string) error {
	// TODO: Implement shell
	fmt.Println("Starting interactive shell...")
	return nil
}
