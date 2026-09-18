package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/apimgr/jokes/src/models"
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
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	category := c.Query("category")
	exclude := c.Query("exclude")
	limitTo := c.Query("limitTo")

	var jokes []models.Joke

	// Handle limitTo parameter
	if limitTo != "" {
		limitToCategories := parseCategoryList(limitTo)
		if !validateCategories(c, limitToCategories) {
			return
		}
		jokes = models.GetAllJokes()
		jokes = models.FilterJokesByCategories(jokes, limitToCategories, nil)
		if len(jokes) == 0 {
			c.JSON(http.StatusNotFound, Response{
				Type:  "error",
				Value: "No jokes found for the specified limitTo categories",
			})
			return
		}
	} else if category != "" {
		if !models.IsCategoryValid(category) {
			jokes = models.GetAllJokes()
			if exclude == "" {
				exclude = "explicit"
			} else if !strings.Contains(exclude, "explicit") {
				exclude += ",explicit"
			}
		} else {
			jokes = models.GetJokesByCategory(category)
		}
	} else {
		jokes = models.GetAllJokes()
	}

	// Handle exclude parameter
	if exclude != "" {
		excludeCategories := strings.Split(exclude, ",")
		for i := range excludeCategories {
			excludeCategories[i] = strings.TrimSpace(excludeCategories[i])
		}
		if !validateCategories(c, excludeCategories) {
			return
		}
		jokes = models.FilterJokesByCategories(jokes, nil, excludeCategories)
	}

	if len(jokes) == 0 {
		c.JSON(http.StatusNotFound, Response{
			Type:  "error",
			Value: "No jokes found matching the criteria",
		})
		return
	}

	// Pick random joke from filtered list
	joke := jokes[0]
	if len(jokes) > 1 {
		jokes = models.GetRandomJokes(1)
		if len(jokes) > 0 {
			// Filter again to ensure it matches criteria
			filtered := models.FilterJokesByCategories(jokes, nil, strings.Split(exclude, ","))
			if len(filtered) > 0 {
				joke = filtered[0]
			}
		}
	}

	// Replace names if requested
	if firstName != "" || lastName != "" {
		joke.Joke = models.ReplaceNameInJoke(joke.Joke, firstName, lastName)
	}

	c.JSON(http.StatusOK, Response{
		Type:  "success",
		Value: joke,
	})
}

