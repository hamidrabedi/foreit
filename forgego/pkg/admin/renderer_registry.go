package admin

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const RendererAPIVersion = "1.0"

type RendererRegistry struct {
	version        string
	fallback       bool
	fieldRenderers map[string]FieldRenderer
	customRenderers map[string]FieldRenderer
	mu             sync.RWMutex
}

func NewRendererRegistry(version string, fallback bool) *RendererRegistry {
	registry := &RendererRegistry{
		version:        version,
		fallback:       fallback,
		fieldRenderers: make(map[string]FieldRenderer),
		customRenderers: make(map[string]FieldRenderer),
	}

	registry.registerBuiltins()
	return registry
}
func (r *RendererRegistry) registerBuiltins() {
	r.fieldRenderers["string"] = &StringFieldRenderer{}
	r.fieldRenderers["text"] = &StringFieldRenderer{}
	r.fieldRenderers["int"] = &IntFieldRenderer{}
	r.fieldRenderers["number"] = &IntFieldRenderer{}
	r.fieldRenderers["float"] = &FloatFieldRenderer{}
	r.fieldRenderers["float64"] = &FloatFieldRenderer{}
	r.fieldRenderers["bool"] = &BoolFieldRenderer{}
	r.fieldRenderers["boolean"] = &BoolFieldRenderer{}
	r.fieldRenderers["datetime"] = &DateTimeFieldRenderer{}
	r.fieldRenderers["date"] = &DateTimeFieldRenderer{}
	r.fieldRenderers["foreignkey"] = &ForeignKeyRenderer{}
	r.fieldRenderers["relation"] = &ForeignKeyRenderer{}
}

func (r *RendererRegistry) RegisterRenderer(renderer FieldRenderer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := renderer.Info()

	if !r.isVersionCompatible(info.Version) {
		return fmt.Errorf("renderer %s version %s incompatible with registry version %s",
			info.Name, info.Version, r.version)
	}

	if err := r.validateRenderer(renderer); err != nil {
		return fmt.Errorf("renderer %s validation failed: %w", info.Name, err)
	}

	r.customRenderers[info.Name] = renderer
	return nil
}
func (r *RendererRegistry) validateRenderer(renderer FieldRenderer) error {
	ctx := RenderContext{
		Field: &FieldMeta{
			Name: "test",
			Type: FieldTypeText,
		},
		Value: "test",
	}

	html, err := renderer.RenderHTML(ctx)
	if err != nil {
		return err
	}

	return r.validateHTMLSafety(string(html))
}

func (r *RendererRegistry) GetRenderer(fieldType string) (FieldRenderer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fieldType = strings.ToLower(fieldType)

	if custom, ok := r.customRenderers[fieldType]; ok {
		return custom, nil
	}

	if builtin, ok := r.fieldRenderers[fieldType]; ok {
		return builtin, nil
	}

	if r.fallback {
		if defaultRenderer, ok := r.fieldRenderers["string"]; ok {
			return defaultRenderer, nil
		}
	}

	return nil, fmt.Errorf("no renderer found for type %s", fieldType)
}

func (r *RendererRegistry) GetRendererSafe(fieldType string) FieldRenderer {
	renderer, err := r.GetRenderer(fieldType)
	if err != nil {
		log.Printf("Renderer error for %s: %v", fieldType, err)

		if r.fallback {
			if defaultRenderer, ok := r.fieldRenderers["string"]; ok {
				return defaultRenderer
			}
		}

		panic(err)
	}

	return renderer
}

func (r *RendererRegistry) isVersionCompatible(rendererVersion string) bool {
	if rendererVersion == r.version {
		return true
	}

	registryMajor := strings.Split(r.version, ".")[0]
	rendererMajor := strings.Split(rendererVersion, ".")[0]
	return registryMajor == rendererMajor
}

func (r *RendererRegistry) validateHTMLSafety(htmlStr string) error {
	htmlLower := strings.ToLower(htmlStr)

	if strings.Contains(htmlLower, "<script") {
		return fmt.Errorf("HTML contains unsafe elements (script tags)")
	}

	if strings.Contains(htmlLower, "<iframe") {
		return fmt.Errorf("HTML contains unsafe elements (iframe tags)")
	}
	dangerousAttrs := []string{"onclick", "onerror", "onload", "onmouseover", "onfocus", "onblur"}
	for _, attr := range dangerousAttrs {
		if strings.Contains(htmlLower, attr+"=") {
			return fmt.Errorf("HTML contains unsafe event handlers (%s)", attr)
		}
	}

	return nil
}

