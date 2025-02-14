package server

import (
	"cacheapp/cache"
	"io"

	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

var (
	cacheImpl cache.Cache
	mutex     sync.RWMutex
)

func init() {
	cacheImpl = cache.NewCache()
}

func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	bucket := parts[2]
	key := parts[3]

	// Read the request body
	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := cacheImpl.Set(bucket, key, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	bucket := parts[2]
	key := parts[3]

	value, err := cacheImpl.Get(bucket, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Determine the content type
	contentType := http.DetectContentType(value)
	w.Header().Set("Content-Type", contentType)

	// Write the value to the response
	w.Write(value)
}

func GetCacheStatsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	stats := map[string]interface{}{
		"total_buckets": len(cacheImpl.(*cache.InMemoryCache).Data),
		"total_keys":    countTotalKeys(),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func countTotalKeys() int {
	totalKeys := 0
	for _, bucket := range cacheImpl.(*cache.InMemoryCache).Data {
		totalKeys += len(bucket)
	}
	return totalKeys
}
