package detection

import (
	"fmt"
	"strings"

	"github.com/forgego/forge/pkg/generator"
)

// getTableName gets the table name from a model definition
func getTableName(def *generator.ModelDefinition) string {
	if def.Meta.TableName != "" {
		return def.Meta.TableName
	}
	return fmt.Sprintf("%ss", toSnakeCase(def.Name))
}

// findTargetTable finds the target table name for a relation
func findTargetTable(targetModel string, allDefs []*generator.ModelDefinition) string {
	for _, def := range allDefs {
		if strings.EqualFold(def.Name, targetModel) {
			return getTableName(def)
		}
	}
	return ""
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return string(result)
}

