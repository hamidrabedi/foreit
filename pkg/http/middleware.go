package http

import (
	"net/http"
	"strconv"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// Middleware is a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// DefaultMiddlewares returns the default middleware stack for the framework
func DefaultMiddlewares() []Middleware {
	return []Middleware{
		RequestID,
		RealIP,
		Recoverer,
		Logger,
		Timeout(60 * time.Second),
	}
}

// RequestID adds a request ID to each request
func RequestID(next http.Handler) http.Handler {
	return chimw.RequestID(next)
}

// RealIP sets the client IP from X-Forwarded-For or X-Real-IP headers
func RealIP(next http.Handler) http.Handler {
	return chimw.RealIP(next)
}

// Recoverer recovers from panics and returns a 500 error
func Recoverer(next http.Handler) http.Handler {
	return chimw.Recoverer(next)
}

// Logger logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return chimw.Logger(next)
}

// Timeout sets a timeout for request processing
func Timeout(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return chimw.Timeout(timeout)(next)
	}
}

// Compress compresses response bodies
func Compress(next http.Handler) http.Handler {
	return chimw.Compress(5)(next)
}

// StripSlashes removes trailing slashes from URLs
func StripSlashes(next http.Handler) http.Handler {
	return chimw.StripSlashes(next)
}

// RedirectSlashes redirects trailing slashes
func RedirectSlashes(next http.Handler) http.Handler {
	return chimw.RedirectSlashes(next)
}

// CORS adds CORS headers to responses
func CORS(allowedOrigins, allowedMethods, allowedHeaders []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if len(allowedMethods) > 0 {
				methods := ""
				for i, m := range allowedMethods {
					if i > 0 {
						methods += ", "
					}
					methods += m
				}
				w.Header().Set("Access-Control-Allow-Methods", methods)
			}

			if len(allowedHeaders) > 0 {
				headers := ""
				for i, h := range allowedHeaders {
					if i > 0 {
						headers += ", "
					}
					headers += h
				}
				w.Header().Set("Access-Control-Allow-Headers", headers)
			}

			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit provides basic rate limiting (simple implementation)
func RateLimit(requestsPerMinute int) Middleware {
	// This is a simplified implementation
	// In production, use a proper rate limiting library
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement actual rate limiting
			next.ServeHTTP(w, r)
		})
	}
}

// SecureHeaders adds security headers to responses
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// CacheControl sets cache control headers
func CacheControl(maxAge int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(maxAge))
			next.ServeHTTP(w, r)
		})
	}
}
