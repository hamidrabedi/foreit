package models

import (
	"context"
	"fmt"
)

// QuerySetImpl is a generic, immutable queryset implementation
type QuerySetImpl[T Model] struct {
	repo     Repository[T]
	meta     *ModelMeta
	filters  []FilterCondition
	excludes []FilterCondition
	orderBy  []OrderField
	limitVal *int
	offsetVal *int
}

// NewQuerySetImpl creates a new type-safe queryset
func NewQuerySetImpl[T Model](repo Repository[T], meta *ModelMeta) *QuerySetImpl[T] {
	return &QuerySetImpl[T]{
		repo:    repo,
		meta:    meta,
		filters: make([]FilterCondition, 0),
		excludes: make([]FilterCondition, 0),
		orderBy: make([]OrderField, 0),
	}
}

// Filter adds a filter condition (returns new immutable QuerySet)
func (qs *QuerySetImpl[T]) Filter(condition Q) QuerySet[T] {
	newQS := qs.clone()
	
	for field, value := range condition {
		// Parse Django-style lookups (field__lookup)
		fieldName, lookup := parseLookup(field)
		newQS.filters = append(newQS.filters, FilterCondition{
			Field:    fieldName,
			Operator: lookup,
			Value:    value,
		})
	}
	
	return newQS
}

// Exclude adds an exclude condition (returns new immutable QuerySet)
func (qs *QuerySetImpl[T]) Exclude(condition Q) QuerySet[T] {
	newQS := qs.clone()
	
	for field, value := range condition {
		fieldName, lookup := parseLookup(field)
		newQS.excludes = append(newQS.excludes, FilterCondition{
			Field:    fieldName,
			Operator: lookup,
			Value:    value,
		})
	}
	
	return newQS
}

// OrderBy adds ordering (returns new immutable QuerySet)
func (qs *QuerySetImpl[T]) OrderBy(fields ...string) QuerySet[T] {
	newQS := qs.clone()
	
	for _, field := range fields {
		desc := false
		if len(field) > 0 && field[0] == '-' {
			desc = true
			field = field[1:]
		}
		newQS.orderBy = append(newQS.orderBy, OrderField{
			Field: field,
			Desc:  desc,
		})
	}
	
	return newQS
}

// Limit limits results (returns new immutable QuerySet)
func (qs *QuerySetImpl[T]) Limit(n int) QuerySet[T] {
	newQS := qs.clone()
	newQS.limitVal = &n
	return newQS
}

// Offset sets offset (returns new immutable QuerySet)
func (qs *QuerySetImpl[T]) Offset(n int) QuerySet[T] {
	newQS := qs.clone()
	newQS.offsetVal = &n
	return newQS
}

// Get retrieves a single result
func (qs *QuerySetImpl[T]) Get(ctx context.Context) (*T, error) {
	results, err := qs.All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, ErrDoesNotExist
	}
	
	if len(results) > 1 {
		return nil, ErrMultipleResults
	}
	
	return results[0], nil
}

// All retrieves all results (type-safe!)
func (qs *QuerySetImpl[T]) All(ctx context.Context) ([]*T, error) {
	return qs.repo.Query(ctx, qs)
}

// Count returns the count
func (qs *QuerySetImpl[T]) Count(ctx context.Context) (int, error) {
	return qs.repo.Count(ctx, qs)
}

// Exists checks if any results exist
func (qs *QuerySetImpl[T]) Exists(ctx context.Context) (bool, error) {
	count, err := qs.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// First returns the first result (type-safe!)
func (qs *QuerySetImpl[T]) First(ctx context.Context) (*T, error) {
	limited := qs.Limit(1)
	return limited.All(ctx)
}

// Last returns the last result (type-safe!)
func (qs *QuerySetImpl[T]) Last(ctx context.Context) (*T, error) {
	reversed := qs.OrderBy("-id").Limit(1)
	results, err := reversed.All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, ErrDoesNotExist
	}
	
	return results[0], nil
}

// clone creates an immutable copy
func (qs *QuerySetImpl[T]) clone() *QuerySetImpl[T] {
	return &QuerySetImpl[T]{
		repo:     qs.repo,
		meta:     qs.meta,
		filters:  append([]FilterCondition{}, qs.filters...),
		excludes: append([]FilterCondition{}, qs.excludes...),
		orderBy:  append([]OrderField{}, qs.orderBy...),
		limitVal: qs.limitVal,
		offsetVal: qs.offsetVal,
	}
}

// parseLookup parses Django-style field lookups (field__lookup)
func parseLookup(field string) (string, string) {
	// Simple implementation - can be enhanced
	if idx := findLast(field, "__"); idx != -1 {
		return field[:idx], field[idx+2:]
	}
	return field, "exact"
}

func findLast(s, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

