package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/forgego/forge/admin"
	httplib "github.com/forgego/forge/server"
)

// Handler provides HTTP handlers for admin views
type Handler struct {
	registry *admin.Registry
}

// NewHandler creates a new admin HTTP handler
func NewHandler(registry *admin.Registry) *Handler {
	return &Handler{
		registry: registry,
	}
}

// HandleList handles list view requests
func (h *Handler) HandleList(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get admin from registry
		adminInterface, err := h.registry.Get(modelName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
			return
		}

		// Create list view using type-safe admin
		// We need to use reflection to get the actual admin instance
		ctx := r.Context()

		// For now, we'll use a generic handler that works with AdminInterface
		// In production, we'd want to use type assertions or a registry that stores type info
		h.handleListGeneric(w, r, ctx, adminInterface)
	}
}

// handleListGeneric handles list view using AdminInterface
func (h *Handler) handleListGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface) {
	// Get pagination params
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

	// Get type-safe handler
	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call type-safe handler
	data, err := handler.HandleList(ctx, page, pageSize, search, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleDetail handles detail view requests
func (h *Handler) HandleDetail(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := httplib.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		adminInterface, err := h.registry.Get(modelName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
			return
		}

		ctx := r.Context()
		h.handleDetailGeneric(w, r, ctx, adminInterface, id)
	}
}

// handleDetailGeneric handles detail view using AdminInterface
func (h *Handler) handleDetailGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface, id int64) {
	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := handler.HandleDetail(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleCreate handles create form requests
func (h *Handler) HandleCreate(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.handleCreateGet(w, r, modelName)
		} else if r.Method == http.MethodPost {
			h.handleCreatePost(w, r, modelName)
		}
	}
}

// handleCreateGet handles GET request for create form
func (h *Handler) handleCreateGet(w http.ResponseWriter, r *http.Request, modelName string) {
	adminInterface, err := h.registry.Get(modelName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Create form - implementation in progress",
		"model":   adminInterface.ModelName(),
	})
}

// handleCreatePost handles POST request for create form
func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request, modelName string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	adminInterface, err := h.registry.Get(modelName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	h.handleCreatePostGeneric(w, r, ctx, adminInterface)
}

// handleCreatePostGeneric handles POST for create using AdminInterface
func (h *Handler) handleCreatePostGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface) {
	// Convert form to map
	formData := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	instance, err := handler.HandleCreate(ctx, formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect to detail page
	id := getIDFromInstance(instance)
	if id > 0 {
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/%d/", adminInterface.ModelName(), id), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/", adminInterface.ModelName()), http.StatusSeeOther)
	}
}

// HandleUpdate handles update form requests
func (h *Handler) HandleUpdate(modelName string) http.HandlerFunc {
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
func (h *Handler) handleUpdateGet(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	adminInterface, err := h.registry.Get(modelName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Update form - implementation in progress",
		"model":   adminInterface.ModelName(),
		"id":      id,
	})
}

// handleUpdatePost handles POST request for update form
func (h *Handler) handleUpdatePost(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	adminInterface, err := h.registry.Get(modelName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	h.handleUpdatePostGeneric(w, r, ctx, adminInterface, id)
}

// handleUpdatePostGeneric handles POST for update using AdminInterface
func (h *Handler) handleUpdatePostGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface, id int64) {
	// Convert form to map
	formData := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = handler.HandleUpdate(ctx, id, formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect to detail page
	http.Redirect(w, r, fmt.Sprintf("/admin/%s/%d/", adminInterface.ModelName(), id), http.StatusSeeOther)
}

// HandleDelete handles delete requests
func (h *Handler) HandleDelete(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := httplib.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		adminInterface, err := h.registry.Get(modelName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
			return
		}

		ctx := r.Context()
		h.handleDeleteGeneric(w, r, ctx, adminInterface, id)
	}
}

// handleDeleteGeneric handles delete using AdminInterface
func (h *Handler) handleDeleteGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface, id int64) {
	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := handler.HandleDelete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to list
	http.Redirect(w, r, fmt.Sprintf("/admin/%s/", adminInterface.ModelName()), http.StatusSeeOther)
}

// HandleIndex handles admin index/dashboard
func (h *Handler) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAdmins := h.registry.GetAll()

		models := make([]map[string]interface{}, 0, len(allAdmins))
		for name := range allAdmins {
			models = append(models, map[string]interface{}{
				"name":        name,
				"verboseName": name,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"title":  "Admin Dashboard",
			"models": models,
		})
	}
}

// HandleExport handles export requests
func (h *Handler) HandleExport(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formatStr := httplib.GetQueryString(r, "format", "csv")
		format := admin.ExportFormat(formatStr)

		adminInterface, err := h.registry.Get(modelName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
			return
		}

		ctx := r.Context()
		h.handleExportGeneric(w, r, ctx, adminInterface, format)
	}
}

// handleExportGeneric handles export using AdminInterface
func (h *Handler) handleExportGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface, format admin.ExportFormat) {
	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The AdminHandler.HandleExport now returns (interface{}, error)
	// We need to cast it to admin.ExportView and then call Export on it.
	exportResult, err := handler.HandleExport(ctx, string(format))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exportView, ok := exportResult.(*admin.ExportView[interface{}]) // Use interface{} for generic type
	if !ok {
		http.Error(w, "Internal server error: invalid export view type", http.StatusInternalServerError)
		return
	}

	if err := exportView.Export(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleBulkAction handles bulk action requests
func (h *Handler) HandleBulkAction(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		action := r.FormValue("action")
		selectedIDs := r.Form["selected_ids"]

		adminInterface, err := h.registry.Get(modelName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
			return
		}

		ctx := r.Context()
		h.handleBulkActionGeneric(w, r, ctx, adminInterface, action, selectedIDs)
	}
}

// handleBulkActionGeneric handles bulk action using AdminInterface
func (h *Handler) handleBulkActionGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface, action string, selectedIDs []string) {
	ids := make([]int64, 0, len(selectedIDs))
	for _, idStr := range selectedIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid ID: %s", idStr), http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}

	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := handler.HandleBulkAction(ctx, action, ids); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/%s/", adminInterface.ModelName()), http.StatusSeeOther)
}

// HandleAutocomplete handles autocomplete requests
func (h *Handler) HandleAutocomplete(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := httplib.GetQueryString(r, "q", "")
		limit := httplib.GetQueryInt(r, "limit", 10)

		adminInterface, err := h.registry.Get(modelName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Model %s not found", modelName), http.StatusNotFound)
			return
		}

		ctx := r.Context()
		h.handleAutocompleteGeneric(w, r, ctx, adminInterface, search, limit)
	}
}

// handleAutocompleteGeneric handles autocomplete using AdminInterface
func (h *Handler) handleAutocompleteGeneric(w http.ResponseWriter, r *http.Request, ctx context.Context, adminInterface admin.AdminInterface, search string, limit int) {
	handler, err := GetAdminHandler(adminInterface.ModelName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := handler.HandleAutocomplete(ctx, search, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// getIDFromInstance extracts ID from an instance using reflection
func getIDFromInstance(instance interface{}) int64 {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanInterface() {
		if id, ok := idField.Interface().(int64); ok {
			return id
		}
	}

	return 0
}
