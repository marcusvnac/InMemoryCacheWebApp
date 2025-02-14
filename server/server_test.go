package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMutex sync.RWMutex
)

func MustSetKeyUsingHandle(t *testing.T, bucket, key string, value []byte) {
	t.Helper()

	req, err := http.NewRequest("POST", fmt.Sprintf("/cache/%s/%s", bucket, key), bytes.NewBuffer(value))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SetKeyHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestSetKeyHandler(t *testing.T) {
	MustSetKeyUsingHandle(t, "bucket1", "key1", []byte("value1"))
}

func TestGetKeyHandler(t *testing.T) {
	MustSetKeyUsingHandle(t, "bucket1", "key1", []byte("value2"))

	req, err := http.NewRequest("GET", "/cache/bucket1/key1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetKeyHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, []byte("value2"), rr.Body.Bytes())
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
	MustSetKeyUsingHandle(t, "bucket1", "key1", []byte("value3"))
	MustSetKeyUsingHandle(t, "bucket1", "key2", []byte("value4"))

	req, err := http.NewRequest("GET", "/cache/stats", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCacheStatsHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var stats map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&stats)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(stats["total_buckets"].(float64)))
	assert.Equal(t, 2, int(stats["total_keys"].(float64)))
}
