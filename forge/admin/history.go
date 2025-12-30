package admin

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// LogEntry represents a change log entry
type LogEntry struct {
	ID          int64
	ActionTime  time.Time
	UserID      int64
	User        interface{} // User object
	ContentType string      // Model name
	ObjectID    int64
	ObjectRepr  string // String representation of object
	ActionFlag  ActionFlag
	ChangeMessage string
}

// ActionFlag represents the type of action
type ActionFlag int

const (
	ActionAddition ActionFlag = iota + 1
	ActionChange
	ActionDeletion
)

// String returns the string representation of the action flag
func (a ActionFlag) String() string {
	switch a {
	case ActionAddition:
		return "Addition"
	case ActionChange:
		return "Change"
	case ActionDeletion:
		return "Deletion"
	default:
		return "Unknown"
	}
}

// HistoryManager manages change history
type HistoryManager interface {
	LogAction(ctx context.Context, entry *LogEntry) error
	GetHistory(ctx context.Context, contentType string, objectID int64) ([]*LogEntry, error)
	GetHistoryForUser(ctx context.Context, userID int64) ([]*LogEntry, error)
}

// DefaultHistoryManager is a default implementation
type DefaultHistoryManager struct {
	// In a real implementation, this would have a database connection
	// For now, it's a placeholder
}

// LogAction logs an action
func (m *DefaultHistoryManager) LogAction(ctx context.Context, entry *LogEntry) error {
	// Default implementation - would save to database
	// This is a placeholder
	return nil
}

// GetHistory retrieves history for an object
func (m *DefaultHistoryManager) GetHistory(ctx context.Context, contentType string, objectID int64) ([]*LogEntry, error) {
	// Default implementation - would query database
	// This is a placeholder
	return []*LogEntry{}, nil
}

// GetHistoryForUser retrieves history for a user
func (m *DefaultHistoryManager) GetHistoryForUser(ctx context.Context, userID int64) ([]*LogEntry, error) {
	// Default implementation - would query database
	// This is a placeholder
	return []*LogEntry{}, nil
}

// LogChange logs a change to an object
func LogChange[T any](ctx context.Context, admin *Admin[T], obj *T, action ActionFlag, user interface{}, message string) error {
	// Get history manager from config or use default
	var historyManager HistoryManager
	if admin.Config() != nil {
		// Would get from config if available
		historyManager = &DefaultHistoryManager{}
	} else {
		historyManager = &DefaultHistoryManager{}
	}

	// Get object ID
	objID := getObjectID(obj)
	if objID == 0 {
		return fmt.Errorf("object does not have an ID field")
	}

	// Get user ID
	userID := getUserID(user)

	// Create log entry
	entry := &LogEntry{
		ActionTime:     time.Now(),
		UserID:         userID,
		User:           user,
		ContentType:    admin.ModelName(),
		ObjectID:       objID,
		ObjectRepr:     fmt.Sprintf("%v", obj),
		ActionFlag:     action,
		ChangeMessage:  message,
	}

	return historyManager.LogAction(ctx, entry)
}

// GetObjectHistory retrieves change history for an object
func GetObjectHistory[T any](ctx context.Context, admin *Admin[T], objID int64) ([]*LogEntry, error) {
	var historyManager HistoryManager = &DefaultHistoryManager{}
	
	return historyManager.GetHistory(ctx, admin.ModelName(), objID)
}

// getObjectID extracts ID from an object using reflection
func getObjectID(obj interface{}) int64 {
	// Use reflection to get ID field
	// This is a simplified version
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanInterface() {
		if id, ok := idField.Interface().(int64); ok {
			return id
		}
	}

	return 0
}

// getUserID extracts user ID from user object
func getUserID(user interface{}) int64 {
	if user == nil {
		return 0
	}

	val := reflect.ValueOf(user)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanInterface() {
		if id, ok := idField.Interface().(int64); ok {
			return id
		}
	}

	return 0
}
