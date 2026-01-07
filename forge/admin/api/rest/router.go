package rest

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/forgego/forge/admin/core"
	"github.com/forgego/forge/media"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Router handles all admin API routes
type Router struct {
	registry    *core.Registry
	prefix      string
	mediaEngine *media.Engine
}

// NewRouter creates a new admin API router
func NewRouter(registry *core.Registry) *Router {
	return &Router{
		registry: registry,
		prefix:   "/api",
	}
}

// SetMediaEngine configures upload/media handling.
func (r *Router) SetMediaEngine(engine *media.Engine) {
	r.mediaEngine = engine
}

// RegisterRoutes registers all admin routes
func (r *Router) RegisterRoutes(router chi.Router) {
	router.Route(r.prefix, func(sub chi.Router) {
		// Middleware
		sub.Use(middleware.Logger)
		sub.Use(middleware.Recoverer)
		sub.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		// Metadata endpoints
		sub.Get("/meta", r.handleMetaList)
		sub.Get("/meta/{model}", r.handleMetaDetail)

		// Auth endpoints
		sub.Post("/login", r.handleLogin)

		// Global search
		sub.Get("/search", r.handleGlobalSearch)

		// Model routes (registered dynamically)
		r.registerModelRoutes(sub)
	})
}

// registerModelRoutes registers routes for each model
func (r *Router) registerModelRoutes(router chi.Router) {
	for name, admin := range r.registry.GetAll() {
		basePath := fmt.Sprintf("/%s", name)

		router.Route(basePath, func(sub chi.Router) {
			// List and create
			sub.Get("/", r.handleList(admin))
			sub.Post("/", r.handleCreate(admin))

			// Detail, update, delete
			sub.Route("/{id}", func(subDetail chi.Router) {
				subDetail.Get("/", r.handleDetail(admin))
				subDetail.Patch("/", r.handleUpdate(admin))
				subDetail.Put("/", r.handleReplace(admin))
				subDetail.Get("/history", r.handleHistory(admin))
				subDetail.Delete("/", r.handleDelete(admin))
			})

			// Bulk operations
			sub.Post("/bulk-create", r.handleBulkCreate(admin))
			sub.Post("/bulk-update", r.handleBulkUpdate(admin))
			sub.Delete("/bulk-delete", r.handleBulkDelete(admin))

			// Actions
			sub.Post("/action/{action}", r.handleAction(admin))

			// Autocomplete
			sub.Get("/autocomplete", r.handleAutocomplete(admin))

			// Upload
			sub.Post("/upload", r.handleUpload(admin))

			// Export
			sub.Get("/export", r.handleExport(admin))
		})
	}
}

// handleMetaList returns list of all models
func (r *Router) handleMetaList(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user := getUserFromContext(ctx)

	allAdmins := r.registry.GetAll()
	models := make([]core.ModelListMetadata, 0, len(allAdmins))

	for name, admin := range allAdmins {
		// Check module permission
		if !admin.HasModulePermission(ctx, user) {
			continue
		}

		meta, err := admin.GetMetadata(ctx, user)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
			return
		}

		count := int64(0)
		if resp, err := admin.ListObjects(ctx, core.ListParams{Page: 1, PageSize: 1}); err == nil {
			count = resp.Count
		}

		models = append(models, core.ModelListMetadata{
			Name:              name,
			VerboseName:       meta.VerboseName,
			VerboseNamePlural: meta.VerboseNamePlural,
			Icon:              meta.Icon,
			Count:             count,
			Permissions:       meta.Permissions,
		})
	}

	allPlugins := r.registry.GetAllPlugins()
	plugins := make([]core.PluginMetadata, 0, len(allPlugins))
	for _, p := range allPlugins {
		plugins = append(plugins, p.GetMetadata())
	}

	dashboard := r.registry.GetDashboardConfig()
	dashboard.Widgets = append(dashboard.Widgets, collectPluginWidgets(allPlugins)...)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"models":       models,
		"plugins":      plugins,
		"custom_pages": r.registry.GetCustomPages(),
		"menu_entries": r.registry.GetMenuEntries(),
		"dashboard":    dashboard,
	})
}

// handleLogin handles admin login
func (r *Router) handleLogin(w http.ResponseWriter, req *http.Request) {
	// Simple stub for now - just return success
	// Production auth should verify credentials
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": "dummy-token",
		"user":  map[string]string{"name": "Admin", "role": "superuser"},
	})
}

