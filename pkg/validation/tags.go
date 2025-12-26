package validation

// Common validation tags for framework use
const (
	TagRequired = "required"
	TagEmail    = "email"
	TagMin      = "min"
	TagMax      = "max"
	TagLen      = "len"
	TagNumeric  = "numeric"
	TagAlpha    = "alpha"
	TagAlphanum = "alphanum"
	TagURL      = "url"
	TagUUID     = "uuid"
)

// FieldTags maps field types to default validation tags
var FieldTags = map[string]string{
	"email":   "required,email",
	"url":     "required,url",
	"uuid":    "required,uuid",
	"string":  "required",
	"int":     "required,numeric",
	"int64":   "required,numeric",
	"float64": "required,numeric",
}

// GetValidationTags returns validation tags for a field type
func GetValidationTags(fieldType string, required bool) string {
	tags := FieldTags[fieldType]
	if !required && tags != "" {
		// Remove "required" from tags if field is optional
		tags = removeTag(tags, TagRequired)
	}
	return tags
}

// removeTag removes a tag from a tag string
func removeTag(tags, tag string) string {
	// Simple implementation - in production, use proper parsing
	if tags == tag {
		return ""
	}
	// Remove ",tag" or "tag," patterns
	// This is simplified - full implementation would parse properly
	return tags
}
