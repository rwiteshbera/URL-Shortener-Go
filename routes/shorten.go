package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/rwiteshbera/URL-Shortener-Go/database"
	"github.com/rwiteshbera/URL-Shortener-Go/helpers"
)

type request struct {
	URL    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/api/v1", func(ctx *gin.Context) {
		var req request

		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Implement rate limiting
		ipDatabase := database.CreateClient(1)
		defer ipDatabase.Close()

		value, err := ipDatabase.Get(database.Ctx, ctx.ClientIP()).Result()
		if err == redis.Nil {
			_ = ipDatabase.Set(database.Ctx, ctx.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err() // 30 minutes to expire
		} else {
			valueInt, _ := strconv.Atoi(value)

			if valueInt <= 0 {
				limit, _ := ipDatabase.TTL(database.Ctx, ctx.ClientIP()).Result()
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Rate limit exceededs", "limit": limit})
				return
			}
		}

		// Check if the input is an actual url
		if !govalidator.IsURL(req.URL) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
			return
		}

		// Check for domain error
		if !helpers.RemoveDomainError(req.URL) {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "invalid domain"})
			return
		}

		// enfore https, SSL
		req.URL = helpers.EnforceHTTP(req.URL)

		// Generate shorten URL
		id := uuid.New().String()[:6]

		urlDatabase := database.CreateClient(0)
		defer urlDatabase.Close()

		if req.Expiry == 0 {
			req.Expiry = 24 // 24 Hours (Default)
		}

		isPresent, err := urlDatabase.Exists(ctx, req.URL).Result()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if isPresent == 1 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "custom link is already used"})
			return
		}

		err = urlDatabase.Set(database.Ctx, id, req.URL, req.Expiry*time.Hour).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ipDatabase.Decr(database.Ctx, ctx.ClientIP())

		resp := response{
			URL:             req.URL,
			CustomShort:     "",
			Expiry:          req.Expiry,
			XRateRemaining:  10,
			XRateLimitReset: 30,
		}

		val, _ := ipDatabase.Get(database.Ctx, ctx.ClientIP()).Result()
		resp.XRateRemaining, _ = strconv.Atoi(val)

		ttl, _ := ipDatabase.TTL(database.Ctx, ctx.ClientIP()).Result()
		resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

		resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

		ctx.JSON(http.StatusOK, gin.H{"response": resp})
	})
}
