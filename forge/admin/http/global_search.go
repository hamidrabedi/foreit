package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/forgego/forge/admin"
)

// GlobalSearchResult represents a search result from a model
type GlobalSearchResult struct {
	ModelName  string
	ModelLabel string
	Results    []SearchResultItem
	TotalCount int
	HasMore    bool
	ViewURL    string
}

// SearchResultItem represents a single search result
type SearchResultItem struct {
	ID          interface{}
	Title       string
	Description string
	URL         string
	Fields      map[string]interface{}
}

// HandleGlobalSearch handles global search across all registered models
func (h *CoreHandler) HandleGlobalSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if query == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		limit := 10 // Results per model
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && l > 0 {
				// limit is set
			}
		}

		ctx := r.Context()
		results := h.searchAllModels(ctx, query, limit)

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"query":   query,
			"results": results,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// searchAllModels searches across all registered models
func (h *CoreHandler) searchAllModels(ctx context.Context, query string, limit int) []GlobalSearchResult {
	registry := admin.GetGlobalRegistry()
	results := make([]GlobalSearchResult, 0)

	// Get all registered admins
	allAdmins := registry.GetAll()

	for modelName, adminInstance := range allAdmins {
		// Get handler for this model
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			continue // Skip if we can't get handler
		}

		// Search in this model
		modelResults, err := h.searchModel(ctx, handler, modelName, query, limit)
		if err != nil {
			continue // Skip on error
		}

		if len(modelResults) > 0 {
			modelName := adminInstance.ModelName()
			results = append(results, GlobalSearchResult{
				ModelName:  modelName,
				ModelLabel: modelName, // Would get from config
				Results:    modelResults,
				TotalCount: len(modelResults),
				HasMore:    len(modelResults) >= limit,
				ViewURL:    fmt.Sprintf("/admin/%s/?search=%s", modelName, query),
			})
		}
	}

	return results
}

// searchModel performs a search on a specific admin model
func (h *CoreHandler) searchModel(ctx context.Context, handler AdminHandler, modelName string, query string, limit int) ([]SearchResultItem, error) {
	// Use type assertion to get the admin instance
	// This is a simplified approach - in production, you'd use a more type-safe method
	adminInstance := h.getAdminFromHandler(handler)
	if adminInstance == nil {
		return []SearchResultItem{}, nil
	}

	// Get config via reflection
	configVal := reflect.ValueOf(adminInstance).MethodByName("Config").Call(nil)[0]
	if configVal.IsNil() {
		return []SearchResultItem{}, nil
	}

	config := configVal.Interface()
	searchFieldsVal := reflect.ValueOf(config).Elem().FieldByName("SearchFields")
	if !searchFieldsVal.IsValid() || searchFieldsVal.Len() == 0 {
		return []SearchResultItem{}, nil
	}

	// For now, return empty results
	// Full implementation would:
	// 1. Get queryset from admin
	// 2. Apply search using ListView's applySearch
	// 3. Limit and get objects
	// 4. Convert to search results

	return []SearchResultItem{}, nil
}

// getAdminFromHandler extracts admin instance from handler
func (h *CoreHandler) getAdminFromHandler(handler AdminHandler) interface{} {
	// Use reflection to get admin field
	val := reflect.ValueOf(handler)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	adminField := val.FieldByName("admin")
	if adminField.IsValid() {
		return adminField.Interface()
	}
	return nil
}
