package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	models "github.com/apimgr/jokes/src/model"
	"github.com/apimgr/jokes/src/swagger"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// HealthCheckHTML handles health check endpoint returning HTML
func HealthCheckHTML(c *gin.Context) {
	jokesCount := models.GetJokesCount()
	categories := models.GetCategories()

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Health Check - Jokes API</title>
    <style>
        :root { --bg: #282a36; --fg: #f8f8f2; --green: #50fa7b; --purple: #bd93f9; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); margin: 0; padding: 2rem; }
        .container { max-width: 600px; margin: 0 auto; }
        h1 { color: var(--purple); }
        .status { background: var(--green); color: #000; padding: 0.5rem 1rem; border-radius: 4px; display: inline-block; font-weight: bold; }
        .stats { margin-top: 2rem; }
        .stat { margin: 0.5rem 0; }
        .label { color: var(--purple); }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎭 Jokes API</h1>
        <p class="status">✅ Healthy</p>
        <div class="stats">
            <p class="stat"><span class="label">Version:</span> 1.0.0</p>
            <p class="stat"><span class="label">Jokes Loaded:</span> ` + strconv.Itoa(jokesCount) + `</p>
            <p class="stat"><span class="label">Categories:</span> ` + strconv.Itoa(len(categories)) + `</p>
            <p class="stat"><span class="label">Status:</span> Running</p>
        </div>
        <p style="margin-top: 2rem; color: #6272a4;">
            <a href="/api/v1/healthz" style="color: var(--purple);">JSON version</a>
        </p>
    </div>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// HealthCheckJSON handles health check endpoint returning JSON
func HealthCheckJSON(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":       "healthy",
		"version":      "1.0.0",
		"jokes_loaded": models.GetJokesCount(),
		"categories":   len(models.GetCategories()),
	})
}

// GetRandomJoke returns a random joke
func GetRandomJoke(c *gin.Context) {
	jokes, ok := jokeCandidates(c)
	if !ok {
		return
	}
	selected := models.RandomJokesFrom(jokes, 1)
	if len(selected) == 0 {
		c.JSON(http.StatusNotFound, Response{Type: "error", Value: "No jokes found matching the criteria"})
		return
	}
	joke := replaceJokeNames(selected[0], c.Query("firstName"), c.Query("lastName"))
	c.JSON(http.StatusOK, Response{Type: "success", Value: joke})
}

// GetRandomJokes returns multiple random jokes
func GetRandomJokes(c *gin.Context) {
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil || count < 1 || count > 100 {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "Count must be a number between 1 and 100"})
		return
	}
	jokes, ok := jokeCandidates(c)
	if !ok {
		return
	}
	jokes = models.RandomJokesFrom(jokes, count)
	for i := range jokes {
		jokes[i] = replaceJokeNames(jokes[i], c.Query("firstName"), c.Query("lastName"))
	}
	if len(jokes) == 0 {
		c.JSON(http.StatusNotFound, Response{Type: "error", Value: "No jokes found matching the criteria"})
		return
	}
	c.JSON(http.StatusOK, Response{Type: "success", Value: jokes})
}

func jokeCandidates(c *gin.Context) ([]models.Joke, bool) {
	jokes := models.GetAllJokes()
	category := strings.TrimSpace(c.Query("category"))
	if category != "" {
		if !validateCategories(c, []string{category}) {
			return nil, false
		}
		jokes = models.GetJokesByCategory(category)
	}
	limitTo := strings.TrimSpace(c.Query("limitTo"))
	if limitTo != "" {
		include := parseCategoryList(limitTo)
		if !validateCategories(c, include) {
			return nil, false
		}
		jokes = models.FilterJokesByCategories(jokes, include, nil)
	}
	excludeValue := strings.TrimSpace(c.Query("exclude"))
	if excludeValue != "" {
		exclude := splitCategoryList(excludeValue)
		if !validateCategories(c, exclude) {
			return nil, false
		}
		jokes = models.FilterJokesByCategories(jokes, nil, exclude)
	}
	return jokes, true
}

func splitCategoryList(value string) []string {
	values := strings.Split(strings.Trim(value, "[]"), ",")
	result := make([]string, 0, len(values))
	for _, item := range values {
		if item = strings.TrimSpace(item); item != "" {
			result = append(result, item)
		}
	}
	return result
}