// GetRandomJokes returns multiple random jokes
func GetRandomJokes(c *gin.Context) {
	countStr := c.Param("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 || count > 100 {
		c.JSON(http.StatusBadRequest, Response{
			Type:  "error",
			Value: "Count must be a number between 1 and 100",
		})
		return
	}

	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	exclude := c.Query("exclude")
	limitTo := c.Query("limitTo")

	jokes := models.GetRandomJokes(count * 2) // Get extra for filtering

	// Handle limitTo parameter
	if limitTo != "" {
		limitToCategories := parseCategoryList(limitTo)
		if !validateCategories(c, limitToCategories) {
			return
		}
		jokes = models.FilterJokesByCategories(jokes, limitToCategories, nil)
	}

	// Handle exclude parameter
	if exclude != "" {
		excludeCategories := strings.Split(exclude, ",")
		for i := range excludeCategories {
			excludeCategories[i] = strings.TrimSpace(excludeCategories[i])
		}
		if !validateCategories(c, excludeCategories) {
			return
		}
		jokes = models.FilterJokesByCategories(jokes, nil, excludeCategories)
	}

	if len(jokes) == 0 {
		c.JSON(http.StatusNotFound, Response{
			Type:  "error",
			Value: "No jokes found matching the criteria",
		})
		return
	}

	// Limit to requested count
	if len(jokes) > count {
		jokes = jokes[:count]
	}

	// Replace names if requested
	if firstName != "" || lastName != "" {
		for i := range jokes {
			jokes[i].Joke = models.ReplaceNameInJoke(jokes[i].Joke, firstName, lastName)
		}
	}

	c.JSON(http.StatusOK, Response{
		Type:  "success",
		Value: jokes,
	})
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

	if id > models.GetJokesCount() {
		c.JSON(http.StatusNotFound, Response{
			Type:  "error",
			Value: "Joke ID " + idStr + " not found. Valid range: 1-" + strconv.Itoa(models.GetJokesCount()),
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

	if firstName != "" || lastName != "" {
		joke.Joke = models.ReplaceNameInJoke(joke.Joke, firstName, lastName)
	}

	c.JSON(http.StatusOK, Response{
		Type:  "success",
		Value: joke,
	})
}

// GetAllJokes returns all jokes with optional filtering
func GetAllJokes(c *gin.Context) {
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	exclude := c.Query("exclude")
	limitTo := c.Query("limitTo")

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

	// Replace names if requested
	if firstName != "" || lastName != "" {
		for i := range jokes {
			jokes[i].Joke = models.ReplaceNameInJoke(jokes[i].Joke, firstName, lastName)
		}
	}

	c.JSON(http.StatusOK, Response{
		Type: "success",
		Value: gin.H{
			"jokes": jokes,
			"meta": gin.H{
				"total_in_database":      models.GetJokesCount(),
				"returned":               len(jokes),
				"limited_to_categories":  limitToCategories,
				"excluded_categories":    excludeCategories,
			},
		},
	})
}

// GetCategories returns all available categories
func GetCategories(c *gin.Context) {
	categories := models.GetCategories()
	c.JSON(http.StatusOK, Response{
		Type:  "success",
		Value: categories,
	})
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
			"/healthz":                      "Health check endpoint",
			"/api/v1/jokes/random":          "Get a random joke",
			"/api/v1/jokes/random/:count":   "Get multiple random jokes (1-100)",
			"/api/v1/jokes/:id":             "Get a specific joke by ID",
			"/api/v1/jokes/all":             "Get all jokes with optional filtering",
			"/api/v1/jokes/categories":      "Get all available categories",
			"/api/v1/jokes/count":           "Get total number of jokes and category stats",
		},
		"query_parameters": gin.H{
			"firstName": "Replace 'Chuck' with this name",
			"lastName":  "Replace 'Norris' with this name",
			"limitTo":   "Array of categories to limit to. Format: [category] or [category1,category2]",
			"exclude":   "Comma-separated categories to exclude",
			"category":  "Filter by category",
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

// ServeSwaggerUI serves the Swagger UI page
func ServeSwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Jokes API - OpenAPI Documentation</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
    <style>
        body { margin: 0; background: #282a36; }
        .swagger-ui { max-width: 1200px; margin: 0 auto; }
        .swagger-ui .topbar { display: none; }
        .swagger-ui .info .title { color: #f8f8f2; }
        .swagger-ui .scheme-container { background: #44475a; }
        .swagger-ui select { background: #282a36; color: #f8f8f2; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/openapi.json',
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
            layout: "BaseLayout"
        });
    </script>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
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

// ServeGraphiQL serves the GraphiQL playground
func ServeGraphiQL(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Jokes API - GraphQL Playground</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphiql@3/graphiql.min.css">
    <style>
        body { margin: 0; height: 100vh; }
        #graphiql { height: 100vh; }
    </style>
</head>
<body>
    <div id="graphiql"></div>
    <script crossorigin src="https://cdn.jsdelivr.net/npm/react@18/umd/react.production.min.js"></script>
    <script crossorigin src="https://cdn.jsdelivr.net/npm/react-dom@18/umd/react-dom.production.min.js"></script>
    <script crossorigin src="https://cdn.jsdelivr.net/npm/graphiql@3/graphiql.min.js"></script>
    <script>
        const fetcher = GraphiQL.createFetcher({ url: '/graphql' });
        ReactDOM.createRoot(document.getElementById('graphiql')).render(
            React.createElement(GraphiQL, {
                fetcher,
                defaultQuery: ` + "`" + `# Welcome to the Jokes API GraphQL Playground!
# Try these queries:

query RandomJoke {
  randomJoke {
    id
    joke
    categories
  }
}

query AllCategories {
  categories
}

query Stats {
  stats {
    total
    categories {
      name
      count
    }
  }
}
` + "`" + `
            })
        );
    </script>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// HandleGraphQL handles GraphQL queries via POST
func HandleGraphQL(c *gin.Context) {
	var request struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{{"message": "Invalid request"}}})
		return
	}

	// Simple GraphQL query parser
	result := executeGraphQL(request.Query, request.Variables)
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

// BearerTokenMiddleware checks for valid bearer token in Authorization header
func BearerTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check for Bearer prefix
		if len(auth) < 7 || auth[:7] != "Bearer " {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token := auth[7:]
		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check against config token (would normally validate against stored tokens)
		// For now, just check if token is non-empty (actual validation would use config.Server.Admin.APIToken)
		c.Set("api_token", token)
		c.Next()
	}
}

// Admin API handlers

// GetAdminConfig returns current configuration (API)
func GetAdminConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"server": gin.H{
			"address": "[::]",
			"port":    "from config",
			"pidfile": true,
		},
		"web": gin.H{
			"theme": "dark",
			"cors":  "*",
		},
		"rate_limit": gin.H{
			"enabled":  true,
			"requests": 120,
			"window":   60,
		},
	})
}

// GetAdminStats returns server statistics (API)
func GetAdminStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"jokes":      models.GetJokesCount(),
		"categories": len(models.GetCategories()),
		"status":     "online",
		"version":    "1.0.0",
	})
}
