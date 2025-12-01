package cache

import (
	"sync"
	"time"
)

type TimedCache struct {
	mu          sync.Mutex
	cache       map[string]interface{}
	expirations map[string]time.Time
}

func NewTimedCache() *TimedCache {
	return &TimedCache{
		cache:       make(map[string]interface{}),
		expirations: make(map[string]time.Time),
	}
}

func (tc *TimedCache) Set(key string, value interface{}, duration time.Duration) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cache[key] = value
	expiration := time.Now().Add(duration)
	tc.expirations[key] = expiration

	time.AfterFunc(duration, func() {
		tc.mu.Lock()
		defer tc.mu.Unlock()
		if tc.expirations[key].Before(time.Now()) {
			delete(tc.cache, key)
			delete(tc.expirations, key)
		}
	})
}

func (tc *TimedCache) Get(key string) (interface{}, bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	value, exists := tc.cache[key]
	if !exists || tc.expirations[key].Before(time.Now()) {
		return "", false
	}
	return value, true
}
