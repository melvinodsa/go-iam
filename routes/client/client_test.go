package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/config"
	"github.com/melvinodsa/go-iam/providers"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/cache"
	"github.com/melvinodsa/go-iam/utils/test"
	"github.com/melvinodsa/go-iam/utils/test/server"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	// setup
	os.Setenv("JWT_SECRET", "abcd")
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("create client successfully", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("POST", "/client/v1", strings.NewReader(`{
			"name": "Test Client",
			"description": "test client",
			"default_auth_provider_id": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("create client bad request", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("POST", "/client/v1", strings.NewReader(`{
			"name": 123,
			"description": "test client"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("service account or auth provider not provided", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("POST", "/client/v1", strings.NewReader(`{
			"name": "Test Client",
			"description": "test client"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("error creating client", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("POST", "/client/v1", strings.NewReader(`{
			"name": "Test Client",
			"description": "test client",
			"default_auth_provider_id": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestGet(t *testing.T) {
	// setup
	os.Setenv("JWT_SECRET", "abcd")
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("get client successfully", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(&sdk.Client{
			Id:          "client-id",
			Name:        "Test Client",
			Description: "test client",
		}, nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("GET", "/client/v1/client-id", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("client not found", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, sdk.ErrClientNotFound).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("GET", "/client/v1/client-id", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("error getting client", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("GET", "/client/v1/client-id", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestFetchhAll(t *testing.T) {
	// setup
	os.Setenv("JWT_SECRET", "abcd")
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("fetch all clients successfully", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("GetAll", mock.Anything, mock.Anything).Return([]sdk.Client{
			{
				Id:          "client-id-1",
				Name:        "Test Client 1",
				Description: "test client 1",
			},
			{
				Id:          "client-id-2",
				Name:        "Test Client 2",
				Description: "test client 2",
			},
		}, nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("GET", "/client/v1", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ClientsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, 2, len(resp.Data))
	})

	t.Run("error fetching all clients", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("GetAll", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("GET", "/client/v1", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ClientsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestUpdate(t *testing.T) {
	// setup
	os.Setenv("JWT_SECRET", "abcd")
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("update client successfully", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id", strings.NewReader(`{
			"name": "Updated Client",
			"description": "updated client",
			"default_auth_provider_id": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("update client bad request", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id", strings.NewReader(`{
			"name": 123,
			"description": "updated client"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("service account or auth provider not provided", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id", strings.NewReader(`{
			"name": "Updated Client",
			"description": "updated client"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("error updating client", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error")).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id", strings.NewReader(`{
			"name": "Updated Client",
			"description": "updated client",
			"default_auth_provider_id": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("client not found", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(sdk.ErrClientNotFound).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id", strings.NewReader(`{
			"name": "Updated Client",
			"description": "updated client",
			"default_auth_provider_id": "google"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestRegenerateSecret(t *testing.T) {
	// setup
	os.Setenv("JWT_SECRET", "abcd")
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("regenerate client secret successfully", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("RegenerateSecret", mock.Anything, mock.Anything).Return(&sdk.Client{
			Id:          "client-id",
			Name:        "Test Client",
			Description: "test client",
		}, nil).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id/regenerate-secret", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error regenerating client secret", func(t *testing.T) {
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
		// client mock

		mockClientSvc := services.MockClientService{}
		mockClientSvc.On("RegenerateSecret", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
		mockClientSvc.On("GetGoIamClients", mock.Anything, mock.Anything).Return([]sdk.Client{}, nil)
		mockClientSvc.On("Subscribe", mock.Anything, mock.Anything).Return()

		svcs.Clients = &mockClientSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/client")

		req, _ := http.NewRequest("PUT", "/client/v1/client-id/regenerate-secret", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ClientResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}
