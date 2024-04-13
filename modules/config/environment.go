package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type Environment struct {
	MongoConfig                  *MongoClientConfig
	SmtpConfig                   *SmtpClientConfig
	ApplicationConfiguration     *ApplicationConfiguration
	EmailVerificationRedirectUrl string
	InjectBrowserReload          bool
}

func (env *Environment) String() string {
	return fmt.Sprintf("MongoConfig: %v, SmtpConfig: %v, ApplicationConfiguration: %v, EmailVerificationRedirectUrl: %s, InjectBrowserReload: %t", env.MongoConfig, env.SmtpConfig, env.ApplicationConfiguration, env.EmailVerificationRedirectUrl, env.InjectBrowserReload)

}

type EnvironmentLoader interface {
	Load(...string) error
}

// Goroutine-safe singleton reference to the parsed environment variables per https://refactoring.guru/design-patterns/singleton/go/example
var env *Environment
var lock = &sync.Mutex{}

// GetEnvironment returns the singleton instance of the environment variables. This is the preferred way of accessing the application's configuration
func GetEnvironment(load ...func(...string) error) *Environment {

	if env == nil {
		lock.Lock()
		defer lock.Unlock()
		// Second existential check in case env was initialized before the lock was aquired.
		if env == nil {
			log.Info("Enivironment uninitialized, loading environment variables from .env file")
			var err error
			if len(load) == 0 {
				err = godotenv.Load(".env")
			} else if len(load) == 1 {
				err = load[0](".env")
			} else {
				log.Panic("Too many arguments passed to GetEnvironment, expected 1")
			}
			if err != nil {
				log.Panicf("Error loading .env file\n Error: %v", err)
			}

			verificationRedirectUrl, ok := os.LookupEnv("EMAIL_VERIFICATION_REDIRECT_URL")

			if !ok {
				log.Panicf("Missing required environment variable EMAIL_VERIFICATION_REDIRECT_URL")
			}

			env = &Environment{
				MongoConfig:                  loadMongoVars(),
				SmtpConfig:                   loadSmtpVars(),
				ApplicationConfiguration:     loadApplicationConfiguration(),
				EmailVerificationRedirectUrl: verificationRedirectUrl,
				InjectBrowserReload:          os.Getenv("INJECT_BROWSER_RELOAD") == "true",
			}
			log.WithField("environment", env).Debug("Environment loaded")
		}
	}
	return env
}

type ApplicationConfiguration struct {
	Port                uint32
	LogLevel            log.Level
	LoggerReportCaller  bool
	HttpPrintDebugError bool
	CookieSecret        string
	EmailOtpSecret      string
}

func (ac *ApplicationConfiguration) String() string {
	return fmt.Sprintf("Port: %d, LogLevel: %s, LoggerReportCaller: %t, HttpPrintDebugError: %t, CookieSecret: <masked>, EmailOtpSecret: <masked>", ac.Port, ac.LogLevel, ac.LoggerReportCaller, ac.HttpPrintDebugError)
}

// Loads general application configurations and packages them into a struct. Handles default values,
// parsing strings into valid number types, and panics if any required variables are missing.
func loadApplicationConfiguration() *ApplicationConfiguration {
	var logLevel log.Level
	logLevelStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		log.Info("No LOG_LEVEL environment variable found, defaulting to INFO")
	} else {
		var err error
		logLevel, err = log.ParseLevel(logLevelStr)
		if err != nil {
			log.Panicf("LOG_LEVEL must be one of logrus.Level, got %s\nError: %v", os.Getenv("LOG_LEVEL"), err)
		}
	}
	// Note: If we want to use anything outside of the default INFO level log, we need to set it immediately after loading the variable.
	// Any log messages logged before this method is called outside of the INFO scope will be ignored
	log.SetLevel(logLevel)

	var port uint32
	portStr, ok := os.LookupEnv("PORT")
	if !ok {
		log.Info("No PORT environment variable found, defaulting to 8080")
		port = 8080
	} else {
		p, err := strconv.ParseUint(portStr, 10, 32)
		if err != nil {
			log.Panicf("LOG_LEVEL must be an integer, got %s\nError: %v", os.Getenv("LOG_LEVEL"), err)
		}
		port = uint32(p)
	}

	return &ApplicationConfiguration{
		Port:                port,
		LogLevel:            logLevel,
		LoggerReportCaller:  os.Getenv("LOGGER_REPORT_CALLER") == "true",
		HttpPrintDebugError: os.Getenv("HTTP_PRINT_DEBUG_ERROR") == "true",
		CookieSecret:        getEnvironmentVariableOrPanic("COOKIE_SECRET"),
		EmailOtpSecret:      getEnvironmentVariableOrPanic("EMAIL_OTP_SECRET"),
	}
}

type MongoClientConfig struct {
	Username string
	Password string
	Url      string
}

func (config *MongoClientConfig) String() string {
	return fmt.Sprintf("Username: %s, Password: <masked>, Url: %s", config.Username, config.Url)
}

// Parse MongoConfig and return a connection string for use with the MongoDB driver
func (config *MongoClientConfig) GetMongoConnectionString() string {
	return fmt.Sprintf("mongodb+srv://%s:%s@%s", config.Username, config.Password, config.Url)
}

// Load MongoDB configuration from environment variables. Panics if any are missing.
func loadMongoVars() *MongoClientConfig {
	mongoUserName := getEnvironmentVariableOrPanic("MONGO_USERNAME")
	mongoPassword := getEnvironmentVariableOrPanic("MONGO_PASSWORD")
	mongoUrl := getEnvironmentVariableOrPanic("MONGO_URL")

	return &MongoClientConfig{
		Username: mongoUserName,
		Password: mongoPassword,
		Url:      mongoUrl,
	}

}

type SmtpClientConfig struct {
	Enabled     bool
	Username    string
	DisplayName string
	Password    string
	Host        string
	Port        uint16
}

func (config *SmtpClientConfig) String() string {
	return fmt.Sprintf("Enabled: %t, Username: %s, DisplayName: %s, Password: <masked>, Host: %s, Port: %d", config.Enabled, config.Username, config.DisplayName, config.Host, config.Port)
}

// Loads SMTP configuration from environment variables. If SMTP_ENABLED is set to false, the SMTP configuration is not required and will not be loaded.
func loadSmtpVars() *SmtpClientConfig {
	smtpEnabled := os.Getenv("SMTP_ENABLED") == "true"

	if !smtpEnabled {
		return &SmtpClientConfig{}
	}

	smtpUserName := getEnvironmentVariableOrPanic("SMTP_USERNAME")
	smtpDisplayName := getEnvironmentVariableOrPanic("SMTP_DISPLAY_NAME")
	smtpPassword := getEnvironmentVariableOrPanic("SMTP_PASSWORD")
	smtpHost := getEnvironmentVariableOrPanic("SMTP_HOST")
	smtpPortStr := getEnvironmentVariableOrPanic("SMTP_PORT")
	var smtpPort uint16
	if smtpEnabled {
		s, err := strconv.ParseUint(smtpPortStr, 10, 16)
		if err != nil {
			log.Panicf("SMTP_PORT must be an integer, got %s\nError: %v", smtpPortStr, err)
		}
		smtpPort = uint16(s)
	}

	return &SmtpClientConfig{
		Enabled:     smtpEnabled,
		Username:    smtpUserName,
		DisplayName: smtpDisplayName,
		Password:    smtpPassword,
		Host:        smtpHost,
		Port:        smtpPort,
	}
}

// Panics if the environment variable is not set
func getEnvironmentVariableOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Panicf("Missing required environment variable %s", key)
	}
	return value
}
