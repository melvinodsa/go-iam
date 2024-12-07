package config

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/melvinodsa/go-iam/sdk"
)

type AppConfig struct {
	Server     Server
	Deployment Deployment
	Logger     Logger
	DB         DB
	Encrypter  Encrypter
	Redis      Redis
}

func NewAppConfig() *AppConfig {
	cnf := &AppConfig{}
	cnf.Load()
	return cnf
}

type keyType struct {
	key string
}

var configKey = keyType{"config"}

func (a *AppConfig) Handle(c *fiber.Ctx) error {
	c.Locals(configKey, a)
	return c.Next()
}

func GetAppConfig(c *fiber.Ctx) AppConfig {
	return c.Locals(configKey).(AppConfig)
}

func (a *AppConfig) Load() {
	/*
	 * load env file
	 * load each config one by one
	 */
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found. Using default environment values")
	}
	a.LoadServerConfig()
	a.LoadDeploymentConfig()
	a.LoadLoggerConfig()
	a.LoadDBConfig()
	a.LoadEncrypterConfig()
	a.LoadRedisConfig()
}

// LoadServerConfig load server config
func (a *AppConfig) LoadServerConfig() {
	// load the default values
	// then load from env variables
	a.Server.Host = "localhost"
	a.Server.Port = "3000"

	host := os.Getenv("SERVER_HOST")
	if host != "" {
		a.Server.Host = host
	}
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		a.Server.Port = port
	}
}

// LoadDeploymentConfig loads the deployment config
func (a *AppConfig) LoadDeploymentConfig() {
	// load the default values
	// then load from env variables
	a.Deployment.Environment = "development"
	a.Deployment.Name = "Cuttle.ai Demo"

	environment := os.Getenv("DEPLOYMENT_ENVIRONMENT")
	if environment != "" {
		a.Deployment.Environment = environment
	}

	name := os.Getenv("DEPLOYMENT_NAME")
	if name != "" {
		a.Deployment.Name = name
	}
}

// LoadLoggerConfig loads logger config
func (a *AppConfig) LoadLoggerConfig() {
	// load the default values
	// then load from env variables
	level := log.LevelInfo

	levelStr := os.Getenv("LOGGER_LEVEL")
	if levelStr != "" {
		lvl, err := strconv.Atoi(levelStr)
		if err == nil {
			level = log.Level(lvl)
		}
	}

	lg := NewLogger(level)
	a.Logger = *lg
}

// LoadDBConfig loads db config
func (a *AppConfig) LoadDBConfig() {
	// load the default values
	// then load from env variables
	a.DB.host = "mongodb://test:test@127.0.0.1"
	host := os.Getenv("DB_HOST")
	if host != "" {
		a.DB.host = host
	}
}

func (a *AppConfig) LoadEncrypterConfig() {
	// load the default values
	// then load from env variables
	defaultKeyStr := "0000000000000000000000000000000000000000000000000000000000000000"
	keyStr := os.Getenv("ENCRYPTER_KEY")
	if keyStr != "" {
		defaultKeyStr = keyStr
	}
	key, err := hex.DecodeString(defaultKeyStr)
	if err != nil {
		panic(fmt.Errorf("error decoding encrypter key: %w", err))
	}
	//goland:noinspection GoRedundantConversion
	a.Encrypter.key = sdk.MaskedBytes(key)
}

func (a *AppConfig) LoadRedisConfig() {
	// load the default values
	// then load from env variables
	a.Redis.Host = "localhost:6379"
	host := os.Getenv("REDIS_HOST")
	if host != "" {
		a.Redis.Host = host
	}
	a.Redis.DB = 0
	dbStr := os.Getenv("REDIS_DB")
	if dbStr != "" {
		db, err := strconv.Atoi(dbStr)
		if err == nil {
			a.Redis.DB = db
		} else {
			panic(fmt.Errorf("error converting redis db to int: %w", err))
		}
	}

	password := os.Getenv("REDIS_PASSWORD")
	if password != "" {
		//goland:noinspection GoRedundantConversion
		a.Redis.Password = sdk.MaskedBytes([]byte(password))
	}
}
