package service

import (
	"sync"
	"time"
)

const cityCacheDuration = time.Hour

type cityCache[T any] struct {
	mu      sync.Mutex
	entries map[string]*cacheEntry[T]
	now     func() time.Time
}

type cacheEntry[T any] struct {
	mu       sync.Mutex
	value    T
	storedAt time.Time
	stored   bool
}

func newCityCache[T any]() *cityCache[T] {
	return &cityCache[T]{
		entries: map[string]*cacheEntry[T]{},
		now:     time.Now,
	}
}

func (c *cityCache[T]) get(cityID string, fetch func() (T, error)) (T, time.Time, error) {
	c.mu.Lock()
	entry, found := c.entries[cityID]
	if !found {
		entry = &cacheEntry[T]{}
		c.entries[cityID] = entry
	}
	c.mu.Unlock()

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.stored && c.now().Sub(entry.storedAt) < cityCacheDuration {
		return entry.value, entry.storedAt, nil
	}

	value, err := fetch()
	if err != nil {
		var zero T
		return zero, time.Time{}, err
	}

	entry.value = value
	entry.storedAt = c.now()
	entry.stored = true
	return value, entry.storedAt, nil
}
