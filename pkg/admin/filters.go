package admin

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/query"
)

// FilterOption represents a filter option
type FilterOption struct {
	Label string
	Value interface{}
	Count int64
}

// FilterField represents a filterable field
type FilterField struct {
	Name    string
	Label   string
	Type    string
	Options []FilterOption
	Active  interface{} // Currently selected value
}

// applyFilters applies list filters to a queryset (reflection-based)
func applyFilters(r *http.Request, model *AdminModel, qs reflect.Value) reflect.Value {
	if len(model.ListFilter) == 0 {
		return qs
	}

	for _, filterField := range model.ListFilter {
		fieldName, ok := filterField.(string)
		if !ok {
			continue
		}

		// Get filter value from query params
		filterValue := httplib.GetQueryString(r, fieldName, "")
		if filterValue == "" {
			continue
		}

		// Apply filter based on field type
		// This is simplified - full implementation would check field type from schema
		filterExpr := query.NewFieldQueryExpr(fieldName, query.OpEquals, filterValue)
		
		// Call Filter method on queryset
		filterMethod := qs.MethodByName("Filter")
		if filterMethod.IsValid() {
			results := filterMethod.Call([]reflect.Value{reflect.ValueOf(filterExpr)})
			if len(results) > 0 {
				qs = results[0]
			}
		}
	}

	return qs
}

// getFilterFields generates filter fields for the list view
func getFilterFields(r *http.Request, model *AdminModel, manager interface{}) []FilterField {
	var filterFields []FilterField

	if len(model.ListFilter) == 0 {
		return filterFields
	}

	ctx := r.Context()
	managerValue := reflect.ValueOf(manager)

	for _, filterField := range model.ListFilter {
		fieldName, ok := filterField.(string)
		if !ok {
			continue
		}

		// Get all distinct values for this field
		// This is simplified - full implementation would query distinct values
		filterField := FilterField{
			Name:  fieldName,
			Label: fieldName,
			Type:  "select", // Default to select
			Active: httplib.GetQueryString(r, fieldName, ""),
		}

		// Try to get distinct values from manager
		allMethod := managerValue.MethodByName("All")
		if allMethod.IsValid() {
			results := allMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if len(results) >= 2 && results[1].IsNil() {
				objects := results[0].Interface()
				options := getDistinctValues(objects, fieldName)
				filterField.Options = options
			}
		}

		filterFields = append(filterFields, filterField)
	}

	return filterFields
}

// getDistinctValues extracts distinct values from objects for a field
func getDistinctValues(objects interface{}, fieldName string) []FilterOption {
	var options []FilterOption
	seen := make(map[interface{}]bool)
	counts := make(map[interface{}]int64)

	objectsValue := reflect.ValueOf(objects)
	if objectsValue.Kind() != reflect.Slice {
		return options
	}

	for i := 0; i < objectsValue.Len(); i++ {
		obj := objectsValue.Index(i)
		if obj.Kind() == reflect.Ptr {
			obj = obj.Elem()
		}

		field := obj.FieldByName(fieldName)
		if !field.IsValid() {
			// Try lowercase
			field = obj.FieldByName(toTitleCase(fieldName))
		}

		if field.IsValid() {
			value := field.Interface()
			if !seen[value] {
				seen[value] = true
				options = append(options, FilterOption{
					Label: fmt.Sprintf("%v", value),
					Value: value,
					Count: 1,
				})
			} else {
				// Increment count
				for i := range options {
					if options[i].Value == value {
						options[i].Count++
						break
					}
				}
			}
			counts[value]++
		}
	}

	return options
}


// parseFilterValue parses a filter value from string
func parseFilterValue(valueStr string, fieldType string) interface{} {
	switch fieldType {
	case "bool":
		if valueStr == "true" || valueStr == "1" {
			return true
		}
		return false
	case "int", "int64":
		if val, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			return val
		}
	case "float", "float64":
		if val, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return val
		}
	}
	return valueStr
}

