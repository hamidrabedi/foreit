package admin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// BulkAction represents a bulk action that can be performed on multiple objects
type BulkAction struct {
	Name        string
	Label       string
	Description string
	Handler     func(r *http.Request, ids []int64) error
}

// handleBulkAction handles bulk actions on selected objects
func handleBulkAction(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canChange(r, modelName) && !canDelete(r, modelName) {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get selected object IDs
		selectedIDs := r.Form["selected_ids"]
		if len(selectedIDs) == 0 {
			http.Error(w, "No objects selected", http.StatusBadRequest)
			return
		}

		// Parse IDs
		ids := make([]int64, 0, len(selectedIDs))
		for _, idStr := range selectedIDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid ID: %s", idStr), http.StatusBadRequest)
				return
			}
			ids = append(ids, id)
		}

		// Get action
		action := r.FormValue("action")
		if action == "" {
			http.Error(w, "No action specified", http.StatusBadRequest)
			return
		}

		// Execute action
		ctx := r.Context()
		var err error

		switch action {
		case "delete":
			if !canDelete(r, modelName) {
				http.Error(w, "Permission denied: You do not have permission to delete objects.", http.StatusForbidden)
				return
			}
			err = bulkDelete(ctx, model, ids)
		case "activate":
			if !canChange(r, modelName) {
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}
			err = bulkActivate(ctx, model, ids)
		case "deactivate":
			if !canChange(r, modelName) {
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}
			err = bulkDeactivate(ctx, model, ids)
		default:
			http.Error(w, fmt.Sprintf("Unknown action: %s", action), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute action: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect back to list
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/", modelName), http.StatusSeeOther)
	}
}

// bulkDelete deletes multiple objects
func bulkDelete(ctx interface{}, model *AdminModel, ids []int64) error {
	manager, err := getManagerForModel(model)
	if err != nil {
		return fmt.Errorf("failed to get manager: %w", err)
	}

	ops := NewManagerOps()
	ctxTyped, ok := ctx.(context.Context)
	if !ok {
		return fmt.Errorf("invalid context type")
	}

	for _, id := range ids {
		// Get object first
		instance, err := ops.GetInstance(ctxTyped, manager, id)
		if err != nil {
			continue // Skip if not found
		}

		// Delete object
		if err := ops.DeleteInstance(ctxTyped, manager, instance); err != nil {
			return fmt.Errorf("failed to delete object %d: %w", id, err)
		}
	}

	return nil
}

// bulkActivate activates multiple objects (if they have IsActive field)
func bulkActivate(ctx interface{}, model *AdminModel, ids []int64) error {
	return bulkUpdateField(ctx, model, ids, "is_active", true)
}

// bulkDeactivate deactivates multiple objects (if they have IsActive field)
func bulkDeactivate(ctx interface{}, model *AdminModel, ids []int64) error {
	return bulkUpdateField(ctx, model, ids, "is_active", false)
}

// bulkUpdateField updates a field for multiple objects
func bulkUpdateField(ctx interface{}, model *AdminModel, ids []int64, fieldName string, value interface{}) error {
	manager, err := getManagerForModel(model)
	if err != nil {
		return fmt.Errorf("failed to get manager: %w", err)
	}

	managerValue := reflect.ValueOf(manager)
	getMethod := managerValue.MethodByName("Get")
	updateMethod := managerValue.MethodByName("Update")

	for _, id := range ids {
		// Get object
		if !getMethod.IsValid() {
			return fmt.Errorf("manager does not have Get method")
		}

		results := getMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(id),
		})

		if len(results) < 2 || !results[1].IsNil() {
			continue // Skip if not found
		}

		instanceValue := reflect.ValueOf(results[0].Interface())
		if instanceValue.Kind() == reflect.Ptr {
			instanceValue = instanceValue.Elem()
		}

		// Update field
		field := instanceValue.FieldByName(toTitleCase(fieldName))
		if !field.IsValid() {
			// Try lowercase
			field = instanceValue.FieldByName(fieldName)
		}

		if field.IsValid() && field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}

		// Save object
		if !updateMethod.IsValid() {
			return fmt.Errorf("manager does not have Update method")
		}

		updateResults := updateMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(instanceValue.Addr().Interface()),
		})

		if len(updateResults) > 0 && !updateResults[0].IsNil() {
			if err, ok := updateResults[0].Interface().(error); ok {
				return fmt.Errorf("failed to update object %d: %w", id, err)
			}
		}
	}

	return nil
}

// toTitleCase converts snake_case or lowercase to TitleCase
func toTitleCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "")
}

