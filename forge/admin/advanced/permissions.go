package advanced

import (
	"context"
	"fmt"

	"github.com/forgego/forge/admin"
)

// PermissionManager provides advanced permission management
type PermissionManager struct {
	checker admin.PermissionChecker
}

// NewPermissionManager creates a new permission manager
func NewPermissionManager(checker admin.PermissionChecker) *PermissionManager {
	return &PermissionManager{
		checker: checker,
	}
}

// CheckPermission checks if user has a permission
func (pm *PermissionManager) CheckPermission(ctx context.Context, user interface{}, permission string) bool {
	return pm.checker.HasPermission(ctx, user, permission)
}

// CheckPermissions checks if user has all permissions
func (pm *PermissionManager) CheckPermissions(ctx context.Context, user interface{}, permissions []string) bool {
	return pm.checker.HasPermissions(ctx, user, permissions)
}

// CheckAnyPermission checks if user has any permission
func (pm *PermissionManager) CheckAnyPermission(ctx context.Context, user interface{}, permissions []string) bool {
	return pm.checker.HasAnyPermission(ctx, user, permissions)
}

// GetPermissionName generates a permission name for a model and action
func GetPermissionName(appLabel, modelName, action string) string {
	return fmt.Sprintf("%s.%s_%s", appLabel, action, modelName)
}

// Permission represents a permission
type Permission struct {
	Name        string
	Label       string
	Description string
	AppLabel    string
	ModelName   string
	Action      string
}

// NewPermission creates a new permission
func NewPermission(appLabel, modelName, action, label string) *Permission {
	return &Permission{
		Name:      GetPermissionName(appLabel, modelName, action),
		Label:     label,
		AppLabel:  appLabel,
		ModelName: modelName,
		Action:    action,
	}
}

// PermissionSet represents a set of permissions
type PermissionSet struct {
	permissions map[string]*Permission
}

// NewPermissionSet creates a new permission set
func NewPermissionSet() *PermissionSet {
	return &PermissionSet{
		permissions: make(map[string]*Permission),
	}
}

// AddPermission adds a permission to the set
func (ps *PermissionSet) AddPermission(perm *Permission) {
	ps.permissions[perm.Name] = perm
}

// GetPermission gets a permission by name
func (ps *PermissionSet) GetPermission(name string) (*Permission, bool) {
	perm, ok := ps.permissions[name]
	return perm, ok
}

// GetAllPermissions returns all permissions
func (ps *PermissionSet) GetAllPermissions() map[string]*Permission {
	return ps.permissions
}
