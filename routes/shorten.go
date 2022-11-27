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
			valueInt, _ := strconv.Atoi(value)

			if valueInt <= 0 {
				limit, _ := r.TTL(database.Ctx, ctx.ClientIP()).Result()
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Rate limit exceeded", "limit": limit})
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

		r = database.CreateClient(0)
		defer r.Close()

		if req.Expiry == 0 {
			req.Expiry = 24 // 24 Hours (Default)
		}

		err = r.Set(database.Ctx, id, req.URL, req.Expiry*3600*time.Second).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		r.Decr(database.Ctx, ctx.ClientIP())

		resp := response{
			URL:             req.URL,
			CustomShort:     "",
			Expiry:          req.Expiry,
			XRateRemaining:  10,
			XRateLimitReset: 30,
		}

		val, _ := r.Get(database.Ctx, ctx.ClientIP()).Result()
		resp.XRateRemaining, _ = strconv.Atoi(val)

		ttl, _ := r.TTL(database.Ctx, ctx.ClientIP()).Result()
		resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

		resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

		ctx.JSON(http.StatusOK, gin.H{"response": resp})
	})
}
