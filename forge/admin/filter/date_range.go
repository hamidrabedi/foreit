package filter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/filter/filters"
	"github.com/forgego/forge/orm"
)

// DateRangeFilter filters date fields with start and end dates
type DateRangeFilter[T any] struct {
	*filters.DateFilter[T]
	startField string
	endField   string
}

// NewDateRangeFilter creates a new date range filter
func NewDateRangeFilter[T any](fieldPath string) *DateRangeFilter[T] {
	baseFilter := filters.NewDateFilter[T](fieldPath)
	return &DateRangeFilter[T]{
		DateFilter: baseFilter,
		startField: fieldPath + "_start",
		endField:   fieldPath + "_end",
	}
}

// Parse parses query parameters for date range (start and end)
func (f *DateRangeFilter[T]) Parse(value string) (interface{}, error) {
	// DateRangeFilter expects two values separated by comma or as separate params
	// Format: "start,end" or separate params
	if value == "" {
		return nil, nil
	}

	parts := strings.Split(value, ",")
	if len(parts) == 2 {
		start, err1 := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
		end, err2 := time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid date range format, expected 'YYYY-MM-DD,YYYY-MM-DD'")
		}
		return []time.Time{start, end}, nil
	}

	// Single value - treat as start date
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	return []time.Time{t, time.Time{}}, nil
}

// Apply applies the date range filter to a queryset
func (f *DateRangeFilter[T]) Apply(ctx context.Context, qs orm.QuerySet[T], value interface{}) (orm.QuerySet[T], error) {
	if value == nil {
		return qs, nil
	}

	rangeValues, ok := value.([]time.Time)
	if !ok {
		return nil, fmt.Errorf("date range filter value must be []time.Time, got %T", value)
	}

	if len(rangeValues) < 2 {
		return qs, nil
	}

	fieldPath := f.GetFieldPath()
	fieldExpr := orm.NewField[time.Time](fieldPath, "")

	// Apply start date (gte)
	if !rangeValues[0].IsZero() {
		expr := fieldExpr.Gte(rangeValues[0])
		qs = qs.Filter(expr)
	}

	// Apply end date (lte)
	if !rangeValues[1].IsZero() {
		expr := fieldExpr.Lte(rangeValues[1])
		qs = qs.Filter(expr)
	}

	return qs, nil
}

// ToExpression converts the date range to ORM expressions
func (f *DateRangeFilter[T]) ToExpression(fieldPath string, value interface{}) (orm.Expression, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot create expression for nil value")
	}

	rangeValues, ok := value.([]time.Time)
	if !ok || len(rangeValues) < 2 {
		return nil, fmt.Errorf("date range filter value must be []time.Time with 2 values")
	}

	fieldExpr := orm.NewField[time.Time](fieldPath, "")

	// Create range expression
	if !rangeValues[0].IsZero() && !rangeValues[1].IsZero() {
		return fieldExpr.Range(rangeValues[0], rangeValues[1]), nil
	}

	// If only start date, use gte
	if !rangeValues[0].IsZero() {
		return fieldExpr.Gte(rangeValues[0]), nil
	}

	// If only end date, use lte
	if !rangeValues[1].IsZero() {
		return fieldExpr.Lte(rangeValues[1]), nil
	}

	return nil, fmt.Errorf("at least one date must be provided for date range")
}

// GetStartField returns the start field name
func (f *DateRangeFilter[T]) GetStartField() string {
	return f.startField
}

// GetEndField returns the end field name
func (f *DateRangeFilter[T]) GetEndField() string {
	return f.endField
}

// GetWidget returns the widget for this filter
func (f *DateRangeFilter[T]) GetWidget() filter.Widget {
	return &DateRangeWidget{
		startField: f.startField,
		endField:   f.endField,
	}
}

// DateRangeWidget is a widget for date range filters
type DateRangeWidget struct {
	startField string
	endField   string
}

// Type returns the widget type
func (w *DateRangeWidget) Type() string {
	return "date_range"
}

// Render renders the widget HTML with two date inputs
func (w *DateRangeWidget) Render(name string, value interface{}, attrs map[string]string) (string, error) {
	var startValue, endValue string

	if value != nil {
		if rangeValues, ok := value.([]time.Time); ok && len(rangeValues) >= 2 {
			if !rangeValues[0].IsZero() {
				startValue = rangeValues[0].Format("2006-01-02")
			}
			if !rangeValues[1].IsZero() {
				endValue = rangeValues[1].Format("2006-01-02")
			}
		}
	}

	html := fmt.Sprintf(`
		<div class="date-range-filter">
			<label>From:</label>
			<input type="date" name="%s" value="%s" class="form-control date-range-start" />
			<label>To:</label>
			<input type="date" name="%s" value="%s" class="form-control date-range-end" />
		</div>`,
		w.startField, startValue,
		w.endField, endValue,
	)

	return html, nil
}

// Parse parses the widget values
func (w *DateRangeWidget) Parse(value string) (interface{}, error) {
	// This will be called for each field separately
	// The actual parsing happens in the filter's Parse method
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	return t, nil
}
