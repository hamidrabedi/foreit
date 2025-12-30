package fieldtypes

import (
	"time"

	"github.com/forgego/forge/admin/fields"
	"github.com/forgego/forge/admin/widgets"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// DateField is a specialized field for date/time values
type DateField[T any] struct {
	*fields.FieldExpr[T, time.Time]
}

// NewDateField creates a new date field
func NewDateField[T any](
	name string,
	schemaField schema.Field,
	accessor *orm.FieldAccessor[T],
) (*DateField[T], error) {
	var widget widgets.Widget
	switch schemaField.Type {
	case schema.TypeDate:
		widget = widgets.NewDateInput()
	case schema.TypeDateTime, schema.TypeTime:
		widget = widgets.NewDateTimeInput()
	default:
		widget = widgets.NewDateInput()
	}

	fieldExpr, err := fields.NewFieldExpr[T, time.Time](name, schemaField, accessor, widget)
	if err != nil {
		return nil, err
	}

	return &DateField[T]{
		FieldExpr: fieldExpr,
	}, nil
}
