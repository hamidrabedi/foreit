package cache

import (
	"context"
	"sync"
	"time"
)

// Store represents a cache store
type Store interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (interface{}, error)
	
	// Set stores a value in cache
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	
	// Delete deletes a value from cache
	Delete(ctx context.Context, key string) error
	
	// Clear clears all values from cache
	Clear(ctx context.Context) error
	
	// Has checks if a key exists
	Has(ctx context.Context, key string) (bool, error)
}

// MemoryStore is an in-memory cache store
type MemoryStore struct {
	data  map[string]*cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
	tags      []string
}

// NewMemoryStore creates a new in-memory cache store
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		data: make(map[string]*cacheItem),
	}
	
	// Start cleanup goroutine
	go store.cleanup()
	
	return store
}

// Get retrieves a value
func (s *MemoryStore) Get(ctx context.Context, key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	item, ok := s.data[key]
	if !ok {
		return nil, ErrNotFound
	}
	
	// Check expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(s.data, key)
		return nil, ErrNotFound
	}
	
	return item.value, nil
}

// Set stores a value
func (s *MemoryStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	item := &cacheItem{
		value: value,
	}
	
	if ttl > 0 {
		item.expiresAt = time.Now().Add(ttl)
	}
	
	s.data[key] = item
	return nil
}

// Delete deletes a value
func (s *MemoryStore) Delete(ctx context.Context, key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	delete(s.data, key)
	return nil
}

// Clear clears all values
func (s *MemoryStore) Clear(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.data = make(map[string]*cacheItem)
	return nil
}

// Has checks if a key exists
func (s *MemoryStore) Has(ctx context.Context, key string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	item, ok := s.data[key]
	if !ok {
		return false, nil
	}
	
	// Check expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return false, nil
	}
	
	return true, nil
}

// cleanup periodically removes expired items
func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for key, item := range s.data {
			if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
				delete(s.data, key)
			}
		}
		s.mutex.Unlock()
	}
}

// Errors
var (
	ErrNotFound = &Error{Code: "not_found", Message: "Cache key not found"}
)

// Error represents a cache error
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// Default store
var defaultStore Store = NewMemoryStore()

// SetDefaultStore sets the default cache store
func SetDefaultStore(store Store) {
	defaultStore = store
}

// Get retrieves a value from the default store
func Get(ctx context.Context, key string) (interface{}, error) {
	return defaultStore.Get(ctx, key)
}

// Set stores a value in the default store
func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return defaultStore.Set(ctx, key, value, ttl)
}

// Delete deletes a value from the default store
func Delete(ctx context.Context, key string) error {
	return defaultStore.Delete(ctx, key)
}

// Has checks if a key exists in the default store
func Has(ctx context.Context, key string) (bool, error) {
	return defaultStore.Has(ctx, key)
}

// Remember gets a value or computes and stores it
func Remember(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache
	if value, err := Get(ctx, key); err == nil {
		return value, nil
	}
	
	// Compute value
	value, err := fn()
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	Set(ctx, key, value, ttl)
	
	return value, nil
}

