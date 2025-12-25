package main

import (
	"fmt"
	"os"

	"github.com/gogo/cmd/go-admin-cli/commands"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gogo",
		Short: "Gogo CLI",
		Long:  "A Django-like admin engine for Go with Ent, Fiber, and Refine support",
		Version: version,
	}

	// Add commands
	rootCmd.AddCommand(commands.NewStartProjectCommand())
	rootCmd.AddCommand(commands.NewStartAppCommand())
	rootCmd.AddCommand(commands.NewMigrateCommand())
	rootCmd.AddCommand(commands.NewCreateSuperuserCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

