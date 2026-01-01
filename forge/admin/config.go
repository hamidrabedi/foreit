package admin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/forgego/forge/orm"
)

// Config represents the configuration for an admin instance
// It can reference schema fields by name (schema-first) or use manual field expressions
type Config[T any] struct {
	// Basic information (can override schema meta)
	VerboseName       string
	VerboseNamePlural string

	// List view configuration
	// Can use field names (strings) that reference schema fields, or FieldExpr for custom
	ListDisplay         []interface{} // Can be string (field name) or FieldExpr
	ListDisplayLinks    []interface{} // Fields that link to detail view
	ListEditable        []interface{} // Fields editable in list view
	ListFilter          []interface{} // Can be string (field name) or Filter
	SearchFields        []interface{} // Can be string (field name) or FieldExpr
	SearchHelpText      string
	DateHierarchy       interface{} // Field name (string) or FieldExpr
	Ordering            []Ordering[T]
	ListPerPage         int
	ListMaxShowAll      int
	ListSelectRelated   []string // Field names for select_related optimization
	ListPrefetchRelated []string // Field names for prefetch_related optimization
	ShowFullResultCount bool
	PreserveFilters     bool
	EmptyValueDisplay   string
	SortableBy          []interface{} // Can be string (field name) or FieldExpr

	// Form configuration
	Fields                []interface{} // Can be string (field name) or FieldExpr
	Exclude               []interface{} // Field names to exclude
	ReadOnlyFields        []interface{} // Field names or FieldExpr
	PrepopulatedFields    map[string][]string
	RawIDFields           []interface{} // Field names or FieldExpr
	AutocompleteFields    []interface{} // Field names or FieldExpr
	RadioFields           map[string]RadioLayout
	Fieldsets             []Fieldset[T]
	Form                  FormGenerator[T]
	GetForm               func(ctx context.Context, instance *T, isNew bool) (Form[T], error)
	GetFields             func(ctx context.Context, instance *T, isNew bool) []interface{}
	GetFieldsets          func(ctx context.Context, instance *T, isNew bool) []Fieldset[T]
	GetReadOnlyFields     func(ctx context.Context, instance *T, isNew bool) []interface{}
	GetPrepopulatedFields func(ctx context.Context, instance *T, isNew bool) map[string][]string

	// Related models
	Inlines          []Inline[T, interface{}]
	RelationManagers []interface{} // RelationManager instances

	// Actions
	Actions []Action[T]

	// Permissions
	PermissionChecker   PermissionChecker
	HasAddPermission    func(ctx context.Context, admin *Admin[T], user interface{}) bool
	HasChangePermission func(ctx context.Context, admin *Admin[T], user interface{}, obj *T) bool
	HasDeletePermission func(ctx context.Context, admin *Admin[T], user interface{}, obj *T) bool
	HasViewPermission   func(ctx context.Context, admin *Admin[T], user interface{}, obj *T) bool
	HasModulePermission func(ctx context.Context, admin *Admin[T], user interface{}) bool

	// View customization hooks
	GetQueryset func(ctx context.Context, admin *Admin[T], qs orm.QuerySet[T]) (orm.QuerySet[T], error)
	SaveModel   func(ctx context.Context, admin *Admin[T], instance *T, formData FormData, isNew bool) error
	DeleteModel func(ctx context.Context, admin *Admin[T], instance *T) error

	// View hooks for customizing views
	ChangelistViewHook ChangelistViewHook[T]
	ChangeViewHook     ChangeViewHook[T]
	AddViewHook        AddViewHook[T]
	DeleteViewHook     DeleteViewHook[T]
	HistoryViewHook    HistoryViewHook[T]
	ResponseAddHook    ResponseAddHook[T]
	ResponseChangeHook ResponseChangeHook[T]
	ResponseDeleteHook ResponseDeleteHook[T]
	GetURLsHook        GetURLsHook[T]

	// Advanced features
	SaveAs         bool
	SaveAsContinue bool
	SaveOnTop      bool
	ViewOnSite     func(ctx context.Context, instance *T) string
	ShowChangeLink bool

	// Form field customization
	FormFieldOverrides     map[string]Widget
	FormFieldForForeignKey func(ctx context.Context, admin *Admin[T], fieldName string, instance *T) Widget
	FormFieldForManyToMany func(ctx context.Context, admin *Admin[T], fieldName string, instance *T) Widget
	FormFieldForDBField    func(ctx context.Context, admin *Admin[T], fieldName string, instance *T) Widget
}

