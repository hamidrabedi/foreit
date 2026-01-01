package dashboard

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/forgego/forge/admin"
)

// Dashboard represents the admin dashboard
type Dashboard struct {
	Title   string
	Widgets []Widget
}

// Widget represents a dashboard widget
type Widget interface {
	Render(ctx context.Context) (WidgetData, error)
	Type() string
	Title() string
}

// WidgetData contains data for rendering a widget
type WidgetData struct {
	Type    string
	Title   string
	Content interface{}
	Config  map[string]interface{}
}

// StatsWidget displays statistics
type StatsWidget struct {
	TitleText  string
	ModelName  string
	Field      string
	Function   string // "count", "sum", "avg", "min", "max"
	TimeRange  string // "today", "week", "month", "year", "all"
	Comparison bool   // Show comparison to previous period
}

// NewStatsWidget creates a new stats widget
func NewStatsWidget(title, modelName, field, function string) *StatsWidget {
	return &StatsWidget{
		TitleText: title,
		ModelName: modelName,
		Field:     field,
		Function:  function,
		TimeRange: "all",
	}
}

// Type returns the widget type
func (w *StatsWidget) Type() string {
	return "stats"
}

// Title returns the widget title
func (w *StatsWidget) Title() string {
	return w.TitleText
}

// Render renders the widget
func (w *StatsWidget) Render(ctx context.Context) (WidgetData, error) {
	// Get admin instance from registry
	registry := admin.GetGlobalRegistry()
	adminInstance, err := registry.Get(w.ModelName)
	if err != nil {
		// Return placeholder if model not found (allows dashboard to work even if model not registered)
		return WidgetData{
			Type:  "stats",
			Title: w.TitleText,
			Content: map[string]interface{}{
				"value":      0,
				"label":      w.TitleText,
				"trend":      "N/A",
				"comparison": false,
			},
			Config: map[string]interface{}{
				"model":     w.ModelName,
				"field":     w.Field,
				"function":  w.Function,
				"timeRange": w.TimeRange,
			},
		}, nil
	}

	// Get manager interface
	managerInterface := adminInstance.ManagerInterface()
	if managerInterface == nil {
		return WidgetData{}, fmt.Errorf("admin %s has no manager", w.ModelName)
	}

	// Use reflection to call All() method on manager
	managerValue := reflect.ValueOf(managerInterface)
	allMethod := managerValue.MethodByName("All")
	if !allMethod.IsValid() {
		return WidgetData{}, fmt.Errorf("manager for %s does not have All() method", w.ModelName)
	}

	// Call All(ctx) to get all instances
	results := allMethod.Call([]reflect.Value{reflect.ValueOf(ctx)})
	if len(results) != 2 {
		return WidgetData{}, fmt.Errorf("All() method returned unexpected number of values")
	}

	// Check for error
	if !results[1].IsNil() {
		err := results[1].Interface().(error)
		return WidgetData{}, fmt.Errorf("failed to query %s: %w", w.ModelName, err)
	}

	// Get the slice of instances
	instances := results[0].Interface()
	instancesValue := reflect.ValueOf(instances)

	// Calculate value based on function
	var value interface{}
	switch w.Function {
	case "count":
		value = instancesValue.Len()
	case "sum", "avg", "min", "max":
		// For numeric aggregations, we'd need to iterate and calculate
		// For now, return count as placeholder for these functions
		// Full implementation would require field accessor and type conversion
		value = instancesValue.Len()
	default:
		value = instancesValue.Len()
	}

	// Calculate trend if comparison is enabled
	trend := "N/A"
	if w.Comparison {
		// TODO: Compare with previous period
		// This would require querying with time filters
		trend = "+0%"
	}

	return WidgetData{
		Type:  "stats",
		Title: w.TitleText,
		Content: map[string]interface{}{
			"value":      value,
			"label":      w.TitleText,
			"trend":      trend,
			"comparison": w.Comparison,
		},
		Config: map[string]interface{}{
			"model":     w.ModelName,
			"field":     w.Field,
			"function":  w.Function,
			"timeRange": w.TimeRange,
		},
	}, nil
}

// ChartWidget displays a chart
type ChartWidget struct {
	TitleText string
	ModelName string
	ChartType  string // "line", "bar", "pie"
	Field     string
	GroupBy   string
	TimeRange string
}

