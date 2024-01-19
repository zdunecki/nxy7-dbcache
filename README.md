# Take home assignment
This library provides generic implementation of `Cache[T]` which can be accessed concurrently and is guaranteed to make at most
one call to underlying data source for any given key.
To function correctly this cache requires `DataRetriever[T]` that implements `Get(key string) (*T, error)`.
