package auth

import (
	"github.com/gofiber/fiber/v2"
)

// Role represents a user role
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleGuest    Role = "guest"
	RoleModerator Role = "moderator"
)

// RoleChecker checks if a user has a role
type RoleChecker interface {
	HasRole(user interface{}, role Role) bool
	HasAnyRole(user interface{}, roles ...Role) bool
	HasAllRoles(user interface{}, roles ...Role) bool
}

// RoleUser interface for users with roles
type RoleUser interface {
	GetRoles() []Role
}

// HasRole checks if a user has a specific role
func HasRole(user interface{}, role Role) bool {
	if roleUser, ok := user.(RoleUser); ok {
		roles := roleUser.GetRoles()
		for _, r := range roles {
			if r == role {
				return true
			}
		}
	}
	return false
}

// HasAnyRole checks if a user has any of the specified roles
func HasAnyRole(user interface{}, roles ...Role) bool {
	if roleUser, ok := user.(RoleUser); ok {
		userRoles := roleUser.GetRoles()
		for _, userRole := range userRoles {
			for _, requiredRole := range roles {
				if userRole == requiredRole {
					return true
				}
			}
		}
	}
	return false
}

// HasAllRoles checks if a user has all of the specified roles
func HasAllRoles(user interface{}, roles ...Role) bool {
	if roleUser, ok := user.(RoleUser); ok {
		userRoles := roleUser.GetRoles()
		roleMap := make(map[Role]bool)
		for _, r := range userRoles {
			roleMap[r] = true
		}
		
		for _, requiredRole := range roles {
			if !roleMap[requiredRole] {
				return false
			}
		}
		return true
	}
	return false
}

// RequireRole middleware requires a specific role
func RequireRole(role Role) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		
		if !HasRole(user, role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
		
		return c.Next()
	}
}

// RequireAnyRole middleware requires any of the specified roles
func RequireAnyRole(roles ...Role) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		
		if !HasAnyRole(user, roles...) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}
		
		return c.Next()
	}
}

