package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	store map[string]cacheEntry
	mutex sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		store: map[string]cacheEntry{},
	}
	go cache.reapLoop(interval)
	return &cache
}

func (cache *Cache) reapLoop(interval time.Duration) {
	timer := time.NewTicker(interval)
	for {
		select {
		case <-timer.C:
			cache.mutex.Lock()
			for key, entry := range cache.store {
				if time.Since(entry.createdAt) > interval {
					delete(cache.store, key)
				}
			}
			cache.mutex.Unlock()
		default:
			continue
		}
	}
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.store[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, ok := cache.store[key]
	if !ok {
		return nil, false
	}
	return entry.val, ok
}
