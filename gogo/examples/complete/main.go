package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gogo/pkg/auth"
	"github.com/gogo/pkg/gogo"
	"github.com/gogo/pkg/routing"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create Gogo application
	app, err := gogo.New(&gogo.AppConfig{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost/dbname"),
		Port:        getEnvInt("PORT", 8080),
		SecretKey:   getEnv("SECRET_KEY", "change-me-in-production"),
		Debug:       getEnvBool("DEBUG", false),
		
		// Enable all modules
		EnableConsole: true,
		EnableWorkers: true,
		EnableCache:   true,
		EnableSessions: true,
		EnableI18n:    true,
		EnableStatic:  true,
		
		// Configure paths
		ConsolePath:  "/admin",
		APIPath:      "/api/v1",
		StaticPath:   "/static",
		StaticRoot:   "./public",
		LocalesPath:  "./locales",
		DefaultLocale: "en",
		
		// Redis configuration (for workers and cache)
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		WorkersConcurrency: getEnvInt("WORKERS_CONCURRENCY", 10),
	})
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	
	// Setup authentication
	jwtAuth := auth.NewJWT(getEnv("SECRET_KEY", "change-me"))
	app.SetAuth(jwtAuth, false) // Optional auth
	
	// Register policies
	// auth.Register[models.Post](&PostPolicy{})
	
	// Register resources
	// app.RegisterResource("users", &UserResource{})
	
	// Register consoles
	// app.RegisterConsole[models.User](&UserConsole{})
	
	// Add custom routes
	router := app.Router()
	router.Get("/", homeHandler, routing.Name("home"))
	router.Get("/health", healthHandler, routing.Name("health"))
	
	// Graceful shutdown
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", getEnvInt("PORT", 8080))); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	
	log.Println("Server stopped")
}

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to Gogo Framework",
		"version": "0.1.0",
		"modules": []string{
			"orm", "settings", "endpoints", "routing",
			"pipeline", "auth", "console", "workers",
			"cache", "sessions", "i18n", "static", "utils",
		},
	})
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

