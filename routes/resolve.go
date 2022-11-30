package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/rwiteshbera/URL-Shortener-Go/database"
)

func ResolveURL(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/:id", func(ctx *gin.Context) {
		url := ctx.Param("id")

		urlDatabase := database.CreateClient(0)
		defer urlDatabase.Close()

		value, err := urlDatabase.Get(database.Ctx, url).Result()
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "short not found"})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "cannot connect to database"})
			return
		}
		ctx.Redirect(http.StatusMovedPermanently, value)
	})
}
