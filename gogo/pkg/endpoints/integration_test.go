package endpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestResourceIntegration tests the full resource lifecycle
func TestResourceIntegration(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", nil)
		return c.Next()
	})
	
	router := NewRouter(app, "/api/v1")
	
	// Create test repository
	repo := NewMockRepository[*TestUser, interface{}]()
	resource := NewResource[*TestUser, interface{}](repo)
	
	// Register resource
	router.RegisterResource("users", resource)
	
	// Test Create
	t.Run("Create", func(t *testing.T) {
		userData := map[string]interface{}{
			"name":  "Test User",
			"email": "test@example.com",
		}
		
		body, _ := json.Marshal(userData)
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})
	
	// Test Index
	t.Run("Index", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
	
	// Test Show
	t.Run("Show", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
		
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		
		// Should work if we have data, or 404 if not
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 200 or 404, got %d", resp.StatusCode)
		}
	})
}

// TestQueryProcessing tests query parameter processing
func TestQueryProcessing(t *testing.T) {
	app := fiber.New()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		page, pageSize := ParsePagination(c)
		sortBy, sortOrder := ParseSorting(c)
		
		return c.JSON(fiber.Map{
			"page":       page,
			"page_size":  pageSize,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		})
	})
	
	req := httptest.NewRequest("GET", "/test?page=2&page_size=10&sort_by=name&sort_order=desc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestErrorHandling tests error handling
func TestErrorHandling(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		return HandleError(c, ErrNotFound)
	})
	
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Next()
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

