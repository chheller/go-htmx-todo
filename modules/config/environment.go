package config

import (
	"log"
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
