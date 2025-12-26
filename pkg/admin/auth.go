package admin

import (
	"fmt"
	"net/http"

	"github.com/forgego/forge/pkg/admin/templates"
	"github.com/forgego/forge/pkg/auth"
	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/security"
)

// handleLogin handles GET request for login page
func handleLogin(sessionManager *security.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If already authenticated, redirect to admin
		if auth.IsAuthenticated(sessionManager, r) {
			http.Redirect(w, r, "/admin/", http.StatusFound)
			return
		}

		// Get next parameter for redirect after login
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/admin/"
		}

		// Render login template
		data := map[string]interface{}{
			"Title": "Log in",
			"Next":  next,
			"Error": "",
		}

		if err := renderTemplate(w, "login.html", data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render login: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// handleLoginPost handles POST request for login
func handleLoginPost(sessionManager *security.SessionManager, userManager *models.UserManagerImpl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		usernameOrEmail := r.FormValue("username")
		password := r.FormValue("password")
		next := r.FormValue("next")

		if next == "" {
			next = "/admin/"
		}

		// Validate inputs
		if usernameOrEmail == "" || password == "" {
			data := map[string]interface{}{
				"Title": "Log in",
				"Next":  next,
				"Error": "Please enter both username and password.",
			}
			renderTemplate(w, "login.html", data)
			return
		}

		// Authenticate user
		ctx := r.Context()
		user, err := userManager.Authenticate(ctx, usernameOrEmail, password)
		if err != nil {
			data := map[string]interface{}{
				"Title": "Log in",
				"Next":  next,
				"Error": "Please enter a correct username/email and password. Note that both fields may be case-sensitive.",
			}
			renderTemplate(w, "login.html", data)
			return
		}

		// Check if user is staff (required for admin access)
		if !user.IsStaff {
			data := map[string]interface{}{
				"Title": "Log in",
				"Next":  next,
				"Error": "You don't have permission to access the admin site.",
			}
			renderTemplate(w, "login.html", data)
			return
		}

		// Create session
		if err := auth.LoginUser(sessionManager, r, user.ID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect to next URL or admin index
		http.Redirect(w, r, next, http.StatusFound)
	}
}

// handleLogout handles POST request for logout
func handleLogout(sessionManager *security.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Destroy session
		if err := auth.LogoutUser(sessionManager, r); err != nil {
			http.Error(w, fmt.Sprintf("Failed to logout: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect to login page
		loginURL := "/admin/login/"
		http.Redirect(w, r, loginURL, http.StatusFound)
	}
}

// getUserFromRequest retrieves the current user from the request context
func getUserFromRequest(r *http.Request, userManager *models.UserManagerImpl) (*models.User, error) {
	// Get user ID from session (this should be set by middleware)
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		return nil, fmt.Errorf("user not authenticated")
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID type")
	}

	// Load user from database
	ctx := r.Context()
	user, err := userManager.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return user, nil
}

// renderTemplate is a helper to render admin templates
func renderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Execute the template
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

