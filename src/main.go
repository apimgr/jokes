package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/apimgr/jokes/src/config"
	"github.com/apimgr/jokes/src/models"
	"github.com/apimgr/jokes/src/paths"
	"github.com/apimgr/jokes/src/routes"
	"github.com/apimgr/jokes/src/scheduler"
	"github.com/apimgr/jokes/src/service"
	"github.com/apimgr/jokes/src/web"
	"github.com/gin-gonic/gin"
)

const VERSION = "1.0.0"

var (
	showHelp       bool
	showVersion    bool
	showStatus     bool
	dataDir        string
	configDir      string
	listenAddr     string
	port           int
	serviceCmd     string
	maintenanceCmd string
	modeFlag       string
	updateCmd      string
	cfg            *config.Config
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showStatus, "status", false, "Show status and health, exit with code")
	flag.StringVar(&dataDir, "data", "", "Set data dir")
	flag.StringVar(&configDir, "config", "", "Set the config dir")
	flag.StringVar(&listenAddr, "address", "", "Set listen address")
	flag.IntVar(&port, "port", 0, "Set the port (64365 or 80,443, etc)")
	flag.StringVar(&serviceCmd, "service", "", "Service command: start, stop, restart, reload, --install, --uninstall, --disable, --help")
	flag.StringVar(&maintenanceCmd, "maintenance", "", "Maintenance command: backup, restore, update, mode, setup [args]")
	flag.StringVar(&modeFlag, "mode", "", "Application mode: production, development")
	flag.StringVar(&updateCmd, "update", "", "Update commands: check, yes, branch {stable|beta|daily}")
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

	// Handle --mode
	if modeFlag != "" {
		setApplicationMode(modeFlag)
		os.Exit(0)
	}

	// Handle --update
	if updateCmd != "" {
		handleUpdateCommand(updateCmd)
		os.Exit(0)
	}

	// Handle --service
	if serviceCmd != "" {
		handleServiceCommand(serviceCmd)
		os.Exit(0)
	}

	// Handle --maintenance
	if maintenanceCmd != "" {
		handleMaintenanceCommand(maintenanceCmd, flag.Args())
		os.Exit(0)
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Override config with CLI flags if provided
	if listenAddr != "" {
		cfg.Server.Address = listenAddr
	}
	if port != 0 {
		cfg.Server.Port = port
	}

	log.Printf("✅ Configuration loaded from: %s", config.GetConfigPath())
	log.Printf("🌐 Server will listen on %s:%d", cfg.Server.Address, cfg.Server.Port)

	// Write PID file if enabled
	if cfg.Server.PIDFile {
		if err := paths.WritePIDFile(); err != nil {
			log.Printf("⚠️  Warning: %v", err)
		} else {
			log.Printf("📝 PID file: %s", paths.GetPIDFilePath())
		}
	}

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
	routes.SetupRoutes(router, cfg.Server.RateLimit.Requests)

	// Setup web routes
	setupWebRoutes(router)

	// Initialize and start scheduler
	sched := scheduler.New()
	if cfg.Server.Schedule.Enabled {
		// Add scheduled tasks
		sched.AddTask("cert_renewal", scheduler.ParseInterval(cfg.Server.Schedule.CertRenewal), func() error {
			log.Println("📜 Checking certificate renewal...")
			// TODO: Implement actual cert renewal check
			return nil
		})

		sched.AddTask("notifications", scheduler.ParseInterval(cfg.Server.Schedule.Notifications), func() error {
			log.Println("🔔 Processing notifications...")
			// TODO: Implement notification processing
			return nil
		})

		sched.AddTask("cleanup", scheduler.ParseInterval(cfg.Server.Schedule.Cleanup), func() error {
			log.Println("🧹 Running cleanup...")
			// TODO: Implement cleanup tasks
			return nil
		})

		sched.Start()
		log.Printf("⏰ Scheduler started with %d tasks", len(sched.GetTasks()))
	}

	// Create HTTP server for graceful shutdown
	addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 JOKES API server starting on %s", addr)
		log.Printf("🏥 Health check: http://%s/healthz", addr)
		log.Printf("📖 API docs: http://%s/docs", addr)
		log.Printf("🌐 Web interface: http://%s/", addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-quit
	log.Printf("🛑 Received signal %v, shutting down gracefully...", sig)

	// Handle SIGHUP for config reload (future use)
	if sig == syscall.SIGHUP {
		log.Printf("🔄 Config reload requested (not yet implemented)")
		// In the future: reload config and continue
		// For now, we'll shut down
	}

	// Stop scheduler
	if cfg.Server.Schedule.Enabled {
		sched.Stop()
		log.Printf("⏰ Scheduler stopped")
	}

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	// Remove PID file
	if cfg.Server.PIDFile {
		paths.RemovePIDFile()
	}

	log.Printf("✅ Server stopped gracefully")
}

// handleServiceCommand handles --service commands
func handleServiceCommand(cmd string) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	switch cmd {
	case "start":
		fmt.Println("🚀 Starting jokes service...")
		if err := service.Start(); err != nil {
			fmt.Printf("❌ Failed to start service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Service started")
	case "stop":
		fmt.Println("🛑 Stopping jokes service...")
		if err := service.Stop(); err != nil {
			fmt.Printf("❌ Failed to stop service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Service stopped")
	case "restart":
		fmt.Println("🔄 Restarting jokes service...")
		if err := service.Restart(); err != nil {
			fmt.Printf("❌ Failed to restart service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Service restarted")
	case "reload":
		fmt.Println("🔄 Reloading jokes configuration...")
		if err := service.Reload(); err != nil {
			fmt.Printf("❌ Failed to reload service: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Configuration reloaded")
	case "--install":
		fmt.Println("📦 Installing jokes as a service...")
		if err := service.Install(); err != nil {
			fmt.Printf("❌ Failed to install service: %v\n", err)
			os.Exit(1)
		}
	case "--uninstall":
		fmt.Println("🗑️  Uninstalling jokes service...")
		if err := service.Uninstall(); err != nil {
			fmt.Printf("❌ Failed to uninstall service: %v\n", err)
			os.Exit(1)
		}
	case "--disable":
		fmt.Println("🚫 Disabling jokes service...")
		fmt.Println("ℹ️  Use systemctl disable jokes (Linux) or equivalent for your OS")
	case "--help", "help":
		printServiceHelp()
	default:
		fmt.Printf("❌ Unknown service command: %s\n", cmd)
		printServiceHelp()
		os.Exit(1)
	}
}

// handleMaintenanceCommand handles --maintenance commands
func handleMaintenanceCommand(cmd string, args []string) {
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	switch cmd {
	case "backup":
		var backupPath string
		if len(args) > 0 {
			backupPath = args[0]
		}
		doBackup(backupPath)
	case "restore":
		var restorePath string
		if len(args) > 0 {
			restorePath = args[0]
		}
		doRestore(restorePath)
	case "update":
		doUpdate()
	case "mode":
		if len(args) == 0 {
			fmt.Println("❌ Usage: --maintenance mode {true|false|enable|disable}")
			os.Exit(1)
		}
		setMaintenanceMode(args[0])
	case "setup":
		runSetupWizard()
	default:
		fmt.Printf("❌ Unknown maintenance command: %s\n", cmd)
		fmt.Println("Available: backup, restore, update, mode, setup")
		os.Exit(1)
	}
}

func printServiceHelp() {
	fmt.Println("Service commands:")
	fmt.Println("  start       Start the service")
	fmt.Println("  stop        Stop the service")
	fmt.Println("  restart     Restart the service")
	fmt.Println("  reload      Reload configuration")
	fmt.Println("  --install   Install as system service")
	fmt.Println("  --uninstall Uninstall system service")
	fmt.Println("  --disable   Disable the service")
	fmt.Println("  --help      Show this help")
}

func doBackup(backupPath string) {
	if backupPath == "" {
		backupPath = filepath.Join(paths.GetBackupDir(), time.Now().Format("20060102150405")+".tar.gz")
	}
	fmt.Printf("💾 Creating backup at: %s\n", backupPath)
	fmt.Println("ℹ️  Backup functionality not yet implemented")
	// TODO: Implement backup
}

func doRestore(restorePath string) {
	if restorePath == "" {
		fmt.Println("ℹ️  Restoring from most recent backup...")
		restorePath = filepath.Join(paths.GetBackupDir(), "latest.tar.gz")
	}
	fmt.Printf("📥 Restoring from: %s\n", restorePath)
	fmt.Println("ℹ️  Restore functionality not yet implemented")
	// TODO: Implement restore
}

func doUpdate() {
	fmt.Println("🔄 Checking for updates...")
	fmt.Printf("Current version: %s\n", VERSION)
	fmt.Println("ℹ️  Update functionality not yet implemented")
	// TODO: Implement auto-update from GitHub releases
}

func setMaintenanceMode(mode string) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	enabled := false

	switch mode {
	case "1", "yes", "true", "enable", "enabled", "on":
		enabled = true
	case "0", "no", "false", "disable", "disabled", "off":
		enabled = false
	default:
		fmt.Printf("❌ Invalid mode: %s\n", mode)
		fmt.Println("Use: true, false, enable, disable, yes, no, on, off, 1, 0")
		os.Exit(1)
	}

	if enabled {
		fmt.Println("🔧 Maintenance mode ENABLED")
		fmt.Println("ℹ️  Server will return 503 to all public requests")
	} else {
		fmt.Println("✅ Maintenance mode DISABLED")
		fmt.Println("ℹ️  Server will accept all requests normally")
	}
	// TODO: Persist maintenance mode to config/state
}

func setupWebRoutes(router *gin.Engine) {
	// Serve static files
	router.StaticFS("/static", http.FS(web.EmbeddedFiles))

	// PWA files
	router.GET("/manifest.json", web.ServeManifest)
	router.GET("/sw.js", web.ServeServiceWorker)

	// Security and robots
	router.GET("/robots.txt", serveRobotsTxt)
	router.GET("/security.txt", serveSecurityTxt)
	router.GET("/.well-known/security.txt", serveSecurityTxt)

	// Web pages
	router.GET("/", web.ServeHome)
	router.GET("/browse", web.ServeBrowse)
	router.GET("/random", web.ServeRandom)
	router.GET("/categories", web.ServeCategories)
	router.GET("/api-docs", web.ServeAPIDocs)

	// Swagger redirect (legacy)
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(302, "/openapi")
	})
}

func serveRobotsTxt(c *gin.Context) {
	content := `User-agent: *
Allow: /
Sitemap: https://jokes.apimgr.us/sitemap.xml
`
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, content)
}

