package views

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"

	"github.com/forgego/forge/admin"
	adminschema "github.com/forgego/forge/admin/schema"
	adminutils "github.com/forgego/forge/admin/utils"
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
	Instance            *T
	IsNew               bool
	Fields              []FormFieldData
	Fieldsets           []FieldsetData
	Inlines             []InlineData
	Errors              map[string][]string
	HasAddPermission    bool
	HasViewPermission   bool
	HasChangePermission bool
	HasDeletePermission bool
	PrepopulatedFields  map[string][]string
}

// InlineData contains data for an inline formset
type InlineData struct {
	Name       string
	Title      string
	Style      string // "tabular" or "stacked"
	Fields     []InlineFieldData
	Rows       []InlineRowData
	ExtraForms []InlineRowData
	CanAddMore bool
	Extra      int
	MaxNum     int
}

// InlineRowData contains data for a single row in an inline formset
type InlineRowData struct {
	ID     interface{}
	Fields []FormFieldData
	Delete bool
}

// InlineFieldData contains data for a field in an inline
type InlineFieldData struct {
	Name  string
	Label string
}

// FormFieldData contains data for a form field
type FormFieldData struct {
	Name     string
	Label    string
	Value    interface{}
	Widget   admin.Widget
	HelpText string
	Required bool
	ReadOnly bool
	Errors   []string
}

// Render renders the widget for the field
func (f FormFieldData) Render() template.HTML {
	if f.Widget == nil {
		return ""
	}
	return f.Widget.Render(f.Name, f.Value, nil)
}

// FieldsetData contains data for a fieldset
type FieldsetData struct {
	Name        string
	Fields      []FormFieldData
	Collapsed   bool
	Description string
}

// Render renders the form view and returns the data
func (fv *FormView[T]) Render(ctx context.Context, r *http.Request, user interface{}, instance *T, isNew bool, errors map[string][]string) (*FormData[T], error) {
	// Check view or change permission
	if isNew {
		if !fv.admin.HasAddPermission(ctx, user) {
			return nil, fmt.Errorf("permission denied")
		}
	} else {
		if !fv.admin.HasViewPermission(ctx, user, instance) && !fv.admin.HasChangePermission(ctx, user, instance) {
			return nil, fmt.Errorf("permission denied")
		}
	}

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
					return fv.Render(ctx, r, user, instance, isNew, errors)
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
					return fv.Render(ctx, r, user, instance, isNew, errors)
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

	// Build field data with errors
	fieldData := make([]FormFieldData, 0)
	for _, field := range fields {
		if name := fv.getFieldName(field); name != "" {
			fd := fv.buildFieldData(ctx, instance, isNew, name, readOnlyFields[name])
			if errs, ok := errors[name]; ok {
				fd.Errors = append(fd.Errors, errs...)
			}
			fieldData = append(fieldData, fd)
		}
	}

	// Build fieldsets with errors
	fieldsetData := make([]FieldsetData, 0)
	for _, fs := range fieldsets {
		fsd := fv.buildFieldsetData(ctx, instance, isNew, fs, readOnlyFields, errors)
		fieldsetData = append(fieldsetData, fsd)
	}

	// Build inlines
	inlineData := fv.buildInlineData(ctx, config, instance, isNew)

	return &FormData[T]{
		Instance:            instance,
		IsNew:               isNew,
		Fields:              fieldData,
		Fieldsets:           fieldsetData,
		Inlines:             inlineData,
		Errors:              errors,
		HasAddPermission:    fv.admin.HasAddPermission(ctx, user),
		HasViewPermission:   fv.admin.HasViewPermission(ctx, user, instance),
		HasChangePermission: fv.admin.HasChangePermission(ctx, user, instance),
		HasDeletePermission: fv.admin.HasDeletePermission(ctx, user, instance),
		PrepopulatedFields:  fv.getPrepopulatedFields(config, instance, isNew),
	}, nil
}

func (fv *FormView[T]) getPrepopulatedFields(config *admin.Config[T], instance *T, isNew bool) map[string][]string {
	if config == nil {
		return nil
	}
	if config.GetPrepopulatedFields != nil {
		return config.GetPrepopulatedFields(context.Background(), instance, isNew)
	}
	return config.PrepopulatedFields
}

// ValidationError represents validation failures
type ValidationError struct {
	Errors map[string][]string
}

func (v ValidationError) Error() string {
	return "form validation failed"
}

