package web

import (
	"net/http"
	"time"

	"github.com/apimgr/jokes/src/models"
	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title            string
	Description      string
	Theme            string
	Page             string
	Year             int
	Version          string
	TotalJokes       int
	TotalCategories  int
	RandomJoke       *models.Joke
}

// ServeHome serves the home page
func ServeHome(c *gin.Context) {
	theme := getTheme(c)
	randomJoke := models.GetRandomJoke()

	data := PageData{
		Title:           "Home",
		Description:     "5,160+ jokes across 16 categories - freely available for everyone!",
		Theme:           theme,
		Page:            "home",
		Year:            time.Now().Year(),
		Version:         "1.0.0",
		TotalJokes:      models.GetJokesCount(),
		TotalCategories: len(models.GetCategories()),
		RandomJoke:      randomJoke,
	}

	c.HTML(http.StatusOK, "index.html", data)
}

// ServeBrowse serves the browse page
func ServeBrowse(c *gin.Context) {
	theme := getTheme(c)

	data := PageData{
		Title:       "Browse Jokes",
		Description: "Browse all jokes in our collection",
		Theme:       theme,
		Page:        "browse",
		Year:        time.Now().Year(),
		Version:     "1.0.0",
	}

	c.HTML(http.StatusOK, "browse.html", data)
}

// ServeRandom serves the random joke page
func ServeRandom(c *gin.Context) {
	theme := getTheme(c)
	randomJoke := models.GetRandomJoke()

	data := PageData{
		Title:       "Random Joke",
		Description: "Get a random joke",
		Theme:       theme,
		Page:        "random",
		Year:        time.Now().Year(),
		Version:     "1.0.0",
		RandomJoke:  randomJoke,
	}

	c.HTML(http.StatusOK, "random.html", data)
}

// ServeCategories serves the categories page
func ServeCategories(c *gin.Context) {
	theme := getTheme(c)

	data := PageData{
		Title:       "Categories",
		Description: "Browse jokes by category",
		Theme:       theme,
		Page:        "categories",
		Year:        time.Now().Year(),
		Version:     "1.0.0",
	}

	c.HTML(http.StatusOK, "categories.html", data)
}

// ServeAPIDocs serves the API documentation page
func ServeAPIDocs(c *gin.Context) {
	theme := getTheme(c)

	data := PageData{
		Title:       "API Documentation",
		Description: "Complete API documentation for Jokes API",
		Theme:       theme,
		Page:        "api-docs",
		Year:        time.Now().Year(),
		Version:     "1.0.0",
	}

	c.HTML(http.StatusOK, "api-docs.html", data)
}

// ServeManifest serves the PWA manifest
func ServeManifest(c *gin.Context) {
	c.Header("Content-Type", "application/manifest+json")
	c.FileFromFS("static/manifest.json", http.FS(EmbeddedFiles))
}

// ServeServiceWorker serves the service worker
func ServeServiceWorker(c *gin.Context) {
	c.Header("Content-Type", "application/javascript")
	c.Header("Service-Worker-Allowed", "/")
	c.FileFromFS("static/sw.js", http.FS(EmbeddedFiles))
}

// Helper function to get theme from cookie or default to dark
func getTheme(c *gin.Context) string {
	theme, err := c.Cookie("theme")
	if err != nil || (theme != "dark" && theme != "light") {
		return "dark"
	}
	return theme
}
