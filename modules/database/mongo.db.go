package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() *mongo.Client {
	mongoUserName := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoUrl := os.Getenv("MONGO_URL")
	mongoConnectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s", mongoUserName, mongoPassword, mongoUrl)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	opts := options.Client().ApplyURI(mongoConnectionString).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	return client
}
