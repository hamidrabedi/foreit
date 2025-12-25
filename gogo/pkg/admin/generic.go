package admin

import (
	"context"
	"github.com/gogo/pkg/models"
)

// ModelAdmin is a generic, type-safe admin interface
type ModelAdmin[T models.Model] interface {
	// List configuration
	ListDisplay() []FieldRef[T]
	ListDisplayLinks() []FieldRef[T]
	ListEditable() []FieldRef[T]
	ListFilter() []FilterSpec[T]
	SearchFields() []FieldRef[T]
	
	// Form configuration
	Fields() []FieldRef[T]
	Exclude() []FieldRef[T]
	ReadonlyFields() []FieldRef[T]
	
	// Permissions
	HasAddPermission(ctx context.Context, user interface{}) bool
	HasChangePermission(ctx context.Context, user interface{}, obj *T) bool
	HasDeletePermission(ctx context.Context, user interface{}, obj *T) bool
	HasViewPermission(ctx context.Context, user interface{}, obj *T) bool
	
	// Custom methods
	SaveModel(ctx context.Context, obj *T, form interface{}, change bool) error
	DeleteModel(ctx context.Context, obj *T) error
	GetQueryset(ctx context.Context) models.QuerySet[T]
}

// FieldRef provides type-safe field references for admin
type FieldRef[T models.Model] struct {
	name string
}

// NewFieldRef creates a type-safe admin field reference
func NewFieldRef[T models.Model](name string) *FieldRef[T] {
	return &FieldRef[T]{name: name}
}

// Name returns the field name
func (f *FieldRef[T]) Name() string {
	return f.name
}

// FilterSpec is a type-safe filter specification
type FilterSpec[T models.Model] struct {
	Field    *FieldRef[T]
	Type     FilterType
	Lookups  []Lookup
	Choices  []Choice
	Label    string
}

// BaseModelAdmin provides type-safe default implementations
type BaseModelAdmin[T models.Model] struct {
	modelName string
	
	// List configuration
	listDisplay      []*FieldRef[T]
	listDisplayLinks []*FieldRef[T]
	listEditable     []*FieldRef[T]
	listFilter       []FilterSpec[T]
	searchFields     []*FieldRef[T]
	
	// Form configuration
	fields         []*FieldRef[T]
	exclude        []*FieldRef[T]
	readonlyFields []*FieldRef[T]
}

// NewBaseModelAdmin creates a new type-safe model admin
func NewBaseModelAdmin[T models.Model](modelName string) *BaseModelAdmin[T] {
	return &BaseModelAdmin[T]{
		modelName:       modelName,
		listDisplay:     make([]*FieldRef[T], 0),
		listDisplayLinks: make([]*FieldRef[T], 0),
		listEditable:    make([]*FieldRef[T], 0),
		listFilter:      make([]FilterSpec[T], 0),
		searchFields:    make([]*FieldRef[T], 0),
		fields:          make([]*FieldRef[T], 0),
		exclude:         make([]*FieldRef[T], 0),
		readonlyFields:  make([]*FieldRef[T], 0),
	}
}

// ListDisplay returns list display fields (type-safe!)
func (a *BaseModelAdmin[T]) ListDisplay() []FieldRef[T] {
	result := make([]FieldRef[T], len(a.listDisplay))
	for i, ref := range a.listDisplay {
		result[i] = *ref
	}
	return result
}

// SetListDisplay sets list display fields (type-safe!)
func (a *BaseModelAdmin[T]) SetListDisplay(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.listDisplay = fields
	return a
}

// ListDisplayLinks returns list display links
func (a *BaseModelAdmin[T]) ListDisplayLinks() []FieldRef[T] {
	if len(a.listDisplayLinks) > 0 {
		result := make([]FieldRef[T], len(a.listDisplayLinks))
		for i, ref := range a.listDisplayLinks {
			result[i] = *ref
		}
		return result
	}
	if len(a.listDisplay) > 0 {
		return []FieldRef[T]{*a.listDisplay[0]}
	}
	return []FieldRef[T]{*NewFieldRef[T]("id")}
}

// SetListDisplayLinks sets list display links
func (a *BaseModelAdmin[T]) SetListDisplayLinks(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.listDisplayLinks = fields
	return a
}

