package config

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppConfig(t *testing.T) {
	// Clean environment before test
	cleanEnv()

	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	config := NewAppConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "3000", config.Server.Port)
	assert.False(t, config.Server.EnableRedis)
	assert.Equal(t, int64(1440), config.Server.TokenCacheTTLInMinutes)
	assert.Equal(t, int64(1), config.Server.AuthProviderRefetchIntervalInMinutes)
	assert.Equal(t, "development", config.Deployment.Environment)
	assert.Equal(t, "Go IAM Demo", config.Deployment.Name)
}

func TestAppConfig_Handle(t *testing.T) {
	config := &AppConfig{
		Server: Server{Host: "testhost", Port: "8080"},
	}
	app := fiber.New()

	app.Use(config.Handle)
	app.Get("/test", func(c *fiber.Ctx) error {
		storedConfig := GetAppConfig(c)
		return c.JSON(fiber.Map{
			"status": "ok",
			"host":   storedConfig.Server.Host,
			"port":   storedConfig.Server.Port,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAppConfig(t *testing.T) {
	config := &AppConfig{
		Server: Server{Host: "testhost", Port: "8080"},
	}

	app := fiber.New()
	app.Use(config.Handle)
	app.Get("/test", func(c *fiber.Ctx) error {
		retrievedConfig := GetAppConfig(c)
		return c.JSON(fiber.Map{
			"host": retrievedConfig.Server.Host,
			"port": retrievedConfig.Server.Port,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAppConfig_LoadServerConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Server
	}{
		{
			name:    "Default values",
			envVars: map[string]string{},
			expected: Server{
				Host:                                 "localhost",
				Port:                                 "3000",
				EnableRedis:                          false,
				TokenCacheTTLInMinutes:               1440,
				AuthProviderRefetchIntervalInMinutes: 1,
			},
		},
		{
			name: "Custom host and port",
			envVars: map[string]string{
				"SERVER_HOST": "0.0.0.0",
				"SERVER_PORT": "8080",
			},
			expected: Server{
				Host:                                 "0.0.0.0",
				Port:                                 "8080",
				EnableRedis:                          false,
				TokenCacheTTLInMinutes:               1440,
				AuthProviderRefetchIntervalInMinutes: 1,
			},
		},
		{
			name: "Enable Redis",
			envVars: map[string]string{
				"ENABLE_REDIS": "true",
			},
			expected: Server{
				Host:                                 "localhost",
				Port:                                 "3000",
				EnableRedis:                          true,
				TokenCacheTTLInMinutes:               1440,
				AuthProviderRefetchIntervalInMinutes: 1,
			},
		},
		{
			name: "Custom TTL and intervals",
			envVars: map[string]string{
				"TOKEN_CACHE_TTL_IN_MINUTES":                "720",
				"AUTH_PROVIDER_REFETCH_INTERVAL_IN_MINUTES": "5",
			},
			expected: Server{
				Host:                                 "localhost",
				Port:                                 "3000",
				EnableRedis:                          false,
				TokenCacheTTLInMinutes:               720,
				AuthProviderRefetchIntervalInMinutes: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv()
			setEnvVars(tt.envVars)
			defer cleanEnv()

			config := &AppConfig{}
			config.LoadServerConfig()

			assert.Equal(t, tt.expected, config.Server)
		})
	}
}

func TestAppConfig_LoadServerConfig_InvalidTTL(t *testing.T) {
	cleanEnv()
	err := os.Setenv("TOKEN_CACHE_TTL_IN_MINUTES", "invalid")
	assert.NoError(t, err)
	defer cleanEnv()

	config := &AppConfig{}

	assert.Panics(t, func() {
		config.LoadServerConfig()
	})
}

func TestAppConfig_LoadServerConfig_InvalidInterval(t *testing.T) {
	cleanEnv()
	err := os.Setenv("AUTH_PROVIDER_REFETCH_INTERVAL_IN_MINUTES", "invalid")
	assert.NoError(t, err)
	defer cleanEnv()

	config := &AppConfig{}

	assert.Panics(t, func() {
		config.LoadServerConfig()
	})
}

func TestAppConfig_LoadDeploymentConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Deployment
	}{
		{
			name:    "Default values",
			envVars: map[string]string{},
			expected: Deployment{
				Environment: "development",
				Name:        "Go IAM Demo",
			},
		},
		{
			name: "Custom values",
			envVars: map[string]string{
				"DEPLOYMENT_ENVIRONMENT": "production",
				"DEPLOYMENT_NAME":        "My App",
			},
			expected: Deployment{
				Environment: "production",
				Name:        "My App",
			},
		},
		{
			name: "Only environment",
			envVars: map[string]string{
				"DEPLOYMENT_ENVIRONMENT": "staging",
			},
			expected: Deployment{
				Environment: "staging",
				Name:        "Go IAM Demo",
			},
		},
		{
			name: "Only name",
			envVars: map[string]string{
				"DEPLOYMENT_NAME": "Test App",
			},
			expected: Deployment{
				Environment: "development",
				Name:        "Test App",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv()
			setEnvVars(tt.envVars)
			defer cleanEnv()

			config := &AppConfig{}
			config.LoadDeploymentConfig()

			assert.Equal(t, tt.expected, config.Deployment)
		})
	}
}

func TestAppConfig_LoadLoggerConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedLvl log.Level
	}{
		{
			name:        "Default level",
			envVars:     map[string]string{},
			expectedLvl: log.LevelInfo,
		},
		{
			name: "Debug level",
			envVars: map[string]string{
				"LOGGER_LEVEL": strconv.Itoa(int(log.LevelDebug)),
			},
			expectedLvl: log.LevelDebug,
		},
		{
			name: "Error level",
			envVars: map[string]string{
				"LOGGER_LEVEL": strconv.Itoa(int(log.LevelError)),
			},
			expectedLvl: log.LevelError,
		},
		{
			name: "Invalid level falls back to default",
			envVars: map[string]string{
				"LOGGER_LEVEL": "invalid",
			},
			expectedLvl: log.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv()
			setEnvVars(tt.envVars)
			defer cleanEnv()

			config := &AppConfig{}
			config.LoadLoggerConfig()

			assert.Equal(t, tt.expectedLvl, config.Logger.Level)
		})
	}
}

func TestAppConfig_LoadDBConfig(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedHost string
	}{
		{
			name:         "Default host",
			envVars:      map[string]string{},
			expectedHost: "mongodb://test:test@127.0.0.1",
		},
		{
			name: "Custom host",
			envVars: map[string]string{
				"DB_HOST": "mongodb://user:pass@production-db:27017",
			},
			expectedHost: "mongodb://user:pass@production-db:27017",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv()
			setEnvVars(tt.envVars)
			defer cleanEnv()

			config := &AppConfig{}
			config.LoadDBConfig()

			assert.Equal(t, tt.expectedHost, config.DB.Host())
		})
	}
}

func TestAppConfig_LoadEncrypterConfig(t *testing.T) {
	tests := []struct {
		name      string
		envVars   map[string]string
		shouldErr bool
	}{
		{
			name:      "Default key",
			envVars:   map[string]string{},
			shouldErr: false,
		},
		{
			name: "Custom valid key",
			envVars: map[string]string{
				"ENCRYPTER_KEY": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			},
			shouldErr: false,
		},
		{
			name: "Invalid hex key",
			envVars: map[string]string{
				"ENCRYPTER_KEY": "invalid-hex-string",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv()
			setEnvVars(tt.envVars)
			defer cleanEnv()

			config := &AppConfig{}

			if tt.shouldErr {
				assert.Panics(t, func() {
					config.LoadEncrypterConfig()
				})
			} else {
				assert.NotPanics(t, func() {
					config.LoadEncrypterConfig()
				})
				assert.NotNil(t, config.Encrypter.Key())
			}
		})
	}
}

func TestAppConfig_LoadRedisConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Redis
	}{
		{
			name:    "Default values",
			envVars: map[string]string{},
			expected: Redis{
				Host:     "localhost:6379",
				Password: nil,
				DB:       0,
			},
		},
		{
			name: "Custom host and DB",
			envVars: map[string]string{
				"REDIS_HOST": "redis-server:6379",
				"REDIS_DB":   "2",
			},
			expected: Redis{
				Host:     "redis-server:6379",
				Password: nil,
				DB:       2,
			},
		},
		{
			name: "With password",
			envVars: map[string]string{
				"REDIS_PASSWORD": "secret123",
			},
			expected: Redis{
				Host:     "localhost:6379",
				Password: sdk.MaskedBytes([]byte("secret123")),
				DB:       0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanEnv()
			setEnvVars(tt.envVars)
			defer cleanEnv()

			config := &AppConfig{}
			config.LoadRedisConfig()

			assert.Equal(t, tt.expected.Host, config.Redis.Host)
			assert.Equal(t, tt.expected.DB, config.Redis.DB)
			if tt.expected.Password == nil {
				assert.Nil(t, config.Redis.Password)
			} else {
				assert.Equal(t, tt.expected.Password, config.Redis.Password)
			}
		})
	}
}

