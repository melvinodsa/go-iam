package authprovider

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/melvinodsa/go-iam/utils/test/server"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("create auth provider successfully", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("POST", "/authprovider/v1", strings.NewReader(`{
			"name": "Test Auth Provider",
			"provider": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("invalid body", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("POST", "/authprovider/v1", strings.NewReader(`{
			"name": "Test Auth Provider",
			"provider": "google",
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})

	t.Run("service error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Create", mock.Anything, mock.Anything).Return(errors.New("some Error")).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("POST", "/authprovider/v1", strings.NewReader(`{
			"name": "Test Auth Provider",
			"provider": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})
}

func TestGet(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("get auth provider successfully", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Get", mock.Anything, "authprovider-123", false).Return(&sdk.AuthProvider{
			Id:        "authprovider-123",
			Name:      "Test Auth Provider",
			Provider:  "google",
			ProjectId: "project-123",
		}, nil).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("GET", "/authprovider/v1/authprovider-123", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("auth provider not found", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Get", mock.Anything, "authprovider-123", false).Return(&sdk.AuthProvider{}, sdk.ErrAuthProviderNotFound).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("GET", "/authprovider/v1/authprovider-123", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})

	t.Run("service error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Get", mock.Anything, "authprovider-123", false).Return(&sdk.AuthProvider{}, errors.New("some error")).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("GET", "/authprovider/v1/authprovider-123", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})
}

func TestFetchAll(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("fetch all auth providers successfully", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("GetAll", mock.Anything, mock.Anything).Return([]sdk.AuthProvider{
			{
				Id:        "authprovider-123",
				Name:      "Test Auth Provider",
				Provider:  "google",
				ProjectId: "project-123",
			},
			{
				Id:        "authprovider-456",
				Name:      "Test Auth Provider 2",
				Provider:  "github",
				ProjectId: "project-123",
			},
		}, nil).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("GET", "/authprovider/v1", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.AuthProvidersResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, 2, len(resp.Data))
	})

	t.Run("service error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("GetAll", mock.Anything, mock.Anything).Return([]sdk.AuthProvider{}, errors.New("some error")).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		req, _ := http.NewRequest("GET", "/authprovider/v1", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.AuthProvidersResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})
}

func TestUpdate(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("update auth provider successfully", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		id := "authprovider-" + uuid.New().String()
		req, _ := http.NewRequest("PUT", "/authprovider/v1/"+id, strings.NewReader(`{
			"name": "Test Auth Provider Updated",
			"provider": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("invalid body", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/authprovider")

		id := "authprovider-" + uuid.New().String()
		req, _ := http.NewRequest("PUT", "/authprovider/v1/"+id, strings.NewReader(`{
			"name": "Test Auth Provider Updated",
			"provider": "google",
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})

	t.Run("auth provider not found", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Update", mock.Anything, mock.Anything).Return(sdk.ErrAuthProviderNotFound).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/authprovider")

		id := "authprovider-" + uuid.New().String()
		req, _ := http.NewRequest("PUT", "/authprovider/v1/"+id, strings.NewReader(`{
			"name": "Test Auth Provider Updated",
			"provider": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})
	t.Run("service error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ReadBufferSize: 8192,
		})

		d := test.SetupMockDB()
		cs := cache.NewMockService()
		svcs, err := server.GetServices(*cnf, cs, d)
		if err != nil {
			t.Errorf("error getting services: %s", err)
			return
		}
		// auth provider mock

		mockAuthProviderSvc := services.MockAuthProviderService{}
		mockAuthProviderSvc.On("Update", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.AuthProviders = &mockAuthProviderSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/authprovider")

		id := "authprovider-" + uuid.New().String()
		req, _ := http.NewRequest("PUT", "/authprovider/v1/"+id, strings.NewReader(`{
			"name": "Test Auth Provider Updated",
			"provider": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.AuthProviderResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.False(t, resp.Success)
	})
}
