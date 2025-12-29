package admin

import (
	"context"
	"fmt"
)

// PermissionChecker checks permissions for admin operations
type PermissionChecker interface {
	HasPermission(ctx context.Context, user interface{}, perm string) bool
	HasPermissions(ctx context.Context, user interface{}, perms []string) bool
	HasAnyPermission(ctx context.Context, user interface{}, perms []string) bool
}

// DefaultPermissionChecker provides default permission checking
type DefaultPermissionChecker struct{}

// NewDefaultPermissionChecker creates a new default permission checker
func NewDefaultPermissionChecker() PermissionChecker {
	return &DefaultPermissionChecker{}
}

// HasPermission checks if user has a permission
func (c *DefaultPermissionChecker) HasPermission(ctx context.Context, user interface{}, perm string) bool {
	// Default implementation - always allow
	// In production, this would check against a permissions system
	return true
}

// HasPermissions checks if user has all permissions
func (c *DefaultPermissionChecker) HasPermissions(ctx context.Context, user interface{}, perms []string) bool {
	for _, perm := range perms {
		if !c.HasPermission(ctx, user, perm) {
			return false
		}
	}
	return true
}

// HasAnyPermission checks if user has any permission
func (c *DefaultPermissionChecker) HasAnyPermission(ctx context.Context, user interface{}, perms []string) bool {
	for _, perm := range perms {
		if c.HasPermission(ctx, user, perm) {
			return true
		}
	}
	return false
}

// Permission names
const (
	PermAdd    = "add"
	PermChange = "change"
	PermDelete = "delete"
	PermView   = "view"
)

// GetPermissionName returns permission name for a model and action
func GetPermissionName(modelName, action string) string {
	return fmt.Sprintf("%s.%s_%s", "admin", action, modelName)
}