// ListEditable returns list editable fields
func (a *BaseModelAdmin[T]) ListEditable() []FieldRef[T] {
	result := make([]FieldRef[T], len(a.listEditable))
	for i, ref := range a.listEditable {
		result[i] = *ref
	}
	return result
}

// SetListEditable sets list editable fields
func (a *BaseModelAdmin[T]) SetListEditable(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.listEditable = fields
	return a
}

// ListFilter returns filter specifications
func (a *BaseModelAdmin[T]) ListFilter() []FilterSpec[T] {
	return a.listFilter
}

// SetListFilter sets list filters (type-safe!)
func (a *BaseModelAdmin[T]) SetListFilter(filters ...FilterSpec[T]) *BaseModelAdmin[T] {
	a.listFilter = filters
	return a
}

// SearchFields returns searchable fields
func (a *BaseModelAdmin[T]) SearchFields() []FieldRef[T] {
	result := make([]FieldRef[T], len(a.searchFields))
	for i, ref := range a.searchFields {
		result[i] = *ref
	}
	return result
}

// SetSearchFields sets search fields (type-safe!)
func (a *BaseModelAdmin[T]) SetSearchFields(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.searchFields = fields
	return a
}

// Fields returns form fields
func (a *BaseModelAdmin[T]) Fields() []FieldRef[T] {
	result := make([]FieldRef[T], len(a.fields))
	for i, ref := range a.fields {
		result[i] = *ref
	}
	return result
}

// SetFields sets form fields (type-safe!)
func (a *BaseModelAdmin[T]) SetFields(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.fields = fields
	return a
}

// Exclude returns excluded fields
func (a *BaseModelAdmin[T]) Exclude() []FieldRef[T] {
	result := make([]FieldRef[T], len(a.exclude))
	for i, ref := range a.exclude {
		result[i] = *ref
	}
	return result
}

// SetExclude sets excluded fields
func (a *BaseModelAdmin[T]) SetExclude(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.exclude = fields
	return a
}

// ReadonlyFields returns readonly fields
func (a *BaseModelAdmin[T]) ReadonlyFields() []FieldRef[T] {
	result := make([]FieldRef[T], len(a.readonlyFields))
	for i, ref := range a.readonlyFields {
		result[i] = *ref
	}
	return result
}

// SetReadonlyFields sets readonly fields
func (a *BaseModelAdmin[T]) SetReadonlyFields(fields ...*FieldRef[T]) *BaseModelAdmin[T] {
	a.readonlyFields = fields
	return a
}

// Permission methods (can be overridden)
func (a *BaseModelAdmin[T]) HasAddPermission(ctx context.Context, user interface{}) bool {
	return true
}

func (a *BaseModelAdmin[T]) HasChangePermission(ctx context.Context, user interface{}, obj *T) bool {
	return true
}

func (a *BaseModelAdmin[T]) HasDeletePermission(ctx context.Context, user interface{}, obj *T) bool {
	return true
}

func (a *BaseModelAdmin[T]) HasViewPermission(ctx context.Context, user interface{}, obj *T) bool {
	return true
}

// Custom methods (can be overridden)
func (a *BaseModelAdmin[T]) SaveModel(ctx context.Context, obj *T, form interface{}, change bool) error {
	return nil
}

func (a *BaseModelAdmin[T]) DeleteModel(ctx context.Context, obj *T) error {
	return nil
}

func (a *BaseModelAdmin[T]) GetQueryset(ctx context.Context) models.QuerySet[T] {
	return nil
}

// ModelName returns the model name
func (a *BaseModelAdmin[T]) ModelName() string {
	return a.modelName
}

// Example usage:
// type User struct { ... }
//
// var UserEmail = admin.NewFieldRef[*User]("email")
// var UserName = admin.NewFieldRef[*User]("name")
// var UserID = admin.NewFieldRef[*User]("id")
//
// var UserAdmin = admin.NewBaseModelAdmin[*User]("User").
//     SetListDisplay(UserID, UserEmail, UserName).
//     SetSearchFields(UserEmail, UserName).
//     SetListFilter(
//         admin.FilterSpec[*User]{
//             Field: UserEmail,
//             Type:  admin.FilterTypeChoice,
//         },
//     )

