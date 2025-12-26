package http

import (
	"encoding/json"
	"net/http"
)

// Response provides helper methods for HTTP responses
type Response struct {
	http.ResponseWriter
	statusCode int
}

// NewResponse wraps an http.ResponseWriter with helper methods
func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// Status sets the HTTP status code
func (r *Response) Status(code int) *Response {
	r.statusCode = code
	return r
}

// JSON sends a JSON response
func (r *Response) JSON(data interface{}) error {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(r.statusCode)
	return json.NewEncoder(r).Encode(data)
}

// JSONError sends a JSON error response
func (r *Response) JSONError(message string, code int) error {
	return r.Status(code).JSON(map[string]interface{}{
		"error":   true,
		"message": message,
		"code":    code,
	})
}

// Text sends a plain text response
func (r *Response) Text(text string) error {
	r.Header().Set("Content-Type", "text/plain")
	r.WriteHeader(r.statusCode)
	_, err := r.Write([]byte(text))
	return err
}

// HTML sends an HTML response
func (r *Response) HTML(html string) error {
	r.Header().Set("Content-Type", "text/html")
	r.WriteHeader(r.statusCode)
	_, err := r.Write([]byte(html))
	return err
}

// Redirect sends a redirect response
func (r *Response) Redirect(url string, code int) {
	http.Redirect(r, r.Request(), url, code)
}

// Cookie sets a cookie
func (r *Response) Cookie(cookie *http.Cookie) {
	http.SetCookie(r, cookie)
}

// Header sets a response header
func (r *Response) SetHeader(key, value string) {
	r.Header().Set(key, value)
}

// WriteHeader writes the status code
func (r *Response) WriteHeader(code int) {
	if r.statusCode == 0 {
		r.statusCode = code
	}
	r.ResponseWriter.WriteHeader(r.statusCode)
}

// Request returns the associated request (if available)
func (r *Response) Request() *http.Request {
	// In a real implementation, we'd store the request
	// For now, this is a placeholder
	return nil
}
