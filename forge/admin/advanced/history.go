package advanced

import (
	"context"
	"time"

	"github.com/forgego/forge/admin"
	adminutils "github.com/forgego/forge/admin/utils"
)

// HistoryManager manages change history for admin models
type HistoryManager[T any] struct {
	admin *admin.Admin[T]
	store HistoryStore
}

// HistoryStore is the interface for storing history entries
type HistoryStore interface {
	LogAction(ctx context.Context, entry *HistoryEntry) error
	GetObjectHistory(ctx context.Context, objectType string, objectID int64) ([]*HistoryEntry, error)
	GetUserHistory(ctx context.Context, userID interface{}, limit int) ([]*HistoryEntry, error)
}

// HistoryEntry represents a history entry
type HistoryEntry struct {
	ID         int64
	ObjectType string
	ObjectID   int64
	Action     ActionFlag
	UserID     interface{}
	UserName   string
	Changes    map[string]ChangeDetail
	Message    string
	Timestamp  time.Time
}

// ActionFlag represents the type of action
type ActionFlag string

const (
	ActionAdd    ActionFlag = "ADD"
	ActionChange ActionFlag = "CHANGE"
	ActionDelete ActionFlag = "DELETE"
	ActionView   ActionFlag = "VIEW"
)

// ChangeDetail represents a field change
type ChangeDetail struct {
	Field    string
	OldValue interface{}
	NewValue interface{}
}

// NewHistoryManager creates a new history manager
func NewHistoryManager[T any](admin *admin.Admin[T], store HistoryStore) *HistoryManager[T] {
	return &HistoryManager[T]{
		admin: admin,
		store: store,
	}
}

// LogChange logs a change to an object
func (hm *HistoryManager[T]) LogChange(ctx context.Context, instance *T, action ActionFlag, user interface{}, changes map[string]ChangeDetail) error {
	// Get object ID using utils
	objectID := adminutils.GetIDFromInstance(instance)

	entry := &HistoryEntry{
		ObjectType: hm.admin.ModelName(),
		ObjectID:   objectID,
		Action:     action,
		UserID:     user,
		Changes:    changes,
		Timestamp:  time.Now(),
	}

	return hm.store.LogAction(ctx, entry)
}

// GetHistory gets history for an object
func (hm *HistoryManager[T]) GetHistory(ctx context.Context, instance *T) ([]*HistoryEntry, error) {
	objectID := adminutils.GetIDFromInstance(instance)
	return hm.store.GetObjectHistory(ctx, hm.admin.ModelName(), objectID)
}

// DefaultHistoryStore is a default in-memory history store
type DefaultHistoryStore struct {
	entries []*HistoryEntry
	nextID  int64
}

// NewDefaultHistoryStore creates a new default history store
func NewDefaultHistoryStore() *DefaultHistoryStore {
	return &DefaultHistoryStore{
		entries: make([]*HistoryEntry, 0),
		nextID:  1,
	}
}

// LogAction logs an action
func (ds *DefaultHistoryStore) LogAction(ctx context.Context, entry *HistoryEntry) error {
	entry.ID = ds.nextID
	ds.nextID++
	ds.entries = append(ds.entries, entry)
	return nil
}

// GetObjectHistory gets history for an object
func (ds *DefaultHistoryStore) GetObjectHistory(ctx context.Context, objectType string, objectID int64) ([]*HistoryEntry, error) {
	result := make([]*HistoryEntry, 0)
	for _, entry := range ds.entries {
		if entry.ObjectType == objectType && entry.ObjectID == objectID {
			result = append(result, entry)
		}
	}
	return result, nil
}

// GetUserHistory gets history for a user
func (ds *DefaultHistoryStore) GetUserHistory(ctx context.Context, userID interface{}, limit int) ([]*HistoryEntry, error) {
	result := make([]*HistoryEntry, 0)
	count := 0
	for i := len(ds.entries) - 1; i >= 0 && count < limit; i-- {
		entry := ds.entries[i]
		if entry.UserID == userID {
			result = append(result, entry)
			count++
		}
	}
	return result, nil
}
