package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	routes "github.com/rwiteshbera/URL-Shortener-Go/routes"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Unable to load .env in main.go!")
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "6001"
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	routes.ResolveURL(router)
	routes.ShortenURL(router)

	if err := router.Run(":" + PORT); err != nil {
		log.Fatalln("Failed to start the server!")
	}

	fmt.Println("Server is listening on " + PORT)
}