// handleMetaDetail returns detailed metadata for a model
func (r *Router) handleMetaDetail(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	modelName := chi.URLParam(req, "model")
	user := getUserFromContext(ctx)

	admin, err := r.registry.Get(modelName)
	if err != nil {
		respondError(w, http.StatusNotFound, "model_not_found", err.Error(), nil)
		return
	}

	// Check module permission
	if !admin.HasModulePermission(ctx, user) {
		respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to access this model", nil)
		return
	}

	meta, err := admin.GetMetadata(ctx, user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal_error", err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, meta)
}

// handleList returns a paginated list of objects
func (r *Router) handleList(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		// Check view permission
		if !admin.HasViewPermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to view this model", nil)
			return
		}

		// Parse query parameters
		params := parseListParams(req)

		// Call implementation
		response, err := admin.ListObjects(ctx, params)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "list_failed", err.Error(), nil)
			return
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// handleDetail returns a single object
func (r *Router) handleDetail(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)
		idStr := chi.URLParam(req, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id", "Invalid ID format", nil)
			return
		}

		// Check view permission
		if !admin.HasViewPermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to view this object", nil)
			return
		}

		// Call implementation
		obj, err := admin.GetObject(ctx, id)
		if err != nil {
			respondError(w, http.StatusNotFound, "not_found", "Object not found", nil)
			return
		}

		respondJSON(w, http.StatusOK, obj)
	}
}

// handleCreate creates a new object
func (r *Router) handleCreate(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		// Check add permission
		if !admin.HasAddPermission(ctx, user) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to add this model", nil)
			return
		}

		// Parse body
		var data map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		// Call implementation
		obj, err := admin.CreateObject(ctx, data)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "create_failed", err.Error(), nil)
			return
		}

		respondJSON(w, http.StatusCreated, obj)
	}
}

// handleUpdate partially updates an object
func (r *Router) handleUpdate(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)
		idStr := chi.URLParam(req, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id", "Invalid ID format", nil)
			return
		}

		// Check change permission
		if !admin.HasChangePermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to change this object", nil)
			return
		}

		// Parse body
		var data map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		// Call implementation
		obj, err := admin.UpdateObject(ctx, id, data)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "update_failed", err.Error(), nil)
			return
		}

		respondJSON(w, http.StatusOK, obj)
	}
}

// handleReplace fully replaces an object
func (r *Router) handleReplace(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)
		idStr := chi.URLParam(req, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id", "Invalid ID format", nil)
			return
		}

		// Check change permission
		if !admin.HasChangePermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to change this object", nil)
			return
		}

		var data map[string]interface{}
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		obj, err := admin.UpdateObject(ctx, id, data)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "update_failed", err.Error(), nil)
			return
		}

		respondJSON(w, http.StatusOK, obj)
	}
}

// handleHistory returns history entries for a specific object.
func (r *Router) handleHistory(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)
		idStr := chi.URLParam(req, "id")

		if idStr == "" {
			respondError(w, http.StatusBadRequest, "invalid_id", "Invalid ID format", nil)
			return
		}

		if !admin.HasViewPermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to view this object", nil)
			return
		}

		history, err := admin.GetHistory(ctx, idStr)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "history_failed", err.Error(), nil)
			return
		}

		respondJSON(w, http.StatusOK, history)
	}
}

// handleDelete deletes an object
func (r *Router) handleDelete(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)
		idStr := chi.URLParam(req, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_id", "Invalid ID format", nil)
			return
		}

		// Check delete permission
		if !admin.HasDeletePermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to delete this object", nil)
			return
		}

		// Call implementation
		err = admin.DeleteObject(ctx, id)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "delete_failed", err.Error(), nil)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handleBulkCreate creates multiple objects
func (r *Router) handleBulkCreate(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		// Check add permission
		if !admin.HasAddPermission(ctx, user) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to add this model", nil)
			return
		}

		var payload interface{}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		records, err := coerceBulkCreatePayload(payload)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		created := 0
		objects := make([]interface{}, 0, len(records))
		errors := []core.BulkActionError{}

		for idx, data := range records {
			obj, err := admin.CreateObject(ctx, data)
			if err != nil {
				errors = append(errors, core.BulkActionError{
					ID:      int64(idx),
					Message: err.Error(),
				})
				continue
			}
			created++
			objects = append(objects, obj)
		}

		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"created": created,
			"objects": objects,
			"errors":  errors,
		})
	}
}

