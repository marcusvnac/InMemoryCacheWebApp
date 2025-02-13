package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"cacheapp/cache"

	"github.com/stretchr/testify/assert"
)

var (
	testCache cache.Cache
	testMutex sync.RWMutex
)

func init() {
	testCache = cache.NewCache()
}

func TestSetKeyHandler(t *testing.T) {
	reqBody := []byte(`"value1"`)
	req, err := http.NewRequest("POST", "/cache/bucket1/key1", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SetKeyHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	value, err := testCache.Get("bucket1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)
}

func TestGetKeyHandler(t *testing.T) {
	testCache.Set("bucket1", "key1", []byte("value1"))

	req, err := http.NewRequest("GET", "/cache/bucket1/key1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetKeyHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var value []byte
	err = json.NewDecoder(rr.Body).Decode(&value)
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value)
}

func TestGetKeyHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/cache/bucket1/key2", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetKeyHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetCacheStatsHandler(t *testing.T) {
	testCache.Set("bucket1", "key1", []byte("value1"))
	testCache.Set("bucket2", "key2", []byte("value2"))

	req, err := http.NewRequest("GET", "/cache/stats", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCacheStatsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var stats map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&stats)
	assert.NoError(t, err)
	assert.Equal(t, 2, int(stats["total_buckets"].(float64)))
	assert.Equal(t, 2, int(stats["total_keys"].(float64)))
}
