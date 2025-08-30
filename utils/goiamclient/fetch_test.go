package goaiamclient

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClientService is a mock implementation of client.Service
type MockClientService struct {
	mock.Mock
}

func (m *MockClientService) GetAll(ctx context.Context, queryParams sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockClientService) GetGoIamClients(ctx context.Context, params sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockClientService) Get(ctx context.Context, id string, dontCheckProjects bool) (*sdk.Client, error) {
	args := m.Called(ctx, id, dontCheckProjects)
	return args.Get(0).(*sdk.Client), args.Error(1)
}

func (m *MockClientService) Create(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientService) Update(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

// Mock methods for utils.Emitter interface
func (m *MockClientService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Client], sdk.Client]) {
	m.Called(eventName, subscriber)
}

func (m *MockClientService) Emit(event utils.Event[sdk.Client]) {
	m.Called(event)
}

func TestGetGoIamClient_Success(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create test client
	now := time.Now()
	expectedClient := sdk.Client{
		Id:                    "test-client-id",
		Name:                  "Test Go IAM Client",
		Description:           "Test client for Go IAM",
		Secret:                "test-secret",
		Tags:                  []string{"go-iam", "test"},
		RedirectURLs:          []string{"http://localhost:3000/callback"},
		Scopes:                []string{"read", "write"},
		ProjectId:             "test-project-id",
		DefaultAuthProviderId: "test-auth-provider-id",
		GoIamClient:           true,
		Enabled:               true,
		CreatedAt:             &now,
		CreatedBy:             "test-user",
		UpdatedAt:             &now,
		UpdatedBy:             "test-user",
	}

	// Set up mock expectations
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{expectedClient}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedClient.Id, result.Id)
	assert.Equal(t, expectedClient.Name, result.Name)
	assert.Equal(t, expectedClient.Description, result.Description)
	assert.Equal(t, expectedClient.Secret, result.Secret)
	assert.Equal(t, expectedClient.Tags, result.Tags)
	assert.Equal(t, expectedClient.RedirectURLs, result.RedirectURLs)
	assert.Equal(t, expectedClient.Scopes, result.Scopes)
	assert.Equal(t, expectedClient.ProjectId, result.ProjectId)
	assert.Equal(t, expectedClient.DefaultAuthProviderId, result.DefaultAuthProviderId)
	assert.True(t, result.GoIamClient)
	assert.True(t, result.Enabled)
	assert.Equal(t, expectedClient.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedClient.CreatedBy, result.CreatedBy)
	assert.Equal(t, expectedClient.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, expectedClient.UpdatedBy, result.UpdatedBy)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_MultipleClients(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create multiple test clients
	now := time.Now()
	client1 := sdk.Client{
		Id:          "client-1",
		Name:        "First Go IAM Client",
		GoIamClient: true,
		CreatedAt:   &now,
	}

	client2 := sdk.Client{
		Id:          "client-2",
		Name:        "Second Go IAM Client",
		GoIamClient: true,
		CreatedAt:   &now,
	}

	// Set up mock expectations - should return first client
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{client1, client2}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions - should return the first client
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, client1.Id, result.Id)
	assert.Equal(t, client1.Name, result.Name)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_ServiceError(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Set up mock expectations to return an error
	expectedError := errors.New("database connection failed")
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{}, expectedError)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.Nil(t, result)
	assert.Error(t, err)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_NoClientsFound(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Set up mock expectations to return empty slice
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.Nil(t, result)
	assert.NoError(t, err)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_EmptyClient(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create empty client (minimal required fields)
	emptyClient := sdk.Client{
		Id:          "",
		Name:        "",
		GoIamClient: true,
		Enabled:     false,
	}

	// Set up mock expectations
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{emptyClient}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Id)
	assert.Equal(t, "", result.Name)
	assert.True(t, result.GoIamClient)
	assert.False(t, result.Enabled)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_NilFields(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create client with nil timestamp fields
	clientWithNilFields := sdk.Client{
		Id:          "test-id",
		Name:        "Test Client",
		GoIamClient: true,
		CreatedAt:   nil,
		UpdatedAt:   nil,
	}

	// Set up mock expectations
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{clientWithNilFields}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.Id)
	assert.Equal(t, "Test Client", result.Name)
	assert.Nil(t, result.CreatedAt)
	assert.Nil(t, result.UpdatedAt)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_CorrectQueryParams(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create test client
	testClient := sdk.Client{
		Id:          "test-id",
		GoIamClient: true,
	}

	// Set up mock expectations with exact parameter matching
	expectedParams := sdk.ClientQueryParams{
		GoIamClient: true,
	}

	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), expectedParams).Return([]sdk.Client{testClient}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test-id", result.Id)

	// Verify that the correct parameters were passed
	mockService.AssertExpectations(t)

	// Additional verification that GetGoIamClients was called with correct params
	calls := mockService.Calls
	assert.Len(t, calls, 1)
	assert.Equal(t, "GetGoIamClients", calls[0].Method)

	// Check the query params argument
	actualParams := calls[0].Arguments[1].(sdk.ClientQueryParams)
	assert.True(t, actualParams.GoIamClient)
	assert.Empty(t, actualParams.ProjectIds)
	assert.False(t, actualParams.SortByUpdatedAt)
}

