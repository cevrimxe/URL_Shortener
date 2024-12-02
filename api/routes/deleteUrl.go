package routes

import (
	"net/http"

	"github.com/cevrimxe/url-shortener/api/database"
	"github.com/gin-gonic/gin"
)

func DeleteURL(c *gin.Context) {
	shortID := c.Param("shortID")

	r := database.CreateClient(0)
	defer r.Close()

	err := r.Del(database.Ctx, shortID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete URL from Database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}
