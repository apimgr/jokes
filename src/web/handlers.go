package web

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	models "github.com/apimgr/jokes/src/model"
	"github.com/gin-gonic/gin"
)

var AssetVersion = "1.0.1"

// SetVersion supplies the linker-stamped build version to HTML and static routes.
func SetVersion(version string) {
	if version != "" {
		AssetVersion = version
	}
}

type CategoryData struct {
	Name  string
	Count int
}

type PageData struct {
	Title           string
	Description     string
	Content         string
	Theme           string
	Page            string
	Year            int
	Version         string
	TotalJokes      int
	TotalCategories int
	RandomJoke      *models.Joke
	Jokes           []models.Joke
	Categories      []CategoryData
	Query           string
	Category        string
	PageNumber      int
	Limit           int
	PreviousPage    int
	NextPage        int
	HasPrevious     bool
	HasNext         bool
}

// ServeHome serves the home page
func ServeHome(c *gin.Context) {
	theme := getTheme(c)
	randomJoke := models.GetRandomJoke()

	data := commonPageData(c, "Home", "Curated jokes across multiple categories", "home")
	data.Theme = theme
	data.RandomJoke = randomJoke

	c.HTML(http.StatusOK, "base.html", data)
}

// ServeServerInfo renders public server information pages defined by IDEA.md.
func ServeServerInfo(c *gin.Context) {
	page := path.Base(c.Request.URL.Path)
	content := map[string]struct{ title, description, body string }{
		"about":   {"About", "About the Jokes API", "A public, read-only API and web catalog for a curated joke collection. The service embeds its dataset and does not require accounts or API keys."},
		"help":    {"Help", "How to use Jokes", "Browse the collection, search by keyword, or use the REST and GraphQL endpoints. Plain-text responses are available with the text endpoints documented on the API page."},
		"healthz": {"Service health", "Jokes service health", "The service is responding to HTTP requests."},
		"privacy": {"Privacy", "Jokes privacy information", "The joke dataset contains no personal information. Requests are processed to provide the requested page or API response; deployment operators may retain server access logs according to their own configuration."},
		"terms":   {"Terms", "Jokes service terms", "The service is provided as a public read-only reference. Joke content is presented from the curated dataset; users are responsible for complying with applicable laws when using the service."},
	}
	pageContent, ok := content[page]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	data := commonPageData(c, pageContent.title, pageContent.description, "server-info")
	data.Content = pageContent.body
	c.HTML(http.StatusOK, "base.html", data)
}

// ServeBrowse serves a server-rendered, searchable, paginated joke list.
func ServeBrowse(c *gin.Context) {
	page := parsePage(c.Query("page"))
	limit := parseLimit(c.Query("limit"))
	query := strings.TrimSpace(c.Query("q"))
	category := strings.TrimSpace(c.Query("category"))
	jokes := models.GetAllJokes()
	if category != "" {
		if !models.IsCategoryValid(category) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		jokes = models.GetJokesByCategory(category)
	}
	jokes = models.SearchJokes(jokes, query)
	total := len(jokes)
	pageCount := (total + limit - 1) / limit
	if pageCount > 0 && page > pageCount {
		page = pageCount
	}
	data := commonPageData(c, "Browse Jokes", "Browse and search the curated joke collection", "browse")
	data.Jokes = models.Paginate(jokes, (page-1)*limit, limit)
	for _, name := range models.GetCategories() {
		data.Categories = append(data.Categories, CategoryData{Name: name, Count: len(models.GetJokesByCategory(name))})
	}
	data.Query, data.Category = query, category
	data.PageNumber, data.Limit = page, limit
	data.HasPrevious, data.HasNext = page > 1, page < pageCount
	data.PreviousPage, data.NextPage = page-1, page+1
	c.HTML(http.StatusOK, "base.html", data)
}

// ServeCategory renders jokes for one known category.
func ServeCategory(c *gin.Context) {
	query := c.Request.URL.Query()
	query.Set("category", c.Param("category"))
	c.Request.URL.RawQuery = query.Encode()
	ServeBrowse(c)
}

