package cli

import (
	"github.com/forgego/forge/pkg/cli/cmd"
	"github.com/forgego/forge/pkg/cli/commands"
	"github.com/spf13/cobra"
)

// BuildRootCommand builds and returns the root cobra command with all registered commands
func BuildRootCommand() *cobra.Command {
	// Register all commands
	commands.RegisterAllCommands()

	// Build and return the root command
	registry := cmd.GetRegistry()
	return registry.BuildRootCommand()
}

