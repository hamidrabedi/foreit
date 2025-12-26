package admin

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

// getStaticFS returns the embedded static filesystem
func getStaticFS() fs.FS {
	// Return a subdirectory of the embedded FS
	fsys, err := fs.Sub(staticFS, "static")
	if err != nil {
		// Fallback to root if sub fails
		return staticFS
	}
	return fsys
}

// serveStaticFile serves a file from the embedded filesystem
func serveStaticFile(w http.ResponseWriter, r *http.Request, path string) {
	fsys := getStaticFS()
	
	// Open the file from embedded filesystem
	file, err := fsys.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()
	
	// Get file info for content type
	info, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Check if it's a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}
	
	// Set appropriate content type based on file extension
	if len(path) >= 4 && path[len(path)-4:] == ".css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if len(path) >= 3 && path[len(path)-3:] == ".js" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	}
	
	// Read file content
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Set content length
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	
	// Write file content
	w.Write(data)
}

