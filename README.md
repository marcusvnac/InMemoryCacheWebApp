# InMemoryCacheWebApp

This is a Go WebAPI application featuring an in-memory cache implementation. This exercise is part of a code assessment for Minerva.

It implements an HTTP WebAPI application with an in-memory cache that can hold a maximum of 255 keys, satisfying the interface provided below.

```go
go
type (

    Cache interface {

        // Set sets the value to the provided key in the given bucket.
        // Applying any provided options during the operation.
        // An error is returned if operation fails.
        Set(bucket string, key string, value []byte, opts ...Option) error
        
        // Get returns the value associated with the given key in the bucket.
        // Applying any provided options during the operation.
        // An error is returned if operation fails.
        Get(bucket, key string, opts ...Option) ([]byte, error)
        
        // Delete removes the key and value from the bucket.
        // Applying any provided options during the operation.
        // An error is returned if operation fails.
        Delete(bucket, key string, opts ...Option) error
    }



    Options struct {
        ttl time.Duration

        // evictionPolicy controls how keys should be removed from the cache.
        // Options: Oldest, Newest, LRU, MRU
        evictionPolicy string
    }

    Option func(o Options) error
)
```
