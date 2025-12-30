package testing

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// TestAdminHelper provides helpers for testing admin functionality
type TestAdminHelper[T any] struct {
	admin *admin.Admin[T]
}

// NewTestAdminHelper creates a new test admin helper
func NewTestAdminHelper[T any](
	schemaInstance schema.Schema,
	manager *orm.Manager[T],
	config *admin.Config[T],
) (*TestAdminHelper[T], error) {
	admin, err := admin.Register(schemaInstance, manager, config)
	if err != nil {
		return nil, err
	}

	return &TestAdminHelper[T]{
		admin: admin,
	}, nil
}

// Admin returns the admin instance
func (h *TestAdminHelper[T]) Admin() *admin.Admin[T] {
	return h.admin
}

// MakeRequest creates a test HTTP request
func (h *TestAdminHelper[T]) MakeRequest(method, path string, body []byte) *http.Request {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	return httptest.NewRequest(method, path, bodyReader)
}

// MakeGetRequest creates a GET request
func (h *TestAdminHelper[T]) MakeGetRequest(path string) *http.Request {
	return h.MakeRequest("GET", path, nil)
}

// MakePostRequest creates a POST request
func (h *TestAdminHelper[T]) MakePostRequest(path string, body []byte) *http.Request {
	return h.MakeRequest("POST", path, body)
}

// ExecuteRequest executes a request and returns the response
func (h *TestAdminHelper[T]) ExecuteRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

// AssertResponseStatus asserts the response status
func AssertResponseStatus(t testing.T, rr *httptest.ResponseRecorder, expected int) {
	if rr.Code != expected {
		t.Errorf("expected status %d, got %d", expected, rr.Code)
	}
}

// AssertResponseBodyContains asserts the response body contains text
func AssertResponseBodyContains(t testing.T, rr *httptest.ResponseRecorder, text string) {
	if !strings.Contains(rr.Body.String(), text) {
		t.Errorf("expected response body to contain %q, got %q", text, rr.Body.String())
	}
}
