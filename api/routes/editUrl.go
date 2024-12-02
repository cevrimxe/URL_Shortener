package routes

import (
	"net/http"
	"time"

	"github.com/cevrimxe/url-shortener/api/database"
	"github.com/cevrimxe/url-shortener/api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func EditURL(c *gin.Context) {
	shortID := c.Param("shortID")
	var body models.Request

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortID).Result()
	if err == redis.Nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Update the URL

	err = r.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update URL in Database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL updated successfully"})
}
