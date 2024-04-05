package database

import (
	"context"
	"fmt"
	"time"

	"github.com/chheller/go-htmx-todo/modules/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() *mongo.Client {
	env := config.GetEnvironment()
	mongoConnectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s", env.MongoUserName, env.MongoPassword, env.MongoUrl)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	opts := options.Client().ApplyURI(mongoConnectionString).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	return client
}
