package testhelpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TempWorkdir creates a temporary working directory for tests
func TempWorkdir(t *testing.T, prefix string) (string, func()) {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to clean up temp dir: %v", err)
		}
	}

	return dir, cleanup
}

// CaptureGeneratedMigrations scans a directory for migration files
func CaptureGeneratedMigrations(dir string) ([]string, error) {
	var migrations []string
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// Look for migration files (typically NNNN_name.up.sql or similar)
		name := entry.Name()
		if len(name) >= 5 && name[0] >= '0' && name[0] <= '9' {
			migrations = append(migrations, filepath.Join(dir, name))
		}
	}

	return migrations, nil
}

// ReadFileString reads a file as a string for testing
func ReadFileString(t *testing.T, path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}

// WriteFileString writes a string to a file for testing
func WriteFileString(t *testing.T, path string, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
