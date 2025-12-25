package main

import (
	"github.com/forgego/forge/cmd/cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "forge",
		Short: "Forge framework CLI",
		Long:  "Command-line interface for Forge framework",
	}
	
	// Register commands
	commands.RegisterStartProjectCommand(rootCmd)
	commands.RegisterGenerateCommand(rootCmd)
	
	rootCmd.Execute()
}

