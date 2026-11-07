package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	if Templates == nil {
		if err := InitTemplates(); err != nil {
			t.Fatalf("InitTemplates() error = %v", err)
		}
	}
	router := gin.New()
	router.Use(gin.Recovery())
	router.SetHTMLTemplate(Templates)
	return router
}

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name      string
		cookie    *http.Cookie
		wantTheme string
	}{
		{"dark cookie", &http.Cookie{Name: "theme", Value: "dark"}, "dark"},
		{"light cookie", &http.Cookie{Name: "theme", Value: "light"}, "light"},
		{"no cookie", nil, "auto"},
		{"invalid cookie value", &http.Cookie{Name: "theme", Value: "purple"}, "auto"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			c.Request = req

			got := getTheme(c)
			if got != tt.wantTheme {
				t.Errorf("getTheme() = %q, want %q", got, tt.wantTheme)
			}
		})
	}
}

func TestServeHome(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	router.GET("/", ServeHome)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ServeHome status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Jokes API") {
		t.Errorf("ServeHome body missing expected content, got: %s", w.Body.String())
	}
}

func TestServeBrowse(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	router.GET("/browse", ServeBrowse)

	req := httptest.NewRequest(http.MethodGet, "/browse?q=infinity", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ServeBrowse status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "counted to infinity") {
		t.Errorf("browse results do not contain expected search match: %s", w.Body.String())
	}
}

func TestServeRandom(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	router.GET("/random", ServeRandom)

	req := httptest.NewRequest(http.MethodGet, "/random", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Random joke") {
		t.Fatalf("ServeRandom response = %d %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Chuck Norris") {
		t.Errorf("random page missing joke content: %s", w.Body.String())
	}
}

func TestServeCategories(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	router.GET("/categories", ServeCategories)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Categories") {
		t.Fatalf("ServeCategories response = %d %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "dev") {
		t.Errorf("categories page missing category content: %s", w.Body.String())
	}
}

func TestServeAPIDocs(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)
	router.GET("/api-docs", ServeAPIDocs)

	req := httptest.NewRequest(http.MethodGet, "/api-docs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "/api/v1/jokes/search") {
		t.Fatalf("ServeAPIDocs response = %d %q", w.Code, w.Body.String())
	}
}

func TestServeManifest(t *testing.T) {
	router := gin.New()
	router.GET("/manifest.json", ServeManifest)

	req := httptest.NewRequest(http.MethodGet, "/manifest.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ServeManifest status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/manifest+json" {
		t.Errorf("ServeManifest Content-Type = %q, want application/manifest+json", ct)
	}
	if w.Body.Len() == 0 {
		t.Error("ServeManifest body is empty")
	}
}

func TestServeServiceWorker(t *testing.T) {
	router := gin.New()
	router.GET("/sw.js", ServeServiceWorker)

	req := httptest.NewRequest(http.MethodGet, "/sw.js", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ServeServiceWorker status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/javascript") {
		t.Errorf("ServeServiceWorker Content-Type = %q, want application/javascript", ct)
	}
	if sw := w.Header().Get("Service-Worker-Allowed"); sw != "/" {
		t.Errorf("ServeServiceWorker Service-Worker-Allowed = %q, want /", sw)
	}
	if w.Body.Len() == 0 {
		t.Error("ServeServiceWorker body is empty")
	}
}
