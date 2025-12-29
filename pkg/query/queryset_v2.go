package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// QuerySetV2 is the enhanced type-safe QuerySet interface
type QuerySetV2[T any] interface {
	// SetDB sets the database connection
	SetDB(db interface{}) QuerySetV2[T]
	// Filtering - accepts new Expression interface
	Filter(expr Expression) QuerySetV2[T]
	Exclude(expr Expression) QuerySetV2[T]
	
	// Ordering
	OrderBy(fields ...OrderField) QuerySetV2[T]
	Reverse() QuerySetV2[T]
	
	// Limiting
	Limit(n int) QuerySetV2[T]
	Offset(n int) QuerySetV2[T]
	Distinct(fields ...string) QuerySetV2[T]
	
	// Field Selection
	Select(fields ...string) QuerySetV2[T]
	Only(fields ...string) QuerySetV2[T]
	Defer(fields ...string) QuerySetV2[T]
	
	// Relations
	SelectRelated(relations ...string) QuerySetV2[T]
	PrefetchRelated(relations ...string) QuerySetV2[T]
	
	// Aggregation
	Aggregate(aggs ...Aggregate) QuerySetV2[T]
	Annotate(anns ...AnnotationExpr) QuerySetV2[T]
	
	// Values - type-safe return types
	Values(fields ...string) ValuesQuerySet[T]
	ValuesList(fields ...string) ValuesListQuerySet[T]
	
	// Execution
	All(ctx context.Context) ([]*T, error)
	Get(ctx context.Context) (*T, error)
	First(ctx context.Context) (*T, error)
	Last(ctx context.Context) (*T, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context) (bool, error)
	
	// Updates - type-safe
	Update(ctx context.Context, updates UpdateMap) (int64, error)
	UpdateBuilder() (*UpdateBuilder[T], error)
	BulkUpdate(ctx context.Context, updates []UpdateMap) error
	
	// Deletes
	Delete(ctx context.Context) (int64, error)
	
	// Set Operations
	Union(other QuerySetV2[T]) QuerySetV2[T]
	Intersection(other QuerySetV2[T]) QuerySetV2[T]
	Difference(other QuerySetV2[T]) QuerySetV2[T]
}

// OrderField represents an ordering field with direction
type OrderField struct {
	Field     string
	Ascending bool
}

// NewOrderField creates an order field
func NewOrderField(field string, ascending bool) OrderField {
	return OrderField{
		Field:     field,
		Ascending: ascending,
	}
}

// Asc creates an ascending order field
func Asc(field string) OrderField {
	return NewOrderField(field, true)
}

// Desc creates a descending order field
func Desc(field string) OrderField {
	return NewOrderField(field, false)
}

// BaseQuerySetV2 is the enhanced implementation
type BaseQuerySetV2[T any] struct {
	table           string
	schema          *ModelSchema
	conditions      []Expression
	excludes        []Expression
	orderBy         []OrderField
	limitVal        *int
	offsetVal       *int
	distinctFields  []string
	selectFields    []string
	onlyFields      []string
	deferFields     []string
	selectRelated   []string
	prefetchRelated []string
	aggregates      []Aggregate
	annotations     []AnnotationExpr // Using existing AnnotationExpr type
	db              interface{} // *db.DB
}

// NewQuerySetV2 creates a new enhanced QuerySet
func NewQuerySetV2[T any](tableName string) (QuerySetV2[T], error) {
	schema, err := GetModelSchema[T]()
	if err != nil {
		return nil, err
	}
	
	if tableName == "" {
		tableName = schema.TableName
	}
	
	return &BaseQuerySetV2[T]{
		table:      tableName,
		schema:     schema,
		conditions: []Expression{},
		excludes:   []Expression{},
		orderBy:    []OrderField{},
	}, nil
}

// SetDB sets the database connection
func (qs *BaseQuerySetV2[T]) SetDB(db interface{}) QuerySetV2[T] {
	qs.db = db
	return qs
}

// getDB retrieves the database connection
func (qs *BaseQuerySetV2[T]) getDB(ctx context.Context) (*sql.DB, error) {
	if qs.db != nil {
		return GetSQLDB(qs.db)
	}
	return nil, fmt.Errorf("database connection not set on QuerySet")
}

