package admin

import (
	"context"
	"fmt"
)

// ModelAdmin defines the interface for admin customization (Django-inspired)
// Everything can be overridden
type ModelAdmin interface {
	// List configuration
	ListDisplay() []string
	ListDisplayLinks() []string
	ListEditable() []string
	ListFilter() []FilterSpec
	SearchFields() []string
	SearchHelpText() string
	ListPerPage() int
	ListMaxShowAll() int
	ShowFullResultCount() bool
	
	// Form configuration
	Fields() []string
	Exclude() []string
	Fieldsets() []Fieldset
	ReadonlyFields() []string
	AutocompleteFields() []string
	
	// Actions
	Actions() []Action
	ActionCheckbox() bool
	
	// Permissions
	HasAddPermission(ctx context.Context, user interface{}) bool
	HasChangePermission(ctx context.Context, user interface{}, obj interface{}) bool
	HasDeletePermission(ctx context.Context, user interface{}, obj interface{}) bool
	HasViewPermission(ctx context.Context, user interface{}, obj interface{}) bool
	
	// Custom methods
	SaveModel(ctx context.Context, obj interface{}, form interface{}, change bool) error
	DeleteModel(ctx context.Context, obj interface{}) error
	GetQueryset(ctx context.Context) (interface{}, error)
	
	// Custom views
	ListView(ctx context.Context) (interface{}, error)
	AddView(ctx context.Context) (interface{}, error)
	ChangeView(ctx context.Context, id interface{}) (interface{}, error)
	DeleteView(ctx context.Context, id interface{}) (interface{}, error)
	HistoryView(ctx context.Context, id interface{}) (interface{}, error)
	
	// Inlines
	Inlines() []InlineAdmin
	
	// Ordering
	Ordering() []string
	GetOrdering(ctx context.Context) []string
	
	// Date hierarchy
	DateHierarchy() string
	
	// Empty value display
	EmptyValueDisplay() string
	
	// Model name
	ModelName() string
	VerboseName() string
	VerboseNamePlural() string
}

// BaseModelAdmin provides default implementations that can be overridden
type BaseModelAdmin struct {
	modelName string
	
	// List configuration
	listDisplay        []string
	listDisplayLinks    []string
	listEditable       []string
	listFilter         []FilterSpec
	searchFields       []string
	searchHelpText     string
	listPerPage        int
	listMaxShowAll     int
	showFullResultCount bool
	
	// Form configuration
	fields            []string
	exclude           []string
	fieldsets         []Fieldset
	readonlyFields    []string
	autocompleteFields []string
	
	// Actions
	actions           []Action
	actionCheckbox    bool
	
	// Inlines
	inlines           []InlineAdmin
	
	// Ordering
	ordering          []string
	
	// Date hierarchy
	dateHierarchy     string
	
	// Empty value display
	emptyValueDisplay string
}

// NewBaseModelAdmin creates a new base model admin
func NewBaseModelAdmin(modelName string) *BaseModelAdmin {
	return &BaseModelAdmin{
		modelName:          modelName,
		listPerPage:        100,
		listMaxShowAll:     200,
		showFullResultCount: true,
		actionCheckbox:     true,
		emptyValueDisplay:  "-",
	}
}

// ListDisplay returns fields to display in list view
func (a *BaseModelAdmin) ListDisplay() []string {
	return a.listDisplay
}

// SetListDisplay sets list display fields
func (a *BaseModelAdmin) SetListDisplay(fields ...string) {
	a.listDisplay = fields
}

// ListDisplayLinks returns fields that should be links
func (a *BaseModelAdmin) ListDisplayLinks() []string {
	if len(a.listDisplayLinks) > 0 {
		return a.listDisplayLinks
	}
	if len(a.listDisplay) > 0 {
		return []string{a.listDisplay[0]}
	}
	return []string{"id"}
}

// SetListDisplayLinks sets list display links
func (a *BaseModelAdmin) SetListDisplayLinks(fields ...string) {
	a.listDisplayLinks = fields
}

// ListEditable returns fields that are editable in list view
func (a *BaseModelAdmin) ListEditable() []string {
	return a.listEditable
}

// SetListEditable sets list editable fields
func (a *BaseModelAdmin) SetListEditable(fields ...string) {
	a.listEditable = fields
}

// ListFilter returns filter specifications
func (a *BaseModelAdmin) ListFilter() []FilterSpec {
	return a.listFilter
}

// SetListFilter sets list filters
func (a *BaseModelAdmin) SetListFilter(filters ...FilterSpec) {
	a.listFilter = filters
}

