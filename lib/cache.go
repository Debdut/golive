package lib

import (
	"fmt"
	"net/http"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		incrementRequest()
		h.ServeHTTP(w, r)
		fmt.Printf("No cache: %s (Method: %s)\n", r.URL.Path, r.Method)
	})
}

func useCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incrementRequest()
		h.ServeHTTP(w, r)
		fmt.Printf("Use cache: %s (Method: %s)\n", r.URL.Path, r.Method)
	})
}