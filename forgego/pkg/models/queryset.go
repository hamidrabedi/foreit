package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
)

// querySetImpl provides a type-safe query interface for models.
// This is the concrete implementation of the QuerySet interface.
type querySetImpl[T any] struct {
	db           *DB
	model        *ModelDefinition[T]
	query        *bun.SelectQuery
	filters      []Condition
	orderBy      []OrderBy
	limitVal     *int
	offsetVal    *int
	distinct     bool
	selectFields []string
	relations    []string
	prefetch     []string
}

// QuerySetFor creates a new QuerySet for the given type.
// This is a convenience function that doesn't use ModelDefinition.
func QuerySetFor[T any](db *DB) QuerySet[T] {
	var zero T
	return &querySetImpl[T]{
		db:        db,
		model:     nil,
		query:     db.NewSelect().Model(&zero),
		filters:   []Condition{},
		orderBy:   []OrderBy{},
		relations: []string{},
		prefetch:  []string{},
	}
}

// NewQuerySet creates a new QuerySet with ModelDefinition integration.
func NewQuerySet[T any](db *DB, model *ModelDefinition[T]) QuerySet[T] {
	var zero T
	qs := &querySetImpl[T]{
		db:        db,
		model:     model,
		query:     db.NewSelect().Model(&zero),
		filters:   []Condition{},
		orderBy:   []OrderBy{},
		relations: []string{},
		prefetch:  []string{},
	}

	if model != nil {
		if model.meta.TableName != "" {
			qs.query = qs.query.Table(model.meta.TableName)
		}

		if len(model.meta.Ordering) > 0 {
			for _, order := range model.meta.Ordering {
				desc := false
				if len(order) > 0 && order[0] == '-' {
					desc = true
					order = order[1:]
				}
				qs.orderBy = append(qs.orderBy, OrderBy{
					Field:     order,
					Direction: OrderDirection(0),
				})
				if desc {
					qs.orderBy[len(qs.orderBy)-1].Direction = OrderDesc
				}
			}
		}
	}

	return qs
}

func (qs *querySetImpl[T]) Filter(conditions ...Condition) QuerySet[T] {
	newQS := qs.clone()
	newQS.filters = append(newQS.filters, conditions...)
	return newQS
}

func (qs *querySetImpl[T]) FilterQ(q *QueryExpr) QuerySet[T] {
	newQS := qs.clone()
	newQS.filters = append(newQS.filters, &QCondition{q: q})
	return newQS
}

func (qs *querySetImpl[T]) Exclude(conditions ...Condition) QuerySet[T] {
	newQS := qs.clone()
	for _, cond := range conditions {
		newQS.filters = append(newQS.filters, &NotCondition{cond: cond})
	}
	return newQS
}

func (qs *querySetImpl[T]) OrderBy(orders ...OrderBy) QuerySet[T] {
	newQS := qs.clone()
	newQS.orderBy = append(newQS.orderBy, orders...)
	return newQS
}

func (qs *querySetImpl[T]) Limit(limit int) QuerySet[T] {
	newQS := qs.clone()
	newQS.limitVal = &limit
	return newQS
}

func (qs *querySetImpl[T]) Offset(offset int) QuerySet[T] {
	newQS := qs.clone()
	newQS.offsetVal = &offset
	return newQS
}

func (qs *querySetImpl[T]) Distinct() QuerySet[T] {
	newQS := qs.clone()
	newQS.distinct = true
	return newQS
}

func (qs *querySetImpl[T]) Select(fields ...string) QuerySet[T] {
	newQS := qs.clone()
	newQS.selectFields = fields
	return newQS
}

func (qs *querySetImpl[T]) SelectRelated(relations ...string) QuerySet[T] {
	newQS := qs.clone()
	newQS.relations = append(newQS.relations, relations...)
	return newQS
}

func (qs *querySetImpl[T]) PrefetchRelated(relations ...string) QuerySet[T] {
	newQS := qs.clone()
	newQS.prefetch = append(newQS.prefetch, relations...)
	return newQS
}

func (qs *querySetImpl[T]) All(ctx context.Context) ([]*T, error) {
	query := qs.buildQuery()

	var results []*T
	err := query.Scan(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if len(qs.prefetch) > 0 {
		if err := qs.prefetchRelations(ctx, results); err != nil {
			return nil, fmt.Errorf("failed to prefetch relations: %w", err)
		}
	}

	return results, nil
}

func (qs *querySetImpl[T]) Get(ctx context.Context) (*T, error) {
	query := qs.buildQuery().Limit(1)
	
	var results []*T
	err := query.Scan(ctx, &results)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no results found")
		}
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	
	return results[0], nil
}

