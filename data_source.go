package dbcache

// Any object capable of retrieving data given it's key
type DataSource[T any] interface {
	// Retrieves object given some key.
	// This function should return (nil, nil) if request to underlying data store was successful
	// but no data for given key was found.
	Get(key string) (*T, error)
}
