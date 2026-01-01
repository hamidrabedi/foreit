package i18n

import (
	"context"
	"fmt"
	"sync"
)

// Translator translates strings to different languages
type Translator interface {
	Translate(ctx context.Context, key string, args ...interface{}) string
	GetLanguage(ctx context.Context) string
	SetLanguage(ctx context.Context, lang string) error
}

// MessageCatalog contains translations for a language
type MessageCatalog map[string]string

// I18nManager manages internationalization
type I18nManager struct {
	catalogs map[string]MessageCatalog
	defaultLang string
	mu       sync.RWMutex
}

// NewI18nManager creates a new i18n manager
func NewI18nManager(defaultLang string) *I18nManager {
	return &I18nManager{
		catalogs:    make(map[string]MessageCatalog),
		defaultLang: defaultLang,
	}
}

// RegisterLanguage registers translations for a language
func (m *I18nManager) RegisterLanguage(lang string, catalog MessageCatalog) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.catalogs[lang] = catalog
}

// Translate translates a key to the specified language
func (m *I18nManager) Translate(lang, key string, args ...interface{}) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalog, ok := m.catalogs[lang]
	if !ok {
		// Fallback to default language
		catalog, ok = m.catalogs[m.defaultLang]
		if !ok {
			// Return key if no translation found
			return key
		}
	}

	message, ok := catalog[key]
	if !ok {
		// Return key if translation not found
		return key
	}

	// Format message with args if provided
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}

	return message
}

// GetLanguage gets the current language from context
func (m *I18nManager) GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value("language").(string); ok {
		return lang
	}
	return m.defaultLang
}

// SetLanguage sets the language in context
func (m *I18nManager) SetLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, "language", lang)
}

// DefaultTranslator provides a simple translator implementation
type DefaultTranslator struct {
	manager *I18nManager
}

// NewDefaultTranslator creates a new default translator
func NewDefaultTranslator(manager *I18nManager) *DefaultTranslator {
	return &DefaultTranslator{
		manager: manager,
	}
}

// Translate translates a key
func (t *DefaultTranslator) Translate(ctx context.Context, key string, args ...interface{}) string {
	lang := t.manager.GetLanguage(ctx)
	return t.manager.Translate(lang, key, args...)
}

// GetLanguage gets the current language
func (t *DefaultTranslator) GetLanguage(ctx context.Context) string {
	return t.manager.GetLanguage(ctx)
}

// SetLanguage sets the language
func (t *DefaultTranslator) SetLanguage(ctx context.Context, lang string) error {
	ctx = t.manager.SetLanguage(ctx, lang)
	return nil
}

// Common translation keys
const (
	KeyAdd    = "admin.add"
	KeyChange = "admin.change"
	KeyDelete = "admin.delete"
	KeyView   = "admin.view"
	KeySave   = "admin.save"
	KeyCancel = "admin.cancel"
	KeySearch = "admin.search"
	KeyFilter = "admin.filter"
)

// DefaultEnglishCatalog returns default English translations
func DefaultEnglishCatalog() MessageCatalog {
	return MessageCatalog{
		KeyAdd:    "Add",
		KeyChange: "Change",
		KeyDelete: "Delete",
		KeyView:   "View",
		KeySave:   "Save",
		KeyCancel: "Cancel",
		KeySearch: "Search",
		KeyFilter: "Filter",
	}
}
