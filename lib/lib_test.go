package lib

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestIp(t *testing.T) {
	_, err := GetLocalIP()
	if err != nil {
		t.Errorf("GetLocalIP isn't working: %s", err)
	}
}
func TestServerWithoutCache(t *testing.T) {
	go StartServer(".", ":80", "", "", "", false, map[string]string{})
	resp, err := http.Get("http://127.0.0.1/")
	if err != nil {
		t.Errorf("Couldn't get for the test %s", err)
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.Header.Get("Cache-Control") == "" {
		t.Error("Server not setting cache control")
	}
	if resp.StatusCode != 200 {
		t.Error("Server not returning 200")
	}
	if body == nil {
		t.Error("Body is nil, server not working")
	}
}

func TestInitialPrint(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printStartMessage("/", ":80", "")

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "golive") {
		t.Error("Did not find golive")
	}

	if !strings.Contains(string(out), "Serving: /") {
		t.Error("Did not find serving message")
	}

	if !strings.Contains(string(out), "http://localhost:80/") {
		t.Error("Did not find local network print")
	}
}

func TestPrintServerInfo(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printServerInformation("/", ":80", "")

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "Net: http://") {
		t.Error("Did not net address")
	}

	if !strings.Contains(string(out), "Requests") {
		t.Error("Did not find requests")
	}
}
