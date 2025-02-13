package cache

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	cache := NewCache()

	err := cache.Set("bucket1", "key1", []byte("value1"))
	assert.NoError(t, err)

	value, err := cache.Get("bucket1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)
}

func TestSetWithTTL(t *testing.T) {
	cache := NewCache()

	err := cache.Set("bucket1", "key1", []byte("value1"), func(o *Options) error {
		o.ttl = 1 * time.Second
		return nil
	})
	assert.NoError(t, err)

	value, err := cache.Get("bucket1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)

	time.Sleep(2 * time.Second)

	value, err = cache.Get("bucket1", "key1")
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestDelete(t *testing.T) {
	cache := NewCache()

	err := cache.Set("bucket1", "key1", []byte("value1"))
	assert.NoError(t, err)

	err = cache.Delete("bucket1", "key1")
	assert.NoError(t, err)

	value, err := cache.Get("bucket1", "key1")
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestCacheFull(t *testing.T) {
	cache := NewCache()

	for i := 0; i < 255; i++ {
		err := cache.Set("bucket1", strconv.Itoa(i), []byte("value"))
		assert.NoError(t, err)
	}

	err := cache.Set("bucket1", "key256", []byte("value"))
	assert.Error(t, err)
	assert.Equal(t, "cache is full", err.Error())
}
