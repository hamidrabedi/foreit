package admin

import (
	"context"
	"testing"
	"time"

	query "github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// TestUser is a test model
type TestUser struct {
	ID        int64
	Username  string
	Email     string
	IsActive  bool
	CreatedAt time.Time
}

// TestAdmin_Register tests admin registration
func TestAdmin_Register(t *testing.T) {
	// Create a test manager
	manager := &query.Manager[*TestUser]{}

	// Create admin config
	config := &Config[*TestUser]{
		VerboseName:       "User",
		VerboseNamePlural: "Users",
		ListPerPage:       25,
	}

	// Register admin
	admin := Register(&TestUser{}, manager, config)

	// Verify registration
	if admin == nil {
		t.Fatal("Admin should not be nil")
	}

	if admin.ModelName() != "TestUser" {
		t.Errorf("Expected ModelName 'TestUser', got '%s'", admin.ModelName())
	}

	if admin.Manager() != manager {
		t.Error("Manager should match")
	}

	if admin.Config() != config {
		t.Error("Config should match")
	}
}

// TestFieldExpr_GetSet tests field expression get/set
func TestFieldExpr_GetSet(t *testing.T) {
	user := &TestUser{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// Create field expressions
	usernameField := NewFieldExpr(
		"username",
		func(u *TestUser) interface{} { return u.Username },
		func(u *TestUser, v interface{}) { u.Username = v.(string) },
		schema.Field{Name: "username", Type: schema.TypeString},
	)

	emailField := NewFieldExpr(
		"email",
		func(u *TestUser) interface{} { return u.Email },
		func(u *TestUser, v interface{}) { u.Email = v.(string) },
		schema.Field{Name: "email", Type: schema.TypeString},
	)

	isActiveField := NewFieldExpr(
		"is_active",
		func(u *TestUser) interface{} { return u.IsActive },
		func(u *TestUser, v interface{}) { u.IsActive = v.(bool) },
		schema.Field{Name: "is_active", Type: schema.TypeBool},
	)

	// Test Get
	if usernameField.Get(user) != "testuser" {
		t.Errorf("Expected username 'testuser', got '%v'", usernameField.Get(user))
	}

	if emailField.Get(user) != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%v'", emailField.Get(user))
	}

	if !isActiveField.Get(user).(bool) {
		t.Error("Expected IsActive to be true")
	}

	// Test Set
	usernameField.Set(user, "newuser")
	if user.Username != "newuser" {
		t.Errorf("Expected username 'newuser', got '%s'", user.Username)
	}

	emailField.Set(user, "new@example.com")
	if user.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got '%s'", user.Email)
	}

	isActiveField.Set(user, false)
	if user.IsActive {
		t.Error("Expected IsActive to be false")
	}
}

// TestFieldExpr_Name tests field name
func TestFieldExpr_Name(t *testing.T) {
	field := NewFieldExpr(
		"username",
		func(u *TestUser) interface{} { return u.Username },
		func(u *TestUser, v interface{}) { u.Username = v.(string) },
		schema.Field{Name: "username", Type: schema.TypeString},
	)

	if field.Name() != "username" {
		t.Errorf("Expected field name 'username', got '%s'", field.Name())
	}
}

// TestRegistry_RegisterGet tests registry operations
func TestRegistry_RegisterGet(t *testing.T) {
	// Use global registry
	registry := GetGlobalRegistry()

	// Create test admin
	manager := &query.Manager[*TestUser]{}
	config := &Config[*TestUser]{}
	_ = Register(&TestUser{}, manager, config)

	// Test Get
	retrieved, err := registry.Get("TestUser")
	if err != nil {
		t.Fatalf("Failed to get admin: %v", err)
	}

	if retrieved.ModelName() != "TestUser" {
		t.Errorf("Expected ModelName 'TestUser', got '%s'", retrieved.ModelName())
	}

	// Test Get non-existent
	_, err = registry.Get("NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent admin")
	}
}

// TestRegistry_GetAll tests getting all admins
func TestRegistry_GetAll(t *testing.T) {
	// Use global registry
	registry := GetGlobalRegistry()

	// Register multiple admins
	manager1 := &query.Manager[*TestUser]{}
	_ = Register(&TestUser{}, manager1, &Config[*TestUser]{})

	// Get all
	all := registry.GetAll()

	if len(all) < 1 {
		t.Error("Expected at least one admin")
	}

	if _, ok := all["TestUser"]; !ok {
		t.Error("Expected TestUser admin to be registered")
	}
}