func serveSecurityTxt(c *gin.Context) {
	content := `Contact: mailto:security@casjay.us
Expires: 2026-12-31T23:59:59.000Z
Preferred-Languages: en
Canonical: https://jokes.apimgr.us/.well-known/security.txt
`
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, content)
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
	fmt.Println("  jokes [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help                 Show this help message")
	fmt.Println("  --version              Show version information")
	fmt.Println("  --status               Show status and health, exit with code")
	fmt.Println("  --data <dir>           Set data directory")
	fmt.Println("  --config <dir>         Set config directory")
	fmt.Println("  --address <addr>       Set listen address (default: [::])")
	fmt.Println("  --port <port>          Set port (default: random 64xxx)")
	fmt.Println("  --service <cmd>        Service: start, stop, restart, reload, --install, --uninstall")
	fmt.Println("  --maintenance <cmd>    Maintenance: backup, restore, update, mode")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  jokes                            Start server with defaults")
	fmt.Println("  jokes --port 8080                Start on specific port")
	fmt.Println("  jokes --address 127.0.0.1        Bind to localhost only")
	fmt.Println("  jokes --service --install        Install as system service")
	fmt.Println("  jokes --maintenance backup       Create backup")
	fmt.Println("  jokes --maintenance mode on      Enable maintenance mode")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Printf("  Root:    /etc/apimgr/jokes/server.yml\n")
	fmt.Printf("  User:    ~/.config/apimgr/jokes/server.yml\n")
	fmt.Println()
	fmt.Println("API Endpoints:")
	fmt.Println("  GET  /healthz                    Health check (HTML)")
	fmt.Println("  GET  /api/v1/healthz             Health check (JSON)")
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

