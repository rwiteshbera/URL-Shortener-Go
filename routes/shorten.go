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

// A request is a struct with two fields, URL and Expiry, where URL is a string and Expiry is a
// time.Duration.
// @property {string} URL - The URL of the page to be shortened.
// @property Expiry - The time duration for which the URL should be stored.
type request struct {
	URL    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}

// `response` is a type that has a `string` field called `URL`, a `string` field called `CustomShort`,
// a `time.Duration` field called `Expiry`, an `int` field called `XRateRemaining`, and a
// `time.Duration` field called `XRateLimitReset`.
// @property {string} URL - The URL that was provided by user.
// @property {string} CustomShort - The custom short URL.
// @property Expiry - The time in seconds that the link will be active for.
// @property {int} XRateRemaining - The number of requests remaining in the current rate limit window.
// @property XRateLimitReset - The time in seconds until the rate limit resets.
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// It takes a URL as input, checks if it is valid, checks if the domain is valid, checks if the URL is using HTTPS, generates a short URL, and returns the short URL
func ShortenURL(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/shorten", func(ctx *gin.Context) {
		var req request

		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Implement rate limiting
		// Creating a new redis client.
		ipDatabase := database.CreateClient(1)
		defer ipDatabase.Close()

		// Getting the value of the key `ctx.ClientIP()` from the database.
		value, err := ipDatabase.Get(database.Ctx, ctx.ClientIP()).Result()
		if err == redis.Nil {
			// Setting the value of the key `ctx.ClientIP()` to the value of the environment variable
			// `API_QUOTA` and the time to live of the key `ctx.ClientIP()` to 30 minutes.
			_ = ipDatabase.Set(database.Ctx, ctx.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
		} else {
			// Converting the string value of the key `ctx.ClientIP()` to an integer.
			valueInt, _ := strconv.Atoi(value)

			if valueInt <= 0 {
				// Getting the time to live of the key `ctx.ClientIP()` from the database.
				limit, _ := ipDatabase.TTL(database.Ctx, ctx.ClientIP()).Result()
				ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded", "limit": limit})
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
		// Checking if the URL is using HTTPS or not. If not, it is adding HTTPS to the URL.
		req.URL = helpers.EnforceHTTP(req.URL)

		// Generate shorten URL
		id := uuid.New().String()[:6]

		urlDatabase := database.CreateClient(0)
		defer urlDatabase.Close()

		if req.Expiry == 0 {
			req.Expiry = 24 // 24 Hours (Default)
		}

		// Checking if the URL is already present in the database.
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

		// Getting the value of the key `ctx.ClientIP()` from the database.
		val, _ := ipDatabase.Get(database.Ctx, ctx.ClientIP()).Result()

		// Converting the string value of the key `ctx.ClientIP()` to an integer.
		resp.XRateRemaining, _ = strconv.Atoi(val)

		// Getting the time to live of the key `ctx.ClientIP()` from the database.
		ttl, _ := ipDatabase.TTL(database.Ctx, ctx.ClientIP()).Result()
		resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

		resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

		ctx.JSON(http.StatusOK, gin.H{"response": resp})
	})
}