func replaceJokeNames(joke models.Joke, firstName, lastName string) models.Joke {
	if firstName == "" && lastName == "" {
		return joke
	}
	joke.Joke = models.ReplaceNameInJoke(joke.Joke, firstName, lastName)
	joke.Setup = models.ReplaceNameInJoke(joke.Setup, firstName, lastName)
	joke.Delivery = models.ReplaceNameInJoke(joke.Delivery, firstName, lastName)
	return joke
}

// GetJokeByID returns a specific joke by ID
func GetJokeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Type:  "error",
			Value: "Invalid joke ID - must be a number",
		})
		return
	}

	if id < 1 {
		c.JSON(http.StatusBadRequest, Response{
			Type:  "error",
			Value: "Invalid joke ID - must be greater than 0",
		})
		return
	}

	joke := models.GetJokeByID(id)
	if joke == nil {
		c.JSON(http.StatusNotFound, Response{
			Type:  "error",
			Value: "Joke not found",
		})
		return
	}

	firstName := c.Query("firstName")
	lastName := c.Query("lastName")

	updated := replaceJokeNames(*joke, firstName, lastName)
	c.JSON(http.StatusOK, Response{
		Type:  "success",
		Value: updated,
	})
}

// GetAllJokes returns all jokes with optional filtering
func GetAllJokes(c *gin.Context) {
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	exclude := c.Query("exclude")
	limitTo := c.Query("limitTo")
	query := strings.TrimSpace(c.Query("q"))
	if len(query) > 200 {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "q must be 200 characters or fewer"})
		return
	}
	page, pageOK := boundedIntQuery(c, "page", 1, 1, 1_000_000)
	limit, limitOK := boundedIntQuery(c, "limit", 25, 1, 100)
	if !pageOK || !limitOK {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "page must be positive and limit must be between 1 and 100"})
		return
	}

	jokes := models.GetAllJokes()

	var limitToCategories []string
	if limitTo != "" {
		limitToCategories = parseCategoryList(limitTo)
		if !validateCategories(c, limitToCategories) {
			return
		}
	}

	var excludeCategories []string
	if exclude != "" {
		excludeCategories = strings.Split(exclude, ",")
		for i := range excludeCategories {
			excludeCategories[i] = strings.TrimSpace(excludeCategories[i])
		}
		if !validateCategories(c, excludeCategories) {
			return
		}
	}

	jokes = models.FilterJokesByCategories(jokes, limitToCategories, excludeCategories)
	jokes = models.SearchJokes(jokes, query)
	total := len(jokes)
	jokes = models.Paginate(jokes, (page-1)*limit, limit)

	// Replace names if requested
	if firstName != "" || lastName != "" {
		for i := range jokes {
			jokes[i] = replaceJokeNames(jokes[i], firstName, lastName)
		}
	}

	c.JSON(http.StatusOK, Response{
		Type: "success",
		Value: gin.H{
			"jokes": jokes,
			"meta": gin.H{
				"total_in_database":     models.GetJokesCount(),
				"total_matching":        total,
				"returned":              len(jokes),
				"page":                  page,
				"limit":                 limit,
				"limited_to_categories": limitToCategories,
				"excluded_categories":   excludeCategories,
			},
		},
	})
}

// SearchJokes returns a bounded page of case-insensitive text matches.
func SearchJokes(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "q is required"})
		return
	}
	if len(query) > 200 {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "q must be 200 characters or fewer"})
		return
	}
	jokes := models.GetAllJokes()
	if category := strings.TrimSpace(c.Query("category")); category != "" {
		if !validateCategories(c, []string{category}) {
			return
		}
		jokes = models.GetJokesByCategory(category)
	}
	jokes = models.SearchJokes(jokes, query)
	total := len(jokes)
	page, pageOK := boundedIntQuery(c, "page", 1, 1, 1_000_000)
	limit, limitOK := boundedIntQuery(c, "limit", 25, 1, 100)
	if !pageOK || !limitOK {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "page must be positive and limit must be between 1 and 100"})
		return
	}
	jokes = models.Paginate(jokes, (page-1)*limit, limit)
	c.JSON(http.StatusOK, Response{Type: "success", Value: gin.H{"jokes": jokes, "meta": gin.H{"query": query, "total_matching": total, "page": page, "limit": limit}}})
}

