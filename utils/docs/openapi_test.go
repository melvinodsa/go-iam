package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestApiWrapper_Struct(t *testing.T) {
	api := ApiWrapper{
		Path:        "/users/:id",
		Method:      "GET",
		Name:        "Get User",
		Description: "Retrieve user by ID",
		Tags:        []string{"users"},
		RequestBody: &ApiRequestBody{
			Description: "User data",
			Content:     map[string]interface{}{"name": "John"},
		},
		Response: &ApiResponse{
			Description: "User response",
			Content:     map[string]interface{}{"id": 1},
		},
		Parameters: []ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "User ID",
				Required:    true,
			},
		},
		UnAuthenticated:      false,
		ProjectIDNotRequired: false,
	}

	assert.Equal(t, "/users/:id", api.Path)
	assert.Equal(t, "GET", api.Method)
	assert.Equal(t, "Get User", api.Name)
	assert.Equal(t, "Retrieve user by ID", api.Description)
	assert.Equal(t, []string{"users"}, api.Tags)
	assert.NotNil(t, api.RequestBody)
	assert.NotNil(t, api.Response)
	assert.Len(t, api.Parameters, 1)
	assert.False(t, api.UnAuthenticated)
	assert.False(t, api.ProjectIDNotRequired)
}

func TestApiParameter_Struct(t *testing.T) {
	param := ApiParameter{
		Name:        "limit",
		In:          "query",
		Description: "Number of items to return",
		Required:    false,
	}

	assert.Equal(t, "limit", param.Name)
	assert.Equal(t, "query", param.In)
	assert.Equal(t, "Number of items to return", param.Description)
	assert.False(t, param.Required)
}

func TestApiRequestBody_Struct(t *testing.T) {
	requestBody := ApiRequestBody{
		Description: "Request data",
		Content:     map[string]interface{}{"key": "value"},
	}

	assert.Equal(t, "Request data", requestBody.Description)
	assert.NotNil(t, requestBody.Content)
}

func TestApiResponse_Struct(t *testing.T) {
	response := ApiResponse{
		Description: "Response data",
		Content:     map[string]interface{}{"status": "success"},
	}

	assert.Equal(t, "Response data", response.Description)
	assert.NotNil(t, response.Content)
}

func TestRegisterApi(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	api1 := ApiWrapper{
		Path:   "/test",
		Method: "GET",
		Name:   "Test API 1",
	}

	api2 := ApiWrapper{
		Path:   "/test",
		Method: "POST",
		Name:   "Test API 2",
	}

	RegisterApi(api1)
	assert.Len(t, apis["/test"], 1)
	assert.Equal(t, "Test API 1", apis["/test"][0].Name)

	RegisterApi(api2)
	assert.Len(t, apis["/test"], 2)
	assert.Equal(t, "Test API 1", apis["/test"][0].Name)
	assert.Equal(t, "Test API 2", apis["/test"][1].Name)
}

func TestGenerateOpenAPI_NoAPIs(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	_, err := GenerateOpenAPI()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no APIs registered")
}

