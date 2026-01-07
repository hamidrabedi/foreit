package admin

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/forgego/forge/admin/api/rest"
	"github.com/forgego/forge/admin/core"
	"github.com/forgego/forge/admin/ui"
	"github.com/forgego/forge/media"
	"github.com/forgego/forge/server"
)

// Site represents an admin site instance
type Site struct {
	Name       string
	Title      string
	Header     string
	IndexTitle string
	SiteURL    string
	registry   *core.Registry
	uiConfig   UIConfig
}

// UISource defines where the Admin UI assets come from
type UISource string

const (
	UISourceEmbedded UISource = "embedded" // Use UI files embedded in the forge binary
	UISourceStatic   UISource = "static"   // Use UI files from a local directory
	UISourceExternal UISource = "external" // Use UI from an external URL (e.g. dev server)
)

// UIConfig configures the Admin UI
type UIConfig struct {
	Source    UISource
	StaticDir string
	EmbedFS   fs.FS
	Prefix    string // URL prefix for the admin site, e.g. "/admin"
	BasePath  string // UI router base path, e.g. "/admin"
	UploadDir string // Filesystem directory for uploads
	UploadURL string // Public URL prefix for uploads
}

func DefaultUIConfig() UIConfig {
	return UIConfig{
		Source:    UISourceEmbedded,
		Prefix:    "",
		BasePath:  "/admin",
		EmbedFS:   ui.GetFS(),
		UploadDir: "uploads",
	}
}

// NewSite creates a new admin site
func NewSite(name string) *Site {
	return &Site{
		Name:       name,
		Title:      "Admin",
		Header:     "Administration",
		IndexTitle: "Site Administration",
		registry:   core.NewRegistry(),
		uiConfig:   DefaultUIConfig(),
	}
}

// WithUIConfig sets the UI configuration for the site
func (s *Site) WithUIConfig(config UIConfig) *Site {
	s.uiConfig = config
	return s
}

// GetUIConfig returns the current UI configuration
func (s *Site) GetUIConfig() UIConfig {
	return s.uiConfig
}

// GetRegistry returns the registry for this site
func (s *Site) GetRegistry() *core.Registry {
	return s.registry
}

// RegisterPlugin registers a plugin with the site
func (s *Site) RegisterPlugin(ctx context.Context, p core.Plugin) error {
	if err := s.registry.RegisterPlugin(p); err != nil {
		return err
	}
	return p.Initialize(ctx, s)
}

// RegisterCustomPage registers a custom admin page at the site level.
func (s *Site) RegisterCustomPage(page core.CustomPage) {
	s.registry.RegisterCustomPage(page)
}

// RegisterMenuEntry registers a custom menu entry at the site level.
func (s *Site) RegisterMenuEntry(entry core.MenuEntry) {
	s.registry.RegisterMenuEntry(entry)
}

// SetDashboardConfig sets the dashboard configuration for the site.
func (s *Site) SetDashboardConfig(config core.DashboardConfig) {
	s.registry.SetDashboardConfig(config)
}

// IndexView handles the admin index/dashboard
func (s *Site) IndexView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allAdmins := s.registry.GetAll()

		models := make([]map[string]interface{}, 0, len(allAdmins))
		for name, admin := range allAdmins {
			models = append(models, map[string]interface{}{
				"name":        name,
				"verboseName": name, // Should ideally come from metadata
				"modelType":   admin.ModelType().String(),
			})
		}

		allPlugins := s.registry.GetAllPlugins()
		plugins := make([]core.PluginMetadata, 0, len(allPlugins))
		for _, p := range allPlugins {
			plugins = append(plugins, p.GetMetadata())
		}

		w.Header().Set("Content-Type", "application/json")
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"models":  models,
			"plugins": plugins,
		})
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Handler returns the HTTP handler for the admin site
func (s *Site) Handler() http.Handler {
	r := chi.NewRouter()

	// UI runtime config
	prefix := s.uiConfig.Prefix
	basePath := s.uiBasePath()
	configRoute := "/config.js"
	if prefix != "" {
		configRoute = prefix + "/config.js"
	}
	r.Get(configRoute, s.configHandler())
	if prefix == "" && basePath != "" {
		r.Get(basePath+"/config.js", s.configHandler())
	}

	// 1. Register API Routes
	apiRouter := rest.NewRouter(s.registry)
	mediaEngine := media.New(media.Config{
		UploadDir:     s.uiConfig.UploadDir,
		UploadURL:     s.uiUploadURL(),
		MaxUploadSize: 10 << 20,
	})
	apiRouter.SetMediaEngine(mediaEngine)
	apiRouter.RegisterRoutes(r)

	// Register Plugin Routes
	for _, p := range s.registry.GetAllPlugins() {
		p.RegisterRoutes(r)
	}

	uploadHandler := mediaEngine.MediaHandler()
	if uploadHandler != nil {
		r.Mount(mediaEngine.UploadURL(), uploadHandler)
	}

	// 2. Serve Static UI Assets
	// Determine the route pattern
	routePattern := "/*"
	if prefix != "" {
		routePattern = prefix + "/*"
	}

	uiHandler := s.uiHandler(prefix, basePath)
	if uiHandler != nil {
		r.Handle(routePattern, uiHandler)
		if prefix == "" {
			r.Handle("/", uiHandler)
		} else {
			r.Handle(prefix, uiHandler)
			r.Handle(prefix+"/", uiHandler)
		}
	}

	return r
}

