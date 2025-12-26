package security

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

// SessionManager wraps scs.SessionManager with framework-specific methods
type SessionManager struct {
	*scs.SessionManager
}

// NewSessionManager creates a new session manager
func NewSessionManager(secretKey []byte) *SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "forge_session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = false // Set to true in production
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Path = "/"

	return &SessionManager{
		SessionManager: sessionManager,
	}
}

// Middleware returns the session middleware
func (sm *SessionManager) Middleware() func(http.Handler) http.Handler {
	return sm.SessionManager.LoadAndSave
}

// Get retrieves a value from the session
func (sm *SessionManager) Get(r *http.Request, key string) interface{} {
	return sm.SessionManager.Get(r.Context(), key)
}

// Put stores a value in the session
func (sm *SessionManager) Put(r *http.Request, key string, val interface{}) {
	sm.SessionManager.Put(r.Context(), key, val)
}

// Remove removes a value from the session
func (sm *SessionManager) Remove(r *http.Request, key string) {
	sm.SessionManager.Remove(r.Context(), key)
}

// Exists checks if a key exists in the session
func (sm *SessionManager) Exists(r *http.Request, key string) bool {
	return sm.SessionManager.Exists(r.Context(), key)
}

// Destroy destroys the session
func (sm *SessionManager) Destroy(r *http.Request) error {
	return sm.SessionManager.Destroy(r.Context())
}

// RenewToken renews the session token
func (sm *SessionManager) RenewToken(r *http.Request) error {
	return sm.SessionManager.RenewToken(r.Context())
}
