package advanced

import (
	"context"
	"fmt"

	"github.com/forgego/forge/admin"
)

// AdminAction is a type alias for admin.Action to avoid generic resolution issues
type AdminAction[T any] = admin.Action[T]

// ActionManager manages bulk actions for admin
type ActionManager[T any] struct {
	admin   *admin.Admin[T]
	actions map[string]AdminAction[T]
}

// NewActionManager creates a new action manager
func NewActionManager[T any](admin *admin.Admin[T]) *ActionManager[T] {
	return &ActionManager[T]{
		admin:   admin,
		actions: make(map[string]AdminAction[T]),
	}
}

// RegisterAction registers a bulk action
func (am *ActionManager[T]) RegisterAction(action AdminAction[T]) {
	am.actions[action.Name] = action
}

// GetAction gets an action by name
func (am *ActionManager[T]) GetAction(name string) (AdminAction[T], bool) {
	action, ok := am.actions[name]
	return action, ok
}

// GetAllActions returns all registered actions
func (am *ActionManager[T]) GetAllActions() map[string]AdminAction[T] {
	return am.actions
}

// ExecuteAction executes a bulk action on selected instances
func (am *ActionManager[T]) ExecuteAction(ctx context.Context, actionName string, instances []*T, user interface{}) error {
	action, ok := am.GetAction(actionName)
	if !ok {
		return fmt.Errorf("action %s not found", actionName)
	}

	// Check permissions
	if len(action.Permissions) > 0 {
		config := am.admin.Config()
		if config != nil && config.PermissionChecker != nil {
			for _, perm := range action.Permissions {
				if !config.PermissionChecker.HasPermission(ctx, user, perm) {
					return fmt.Errorf("permission denied for action %s", actionName)
				}
			}
		}
	}

	// Execute action handler
	return action.Handler(ctx, instances)
}

// BuiltinActions provides built-in actions
type BuiltinActions[T any] struct {
	admin *admin.Admin[T]
}

// NewBuiltinActions creates built-in actions
func NewBuiltinActions[T any](admin *admin.Admin[T]) *BuiltinActions[T] {
	return &BuiltinActions[T]{
		admin: admin,
	}
}

// DeleteSelected creates a delete selected action
func (ba *BuiltinActions[T]) DeleteSelected() AdminAction[T] {
	return admin.NewAction[T](
		"delete_selected",
		"Delete selected items",
		func(ctx context.Context, instances []*T) error {
			for _, instance := range instances {
				if err := ba.admin.DeleteModel(ctx, instance); err != nil {
					return fmt.Errorf("failed to delete instance: %w", err)
				}
			}
			return nil
		},
	).WithDescription("Delete the selected items permanently")
}

// ExportSelected creates an export selected action
func (ba *BuiltinActions[T]) ExportSelected(format string) AdminAction[T] {
	return admin.NewAction[T](
		"export_selected",
		"Export selected items",
		func(ctx context.Context, instances []*T) error {
			// Convert instances to exportable format
			exportData := make([]map[string]interface{}, 0, len(instances))
			for _, inst := range instances {
				// Convert instance to map (simplified - would use proper serialization)
				instanceMap := make(map[string]interface{})
				// In a real implementation, would serialize fields properly
				_ = inst // Mark instance as used
				exportData = append(exportData, instanceMap)
			}
			
			// Export functionality is handled by HTTP handler
			// This action just marks items for export
			_ = exportData
			_ = format
			return nil
		},
	).WithDescription(fmt.Sprintf("Export selected items as %s", format))
}
