package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache interface {
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

type Options struct {
	ttl            time.Duration
	evictionPolicy string
}

type Option func(o *Options) error

type cacheItem struct {
	value      []byte
	expiration time.Time
}

type InMemoryCache struct {
	Data  map[string]map[string]cacheItem
	mutex sync.RWMutex
}

func NewCache() Cache {
	return &InMemoryCache{
		Data: make(map[string]map[string]cacheItem),
	}
}

func (c *InMemoryCache) Set(bucket string, key string, value []byte, opts ...Option) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.Data[bucket]; !exists {
		c.Data[bucket] = make(map[string]cacheItem)
	} else if len(c.Data[bucket]) >= 255 {
		return errors.New("cache is full")
	}

	// process options funcs
	options := Options{}
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return err
		}
	}

	expiration := time.Time{}
	if options.ttl > 0 {
		expiration = time.Now().Add(options.ttl)
	}

	c.Data[bucket][key] = cacheItem{
		value:      value,
		expiration: expiration,
	}

	return nil
}

func (c *InMemoryCache) Get(bucket, key string, opts ...Option) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if bucketData, exists := c.Data[bucket]; exists {
		if item, exists := bucketData[key]; exists {
			if item.expiration.IsZero() || item.expiration.After(time.Now()) {
				return item.value, nil
			}
			delete(bucketData, key)
		}
	}

	return nil, errors.New("key not found")
}

func (c *InMemoryCache) Delete(bucket, key string, opts ...Option) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if bucketData, exists := c.Data[bucket]; exists {
		if _, exists := bucketData[key]; exists {
			delete(bucketData, key)
			return nil
		}
	}

	return errors.New("key not found")
}
