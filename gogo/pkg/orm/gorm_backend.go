package orm

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// GORMBackend wraps GORM as the database backend
type GORMBackend struct {
	db *gorm.DB
}

// NewGORMBackend creates a new GORM backend
func NewGORMBackend(db *gorm.DB) *GORMBackend {
	return &GORMBackend{db: db}
}

// GetDB returns the underlying GORM DB
func (b *GORMBackend) GetDB() *gorm.DB {
	return b.db
}

// Query executes a query and scans results
func (b *GORMBackend) Query(ctx context.Context, dest interface{}, query *QueryPlan) error {
	gormQuery := b.db.WithContext(ctx)
	
	// Apply filters
	for _, filter := range query.Filters {
		gormQuery = b.applyFilter(gormQuery, filter)
	}
	
	// Apply excludes
	for _, exclude := range query.Excludes {
		gormQuery = b.applyExclude(gormQuery, exclude)
	}
	
	// Apply ordering
	if len(query.Orders) > 0 {
		for _, order := range query.Orders {
			if order.Desc {
				gormQuery = gormQuery.Order(fmt.Sprintf("%s DESC", order.Field))
			} else {
				gormQuery = gormQuery.Order(fmt.Sprintf("%s ASC", order.Field))
			}
		}
	}
	
	// Apply limit/offset
	if query.Limit != nil {
		gormQuery = gormQuery.Limit(*query.Limit)
	}
	if query.Offset != nil {
		gormQuery = gormQuery.Offset(*query.Offset)
	}
	
	// Execute query
	return gormQuery.Find(dest).Error
}

// Count executes a count query
func (b *GORMBackend) Count(ctx context.Context, model interface{}, query *QueryPlan) (int64, error) {
	gormQuery := b.db.WithContext(ctx).Model(model)
	
	// Apply filters
	for _, filter := range query.Filters {
		gormQuery = b.applyFilter(gormQuery, filter)
	}
	
	// Apply excludes
	for _, exclude := range query.Excludes {
		gormQuery = b.applyExclude(gormQuery, exclude)
	}
	
	var count int64
	err := gormQuery.Count(&count).Error
	return count, err
}

// Create creates a record
func (b *GORMBackend) Create(ctx context.Context, model interface{}) error {
	return b.db.WithContext(ctx).Create(model).Error
}

// Update updates a record
func (b *GORMBackend) Update(ctx context.Context, model interface{}) error {
	return b.db.WithContext(ctx).Save(model).Error
}

// Delete deletes a record
func (b *GORMBackend) Delete(ctx context.Context, model interface{}) error {
	return b.db.WithContext(ctx).Delete(model).Error
}

// Get retrieves a record by ID
func (b *GORMBackend) Get(ctx context.Context, model interface{}, id interface{}) error {
	return b.db.WithContext(ctx).First(model, id).Error
}

// applyFilter applies a filter condition to GORM query
func (b *GORMBackend) applyFilter(query *gorm.DB, filter FilterCondition) *gorm.DB {
	switch filter.Operator {
	case "=", "exact":
		return query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case "!=", "ne":
		return query.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case ">", "gt":
		return query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case ">=", "gte":
		return query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case "<", "lt":
		return query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case "<=", "lte":
		return query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case "in":
		return query.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)
	case "contains", "like":
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case "icontains", "ilike":
		return query.Where(fmt.Sprintf("%s ILIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case "startswith":
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%v%%", filter.Value))
	case "endswith":
		return query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v", filter.Value))
	case "isnull":
		if filter.Value.(bool) {
			return query.Where(fmt.Sprintf("%s IS NULL", filter.Field))
		}
		return query.Where(fmt.Sprintf("%s IS NOT NULL", filter.Field))
	}
	return query
}

// applyExclude applies an exclude condition
func (b *GORMBackend) applyExclude(query *gorm.DB, exclude FilterCondition) *gorm.DB {
	switch exclude.Operator {
	case "=", "exact":
		return query.Not(fmt.Sprintf("%s = ?", exclude.Field), exclude.Value)
	case "!=", "ne":
		return query.Where(fmt.Sprintf("%s = ?", exclude.Field), exclude.Value)
	case "in":
		return query.Not(fmt.Sprintf("%s IN ?", exclude.Field), exclude.Value)
	default:
		return query.Not(fmt.Sprintf("%s = ?", exclude.Field), exclude.Value)
	}
}

// QueryPlan represents a query plan
type QueryPlan struct {
	Table    string
	Filters  []FilterCondition
	Excludes []FilterCondition
	Orders   []OrderField
	Limit    *int
	Offset   *int
}

// FilterCondition represents a filter
type FilterCondition struct {
	Field    string
	Operator string
	Value    interface{}
}

// OrderField represents ordering
type OrderField struct {
	Field string
	Desc  bool
}

