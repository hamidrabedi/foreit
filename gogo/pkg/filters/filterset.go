package filters

import (
	"context"
	"github.com/gogo/pkg/models"
)

// FilterSet provides django-filters-like functionality (type-safe!)
type FilterSet[T models.Model] struct {
	manager models.Manager[T]
	filters map[string]*Filter[T]
}

// Filter represents a type-safe filter field
type Filter[T models.Model] struct {
	field    *models.FieldRef[T]
	lookups  []Lookup
	required bool
	label    string
	helpText string
}

// Lookup represents a filter lookup type
type Lookup string

const (
	LookupExact      Lookup = "exact"
	LookupIExact     Lookup = "iexact"
	LookupContains   Lookup = "contains"
	LookupIContains  Lookup = "icontains"
	LookupIn         Lookup = "in"
	LookupGt         Lookup = "gt"
	LookupGte        Lookup = "gte"
	LookupLt         Lookup = "lt"
	LookupLte        Lookup = "lte"
	LookupStartswith Lookup = "startswith"
	LookupIStartswith Lookup = "istartswith"
	LookupEndswith   Lookup = "endswith"
	LookupIEndswith  Lookup = "iendswith"
	LookupRange      Lookup = "range"
	LookupDate       Lookup = "date"
	LookupYear       Lookup = "year"
	LookupMonth      Lookup = "month"
	LookupDay        Lookup = "day"
	LookupIsNull     Lookup = "isnull"
	LookupRegex      Lookup = "regex"
	LookupIRegex     Lookup = "iregex"
)

// NewFilterSet creates a new type-safe filter set
func NewFilterSet[T models.Model](manager models.Manager[T]) *FilterSet[T] {
	return &FilterSet[T]{
		manager: manager,
		filters: make(map[string]*Filter[T]),
	}
}

// AddFilter adds a type-safe filter
func (fs *FilterSet[T]) AddFilter(name string, field *models.FieldRef[T], lookups ...Lookup) *FilterSet[T] {
	fs.filters[name] = &Filter[T]{
		field:   field,
		lookups: lookups,
	}
	return fs
}

// Filter filters queryset based on request parameters
func (fs *FilterSet[T]) Filter(ctx context.Context, params map[string]interface{}) models.QuerySet[T] {
	qs := fs.manager.Filter(ctx)
	
	for name, value := range params {
		if filter, ok := fs.filters[name]; ok {
			qs = fs.applyFilter(qs, filter, value)
		}
	}
	
	return qs
}

// applyFilter applies a filter to queryset
func (fs *FilterSet[T]) applyFilter(qs models.QuerySet[T], filter *Filter[T], value interface{}) models.QuerySet[T] {
	// Determine lookup from value type or explicit lookup
	lookup := LookupExact
	
	// If value is a map, extract lookup and value
	if valueMap, ok := value.(map[string]interface{}); ok {
		if l, ok := valueMap["lookup"].(string); ok {
			lookup = Lookup(l)
		}
		if v, ok := valueMap["value"]; ok {
			value = v
		}
	}
	
	// Apply filter based on lookup type
	switch lookup {
	case LookupExact:
		return qs.Filter(models.Q{filter.field.Name(): value})
	case LookupIExact:
		return qs.Filter(models.Q{filter.field.Name() + "__iexact": value})
	case LookupContains:
		if strField, ok := filter.field.(*models.FieldRef[string]); ok {
			return qs.Filter(strField.Contains(value.(string)))
		}
	case LookupIContains:
		if strField, ok := filter.field.(*models.FieldRef[string]); ok {
			return qs.Filter(strField.IContains(value.(string)))
		}
	case LookupIn:
		return qs.Filter(models.Q{filter.field.Name() + "__in": value})
	case LookupGt:
		return qs.Filter(models.Q{filter.field.Name() + "__gt": value})
	case LookupGte:
		return qs.Filter(models.Q{filter.field.Name() + "__gte": value})
	case LookupLt:
		return qs.Filter(models.Q{filter.field.Name() + "__lt": value})
	case LookupLte:
		return qs.Filter(models.Q{filter.field.Name() + "__lte": value})
	case LookupStartswith:
		if strField, ok := filter.field.(*models.FieldRef[string]); ok {
			return qs.Filter(strField.StartsWith(value.(string)))
		}
	case LookupEndswith:
		if strField, ok := filter.field.(*models.FieldRef[string]); ok {
			return qs.Filter(strField.EndsWith(value.(string)))
		}
	case LookupRange:
		if rangeVal, ok := value.([]interface{}); ok && len(rangeVal) == 2 {
			return qs.Filter(models.Q{
				filter.field.Name() + "__gte": rangeVal[0],
				filter.field.Name() + "__lte": rangeVal[1],
			})
		}
	case LookupIsNull:
		return qs.Filter(filter.field.IsNull())
	}
	
	return qs
}

// GetFilters returns all registered filters
func (fs *FilterSet[T]) GetFilters() map[string]*Filter[T] {
	return fs.filters
}

// Example usage:
// type User struct { ... }
// var UserEmail = models.NewFieldRef[string]("email")
// var UserAge = models.NewFieldRef[int]("age")
//
// filterSet := filters.NewFilterSet[*User](userManager).
//     AddFilter("email", UserEmail, LookupExact, LookupIContains).
//     AddFilter("age", UserAge, LookupExact, LookupGt, LookupLt, LookupRange)
//
// qs := filterSet.Filter(ctx, map[string]interface{}{
//     "email__icontains": "@example.com",
//     "age__gte": 18,
// })
//
// users, _ := qs.All(ctx)  // Returns []*User, type-safe!