// GetJokesByCategory returns one bounded page of a category's jokes.
func GetJokesByCategory(c *gin.Context) {
	category := c.Param("category")
	if !validateCategories(c, []string{category}) {
		return
	}
	jokes := models.GetJokesByCategory(category)
	total := len(jokes)
	page, pageOK := boundedIntQuery(c, "page", 1, 1, 1_000_000)
	limit, limitOK := boundedIntQuery(c, "limit", 25, 1, 100)
	if !pageOK || !limitOK {
		c.JSON(http.StatusBadRequest, Response{Type: "error", Value: "page must be positive and limit must be between 1 and 100"})
		return
	}
	jokes = models.Paginate(jokes, (page-1)*limit, limit)
	c.JSON(http.StatusOK, Response{Type: "success", Value: gin.H{"category": category, "jokes": jokes, "meta": gin.H{"total_matching": total, "page": page, "limit": limit}}})
}

func boundedIntQuery(c *gin.Context, name string, fallback, minimum, maximum int) (int, bool) {
	value := c.Query(name)
	if value == "" {
		return fallback, true
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < minimum || parsed > maximum {
		return 0, false
	}
	return parsed, true
}

// GetCategories returns all available categories
func GetCategories(c *gin.Context) {
	categories := models.GetCategories()
	result := make([]gin.H, 0, len(categories))
	for _, category := range categories {
		result = append(result, gin.H{"category": category, "count": len(models.GetJokesByCategory(category))})
	}
	c.JSON(http.StatusOK, Response{Type: "success", Value: result})
}

// GetCount returns joke count and statistics
func GetCount(c *gin.Context) {
	categories := models.GetCategories()
	categoryStats := []gin.H{}

	for _, cat := range categories {
		jokes := models.GetJokesByCategory(cat)
		categoryStats = append(categoryStats, gin.H{
			"name":  cat,
			"count": len(jokes),
		})
	}

	c.JSON(http.StatusOK, Response{
		Type: "success",
		Value: gin.H{
			"total":      models.GetJokesCount(),
			"categories": categoryStats,
		},
	})
}

// GetDocs returns API documentation
func GetDocs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title":       "JOKES API Documentation",
		"version":     "1.0.0",
		"description": "A comprehensive jokes API with 5000+ jokes across multiple categories",
		"base_url":    c.Request.Host,
		"endpoints": gin.H{
			"/healthz":                            "Health check endpoint",
			"/api/v1/jokes/random":                "Get a random joke",
			"/api/v1/jokes/random/:count":         "Get multiple random jokes (1-100)",
			"/api/v1/jokes/:id":                   "Get a specific joke by ID",
			"/api/v1/jokes/all":                   "Browse jokes with filtering and pagination",
			"/api/v1/jokes/search":                "Search jokes by text",
			"/api/v1/jokes/categories/{category}": "Browse a category",
			"/api/v1/jokes/categories":            "Get all available categories",
			"/api/v1/jokes/count":                 "Get total number of jokes and category stats",
		},
		"query_parameters": gin.H{
			"firstName": "Replace 'Chuck' with this name",
			"lastName":  "Replace 'Norris' with this name",
			"limitTo":   "Array of categories to limit to. Format: [category] or [category1,category2]",
			"exclude":   "Comma-separated categories to exclude",
			"category":  "Filter by category",
			"q":         "Search joke text or category names (up to 200 characters)",
			"page":      "One-based result page (default 1)",
			"limit":     "Results per page, 1-100 (default 25)",
		},
		"categories": models.GetCategories(),
	})
}

// GetRoot returns root API information
func GetRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the JOKES API",
		"version": "1.0.0",
		"stats": gin.H{
			"total_jokes": models.GetJokesCount(),
			"categories":  len(models.GetCategories()),
		},
		"endpoints": gin.H{
			"/healthz":                    "Health check endpoint",
			"/docs":                       "API documentation",
			"/api/v1/jokes/random":        "Get a random joke",
			"/api/v1/jokes/random/:count": "Get multiple random jokes (1-100)",
			"/api/v1/jokes/:id":           "Get a specific joke by ID",
			"/api/v1/jokes/all":           "Get all jokes with optional filtering",
			"/api/v1/jokes/categories":    "Get all available categories",
			"/api/v1/jokes/count":         "Get total number of jokes and category stats",
		},
	})
}

