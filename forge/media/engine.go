package media

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMaxUploadSize = 10 << 20 // 10MB
)

// Config defines directories and URL prefixes for static and upload handling.
type Config struct {
	StaticDir     string
	StaticURL     string
	UploadDir     string
	UploadURL     string
	MaxUploadSize int64
}

// UploadResult is returned after a successful upload.
type UploadResult struct {
	URL        string    `json:"url"`
	Path       string    `json:"path"`
	Filename   string    `json:"filename"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Engine provides static file handlers and upload processing.
type Engine struct {
	cfg Config
}

// New creates a new media engine with defaults applied.
func New(cfg Config) *Engine {
	if cfg.MaxUploadSize == 0 {
		cfg.MaxUploadSize = defaultMaxUploadSize
	}
	return &Engine{cfg: cfg}
}

// StaticHandler returns a handler that serves static files if configured.
func (e *Engine) StaticHandler() http.Handler {
	if e.cfg.StaticDir == "" {
		return nil
	}
	prefix := e.StaticURL()
	handler := http.FileServer(http.Dir(e.cfg.StaticDir))
	if prefix == "" {
		return handler
	}
	return http.StripPrefix(prefix, handler)
}

// MediaHandler serves uploaded files from UploadDir if configured.
func (e *Engine) MediaHandler() http.Handler {
	if e.cfg.UploadDir == "" {
		return nil
	}
	prefix := normalizePath(e.cfg.UploadURL)
	if prefix == "" {
		prefix = "/media"
	}
	handler := http.FileServer(http.Dir(e.cfg.UploadDir))
	return http.StripPrefix(prefix, handler)
}

// UploadURL returns the normalized upload URL prefix.
func (e *Engine) UploadURL() string {
	prefix := normalizePath(e.cfg.UploadURL)
	if prefix == "" {
		return "/media"
	}
	return prefix
}

// StaticURL returns the normalized static URL prefix.
func (e *Engine) StaticURL() string {
	if e.cfg.StaticDir == "" {
		return ""
	}
	prefix := normalizePath(e.cfg.StaticURL)
	if prefix == "" {
		return "/static"
	}
	return prefix
}

// NormalizeStaticURL is an alias for StaticURL.
func (e *Engine) NormalizeStaticURL() string {
	return e.StaticURL()
}

// SaveUploadFromRequest parses a multipart request and stores the file.
func (e *Engine) SaveUploadFromRequest(req *http.Request, formField, subdir string) (*UploadResult, error) {
	if e.cfg.UploadDir == "" {
		return nil, errors.New("uploads are not configured")
	}
	if formField == "" {
		formField = "file"
	}

	if err := req.ParseMultipartForm(e.cfg.MaxUploadSize); err != nil {
		return nil, fmt.Errorf("invalid upload payload: %w", err)
	}

	file, header, err := req.FormFile(formField)
	if err != nil {
		return nil, fmt.Errorf("file is required: %w", err)
	}
	defer file.Close()

	safeSubdir, err := sanitizeSubdir(subdir)
	if err != nil {
		return nil, err
	}

	storageDir := filepath.Join(e.cfg.UploadDir, safeSubdir)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to prepare upload directory: %w", err)
	}

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(storageDir, filename)

	dest, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dest.Close()

	contentType, reader, size, err := sniffAndWrap(file, header.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	written, err := io.Copy(dest, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to write upload: %w", err)
	}
	if size > 0 {
		written = size
	}

	uploadURL := strings.TrimRight(e.UploadURL(), "/")
	relativePath := filename
	if safeSubdir != "" {
		relativePath = filepath.ToSlash(filepath.Join(safeSubdir, filename))
	}
	fileURL := fmt.Sprintf("%s/%s", uploadURL, relativePath)

	return &UploadResult{
		URL:        fileURL,
		Path:       relativePath,
		Filename:   header.Filename,
		Size:       written,
		MimeType:   contentType,
		UploadedAt: time.Now(),
	}, nil
}

func normalizePath(input string) string {
	if input == "" || input == "/" {
		return ""
	}
	path := strings.TrimSpace(input)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimRight(path, "/")
}

func sanitizeSubdir(subdir string) (string, error) {
	if subdir == "" {
		return "", nil
	}
	clean := filepath.Clean(subdir)
	if strings.Contains(clean, "..") {
		return "", errors.New("invalid upload path")
	}
	clean = strings.TrimPrefix(clean, string(filepath.Separator))
	return clean, nil
}

func sniffAndWrap(file io.Reader, hintedType string) (string, io.Reader, int64, error) {
	if hintedType != "" {
		return hintedType, file, 0, nil
	}

	buf := make([]byte, 512)
	n, err := io.ReadFull(file, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", nil, 0, fmt.Errorf("failed to read upload: %w", err)
	}
	buf = buf[:n]
	detected := http.DetectContentType(buf)
	reader := io.MultiReader(bytes.NewReader(buf), file)
	return detected, reader, 0, nil
}

