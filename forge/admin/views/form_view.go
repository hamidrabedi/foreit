package views

import (
	"context"
	"fmt"
	"net/http"

	"github.com/forgego/forge/admin"
	adminschema "github.com/forgego/forge/admin/schema"
	"github.com/forgego/forge/schema"
)

// FormView represents a type-safe form view for create/update
type FormView[T any] struct {
	*BaseView[T]
}

// NewFormView creates a new form view
func NewFormView[T any](admin *admin.Admin[T]) *FormView[T] {
	return &FormView[T]{
		BaseView: NewBaseView(admin),
	}
}

// FormData contains data for rendering the form view
type FormData[T any] struct {
	Instance      *T
	IsNew         bool
	Fields        []FormFieldData
	Fieldsets     []FieldsetData
	Errors        map[string][]string
	HasAddPermission  bool
	HasChangePermission bool
	HasDeletePermission bool
}

// FormFieldData contains data for a form field
type FormFieldData struct {
	Name        string
	Label       string
	Value       interface{}
	Widget      admin.Widget
	HelpText    string
	Required    bool
	ReadOnly    bool
	Errors      []string
}

// FieldsetData contains data for a fieldset
type FieldsetData struct {
	Name        string
	Fields      []FormFieldData
	Collapsed   bool
	Description string
}

// Render renders the form view and returns the data
func (fv *FormView[T]) Render(ctx context.Context, r *http.Request, user interface{}, instance *T, isNew bool) (*FormData[T], error) {
	config := fv.admin.Config()

	// Check for view hooks
	if config != nil {
		if isNew && config.AddViewHook != nil {
			customView, err := config.AddViewHook(ctx, fv.admin, r)
			if err != nil {
				return nil, fmt.Errorf("add view hook error: %w", err)
			}
			if customView != nil {
				// Type assert to FormView[T]
				if fv, ok := customView.(*FormView[T]); ok {
					return fv.Render(ctx, r, user, instance, isNew)
				}
			}
		} else if !isNew && config.ChangeViewHook != nil {
			customView, err := config.ChangeViewHook(ctx, fv.admin, instance, r)
			if err != nil {
				return nil, fmt.Errorf("change view hook error: %w", err)
			}
			if customView != nil {
				// Type assert to FormView[T]
				if fv, ok := customView.(*FormView[T]); ok {
					return fv.Render(ctx, r, user, instance, isNew)
				}
			}
		}
	}

	// Get fields from config or auto-discover from schema
	fields := fv.getFields(config, instance, isNew)

	// Get fieldsets
	fieldsets := fv.getFieldsets(config, fields, instance, isNew)

	// Get read-only fields
	readOnlyFields := fv.getReadOnlyFields(config, instance, isNew)

	// Build form field data
	formFields := fv.buildFormFields(fields, readOnlyFields, instance)

	// Group fields into fieldsets
	fieldsetData := fv.groupFieldsIntoFieldsets(fieldsets, formFields)

	return &FormData[T]{
		Instance:            instance,
		IsNew:               isNew,
		Fields:              formFields,
		Fieldsets:           fieldsetData,
		Errors:              make(map[string][]string),
		HasAddPermission:    fv.admin.HasAddPermission(ctx, user),
		HasChangePermission: fv.admin.HasChangePermission(ctx, user, instance),
		HasDeletePermission: fv.admin.HasDeletePermission(ctx, user, instance),
	}, nil
}

// Save saves the form data to the model instance
func (fv *FormView[T]) Save(ctx context.Context, r *http.Request, instance *T, isNew bool) error {
	// Parse form data
	formData := admin.FormData{}
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	for key, values := range r.Form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	// Apply prepopulated fields
	prepopHandler := NewPrepopulatedFieldHandler(fv.admin)
	if err := prepopHandler.PopulateFields(ctx, instance, formData, isNew); err != nil {
		return fmt.Errorf("failed to populate fields: %w", err)
	}

	// Use admin's SaveModel method
	return fv.admin.SaveModel(ctx, instance, formData, isNew)
}

// getFields gets fields from config or auto-discovers from schema
func (fv *FormView[T]) getFields(config *admin.Config[T], instance *T, isNew bool) []interface{} {
	if config != nil && config.GetFields != nil {
		return config.GetFields(context.Background(), instance, isNew)
	}

	if config != nil && len(config.Fields) > 0 {
		return config.Fields
	}

	// Auto-discover from schema
	fields := fv.admin.Fields()
	fieldMapper := adminschema.NewFieldMapper()
	result := make([]interface{}, 0)
	for _, field := range fields {
		if fieldMapper.ShouldDisplayInForm(field.SchemaField) {
			result = append(result, field.Name)
		}
	}
	return result
}

// getFieldsets gets fieldsets from config
func (fv *FormView[T]) getFieldsets(config *admin.Config[T], fields []interface{}, instance *T, isNew bool) []admin.Fieldset[T] {
	if config != nil && config.GetFieldsets != nil {
		return config.GetFieldsets(context.Background(), instance, isNew)
	}

	if config != nil && len(config.Fieldsets) > 0 {
		return config.Fieldsets
	}

	// Default: single fieldset with all fields
	return []admin.Fieldset[T]{
		admin.NewFieldset[T]("", fields...),
	}
}

