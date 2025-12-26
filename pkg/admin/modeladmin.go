package admin

import (
	"context"
	"reflect"
	"time"

	"github.com/forgego/forge/pkg/query"
)

// ModelAdminConfig provides comprehensive Django-style ModelAdmin configuration
type ModelAdminConfig struct {
	// List display
	ListDisplay   []interface{} // Fields to display in list view
	ListFilter    []interface{} // Fields to filter by
	SearchFields  []interface{} // Fields to search in
	DateHierarchy string        // Date field for grouping (e.g., "created_at")
	Ordering      []string      // Default ordering (e.g., ["-created_at", "name"])
	ListPerPage   int           // Items per page (default: 20)
	
	// Form configuration
	Fieldsets     []Fieldset   // Group fields in forms
	ReadOnlyFields []interface{} // Fields that are read-only
	AutocompleteFields []string // Foreign key fields with autocomplete
	RawIDFields   []string      // Foreign key fields shown as raw ID
	
	// Actions
	Actions       []AdminAction // Custom admin actions
	
	// Customization
	VerboseName   string        // Human-readable name
	VerboseNamePlural string    // Plural name
	SaveOnTop     bool          // Show save buttons on top
	SaveAs        bool          // Show "Save as new" button
	SaveAsContinue bool         // Show "Save and continue editing" button
	SaveAndAddAnother bool      // Show "Save and add another" button
	
	// Custom methods
	GetQueryset   func(ctx context.Context, manager interface{}) (interface{}, error) // Custom queryset
	SaveModel     func(ctx context.Context, instance interface{}, form map[string]interface{}, isNew bool) error // Custom save
	DeleteModel   func(ctx context.Context, instance interface{}) error // Custom delete
	
	// Export
	ExportFormats []string // Supported export formats (e.g., ["csv", "json"])
}

// Fieldset groups form fields together
type Fieldset struct {
	Name      string      // Fieldset name/title
	Fields    []string    // Fields in this fieldset
	Collapsed bool        // Whether fieldset is collapsed by default
	Classes   []string    // CSS classes
}

// AdminAction represents a custom admin action
type AdminAction struct {
	Name        string
	Label       string
	Description string
	Handler     func(ctx context.Context, instances []interface{}) error
	Permissions []string // Required permissions
}

// ApplyModelAdminConfig applies ModelAdminConfig to AdminModel
func ApplyModelAdminConfig(model *AdminModel, config ModelAdminConfig) {
	if len(config.ListDisplay) > 0 {
		model.ListDisplay = config.ListDisplay
	}
	if len(config.ListFilter) > 0 {
		model.ListFilter = config.ListFilter
	}
	if len(config.SearchFields) > 0 {
		model.SearchFields = config.SearchFields
	}
	if len(config.ReadOnlyFields) > 0 {
		model.ReadOnlyFields = config.ReadOnlyFields
	}
	
	// Store extended config in a map for later retrieval
	if model.ExtendedConfig == nil {
		model.ExtendedConfig = make(map[string]interface{})
	}
	model.ExtendedConfig["date_hierarchy"] = config.DateHierarchy
	model.ExtendedConfig["ordering"] = config.Ordering
	model.ExtendedConfig["list_per_page"] = config.ListPerPage
	model.ExtendedConfig["fieldsets"] = config.Fieldsets
	model.ExtendedConfig["autocomplete_fields"] = config.AutocompleteFields
	model.ExtendedConfig["raw_id_fields"] = config.RawIDFields
	model.ExtendedConfig["actions"] = config.Actions
	model.ExtendedConfig["verbose_name"] = config.VerboseName
	model.ExtendedConfig["verbose_name_plural"] = config.VerboseNamePlural
	model.ExtendedConfig["save_on_top"] = config.SaveOnTop
	model.ExtendedConfig["save_as"] = config.SaveAs
	model.ExtendedConfig["save_as_continue"] = config.SaveAsContinue
	model.ExtendedConfig["save_and_add_another"] = config.SaveAndAddAnother
	model.ExtendedConfig["get_queryset"] = config.GetQueryset
	model.ExtendedConfig["save_model"] = config.SaveModel
	model.ExtendedConfig["delete_model"] = config.DeleteModel
	model.ExtendedConfig["export_formats"] = config.ExportFormats
}

// GetDateHierarchy returns the date hierarchy field for a model
func GetDateHierarchy(model *AdminModel) string {
	if model.ExtendedConfig == nil {
		return ""
	}
	if dh, ok := model.ExtendedConfig["date_hierarchy"].(string); ok {
		return dh
	}
	return ""
}

// GetOrdering returns the default ordering for a model
func GetOrdering(model *AdminModel) []string {
	if model.ExtendedConfig == nil {
		return nil
	}
	if ordering, ok := model.ExtendedConfig["ordering"].([]string); ok {
		return ordering
	}
	return nil
}

// GetListPerPage returns items per page for a model
func GetListPerPage(model *AdminModel) int {
	if model.ExtendedConfig == nil {
		return 20 // Default
	}
	if lpp, ok := model.ExtendedConfig["list_per_page"].(int); ok && lpp > 0 {
		return lpp
	}
	return 20
}

