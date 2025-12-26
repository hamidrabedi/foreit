package security

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRF wraps gorilla/csrf with framework-specific methods
type CSRF struct {
	protect func(http.Handler) http.Handler
	secret  []byte
}

// NewCSRF creates a new CSRF protector
func NewCSRF(secretKey []byte, opts ...csrf.Option) *CSRF {
	protect := csrf.Protect(secretKey, opts...)
	return &CSRF{
		protect: protect,
		secret:  secretKey,
	}
}

// Middleware returns the CSRF protection middleware
func (c *CSRF) Middleware() func(http.Handler) http.Handler {
	return c.protect
}

// Token returns the CSRF token for the current request
func (c *CSRF) Token(r *http.Request) string {
	return csrf.Token(r)
}

// TemplateField returns a hidden input field for forms
func (c *CSRF) TemplateField(r *http.Request) string {
	return string(csrf.TemplateField(r))
}

// DefaultOptions returns default CSRF options suitable for framework
func DefaultOptions() []csrf.Option {
	return []csrf.Option{
		csrf.Secure(false), // Set to true in production with HTTPS
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Path("/"),
	}
}

// ProductionOptions returns CSRF options for production
func ProductionOptions() []csrf.Option {
	return []csrf.Option{
		csrf.Secure(true), // Require HTTPS
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Path("/"),
	}
}
