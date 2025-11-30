package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/apimgr/jokes/src/models"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// HealthCheck handles health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":       "healthy",
		"timestamp":    c.GetTime("timestamp"),
		"version":      "1.0.0",
		"jokes_loaded": models.GetJokesCount(),
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
