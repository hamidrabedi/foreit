package cache

import (
	"context"
	"sync"
)

// TagStore extends Store with tag support
type TagStore interface {
	Store
	
	// Tag tags a key with tags
	Tag(ctx context.Context, key string, tags ...string) error
	
	// Invalidate invalidates all keys with a tag
	Invalidate(ctx context.Context, tag string) error
	
	// InvalidateAll invalidates all keys with any of the tags
	InvalidateAll(ctx context.Context, tags ...string) error
}

// TaggedMemoryStore extends MemoryStore with tags
type TaggedMemoryStore struct {
	*MemoryStore
	tagIndex map[string]map[string]bool // tag -> keys
	tagMutex sync.RWMutex
}

// NewTaggedMemoryStore creates a new tagged memory store
func NewTaggedMemoryStore() *TaggedMemoryStore {
	return &TaggedMemoryStore{
		MemoryStore: NewMemoryStore(),
		tagIndex:    make(map[string]map[string]bool),
	}
}

// Tag tags a key with tags
func (s *TaggedMemoryStore) Tag(ctx context.Context, key string, tags ...string) error {
	s.tagMutex.Lock()
	defer s.tagMutex.Unlock()
	
	// Get item to store tags
	s.mutex.Lock()
	item, ok := s.data[key]
	if !ok {
		item = &cacheItem{}
		s.data[key] = item
	}
	item.tags = tags
	s.mutex.Unlock()
	
	// Update tag index
	for _, tag := range tags {
		if s.tagIndex[tag] == nil {
			s.tagIndex[tag] = make(map[string]bool)
		}
		s.tagIndex[tag][key] = true
	}
	
	return nil
}

// Set overrides Set to handle tags
func (s *TaggedMemoryStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if err := s.MemoryStore.Set(ctx, key, value, ttl); err != nil {
		return err
	}
	
	// Preserve existing tags
	s.mutex.RLock()
	item, ok := s.data[key]
	s.mutex.RUnlock()
	
	if ok && len(item.tags) > 0 {
		return s.Tag(ctx, key, item.tags...)
	}
	
	return nil
}

// Invalidate invalidates all keys with a tag
func (s *TaggedMemoryStore) Invalidate(ctx context.Context, tag string) error {
	s.tagMutex.RLock()
	keys, ok := s.tagIndex[tag]
	if !ok {
		s.tagMutex.RUnlock()
		return nil
	}
	
	keysCopy := make([]string, 0, len(keys))
	for key := range keys {
		keysCopy = append(keysCopy, key)
	}
	s.tagMutex.RUnlock()
	
	// Delete all keys
	for _, key := range keysCopy {
		s.Delete(ctx, key)
	}
	
	// Remove from tag index
	s.tagMutex.Lock()
	delete(s.tagIndex, tag)
	s.tagMutex.Unlock()
	
	return nil
}

// InvalidateAll invalidates all keys with any of the tags
func (s *TaggedMemoryStore) InvalidateAll(ctx context.Context, tags ...string) error {
	for _, tag := range tags {
		if err := s.Invalidate(ctx, tag); err != nil {
			return err
		}
	}
	return nil
}

// TagSet tags a key when setting
func TagSet(ctx context.Context, key string, value interface{}, ttl time.Duration, tags ...string) error {
	store, ok := defaultStore.(TagStore)
	if !ok {
		// Fallback to regular set
		return Set(ctx, key, value, ttl)
	}
	
	if err := store.Set(ctx, key, value, ttl); err != nil {
		return err
	}
	
	return store.Tag(ctx, key, tags...)
}

// TagInvalidate invalidates a tag
func TagInvalidate(ctx context.Context, tag string) error {
	store, ok := defaultStore.(TagStore)
	if !ok {
		return nil // No-op if store doesn't support tags
	}
	return store.Invalidate(ctx, tag)
}