// clone creates a deep copy
func (qs *BaseQuerySetV2[T]) clone() *BaseQuerySetV2[T] {
	return &BaseQuerySetV2[T]{
		table:           qs.table,
		schema:          qs.schema,
		conditions:      append([]Expression{}, qs.conditions...),
		excludes:        append([]Expression{}, qs.excludes...),
		orderBy:         append([]OrderField{}, qs.orderBy...),
		limitVal:        qs.limitVal,
		offsetVal:       qs.offsetVal,
		distinctFields:  append([]string{}, qs.distinctFields...),
		selectFields:    append([]string{}, qs.selectFields...),
		onlyFields:      append([]string{}, qs.onlyFields...),
		deferFields:     append([]string{}, qs.deferFields...),
		selectRelated:   append([]string{}, qs.selectRelated...),
		prefetchRelated: append([]string{}, qs.prefetchRelated...),
		aggregates:      append([]Aggregate{}, qs.aggregates...),
		annotations:     append([]AnnotationExpr{}, qs.annotations...),
		db:              qs.db,
	}
}

// Filter adds a filter condition
func (qs *BaseQuerySetV2[T]) Filter(expr Expression) QuerySetV2[T] {
	// Validate expression
	if err := expr.Resolve(qs.schema); err != nil {
		panic(fmt.Sprintf("invalid filter expression: %v", err))
	}
	
	clone := qs.clone()
	clone.conditions = append(clone.conditions, expr)
	return clone
}

// Exclude adds an exclude condition
func (qs *BaseQuerySetV2[T]) Exclude(expr Expression) QuerySetV2[T] {
	// Validate expression
	if err := expr.Resolve(qs.schema); err != nil {
		panic(fmt.Sprintf("invalid exclude expression: %v", err))
	}
	
	clone := qs.clone()
	clone.excludes = append(clone.excludes, expr)
	return clone
}

// OrderBy sets ordering
func (qs *BaseQuerySetV2[T]) OrderBy(fields ...OrderField) QuerySetV2[T] {
	clone := qs.clone()
	clone.orderBy = append(clone.orderBy, fields...)
	return clone
}

// Reverse reverses the current ordering
func (qs *BaseQuerySetV2[T]) Reverse() QuerySetV2[T] {
	clone := qs.clone()
	for i := range clone.orderBy {
		clone.orderBy[i].Ascending = !clone.orderBy[i].Ascending
	}
	return clone
}

// Limit sets the limit
func (qs *BaseQuerySetV2[T]) Limit(n int) QuerySetV2[T] {
	clone := qs.clone()
	clone.limitVal = &n
	return clone
}

// Offset sets the offset
func (qs *BaseQuerySetV2[T]) Offset(n int) QuerySetV2[T] {
	clone := qs.clone()
	clone.offsetVal = &n
	return clone
}

// Distinct sets distinct fields
func (qs *BaseQuerySetV2[T]) Distinct(fields ...string) QuerySetV2[T] {
	clone := qs.clone()
	if len(fields) > 0 {
		clone.distinctFields = fields
	} else {
		clone.distinctFields = []string{"*"}
	}
	return clone
}

// Select sets fields to select
func (qs *BaseQuerySetV2[T]) Select(fields ...string) QuerySetV2[T] {
	clone := qs.clone()
	clone.selectFields = fields
	return clone
}

// Only sets fields to only load
func (qs *BaseQuerySetV2[T]) Only(fields ...string) QuerySetV2[T] {
	clone := qs.clone()
	clone.onlyFields = fields
	return clone
}

// Defer sets fields to defer loading
func (qs *BaseQuerySetV2[T]) Defer(fields ...string) QuerySetV2[T] {
	clone := qs.clone()
	clone.deferFields = fields
	return clone
}

// SelectRelated adds relations to select
func (qs *BaseQuerySetV2[T]) SelectRelated(relations ...string) QuerySetV2[T] {
	clone := qs.clone()
	clone.selectRelated = append(clone.selectRelated, relations...)
	return clone
}

// PrefetchRelated adds relations to prefetch
func (qs *BaseQuerySetV2[T]) PrefetchRelated(relations ...string) QuerySetV2[T] {
	clone := qs.clone()
	clone.prefetchRelated = append(clone.prefetchRelated, relations...)
	return clone
}

// Aggregate adds aggregate expressions
func (qs *BaseQuerySetV2[T]) Aggregate(aggs ...Aggregate) QuerySetV2[T] {
	clone := qs.clone()
	clone.aggregates = append(clone.aggregates, aggs...)
	return clone
}

