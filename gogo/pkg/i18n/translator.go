package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Translator handles translations
type Translator struct {
	translations map[string]map[string]string
	defaultLocale string
	mutex        sync.RWMutex
}

// NewTranslator creates a new translator
func NewTranslator(defaultLocale string) *Translator {
	return &Translator{
		translations: make(map[string]map[string]string),
		defaultLocale: defaultLocale,
	}
}

// Load loads translations from a directory
func (t *Translator) Load(dir string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Look for JSON files in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		locale := strings.TrimSuffix(entry.Name(), ".json")
		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		t.translations[locale] = translations
	}

	return nil
}

// T translates a key to the specified locale
func (t *Translator) T(locale, key string, args ...interface{}) string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Try specified locale
	if translations, ok := t.translations[locale]; ok {
		if value, ok := translations[key]; ok {
			return formatString(value, args...)
		}
	}

	// Fallback to default locale
	if translations, ok := t.translations[t.defaultLocale]; ok {
		if value, ok := translations[key]; ok {
			return formatString(value, args...)
		}
	}

	// Return key if not found
	return formatString(key, args...)
}

// formatString formats a string with arguments
func formatString(template string, args ...interface{}) string {
	if len(args) == 0 {
		return template
	}

	// Simple placeholder replacement
	// In production, use a proper template engine
	result := template
	for i, arg := range args {
		placeholder := fmt.Sprintf("{%d}", i)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", arg))
	}

	return result
}

// GetLocale extracts locale from Fiber context
func GetLocale(c *fiber.Ctx) string {
	// Try Accept-Language header
	acceptLang := c.Get("Accept-Language")
	if acceptLang != "" {
		// Parse Accept-Language (simplified)
		parts := strings.Split(acceptLang, ",")
		if len(parts) > 0 {
			locale := strings.TrimSpace(strings.Split(parts[0], ";")[0])
			return locale
		}
	}

	// Try locale from query param
	if locale := c.Query("locale"); locale != "" {
		return locale
	}

	// Try locale from session
	if session := c.Locals("session"); session != nil {
		// Would get from session if available
	}

	return "en" // Default
}

// Default translator instance
var defaultTranslator *Translator

// Load loads translations into the default translator
func Load(dir string, defaultLocale string) error {
	defaultTranslator = NewTranslator(defaultLocale)
	return defaultTranslator.Load(dir)
}

// T translates using the default translator
func T(c *fiber.Ctx, key string, args ...interface{}) string {
	if defaultTranslator == nil {
		return key
	}

	locale := GetLocale(c)
	return defaultTranslator.T(locale, key, args...)
}

// Middleware adds translation support to requests
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		locale := GetLocale(c)
		c.Locals("locale", locale)
		return c.Next()
	}
}

