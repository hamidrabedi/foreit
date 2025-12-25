package admin

import (
	"html/template"

	"github.com/gofiber/fiber/v2"
)

// RendererInfo provides metadata about a renderer
type RendererInfo struct {
	Name    string
	Version string
	Author  string // for audit trails
}

// FieldRenderer is the core interface - simple and focused
type FieldRenderer interface {
	// Info returns metadata about the renderer
	Info() RendererInfo

	// RenderHTML returns HTML for the field
	// Framework will add HTMX attributes based on metadata
	RenderHTML(ctx RenderContext) (template.HTML, error)

	// Validate validates the field value
	Validate(value interface{}) error
}

// RenderContext provides data needed for rendering
type RenderContext struct {
	Model     *ModelMeta
	Field     *FieldMeta
	Value     interface{}
	User      interface{}
	Request   *fiber.Ctx
	Metadata  map[string]interface{}
	HTMXAttrs map[string]string // Explicit HTMX attributes from metadata
}

