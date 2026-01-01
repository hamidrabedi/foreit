package storage

import (
	"context"
	"io"
	"time"
)

// FileInfo contains metadata about a stored file
type FileInfo struct {
	Path      string
	URL       string
	Size      int64
	MimeType  string
	CreatedAt time.Time
}

// Storage is the interface for file storage backends
type Storage interface {
	// Save saves a file and returns its path and URL
	Save(ctx context.Context, name string, reader io.Reader, size int64, mimeType string) (*FileInfo, error)

	// Get retrieves a file by path
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete deletes a file by path
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)

	// URL returns the public URL for a file path
	URL(path string) string

	// GenerateThumbnail generates a thumbnail for an image
	GenerateThumbnail(ctx context.Context, imagePath string, width, height int) (string, error)
}

// StorageConfig contains configuration for storage
type StorageConfig struct {
	// BasePath is the base directory for local storage
	BasePath string

	// BaseURL is the base URL for serving files
	BaseURL string

	// MaxFileSize is the maximum file size in bytes
	MaxFileSize int64

	// AllowedMimeTypes is a list of allowed MIME types (empty = all)
	AllowedMimeTypes []string

	// GenerateThumbnails enables automatic thumbnail generation for images
	GenerateThumbnails bool

	// ThumbnailSizes is a list of thumbnail sizes to generate
	ThumbnailSizes []ThumbnailSize
}

// ThumbnailSize represents a thumbnail size
type ThumbnailSize struct {
	Width  int
	Height int
	Name   string // e.g., "small", "medium", "large"
}

// DefaultStorageConfig returns a default storage configuration
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		BasePath:          "./uploads",
		BaseURL:           "/uploads",
		MaxFileSize:       10 * 1024 * 1024, // 10MB
		GenerateThumbnails: true,
		ThumbnailSizes: []ThumbnailSize{
			{Width: 150, Height: 150, Name: "thumb"},
			{Width: 300, Height: 300, Name: "small"},
			{Width: 800, Height: 800, Name: "medium"},
		},
	}
}
