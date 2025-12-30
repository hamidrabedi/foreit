package admin

import (
	"context"
	"fmt"
	"strings"
)

// PrepopulatedFieldHandler handles auto-population of fields
type PrepopulatedFieldHandler[T any] struct {
	admin *Admin[T]
}

// NewPrepopulatedFieldHandler creates a new prepopulated field handler
func NewPrepopulatedFieldHandler[T any](admin *Admin[T]) *PrepopulatedFieldHandler[T] {
	return &PrepopulatedFieldHandler[T]{
		admin: admin,
	}
}

// PopulateFields populates fields based on prepopulated_fields configuration
func (h *PrepopulatedFieldHandler[T]) PopulateFields(ctx context.Context, instance *T, formData FormData) error {
	config := h.admin.Config()
	if config == nil {
		return nil
	}

	// Get prepopulated fields (could be dynamic)
	prepopulatedFields := config.PrepopulatedFields
	if config.GetPrepopulatedFields != nil {
		prepopulatedFields = config.GetPrepopulatedFields(ctx, instance, false)
	}

	if len(prepopulatedFields) == 0 {
		return nil
	}

	// For each target field, populate from source fields
	for targetField, sourceFields := range prepopulatedFields {
		if len(sourceFields) == 0 {
			continue
		}

		// Get values from source fields
		var values []string
		for _, sourceField := range sourceFields {
			if val, ok := formData[sourceField]; ok {
				if strVal, ok := val.(string); ok {
					values = append(values, strVal)
				} else {
					values = append(values, fmt.Sprintf("%v", val))
				}
			}
		}

		// Combine source values (typically for slug generation)
		populatedValue := h.combineValues(values)

		// Set the target field value
		formData[targetField] = populatedValue
	}

	return nil
}

// combineValues combines source field values into a single value
// This is typically used for slug generation (e.g., "my-article-title")
func (h *PrepopulatedFieldHandler[T]) combineValues(values []string) string {
	// Join values with spaces, then slugify
	combined := strings.Join(values, " ")
	return slugify(combined)
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Replace spaces and special characters with hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		} else if r == ' ' || r == '_' {
			// Replace spaces and underscores with hyphens
			if result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
		}
		// Skip other special characters
	}

	// Remove leading/trailing hyphens and collapse multiple hyphens
	resultStr := result.String()
	resultStr = strings.Trim(resultStr, "-")
	
	// Collapse multiple consecutive hyphens
	for strings.Contains(resultStr, "--") {
		resultStr = strings.ReplaceAll(resultStr, "--", "-")
	}

	return resultStr
}
