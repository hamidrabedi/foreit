package admin

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// QueryBuilder builds Ent queries from request parameters
type QueryBuilder struct {
	modelMeta *ModelMeta
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(meta *ModelMeta) *QueryBuilder {
	return &QueryBuilder{
		modelMeta: meta,
	}
}

// QueryParams represents parsed query parameters
type QueryParams struct {
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Filters   map[string]interface{} // Can be simple values or FilterSpec maps
	Search    string
	SearchFields []string // Fields to search in
}

// FilterSpec represents a filter with operator
type FilterSpec struct {
	Field    string
	Operator string // eq, ne, gt, gte, lt, lte, contains, in, between, isnull
	Value    interface{}
}

// ParseQueryParams parses query parameters from a request
func (qb *QueryBuilder) ParseQueryParams(params map[string]string) *QueryParams {
	qp := &QueryParams{
		Page:      1,
		PageSize:  20,
		SortBy:    "id",
		SortOrder: "asc",
		Filters:    make(map[string]interface{}),
	}
	
	// Parse pagination
	if pageStr, ok := params["page"]; ok {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			qp.Page = page
		}
	}
	
	if pageSizeStr, ok := params["page_size"]; ok {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			if pageSize > 100 {
				pageSize = 100 // Max page size
			}
			qp.PageSize = pageSize
		}
	}
	
	// Parse sorting
	if sortBy, ok := params["sort_by"]; ok {
		if qb.isFieldSortable(sortBy) {
			qp.SortBy = sortBy
		}
	}
	
	if sortOrder, ok := params["sort_order"]; ok {
		if sortOrder == "desc" || sortOrder == "asc" {
			qp.SortOrder = sortOrder
		}
	}
	
	// Parse filters with support for operators (e.g., filter_field__gte, filter_field__contains)
	for key, value := range params {
		// Skip pagination and sorting params
		if key == "page" || key == "page_size" || key == "sort_by" || key == "sort_order" || key == "search" {
			continue
		}
		
		// Check if this is a filter
		if strings.HasPrefix(key, "filter_") {
			fieldPart := strings.TrimPrefix(key, "filter_")
			
			// Parse field name and operator (e.g., "created_at__gte" -> field="created_at", op="gte")
			parts := strings.Split(fieldPart, "__")
			fieldName := parts[0]
			operator := "eq" // default operator
			
			if len(parts) > 1 {
				operator = parts[1]
			}
			
			if qb.isFieldFilterable(fieldName) {
				// Store filter with operator information
				filterKey := fieldName
				if operator != "eq" {
					filterKey = fmt.Sprintf("%s__%s", fieldName, operator)
				}
				qp.Filters[filterKey] = map[string]interface{}{
					"field":    fieldName,
					"operator": operator,
					"value":    value,
				}
			}
		} else if qb.isFieldFilterable(key) {
			// Simple filter without operator (defaults to eq)
			qp.Filters[key] = map[string]interface{}{
				"field":    key,
				"operator": "eq",
				"value":    value,
			}
		}
	}
	
	// Parse search
	if search, ok := params["search"]; ok {
		qp.Search = search
	}
	
	// Parse search fields (comma-separated)
	if searchFields, ok := params["search_fields"]; ok {
		fields := strings.Split(searchFields, ",")
		qp.SearchFields = make([]string, 0, len(fields))
		for _, field := range fields {
			field = strings.TrimSpace(field)
			if field != "" {
				qp.SearchFields = append(qp.SearchFields, field)
			}
		}
	}
	
	return qp
}

// isFieldFilterable checks if a field is filterable
func (qb *QueryBuilder) isFieldFilterable(fieldName string) bool {
	// Check if field is in filterable fields list
	if len(qb.modelMeta.Options.FilterableFields) > 0 {
		for _, f := range qb.modelMeta.Options.FilterableFields {
			if f == fieldName {
				return true
			}
		}
		return false
	}
	
	// Default: check field metadata
	for _, field := range qb.modelMeta.Fields {
		if field.Name == fieldName && field.Filterable {
			return true
		}
	}
	
	return false
}

// isFieldSortable checks if a field is sortable
func (qb *QueryBuilder) isFieldSortable(fieldName string) bool {
	// Check if field is in sortable fields list
	if len(qb.modelMeta.Options.SortableFields) > 0 {
		for _, f := range qb.modelMeta.Options.SortableFields {
			if f == fieldName {
				return true
			}
		}
		return false
	}
	
	// Default: check field metadata
	for _, field := range qb.modelMeta.Fields {
		if field.Name == fieldName && field.Sortable {
			return true
		}
	}
	
	return false
}

// BuildEntQuery builds an Ent query from query parameters
// This is a placeholder that will be implemented with actual Ent client
func (qb *QueryBuilder) BuildEntQuery(client interface{}, params *QueryParams) (interface{}, error) {
	// This will be implemented to build actual Ent queries
	// For now, return the params for the handler to use
	return params, nil
}

// ApplyFilters applies filters to a query
func (qb *QueryBuilder) ApplyFilters(query interface{}, filters map[string]interface{}) (interface{}, error) {
	// This will apply filters to the Ent query
	// Implementation depends on Ent's query API
	return query, nil
}

// ApplySorting applies sorting to a query
func (qb *QueryBuilder) ApplySorting(query interface{}, sortBy string, sortOrder string) (interface{}, error) {
	// This will apply sorting to the Ent query
	return query, nil
}

// ApplyPagination applies pagination to a query
func (qb *QueryBuilder) ApplyPagination(query interface{}, page int, pageSize int) (interface{}, error) {
	// This will apply pagination to the Ent query
	return query, nil
}

// ApplySearch applies search to a query
func (qb *QueryBuilder) ApplySearch(query interface{}, search string, searchFields []string) (interface{}, error) {
	// This will apply search to the Ent query
	// Search across multiple fields
	return query, nil
}

// ConvertValue converts a string value to the appropriate type for a field
func (qb *QueryBuilder) ConvertValue(fieldName string, value string) (interface{}, error) {
	// Find the field
	var fieldMeta *FieldMeta
	for i := range qb.modelMeta.Fields {
		if qb.modelMeta.Fields[i].Name == fieldName {
			fieldMeta = &qb.modelMeta.Fields[i]
			break
		}
	}
	
	if fieldMeta == nil {
		return value, nil // Return as string if field not found
	}
	
	// Convert based on field type
	switch fieldMeta.Type {
	case FieldTypeNumber:
		// Try int first, then float
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal, nil
		}
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal, nil
		}
		return nil, fmt.Errorf("invalid number: %s", value)
		
	case FieldTypeBoolean:
		return strings.ToLower(value) == "true" || value == "1", nil
		
	case FieldTypeDate, FieldTypeDateTime, FieldTypeTime:
		// Try common date formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, value); err == nil {
				return t, nil
			}
		}
		return nil, fmt.Errorf("invalid date format: %s", value)
		
	default:
		return value, nil
	}
}

// PaginationResult represents pagination information
type PaginationResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// CalculatePagination calculates pagination info
func CalculatePagination(page, pageSize int, total int64) PaginationResult {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