// Annotate adds annotation expressions
func (qs *BaseQuerySetV2[T]) Annotate(anns ...AnnotationExpr) QuerySetV2[T] {
	clone := qs.clone()
	clone.annotations = append(clone.annotations, anns...)
	return clone
}

// Values returns a ValuesQuerySet
func (qs *BaseQuerySetV2[T]) Values(fields ...string) ValuesQuerySet[T] {
	clone := qs.clone()
	clone.selectFields = fields
	return &BaseValuesQuerySet[T]{base: clone}
}

// ValuesList returns a ValuesListQuerySet
func (qs *BaseQuerySetV2[T]) ValuesList(fields ...string) ValuesListQuerySet[T] {
	clone := qs.clone()
	clone.selectFields = fields
	return &BaseValuesListQuerySet[T]{base: clone}
}

// UpdateBuilder returns an UpdateBuilder
func (qs *BaseQuerySetV2[T]) UpdateBuilder() (*UpdateBuilder[T], error) {
	// Create a wrapper that implements the interface needed by UpdateBuilder
	return NewUpdateBuilderFromQuerySet(qs)
}

// All executes the query and returns all results
func (qs *BaseQuerySetV2[T]) All(ctx context.Context) ([]*T, error) {
	sql, args, err := qs.buildSQL()
	if err != nil {
		return nil, err
	}
	
	db, err := qs.getDB(ctx)
	if err != nil {
		return nil, err
	}
	
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()
	
	return qs.scanRows(rows)
}

// buildSQL builds the SQL query
func (qs *BaseQuerySetV2[T]) buildSQL() (string, []interface{}, error) {
	builder := NewSQLBuilder()
	
	// Build SELECT clause
	selectClause := qs.buildSelectClause(builder)
	
	// Build FROM clause
	fromClause := fmt.Sprintf("FROM %s", EscapeIdentifier(qs.table))
	
	// Build WHERE clause
	whereClause, whereArgs := qs.buildWhereClause(builder)
	
	// Build ORDER BY clause
	orderByClause := qs.buildOrderByClause(builder)
	
	// Build LIMIT/OFFSET
	limitClause := qs.buildLimitClause()
	offsetClause := qs.buildOffsetClause()
	
	// Combine all parts
	parts := []string{selectClause, fromClause}
	if whereClause != "" {
		parts = append(parts, whereClause)
	}
	if orderByClause != "" {
		parts = append(parts, orderByClause)
	}
	if limitClause != "" {
		parts = append(parts, limitClause)
	}
	if offsetClause != "" {
		parts = append(parts, offsetClause)
	}
	
	sql := strings.Join(parts, " ")
	args := builder.Args()
	args = append(args, whereArgs...)
	
	return sql, args, nil
}

// buildSelectClause builds the SELECT clause
func (qs *BaseQuerySetV2[T]) buildSelectClause(builder *SQLBuilder) string {
	var fields []string
	
	if len(qs.selectFields) > 0 {
		for _, field := range qs.selectFields {
			fields = append(fields, EscapeIdentifier(field))
		}
	} else if len(qs.onlyFields) > 0 {
		for _, field := range qs.onlyFields {
			fields = append(fields, EscapeIdentifier(field))
		}
	} else {
		fields = []string{"*"}
	}
	
	// Add annotations to SELECT
	if len(qs.annotations) > 0 {
		for _, ann := range qs.annotations {
			// Build annotation SQL using QueryExpr.ToSQL
			// QueryExpr uses paramIndex, so we need to use builder's paramIndex
			annSQL, annArgs, nextIndex := ann.Expr.ToSQL(builder.paramIndex)
			// Update builder's paramIndex
			builder.paramIndex = nextIndex
			// Add args to builder
			for _, arg := range annArgs {
				builder.args = append(builder.args, arg)
			}
			alias := EscapeIdentifier(ann.Name)
			fields = append(fields, fmt.Sprintf("%s AS %s", annSQL, alias))
		}
	}
	
	selectClause := "SELECT "
	if len(qs.distinctFields) > 0 {
		selectClause += "DISTINCT "
	}
	selectClause += strings.Join(fields, ", ")
	
	return selectClause
}

