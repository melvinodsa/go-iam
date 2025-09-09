package health

import (
	"encoding/json"
	"net/http"
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

func TestHealth(t *testing.T) {
	cnf := config.NewAppConfig()
	log.Infow("Loaded Configurations",
		"host", cnf.Server.Host,
		"port", cnf.Server.Port,
		"env", cnf.Deployment.Environment,
		"app_name", cnf.Deployment.Name,
	)

	t.Run("health check returns success", func(t *testing.T) {
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

		// Mock the project service to return empty projects for health check
		mockProjectSvc := &services.MockProjectService{}
		mockProjectSvc.On("GetAll", mock.Anything).Return([]sdk.Project{}, nil)
		// Mock GetByName to return that default project doesn't exist (which is fine for health check)mockProjectSvc := services.MockProjectService{}
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		svcs.Projects = mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/health")

		req, _ := http.NewRequest("GET", "/health/v1", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)

		var resp HealthResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		assert.Equal(t, "Health check completed successfully", resp.Message)
		assert.NotEmpty(t, resp.Timestamp)
		assert.NotNil(t, resp.Data)
		assert.Contains(t, []string{"healthy", "degraded"}, resp.Data.Status)
		assert.Equal(t, "1.0.0", resp.Data.Version)
		assert.NotEmpty(t, resp.Data.Uptime)
		assert.NotNil(t, resp.Data.Components)
		assert.Contains(t, resp.Data.Components, "database")
		assert.Contains(t, resp.Data.Components, "cache")
	})

	t.Run("health check returns correct component status", func(t *testing.T) {
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

		// Mock the project service to return empty projects for health check
		mockProjectSvc := &services.MockProjectService{}
		mockProjectSvc.On("GetAll", mock.Anything).Return([]sdk.Project{}, nil)
		// Mock GetByName to return that default project doesn't exist (which is fine for health check)
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		svcs.Projects = mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/health")

		req, _ := http.NewRequest("GET", "/health/v1", nil)
		res, err := app.Test(req, -1)
		assert.Equalf(t, 200, res.StatusCode, "Expected status code 200")
		assert.Nil(t, err)

		var resp HealthResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		// With mock services, components should be healthy
		assert.Equal(t, "healthy", resp.Data.Components["database"])
		assert.Equal(t, "healthy", resp.Data.Components["cache"])
		assert.Equal(t, "healthy", resp.Data.Status)
	})

	t.Run("health check handles different status codes", func(t *testing.T) {
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

		// Mock the project service to return empty projects for health check
		mockProjectSvc := &services.MockProjectService{}
		mockProjectSvc.On("GetAll", mock.Anything).Return([]sdk.Project{}, nil)
		// Mock GetByName to return that default project doesn't exist (which is fine for health check)
		mockProjectSvc.On("GetByName", mock.Anything, mock.Anything).Return(&sdk.Project{
			Id:          "project-id",
			Name:        "Test Project",
			Description: "test project",
		}, nil)
		svcs.Projects = mockProjectSvc

		prv := server.SetupTestServer(app, cnf, svcs, cs, d)
		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/health")

		req, _ := http.NewRequest("GET", "/health/v1", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)

		var resp HealthResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.NotNil(t, resp)

		// Status code should be 200 for healthy systems
		if resp.Data.Status == "healthy" {
			assert.Equal(t, 200, res.StatusCode)
		} else if resp.Data.Status == "degraded" {
			assert.Equal(t, 206, res.StatusCode) // StatusPartialContent
		} else if resp.Data.Status == "unhealthy" {
			assert.Equal(t, 503, res.StatusCode) // StatusServiceUnavailable
		}
	})
}
