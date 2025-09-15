// Package config provides configuration management for the Go IAM API server.
// It handles loading configuration from environment variables and .env files
// for all application components including server, database, encryption, JWT, and more.
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

// AppConfig holds all configuration settings for the Go IAM application.
// It includes settings for server, deployment, logging, database, encryption,
// Redis, JWT, and service account configurations.
type AppConfig struct {
	Server         Server         // HTTP server configuration
	Deployment     Deployment     // Deployment environment settings
	Logger         Logger         // Logging configuration
	DB             DB             // Database connection settings
	Encrypter      Encrypter      // Encryption key configuration
	Redis          Redis          // Redis cache configuration
	Jwt            Jwt            // JWT token configuration
	ServiceAccount ServiceAccount // Service account token settings
}

// NewAppConfig creates a new AppConfig instance and loads all configuration
// from environment variables and .env files. This is the primary entry point
// for initializing application configuration.
func NewAppConfig() *AppConfig {
	cnf := &AppConfig{}
	cnf.Load()
	return cnf
}

type keyType struct {
	key string
}

var configKey = keyType{"config"}

// Handle is a Fiber middleware that stores the AppConfig in the request context.
// This allows handlers to access configuration using GetAppConfig().
//
// Usage:
//
//	app.Use(config.Handle)
func (a *AppConfig) Handle(c *fiber.Ctx) error {
	c.Locals(configKey, *a)
	return c.Next()
}

// GetAppConfig retrieves the AppConfig from the Fiber context.
// This function should be called from handlers that need access to configuration.
// The config must have been previously stored using the Handle middleware.
//
// Returns the AppConfig instance stored in the context.
func GetAppConfig(c *fiber.Ctx) AppConfig {
	return c.Locals(configKey).(AppConfig)
}

// Load reads configuration from environment variables and .env files.
// It loads all configuration sections including server, deployment, logger,
// database, encrypter, and Redis settings. This method is called automatically
// by NewAppConfig() and should not typically be called directly.
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
	a.LoadJwtConfig()
	a.LoadServiceAccountConfig()
}

// LoadServerConfig loads server-specific configuration from environment variables.
// It sets up HTTP server host, port, Redis enablement, token cache TTL,
// and auth provider refetch interval settings.
//
// Environment variables:
//   - SERVER_HOST: HTTP server host (default: localhost)
//   - SERVER_PORT: HTTP server port (default: 3000)
//   - ENABLE_REDIS: Enable Redis caching (default: false)
//   - TOKEN_CACHE_TTL_IN_MINUTES: Token cache TTL in minutes (default: 1440)
//   - AUTH_PROVIDER_REFETCH_INTERVAL_IN_MINUTES: Auth provider refresh interval (default: 1)
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
	enableRedis := os.Getenv("ENABLE_REDIS")
	if enableRedis == "true" {
		a.Server.EnableRedis = true
	}
	tokenCacheTTL := os.Getenv("TOKEN_CACHE_TTL_IN_MINUTES")
	if tokenCacheTTL != "" {
		ttl, err := strconv.ParseInt(tokenCacheTTL, 10, 64)
		if err == nil {
			a.Server.TokenCacheTTLInMinutes = ttl
		} else {
			panic(fmt.Errorf("error converting token cache ttl to int: %w", err))
		}
	} else {
		a.Server.TokenCacheTTLInMinutes = 1440 // default to 1440 minutes - 24 hours
	}
	authProviderRefetchInterval := os.Getenv("AUTH_PROVIDER_REFETCH_INTERVAL_IN_MINUTES")
	if authProviderRefetchInterval != "" {
		interval, err := strconv.ParseInt(authProviderRefetchInterval, 10, 64)
		if err == nil {
			a.Server.AuthProviderRefetchIntervalInMinutes = interval
		} else {
			panic(fmt.Errorf("error converting auth provider refetch interval to int: %w", err))
		}
	} else {
		a.Server.AuthProviderRefetchIntervalInMinutes = 1 // default to 1 minute
	}
	log.Infow("Loaded Server Configurations",
		"host", a.Server.Host,
		"port", a.Server.Port,
		"enable_redis", a.Server.EnableRedis,
		"token_cache_ttl", a.Server.TokenCacheTTLInMinutes,
	)
}

