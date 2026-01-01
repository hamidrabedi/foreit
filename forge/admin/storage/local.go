package storage

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	_ "image/gif"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorage implements Storage using the local filesystem
type LocalStorage struct {
	config *StorageConfig
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(config *StorageConfig) (*LocalStorage, error) {
	if config == nil {
		config = DefaultStorageConfig()
	}

	// Ensure base directory exists
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		config: config,
	}, nil
}

// Save saves a file to local storage
func (ls *LocalStorage) Save(ctx context.Context, name string, reader io.Reader, size int64, mimeType string) (*FileInfo, error) {
	// Validate file size
	if size > ls.config.MaxFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum %d", size, ls.config.MaxFileSize)
	}

	// Validate MIME type if restrictions are set
	if len(ls.config.AllowedMimeTypes) > 0 {
		allowed := false
		for _, allowedType := range ls.config.AllowedMimeTypes {
			if mimeType == allowedType || strings.HasPrefix(mimeType, allowedType+"/") {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("MIME type %s is not allowed", mimeType)
		}
	}

	// Generate unique filename
	filename := ls.generateFilename(name)
	filePath := filepath.Join(ls.config.BasePath, filename)

	// Create directory if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data
	written, err := io.Copy(file, reader)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	if written != size {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("file size mismatch: expected %d, wrote %d", size, written)
	}

	info := &FileInfo{
		Path:      filename,
		URL:       ls.URL(filename),
		Size:      written,
		MimeType:  mimeType,
		CreatedAt: time.Now(),
	}

	// Generate thumbnails for images
	if ls.config.GenerateThumbnails && strings.HasPrefix(mimeType, "image/") {
		for _, thumbSize := range ls.config.ThumbnailSizes {
			_, err := ls.GenerateThumbnail(ctx, filename, thumbSize.Width, thumbSize.Height)
			if err != nil {
				// Log error but don't fail the upload
				// In production, you'd want proper logging here
				_ = err
			}
		}
	}

	return info, nil
}

// Get retrieves a file from local storage
func (ls *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(ls.config.BasePath, path)
	
	// Security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(ls.config.BasePath)) {
		return nil, fmt.Errorf("invalid file path")
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete deletes a file from local storage
func (ls *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(ls.config.BasePath, path)
	
	// Security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(ls.config.BasePath)) {
		return fmt.Errorf("invalid file path")
	}

	// Delete main file
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete thumbnails
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	for _, thumbSize := range ls.config.ThumbnailSizes {
		thumbPath := filepath.Join(ls.config.BasePath, base+"_"+thumbSize.Name+ext)
		os.Remove(thumbPath) // Ignore errors for thumbnails
	}

	return nil
}

// Exists checks if a file exists
func (ls *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(ls.config.BasePath, path)
	
	// Security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(ls.config.BasePath)) {
		return false, fmt.Errorf("invalid file path")
	}

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// URL returns the public URL for a file path
func (ls *LocalStorage) URL(path string) string {
	return strings.TrimSuffix(ls.config.BaseURL, "/") + "/" + strings.TrimPrefix(path, "/")
}

// GenerateThumbnail generates a thumbnail for an image
func (ls *LocalStorage) GenerateThumbnail(ctx context.Context, imagePath string, width, height int) (string, error) {
	fullPath := filepath.Join(ls.config.BasePath, imagePath)
	
	// Security: prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(fullPath), filepath.Clean(ls.config.BasePath)) {
		return "", fmt.Errorf("invalid file path")
	}

	// Open image
	file, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image using simple nearest-neighbor for now
	// In production, you'd want to use a proper image library like github.com/nfnt/resize
	resized := resizeImage(img, width, height)

	// Generate thumbnail filename
	ext := filepath.Ext(imagePath)
	base := strings.TrimSuffix(imagePath, ext)
	thumbName := fmt.Sprintf("%s_%dx%d%s", base, width, height, ext)
	thumbPath := filepath.Join(ls.config.BasePath, thumbName)

	// Create thumbnail file
	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return "", fmt.Errorf("failed to create thumbnail: %w", err)
	}
	defer thumbFile.Close()

	// Encode thumbnail
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(thumbFile, resized, &jpeg.Options{Quality: 85})
	case "png":
		err = png.Encode(thumbFile, resized)
	default:
		return "", fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return thumbName, nil
}

// resizeImage resizes an image (simple implementation)
// For production, use a proper library like github.com/nfnt/resize
func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// Calculate scaling factors
	scaleX := float64(width) / float64(srcWidth)
	scaleY := float64(height) / float64(srcHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Calculate new dimensions maintaining aspect ratio
	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	// Create new image
	rgba := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Simple nearest-neighbor resize
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			if srcX < srcWidth && srcY < srcHeight {
				rgba.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
			}
		}
	}

	return rgba
}

// generateFilename generates a unique filename
func (ls *LocalStorage) generateFilename(originalName string) string {
	// Extract extension
	ext := filepath.Ext(originalName)
	base := strings.TrimSuffix(filepath.Base(originalName), ext)

	// Sanitize base name
	base = sanitizeFilename(base)
	if base == "" {
		base = "file"
	}

	// Generate unique name with timestamp
	timestamp := time.Now().Format("20060102_150405")
	random := fmt.Sprintf("%d", time.Now().UnixNano()%10000)

	return fmt.Sprintf("%s_%s_%s%s", base, timestamp, random, ext)
}

// sanitizeFilename sanitizes a filename
func sanitizeFilename(name string) string {
	// Remove dangerous characters
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, " ", "_")
	
	// Keep only alphanumeric, dash, underscore, and dot
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}
