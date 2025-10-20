package health

import (
	"encoding/json"
	"net/http"
	"os"
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

func TestHealth(t *testing.T) {
	err := os.Setenv("JWT_SECRET", "abcd")
	require.NoError(t, err)
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
		switch resp.Data.Status {
		case "healthy":
			assert.Equal(t, 200, res.StatusCode)
		case "degraded":
			assert.Equal(t, 206, res.StatusCode) // StatusPartialContent
		case "unhealthy":
			assert.Equal(t, 503, res.StatusCode) // StatusServiceUnavailable
		}
	})

	t.Run("health check with database unhealthy", func(t *testing.T) {
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

		// Mock the project service to return error for health check
		mockProjectSvc := &services.MockProjectService{}
		mockProjectSvc.On("GetAll", mock.Anything).Return([]sdk.Project(nil), assert.AnError)
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
		assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

		var resp HealthResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "unhealthy", resp.Data.Status)
		assert.Equal(t, "unhealthy", resp.Data.Components["database"])
		assert.Equal(t, "healthy", resp.Data.Components["cache"])
	})

	t.Run("health check with database unavailable", func(t *testing.T) {
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

		prv := &providers.Provider{
			S: svcs,
			D: nil, // Database unavailable
			C: cs,
		}

		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/health")

		req, _ := http.NewRequest("GET", "/health/v1", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusPartialContent, res.StatusCode)

		var resp HealthResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "degraded", resp.Data.Status)
		assert.Equal(t, "unavailable", resp.Data.Components["database"])
		assert.Equal(t, "healthy", resp.Data.Components["cache"])
	})

	t.Run("health check with cache unavailable", func(t *testing.T) {
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
		svcs.Projects = mockProjectSvc

		prv := &providers.Provider{
			S: svcs,
			D: d,
			C: nil, // Cache unavailable
		}

		app.Use(providers.Handle(prv))
		RegisterRoutes(app, "/health")

		req, _ := http.NewRequest("GET", "/health/v1", nil)
		res, err := app.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusPartialContent, res.StatusCode)

		var resp HealthResponse
		err = json.NewDecoder(res.Body).Decode(&resp)
		assert.Nil(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "degraded", resp.Data.Status)
		assert.Equal(t, "healthy", resp.Data.Components["database"])
		assert.Equal(t, "unavailable", resp.Data.Components["cache"])
	})
}

func TestHealthRoute(t *testing.T) {
	app := fiber.New()
	HealthRoute(app, "/health")

	// Check if routes are registered
	routes := app.GetRoutes()
	assert.NotEmpty(t, routes, "Routes should be registered")

	// Check that at least one GET route exists
	getRouteFound := false
	for _, route := range routes {
		if route.Method == "GET" {
			getRouteFound = true
			break
		}
	}
	assert.True(t, getRouteFound, "GET route should be registered")
}

func TestHealthResponse_Structs(t *testing.T) {
	t.Run("HealthResponse struct initialization", func(t *testing.T) {
		response := HealthResponse{
			Success:   true,
			Message:   "test message",
			Timestamp: "2023-01-01T00:00:00Z",
			Data: HealthData{
				Status:  "healthy",
				Version: "1.0.0",
				Uptime:  "1h0m0s",
				Components: map[string]string{
					"database": "healthy",
					"cache":    "healthy",
				},
			},
		}

		assert.True(t, response.Success)
		assert.Equal(t, "test message", response.Message)
		assert.Equal(t, "2023-01-01T00:00:00Z", response.Timestamp)
		assert.Equal(t, "healthy", response.Data.Status)
		assert.Equal(t, "1.0.0", response.Data.Version)
		assert.Equal(t, "1h0m0s", response.Data.Uptime)
		assert.Equal(t, "healthy", response.Data.Components["database"])
		assert.Equal(t, "healthy", response.Data.Components["cache"])
	})

	t.Run("HealthData struct initialization", func(t *testing.T) {
		data := HealthData{
			Status:  "degraded",
			Version: "2.0.0",
			Uptime:  "2h30m15s",
			Components: map[string]string{
				"database": "healthy",
				"cache":    "unavailable",
			},
		}

		assert.Equal(t, "degraded", data.Status)
		assert.Equal(t, "2.0.0", data.Version)
		assert.Equal(t, "2h30m15s", data.Uptime)
		assert.Equal(t, "healthy", data.Components["database"])
		assert.Equal(t, "unavailable", data.Components["cache"])
	})
}
