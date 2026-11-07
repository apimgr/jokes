package route

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	models "github.com/apimgr/jokes/src/model"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func loadFixture(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "jokes.json")

	data := models.JokesData{
		Jokes: []models.Joke{
			{ID: 1, Joke: "Chuck Norris counted to infinity twice.", Categories: []string{"dev"}},
			{ID: 2, Joke: "Chuck Norris can divide by zero.", Categories: []string{"dev"}},
		},
		Categories: []string{"dev"},
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}
	if err := models.LoadJokes(path); err != nil {
		t.Fatalf("LoadJokes() error = %v", err)
	}
}

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	router := gin.New()
	router.Use(gin.Recovery())
	SetupRoutes(router, 120)
	return router
}

func TestSetupRoutesRegistersEndpoints(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"versioned server healthz", http.MethodGet, "/api/v1/server/healthz"},
		{"openapi json", http.MethodGet, "/openapi.json"},
		{"jokes count", http.MethodGet, "/api/v1/jokes/count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusNotFound {
				var body map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &body); err == nil && body["value"] == "Endpoint not found" {
					t.Errorf("%s %s hit the custom 404 handler, want a registered route", tt.method, tt.path)
				}
			}
		})
	}
}

func TestSetupRoutesHasNoAdminManagement(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("/admin/dashboard status = %d, want 404 because IDEA.md excludes management routes", w.Code)
	}
}

func TestSetupRoutesCustom404(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/this/path/does/not/exist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown path status = %d, want %d", w.Code, http.StatusNotFound)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal 404 body: %v", err)
	}
	if body["type"] != "error" || body["value"] != "Endpoint not found" {
		t.Errorf("404 body = %v, want {type: error, value: Endpoint not found}", body)
	}
}

func TestAPIRateLimitSetsRetryAfter(t *testing.T) {
	apiRateLimiter.Lock()
	apiRateLimiter.clients = make(map[string]rateWindow)
	apiRateLimiter.Unlock()
	t.Cleanup(func() {
		apiRateLimiter.Lock()
		apiRateLimiter.clients = make(map[string]rateWindow)
		apiRateLimiter.Unlock()
	})
	router := gin.New()
	router.Use(apiRateLimit(1, time.Minute))
	router.GET("/limited", func(c *gin.Context) { c.Status(http.StatusOK) })
	for i, wantStatus := range []int{http.StatusOK, http.StatusTooManyRequests} {
		req := httptest.NewRequest(http.MethodGet, "/limited", nil)
		req.RemoteAddr = "192.0.2.10:1234"
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)
		if response.Code != wantStatus {
			t.Fatalf("request %d status = %d, want %d", i+1, response.Code, wantStatus)
		}
		if i == 1 && response.Header().Get("Retry-After") == "" {
			t.Fatal("rate limit response missing Retry-After")
		}
	}
}

func TestTextExtensionRoutesIncludingParameters(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	for _, path := range []string{
		"/api/v1/jokes/random.txt",
		"/api/v1/jokes/random/2.txt",
		"/api/v1/jokes/categories/dev.txt",
		"/api/v1/jokes/1.txt",
	} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, path, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("GET %s status = %d, body=%s", path, response.Code, response.Body.String())
			}
			if contentType := response.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
				t.Errorf("GET %s content type = %q", path, contentType)
			}
		})
	}
}

func TestPlainTextResponseFormatting(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{"value string", `{"value":"plain"}`, "plain\n"},
		{"message fallback", `{"message":"fallback"}`, "fallback\n"},
		{"single joke", `{"value":{"joke":"one line"}}`, "one line\n"},
		{"two part joke", `{"value":{"setup":"setup","delivery":"punchline"}}`, "setup\npunchline\n"},
		{"joke list", `{"value":{"jokes":[{"joke":"first"},{"setup":"second","delivery":"last"}]}}`, "first\nsecond\nlast\n"},
		{"invalid json passthrough", "not-json", "not-json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := responseText([]byte(tt.body)); got != tt.want {
				t.Errorf("responseText() = %q, want %q", got, tt.want)
			}
		})
	}
	if got := requestIP("192.0.2.1:1234"); got != "192.0.2.1" {
		t.Errorf("requestIP(host:port) = %q", got)
	}
	if got := requestIP("2001:db8::1"); got != "2001:db8::1" {
		t.Errorf("requestIP(ip) = %q", got)
	}
	if got := requestIP("untrusted-value"); got != "unknown" {
		t.Errorf("requestIP(invalid) = %q", got)
	}
}

func TestNegotiatedTextPreservesJSONForBrowserRequests(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/jokes/count", nil)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "browser-test")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", response.Code, response.Body.String())
	}
	if contentType := response.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("content type = %q, want JSON", contentType)
	}
}
