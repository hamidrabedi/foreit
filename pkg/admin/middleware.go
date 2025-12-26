package admin

import (
	"context"
	"net/http"

	"github.com/forgego/forge/pkg/auth"
	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/security"
)

// Context keys for admin package
type adminContextKey string

const (
	// AdminUserKey is the context key for the authenticated admin user
	AdminUserKey adminContextKey = "admin_user"
	// AdminUserIDKey is the context key for the authenticated admin user ID
	AdminUserIDKey adminContextKey = "admin_user_id"
)

// AdminAuthMiddleware creates middleware that requires authentication and staff status
func AdminAuthMiddleware(sessionManager *security.SessionManager, userManager *models.UserManagerImpl) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for login/logout routes
			if r.URL.Path == "/admin/login/" || r.URL.Path == "/admin/logout/" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if user is authenticated
			if !auth.IsAuthenticated(sessionManager, r) {
				// Redirect to login with next parameter
				nextURL := r.URL.Path
				if r.URL.RawQuery != "" {
					nextURL += "?" + r.URL.RawQuery
				}
				http.Redirect(w, r, "/admin/login/?next="+nextURL, http.StatusFound)
				return
			}

			// Get user ID from session
			userID := auth.GetUserID(sessionManager, r)
			if userID == nil {
				http.Redirect(w, r, "/admin/login/", http.StatusFound)
				return
			}

			// Convert userID to int64
			var id int64
			switch v := userID.(type) {
			case int64:
				id = v
			case int:
				id = int64(v)
			case int32:
				id = int64(v)
			default:
				http.Redirect(w, r, "/admin/login/", http.StatusFound)
				return
			}

			// Load user from database
			ctx := r.Context()
			user, err := userManager.Get(ctx, id)
			if err != nil {
				// User not found, clear session and redirect
				_ = auth.LogoutUser(sessionManager, r)
				http.Redirect(w, r, "/admin/login/", http.StatusFound)
				return
			}

			// Check if user is active
			if !user.IsActive {
				_ = auth.LogoutUser(sessionManager, r)
				http.Redirect(w, r, "/admin/login/", http.StatusFound)
				return
			}

			// Check if user is staff (required for admin access)
			if !user.IsStaff {
				_ = auth.LogoutUser(sessionManager, r)
				http.Redirect(w, r, "/admin/login/", http.StatusFound)
				return
			}

			// Add user to context using typed keys
			ctx = context.WithValue(ctx, AdminUserKey, user)
			ctx = context.WithValue(ctx, AdminUserIDKey, user.ID)
			r = r.WithContext(ctx)

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GetUser retrieves the current user from the request context
func GetUser(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(AdminUserKey).(*models.User)
	return user, ok
}

// GetUserID retrieves the current user ID from the request context
func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(AdminUserIDKey).(int64)
	return userID, ok
}

// IsAuthenticated checks if the user is authenticated (from context)
func IsAuthenticated(r *http.Request) bool {
	_, ok := GetUser(r)
	return ok
}

// IsStaff checks if the user is staff (from context)
func IsStaff(r *http.Request) bool {
	user, ok := GetUser(r)
	if !ok {
		return false
	}
	return user.IsStaff
}

// IsSuperuser checks if the user is superuser (from context)
func IsSuperuser(r *http.Request) bool {
	user, ok := GetUser(r)
	if !ok {
		return false
	}
	return user.IsSuperuser
}

