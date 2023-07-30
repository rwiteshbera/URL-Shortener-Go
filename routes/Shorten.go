package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/rwiteshbera/URL-Shortener-Go/database"
	"github.com/rwiteshbera/URL-Shortener-Go/helpers"
	"github.com/rwiteshbera/URL-Shortener-Go/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (config *ENV) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var request models.RequestInfo

		// Decode JSON Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// If there is no expiry mentioned, default will be 1 year
		if request.Expiry == 0 {
			request.Expiry = 365 * 24
		}

		// Validate original URL : Check if the input is an actual url
		if !govalidator.IsURL(request.OriginalURL) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		// Create MongoDB Client
		mongoClient, err := database.CreateMongoClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				log.Panic(err.Error())
			}
		}()

		// Generate Shorten URL : Base62 Encoding
		shortId := helpers.EncodeBase62()

		payload := &models.URL{
			OriginalURL:    helpers.RemoveDomainError(request.OriginalURL),
			ShortURL:       shortId,
			ExpirationDate: primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(request.Expiry) * time.Hour)),
			CreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		}

		// Insert Data
		urlCollection := mongoClient.Database(config.MONGO_DATABASE).Collection(config.URL_COLLECTION)
		_, err = urlCollection.InsertOne(context.TODO(), payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(&models.ResponseInfo{
			OriginalURL:    request.OriginalURL,
			ShortURL:       helpers.EnforceHTTP(config.DOMAIN + "/" + payload.ShortURL),
			CreatedAt:      payload.CreatedAt.Time().Format(time.RFC1123),
			ExpirationDate: payload.ExpirationDate.Time().Format(time.RFC1123),
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