// SearchFields returns searchable fields
func (a *BaseModelAdmin) SearchFields() []string {
	return a.searchFields
}

// SetSearchFields sets search fields
func (a *BaseModelAdmin) SetSearchFields(fields ...string) {
	a.searchFields = fields
}

// SearchHelpText returns search help text
func (a *BaseModelAdmin) SearchHelpText() string {
	return a.searchHelpText
}

// SetSearchHelpText sets search help text
func (a *BaseModelAdmin) SetSearchHelpText(text string) {
	a.searchHelpText = text
}

// ListPerPage returns items per page
func (a *BaseModelAdmin) ListPerPage() int {
	return a.listPerPage
}

// SetListPerPage sets items per page
func (a *BaseModelAdmin) SetListPerPage(n int) {
	a.listPerPage = n
}

// ListMaxShowAll returns max items to show all
func (a *BaseModelAdmin) ListMaxShowAll() int {
	return a.listMaxShowAll
}

// SetListMaxShowAll sets max show all
func (a *BaseModelAdmin) SetListMaxShowAll(n int) {
	a.listMaxShowAll = n
}

// ShowFullResultCount returns whether to show full result count
func (a *BaseModelAdmin) ShowFullResultCount() bool {
	return a.showFullResultCount
}

// SetShowFullResultCount sets show full result count
func (a *BaseModelAdmin) SetShowFullResultCount(b bool) {
	a.showFullResultCount = b
}

// Fields returns form fields
func (a *BaseModelAdmin) Fields() []string {
	return a.fields
}

// SetFields sets form fields
func (a *BaseModelAdmin) SetFields(fields ...string) {
	a.fields = fields
}

// Exclude returns excluded fields
func (a *BaseModelAdmin) Exclude() []string {
	return a.exclude
}

// SetExclude sets excluded fields
func (a *BaseModelAdmin) SetExclude(fields ...string) {
	a.exclude = fields
}

// Fieldsets returns fieldset configuration
func (a *BaseModelAdmin) Fieldsets() []Fieldset {
	return a.fieldsets
}

// SetFieldsets sets fieldsets
func (a *BaseModelAdmin) SetFieldsets(fieldsets ...Fieldset) {
	a.fieldsets = fieldsets
}

// ReadonlyFields returns readonly fields
func (a *BaseModelAdmin) ReadonlyFields() []string {
	return a.readonlyFields
}

// SetReadonlyFields sets readonly fields
func (a *BaseModelAdmin) SetReadonlyFields(fields ...string) {
	a.readonlyFields = fields
}

// AutocompleteFields returns autocomplete fields
func (a *BaseModelAdmin) AutocompleteFields() []string {
	return a.autocompleteFields
}

// SetAutocompleteFields sets autocomplete fields
func (a *BaseModelAdmin) SetAutocompleteFields(fields ...string) {
	a.autocompleteFields = fields
}

// Actions returns custom actions
func (a *BaseModelAdmin) Actions() []Action {
	return a.actions
}

// SetActions sets actions
func (a *BaseModelAdmin) SetActions(actions ...Action) {
	a.actions = actions
}

// ActionCheckbox returns whether to show action checkbox
func (a *BaseModelAdmin) ActionCheckbox() bool {
	return a.actionCheckbox
}

// SetActionCheckbox sets action checkbox
func (a *BaseModelAdmin) SetActionCheckbox(b bool) {
	a.actionCheckbox = b
}

// Inlines returns inline admins
func (a *BaseModelAdmin) Inlines() []InlineAdmin {
	return a.inlines
}

// SetInlines sets inlines
func (a *BaseModelAdmin) SetInlines(inlines ...InlineAdmin) {
	a.inlines = inlines
}

// Ordering returns default ordering
func (a *BaseModelAdmin) Ordering() []string {
	return a.ordering
}

// SetOrdering sets ordering
func (a *BaseModelAdmin) SetOrdering(fields ...string) {
	a.ordering = fields
}

// GetOrdering returns ordering for a request
func (a *BaseModelAdmin) GetOrdering(ctx context.Context) []string {
	return a.ordering
}

// DateHierarchy returns date hierarchy field
func (a *BaseModelAdmin) DateHierarchy() string {
	return a.dateHierarchy
}

// SetDateHierarchy sets date hierarchy
func (a *BaseModelAdmin) SetDateHierarchy(field string) {
	a.dateHierarchy = field
}

// EmptyValueDisplay returns empty value display
func (a *BaseModelAdmin) EmptyValueDisplay() string {
	return a.emptyValueDisplay
}

