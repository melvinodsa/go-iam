package resource

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
	"github.com/stretchr/testify/require"
)

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

	t.Run("fetch resource successfully", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Get", mock.Anything, "0001").Return(&sdk.Resource{
			ID:   "0001",
			Name: "test",
		}, nil).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("GET", "/resource/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "0001", resp.Data.ID)
	})

	t.Run("resource not found", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Get", mock.Anything, "0001").Return(&sdk.Resource{}, sdk.ErrResourceNotFound).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("GET", "/resource/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("internal error", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Get", mock.Anything, "0001").Return(&sdk.Resource{}, errors.New("some error")).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("GET", "/resource/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

}

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

	t.Run("create resource successfully", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("POST", "/resource/v1", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error in creating resource", func(t *testing.T) {
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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Create", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("POST", "/resource/v1", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("bad request", func(t *testing.T) {
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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("POST", "/resource/v1", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestSearch(t *testing.T) {

	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("fetch all resources successfully", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Search", mock.Anything, mock.Anything).Return(&sdk.ResourceList{Resources: []sdk.Resource{
			{
				ID:   "0001",
				Name: "Test Resource 1",
			},
			{
				ID:   "0002",
				Name: "Test Resource 2",
			},
		}}, nil).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("GET", "/resource/v1/search?skip=0&limit=10", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ResourcesResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Len(t, resp.Data.Resources, 2)
	})

	t.Run("error in fetching resources", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Search", mock.Anything, mock.Anything).Return(&sdk.ResourceList{}, errors.New("some error")).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("GET", "/resource/v1/search?skip=0&limit=10", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
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

	t.Run("update resource successfully", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")
		req, _ := http.NewRequest("PUT", "/resource/v1/0001", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error in updating resource", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Update", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("PUT", "/resource/v1/0001", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("resource not found", func(t *testing.T) {
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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Update", mock.Anything, mock.Anything).Return(sdk.ErrResourceNotFound).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("PUT", "/resource/v1/0001", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("bad request", func(t *testing.T) {

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
		// resource mock

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("PUT", "/resource/v1/0001", strings.NewReader(`{
			"name": "Test Resource",
			"key": "test"
		}`))
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestDelete(t *testing.T) {

	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("delete resource successfully", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")
		req, _ := http.NewRequest("DELETE", "/resource/v1/0001", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("error in updating resource", func(t *testing.T) {

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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Delete", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("DELETE", "/resource/v1/0001", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("resource not found", func(t *testing.T) {
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
		// resource mock

		mockResourceSvc := services.MockResourceService{}
		mockResourceSvc.On("Delete", mock.Anything, mock.Anything).Return(sdk.ErrResourceNotFound).Once()

		svcs.Resources = &mockResourceSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/resource")

		req, _ := http.NewRequest("DELETE", "/resource/v1/0001", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.ResourceResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
