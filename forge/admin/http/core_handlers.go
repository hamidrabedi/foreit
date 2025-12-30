package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/forgego/forge/admin"
	adminutils "github.com/forgego/forge/admin/utils"
	admintemplates "github.com/forgego/forge/admin/templates"
	"github.com/forgego/forge/server"
	httplib "github.com/forgego/forge/server"
)

// CoreHandler provides HTTP handlers for the new core admin system
type CoreHandler struct {
	registry       *admin.Registry
	renderer       *admintemplates.Renderer
	sessionManager *server.SessionManager
}

// NewCoreHandler creates a new core HTTP handler
func NewCoreHandler(registry *admin.Registry, renderer *admintemplates.Renderer, sessionManager *server.SessionManager) *CoreHandler {
	return &CoreHandler{
		registry:       registry,
		renderer:       renderer,
		sessionManager: sessionManager,
	}
}

// HandleList handles list view requests using the new core system
func (h *CoreHandler) HandleList(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use type registry to get type-safe admin
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		page := httplib.GetQueryInt(r, "page", 1)
		pageSize := httplib.GetQueryInt(r, "page_size", 20)
		search := httplib.GetQueryString(r, "search", "")

		// Get filters from query params
		filters := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if key != "page" && key != "page_size" && key != "search" {
				if len(values) > 0 {
					filters[key] = values[0]
				}
			}
		}

		// Call type-safe handler
		data, err := handler.HandleList(ctx, page, pageSize, search, filters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render template or return JSON based on Accept header
		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		} else {
			// Render HTML template
			h.renderListView(w, r, modelName, data)
		}
	}
}

// renderListView renders the list view template
func (h *CoreHandler) renderListView(w http.ResponseWriter, r *http.Request, modelName string, data interface{}) {
	_ = r // May be used for user context later
	templateData := map[string]interface{}{
		"Title":      fmt.Sprintf("%s List", modelName),
		"SiteTitle":  "Admin",
		"SiteURL":    "/admin",
		"ModelName":  modelName,
		"Data":       data,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": ""},
		},
	}

	if err := h.renderer.Render(w, "list", templateData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

// HandleDetail handles detail view requests
func (h *CoreHandler) HandleDetail(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := httplib.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		data, err := handler.HandleDetail(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		} else {
			h.renderDetailView(w, r, modelName, id, data)
		}
	}
}

// renderDetailView renders the detail view template
func (h *CoreHandler) renderDetailView(w http.ResponseWriter, r *http.Request, modelName string, id int64, data interface{}) {
	templateData := map[string]interface{}{
		"Title":     fmt.Sprintf("%s #%d", modelName, id),
		"SiteTitle": "Admin",
		"SiteURL":   "/admin",
		"ModelName": modelName,
		"InstanceID": id,
		"Data":      data,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": fmt.Sprintf("/admin/%s/", modelName)},
			{"Label": fmt.Sprintf("#%d", id), "URL": ""},
		},
	}

	if err := h.renderer.Render(w, "detail", templateData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

// HandleCreate handles create form requests
func (h *CoreHandler) HandleCreate(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.handleCreateGet(w, r, modelName)
		} else if r.Method == http.MethodPost {
			h.handleCreatePost(w, r, modelName)
		}
	}
}

