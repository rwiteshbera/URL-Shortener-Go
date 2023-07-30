package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rwiteshbera/URL-Shortener-Go/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (config *ENV) PurgeDatabase(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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

		urlCollection := mongoClient.Database(config.MONGO_DATABASE).Collection(config.URL_COLLECTION)
		res, err := urlCollection.DeleteMany(context.TODO(), bson.D{{"expiration_date", bson.D{{"$lt", primitive.NewDateTimeFromTime(time.Now())}}}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", res.DeletedCount)))
	}
}
