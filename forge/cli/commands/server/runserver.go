package server

import (
	"fmt"
	"net/http"

	"github.com/forgego/forge/cli/core"
	"github.com/forgego/forge/config"
	httplib "github.com/forgego/forge/server"
	"github.com/forgego/forge/log"
	"github.com/spf13/cobra"
)

// RunServerCommand creates the runserver command
type RunServerCommand struct{}

// NewRunServerCommand creates a new instance of RunServerCommand
func NewRunServerCommand() *RunServerCommand {
	return &RunServerCommand{}
}

// Definition returns the cobra command definition
func (c *RunServerCommand) Definition() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runserver",
		Short: "Run development server",
		Long:  "Start the development web server",
	}
	cmd.Flags().String("port", "", "Port to run server on (overrides config)")
	return cmd
}

// Execute runs the command logic
func (c *RunServerCommand) Execute(ctx *core.Context, args []string) error {
	settings := config.LoadSettings(ctx.Config)

	// Override port if provided via flag
	if port, err := ctx.Cmd.Flags().GetString("port"); err == nil && port != "" {
		settings.Server.Port = port
	}

	// Create logger if not already set
	if ctx.Logger == nil {
		logger, err := log.NewLogger(settings.App.Debug)
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}
		ctx.WithLogger(logger)
	}

	// Create server with all middleware
	server, err := httplib.NewServer(ctx.Config, settings, ctx.Logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Register routes
	server.RegisterRoutes(func(router *httplib.Router) {
		registerRoutes(router, ctx.Logger)
	})

	// Start server
	return server.Start()
}

// registerRoutes registers all framework routes
func registerRoutes(router *httplib.Router, logger *log.Logger) {
	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// nolint:errcheck // HTTP response write errors can't be handled meaningfully
		_, _ = w.Write([]byte("OK"))
	})

	// TODO: Register admin routes
	// Example:
	// adminRegistry := admin.GetGlobalRegistry()
	// adminRouter := adminhttp.NewRouter(adminRegistry)
	// adminRouter.RegisterRoutes(router, settings.Admin.Path)

	// TODO: Register API routes
}