// handleBulkUpdate updates multiple objects
func (r *Router) handleBulkUpdate(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		// Check change permission
		if !admin.HasChangePermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to change this model", nil)
			return
		}

		var payload struct {
			IDs  []interface{}          `json:"ids"`
			Data map[string]interface{} `json:"data"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		updated := 0
		errors := []core.BulkActionError{}
		for _, rawID := range payload.IDs {
			obj, err := admin.UpdateObject(ctx, rawID, payload.Data)
			if err != nil {
				errors = append(errors, core.BulkActionError{
					ID:      toInt64(rawID),
					Message: err.Error(),
				})
				continue
			}
			_ = obj
			updated++
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"updated": updated,
			"errors":  errors,
		})
	}
}

// handleBulkDelete deletes multiple objects
func (r *Router) handleBulkDelete(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		// Check delete permission
		if !admin.HasDeletePermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to delete this model", nil)
			return
		}

		var payload struct {
			IDs []interface{} `json:"ids"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		deleted := 0
		errors := []core.BulkActionError{}
		for _, rawID := range payload.IDs {
			if err := admin.DeleteObject(ctx, rawID); err != nil {
				errors = append(errors, core.BulkActionError{
					ID:      toInt64(rawID),
					Message: err.Error(),
				})
				continue
			}
			deleted++
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"deleted": deleted,
			"errors":  errors,
		})
	}
}

// handleAction executes a bulk action
func (r *Router) handleAction(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)
		actionName := chi.URLParam(req, "action")

		// Check change permission (actions typically require change permission)
		if !admin.HasChangePermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to execute actions", nil)
			return
		}

		// Parse request
		var request core.BulkActionRequest
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			respondError(w, http.StatusBadRequest, "invalid_body", err.Error(), nil)
			return
		}

		// Call implementation
		response, err := admin.ExecuteAction(ctx, actionName, request.IDs, request.Params)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "action_failed", err.Error(), nil)
			return
		}

		respondJSON(w, http.StatusOK, response)
	}
}

// handleAutocomplete returns autocomplete suggestions
func (r *Router) handleAutocomplete(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		// Check view permission
		if !admin.HasViewPermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to view this model", nil)
			return
		}

		field := req.URL.Query().Get("field")
		query := req.URL.Query().Get("q")
		limit := 10

		if limitStr := req.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		results, err := admin.Autocomplete(ctx, query, limit)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "autocomplete_failed", err.Error(), nil)
			return
		}

		_ = field // Currently not using field-specific autocomplete, but available for future

		respondJSON(w, http.StatusOK, core.AutocompleteResponse{
			Results: results,
		})
	}
}

// handleUpload accepts a single file upload for a model.
func (r *Router) handleUpload(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		if !admin.HasChangePermission(ctx, user, nil) && !admin.HasAddPermission(ctx, user) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to upload files for this model", nil)
			return
		}

		if r.mediaEngine == nil {
			respondError(w, http.StatusBadRequest, "upload_disabled", "File uploads are not configured", nil)
			return
		}

		result, err := r.mediaEngine.SaveUploadFromRequest(req, "file", admin.ModelName())
		if err != nil {
			respondError(w, http.StatusBadRequest, "upload_failed", err.Error(), nil)
			return
		}

		resp := core.UploadResponse{
			URL:        result.URL,
			Filename:   result.Filename,
			Size:       result.Size,
			MimeType:   result.MimeType,
			UploadedAt: result.UploadedAt,
		}
		respondJSON(w, http.StatusOK, resp)
	}
}

// handleGlobalSearch handles global search across all models
func (r *Router) handleGlobalSearch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	query := req.URL.Query().Get("q")

	if query == "" {
		respondError(w, http.StatusBadRequest, "missing_query", "Search query is required", nil)
		return
	}

	allAdmins := r.registry.GetAll()
	results := make([]core.SearchResultGroup, 0)

	for name, admin := range allAdmins {
		items, err := admin.Autocomplete(ctx, query, 5)
		if err != nil || len(items) == 0 {
			continue
		}

		group := core.SearchResultGroup{
			Model: name,
			Count: len(items),
			Items: make([]core.SearchResultItem, 0, len(items)),
		}

		for _, item := range items {
			group.Items = append(group.Items, core.SearchResultItem{
				ID:    item.Value,
				Title: item.Label,
				URL:   fmt.Sprintf("/%s/%v", name, item.Value),
			})
		}

		results = append(results, group)
	}

	respondJSON(w, http.StatusOK, core.SearchResponse{
		Results: results,
	})
}

