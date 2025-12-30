package admin

import (
	"context"

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
	Inlines []Inline[T, interface{}]

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
	field       interface{} // Field name (string) or FieldExpr
	descending  bool
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
	parentField interface{} // Field name (string) or FieldExpr
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
