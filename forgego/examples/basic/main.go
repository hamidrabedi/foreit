package main

import (
	"log"

	"github.com/forgego/forge/pkg/app"
	"github.com/forgego/forge/pkg/settings"
	"github.com/gofiber/fiber/v2"
)

// MyAppConfig represents custom application settings
type MyAppConfig struct {
	APIKey    string `mapstructure:"api_key" validate:"required"`
	MaxUsers  int    `mapstructure:"max_users" default:"100"`
	EnableFeature bool `mapstructure:"enable_feature" default:"false"`
}

// Settings embeds AppSettings and adds custom configs
// This is how users extend the framework settings
type Settings struct {
	app.AppSettings  // Embed all framework settings
	
	// Add your own custom settings
	MyApp MyAppConfig `mapstructure:"myapp"`
}

func main() {
	// Load settings from config file or environment variables
	appSettings, err := settings.Load[Settings]()
	if err != nil {
		log.Fatalf("Failed to load settings: %v", err)
	}

	// Access framework settings (type-safe)
	log.Printf("Database URL: %s", appSettings.Database.URL)
	log.Printf("Server Port: %d", appSettings.Server.Port)
	log.Printf("Admin Path: %s", appSettings.Admin.Path)
	
	// Access your custom settings
	log.Printf("API Key: %s", appSettings.MyApp.APIKey)
	log.Printf("Max Users: %d", appSettings.MyApp.MaxUsers)
	log.Printf("Feature Enabled: %v", appSettings.MyApp.EnableFeature)

	// Create app - framework features are automatically initialized based on settings
	// If initialization fails, an error is returned
	application, err := app.New(&appSettings.AppSettings)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Add your own routes
	application.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Forge!",
			"api_key": appSettings.MyApp.APIKey,
		})
	})

	// Start the server
	if err := application.Listen(""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
