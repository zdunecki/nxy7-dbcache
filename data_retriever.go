package livesessiondbcache

type DataRetriever[T any] interface {
	Get(key string) (T, error)
}
