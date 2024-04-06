package database

import (
	"context"
	"fmt"
	"time"

	"github.com/chheller/go-htmx-todo/modules/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClientConfig struct {
	Username string
	Password string
	Url      string
}

func (config *MongoClientConfig) GetMongoConnectionString() string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s", config.Username, config.Password, config.Url)
}
func GetMongoClient() *mongo.Client {
	env := config.GetEnvironment()
	mongoConnectionString := env.MongoConfig.GetMongoConnectionString()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(mongoConnectionString).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	return client
}
