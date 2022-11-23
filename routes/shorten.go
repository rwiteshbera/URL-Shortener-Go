package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	incomingRoutes.POST("/api", func(ctx *gin.Context) {
		var req request

		err := ctx.BindJSON(&req)

		req.Expiry = 24 * time.Hour
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Implement rate limiting

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

		// enfore https, SSl
		req.URL = helpers.EnforceHTTP(req.URL)

		ctx.JSON(http.StatusOK, gin.H{"req": req})
	})
}
