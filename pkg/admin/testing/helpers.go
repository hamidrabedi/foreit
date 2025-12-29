package testing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

adminv2 "github.com/forgego/forge/pkg/admin"
adminhttp "github.com/forgego/forge/pkg/admin/http"
	httplib "github.com/forgego/forge/pkg/http"
	"github.com/forgego/forge/pkg/query"
)

// TestAdminClient is a test client for admin HTTP endpoints
type TestAdminClient struct {
	Handler  http.Handler
	Registry *adminv2.Registry
	BaseURL  string
}

// NewTestAdminClient creates a new admin test client
func NewTestAdminClient(registry *adminv2.Registry) *TestAdminClient {
	router := httplib.NewRouter()
	adminRouter := adminhttp.NewRouter(registry)
	adminRouter.RegisterRoutes(router, "/admin")

	return &TestAdminClient{
		Handler:  router,
		Registry: registry,
		BaseURL:  "/admin",
	}
}

// Get performs a GET request to an admin endpoint
func (c *TestAdminClient) Get(path string) *AdminResponse {
	return c.Request("GET", path, nil, nil)
}

// Post performs a POST request to an admin endpoint
func (c *TestAdminClient) Post(path string, formData map[string]string) *AdminResponse {
	return c.Request("POST", path, formData, nil)
}

// Request performs an HTTP request
func (c *TestAdminClient) Request(method, path string, formData map[string]string, headers map[string]string) *AdminResponse {
	req := httptest.NewRequest(method, c.BaseURL+path, nil)

	// Set form data if provided
	if formData != nil {
		req.ParseForm()
		for k, v := range formData {
			req.Form.Set(k, v)
		}
	}

	// Set headers
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	recorder := httptest.NewRecorder()
	c.Handler.ServeHTTP(recorder, req)

	return &AdminResponse{
		StatusCode: recorder.Code,
		Body:       recorder.Body.Bytes(),
		Headers:    recorder.Header(),
	}
}

// AdminResponse represents an admin HTTP response
type AdminResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// JSON parses the response body as JSON
func (r *AdminResponse) JSON() map[string]interface{} {
	var data map[string]interface{}
	json.Unmarshal(r.Body, &data)
	return data
}

// Status returns the status code
func (r *AdminResponse) Status() int {
	return r.StatusCode
}

// BodyString returns the response body as string
func (r *AdminResponse) BodyString() string {
	return string(r.Body)
}

// TestManager is a mock manager for testing
type TestManager[T any] struct {
	objects []*T
	nextID  int64
}

// NewTestManager creates a new test manager
func NewTestManager[T any]() *TestManager[T] {
	return &TestManager[T]{
		objects: make([]*T, 0),
		nextID:  1,
	}
}

// Get retrieves an object by ID
func (m *TestManager[T]) Get(ctx context.Context, id int64) (*T, error) {
	for _, obj := range m.objects {
		if getID(obj) == id {
			return obj, nil
		}
	}
	return nil, fmt.Errorf("object not found")
}

// All retrieves all objects
func (m *TestManager[T]) All(ctx context.Context) ([]*T, error) {
	return m.objects, nil
}

// Create creates a new object
func (m *TestManager[T]) Create(ctx context.Context, obj *T) error {
	setID(obj, m.nextID)
	m.nextID++
	m.objects = append(m.objects, obj)
	return nil
}

// Update updates an object
func (m *TestManager[T]) Update(ctx context.Context, obj *T) error {
	id := getID(obj)
	for i, o := range m.objects {
		if getID(o) == id {
			m.objects[i] = obj
			return nil
		}
	}
	return fmt.Errorf("object not found")
}

// Delete deletes an object
func (m *TestManager[T]) Delete(ctx context.Context, obj *T) error {
	id := getID(obj)
	for i, o := range m.objects {
		if getID(o) == id {
			m.objects = append(m.objects[:i], m.objects[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("object not found")
}

// Filter returns a queryset (simplified for testing)
func (m *TestManager[T]) Filter(expr ...query.QueryExpr) query.QuerySet[T] {
	// Return a simple queryset that wraps the manager
	// This is a simplified version for testing
	return &TestQuerySet[T]{manager: m}
}

// TestQuerySet is a test queryset
type TestQuerySet[T any] struct {
	manager *TestManager[T]
}

func (qs *TestQuerySet[T]) All(ctx context.Context) ([]*T, error) {
	return qs.manager.All(ctx)
}

func (qs *TestQuerySet[T]) Count(ctx context.Context) (int64, error) {
	objects, _ := qs.manager.All(ctx)
	return int64(len(objects)), nil
}

func (qs *TestQuerySet[T]) Filter(expr query.QueryExpr) query.QuerySet[T] {
	return qs // Simplified - just return self
}

func (qs *TestQuerySet[T]) Exclude(expr query.QueryExpr) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) OrderBy(fields ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Reverse() query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Limit(n int) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Offset(n int) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Distinct() query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Select(fields ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) SelectRelated(relations ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) PrefetchRelated(relations ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Only(fields ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Defer(fields ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Aggregate(aggregates ...query.Aggregate) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Annotate(annotations ...query.AnnotationExpr) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Values(fields ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) ValuesList(fields ...string) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Get(ctx context.Context) (*T, error) {
	objects, err := qs.manager.All(ctx)
	if err != nil || len(objects) == 0 {
		return nil, fmt.Errorf("object not found")
	}
	return objects[0], nil
}

func (qs *TestQuerySet[T]) First(ctx context.Context) (*T, error) {
	return qs.Get(ctx)
}

func (qs *TestQuerySet[T]) Last(ctx context.Context) (*T, error) {
	objects, err := qs.manager.All(ctx)
	if err != nil || len(objects) == 0 {
		return nil, fmt.Errorf("object not found")
	}
	return objects[len(objects)-1], nil
}

func (qs *TestQuerySet[T]) Exists(ctx context.Context) (bool, error) {
	objects, err := qs.manager.All(ctx)
	return len(objects) > 0, err
}

func (qs *TestQuerySet[T]) Update(ctx context.Context, fields map[string]interface{}) (int64, error) {
	return 0, fmt.Errorf("not implemented in test")
}

func (qs *TestQuerySet[T]) BulkUpdate(ctx context.Context, updates []map[string]interface{}) error {
	return fmt.Errorf("not implemented in test")
}

func (qs *TestQuerySet[T]) BulkCreate(ctx context.Context, instances []*T) error {
	return fmt.Errorf("not implemented in test")
}

func (qs *TestQuerySet[T]) Delete(ctx context.Context) (int64, error) {
	return 0, fmt.Errorf("not implemented in test")
}

func (qs *TestQuerySet[T]) Union(other query.QuerySet[T]) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Intersection(other query.QuerySet[T]) query.QuerySet[T] {
	return qs
}

func (qs *TestQuerySet[T]) Difference(other query.QuerySet[T]) query.QuerySet[T] {
	return qs
}

// Helper functions for ID manipulation (using reflection)
func getID(obj interface{}) int64 {
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

func setID(obj interface{}, id int64) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		idField.Set(reflect.ValueOf(id))
	}
}

// WithTestDB creates a test database connection
func WithTestDB(t *testing.T, fn func(*sql.DB)) {
	// Use in-memory SQLite for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	fn(db)
}

// CreateTestAdmin creates a test admin instance
func CreateTestAdmin[T any](model T, manager *query.Manager[T], config *adminv2.Config[T]) *adminv2.Admin[T] {
	return adminv2.Register(model, manager, config)
}
