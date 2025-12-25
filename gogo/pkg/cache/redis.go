package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements cache store using Redis
type RedisStore struct {
	client *redis.Client
	ctx    context.Context // Kept for backward compatibility, but methods use context parameter
}

// NewRedisStore creates a new Redis cache store
func NewRedisStore(addr string, password string, db int) (Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	ctx := context.Background()
	
	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	
	return &RedisStore{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache
func (s *RedisStore) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	
	// Try to unmarshal as JSON
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// If not JSON, return as string
		return val, nil
	}
	
	return result, nil
}

// Set stores a value in cache
func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data []byte
	var err error
	
	// Try to marshal as JSON
	if data, err = json.Marshal(value); err != nil {
		// If not JSON-serializable, convert to string
		data = []byte(fmt.Sprintf("%v", value))
	}
	
	return s.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes a value from cache
func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// Clear clears all cache
func (s *RedisStore) Clear(ctx context.Context) error {
	return s.client.FlushDB(ctx).Err()
}

// Has checks if a key exists
func (s *RedisStore) Has(ctx context.Context, key string) (bool, error) {
	count, err := s.client.Exists(ctx, key).Result()
	return count > 0, err
}

// GetClient returns the underlying Redis client
func (s *RedisStore) GetClient() *redis.Client {
	return s.client
}

// Close closes the Redis connection
func (s *RedisStore) Close() error {
	return s.client.Close()
}
