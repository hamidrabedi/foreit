package fieldtypes

import (
	"github.com/forgego/forge/admin/fields"
	"github.com/forgego/forge/admin/widgets"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// NumberField is a specialized field for numeric values
type NumberField[T any] struct {
	*fields.FieldExpr[T, int64]
}

// NewNumberField creates a new number field
func NewNumberField[T any](
	name string,
	schemaField schema.Field,
	accessor *orm.FieldAccessor[T],
) (*NumberField[T], error) {
	widget := widgets.NewNumberInput()

	fieldExpr, err := fields.NewFieldExpr[T, int64](name, schemaField, accessor, widget)
	if err != nil {
		return nil, err
	}

	return &NumberField[T]{
		FieldExpr: fieldExpr,
	}, nil
}
