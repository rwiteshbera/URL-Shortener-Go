package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/rwiteshbera/URL-Shortener-Go/database"
	"github.com/rwiteshbera/URL-Shortener-Go/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
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

		req.Expiry = 24 * time.Hour
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Implement rate limiting
		r := database.CreateClient(1)
		defer r.Close()

		value, err := r.Get(database.Ctx, ctx.ClientIP()).Result()
		if err == redis.Nil {
			_ = r.Set(database.Ctx, ctx.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err() // 30 minutes to expire
		} else {
			val, _ := r.Get(database.Ctx, ctx.ClientIP()).Result()
			valueInt, _ := strconv.Atoi(val)

			if valueInt <= 0 {
				limit, _ := r.TTL(database.Ctx, ctx.ClientIP()).Result()
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Rate limit exceeded", "limit": limit})
			}
		}

		// Check if the input is an actual url
		if !goValidator.isURL(req.URL) {
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

		r.Decr(database.Ctx, ctx.ClientIP())

		ctx.JSON(http.StatusOK, gin.H{"req": req})
	})
}
