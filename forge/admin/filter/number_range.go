package filter

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/filter/filters"
	"github.com/forgego/forge/orm"
)

// NumberRangeFilter filters numeric fields with min and max values
type NumberRangeFilter[T any] struct {
	*filters.NumberFilter[T]
	minField string
	maxField string
}

// NewNumberRangeFilter creates a new number range filter
func NewNumberRangeFilter[T any](fieldPath string) *NumberRangeFilter[T] {
	baseFilter := filters.NewNumberFilter[T](fieldPath)
	return &NumberRangeFilter[T]{
		NumberFilter: baseFilter,
		minField:     fieldPath + "_min",
		maxField:     fieldPath + "_max",
	}
}

// Parse parses query parameters for number range (min and max)
func (f *NumberRangeFilter[T]) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	// Format: "min,max" or separate params
	parts := strings.Split(value, ",")
	if len(parts) == 2 {
		min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid number range format, expected 'min,max'")
		}
		return []float64{min, max}, nil
	}

	// Single value - treat as min
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %w", err)
	}
	return []float64{num, 0}, nil
}

// Apply applies the number range filter to a queryset
func (f *NumberRangeFilter[T]) Apply(ctx context.Context, qs orm.QuerySet[T], value interface{}) (orm.QuerySet[T], error) {
	if value == nil {
		return qs, nil
	}

	rangeValues, ok := value.([]float64)
	if !ok {
		return nil, fmt.Errorf("number range filter value must be []float64, got %T", value)
	}

	if len(rangeValues) < 2 {
		return qs, nil
	}

	fieldPath := f.GetFieldPath()
	fieldExpr := orm.NewField[float64](fieldPath, "")

	// Apply min value (gte)
	if rangeValues[0] != 0 {
		expr := fieldExpr.Gte(rangeValues[0])
		qs = qs.Filter(expr)
	}

	// Apply max value (lte)
	if rangeValues[1] != 0 {
		expr := fieldExpr.Lte(rangeValues[1])
		qs = qs.Filter(expr)
	}

	return qs, nil
}

// ToExpression converts the number range to ORM expressions
func (f *NumberRangeFilter[T]) ToExpression(fieldPath string, value interface{}) (orm.Expression, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot create expression for nil value")
	}

	rangeValues, ok := value.([]float64)
	if !ok || len(rangeValues) < 2 {
		return nil, fmt.Errorf("number range filter value must be []float64 with 2 values")
	}

	fieldExpr := orm.NewField[float64](fieldPath, "")

	// Create range expression
	if rangeValues[0] != 0 && rangeValues[1] != 0 {
		return fieldExpr.Range(rangeValues[0], rangeValues[1]), nil
	}

	// If only min value, use gte
	if rangeValues[0] != 0 {
		return fieldExpr.Gte(rangeValues[0]), nil
	}

	// If only max value, use lte
	if rangeValues[1] != 0 {
		return fieldExpr.Lte(rangeValues[1]), nil
	}

	return nil, fmt.Errorf("at least one number must be provided for number range")
}

// GetMinField returns the min field name
func (f *NumberRangeFilter[T]) GetMinField() string {
	return f.minField
}

// GetMaxField returns the max field name
func (f *NumberRangeFilter[T]) GetMaxField() string {
	return f.maxField
}

// GetWidget returns the widget for this filter
func (f *NumberRangeFilter[T]) GetWidget() filter.Widget {
	return &NumberRangeWidget{
		minField: f.minField,
		maxField: f.maxField,
	}
}

// NumberRangeWidget is a widget for number range filters
type NumberRangeWidget struct {
	minField string
	maxField string
}

// Type returns the widget type
func (w *NumberRangeWidget) Type() string {
	return "number_range"
}

// Render renders the widget HTML with two number inputs
func (w *NumberRangeWidget) Render(name string, value interface{}, attrs map[string]string) (string, error) {
	var minValue, maxValue string

	if value != nil {
		if rangeValues, ok := value.([]float64); ok && len(rangeValues) >= 2 {
			if rangeValues[0] != 0 {
				minValue = fmt.Sprintf("%.2f", rangeValues[0])
			}
			if rangeValues[1] != 0 {
				maxValue = fmt.Sprintf("%.2f", rangeValues[1])
			}
		}
	}

	html := fmt.Sprintf(`
		<div class="number-range-filter">
			<label>Min:</label>
			<input type="number" name="%s" value="%s" step="any" class="form-control number-range-min" />
			<label>Max:</label>
			<input type="number" name="%s" value="%s" step="any" class="form-control number-range-max" />
		</div>`,
		w.minField, minValue,
		w.maxField, maxValue,
	)

	return html, nil
}

// Parse parses the widget values
func (w *NumberRangeWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %w", err)
	}

	return num, nil
}
