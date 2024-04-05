package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	MongoUserName string
	MongoPassword string
	MongoUrl      string
}

var env Environment
var isInitialized = false

func GetEnvironment() Environment {
	if !isInitialized {
		godotenv.Load(".env")
		MongoUserName, ok := os.LookupEnv("MONGO_USERNAME")
		if !ok {
			panic("Missing required environment variable MONGO_USERNAME")
		}
		MongoPassword, ok := os.LookupEnv("MONGO_PASSWORD")
		if !ok {
			panic("Missing required environment variable MONGO_PASSWORD")
		}
		MongoUrl, ok := os.LookupEnv("MONGO_URL")
		if !ok {
			panic("Missing required environment variable MONGO_URL")
		}
		env = Environment{
			MongoUserName,
			MongoPassword,
			MongoUrl,
		}
	}
	return env
}
