package livesessiondbcache

// Any object capable of retrieving data given it's key
type DataSource[T any] interface {
	Get(key string) (*T, error)
}
