package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cevrimxe/url-shortener/api/database"
	"github.com/cevrimxe/url-shortener/api/models"
	"github.com/cevrimxe/url-shortener/api/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func ShortenURL(c *gin.Context) {
	if os.Getenv("DOMAIN") == "" || os.Getenv("API_QUOTA") == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	apiQuota := os.Getenv("API_QUOTA")
	if apiQuota == "" {
		apiQuota = "10" // Default value if not set
	}

	_, err := strconv.Atoi(apiQuota)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid API_QUOTA value"})
		return
	}

	var body models.Request

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Parse JSON Body"})
		return
	}
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.ClientIP()).Result()

	if err == redis.Nil {
		err = r2.Set(database.Ctx, c.ClientIP(), apiQuota, 30*60*time.Second).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit setup failed"})
			return
		}
		val = apiQuota // Set val to apiQuota for first-time users
	}

	if len(body.URL) > 2048 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL too long"})
		return
	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	if !utils.IsDifferentDomain(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Cannot Shorten Same Domain"})
		return
	}
	body.URL = utils.EnsureHttpPrefix(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "URL Custom Short Already in Use"})
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	if body.Expiry > 168 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expiry too long, maximum is 168 hours"})
		return
	}
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to Store URL in Database"})
		return
	}

	resp := models.Response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.ClientIP())

	val, _ = r2.Get(database.Ctx, c.ClientIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.ClientIP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	c.JSON(http.StatusOK, resp)
}
