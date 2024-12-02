package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cevrimxe/url-shortener/api/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type TagRequest struct {
	ShortID string `json:"shortID"`
	Tag     string `json:"tag"`
}

func AddTag(c *gin.Context) {
	var TagRequest TagRequest

	if err := c.ShouldBindJSON(&TagRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
		return
	}
	shortID := TagRequest.ShortID
	tag := TagRequest.Tag

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortID).Result()
	if err == redis.Nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		data = make(map[string]interface{})
		data["data"] = val
	}

	var tags []string
	if existingTags, ok := data["tags"].([]interface{}); ok {
		for _, t := range existingTags {
			if strTag, ok := t.(string); ok {
				tags = append(tags, strTag)
			}
		}
	}

	for _, existingTag := range tags {
		if existingTag == tag {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag already exists"})
			return
		}
	}

	tags = append(tags, tag)
	data["tags"] = tags

	updatedData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to Marshal Data"})
		return
	}

	err = r.Set(database.Ctx, shortID, string(updatedData), 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to Update URL in Database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added successfully"})
}