// SetEmptyValueDisplay sets empty value display
func (a *BaseModelAdmin) SetEmptyValueDisplay(s string) {
	a.emptyValueDisplay = s
}

// ModelName returns the model name
func (a *BaseModelAdmin) ModelName() string {
	return a.modelName
}

// VerboseName returns verbose name
func (a *BaseModelAdmin) VerboseName() string {
	return a.modelName
}

// VerboseNamePlural returns plural verbose name
func (a *BaseModelAdmin) VerboseNamePlural() string {
	return fmt.Sprintf("%ss", a.modelName)
}

// Permission methods (can be overridden)
func (a *BaseModelAdmin) HasAddPermission(ctx context.Context, user interface{}) bool {
	return true
}

func (a *BaseModelAdmin) HasChangePermission(ctx context.Context, user interface{}, obj interface{}) bool {
	return true
}

func (a *BaseModelAdmin) HasDeletePermission(ctx context.Context, user interface{}, obj interface{}) bool {
	return true
}

func (a *BaseModelAdmin) HasViewPermission(ctx context.Context, user interface{}, obj interface{}) bool {
	return true
}

// Custom methods (can be overridden)
func (a *BaseModelAdmin) SaveModel(ctx context.Context, obj interface{}, form interface{}, change bool) error {
	return nil
}

func (a *BaseModelAdmin) DeleteModel(ctx context.Context, obj interface{}) error {
	return nil
}

func (a *BaseModelAdmin) GetQueryset(ctx context.Context) (interface{}, error) {
	return nil, nil
}

// Custom views (can be overridden)
func (a *BaseModelAdmin) ListView(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func (a *BaseModelAdmin) AddView(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func (a *BaseModelAdmin) ChangeView(ctx context.Context, id interface{}) (interface{}, error) {
	return nil, nil
}

func (a *BaseModelAdmin) DeleteView(ctx context.Context, id interface{}) (interface{}, error) {
	return nil, nil
}

func (a *BaseModelAdmin) HistoryView(ctx context.Context, id interface{}) (interface{}, error) {
	return nil, nil
}

// FilterSpec represents a filter specification
type FilterSpec struct {
	Field    string
	Type     FilterType
	Lookups  []Lookup
	Choices  []Choice
	Label    string
}

// FilterType represents filter type
type FilterType string

const (
	FilterTypeList      FilterType = "list"
	FilterTypeDate      FilterType = "date"
	FilterTypeDateTime  FilterType = "datetime"
	FilterTypeBoolean   FilterType = "boolean"
	FilterTypeChoice    FilterType = "choice"
	FilterTypeRelated   FilterType = "related"
	FilterTypeRelatedOnly FilterType = "related_only"
)

// Lookup represents a lookup type
type Lookup string

const (
	LookupExact      Lookup = "exact"
	LookupIExact     Lookup = "iexact"
	LookupContains   Lookup = "contains"
	LookupIContains  Lookup = "icontains"
	LookupIn         Lookup = "in"
	LookupGt         Lookup = "gt"
	LookupGte        Lookup = "gte"
	LookupLt         Lookup = "lt"
	LookupLte        Lookup = "lte"
	LookupStartswith Lookup = "startswith"
	LookupIStartswith Lookup = "istartswith"
	LookupEndswith   Lookup = "endswith"
	LookupIEndswith  Lookup = "iendswith"
	LookupRange      Lookup = "range"
	LookupDate       Lookup = "date"
	LookupYear       Lookup = "year"
	LookupMonth      Lookup = "month"
	LookupDay        Lookup = "day"
	LookupWeek       Lookup = "week"
	LookupWeekDay    Lookup = "week_day"
	LookupTime       Lookup = "time"
	LookupHour       Lookup = "hour"
	LookupMinute     Lookup = "minute"
	LookupSecond     Lookup = "second"
	LookupIsNull     Lookup = "isnull"
	LookupRegex      Lookup = "regex"
	LookupIRegex     Lookup = "iregex"
)

// Choice represents a choice option
type Choice struct {
	Value interface{}
	Label string
}

// Fieldset represents a fieldset in forms
type Fieldset struct {
	Name      string
	Fields    []string
	Classes   []string
	Collapsed bool
}

// Action represents a custom admin action
type Action struct {
	Name        string
	Label       string
	Description string
	Handler     func(ctx context.Context, queryset interface{}) error
	Permissions []string
}

// InlineAdmin represents an inline admin
type InlineAdmin interface {
	Model() string
	Fields() []string
	Extra() int
	MaxNum() int
	MinNum() int
	CanDelete() bool
	VerboseName() string
	VerboseNamePlural() string
}

