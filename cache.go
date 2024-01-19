package livesessiondbcache

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

// Generic cache that can store any value given some string key
// it's safe to access concurrently and guarantees that data associated with a key
// will be requested at most once
type Cache[T any] struct {
	// sync map that stores data retrieved using dataRetriever
	m sync.Map

	// because sync.Map doesn't provide lazy initialization we need another
	// structure to guarantee that we'll only request user data once per user
	s singleflight.Group

	dataRetriever DataRetriever[T]
}

func MakeCache[T any](dataRetriever DataRetriever[T]) *Cache[T] {
	var cache Cache[T]
	cache.dataRetriever = dataRetriever
	return &cache
}

func (c *Cache[T]) Get(key string) (*T, error) {
	val, ok := c.m.Load(key)
	if ok {
		as_t := val.(*T)
		return as_t, nil
	}

	v, err, _ := c.s.Do(key, func() (any, error) {
		val, err := c.dataRetriever.Get(key)
		if err == nil {
			c.m.Store(key, val)
		}
		return val, err
	})

	if err != nil {
		return nil, err
	}

	as_t := v.(*T)
	return as_t, nil
}
