package filter

import (
	"fmt"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/filter/filters"
	"github.com/forgego/forge/schema"
)

// BuildFiltersFromSchema builds filters automatically from schema fields
func BuildFiltersFromSchema[T any](
	schemaInstance schema.Schema,
	filterset *filter.FilterSet[T],
) error {
	fields := schemaInstance.Fields()

	for _, field := range fields {
		adminFilter := createFilterForField[T](field)
		if adminFilter != nil {
			filterset.AddFilter(field.Name, adminFilter)
		}
	}

	return nil
}

// createFilterForField creates an appropriate filter for a schema field
func createFilterForField[T any](field schema.Field) filter.Filter[T] {
	switch field.Type {
	case schema.TypeBool:
		// Boolean filter
		boolFilter := filters.NewBooleanFilter[T](field.Name)
		return boolFilter

	case schema.TypeString, schema.TypeText, schema.TypeEmail, schema.TypeURL:
		// Char/text filter
		charFilter := filters.NewCharFilter[T](field.Name)
		return charFilter

	case schema.TypeDateTime, schema.TypeDate, schema.TypeTime:
		// Date filter
		dateFilter := filters.NewDateFilter[T](field.Name)
		return dateFilter

	case schema.TypeInt64, schema.TypeInt32:
		// Number filter
		numberFilter := filters.NewNumberFilter[T](field.Name)
		return numberFilter

	case schema.TypeFloat32, schema.TypeFloat64, schema.TypeDecimal:
		// Number filter for floats
		numberFilter := filters.NewNumberFilter[T](field.Name)
		return numberFilter

	case schema.TypeForeignKey, schema.TypeOneToOne:
		// Related filter
		// This would need the target model type, which we don't have here
		// For now, return a choice filter if choices are available
		if len(field.Choices) > 0 {
			choices := convertSchemaChoicesToFilterChoices(field.Choices)
			choiceFilter := filters.NewChoiceFilter[T](field.Name, choices)
			return choiceFilter
		}
		// For ForeignKey/OneToOne fields without choices, create a related filter
		// The filter will work with the foreign key ID value
		// In a full implementation, this would query the related model for options
		// For now, we create a number filter that accepts the related object's ID
		relatedFilter := filters.NewNumberFilter[T](field.Name)
		return relatedFilter

	case schema.TypeManyToMany:
		// ManyToMany filter - similar to ForeignKey
		if len(field.Choices) > 0 {
			choices := convertSchemaChoicesToFilterChoices(field.Choices)
			choiceFilter := filters.NewChoiceFilter[T](field.Name, choices)
			return choiceFilter
		}
		return nil

	default:
		// Unknown type - skip
		return nil
	}
}

// BuildFiltersFromConfig builds filters based on admin configuration
func BuildFiltersFromConfig[T any](
	config *FilterConfig[T],
	schemaInstance schema.Schema,
	filterset *filter.FilterSet[T],
) error {
	// If no enabled filters specified, build all from schema
	if len(config.EnabledFilters) == 0 {
		return BuildFiltersFromSchema[T](schemaInstance, filterset)
	}

	// Build only specified filters
	fields := schemaInstance.Fields()
	fieldMap := make(map[string]schema.Field)
	for _, field := range fields {
		fieldMap[field.Name] = field
	}

	for _, filterName := range config.EnabledFilters {
		field, ok := fieldMap[filterName]
		if !ok {
			return fmt.Errorf("field %s not found in schema", filterName)
		}

		adminFilter := createFilterForField[T](field)
		if adminFilter != nil {
			filterset.AddFilter(filterName, adminFilter)
		}
	}

	return nil
}

// convertSchemaChoicesToFilterChoices converts schema.Choice to filter.Choice
func convertSchemaChoicesToFilterChoices(schemaChoices []schema.Choice) []filters.Choice {
	choices := make([]filters.Choice, len(schemaChoices))
	for i, sc := range schemaChoices {
		choices[i] = filters.Choice{
			Label: sc.Label,
			Value: sc.Value,
		}
	}
	return choices
}
