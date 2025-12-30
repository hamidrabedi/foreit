package fields

import (
	"github.com/forgego/forge/admin/widgets"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// Field represents a type-safe admin field that integrates with schema and ORM
// T is the model type, V is the field value type
type Field[T any, V any] interface {
	// Name returns the field name
	Name() string

	// Get gets the field value from an instance (type-safe)
	Get(instance *T) V

	// Set sets the field value on an instance (type-safe)
	Set(instance *T, value V) error

	// GetFieldExpression returns the ORM field expression for queries
	GetFieldExpression() orm.FieldExpression[V]

	// GetSchemaField returns the schema field definition
	GetSchemaField() schema.Field

	// GetWidget returns the widget for this field
	GetWidget() widgets.Widget

	// IsRequired returns if the field is required
	IsRequired() bool

	// IsReadOnly returns if the field is read-only
	IsReadOnly() bool

	// GetVerboseName returns the human-readable name
	GetVerboseName() string

	// GetHelpText returns the help text
	GetHelpText() string
}

// BaseField provides base implementation for fields
type BaseField[T any, V any] struct {
	name          string
	schemaField   schema.Field
	fieldExpr     orm.FieldExpression[V]
	widget        widgets.Widget
	accessor      *orm.FieldAccessor[T]
	readOnly      bool
	verboseName   string
	helpText      string
}

// NewBaseField creates a new base field
func NewBaseField[T any, V any](
	name string,
	schemaField schema.Field,
	fieldExpr orm.FieldExpression[V],
	accessor *orm.FieldAccessor[T],
	widget widgets.Widget,
) *BaseField[T, V] {
	verboseName := schemaField.VerboseName
	if verboseName == "" {
		verboseName = name
	}

	return &BaseField[T, V]{
		name:        name,
		schemaField: schemaField,
		fieldExpr:   fieldExpr,
		accessor:    accessor,
		widget:      widget,
		readOnly:    !schemaField.Editable,
		verboseName: verboseName,
		helpText:    schemaField.HelpText,
	}
}

// Name returns the field name
func (f *BaseField[T, V]) Name() string {
	return f.name
}

// GetSchemaField returns the schema field
func (f *BaseField[T, V]) GetSchemaField() schema.Field {
	return f.schemaField
}

// GetFieldExpression returns the ORM field expression
func (f *BaseField[T, V]) GetFieldExpression() orm.FieldExpression[V] {
	return f.fieldExpr
}

// GetWidget returns the widget
func (f *BaseField[T, V]) GetWidget() widgets.Widget {
	return f.widget
}

// IsRequired returns if the field is required
func (f *BaseField[T, V]) IsRequired() bool {
	return f.schemaField.Required
}

// IsReadOnly returns if the field is read-only
func (f *BaseField[T, V]) IsReadOnly() bool {
	return f.readOnly
}

// GetVerboseName returns the verbose name
func (f *BaseField[T, V]) GetVerboseName() string {
	return f.verboseName
}

// GetHelpText returns the help text
func (f *BaseField[T, V]) GetHelpText() string {
	return f.helpText
}

// SetReadOnly sets the read-only status
func (f *BaseField[T, V]) SetReadOnly(readOnly bool) {
	f.readOnly = readOnly
}

// SetWidget sets the widget
func (f *BaseField[T, V]) SetWidget(widget widgets.Widget) {
	f.widget = widget
}
