package admin

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	httplib "github.com/forgego/forge/pkg/http"
)

// ExportFormat represents an export format
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
)

// handleExport handles export requests
func handleExport(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canView(r, modelName) {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}

		format := ExportFormat(strings.ToLower(httplib.GetQueryString(r, "format", "csv")))
		
		ctx := r.Context()
		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		ops := NewManagerOps()
		instances, err := ops.GetAllInstances(ctx, manager)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get instances: %v", err), http.StatusInternalServerError)
			return
		}

		// Get list display fields
		listDisplay := model.ListDisplay
		if len(listDisplay) == 0 {
			// Default: use all fields from first instance
			if len(instances) > 0 {
				instanceValue := reflect.ValueOf(instances[0])
				if instanceValue.Kind() == reflect.Ptr {
					instanceValue = instanceValue.Elem()
				}
				typ := instanceValue.Type()
				for i := 0; i < typ.NumField(); i++ {
					field := typ.Field(i)
					if field.IsExported() {
						listDisplay = append(listDisplay, field.Name)
					}
				}
			}
		}

		switch format {
		case ExportFormatCSV:
			exportCSV(w, instances, listDisplay, model)
		case ExportFormatJSON:
			exportJSON(w, instances, listDisplay, model)
		default:
			http.Error(w, fmt.Sprintf("Unsupported export format: %s", format), http.StatusBadRequest)
		}
	}
}

// exportCSV exports data as CSV
func exportCSV(w http.ResponseWriter, instances []interface{}, listDisplay []interface{}, model *AdminModel) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", model.Name))
	
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	headers := make([]string, 0, len(listDisplay))
	for _, field := range listDisplay {
		if fieldName, ok := field.(string); ok {
			headers = append(headers, fieldName)
		}
	}
	if err := writer.Write(headers); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write CSV header: %v", err), http.StatusInternalServerError)
		return
	}

	// Write rows
	for _, instance := range instances {
		row := make([]string, 0, len(listDisplay))
		instanceValue := reflect.ValueOf(instance)
		if instanceValue.Kind() == reflect.Ptr {
			instanceValue = instanceValue.Elem()
		}

		for _, field := range listDisplay {
			if fieldName, ok := field.(string); ok {
				fieldValue := instanceValue.FieldByName(fieldName)
				var value string
				if fieldValue.IsValid() && fieldValue.CanInterface() {
					value = fmt.Sprintf("%v", fieldValue.Interface())
				}
				row = append(row, value)
			}
		}
		
		if err := writer.Write(row); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write CSV row: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// exportJSON exports data as JSON
func exportJSON(w http.ResponseWriter, instances []interface{}, listDisplay []interface{}, model *AdminModel) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", model.Name))

	// Convert instances to maps
	results := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		instanceValue := reflect.ValueOf(instance)
		if instanceValue.Kind() == reflect.Ptr {
			instanceValue = instanceValue.Elem()
		}

		obj := make(map[string]interface{})
		for _, field := range listDisplay {
			if fieldName, ok := field.(string); ok {
				fieldValue := instanceValue.FieldByName(fieldName)
				if fieldValue.IsValid() && fieldValue.CanInterface() {
					obj[fieldName] = fieldValue.Interface()
				}
			}
		}
		results = append(results, obj)
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

