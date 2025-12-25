package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Store represents a session store
type Store interface {
	Get(ctx *fiber.Ctx, key string) (interface{}, error)
	Set(ctx *fiber.Ctx, key string, value interface{}) error
	Delete(ctx *fiber.Ctx, key string) error
	Clear(ctx *fiber.Ctx) error
	Regenerate(ctx *fiber.Ctx) error
}

// MemoryStore is an in-memory session store
type MemoryStore struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	config   Config
}

// Session represents a session
type Session struct {
	ID        string
	Data      map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Config configures a session store
type Config struct {
	CookieName     string
	CookiePath     string
	CookieDomain   string
	CookieSecure   bool
	CookieHTTPOnly bool
	Lifetime       time.Duration
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		CookieName:     "session",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		Lifetime:       24 * time.Hour,
	}
}

// NewMemoryStore creates a new in-memory session store
func NewMemoryStore(config Config) *MemoryStore {
	if config.CookieName == "" {
		config = DefaultConfig()
	}
	
	store := &MemoryStore{
		sessions: make(map[string]*Session),
		config:   config,
	}
	
	// Start cleanup goroutine
	go store.cleanup()
	
	return store
}

// Get retrieves a session value
func (s *MemoryStore) Get(ctx *fiber.Ctx, key string) (interface{}, error) {
	session, err := s.getSession(ctx)
	if err != nil {
		return nil, err
	}
	
	if session == nil {
		return nil, nil
	}
	
	return session.Data[key], nil
}

// Set sets a session value
func (s *MemoryStore) Set(ctx *fiber.Ctx, key string, value interface{}) error {
	session, err := s.getOrCreateSession(ctx)
	if err != nil {
		return err
	}
	
	session.Data[key] = value
	return s.saveSession(ctx, session)
}

// Delete deletes a session value
func (s *MemoryStore) Delete(ctx *fiber.Ctx, key string) error {
	session, err := s.getSession(ctx)
	if err != nil {
		return err
	}
	
	if session != nil {
		delete(session.Data, key)
		return s.saveSession(ctx, session)
	}
	
	return nil
}

// Clear clears all session data
func (s *MemoryStore) Clear(ctx *fiber.Ctx) error {
	sessionID := s.getSessionID(ctx)
	if sessionID != "" {
		s.mutex.Lock()
		delete(s.sessions, sessionID)
		s.mutex.Unlock()
		
		// Clear cookie
		ctx.Cookie(&fiber.Cookie{
			Name:     s.config.CookieName,
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			Path:     s.config.CookiePath,
			Domain:   s.config.CookieDomain,
			Secure:   s.config.CookieSecure,
			HTTPOnly: s.config.CookieHTTPOnly,
		})
	}
	return nil
}

// Regenerate regenerates the session ID
func (s *MemoryStore) Regenerate(ctx *fiber.Ctx) error {
	oldID := s.getSessionID(ctx)
	
	// Create new session
	session, err := s.getOrCreateSession(ctx)
	if err != nil {
		return err
	}
	
	// Delete old session if exists
	if oldID != "" && oldID != session.ID {
		s.mutex.Lock()
		delete(s.sessions, oldID)
		s.mutex.Unlock()
	}
	
	return s.saveSession(ctx, session)
}

// getSession retrieves a session
func (s *MemoryStore) getSession(ctx *fiber.Ctx) (*Session, error) {
	sessionID := s.getSessionID(ctx)
	if sessionID == "" {
		return nil, nil
	}
	
	s.mutex.RLock()
	session, ok := s.sessions[sessionID]
	s.mutex.RUnlock()
	
	if !ok {
		return nil, nil
	}
	
	// Check expiration
	if !session.ExpiresAt.IsZero() && time.Now().After(session.ExpiresAt) {
		s.mutex.Lock()
		delete(s.sessions, sessionID)
		s.mutex.Unlock()
		return nil, nil
	}
	
	return session, nil
}

// getOrCreateSession gets or creates a session
func (s *MemoryStore) getOrCreateSession(ctx *fiber.Ctx) (*Session, error) {
	session, err := s.getSession(ctx)
	if err != nil || session != nil {
		return session, err
	}
	
	// Create new session
	sessionID := generateSessionID()
	session = &Session{
		ID:        sessionID,
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.config.Lifetime),
	}
	
	return session, nil
}

// saveSession saves a session
func (s *MemoryStore) saveSession(ctx *fiber.Ctx, session *Session) error {
	s.mutex.Lock()
	s.sessions[session.ID] = session
	s.mutex.Unlock()
	
	// Set cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     s.config.CookieName,
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		Path:     s.config.CookiePath,
		Domain:   s.config.CookieDomain,
		Secure:   s.config.CookieSecure,
		HTTPOnly: s.config.CookieHTTPOnly,
	})
	
	return nil
}

// getSessionID gets the session ID from cookie
func (s *MemoryStore) getSessionID(ctx *fiber.Ctx) string {
	return ctx.Cookies(s.config.CookieName)
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// cleanup periodically removes expired sessions
func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for id, session := range s.sessions {
			if !session.ExpiresAt.IsZero() && now.After(session.ExpiresAt) {
				delete(s.sessions, id)
			}
		}
		s.mutex.Unlock()
	}
}

// Middleware creates session middleware
func Middleware(store Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("session", store)
		return c.Next()
	}
}

// Get retrieves a session value
func Get(c *fiber.Ctx, key string) (interface{}, error) {
	store := c.Locals("session").(Store)
	return store.Get(c, key)
}

// Set sets a session value
func Set(c *fiber.Ctx, key string, value interface{}) error {
	store := c.Locals("session").(Store)
	return store.Set(c, key, value)
}

// Delete deletes a session value
func Delete(c *fiber.Ctx, key string) error {
	store := c.Locals("session").(Store)
	return store.Delete(c, key)
}

// Clear clears all session data
func Clear(c *fiber.Ctx) error {
	store := c.Locals("session").(Store)
	return store.Clear(c)
}

