package endpoints

import (
	"strconv"
	"strings"
	
	"github.com/gofiber/fiber/v2"
)

func ApplyQuery[Q any](query Q, c *fiber.Ctx) Q {
	queryParams := c.Queries()
	
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

