package admin

import (
	"context"

	query "github.com/forgego/forge/orm"
)

// Config represents the configuration for an admin instance
type Config[T any] struct {
	// Basic information
	VerboseName       string
	VerboseNamePlural string

	// List view configuration
	ListDisplay      []FieldExpr[T, interface{}]
	ListDisplayLinks []FieldExpr[T, interface{}] // Fields that link to detail view
	ListEditable     []FieldExpr[T, interface{}] // Fields editable in list view
	ListFilter       []Filter[T]
	SearchFields     []FieldExpr[T, interface{}]
	SearchHelpText   string
	DateHierarchy    FieldExpr[T, interface{}] // Field for date-based navigation
	Ordering         []Ordering[T]
	ListPerPage      int
	ListMaxShowAll   int // Maximum items to show in "show all" mode
	ListSelectRelated []string // Fields for select_related optimization
	ListPrefetchRelated []string // Fields for prefetch_related optimization
	ShowFullResultCount bool // Show full count vs filtered count
	PreserveFilters     bool // Preserve filters when navigating
	EmptyValueDisplay    string // Display value for empty/null fields
	SortableBy           []FieldExpr[T, interface{}] // Fields that can be sorted

	// Form configuration
	Fields              []FieldExpr[T, interface{}] // Explicit field ordering
	Exclude             []FieldExpr[T, interface{}] // Fields to exclude
	ReadOnlyFields      []FieldExpr[T, interface{}]
	PrepopulatedFields  map[string][]string // Field -> source fields for auto-population
	RawIDFields         []FieldExpr[T, interface{}] // Use raw ID input for foreign keys
	AutocompleteFields  []FieldExpr[T, interface{}] // Autocomplete for foreign keys
	RadioFields         map[string]RadioLayout // Use radio buttons for choices
	Fieldsets           []Fieldset[T]
	Form                FormGenerator[T] // Custom form generator
	GetForm             func(ctx context.Context, instance *T, isNew bool) (Form[T], error)
	GetFields            func(ctx context.Context, instance *T, isNew bool) []FieldExpr[T, interface{}]
	GetFieldsets         func(ctx context.Context, instance *T, isNew bool) []Fieldset[T]
	GetReadOnlyFields    func(ctx context.Context, instance *T, isNew bool) []FieldExpr[T, interface{}]
	GetPrepopulatedFields func(ctx context.Context, instance *T, isNew bool) map[string][]string

	// Related models
	Inlines []Inline[T, interface{}]

	// Actions
	Actions []Action[T]

	// Permissions
	PermissionChecker PermissionChecker
	HasAddPermission    func(ctx context.Context, admin *Admin[T], user interface{}) bool
	HasChangePermission func(ctx context.Context, admin *Admin[T], user interface{}, obj *T) bool
	HasDeletePermission  func(ctx context.Context, admin *Admin[T], user interface{}, obj *T) bool
	HasViewPermission    func(ctx context.Context, admin *Admin[T], user interface{}, obj *T) bool
	HasModulePermission  func(ctx context.Context, admin *Admin[T], user interface{}) bool

	// View customization hooks
	GetQueryset func(ctx context.Context, admin *Admin[T], qs query.QuerySet[T]) (query.QuerySet[T], error)
	SaveModel   func(ctx context.Context, admin *Admin[T], instance *T, formData FormData, isNew bool) error
	DeleteModel func(ctx context.Context, admin *Admin[T], instance *T) error

	// Advanced features
	SaveAs          bool // Show "Save as new" button
	SaveAsContinue  bool // Continue editing after "save as new"
	SaveOnTop       bool // Show save buttons at top of form
	ViewOnSite      func(ctx context.Context, instance *T) string // Get URL for viewing object on public site
	ShowChangeLink  bool // Show change link in inlines

	// Form field customization
	FormFieldOverrides map[string]Widget // Override widgets for specific fields
	FormFieldForForeignKey func(ctx context.Context, admin *Admin[T], field FieldExpr[T, interface{}], instance *T) Widget
	FormFieldForManyToMany  func(ctx context.Context, admin *Admin[T], field FieldExpr[T, interface{}], instance *T) Widget
	FormFieldForDBField     func(ctx context.Context, admin *Admin[T], field FieldExpr[T, interface{}], instance *T) Widget
}

// RadioLayout specifies the layout for radio buttons
type RadioLayout string

const (
	RadioHorizontal RadioLayout = "horizontal"
	RadioVertical   RadioLayout = "vertical"
)

// FormGenerator is a function type for generating custom forms
type FormGenerator[T any] func(ctx context.Context, instance *T, isNew bool) (Form[T], error)
