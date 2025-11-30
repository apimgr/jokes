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

	// Documentation and health endpoints
	router.GET("/api", handlers.GetRoot)
	router.GET("/docs", handlers.GetDocs)
	router.GET("/healthz", handlers.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Jokes endpoints
		v1.GET("/jokes/random", handlers.GetRandomJoke)
		v1.GET("/jokes/random/:count", handlers.GetRandomJokes)
		v1.GET("/jokes/categories", handlers.GetCategories)
		v1.GET("/jokes/count", handlers.GetCount)
		v1.GET("/jokes/all", handlers.GetAllJokes)
		v1.GET("/jokes/:id", handlers.GetJokeByID)
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"type":  "error",
			"value": "Endpoint not found",
		})
	})
}