// getReadOnlyFields gets read-only fields from config
func (fv *FormView[T]) getReadOnlyFields(config *admin.Config[T], instance *T, isNew bool) map[string]bool {
	readOnlyMap := make(map[string]bool)

	if config != nil && config.GetReadOnlyFields != nil {
		readOnlyFields := config.GetReadOnlyFields(context.Background(), instance, isNew)
		for _, field := range readOnlyFields {
			if name := fv.getFieldName(field); name != "" {
				readOnlyMap[name] = true
			}
		}
		return readOnlyMap
	}

	if config != nil && len(config.ReadOnlyFields) > 0 {
		for _, field := range config.ReadOnlyFields {
			if name := fv.getFieldName(field); name != "" {
				readOnlyMap[name] = true
			}
		}
	}

	return readOnlyMap
}

// buildFormFields builds form field data from field names
func (fv *FormView[T]) buildFormFields(fields []interface{}, readOnlyFields map[string]bool, instance *T) []FormFieldData {
	result := make([]FormFieldData, 0, len(fields))

	for _, field := range fields {
		fieldName := fv.getFieldName(field)
		if fieldName == "" {
			continue
		}

		// Get field info from schema
		adminFields := fv.admin.Fields()
		var fieldInfo *adminschema.FieldInfo
		for i := range adminFields {
			if adminFields[i].Name == fieldName {
				fieldInfo = &adminFields[i]
				break
			}
		}

		if fieldInfo == nil {
			continue
		}

		// Get widget based on field configuration
		config := fv.admin.Config()
		widget := fv.getWidgetForField(config, fieldName, fieldInfo, instance)

		// Get field value from instance (would use ORM field accessor)
		var value interface{} = nil

		formField := FormFieldData{
			Name:     fieldName,
			Label:    fieldInfo.VerboseName,
			Value:    value,
			Widget:   widget,
			HelpText: fieldInfo.HelpText,
			Required: fieldInfo.Required,
			ReadOnly: readOnlyFields[fieldName] || fieldInfo.ReadOnly,
			Errors:   []string{},
		}

		result = append(result, formField)
	}

	return result
}

// groupFieldsIntoFieldsets groups form fields into fieldsets
func (fv *FormView[T]) groupFieldsIntoFieldsets(fieldsets []admin.Fieldset[T], formFields []FormFieldData) []FieldsetData {
	result := make([]FieldsetData, 0, len(fieldsets))

	// Create map of field name to form field
	fieldMap := make(map[string]FormFieldData)
	for _, field := range formFields {
		fieldMap[field.Name] = field
	}

	for _, fieldset := range fieldsets {
		fieldsetFields := make([]FormFieldData, 0)
		for _, field := range fieldset.Fields {
			fieldName := fv.getFieldName(field)
			if formField, ok := fieldMap[fieldName]; ok {
				fieldsetFields = append(fieldsetFields, formField)
			}
		}

		result = append(result, FieldsetData{
			Name:        fieldset.Name,
			Fields:      fieldsetFields,
			Collapsed:   fieldset.Collapsed,
			Description: fieldset.Description,
		})
	}

	return result
}

// getFieldName extracts field name from interface{} (string or FieldExpr)
func (fv *FormView[T]) getFieldName(field interface{}) string {
	if name, ok := field.(string); ok {
		return name
	}
	// If it's a FieldExpr, we'd need to get the name from it
	// FieldExpr is generic, so we can't type assert directly
	// For now, return empty string - would need reflection or different approach
	return ""
}

// getWidgetForField gets the appropriate widget for a field based on configuration
func (fv *FormView[T]) getWidgetForField(config *admin.Config[T], fieldName string, fieldInfo *adminschema.FieldInfo, instance *T) admin.Widget {
	// 1. Check FormFieldOverrides first (highest priority)
	if config != nil && config.FormFieldOverrides != nil {
		if widget, ok := config.FormFieldOverrides[fieldName]; ok {
			return widget
		}
	}

	// 2. Check RawIDFields
	if config != nil && fv.isInFieldList(config.RawIDFields, fieldName) {
		return admin.NewRawIDWidget()
	}

	// 3. Check AutocompleteFields
	if config != nil && fv.isInFieldList(config.AutocompleteFields, fieldName) {
		// Use Select widget for autocomplete (will be enhanced with JS)
		return admin.NewSelect([]admin.Choice[interface{}]{})
	}

	// 4. Check custom widget hooks based on field type
	if config != nil {
		ctx := context.Background()
		// Check if it's a foreign key
		if fieldInfo.Type == schema.TypeForeignKey && config.FormFieldForForeignKey != nil {
			return config.FormFieldForForeignKey(ctx, fv.admin, fieldName, instance)
		}
		// Check if it's a many-to-many
		if fieldInfo.Type == schema.TypeManyToMany && config.FormFieldForManyToMany != nil {
			return config.FormFieldForManyToMany(ctx, fv.admin, fieldName, instance)
		}
		// Check if it's a regular DB field
		if config.FormFieldForDBField != nil {
			return config.FormFieldForDBField(ctx, fv.admin, fieldName, instance)
		}
	}

	// 5. Check RadioFields
	if config != nil && config.RadioFields != nil {
		if _, ok := config.RadioFields[fieldName]; ok {
			// Return radio widget (would need to implement)
			// For now, use select as fallback
			return admin.NewSelect(nil)
		}
	}

	// 6. Use default widget based on field type
	// For now, return text input as default
	// In full implementation, would use widget registry
	return admin.NewTextInput()
}

// isInFieldList checks if a field name is in a list of fields (string or FieldExpr)
func (fv *FormView[T]) isInFieldList(fieldList []interface{}, fieldName string) bool {
	for _, field := range fieldList {
		if fv.getFieldName(field) == fieldName {
			return true
		}
	}
	return false
}

