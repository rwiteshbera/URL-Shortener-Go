package routes

import (
	"os"

	"github.com/joho/godotenv"
)

type ENV struct {
	SERVER_BASE_URL string
	MONGO_URI       string
	MONGO_DATABASE  string
	URL_COLLECTION  string
	DOMAIN          string
	REDIS_ADDRESS   string
	REDIS_PASSWORD  string
}

func LoadConfig() *ENV {
	godotenv.Load()

	config := &ENV{
		SERVER_BASE_URL: os.Getenv("SERVER_BASE_URL"),
		MONGO_URI:       os.Getenv("MONGO_URI"),
		MONGO_DATABASE:  os.Getenv("MONGO_DATABASE"),
		URL_COLLECTION:  os.Getenv("URL_COLLECTION"),
		DOMAIN:          os.Getenv("DOMAIN"),
		REDIS_ADDRESS:   os.Getenv("REDIS_ADDRESS"),
		REDIS_PASSWORD:  os.Getenv("REDIS_PASSWORD"),
	}

	return config
}
