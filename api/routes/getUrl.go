package routes

import (
	"net/http"

	"github.com/cevrimxe/url-shortener/api/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func GetByShortID(c *gin.Context) {
	shortID := c.Param("shortID")

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortID).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": val})
}
