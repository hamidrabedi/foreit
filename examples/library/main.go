package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/forgego/forge/config"
	"github.com/forgego/forge/db"
	"github.com/forgego/forge/log"
	"github.com/forgego/forge/server"
	"github.com/forgego/forge/admin"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()
	settings := config.LoadSettings(cfg)

	// Create logger
	logger, err := log.NewLogger(settings.App.Debug)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Connect to database (optional - server can run without it)
	var database *db.DB
	database, err = db.NewDBFromConfig(cfg)
	if err != nil {
		logger.Warn("Failed to connect to database", log.Error(err))
		logger.Info("Server will start without database connection")
	} else {
		defer database.Close()
		logger.Info("Database connection established")
	}

	// Create server
	srv, err := server.NewServer(cfg, settings, logger)
	if err != nil {
		log.Fatal(err)
	}

	// Register routes
	srv.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Library Management System!")
	})

	// Register admin routes
	if settings.Admin.Enabled {
		// Get session manager from server (if available)
		// For now, pass nil - authentication will be added later
		var sessionManager interface{} = nil
		var dbInterface interface{} = nil
		if database != nil {
			dbInterface = database
		}
		// Note: Admin routes registration would go here
		// admin.RegisterAdminRoutes(srv.Router, settings.Admin.Path, sessionManager, dbInterface)
		_ = sessionManager
		_ = dbInterface
	}

	// Run query examples if database is available
	if database != nil {
		fmt.Println("\n=== Running Query Examples ===")
		RunQueryExamples(database)
	}

	// Start server
	fmt.Printf("Starting server on %s:%s\n", settings.Server.Host, settings.Server.Port)
	fmt.Printf("Admin interface available at http://%s:%s%s\n", settings.Server.Host, settings.Server.Port, settings.Admin.Path)
	
	addr := fmt.Sprintf("%s:%s", settings.Server.Host, settings.Server.Port)
	if err := http.ListenAndServe(addr, srv.Router); err != nil {
		log.Fatal(err)
	}
}
