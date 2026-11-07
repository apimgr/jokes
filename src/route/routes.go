package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	handlers "github.com/apimgr/jokes/src/handler"
	"github.com/gin-gonic/gin"
)

type rateWindow struct {
	started time.Time
	count   int
}

var apiRateLimiter = struct {
	sync.Mutex
	clients map[string]rateWindow
}{clients: make(map[string]rateWindow)}

// SetupRoutes configures public API routes and applies shared protections.
func SetupRoutes(router *gin.Engine, rateLimit int) {
	SetupRoutesWithLimits(router, true, rateLimit, time.Minute)
}

// SetupRoutesWithLimits configures public routes and the operator-selected limiter.
func SetupRoutesWithLimits(router *gin.Engine, enabled bool, rateLimit int, window time.Duration) {
	if rateLimit <= 0 {
		rateLimit = 120
	}
	if window <= 0 {
		window = time.Minute
	}
	router.Use(func(c *gin.Context) {
		c.Set("timestamp", time.Now().UTC().Format(time.RFC3339))
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'; img-src 'self' data:; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		}
		c.Next()
	})
	if enabled {
		router.Use(apiRateLimit(rateLimit, window))
	}

	// API routes expose only public, read-only data. There is deliberately no
	// admin route group because IDEA.md excludes web and API management.
	v1 := router.Group("/api/v1")
	{
		v1.GET("/server/healthz", handlers.HealthCheckJSON)
		v1.GET("/metrics", handlers.ServeMetrics)
		v1.GET("/openapi", handlers.ServeSwaggerUI)
		v1.GET("/openapi.json", handlers.ServeOpenAPIJSON)
		v1.GET("/openapi.yaml", handlers.ServeOpenAPIYAML)
		v1.GET("/graphql", handlers.ServeGraphiQL)
		v1.POST("/graphql", handlers.HandleGraphQL)
		registerAPIGet(v1, "/jokes/random", handlers.GetRandomJoke)
		registerAPIGet(v1, "/jokes/random/:count", handlers.GetRandomJokes)
		registerAPIGet(v1, "/jokes/search", handlers.SearchJokes)
		registerAPIGet(v1, "/jokes/categories", handlers.GetCategories)
		registerAPIGet(v1, "/jokes/categories/:category", handlers.GetJokesByCategory)
		registerAPIGet(v1, "/jokes/count", handlers.GetCount)
		registerAPIGet(v1, "/jokes/all", handlers.GetAllJokes)
		registerAPIGet(v1, "/jokes/:id", handlers.GetJokeByID)
		v1.GET("/server/swagger", handlers.ServeOpenAPIJSON)
		v1.POST("/server/graphql", handlers.HandleGraphQL)
	}

	router.GET("/api/healthz", handlers.HealthCheckJSON)
	router.GET("/metrics", handlers.ServeMetrics)
	router.GET("/openapi", handlers.ServeSwaggerUI)
	router.GET("/openapi.json", handlers.ServeOpenAPIJSON)
	router.GET("/openapi.yaml", handlers.ServeOpenAPIYAML)
	router.GET("/graphql", handlers.ServeGraphiQL)
	router.POST("/graphql", handlers.HandleGraphQL)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"type":  "error",
			"value": "Endpoint not found",
		})
	})
}

func registerAPIGet(group *gin.RouterGroup, route string, handler gin.HandlerFunc) {
	wrapped := negotiatedText(handler)
	group.GET(route, wrapped)
	if !strings.Contains(route, ":") {
		group.GET(route+".txt", wrapped)
	}
}

type bufferedResponse struct {
	gin.ResponseWriter
	status  int
	body    bytes.Buffer
	written bool
}

func (w *bufferedResponse) WriteHeader(status int) {
	if !w.written {
		w.status = status
	}
}

func (w *bufferedResponse) WriteHeaderNow() {}

func (w *bufferedResponse) Write(body []byte) (int, error) {
	w.written = true
	return w.body.Write(body)
}

