package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/forgego/forge/admin/advanced"
	httplib "github.com/forgego/forge/server"
)

// HistoryData represents data for the history template
type HistoryData struct {
	ModelName  string
	InstanceID int64
	History    []HistoryItemData
}

// HistoryItemData represents a single history item for display
type HistoryItemData struct {
	Timestamp   string
	UserName    string
	Action      string
	ActionClass string
	Description string
}

// HandleHistory handles history view requests
func (h *CoreHandler) HandleHistory(modelName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get instance ID from URL
		idStr := httplib.GetQueryString(r, "id", "")
		if idStr == "" {
			http.Error(w, "Instance ID required", http.StatusBadRequest)
			return
		}

		instanceID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid instance ID", http.StatusBadRequest)
			return
		}

		// Get history store from handler
		handler, err := GetAdminHandler(modelName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		// For now, return empty history
		// In a full implementation, this would call the history store
		history := make([]*advanced.HistoryEntry, 0)
		
		// If handler implements history interface, get history
		if historyHandler, ok := handler.(interface {
			GetHistory(ctx context.Context, instanceID int64) ([]*advanced.HistoryEntry, error)
		}); ok {
			history, err = historyHandler.GetHistory(ctx, instanceID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get history: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Convert history entries to display format
		historyItems := make([]HistoryItemData, len(history))
		for i, entry := range history {
			historyItems[i] = HistoryItemData{
				Timestamp:   entry.Timestamp.Format("Jan 02, 2006 15:04"),
				UserName:    entry.UserName,
				Action:      string(entry.Action),
				ActionClass: getActionClass(entry.Action),
				Description: formatHistoryDescription(entry),
			}
		}

		// Render history template
		data := HistoryData{
			ModelName:  modelName,
			InstanceID: instanceID,
			History:    historyItems,
		}

		if err := h.renderer.Render(w, "history", map[string]interface{}{
			"ModelName":  data.ModelName,
			"InstanceID": data.InstanceID,
			"History":    data.History,
		}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to render history: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// getActionClass returns the CSS class for an action
func getActionClass(action advanced.ActionFlag) string {
	switch action {
	case advanced.ActionAdd:
		return "success"
	case advanced.ActionChange:
		return "info"
	case advanced.ActionDelete:
		return "danger"
	case advanced.ActionView:
		return "secondary"
	default:
		return "secondary"
	}
}

// formatHistoryDescription formats a history entry description
func formatHistoryDescription(entry *advanced.HistoryEntry) string {
	if entry.Message != "" {
		return entry.Message
	}

	// Generate description from changes
	if len(entry.Changes) == 0 {
		switch entry.Action {
		case advanced.ActionAdd:
			return "Object created"
		case advanced.ActionDelete:
			return "Object deleted"
		case advanced.ActionView:
			return "Object viewed"
		default:
			return "No changes recorded"
		}
	}

	// Format changes
	desc := "Changed: "
	first := true
	for fieldName := range entry.Changes {
		if !first {
			desc += ", "
		}
		desc += fieldName
		first = false
	}

	return desc
}

// adminHandler interface extension for history
type adminHandlerWithHistory interface {
	GetHistory(ctx interface{}, instanceID int64) ([]*advanced.HistoryEntry, error)
}

// LogHistory logs a history entry for an admin action
func (h *CoreHandler) LogHistory(ctx interface{}, modelName string, instanceID int64, action advanced.ActionFlag, userID interface{}, changes map[string]advanced.ChangeDetail) error {
	handler, err := GetAdminHandler(modelName)
	if err != nil {
		return err
	}

	if historyHandler, ok := handler.(adminHandlerWithHistory); ok {
		entry := &advanced.HistoryEntry{
			ObjectType: modelName,
			ObjectID:   instanceID,
			Action:     action,
			UserID:     userID,
			Changes:    changes,
		}
		_ = entry // Use entry with history handler
		_ = historyHandler
		// TODO: Implement actual logging through history handler
	}

	return nil
}
