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
	// Get instances
	instances := make([]*T, 0, len(ids))
	for _, id := range ids {
		instance, err := h.admin.Manager().Get(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get instance %d: %w", id, err)
		}
		instances = append(instances, instance)
	}

	// Find and execute action from config
	config := h.admin.Config()
	if config != nil {
		for _, act := range config.Actions {
			if act.Name == actionName {
				return act.Handler(ctx, instances)
			}
		}
	}

	return fmt.Errorf("action %s not found", actionName)
}