func (w *bufferedResponse) WriteString(body string) (int, error) {
	w.written = true
	return w.body.WriteString(body)
}

func (w *bufferedResponse) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *bufferedResponse) Size() int     { return w.body.Len() }
func (w *bufferedResponse) Written() bool { return w.written }

func negotiatedText(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for i := range c.Params {
			c.Params[i].Value = strings.TrimSuffix(c.Params[i].Value, ".txt")
		}
		if !isTextRequest(c.Request) {
			handler(c)
			return
		}
		original := c.Writer
		buffer := &bufferedResponse{ResponseWriter: original}
		c.Writer = buffer
		handler(c)
		c.Writer = original
		original.Header().Set("Content-Type", "text/plain; charset=utf-8")
		original.Header().Del("Content-Length")
		original.WriteHeader(buffer.Status())
		_, _ = original.Write([]byte(responseText(buffer.body.Bytes())))
	}
}

func isTextRequest(r *http.Request) bool {
	if strings.HasSuffix(r.URL.Path, ".txt") || strings.Contains(r.Header.Get("Accept"), "text/plain") {
		return true
	}
	userAgent := strings.ToLower(r.UserAgent())
	for _, name := range []string{"curl/", "wget/", "httpie/", "python-requests/", "go-http-client/"} {
		if strings.Contains(userAgent, name) {
			return true
		}
	}
	return false
}

func responseText(body []byte) string {
	var envelope map[string]any
	if err := json.Unmarshal(body, &envelope); err != nil {
		return string(body)
	}
	value := envelope["value"]
	if value == nil {
		value = envelope["message"]
	}
	return formatTextValue(value)
}

func formatTextValue(value any) string {
	switch result := value.(type) {
	case string:
		return result + "\n"
	case []any:
		lines := make([]string, 0, len(result))
		for _, item := range result {
			lines = append(lines, strings.TrimSuffix(formatTextValue(item), "\n"))
		}
		return strings.Join(lines, "\n") + "\n"
	case map[string]any:
		if jokes, ok := result["jokes"].([]any); ok {
			return formatTextValue(jokes)
		}
		if setup, ok := result["setup"].(string); ok {
			delivery, _ := result["delivery"].(string)
			return setup + "\n" + delivery + "\n"
		}
		if joke, ok := result["joke"].(string); ok {
			return joke + "\n"
		}
		encoded, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Sprint(value) + "\n"
		}
		return string(encoded) + "\n"
	default:
		return fmt.Sprint(value) + "\n"
	}
}

func apiRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := requestIP(c.Request.RemoteAddr)
		now := time.Now()
		apiRateLimiter.Lock()
		if len(apiRateLimiter.clients) >= 100_000 {
			for key, entry := range apiRateLimiter.clients {
				if now.Sub(entry.started) >= window {
					delete(apiRateLimiter.clients, key)
				}
			}
		}
		bucket, exists := apiRateLimiter.clients[ip]
		if !exists && len(apiRateLimiter.clients) >= 100_000 {
			apiRateLimiter.Unlock()
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"ok": false, "error": "RATE_LIMIT_CAPACITY", "message": "Rate limiter is at capacity"})
			return
		}
		if bucket.started.IsZero() || now.Sub(bucket.started) >= window {
			bucket = rateWindow{started: now}
		}
		bucket.count++
		apiRateLimiter.clients[ip] = bucket
		blocked := bucket.count > limit
		remaining := int((window - now.Sub(bucket.started)).Seconds())
		apiRateLimiter.Unlock()
		if blocked {
			if remaining < 1 {
				remaining = 1
			}
			c.Header("Retry-After", strconv.Itoa(remaining))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"ok": false, "error": "RATE_LIMITED", "message": "Too many requests",
			})
			return
		}
		c.Next()
	}
}

func requestIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return ip
	}
	if parsed := net.ParseIP(remoteAddr); parsed != nil {
		return parsed.String()
	}
	return "unknown"
}