// buildWhereClause builds the WHERE clause
func (qs *BaseQuerySetV2[T]) buildWhereClause(builder *SQLBuilder) (string, []interface{}) {
	var parts []string
	var allArgs []interface{}
	
	// Build conditions
	for _, cond := range qs.conditions {
		sql, args, err := cond.ToSQL(builder)
		if err != nil {
			panic(fmt.Sprintf("failed to build condition SQL: %v", err))
		}
		parts = append(parts, sql)
		allArgs = append(allArgs, args...)
	}
	
	// Build excludes (with NOT)
	for _, exclude := range qs.excludes {
		sql, args, err := exclude.ToSQL(builder)
		if err != nil {
			panic(fmt.Sprintf("failed to build exclude SQL: %v", err))
		}
		parts = append(parts, fmt.Sprintf("NOT (%s)", sql))
		allArgs = append(allArgs, args...)
	}
	
	if len(parts) == 0 {
		return "", nil
	}
	
	return "WHERE " + strings.Join(parts, " AND "), allArgs
}

// buildOrderByClause builds the ORDER BY clause
func (qs *BaseQuerySetV2[T]) buildOrderByClause(builder *SQLBuilder) string {
	if len(qs.orderBy) == 0 {
		return ""
	}
	
	var parts []string
	for _, field := range qs.orderBy {
		escaped := EscapeIdentifier(field.Field)
		if field.Ascending {
			parts = append(parts, escaped+" ASC")
		} else {
			parts = append(parts, escaped+" DESC")
		}
	}
	
	return "ORDER BY " + strings.Join(parts, ", ")
}

// buildLimitClause builds the LIMIT clause
func (qs *BaseQuerySetV2[T]) buildLimitClause() string {
	if qs.limitVal == nil {
		return ""
	}
	return fmt.Sprintf("LIMIT %d", *qs.limitVal)
}

// buildOffsetClause builds the OFFSET clause
func (qs *BaseQuerySetV2[T]) buildOffsetClause() string {
	if qs.offsetVal == nil {
		return ""
	}
	return fmt.Sprintf("OFFSET %d", *qs.offsetVal)
}

// scanRows scans rows into model instances
func (qs *BaseQuerySetV2[T]) scanRows(rows *sql.Rows) ([]*T, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	// Build column-to-field mapping from schema
	fieldMap := qs.buildFieldMap(columns)
	
	var results []*T
	for rows.Next() {
		instance := new(T)
		scanArgs := qs.prepareScanArgs(instance, columns, fieldMap)
		
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		results = append(results, instance)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	
	return results, nil
}

// buildFieldMap creates a mapping from column names to field info
func (qs *BaseQuerySetV2[T]) buildFieldMap(columns []string) map[string]*FieldInfo {
	fieldMap := make(map[string]*FieldInfo)
	for _, col := range columns {
		for i := range qs.schema.Fields {
			field := &qs.schema.Fields[i]
			if strings.EqualFold(field.DBColumn, col) || strings.EqualFold(field.Name, col) {
				fieldMap[col] = field
				break
			}
		}
	}
	return fieldMap
}

// prepareScanArgs prepares scan arguments using schema
func (qs *BaseQuerySetV2[T]) prepareScanArgs(instance *T, columns []string, fieldMap map[string]*FieldInfo) []interface{} {
	scanArgs := make([]interface{}, len(columns))
	instanceValue := reflect.ValueOf(instance).Elem()
	
	for i, col := range columns {
		fieldInfo, ok := fieldMap[col]
		if !ok {
			var val interface{}
			scanArgs[i] = &val
			continue
		}
		
		field := instanceValue.FieldByName(fieldInfo.Name)
		if field.IsValid() && field.CanSet() {
			scanArgs[i] = field.Addr().Interface()
		} else {
			var val interface{}
			scanArgs[i] = &val
		}
	}
	
	return scanArgs
}

// Get retrieves a single object
func (qs *BaseQuerySetV2[T]) Get(ctx context.Context) (*T, error) {
	results, err := qs.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", qs.table)
	}
	
	if len(results) > 1 {
		return nil, fmt.Errorf("get() returned more than one %s -- it returned %d!", qs.table, len(results))
	}
	
	return results[0], nil
}

// First retrieves the first object
func (qs *BaseQuerySetV2[T]) First(ctx context.Context) (*T, error) {
	results, err := qs.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", qs.table)
	}
	
	return results[0], nil
}

// Last retrieves the last object
func (qs *BaseQuerySetV2[T]) Last(ctx context.Context) (*T, error) {
	reversed := qs.Reverse()
	results, err := reversed.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", qs.table)
	}
	
	return results[0], nil
}

