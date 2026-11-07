package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apimgr/jokes/src/web"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestEffectiveVersion(t *testing.T) {
	origVersion := Version
	t.Cleanup(func() { Version = origVersion })

	t.Run("falls back to VERSION const when unstamped", func(t *testing.T) {
		Version = ""
		if got := effectiveVersion(); got != VERSION {
			t.Errorf("effectiveVersion() = %q, want %q", got, VERSION)
		}
	})

	t.Run("uses linker-stamped Version when set", func(t *testing.T) {
		Version = "9.9.9"
		if got := effectiveVersion(); got != "9.9.9" {
			t.Errorf("effectiveVersion() = %q, want 9.9.9", got)
		}
	})
}

func TestGetJokesPath(t *testing.T) {
	t.Run("explicit dataDir is joined with jokes.json", func(t *testing.T) {
		got := getJokesPath("/custom/data")
		want := filepath.Join("/custom/data", "jokes.json")
		if got != want {
			t.Errorf("getJokesPath(%q) = %q, want %q", "/custom/data", got, want)
		}
	})

	t.Run("empty dataDir falls back to a default path", func(t *testing.T) {
		got := getJokesPath("")
		if got == "" {
			t.Error("getJokesPath(\"\") returned empty string")
		}
		if !strings.HasSuffix(got, "jokes.json") {
			t.Errorf("getJokesPath(\"\") = %q, want a path ending in jokes.json", got)
		}
	})

	t.Run("resolves an existing file among the candidate paths", func(t *testing.T) {
		dir := t.TempDir()
		origWD, err := os.Getwd()
		if err != nil {
			t.Fatalf("Getwd() error = %v", err)
		}
		t.Cleanup(func() { os.Chdir(origWD) })

		if err := os.MkdirAll(filepath.Join(dir, "data"), 0755); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		jokesFile := filepath.Join(dir, "data", "jokes.json")
		if err := os.WriteFile(jokesFile, []byte("{}"), 0644); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("Chdir() error = %v", err)
		}

		got := getJokesPath("")
		if got != "data/jokes.json" {
			t.Errorf("getJokesPath(\"\") = %q, want data/jokes.json", got)
		}
	})
}

func TestServeRobotsTxt(t *testing.T) {
	router := gin.New()
	router.GET("/robots.txt", serveRobotsTxt)

	req := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("serveRobotsTxt status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("serveRobotsTxt Content-Type = %q, want text/plain prefix", ct)
	}
	if !strings.Contains(w.Body.String(), "User-agent: *") {
		t.Errorf("serveRobotsTxt body missing expected content, got: %s", w.Body.String())
	}
}

func TestServeSecurityTxt(t *testing.T) {
	router := gin.New()
	router.GET("/security.txt", serveSecurityTxt)

	req := httptest.NewRequest(http.MethodGet, "/security.txt", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("serveSecurityTxt status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("serveSecurityTxt Content-Type = %q, want text/plain prefix", ct)
	}
	if !strings.Contains(w.Body.String(), "Contact: mailto:security@casjay.us") {
		t.Errorf("serveSecurityTxt body missing expected content, got: %s", w.Body.String())
	}
}

func TestSetupWebRoutes(t *testing.T) {
	if err := web.InitTemplates(); err != nil {
		t.Fatalf("InitTemplates() error = %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.SetHTMLTemplate(web.Templates)
	setupWebRoutes(router)

	tests := []struct {
		name string
		path string
		want int
	}{
		{"robots.txt registered", "/robots.txt", http.StatusOK},
		{"security.txt registered", "/security.txt", http.StatusOK},
		{"well-known security.txt registered", "/.well-known/security.txt", http.StatusOK},
		{"manifest registered", "/manifest.json", http.StatusOK},
		{"service worker registered", "/sw.js", http.StatusOK},
		{"API docs registered", "/api-docs", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.want {
				t.Errorf("%s status = %d, want %d", tt.path, w.Code, tt.want)
			}
		})
	}
}

func TestPrintHelpDocumentsCLIOptions(t *testing.T) {
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = writer
	printHelp()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = original
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Jokes API", "Usage:", "--help", "--version", "--maintenance"} {
		if !strings.Contains(string(output), want) {
			t.Errorf("printHelp() output missing %q", want)
		}
	}
}

func TestGuaranteedRecoverySanitizesPanics(t *testing.T) {
	for _, tc := range []struct {
		name   string
		header string
		want   string
	}{
		{"api response is JSON", "application/json", `"error":"SERVER_ERROR"`},
		{"browser response is HTML", "text/html", "Something went wrong"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(guaranteedRecovery())
			router.GET("/", func(*gin.Context) { panic("private panic detail") })
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set("Accept", tc.header)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want 500", response.Code)
			}
			if !strings.Contains(response.Body.String(), tc.want) || strings.Contains(response.Body.String(), "private panic detail") {
				t.Fatalf("panic response was not sanitized: %s", response.Body.String())
			}
		})
	}
}
