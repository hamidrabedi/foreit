package static

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Config configures static file serving
type Config struct {
	Root      string
	Prefix    string
	Index     string
	Compress  bool
	MaxAge    int
	Immutable bool
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Root:     "./public",
		Prefix:   "/static",
		Index:    "index.html",
		Compress: true,
		MaxAge:   31536000, // 1 year
	}
}

// New creates a new static file handler
func New(config Config) fiber.Handler {
	if config.Root == "" {
		config = DefaultConfig()
	}

	return func(c *fiber.Ctx) error {
		// Get file path from URL
		path := c.Path()
		
		// Remove prefix if present
		if config.Prefix != "" && strings.HasPrefix(path, config.Prefix) {
			path = strings.TrimPrefix(path, config.Prefix)
		}
		
		// Clean path
		path = filepath.Clean(path)
		if path == "." || path == "/" {
			path = config.Index
		}
		
		// Build full file path
		fullPath := filepath.Join(config.Root, path)
		
		// Security: ensure path is within root
		rootAbs, _ := filepath.Abs(config.Root)
		fileAbs, _ := filepath.Abs(fullPath)
		if !strings.HasPrefix(fileAbs, rootAbs) {
			return c.SendStatus(fiber.StatusForbidden)
		}
		
		// Check if file exists
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return err
		}
		
		// Don't serve directories
		if info.IsDir() {
			return c.SendStatus(fiber.StatusForbidden)
		}
		
		// Set content type
		ext := filepath.Ext(fullPath)
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			c.Set("Content-Type", mimeType)
		}
		
		// Set cache headers
		if config.MaxAge > 0 {
			c.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", config.MaxAge))
			if config.Immutable {
				c.Set("Cache-Control", fmt.Sprintf("public, max-age=%d, immutable", config.MaxAge))
			}
		}
		
		// Set ETag
		etag := generateETag(info)
		c.Set("ETag", etag)
		
		// Check If-None-Match
		if match := c.Get("If-None-Match"); match == etag {
			return c.SendStatus(fiber.StatusNotModified)
		}
		
		// Serve file
		return c.SendFile(fullPath, config.Compress)
	}
}

// generateETag generates an ETag from file info
func generateETag(info os.FileInfo) string {
	// Simple ETag based on modification time and size
	return fmt.Sprintf(`"%x-%x"`, info.ModTime().Unix(), info.Size())
}

// ServeFile serves a single file
func ServeFile(path string, compress bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendFile(path, compress)
	}
}

