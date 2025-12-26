package auth

import (
	"net/http"

	"github.com/forgego/forge/pkg/security"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(sessionManager *security.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user is authenticated
			userID := sessionManager.Get(r, "user_id")
			if userID == nil {
				// Redirect to login or return 401
				http.Redirect(w, r, "/admin/login", http.StatusFound)
				return
			}

			// User is authenticated, continue
			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuthMiddleware creates middleware that sets user if authenticated but doesn't require it
func OptionalAuthMiddleware(sessionManager *security.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user is authenticated
			userID := sessionManager.Get(r, "user_id")
			if userID != nil {
				// Set user in context for optional use
				// This allows views to check authentication but not require it
			}

			next.ServeHTTP(w, r)
		})
	}
}

// LoginUser logs in a user by setting session
func LoginUser(sessionManager *security.SessionManager, r *http.Request, userID interface{}) error {
	sessionManager.Put(r, "user_id", userID)
	return nil
}

// LogoutUser logs out a user by destroying session
func LogoutUser(sessionManager *security.SessionManager, r *http.Request) error {
	return sessionManager.Destroy(r)
}

// GetUserID retrieves the current user ID from session
func GetUserID(sessionManager *security.SessionManager, r *http.Request) interface{} {
	return sessionManager.Get(r, "user_id")
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(sessionManager *security.SessionManager, r *http.Request) bool {
	return sessionManager.Exists(r, "user_id")
}