// TestBooleanFilter tests boolean filter
func TestBooleanFilter(t *testing.T) {
	// Create bool field expression directly for filter
	boolField := NewFieldExpr(
		"is_active",
		func(u *TestUser) bool { return u.IsActive },
		func(u *TestUser, v bool) { u.IsActive = v },
		schema.Field{Name: "is_active", Type: schema.TypeBool},
	)

	filter := NewBooleanFilter(boolField)

	if filter.Name() != "is_active" {
		t.Errorf("Expected filter name 'is_active', got '%s'", filter.Name())
	}

	// Test GetOptions
	ctx := context.Background()
	options, err := filter.GetOptions(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get options: %v", err)
	}

	if len(options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(options))
	}
}

// TestChoiceFilter tests choice filter
func TestChoiceFilter(t *testing.T) {
	// Create string field expression directly for filter
	stringField := NewFieldExpr(
		"status",
		func(u *TestUser) string { return "" },
		func(u *TestUser, v string) {},
		schema.Field{Name: "status", Type: schema.TypeString},
	)

	choices := []Choice[string]{
		{Label: "Active", Value: "active"},
		{Label: "Inactive", Value: "inactive"},
	}

	filter := NewChoiceFilter(stringField, choices)

	if filter.Name() != "status" {
		t.Errorf("Expected filter name 'status', got '%s'", filter.Name())
	}

	// Test GetOptions
	ctx := context.Background()
	options, err := filter.GetOptions(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to get options: %v", err)
	}

	if len(options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(options))
	}
}

// TestAction tests action creation
func TestAction(t *testing.T) {
	action := NewAction(
		"activate",
		"Activate selected users",
		func(ctx context.Context, users []*TestUser) error {
			for _, user := range users {
				user.IsActive = true
			}
			return nil
		},
	)

	if action.Name != "activate" {
		t.Errorf("Expected action name 'activate', got '%s'", action.Name)
	}

	if action.Label != "Activate selected users" {
		t.Errorf("Expected action label 'Activate selected users', got '%s'", action.Label)
	}

	// Test handler
	ctx := context.Background()
	users := []*TestUser{
		{IsActive: false},
		{IsActive: false},
	}

	err := action.Handler(ctx, users)
	if err != nil {
		t.Fatalf("Action handler failed: %v", err)
	}

	for _, user := range users {
		if !user.IsActive {
			t.Error("Expected all users to be activated")
		}
	}
}

// TestFieldset tests fieldset creation
func TestFieldset(t *testing.T) {
	field1 := NewFieldExpr(
		"username",
		func(u *TestUser) interface{} { return u.Username },
		func(u *TestUser, v interface{}) { u.Username = v.(string) },
		schema.Field{Name: "username", Type: schema.TypeString},
	)

	field2 := NewFieldExpr(
		"email",
		func(u *TestUser) interface{} { return u.Email },
		func(u *TestUser, v interface{}) { u.Email = v.(string) },
		schema.Field{Name: "email", Type: schema.TypeString},
	)

	fieldset := NewFieldset("Account Information", field1, field2)

	if fieldset.Name != "Account Information" {
		t.Errorf("Expected fieldset name 'Account Information', got '%s'", fieldset.Name)
	}

	if len(fieldset.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fieldset.Fields))
	}
}

// TestOrdering tests ordering
func TestOrdering(t *testing.T) {
	field := NewFieldExpr(
		"username",
		func(u *TestUser) interface{} { return u.Username },
		func(u *TestUser, v interface{}) { u.Username = v.(string) },
		schema.Field{Name: "username", Type: schema.TypeString},
	)

	// Test ascending
	orderingAsc := OrderBy(field).Asc()
	if !orderingAsc.IsAscending() {
		t.Error("Expected ascending order")
	}

	// Test descending
	orderingDesc := OrderBy(field).Desc()
	if !orderingDesc.IsDescending() {
		t.Error("Expected descending order")
	}
}

// TestGetGlobalRegistry tests global registry access
func TestGetGlobalRegistry(t *testing.T) {
	registry := GetGlobalRegistry()
	if registry == nil {
		t.Fatal("Global registry should not be nil")
	}
}
