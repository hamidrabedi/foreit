package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	adminutils "github.com/forgego/forge/admin/utils"
	"github.com/forgego/forge/orm"
	httplib "github.com/forgego/forge/server"
)

// HandleAutocomplete handles autocomplete requests for foreign key fields
func (h *CoreHandler) HandleAutocomplete(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		field := httplib.GetQueryString(r, "field", "")
		search := httplib.GetQueryString(r, "search", "")
		limitStr := httplib.GetQueryString(r, "limit", "10")

		if field == "" {
			http.Error(w, "Field parameter required", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}

		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		results, err := handler.HandleAutocomplete(ctx, search, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

// HandleAutocomplete implementation for adminHandler
func (h *adminHandler[T]) HandleAutocomplete(ctx context.Context, search string, limit int) ([]map[string]interface{}, error) {
	// Get base queryset
	qs, err := h.admin.GetQueryset(ctx)
	if err != nil {
		return nil, err
	}

	// Apply search if provided
	if search != "" {
		// Get field accessor for type-safe field expressions
		fa, err := h.admin.Manager().GetFieldAccessor()
		if err != nil {
			return nil, fmt.Errorf("failed to get field accessor: %w", err)
		}

		// Get model schema to find string fields for search
		schema, err := orm.GetModelSchema[T]()
		if err != nil {
			return nil, fmt.Errorf("failed to get model schema: %w", err)
		}

		// Try common searchable fields (name, title, username, email)
		searchFields := []string{"name", "title", "username", "email"}
		var searchExpr orm.Expression

		for _, fieldName := range searchFields {
			fieldInfo := schema.GetField(fieldName)
			if fieldInfo != nil && fieldInfo.Type.Kind() == reflect.String {
				field := orm.FieldFor[T, string](fa, fieldName)
				containsExpr := field.Contains(search)
				if searchExpr == nil {
					searchExpr = containsExpr
				} else {
					// Combine with OR
					q := orm.NewQ(searchExpr)
					searchExpr = q.Or(orm.NewQ(containsExpr))
				}
			}
		}

		if searchExpr != nil {
			qs = qs.Filter(searchExpr)
		}
	}

	// Limit results
	qs = qs.Limit(limit)

	// Get instances
	instances, err := qs.All(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to autocomplete format
	results := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		// Get ID and display value
		result := map[string]interface{}{
			"id":    adminutils.GetIDFromInstance(instance),
			"label": fmt.Sprintf("%v", instance),
		}
		results = append(results, result)
	}

	return results, nil
}
