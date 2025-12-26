package migrations

import (
	"github.com/forgego/forge/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

// MigrateGroup represents the migrate command group
type MigrateGroup struct {
	commands []cmd.Command
}

// NewMigrateGroup creates a new migrate command group
func NewMigrateGroup() *MigrateGroup {
	return &MigrateGroup{
		commands: []cmd.Command{
			NewUpCommand(),
			NewRollbackCommand(),
			NewStatusCommand(),
			NewShowCommand(),
			NewLintCommand(),
			NewMakeMigrationsCommand(),
			NewForceCommand(),
			NewSquashCommand(),
			NewFakeCommand(),
		},
	}
}

// Definition returns the cobra command definition for the group
func (g *MigrateGroup) Definition() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migration management commands",
		Long:  "Commands for managing database migrations",
	}
}

// Execute runs the group command (usually just shows help)
func (g *MigrateGroup) Execute(ctx *cmd.Context, args []string) error {
	// Group commands typically just show help
	return ctx.Cmd.Help()
}

// Commands returns subcommands in this group
func (g *MigrateGroup) Commands() []cmd.Command {
	return g.commands
}

