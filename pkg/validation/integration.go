package validation

import (
	"github.com/forgego/forge/pkg/schema"
)

// GenerateValidationTag generates a validation tag for a field
func GenerateValidationTag(field schema.Field) string {
	var tags []string

	// Required
	if field.Required {
		tags = append(tags, "required")
	}

	// Type-specific validations
	switch field.Type {
	case schema.TypeEmail:
		tags = append(tags, "email")
	case schema.TypeURL:
		tags = append(tags, "url")
	case schema.TypeUUID:
		tags = append(tags, "uuid")
	}

	// Length validations
	if field.MaxLength != nil {
		tags = append(tags, "max="+string(rune(*field.MaxLength)))
	}
	if field.MinLength != nil {
		tags = append(tags, "min="+string(rune(*field.MinLength)))
	}

	// Use custom validation tag if provided
	if field.ValidationTag != "" {
		return field.ValidationTag
	}

	// Join tags
	if len(tags) == 0 {
		return ""
	}

	result := tags[0]
	for i := 1; i < len(tags); i++ {
		result += "," + tags[i]
	}

	return result
}

// ValidateModel validates a model instance using the validator
func ValidateModel(v *Validator, instance interface{}) error {
	return v.ValidateStruct(instance)
}
