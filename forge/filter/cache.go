package filter

import (
	"sync"
	"time"
)

// CacheEntry represents a cache entry
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// FilterCache caches filter-related data
type FilterCache struct {
	parsedTrees map[string]*CacheEntry
	compiledSQL map[string]*CacheEntry
	metadata    map[string]*CacheEntry
	mu          sync.RWMutex
	defaultTTL  time.Duration
	stop        chan struct{}
	closeOnce   sync.Once
}

// NewFilterCache creates a new filter cache
func NewFilterCache(defaultTTL time.Duration) *FilterCache {
	cache := &FilterCache{
		parsedTrees: make(map[string]*CacheEntry),
		compiledSQL: make(map[string]*CacheEntry),
		metadata:    make(map[string]*CacheEntry),
		defaultTTL:  defaultTTL,
		stop:        make(chan struct{}),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Close stops the background cleanup goroutine. Safe to call multiple times.
func (c *FilterCache) Close() error {
	c.closeOnce.Do(func() {
		close(c.stop)
	})
	return nil
}

// GetParsedTree gets a cached parsed filter tree
func (c *FilterCache) GetParsedTree(key string) (*FilterNode, bool) {
	c.mu.RLock()
	entry, ok := c.parsedTrees[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		if cur, ok := c.parsedTrees[key]; ok && time.Now().After(cur.ExpiresAt) {
			delete(c.parsedTrees, key)
		}
		c.mu.Unlock()
		return nil, false
	}

	node, ok := entry.Value.(*FilterNode)
	return node, ok
}

// SetParsedTree caches a parsed filter tree
func (c *FilterCache) SetParsedTree(key string, node *FilterNode, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.parsedTrees[key] = &CacheEntry{
		Value:     node,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// GetCompiledSQL gets cached compiled SQL
func (c *FilterCache) GetCompiledSQL(key string) (string, bool) {
	c.mu.RLock()
	entry, ok := c.compiledSQL[key]
	c.mu.RUnlock()
	if !ok {
		return "", false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		if cur, ok := c.compiledSQL[key]; ok && time.Now().After(cur.ExpiresAt) {
			delete(c.compiledSQL, key)
		}
		c.mu.Unlock()
		return "", false
	}

	sql, ok := entry.Value.(string)
	return sql, ok
}

// SetCompiledSQL caches compiled SQL
func (c *FilterCache) SetCompiledSQL(key string, sql string, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.compiledSQL[key] = &CacheEntry{
		Value:     sql,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// GetMetadata gets cached filter metadata
func (c *FilterCache) GetMetadata(key string) (map[string]interface{}, bool) {
	c.mu.RLock()
	entry, ok := c.metadata[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		if cur, ok := c.metadata[key]; ok && time.Now().After(cur.ExpiresAt) {
			delete(c.metadata, key)
		}
		c.mu.Unlock()
		return nil, false
	}

	metadata, ok := entry.Value.(map[string]interface{})
	return metadata, ok
}

// SetMetadata caches filter metadata
func (c *FilterCache) SetMetadata(key string, metadata map[string]interface{}, ttl time.Duration) {
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.metadata[key] = &CacheEntry{
		Value:     metadata,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Clear clears all cache entries
func (c *FilterCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.parsedTrees = make(map[string]*CacheEntry)
	c.compiledSQL = make(map[string]*CacheEntry)
	c.metadata = make(map[string]*CacheEntry)
}

// cleanup periodically removes expired entries
func (c *FilterCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()

			// Clean parsed trees
			for k, entry := range c.parsedTrees {
				if now.After(entry.ExpiresAt) {
					delete(c.parsedTrees, k)
				}
			}

			// Clean compiled SQL
			for k, entry := range c.compiledSQL {
				if now.After(entry.ExpiresAt) {
					delete(c.compiledSQL, k)
				}
			}

			// Clean metadata
			for k, entry := range c.metadata {
				if now.After(entry.ExpiresAt) {
					delete(c.metadata, k)
				}
			}

			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}
