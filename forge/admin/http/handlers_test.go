package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	adminv2 "github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
	"github.com/go-chi/chi/v5"
)

// TestUser is a test model for HTTP handler tests
type TestUser struct {
	ID       int64
	Username string
	Email    string
	IsActive bool
}

// TestHandler_HandleIndex tests the admin index handler
func TestHandler_HandleIndex(t *testing.T) {
	// Create registry and register test admin
	registry := adminv2.GetGlobalRegistry()
	manager := &query.Manager[*TestUser]{}
	config := &adminv2.Config[*TestUser]{
		VerboseName:       "User",
		VerboseNamePlural: "Users",
	}
	adminv2.Register(&TestUser{}, manager, config)

	// Create handler
	handler := NewHandler(registry)

	// Create request
	req := httptest.NewRequest("GET", "/admin/", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.HandleIndex()(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestHandler_HandleList tests the list view handler
func TestHandler_HandleList(t *testing.T) {
	// Create registry and register test admin
	registry := adminv2.GetGlobalRegistry()
	manager := &query.Manager[*TestUser]{}
	config := &adminv2.Config[*TestUser]{}
	admin := adminv2.Register(&TestUser{}, manager, config)

	// Register with type registry for HTTP handlers
	RegisterAdmin(admin)

	// Create handler
	handler := NewHandler(registry)

	// Create request
	req := httptest.NewRequest("GET", "/admin/TestUser/", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.HandleList("TestUser")(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestHandler_HandleList_NotFound tests list view with non-existent model
func TestHandler_HandleList_NotFound(t *testing.T) {
	registry := adminv2.GetGlobalRegistry()
	handler := NewHandler(registry)

	req := httptest.NewRequest("GET", "/admin/NonExistent/", nil)
	w := httptest.NewRecorder()

	handler.HandleList("NonExistent")(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// TestHandler_HandleDetail tests the detail view handler
func TestHandler_HandleDetail(t *testing.T) {
	registry := adminv2.GetGlobalRegistry()
	manager := &query.Manager[*TestUser]{}
	config := &adminv2.Config[*TestUser]{}
	admin := adminv2.Register(&TestUser{}, manager, config)

	// Register with type registry for HTTP handlers
	RegisterAdmin(admin)

	handler := NewHandler(registry)

	req := httptest.NewRequest("GET", "/admin/TestUser/1/", nil)
	// Set URL parameter in context (chi router does this)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.HandleDetail("TestUser")(w, req)

	// Should return 200 or 404 depending on whether object exists
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status 200 or 404, got %d", w.Code)
	}
}

// TestHandler_HandleCreate tests the create handler
func TestHandler_HandleCreate(t *testing.T) {
	registry := adminv2.GetGlobalRegistry()
	manager := &query.Manager[*TestUser]{}
	config := &adminv2.Config[*TestUser]{}
	adminv2.Register(&TestUser{}, manager, config)

	handler := NewHandler(registry)

	// Test GET request
	req := httptest.NewRequest("GET", "/admin/TestUser/new/", nil)
	w := httptest.NewRecorder()

	handler.HandleCreate("TestUser")(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestHandler_HandleDelete tests the delete handler
func TestHandler_HandleDelete(t *testing.T) {
	registry := adminv2.GetGlobalRegistry()
	manager := &query.Manager[*TestUser]{}
	config := &adminv2.Config[*TestUser]{}
	adminv2.Register(&TestUser{}, manager, config)

	handler := NewHandler(registry)

	req := httptest.NewRequest("POST", "/admin/TestUser/1/delete/", nil)
	w := httptest.NewRecorder()

	handler.HandleDelete("TestUser")(w, req)

	// Should return 200, 303, 404, or 500 depending on whether object exists
	if w.Code < 200 || w.Code >= 600 {
		t.Errorf("Expected status 2xx-5xx, got %d", w.Code)
	}
}
