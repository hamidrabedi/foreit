package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/forgego/forge/pkg/errors"
	"github.com/forgego/forge/pkg/schema"
)

// fieldMappingCache caches field mappings per type to reduce reflection
// Cache maps reflect.Type -> (column name -> field name)
var (
	fieldMappingCache = make(map[reflect.Type]map[string]string)
	fieldMappingMu    sync.RWMutex
)

// QuerySet is the base interface for type-safe query sets
// Generated QuerySets will implement this interface
type QuerySet[T any] interface {
	// Filtering
	Filter(expr QueryExpr) QuerySet[T]
	Exclude(expr QueryExpr) QuerySet[T]

	// Ordering
	OrderBy(fields ...string) QuerySet[T]
	Reverse() QuerySet[T]

	// Limiting
	Limit(n int) QuerySet[T]
	Offset(n int) QuerySet[T]
	Distinct() QuerySet[T]

	// Selection
	Select(fields ...string) QuerySet[T]
	SelectRelated(relations ...string) QuerySet[T]
	PrefetchRelated(relations ...string) QuerySet[T]
	Only(fields ...string) QuerySet[T]
	Defer(fields ...string) QuerySet[T]

	// Aggregation
	Aggregate(aggregates ...Aggregate) QuerySet[T]
	Annotate(annotations ...AnnotationExpr) QuerySet[T]

	// Values
	Values(fields ...string) QuerySet[T]
	ValuesList(fields ...string) QuerySet[T]

	// Queries
	All(ctx context.Context) ([]*T, error)
	Get(ctx context.Context) (*T, error)
	First(ctx context.Context) (*T, error)
	Last(ctx context.Context) (*T, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context) (bool, error)

	// Updates
	Update(ctx context.Context, fields map[string]interface{}) (int64, error)
	BulkUpdate(ctx context.Context, updates []map[string]interface{}) error
	BulkCreate(ctx context.Context, instances []*T) error

	// Deletes
	Delete(ctx context.Context) (int64, error)

	// Unions
	Union(other QuerySet[T]) QuerySet[T]
	Intersection(other QuerySet[T]) QuerySet[T]
	Difference(other QuerySet[T]) QuerySet[T]
}

// BaseQuerySet provides a base implementation for generated QuerySets
// This is exported so generated code can embed it
type BaseQuerySet[T any] struct {
	table           string
	conditions      []QueryExpr
	excludes        []QueryExpr
	orderBy         []string
	limitVal        *int
	offsetVal       *int
	distinct        bool
	selectFields    []string
	selectRelated   []string
	prefetchRelated []string
	onlyFields      []string
	deferFields     []string
	aggregates      []Aggregate
	annotations     []AnnotationExpr
	// Database connection - can be set explicitly or retrieved from context
	db interface{} // *db.DB - using interface{} to avoid circular import
}

// Aggregate represents an aggregate function
type Aggregate struct {
	Name  string
	Field string
	Func  string // COUNT, SUM, AVG, MAX, MIN, etc.
}

// Annotation represents a computed field annotation
type Annotation struct {
	Name string
	Expr QueryExpr
}

// NewBaseQuerySet creates a new base query set
func NewBaseQuerySet[T any](table string) *BaseQuerySet[T] {
	return &BaseQuerySet[T]{
		table:           table,
		conditions:      []QueryExpr{},
		excludes:        []QueryExpr{},
		orderBy:         []string{},
		selectRelated:   []string{},
		prefetchRelated: []string{},
		aggregates:      []Aggregate{},
		annotations:     []AnnotationExpr{},
	}
}

// SetDB sets the database connection for this QuerySet
func (b *BaseQuerySet[T]) SetDB(db interface{}) {
	b.db = db
}

// getDB retrieves the database connection from QuerySet or context
func (b *BaseQuerySet[T]) getDB(ctx context.Context) (*sql.DB, error) {
	if b.db != nil {
		return GetSQLDB(b.db)
	}
	return nil, errors.NewNotImplementedError("database connection not set on QuerySet")
}

