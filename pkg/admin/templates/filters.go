package templates

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
)

// FuncMap returns a template function map with sprig functions and custom helpers
func FuncMap() template.FuncMap {
	fm := sprig.FuncMap()

	// Add custom framework functions
	fm["title"] = strings.Title
	fm["lower"] = strings.ToLower
	fm["upper"] = strings.ToUpper
	fm["truncate"] = truncate
	fm["formatDate"] = formatDate
	fm["formatDateTime"] = formatDateTime
	fm["pluralize"] = pluralize

	// Add basic math functions (in case sprig doesn't have them)
	fm["add"] = func(a, b int) int { return a + b }
	fm["sub"] = func(a, b int) int { return a - b }

	// Add index helper for map access
	fm["index"] = func(data, key interface{}) interface{} {
		// Helper for accessing map fields in templates
		if m, ok := data.(map[string]interface{}); ok {
			if k, ok := key.(string); ok {
				return m[k]
			}
		}
		return nil
	}

	// Add printf for formatting
	fm["printf"] = fmt.Sprintf

	return fm
}

// truncate truncates a string to a specified length
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// formatDate formats a time.Time as a date string
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// formatDateTime formats a time.Time as a datetime string
func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// pluralize returns plural form of a word (simple implementation)
func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	if plural == "" {
		return singular + "s"
	}
	return plural
}