// RadioLayout specifies the layout for radio buttons
type RadioLayout string

const (
	RadioHorizontal RadioLayout = "horizontal"
	RadioVertical   RadioLayout = "vertical"
)

// FormGenerator is a function type for generating custom forms
type FormGenerator[T any] func(ctx context.Context, instance *T, isNew bool) (Form[T], error)

// Form represents a form for create/update
type Form[T any] interface {
	Fields() []FormField[T]
	AddField(field FormField[T])
	Validate() error
}

// FormField represents a form field
type FormField[T any] interface {
	Name() string
	Value() interface{}
	SetValue(value interface{}) error
	IsRequired() bool
	IsReadOnly() bool
	Widget() Widget
}

// Ordering represents field ordering
type Ordering[T any] struct {
	field      interface{} // Field name (string) or FieldExpr
	descending bool
}

// OrderBy creates an ordering
func OrderBy[T any](field interface{}) *Ordering[T] {
	return &Ordering[T]{
		field:      field,
		descending: false,
	}
}

// Desc makes the ordering descending
func (o *Ordering[T]) Desc() *Ordering[T] {
	o.descending = true
	return o
}

// Asc makes the ordering ascending
func (o *Ordering[T]) Asc() *Ordering[T] {
	o.descending = false
	return o
}

// Field returns the field
func (o *Ordering[T]) Field() interface{} {
	return o.field
}

// IsDescending returns if ordering is descending
func (o *Ordering[T]) IsDescending() bool {
	return o.descending
}

// Fieldset represents a fieldset for grouping form fields
type Fieldset[T any] struct {
	Name        string
	Fields      []interface{} // Field names or FieldExpr
	Collapsed   bool
	Description string
}

// NewFieldset creates a new fieldset
func NewFieldset[T any](name string, fields ...interface{}) Fieldset[T] {
	return Fieldset[T]{
		Name:   name,
		Fields: fields,
	}
}

// WithCollapsed sets the fieldset as collapsed
func (f Fieldset[T]) WithCollapsed(collapsed bool) Fieldset[T] {
	f.Collapsed = collapsed
	return f
}

// WithDescription sets the fieldset description
func (f Fieldset[T]) WithDescription(desc string) Fieldset[T] {
	f.Description = desc
	return f
}

// Inline represents inline editing for related models
type Inline[T any, R any] struct {
	model       R
	manager     *orm.Manager[R]
	parentField interface{}   // Field name (string) or FieldExpr
	fields      []interface{} // Field names or FieldExpr
	extra       int
	maxNum      int
	style       InlineStyle
}

// InlineStyle specifies the display style
type InlineStyle string

const (
	InlineTabular InlineStyle = "tabular"
	InlineStacked InlineStyle = "stacked"
)

// TabularInline creates a tabular inline
func TabularInline[T any, R any](
	model R,
	manager *orm.Manager[R],
	parentField interface{},
	fields []interface{},
) Inline[T, R] {
	return Inline[T, R]{
		model:       model,
		manager:     manager,
		parentField: parentField,
		fields:      fields,
		extra:       1,
		style:       InlineTabular,
	}
}

// StackedInline creates a stacked inline
func StackedInline[T any, R any](
	model R,
	manager *orm.Manager[R],
	parentField interface{},
	fields []interface{},
) Inline[T, R] {
	return Inline[T, R]{
		model:       model,
		manager:     manager,
		parentField: parentField,
		fields:      fields,
		extra:       1,
		style:       InlineStacked,
	}
}

// WithExtra sets the number of extra empty forms
func (i Inline[T, R]) WithExtra(extra int) Inline[T, R] {
	i.extra = extra
	return i
}

// WithMaxNum sets the maximum number of forms
func (i Inline[T, R]) WithMaxNum(maxNum int) Inline[T, R] {
	i.maxNum = maxNum
	return i
}