func (b *BaseQuerySet[T]) Filter(expr QueryExpr) QuerySet[T] {
	b.conditions = append(b.conditions, expr)
	return b
}

func (b *BaseQuerySet[T]) Exclude(expr QueryExpr) QuerySet[T] {
	b.excludes = append(b.excludes, expr)
	return b
}

func (b *BaseQuerySet[T]) OrderBy(fields ...string) QuerySet[T] {
	b.orderBy = append(b.orderBy, fields...)
	return b
}

func (b *BaseQuerySet[T]) Reverse() QuerySet[T] {
	// Reverse the order by fields
	for i, field := range b.orderBy {
		if strings.HasPrefix(field, "-") {
			b.orderBy[i] = strings.TrimPrefix(field, "-")
		} else {
			b.orderBy[i] = "-" + field
		}
	}
	return b
}

func (b *BaseQuerySet[T]) Limit(n int) QuerySet[T] {
	b.limitVal = &n
	return b
}

func (b *BaseQuerySet[T]) Offset(n int) QuerySet[T] {
	b.offsetVal = &n
	return b
}

func (b *BaseQuerySet[T]) Distinct() QuerySet[T] {
	b.distinct = true
	return b
}

func (b *BaseQuerySet[T]) Select(fields ...string) QuerySet[T] {
	b.selectFields = fields
	return b
}

func (b *BaseQuerySet[T]) SelectRelated(relations ...string) QuerySet[T] {
	b.selectRelated = append(b.selectRelated, relations...)
	return b
}

func (b *BaseQuerySet[T]) PrefetchRelated(relations ...string) QuerySet[T] {
	b.prefetchRelated = append(b.prefetchRelated, relations...)
	return b
}

func (b *BaseQuerySet[T]) Only(fields ...string) QuerySet[T] {
	b.onlyFields = fields
	return b
}

func (b *BaseQuerySet[T]) Defer(fields ...string) QuerySet[T] {
	b.deferFields = fields
	return b
}

func (b *BaseQuerySet[T]) Aggregate(aggregates ...Aggregate) QuerySet[T] {
	b.aggregates = append(b.aggregates, aggregates...)
	return b
}

// ExecuteAggregates executes aggregate queries and returns results as a map
func (b *BaseQuerySet[T]) ExecuteAggregates(ctx context.Context) (map[string]interface{}, error) {
	if len(b.aggregates) == 0 {
		return nil, fmt.Errorf("no aggregates specified")
	}

	db, err := b.getDB(ctx)
	if err != nil {
		return nil, err
	}

	// Build SELECT clause with aggregates using SQL builder
	builder := NewSQLBuilder()
	var selectParts []string
	for _, agg := range b.aggregates {
		fieldExpr := agg.Field
		if fieldExpr == "*" || fieldExpr == "" {
			fieldExpr = "*"
		} else {
			// Escape field identifier
			fieldExpr = EscapeIdentifier(fieldExpr)
		}
		// Escape aggregate name alias
		alias := EscapeIdentifier(agg.Name)
		selectParts = append(selectParts, fmt.Sprintf("%s(%s) AS %s", agg.Func, fieldExpr, alias))
	}

	// Build FROM clause with escaped table name
	escapedTable := EscapeIdentifier(b.table)
	fromClause := "FROM " + escapedTable

	// Build WHERE clause using SQL builder
	whereClause, whereArgs := builder.BuildWhere(b.conditions, b.excludes)

	sqlQuery := "SELECT " + strings.Join(selectParts, ", ") + " " + fromClause
	if whereClause != "" {
		sqlQuery += " " + whereClause
	}
	
	args := whereArgs

	// Execute query
	row := db.QueryRowContext(ctx, sqlQuery, args...)

	// Scan results into map
	result := make(map[string]interface{})
	scanArgs := make([]interface{}, len(b.aggregates))
	for i := range scanArgs {
		var val interface{}
		scanArgs[i] = &val
	}

	if err := row.Scan(scanArgs...); err != nil {
		return nil, fmt.Errorf("aggregate query failed: %w", err)
	}

	// Map results to aggregate names
	for i, agg := range b.aggregates {
		if i < len(scanArgs) {
			valPtr := scanArgs[i].(*interface{})
			result[agg.Name] = *valPtr
		}
	}

	return result, nil
}

