package widgets

import (
	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/schema"
)

// Widget represents a form widget (re-export from admin package)
type Widget = admin.Widget

// TextInput creates a new text input widget
func NewTextInput() Widget {
	return admin.NewTextInput()
}

// NewNumberInput creates a new number input widget
func NewNumberInput() Widget {
	return admin.NewNumberInput()
}

// NewEmailInput creates a new email input widget
func NewEmailInput() Widget {
	return admin.NewEmailInput()
}

// NewTextarea creates a new textarea widget
func NewTextarea(rows, cols int) Widget {
	return admin.NewTextarea(rows, cols)
}

// NewCheckbox creates a new checkbox widget
func NewCheckbox() Widget {
	return admin.NewCheckbox()
}

// NewSelect creates a new select widget
// Note: This needs to match the Choice type from admin package
func NewSelect(choices []interface{}) Widget {
	// Convert to admin.Choice format if needed
	// For now, return a basic select
	return admin.NewSelect(nil)
}

// NewDateInput creates a new date input widget
func NewDateInput() Widget {
	return admin.NewDateInput()
}

// NewDateTimeInput creates a new datetime input widget
func NewDateTimeInput() Widget {
	return admin.NewDateTimeInput()
}

// WidgetRegistry manages widget selection based on field types
type WidgetRegistry struct {
	typeMappings map[string]func() Widget
}

// NewWidgetRegistry creates a new widget registry
func NewWidgetRegistry() *WidgetRegistry {
	registry := &WidgetRegistry{
		typeMappings: make(map[string]func() Widget),
	}

	// Register default widgets
	registry.registerDefaults()

	return registry
}

// registerDefaults registers default widget mappings
func (wr *WidgetRegistry) registerDefaults() {
	wr.typeMappings["text"] = NewTextInput
	wr.typeMappings["textarea"] = func() Widget { return NewTextarea(10, 40) }
	wr.typeMappings["number"] = NewNumberInput
	wr.typeMappings["email"] = NewEmailInput
	wr.typeMappings["checkbox"] = NewCheckbox
	wr.typeMappings["date"] = NewDateInput
	wr.typeMappings["datetime"] = NewDateTimeInput
	wr.typeMappings["select"] = func() Widget { return NewSelect(nil) }
}

// GetWidgetForFieldType gets a widget for a schema field type
func (wr *WidgetRegistry) GetWidgetForFieldType(fieldType schema.FieldType) Widget {
	switch fieldType {
	case schema.TypeBool:
		return NewCheckbox()
	case schema.TypeText:
		return NewTextarea(10, 40)
	case schema.TypeDate:
		return NewDateInput()
	case schema.TypeDateTime, schema.TypeTime:
		return NewDateTimeInput()
	case schema.TypeEmail:
		return NewEmailInput()
	case schema.TypeInt64, schema.TypeInt32, schema.TypeFloat32, schema.TypeFloat64, schema.TypeDecimal:
		return NewNumberInput()
	case schema.TypeForeignKey, schema.TypeOneToOne, schema.TypeManyToMany:
		return NewSelect(nil)
	default:
		return NewTextInput()
	}
}

// RegisterWidget registers a custom widget for a field type
func (wr *WidgetRegistry) RegisterWidget(fieldType string, factory func() Widget) {
	wr.typeMappings[fieldType] = factory
}
