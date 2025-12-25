package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

// Authenticator authenticates requests
type Authenticator interface {
	Authenticate(c *fiber.Ctx) (interface{}, error)
}

// User represents an authenticated user
type User interface {
	GetID() interface{}
	IsAuthenticated() bool
}

// JWT authenticates using JWT tokens
type JWT struct {
	Secret     string
	HeaderName string
	Extractor  func(*fiber.Ctx) string
}

// NewJWT creates a new JWT authenticator
func NewJWT(secret string) *JWT {
	return &JWT{
		Secret:     secret,
		HeaderName: "Authorization",
		Extractor: func(c *fiber.Ctx) string {
			auth := c.Get("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				return auth[7:]
			}
			return ""
		},
	}
}

// Authenticate authenticates a request using JWT
func (j *JWT) Authenticate(c *fiber.Ctx) (interface{}, error) {
	token := j.Extractor(c)
	if token == "" {
		return nil, errors.New("no token provided")
	}
	
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	
	return claims, nil
}

// Session authenticates using sessions
type Session struct {
	SessionKey string
}

// NewSession creates a new session authenticator
func NewSession(sessionKey string) *Session {
	return &Session{
		SessionKey: sessionKey,
	}
}

// Authenticate authenticates a request using session
func (s *Session) Authenticate(c *fiber.Ctx) (interface{}, error) {
	// Get user from session
	user := c.Locals(s.SessionKey)
	if user == nil {
		return nil, errors.New("not authenticated")
	}
	return user, nil
}

// APIKey authenticates using API keys
type APIKey struct {
	HeaderName string
	Validator  func(string) (interface{}, error)
}

// NewAPIKey creates a new API key authenticator
func NewAPIKey(validator func(string) (interface{}, error)) *APIKey {
	return &APIKey{
		HeaderName: "X-API-Key",
		Validator:  validator,
	}
}

// Authenticate authenticates a request using API key
func (a *APIKey) Authenticate(c *fiber.Ctx) (interface{}, error) {
	key := c.Get(a.HeaderName)
	if key == "" {
		return nil, errors.New("no API key provided")
	}
	
	return a.Validator(key)
}

// Middleware creates authentication middleware
func Middleware(authenticator Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := authenticator.Authenticate(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		
		c.Locals("user", user)
		return c.Next()
	}
}

// Optional creates optional authentication middleware
func Optional(authenticator Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, _ := authenticator.Authenticate(c)
		if user != nil {
			c.Locals("user", user)
		}
		return c.Next()
	}
}

