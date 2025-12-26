package admin

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/forgego/forge/pkg/admin/templates"
	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/schema"
)

// FormData contains data for rendering admin forms
type FormData struct {
	Model     *AdminModel
	Instance  interface{}
	Fields    []FormField
	Fieldsets []FormFieldsetData // Grouped fields (if fieldsets configured)
	Errors    map[string]string
	BaseURL   string
	IsCreate  bool
	SaveOnTop bool
	SaveAs    bool
	SaveAsContinue bool
	SaveAndAddAnother bool
}

// FormField represents a form field
type FormField struct {
	Name        string
	Label       string
	Type        string
	Value       interface{}
	Required    bool
	HelpText    string
	Errors      []string
	Choices     []schema.Choice
	ReadOnly    bool
}

// handleModelCreate handles GET request for create form
func handleModelCreate(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canAdd(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to add this object.", http.StatusForbidden)
			return
		}

		// Create new instance
		instance := CreateNewInstance(model.Model)

		// Generate form fields from model
		fields := generateFormFields(model, instance, true)

		data := FormData{
			Model:    model,
			Instance: instance,
			Fields:   fields,
			Errors:   make(map[string]string),
			BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
			IsCreate: true,
		}

		if err := renderFormTemplate(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render form: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// handleModelCreatePost handles POST request for create form
func handleModelCreatePost(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canAdd(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to add this object.", http.StatusForbidden)
			return
		}

		ctx := r.Context()

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
			return
		}

		// Create new instance
		instance := CreateNewInstance(model.Model)

		// Populate instance from form
		errors := populateInstanceFromForm(model, instance, r.Form)

		// If there are errors, re-render form
		if len(errors) > 0 {
			fields := generateFormFields(model, instance, true)
			data := FormData{
				Model:    model,
				Instance: instance,
				Fields:   fields,
				Errors:   errors,
				BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
				IsCreate: true,
			}
			if err := renderFormTemplate(w, data); err != nil {
				http.Error(w, fmt.Sprintf("Failed to render form: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		// Save instance using manager
		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		// Get Create method
		managerValue := reflect.ValueOf(manager)
		createMethod := managerValue.MethodByName("Create")
		if !createMethod.IsValid() {
			http.Error(w, "Manager does not have Create method", http.StatusInternalServerError)
			return
		}

		// Call Create
		results := createMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(instance),
		})

		if len(results) > 0 && !results[0].IsNil() {
			if err, ok := results[0].Interface().(error); ok && err != nil {
				// Re-render form with error
				fields := generateFormFields(model, instance, true)
				data := FormData{
					Model:    model,
					Instance: instance,
					Fields:   fields,
					Errors:   map[string]string{"_error": err.Error()},
					BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
					IsCreate: true,
				}
				if err := renderFormTemplate(w, data); err != nil {
					http.Error(w, fmt.Sprintf("Failed to render form: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		// Redirect to detail page
		id := getIDFromInstance(instance)
		if id > 0 {
			http.Redirect(w, r, fmt.Sprintf("/admin/%s/%d/", modelName, id), http.StatusSeeOther)
			return
		}

		// Redirect to list if no ID
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/", modelName), http.StatusSeeOther)
	}
}

// handleModelUpdate handles GET request for update form
func handleModelUpdate(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canChange(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to change this object.", http.StatusForbidden)
			return
		}

		ctx := r.Context()

		// Get ID from URL
		idStr := httplib.GetParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Get instance using manager
		modelType := reflect.TypeOf(model.Model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}
		instanceValue := reflect.New(modelType)
		instance := instanceValue.Interface()

		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		// Get Get method
		managerValue := reflect.ValueOf(manager)
		getMethod := managerValue.MethodByName("Get")
		if !getMethod.IsValid() {
			http.Error(w, "Manager does not have Get method", http.StatusInternalServerError)
			return
		}

		// Call Get
		results := getMethod.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(id),
		})

		if len(results) < 2 {
			http.Error(w, "Invalid Get method signature", http.StatusInternalServerError)
			return
		}

		if !results[1].IsNil() {
			if err, ok := results[1].Interface().(error); ok && err != nil {
				http.Error(w, fmt.Sprintf("Failed to get instance: %v", err), http.StatusNotFound)
				return
			}
		}

		instance = results[0].Interface()

		// Generate form fields
		fields := generateFormFields(model, instance, false)

		data := FormData{
			Model:    model,
			Instance: instance,
			Fields:   fields,
			Errors:   make(map[string]string),
			BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
			IsCreate: false,
		}

		if err := renderFormTemplate(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render form: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// handleModelUpdatePost handles POST request for update form
func handleModelUpdatePost(modelName string, model *AdminModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		if !canChange(r, modelName) {
			http.Error(w, "Permission denied: You do not have permission to change this object.", http.StatusForbidden)
			return
		}

		ctx := r.Context()

		// Get ID from URL
		idStr := httplib.GetParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
			return
		}

		manager, err := getManagerForModel(model)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get manager: %v", err), http.StatusInternalServerError)
			return
		}

		ops := NewManagerOps()
		instance, err := ops.GetInstance(ctx, manager, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get instance: %v", err), http.StatusNotFound)
			return
		}

		// Populate instance from form
		errors := populateInstanceFromForm(model, instance, r.Form)

		// If there are errors, re-render form
		if len(errors) > 0 {
			fields := generateFormFields(model, instance, false)
			data := FormData{
				Model:    model,
				Instance: instance,
				Fields:   fields,
				Errors:   errors,
				BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
				IsCreate: false,
			}
			if err := renderFormTemplate(w, data); err != nil {
				http.Error(w, fmt.Sprintf("Failed to render form: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		// Save instance using manager
		if err := ops.UpdateInstance(ctx, manager, instance); err != nil {
			// Re-render form with error
			fields := generateFormFields(model, instance, false)
			data := FormData{
				Model:    model,
				Instance: instance,
				Fields:   fields,
				Errors:   map[string]string{"_error": err.Error()},
				BaseURL:  fmt.Sprintf("/admin/%s/", modelName),
				IsCreate: false,
			}
			if err := renderFormTemplate(w, data); err != nil {
				http.Error(w, fmt.Sprintf("Failed to render form: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		// Redirect to detail page
		http.Redirect(w, r, fmt.Sprintf("/admin/%s/%d/", modelName, id), http.StatusSeeOther)
	}
}

// generateFormFields generates form fields from model schema
func generateFormFields(model *AdminModel, instance interface{}, isCreate bool) []FormField {
	var fields []FormField

	// Try to get fields from schema
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	// Check if instance implements schema.Schema
	if schemaInstance, ok := instance.(schema.Schema); ok {
		schemaFields := schemaInstance.Fields()
		for _, schemaField := range schemaFields {
			// Skip primary key and auto-increment fields on create
			if isCreate && (schemaField.PrimaryKey || schemaField.AutoIncrement) {
				continue
			}

			// Skip read-only fields
			if !schemaField.Editable && !isCreate {
				continue
			}

			// Get field value from instance
			fieldValue := instanceValue.FieldByName(schemaField.Name)
			var value interface{}
			if fieldValue.IsValid() {
				value = fieldValue.Interface()
			}

			// Determine HTML input type
			inputType := getInputTypeForField(schemaField)

			field := FormField{
				Name:     schemaField.Name,
				Label:    getFieldLabel(schemaField),
				Type:     inputType,
				Value:    value,
				Required: schemaField.Required,
				HelpText: schemaField.HelpText,
				Choices:  schemaField.Choices,
				ReadOnly: !schemaField.Editable && !isCreate,
			}

			fields = append(fields, field)
		}
	} else {
		// Fallback: use reflection to get struct fields
		typ := instanceValue.Type()
		for i := 0; i < typ.NumField(); i++ {
			structField := typ.Field(i)
			if !structField.IsExported() {
				continue
			}

			// Skip BaseSchema
			if structField.Name == "BaseSchema" {
				continue
			}

			fieldValue := instanceValue.Field(i)
			field := FormField{
				Name:     structField.Name,
				Label:    structField.Name,
				Type:     "text",
				Value:    fieldValue.Interface(),
				Required: false,
			}

			fields = append(fields, field)
		}
	}

	return fields
}

// getInputTypeForField determines HTML input type from schema field
func getInputTypeForField(field schema.Field) string {
	switch field.Type {
	case schema.TypeBool:
		return "checkbox"
	case schema.TypeInt64, schema.TypeInt32:
		return "number"
	case schema.TypeFloat64, schema.TypeDecimal:
		return "number"
	case schema.TypeEmail:
		return "email"
	case schema.TypeURL:
		return "url"
	case schema.TypeTime, schema.TypeDate, schema.TypeDateTime:
		return "datetime-local"
	case schema.TypeText:
		return "textarea"
	case schema.TypeUUID:
		return "text"
	default:
		if len(field.Choices) > 0 {
			return "select"
		}
		return "text"
	}
}

// getFieldLabel gets the display label for a field
func getFieldLabel(field schema.Field) string {
	if field.VerboseName != "" {
		return field.VerboseName
	}
	return field.Name
}

// populateInstanceFromForm populates instance from form values
func populateInstanceFromForm(model *AdminModel, instance interface{}, formValues map[string][]string) map[string]string {
	errors := make(map[string]string)
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	// Get schema fields if available
	var schemaFields []schema.Field
	if schemaInstance, ok := instance.(schema.Schema); ok {
		schemaFields = schemaInstance.Fields()
	}

	// Create field map
	fieldMap := make(map[string]schema.Field)
	for _, field := range schemaFields {
		fieldMap[field.Name] = field
	}

	// Populate fields
	for fieldName, values := range formValues {
		if len(values) == 0 {
			continue
		}

		field := instanceValue.FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		value := values[0]

		// Get schema field info
		schemaField, hasSchema := fieldMap[fieldName]

		// Skip read-only fields
		if hasSchema && !schemaField.Editable {
			continue
		}

		// Parse value based on field type
		if hasSchema {
			parsedValue, err := parseFormValue(value, schemaField.Type)
			if err != nil {
				errors[fieldName] = err.Error()
				continue
			}
			field.Set(reflect.ValueOf(parsedValue))
		} else {
			// Fallback: try to set as string
			field.SetString(value)
		}
	}

	return errors
}

// parseFormValue parses form value based on field type
func parseFormValue(value string, fieldType schema.FieldType) (interface{}, error) {
	const trueValue = "true"
	switch fieldType {
	case schema.TypeBool:
		return value == "on" || value == trueValue || value == "1", nil
	case schema.TypeInt64:
		return strconv.ParseInt(value, 10, 64)
	case schema.TypeInt32:
		val, err := strconv.ParseInt(value, 10, 32)
		return int32(val), err
	case schema.TypeFloat64:
		return strconv.ParseFloat(value, 64)
	case schema.TypeTime, schema.TypeDate, schema.TypeDateTime:
		// Parse datetime - simplified version
		if value == "" {
			return time.Time{}, nil
		}
		return time.Parse("2006-01-02T15:04", value)
	default:
		return value, nil
	}
}

// renderFormTemplate renders the admin form template
func renderFormTemplate(w http.ResponseWriter, data FormData) error {
	// Load templates
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Format field values for display (HTML escaping is automatic with html/template)
	formattedFields := make([]FormField, len(data.Fields))
	for i, field := range data.Fields {
		formattedFields[i] = field
		// Format value for display
		if field.Value != nil {
			switch v := field.Value.(type) {
			case time.Time:
				if !v.IsZero() {
					formattedFields[i].Value = v.Format("2006-01-02T15:04")
				} else {
					formattedFields[i].Value = ""
				}
			case bool:
				// Boolean values are handled in checkbox template
			default:
				// Other values are displayed as-is (html/template will escape)
				formattedFields[i].Value = fmt.Sprintf("%v", v)
			}
		}
	}

	// Generate fieldsets
	fieldsets := GenerateFieldsets(data.Model, data.Instance, data.IsCreate)
	
	// Convert fieldsets to template format with properly formatted fields
	fieldsetTemplates := make([]map[string]interface{}, len(fieldsets))
	for i, fs := range fieldsets {
		// Convert FormField to template field format
		templateFields := make([]map[string]interface{}, len(fs.Fields))
		for j, field := range fs.Fields {
			fieldHTML := renderFieldHTML(field)
			contents := ""
			if field.ReadOnly {
				contents = fmt.Sprintf("%v", field.Value)
			}
			
			templateFields[j] = map[string]interface{}{
				"ID":         field.Name,
				"Name":      field.Name,
				"Label":     field.Label,
				"Type":      field.Type,
				"Value":     field.Value,
				"Required":  field.Required,
				"HelpText":  field.HelpText,
				"Errors":    field.Errors,
				"ReadOnly":  field.ReadOnly,
				"IsCheckbox": field.Type == "checkbox",
				"IsReadonly": field.ReadOnly,
				"IsVisible":  true,
				"Field":      template.HTML(fieldHTML),
				"Contents":   contents,
			}
		}
		
		fieldsetTemplates[i] = map[string]interface{}{
			"ID":            i,
			"Name":          fs.Name,
			"Fields":        templateFields,
			"IsCollapsible": fs.Name != "",
			"Collapsed":     fs.Collapsed,
			"Classes":       strings.Join(fs.Classes, " "),
			"Description":   "",
		}
	}
	
	// Generate submit button data
	submitButtons := generateSubmitButtons(data)
	
	// Get permissions (defaults - should be passed from request context)
	canChange := true
	canDelete := true
	
	// Prepare template data
	templateData := map[string]interface{}{
		"Title":        fmt.Sprintf("%s %s", map[bool]string{true: "Create", false: "Edit"}[data.IsCreate], data.Model.Name),
		"Model":        data.Model,
		"Instance":     data.Instance,
		"Fields":       formattedFields,
		"Fieldsets":    fieldsetTemplates,
		"Errors":       data.Errors,
		"BaseURL":      data.BaseURL,
		"IsCreate":     data.IsCreate,
		"Models":       GetAllModels(), // For navigation
		"SaveOnTop":    data.SaveOnTop,
		"SubmitButtons": submitButtons,
		"CanChange":    canChange,
		"CanDelete":    canDelete,
	}

	// Execute form template
	// html/template automatically escapes all values for XSS protection
	return tmpl.ExecuteTemplate(w, "form", templateData)
}

// generateSubmitButtons generates submit button configuration
func generateSubmitButtons(data FormData) map[string]interface{} {
	baseURL := data.BaseURL
	instanceID := ""
	if !data.IsCreate {
		// Try to get ID from instance
		instanceValue := reflect.ValueOf(data.Instance)
		if instanceValue.Kind() == reflect.Ptr {
			instanceValue = instanceValue.Elem()
		}
		idField := instanceValue.FieldByName("ID")
		if idField.IsValid() {
			instanceID = fmt.Sprintf("%v", idField.Interface())
		}
	}
	
	showSave := true
	showSaveAsNew := !data.IsCreate && data.SaveAs
	showSaveAndAddAnother := data.SaveAndAddAnother
	showSaveAndContinue := data.SaveAsContinue
	showClose := !data.IsCreate
	showDeleteLink := !data.IsCreate
	
	changelistURL := baseURL
	deleteURL := ""
	if instanceID != "" {
		deleteURL = fmt.Sprintf("%s%s/delete/", baseURL, instanceID)
	}
	
	return map[string]interface{}{
		"ShowSave":            showSave,
		"ShowSaveAsNew":       showSaveAsNew,
		"ShowSaveAndAddAnother": showSaveAndAddAnother,
		"ShowSaveAndContinue": showSaveAndContinue,
		"ShowClose":           showClose,
		"ShowDeleteLink":      showDeleteLink,
		"CanChange":           true, // Should come from permissions
		"ChangelistURL":       changelistURL,
		"DeleteURL":           deleteURL,
		"Original":            !data.IsCreate,
	}
}

// renderFieldHTML renders HTML for a form field
func renderFieldHTML(field FormField) string {
	var html strings.Builder
	valueStr := ""
	if field.Value != nil {
		switch v := field.Value.(type) {
		case bool:
			if v {
				valueStr = "checked"
			}
		case time.Time:
			if !v.IsZero() {
				valueStr = v.Format("2006-01-02T15:04")
			}
		default:
			valueStr = fmt.Sprintf("%v", v)
		}
	}
	
	switch field.Type {
	case "textarea":
		html.WriteString(fmt.Sprintf(`<textarea name="%s" id="%s"`, field.Name, field.Name))
		if field.Required {
			html.WriteString(` required`)
		}
		if field.ReadOnly {
			html.WriteString(` readonly`)
		}
		html.WriteString(`>`)
		if valueStr != "" {
			html.WriteString(template.HTMLEscapeString(valueStr))
		}
		html.WriteString(`</textarea>`)
	case "select":
		html.WriteString(fmt.Sprintf(`<select name="%s" id="%s"`, field.Name, field.Name))
		if field.Required {
			html.WriteString(` required`)
		}
		if field.ReadOnly {
			html.WriteString(` disabled`)
		}
		html.WriteString(`>`)
		html.WriteString(`<option value="">---</option>`)
		for _, choice := range field.Choices {
			selected := ""
			if fmt.Sprintf("%v", choice.Value) == valueStr {
				selected = " selected"
			}
			html.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`,
				template.HTMLEscapeString(fmt.Sprintf("%v", choice.Value)),
				selected,
				template.HTMLEscapeString(choice.Label)))
		}
		html.WriteString(`</select>`)
	case "checkbox":
		html.WriteString(fmt.Sprintf(`<input type="checkbox" name="%s" id="%s"`, field.Name, field.Name))
		if valueStr == "checked" || valueStr == "true" {
			html.WriteString(` checked`)
		}
		if field.ReadOnly {
			html.WriteString(` disabled`)
		}
		html.WriteString(`>`)
	default:
		html.WriteString(fmt.Sprintf(`<input type="%s" name="%s" id="%s"`, field.Type, field.Name, field.Name))
		if valueStr != "" && field.Type != "checkbox" {
			html.WriteString(fmt.Sprintf(` value="%s"`, template.HTMLEscapeString(valueStr)))
		}
		if field.Required {
			html.WriteString(` required`)
		}
		if field.ReadOnly {
			html.WriteString(` readonly`)
		}
		html.WriteString(`>`)
	}
	
	return html.String()
}

