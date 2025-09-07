package project

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

func TestProject(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("create project successfully", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("POST", "/project/v1", strings.NewReader(`{
			"name": "Test Project",
			"description": "test project"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 201, res.StatusCode, "Expected status code 201")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error creating project", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("Create", mock.Anything, mock.Anything).Return(errors.New("project already exists")).Once()
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("POST", "/project/v1", strings.NewReader(`{
			"name": "Test Project",
			"description": "test project"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("bad payload", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("POST", "/project/v1", strings.NewReader(`{
			"name": 123,
			"description": "test project"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
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

	t.Run("get project successfully", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("Get", mock.Anything, "project-id").Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil).Once()

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("GET", "/project/v1/project-id", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error getting project", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("Get", mock.Anything, "project-id").Return(&sdk.Project{}, errors.New("project not found")).Once()

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("GET", "/project/v1/project-id", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("project not found", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("Get", mock.Anything, "project-id").Return(&sdk.Project{}, sdk.ErrProjectNotFound).Once()

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("GET", "/project/v1/project-id", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}

func TestFetchAll(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("fetch all projects successfully", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("GetAll", mock.Anything).Return([]sdk.Project{
			{
				Id:          "project-id",
				Name:        "Test Project",
				Description: "test project",
			},
		}, nil).Once()

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("GET", "/project/v1", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ProjectsResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error fetching all projects", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("GetAll", mock.Anything).Return([]sdk.Project{}, errors.New("database error")).Once()

		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("GET", "/project/v1", nil)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ProjectsResponse
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

	t.Run("update project successfully", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("PUT", "/project/v1/project-id", strings.NewReader(`{
			"name": "Updated Project",
			"description": "updated project description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Data)
	})

	t.Run("error updating project", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed")).Once()
		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("PUT", "/project/v1/project-id", strings.NewReader(`{
			"name": "Updated Project",
			"description": "updated project description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 500, res.StatusCode, "Expected status code 500")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("bad payload", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("PUT", "/project/v1/project-id", strings.NewReader(`{
			"name": 123,
			"description": "updated project description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 400, res.StatusCode, "Expected status code 400")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("project not found", func(t *testing.T) {
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
		// project mock

		mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		mockProjectSvc.On("Update", mock.Anything, mock.Anything).Return(sdk.ErrProjectNotFound).Once()
		svcs.Projects = &mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)

		app.Use(providers.Handle(prv))

		RegisterRoutes(app, "/project")

		req, _ := http.NewRequest("PUT", "/project/v1/project-id", strings.NewReader(`{
			"name": "Updated Project",
			"description": "updated project description"
		}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)
		assert.Equalf(t, 404, res.StatusCode, "Expected status code 404")
		assert.Nil(t, err)
		var resp sdk.ProjectResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
	})
}
