package models

import (
	"context"
	"fmt"
)

// QuerySetImpl implements the QuerySet interface
type QuerySetImpl struct {
	modelType  string
	manager    Manager
	filters    []FilterCondition
	excludes   []FilterCondition
	orderBy    []OrderField
	limitVal   *int
	offsetVal  *int
}

// FilterCondition represents a filter condition
type FilterCondition struct {
	Field    string
	Operator string
	Value    interface{}
}

// OrderField represents an ordering field
type OrderField struct {
	Field string
	Desc  bool
}

// NewQuerySet creates a new queryset
func NewQuerySet(modelType string, manager Manager) *QuerySetImpl {
	return &QuerySetImpl{
		modelType: modelType,
		manager:   manager,
		filters:   make([]FilterCondition, 0),
		excludes:  make([]FilterCondition, 0),
		orderBy:   make([]OrderField, 0),
	}
}

// Filter adds a filter condition
func (qs *QuerySetImpl) Filter(condition interface{}) QuerySet {
	newQS := qs.clone()
	
	switch cond := condition.(type) {
	case map[string]interface{}:
		for field, value := range cond {
			newQS.filters = append(newQS.filters, FilterCondition{
				Field:    field,
				Operator: "=",
				Value:    value,
			})
		}
	case Q:
		newQS.filters = append(newQS.filters, cond.toFilterConditions()...)
	default:
		panic(fmt.Sprintf("unsupported filter condition type: %T", condition))
	}
	
	return newQS
}

// Exclude adds an exclude condition
func (qs *QuerySetImpl) Exclude(condition interface{}) QuerySet {
	newQS := qs.clone()
	
	switch cond := condition.(type) {
	case map[string]interface{}:
		for field, value := range cond {
			newQS.excludes = append(newQS.excludes, FilterCondition{
				Field:    field,
				Operator: "=",
				Value:    value,
			})
		}
	case Q:
		newQS.excludes = append(newQS.excludes, cond.toFilterConditions()...)
	default:
		panic(fmt.Sprintf("unsupported exclude condition type: %T", condition))
	}
	
	return newQS
}

// OrderBy adds ordering
func (qs *QuerySetImpl) OrderBy(fields ...string) QuerySet {
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

// Limit limits the number of results
func (qs *QuerySetImpl) Limit(n int) QuerySet {
	newQS := qs.clone()
	newQS.limitVal = &n
	return newQS
}

// Offset sets the offset
func (qs *QuerySetImpl) Offset(n int) QuerySet {
	newQS := qs.clone()
	newQS.offsetVal = &n
	return newQS
}

// Get retrieves a single result
func (qs *QuerySetImpl) Get(ctx context.Context) (Model, error) {
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

// All retrieves all results
func (qs *QuerySetImpl) All(ctx context.Context) ([]Model, error) {
	if querysetManager, ok := qs.manager.(QuerySetManager); ok {
		return querysetManager.QuerySetAll(ctx, qs)
	}
	return nil, fmt.Errorf("manager does not support queryset operations")
}

// Count returns the count
func (qs *QuerySetImpl) Count(ctx context.Context) (int, error) {
	if querysetManager, ok := qs.manager.(QuerySetManager); ok {
		return querysetManager.QuerySetCount(ctx, qs)
	}
	return 0, fmt.Errorf("manager does not support queryset operations")
}

// Exists checks if any results exist
func (qs *QuerySetImpl) Exists(ctx context.Context) (bool, error) {
	count, err := qs.Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// First returns the first result
func (qs *QuerySetImpl) First(ctx context.Context) (Model, error) {
	limited := qs.Limit(1)
	results, err := limited.All(ctx)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return nil, ErrDoesNotExist
	}
	
	return results[0], nil
}

// Last returns the last result
func (qs *QuerySetImpl) Last(ctx context.Context) (Model, error) {
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

// clone creates a copy of the queryset
func (qs *QuerySetImpl) clone() *QuerySetImpl {
	return &QuerySetImpl{
		modelType: qs.modelType,
		manager:   qs.manager,
		filters:   append([]FilterCondition{}, qs.filters...),
		excludes:  append([]FilterCondition{}, qs.excludes...),
		orderBy:   append([]OrderField{}, qs.orderBy...),
		limitVal:  qs.limitVal,
		offsetVal: qs.offsetVal,
	}
}

// QuerySetManager extends Manager with queryset support
type QuerySetManager interface {
	Manager
	QuerySetAll(ctx context.Context, qs *QuerySetImpl) ([]Model, error)
	QuerySetCount(ctx context.Context, qs *QuerySetImpl) (int, error)
}

var (
	ErrDoesNotExist    = errors.New("model does not exist")
	ErrMultipleResults = errors.New("multiple results returned")
)

