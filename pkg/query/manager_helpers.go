package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/forgego/forge/pkg/db"
	"github.com/forgego/forge/pkg/errors"
	"github.com/forgego/forge/pkg/schema"
)

// ExecuteHooks executes model hooks in the correct order
func ExecuteHooks(ctx context.Context, instance interface{}, hookType string) error {
	// Check if instance implements Schema interface
	schemaInstance, ok := instance.(schema.Schema)
	if !ok {
		return nil // No hooks if not a schema
	}

	hooks := schemaInstance.Hooks()
	if hooks == nil {
		return nil
	}

	switch hookType {
	case "BeforeCreate":
		if hooks.BeforeSave != nil {
			if err := hooks.BeforeSave(ctx, instance); err != nil {
				return fmt.Errorf("BeforeSave hook failed: %w", err)
			}
		}
		if hooks.BeforeCreate != nil {
			if err := hooks.BeforeCreate(ctx, instance); err != nil {
				return fmt.Errorf("BeforeCreate hook failed: %w", err)
			}
		}
	case "AfterCreate":
		if hooks.AfterCreate != nil {
			if err := hooks.AfterCreate(ctx, instance); err != nil {
				return fmt.Errorf("AfterCreate hook failed: %w", err)
			}
		}
		if hooks.AfterSave != nil {
			if err := hooks.AfterSave(ctx, instance); err != nil {
				return fmt.Errorf("AfterSave hook failed: %w", err)
			}
		}
	case "BeforeUpdate":
		if hooks.BeforeSave != nil {
			if err := hooks.BeforeSave(ctx, instance); err != nil {
				return fmt.Errorf("BeforeSave hook failed: %w", err)
			}
		}
		if hooks.BeforeUpdate != nil {
			if err := hooks.BeforeUpdate(ctx, instance); err != nil {
				return fmt.Errorf("BeforeUpdate hook failed: %w", err)
			}
		}
	case "AfterUpdate":
		if hooks.AfterUpdate != nil {
			if err := hooks.AfterUpdate(ctx, instance); err != nil {
				return fmt.Errorf("AfterUpdate hook failed: %w", err)
			}
		}
		if hooks.AfterSave != nil {
			if err := hooks.AfterSave(ctx, instance); err != nil {
				return fmt.Errorf("AfterSave hook failed: %w", err)
			}
		}
	case "BeforeDelete":
		if hooks.BeforeDelete != nil {
			if err := hooks.BeforeDelete(ctx, instance); err != nil {
				return fmt.Errorf("BeforeDelete hook failed: %w", err)
			}
		}
	case "AfterDelete":
		if hooks.AfterDelete != nil {
			if err := hooks.AfterDelete(ctx, instance); err != nil {
				return fmt.Errorf("AfterDelete hook failed: %w", err)
			}
		}
	}

	return nil
}

// ValidateInstance runs Clean validation hook if available
func ValidateInstance(instance interface{}) error {
	schemaInstance, ok := instance.(schema.Schema)
	if !ok {
		return nil
	}

	hooks := schemaInstance.Hooks()
	if hooks == nil || hooks.Clean == nil {
		return nil
	}

	if err := hooks.Clean(instance); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// BuildInsertSQL builds an INSERT SQL statement from a model instance
func BuildInsertSQL(instance interface{}, tableName string) (sql string, values []interface{}, columns []string, err error) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	if instanceValue.Kind() != reflect.Struct {
		return "", nil, nil, fmt.Errorf("instance must be a struct")
	}

	// Get schema to understand fields
	schemaInstance, ok := instance.(schema.Schema)
	if !ok {
		return "", nil, nil, fmt.Errorf("instance must implement schema.Schema")
	}

	// Iterate schema fields first, then get values only for those fields
	schemaFields := schemaInstance.Fields()

	var insertColumns []string
	var insertPlaceholders []string
	var insertValues []interface{}
	columnIndex := 1

	for _, schemaField := range schemaFields {
		// Skip primary key if auto-increment
		if schemaField.PrimaryKey && schemaField.AutoIncrement {
			continue
		}

		// Get field value using helper function
		fieldValue, err := getFieldValueByName(instance, schemaField.Name)
		if err != nil {
			// Field not found or not accessible - skip it
			continue
		}

		// Check if value is zero and field is not required
		fieldValueReflect := reflect.ValueOf(fieldValue)
		if !schemaField.Required && fieldValueReflect.IsZero() {
			continue
		}

		// Get column name from schema (DBColumn or Name)
		columnName := schemaField.DBColumn
		if columnName == "" {
			columnName = schemaField.Name
		}

		insertColumns = append(insertColumns, columnName)
		insertPlaceholders = append(insertPlaceholders, fmt.Sprintf("$%d", columnIndex))
		insertValues = append(insertValues, fieldValue)
		columnIndex++
	}

	if len(insertColumns) == 0 {
		return "", nil, nil, fmt.Errorf("no fields to insert")
	}

	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(insertColumns, ", "),
		strings.Join(insertPlaceholders, ", "),
	)

	return insertSQL, insertValues, insertColumns, nil
}

