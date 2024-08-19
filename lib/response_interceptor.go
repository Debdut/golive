package lib

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

type responseInterceptor struct {
	http.ResponseWriter
	body          *bytes.Buffer
	statusCode    int
	headerWritten bool
}

func (rec *responseInterceptor) WriteHeader(statusCode int) {
	if !rec.headerWritten {
		rec.statusCode = statusCode
		rec.ResponseWriter.WriteHeader(statusCode)
		rec.headerWritten = true
	}
}

func (rec *responseInterceptor) Write(p []byte) (int, error) {
	if !rec.headerWritten {
		rec.WriteHeader(http.StatusOK) // Sets default status to 200 OK if not set
	}
	if rec.body == nil {
		rec.body = &bytes.Buffer{}
	}
	return rec.body.Write(p)
}

func (rec *responseInterceptor) injectReloadScript() {
	if rec.body != nil {
		bodyString := rec.body.String()
		reloadScript := `<script>
            const source = new EventSource('/reload');
            source.onmessage = function(event) {
                if (event.data === 'reload') {
                    window.location.reload();
                }
            };
        </script></body>`
		modifiedBody := strings.Replace(bodyString, "</body>", reloadScript, 1)

		// Write the modified content to the original ResponseWriter
		if !rec.headerWritten {
			rec.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
			rec.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(modifiedBody)))
			rec.ResponseWriter.WriteHeader(rec.statusCode)
			rec.headerWritten = true
		}
		rec.ResponseWriter.Write([]byte(modifiedBody))
	} else if !rec.headerWritten {
		// If there's no body and the header hasn't been written, just write the status code
		rec.ResponseWriter.WriteHeader(rec.statusCode)
		rec.headerWritten = true
	}
}
