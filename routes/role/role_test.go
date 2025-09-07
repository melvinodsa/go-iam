package role

import (
	"encoding/json"
	"errors"
	"net/http"
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
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("create role successfully", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("POST", "/role/v1", strings.NewReader(`{
			"name": "Test Role",
			"description": "test role"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("create role error", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("Create", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("POST", "/role/v1", strings.NewReader(`{
			"name": "Test Role",
			"description": "test role"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("create role invalid request body", func(t *testing.T) {
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

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("POST", "/role/v1", strings.NewReader(`{
			"name": 123,
			"description": "test role"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestSearch(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("search roles successfully", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("GetAll", mock.Anything, mock.Anything).Return(&sdk.RoleList{
			Roles: []sdk.Role{
				{
					Id:          "role1",
					Name:        "Role 1",
					Description: "First role",
				},
				{
					Id:          "role2",
					Name:        "Role 2",
					Description: "Second role",
				},
			},
			Total: 2,
			Skip:  0,
			Limit: 10,
		}, nil).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("GET", "/role/v1?skip=0&limit=10", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.RoleListResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, int64(2), resp.Data.Total)
	})

	t.Run("search roles error", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("GetAll", mock.Anything, mock.Anything).Return(&sdk.RoleList{}, errors.New("some error")).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("GET", "/role/v1?skip=0&limit=10", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.RoleListResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestGet(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("get role successfully", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("GetById", mock.Anything, mock.Anything).Return(&sdk.Role{
			Id:          "role1",
			Name:        "Role 1",
			Description: "First role",
		}, nil).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("GET", "/role/v1/role1", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "role1", resp.Data.Id)
	})

	t.Run("get role error", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("GetById", mock.Anything, mock.Anything).Return(&sdk.Role{}, errors.New("some error")).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("GET", "/role/v1/role1", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("role not found", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("GetById", mock.Anything, mock.Anything).Return(&sdk.Role{}, sdk.ErrRoleNotFound).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("GET", "/role/v1/role1", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestUpdate(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("update role successfully", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("PUT", "/role/v1/role1", strings.NewReader(`{
			"name": "Updated Role",
			"description": "Updated description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("update role error", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("Update", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("PUT", "/role/v1/role1", strings.NewReader(`{
			"name": "Updated Role",
			"description": "Updated description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update role not found", func(t *testing.T) {
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
		// role mock

		mockRoleSvc := services.MockRoleService{}
		mockRoleSvc.On("Update", mock.Anything, mock.Anything).Return(sdk.ErrRoleNotFound).Once()

		svcs.Role = &mockRoleSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("PUT", "/role/v1/role1", strings.NewReader(`{
			"name": "Updated Role",
			"description": "Updated description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update role bad payload error", func(t *testing.T) {
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
		// role mock

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/role")

		req, _ := http.NewRequest("PUT", "/role/v1/role1", strings.NewReader(`{
			"name": 23,
			"description": "Updated description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.RoleResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