// Save saves the form data to the model instance
func (fv *FormView[T]) Save(ctx context.Context, r *http.Request, instance *T, isNew bool, externalData admin.FormData) error {
	formData := externalData
	if formData == nil {
		// Parse form data from request
		formData = admin.FormData{}
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %w", err)
		}

		for key, values := range r.Form {
			if len(values) > 0 {
				formData[key] = values[0]
			}
		}
	}

	// Apply prepopulated fields
	prepopHandler := NewPrepopulatedFieldHandler(fv.admin)
	if err := prepopHandler.PopulateFields(ctx, instance, formData, isNew); err != nil {
		return fmt.Errorf("failed to populate fields: %w", err)
	}

	// Populate instance from form data
	fv.populateInstance(instance, formData)

	// Validate (simplified for now)
	formErrors := admin.ValidateForm(fv.admin, instance, formData, isNew)
	if len(formErrors) > 0 {
		errors := make(map[string][]string)
		for k, v := range formErrors {
			errors[k] = []string{v}
		}
		return ValidationError{Errors: errors}
	}

	// Use admin's SaveModel method
	if err := fv.admin.SaveModel(ctx, instance, formData, isNew); err != nil {
		return err
	}

	// Save inlines (only if not external data, or if specifically desired)
	// Usually externalData for list_editable doesn't include inlines
	if externalData == nil {
		return fv.saveInlines(ctx, r, instance)
	}
	return nil
}

// populateInstance populates an instance from form data
func (fv *FormView[T]) populateInstance(instance *T, formData admin.FormData) {
	discoveredFields := fv.admin.Fields()
	for _, fieldInfo := range discoveredFields {
		if value, ok := formData[fieldInfo.Name]; ok {
			// Skip special fields or read-only
			if fieldInfo.AutoNow || fieldInfo.AutoNowAdd || fieldInfo.AutoIncrement {
				continue
			}

			// Convert and set field using admin utility
			adminutils.SetFieldValue(instance, fieldInfo.Name, fv.convertFieldValue(value.(string), fieldInfo))
		}
	}
}

// convertFieldValue converts string values to appropriate types
func (fv *FormView[T]) convertFieldValue(value string, fieldInfo adminschema.FieldInfo) interface{} {
	switch fieldInfo.Type {
	case schema.TypeInt64:
		i, _ := strconv.ParseInt(value, 10, 64)
		return i
	case schema.TypeInt32:
		i, _ := strconv.ParseInt(value, 10, 32)
		return int32(i)
	case schema.TypeBool:
		return value == "on" || value == "true"
	case schema.TypeFloat64:
		f, _ := strconv.ParseFloat(value, 64)
		return f
	case schema.TypeString, schema.TypeText:
		return value
	}
	return value
}

// saveInlines saves all configured inlines for a parent instance
func (fv *FormView[T]) saveInlines(ctx context.Context, r *http.Request, instance *T) error {
	config := fv.admin.Config()
	if config == nil || len(config.Inlines) == 0 {
		return nil
	}

	for _, inline := range config.Inlines {
		val := reflect.ValueOf(inline)
		if ii, ok := val.Interface().(admin.InlineInterface[T]); ok {
			if err := ii.SaveInstances(ctx, instance, r); err != nil {
				return err
			}
		}
	}
	return nil
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

// buildFieldData builds form field data for a single field
func (fv *FormView[T]) buildFieldData(ctx context.Context, instance *T, isNew bool, fieldName string, isReadOnly bool) FormFieldData {
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
		return FormFieldData{Name: fieldName, Label: fieldName, ReadOnly: true} // Fallback
	}

	// Get widget based on field configuration
	config := fv.admin.Config()
	widget := fv.getWidgetForField(config, fieldName, fieldInfo, instance)

	// Get field value from instance
	var value interface{} = nil
	if instance != nil {
		val := reflect.ValueOf(instance)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		fieldVal := val.FieldByName(fieldName)
		if fieldVal.IsValid() && fieldVal.CanInterface() {
			value = fieldVal.Interface()
		}
	}

	return FormFieldData{
		Name:     fieldName,
		Label:    fieldInfo.VerboseName,
		Value:    value,
		Widget:   widget,
		HelpText: fieldInfo.HelpText,
		Required: fieldInfo.Required,
		ReadOnly: isReadOnly || fieldInfo.ReadOnly,
		Errors:   []string{},
	}
}

