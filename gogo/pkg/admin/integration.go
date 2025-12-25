package admin

import (
	"context"
	"github.com/gogo/pkg/models"
	"gorm.io/gorm"
)

// Admin integrates admin with type-safe queries
type Admin[T any] struct {
	*BaseModelAdmin[T]
	db       *gorm.DB
	manager  *models.ManagerImpl[T]
}

// NewAdmin creates a new admin instance
func NewAdmin[T any](db *gorm.DB, modelName string) *Admin[T] {
	model := *new(T)
	manager := models.NewManager(db, model)
	
	return &Admin[T]{
		BaseModelAdmin: NewBaseModelAdmin[T](modelName),
		db:             db,
		manager:        manager,
	}
}

// GetQueryset returns a type-safe queryset
func (a *Admin[T]) GetQueryset(ctx context.Context) *models.QuerySetImpl[T] {
	return a.manager.Filter(ctx)
}

// ApplyFilters applies admin filters to a queryset
// Note: This requires mapping admin FieldRef to models FieldRef
// For full type safety, use models.FieldRef directly in FilterSpec
func (a *Admin[T]) ApplyFilters(ctx context.Context, queryset *models.QuerySetImpl[T], filters map[string]interface{}) *models.QuerySetImpl[T] {
	for _, spec := range a.ListFilter() {
		fieldName := spec.Field.Name()
		if value, ok := filters[fieldName]; ok {
			queryset = a.applyFilter(queryset, spec, value)
		}
	}
	return queryset
}

// ApplySearch applies search to a queryset using type-safe field references
// Note: For full type safety, SearchFields should use models.FieldRef
func (a *Admin[T]) ApplySearch(ctx context.Context, queryset *models.QuerySetImpl[T], searchTerm string) *models.QuerySetImpl[T] {
	if searchTerm == "" {
		return queryset
	}
	
	// Build OR conditions for all search fields
	// This is a simplified version - full implementation would map
	// admin FieldRef to models FieldRef for type safety
	var combined *models.Condition[T]
	first := true
	
	for _, field := range a.SearchFields() {
		// In a full implementation, this would use models.FieldRef
		// to build type-safe conditions
		condition := a.buildSearchCondition(field.Name(), searchTerm)
		if condition != nil {
			if first {
				combined = condition
				first = false
			} else {
				combined = combined.Or(condition)
			}
		}
	}
	
	if combined != nil {
		queryset = queryset.Filter(combined)
	}
	
	return queryset
}

// applyFilter applies a single filter spec to queryset
func (a *Admin[T]) applyFilter(queryset *models.QuerySetImpl[T], spec FilterSpec[T], value interface{}) *models.QuerySetImpl[T] {
	// Map admin FieldRef to models FieldRef and build condition
	// This is a placeholder - full implementation would use type-safe field refs
	return queryset
}

// buildSearchCondition builds a search condition for a field
// Note: For full type safety, this should accept models.FieldRef
func (a *Admin[T]) buildSearchCondition(fieldName, searchTerm string) *models.Condition[T] {
	// Placeholder - would use models.StringFieldRef.ApplyIContains
	// This requires mapping field names to actual FieldRef instances
	return nil
}

