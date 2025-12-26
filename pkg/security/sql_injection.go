package security

import (
	"fmt"
	"regexp"
	"strings"
)

// SQLInjection provides SQL injection prevention
type SQLInjection struct{}

// NewSQLInjection creates a new SQL injection protector
func NewSQLInjection() *SQLInjection {
	return &SQLInjection{}
}

// ValidateInput validates input to prevent SQL injection
func (s *SQLInjection) ValidateInput(input string) error {
	// Check for SQL injection patterns
	dangerousPatterns := []string{
		";",
		"--",
		"/*",
		"*/",
		"xp_",
		"sp_",
		"exec",
		"execute",
		"union",
		"select",
		"insert",
		"update",
		"delete",
		"drop",
		"create",
		"alter",
		"truncate",
	}

	inputLower := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(inputLower, pattern) {
			return fmt.Errorf("potentially dangerous input detected: %s", pattern)
		}
	}

	return nil
}

// SanitizeIdentifier sanitizes SQL identifiers (table/column names)
func (s *SQLInjection) SanitizeIdentifier(identifier string) string {
	// Only allow alphanumeric and underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return reg.ReplaceAllString(identifier, "")
}

// EnsureParameterized ensures a query uses parameterized queries
func (s *SQLInjection) EnsureParameterized(query string) error {
	// Check for string concatenation in SQL
	if strings.Contains(query, "+") || strings.Contains(query, "||") {
		return fmt.Errorf("query appears to use string concatenation - use parameterized queries instead")
	}

	// Check for fmt.Sprintf patterns
	if strings.Contains(query, "%s") || strings.Contains(query, "%d") {
		return fmt.Errorf("query appears to use string formatting - use parameterized queries instead")
	}

	return nil
}

// LogQuery logs a query for security auditing
func (s *SQLInjection) LogQuery(query string, args []interface{}) {
	// TODO: Implement query logging for security auditing
	// This should log to a secure audit log
}
