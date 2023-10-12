package wikipedia

import (
	"crypto/sha256"
	"errors"
	"time"
)

// Find and delete string s in string slice
func FindAndDel(arr []string, s string) []string {
	index := 0
	for i, v := range arr {
		if v == s {
			index = i
			break
		}
	}
	return append(arr[:index], arr[index+1:]...)
}

func MakeWikiCache(expiration time.Duration, maxMemory int) *Cache {
	if expiration != 0 {
		expiration = (12 * time.Hour)
	}

	if maxMemory != 0 {
		maxMemory = 500
	}

	c := &Cache{
		Memory:         map[string]RequestResult{},
		MaxMemory:      maxMemory,
		Expiration:     expiration,
		HashedKeyQueue: make([]string, 0, maxMemory),
		CreatedTime:    map[string]time.Time{},
	}

	return c
}

// Cache to store request result
type Cache struct {
	Memory         map[string]RequestResult // Map store request result
	HashedKeyQueue []string                 // Key queue. Delete the first item if reach max cache
	CreatedTime    map[string]time.Time     // Map store created time
	Expiration     time.Duration            // Cache expiration
	MaxMemory      int                      // Max cache memory
}

// Hash a string into SHA256
func HashCacheKey(s string) string {
	hasher := sha256.New()
	hasher.Write([]byte(s))

	return string(hasher.Sum(nil))
}

// Get Cache current number of cache
func (cache Cache) GetLen() int {
	return len(cache.HashedKeyQueue)
}

// Add result into the Cache
func (cache *Cache) Add(s string, res RequestResult) {
	if len(cache.Memory) >= cache.MaxMemory {
		cache.Pop()
	}

	key := HashCacheKey(s)
	if cache.Memory == nil {
		cache.Memory = map[string]RequestResult{}
		cache.CreatedTime = map[string]time.Time{}
		cache.HashedKeyQueue = make([]string, 0, cache.MaxMemory)
	}
	if _, ok := cache.Memory[key]; !ok {
		cache.Memory[key] = res
		cache.CreatedTime[key] = time.Now()
		cache.HashedKeyQueue = append(cache.HashedKeyQueue, key)
	}
}

func (cache *Cache) Get(s string) (RequestResult, error) {
	key := HashCacheKey(s)
	if value, ok := cache.Memory[key]; ok {
		if time.Since(cache.CreatedTime[key]) <= cache.Expiration {
			cache.HashedKeyQueue = FindAndDel(cache.HashedKeyQueue, key)
			cache.HashedKeyQueue = append(cache.HashedKeyQueue, key)
			return value, nil
		} else {
			cache.HashedKeyQueue = FindAndDel(cache.HashedKeyQueue, key)
			delete(cache.Memory, key)
			return RequestResult{}, errors.New("the data is outdated")
		}
	}
	return RequestResult{}, errors.New("cache key not exist")
}

// Delete the first key in the Cache
func (cache *Cache) Pop() {
	if len(cache.HashedKeyQueue) == 0 {
		return
	}
	delete(cache.Memory, cache.HashedKeyQueue[0])
	cache.HashedKeyQueue = cache.HashedKeyQueue[1:]
}

// Clear the whole Cache
func (cache *Cache) Clear() {
	*cache = Cache{}
	// This line to avoid declare but not used error
	_ = cache
}