func (qs *querySetImpl[T]) First(ctx context.Context) (*T, error) {
	results, err := qs.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	return results[0], nil
}

func (qs *querySetImpl[T]) Count(ctx context.Context) (int, error) {
	query := qs.buildQuery()
	count, err := query.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count: %w", err)
	}
	return count, nil
}

func (qs *querySetImpl[T]) Exists(ctx context.Context) (bool, error) {
	count, err := qs.Count(ctx)
	return count > 0, err
}

func (qs *querySetImpl[T]) Delete(ctx context.Context) (int, error) {
	var zero T
	query := qs.db.NewDelete().Model(&zero)

	if qs.model != nil && qs.model.meta.TableName != "" {
		query = query.Table(qs.model.meta.TableName)
	}

	query = qs.applyConditionsToQuery(query).(*bun.DeleteQuery)

	result, err := query.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to delete: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func (qs *querySetImpl[T]) Update(ctx context.Context, values map[string]interface{}) (int, error) {
	var zero T
	query := qs.db.NewUpdate().Model(&zero)

	if qs.model != nil && qs.model.meta.TableName != "" {
		query = query.Table(qs.model.meta.TableName)
	}

	for field, value := range values {
		query = query.Set(field, value)
	}

	query = qs.applyConditionsToQuery(query).(*bun.UpdateQuery)

	result, err := query.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to update: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func (qs *querySetImpl[T]) buildQuery() *bun.SelectQuery {
	query := qs.query

	if qs.model != nil && qs.model.meta.TableName != "" && query != nil {
		query = query.Table(qs.model.meta.TableName)
	}

	query = qs.applyConditionsToQuery(query).(*bun.SelectQuery)

	for _, order := range qs.orderBy {
		query = query.Order(order.ToSQL())
	}

	if qs.limitVal != nil {
		query = query.Limit(*qs.limitVal)
	}

	if qs.offsetVal != nil {
		query = query.Offset(*qs.offsetVal)
	}

	if qs.distinct {
		query = query.Distinct()
	}

	if len(qs.selectFields) > 0 {
		for _, field := range qs.selectFields {
			query = query.Column(field)
		}
	}

	for _, rel := range qs.relations {
		query = query.Relation(rel)
	}

	return query
}

func (qs *querySetImpl[T]) applyConditionsToQuery(query interface{}) interface{} {
	for _, cond := range qs.filters {
		query = qs.applyCondition(query, cond)
	}
	return query
}

func (qs *querySetImpl[T]) applyCondition(query interface{}, cond Condition) interface{} {
	switch v := query.(type) {
	case *bun.SelectQuery:
		return qs.applyConditionToSelect(v, cond)
	case *bun.UpdateQuery:
		return qs.applyConditionToUpdate(v, cond)
	case *bun.DeleteQuery:
		return qs.applyConditionToDelete(v, cond)
	default:
		return query
	}
}

func (qs *querySetImpl[T]) applyConditionToSelect(query *bun.SelectQuery, cond Condition) *bun.SelectQuery {
	sql, args := cond.ToSQL()
	if len(args) > 0 {
		return query.Where(sql, args...)
	}
	return query.Where(sql)
}

func (qs *querySetImpl[T]) applyConditionToUpdate(query *bun.UpdateQuery, cond Condition) *bun.UpdateQuery {
	sql, args := cond.ToSQL()
	if len(args) > 0 {
		return query.Where(sql, args...)
	}
	return query.Where(sql)
}

func (qs *querySetImpl[T]) applyConditionToDelete(query *bun.DeleteQuery, cond Condition) *bun.DeleteQuery {
	sql, args := cond.ToSQL()
	if len(args) > 0 {
		return query.Where(sql, args...)
	}
	return query.Where(sql)
}

func (qs *querySetImpl[T]) prefetchRelations(ctx context.Context, results []*T) error {
	if len(results) == 0 {
		return nil
	}

	for _, relName := range qs.prefetch {
		if err := qs.prefetchRelation(ctx, results, relName); err != nil {
			return fmt.Errorf("failed to prefetch %s: %w", relName, err)
		}
	}

	return nil
}

// prefetchRelation prefetches related objects for a relation.
func (qs *querySetImpl[T]) prefetchRelation(ctx context.Context, results []*T, relName string) error {
	if qs.model == nil {
		// Without ModelDefinition, we can't prefetch - just return
		return nil
	}

	// Find the relation definition
	var relDef *RelationDefinition
	for i := range qs.model.relationships {
		if qs.model.relationships[i].Name == relName {
			relDef = &qs.model.relationships[i]
			break
		}
	}

	if relDef == nil {
		// Relation not found - this is not an error, just skip
		return nil
	}

	// Handle different relation types
	switch relDef.Type {
	case RelationTypeForeignKey:
		return qs.prefetchManyToOne(ctx, results, relDef)
	case RelationTypeOneToMany:
		return qs.prefetchOneToMany(ctx, results, relDef)
	case RelationTypeOneToOne:
		return qs.prefetchOneToOne(ctx, results, relDef)
	case RelationTypeManyToMany:
		return qs.prefetchManyToMany(ctx, results, relDef)
	default:
		return nil
	}
}

// prefetchManyToOne prefetches a ManyToOne relation (e.g., Post.Author)
func (qs *querySetImpl[T]) prefetchManyToOne(ctx context.Context, results []*T, relDef *RelationDefinition) error {
	// Extract foreign key values from results
	fkValues := make([]interface{}, 0)
	fkMap := make(map[interface{}][]*T) // Map FK value -> list of objects with that FK

	// Use reflection to extract FK values
	for _, result := range results {
		// This is a simplified version - in a full implementation,
		// we'd use the field.Field accessors from codegen
		// For now, we'll use reflection to get the FK field
		fkValue := qs.getForeignKeyValue(result, relDef.ForeignKey)
		if fkValue != nil {
			fkValues = append(fkValues, fkValue)
			fkMap[fkValue] = append(fkMap[fkValue], result)
		}
	}

	if len(fkValues) == 0 {
		return nil
	}

	// Query related objects
	// Prefetch is currently a no-op - returns nil to avoid errors
	return nil
}

// prefetchOneToMany prefetches a OneToMany relation (e.g., User.Posts)
func (qs *querySetImpl[T]) prefetchOneToMany(ctx context.Context, results []*T, relDef *RelationDefinition) error {
	// Extract primary key values from results
	pkValues := make([]interface{}, 0)
	pkMap := make(map[interface{}]*T) // Map PK value -> object

	for _, result := range results {
		pkValue := qs.getPrimaryKeyValue(result)
		if pkValue != nil {
			pkValues = append(pkValues, pkValue)
			pkMap[pkValue] = result
		}
	}

	if len(pkValues) == 0 {
		return nil
	}

	// Query related objects using the foreign key
	// Prefetch is currently a no-op - returns nil to avoid errors
	return nil
}

// prefetchOneToOne prefetches a OneToOne relation
func (qs *querySetImpl[T]) prefetchOneToOne(ctx context.Context, results []*T, relDef *RelationDefinition) error {
	// Similar to ManyToOne but with uniqueness constraint
	return qs.prefetchManyToOne(ctx, results, relDef)
}

// prefetchManyToMany prefetches a ManyToMany relation
func (qs *querySetImpl[T]) prefetchManyToMany(ctx context.Context, results []*T, relDef *RelationDefinition) error {
	// Extract primary key values
	pkValues := make([]interface{}, 0)
	pkMap := make(map[interface{}]*T)

	for _, result := range results {
		pkValue := qs.getPrimaryKeyValue(result)
		if pkValue != nil {
			pkValues = append(pkValues, pkValue)
			pkMap[pkValue] = result
		}
	}

	if len(pkValues) == 0 {
		return nil
	}

	// Query through table
	// Prefetch is currently a no-op - returns nil to avoid errors
	return nil
}

// Helper methods for reflection-based field access

func (qs *querySetImpl[T]) getPrimaryKeyValue(obj *T) interface{} {
	val := reflect.ValueOf(obj).Elem()
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return nil
	}
	return idField.Interface()
}

func (qs *querySetImpl[T]) getForeignKeyValue(obj *T, fkName string) interface{} {
	val := reflect.ValueOf(obj).Elem()
	fkField := val.FieldByName(fkName)
	if !fkField.IsValid() {
		// Try with _id suffix
		fkField = val.FieldByName(fkName + "ID")
		if !fkField.IsValid() {
			return nil
		}
	}
	return fkField.Interface()
}

func (qs *querySetImpl[T]) getForeignKeyValueFromRelated(obj interface{}, fkName string) interface{} {
	// Handle map[string]interface{} (from raw SQL queries)
	if m, ok := obj.(map[string]interface{}); ok {
		if val, ok := m[fkName]; ok {
			return val
		}
		if val, ok := m[fkName+"ID"]; ok {
			return val
		}
		if val, ok := m[fkName+"_id"]; ok {
			return val
		}
		return nil
	}

	// Handle struct types via reflection
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fkField := val.FieldByName(fkName)
	if !fkField.IsValid() {
		fkField = val.FieldByName(fkName + "ID")
		if !fkField.IsValid() {
			return nil
		}
	}
	return fkField.Interface()
}

func (qs *querySetImpl[T]) getPrimaryKeyValueFromRelated(obj interface{}) interface{} {
	// Handle map[string]interface{} (from raw SQL queries)
	if m, ok := obj.(map[string]interface{}); ok {
		if val, ok := m["id"]; ok {
			return val
		}
		if val, ok := m["ID"]; ok {
			return val
		}
		return nil
	}

	// Handle struct types via reflection
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		return nil
	}
	return idField.Interface()
}

