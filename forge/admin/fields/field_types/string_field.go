package fieldtypes

import (
	"github.com/forgego/forge/admin/fields"
	"github.com/forgego/forge/admin/widgets"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// StringField is a specialized field for string values
type StringField[T any] struct {
	*fields.FieldExpr[T, string]
}

// NewStringField creates a new string field
func NewStringField[T any](
	name string,
	schemaField schema.Field,
	accessor *orm.FieldAccessor[T],
) (*StringField[T], error) {
	widget := widgets.NewTextInput()
	if schemaField.Type == schema.TypeText {
		widget = widgets.NewTextarea(10, 40)
	} else if schemaField.Type == schema.TypeEmail {
		widget = widgets.NewEmailInput()
	}

	fieldExpr, err := fields.NewFieldExpr[T, string](name, schemaField, accessor, widget)
	if err != nil {
		return nil, err
	}

	return &StringField[T]{
		FieldExpr: fieldExpr,
	}, nil
}
