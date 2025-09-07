package user

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

func TestGetById(t *testing.T) {

	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("fetch user successfully", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("GetById", mock.Anything, "0001").Return(&sdk.User{
			Id:    "0001",
			Email: "",
		}, nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "0001", resp.Data.Id)
	})

	t.Run("user not found", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("GetById", mock.Anything, "0001").Return(&sdk.User{}, sdk.ErrUserNotFound).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("GetById", mock.Anything, "0001").Return(&sdk.User{}, errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1/0001", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

}

func TestCreate(t *testing.T) {

	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("create user successfully", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("POST", "/user/v1", strings.NewReader(`{
			"name": "Test User",
			"email": "testuser@example.com"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Equal(t, "testuser@example.com", resp.Data.Email)
	})

	t.Run("error in creating user", func(t *testing.T) {
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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("Create", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("POST", "/user/v1", strings.NewReader(`{
			"name": "Test User",
			"email": "testuser@example.com"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("POST", "/user/v1", strings.NewReader(`{
			"name": "Test User",
			"email": "not-an-email"
		}`))
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestGetAll(t *testing.T) {

	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("fetch all users successfully", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("GetAll", mock.Anything, mock.Anything).Return(&sdk.UserList{Users: []sdk.User{
			{
				Id:    "0001",
				Email: "",
			},
			{
				Id:    "0002",
				Email: "",
			},
		}}, nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1?skip=0&limit=10", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.UserListResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
		assert.Len(t, resp.Data.Users, 2)
	})

	t.Run("error in fetching users", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("GetAll", mock.Anything, mock.Anything).Return(&sdk.UserList{}, errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("GET", "/user/v1?skip=0&limit=10", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
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

	t.Run("update user successfully", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("Update", mock.Anything, mock.Anything).Return(nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001", strings.NewReader(`{
			"name": "Updated User",
			"email": "updated@example.com"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error in updating user", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("Update", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("PUT", "/user/v1/0001", strings.NewReader(`{
			"name": "Updated User",
			"email": "updated@example.com"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})

	t.Run("user not found", func(t *testing.T) {
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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("Update", mock.Anything, mock.Anything).Return(sdk.ErrUserNotFound).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("PUT", "/user/v1/0001", strings.NewReader(`{
			"name": "Updated User",
			"email": "updated@example.com"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
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
		// user mock

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")

		req, _ := http.NewRequest("PUT", "/user/v1/0001", strings.NewReader(`{
			"name": "Updated User",
			"email": "updated@example.com"
		}`))
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Nil(t, resp.Data)
	})
}

func TestUpdateRoles(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("update user role successfully", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemoveRoleFromUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUserSvc.On("AddRoleToUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/roles", strings.NewReader(`{
			"to_be_added": ["0001"],
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user role not found while removing", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemoveRoleFromUser", mock.Anything, mock.Anything, mock.Anything).Return(sdk.ErrRoleNotFound).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/roles", strings.NewReader(`{
			"to_be_added": ["0001"],
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user role error while removing role", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemoveRoleFromUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/roles", strings.NewReader(`{
			"to_be_added": ["0001"],
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user role not found while adding", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemoveRoleFromUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUserSvc.On("AddRoleToUser", mock.Anything, mock.Anything, mock.Anything).Return(sdk.ErrRoleNotFound).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/roles", strings.NewReader(`{
			"to_be_added": ["0001"],
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user role error while adding", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemoveRoleFromUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUserSvc.On("AddRoleToUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/roles", strings.NewReader(`{
			"to_be_added": ["0001"],
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
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

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/roles", strings.NewReader(`{
			"to_be_added": ["0001"],
			"to_be_removed": ["0002"]
		}`))
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestUpdatePolicies(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("update user policy successfully", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemovePolicyFromUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUserSvc.On("AddPolicyToUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/policies", strings.NewReader(`{
			"to_be_added": {"0001": {"name": "0001", "mapping": {"firstArg": {"arguments": {"static": "1234"}}}}},
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user policy user not found while removing", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemovePolicyFromUser", mock.Anything, mock.Anything, mock.Anything).Return(sdk.ErrUserNotFound).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/policies", strings.NewReader(`{
			"to_be_added": {"0001": {"name": "0001", "mapping": {"firstArg": {"arguments": {"static": "1234"}}}}},
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user policy error while removing policy", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemovePolicyFromUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/policies", strings.NewReader(`{
			"to_be_added": {"0001": {"name": "0001", "mapping": {"firstArg": {"arguments": {"static": "1234"}}}}},
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user policy user not found while adding", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemovePolicyFromUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUserSvc.On("AddPolicyToUser", mock.Anything, mock.Anything, mock.Anything).Return(sdk.ErrUserNotFound).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/policies", strings.NewReader(`{
			"to_be_added": {"0001": {"name": "0001", "mapping": {"firstArg": {"arguments": {"static": "1234"}}}}},
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("update user policy error while adding", func(t *testing.T) {

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
		// user mock

		mockUserSvc := services.MockUserService{}
		mockUserSvc.On("RemovePolicyFromUser", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockUserSvc.On("AddPolicyToUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error")).Once()

		svcs.User = &mockUserSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/policies", strings.NewReader(`{
			"to_be_added": {"0001": {"name": "0001", "mapping": {"firstArg": {"arguments": {"static": "1234"}}}}},
			"to_be_removed": ["0002"]
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
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

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/user")
		req, _ := http.NewRequest("PUT", "/user/v1/0001/policies", strings.NewReader(`{
			"to_be_added": {"0001": {"name": "0001", "mapping": {"firstArg": {"arguments": {"static": "1234"}}}}},
			"to_be_removed": ["0002"]
		}`))
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.UserResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
