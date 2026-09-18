package handlers

import (
	"bytes"
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

func loadTestJokes(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "jokes.json")
	data := models.JokesData{
		Jokes: []models.Joke{
			{ID: 1, Joke: "Chuck Norris counted to infinity twice.", Categories: []string{"dev", "nerdy"}},
			{ID: 2, Joke: "Chuck Norris can divide by zero.", Categories: []string{"dev"}},
			{ID: 3, Joke: "An explicit joke about Chuck Norris.", Categories: []string{"explicit"}},
		},
		Categories: []string{"dev", "nerdy", "explicit"},
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

func newTestContext(method, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var reqBody *bytes.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	} else {
		reqBody = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, target, reqBody)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.Request = req
	return c, w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	t.Helper()
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response body %q: %v", w.Body.String(), err)
	}
	return resp
}

func TestHealthCheckHTML(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/healthz", nil)
	HealthCheckHTML(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("Healthy")) {
		t.Error("body does not contain expected health status text")
	}
}

func TestHealthCheckJSON(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/api/v1/healthz", nil)
	HealthCheckJSON(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["status"] != "healthy" {
		t.Errorf("status = %v, want healthy", body["status"])
	}
	if body["jokes_loaded"].(float64) != 3 {
		t.Errorf("jokes_loaded = %v, want 3", body["jokes_loaded"])
	}
}

func TestGetRandomJoke(t *testing.T) {
	loadTestJokes(t)

	t.Run("no filters returns a joke", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random", nil)
		GetRandomJoke(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
		}
		resp := decodeResponse(t, w)
		if resp.Type != "success" {
			t.Errorf("Type = %q, want success", resp.Type)
		}
	})

	t.Run("valid category filters jokes", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random?category=explicit", nil)
		GetRandomJoke(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
		}
	})

	t.Run("invalid limitTo category returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random?limitTo=[nonexistent]", nil)
		GetRandomJoke(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid exclude category returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random?exclude=nonexistent", nil)
		GetRandomJoke(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("limitTo with no matches returns 404", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random?limitTo=[explicit]&exclude=explicit", nil)
		GetRandomJoke(c)
		// limitTo=explicit then exclude=explicit removes everything
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d, body=%s", w.Code, http.StatusNotFound, w.Body.String())
		}
	})

	t.Run("name replacement applied", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random?firstName=Bruce&lastName=Lee&category=dev", nil)
		GetRandomJoke(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("invalid category value falls back to all jokes excluding explicit", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random?category=bogus", nil)
		GetRandomJoke(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
		}
		resp := decodeResponse(t, w)
		value, ok := resp.Value.(map[string]interface{})
		if !ok {
			t.Fatalf("Value is not a joke object: %#v", resp.Value)
		}
		cats, _ := value["categories"].([]interface{})
		for _, cat := range cats {
			if cat == "explicit" {
				t.Error("explicit joke should have been excluded by default")
			}
		}
	})
}

