package filter

import (
	"context"
	"fmt"
	"reflect"

	"github.com/forgego/forge/filter"
	"github.com/forgego/forge/orm"
)

// RelatedModelFilter filters foreign key and related model fields
type RelatedModelFilter[T any, R any] struct {
	*filter.BaseFilter[T]
	relatedManager *orm.Manager[R]
	displayField   string // Field to display in the dropdown (e.g., "name", "title")
	searchFields   []string // Fields to search in the related model
}

// NewRelatedModelFilter creates a new related model filter
func NewRelatedModelFilter[T any, R any](
	fieldPath string,
	relatedManager *orm.Manager[R],
	displayField string,
) *RelatedModelFilter[T, R] {
	return &RelatedModelFilter[T, R]{
		BaseFilter:     filter.NewBaseFilter[T](fieldPath, "exact"),
		relatedManager: relatedManager,
		displayField:   displayField,
		searchFields:   []string{displayField}, // Default to display field
	}
}

// WithSearchFields sets the fields to search in the related model
func (f *RelatedModelFilter[T, R]) WithSearchFields(fields ...string) *RelatedModelFilter[T, R] {
	f.searchFields = fields
	return f
}

// Parse parses a query parameter value (expects related model ID)
func (f *RelatedModelFilter[T, R]) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	// Try to parse as int64 (ID)
	var id int64
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		return id, nil
	}

	// Could also support searching by display field value
	// For now, return the value as-is and let Apply handle it
	return value, nil
}

// Apply applies the related model filter to a queryset
func (f *RelatedModelFilter[T, R]) Apply(ctx context.Context, qs orm.QuerySet[T], value interface{}) (orm.QuerySet[T], error) {
	if value == nil {
		return qs, nil
	}

	fieldPath := f.GetFieldPath()
	fieldExpr := orm.NewField[int64](fieldPath, "")

	// If value is int64, use exact match
	if id, ok := value.(int64); ok {
		expr := fieldExpr.Eq(id)
		return qs.Filter(expr), nil
	}

	// If value is string, try to find related model by display field
	if strValue, ok := value.(string); ok {
		// Use Manager.Filter to create queryset
		displayFieldExpr := orm.NewField[string](f.displayField, "")
		searchExpr := displayFieldExpr.IContains(strValue)
		relatedQs, err := f.relatedManager.Filter(searchExpr)
		if err != nil {
			return nil, fmt.Errorf("failed to filter related model: %w", err)
		}

		relatedInstances, err := relatedQs.All(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to query related model: %w", err)
		}

		if len(relatedInstances) == 0 {
			// No matches, return empty queryset
			return qs.Filter(orm.NewField[int64](fieldPath, "").Eq(-1)), nil
		}

		// Extract IDs
		ids := make([]int64, 0, len(relatedInstances))
		for _, inst := range relatedInstances {
			id := f.getID(inst)
			if id > 0 {
				ids = append(ids, id)
			}
		}

		if len(ids) == 0 {
			return qs.Filter(orm.NewField[int64](fieldPath, "").Eq(-1)), nil
		}

		// Filter by IDs
		expr := fieldExpr.In(ids...)
		return qs.Filter(expr), nil
	}

	return nil, fmt.Errorf("related model filter value must be int64 or string, got %T", value)
}

// ToExpression converts the filter value to an ORM expression
func (f *RelatedModelFilter[T, R]) ToExpression(fieldPath string, value interface{}) (orm.Expression, error) {
	if value == nil {
		return nil, fmt.Errorf("cannot create expression for nil value")
	}

	fieldExpr := orm.NewField[int64](fieldPath, "")

	if id, ok := value.(int64); ok {
		return fieldExpr.Eq(id), nil
	}

	return nil, fmt.Errorf("related model filter value must be int64, got %T", value)
}

// GetOptions returns filter options by querying the related model
func (f *RelatedModelFilter[T, R]) GetOptions(ctx context.Context, qs orm.QuerySet[T]) ([]filter.FilterOption, error) {
	// Get all related instances using Manager.All
	relatedInstances, err := f.relatedManager.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query related model: %w", err)
	}

	options := make([]filter.FilterOption, 0, len(relatedInstances))
	for _, inst := range relatedInstances {
		id := f.getID(inst)
		display := f.getDisplayValue(inst)
		options = append(options, filter.FilterOption{
			Label: display,
			Value: id,
		})
	}

	return options, nil
}

// getID extracts the ID from a related model instance
func (f *RelatedModelFilter[T, R]) getID(instance *R) int64 {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Try ID field
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		idField = val.FieldByName("Id")
	}
	if !idField.IsValid() {
		idField = val.FieldByName("id")
	}

	if idField.IsValid() {
		switch idField.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return idField.Int()
		}
	}

	return 0
}

// getDisplayValue extracts the display value from a related model instance
func (f *RelatedModelFilter[T, R]) getDisplayValue(instance *R) string {
	val := reflect.ValueOf(instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Try display field
	displayField := val.FieldByName(f.displayField)
	if !displayField.IsValid() {
		// Try capitalized version
		capitalized := f.capitalize(f.displayField)
		displayField = val.FieldByName(capitalized)
	}

	if displayField.IsValid() {
		switch displayField.Kind() {
		case reflect.String:
			return displayField.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", displayField.Int())
		}
	}

	// Fallback to ID
	id := f.getID(instance)
	return fmt.Sprintf("ID: %d", id)
}

// capitalize capitalizes the first letter of a string
func (f *RelatedModelFilter[T, R]) capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

// GetWidget returns the widget for this filter
func (f *RelatedModelFilter[T, R]) GetWidget() filter.Widget {
	return &RelatedModelWidget{
		fieldName:    f.GetFieldPath(),
		displayField: f.displayField,
	}
}

// RelatedModelWidget is a widget for related model filters
type RelatedModelWidget struct {
	fieldName    string
	displayField string
}

// Type returns the widget type
func (w *RelatedModelWidget) Type() string {
	return "related_model"
}

// Render renders the widget HTML as a select dropdown
func (w *RelatedModelWidget) Render(name string, value interface{}, attrs map[string]string) (string, error) {
	valueStr := ""
	if value != nil {
		valueStr = fmt.Sprintf("%v", value)
	}

	// This will be populated by JavaScript with options from GetOptions
	html := fmt.Sprintf(`
		<select name="%s" class="form-control related-model-select" data-field="%s" data-display-field="%s">
			<option value="">All</option>
			<option value="%s" selected>Loading...</option>
		</select>`,
		name, w.fieldName, w.displayField, valueStr,
	)

	return html, nil
}

// Parse parses the widget value
func (w *RelatedModelWidget) Parse(value string) (interface{}, error) {
	if value == "" {
		return nil, nil
	}

	var id int64
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		return id, nil
	}

	return value, nil
}
