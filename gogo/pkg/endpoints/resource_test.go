package endpoints

import (
	"context"
	"testing"
)

// Mock repository for testing
type MockRepository[T any, Q any] struct {
	items []*T
	nextID int
}

func NewMockRepository[T any, Q any]() *MockRepository[T, Q] {
	return &MockRepository[T, Q]{
		items:  make([]*T, 0),
		nextID: 1,
	}
}

func (r *MockRepository[T, Q]) Query() Q {
	var zero Q
	return zero
}

func (r *MockRepository[T, Q]) GetByID(ctx context.Context, id interface{}) (*T, error) {
	// Simple mock - in real test, would find by ID
	if len(r.items) > 0 {
		return r.items[0], nil
	}
	return nil, ErrNotFound
}

func (r *MockRepository[T, Q]) All(ctx context.Context, query Q) ([]*T, error) {
	return r.items, nil
}

func (r *MockRepository[T, Q]) Count(ctx context.Context, query Q) (int, error) {
	return len(r.items), nil
}

func (r *MockRepository[T, Q]) Create(ctx context.Context, data *T) (*T, error) {
	r.items = append(r.items, data)
	return data, nil
}

func (r *MockRepository[T, Q]) Update(ctx context.Context, id interface{}, data *T) (*T, error) {
	if len(r.items) > 0 {
		r.items[0] = data
		return data, nil
	}
	return nil, ErrNotFound
}

func (r *MockRepository[T, Q]) Delete(ctx context.Context, id interface{}) error {
	if len(r.items) > 0 {
		r.items = r.items[1:]
		return nil
	}
	return ErrNotFound
}

// Test model
type TestUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestBaseResource_Index(t *testing.T) {
	repo := NewMockRepository[*TestUser, interface{}]()
	resource := NewResource[*TestUser, interface{}](repo)
	
	// Add test data
	user1 := &TestUser{ID: 1, Name: "John", Email: "john@example.com"}
	user2 := &TestUser{ID: 2, Name: "Jane", Email: "jane@example.com"}
	repo.items = []*TestUser{user1, user2}
	
	// Test would use actual Fiber context
	// For now, just verify structure
	if resource.Repo == nil {
		t.Error("Repository should not be nil")
	}
}

func TestBaseResource_Create(t *testing.T) {
	repo := NewMockRepository[*TestUser, interface{}]()
	resource := NewResource[*TestUser, interface{}](repo)
	
	// Test would create a user
	// For now, just verify structure
	if resource.Serializer == nil {
		t.Error("Serializer should not be nil")
	}
}

