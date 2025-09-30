package server

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupServer(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "test-secret")
	require.NoError(t, err)

	app := fiber.New()

	t.Run("SetupServer initializes app config", func(t *testing.T) {
		cnf := SetupServer(app)
		assert.NotNil(t, cnf)
		assert.NotNil(t, cnf.Server)
		assert.NotNil(t, cnf.DB)
		assert.NotNil(t, cnf.Jwt)
		assert.NotNil(t, cnf.Encrypter)
		assert.NotNil(t, cnf.Redis)
		assert.NotNil(t, cnf.Deployment)
		assert.NotNil(t, cnf.Logger)
	})
}