package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandleListEditable handles saving inline edits from list view
// Uses type registry for type-safe operations
func (h *CoreHandler) HandleListEditable(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var request struct {
			Edits []struct {
				ObjectID  int64       `json:"object_id"`
				FieldName string      `json:"field_name"`
				Value     interface{} `json:"value"`
			} `json:"edits"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get type-safe admin handler using type registry
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		var user interface{}
		// Try to get user from session, but allow nil for testing
		defer func() {
			if r := recover(); r != nil {
				// Session not initialized, use nil user
				user = nil
			}
		}()
		user = h.sessionManager.Get(r, "user")

		// Process each edit
		for _, edit := range request.Edits {
			// Get the instance
			instanceData, err := handler.HandleDetail(ctx, w, r, user, edit.ObjectID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get instance %d: %v", edit.ObjectID, err), http.StatusBadRequest)
				return
			}

			// Extract instance from data
			instanceMap, ok := instanceData.(map[string]interface{})
			if !ok {
				http.Error(w, "Invalid instance data", http.StatusInternalServerError)
				return
			}

			instance := instanceMap["instance"]
			if instance == nil {
				http.Error(w, fmt.Sprintf("Instance %d not found", edit.ObjectID), http.StatusNotFound)
				return
			}

			// Create form data with the edit
			formData := map[string]interface{}{
				edit.FieldName: edit.Value,
			}

			// Update the instance
			_, err = handler.HandleUpdate(ctx, w, r, user, edit.ObjectID, formData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to update instance %d: %v", edit.ObjectID, err), http.StatusBadRequest)
				return
			}
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Edits saved successfully",
		})
	}
}