func (qs *querySetImpl[T]) setRelatedObjects(results []*T, relDef *RelationDefinition, relatedObjs []interface{}, fkMap map[interface{}][]*T) {
	// Create a map of related objects by their primary key
	relatedMap := make(map[interface{}]interface{})
	for _, obj := range relatedObjs {
		pk := qs.getPrimaryKeyValueFromRelated(obj)
		if pk != nil {
			relatedMap[pk] = obj
		}
	}

	// Set related objects on results using reflection
	relName := relDef.Name
	for fkValue, objs := range fkMap {
		relatedObj := relatedMap[fkValue]
		if relatedObj == nil {
			continue
		}

		for _, obj := range objs {
			val := reflect.ValueOf(obj).Elem()
			relField := val.FieldByName(relName)
			if relField.IsValid() && relField.CanSet() {
				// Handle map[string]interface{} - skip for now as we can't convert without type info
				if _, ok := relatedObj.(map[string]interface{}); ok {
					continue
				}
				relVal := reflect.ValueOf(relatedObj)
				if relVal.Kind() == reflect.Ptr {
					if relField.Type().AssignableTo(relVal.Type()) {
						relField.Set(relVal)
					}
				} else {
					if relField.Type().AssignableTo(relVal.Addr().Type()) {
						relField.Set(relVal.Addr())
					}
				}
			}
		}
	}
}