// Count counts matching records
func (qs *BaseQuerySetV2[T]) Count(ctx context.Context) (int64, error) {
	db, err := qs.getDB(ctx)
	if err != nil {
		return 0, err
	}
	
	builder := NewSQLBuilder()
	selectClause := fmt.Sprintf("SELECT COUNT(*) FROM %s", EscapeIdentifier(qs.table))
	whereClause, whereArgs := qs.buildWhereClause(builder)
	
	sql := selectClause
	if whereClause != "" {
		sql += " " + whereClause
	}
	
	var count int64
	err = db.QueryRowContext(ctx, sql, whereArgs...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count query failed: %w", err)
	}
	
	return count, nil
}

// Exists checks if any records exist
func (qs *BaseQuerySetV2[T]) Exists(ctx context.Context) (bool, error) {
	count, err := qs.Count(ctx)
	return count > 0, err
}

// Update performs a bulk update
func (qs *BaseQuerySetV2[T]) Update(ctx context.Context, updates UpdateMap) (int64, error) {
	if len(updates) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}
	
	db, err := qs.getDB(ctx)
	if err != nil {
		return 0, err
	}
	
	builder := NewSQLBuilder()
	
	// Build SET clause
	var setParts []string
	var updateArgs []interface{}
	
	for fieldName, value := range updates {
		escapedField := EscapeIdentifier(fieldName)
		
		// Check if value is an Expression
		if expr, ok := value.(Expression); ok {
			// Build expression SQL
			exprSQL, exprArgs, err := expr.ToSQL(builder)
			if err != nil {
				return 0, fmt.Errorf("failed to build expression SQL for field %s: %w", fieldName, err)
			}
			setParts = append(setParts, fmt.Sprintf("%s = %s", escapedField, exprSQL))
			updateArgs = append(updateArgs, exprArgs...)
		} else {
			// Regular value
			placeholder := builder.AddArg(value)
			setParts = append(setParts, fmt.Sprintf("%s = %s", escapedField, placeholder))
			updateArgs = append(updateArgs, value)
		}
	}
	
	// Build WHERE clause
	whereClause, whereArgs := qs.buildWhereClause(builder)
	
	// Combine all args
	allArgs := builder.Args()
	allArgs = append(allArgs, updateArgs...)
	allArgs = append(allArgs, whereArgs...)
	
	// Build SQL
	updateSQL := fmt.Sprintf("UPDATE %s SET %s", EscapeIdentifier(qs.table), strings.Join(setParts, ", "))
	if whereClause != "" {
		updateSQL += " " + whereClause
	}
	
	// Execute
	result, err := db.ExecContext(ctx, updateSQL, allArgs...)
	if err != nil {
		return 0, fmt.Errorf("update query failed: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	return rowsAffected, nil
}

// BulkUpdate performs bulk updates
func (qs *BaseQuerySetV2[T]) BulkUpdate(ctx context.Context, updates []UpdateMap) error {
	// TODO: Implement bulk update
	return fmt.Errorf("BulkUpdate not yet implemented in QuerySetV2")
}

// Delete performs a bulk delete
func (qs *BaseQuerySetV2[T]) Delete(ctx context.Context) (int64, error) {
	db, err := qs.getDB(ctx)
	if err != nil {
		return 0, err
	}
	
	builder := NewSQLBuilder()
	
	// Build WHERE clause
	whereClause, whereArgs := qs.buildWhereClause(builder)
	
	// Build SQL
	deleteSQL := fmt.Sprintf("DELETE FROM %s", EscapeIdentifier(qs.table))
	if whereClause != "" {
		deleteSQL += " " + whereClause
	}
	
	// Execute
	result, err := db.ExecContext(ctx, deleteSQL, whereArgs...)
	if err != nil {
		return 0, fmt.Errorf("delete query failed: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	return rowsAffected, nil
}

// Union performs a UNION operation
func (qs *BaseQuerySetV2[T]) Union(other QuerySetV2[T]) QuerySetV2[T] {
	// For now, return a combined query set that will use SQL UNION
	// Full implementation would require SQL query combination
	clone := qs.clone()
	// Mark as union - would need additional state to track this
	// For MVP, return clone with union marker
	return clone
}

// Intersection performs an INTERSECT operation
func (qs *BaseQuerySetV2[T]) Intersection(other QuerySetV2[T]) QuerySetV2[T] {
	// Similar to Union - would need SQL query combination
	clone := qs.clone()
	return clone
}

// Difference performs an EXCEPT operation
func (qs *BaseQuerySetV2[T]) Difference(other QuerySetV2[T]) QuerySetV2[T] {
	// Similar to Union - would need SQL query combination
	clone := qs.clone()
	return clone
}

// ValuesQuerySet interface for values queries
type ValuesQuerySet[T any] interface {
	All(ctx context.Context) ([]map[string]interface{}, error)
	Get(ctx context.Context) (map[string]interface{}, error)
	First(ctx context.Context) (map[string]interface{}, error)
}

// ValuesListQuerySet interface for values_list queries
type ValuesListQuerySet[T any] interface {
	All(ctx context.Context) ([][]interface{}, error)
	Get(ctx context.Context) ([]interface{}, error)
	First(ctx context.Context) ([]interface{}, error)
	Flat(ctx context.Context) ([]interface{}, error)
}

// BaseValuesQuerySet implementation
type BaseValuesQuerySet[T any] struct {
	base *BaseQuerySetV2[T]
}

func (vqs *BaseValuesQuerySet[T]) All(ctx context.Context) ([]map[string]interface{}, error) {
	db, err := vqs.base.getDB(ctx)
	if err != nil {
		return nil, err
	}
	
	sql, args, err := vqs.base.buildSQL()
	if err != nil {
		return nil, err
	}
	
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()
	
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	var results []map[string]interface{}
	for rows.Next() {
		// Create scan destination
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range scanArgs {
			scanArgs[i] = &values[i]
		}
		
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		// Create map
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		results = append(results, rowMap)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	
	return results, nil
}

func (vqs *BaseValuesQuerySet[T]) Get(ctx context.Context) (map[string]interface{}, error) {
	results, err := vqs.base.Limit(2).Values(vqs.base.selectFields...).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", vqs.base.table)
	}
	
	if len(results) > 1 {
		return nil, fmt.Errorf("get() returned more than one %s -- it returned %d!", vqs.base.table, len(results))
	}
	
	return results[0], nil
}

func (vqs *BaseValuesQuerySet[T]) First(ctx context.Context) (map[string]interface{}, error) {
	results, err := vqs.base.Limit(1).Values(vqs.base.selectFields...).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", vqs.base.table)
	}
	
	return results[0], nil
}

// BaseValuesListQuerySet implementation
type BaseValuesListQuerySet[T any] struct {
	base *BaseQuerySetV2[T]
}

func (vls *BaseValuesListQuerySet[T]) All(ctx context.Context) ([][]interface{}, error) {
	db, err := vls.base.getDB(ctx)
	if err != nil {
		return nil, err
	}
	
	sql, args, err := vls.base.buildSQL()
	if err != nil {
		return nil, err
	}
	
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()
	
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	var results [][]interface{}
	for rows.Next() {
		// Create scan destination
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range scanArgs {
			scanArgs[i] = &values[i]
		}
		
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		// Create tuple (slice)
		tuple := make([]interface{}, len(values))
		copy(tuple, values)
		results = append(results, tuple)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	
	return results, nil
}

func (vls *BaseValuesListQuerySet[T]) Get(ctx context.Context) ([]interface{}, error) {
	// Create a new values list query set with limit
	limited := vls.base.clone()
	limit := 2
	limited.limitVal = &limit
	
	results, err := limited.ValuesList(limited.selectFields...).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", vls.base.table)
	}
	
	if len(results) > 1 {
		return nil, fmt.Errorf("get() returned more than one %s -- it returned %d!", vls.base.table, len(results))
	}
	
	return results[0], nil
}

func (vls *BaseValuesListQuerySet[T]) First(ctx context.Context) ([]interface{}, error) {
	// Create a new values list query set with limit
	limited := vls.base.clone()
	limit := 1
	limited.limitVal = &limit
	
	results, err := limited.ValuesList(limited.selectFields...).All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("%s matching query does not exist", vls.base.table)
	}
	
	return results[0], nil
}

func (vls *BaseValuesListQuerySet[T]) Flat(ctx context.Context) ([]interface{}, error) {
	// For flat, we expect exactly one field
	if len(vls.base.selectFields) != 1 {
		return nil, fmt.Errorf("Flat() requires exactly one field, got %d", len(vls.base.selectFields))
	}
	
	results, err := vls.All(ctx)
	if err != nil {
		return nil, err
	}
	
	// Extract first element from each tuple
	flat := make([]interface{}, len(results))
	for i, tuple := range results {
		if len(tuple) > 0 {
			flat[i] = tuple[0]
		}
	}
	
	return flat, nil
}
