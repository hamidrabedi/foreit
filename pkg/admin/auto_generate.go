package admin

import (
	"github.com/forgego/forge/pkg/schema"
)

// autoGenerateAdminFromSchema auto-generates admin config from schema
func autoGenerateAdminFromSchema(adminModel *AdminModel, schemaInstance schema.Schema) *AdminModel {
	// Get fields from schema
	fields := schemaInstance.Fields()
	meta := schemaInstance.Meta()

	// Auto-generate list display
	if len(adminModel.ListDisplay) == 0 {
		adminModel.ListDisplay = autoGenerateListDisplay(fields, meta)
	}

	// Auto-generate list filter
	if len(adminModel.ListFilter) == 0 {
		adminModel.ListFilter = autoGenerateListFilter(fields)
	}

	// Auto-generate search fields
	if len(adminModel.SearchFields) == 0 {
		adminModel.SearchFields = autoGenerateSearchFields(fields)
	}

	// Auto-generate read-only fields
	if len(adminModel.ReadOnlyFields) == 0 {
		adminModel.ReadOnlyFields = autoGenerateReadOnlyFields(fields)
	}

	return adminModel
}

// autoGenerateListDisplay generates list display fields
func autoGenerateListDisplay(fields []schema.Field, meta schema.Meta) []interface{} {
	var listDisplay []interface{}

	// Use Meta.OrderBy if available, otherwise use first few fields
	if len(meta.OrderBy) > 0 {
		for _, fieldName := range meta.OrderBy {
			// Remove leading dash for descending order
			if len(fieldName) > 0 && fieldName[0] == '-' {
				fieldName = fieldName[1:]
			}
			listDisplay = append(listDisplay, fieldName)
		}
	} else {
		// Default: use first 5 non-primary key fields
		count := 0
		for _, field := range fields {
			if !field.PrimaryKey && count < 5 {
				listDisplay = append(listDisplay, field.Name)
				count++
			}
		}
	}

	// Always include ID if not already included
	hasID := false
	for _, item := range listDisplay {
		if name, ok := item.(string); ok && (name == "id" || name == "ID") {
			hasID = true
			break
		}
	}
	if !hasID {
		// Prepend ID
		listDisplay = append([]interface{}{"id"}, listDisplay...)
	}

	return listDisplay
}

// autoGenerateListFilter generates list filter fields
func autoGenerateListFilter(fields []schema.Field) []interface{} {
	var filters []interface{}

	for _, field := range fields {
		// Add boolean fields
		if field.Type == schema.TypeBool {
			filters = append(filters, field.Name)
		}
		// Add foreign key fields (would need relation info)
		// For now, skip
	}

	return filters
}

// autoGenerateSearchFields generates search fields
func autoGenerateSearchFields(fields []schema.Field) []interface{} {
	var searchFields []interface{}

	for _, field := range fields {
		// Add string and text fields
		if field.Type == schema.TypeString || field.Type == schema.TypeText {
			searchFields = append(searchFields, field.Name)
		}
	}

	return searchFields
}

// autoGenerateReadOnlyFields generates read-only fields
func autoGenerateReadOnlyFields(fields []schema.Field) []interface{} {
	var readOnly []interface{}

	for _, field := range fields {
		// Primary keys are read-only
		if field.PrimaryKey {
			readOnly = append(readOnly, field.Name)
		}
		// Auto-now and auto-now-add fields are read-only
		if field.AutoNow || field.AutoNowAdd {
			readOnly = append(readOnly, field.Name)
		}
	}

	return readOnly
}

