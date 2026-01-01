package advanced

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockModel for testing
type MockModel struct {
	ID   int64
	Name string
}

func (m *MockModel) GetID() int64 {
	return m.ID
}

// MockAdmin for testing
type MockAdmin struct {
	modelName string
}

func (m *MockAdmin) ModelName() string {
	return m.modelName
}

func TestDefaultHistoryStore(t *testing.T) {
	store := NewDefaultHistoryStore()
	ctx := context.Background()

	t.Run("LogAction", func(t *testing.T) {
		entry := &HistoryEntry{
			ObjectType: "user",
			ObjectID:   1,
			Action:     ActionAdd,
			UserID:     "admin",
			UserName:   "Admin User",
			Changes:    make(map[string]ChangeDetail),
			Message:    "Created user",
			Timestamp:  time.Now(),
		}

		err := store.LogAction(ctx, entry)
		require.NoError(t, err)
		assert.Equal(t, int64(1), entry.ID)
	})

	t.Run("GetObjectHistory", func(t *testing.T) {
		// Log multiple entries
		for i := 1; i <= 3; i++ {
			entry := &HistoryEntry{
				ObjectType: "user",
				ObjectID:   1,
				Action:     ActionChange,
				UserID:     "admin",
				UserName:   "Admin User",
				Timestamp:  time.Now(),
			}
			store.LogAction(ctx, entry)
		}

		history, err := store.GetObjectHistory(ctx, "user", 1)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(history), 3)
	})

	t.Run("GetUserHistory", func(t *testing.T) {
		// Log entries for different objects
		for i := 1; i <= 5; i++ {
			entry := &HistoryEntry{
				ObjectType: "post",
				ObjectID:   int64(i),
				Action:     ActionAdd,
				UserID:     "user1",
				UserName:   "Test User",
				Timestamp:  time.Now(),
			}
			store.LogAction(ctx, entry)
		}

		history, err := store.GetUserHistory(ctx, "user1", 3)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(history), 3)
		
		// Verify order (most recent first)
		if len(history) > 1 {
			assert.True(t, history[0].Timestamp.After(history[1].Timestamp) ||
				history[0].Timestamp.Equal(history[1].Timestamp))
		}
	})

	t.Run("MultipleObjectTypes", func(t *testing.T) {
		// Log entries for different object types
		store.LogAction(ctx, &HistoryEntry{
			ObjectType: "article",
			ObjectID:   1,
			Action:     ActionAdd,
			Timestamp:  time.Now(),
		})
		
		store.LogAction(ctx, &HistoryEntry{
			ObjectType: "comment",
			ObjectID:   1,
			Action:     ActionAdd,
			Timestamp:  time.Now(),
		})

		articleHistory, err := store.GetObjectHistory(ctx, "article", 1)
		require.NoError(t, err)
		
		commentHistory, err := store.GetObjectHistory(ctx, "comment", 1)
		require.NoError(t, err)

		// Should be separate
		assert.NotEqual(t, len(articleHistory), 0)
		assert.NotEqual(t, len(commentHistory), 0)
	})
}

func TestHistoryEntry(t *testing.T) {
	t.Run("ActionFlags", func(t *testing.T) {
		assert.Equal(t, ActionFlag("ADD"), ActionAdd)
		assert.Equal(t, ActionFlag("CHANGE"), ActionChange)
		assert.Equal(t, ActionFlag("DELETE"), ActionDelete)
		assert.Equal(t, ActionFlag("VIEW"), ActionView)
	})

	t.Run("ChangeDetail", func(t *testing.T) {
		change := ChangeDetail{
			Field:    "username",
			OldValue: "old_user",
			NewValue: "new_user",
		}

		assert.Equal(t, "username", change.Field)
		assert.Equal(t, "old_user", change.OldValue)
		assert.Equal(t, "new_user", change.NewValue)
	})
}

func TestHistoryManager(t *testing.T) {
	store := NewDefaultHistoryStore()
	admin := &MockAdmin{modelName: "TestModel"}
	
	// Create a generic history manager (we can't use the actual generic type in tests easily)
	// So we test the store directly
	ctx := context.Background()

	t.Run("LogChangeWithDetails", func(t *testing.T) {
		changes := map[string]ChangeDetail{
			"title": {
				Field:    "title",
				OldValue: "Old Title",
				NewValue: "New Title",
			},
			"status": {
				Field:    "status",
				OldValue: "draft",
				NewValue: "published",
			},
		}

		entry := &HistoryEntry{
			ObjectType: admin.ModelName(),
			ObjectID:   1,
			Action:     ActionChange,
			UserID:     "admin",
			Changes:    changes,
			Timestamp:  time.Now(),
		}

		err := store.LogAction(ctx, entry)
		require.NoError(t, err)

		// Retrieve and verify
		history, err := store.GetObjectHistory(ctx, admin.ModelName(), 1)
		require.NoError(t, err)
		require.NotEmpty(t, history)

		lastEntry := history[len(history)-1]
		assert.Len(t, lastEntry.Changes, 2)
		assert.Contains(t, lastEntry.Changes, "title")
		assert.Contains(t, lastEntry.Changes, "status")
	})

	t.Run("LogDifferentActions", func(t *testing.T) {
		actions := []ActionFlag{ActionAdd, ActionChange, ActionDelete, ActionView}
		
		for i, action := range actions {
			entry := &HistoryEntry{
				ObjectType: "TestObject",
				ObjectID:   int64(i + 1),
				Action:     action,
				UserID:     "admin",
				Timestamp:  time.Now(),
			}
			err := store.LogAction(ctx, entry)
			require.NoError(t, err)
		}

		// Verify all actions were logged
		for i := range actions {
			history, err := store.GetObjectHistory(ctx, "TestObject", int64(i+1))
			require.NoError(t, err)
			require.NotEmpty(t, history)
		}
	})
}

func BenchmarkHistoryStore(b *testing.B) {
	store := NewDefaultHistoryStore()
	ctx := context.Background()

	b.Run("LogAction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			entry := &HistoryEntry{
				ObjectType: "benchmark",
				ObjectID:   int64(i),
				Action:     ActionChange,
				UserID:     "admin",
				Timestamp:  time.Now(),
			}
			store.LogAction(ctx, entry)
		}
	})

	// Prepopulate for read benchmarks
	for i := 0; i < 1000; i++ {
		store.LogAction(ctx, &HistoryEntry{
			ObjectType: "benchmark",
			ObjectID:   int64(i % 100),
			Action:     ActionChange,
			UserID:     "admin",
			Timestamp:  time.Now(),
		})
	}

	b.Run("GetObjectHistory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			store.GetObjectHistory(ctx, "benchmark", int64(i%100))
		}
	})

	b.Run("GetUserHistory", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			store.GetUserHistory(ctx, "admin", 10)
		}
	})
}
