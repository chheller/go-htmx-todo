package config

import (
	"log"
	"os"

	"github.com/chheller/go-htmx-todo/modules/database"
	"github.com/joho/godotenv"
)

type Environment struct {
	MongoConfig database.MongoClientConfig
}

var env Environment
var isInitialized = false

type EnvironmentLoader interface {
	Load(...string) error
}

func GetEnvironment(load ...func(...string) error) Environment {

	if !isInitialized {
		isInitialized = true
		log.Println("Enivironment uninitialized, loading environment variables from .env file")
		var err error
		if len(load) == 0 {
			err = godotenv.Load(".env")
		} else if len(load) == 1 {
			err = load[0](".env")
		} else {
			panic("Too many arguments passed to GetEnvironment, expected 1")
		}
		if err != nil {
			log.Fatal("Error loading .env file, falling back to environment variables")
		}

		mongoUserName, ok := os.LookupEnv("MONGO_USERNAME")
		if !ok {
			panic("Missing required environment variable MONGO_USERNAME")
		}
		mongoPassword, ok := os.LookupEnv("MONGO_PASSWORD")
		if !ok {
			panic("Missing required environment variable MONGO_PASSWORD")
		}
		mongoUrl, ok := os.LookupEnv("MONGO_URL")
		if !ok {
			panic("Missing required environment variable MONGO_URL")
		}
		env = Environment{
			MongoConfig: database.MongoClientConfig{
				Username: mongoUserName,
				Password: mongoPassword,
				Url:      mongoUrl,
			},
		}
	}
	return env
}
