package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/flf2ko/playground/go-api-sample/database"
	"github.com/flf2ko/playground/go-api-sample/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize handlers
	linkHandler := handlers.NewLinkHandler(db)

	// Setup Gin router
	router := setupRouter(linkHandler)

	// Setup graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(linkHandler *handlers.LinkHandler) *gin.Engine {
	// Set Gin mode
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.ContextWithFallback = true

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API server is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/fetch-json", linkHandler.FetchJSON)
		api.GET("/records", linkHandler.GetRecords)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "go-api-sample",
			"version":     "1.0.0",
			"description": "A simple JSON fetcher API",
			"endpoints": gin.H{
				"health":     "GET /health",
				"fetch-json": "GET /api/v1/fetch-json?link=<url>",
				"records":    "GET /api/v1/records",
			},
		})
	})

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