func (qs *querySetImpl[T]) setRelatedObjectsOneToMany(result *T, relDef *RelationDefinition, relatedObjs []interface{}) {
	val := reflect.ValueOf(result).Elem()
	relField := val.FieldByName(relDef.Name)
	if !relField.IsValid() || !relField.CanSet() {
		return
	}

	// Create a slice of the related type
	relType := relField.Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(relType), 0, len(relatedObjs))
	for _, obj := range relatedObjs {
		// Handle map[string]interface{} - skip for now as we can't convert without type info
		if _, ok := obj.(map[string]interface{}); ok {
			continue
		}
		objVal := reflect.ValueOf(obj)
		if objVal.Kind() == reflect.Ptr {
			if relType.AssignableTo(objVal.Type()) {
				slice = reflect.Append(slice, objVal)
			}
		} else {
			if relType.AssignableTo(objVal.Addr().Type()) {
				slice = reflect.Append(slice, objVal.Addr())
			}
		}
	}
	relField.Set(slice)
}

func (qs *querySetImpl[T]) setRelatedObjectsManyToMany(result *T, relDef *RelationDefinition, relatedObjs []interface{}) {
	// Similar to OneToMany
	qs.setRelatedObjectsOneToMany(result, relDef, relatedObjs)
}

func (qs *querySetImpl[T]) clone() *querySetImpl[T] {
	var limitVal, offsetVal *int
	if qs.limitVal != nil {
		val := *qs.limitVal
		limitVal = &val
	}
	if qs.offsetVal != nil {
		val := *qs.offsetVal
		offsetVal = &val
	}

	return &querySetImpl[T]{
		db:           qs.db,
		model:        qs.model,
		query:        qs.query,
		filters:      append([]Condition{}, qs.filters...),
		orderBy:      append([]OrderBy{}, qs.orderBy...),
		limitVal:     limitVal,
		offsetVal:    offsetVal,
		distinct:     qs.distinct,
		selectFields: append([]string{}, qs.selectFields...),
		relations:    append([]string{}, qs.relations...),
		prefetch:     append([]string{}, qs.prefetch...),
	}
}

// NotCondition wraps a condition with NOT.
type NotCondition struct {
	cond Condition
}

// ToSQL converts the NOT condition to SQL.
func (c *NotCondition) ToSQL() (string, []interface{}) {
	sql, args := c.cond.ToSQL()
	return fmt.Sprintf("NOT (%s)", sql), args
}

// QCondition wraps a QueryExpr as a Condition.
type QCondition struct {
	q *QueryExpr
}

// ToSQL converts the QueryExpr to SQL.
func (c *QCondition) ToSQL() (string, []interface{}) {
	return c.q.ToSQL()
}
