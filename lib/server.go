package lib

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var requests uint64
var epoch = time.Unix(0, 0).Format(time.RFC1123)
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}
var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func incrementRequest() {
	atomic.AddUint64(&requests, 1)
}

// StartServer starts up the file server
func StartServer(dir, httpPort, httpsPort, certFile, keyFile string, cache bool) error {
	go Printer(dir, httpPort)
	fs := fileServer(dir)

	if cache {
		http.Handle("/", useCache(fs))
	} else {
		http.Handle("/", noCache(fs))
	}

	reloadChan := make(chan bool)
	go watchForChanges(dir, reloadChan)

	// SSE endpoint for client to listen for reload events
	http.Handle("/reload", reloadHandler(reloadChan))

	var wg sync.WaitGroup
	var httpErr, httpsErr error

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("Starting HTTP server on port %s\n", httpPort)
		httpErr = http.ListenAndServe(httpPort, nil)
		if httpErr != nil {
			fmt.Printf("HTTP server error: %v\n", httpErr)
		}
	}()

	// Start HTTPS server if cert and key are provided
	if certFile != "" && keyFile != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Starting HTTPS server on port %s\n", httpsPort)
			httpsErr = http.ListenAndServeTLS(httpsPort, certFile, keyFile, nil)
			if httpsErr != nil {
				fmt.Printf("HTTPS server error: %v\n", httpsErr)
			}
		}()
	}

	wg.Wait()

	if httpErr != nil {
		return fmt.Errorf("error starting HTTP server: %v", httpErr)
	}
	if httpsErr != nil {
		return fmt.Errorf("error starting HTTPS server: %v", httpsErr)
	}

	return nil
}
