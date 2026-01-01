package filter

import (
	"context"
	"fmt"
	"strings"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/filter/filters"
	"github.com/forgego/forge/orm"
)

// AdminMultipleChoiceFilter is an admin-specific wrapper for multiple choice filters
// It provides enhanced UI features for the admin interface
type AdminMultipleChoiceFilter[T any] struct {
	*filters.MultipleChoiceFilter[T]
	allowMultiple bool
}

// NewAdminMultipleChoiceFilter creates a new admin multiple choice filter
func NewAdminMultipleChoiceFilter[T any](
	fieldPath string,
	choices []filters.Choice,
) *AdminMultipleChoiceFilter[T] {
	baseFilter := filters.NewMultipleChoiceFilter[T](fieldPath, choices)
	return &AdminMultipleChoiceFilter[T]{
		MultipleChoiceFilter: baseFilter,
		allowMultiple:        true,
	}
}

// Parse parses query parameters for multiple values
func (f *AdminMultipleChoiceFilter[T]) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	// Support comma-separated values or array format
	values := strings.Split(value, ",")
	result := make([]interface{}, 0, len(values))

	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		result = append(result, v)
	}

	if len(result) == 0 {
		return nil, nil
	}

	if len(result) == 1 && !f.allowMultiple {
		return result[0], nil
	}

	return result, nil
}

// Apply applies the multiple choice filter to a queryset
func (f *AdminMultipleChoiceFilter[T]) Apply(ctx context.Context, qs orm.QuerySet[T], value interface{}) (orm.QuerySet[T], error) {
	if value == nil {
		return qs, nil
	}

	fieldPath := f.GetFieldPath()
	fieldExpr := orm.NewField[string](fieldPath, "")

	// Handle single value
	if singleValue, ok := value.(string); ok {
		expr := fieldExpr.Eq(singleValue)
		return qs.Filter(expr), nil
	}

	// Handle multiple values
	if values, ok := value.([]interface{}); ok && len(values) > 0 {
		// Convert to []string
		strValues := make([]string, 0, len(values))
		for _, v := range values {
			if str, ok := v.(string); ok {
				strValues = append(strValues, str)
			}
		}

		if len(strValues) > 0 {
			expr := fieldExpr.In(strValues...)
			return qs.Filter(expr), nil
		}
	}

	return qs, nil
}

// ToExpression converts the filter value to an ORM expression
func (f *AdminMultipleChoiceFilter[T]) ToExpression(fieldPath string, value interface{}) (orm.Expression, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot create expression for nil value")
	}

	fieldExpr := orm.NewField[string](fieldPath, "")

	// Handle single value
	if singleValue, ok := value.(string); ok {
		return fieldExpr.Eq(singleValue), nil
	}

	// Handle multiple values
	if values, ok := value.([]interface{}); ok && len(values) > 0 {
		strValues := make([]string, 0, len(values))
		for _, v := range values {
			if str, ok := v.(string); ok {
				strValues = append(strValues, str)
			}
		}

		if len(strValues) > 0 {
			return fieldExpr.In(strValues...), nil
		}
	}

	return nil, fmt.Errorf("multiple choice filter value must be string or []interface{}, got %T", value)
}

// GetWidget returns the widget for this filter
func (f *AdminMultipleChoiceFilter[T]) GetWidget() filter.Widget {
	// Get choices from the underlying filter via GetOptions
	return &AdminMultipleChoiceWidget{
		allowMultiple: f.allowMultiple,
		fieldPath:     f.GetFieldPath(),
	}
}

// AdminMultipleChoiceWidget is a widget for admin multiple choice filters
type AdminMultipleChoiceWidget struct {
	allowMultiple bool
	fieldPath     string
}

// Type returns the widget type
func (w *AdminMultipleChoiceWidget) Type() string {
	return "admin_multiple_choice"
}

// Render renders the widget HTML as a multi-select or checkbox group
// Note: Choices will be populated via JavaScript from GetOptions
func (w *AdminMultipleChoiceWidget) Render(name string, value interface{}, attrs map[string]string) (string, error) {
	var selectedValues []string

	if value != nil {
		if singleValue, ok := value.(string); ok {
			selectedValues = []string{singleValue}
		} else if values, ok := value.([]interface{}); ok {
			selectedValues = make([]string, 0, len(values))
			for _, v := range values {
				if str, ok := v.(string); ok {
					selectedValues = append(selectedValues, str)
				}
			}
		}
	}

	selectedStr := strings.Join(selectedValues, ",")

	if w.allowMultiple {
		// Multi-select dropdown - will be populated by JavaScript
		html := fmt.Sprintf(
			`<select name="%s" multiple class="form-control admin-multiple-choice-select" size="5" data-field="%s" data-selected="%s">
				<option value="">Loading options...</option>
			</select>`,
			name, w.fieldPath, selectedStr,
		)
		return html, nil
	}

	// Checkbox group - will be populated by JavaScript
	html := fmt.Sprintf(
		`<div class="admin-multiple-choice-checkboxes" data-field="%s" data-selected="%s">
			<!-- Options will be loaded via JavaScript -->
		</div>`,
		w.fieldPath, selectedStr,
	)

	return html, nil
}

// Parse parses the widget values
func (w *AdminMultipleChoiceWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	if w.allowMultiple {
		// Multiple values separated by comma
		values := strings.Split(value, ",")
		result := make([]interface{}, 0, len(values))
		for _, v := range values {
			v = strings.TrimSpace(v)
			if v != "" {
				result = append(result, v)
			}
		}
		return result, nil
	}

	return value, nil
}
