package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunFetchesRandomJoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/jokes/random" {
			t.Errorf("path = %q, want random joke endpoint", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("User-Agent"), "jokes-cli/") {
			t.Errorf("User-Agent = %q, want jokes-cli prefix", r.Header.Get("User-Agent"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"type":"success","value":{"id":7,"type":"single","joke":"A safe joke","categories":["general"]}}`))
	}))
	defer server.Close()

	var stdout, stderr strings.Builder
	if err := run([]string{"--server", server.URL}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v, stderr=%q", err, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "A safe joke") {
		t.Errorf("output = %q, want joke text", got)
	}
}

func TestRunPrintsTwoPartJoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"type":"success","value":{"id":8,"type":"twopart","setup":"Question?","delivery":"Answer.","categories":["general"]}}`))
	}))
	defer server.Close()

	var stdout, stderr strings.Builder
	if err := run([]string{"--server", server.URL}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v, stderr=%q", err, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "Question?\nAnswer.") {
		t.Errorf("output = %q, want setup and delivery", got)
	}
}

func TestRunRejectsInsecureRemoteServer(t *testing.T) {
	var stdout, stderr strings.Builder
	err := run([]string{"--server", "http://example.com"}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "must use HTTPS") {
		t.Errorf("run() error = %v, want HTTPS validation error", err)
	}
}

func TestResolveServerEnvironmentOverride(t *testing.T) {
	t.Setenv("JOKES_SERVER_PRIMARY", "https://api.example.test")
	if got := resolveServer(); got != "https://api.example.test" {
		t.Errorf("resolveServer() = %q, want environment URL", got)
	}
}