func TestAppConfig_LoadRedisConfig_InvalidDB(t *testing.T) {
	cleanEnv()
	err := os.Setenv("REDIS_DB", "invalid")
	assert.NoError(t, err)
	defer cleanEnv()

	config := &AppConfig{}

	assert.Panics(t, func() {
		config.LoadRedisConfig()
	})
}

func TestAppConfig_LoadJwtConfig(t *testing.T) {
	t.Run("Valid JWT secret", func(t *testing.T) {
		cleanEnv()
		err := os.Setenv("JWT_SECRET", "my-secret-key")
		assert.NoError(t, err)
		defer cleanEnv()

		config := &AppConfig{}
		assert.NotPanics(t, func() {
			config.LoadJwtConfig()
		})
		assert.NotNil(t, config.Jwt.Secret())
		assert.Equal(t, sdk.MaskedBytes("my-secret-key"), config.Jwt.Secret())
	})

	t.Run("Missing JWT secret", func(t *testing.T) {
		cleanEnv()
		defer cleanEnv()

		config := &AppConfig{}
		assert.Panics(t, func() {
			config.LoadJwtConfig()
		})
	})
}

func TestAppConfig_Load(t *testing.T) {
	cleanEnv()

	// Set some environment variables
	envVars := map[string]string{
		"SERVER_HOST":            "testhost",
		"SERVER_PORT":            "9000",
		"DEPLOYMENT_ENVIRONMENT": "test",
		"DB_HOST":                "mongodb://testdb",
		"JWT_SECRET":             "abcd",
		"ENCRYPTER_KEY":          "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"REDIS_HOST":             "testredis:6379",
	}
	setEnvVars(envVars)
	defer cleanEnv()

	config := &AppConfig{}
	config.Load()

	// Verify all configs are loaded
	assert.Equal(t, "testhost", config.Server.Host)
	assert.Equal(t, "9000", config.Server.Port)
	assert.Equal(t, "test", config.Deployment.Environment)
	assert.Equal(t, "mongodb://testdb", config.DB.Host())
	assert.NotNil(t, config.Encrypter.Key())
	assert.Equal(t, "testredis:6379", config.Redis.Host)
	assert.Equal(t, log.LevelInfo, config.Logger.Level)
}

