package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"enterprise-demo/app/enterprise"
	"github.com/forgego/forge/config"
	"github.com/forgego/forge/db"
	httplib "github.com/forgego/forge/server"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()
	cfg.SetConfigFile("config/config.yaml")
	if err := cfg.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config file: %v (using defaults)", err)
	}

	// Setup database
	database, err := db.NewDBFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	orgRepo, err := enterprise.NewOrganizationRepository(database)
	if err != nil {
		log.Printf("Warning: Failed to create organization repository: %v", err)
	}

	// Additional repositories can be initialized here as needed
	_, err = enterprise.NewEmployeeRepository(database)
	if err != nil {
		log.Printf("Warning: Failed to create employee repository: %v", err)
	}

	_, err = enterprise.NewProjectRepository(database)
	if err != nil {
		log.Printf("Warning: Failed to create project repository: %v", err)
	}

	// Initialize services
	var orgService *enterprise.OrganizationService
	if orgRepo != nil {
		orgService = enterprise.NewOrganizationService(orgRepo)
	}

	// Setup router
	router := httplib.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Example endpoint using repositories
	router.HandleFunc("/api/organizations", func(w http.ResponseWriter, r *http.Request) {
		if orgRepo == nil {
			http.Error(w, "Repository not initialized", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		orgs, err := orgRepo.GetActiveOrganizations(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"count": %d}`, len(orgs))))
	})

	// Example endpoint using services
	router.HandleFunc("/api/organizations/stats", func(w http.ResponseWriter, r *http.Request) {
		if orgService == nil {
			http.Error(w, "Service not initialized", http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		stats, err := orgService.GetOrganizationStats(ctx, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"employees": %d, "projects": %d}`,
			stats.EmployeeCount, stats.ProjectCount)))
	})

	// Run query examples endpoint (for demonstration)
	router.HandleFunc("/api/demo/queries", func(w http.ResponseWriter, r *http.Request) {
		enterprise.RunAllQueryExamples(database)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Query examples executed - check logs"))
	})

	// Start server
	port := cfg.GetString("server.port", "8080")
	log.Printf("Starting server on :%s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET /health")
	log.Printf("  GET /api/organizations")
	log.Printf("  GET /api/organizations/stats")
	log.Printf("  GET /api/demo/queries")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