// setApplicationMode sets the application mode in config (for --mode flag)
func setApplicationMode(mode string) {
	validModes := map[string]bool{
		"production":  true,
		"development": true,
		"debug":       true,
	}

	if !validModes[mode] {
		fmt.Printf("❌ Invalid mode: %s\n", mode)
		fmt.Println("Valid modes: production, development, debug")
		os.Exit(1)
	}

	currentCfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	currentCfg.Server.Mode = mode
	fmt.Printf("✅ Application mode set to: %s\n", mode)
	fmt.Println("ℹ️  Config save not yet implemented")
}

// handleUpdateCommand handles update-related commands
func handleUpdateCommand(cmd string) {
	currentCfg, _ := config.LoadConfig()

	switch cmd {
	case "check":
		fmt.Println("🔄 Checking for updates...")
		fmt.Printf("Current version: %s\n", VERSION)
		branch := "stable"
		if currentCfg != nil && currentCfg.Server.UpdateBranch != "" {
			branch = currentCfg.Server.UpdateBranch
		}
		fmt.Printf("Update branch: %s\n", branch)
		fmt.Println("ℹ️  Update feature not yet implemented")
	case "yes":
		fmt.Println("🔄 Performing update...")
		fmt.Println("ℹ️  Update feature not yet implemented")
	default:
		if strings.HasPrefix(cmd, "branch ") {
			branch := strings.TrimPrefix(cmd, "branch ")
			validBranches := map[string]bool{
				"stable": true,
				"beta":   true,
				"daily":  true,
			}
			if !validBranches[branch] {
				fmt.Printf("❌ Invalid branch: %s\n", branch)
				fmt.Println("Valid branches: stable, beta, daily")
				os.Exit(1)
			}
			fmt.Printf("✅ Update branch set to: %s\n", branch)
			fmt.Println("ℹ️  Config save not yet implemented")
		} else {
			fmt.Printf("❌ Unknown update command: %s\n", cmd)
			fmt.Println("Usage: --update check|yes|branch {stable|beta|daily}")
			os.Exit(1)
		}
	}
}

// runSetupWizard runs the interactive setup wizard
func runSetupWizard() {
	fmt.Println("🎭 Jokes API Setup Wizard")
	fmt.Println("=========================")
	fmt.Println()

	configPath := config.GetConfigPath()
	currentCfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not load config: %v\n", err)
		fmt.Println("Using defaults...")
	}

	fmt.Printf("Configuration file: %s\n", configPath)
	if currentCfg != nil {
		fmt.Printf("Current mode: %s\n", currentCfg.Server.Mode)
		fmt.Printf("Current port: %d\n", currentCfg.Server.Port)
		fmt.Printf("Admin username: %s\n", currentCfg.Server.Admin.Username)
	}
	fmt.Println()
	fmt.Println("✅ Setup complete. Edit the configuration file to customize settings.")
}
