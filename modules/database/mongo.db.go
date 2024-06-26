package database

import (
	"context"
	"time"

	"github.com/chheller/go-htmx-todo/modules/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() *mongo.Client {
	env := config.GetEnvironment()
	mongoConnectionString := env.MongoConfig.GetMongoConnectionString()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(mongoConnectionString).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logrus.WithError(err).Panic("Failed to connect to database")
	}
	logrus.WithField("config", env.MongoConfig).Info("Connected to database")
	return client
}
