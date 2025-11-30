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

// Admin handlers

var adminSessions = make(map[string]string) // token -> username

// ServeAdminLogin serves the admin login page
func ServeAdminLogin(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Login - Jokes API</title>
    <style>
        :root { --bg: #282a36; --fg: #f8f8f2; --purple: #bd93f9; --green: #50fa7b; --red: #ff5555; --card: #44475a; }
        * { box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); margin: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
        .login-card { background: var(--card); padding: 2rem; border-radius: 8px; width: 100%; max-width: 400px; }
        h1 { color: var(--purple); margin-top: 0; text-align: center; }
        .form-group { margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; color: var(--purple); }
        input { width: 100%; padding: 0.75rem; border: 1px solid var(--purple); border-radius: 4px; background: var(--bg); color: var(--fg); font-size: 1rem; }
        input:focus { outline: none; border-color: var(--green); }
        button { width: 100%; padding: 0.75rem; background: var(--purple); color: var(--bg); border: none; border-radius: 4px; font-size: 1rem; cursor: pointer; margin-top: 1rem; }
        button:hover { background: var(--green); }
        .error { color: var(--red); text-align: center; margin-top: 1rem; }
        .back { text-align: center; margin-top: 1rem; }
        .back a { color: var(--purple); }
    </style>
</head>
<body>
    <div class="login-card">
        <h1>🎭 Admin Login</h1>
        <form method="POST" action="/admin/login">
            <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" name="username" required autocomplete="username">
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required autocomplete="current-password">
            </div>
            <button type="submit">Login</button>
        </form>
        <p class="back"><a href="/">← Back to Home</a></p>
    </div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// HandleAdminLogin handles admin login POST
func HandleAdminLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Default credentials (should be from config in production)
	if username == "administrator" && password == "admin" {
		// Generate session token
		token := generateToken()
		adminSessions[token] = username

		// Set cookie
		c.SetCookie("admin_token", token, 86400*30, "/", "", false, true)
		c.Redirect(302, "/admin/dashboard")
		return
	}

	c.HTML(http.StatusUnauthorized, "", gin.H{"error": "Invalid credentials"})
	ServeAdminLogin(c)
}

// AdminAuthMiddleware checks for valid admin session
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("admin_token")
		if err != nil || adminSessions[token] == "" {
			c.Redirect(302, "/admin")
			c.Abort()
			return
		}
		c.Set("admin_user", adminSessions[token])
		c.Next()
	}
}

// ServeAdminDashboard serves the admin dashboard
func ServeAdminDashboard(c *gin.Context) {
	user := c.GetString("admin_user")
	jokesCount := models.GetJokesCount()
	categories := models.GetCategories()

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Dashboard - Jokes API</title>
    <style>
        :root { --bg: #282a36; --fg: #f8f8f2; --purple: #bd93f9; --green: #50fa7b; --cyan: #8be9fd; --card: #44475a; }
        * { box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); margin: 0; }
        .header { background: var(--card); padding: 1rem 2rem; display: flex; justify-content: space-between; align-items: center; }
        .header h1 { margin: 0; color: var(--purple); }
        .header a { color: var(--fg); text-decoration: none; }
        .container { max-width: 1200px; margin: 0 auto; padding: 2rem; }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; }
        .card { background: var(--card); padding: 1.5rem; border-radius: 8px; }
        .card h2 { color: var(--purple); margin-top: 0; font-size: 1rem; }
        .card .value { font-size: 2rem; font-weight: bold; color: var(--green); }
        .nav { margin-top: 2rem; }
        .nav a { display: inline-block; background: var(--purple); color: var(--bg); padding: 0.5rem 1rem; border-radius: 4px; text-decoration: none; margin-right: 0.5rem; margin-bottom: 0.5rem; }
        .nav a:hover { background: var(--green); }
        .section { margin-top: 2rem; }
        .section h2 { color: var(--cyan); border-bottom: 1px solid var(--card); padding-bottom: 0.5rem; }
        table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
        th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid var(--card); }
        th { color: var(--purple); }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎭 Admin Dashboard</h1>
        <div>
            <span>Welcome, ` + user + `</span> |
            <a href="/admin/logout">Logout</a> |
            <a href="/">View Site</a>
        </div>
    </div>
    <div class="container">
        <div class="grid">
            <div class="card">
                <h2>Total Jokes</h2>
                <div class="value">` + strconv.Itoa(jokesCount) + `</div>
            </div>
            <div class="card">
                <h2>Categories</h2>
                <div class="value">` + strconv.Itoa(len(categories)) + `</div>
            </div>
            <div class="card">
                <h2>API Status</h2>
                <div class="value" style="color: var(--green);">Online</div>
            </div>
            <div class="card">
                <h2>Version</h2>
                <div class="value">1.0.0</div>
            </div>
        </div>

        <div class="nav">
            <a href="/admin/settings">⚙️ Settings</a>
            <a href="/admin/logs">📋 Logs</a>
            <a href="/admin/backup">💾 Backup</a>
            <a href="/healthz">🏥 Health</a>
            <a href="/openapi">📖 API Docs</a>
            <a href="/metrics">📊 Metrics</a>
        </div>

        <div class="section">
            <h2>Category Statistics</h2>
            <table>
                <tr><th>Category</th><th>Jokes</th></tr>`

	for _, cat := range categories {
		jokes := models.GetJokesByCategory(cat)
		html += `<tr><td>` + cat + `</td><td>` + strconv.Itoa(len(jokes)) + `</td></tr>`
	}

	html += `
            </table>
        </div>
    </div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// HandleAdminLogout handles admin logout
