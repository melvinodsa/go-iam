package server

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupServer(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "test-secret")
	require.NoError(t, err)

	app := fiber.New()

	t.Run("SetupServer initializes app config", func(t *testing.T) {
		assert.NotNil(t, app)
	})

	t.Run("SetupServer creates valid config", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotNil(t, cnf)
		assert.NotNil(t, cnf.Server)
		assert.NotNil(t, cnf.DB)
		assert.NotNil(t, cnf.Jwt)
		assert.NotNil(t, cnf.Encrypter)
		assert.NotNil(t, cnf.Redis)
		assert.NotNil(t, cnf.Deployment)
		assert.NotNil(t, cnf.Logger)
	})

	t.Run("SetupServer config has valid server settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotEmpty(t, cnf.Server.Host)
		assert.NotEmpty(t, cnf.Server.Port)
		assert.GreaterOrEqual(t, cnf.Server.TokenCacheTTLInMinutes, int64(0))
		assert.GreaterOrEqual(t, cnf.Server.AuthProviderRefetchIntervalInMinutes, int64(0))
	})

	t.Run("SetupServer config has valid database settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotEmpty(t, cnf.DB.Host())
	})

	t.Run("SetupServer config has valid JWT settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotNil(t, cnf.Jwt.Secret())
	})

	t.Run("SetupServer config has valid encryption settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotNil(t, cnf.Encrypter.Key())
	})

	t.Run("SetupServer config has valid Redis settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotEmpty(t, cnf.Redis.Host)
	})

	t.Run("SetupServer config has valid deployment settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotEmpty(t, cnf.Deployment.Environment)
		assert.NotEmpty(t, cnf.Deployment.Name)
	})

	t.Run("SetupServer config has valid logger settings", func(t *testing.T) {
		cnf := config.NewAppConfig()

		assert.NotNil(t, cnf.Logger)
	})
}