func (b *BaseQuerySet[T]) Annotate(annotations ...AnnotationExpr) QuerySet[T] {
	b.annotations = append(b.annotations, annotations...)
	return b
}

func (b *BaseQuerySet[T]) Values(fields ...string) QuerySet[T] {
	// Values returns dictionaries instead of model instances
	// This is a simplified implementation - full version would return []map[string]interface{}
	b.selectFields = fields
	return b
}

func (b *BaseQuerySet[T]) ValuesList(fields ...string) QuerySet[T] {
	// ValuesList returns tuples instead of model instances
	// This is a simplified implementation - full version would return [][]interface{}
	b.selectFields = fields
	return b
}

func (b *BaseQuerySet[T]) All(ctx context.Context) ([]*T, error) {
	db, err := b.getDB(ctx)
	if err != nil {
		return nil, err
	}

	sqlQuery, args := b.buildSQL()
	
	rows, err := db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var results []*T
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Get type once for all rows
	var instance T
	typ := reflect.TypeOf(instance)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Build column-to-field mapping from schema
	columnToFieldMap, err := getColumnToFieldMapping(typ)
	if err != nil {
		// Fallback to reflection-based approach if schema not available
		columnToFieldMap = nil
	}

	for rows.Next() {
		// Create a new instance of T
		instanceValue := reflect.New(typ).Elem()
		
		// Prepare scan destination using schema-based or reflection-based mapping
		scanArgs := make([]interface{}, len(columns))
		for i, colName := range columns {
			if columnToFieldMap != nil {
				// Use schema-based mapping
				if fieldName, ok := columnToFieldMap[strings.ToLower(colName)]; ok {
					field := instanceValue.FieldByName(fieldName)
					if field.IsValid() && field.CanSet() {
						scanArgs[i] = field.Addr().Interface()
					} else {
						var val interface{}
						scanArgs[i] = &val
					}
				} else {
					// Use a generic interface{} for unmapped columns
					var val interface{}
					scanArgs[i] = &val
				}
			} else {
				// Fallback to reflection-based approach
				if fieldIdx, ok := getFieldIndex(typ, colName); ok {
					field := instanceValue.Field(fieldIdx)
					if field.CanSet() {
						scanArgs[i] = field.Addr().Interface()
					} else {
						var val interface{}
						scanArgs[i] = &val
					}
				} else {
					// Use a generic interface{} for unmapped columns
					var val interface{}
					scanArgs[i] = &val
				}
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, instanceValue.Addr().Interface().(*T))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// getColumnToFieldMapping builds a column-to-field mapping from schema metadata
// Returns map[columnName]fieldName (both lowercase for case-insensitive matching)
// Uses caching to reduce overhead
func getColumnToFieldMapping(typ reflect.Type) (map[string]string, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type must be a struct")
	}

	// Check cache first
	fieldMappingMu.RLock()
	columnMap, exists := fieldMappingCache[typ]
	fieldMappingMu.RUnlock()

	if exists {
		return columnMap, nil
	}

	// Build cache entry from schema
	fieldMappingMu.Lock()
	defer fieldMappingMu.Unlock()

	// Double-check after acquiring write lock
	if columnMap, exists := fieldMappingCache[typ]; exists {
		return columnMap, nil
	}

	// Try to get schema from a zero value instance
	instanceValue := reflect.New(typ).Elem()
	schemaInstance, ok := instanceValue.Interface().(schema.Schema)
	if !ok {
		// Not a schema model - return error to use fallback
		return nil, fmt.Errorf("type does not implement schema.Schema")
	}

	// Build mapping from schema fields
	fields := schemaInstance.Fields()
	columnMap = make(map[string]string, len(fields))
	for _, field := range fields {
		// Get column name from schema (DBColumn or Name)
		colName := field.DBColumn
		if colName == "" {
			colName = field.Name
		}
		// Store both column name and field name (lowercase for case-insensitive matching)
		columnMap[strings.ToLower(colName)] = field.Name
		// Also map field name to itself (in case column name matches field name)
		columnMap[strings.ToLower(field.Name)] = field.Name
	}

	fieldMappingCache[typ] = columnMap
	return columnMap, nil
}

// getFieldIndex finds a struct field index by database column name (case-insensitive)
// Uses caching to reduce reflection overhead (fallback when schema not available)
func getFieldIndex(typ reflect.Type, columnName string) (int, bool) {
	if typ.Kind() != reflect.Struct {
		return -1, false
	}

	// Try schema-based mapping first
	columnMap, err := getColumnToFieldMapping(typ)
	if err == nil && columnMap != nil {
		// Schema-based mapping available - find field index by name
		if fieldName, ok := columnMap[strings.ToLower(columnName)]; ok {
			// Find field index by name
			for i := 0; i < typ.NumField(); i++ {
				if typ.Field(i).Name == fieldName {
					return i, true
				}
			}
		}
	}

	// Fallback: build reflection-based mapping
	// This is the old approach for non-schema models
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Check field name
		if strings.EqualFold(field.Name, columnName) {
			return i, true
		}
		// Check db tag
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			tagParts := strings.Split(dbTag, ",")
			if tagParts[0] != "" && tagParts[0] != "-" {
				if strings.EqualFold(tagParts[0], columnName) {
					return i, true
				}
			}
		}
	}

	return -1, false
}