func TestGetGoIamClient_ContextPassing(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create test client
	testClient := sdk.Client{
		Id:          "test-id",
		GoIamClient: true,
	}

	// Set up mock with context verification
	mockService.On("GetGoIamClients", mock.MatchedBy(func(ctx context.Context) bool {
		// Verify that a context is passed (should be background context)
		return ctx != nil
	}), mock.AnythingOfType("sdk.ClientQueryParams")).Return([]sdk.Client{testClient}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Assertions
	assert.NotNil(t, result)
	assert.NoError(t, err)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_ClientWithAllFields(t *testing.T) {
	// Create mock service
	mockService := new(MockClientService)

	// Create comprehensive test client with all fields populated
	now := time.Now()
	fullClient := sdk.Client{
		Id:                    "full-client-id",
		Name:                  "Full Test Client",
		Description:           "A comprehensive test client with all fields",
		Secret:                "super-secret-key",
		Tags:                  []string{"production", "go-iam", "auth"},
		RedirectURLs:          []string{"https://app.example.com/callback", "http://localhost:3000/auth"},
		Scopes:                []string{"read", "write", "admin", "delete"},
		ProjectId:             "project-123",
		DefaultAuthProviderId: "auth-provider-456",
		GoIamClient:           true,
		Enabled:               true,
		CreatedAt:             &now,
		CreatedBy:             "admin-user",
		UpdatedAt:             &now,
		UpdatedBy:             "admin-user",
	}

	// Set up mock expectations
	mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
		GoIamClient: true,
	}).Return([]sdk.Client{fullClient}, nil)

	// Call the function
	result, err := GetGoIamClient(mockService)

	// Comprehensive assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, fullClient.Id, result.Id)
	assert.Equal(t, fullClient.Name, result.Name)
	assert.Equal(t, fullClient.Description, result.Description)
	assert.Equal(t, fullClient.Secret, result.Secret)
	assert.Equal(t, fullClient.Tags, result.Tags)
	assert.Equal(t, len(fullClient.Tags), len(result.Tags))
	assert.Contains(t, result.Tags, "production")
	assert.Contains(t, result.Tags, "go-iam")
	assert.Contains(t, result.Tags, "auth")
	assert.Equal(t, fullClient.RedirectURLs, result.RedirectURLs)
	assert.Equal(t, len(fullClient.RedirectURLs), len(result.RedirectURLs))
	assert.Contains(t, result.RedirectURLs, "https://app.example.com/callback")
	assert.Contains(t, result.RedirectURLs, "http://localhost:3000/auth")
	assert.Equal(t, fullClient.Scopes, result.Scopes)
	assert.Equal(t, len(fullClient.Scopes), len(result.Scopes))
	assert.Contains(t, result.Scopes, "read")
	assert.Contains(t, result.Scopes, "write")
	assert.Contains(t, result.Scopes, "admin")
	assert.Contains(t, result.Scopes, "delete")
	assert.Equal(t, fullClient.ProjectId, result.ProjectId)
	assert.Equal(t, fullClient.DefaultAuthProviderId, result.DefaultAuthProviderId)
	assert.True(t, result.GoIamClient)
	assert.True(t, result.Enabled)
	assert.Equal(t, fullClient.CreatedAt.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, fullClient.CreatedBy, result.CreatedBy)
	assert.Equal(t, fullClient.UpdatedAt.Unix(), result.UpdatedAt.Unix())
	assert.Equal(t, fullClient.UpdatedBy, result.UpdatedBy)

	// Verify mock expectations
	mockService.AssertExpectations(t)
}

func TestGetGoIamClient_ErrorTypes(t *testing.T) {
	testCases := []struct {
		name           string
		serviceError   error
		expectedResult *sdk.Client
	}{
		{
			name:           "Database error",
			serviceError:   errors.New("database connection failed"),
			expectedResult: nil,
		},
		{
			name:           "Network error",
			serviceError:   errors.New("network timeout"),
			expectedResult: nil,
		},
		{
			name:           "Authentication error",
			serviceError:   errors.New("unauthorized access"),
			expectedResult: nil,
		},
		{
			name:           "Generic error",
			serviceError:   errors.New("something went wrong"),
			expectedResult: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockClientService)

			// Set up mock expectations
			mockService.On("GetGoIamClients", mock.AnythingOfType("context.backgroundCtx"), sdk.ClientQueryParams{
				GoIamClient: true,
			}).Return([]sdk.Client{}, tc.serviceError)

			// Call the function
			result, err := GetGoIamClient(mockService)

			// Assertions
			assert.Equal(t, tc.expectedResult, result)
			assert.Error(t, err)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
