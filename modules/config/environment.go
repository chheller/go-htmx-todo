package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Environment struct {
	MongoConfig                  *MongoClientConfig
	SmtpConfig                   *SmtpClientConfig
	EmailVerificationRedirectUrl string
}
type EnvironmentLoader interface {
	Load(...string) error
}

// Goroutine-safe singleton reference to the parsed environment variables per https://refactoring.guru/design-patterns/singleton/go/example
var env *Environment
var lock = &sync.Mutex{}

func GetEnvironment(load ...func(...string) error) *Environment {

	if env == nil {
		lock.Lock()
		defer lock.Unlock()
		// Second existential check in case env was initialized before the lock was aquired.
		if env == nil {
			log.Println("Enivironment uninitialized, loading environment variables from .env file")
			var err error
			if len(load) == 0 {
				err = godotenv.Load(".env")
			} else if len(load) == 1 {
				err = load[0](".env")
			} else {
				log.Panic("Too many arguments passed to GetEnvironment, expected 1")
			}
			if err != nil {
				log.Panic("Error loading .env file")
			}

			verificationRedirectUrl, ok := os.LookupEnv("EMAIL_VERIFICATION_REDIRECT_URL")

			if !ok {
				log.Panicf("Missing required environment variable EMAIL_VERIFICATION_REDIRECT_URL")
			}

			env = &Environment{
				MongoConfig:                  loadMongoVars(),
				SmtpConfig:                   loadSmtpVars(),
				EmailVerificationRedirectUrl: verificationRedirectUrl,
			}
		}
	}
	return env
}

type MongoClientConfig struct {
	Username string
	Password string
	Url      string
}

func (config *MongoClientConfig) GetMongoConnectionString() string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s", config.Username, config.Password, config.Url)
}

func loadMongoVars() *MongoClientConfig {
	mongoUserName, ok := os.LookupEnv("MONGO_USERNAME")
	if !ok {
		log.Panic("Missing required environment variable MONGO_USERNAME")
	}
	mongoPassword, ok := os.LookupEnv("MONGO_PASSWORD")
	if !ok {
		log.Panic("Missing required environment variable MONGO_PASSWORD")
	}
	mongoUrl, ok := os.LookupEnv("MONGO_URL")
	if !ok {
		log.Panic("Missing required environment variable MONGO_URL")
	}

	return &MongoClientConfig{
		Username: mongoUserName,
		Password: mongoPassword,
		Url:      mongoUrl,
	}

}

type SmtpClientConfig struct {
	Username    string
	DisplayName string
	Password    string
	Host        string
	Port        uint16
}

func loadSmtpVars() *SmtpClientConfig {
	smtpUserName, ok := os.LookupEnv("SMTP_USERNAME")
	if !ok {
		log.Panic("Missing required environment variable SMTP_USERNAME")
	}
	smtpDisplayName, ok := os.LookupEnv("SMTP_DISPLAY_NAME")
	if !ok {
		log.Panic("Missing required environment variable SMTP_DISPLAY_NAME")
	}
	smtpPassword, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		log.Panic("Missing required environment variable SMTP_PASSWORD")
	}
	smtpHost, ok := os.LookupEnv("SMTP_HOST")
	if !ok {
		log.Panic("Missing required environment variable SMTP_HOST")
	}
	smtpPort, ok := os.LookupEnv("SMTP_PORT")
	if !ok {
		log.Panic("Missing required environment variable SMTP_PORT")
	}
	smtpPortParsed, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Panicf("SMTP_PORT must be an integer, got %s", smtpPort)
	}

	return &SmtpClientConfig{
		Username:    smtpUserName,
		DisplayName: smtpDisplayName,
		Password:    smtpPassword,
		Host:        smtpHost,
		Port:        uint16(smtpPortParsed),
	}
}
