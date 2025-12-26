// Package forge provides HTTP server functionality
package forge

import (
	pkgHTTP "github.com/forgego/forge/pkg/http"
)

// Server wraps the pkg HTTP server
type Server = pkgHTTP.Server

// Router wraps the pkg HTTP router
type Router = pkgHTTP.Router

// NewServer creates a new framework server
func NewServer(cfg *Config, settings *Settings, logger *Logger) (*Server, error) {
	return pkgHTTP.NewServer(cfg, settings, logger)
}