func HandleAdminLogout(c *gin.Context) {
	token, _ := c.Cookie("admin_token")
	delete(adminSessions, token)
	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	c.Redirect(302, "/admin")
}

// ServeAdminSettings serves the admin settings page
func ServeAdminSettings(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings - Jokes API Admin</title>
    <style>
        :root { --bg: #282a36; --fg: #f8f8f2; --purple: #bd93f9; --green: #50fa7b; --card: #44475a; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); margin: 0; padding: 2rem; }
        .container { max-width: 800px; margin: 0 auto; }
        h1 { color: var(--purple); }
        .back { margin-bottom: 1rem; }
        .back a { color: var(--purple); }
        .section { background: var(--card); padding: 1.5rem; border-radius: 8px; margin-bottom: 1rem; }
        .section h2 { color: var(--purple); margin-top: 0; }
        .form-group { margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; }
        input, select { width: 100%; padding: 0.5rem; background: var(--bg); color: var(--fg); border: 1px solid var(--purple); border-radius: 4px; }
        button { background: var(--purple); color: var(--bg); padding: 0.5rem 1rem; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: var(--green); }
    </style>
</head>
<body>
    <div class="container">
        <p class="back"><a href="/admin/dashboard">← Back to Dashboard</a></p>
        <h1>⚙️ Settings</h1>

        <div class="section">
            <h2>Server Settings</h2>
            <div class="form-group">
                <label>Listen Address</label>
                <input type="text" value="[::]" disabled>
            </div>
            <div class="form-group">
                <label>Port</label>
                <input type="text" value="(from config)" disabled>
            </div>
        </div>

        <div class="section">
            <h2>Web Settings</h2>
            <div class="form-group">
                <label>Theme</label>
                <select disabled><option>Dark (Dracula)</option><option>Light</option></select>
            </div>
            <div class="form-group">
                <label>CORS</label>
                <input type="text" value="*" disabled>
            </div>
        </div>

        <div class="section">
            <h2>Rate Limiting</h2>
            <div class="form-group">
                <label>Requests per minute</label>
                <input type="number" value="120" disabled>
            </div>
        </div>

        <p style="color: #6272a4;">Settings are read from server.yml configuration file.</p>
    </div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// ServeAdminLogs serves the admin logs page
func ServeAdminLogs(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Logs - Jokes API Admin</title>
    <style>
        :root { --bg: #282a36; --fg: #f8f8f2; --purple: #bd93f9; --card: #44475a; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); margin: 0; padding: 2rem; }
        .container { max-width: 1000px; margin: 0 auto; }
        h1 { color: var(--purple); }
        .back { margin-bottom: 1rem; }
        .back a { color: var(--purple); }
        .log-viewer { background: var(--card); padding: 1rem; border-radius: 8px; font-family: monospace; font-size: 0.85rem; max-height: 500px; overflow: auto; }
    </style>
</head>
<body>
    <div class="container">
        <p class="back"><a href="/admin/dashboard">← Back to Dashboard</a></p>
        <h1>📋 Logs</h1>
        <div class="log-viewer">
            <p style="color: #6272a4;">Log viewer not yet implemented.</p>
            <p style="color: #6272a4;">Logs are written to:</p>
            <p>- access.log (Apache format)</p>
            <p>- server.log (Application logs)</p>
        </div>
    </div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// ServeAdminBackup serves the admin backup page
func ServeAdminBackup(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Backup - Jokes API Admin</title>
    <style>
        :root { --bg: #282a36; --fg: #f8f8f2; --purple: #bd93f9; --green: #50fa7b; --card: #44475a; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: var(--bg); color: var(--fg); margin: 0; padding: 2rem; }
        .container { max-width: 800px; margin: 0 auto; }
        h1 { color: var(--purple); }
        .back { margin-bottom: 1rem; }
        .back a { color: var(--purple); }
        .section { background: var(--card); padding: 1.5rem; border-radius: 8px; margin-bottom: 1rem; }
        .section h2 { color: var(--purple); margin-top: 0; }
        button { background: var(--purple); color: var(--bg); padding: 0.75rem 1.5rem; border: none; border-radius: 4px; cursor: pointer; margin-right: 0.5rem; }
        button:hover { background: var(--green); }
    </style>
</head>
<body>
    <div class="container">
        <p class="back"><a href="/admin/dashboard">← Back to Dashboard</a></p>
        <h1>💾 Backup & Restore</h1>

        <div class="section">
            <h2>Create Backup</h2>
            <p>Create a backup of all configuration and data.</p>
            <button onclick="alert('Use CLI: jokes --maintenance backup')">Create Backup</button>
        </div>

        <div class="section">
            <h2>Restore from Backup</h2>
            <p>Restore configuration and data from a backup file.</p>
            <button onclick="alert('Use CLI: jokes --maintenance restore [file]')">Restore</button>
        </div>

        <div class="section">
            <h2>CLI Commands</h2>
            <pre style="background: var(--bg); padding: 1rem; border-radius: 4px; overflow-x: auto;">
jokes --maintenance backup
jokes --maintenance backup /path/to/backup.tar.gz
jokes --maintenance restore
jokes --maintenance restore /path/to/backup.tar.gz
            </pre>
        </div>
    </div>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
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

// Helper function to generate session token
func generateToken() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[i%len(chars)]
	}
	return string(b)
}
