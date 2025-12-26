package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/forgego/forge/pkg/admin/templates"
	"github.com/forgego/forge/pkg/db"
	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/logging"
	"github.com/forgego/forge/pkg/models"
	"github.com/forgego/forge/pkg/security"
	"go.uber.org/zap"
)

// RegisterAdminRoutes registers admin routes on the router
// sessionManager and database are optional - if provided, login/logout routes and auth middleware will be registered
func RegisterAdminRoutes(router *httplib.Router, path string, sessionManager interface{}, database interface{}) {
	router.Route(path, func(r *httplib.Router) {
		var userManager *models.UserManagerImpl
		
		// Login/logout routes (if session manager and database provided)
		if sessionManager != nil && database != nil {
			sm, ok1 := sessionManager.(*security.SessionManager)
			db, ok2 := database.(*db.DB)
			if ok1 && ok2 {
				userManager = models.NewUserManager(db)
				r.Get("/login/", handleLogin(sm))
				r.Post("/login/", handleLoginPost(sm, userManager))
				r.Post("/logout/", handleLogout(sm))
				
				// Auto-register User model for admin with manager
				userModel := &models.User{}
				if err := RegisterModelWithManager(userModel, userManager); err != nil {
					// Log error but continue - use no-op logger as fallback since we don't have logger instance
					logger := logging.NewNopLogger()
					logger.Warn("Failed to auto-register User model for admin",
						zap.Error(err),
					)
				}
				
				// Apply auth middleware to all other routes
				r.Use(AdminAuthMiddleware(sm, userManager))
			}
		}
		
		// Static files
		r.Get("/static/*", handleStaticFiles)
		
		// Admin index
		r.Get("/", handleAdminIndex)
		
		// Model routes
		adminModels := GetAllModels()
		for modelName, model := range adminModels {
			registerModelRoutes(r, modelName, model)
		}
	})
}

// registerModelRoutes registers routes for a specific model
func registerModelRoutes(router *httplib.Router, modelName string, model *AdminModel) {
	router.Route("/"+modelName, func(r *httplib.Router) {
		// List view
		r.Get("/", handleModelList(modelName, model))
		
		// Export
		r.Get("/export/", handleExport(modelName, model))
		
		// Bulk actions
		r.Post("/bulk-action/", handleBulkAction(modelName, model))
		
		// Create view
		r.Get("/new/", handleModelCreate(modelName, model))
		r.Post("/new/", handleModelCreatePost(modelName, model))
		
		// Detail view
		r.Get("/{id}/", handleModelDetail(modelName, model))
		
		// Update view
		r.Get("/{id}/change/", handleModelUpdate(modelName, model))
		r.Post("/{id}/change/", handleModelUpdatePost(modelName, model))
		
		// Inline field editing
		r.Get("/{id}/edit-field/{field}/", handleInlineFieldEdit(modelName, model))
		r.Post("/{id}/edit-field/{field}/", handleInlineFieldEditPost(modelName, model))
		
		// Delete view
		r.Delete("/{id}/delete/", handleModelDelete(modelName, model))
		r.Post("/{id}/delete/", handleModelDelete(modelName, model))
	})
}

// handleAdminIndex handles the admin index/dashboard
func handleAdminIndex(w http.ResponseWriter, r *http.Request) {
	// Get current user from context (if authenticated)
	var username string
	user, ok := GetUser(r)
	if ok {
		username = user.Username
	} else {
		username = "Guest"
	}

	// Get all registered models
	adminModels := GetAllModels()
	
	// Prepare model list with actual counts
	modelList := make([]map[string]interface{}, 0, len(adminModels))
	ctx := r.Context()
	
	for name, model := range adminModels {
		count := int64(0)
		
		// Try to get count from manager if available
		if model.Manager != nil {
			ops := NewManagerOps()
			if countVal, err := ops.GetCount(ctx, model.Manager); err == nil {
				count = countVal
			}
		}
		
		modelList = append(modelList, map[string]interface{}{
			"Name":        name,
			"VerboseName": model.Name,
			"Count":       count,
		})
	}

	// Prepare template data
	data := map[string]interface{}{
		"Title":    "Admin Dashboard",
		"Username": username,
		"User":     user,
		"Models":   modelList,
	}

	// Render template
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load templates: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "index", data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render index: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleModelList is now in list.go

// handleModelCreate and handleModelCreatePost are now in forms.go

func handleModelDetail(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canView(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to view this object.", http.StatusForbidden)
			return
		}

		ctx := r.Context()

		// Get ID from URL
		idStr := httplib.GetParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		ops := NewManagerOps()
		instance, err := ops.GetInstance(ctx, manager, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get instance: %v", err), http.StatusNotFound)
			return
		}

		// Generate form fields for read-only display
		fields := generateFormFields(model, instance, false)
		for i := range fields {
			fields[i].ReadOnly = true
		}

		data := FormData{
			Model:    model,
			Instance: instance,
			Fields:   fields,
			Errors:   make(map[string]string),
			BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
			IsCreate: false,
		}

		if err := renderFormTemplate(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render detail: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// handleModelUpdate and handleModelUpdatePost are now in forms.go

func handleModelDelete(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canDelete(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to delete this object.", http.StatusForbidden)
			return
		}

		ctx := r.Context()

		// Get ID from URL
		idStr := httplib.GetParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		// Get instance first, then delete it
		ops := NewManagerOps()
		instance, err := ops.GetInstance(ctx, manager, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get instance: %v", err), http.StatusNotFound)
			return
		}

		if err := ops.DeleteInstance(ctx, manager, instance); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete: %v", err), http.StatusInternalServerError)
			return
		}

		// Return empty response for HTMX swap (removes the row)
		w.WriteHeader(http.StatusOK)
	}
}

// handleStaticFiles serves static files from the embedded filesystem
func handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Get the file path from URL (chi router wildcard)
	path := httplib.URLParam(r, "*")
	if path == "" {
		http.NotFound(w, r)
		return
	}
	
	// Remove leading slash if present
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	
	// Serve from embedded filesystem
	serveStaticFile(w, r, path)
}

