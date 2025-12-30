package http

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	httplib "github.com/forgego/forge/server"
)

// ExportHandler handles export functionality
type ExportHandler struct {
	handler *CoreHandler
}

// HandleExport handles export requests
func (h *CoreHandler) HandleExport(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		format := httplib.GetQueryString(r, "format", "csv")
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		data, err := handler.HandleExport(ctx, format)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch format {
		case "csv":
			h.exportCSV(w, data)
		case "json":
			h.exportJSON(w, data)
		case "xlsx":
			http.Error(w, "Excel export not yet implemented", http.StatusNotImplemented)
		default:
			http.Error(w, fmt.Sprintf("Unsupported format: %s", format), http.StatusBadRequest)
		}
	}
}

// exportCSV exports data as CSV
func (h *CoreHandler) exportCSV(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=export.csv")

	// Convert data to CSV format
	// This is a simplified version - would need proper type handling
	if exportData, ok := data.(map[string]interface{}); ok {
		if objects, ok := exportData["objects"].([]interface{}); ok {
			writer := csv.NewWriter(w)
			defer writer.Flush()

			// Write header
			if len(objects) > 0 {
				if firstObj, ok := objects[0].(map[string]interface{}); ok {
					headers := make([]string, 0, len(firstObj))
					for key := range firstObj {
						headers = append(headers, key)
					}
					writer.Write(headers)

					// Write rows
					for _, obj := range objects {
						if objMap, ok := obj.(map[string]interface{}); ok {
							row := make([]string, 0, len(headers))
							for _, header := range headers {
								val := ""
								if v, ok := objMap[header]; ok {
									val = fmt.Sprintf("%v", v)
								}
								row = append(row, val)
							}
							writer.Write(row)
						}
					}
				}
			}
		}
	}
}

// exportJSON exports data as JSON
func (h *CoreHandler) exportJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=export.json")
	json.NewEncoder(w).Encode(data)
}
