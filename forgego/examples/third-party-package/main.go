package main

import (
	"log"

	"github.com/forgego/forge/pkg/app"
	"github.com/forgego/forge/pkg/settings"

	// Import your third-party package
	"./mypackage"
)

// Settings embeds AppSettings and includes third-party package configs
// This shows how users can integrate multiple packages
type Settings struct {
	app.AppSettings  // Embed all framework settings
	
	// Add third-party package configs
	MyPackage mypackage.MyPackageConfig `mapstructure:"mypackage"`
}

func main() {
	// Load settings from config file or environment variables
	appSettings, err := settings.Load[Settings]()
	if err != nil {
		log.Fatalf("Failed to load settings: %v", err)
	}

	// Initialize third-party package with its config
	myPackage := mypackage.New(&appSettings.MyPackage)
	if err := myPackage.Start(); err != nil {
		log.Fatalf("Failed to start MyPackage: %v", err)
	}

	// Create app with framework settings
	application, err := app.New(&appSettings.AppSettings)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	application.Fiber().Use(mypackage.Middleware())

	// Start the server
	if err := application.Listen(""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

