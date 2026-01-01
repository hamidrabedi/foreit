package views

import (
	"context"
	"fmt"
	"net/http"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
)

// applyGrouping applies grouping to the list view if configured
func (lv *ListView[T]) applyGrouping(ctx context.Context, qs orm.QuerySet[T], config *admin.Config[T]) ([]GroupedRowData, error) {
	// Check if grouping is enabled
	// This would be configured in Config, but for now we'll check URL param
	// In a full implementation, Config would have a GroupBy field
	
	// For now, return empty - grouping would require:
	// 1. Config option for GroupBy field
	// 2. Query modification to group by field
	// 3. Aggregation functions
	// 4. Row reorganization
	
	return []GroupedRowData{}, nil
}

// calculateAggregates calculates aggregate values for a queryset
func (lv *ListView[T]) calculateAggregates(ctx context.Context, qs orm.QuerySet[T], config *admin.Config[T]) ([]AggregateData, error) {
	// This would calculate aggregates like sum, count, avg, min, max
	// For now, return empty
	// Full implementation would:
	// 1. Check Config for aggregate fields
	// 2. Use ORM aggregation functions
	// 3. Return aggregate results
	
	return []AggregateData{}, nil
}

// groupRows groups rows by a field value
func (lv *ListView[T]) groupRows(rows []ListRowData, groupByField string, objects []*T) []GroupedRowData {
	if groupByField == "" {
		return []GroupedRowData{}
	}

	// Group rows by field value
	groups := make(map[interface{}][]ListRowData)
	groupLabels := make(map[interface{}]string)

	for i, row := range rows {
		if i >= len(objects) {
			break
		}

		obj := objects[i]
		groupValue := lv.getFieldValue(obj, groupByField)
		groupLabel := fmt.Sprintf("%v", groupValue)

		if groupValue == nil {
			groupValue = "(None)"
			groupLabel = "(None)"
		}

		groups[groupValue] = append(groups[groupValue], row)
		groupLabels[groupValue] = groupLabel
	}

	// Convert to GroupedRowData
	result := make([]GroupedRowData, 0, len(groups))
	for groupValue, groupRows := range groups {
		result = append(result, GroupedRowData{
			GroupValue: groupValue,
			GroupLabel: groupLabels[groupValue],
			Rows:       groupRows,
			IsGroup:    true,
		})
	}

	return result
}

// calculateSubtotals calculates subtotals for grouped rows
func (lv *ListView[T]) calculateSubtotals(groupedRows []GroupedRowData, aggregateFields []string) {
	for i := range groupedRows {
		group := &groupedRows[i]
		group.Subtotal = make(map[string]AggregateValue)

		// Calculate subtotals for each aggregate field
		for _, field := range aggregateFields {
			// This would calculate sum, count, avg, etc.
			// For now, just count
			group.Subtotal[field] = AggregateValue{
				Type:  "count",
				Value: len(group.Rows),
			}
		}
	}
}

// getGroupByField gets the field to group by from config or request
func (lv *ListView[T]) getGroupByField(r *http.Request, config *admin.Config[T]) string {
	// Check URL parameter first
	if groupBy := r.URL.Query().Get("group_by"); groupBy != "" {
		return groupBy
	}

	// Check config (would need to be added to Config)
	// if config != nil && config.GroupBy != "" {
	//     return config.GroupBy
	// }

	return ""
}