func TestGenerateOpenAPI_Success(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	api := ApiWrapper{
		Path:        "/users/:id",
		Method:      "GET",
		Name:        "Get User",
		Description: "Retrieve user by ID",
		Tags:        []string{"users"},
		RequestBody: &ApiRequestBody{
			Description: "User data",
			Content:     map[string]interface{}{"name": "John"},
		},
		Response: &ApiResponse{
			Description: "User response",
			Content:     map[string]interface{}{"id": 1, "name": "John"},
		},
		Parameters: []ApiParameter{
			{
				Name:        "id",
				In:          "path",
				Description: "User ID",
				Required:    true,
			},
		},
		UnAuthenticated:      false,
		ProjectIDNotRequired: false,
	}

	RegisterApi(api)

	result, err := GenerateOpenAPI()
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify it's valid YAML
	var parsed map[string]interface{}
	err = yaml.Unmarshal(result, &parsed)
	require.NoError(t, err)

	// Verify basic structure
	assert.Equal(t, "3.0.3", parsed["openapi"])

	info, ok := parsed["info"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Go IAM APIs", info["title"])
	assert.Equal(t, "1.0.0", info["version"])

	paths, ok := parsed["paths"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, paths, "/users/{id}")

	components, ok := parsed["components"].(map[string]interface{})
	require.True(t, ok)

	securitySchemes, ok := components["securitySchemes"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, securitySchemes, "BearerAuth")
}

func TestGenerateOpenAPI_WithInvalidAPI(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	api := ApiWrapper{
		Path:   "/test",
		Method: "INVALID",
		Name:   "Test API",
	}

	RegisterApi(api)

	_, err := GenerateOpenAPI()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported HTTP method INVALID")
}

func TestExtractPathParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No parameters",
			input:    "/users",
			expected: "/users",
		},
		{
			name:     "Single parameter",
			input:    "/users/:id",
			expected: "/users/{id}",
		},
		{
			name:     "Multiple parameters",
			input:    "/users/:userId/posts/:postId",
			expected: "/users/{userId}/posts/{postId}",
		},
		{
			name:     "Mixed path",
			input:    "/api/v1/users/:id/profile",
			expected: "/api/v1/users/{id}/profile",
		},
		{
			name:     "Root path",
			input:    "/",
			expected: "/",
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPathParams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateDocsForApi_EmptyAPIs(t *testing.T) {
	_, err := generateDocsForApi([]ApiWrapper{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API path and method must be specified")
}

func TestGenerateDocsForApi_Success(t *testing.T) {
	apis := []ApiWrapper{
		{
			Path:        "/users",
			Method:      "GET",
			Name:        "List Users",
			Description: "Get all users",
			Tags:        []string{"users"},
			RequestBody: &ApiRequestBody{
				Description: "Filter criteria",
				Content:     map[string]interface{}{"limit": 10},
			},
			Response: &ApiResponse{
				Description: "List of users",
				Content:     []map[string]interface{}{{"id": 1, "name": "John"}},
			},
			Parameters: []ApiParameter{
				{
					Name:        "limit",
					In:          "query",
					Description: "Number of items",
					Required:    false,
				},
			},
			UnAuthenticated:      false,
			ProjectIDNotRequired: false,
		},
		{
			Path:                 "/users",
			Method:               "POST",
			Name:                 "Create User",
			Description:          "Create a new user",
			Tags:                 []string{"users"},
			UnAuthenticated:      true,
			ProjectIDNotRequired: true,
		},
	}

	pathItem, err := generateDocsForApi(apis)
	require.NoError(t, err)
	assert.NotNil(t, pathItem)

	// Verify GET operation
	assert.NotNil(t, pathItem.Get)
	assert.Equal(t, "List Users", pathItem.Get.Summary)
	assert.Equal(t, "Get all users", pathItem.Get.Description)
	assert.Equal(t, []string{"users"}, pathItem.Get.Tags)
	assert.NotNil(t, pathItem.Get.RequestBody)
	assert.NotNil(t, pathItem.Get.Responses)
	assert.NotNil(t, pathItem.Get.Security)
	assert.Len(t, pathItem.Get.Parameters, 2) // 1 custom + 1 X-Project-Ids

	// Verify POST operation
	assert.NotNil(t, pathItem.Post)
	assert.Equal(t, "Create User", pathItem.Post.Summary)
	assert.Equal(t, "Create a new user", pathItem.Post.Description)
	assert.Equal(t, []string{"users"}, pathItem.Post.Tags)
	assert.Nil(t, pathItem.Post.Security)      // Unauthenticated
	assert.Len(t, pathItem.Post.Parameters, 0) // No X-Project-Ids since ProjectIDNotRequired is true
}

func TestGenerateDocsForApi_AllHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			apis := []ApiWrapper{
				{
					Path:   "/test",
					Method: method,
					Name:   "Test " + method,
				},
			}

			pathItem, err := generateDocsForApi(apis)
			require.NoError(t, err)
			assert.NotNil(t, pathItem)

			switch method {
			case "GET":
				assert.NotNil(t, pathItem.Get)
			case "POST":
				assert.NotNil(t, pathItem.Post)
			case "PUT":
				assert.NotNil(t, pathItem.Put)
			case "DELETE":
				assert.NotNil(t, pathItem.Delete)
			}
		})
	}
}

func TestGenerateDocsForApi_UnsupportedMethod(t *testing.T) {
	apis := []ApiWrapper{
		{
			Path:   "/test",
			Method: "PATCH",
			Name:   "Test PATCH",
		},
	}

	_, err := generateDocsForApi(apis)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported HTTP method PATCH")
}

func TestGenerateDocsForRequestBody_NilContent(t *testing.T) {
	requestBody := ApiRequestBody{
		Description: "Test",
		Content:     nil,
	}

	_, err := generateDocsForRequestBody(requestBody)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request body content must be provided")
}

