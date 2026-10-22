// Package pokecache caches data from network requests
package pokecache

import (
	"sync"
	"time"
)

type PokeCache struct {
	Mu sync.Mutex
	Entries map[string]CacheEntry
	interval time.Duration
}

type CacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

func NewPokeCache(interval time.Duration) *PokeCache {
	newPokeCache := &PokeCache{
		Mu: sync.Mutex{},
		Entries: map[string]CacheEntry{},
		interval: interval,
	}
	go newPokeCache.reapLoop()
	return newPokeCache
}

func (pc *PokeCache) Add(key string, val []byte) {
	pc.Mu.Lock()
	defer pc.Mu.Unlock()
	newEntry := CacheEntry{
		CreatedAt: time.Now(),
		Val: val,
	}
	pc.Entries[key] = newEntry
} 

func (pc *PokeCache) Get(key string) ([]byte, bool) {
	pc.Mu.Lock()
	defer pc.Mu.Unlock()
	entry, ok := pc.Entries[key]
	if !ok {
		return []byte{}, false
	}
	return entry.Val, true
}

func (pc *PokeCache) reapLoop() {
	ticker := time.NewTicker(pc.interval)
	for range ticker.C {
		var markedForDeletion []string
		for key, entry := range pc.Entries {
			expiryTime := entry.CreatedAt.Add(pc.interval)
			if time.Now().After(expiryTime) {
				markedForDeletion = append(markedForDeletion, key)
			}
		}
		for _, entryToBeDeleted := range markedForDeletion {
			delete(pc.Entries, entryToBeDeleted)
		}
	}
}