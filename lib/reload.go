package lib

import (
	"fmt"
	"net/http"
	"time"
)

func reloadHandler(reloadChan chan bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		timeout := time.After(5 * time.Minute)
		for {
			select {
			case <-reloadChan:
				_, err := fmt.Fprintf(w, "data: reload\n\n")
				if err != nil {
					fmt.Printf("Error writing to SSE stream: %v\n", err)
					return
				}
				flusher.Flush()
			case <-timeout:
				// fmt.Println("SSE connection timed out")
				return
			case <-r.Context().Done():
				// fmt.Println("SSE connection closed by client")
				return
			}
		}
	})
}
