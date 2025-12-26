package server

import (
	"fmt"
	"net/http"

	"github.com/forgego/forge/pkg/admin"
	"github.com/forgego/forge/pkg/cli/cmd"
	"github.com/forgego/forge/pkg/config"
	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
func (c *RunServerCommand) Execute(ctx *cmd.Context, args []string) error {
	settings := config.LoadSettings(ctx.Config)

	// Override port if provided via flag
	if port, err := ctx.Cmd.Flags().GetString("port"); err == nil && port != "" {
		settings.Server.Port = port
	}

	// Create logger if not already set
	if ctx.Logger == nil {
		logger, err := logging.NewLogger(ctx.Config.IsDevelopment())
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
func registerRoutes(router *httplib.Router, logger *logging.Logger) {
	// Health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// nolint:errcheck // HTTP response write errors can't be handled meaningfully
		_, _ = w.Write([]byte("OK"))
	})

	// Register admin routes
	cfg := config.NewConfig()
	settings := config.LoadSettings(cfg)
	if settings.Admin.Enabled {
		admin.RegisterAdminRoutes(router, settings.Admin.Path, nil, nil)
		logger.Info("Admin routes registered", zap.String("path", settings.Admin.Path))
	}

	// TODO: Register API routes
}