func TestGenerateDocsForRequestBody_Success(t *testing.T) {
	requestBody := ApiRequestBody{
		Description: "User data",
		Content:     map[string]interface{}{"name": "John", "email": "john@example.com"},
	}

	result, err := generateDocsForRequestBody(requestBody)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Value)
	assert.Equal(t, "User data", result.Value.Description)
	assert.Contains(t, result.Value.Content, "application/json")
}

func TestGenerateDocsForResponse_NilContent(t *testing.T) {
	response := ApiResponse{
		Description: "Test",
		Content:     nil,
	}

	_, err := generateDocsForResponse(response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "response content must be provided")
}

func TestGenerateDocsForResponse_Success(t *testing.T) {
	response := ApiResponse{
		Description: "User response",
		Content:     map[string]interface{}{"id": 1, "name": "John"},
	}

	result, err := generateDocsForResponse(response)
	require.NoError(t, err)
	assert.NotNil(t, result)

	defaultResponse := result.Map()["default"]
	assert.NotNil(t, defaultResponse)
	assert.Equal(t, "User response", *defaultResponse.Value.Description)
	assert.Contains(t, defaultResponse.Value.Content, "application/json")
}

func TestGenerateDocsForParameters_EmptyParameters(t *testing.T) {
	result, err := generateDocsForParameters([]ApiParameter{})
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGenerateDocsForParameters_Success(t *testing.T) {
	parameters := []ApiParameter{
		{
			Name:        "limit",
			In:          "query",
			Description: "Number of items",
			Required:    false,
		},
		{
			Name:        "id",
			In:          "path",
			Description: "Resource ID",
			Required:    true,
		},
	}

	result, err := generateDocsForParameters(parameters)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// Verify first parameter
	assert.Equal(t, "limit", result[0].Value.Name)
	assert.Equal(t, "query", result[0].Value.In)
	assert.Equal(t, "Number of items", result[0].Value.Description)
	assert.False(t, result[0].Value.Required)

	// Verify second parameter
	assert.Equal(t, "id", result[1].Value.Name)
	assert.Equal(t, "path", result[1].Value.In)
	assert.Equal(t, "Resource ID", result[1].Value.Description)
	assert.True(t, result[1].Value.Required)
}

func TestCreateOpenApiDoc_Success(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	// Register a test API
	api := ApiWrapper{
		Path:        "/test",
		Method:      "GET",
		Name:        "Test API",
		Description: "Test description",
		Response: &ApiResponse{
			Description: "Test response",
			Content:     map[string]interface{}{"status": "ok"},
		},
	}
	RegisterApi(api)

	// Create temporary file
	tempDir := t.TempDir()
	fileName := filepath.Join(tempDir, "test-openapi.yaml")

	err := CreateOpenApiDoc(fileName)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(fileName)
	assert.NoError(t, err)

	// Verify file content
	content, err := os.ReadFile(fileName)
	require.NoError(t, err)
	assert.NotEmpty(t, content)

	// Verify it's valid YAML
	var parsed map[string]interface{}
	err = yaml.Unmarshal(content, &parsed)
	require.NoError(t, err)
}

func TestCreateOpenApiDoc_GenerateError(t *testing.T) {
	// Clear any existing APIs to trigger error
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	tempDir := t.TempDir()
	fileName := filepath.Join(tempDir, "test-openapi.yaml")

	err := CreateOpenApiDoc(fileName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error generating OpenAPI document")
}

func TestCreateOpenApiDoc_WriteError(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	// Register a test API
	api := ApiWrapper{
		Path:   "/test",
		Method: "GET",
		Name:   "Test API",
		Response: &ApiResponse{
			Description: "Test response",
			Content:     map[string]interface{}{"status": "ok"},
		},
	}
	RegisterApi(api)

	// Try to write to invalid path
	err := CreateOpenApiDoc("/invalid/path/that/does/not/exist/test.yaml")
	assert.Error(t, err)
}

func TestGenerateOpenAPI_ComplexScenario(t *testing.T) {
	// Clear any existing APIs
	originalApis := apis
	defer func() { apis = originalApis }()
	apis = make(map[string][]ApiWrapper)

	// Register multiple APIs with different configurations
	apis["/users"] = []ApiWrapper{
		{
			Path:        "/users",
			Method:      "GET",
			Name:        "List Users",
			Description: "Get all users",
			Tags:        []string{"users", "admin"},
			Parameters: []ApiParameter{
				{Name: "page", In: "query", Description: "Page number", Required: false},
				{Name: "limit", In: "query", Description: "Items per page", Required: false},
			},
			Response: &ApiResponse{
				Description: "List of users",
				Content:     []map[string]interface{}{{"id": 1, "name": "John"}},
			},
			UnAuthenticated:      false,
			ProjectIDNotRequired: false,
		},
		{
			Path:        "/users",
			Method:      "POST",
			Name:        "Create User",
			Description: "Create a new user",
			Tags:        []string{"users"},
			RequestBody: &ApiRequestBody{
				Description: "User data",
				Content:     map[string]interface{}{"name": "string", "email": "string"},
			},
			Response: &ApiResponse{
				Description: "Created user",
				Content:     map[string]interface{}{"id": 1, "name": "John"},
			},
			UnAuthenticated:      true,
			ProjectIDNotRequired: true,
		},
	}

	apis["/users/:id"] = []ApiWrapper{
		{
			Path:        "/users/:id",
			Method:      "GET",
			Name:        "Get User",
			Description: "Get user by ID",
			Tags:        []string{"users"},
			Parameters: []ApiParameter{
				{Name: "id", In: "path", Description: "User ID", Required: true},
			},
			Response: &ApiResponse{
				Description: "User details",
				Content:     map[string]interface{}{"id": 1, "name": "John"},
			},
			UnAuthenticated:      false,
			ProjectIDNotRequired: false,
		},
		{
			Path:        "/users/:id",
			Method:      "DELETE",
			Name:        "Delete User",
			Description: "Delete user by ID",
			Tags:        []string{"users", "admin"},
			Parameters: []ApiParameter{
				{Name: "id", In: "path", Description: "User ID", Required: true},
			},
			Response: &ApiResponse{
				Description: "Success response",
				Content:     map[string]interface{}{"message": "User deleted"},
			},
			UnAuthenticated:      false,
			ProjectIDNotRequired: false,
		},
	}

	result, err := GenerateOpenAPI()
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify it's valid YAML
	var parsed map[string]interface{}
	err = yaml.Unmarshal(result, &parsed)
	require.NoError(t, err)

	// Verify paths
	paths, ok := parsed["paths"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, paths, "/users")
	assert.Contains(t, paths, "/users/{id}")

	// Verify users path has GET and POST
	usersPath, ok := paths["/users"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, usersPath, "get")
	assert.Contains(t, usersPath, "post")

	// Verify users/{id} path has GET and DELETE
	userIdPath, ok := paths["/users/{id}"].(map[string]interface{})
	require.True(t, ok)
	assert.Contains(t, userIdPath, "get")
	assert.Contains(t, userIdPath, "delete")
}

func TestGenerateDocsForApi_RequestBodyError(t *testing.T) {
	apis := []ApiWrapper{
		{
			Path:   "/test",
			Method: "POST",
			Name:   "Test API",
			RequestBody: &ApiRequestBody{
				Description: "Test",
				Content:     nil, // This will cause an error
			},
		},
	}

	_, err := generateDocsForApi(apis)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error generating request body")
}

func TestGenerateDocsForApi_ResponseError(t *testing.T) {
	apis := []ApiWrapper{
		{
			Path:   "/test",
			Method: "GET",
			Name:   "Test API",
			Response: &ApiResponse{
				Description: "Test",
				Content:     nil, // This will cause an error
			},
		},
	}

	_, err := generateDocsForApi(apis)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error generating response")
}

func TestConstants(t *testing.T) {
	// Test that the constants are defined and not empty
	assert.NotEmpty(t, intoDescription)
	assert.Contains(t, intoDescription, "Go IAM")
	assert.Contains(t, intoDescription, "Identity and Access Management")
	assert.Contains(t, strings.ToLower(intoDescription), "golang")
}

// Test to ensure global state is properly managed
func TestGlobalApisState(t *testing.T) {
	// Save original state
	originalApis := apis
	defer func() { apis = originalApis }()

	// Clear state
	apis = make(map[string][]ApiWrapper)
	assert.Empty(t, apis)

	// Add API
	api := ApiWrapper{Path: "/test", Method: "GET", Name: "Test"}
	RegisterApi(api)
	assert.Len(t, apis, 1)
	assert.Len(t, apis["/test"], 1)

	// Add another API to same path
	api2 := ApiWrapper{Path: "/test", Method: "POST", Name: "Test POST"}
	RegisterApi(api2)
	assert.Len(t, apis, 1)
	assert.Len(t, apis["/test"], 2)

	// Add API to different path
	api3 := ApiWrapper{Path: "/other", Method: "GET", Name: "Other"}
	RegisterApi(api3)
	assert.Len(t, apis, 2)
	assert.Len(t, apis["/test"], 2)
	assert.Len(t, apis["/other"], 1)
}
