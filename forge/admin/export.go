package admin

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	query "github.com/forgego/forge/orm"
)

// ExportFormat represents an export format
type ExportFormat string

const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
)

// ExportView handles data export
type ExportView[T any] struct {
	admin    *Admin[T]
	format   ExportFormat
	queryset query.QuerySet[T]
}

// NewExportView creates a new export view
func NewExportView[T any](admin *Admin[T], format ExportFormat) *ExportView[T] {
	return &ExportView[T]{
		admin:  admin,
		format: format,
	}
}

// Export exports data in the specified format
func (v *ExportView[T]) Export(ctx context.Context, w http.ResponseWriter) error {
	// Get queryset
	qs, err := v.admin.GetQueryset(ctx)
	if err != nil {
		return fmt.Errorf("failed to get queryset: %w", err)
	}

	// Get all objects (no pagination for export)
	objects, err := qs.All(ctx)
	if err != nil {
		return fmt.Errorf("failed to get objects: %w", err)
	}

	// Get fields to export
	fields := v.getExportFields()

	// Export based on format
	switch v.format {
	case ExportFormatCSV:
		return v.exportCSV(w, objects, fields)
	case ExportFormatJSON:
		return v.exportJSON(w, objects, fields)
	default:
		return fmt.Errorf("unsupported export format: %s", v.format)
	}
}

// getExportFields returns fields to export
func (v *ExportView[T]) getExportFields() []FieldExpr[T, interface{}] {
	if len(v.admin.config.ListDisplay) > 0 {
		return v.admin.config.ListDisplay
	}
	// Return empty - would need schema introspection
	return []FieldExpr[T, interface{}]{}
}

// exportCSV exports data as CSV
func (v *ExportView[T]) exportCSV(w http.ResponseWriter, objects []*T, fields []FieldExpr[T, interface{}]) error {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", v.admin.name))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	headers := make([]string, len(fields))
	for i, field := range fields {
		headers[i] = field.Name()
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, obj := range objects {
		row := make([]string, len(fields))
		for i, field := range fields {
			value := field.Get(obj)
			row[i] = fmt.Sprintf("%v", value)
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// exportJSON exports data as JSON
func (v *ExportView[T]) exportJSON(w http.ResponseWriter, objects []*T, fields []FieldExpr[T, interface{}]) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", v.admin.name))

	// Convert objects to maps
	results := make([]map[string]interface{}, len(objects))
	for i, obj := range objects {
		result := make(map[string]interface{})
		for _, field := range fields {
			value := field.Get(obj)
			result[field.Name()] = value
		}
		results[i] = result
	}

	return json.NewEncoder(w).Encode(results)
}

// GetFieldValue gets a field value from an object using reflection
func GetFieldValue(obj interface{}, fieldName string) interface{} {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	fieldValue := objValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil
	}

	if fieldValue.CanInterface() {
		return fieldValue.Interface()
	}

	return nil
}
