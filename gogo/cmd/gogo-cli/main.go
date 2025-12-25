package main

import (
	"github.com/gogo/cmd/gogo-cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gogo",
		Short: "Gogo framework CLI",
		Long:  "Command-line interface for Gogo framework",
	}
	
	// Register commands
	commands.RegisterStartProjectCommand(rootCmd)
	commands.RegisterStartAppCommand(rootCmd)
	commands.RegisterGenerateCommand(rootCmd)
	commands.RegisterMigrateCommand(rootCmd)
	
	rootCmd.Execute()
}