// ServeSearch serves the search page using the browse results view.
func ServeSearch(c *gin.Context) {
	ServeBrowse(c)
}

// ServeRandom serves a random joke page, optionally scoped to a category.
func ServeRandom(c *gin.Context) {
	theme := getTheme(c)
	category := strings.TrimSpace(c.Query("category"))
	candidates := models.GetAllJokes()
	if category != "" {
		if !models.IsCategoryValid(category) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		candidates = models.GetJokesByCategory(category)
	}
	randomJokes := models.RandomJokesFrom(candidates, 1)
	data := commonPageData(c, "Random Joke", "Get a random joke", "random")
	data.Theme = theme
	data.Category = category
	for _, name := range models.GetCategories() {
		data.Categories = append(data.Categories, CategoryData{Name: name, Count: len(models.GetJokesByCategory(name))})
	}
	if len(randomJokes) != 0 {
		data.RandomJoke = &randomJokes[0]
	}
	c.HTML(http.StatusOK, "base.html", data)
}

// ServeCategories serves category names and their joke counts.
func ServeCategories(c *gin.Context) {
	data := commonPageData(c, "Categories", "Browse jokes by category", "categories")
	for _, name := range models.GetCategories() {
		data.Categories = append(data.Categories, CategoryData{Name: name, Count: len(models.GetJokesByCategory(name))})
	}
	c.HTML(http.StatusOK, "base.html", data)
}

// ServeAPIDocs serves the API documentation page.
func ServeAPIDocs(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", commonPageData(c, "API Documentation", "Jokes API endpoints and examples", "api-docs"))
}

func commonPageData(c *gin.Context, title, description, pageName string) PageData {
	return PageData{
		Title: title, Description: description, Theme: getTheme(c), Page: pageName,
		Year: time.Now().Year(), Version: AssetVersion, TotalJokes: models.GetJokesCount(),
		TotalCategories: len(models.GetCategories()), Limit: 25, PageNumber: 1,
	}
}

func parsePage(value string) int {
	page, err := strconv.Atoi(value)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func parseLimit(value string) int {
	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 {
		return 25
	}
	if limit > 100 {
		return 100
	}
	return limit
}

// ServeStatic serves embedded assets with version-aware cache headers.
func ServeStatic(c *gin.Context) {
	name := strings.TrimPrefix(c.Param("filepath"), "/")
	name = path.Clean(name)
	if name == "." || strings.HasPrefix(name, "../") || path.IsAbs(name) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	filePath := path.Join("static", name)
	file, err := EmbeddedFiles.Open(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || info.IsDir() {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if c.Query("v") == AssetVersion {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		c.Header("Cache-Control", "no-cache")
		c.Header("ETag", `"`+AssetVersion+`"`)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), bytes.NewReader(content))
}

// ServeManifest serves the PWA manifest with update-aware cache headers.
func ServeManifest(c *gin.Context) {
	contents, err := EmbeddedFiles.ReadFile("static/manifest.json")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "application/manifest+json")
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", `"`+AssetVersion+`"`)
	c.Data(http.StatusOK, "application/manifest+json", contents)
}

// ServeServiceWorker serves the current versioned service worker.
func ServeServiceWorker(c *gin.Context) {
	contents, err := EmbeddedFiles.ReadFile("static/sw.js")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	contents = bytes.ReplaceAll(contents, []byte("1.0.1"), []byte(AssetVersion))
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	c.Header("Service-Worker-Allowed", "/")
	c.Header("Cache-Control", "no-cache")
	c.Header("ETag", `"`+AssetVersion+`"`)
	c.Data(http.StatusOK, "application/javascript; charset=utf-8", contents)
}

// Helper function to get theme from cookie or default to dark
func getTheme(c *gin.Context) string {
	theme, err := c.Cookie("theme")
	if err != nil || (theme != "dark" && theme != "light" && theme != "auto") {
		return "auto"
	}
	return theme
}
