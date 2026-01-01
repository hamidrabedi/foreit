package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/forgego/forge/admin"
	admintemplates "github.com/forgego/forge/admin/templates"
	adminutils "github.com/forgego/forge/admin/utils"
	"github.com/forgego/forge/admin/views"
	"github.com/forgego/forge/server"
	httplib "github.com/forgego/forge/server"
	"github.com/gorilla/csrf"
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

// HandleIndex handles the admin root index page
func (h *CoreHandler) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAdmins := h.registry.GetAll()
		models := make([]map[string]interface{}, 0, len(allAdmins))
		for name := range allAdmins {
			models = append(models, map[string]interface{}{
				"Name": name,
				"URL":  fmt.Sprintf("/admin/%s/", name),
			})
		}

		templateData := map[string]interface{}{
			"Title":     "Admin Dashboard",
			"SiteTitle": "Admin",
			"SiteURL":   "/admin",
			"Models":    models,
			"User":      h.sessionManager.Get(r, "user"),
		}

		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(templateData)
			return
		}

		if err := h.renderer.Render(w, "index", templateData); err != nil {
			http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		}
	}
}

// HandleList handles list view requests using the new core system
func (h *CoreHandler) HandleList(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		user := h.sessionManager.Get(r, "user")

		data, err := handler.HandleList(ctx, w, r, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		} else {
			h.renderListView(w, r, modelName, data)
		}
	}
}

func (h *CoreHandler) renderListView(w http.ResponseWriter, r *http.Request, modelName string, data interface{}) {
	templateData := map[string]interface{}{
		"Title":     fmt.Sprintf("%s List", modelName),
		"SiteTitle": "Admin",
		"SiteURL":   "/admin",
		"ModelName": modelName,
		"Data":      data,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": ""},
		},
		"CSRFToken": csrf.TemplateField(r),
		"User":      h.sessionManager.Get(r, "user"),
		"Request":   r,
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
		user := h.sessionManager.Get(r, "user")
		data, err := handler.HandleDetail(ctx, w, r, user, id)
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

func (h *CoreHandler) renderDetailView(w http.ResponseWriter, r *http.Request, modelName string, id int64, data interface{}) {
	templateData := map[string]interface{}{
		"Title":      fmt.Sprintf("%s #%d", modelName, id),
		"SiteTitle":  "Admin",
		"SiteURL":    "/admin",
		"ModelName":  modelName,
		"InstanceID": id,
		"Data":       data,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": fmt.Sprintf("/admin/%s/", modelName)},
			{"Label": fmt.Sprintf("#%d", id), "URL": ""},
		},
		"CSRFToken": csrf.TemplateField(r),
		"User":      h.sessionManager.Get(r, "user"),
	}

	if err := h.renderer.Render(w, "detail", templateData); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
	}
}

func (h *CoreHandler) renderFormView(w http.ResponseWriter, r *http.Request, modelName string, data interface{}, isNew bool) {
	title := fmt.Sprintf("Change %s", modelName)
	if isNew {
		title = fmt.Sprintf("Add %s", modelName)
	}

	templateData := map[string]interface{}{
		"Title":     title,
		"SiteTitle": "Admin",
		"SiteURL":   "/admin",
		"ModelName": modelName,
		"Data":      data,
		"IsNew":     isNew,
		"Breadcrumbs": []map[string]interface{}{
			{"Label": "Home", "URL": "/admin"},
			{"Label": modelName, "URL": fmt.Sprintf("/admin/%s/", modelName)},
			{"Label": title, "URL": ""},
		},
		"CSRFToken": csrf.TemplateField(r),
		"User":      h.sessionManager.Get(r, "user"),
	}

	if err := h.renderer.Render(w, "form", templateData); err != nil {
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

func (h *CoreHandler) handleCreateGet(w http.ResponseWriter, r *http.Request, modelName string) {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ctx := r.Context()
	user := h.sessionManager.Get(r, "user")
	data, _ := handler.HandleCreate(ctx, w, r, user, nil)
	h.renderFormView(w, r, modelName, data, true)
}

func (h *CoreHandler) handleCreatePost(w http.ResponseWriter, r *http.Request, modelName string) {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	user := h.sessionManager.Get(r, "user")
	data, err := handler.HandleCreate(ctx, w, r, user, nil)
	if err != nil {
		if _, ok := err.(views.ValidationError); ok {
			h.renderFormView(w, r, modelName, data, true)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adminutils.Success(ctx, r, fmt.Sprintf("%s was added successfully.", modelName))
	instance := adminutils.GetFieldValue(data, "Instance")
	id := adminutils.GetIDFromInstance(instance)
	if id > 0 {
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/%v/", modelName, id), http.StatusSeeOther)
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

func (h *CoreHandler) handleUpdateGet(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	user := h.sessionManager.Get(r, "user")
	data, err := handler.HandleDetail(ctx, w, r, user, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.renderFormView(w, r, modelName, data, false)
}

func (h *CoreHandler) handleUpdatePost(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	user := h.sessionManager.Get(r, "user")
	data, err := handler.HandleUpdate(ctx, w, r, user, id, nil)
	if err != nil {
		if _, ok := err.(views.ValidationError); ok {
			h.renderFormView(w, r, modelName, data, false)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adminutils.Success(ctx, r, fmt.Sprintf("%s was updated successfully.", modelName))
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

		if r.Method == http.MethodGet {
			h.handleDeleteGet(w, r, modelName, id)
		} else if r.Method == http.MethodPost {
			h.handleDeletePost(w, r, modelName, id)
		}
	}
}

func (h *CoreHandler) handleDeleteGet(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	templateData := map[string]interface{}{
		"Title":     "Are you sure?",
		"ModelName": modelName,
		"ObjectID":  id,
		"CSRFToken": csrf.TemplateField(r),
	}
	h.renderer.Render(w, "delete_confirmation", templateData)
}

func (h *CoreHandler) handleDeletePost(w http.ResponseWriter, r *http.Request, modelName string, id int64) {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	if err := handler.HandleDelete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	adminutils.Success(r.Context(), r, fmt.Sprintf("Successfully deleted %s.", modelName))
	http.Redirect(w, r, fmt.Sprintf("/admin/%s/", modelName), http.StatusSeeOther)
}
