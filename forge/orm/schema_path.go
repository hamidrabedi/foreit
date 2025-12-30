package orm

import (
	"fmt"
	"strings"
)

// ValidatePath validates a deep relation path (e.g., "author__company__country")
// Returns an error with helpful suggestions if the path is invalid
func (ms *ModelSchema) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("field path cannot be empty")
	}

	parts := strings.Split(path, "__")
	if len(parts) == 0 {
		return fmt.Errorf("invalid field path: %s", path)
	}

	currentSchema := ms

	for i, part := range parts {
		if part == "" {
			return fmt.Errorf("empty part in field path: %s", path)
		}

		if i == len(parts)-1 {
			// Last part is a field
			field := currentSchema.GetField(part)
			if field == nil {
				// Provide helpful suggestions
				allFields := make([]string, len(currentSchema.Fields))
				for j, f := range currentSchema.Fields {
					allFields[j] = f.Name
				}
				return fmt.Errorf("field '%s' not found on %s. Did you mean: %v?",
					part, currentSchema.TableName, allFields)
			}
		} else {
			// Intermediate parts are relations
			rel := currentSchema.GetRelation(part)
			if rel == nil {
				// Provide helpful suggestions
				allRelations := make([]string, len(currentSchema.Relations))
				for j, r := range currentSchema.Relations {
					allRelations[j] = r.Name
				}
				return fmt.Errorf("relation '%s' not found on %s. Did you mean: %v?",
					part, currentSchema.TableName, allRelations)
			}

			// Resolve target model schema
			targetSchema, err := GetModelSchemaByName(rel.TargetModel)
			if err != nil {
				return fmt.Errorf("failed to resolve relation '%s': %w", part, err)
			}
			currentSchema = targetSchema
		}
	}

	return nil
}

// ResolvePath resolves a field path and returns the final field info and schema
func (ms *ModelSchema) ResolvePath(path string) (*FieldInfo, *ModelSchema, error) {
	if err := ms.ValidatePath(path); err != nil {
		return nil, nil, err
	}

	parts := strings.Split(path, "__")
	currentSchema := ms

	// Navigate through relations
	for i := 0; i < len(parts)-1; i++ {
		rel := currentSchema.GetRelation(parts[i])
		if rel == nil {
			return nil, nil, fmt.Errorf("relation '%s' not found", parts[i])
		}

		targetSchema, err := GetModelSchemaByName(rel.TargetModel)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to resolve relation '%s': %w", parts[i], err)
		}
		currentSchema = targetSchema
	}

	// Get the final field
	fieldName := parts[len(parts)-1]
	field := currentSchema.GetField(fieldName)
	if field == nil {
		return nil, nil, fmt.Errorf("field '%s' not found", fieldName)
	}

	return field, currentSchema, nil
}

// GetAllowedFields returns allowed fields for a role (for security whitelisting)
// This can be extended with schema metadata/tags for role-based filtering
func (ms *ModelSchema) GetAllowedFields(role string) []string {
	// For now, return all fields
	// This can be extended to check schema metadata or tags
	fields := make([]string, len(ms.Fields))
	for i, f := range ms.Fields {
		fields[i] = f.Name
	}
	return fields
}

// GetAllowedLookups returns allowed lookups for a field (for security whitelisting)
func (ms *ModelSchema) GetAllowedLookups(fieldPath string) ([]string, error) {
	field, _, err := ms.ResolvePath(fieldPath)
	if err != nil {
		return nil, err
	}

	// Return default lookups based on field type
	// This can be extended with schema metadata
	lookups := getDefaultLookupsForType(field.Type)
	return lookups, nil
}

// getDefaultLookupsForType returns default lookups for a field type
func getDefaultLookupsForType(fieldType interface{}) []string {
	// This is a simplified version
	// In practice, this would check the actual reflect.Type
	return []string{"exact", "in"}
}

// GetModelSchemaByName is now implemented in schema_registry.go
// This function is kept for backward compatibility but delegates to the registry

// GetRelationDepth calculates the depth of a relation path
func (ms *ModelSchema) GetRelationDepth(path string) (int, error) {
	parts := strings.Split(path, "__")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid path")
	}

	// Count relations (all parts except the last)
	depth := len(parts) - 1
	return depth, nil
}
