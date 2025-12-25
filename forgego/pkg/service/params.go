package service

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
)

// ParseListParams converts HTTP query params to typed ListParams.
// Uses typed Condition AST to prevent SQL injection - all predicates are validated.
func ParseListParams(c *fiber.Ctx) ListParams {
	params := ListParams{
		Page:     1,
		PageSize: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			if pageSize > 200 {
				pageSize = 200
			}
			params.PageSize = pageSize
		}
	}

	predicates := []models.Condition{}
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		keyStr := string(key)
		valueStr := string(value)

		if keyStr == "page" || keyStr == "page_size" || keyStr == "sort_by" || keyStr == "sort_order" || keyStr == "search" {
			return
		}

		if strings.HasPrefix(keyStr, "filter_") {
			field := strings.TrimPrefix(keyStr, "filter_")
			parts := strings.Split(field, "__")
			fieldName := parts[0]
			operator := "eq"

			if len(parts) > 1 {
				operator = parts[1]
			}

			var condition models.Condition
			switch operator {
			case "eq":
				condition = newStringCondition(fieldName, "=", valueStr)
			case "ne":
				condition = newStringCondition(fieldName, "!=", valueStr)
			case "gt":
				condition = newStringCondition(fieldName, ">", valueStr)
			case "gte":
				condition = newStringCondition(fieldName, ">=", valueStr)
			case "lt":
				condition = newStringCondition(fieldName, "<", valueStr)
			case "lte":
				condition = newStringCondition(fieldName, "<=", valueStr)
			case "contains", "icontains":
				condition = newStringCondition(fieldName, "ILIKE", "%"+valueStr+"%")
			case "in":
				values := strings.Split(valueStr, ",")
				condition = newInCondition(fieldName, values)
			case "isnull":
				isNull := valueStr == "true" || valueStr == "1"
				condition = newIsNullCondition(fieldName, isNull)
			}

			if condition != nil {
				predicates = append(predicates, condition)
			}
		}
	})
	params.Predicates = predicates

	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	if searchFields := c.Query("search_fields"); searchFields != "" {
		fields := strings.Split(searchFields, ",")
		params.SearchFields = make([]string, 0, len(fields))
		for _, field := range fields {
			field = strings.TrimSpace(field)
			if field != "" {
				params.SearchFields = append(params.SearchFields, field)
			}
		}
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		sortOrder := c.Query("sort_order")
		if sortOrder == "" {
			sortOrder = "ASC"
		}
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "ASC"
		}
		params.Sort = []OrderBy{{Field: sortBy, Direction: sortOrder}}
	}

	if cursor := c.Query("cursor"); cursor != "" {
		params.Cursor = &cursor
		if limitStr := c.Query("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				params.Limit = limit
			}
		}
	}

	return params
}

