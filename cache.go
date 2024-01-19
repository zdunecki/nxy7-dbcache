package livesessiondbcache

import "sync"

type Cache[T any] struct {
	data          sync.Map
	dataRetriever DataRetriever[T]
}
