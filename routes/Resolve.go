package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rwiteshbera/URL-Shortener-Go/helpers"

	"github.com/gorilla/mux"
	"github.com/rwiteshbera/URL-Shortener-Go/database"
	"github.com/rwiteshbera/URL-Shortener-Go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (config *ENV) ResolveURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		shortId := mux.Vars(r)["id"]

		// Create Redis Client
		urlCache := database.CreateRedisClient(0)
		defer func(urlCache *redis.Client) {
			err := urlCache.Close()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}(urlCache)

		// Check Redis Cache for url
		value, _ := urlCache.Get(context.TODO(), shortId).Result()

		// If Original url is not in cache
		if value == "" {
			// Create MongoDB Client
			mongoClient, err := database.CreateMongoClient()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func(mongoClient *mongo.Client) {
				if err = mongoClient.Disconnect(context.TODO()); err != nil {
					fmt.Println(err.Error())
					return
				}
			}(mongoClient)

			// Find Data
			var resultURL models.URL
			urlCollection := mongoClient.Database(config.MONGO_DATABASE).Collection(config.URL_COLLECTION)
			err = urlCollection.FindOne(context.TODO(), bson.D{{"short_url", shortId}}).Decode(&resultURL)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					http.Error(w, "invalid short url", http.StatusBadRequest)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Set the cache
			// If the cache expiry time is more than 1 hour, then set TTL to 1 hour, else set the TTL 15 miutes
			if time.Until(resultURL.ExpirationDate) > time.Hour {
				urlCache.Set(context.TODO(), resultURL.ShortURL, resultURL.OriginalURL, time.Hour)
			} else {
				urlCache.Set(context.TODO(), resultURL.ShortURL, resultURL.OriginalURL, 15*time.Minute)
			}

			// Redirect the URL
			http.Redirect(w, r, helpers.EnforceHTTP(resultURL.OriginalURL), http.StatusPermanentRedirect)
		} else {
			// Redirect the URL
			http.Redirect(w, r, helpers.EnforceHTTP(value), http.StatusPermanentRedirect)
		}
	}
}
