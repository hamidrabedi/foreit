package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/forgego/forge/pkg/config"
	"github.com/forgego/forge/pkg/logging"
	"github.com/forgego/forge/pkg/security"
	"go.uber.org/zap"
)

// Server wraps http.Server with framework-specific functionality
type Server struct {
	*http.Server
	router   *Router
	logger   *logging.Logger
	config   *config.Config
	settings *config.Settings
}

// NewServer creates a new framework server
func NewServer(cfg *config.Config, settings *config.Settings, logger *logging.Logger) (*Server, error) {
	// Create router
	router := NewRouter()

	// Add logging middleware
	router.Use(logging.Middleware(logger))

	// Add session middleware if configured
	if settings.Security.SessionSecret != "" {
		sessionManager := security.NewSessionManager([]byte(settings.Security.SessionSecret))
		router.Use(sessionManager.Middleware())
	}

	// Add CSRF middleware for non-GET requests
	if settings.Security.CSRFSecretKey != "" {
		csrfProtect := security.NewCSRF(
			[]byte(settings.Security.CSRFSecretKey),
			security.DefaultOptions()...,
		)
		router.Use(csrfProtect.Middleware())
	}

	// Create server
	addr := settings.Server.Host + ":" + settings.Server.Port
	server := &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  time.Duration(settings.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(settings.Server.WriteTimeout) * time.Second,
		},
		router:   router,
		logger:   logger,
		config:   cfg,
		settings: settings,
	}

	return server, nil
}

// RegisterRoutes registers routes on the server's router
func (s *Server) RegisterRoutes(fn func(*Router)) {
	fn(s.router)
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.Info("Starting server",
		zap.String("address", s.Addr),
		zap.String("environment", s.settings.App.Env),
	)

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	return s.Server.Shutdown(ctx)
}
