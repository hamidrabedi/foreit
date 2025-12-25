package admin

import (
	"testing"
)

// TestRegistry_Register tests model registration
func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	
	// Test registering a model
	type TestModel struct {
		ID   int
		Name string
	}
	
	err := registry.Register(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to register model: %v", err)
	}
	
	// Test retrieving registered model
	meta, err := registry.GetModel("TestModel")
	if err != nil {
		t.Fatalf("Failed to get model: %v", err)
	}
	
	if meta.Name != "TestModel" {
		t.Errorf("Expected model name 'TestModel', got '%s'", meta.Name)
	}
	
	// Test registering duplicate model
	err = registry.Register(&TestModel{})
	if err == nil {
		t.Error("Expected error when registering duplicate model")
	}
}

// TestRegistry_GetAllModels tests getting all models
func TestRegistry_GetAllModels(t *testing.T) {
	registry := NewRegistry()
	
	type Model1 struct {
		ID int
	}
	
	type Model2 struct {
		ID int
	}
	
	registry.Register(&Model1{})
	registry.Register(&Model2{})
	
	models := registry.GetAllModels()
	if len(models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(models))
	}
	
	if _, ok := models["Model1"]; !ok {
		t.Error("Model1 not found")
	}
	
	if _, ok := models["Model2"]; !ok {
		t.Error("Model2 not found")
	}
}

// TestRegistry_ThreadSafety tests thread safety
func TestRegistry_ThreadSafety(t *testing.T) {
	registry := NewRegistry()
	
	type TestModel struct {
		ID int
	}
	
	// Register model
	registry.Register(&TestModel{})
	
	// Test concurrent access (basic test)
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = registry.GetModel("TestModel")
			done <- true
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestModelMeta_Permissions tests default permissions
func TestModelMeta_Permissions(t *testing.T) {
	registry := NewRegistry()
	
	type TestModel struct {
		ID int
	}
	
	err := registry.Register(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to register model: %v", err)
	}
	
	meta, err := registry.GetModel("TestModel")
	if err != nil {
		t.Fatalf("Failed to get model: %v", err)
	}
	
	// Check default permissions
	if !meta.Permissions.CanList {
		t.Error("Expected CanList to be true by default")
	}
	if !meta.Permissions.CanView {
		t.Error("Expected CanView to be true by default")
	}
	if !meta.Permissions.CanCreate {
		t.Error("Expected CanCreate to be true by default")
	}
	if !meta.Permissions.CanUpdate {
		t.Error("Expected CanUpdate to be true by default")
	}
	if !meta.Permissions.CanDelete {
		t.Error("Expected CanDelete to be true by default")
	}
}

