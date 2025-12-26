package commands

import (
	"github.com/forgego/forge/pkg/cli/cmd"
	"github.com/forgego/forge/pkg/cli/commands/admin"
	"github.com/forgego/forge/pkg/cli/commands/development"
	"github.com/forgego/forge/pkg/cli/commands/generation"
	"github.com/forgego/forge/pkg/cli/commands/migrations"
	"github.com/forgego/forge/pkg/cli/commands/project"
	"github.com/forgego/forge/pkg/cli/commands/server"
)

// RegisterAllCommands registers all CLI commands with the global registry
func RegisterAllCommands() {
	registry := cmd.GetRegistry()

	// Register standalone commands
	registry.RegisterCommand("new", project.NewNewCommand())
	registry.RegisterCommand("generate", generation.NewGenerateCommand())
	registry.RegisterCommand("makemigrations", migrations.NewMakeMigrationsCommand())
	registry.RegisterCommand("createsuperuser", admin.NewCreateSuperUserCommand())
	registry.RegisterCommand("runserver", server.NewRunServerCommand())
	registry.RegisterCommand("shell", development.NewShellCommand())
	registry.RegisterCommand("test", development.NewTestCommand())
	registry.RegisterCommand("auth", project.NewAuthCommand())

	// Register command groups
	registry.RegisterGroup("migrate", migrations.NewMigrateGroup())
	registry.RegisterGroup("add", project.NewAddGroup())
}

