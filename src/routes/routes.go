package routes

import (
	"time"

	"github.com/apimgr/jokes/src/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, rateLimit int) {
	// Middleware to add timestamp
	router.Use(func(c *gin.Context) {
		c.Set("timestamp", time.Now().Format(time.RFC3339))
		c.Next()
	})

	// Root-level endpoints (BASE.md required)
	router.GET("/healthz", handlers.HealthCheckHTML) // HTML version at root
	router.GET("/metrics", handlers.ServeMetrics)    // Prometheus metrics

	// OpenAPI/Swagger endpoints (root level)
	router.GET("/openapi", handlers.ServeSwaggerUI)
	router.GET("/openapi.json", handlers.ServeOpenAPIJSON)
	router.GET("/openapi.yaml", handlers.ServeOpenAPIYAML)

	// GraphQL endpoints (root level)
	router.GET("/graphql", handlers.ServeGraphiQL)
	router.POST("/graphql", handlers.HandleGraphQL)

	// Admin endpoints (root level)
	router.GET("/admin", handlers.ServeAdminLogin)
	router.POST("/admin/login", handlers.HandleAdminLogin)
	router.GET("/admin/logout", handlers.HandleAdminLogout)

	// Protected admin routes
	admin := router.Group("/admin")
	admin.Use(handlers.AdminAuthMiddleware())
	{
		admin.GET("/dashboard", handlers.ServeAdminDashboard)
		admin.GET("/settings", handlers.ServeAdminSettings)
		admin.GET("/logs", handlers.ServeAdminLogs)
		admin.GET("/backup", handlers.ServeAdminBackup)
	}

	// Documentation endpoints
	router.GET("/api", handlers.GetRoot)
	router.GET("/docs", handlers.GetDocs)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check (JSON)
		v1.GET("/healthz", handlers.HealthCheckJSON)

		// Metrics
		v1.GET("/metrics", handlers.ServeMetrics)

		// OpenAPI in v1
		v1.GET("/openapi", handlers.ServeSwaggerUI)
		v1.GET("/openapi.json", handlers.ServeOpenAPIJSON)
		v1.GET("/openapi.yaml", handlers.ServeOpenAPIYAML)

		// GraphQL in v1
		v1.GET("/graphql", handlers.ServeGraphiQL)
		v1.POST("/graphql", handlers.HandleGraphQL)

		// Jokes endpoints
		v1.GET("/jokes/random", handlers.GetRandomJoke)
		v1.GET("/jokes/random/:count", handlers.GetRandomJokes)
		v1.GET("/jokes/categories", handlers.GetCategories)
		v1.GET("/jokes/count", handlers.GetCount)
		v1.GET("/jokes/all", handlers.GetAllJokes)
		v1.GET("/jokes/:id", handlers.GetJokeByID)

		// Admin API (token auth)
		adminAPI := v1.Group("/admin")
		{
			adminAPI.GET("/config", handlers.GetAdminConfig)
			adminAPI.GET("/stats", handlers.GetAdminStats)
		}
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"type":  "error",
			"value": "Endpoint not found",
		})
	})
}