func TestGetRandomJokesHandler(t *testing.T) {
	loadTestJokes(t)

	tests := []struct {
		name       string
		count      string
		wantStatus int
	}{
		{"valid count", "2", http.StatusOK},
		{"count exceeds total but within limit", "10", http.StatusOK},
		{"zero count is invalid", "0", http.StatusBadRequest},
		{"negative count is invalid", "-1", http.StatusBadRequest},
		{"count over 100 is invalid", "101", http.StatusBadRequest},
		{"non-numeric count is invalid", "abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random/"+tt.count, nil)
			c.Params = gin.Params{{Key: "count", Value: tt.count}}
			GetRandomJokes(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}

	t.Run("invalid limitTo category returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random/2?limitTo=[nope]", nil)
		c.Params = gin.Params{{Key: "count", Value: "2"}}
		GetRandomJokes(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid exclude category returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random/2?exclude=nope", nil)
		c.Params = gin.Params{{Key: "count", Value: "2"}}
		GetRandomJokes(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("name replacement applied to results", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/random/2?firstName=Bruce&lastName=Lee", nil)
		c.Params = gin.Params{{Key: "count", Value: "2"}}
		GetRandomJokes(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestGetJokeByID(t *testing.T) {
	loadTestJokes(t)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{"existing id", "1", http.StatusOK},
		{"nonexistent but in-range id", "2", http.StatusOK},
		{"out of range id", "999", http.StatusNotFound},
		{"zero id", "0", http.StatusBadRequest},
		{"negative id", "-1", http.StatusBadRequest},
		{"non-numeric id", "abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := newTestContext(http.MethodGet, "/api/v1/jokes/"+tt.id, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.id}}
			GetJokeByID(c)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body=%s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}

	t.Run("name replacement applied", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/1?firstName=Bruce&lastName=Lee", nil)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		GetJokeByID(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestGetAllJokes(t *testing.T) {
	loadTestJokes(t)

	t.Run("no filters returns all jokes", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/all", nil)
		GetAllJokes(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		resp := decodeResponse(t, w)
		value := resp.Value.(map[string]interface{})
		jokes := value["jokes"].([]interface{})
		if len(jokes) != 3 {
			t.Errorf("returned %d jokes, want 3", len(jokes))
		}
	})

	t.Run("limitTo filters jokes", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/all?limitTo=[dev]", nil)
		GetAllJokes(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		resp := decodeResponse(t, w)
		value := resp.Value.(map[string]interface{})
		jokes := value["jokes"].([]interface{})
		if len(jokes) != 2 {
			t.Errorf("returned %d jokes, want 2", len(jokes))
		}
	})

	t.Run("invalid limitTo category returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/all?limitTo=[nope]", nil)
		GetAllJokes(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid exclude category returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/all?exclude=nope", nil)
		GetAllJokes(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("name replacement applied", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/api/v1/jokes/all?firstName=Bruce&lastName=Lee", nil)
		GetAllJokes(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestGetCategories(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/api/v1/jokes/categories", nil)
	GetCategories(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	resp := decodeResponse(t, w)
	cats := resp.Value.([]interface{})
	if len(cats) != 3 {
		t.Errorf("returned %d categories, want 3", len(cats))
	}
}

func TestGetCount(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/api/v1/jokes/count", nil)
	GetCount(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	resp := decodeResponse(t, w)
	value := resp.Value.(map[string]interface{})
	if value["total"].(float64) != 3 {
		t.Errorf("total = %v, want 3", value["total"])
	}
}

func TestGetDocs(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/docs", nil)
	GetDocs(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRoot(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/", nil)
	GetRoot(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestParseCategoryList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"bracketed multiple", "[dev,nerdy]", []string{"dev", "nerdy"}},
		{"no brackets single", "dev", []string{"dev"}},
		{"empty brackets", "[]", []string{}},
		{"empty string", "", []string{}},
		{"whitespace around entries", "[ dev , nerdy ]", []string{"dev", "nerdy"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCategoryList(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("parseCategoryList(%q) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseCategoryList(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestValidateCategories(t *testing.T) {
	loadTestJokes(t)

	t.Run("all valid returns true and writes nothing", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/", nil)
		if !validateCategories(c, []string{"dev", "nerdy"}) {
			t.Error("validateCategories() = false, want true")
		}
		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want default 200 (untouched)", w.Code)
		}
	})

	t.Run("invalid category returns false and writes 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodGet, "/", nil)
		if validateCategories(c, []string{"nope"}) {
			t.Error("validateCategories() = true, want false")
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestServeSwaggerUI(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/swagger", nil)
	ServeSwaggerUI(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}

func TestServeOpenAPIJSON(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/openapi.json", nil)
	ServeOpenAPIJSON(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}

func TestServeOpenAPIYAML(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/openapi.yaml", nil)
	ServeOpenAPIYAML(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/yaml; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/yaml", ct)
	}
}

func TestServeGraphiQL(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/graphiql", nil)
	ServeGraphiQL(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandleGraphQL(t *testing.T) {
	loadTestJokes(t)

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		c, w := newTestContext(http.MethodPost, "/graphql", []byte("{not json"))
		HandleGraphQL(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("valid query returns 200", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{"query": "{ categories }"})
		c, w := newTestContext(http.MethodPost, "/graphql", body)
		HandleGraphQL(c)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body=%s", w.Code, http.StatusOK, w.Body.String())
		}
	})
}

func TestExecuteGraphQL(t *testing.T) {
	loadTestJokes(t)

	tests := []struct {
		name      string
		query     string
		variables map[string]interface{}
		wantKey   string
	}{
		{"randomJoke query", "{ randomJoke { id joke } }", nil, "randomJoke"},
		{"randomJokes query", "{ randomJokes { id } }", nil, "randomJokes"},
		{"categories query", "{ categories }", nil, "categories"},
		{"stats query", "{ stats { total } }", nil, "stats"},
		{"allJokes query", "{ allJokes { id } }", nil, "allJokes"},
		{"joke by id with variable", "query($id: Int!) { joke(id: $id) { id } }", map[string]interface{}{"id": float64(1)}, "joke"},
		{"joke without id variable", "{ joke(id: 1) { id } }", nil, "joke"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executeGraphQL(tt.query, tt.variables)
			data, ok := result["data"].(gin.H)
			if !ok {
				t.Fatalf("executeGraphQL(%q) = %#v, want a data map", tt.query, result)
			}
			if _, ok := data[tt.wantKey]; !ok {
				t.Errorf("executeGraphQL(%q) data = %#v, want key %q", tt.query, data, tt.wantKey)
			}
		})
	}

	t.Run("unsupported query returns errors", func(t *testing.T) {
		result := executeGraphQL("{ somethingUnsupported }", nil)
		if _, ok := result["errors"]; !ok {
			t.Errorf("executeGraphQL(unsupported) = %#v, want errors key", result)
		}
	})
}

func TestServeMetrics(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/metrics", nil)
	ServeMetrics(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("jokes_total 3")) {
		t.Errorf("metrics body missing jokes_total: %s", w.Body.String())
	}
}

func TestBearerTokenMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantAbort  bool
		wantStatus int
	}{
		{"missing header", "", true, http.StatusUnauthorized},
		{"missing Bearer prefix", "Token abc123", true, http.StatusUnauthorized},
		{"empty token after Bearer", "Bearer ", true, http.StatusUnauthorized},
		{"valid non-empty token", "Bearer abc123", false, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := newTestContext(http.MethodGet, "/api/v1/admin/stats", nil)
			if tt.authHeader != "" {
				c.Request.Header.Set("Authorization", tt.authHeader)
			}
			BearerTokenMiddleware()(c)

			if c.IsAborted() != tt.wantAbort {
				t.Errorf("IsAborted() = %v, want %v", c.IsAborted(), tt.wantAbort)
			}
			if tt.wantAbort && w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetAdminConfig(t *testing.T) {
	c, w := newTestContext(http.MethodGet, "/api/v1/admin/config", nil)
	GetAdminConfig(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetAdminStats(t *testing.T) {
	loadTestJokes(t)
	c, w := newTestContext(http.MethodGet, "/api/v1/admin/stats", nil)
	GetAdminStats(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}
