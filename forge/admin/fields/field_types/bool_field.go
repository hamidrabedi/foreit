package fieldtypes

import (
	"github.com/forgego/forge/admin/fields"
	"github.com/forgego/forge/admin/widgets"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// BoolField is a specialized field for boolean values
type BoolField[T any] struct {
	*fields.FieldExpr[T, bool]
}

// NewBoolField creates a new boolean field
func NewBoolField[T any](
	name string,
	schemaField schema.Field,
	accessor *orm.FieldAccessor[T],
) (*BoolField[T], error) {
	widget := widgets.NewCheckbox()

	fieldExpr, err := fields.NewFieldExpr[T, bool](name, schemaField, accessor, widget)
	if err != nil {
		return nil, err
	}

	return &BoolField[T]{
		FieldExpr: fieldExpr,
	}, nil
}