// getColumnName extracts column name from struct field (db tag or field name)
func getColumnName(field reflect.StructField) string {
	if dbTag := field.Tag.Get("db"); dbTag != "" {
		tagParts := strings.Split(dbTag, ",")
		if tagParts[0] != "" && tagParts[0] != "-" {
			return tagParts[0]
		}
	}
	return field.Name
}

// getFieldValueByName gets a struct field value by name using minimal reflection
// This is used when we know the field name from schema metadata
func getFieldValueByName(instance interface{}, fieldName string) (interface{}, error) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	if instanceValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("instance must be a struct")
	}

	fieldValue := instanceValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	if !fieldValue.CanInterface() {
		return nil, fmt.Errorf("field %s is not accessible", fieldName)
	}

	return fieldValue.Interface(), nil
}

// BuildUpdateSQL builds an UPDATE SQL statement from a model instance
func BuildUpdateSQL(instance interface{}, tableName, idField string) (string, []interface{}, error) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	if instanceValue.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("instance must be a struct")
	}

	// Get schema to understand fields
	schemaInstance, ok := instance.(schema.Schema)
	if !ok {
		return "", nil, fmt.Errorf("instance must implement schema.Schema")
	}

	// Iterate schema fields first, then get values only for those fields
	schemaFields := schemaInstance.Fields()

	var setParts []string
	var values []interface{}
	paramIndex := 1
	var idValue interface{}

	for _, schemaField := range schemaFields {
		// Get column name from schema (DBColumn or Name)
		columnName := schemaField.DBColumn
		if columnName == "" {
			columnName = schemaField.Name
		}

		// Check if this is the ID field
		if strings.EqualFold(schemaField.Name, idField) || strings.EqualFold(columnName, idField) {
			// Get ID value
			val, err := getFieldValueByName(instance, schemaField.Name)
			if err == nil {
				idValue = val
			}
			continue
		}

		// Skip primary key
		if schemaField.PrimaryKey {
			continue
		}

		// Get field value using helper function
		fieldValue, err := getFieldValueByName(instance, schemaField.Name)
		if err != nil {
			// Field not found or not accessible - skip it
			continue
		}

		// Include field in UPDATE
		setParts = append(setParts, fmt.Sprintf("%s = $%d", columnName, paramIndex))
		values = append(values, fieldValue)
		paramIndex++
	}

	if len(setParts) == 0 {
		return "", nil, fmt.Errorf("no fields to update")
	}

	if idValue == nil {
		return "", nil, fmt.Errorf("id field not found or is zero")
	}

	// Add ID to the end for WHERE clause
	values = append(values, idValue)
	sql := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d",
		tableName,
		strings.Join(setParts, ", "),
		idField,
		paramIndex,
	)

	return sql, values, nil
}

// BuildDeleteSQL builds a DELETE SQL statement
func BuildDeleteSQL(tableName, idField string, idValue interface{}) (string, []interface{}) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, idField)
	return sql, []interface{}{idValue}
}

// ExecuteInsert executes an INSERT statement and returns the generated ID
func ExecuteInsert(ctx context.Context, database *db.DB, sql string, args []interface{}) (int64, error) {
	sqldb, err := getSQLDB(database)
	if err != nil {
		return 0, err
	}

	var id int64
	err = sqldb.QueryRowContext(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert failed: %w", err)
	}

	return id, nil
}

// ExecuteUpdate executes an UPDATE statement and returns rows affected
func ExecuteUpdate(ctx context.Context, database *db.DB, sql string, args []interface{}) (int64, error) {
	sqldb, err := getSQLDB(database)
	if err != nil {
		return 0, err
	}

	result, err := sqldb.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("update failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// ExecuteDelete executes a DELETE statement and returns rows affected
func ExecuteDelete(ctx context.Context, database *db.DB, sql string, args []interface{}) (int64, error) {
	sqldb, err := getSQLDB(database)
	if err != nil {
		return 0, err
	}

	result, err := sqldb.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// getSQLDB extracts *sql.DB from *db.DB
func getSQLDB(database *db.DB) (*sql.DB, error) {
	if database == nil {
		return nil, errors.NewNotImplementedError("database connection not set")
	}
	return database.DB, nil
}

// GetIDValue extracts the ID value from an instance
func GetIDValue(instance interface{}, idFieldName string) (interface{}, error) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	if instanceValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("instance must be a struct")
	}

	typ := instanceValue.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := instanceValue.Field(i)

		// Check field name
		if strings.EqualFold(field.Name, idFieldName) {
			if !fieldValue.CanInterface() {
				return nil, fmt.Errorf("id field is not accessible")
			}
			return fieldValue.Interface(), nil
		}

		// Check db tag
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			tagParts := strings.Split(dbTag, ",")
			if strings.EqualFold(tagParts[0], idFieldName) {
				if !fieldValue.CanInterface() {
					return nil, fmt.Errorf("id field is not accessible")
				}
				return fieldValue.Interface(), nil
			}
		}
	}

	return nil, fmt.Errorf("id field '%s' not found", idFieldName)
}

