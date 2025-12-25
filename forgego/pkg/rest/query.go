package rest

import (
	"strconv"
	"strings"
	
	"github.com/gofiber/fiber/v2"
)

func ApplyQuery[Q any](query Q, c *fiber.Ctx) Q {
	queryParams := make(map[string][]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = []string{string(value)}
	})
	
	for key, values := range queryParams {
		if len(values) == 0 {
			continue
		}
		
		value := values[0]
		
		if key == "page" || key == "page_size" || key == "sort_by" || key == "sort_order" {
			continue
		}
		
		if strings.Contains(key, "__") {
			parts := strings.Split(key, "__")
			field := parts[0]
			operator := parts[1]
			query = applyOperatorFilter(query, field, operator, value)
		} else {
			query = applyEqualityFilter(query, key, value)
		}
	}
	
	return query
}

func applyEqualityFilter[Q any](query Q, field, value string) Q {
	return query
}

func applyOperatorFilter[Q any](query Q, field, operator, value string) Q {
	switch operator {
	case "eq":
		return applyEqualityFilter(query, field, value)
	case "ne", "gt", "gte", "lt", "lte", "contains", "icontains", "startswith", "endswith", "isnull":
		return query
	case "in":
		values := strings.Split(value, ",")
		return applyInFilter(query, field, values)
	default:
		return query
	}
}

func applyInFilter[Q any](query Q, field string, values []string) Q {
	return query
}

// ParsePagination parses pagination parameters
func ParsePagination(c *fiber.Ctx) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page", "1"))
	pageSize, _ = strconv.Atoi(c.Query("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	return page, pageSize
}

// ParseSorting parses sorting parameters
func ParseSorting(c *fiber.Ctx) (sortBy, sortOrder string) {
	sortBy = c.Query("sort_by", "")
	sortOrder = c.Query("sort_order", "asc")
	
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	
	return sortBy, sortOrder
}

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

