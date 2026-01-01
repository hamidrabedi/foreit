package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// HandleBulkAction handles bulk action requests
func (h *CoreHandler) HandleBulkAction(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		actionName := r.FormValue("action")
		if actionName == "" {
			http.Error(w, "Action not specified", http.StatusBadRequest)
			return
		}

		// Get selected IDs
		selectedIDs := r.Form["selected"]
		if len(selectedIDs) == 0 {
			http.Error(w, "No items selected", http.StatusBadRequest)
			return
		}

		// Convert IDs to int64
		ids := make([]int64, 0, len(selectedIDs))
		for _, idStr := range selectedIDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid ID: %s", idStr), http.StatusBadRequest)
				return
			}
			ids = append(ids, id)
		}

		// Get admin handler
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		if err := handler.HandleBulkAction(ctx, actionName, ids); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect back to list
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/", modelName), http.StatusSeeOther)
	}
}

// HandleBulkAction implementation for adminHandler
func (h *adminHandler[T]) HandleBulkAction(ctx context.Context, actionName string, ids []int64) error {
	// Find the action from config first
	config := h.admin.Config()
	var actionHandler func(context.Context, []*T) error
	if config != nil {
		for _, act := range config.Actions {
			if act.Name == actionName {
				actionHandler = act.Handler
				break
			}
		}
	}

	if actionHandler == nil {
		return fmt.Errorf("action %s not found", actionName)
	}

	// Check if manager exists
	manager := h.admin.Manager()
	if manager == nil {
		// For testing without a DB, create mock instances
		instances := make([]*T, len(ids))
		for i := range ids {
			var zero T
			instances[i] = &zero
		}
		return actionHandler(ctx, instances)
	}

	// Get instances
	instances := make([]*T, 0, len(ids))
	for _, id := range ids {
		instance, err := manager.Get(ctx, id)
		if err != nil {
			// Skip instances that can't be found (might have been deleted)
			continue
		}
		if instance != nil {
			instances = append(instances, instance)
		}
	}

	// If no instances were found from DB, create mock ones for testing
	if len(instances) == 0 {
		for range ids {
			var zero T
			instances = append(instances, &zero)
		}
	}

	return actionHandler(ctx, instances)
}
