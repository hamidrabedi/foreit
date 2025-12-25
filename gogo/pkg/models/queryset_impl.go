package models

import (
	"context"
	"gorm.io/gorm"
)

// QuerySetImpl provides a type-safe QuerySet implementation
type QuerySetImpl[T any] struct {
	db    *gorm.DB
	model T
}

// NewQuerySet creates a new queryset
func NewQuerySet[T any](db *gorm.DB, model T) *QuerySetImpl[T] {
	return &QuerySetImpl[T]{
		db:    db,
		model: model,
	}
}

// Filter applies a type-safe condition
func (qs *QuerySetImpl[T]) Filter(condition *Condition[T]) *QuerySetImpl[T] {
	newQS := qs.clone()
	newQS.db = condition.Apply(newQS.db)
	return newQS
}

// Exclude applies a negated condition
func (qs *QuerySetImpl[T]) Exclude(condition *Condition[T]) *QuerySetImpl[T] {
	newQS := qs.clone()
	newQS.db = newQS.db.Not(condition.apply)
	return newQS
}

// OrderBy applies ordering using field references (type-safe)
func OrderByField[T any, F any](qs *QuerySetImpl[T], field *FieldRef[F, T], desc bool) *QuerySetImpl[T] {
	newQS := qs.clone()
	order := field.Column()
	if desc {
		order += " DESC"
	} else {
		order += " ASC"
	}
	newQS.db = newQS.db.Order(order)
	return newQS
}

// OrderByColumn applies ordering using a column name directly
func (qs *QuerySetImpl[T]) OrderByColumn(column string, desc bool) *QuerySetImpl[T] {
	newQS := qs.clone()
	if desc {
		newQS.db = newQS.db.Order(column + " DESC")
	} else {
		newQS.db = newQS.db.Order(column + " ASC")
	}
	return newQS
}

// Limit limits the number of results
func (qs *QuerySetImpl[T]) Limit(n int) *QuerySetImpl[T] {
	newQS := qs.clone()
	newQS.db = newQS.db.Limit(n)
	return newQS
}

// Offset sets the offset
func (qs *QuerySetImpl[T]) Offset(n int) *QuerySetImpl[T] {
	newQS := qs.clone()
	newQS.db = newQS.db.Offset(n)
	return newQS
}

// Get retrieves a single result
func (qs *QuerySetImpl[T]) Get(ctx context.Context) (*T, error) {
	var result T
	err := qs.db.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// All retrieves all results
func (qs *QuerySetImpl[T]) All(ctx context.Context) ([]*T, error) {
	var results []*T
	err := qs.db.WithContext(ctx).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Count returns the count
func (qs *QuerySetImpl[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := qs.db.WithContext(ctx).Model(qs.model).Count(&count).Error
	return count, err
}

// Exists checks if any results exist
func (qs *QuerySetImpl[T]) Exists(ctx context.Context) (bool, error) {
	count, err := qs.Count(ctx)
	return count > 0, err
}

// First returns the first result
func (qs *QuerySetImpl[T]) First(ctx context.Context) (*T, error) {
	var result T
	err := qs.db.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Last returns the last result
func (qs *QuerySetImpl[T]) Last(ctx context.Context) (*T, error) {
	var result T
	err := qs.db.WithContext(ctx).Last(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// clone creates a copy of the queryset
func (qs *QuerySetImpl[T]) clone() *QuerySetImpl[T] {
	return &QuerySetImpl[T]{
		db:    qs.db,
		model: qs.model,
	}
}

