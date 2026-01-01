package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	importpkg "github.com/forgego/forge/admin/import"
)

// HandleImport handles bulk import requests
func (h *CoreHandler) HandleImport(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
			return
		}

		// Get file
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
			return
		}
		defer file.Close()

	// Get options
	options := importpkg.DefaultImportOptions()
		if dryRun := r.FormValue("dry_run"); dryRun == "true" {
			options.DryRun = true
		}
		if validate := r.FormValue("validate"); validate == "false" {
			options.Validate = false
		}

		// Get column mapping
		mappingJSON := r.FormValue("column_mapping")
		if mappingJSON != "" {
			var mapping map[string]string
			if err := json.Unmarshal([]byte(mappingJSON), &mapping); err == nil {
				options.ColumnMapping = mapping
			}
		}

		// Get admin handler
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Import file
		ctx := r.Context()
		result, err := handler.ImportFile(ctx, file, header.Filename, options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// ImportFile imports a file (interface method)
type ImportFileHandler interface {
	ImportFile(ctx context.Context, file interface{}, filename string, options importpkg.ImportOptions) (*importpkg.ImportResult, error)
}
