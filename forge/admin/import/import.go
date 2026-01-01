package importpkg

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/forgego/forge/admin"
)

// ImportResult represents the result of an import operation
type ImportResult struct {
	TotalRows    int
	SuccessCount int
	ErrorCount   int
	Errors       []ImportError
	Created      []interface{}
	Updated      []interface{}
}

// ImportError represents an error during import
type ImportError struct {
	Row    int
	Column string
	Value  string
	Error  string
}

// Importer handles bulk import operations
type Importer[T any] struct {
	admin *admin.Admin[T]
}

// NewImporter creates a new importer
func NewImporter[T any](admin *admin.Admin[T]) *Importer[T] {
	return &Importer[T]{
		admin: admin,
	}
}

// ImportCSV imports data from a CSV reader
func (imp *Importer[T]) ImportCSV(ctx context.Context, reader io.Reader, options ImportOptions) (*ImportResult, error) {
	csvReader := csv.NewReader(reader)
	
	// Read header row
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Map columns to fields
	fieldMapping := imp.mapColumnsToFields(headers, options.ColumnMapping)

	result := &ImportResult{
		Errors: make([]ImportError, 0),
		Created: make([]interface{}, 0),
		Updated: make([]interface{}, 0),
	}

	rowNum := 1 // Start at 1 (header is row 0)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		result.TotalRows++
		rowNum++

		// Create instance from row
		instance, err := imp.createInstanceFromRow(row, fieldMapping, options)
		if err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, ImportError{
				Row:   rowNum,
				Error: err.Error(),
			})
			continue
		}

		// Validate instance
		if options.Validate {
			if err := imp.validateInstance(ctx, instance); err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, ImportError{
					Row:   rowNum,
					Error: err.Error(),
				})
				continue
			}
		}

		// Save instance
		if options.DryRun {
			// Don't actually save in dry run mode
			result.SuccessCount++
		} else {
			if err := imp.saveInstance(ctx, instance, options); err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, ImportError{
					Row:   rowNum,
					Error: err.Error(),
				})
				continue
			}
			result.SuccessCount++
			result.Created = append(result.Created, instance)
		}
	}

	return result, nil
}

// ImportOptions contains options for import
type ImportOptions struct {
	ColumnMapping map[string]string // CSV column -> model field
	Validate      bool              // Validate before import
	DryRun        bool              // Don't actually save
	UpdateExisting bool             // Update existing records instead of creating
	MatchFields   []string          // Fields to match for updates
	SkipHeader    bool              // Skip first row (header)
}

// DefaultImportOptions returns default import options
func DefaultImportOptions() ImportOptions {
	return ImportOptions{
		Validate:      true,
		DryRun:        false,
		UpdateExisting: false,
		SkipHeader:    true,
	}
}

// mapColumnsToFields maps CSV columns to model fields
func (imp *Importer[T]) mapColumnsToFields(headers []string, mapping map[string]string) map[int]string {
	result := make(map[int]string)
	
	for i, header := range headers {
		header = strings.TrimSpace(header)
		
		// Check explicit mapping first
		if field, ok := mapping[header]; ok {
			result[i] = field
			continue
		}

		// Try to match header to field name
		// Convert header to field name format (e.g., "First Name" -> "FirstName")
		fieldName := imp.headerToFieldName(header)
		
		// Check if field exists in model
		if imp.fieldExists(fieldName) {
			result[i] = fieldName
		}
	}

	return result
}

// headerToFieldName converts a CSV header to a field name
func (imp *Importer[T]) headerToFieldName(header string) string {
	// Remove spaces and capitalize words
	parts := strings.Fields(header)
	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return result
}

// fieldExists checks if a field exists in the model
func (imp *Importer[T]) fieldExists(fieldName string) bool {
	// Use reflection to check if field exists
	var zero T
	val := reflect.ValueOf(&zero).Elem()
	return val.FieldByName(fieldName).IsValid()
}

// createInstanceFromRow creates an instance from a CSV row
func (imp *Importer[T]) createInstanceFromRow(row []string, fieldMapping map[int]string, options ImportOptions) (*T, error) {
	var instance T
	val := reflect.ValueOf(&instance).Elem()

	for colIndex, fieldName := range fieldMapping {
		if colIndex >= len(row) {
			continue
		}

		value := strings.TrimSpace(row[colIndex])
		if value == "" {
			continue
		}

		field := val.FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		// Convert value to field type
		if err := imp.setFieldValue(field, value); err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldName, err)
		}
	}

	return &instance, nil
}

// setFieldValue sets a field value from a string
func (imp *Importer[T]) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer: %w", err)
		}
		field.SetInt(intVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float: %w", err)
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean: %w", err)
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

// validateInstance validates an instance
func (imp *Importer[T]) validateInstance(ctx context.Context, instance *T) error {
	// Use admin's validation system
	return admin.ValidateInstance(instance)
}

// saveInstance saves an instance
func (imp *Importer[T]) saveInstance(ctx context.Context, instance *T, options ImportOptions) error {
	manager := imp.admin.Manager()
	if manager == nil {
		return fmt.Errorf("admin has no manager")
	}

	if options.UpdateExisting && len(options.MatchFields) > 0 {
		// Try to find existing instance by matching fields
		existing, err := imp.findExistingInstance(ctx, manager, instance, options.MatchFields)
		if err == nil && existing != nil {
			// Update existing instance
			// Copy values from new instance to existing
			if err := imp.copyInstanceValues(existing, instance); err != nil {
				return fmt.Errorf("failed to copy values: %w", err)
			}
			return manager.Update(ctx, existing)
		}
		// If not found or error, fall through to create
	}

	return manager.Create(ctx, instance)
}

// findExistingInstance finds an existing instance by matching fields
func (imp *Importer[T]) findExistingInstance(ctx context.Context, manager interface{ All(context.Context) ([]*T, error) }, instance *T, matchFields []string) (*T, error) {
	// Get all instances and filter in memory
	// In production, you'd want to build a proper ORM query with Filter()
	all, err := manager.All(ctx)
	if err != nil {
		return nil, err
	}

	// Filter in memory by matching fields
	for _, item := range all {
		if imp.matchesInstance(item, instance, matchFields) {
			return item, nil
		}
	}

	return nil, fmt.Errorf("instance not found")
}

// matchesInstance checks if an instance matches another by specified fields
func (imp *Importer[T]) matchesInstance(existing, new *T, matchFields []string) bool {
	existingVal := reflect.ValueOf(existing).Elem()
	newVal := reflect.ValueOf(new).Elem()

	for _, fieldName := range matchFields {
		existingField := existingVal.FieldByName(fieldName)
		newField := newVal.FieldByName(fieldName)

		if !existingField.IsValid() || !newField.IsValid() {
			return false
		}

		if !reflect.DeepEqual(existingField.Interface(), newField.Interface()) {
			return false
		}
	}

	return true
}

// copyInstanceValues copies values from source to destination instance
func (imp *Importer[T]) copyInstanceValues(dest, src *T) error {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	// Get all fields from schema or reflection
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		destField := destVal.Field(i)

		if !destField.CanSet() {
			continue
		}

		// Skip ID and auto fields
		fieldName := srcVal.Type().Field(i).Name
		if fieldName == "ID" || fieldName == "Id" || fieldName == "id" {
			continue
		}

		// Copy value if source is not zero
		if !reflect.DeepEqual(srcField.Interface(), reflect.Zero(srcField.Type()).Interface()) {
			destField.Set(srcField)
		}
	}

	return nil
}