// handleCreateGet handles GET request for create form
func (h *CoreHandler) handleCreateGet(w http.ResponseWriter, r *http.Request, modelName string) {
	_, err := h.registry.Get(modelName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
		return
	}

	// Use views to render form
	// This would need to be implemented with proper type handling
	templateData := map[string]interface{}{
		"Title":     fmt.Sprintf("Add %s", modelName),
		"SiteTitle": "Admin",
		"SiteURL":   "/admin",
		"ModelName": modelName,
		"IsNew":     true,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": fmt.Sprintf("/admin/%s/", modelName)},
			{"Label": "Add", "URL": ""},
		},
	}

	if err := h.renderer.Render(w, "form", templateData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

// handleCreatePost handles POST request for create form
func (h *CoreHandler) handleCreatePost(w http.ResponseWriter, r *http.Request, modelName string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formData := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	ctx := r.Context()
	instance, err := handler.HandleCreate(ctx, formData)
	if err != nil {
		adminutils.Error(ctx, r, fmt.Sprintf("Failed to create: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for ResponseAddHook (using reflection to access hooks)
	// Note: Full hook integration requires type-safe access which is complex
	// For now, hooks are called from within views where type is known

	// Show success message
	adminutils.Success(ctx, r, fmt.Sprintf("%s was added successfully.", modelName))

	// Redirect to detail page
	id := adminutils.GetIDFromInstance(instance)
	if id > 0 {
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/%d/", modelName, id), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/", modelName), http.StatusSeeOther)
	}
}

// HandleUpdate handles update form requests
func (h *CoreHandler) HandleUpdate(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := httplib.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			h.handleUpdateGet(w, r, modelName, id)
		} else if r.Method == http.MethodPost {
			h.handleUpdatePost(w, r, modelName, id)
		}
	}
}

// handleUpdateGet handles GET request for update form
func (h *CoreHandler) handleUpdateGet(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	instance, err := handler.HandleDetail(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	templateData := map[string]interface{}{
		"Title":     fmt.Sprintf("Change %s", modelName),
		"SiteTitle": "Admin",
		"SiteURL":   "/admin",
		"ModelName": modelName,
		"InstanceID": id,
		"Instance": instance,
		"IsNew":     false,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": fmt.Sprintf("/admin/%s/", modelName)},
			{"Label": fmt.Sprintf("#%d", id), "URL": fmt.Sprintf("/admin/%s/%d/", modelName, id)},
			{"Label": "Change", "URL": ""},
		},
	}

	if err := h.renderer.Render(w, "form", templateData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

// handleUpdatePost handles POST request for update form
func (h *CoreHandler) handleUpdatePost(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formData := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	ctx := r.Context()
	_, err := handler.HandleUpdate(ctx, id, formData)
	if err != nil {
		adminutils.Error(ctx, r, fmt.Sprintf("Failed to update: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for ResponseChangeHook (using reflection to access hooks)
	// Note: Full hook integration requires type-safe access which is complex
	// For now, hooks are called from within views where type is known

	// Show success message
	adminutils.Success(ctx, r, fmt.Sprintf("%s was changed successfully.", modelName))

	// Redirect to detail page
	http.Redirect(w, r, fmt.Sprintf("/admin/%s/%d/", modelName, id), http.StatusSeeOther)
}

// HandleDelete handles delete requests
func (h *CoreHandler) HandleDelete(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := httplib.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	ctx := r.Context()

	if err := handler.HandleDelete(ctx, id); err != nil {
		adminutils.Error(ctx, r, fmt.Sprintf("Failed to delete: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Show success message
	adminutils.Success(ctx, r, fmt.Sprintf("%s was deleted successfully.", modelName))

	// Redirect to list page
	http.Redirect(w, r, fmt.Sprintf("/admin/%s/", modelName), http.StatusSeeOther)
	}
}

// HandleIndex handles admin index/dashboard
func (h *CoreHandler) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAdmins := h.registry.GetAll()

		models := make([]map[string]interface{}, 0, len(allAdmins))
		for name, admin := range allAdmins {
			models = append(models, map[string]interface{}{
				"name":        name,
				"verboseName": name,
				"modelType":   admin.ModelType().String(),
				"url":         fmt.Sprintf("/admin/%s/", name),
			})
		}

		templateData := map[string]interface{}{
			"Title":     "Admin Dashboard",
			"SiteTitle": "Admin",
			"SiteURL":   "/admin",
			"Models":    models,
		}

		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(templateData)
		} else {
			if err := h.renderer.Render(w, "index", templateData); err != nil {
				// Fallback to JSON if template not found
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(templateData)
			}
		}
	}
}
