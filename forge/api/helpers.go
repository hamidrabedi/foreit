package api

import (
	"github.com/forgego/forge/api/parsers"
	"github.com/forgego/forge/api/renderers"
)

// SetupDefaultAPI sets up the API with sensible defaults
func SetupDefaultAPI() {
	Initialize()

	// Set default renderers
	SetDefaultRenderers(
		renderers.NewJSONRenderer(),
		renderers.NewHTMLRenderer(),
	)

	// Set default parsers
	SetDefaultParsers(
		parsers.NewJSONParser(),
		parsers.NewFormParser(),
	)
}

// SetDefaultRenderers sets default renderers.
func SetDefaultRenderers(rendererList ...renderers.Renderer) {
	settings := GetSettings()
	settings.DefaultRenderers = rendererList
}

// SetDefaultParsers sets default parsers.
func SetDefaultParsers(parserList ...parsers.Parser) {
	settings := GetSettings()
	settings.DefaultParsers = parserList
}
