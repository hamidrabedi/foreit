package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

// SelectRelated specifies fields to eagerly load via JOIN (like Django's select_related)
type SelectRelated struct {
	fields []string
}

// PrefetchRelated specifies relationships to prefetch (like Django's prefetch_related)
type PrefetchRelated struct {
	relationships []string
}

// QuerySetImpl with eager loading support
type QuerySetWithEager[T any] struct {
	*QuerySetImpl[T]
	selectRelated  []string
	prefetchRelated []string
}

// SelectRelated adds fields to eagerly load via JOIN
func (qs *QuerySetImpl[T]) SelectRelated(fields ...string) *QuerySetWithEager[T] {
	return &QuerySetWithEager[T]{
		QuerySetImpl:   qs,
		selectRelated:  fields,
		prefetchRelated: []string{},
	}
}

// PrefetchRelated adds relationships to prefetch
func (qs *QuerySetImpl[T]) PrefetchRelated(relationships ...string) *QuerySetWithEager[T] {
	return &QuerySetWithEager[T]{
		QuerySetImpl:   qs,
		selectRelated:  []string{},
		prefetchRelated: relationships,
	}
}

// WithEager adds both select_related and prefetch_related
func (qs *QuerySetImpl[T]) WithEager(selectFields []string, prefetchFields []string) *QuerySetWithEager[T] {
	return &QuerySetWithEager[T]{
		QuerySetImpl:   qs,
		selectRelated:  selectFields,
		prefetchRelated: prefetchFields,
	}
}

// All retrieves all results with eager loading
func (qs *QuerySetWithEager[T]) All(ctx context.Context) ([]*T, error) {
	query := qs.db.WithContext(ctx)
	
	// Apply select_related (JOINs)
	for _, field := range qs.selectRelated {
		query = query.Preload(field)
	}
	
	// Apply prefetch_related (separate queries)
	for _, rel := range qs.prefetchRelated {
		query = query.Preload(rel)
	}
	
	var results []*T
	err := query.Find(&results).Error
	return results, err
}

// Get retrieves a single result with eager loading
func (qs *QuerySetWithEager[T]) Get(ctx context.Context) (*T, error) {
	query := qs.db.WithContext(ctx)
	
	// Apply select_related
	for _, field := range qs.selectRelated {
		query = query.Preload(field)
	}
	
	// Apply prefetch_related
	for _, rel := range qs.prefetchRelated {
		query = query.Preload(rel)
	}
	
	var result T
	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// First retrieves the first result with eager loading
func (qs *QuerySetWithEager[T]) First(ctx context.Context) (*T, error) {
	query := qs.db.WithContext(ctx).Limit(1)
	
	// Apply select_related
	for _, field := range qs.selectRelated {
		query = query.Preload(field)
	}
	
	// Apply prefetch_related
	for _, rel := range qs.prefetchRelated {
		query = query.Preload(rel)
	}
	
	var result T
	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Filter continues the chain with eager loading
func (qs *QuerySetWithEager[T]) Filter(condition *Condition[T]) *QuerySetWithEager[T] {
	newQS := qs.QuerySetImpl.Filter(condition)
	return &QuerySetWithEager[T]{
		QuerySetImpl:   newQS,
		selectRelated:  qs.selectRelated,
		prefetchRelated: qs.prefetchRelated,
	}
}

// Exclude continues the chain with eager loading
func (qs *QuerySetWithEager[T]) Exclude(condition *Condition[T]) *QuerySetWithEager[T] {
	newQS := qs.QuerySetImpl.Exclude(condition)
	return &QuerySetWithEager[T]{
		QuerySetImpl:   newQS,
		selectRelated:  qs.selectRelated,
		prefetchRelated: qs.prefetchRelated,
	}
}

// OrderByColumn continues the chain with eager loading
func (qs *QuerySetWithEager[T]) OrderByColumn(column string, desc bool) *QuerySetWithEager[T] {
	newQS := qs.QuerySetImpl.OrderByColumn(column, desc)
	return &QuerySetWithEager[T]{
		QuerySetImpl:   newQS,
		selectRelated:  qs.selectRelated,
		prefetchRelated: qs.prefetchRelated,
	}
}

// Limit continues the chain with eager loading
func (qs *QuerySetWithEager[T]) Limit(n int) *QuerySetWithEager[T] {
	newQS := qs.QuerySetImpl.Limit(n)
	return &QuerySetWithEager[T]{
		QuerySetImpl:   newQS,
		selectRelated:  qs.selectRelated,
		prefetchRelated: qs.prefetchRelated,
	}
}

// Offset continues the chain with eager loading
func (qs *QuerySetWithEager[T]) Offset(n int) *QuerySetWithEager[T] {
	newQS := qs.QuerySetImpl.Offset(n)
	return &QuerySetWithEager[T]{
		QuerySetImpl:   newQS,
		selectRelated:  qs.selectRelated,
		prefetchRelated: qs.prefetchRelated,
	}
}

// EagerLoadConfig configures eager loading behavior
type EagerLoadConfig struct {
	// SelectRelated fields (JOINs)
	Select []string
	// PrefetchRelated relationships (separate queries)
	Prefetch []string
	// MaxDepth limits relationship depth
	MaxDepth int
}

// NewEagerLoadConfig creates a new eager load configuration
func NewEagerLoadConfig() *EagerLoadConfig {
	return &EagerLoadConfig{
		Select:   []string{},
		Prefetch: []string{},
		MaxDepth: 3,
	}
}

// WithSelect adds select_related fields
func (c *EagerLoadConfig) WithSelect(fields ...string) *EagerLoadConfig {
	c.Select = append(c.Select, fields...)
	return c
}

// WithPrefetch adds prefetch_related relationships
func (c *EagerLoadConfig) WithPrefetch(relationships ...string) *EagerLoadConfig {
	c.Prefetch = append(c.Prefetch, relationships...)
	return c
}

// WithMaxDepth sets the maximum relationship depth
func (c *EagerLoadConfig) WithMaxDepth(depth int) *EagerLoadConfig {
	c.MaxDepth = depth
	return c
}

// ApplyEagerLoading applies eager loading configuration to a query
func ApplyEagerLoading[T any](qs *QuerySetImpl[T], config *EagerLoadConfig) *QuerySetWithEager[T] {
	return qs.WithEager(config.Select, config.Prefetch)
}