func (b *BaseQuerySet[T]) Get(ctx context.Context) (*T, error) {
	results, err := b.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundErrorWithMessage(fmt.Sprintf("%s matching query does not exist", b.table))
	}

	if len(results) > 1 {
		return nil, errors.NewMultipleObjectsReturnedError(b.table, len(results))
	}

	return results[0], nil
}

func (b *BaseQuerySet[T]) First(ctx context.Context) (*T, error) {
	results, err := b.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundErrorWithMessage(fmt.Sprintf("%s matching query does not exist", b.table))
	}

	return results[0], nil
}

func (b *BaseQuerySet[T]) Last(ctx context.Context) (*T, error) {
	// Reverse order and get first
	reversed := b.Reverse()
	results, err := reversed.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundErrorWithMessage(fmt.Sprintf("%s matching query does not exist", b.table))
	}

	return results[0], nil
}

func (b *BaseQuerySet[T]) Count(ctx context.Context) (int64, error) {
	db, err := b.getDB(ctx)
	if err != nil {
		return 0, err
	}

	// Build COUNT query using SQL builder
	builder := NewSQLBuilder()
	selectClause := builder.BuildSelect(b.table, []string{}, false)
	// Replace SELECT * with SELECT COUNT(*)
	selectClause = strings.Replace(selectClause, "SELECT *", "SELECT COUNT(*)", 1)
	
	var parts []string
	parts = append(parts, selectClause)

	// Build WHERE clause
	whereClause, _ := builder.BuildWhere(b.conditions, b.excludes)
	if whereClause != "" {
		parts = append(parts, whereClause)
	}

	sqlQuery := strings.Join(parts, " ")
	args := builder.Args()

	var count int64
	err = db.QueryRowContext(ctx, sqlQuery, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count query failed: %w", err)
	}

	return count, nil
}

func (b *BaseQuerySet[T]) Exists(ctx context.Context) (bool, error) {
	count, err := b.Count(ctx)
	return count > 0, err
}

