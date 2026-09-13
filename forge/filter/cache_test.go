package filter

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterCache_Concurrency(t *testing.T) {
	cache := NewFilterCache(1 * time.Minute)
	var wg sync.WaitGroup
	numGoroutines := 50
	iterations := 200

	keys := []string{"key-0", "key-1", "key-2", "key-3", "key-4"}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(gID int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if j%10 == 0 {
					time.Sleep(2 * time.Millisecond)
				}
				key := keys[(gID+j)%len(keys)]
				switch (gID + j) % 6 {
				case 0:
					cache.SetParsedTree(key, &FilterNode{Field: key}, 1*time.Millisecond)
				case 1:
					_, _ = cache.GetParsedTree(key)
				case 2:
					cache.SetCompiledSQL(key, "SELECT 1", 1*time.Millisecond)
				case 3:
					_, _ = cache.GetCompiledSQL(key)
				case 4:
					cache.SetMetadata(key, map[string]interface{}{"k": "v"}, 1*time.Millisecond)
				case 5:
					_, _ = cache.GetMetadata(key)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestFilterCache_Expiry(t *testing.T) {
	cache := NewFilterCache(1 * time.Minute)

	cache.SetParsedTree("tree-key", &FilterNode{Field: "f"}, 10*time.Millisecond)
	cache.SetCompiledSQL("sql-key", "SELECT 1", 10*time.Millisecond)
	cache.SetMetadata("meta-key", map[string]interface{}{"k": "v"}, 10*time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	tree, ok := cache.GetParsedTree("tree-key")
	assert.False(t, ok)
	assert.Nil(t, tree)

	sql, ok := cache.GetCompiledSQL("sql-key")
	assert.False(t, ok)
	assert.Empty(t, sql)

	meta, ok := cache.GetMetadata("meta-key")
	assert.False(t, ok)
	assert.Nil(t, meta)
}

func TestFilterCache_Close(t *testing.T) {
	before := runtime.NumGoroutine()
	cache := NewFilterCache(1 * time.Minute)

	err := cache.Close()
	require.NoError(t, err)

	deadline := time.Now().Add(1 * time.Second)
	stopped := false
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before {
			stopped = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	require.True(t, stopped, "cleanup goroutine did not stop within 1s")
}

func TestFilterCache_CloseTwice(t *testing.T) {
	cache := NewFilterCache(1 * time.Minute)
	err := cache.Close()
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		err = cache.Close()
		assert.NoError(t, err)
	})
}