// GetExtra returns the number of extra empty forms
func (i Inline[T, R]) GetExtra() int {
	return i.extra
}

// GetMaxNum returns the maximum number of forms
func (i Inline[T, R]) GetMaxNum() int {
	return i.maxNum
}

// GetFields returns the fields for this inline
func (i Inline[T, R]) GetFields() []interface{} {
	return i.fields
}

// GetParentField returns the parent field
func (i Inline[T, R]) GetParentField() interface{} {
	return i.parentField
}

// GetStyle returns the inline style
func (i Inline[T, R]) GetStyle() InlineStyle {
	return i.style
}

// GetInstances returns related instances for a parent
func (i Inline[T, R]) GetInstances(ctx context.Context, parent *T) ([]interface{}, error) {
	// Query related instances
	fieldName := i.GetParentFieldName()
	if fieldName == "" {
		return nil, fmt.Errorf("invalid parent field configuration")
	}

	expr := genericExpression{
		field: fieldName,
		op:    orm.OpEquals,
		value: getID(parent),
	}

	qs, err := i.manager.Filter(expr)
	if err != nil {
		return nil, err
	}

	instances, err := qs.All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(instances))
	for idx, inst := range instances {
		result[idx] = inst
	}
	return result, nil
}

// SaveInstances saves related instances from form data
func (i Inline[T, R]) SaveInstances(ctx context.Context, parent *T, r *http.Request) error {
	modelName := i.GetModelName()

	// Determine how many forms were submitted
	// This is a simplified version of Django's ManagementForm
	// We iterate through indices until we don't find any more
	for idx := 0; ; idx++ {
		prefix := fmt.Sprintf("%s-%d", modelName, idx)

		// Check if this row exists in the form (at least one field or the DELETE checkbox)
		found := false
		for key := range r.PostForm {
			if strings.HasPrefix(key, prefix) {
				found = true
				break
			}
		}

		if !found {
			break
		}

		// Check for deletion
		deleteVal := r.PostFormValue(prefix + "-DELETE")
		if deleteVal != "" {
			// If it's an existing instance, delete it
			idStr := deleteVal
			if idStr != "" && idStr != "on" { // value might be the ID
				id, _ := strconv.ParseInt(idStr, 10, 64)
				if id > 0 {
					var zero R
					instance := &zero
					setID(instance, id)
					if err := i.manager.Delete(ctx, instance); err != nil {
						return fmt.Errorf("failed to delete inline instance %d: %w", id, err)
					}
				}
			}
			continue
		}

		// Populate instance
		var zero R
		instance := &zero

		// Check if it's an existing instance (look for hidden ID field if we had one, or use the DELETE's value)
		// For now, let's assume we might have a hidden ID field or we can find it in the prefix
		// Actually, let's look for prefix + "-ID"
		idField := r.PostFormValue(prefix + "-id")
		isNew := true
		if idField != "" {
			id, _ := strconv.ParseInt(idField, 10, 64)
			if id > 0 {
				existing, err := i.manager.Get(ctx, id)
				if err == nil {
					instance = existing
					isNew = false
				}
			}
		}

		// Set parent field
		parentFieldName := i.GetParentFieldName()
		if err := setParentField(instance, parentFieldName, parent); err != nil {
			return fmt.Errorf("failed to set parent field on inline %d: %w", idx, err)
		}

		// Set other fields from form data
		fieldsChanged := false
		for _, field := range i.fields {
			fieldName := ""
			if name, ok := field.(string); ok {
				fieldName = name
			} else {
				// Try Name()
				val := reflect.ValueOf(field)
				method := val.MethodByName("Name")
				if method.IsValid() {
					results := method.Call(nil)
					if len(results) > 0 {
						fieldName = results[0].String()
					}
				}
			}

			if fieldName == "" {
				continue
			}

			if val, ok := r.PostForm[prefix+"-"+fieldName]; ok && len(val) > 0 {
				if err := setFieldValue(instance, fieldName, val[0]); err != nil {
					// Log error or return?
				}
				fieldsChanged = true
			}
		}

		// Save if it's new, or if it changed (always save for now if not deleting)
		if isNew {
			// Only save new if some fields were filled
			if fieldsChanged {
				if err := i.manager.Create(ctx, instance); err != nil {
					return fmt.Errorf("failed to create inline instance %d: %w", idx, err)
				}
			}
		} else {
			if err := i.manager.Update(ctx, instance); err != nil {
				return fmt.Errorf("failed to update inline instance %d: %w", idx, err)
			}
		}
	}

	return nil
}

