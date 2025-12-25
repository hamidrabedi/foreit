package testing

import (
	"context"
	"os"

	"github.com/gogo/pkg/gogo"
	"github.com/gogo/pkg/orm"
	"github.com/gofiber/fiber/v2"
)

// TestApp creates a test application
func TestApp() (*gogo.App, *orm.Client, error) {
	// Use test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://test:test@localhost/test_db?sslmode=disable"
	}
	
	app, err := gogo.New(&gogo.AppConfig{
		DatabaseURL: dbURL,
		Port:        0, // Random port
		SecretKey:   "test-secret-key",
		Debug:       true,
		
		// Disable some modules for testing
		EnableConsole: false,
		EnableWorkers: false,
		EnableCache:   true,
		EnableSessions: true,
		EnableI18n:    false,
		EnableStatic:  false,
	})
	if err != nil {
		return nil, nil, err
	}
	
	client := app.Client()
	
	return app, client, nil
}

// TestClient creates a test database client
func TestClient() (*orm.Client, error) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://test:test@localhost/test_db?sslmode=disable"
	}
	
	return orm.NewClient("postgres", dbURL)
}

// CleanupTestDB cleans up test database
func CleanupTestDB(client *orm.Client) error {
	// This would drop test tables or reset database
	// Implementation depends on test strategy
	return nil
}

// TestRequest creates a test HTTP request
func TestRequest(app *gogo.App, method, path string) *fiber.Request {
	req := &fiber.Request{}
	req.Header.SetMethod(method)
	req.SetRequestURI(path)
	return req
}

// TestContext creates a test context
func TestContext() context.Context {
	return context.Background()
}