func (fv *FormView[T]) buildFieldsetData(ctx context.Context, instance *T, isNew bool, fs admin.Fieldset[T], readOnlyMap map[string]bool, errors map[string][]string) FieldsetData {
	fields := make([]FormFieldData, 0)
	for _, field := range fs.Fields {
		if name := fv.getFieldName(field); name != "" {
			fd := fv.buildFieldData(ctx, instance, isNew, name, readOnlyMap[name])
			if errs, ok := errors[name]; ok {
				fd.Errors = append(fd.Errors, errs...)
			}
			fields = append(fields, fd)
		}
	}
	return FieldsetData{
		Name:        fs.Name,
		Description: fs.Description,
		Collapsed:   fs.Collapsed,
		Fields:      fields,
	}
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

// buildInlineData builds inline data for the form view
func (fv *FormView[T]) buildInlineData(ctx context.Context, config *admin.Config[T], instance *T, isNew bool) []InlineData {
	if config == nil || len(config.Inlines) == 0 {
		return nil
	}

	result := make([]InlineData, 0, len(config.Inlines))

	for _, inline := range config.Inlines {
		val := reflect.ValueOf(inline)

		// Get model type and name
		modelField := val.FieldByName("model")
		modelName := "Related"
		if modelField.IsValid() {
			modelName = modelField.Type().Name()
		}

		style := "tabular"
		styleMethod := val.MethodByName("GetStyle")
		if styleMethod.IsValid() {
			results := styleMethod.Call(nil)
			if len(results) > 0 {
				style = fmt.Sprintf("%v", results[0].Interface())
			}
		}

		fieldsMethod := val.MethodByName("GetFields")
		var fields []InlineFieldData
		if fieldsMethod.IsValid() {
			results := fieldsMethod.Call(nil)
			if len(results) > 0 {
				if fvList, ok := results[0].Interface().([]interface{}); ok {
					for _, f := range fvList {
						fieldName := fv.getFieldName(f)
						if fieldName != "" {
							fields = append(fields, InlineFieldData{
								Name:  fieldName,
								Label: fieldName,
							})
						}
					}
				}
			}
		}

		extra := 1
		extraMethod := val.MethodByName("GetExtra")
		if extraMethod.IsValid() {
			results := extraMethod.Call(nil)
			if len(results) > 0 {
				if e, ok := results[0].Interface().(int); ok {
					extra = e
				}
			}
		}

		// Use InlineInterface for type-neutral access
		var instances []interface{}
		inlineModelName := modelName
		if ii, ok := val.Interface().(admin.InlineInterface[T]); ok {
			inlineModelName = ii.GetModelName()
			if !isNew {
				instances, _ = ii.GetInstances(ctx, instance)
			}
		}

		rows := make([]InlineRowData, 0, len(instances))
		for idx, inst := range instances {
			rows = append(rows, fv.buildInlineRowData(fields, inst, fmt.Sprintf("%s-%d", inlineModelName, idx)))
		}

		extraRows := make([]InlineRowData, 0, extra)
		for i := 0; i < extra; i++ {
			extraRows = append(extraRows, fv.buildInlineRowData(fields, nil, fmt.Sprintf("%s-%d", inlineModelName, len(instances)+i)))
		}

		maxNum := 0
		maxNumMethod := val.MethodByName("GetMaxNum")
		if maxNumMethod.IsValid() {
			results := maxNumMethod.Call(nil)
			if len(results) > 0 {
				if m, ok := results[0].Interface().(int); ok {
					maxNum = m
				}
			}
		}

		result = append(result, InlineData{
			Name:       inlineModelName,
			Title:      inlineModelName,
			Style:      style,
			Fields:     fields,
			Rows:       rows,
			ExtraForms: extraRows,
			CanAddMore: true,
			Extra:      extra,
			MaxNum:     maxNum,
		})
	}

	return result
}

func (fv *FormView[T]) buildInlineRowData(fields []InlineFieldData, instance interface{}, prefix string) InlineRowData {
	rowFields := make([]FormFieldData, 0, len(fields))
	for _, f := range fields {
		var value interface{}
		if instance != nil {
			// Extract value from instance via reflection
			val := reflect.ValueOf(instance)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			fieldVal := val.FieldByName(f.Name)
			if fieldVal.IsValid() {
				value = fieldVal.Interface()
			}
		}

		rowFields = append(rowFields, FormFieldData{
			Name:   fmt.Sprintf("%s-%s", prefix, f.Name),
			Label:  f.Label,
			Value:  value,
			Widget: admin.NewTextInput(), // Default widget
		})
	}

	id := interface{}(nil)
	if instance != nil {
		// Use a generic getID helper or similar
		val := reflect.ValueOf(instance)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		idField := val.FieldByName("ID")
		if !idField.IsValid() {
			idField = val.FieldByName("id")
		}
		if idField.IsValid() {
			id = idField.Interface()
		}
	}

	return InlineRowData{
		ID:     id,
		Fields: rowFields,
	}
}