// GetParentFieldName returns the parent field name
func (i Inline[T, R]) GetParentFieldName() string {
	if name, ok := i.parentField.(string); ok {
		return name
	}
	// Try reflection for Name() method
	val := reflect.ValueOf(i.parentField)
	method := val.MethodByName("Name")
	if method.IsValid() {
		results := method.Call(nil)
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return ""
}

// GetModelName returns the name of the related model
func (i Inline[T, R]) GetModelName() string {
	var zero R
	typ := reflect.TypeOf(zero)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Name()
}

// InlineInterface allows type-neutral access to inlines
type InlineInterface[T any] interface {
	GetStyle() InlineStyle
	GetFields() []interface{}
	GetExtra() int
	GetMaxNum() int
	GetInstances(ctx context.Context, parent *T) ([]interface{}, error)
	SaveInstances(ctx context.Context, parent *T, r *http.Request) error
	GetModelName() string
}

// Internal helpers copied from inlines package to avoid circular dependency
type genericExpression struct {
	field string
	op    orm.Operator
	value interface{}
}

func (e genericExpression) ToSQL(builder *orm.SQLBuilder) (string, []interface{}, error) {
	placeholder := builder.AddArg(e.value)
	return fmt.Sprintf("%s %s %s", orm.EscapeIdentifier(e.field), e.op, placeholder), []interface{}{e.value}, nil
}

func (e genericExpression) Resolve(schema *orm.ModelSchema) error {
	return nil
}

func getID(instance interface{}) interface{} {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	field := val.FieldByName("ID")
	if !field.IsValid() {
		field = val.FieldByName("id")
	}
	if field.IsValid() {
		return field.Interface()
	}
	return nil
}

func setID(instance interface{}, id int64) {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	field := val.FieldByName("ID")
	if !field.IsValid() {
		field = val.FieldByName("id")
	}
	if field.IsValid() && field.CanSet() {
		field.SetInt(id)
	}
}

func setParentField(instance interface{}, fieldName string, parent interface{}) error {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	parentVal := reflect.ValueOf(parent)
	if field.Type().Kind() == reflect.Ptr && parentVal.Type().Kind() != reflect.Ptr {
		parentVal = parentVal.Addr()
	}

	if parentVal.Type().AssignableTo(field.Type()) {
		field.Set(parentVal)
		return nil
	}

	// If it's an ID field
	if field.Kind() == reflect.Int64 {
		parentID := getID(parent)
		if id, ok := parentID.(int64); ok {
			field.SetInt(id)
			return nil
		}
	}

	return fmt.Errorf("cannot assign %T to %s (type %v)", parent, fieldName, field.Type())
}

func setFieldValue(instance interface{}, fieldName string, value string) error {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		// Try lowercase
		field = val.FieldByName(strings.Title(fieldName))
		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("field %s not found or not settable", fieldName)
		}
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(value, 10, 64)
		field.SetInt(i)
	case reflect.Bool:
		field.SetBool(value == "on" || value == "true")
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(value, 64)
		field.SetFloat(f)
	}
	return nil
}

// Action represents a bulk action
type Action[T any] struct {
	Name        string
	Label       string
	Description string
	Handler     func(ctx context.Context, instances []*T) error
	Permissions []string
}

// NewAction creates a new action
func NewAction[T any](
	name string,
	label string,
	handler func(ctx context.Context, instances []*T) error,
) Action[T] {
	return Action[T]{
		Name:    name,
		Label:   label,
		Handler: handler,
	}
}

// WithDescription sets the action description
func (a Action[T]) WithDescription(desc string) Action[T] {
	a.Description = desc
	return a
}

// WithPermissions sets required permissions
func (a Action[T]) WithPermissions(perms ...string) Action[T] {
	a.Permissions = perms
	return a
}