func (s *Site) configHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		basePath := s.uiBasePath()
		uploadURL := s.uiUploadURL()
		apiBase := "/api"
		overridesURL := "/overrides.js"
		if basePath != "" {
			apiBase = basePath + "/api"
			overridesURL = basePath + "/overrides.js"
		}

		payload := map[string]string{
			"basePath":     basePath,
			"apiBase":      apiBase,
			"overridesUrl": overridesURL,
			"uploadUrl":    uploadURL,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = w.Write([]byte("window.__FORGE_ADMIN__ = "))
		_, _ = w.Write(data)
		_, _ = w.Write([]byte(";"))
	}
}

func (s *Site) uiHandler(prefix string, basePath string) http.Handler {
	var (
		staticDir string
		embedFS   fs.FS
	)

	switch s.uiConfig.Source {
	case UISourceStatic:
		staticDir = s.uiConfig.StaticDir
	case UISourceEmbedded:
		embedFS = s.uiConfig.EmbedFS
	}

	if staticDir == "" && embedFS == nil {
		return nil
	}

	var staticHandler http.Handler
	if staticDir != "" {
		staticHandler = server.StaticFiles("", staticDir, server.WithIndexFiles("index.html"))
	} else {
		staticHandler = server.StaticFS("", embedFS, server.WithIndexFiles("index.html"))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		stripPrefix := prefix
		if stripPrefix == "" && basePath != "" && strings.HasPrefix(requestPath, basePath) {
			stripPrefix = basePath
		}
		if stripPrefix != "" {
			if !strings.HasPrefix(requestPath, stripPrefix) {
				http.NotFound(w, r)
				return
			}
			trimmed := strings.TrimPrefix(requestPath, stripPrefix)
			if trimmed == "" {
				trimmed = "/"
			}
			cloned := r.Clone(r.Context())
			cloned.URL.Path = trimmed
			r = cloned
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			staticHandler.ServeHTTP(w, r)
			return
		}

		if isReservedAdminPath(r.URL.Path, "") {
			http.NotFound(w, r)
			return
		}

		filePath, ok := uiFilePath(r.URL.Path, "")
		if !ok {
			http.NotFound(w, r)
			return
		}

		if filePath == "" {
			s.serveIndex(w, r, staticDir, embedFS)
			return
		}

		if s.uiFileExists(staticDir, embedFS, filePath) {
			staticHandler.ServeHTTP(w, r)
			return
		}

		if path.Ext(filePath) != "" || !acceptsHTML(r) {
			http.NotFound(w, r)
			return
		}

		s.serveIndex(w, r, staticDir, embedFS)
	})
}

func (s *Site) uiFileExists(staticDir string, embedFS fs.FS, filePath string) bool {
	if staticDir != "" {
		_, err := os.Stat(filepath.Join(staticDir, filepath.FromSlash(filePath)))
		return err == nil
	}

	if embedFS != nil {
		_, err := fs.Stat(embedFS, filePath)
		return err == nil
	}

	return false
}

func (s *Site) serveIndex(w http.ResponseWriter, r *http.Request, staticDir string, embedFS fs.FS) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var (
		file io.ReadCloser
		err  error
	)

	if staticDir != "" {
		file, err = os.Open(filepath.Join(staticDir, "index.html"))
	} else if embedFS != nil {
		file, err = embedFS.Open("index.html")
	}

	if err != nil || file == nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	_, _ = io.Copy(w, file)
}

func (s *Site) uiBasePath() string {
	if s.uiConfig.BasePath != "" {
		return normalizeBasePath(s.uiConfig.BasePath)
	}
	if s.uiConfig.Prefix != "" {
		return normalizeBasePath(s.uiConfig.Prefix)
	}
	return "/admin"
}

func (s *Site) uiUploadURL() string {
	if s.uiConfig.UploadDir == "" {
		return ""
	}
	if s.uiConfig.UploadURL != "" {
		return normalizeBasePath(s.uiConfig.UploadURL)
	}
	basePath := s.uiBasePath()
	if basePath == "" {
		return "/uploads"
	}
	return basePath + "/uploads"
}

func normalizeBasePath(input string) string {
	pathValue := strings.TrimSpace(input)
	if pathValue == "" || pathValue == "/" {
		return ""
	}
	if !strings.HasPrefix(pathValue, "/") {
		pathValue = "/" + pathValue
	}
	pathValue = strings.TrimRight(pathValue, "/")
	if pathValue == "/" {
		return ""
	}
	return pathValue
}

func acceptsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if accept == "" {
		return true
	}
	return strings.Contains(accept, "text/html")
}

func uiFilePath(urlPath string, prefix string) (string, bool) {
	if prefix != "" {
		if !strings.HasPrefix(urlPath, prefix) {
			return "", false
		}
		urlPath = strings.TrimPrefix(urlPath, prefix)
	}

	urlPath = path.Clean("/" + urlPath)
	if strings.Contains(urlPath, "..") {
		return "", false
	}

	return strings.TrimPrefix(urlPath, "/"), true
}

func isReservedAdminPath(urlPath string, prefix string) bool {
	cleanPath := urlPath
	if prefix != "" && strings.HasPrefix(cleanPath, prefix) {
		cleanPath = strings.TrimPrefix(cleanPath, prefix)
	}
	return strings.HasPrefix(cleanPath, "/api")
}