// Test individual config types

func TestDB_Host(t *testing.T) {
	db := DB{host: "test-host"}
	assert.Equal(t, "test-host", db.Host())
}

func TestEncrypter_Key(t *testing.T) {
	testKey := sdk.MaskedBytes([]byte("test-key"))
	encrypter := Encrypter{key: testKey}
	assert.Equal(t, testKey, encrypter.Key())
}

func TestJwt_Secret(t *testing.T) {
	testSecret := sdk.MaskedBytes("test-secret")
	jwt := Jwt{secret: testSecret}
	assert.Equal(t, testSecret, jwt.Secret())
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger(log.LevelDebug)
	assert.NotNil(t, logger)
	assert.Equal(t, log.LevelDebug, logger.Level)
}

func TestMapperConstants(t *testing.T) {
	assert.Equal(t, "resourceMap", ResourceMapContextKey)
	assert.Equal(t, "roleMap", RoleMapContextKey)
}

// Test edge cases and error conditions

func TestAppConfig_GetAppConfig_Panic(t *testing.T) {
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		// This should panic because no config was set
		assert.Panics(t, func() {
			GetAppConfig(c)
		})
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAppConfig_Integration(t *testing.T) {
	// Test full integration with environment loading
	cleanEnv()

	// Create a temporary .env file
	envContent := `SERVER_HOST=envhost
SERVER_PORT=8080
DEPLOYMENT_ENVIRONMENT=production
DB_HOST=mongodb://prod-db
ENCRYPTER_KEY=abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
REDIS_HOST=prod-redis:6379
REDIS_PASSWORD=prodpass
JWT_SECRET=prod-jwt-secret
`

	// Write temporary .env file
	err := os.WriteFile(".env", []byte(envContent), 0644)
	require.NoError(t, err)
	defer func() {
		err := os.Remove(".env")
		assert.NoError(t, err)
	}()
	defer cleanEnv()

	config := NewAppConfig()

	// Verify .env values were loaded
	assert.Equal(t, "envhost", config.Server.Host)
	assert.Equal(t, "8080", config.Server.Port)
	assert.Equal(t, "production", config.Deployment.Environment)
	assert.Equal(t, "mongodb://prod-db", config.DB.Host())
	assert.Equal(t, "prod-redis:6379", config.Redis.Host)
	assert.Equal(t, sdk.MaskedBytes([]byte("prodpass")), config.Redis.Password)
}

// Benchmark tests

func BenchmarkNewAppConfig(b *testing.B) {
	cleanEnv()
	defer cleanEnv()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewAppConfig()
	}
}

func BenchmarkAppConfig_Handle(b *testing.B) {
	config := &AppConfig{}
	app := fiber.New()
	app.Use(config.Handle)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		_, err := app.Test(req)
		assert.NoError(b, err)
	}
}

// Helper functions

func cleanEnv() {
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "ENABLE_REDIS",
		"TOKEN_CACHE_TTL_IN_MINUTES", "AUTH_PROVIDER_REFETCH_INTERVAL_IN_MINUTES",
		"DEPLOYMENT_ENVIRONMENT", "DEPLOYMENT_NAME",
		"LOGGER_LEVEL", "DB_HOST", "ENCRYPTER_KEY",
		"REDIS_HOST", "REDIS_DB", "REDIS_PASSWORD",
		"JWT_SECRET",
	}

	for _, env := range envVars {
		err := os.Unsetenv(env)
		if err != nil {
			log.Error(err)
		}
	}
}

func setEnvVars(envVars map[string]string) {
	for key, value := range envVars {
		err := os.Setenv(key, value)
		if err != nil {
			log.Error(err)
		}
	}
}
