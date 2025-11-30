package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/apimgr/jokes/src/config"
	"github.com/apimgr/jokes/src/models"
	"github.com/apimgr/jokes/src/routes"
	"github.com/apimgr/jokes/src/web"
	"github.com/gin-gonic/gin"
)

const VERSION = "1.0.0"

var (
	showHelp    bool
	showVersion bool
	showStatus  bool
	dataDir     string
	configDir   string
	listenAddr  string
	port        int
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showStatus, "status", false, "Show status and health, exit with code")
	flag.StringVar(&dataDir, "data", "", "Set data dir")
	flag.StringVar(&configDir, "config", "", "Set the config dir")
	flag.StringVar(&listenAddr, "address", "", "Set listen address")
	flag.IntVar(&port, "port", 0, "Set the port (64365 or 80,443, etc)")
}

func main() {
	flag.Parse()

	// Handle --help
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Handle --version
	if showVersion {
		fmt.Printf("🎭 Jokes API v%s\n", VERSION)
		os.Exit(0)
	}

	// Handle --status
	if showStatus {
		exitCode := checkStatus()
		os.Exit(exitCode)
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Override config with CLI flags if provided
	if listenAddr != "" {
		cfg.Server.Host = listenAddr
	}
	if port != 0 {
		cfg.Server.Port = port
	}

	log.Printf("✅ Configuration loaded from: %s", config.GetConfigPath())
	log.Printf("🌐 Server will listen on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Determine jokes data file path
	jokesPath := getJokesPath(dataDir)
	log.Printf("📚 Loading jokes from: %s", jokesPath)

	// Load jokes data
	if err := models.LoadJokes(jokesPath); err != nil {
		log.Fatalf("❌ Failed to load jokes: %v", err)
	}

	log.Printf("🎉 Loaded %d jokes from database", models.GetJokesCount())

	// Initialize web templates
	if err := web.InitTemplates(); err != nil {
		log.Fatalf("❌ Failed to load templates: %v", err)
	}

	log.Printf("🎨 Templates loaded successfully")

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, cfg.Server.RateLimit)

	// Setup web routes
	setupWebRoutes(router)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 JOKES API server starting on %s", addr)
	log.Printf("🏥 Health check: http://%s/healthz", addr)
	log.Printf("📖 API docs: http://%s/docs", addr)
	log.Printf("🌐 Web interface: http://%s/", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func setupWebRoutes(router *gin.Engine) {
	// Serve static files
	router.StaticFS("/static", http.FS(web.EmbeddedFiles))

	// PWA files
	router.GET("/manifest.json", web.ServeManifest)
	router.GET("/sw.js", web.ServeServiceWorker)

	// Web pages
	router.GET("/", web.ServeHome)
	router.GET("/browse", web.ServeBrowse)
	router.GET("/random", web.ServeRandom)
	router.GET("/categories", web.ServeCategories)
	router.GET("/api-docs", web.ServeAPIDocs)

	// Swagger and GraphQL placeholders (to be implemented)
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/api-docs")
	})
	router.GET("/graphql", func(c *gin.Context) {
		c.Redirect(302, "/api-docs")
	})
}

func getJokesPath(dataDir string) string {
	// If data dir provided via flag, use it
	if dataDir != "" {
		return filepath.Join(dataDir, "jokes.json")
	}

	// Try multiple locations for jokes.json
	possiblePaths := []string{
		"src/data/jokes.json",
		"data/jokes.json",
		"/usr/share/apimgr/jokes/jokes.json",
		"/opt/apimgr/jokes/jokes.json",
	}

	// Check if executable directory has the data
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths = append([]string{
			filepath.Join(exeDir, "src/data/jokes.json"),
			filepath.Join(exeDir, "data/jokes.json"),
		}, possiblePaths...)
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Default to src/data/jokes.json
	return "src/data/jokes.json"
}

func printHelp() {
	fmt.Println("🎭 Jokes API - 5,160+ jokes across 16 categories")
	fmt.Println()
	fmt.Printf("Version: %s\n", VERSION)
	fmt.Println("Website: https://jokes.apimgr.us")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  jokes-api [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help               Show this help message")
	fmt.Println("  --version            Show version information")
	fmt.Println("  --status             Show status and health, exit with code")
	fmt.Println("  --data <dir>         Set data directory")
	fmt.Println("  --config <dir>       Set config directory")
	fmt.Println("  --address <addr>     Set listen address (default: 0.0.0.0)")
	fmt.Println("  --port <port>        Set port (default: random 64xxx or 3009)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  jokes-api")
	fmt.Println("  jokes-api --port 8080")
	fmt.Println("  jokes-api --address 127.0.0.1 --port 3009")
	fmt.Println("  jokes-api --data /path/to/data")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Root:    /etc/apimgr/jokes/server.yaml\n")
	fmt.Printf("  User:    ~/.config/apimgr/jokes/server.yaml\n")
	fmt.Println()
	fmt.Println("API Endpoints:")
	fmt.Println("  GET  /healthz                    Health check")
	fmt.Println("  GET  /api/v1/jokes/random        Random joke")
	fmt.Println("  GET  /api/v1/jokes/{id}          Specific joke")
	fmt.Println("  GET  /api/v1/jokes/all           All jokes")
	fmt.Println("  GET  /api/v1/jokes/categories    List categories")
	fmt.Println()
}

func checkStatus() int {
	fmt.Println("🎭 Jokes API Status Check")
	fmt.Println()

	// Check config file
	configPath := config.GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("✅ Config file: %s\n", configPath)
	} else {
		fmt.Printf("⚠️  Config file not found: %s\n", configPath)
	}

	// Check data file
	jokesPath := getJokesPath(dataDir)
	if _, err := os.Stat(jokesPath); err == nil {
		fmt.Printf("✅ Data file: %s\n", jokesPath)

		// Try to load jokes
		if err := models.LoadJokes(jokesPath); err == nil {
			count := models.GetJokesCount()
			fmt.Printf("✅ Loaded %d jokes\n", count)
			if count > 0 {
				fmt.Println()
				fmt.Printf("📊 Statistics:\n")
				fmt.Printf("   Total jokes: %d\n", count)
				fmt.Printf("   Categories: %d\n", len(models.GetCategories()))
				fmt.Println()
				fmt.Println("✅ Status: Healthy")
				return 0
			}
		} else {
			fmt.Printf("❌ Failed to load jokes: %v\n", err)
			return 1
		}
	} else {
		fmt.Printf("❌ Data file not found: %s\n", jokesPath)
		return 1
	}

	return 0
}
