package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/forgego/forge/pkg/models"
)

// QueryParamFilter provides filtering based on query parameters
// Supports operators: eq, ne, gt, gte, lt, lte, contains, in, isnull
// Format: ?field__operator=value or ?field=value (defaults to eq)
type QueryParamFilter[T any] struct {
	BuildCondition func(fieldName, operator, value string) models.Condition
}

// NewQueryParamFilter creates a new QueryParamFilter
// buildCondition is a function that builds a Condition from field name, operator, and value
func NewQueryParamFilter[T any](buildCondition func(fieldName, operator, value string) models.Condition) *QueryParamFilter[T] {
	return &QueryParamFilter[T]{
		BuildCondition: buildCondition,
	}
}

// Apply applies query parameter filters to the queryset
func (f *QueryParamFilter[T]) Apply(c *fiber.Ctx, qs models.QuerySet[T]) models.QuerySet[T] {
	// Get all query parameters
	query := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		query[string(key)] = string(value)
	})
	
	for key, value := range query {
		// Skip pagination params
		if key == "page" || key == "page_size" || key == "ordering" {
			continue
		}
		
		// Parse field and operator (e.g., "email__contains" -> field="email", op="contains")
		parts := strings.Split(key, "__")
		fieldName := parts[0]
		operator := "eq" // default
		
		if len(parts) > 1 {
			operator = parts[1]
		}
		
		// Build condition using the provided function
		if f.BuildCondition != nil {
			condition := f.BuildCondition(fieldName, operator, value)
			if condition != nil {
				qs = qs.Filter(condition)
			}
		}
	}
	
	return qs
}

// OrderingFilter provides ordering based on query parameter
// Format: ?ordering=field or ?ordering=-field (descending)
type OrderingFilter[T any] struct {
	AllowedFields []string
}

// NewOrderingFilter creates a new OrderingFilter
func NewOrderingFilter[T any](allowedFields ...string) *OrderingFilter[T] {
	return &OrderingFilter[T]{
		AllowedFields: allowedFields,
	}
}

// Apply applies ordering to the queryset
func (f *OrderingFilter[T]) Apply(c *fiber.Ctx, qs models.QuerySet[T]) models.QuerySet[T] {
	ordering := c.Query("ordering")
	if ordering == "" {
		return qs
	}
	
	fields := strings.Split(ordering, ",")
	var orders []models.OrderBy
	
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		
		// Check if descending
		descending := false
		if strings.HasPrefix(field, "-") {
			descending = true
			field = field[1:]
		}
		
		// Check if field is allowed
		if len(f.AllowedFields) > 0 {
			allowed := false
			for _, allowedField := range f.AllowedFields {
				if allowedField == field {
					allowed = true
					break
				}
			}
			if !allowed {
				continue
			}
		}
		
		direction := models.OrderAsc
		if descending {
			direction = models.OrderDesc
		}
		orders = append(orders, models.OrderBy{
			Field:     field,
			Direction: direction,
		})
	}
	
	if len(orders) > 0 {
		qs = qs.OrderBy(orders...)
	}
	
	return qs
}

