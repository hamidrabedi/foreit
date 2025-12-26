package admin

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	httplib "github.com/forgego/forge/pkg/http"
)

// handleInlineFieldEdit handles GET request for inline field editing
func handleInlineFieldEdit(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get ID and field name from URL
		idStr := httplib.GetParam(r, "id")
		fieldName := httplib.GetParam(r, "field")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Get instance
		modelType := reflect.TypeOf(model.Model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		instanceValue := reflect.New(modelType)
		instance := instanceValue.Interface()

		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		// Get instance
		managerValue := reflect.ValueOf(manager)
		getMethod := managerValue.MethodByName("Get")
		if !getMethod.IsValid() {
			http.Error(w, "Manager does not have Get method", http.StatusInternalServerError)
			return
		}

		results := getMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(id),
		})

		if len(results) < 2 || !results[1].IsNil() {
			if err, ok := results[1].Interface().(error); ok && err != nil {
				http.Error(w, fmt.Sprintf("Failed to get instance: %v", err), http.StatusNotFound)
				return
			}
		}

		instance = results[0].Interface()

		// Get field value
		instanceValue = reflect.ValueOf(instance)
		if instanceValue.Kind() == reflect.Ptr {
			instanceValue = instanceValue.Elem()
		}

		fieldValue := instanceValue.FieldByName(fieldName)
		if !fieldValue.IsValid() {
			http.Error(w, "Invalid field", http.StatusBadRequest)
			return
		}

		// Render inline edit input
		fieldType := fieldValue.Type()
		var inputHTML string

		switch fieldType.Kind() {
		case reflect.String:
			currentValue := fieldValue.String()
			inputHTML = fmt.Sprintf(`<input type="text" 
				class="form-control form-control-sm" 
				value="%s"
				hx-post="/admin/%s/%d/edit-field/%s/"
				hx-trigger="blur, keyup[key=='Enter']"
				hx-target="closest td"
				hx-swap="outerHTML"
				name="value">`, currentValue, modelName, id, fieldName)
		case reflect.Int, reflect.Int64, reflect.Int32:
			currentValue := fieldValue.Int()
			inputHTML = fmt.Sprintf(`<input type="number" 
				class="form-control form-control-sm" 
				value="%d"
				hx-post="/admin/%s/%d/edit-field/%s/"
				hx-trigger="blur, keyup[key=='Enter']"
				hx-target="closest td"
				hx-swap="outerHTML"
				name="value">`, currentValue, modelName, id, fieldName)
		case reflect.Bool:
			currentValue := fieldValue.Bool()
			checked := ""
			if currentValue {
				checked = "checked"
			}
			inputHTML = fmt.Sprintf(`<input type="checkbox" 
				class="form-check-input" 
				%s
				hx-post="/admin/%s/%d/edit-field/%s/"
				hx-trigger="change"
				hx-target="closest td"
				hx-swap="outerHTML"
				name="value">`, checked, modelName, id, fieldName)
		default:
			// Fallback to text input
			currentValue := fmt.Sprintf("%v", fieldValue.Interface())
			inputHTML = fmt.Sprintf(`<input type="text" 
				class="form-control form-control-sm" 
				value="%s"
				hx-post="/admin/%s/%d/edit-field/%s/"
				hx-trigger="blur, keyup[key=='Enter']"
				hx-target="closest td"
				hx-swap="outerHTML"
				name="value">`, currentValue, modelName, id, fieldName)
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<td class="editable-cell" style="cursor: pointer;" title="Double-click to edit">%s</td>`, inputHTML)
	}
}

// handleInlineFieldEditPost handles POST request for inline field editing
func handleInlineFieldEditPost(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get ID and field name from URL
		idStr := httplib.GetParam(r, "id")
		fieldName := httplib.GetParam(r, "field")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		valueStr := r.Form.Get("value")

		// Get instance
		modelType := reflect.TypeOf(model.Model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		instanceValue := reflect.New(modelType)
		instance := instanceValue.Interface()

		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		// Get instance
		managerValue := reflect.ValueOf(manager)
		getMethod := managerValue.MethodByName("Get")
		if !getMethod.IsValid() {
			http.Error(w, "Manager does not have Get method", http.StatusInternalServerError)
			return
		}

		results := getMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(id),
		})

		if len(results) < 2 || !results[1].IsNil() {
			if err, ok := results[1].Interface().(error); ok && err != nil {
				http.Error(w, fmt.Sprintf("Failed to get instance: %v", err), http.StatusNotFound)
				return
			}
		}

		instance = results[0].Interface()

		// Update field value
		instanceValue = reflect.ValueOf(instance)
		if instanceValue.Kind() == reflect.Ptr {
			instanceValue = instanceValue.Elem()
		}

		fieldValue := instanceValue.FieldByName(fieldName)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			http.Error(w, "Invalid or readonly field", http.StatusBadRequest)
			return
		}

		// Set field value based on type
		fieldType := fieldValue.Type()
		switch fieldType.Kind() {
		case reflect.String:
			fieldValue.SetString(valueStr)
		case reflect.Int, reflect.Int64, reflect.Int32:
			if intVal, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
				fieldValue.SetInt(intVal)
			}
		case reflect.Bool:
			fieldValue.SetBool(valueStr == "on" || valueStr == "true")
		default:
			// Try to set as string representation
			// This is a simplified approach
		}

		// Update instance using manager
		updateMethod := managerValue.MethodByName("Update")
		if !updateMethod.IsValid() {
			http.Error(w, "Manager does not have Update method", http.StatusInternalServerError)
			return
		}

		updateResults := updateMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(instance),
		})

		if len(updateResults) > 0 && !updateResults[0].IsNil() {
			if err, ok := updateResults[0].Interface().(error); ok && err != nil {
				http.Error(w, fmt.Sprintf("Failed to update: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Return updated cell
		displayValue := fmt.Sprintf("%v", fieldValue.Interface())
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<td class="editable-cell" 
			hx-get="/admin/%s/%d/edit-field/%s/"
			hx-trigger="dblclick"
			hx-target="this"
			hx-swap="outerHTML"
			style="cursor: pointer;"
			title="Double-click to edit">%s</td>`, modelName, id, fieldName, displayValue)
	}
}

