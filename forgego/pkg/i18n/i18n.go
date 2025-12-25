package i18n

import (
	"context"
	"fmt"
	"strings"
)

type Translator interface {
	Translate(ctx context.Context, key string, args ...interface{}) string
	TranslatePlural(ctx context.Context, key string, count int, args ...interface{}) string
	SetLocale(locale string)
	GetLocale() string
}

type Message struct {
	ID          string
	Translation string
	Plural      string
}

type Bundle struct {
	locale   string
	messages map[string]*Message
}

func NewBundle(locale string) *Bundle {
	return &Bundle{
		locale:   locale,
		messages: make(map[string]*Message),
	}
}

func (b *Bundle) AddMessage(key, translation string) {
	b.messages[key] = &Message{
		ID:          key,
		Translation: translation,
	}
}

func (b *Bundle) AddPluralMessage(key, singular, plural string) {
	b.messages[key] = &Message{
		ID:          key,
		Translation: singular,
		Plural:      plural,
	}
}

func (b *Bundle) Translate(key string, args ...interface{}) string {
	msg, ok := b.messages[key]
	if !ok {
		return key
	}
	
	translation := msg.Translation
	if len(args) > 0 {
		translation = fmt.Sprintf(translation, args...)
	}
	
	return translation
}

func (b *Bundle) TranslatePlural(key string, count int, args ...interface{}) string {
	msg, ok := b.messages[key]
	if !ok {
		return key
	}
	
	var translation string
	if count == 1 {
		translation = msg.Translation
	} else {
		if msg.Plural != "" {
			translation = msg.Plural
		} else {
			translation = msg.Translation
		}
	}
	
	if len(args) > 0 {
		translation = fmt.Sprintf(translation, args...)
	} else {
		translation = strings.ReplaceAll(translation, "%d", fmt.Sprintf("%d", count))
	}
	
	return translation
}

type Manager struct {
	bundles map[string]*Bundle
	defaultLocale string
}

func NewManager(defaultLocale string) *Manager {
	return &Manager{
		bundles:       make(map[string]*Bundle),
		defaultLocale: defaultLocale,
	}
}

func (m *Manager) RegisterBundle(locale string, bundle *Bundle) {
	m.bundles[locale] = bundle
}

func (m *Manager) GetBundle(locale string) *Bundle {
	bundle, ok := m.bundles[locale]
	if !ok {
		bundle, ok = m.bundles[m.defaultLocale]
		if !ok {
			return NewBundle(locale)
		}
	}
	return bundle
}

func (m *Manager) Translate(ctx context.Context, locale, key string, args ...interface{}) string {
	bundle := m.GetBundle(locale)
	return bundle.Translate(key, args...)
}

func (m *Manager) TranslatePlural(ctx context.Context, locale, key string, count int, args ...interface{}) string {
	bundle := m.GetBundle(locale)
	return bundle.TranslatePlural(key, count, args...)
}

var defaultManager *Manager

func SetDefaultManager(manager *Manager) {
	defaultManager = manager
}

func T(ctx context.Context, key string, args ...interface{}) string {
	if defaultManager == nil {
		return key
	}
	
	locale := getLocaleFromContext(ctx)
	return defaultManager.Translate(ctx, locale, key, args...)
}

func TP(ctx context.Context, key string, count int, args ...interface{}) string {
	if defaultManager == nil {
		return key
	}
	
	locale := getLocaleFromContext(ctx)
	return defaultManager.TranslatePlural(ctx, locale, key, count, args...)
}

func getLocaleFromContext(ctx context.Context) string {
	if locale, ok := ctx.Value("locale").(string); ok {
		return locale
	}
	return "en"
}

