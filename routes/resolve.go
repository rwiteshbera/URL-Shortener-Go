package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/rwiteshbera/URL-Shortener-Go/database"
)

func ResolveURL(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/api", func(ctx *gin.Context) {
		url := ctx.Param("url")

		r := database.CreateClient(0)
		defer r.Close()

		value, err := r.Get(database.Ctx, url).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "short not found"})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot connect to database"})
			return
		}

		incrementCount := database.CreateClient(1)
		defer incrementCount.Close()

		_ = incrementCount.Incr(database.Ctx, "counter")

		ctx.Redirect(http.StatusMovedPermanently, value)
	})
}