// GetFieldsets returns fieldsets for a model
func GetFieldsets(model *AdminModel) []Fieldset {
	if model.ExtendedConfig == nil {
		return nil
	}
	if fieldsets, ok := model.ExtendedConfig["fieldsets"].([]Fieldset); ok {
		return fieldsets
	}
	return nil
}

// GetAutocompleteFields returns autocomplete fields for a model
func GetAutocompleteFields(model *AdminModel) []string {
	if model.ExtendedConfig == nil {
		return nil
	}
	if fields, ok := model.ExtendedConfig["autocomplete_fields"].([]string); ok {
		return fields
	}
	return nil
}

// GetRawIDFields returns raw ID fields for a model
func GetRawIDFields(model *AdminModel) []string {
	if model.ExtendedConfig == nil {
		return nil
	}
	if fields, ok := model.ExtendedConfig["raw_id_fields"].([]string); ok {
		return fields
	}
	return nil
}

// GetActions returns custom actions for a model
func GetActions(model *AdminModel) []AdminAction {
	if model.ExtendedConfig == nil {
		return nil
	}
	if actions, ok := model.ExtendedConfig["actions"].([]AdminAction); ok {
		return actions
	}
	return nil
}

// GetVerboseName returns the verbose name for a model
func GetVerboseName(model *AdminModel) string {
	if model.ExtendedConfig == nil {
		return model.Name
	}
	if vn, ok := model.ExtendedConfig["verbose_name"].(string); ok && vn != "" {
		return vn
	}
	return model.Name
}

// GetVerboseNamePlural returns the plural verbose name for a model
func GetVerboseNamePlural(model *AdminModel) string {
	if model.ExtendedConfig == nil {
		return model.Name + "s"
	}
	if vnp, ok := model.ExtendedConfig["verbose_name_plural"].(string); ok && vnp != "" {
		return vnp
	}
	return model.Name + "s"
}

// ApplyOrdering applies ordering to a queryset
func ApplyOrdering(qs reflect.Value, ordering []string) reflect.Value {
	if len(ordering) == 0 {
		return qs
	}
	
	orderByMethod := qs.MethodByName("OrderBy")
	if !orderByMethod.IsValid() {
		return qs
	}
	
	// Apply ordering
	for _, orderField := range ordering {
		reverse := false
		fieldName := orderField
		if len(fieldName) > 0 && fieldName[0] == '-' {
			reverse = true
			fieldName = fieldName[1:]
		}
		
		results := orderByMethod.Call([]reflect.Value{reflect.ValueOf(fieldName)})
		if len(results) > 0 {
			qs = results[0]
		}
		
		if reverse {
			reverseMethod := qs.MethodByName("Reverse")
			if reverseMethod.IsValid() {
				results := reverseMethod.Call([]reflect.Value{})
				if len(results) > 0 {
					qs = results[0]
				}
			}
		}
	}
	
	return qs
}

// ApplyDateHierarchy applies date hierarchy filtering
func ApplyDateHierarchy(qs reflect.Value, dateField string, year, month, day int) reflect.Value {
	if dateField == "" {
		return qs
	}
	
	filterMethod := qs.MethodByName("Filter")
	if !filterMethod.IsValid() {
		return qs
	}
	
	// Build date filter
	if year > 0 {
		// Filter by year
		yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		yearEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		
		// Create date range filter
		// This is simplified - full implementation would use proper date range queries
		yearExpr := query.NewFieldQueryExpr(dateField, query.OpGreaterOrEqual, yearStart)
		yearEndExpr := query.NewFieldQueryExpr(dateField, query.OpLess, yearEnd)
		combinedExpr := yearExpr.And(yearEndExpr)
		
		results := filterMethod.Call([]reflect.Value{reflect.ValueOf(combinedExpr)})
		if len(results) > 0 {
			qs = results[0]
		}
		
		if month > 0 {
			monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			var monthEnd time.Time
			if month < 12 {
				monthEnd = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
			} else {
				monthEnd = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
			}
			
			monthExpr := query.NewFieldQueryExpr(dateField, query.OpGreaterOrEqual, monthStart)
			monthEndExpr := query.NewFieldQueryExpr(dateField, query.OpLess, monthEnd)
			combinedExpr := monthExpr.And(monthEndExpr)
			
			results := filterMethod.Call([]reflect.Value{reflect.ValueOf(combinedExpr)})
			if len(results) > 0 {
				qs = results[0]
			}
			
			if day > 0 {
				dayStart := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
				dayEnd := time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.UTC)
				
				dayExpr := query.NewFieldQueryExpr(dateField, query.OpGreaterOrEqual, dayStart)
				dayEndExpr := query.NewFieldQueryExpr(dateField, query.OpLess, dayEnd)
				combinedExpr := dayExpr.And(dayEndExpr)
				
				results := filterMethod.Call([]reflect.Value{reflect.ValueOf(combinedExpr)})
				if len(results) > 0 {
					qs = results[0]
				}
			}
		}
	}
	
	return qs
}

