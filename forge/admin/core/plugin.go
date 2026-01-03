package core

import (
	"context"

	"github.com/go-chi/chi/v5"
)

// Plugin defines the interface for admin system extensions
type Plugin interface {
	// Name returns the unique identifier for this plugin
	Name() string

	// Initialize is called when the plugin is registered with the Site
	Initialize(ctx context.Context, site interface{}) error

	// RegisterRoutes allows the plugin to register its own API endpoints
	// The prefix for these routes will typically be /api/admin/plugins/{name}
	RegisterRoutes(router chi.Router)

	// GetMetadata returns information about UI extensions provided by this plugin
	GetMetadata() PluginMetadata
}

// PluginMetadata contains UI-side extension information
type PluginMetadata struct {
	Name        string       `json:"name"`
	Label       string       `json:"label"`
	Icon        string       `json:"icon,omitempty"`
	Pages       []CustomPage `json:"pages,omitempty"`
	Widgets     []WidgetMeta `json:"widgets,omitempty"`
	MenuEntries []MenuEntry  `json:"menuEntries,omitempty"`
}

// CustomPage represents a non-model page in the admin UI
type CustomPage struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	Path      string `json:"path"`
	Component string `json:"component"` // Corresponds to a name in the React registry
}

// WidgetMeta represents a dashboard widget configuration
type WidgetMeta struct {
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	DefaultSize string                 `json:"defaultSize,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// MenuEntry represents an entry in the sidebar navigation
type MenuEntry struct {
	Label    string      `json:"label"`
	Icon     string      `json:"icon,omitempty"`
	Path     string      `json:"path,omitempty"`
	Children []MenuEntry `json:"children,omitempty"`
	Order    int         `json:"order,omitempty"`
}