func (b *BaseQuerySet[T]) Update(ctx context.Context, fields map[string]interface{}) (int64, error) {
	db, err := b.getDB(ctx)
	if err != nil {
		return 0, err
	}

	if len(fields) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	// Build UPDATE query using SQL builder
	builder := NewSQLBuilder()
	updateClause, _ := builder.BuildUpdate(b.table, fields)
	
	// Build WHERE clause
	whereClause, _ := builder.BuildWhere(b.conditions, b.excludes)
	
	sqlQuery := updateClause
	if whereClause != "" {
		sqlQuery += " " + whereClause
	}
	
	// Get all args from builder (includes both update and where args)
	args := builder.Args()

	result, err := db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("update query failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (b *BaseQuerySet[T]) BulkUpdate(ctx context.Context, updates []map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	db, err := b.getDB(ctx)
	if err != nil {
		return err
	}

	// Get transaction for bulk operations
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Get all field names from first update
	var allFields []string
	fieldSet := make(map[string]bool)
	for _, update := range updates {
		for field := range update {
			if !fieldSet[field] {
				allFields = append(allFields, field)
				fieldSet[field] = true
			}
		}
	}

	// Build bulk update using CASE statements (PostgreSQL)
	// UPDATE table SET field1 = CASE id WHEN ? THEN ? WHEN ? THEN ? END, ...
	var setParts []string
	var args []interface{}
	var allIDs []interface{}
	paramIndex := 1

	// We need IDs to identify which row to update
	// For now, assume 'id' field exists - full implementation would get primary key from schema
	idField := "id"

	// Collect all IDs first
	for _, update := range updates {
		if id, ok := update[idField]; ok {
			allIDs = append(allIDs, id)
		}
	}

	// Build CASE statements for each field
	for _, field := range allFields {
		var whenParts []string
		var fieldArgs []interface{}

		for _, update := range updates {
			if id, ok := update[idField]; ok {
				if val, ok := update[field]; ok {
					whenParts = append(whenParts, fmt.Sprintf("WHEN $%d THEN $%d", paramIndex, paramIndex+1))
					fieldArgs = append(fieldArgs, id, val)
					paramIndex += 2
				}
			}
		}

		if len(whenParts) > 0 {
			setParts = append(setParts, fmt.Sprintf("%s = CASE %s %s END", field, idField, strings.Join(whenParts, " ")))
			args = append(args, fieldArgs...)
		}
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no valid updates")
	}

	// Build WHERE clause for IDs
	var idPlaceholders []string
	for _, id := range allIDs {
		idPlaceholders = append(idPlaceholders, fmt.Sprintf("$%d", paramIndex))
		args = append(args, id)
		paramIndex++
	}

	sqlQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s IN (%s)",
		b.table, strings.Join(setParts, ", "), idField, strings.Join(idPlaceholders, ", "))

	_, err = tx.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("bulk update failed: %w", err)
	}

	return tx.Commit()
}

func (b *BaseQuerySet[T]) BulkCreate(ctx context.Context, instances []*T) error {
	if len(instances) == 0 {
		return fmt.Errorf("no instances to create")
	}

	db, err := b.getDB(ctx)
	if err != nil {
		return err
	}

	// Get fields from first instance using reflection
	firstInstance := instances[0]
	instanceValue := reflect.ValueOf(firstInstance).Elem()
	typ := instanceValue.Type()

	var columns []string
	var columnTypes []reflect.Type
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		// Skip primary key if auto-increment
		if field.Name == "ID" {
			continue
		}
		columns = append(columns, field.Name)
		columnTypes = append(columnTypes, field.Type)
	}

	if len(columns) == 0 {
		return fmt.Errorf("no fields to insert")
	}

	// Build bulk INSERT
	var placeholders []string
	var allValues []interface{}
	paramIndex := 1

		for _, instance := range instances {
			var rowPlaceholders []string
			instanceValue := reflect.ValueOf(instance).Elem()
			for _, col := range columns {
				fieldValue := instanceValue.FieldByName(col)
				if fieldValue.IsValid() && fieldValue.CanInterface() {
					rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", paramIndex))
					allValues = append(allValues, fieldValue.Interface())
					paramIndex++
				}
			}
		placeholders = append(placeholders, "("+strings.Join(rowPlaceholders, ", ")+")")
	}

	sqlQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s RETURNING id",
		b.table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	rows, err := db.QueryContext(ctx, sqlQuery, allValues...)
	if err != nil {
		return fmt.Errorf("bulk create failed: %w", err)
	}
	defer rows.Close()

	// Update IDs on instances
	idx := 0
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan ID: %w", err)
		}
		if idx < len(instances) {
			instanceValue := reflect.ValueOf(instances[idx]).Elem()
			idField := instanceValue.FieldByName("ID")
			if idField.IsValid() && idField.CanSet() {
				idField.SetInt(id)
			}
		}
		idx++
	}

	return rows.Err()
}

