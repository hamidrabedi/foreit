package admin

import (
	"testing"
)

// TestPermissionChecker_CheckBasicPermissions tests basic permission checking
func TestPermissionChecker_CheckBasicPermissions(t *testing.T) {
	registry := NewRegistry()
	checker := NewPermissionChecker(registry)
	
	meta := &ModelMeta{
		Name: "TestModel",
		Permissions: &Permissions{
			CanList:   true,
			CanView:   true,
			CanCreate: false,
			CanUpdate: false,
			CanDelete: false,
		},
	}
	
	ctx := &Context{
		Model:  meta,
		Action: "list",
	}
	
	allowed, err := checker.CheckPermission(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Error("Expected permission to be allowed")
	}
	
	// Test create (should be denied)
	ctx.Action = "create"
	allowed, err = checker.CheckPermission(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if allowed {
		t.Error("Expected permission to be denied")
	}
}

// TestPermissionChecker_CheckRules tests rule-based permission checking
func TestPermissionChecker_CheckRules(t *testing.T) {
	registry := NewRegistry()
	checker := NewPermissionChecker(registry)
	
	meta := &ModelMeta{
		Name: "TestModel",
		Permissions: &Permissions{
			CanList:   true,
			CanView:   true,
			CanCreate: true,
			CanUpdate: true,
			CanDelete: true,
			Rules: map[string]string{
				"list": "true",
				"view": "published == true",
			},
		},
	}
	
	// Test rule that always allows
	ctx := &Context{
		Model:  meta,
		Action: "list",
		User:   &User{ID: "user1"},
	}
	
	allowed, err := checker.CheckPermission(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Error("Expected permission to be allowed (rule: true)")
	}
	
	// Test rule with field check (will need resource)
	ctx.Action = "view"
	ctx.Resource = map[string]interface{}{
		"published": true,
	}
	
	allowed, err = checker.CheckPermission(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Error("Expected permission to be allowed (published == true)")
	}
}

// TestPermissionChecker_EvaluateRule tests rule evaluation
func TestPermissionChecker_EvaluateRule(t *testing.T) {
	registry := NewRegistry()
	checker := NewPermissionChecker(registry)
	
	tests := []struct {
		name     string
		rule     string
		ctx      *Context
		expected bool
	}{
		{
			name: "always true",
			rule: "true",
			ctx: &Context{
				Model: &ModelMeta{Name: "Test"},
			},
			expected: true,
		},
		{
			name: "always false",
			rule: "false",
			ctx: &Context{
				Model: &ModelMeta{Name: "Test"},
			},
			expected: false,
		},
		{
			name: "equality check",
			rule: "status == \"active\"",
			ctx: &Context{
				Model:    &ModelMeta{Name: "Test"},
				Resource: map[string]interface{}{"status": "active"},
			},
			expected: true,
		},
		{
			name: "context variable",
			rule: "@request.auth.id == \"user1\"",
			ctx: &Context{
				Model: &ModelMeta{Name: "Test"},
				User:  &User{ID: "user1"},
			},
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := checker.evaluateRule(tt.rule, tt.ctx)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestPermissionChecker_GetAuthID tests auth ID extraction
func TestPermissionChecker_GetAuthID(t *testing.T) {
	registry := NewRegistry()
	checker := NewPermissionChecker(registry)
	
	tests := []struct {
		name     string
		user     interface{}
		expected string
	}{
		{
			name: "user with ID field",
			user: struct {
				ID string
			}{ID: "user123"},
			expected: "user123",
		},
		{
			name:     "nil user",
			user:     nil,
			expected: "",
		},
		{
			name: "user with UserID field",
			user: struct {
				UserID string
			}{UserID: "user456"},
			expected: "user456",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{User: tt.user}
			result := checker.getAuthID(ctx)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

