package database

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMongoClient() (*mongo.Client, error) {
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		return nil, errors.New("unable to connect with mongodb")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)

	// Create new client
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
