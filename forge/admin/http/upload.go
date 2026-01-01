package http

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/forgego/forge/admin/storage"
)

// UploadHandler handles file uploads
type UploadHandler struct {
	storage storage.Storage
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage storage.Storage) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

// HandleUpload handles file upload requests
func (uh *UploadHandler) HandleUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form (max 32MB)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
			return
		}

		// Get file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get MIME type
		mimeType := header.Header.Get("Content-Type")
		if mimeType == "" {
			// Try to detect from extension
			ext := filepath.Ext(header.Filename)
			mimeType = getMimeTypeFromExtension(ext)
		}

		// Save file
		ctx := r.Context()
		fileInfo, err := uh.storage.Save(ctx, header.Filename, file, header.Size, mimeType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"success": true,
			"file": {
				"path": "%s",
				"url": "%s",
				"size": %d,
				"mimeType": "%s"
			}
		}`, fileInfo.Path, fileInfo.URL, fileInfo.Size, fileInfo.MimeType)
	}
}

// HandleDelete handles file deletion requests
func (uh *UploadHandler) HandleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete && r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get path from query or form
		path := r.URL.Query().Get("path")
		if path == "" {
			if err := r.ParseForm(); err == nil {
				path = r.FormValue("path")
			}
		}

		if path == "" {
			http.Error(w, "Path not specified", http.StatusBadRequest)
			return
		}

		// Delete file
		ctx := r.Context()
		if err := uh.storage.Delete(ctx, path); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete file: %v", err), http.StatusInternalServerError)
			return
		}

		// Return success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"success": true}`)
	}
}

// HandleGet handles file retrieval requests
func (uh *UploadHandler) HandleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			http.Error(w, "Path not specified", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		reader, err := uh.storage.Get(ctx, path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusNotFound)
			return
		}
		defer reader.Close()

		// Set content type
		ext := filepath.Ext(path)
		mimeType := getMimeTypeFromExtension(ext)
		w.Header().Set("Content-Type", mimeType)

		// Stream file
		io.Copy(w, reader)
	}
}

// getMimeTypeFromExtension returns MIME type from file extension
func getMimeTypeFromExtension(ext string) string {
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
	}

	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return "application/octet-stream"
}
