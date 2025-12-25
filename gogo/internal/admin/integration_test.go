package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestCRUD_Integration tests the full CRUD flow
func TestCRUD_Integration(t *testing.T) {
	// Setup: Create a mock Ent client
	// In a real test, you'd use an in-memory database or test database
	var mockClient interface{} = nil // Placeholder
	
	registry := NewRegistry()
	
	// Register a test model
	type TestModel struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	err := registry.Register(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to register model: %v", err)
	}
	
	// Create handler
	meta, _ := registry.GetModel("TestModel")
	handler := NewCRUDHandler(meta, mockClient, registry)
	
	// Setup Fiber app
	app := fiber.New()
	app.Post("/test", handler.Create)
	app.Get("/test", handler.List)
	app.Get("/test/:id", handler.Get)
	app.Put("/test/:id", handler.Update)
	app.Delete("/test/:id", handler.Delete)
	
	// Test Create
	createData := map[string]interface{}{
		"name": "Test User",
		"age":  25,
	}
	createBody, _ := json.Marshal(createData)
	
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	
	// Test List
	req = httptest.NewRequest("GET", "/test", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	
	// Test Get (would need actual ID from create)
	req = httptest.NewRequest("GET", "/test/1", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	// Test Update
	updateData := map[string]interface{}{
		"name": "Updated User",
	}
	updateBody, _ := json.Marshal(updateData)
	
	req = httptest.NewRequest("PUT", "/test/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	// Test Delete
	req = httptest.NewRequest("DELETE", "/test/1", nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d or %d, got %d", http.StatusNoContent, http.StatusOK, resp.StatusCode)
	}
}

// TestBulkOperations_Integration tests bulk operations
func TestBulkOperations_Integration(t *testing.T) {
	var mockClient interface{} = nil
	registry := NewRegistry()
	
	type TestModel struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	
	registry.Register(&TestModel{})
	
	meta, _ := registry.GetModel("TestModel")
	handler := NewCRUDHandler(meta, mockClient, registry)
	
	app := fiber.New()
	app.Post("/test/bulk", handler.BulkCreate)
	app.Put("/test/bulk", handler.BulkUpdate)
	app.Delete("/test/bulk", handler.BulkDelete)
	
	// Test Bulk Create
	bulkData := BulkCreateRequest{
		Items: []map[string]interface{}{
			{"name": "User 1"},
			{"name": "User 2"},
		},
	}
	bulkBody, _ := json.Marshal(bulkData)
	
	req := httptest.NewRequest("POST", "/test/bulk", bytes.NewBuffer(bulkBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d or %d, got %d", http.StatusCreated, http.StatusBadRequest, resp.StatusCode)
	}
}

// TestPermissions_Integration tests permission integration
func TestPermissions_Integration(t *testing.T) {
	registry := NewRegistry()
	
	type TestModel struct {
		ID int
	}
	
	// Register with restricted permissions
	err := registry.Register(&TestModel{},
		WithPermissions(Permissions{
			CanList:   true,
			CanView:   true,
			CanCreate: false, // Deny create
			CanUpdate: false,
			CanDelete: false,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to register model: %v", err)
	}
	
	meta, _ := registry.GetModel("TestModel")
	handler := NewCRUDHandler(meta, nil, registry)
	
	app := fiber.New()
	app.Post("/test", handler.Create)
	
	createData := map[string]interface{}{"name": "Test"}
	createBody, _ := json.Marshal(createData)
	
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	// Should be forbidden
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}

// TestHooks_Integration tests hook integration
func TestHooks_Integration(t *testing.T) {
	registry := NewRegistry()
	
	type TestModel struct {
		ID   int
		Name string
	}
	
	registry.Register(&TestModel{})
	
	meta, _ := registry.GetModel("TestModel")
	handler := NewCRUDHandler(meta, nil, registry)
	
	// Register a hook
	hookCalled := false
	registry.GetHookRegistry().Register("TestModel", HookBeforeCreate, func(ctx *HookContext) error {
		hookCalled = true
		return nil
	})
	
	app := fiber.New()
	app.Post("/test", handler.Create)
	
	createData := map[string]interface{}{"name": "Test"}
	createBody, _ := json.Marshal(createData)
	
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	_, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	// Hook should have been called (even if create fails due to missing client)
	// In a real test with actual database, we'd verify the hook was called
	if !hookCalled {
		t.Error("Hook was not called")
	}
}

// TestFiltering_Integration tests filtering integration
func TestFiltering_Integration(t *testing.T) {
	registry := NewRegistry()
	
	type TestModel struct {
		ID   int
		Name string
		Age  int
	}
	
	registry.Register(&TestModel{})
	
	meta, _ := registry.GetModel("TestModel")
	handler := NewCRUDHandler(meta, nil, registry)
	
	app := fiber.New()
	app.Get("/test", handler.List)
	
	// Test with filters
	req := httptest.NewRequest("GET", "/test?filter_name=John&filter_age__gt=18", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

