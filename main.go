package main

import (
	"log"
	"net/http"

	"github.com/cevrimxe/url-shortener/api/database"
	"github.com/cevrimxe/url-shortener/api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"POST", "GET", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	router.Use(cors.New(config))

	// Add redirect handler for shortened URLs
	router.GET("/:shortID", func(c *gin.Context) {
		shortID := c.Param("shortID")

		r := database.CreateClient(0)
		defer r.Close()

		val, err := r.Get(database.Ctx, shortID).Result()
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, val)
	})

	// Test route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/shorten", routes.ShortenURL)
	}

	router.Run(":8080")
}
