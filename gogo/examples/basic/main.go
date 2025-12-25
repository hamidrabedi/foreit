package main

import (
	"log"

	"github.com/gogo/pkg/gogo"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app, err := gogo.New(&gogo.AppConfig{
		DatabaseURL: "postgres://user:pass@localhost/dbname?sslmode=disable",
		Port:        8080,
		SecretKey:   "change-me-in-production",
		Debug:       true,
		EnableConsole: true,
	})
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Gogo framework is running",
			"console": "/console",
			"api":     "/api/v1",
		})
	})

	log.Println("Starting server on :8080")
	log.Println("Console: http://localhost:8080/console")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
