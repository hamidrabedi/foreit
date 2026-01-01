package permissions

import (
	"context"
	"fmt"
)

// ObjectPermissionChecker checks permissions for specific object instances
type ObjectPermissionChecker interface {
	// CanView checks if user can view a specific object
	CanView(ctx context.Context, user interface{}, obj interface{}) (bool, error)

	// CanChange checks if user can change a specific object
	CanChange(ctx context.Context, user interface{}, obj interface{}) (bool, error)

	// CanDelete checks if user can delete a specific object
	CanDelete(ctx context.Context, user interface{}, obj interface{}) (bool, error)
}

// ObjectPermissionRule represents a rule for object-level permissions
type ObjectPermissionRule struct {
	Name        string
	Description string
	Checker     func(ctx context.Context, user interface{}, obj interface{}) (bool, error)
}

// ObjectPermissionManager manages object-level permissions
type ObjectPermissionManager struct {
	rules []ObjectPermissionRule
}

// NewObjectPermissionManager creates a new object permission manager
func NewObjectPermissionManager() *ObjectPermissionManager {
	return &ObjectPermissionManager{
		rules: make([]ObjectPermissionRule, 0),
	}
}

// AddRule adds a permission rule
func (m *ObjectPermissionManager) AddRule(rule ObjectPermissionRule) {
	m.rules = append(m.rules, rule)
}

// CanView checks if user can view an object
func (m *ObjectPermissionManager) CanView(ctx context.Context, user interface{}, obj interface{}) (bool, error) {
	// Check all rules - all must pass (AND logic)
	for _, rule := range m.rules {
		allowed, err := rule.Checker(ctx, user, obj)
		if err != nil {
			return false, fmt.Errorf("permission check failed for rule %s: %w", rule.Name, err)
		}
		if !allowed {
			return false, nil
		}
	}
	return true, nil
}

// CanChange checks if user can change an object
func (m *ObjectPermissionManager) CanChange(ctx context.Context, user interface{}, obj interface{}) (bool, error) {
	// Similar to CanView
	return m.CanView(ctx, user, obj)
}

// CanDelete checks if user can delete an object
func (m *ObjectPermissionManager) CanDelete(ctx context.Context, user interface{}, obj interface{}) (bool, error) {
	// Similar to CanView
	return m.CanView(ctx, user, obj)
}

// OwnerBasedRule creates a rule based on object ownership
func OwnerBasedRule(ownerField string) ObjectPermissionRule {
	return ObjectPermissionRule{
		Name:        "owner_based",
		Description: fmt.Sprintf("User must be owner (field: %s)", ownerField),
		Checker: func(ctx context.Context, user interface{}, obj interface{}) (bool, error) {
			// Extract owner from object
			// Extract user ID from user
			// Compare them
			// This is a simplified version - full implementation would use reflection
			return true, nil // Placeholder
		},
	}
}

// FieldBasedRule creates a rule based on a field value
func FieldBasedRule(fieldName string, allowedValues []interface{}) ObjectPermissionRule {
	return ObjectPermissionRule{
		Name:        "field_based",
		Description: fmt.Sprintf("Object field %s must be in allowed values", fieldName),
		Checker: func(ctx context.Context, user interface{}, obj interface{}) (bool, error) {
			// Extract field value from object
			// Check if it's in allowed values
			// This is a simplified version
			return true, nil // Placeholder
		},
	}
}

// CustomRule creates a custom permission rule
func CustomRule(name, description string, checker func(ctx context.Context, user interface{}, obj interface{}) (bool, error)) ObjectPermissionRule {
	return ObjectPermissionRule{
		Name:        name,
		Description: description,
		Checker:     checker,
	}
}
