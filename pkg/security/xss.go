package security

import (
	"html"
	"strings"
)

// Note: Template functions from sprig are available via internal/admin/templates/filters.go
// This file provides low-level XSS protection utilities

// XSS provides XSS protection
type XSS struct{}

// NewXSS creates a new XSS protector
func NewXSS() *XSS {
	return &XSS{}
}

// EscapeHTML escapes HTML special characters
func (x *XSS) EscapeHTML(s string) string {
	return html.EscapeString(s)
}

// SanitizeHTML sanitizes HTML content
func (x *XSS) SanitizeHTML(html string) string {
	// TODO: Implement proper HTML sanitization
	// For now, just escape everything
	return x.EscapeHTML(html)
}

// SafeString represents a string that is safe to output without escaping
type SafeString string

// String returns the string value
func (s SafeString) String() string {
	return string(s)
}

// MarkSafe marks a string as safe (trusted content)
func MarkSafe(s string) SafeString {
	return SafeString(s)
}

// ContentSecurityPolicy generates CSP headers
func (x *XSS) ContentSecurityPolicy() map[string]string {
	return map[string]string{
		"Content-Security-Policy": "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'",
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"X-XSS-Protection":        "1; mode=block",
	}
}

// SanitizeInput sanitizes user input
func (x *XSS) SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Escape HTML
	input = x.EscapeHTML(input)

	return input
}
