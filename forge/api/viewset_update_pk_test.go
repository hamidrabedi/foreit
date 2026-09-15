package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	forgehttp "github.com/forgego/forge/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type updateTestItem struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type updateTestItemSerializer struct {
	*BaseSerializer
}

func newUpdateTestItemSerializer() Serializer {
	return &updateTestItemSerializer{
		BaseSerializer: NewBaseSerializer(make(map[string]interface{})),
	}
}

func (s *updateTestItemSerializer) Fields() []string {
	return []string{"id", "title"}
}

type fakeUpdateManager struct {
	mu    sync.Mutex
	items map[int64]*updateTestItem
}

func (m *fakeUpdateManager) Create(ctx context.Context, item interface{}) error {
	return nil
}

func (m *fakeUpdateManager) Get(ctx context.Context, id int64) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	// Return a copy so in-memory pointer mutation doesn't alias map entry directly
	copy := *item
	return &copy, nil
}

func (m *fakeUpdateManager) Update(ctx context.Context, model interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	item := model.(*updateTestItem)
	m.items[item.ID] = item
	return nil
}

func TestUpdate_IgnoresBodyIDAndDoesNotModifyOtherRow(t *testing.T) {
	mgr := &fakeUpdateManager{
		items: map[int64]*updateTestItem{
			1: {ID: 1, Title: "row 1 original"},
			2: {ID: 2, Title: "row 2 original"},
		},
	}

	vs := NewBaseViewSet(
		newUpdateTestItemSerializer,
		mgr,
		&updateTestItem{},
	)

	router := NewRouter("/api")
	router.Register("items", vs)
	handler := forgehttp.NewRouter()
	router.RegisterRoutes(handler)

	body := `{"id": 2, "title": "row 1 updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/items/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	// Row 2 must remain untouched
	assert.Equal(t, "row 2 original", mgr.items[2].Title, "row 2 must be unchanged")
	assert.Equal(t, int64(2), mgr.items[2].ID, "row 2 ID must be 2")

	// Row 1 must have been updated
	assert.Equal(t, "row 1 updated", mgr.items[1].Title, "row 1 title should be updated")
	assert.Equal(t, int64(1), mgr.items[1].ID, "row 1 ID must remain 1")
}