// LoadDeploymentConfig loads deployment environment configuration from environment variables.
// It configures the deployment environment and application name for identification.
//
// Environment variables:
//   - DEPLOYMENT_ENVIRONMENT: Deployment environment name (default: development)
//   - DEPLOYMENT_NAME: Application deployment name (default: Cuttle.ai Demo)
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

// LoadLoggerConfig loads logging configuration from environment variables.
// It sets up the global logger with the specified log level.
//
// Environment variables:
//   - LOGGER_LEVEL: Log level as integer (default: Info level)
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

// LoadDBConfig loads database configuration from environment variables.
// It configures the MongoDB connection string.
//
// Environment variables:
//   - DB_HOST: MongoDB connection string (default: mongodb://test:test@127.0.0.1)
func (a *AppConfig) LoadDBConfig() {
	// load the default values
	// then load from env variables
	a.DB.host = "mongodb://test:test@127.0.0.1"
	host := os.Getenv("DB_HOST")
	if host != "" {
		a.DB.host = host
	}
}

// LoadEncrypterConfig loads encryption configuration from environment variables.
// It sets up the encryption key used for sensitive data encryption.
// The key must be a valid hex-encoded string.
//
// Environment variables:
//   - ENCRYPTER_KEY: Hex-encoded encryption key (default: 64-character zero string)
//
// Panics if the encryption key cannot be decoded from hex.
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

// LoadRedisConfig loads Redis configuration from environment variables.
// It configures Redis connection settings including host, database number, and password.
//
// Environment variables:
//   - REDIS_HOST: Redis server address (default: localhost:6379)
//   - REDIS_DB: Redis database number (default: 0)
//   - REDIS_PASSWORD: Redis password (optional)
//
// Panics if REDIS_DB cannot be converted to integer.
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

// LoadJwtConfig loads JWT configuration from environment variables.
// It sets up the JWT secret key used for token signing and verification.
//
// Environment variables:
//   - JWT_SECRET: JWT secret key (required)
//
// Panics if JWT_SECRET is not provided.
func (a *AppConfig) LoadJwtConfig() {
	// load from env variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is required")
	}
	a.Jwt.secret = sdk.MaskedBytes(secret)
}

// LoadServiceAccountConfig loads service account configuration from environment variables.
// It configures token TTL settings for service account access and refresh tokens.
//
// Environment variables:
//   - SERVICE_ACCOUNT_ACCESS_TOKEN_TTL_MINUTES: Access token TTL in minutes (default: 60)
//   - SERVICE_ACCOUNT_REFRESH_TOKEN_TTL_DAYS: Refresh token TTL in days (default: 30)
func (a *AppConfig) LoadServiceAccountConfig() {
	a.ServiceAccount.AccessTokenTTLInMinutes = 60 // Default 1 hour
	a.ServiceAccount.RefreshTokenTTLInDays = 30   // Default 30 days

	if val := os.Getenv("SERVICE_ACCOUNT_ACCESS_TOKEN_TTL_MINUTES"); val != "" {
		if ttl, err := strconv.ParseInt(val, 10, 64); err == nil {
			a.ServiceAccount.AccessTokenTTLInMinutes = ttl
		}
	}

	if val := os.Getenv("SERVICE_ACCOUNT_REFRESH_TOKEN_TTL_DAYS"); val != "" {
		if ttl, err := strconv.ParseInt(val, 10, 64); err == nil {
			a.ServiceAccount.RefreshTokenTTLInDays = ttl
		}
	}
}
