package main

import (
	"net/http"

	"cacheapp/server"
)

func main() {
	http.HandleFunc("/cache/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			server.SetKeyHandler(w, r)
		} else if r.Method == http.MethodGet {
			server.GetKeyHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/cache/stats", server.GetCacheStatsHandler)

	http.ListenAndServe(":8080", nil)
}