// Helper functions

func parseCategoryList(limitTo string) []string {
	// Remove brackets and parse
	limitTo = strings.Trim(limitTo, "[]")
	limitTo = strings.TrimSpace(limitTo)

	if limitTo == "" {
		return []string{}
	}

	categories := strings.Split(limitTo, ",")
	result := []string{}
	for _, cat := range categories {
		cat = strings.TrimSpace(cat)
		if cat != "" {
			result = append(result, cat)
		}
	}
	return result
}

func validateCategories(c *gin.Context, categories []string) bool {
	invalidCats := []string{}
	for _, cat := range categories {
		if !models.IsCategoryValid(cat) {
			invalidCats = append(invalidCats, cat)
		}
	}

	if len(invalidCats) > 0 {
		c.JSON(http.StatusBadRequest, Response{
			Type:  "error",
			Value: "Invalid categories: " + strings.Join(invalidCats, ", ") + ". Available categories: " + strings.Join(models.GetCategories(), ", "),
		})
		return false
	}

	return true
}

// OpenAPI/Swagger handlers

// ServeSwaggerUI serves self-contained API documentation without third-party assets.
func ServeSwaggerUI(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Jokes API reference</title><style>body{font:1rem/1.5 system-ui,sans-serif;max-width:64rem;margin:2rem auto;padding:0 1rem;color:#20242a}code{background:#eee;padding:.15rem .3rem}li{margin:.7rem 0}</style></head><body><main><h1>Jokes API reference</h1><p>Public read-only API. No account or key is required.</p><ul><li><a href="/openapi.json">Swagger 2.0 JSON specification</a></li><li><a href="/openapi.yaml">Swagger YAML specification</a></li><li><code>GET /api/v1/jokes/random</code> — retrieve one random joke</li><li><code>GET /api/v1/jokes/random/{count}</code> — retrieve 1–100 random jokes</li><li><code>GET /api/v1/jokes/{id}</code> — retrieve a joke by ID</li><li><code>GET /api/v1/jokes/all?page=1&amp;limit=25</code> — browse jokes</li><li><code>GET /api/v1/jokes/search?q=dad</code> — search jokes</li><li><code>GET /api/v1/jokes/categories</code> — list categories</li><li><code>GET /api/v1/jokes/categories/{category}</code> — browse a category</li><li><code>POST /api/v1/graphql</code> — run a GraphQL query</li></ul><p>Use <code>Accept: text/plain</code> or append <code>.txt</code> to supported joke GET routes for plain text.</p><p><a href="/">Home</a></p></main></body></html>`)
}

// ServeOpenAPIJSON serves the OpenAPI specification as JSON
func ServeOpenAPIJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, swagger.SwaggerJSON)
}

// ServeOpenAPIYAML serves the OpenAPI specification as YAML
func ServeOpenAPIYAML(c *gin.Context) {
	// Convert JSON to YAML
	var data interface{}
	if err := json.Unmarshal([]byte(swagger.SwaggerJSON), &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OpenAPI spec"})
		return
	}

	c.Header("Content-Type", "text/yaml; charset=utf-8")
	c.YAML(http.StatusOK, data)
}

// GraphQL handlers

