package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/forgego/forge/schema"
	forgehttp "github.com/forgego/forge/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type nonSerializableTestModel struct {
	schema.BaseSchema
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func (nonSerializableTestModel) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64Field("id"),
		schema.StringField("username"),
		schema.StringField("password", schema.Serialize(false)),
	}
}

func TestNonSerializableFields_ReturnsOnlySerializeFalseFields(t *testing.T) {
	got := NonSerializableFields(&nonSerializableTestModel{})
	require.Equal(t, []string{"password"}, got)
}

type excludeResponseSerializer struct {
	*BaseSerializer
}

func newExcludeResponseSerializer() Serializer {
	return &excludeResponseSerializer{
		BaseSerializer: NewBaseSerializer(make(map[string]interface{})),
	}
}

func (s *excludeResponseSerializer) Fields() []string {
	return []string{"id", "username", "password"}
}

type excludeResponseQuerySet struct {
	mgr    *excludeResponseManager
	limit  int
	offset int
}

func (qs *excludeResponseQuerySet) Count(ctx context.Context) (int64, error) {
	return int64(len(qs.mgr.items)), nil
}

func (qs *excludeResponseQuerySet) Offset(n int) interface{} {
	qs.offset = n
	return qs
}

func (qs *excludeResponseQuerySet) Limit(n int) interface{} {
	qs.limit = n
	return qs
}

func (qs *excludeResponseQuerySet) All(ctx context.Context) ([]*nonSerializableTestModel, error) {
	start := qs.offset
	if start > len(qs.mgr.items) {
		start = len(qs.mgr.items)
	}
	end := start + qs.limit
	if qs.limit <= 0 || end > len(qs.mgr.items) {
		end = len(qs.mgr.items)
	}
	out := make([]*nonSerializableTestModel, 0, end-start)
	for i := start; i < end; i++ {
		cp := *qs.mgr.items[i]
		out = append(out, &cp)
	}
	return out, nil
}

type excludeResponseManager struct {
	mu     sync.Mutex
	items  []*nonSerializableTestModel
	nextID int64
}

func (m *excludeResponseManager) QuerySet() interface{} {
	return &excludeResponseQuerySet{mgr: m}
}

func (m *excludeResponseManager) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m *excludeResponseManager) All(ctx context.Context) ([]*nonSerializableTestModel, error) {
	out := make([]*nonSerializableTestModel, len(m.items))
	for i := range m.items {
		cp := *m.items[i]
		out[i] = &cp
	}
	return out, nil
}

func (m *excludeResponseManager) Get(ctx context.Context, id int64) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, item := range m.items {
		if item.ID == id {
			cp := *item
			return &cp, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *excludeResponseManager) Create(ctx context.Context, model interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := model.(*nonSerializableTestModel)
	if !ok {
		return errors.New("unexpected model type")
	}
	m.nextID++
	item.ID = m.nextID
	cp := *item
	m.items = append(m.items, &cp)
	return nil
}

func (m *excludeResponseManager) Update(ctx context.Context, model interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := model.(*nonSerializableTestModel)
	if !ok {
		return errors.New("unexpected model type")
	}
	for i, existing := range m.items {
		if existing.ID == item.ID {
			cp := *item
			m.items[i] = &cp
			return nil
		}
	}
	return errors.New("not found")
}

func (m *excludeResponseManager) Delete(ctx context.Context, model interface{}) error {
	return nil
}

func TestBaseViewSet_ExcludeResponseFields_OmitsFieldsInListRetrieveCreateUpdate(t *testing.T) {
	mgr := &excludeResponseManager{
		items: []*nonSerializableTestModel{
			{ID: 1, Username: "alice", Password: "secret1"},
		},
		nextID: 1,
	}
	vs := NewBaseViewSet(
		newExcludeResponseSerializer,
		mgr,
		&nonSerializableTestModel{},
	)
	vs.ExcludeResponseFields = NonSerializableFields(&nonSerializableTestModel{})

	router := NewRouter("/api")
	router.Register("users", vs)
	handler := forgehttp.NewRouter()
	router.RegisterRoutes(handler)

	decodeMap := func(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
		t.Helper()
		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &m))
		return m
	}

	// List
	req := httptest.NewRequest(http.MethodGet, "/api/users/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var listResp struct {
		Results []map[string]interface{} `json:"results"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &listResp))
	require.Len(t, listResp.Results, 1)
	assert.NotContains(t, listResp.Results[0], "password")
	assert.Contains(t, listResp.Results[0], "username")

	// Retrieve
	req = httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	m := decodeMap(t, rec)
	assert.NotContains(t, m, "password")
	assert.Contains(t, m, "username")

	// Create
	req = httptest.NewRequest(http.MethodPost, "/api/users/", bytes.NewBufferString(`{"username":"bob","password":"secret2"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	m = decodeMap(t, rec)
	assert.NotContains(t, m, "password")
	assert.Equal(t, "bob", m["username"])

	// Update
	req = httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBufferString(`{"username":"alice2","password":"newsecret"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	m = decodeMap(t, rec)
	assert.NotContains(t, m, "password")
	assert.Equal(t, "alice2", m["username"])
}