// NewChartWidget creates a new chart widget
func NewChartWidget(title, modelName, chartType, field string) *ChartWidget {
	return &ChartWidget{
		TitleText: title,
		ModelName: modelName,
		ChartType: chartType,
		Field:     field,
		TimeRange: "month",
	}
}

// Type returns the widget type
func (w *ChartWidget) Type() string {
	return "chart"
}

// Title returns the widget title
func (w *ChartWidget) Title() string {
	return w.TitleText
}

// Render renders the widget
func (w *ChartWidget) Render(ctx context.Context) (WidgetData, error) {
	// Placeholder chart data
	chartData := map[string]interface{}{
		"labels": []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
		"datasets": []map[string]interface{}{
			{
				"label": w.Title,
				"data":  []int{10, 20, 15, 25, 30, 35},
			},
		},
	}

	return WidgetData{
		Type:  "chart",
		Title: w.TitleText,
		Content: map[string]interface{}{
			"type":  w.ChartType,
			"data":  chartData,
			"config": map[string]interface{}{
				"model":     w.ModelName,
				"field":     w.Field,
				"groupBy":   w.GroupBy,
				"timeRange": w.TimeRange,
			},
		},
		Config: map[string]interface{}{
			"model":     w.ModelName,
			"field":     w.Field,
			"groupBy":   w.GroupBy,
			"timeRange": w.TimeRange,
		},
	}, nil
}

// RecentActivityWidget displays recent activity
type RecentActivityWidget struct {
	TitleText string
	ModelName string
	Limit     int
}

// NewRecentActivityWidget creates a new recent activity widget
func NewRecentActivityWidget(title, modelName string) *RecentActivityWidget {
	return &RecentActivityWidget{
		TitleText: title,
		ModelName: modelName,
		Limit:     10,
	}
}

// Type returns the widget type
func (w *RecentActivityWidget) Type() string {
	return "activity"
}

// Title returns the widget title
func (w *RecentActivityWidget) Title() string {
	return w.TitleText
}

// Render renders the widget
func (w *RecentActivityWidget) Render(ctx context.Context) (WidgetData, error) {
	// Placeholder activity data
	activities := []map[string]interface{}{
		{
			"id":        1,
			"action":    "created",
			"model":     w.ModelName,
			"timestamp": time.Now().Add(-1 * time.Hour),
			"user":      "admin",
		},
		{
			"id":        2,
			"action":    "updated",
			"model":     w.ModelName,
			"timestamp": time.Now().Add(-2 * time.Hour),
			"user":      "admin",
		},
	}

	return WidgetData{
		Type:  "activity",
		Title: w.TitleText,
		Content: map[string]interface{}{
			"activities": activities,
			"limit":      w.Limit,
		},
		Config: map[string]interface{}{
			"model": w.ModelName,
			"limit": w.Limit,
		},
	}, nil
}

// QuickActionsWidget displays quick action buttons
type QuickActionsWidget struct {
	TitleText string
	Actions []QuickAction
}

// QuickAction represents a quick action
type QuickAction struct {
	Label string
	URL   string
	Icon  string
}

// NewQuickActionsWidget creates a new quick actions widget
func NewQuickActionsWidget(title string) *QuickActionsWidget {
	return &QuickActionsWidget{
		TitleText: title,
		Actions: []QuickAction{},
	}
}

// Type returns the widget type
func (w *QuickActionsWidget) Type() string {
	return "quick_actions"
}

// Title returns the widget title
func (w *QuickActionsWidget) Title() string {
	return w.TitleText
}

// Render renders the widget
func (w *QuickActionsWidget) Render(ctx context.Context) (WidgetData, error) {
	return WidgetData{
		Type:  "quick_actions",
		Title: w.TitleText,
		Content: map[string]interface{}{
			"actions": w.Actions,
		},
		Config: map[string]interface{}{},
	}, nil
}

// AddAction adds an action to the widget
func (w *QuickActionsWidget) AddAction(label, url, icon string) {
	w.Actions = append(w.Actions, QuickAction{
		Label: label,
		URL:   url,
		Icon:  icon,
	})
}

// RenderDashboard renders the dashboard
func RenderDashboard(ctx context.Context, widgets []Widget) (*Dashboard, error) {
	dashboard := &Dashboard{
		Title:   "Admin Dashboard",
		Widgets: widgets,
	}

	return dashboard, nil
}
