package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/apimgr/jokes/src/models"
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
		{"root healthz", http.MethodGet, "/healthz"},
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

func TestSetupRoutesAdminProtected(t *testing.T) {
	loadFixture(t)
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("/admin/dashboard without auth should not return 200")
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
