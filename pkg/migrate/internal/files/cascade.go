package files

import "strings"

// NormalizeCascadeAction normalizes cascade action strings to SQL format
func NormalizeCascadeAction(action string) string {
	action = strings.ToUpper(strings.TrimSpace(action))
	switch action {
	case "CASCADE", "RESTRICT", "SET NULL", "NO ACTION":
		return action
	case "PROTECT":
		return "RESTRICT"
	default:
		return "NO ACTION"
	}
}

// DenormalizeCascadeAction converts SQL cascade action back to model format
func DenormalizeCascadeAction(action string) string {
	action = strings.ToUpper(strings.TrimSpace(action))
	switch action {
	case "RESTRICT":
		return "PROTECT"
	case "CASCADE", "SET NULL", "NO ACTION":
		return action
	default:
		return "NO ACTION"
	}
}
