package core

// DashboardConfig defines widget layout and configuration for the admin dashboard.
type DashboardConfig struct {
	Widgets []WidgetMeta            `json:"widgets,omitempty"`
	Layout  []DashboardLayoutWidget `json:"layout,omitempty"`
}

// DashboardLayoutWidget defines a widget layout hint.
type DashboardLayoutWidget struct {
	ID   string `json:"id"`
	Size string `json:"size,omitempty"` // sm, md, lg, full
}

