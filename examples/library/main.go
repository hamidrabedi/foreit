package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/forgego/forge"
)

func main() {
	// Load configuration
	cfg := forge.NewConfig()
	settings := forge.LoadSettings(cfg)

	// Create logger
	logger, err := forge.NewLogger(cfg.IsDevelopment())
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Connect to database (optional - server can run without it)
	var database *forge.DB
	database, err = forge.NewDBFromConfig(cfg)
	if err != nil {
		logger.Warn("Failed to connect to database", forge.Error(err))
		logger.Info("Server will start without database connection")
	} else {
		defer database.Close()
		logger.Info("Database connection established")
	}

	// Create server
	server, err := forge.NewServer(cfg, settings, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Register routes
	server.RegisterRoutes(func(router *forge.Router) {
		// Health check
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Welcome to Library Management System!")
		})

		// Register admin routes
		if settings.Admin.Enabled {
			// Get session manager from server (if available)
			// For now, pass nil - authentication will be added later
			var sessionManager interface{} = nil
			var db interface{} = nil
			if database != nil {
				db = database
			}
			forge.RegisterAdminRoutes(router, settings.Admin.Path, sessionManager, db)
		}
	})

	// Start server
	fmt.Printf("Starting server on %s:%s\n", settings.Server.Host, settings.Server.Port)
	fmt.Printf("Admin interface available at http://%s:%s%s\n", settings.Server.Host, settings.Server.Port, settings.Admin.Path)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

