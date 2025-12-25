package endpoints

import (
	"strconv"
)

// QueryProcessor processes query parameters for filtering, pagination, and sorting
type QueryProcessor struct {
	Filters   map[string]Filter
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Search    string
}

// Filter represents a query filter
type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

// ParseQueryParams parses query parameters from the request
func ParseQueryParams(ctx *Context) *QueryProcessor {
	processor := &QueryProcessor{
		Filters: make(map[string]Filter),
		Page:    1,
		PageSize: 20,
		SortOrder: "asc",
	}
	
	// Parse pagination
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			processor.Page = page
		}
	}
	
	if pageSizeStr := ctx.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			processor.PageSize = pageSize
		}
	}
	
	// Parse sorting
	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		processor.SortBy = sortBy
	}
	
	if sortOrder := ctx.Query("sort_order"); sortOrder != "" {
		if sortOrder == "desc" {
			processor.SortOrder = "desc"
		}
	}
	
	// Parse search
	if search := ctx.Query("search"); search != "" {
		processor.Search = search
	}
	
	// Parse filters (format: filter_field=value or filter_field__operator=value)
	queryParams := ctx.Request.URI().QueryArgs()
	queryParams.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		valueStr := string(value)
		
		// Skip known parameters
		if keyStr == "page" || keyStr == "page_size" || keyStr == "sort_by" || 
		   keyStr == "sort_order" || keyStr == "search" {
			return
		}
		
		// Parse filter syntax: field__operator or just field
		parts := splitFilterKey(keyStr)
		field := parts[0]
		operator := "eq"
		if len(parts) > 1 {
			operator = parts[1]
		}
		
		processor.Filters[field] = Filter{
			Field:    field,
			Operator: operator,
			Value:    valueStr,
		}
	})
	
	return processor
}

// splitFilterKey splits a filter key like "name__contains" into ["name", "contains"]
func splitFilterKey(key string) []string {
	// Handle both "__" and "_" separators
	if idx := len(key); idx > 0 {
		for i := len(key) - 1; i >= 0; i-- {
			if key[i] == '_' && i > 0 && key[i-1] == '_' {
				return []string{key[:i-1], key[i+1:]}
			}
		}
		// Try single underscore
		for i := len(key) - 1; i >= 0; i-- {
			if key[i] == '_' {
				return []string{key[:i], key[i+1:]}
			}
		}
	}
	return []string{key}
}

// PaginationResult contains pagination information
type PaginationResult struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// CalculatePagination calculates pagination metadata
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

