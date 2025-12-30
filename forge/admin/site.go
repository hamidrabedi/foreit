package admin

import (
	"net/http"
)

// Site represents an admin site instance
type Site struct {
	Name       string
	Title      string
	Header     string
	IndexTitle string
	SiteURL    string
	registry   *Registry
}

// NewSite creates a new admin site
func NewSite(name string) *Site {
	return &Site{
		Name:       name,
		Title:      "Admin",
		Header:     "Administration",
		IndexTitle: "Site Administration",
		registry:   NewRegistry(),
	}
}

// GetRegistry returns the registry for this site
func (s *Site) GetRegistry() *Registry {
	return s.registry
}

// IndexView handles the admin index/dashboard
func (s *Site) IndexView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAdmins := s.registry.GetAll()

		models := make([]map[string]interface{}, 0, len(allAdmins))
		for name, admin := range allAdmins {
			models = append(models, map[string]interface{}{
				"name":        name,
				"verboseName": name,
				"modelType":   admin.ModelType().String(),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		// In a real implementation, would render HTML template
		_ = models
	}
}

// NewRegistry creates a new registry (helper function)
func NewRegistry() *Registry {
	return &Registry{
		admins: make(map[string]AdminInterface),
	}
}
