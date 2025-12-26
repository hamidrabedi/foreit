package admin

import (
	"github.com/forgego/forge/pkg/models"
)

// HasPerm checks if a user has a specific permission
// For now, this is a simplified implementation
// Full implementation would check against a permissions table
func HasPerm(user *models.User, perm string) bool {
	// Superusers have all permissions
	if user.IsSuperuser {
		return true
	}
	
	// For MVP, we only check built-in permissions
	switch perm {
	case "admin.view_user", "admin.add_user", "admin.change_user", "admin.delete_user":
		return user.IsStaff
	default:
		// Default: staff users have access to admin
		return user.IsStaff
	}
}

// HasPerms checks if a user has all specified permissions
func HasPerms(user *models.User, perms []string) bool {
	for _, perm := range perms {
		if !HasPerm(user, perm) {
			return false
		}
	}
	return true
}

// HasAnyPerm checks if a user has any of the specified permissions
func HasAnyPerm(user *models.User, perms []string) bool {
	for _, perm := range perms {
		if HasPerm(user, perm) {
			return true
		}
	}
	return false
}

// IsActive checks if user is active (helper function)
func IsActiveUser(user *models.User) bool {
	return user.IsActive
}

