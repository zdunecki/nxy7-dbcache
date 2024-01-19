package livesessiondbcache

// Any object capable of retrieving data given it's key
type DataRetriever[T any] interface {
	Get(key string) (*T, error)
}
