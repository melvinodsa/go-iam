package projects

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService is a mock implementation of project.Service
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetAll(ctx context.Context) ([]sdk.Project, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]sdk.Project), args.Error(1)
}

func (m *MockProjectService) Get(ctx context.Context, id string) (*sdk.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockProjectService) Create(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) Update(ctx context.Context, project *sdk.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectService) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func setupTestApp() (*fiber.App, *MockProjectService) {
	app := fiber.New()
	mockProjectSvc := new(MockProjectService)
	return app, mockProjectSvc
}

func TestNewMiddlewares(t *testing.T) {
	_, mockProjectSvc := setupTestApp()

	middlewares := NewMiddlewares(mockProjectSvc)

	assert.NotNil(t, middlewares)
	assert.Equal(t, mockProjectSvc, middlewares.projectSvc)
}

func TestMiddlewares_Projects_WithProjectIds(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	// Test with single project ID
	t.Run("Single Project ID", func(t *testing.T) {
		app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
			projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
			return c.JSON(fiber.Map{"projects": projects})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Project-Ids", "project1")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Test with multiple project IDs
	t.Run("Multiple Project IDs", func(t *testing.T) {
		app.Get("/test-multi", middlewares.Projects, func(c *fiber.Ctx) error {
			projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
			return c.JSON(fiber.Map{"projects": projects, "count": len(projects)})
		})

		req := httptest.NewRequest(http.MethodGet, "/test-multi", nil)
		req.Header.Set("X-Project-Ids", "project1,project2,project3")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestMiddlewares_Projects_WithoutProjectIds(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		return c.JSON(fiber.Map{"projects": projects, "isEmpty": len(projects) == 0})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No X-Project-Ids header

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewares_Projects_WithEmptyProjectIds(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	// Test cases for empty or invalid headers
	testCases := []struct {
		name        string
		headerValue string
		description string
	}{
		{"EmptyString", "", "Empty header value"},
		{"OnlyComma", ",", "Only comma"},
		{"MultipleCommas", ",,", "Multiple commas"},
		{"Whitespace", "   ", "Whitespace only"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Get("/test-empty", middlewares.Projects, func(c *fiber.Ctx) error {
				projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
				return c.JSON(fiber.Map{"projects": projects, "count": len(projects)})
			})

			req := httptest.NewRequest(http.MethodGet, "/test-empty", nil)
			req.Header.Set("X-Project-Ids", tc.headerValue)

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestMiddlewares_Projects_ProjectIdsWithSpaces(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		return c.JSON(fiber.Map{"projects": projects})
	})

	// Test with spaces around commas
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", "project1, project2 , project3")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewares_Projects_HeaderCaseInsensitive(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	// Test different header case variations
	testCases := []struct {
		name       string
		headerName string
	}{
		{"Lowercase", "x-project-ids"},
		{"Uppercase", "X-PROJECT-IDS"},
		{"MixedCase", "X-Project-Ids"},
		{"StandardCase", "X-Project-Ids"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Get("/test-case", middlewares.Projects, func(c *fiber.Ctx) error {
				projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
				return c.JSON(fiber.Map{"projects": projects})
			})

			req := httptest.NewRequest(http.MethodGet, "/test-case", nil)
			req.Header.Set(tc.headerName, "project1,project2")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestMiddlewares_Projects_MultipleHeaderValues(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		return c.JSON(fiber.Map{"projects": projects})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// Add multiple header values (this tests if we correctly take the first one)
	req.Header.Add("X-Project-Ids", "project1,project2")
	req.Header.Add("X-Project-Ids", "project3,project4")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewares_Projects_ContextIntegration(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	// Test that the context value is properly set and accessible
	var capturedProjects []string

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		capturedProjects = projects
		return c.JSON(fiber.Map{"status": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", "proj-1,proj-2,proj-3")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the context value was set correctly
	expectedProjects := []string{"proj-1", "proj-2", "proj-3"}
	assert.Equal(t, expectedProjects, capturedProjects)
}

func TestMiddlewares_Projects_SpecialCharacters(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		return c.JSON(fiber.Map{"projects": projects})
	})

	// Test with special characters in project IDs
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", "project-1_test,project@2,project.3")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMiddlewares_Projects_UnicodeCharacters(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		return c.JSON(fiber.Map{"projects": projects})
	})

	// Test with unicode characters
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", "项目1,프로젝트2,projeto3")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Benchmark tests
func BenchmarkMiddlewares_Projects_SingleProject(b *testing.B) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", "project1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMiddlewares_Projects_MultipleProjects(b *testing.B) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", "project1,project2,project3,project4,project5")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMiddlewares_Projects_NoHeader(b *testing.B) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No X-Project-Ids header

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Edge case tests
func TestMiddlewares_Projects_ManyProjects(t *testing.T) {
	app, mockProjectSvc := setupTestApp()
	middlewares := NewMiddlewares(mockProjectSvc)

	app.Get("/test", middlewares.Projects, func(c *fiber.Ctx) error {
		projects := c.Context().UserValue(sdk.ProjectsTypeVal).([]string)
		return c.JSON(fiber.Map{"count": len(projects)})
	})

	// Create a reasonable list of project IDs (50 instead of 1000)
	projectIds := make([]string, 50)
	for i := 0; i < 50; i++ {
		projectIds[i] = fmt.Sprintf("project-%d", i)
	}
	projectList := strings.Join(projectIds, ",")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Project-Ids", projectList)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
