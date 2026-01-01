package http

import (
	"net/http"

	"github.com/forgego/forge/admin/dashboard"
)

// HandleDashboard handles the admin dashboard
func (h *CoreHandler) HandleDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Create default widgets
		widgets := []dashboard.Widget{
			dashboard.NewStatsWidget("Total Users", "users", "id", "count"),
			dashboard.NewStatsWidget("Total Orders", "orders", "id", "count"),
			dashboard.NewChartWidget("Sales Over Time", "orders", "line", "total"),
			dashboard.NewRecentActivityWidget("Recent Activity", "orders"),
		}

		// Render dashboard
		dashboardData, err := dashboard.RenderDashboard(ctx, widgets)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render template
		if h.renderer != nil {
			data := map[string]interface{}{
				"Title":   dashboardData.Title,
				"Widgets": dashboardData.Widgets,
			}
			h.renderer.Render(w, "dashboard", data)
		} else {
			http.Error(w, "Template renderer not configured", http.StatusInternalServerError)
		}
	}
}
