package http

import (
	"context"
	"net/http"
)

type contextKey string

const (
	// UserKey is the context key for the authenticated user
	UserKey contextKey = "user"
	// DBKey is the context key for the database connection
	DBKey contextKey = "db"
	// LoggerKey is the context key for the logger
	LoggerKey contextKey = "logger"
	// LocaleKey is the context key for the locale
	LocaleKey contextKey = "locale"
	// RequestIDKey is the context key for the request ID
	RequestIDKey contextKey = "request_id"
)

// GetUser retrieves the user from the request context
func GetUser(r *http.Request) interface{} {
	return r.Context().Value(UserKey)
}

// SetUser sets the user in the request context
func SetUser(r *http.Request, user interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), UserKey, user))
}

// GetDB retrieves the database connection from the request context
func GetDB(r *http.Request) interface{} {
	return r.Context().Value(DBKey)
}

// SetDB sets the database connection in the request context
func SetDB(r *http.Request, db interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), DBKey, db))
}

// GetLogger retrieves the logger from the request context
func GetLogger(r *http.Request) interface{} {
	return r.Context().Value(LoggerKey)
}

// SetLogger sets the logger in the request context
func SetLogger(r *http.Request, logger interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), LoggerKey, logger))
}

// GetLocale retrieves the locale from the request context
func GetLocale(r *http.Request) string {
	if locale, ok := r.Context().Value(LocaleKey).(string); ok {
		return locale
	}
	return "en"
}

// SetLocale sets the locale in the request context
func SetLocale(r *http.Request, locale string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), LocaleKey, locale))
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// SetRequestID sets the request ID in the context
func SetRequestID(r *http.Request, id string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), RequestIDKey, id))
}

// WithContext returns a new request with the given context
func WithContext(r *http.Request, ctx context.Context) *http.Request {
	return r.WithContext(ctx)
}

// GetContext retrieves the request context
func GetContext(r *http.Request) context.Context {
	return r.Context()
}