// ServeGraphiQL serves a browser-based GraphQL query form without external assets.
func ServeGraphiQL(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Jokes GraphQL</title><style>body{font:1rem/1.5 system-ui,sans-serif;max-width:64rem;margin:2rem auto;padding:0 1rem}textarea{box-sizing:border-box;width:100%;min-height:16rem;font:1rem/1.5 monospace}button{padding:.65rem 1rem}</style></head><body><main><h1>Jokes GraphQL</h1><p>Submit a query to the public read-only endpoint. Example: <code>{ categories stats { total } }</code></p><form method="post" action="/graphql"><label for="query">GraphQL query</label><textarea id="query" name="query" required>{ categories }</textarea><button type="submit">Run query</button></form><p>Endpoint: <code>POST /api/v1/graphql</code></p><a href="/">Home</a></main></body></html>`)
}

// HandleGraphQL handles GraphQL queries via POST
func HandleGraphQL(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64*1024)
	var request struct {
		Query     string                 `json:"query" form:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	formRequest := strings.HasPrefix(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded")
	var bindErr error
	if formRequest {
		bindErr = c.ShouldBind(&request)
	} else {
		bindErr = c.ShouldBindJSON(&request)
	}
	if bindErr != nil || strings.TrimSpace(request.Query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{"message": "Invalid request"}}})
		return
	}

	// Simple GraphQL query parser
	result := executeGraphQL(request.Query, request.Variables)
	if formRequest {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!doctype html><html lang="en"><meta charset="utf-8"><title>GraphQL result</title><body><main><h1>Query result</h1><pre>%s</pre><a href="/graphql">Run another query</a></main></body></html>`, template.HTMLEscapeString(mustJSON(result)))
		return
	}
	c.JSON(http.StatusOK, result)
}

// executeGraphQL executes a GraphQL query (simplified implementation)
func executeGraphQL(query string, variables map[string]interface{}) gin.H {
	query = strings.TrimSpace(query)

	// Parse query type
	if strings.Contains(query, "randomJoke") && !strings.Contains(query, "randomJokes") {
		joke := models.GetRandomJokes(1)
		if len(joke) > 0 {
			return gin.H{"data": gin.H{"randomJoke": joke[0]}}
		}
		return gin.H{"data": gin.H{"randomJoke": nil}}
	}

	if strings.Contains(query, "randomJokes") {
		// Extract count (default 5)
		count := 5
		jokes := models.GetRandomJokes(count)
		return gin.H{"data": gin.H{"randomJokes": jokes}}
	}

	if strings.Contains(query, "categories") && !strings.Contains(query, "jokesByCategory") {
		categories := models.GetCategories()
		return gin.H{"data": gin.H{"categories": categories}}
	}

	if strings.Contains(query, "stats") {
		categories := models.GetCategories()
		catStats := []gin.H{}
		for _, cat := range categories {
			jokes := models.GetJokesByCategory(cat)
			catStats = append(catStats, gin.H{"name": cat, "count": len(jokes)})
		}
		return gin.H{"data": gin.H{"stats": gin.H{
			"total":      models.GetJokesCount(),
			"categories": catStats,
		}}}
	}

	if strings.Contains(query, "allJokes") {
		jokes := models.GetAllJokes()
		return gin.H{"data": gin.H{"allJokes": jokes}}
	}

	if strings.Contains(query, "joke(") || strings.Contains(query, "joke (") {
		// Try to extract ID from query or variables
		if id, ok := variables["id"].(float64); ok {
			joke := models.GetJokeByID(int(id))
			return gin.H{"data": gin.H{"joke": joke}}
		}
		return gin.H{"data": gin.H{"joke": nil}}
	}

	return gin.H{"errors": []gin.H{{"message": "Query not supported. Try: randomJoke, randomJokes, categories, stats, allJokes, joke(id)"}}}
}

// Metrics handler

// ServeMetrics serves Prometheus-compatible metrics
func ServeMetrics(c *gin.Context) {
	jokesCount := models.GetJokesCount()
	categoriesCount := len(models.GetCategories())

	metrics := `# HELP jokes_total Total number of jokes in database
# TYPE jokes_total gauge
jokes_total ` + strconv.Itoa(jokesCount) + `

# HELP jokes_categories_total Total number of joke categories
# TYPE jokes_categories_total gauge
jokes_categories_total ` + strconv.Itoa(categoriesCount) + `

# HELP jokes_api_info API version information
# TYPE jokes_api_info gauge
jokes_api_info{version="1.0.0"} 1

# HELP jokes_api_up API health status (1 = up, 0 = down)
# TYPE jokes_api_up gauge
jokes_api_up 1
`
	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.String(http.StatusOK, metrics)
}

func mustJSON(value any) string {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return `{"errors":[{"message":"Unable to encode result"}]}`
	}
	return string(encoded)
}