func (b *BaseQuerySet[T]) Delete(ctx context.Context) (int64, error) {
	db, err := b.getDB(ctx)
	if err != nil {
		return 0, err
	}

	// Build DELETE query
	var args []interface{}
	paramIndex := 1

	// Build WHERE clause
	var whereParts []string
	if len(b.conditions) > 0 {
		for _, cond := range b.conditions {
			sql, condArgs, nextIndex := cond.ToSQL(paramIndex)
			whereParts = append(whereParts, sql)
			args = append(args, condArgs...)
			paramIndex = nextIndex
		}
	}
	if len(b.excludes) > 0 {
		for _, exclude := range b.excludes {
			sql, excludeArgs, nextIndex := exclude.ToSQL(paramIndex)
			whereParts = append(whereParts, "NOT ("+sql+")")
			args = append(args, excludeArgs...)
			paramIndex = nextIndex
		}
	}

	sqlQuery := fmt.Sprintf("DELETE FROM %s", b.table)
	if len(whereParts) > 0 {
		sqlQuery += " WHERE " + strings.Join(whereParts, " AND ")
	}

	result, err := db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("delete query failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (b *BaseQuerySet[T]) Union(other QuerySet[T]) QuerySet[T] {
	// TODO: Implement union - SQL UNION operation
	// This will require combining SQL queries from both QuerySets
	return b
}

func (b *BaseQuerySet[T]) Intersection(other QuerySet[T]) QuerySet[T] {
	// TODO: Implement intersection - SQL INTERSECT operation
	// This will require combining SQL queries from both QuerySets
	return b
}

func (b *BaseQuerySet[T]) Difference(other QuerySet[T]) QuerySet[T] {
	// TODO: Implement difference - SQL EXCEPT operation
	// This will require combining SQL queries from both QuerySets
	return b
}

// buildSQL builds the SQL query from the query set using the SQL builder
func (b *BaseQuerySet[T]) buildSQL() (string, []interface{}) {
	builder := NewSQLBuilder()
	var parts []string

	// Determine fields to select
	var selectFields []string
	if len(b.selectFields) > 0 {
		selectFields = b.selectFields
	} else if len(b.onlyFields) > 0 {
		selectFields = b.onlyFields
	}

	// Build SELECT clause
	selectClause := builder.BuildSelect(b.table, selectFields, b.distinct)
	parts = append(parts, selectClause)

	// JOIN for SelectRelated (basic implementation for ForeignKey relations)
	// For MVP, we'll add JOIN clauses based on relation names
	// Full implementation would require model registry to get relation details
	if len(b.selectRelated) > 0 {
		for _, relation := range b.selectRelated {
			// Basic JOIN: assume relation name matches table name + _id suffix
			// e.g., "author" -> JOIN users ON posts.author_id = users.id
			// This is a simplified implementation - full version would use model registry
			joinTable := relation + "s" // Simple pluralization
			joinClause := fmt.Sprintf("LEFT JOIN %s ON %s.%s_id = %s.id",
				EscapeIdentifier(joinTable),
				EscapeIdentifier(b.table),
				relation,
				EscapeIdentifier(joinTable))
			parts = append(parts, joinClause)
		}
	}

	// Build WHERE clause
	whereClause, _ := builder.BuildWhere(b.conditions, b.excludes)
	if whereClause != "" {
		parts = append(parts, whereClause)
	}

	// Build ORDER BY
	orderByClause := builder.BuildOrderBy(b.orderBy)
	if orderByClause != "" {
		parts = append(parts, orderByClause)
	}

	// Build LIMIT
	limitClause := builder.BuildLimit(b.limitVal)
	if limitClause != "" {
		parts = append(parts, limitClause)
	}

	// Build OFFSET
	offsetClause := builder.BuildOffset(b.offsetVal)
	if offsetClause != "" {
		parts = append(parts, offsetClause)
	}

	return strings.Join(parts, " "), builder.Args()
}