// handleExport exports data for a model.
func (r *Router) handleExport(admin core.AdminInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := getUserFromContext(ctx)

		if !admin.HasViewPermission(ctx, user, nil) {
			respondError(w, http.StatusForbidden, "permission_denied", "You don't have permission to export this model", nil)
			return
		}

		format := strings.ToLower(req.URL.Query().Get("format"))
		if format == "" {
			format = "json"
		}

		meta, err := admin.GetMetadata(ctx, user)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "metadata_failed", err.Error(), nil)
			return
		}

		results, err := fetchAllForExport(ctx, admin)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "export_failed", err.Error(), nil)
			return
		}

		switch format {
		case "csv":
			if err := writeCSVExport(w, meta.ListDisplay, results); err != nil {
				respondError(w, http.StatusInternalServerError, "export_failed", err.Error(), nil)
			}
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(results)
		}
	}
}

func collectPluginWidgets(plugins map[string]core.Plugin) []core.WidgetMeta {
	result := []core.WidgetMeta{}
	for _, p := range plugins {
		meta := p.GetMetadata()
		if len(meta.Widgets) > 0 {
			result = append(result, meta.Widgets...)
		}
	}
	return result
}

func fetchAllForExport(ctx context.Context, admin core.AdminInterface) ([]map[string]interface{}, error) {
	page := 1
	pageSize := 500
	allResults := []map[string]interface{}{}

	for {
		resp, err := admin.ListObjects(ctx, core.ListParams{
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			return nil, err
		}

		batch, err := coerceResults(resp.Results)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, batch...)

		if page >= resp.TotalPages || len(batch) == 0 {
			break
		}
		page++
	}

	return allResults, nil
}

func coerceResults(results interface{}) ([]map[string]interface{}, error) {
	raw, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	var decoded []map[string]interface{}
	if err := json.Unmarshal(raw, &decoded); err == nil {
		return decoded, nil
	}

	var single map[string]interface{}
	if err := json.Unmarshal(raw, &single); err == nil && len(single) > 0 {
		return []map[string]interface{}{single}, nil
	}

	return []map[string]interface{}{}, nil
}

func writeCSVExport(w http.ResponseWriter, columns []string, rows []map[string]interface{}) error {
	if len(columns) == 0 && len(rows) > 0 {
		for key := range rows[0] {
			columns = append(columns, key)
		}
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer := csv.NewWriter(w)
	if err := writer.Write(columns); err != nil {
		return err
	}

	for _, row := range rows {
		record := make([]string, len(columns))
		for i, col := range columns {
			if val, ok := row[col]; ok && val != nil {
				record[i] = fmt.Sprintf("%v", val)
			} else {
				record[i] = ""
			}
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

// Helper types and functions

func parseListParams(req *http.Request) core.ListParams {
	params := core.ListParams{
		Page:     1,
		PageSize: 25,
		Filters:  make(map[string]interface{}),
	}

	if pageStr := req.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if pageSizeStr := req.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			params.PageSize = pageSize
		}
	}

	params.Search = req.URL.Query().Get("search")

	if ordering := req.URL.Query().Get("ordering"); ordering != "" {
		params.Ordering = strings.Split(ordering, ",")
	}

	// Parse filters from query params
	for key, values := range req.URL.Query() {
		if key == "page" || key == "page_size" || key == "search" || key == "ordering" {
			continue
		}
		if len(values) > 0 {
			params.Filters[key] = values[0]
		}
	}

	return params
}

func coerceBulkCreatePayload(payload interface{}) ([]map[string]interface{}, error) {
	switch typed := payload.(type) {
	case []interface{}:
		records := make([]map[string]interface{}, 0, len(typed))
		for _, item := range typed {
			mapped, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("bulk create expects array of objects")
			}
			records = append(records, mapped)
		}
		return records, nil
	case map[string]interface{}:
		if raw, ok := typed["objects"]; ok {
			return coerceBulkCreatePayload(raw)
		}
		return nil, fmt.Errorf("bulk create expects array of objects")
	default:
		return nil, fmt.Errorf("bulk create expects array of objects")
	}
}

func toInt64(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}

type contextKey string

const userContextKey contextKey = "user"

func getUserFromContext(ctx context.Context) interface{} {
	return ctx.Value(userContextKey)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, code string, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(core.ErrorResponse{
		Error: core.ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
