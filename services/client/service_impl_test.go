package client

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/services/project"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore implements Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetAll(ctx context.Context, params sdk.ClientQueryParams) ([]sdk.Client, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]sdk.Client), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*sdk.Client, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.Client), args.Error(1)
}

func (m *MockStore) Create(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, client *sdk.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

// MockProjectService implements project.Service interface for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) GetAll(ctx context.Context) ([]sdk.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sdk.Project), args.Error(1)
}

func (m *MockProjectService) Get(ctx context.Context, id string) (*sdk.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.Project), args.Error(1)
}

func (m *MockProjectService) GetByName(ctx context.Context, name string) (*sdk.Project, error) {
	args := m.Called(ctx, name)
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

// Helper function to create context with projects
func createContextWithProjects(projects []string) context.Context {
	metadata := sdk.Metadata{
		User:       nil,
		ProjectIds: projects,
	}
	return middlewares.AddMetadata(context.Background(), metadata)
}

func TestNewService(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}

	service := NewService(mockStore, mockProjectService)

	assert.NotNil(t, service)
	assert.Implements(t, (*Service)(nil), service)
}

func TestService_GetAll_Success(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1", "project2"})

	queryParams := sdk.ClientQueryParams{
		SortByUpdatedAt: true,
	}

	expectedParams := sdk.ClientQueryParams{
		ProjectIds:      []string{"project1", "project2"},
		SortByUpdatedAt: true,
	}

	expectedClients := []sdk.Client{
		{Id: "client1", Name: "Test Client 1", ProjectId: "project1"},
		{Id: "client2", Name: "Test Client 2", ProjectId: "project2"},
	}

	mockStore.On("GetAll", ctx, expectedParams).Return(expectedClients, nil)

	result, err := service.GetAll(ctx, queryParams)

	assert.NoError(t, err)
	assert.Equal(t, expectedClients, result)
	mockStore.AssertExpectations(t)
}

func TestService_GetAll_StoreError(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1"})

	queryParams := sdk.ClientQueryParams{}
	expectedParams := sdk.ClientQueryParams{
		ProjectIds: []string{"project1"},
	}

	mockStore.On("GetAll", ctx, expectedParams).Return(([]sdk.Client)(nil), errors.New("database error"))

	result, err := service.GetAll(ctx, queryParams)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockStore.AssertExpectations(t)
}

func TestService_Get_Success(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1", "project2"})

	clientId := "client1"
	expectedClient := &sdk.Client{
		Id:        "client1",
		Name:      "Test Client",
		ProjectId: "project1",
	}

	mockStore.On("Get", ctx, clientId).Return(expectedClient, nil)

	result, err := service.Get(ctx, clientId, false)

	assert.NoError(t, err)
	assert.Equal(t, expectedClient, result)
	mockStore.AssertExpectations(t)
}

func TestService_Get_ClientNotInUserProjects(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1"})

	clientId := "client1"
	client := &sdk.Client{
		Id:        "client1",
		Name:      "Test Client",
		ProjectId: "project2", // Not in user's projects
	}

	mockStore.On("Get", ctx, clientId).Return(client, nil)

	result, err := service.Get(ctx, clientId, false)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrClientNotFound, err)
	mockStore.AssertExpectations(t)
}

func TestService_Create_Success(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	testUser := &sdk.User{Id: "user1", Name: "Test User"}
	metadata := sdk.Metadata{
		User:       testUser,
		ProjectIds: []string{"project1", "project2"},
	}
	ctx := middlewares.AddMetadata(context.Background(), metadata)

	client := &sdk.Client{
		Name:        "Test Client",
		Description: "Test Description",
		ProjectId:   "project1",
	}

	// We expect the store to be called with the client having a generated secret
	mockStore.On("Create", ctx, mock.MatchedBy(func(c *sdk.Client) bool {
		return c.Name == "Test Client" &&
			c.Description == "Test Description" &&
			c.ProjectId == "project1" &&
			c.Secret != "" // Secret should be generated
	})).Return(nil)

	err := service.Create(ctx, client)

	assert.NoError(t, err)
	assert.NotEmpty(t, client.Secret) // Secret should be generated
	mockStore.AssertExpectations(t)
}

func TestService_Create_ProjectNotFound(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1"})

	client := &sdk.Client{
		Name:        "Test Client",
		Description: "Test Description",
		ProjectId:   "project2", // Not in user's projects
	}

	err := service.Create(ctx, client)

	assert.Error(t, err)
	assert.Equal(t, project.ErrProjectNotFound, err)
	// Store should not be called
	mockStore.AssertNotCalled(t, "Create")
}

func TestService_Update_Success(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	testUser := &sdk.User{Id: "user1", Name: "Test User"}
	metadata := sdk.Metadata{
		User:       testUser,
		ProjectIds: []string{"project1", "project2"},
	}
	ctx := middlewares.AddMetadata(context.Background(), metadata)

	client := &sdk.Client{
		Id:          "client1",
		Name:        "Updated Client",
		Description: "Updated Description",
		ProjectId:   "project1",
	}

	mockStore.On("Update", ctx, client).Return(nil)

	err := service.Update(ctx, client)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestService_Update_ProjectNotFound(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1"})

	client := &sdk.Client{
		Id:          "client1",
		Name:        "Updated Client",
		Description: "Updated Description",
		ProjectId:   "project2", // Not in user's projects
	}

	err := service.Update(ctx, client)

	assert.Error(t, err)
	assert.Equal(t, project.ErrProjectNotFound, err)
	// Store should not be called
	mockStore.AssertNotCalled(t, "Update")
}

func TestService_GetGoIamClients_Success(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := context.Background()

	params := sdk.ClientQueryParams{
		ProjectIds: []string{"project1"},
	}

	expectedParams := sdk.ClientQueryParams{
		ProjectIds:      []string{"project1"},
		GoIamClient:     true,
		SortByUpdatedAt: true,
	}

	expectedClients := []sdk.Client{
		{Id: "client1", Name: "GoIam Client 1", GoIamClient: true},
	}

	mockStore.On("GetAll", ctx, expectedParams).Return(expectedClients, nil)

	result, err := service.GetGoIamClients(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, expectedClients, result)
	mockStore.AssertExpectations(t)
}

func TestService_GetGoIamClients_StoreError(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := context.Background()

	params := sdk.ClientQueryParams{}
	expectedParams := sdk.ClientQueryParams{
		GoIamClient:     true,
		SortByUpdatedAt: true,
	}

	mockStore.On("GetAll", ctx, expectedParams).Return(([]sdk.Client)(nil), errors.New("database error"))

	result, err := service.GetGoIamClients(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockStore.AssertExpectations(t)
}

func TestService_Get_WithoutProjectCheck(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := context.Background()

	clientId := "client1"
	expectedClient := &sdk.Client{
		Id:        "client1",
		Name:      "Test Client",
		ProjectId: "project3", // Not in user's projects, but should still be returned
	}

	mockStore.On("Get", ctx, clientId).Return(expectedClient, nil)

	result, err := service.Get(ctx, clientId, true) // dontCheckProjects = true

	assert.NoError(t, err)
	assert.Equal(t, expectedClient, result)
	mockStore.AssertExpectations(t)
}

func TestService_Get_StoreError(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	ctx := createContextWithProjects([]string{"project1"})

	clientId := "client1"

	mockStore.On("Get", ctx, clientId).Return((*sdk.Client)(nil), errors.New("database error"))

	result, err := service.Get(ctx, clientId, false)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockStore.AssertExpectations(t)
}

func TestService_Create_StoreError(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	testUser := &sdk.User{Id: "user1", Name: "Test User"}
	metadata := sdk.Metadata{
		User:       testUser,
		ProjectIds: []string{"project1"},
	}
	ctx := middlewares.AddMetadata(context.Background(), metadata)

	client := &sdk.Client{
		Name:        "Test Client",
		Description: "Test Description",
		ProjectId:   "project1",
	}

	mockStore.On("Create", ctx, mock.AnythingOfType("*sdk.Client")).Return(errors.New("database error"))

	err := service.Create(ctx, client)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error while creating client")
	assert.Contains(t, err.Error(), "database error")
	mockStore.AssertExpectations(t)
}

func TestService_Update_StoreError(t *testing.T) {
	mockStore := &MockStore{}
	mockProjectService := &MockProjectService{}
	service := NewService(mockStore, mockProjectService)

	testUser := &sdk.User{Id: "user1", Name: "Test User"}
	metadata := sdk.Metadata{
		User:       testUser,
		ProjectIds: []string{"project1"},
	}
	ctx := middlewares.AddMetadata(context.Background(), metadata)

	client := &sdk.Client{
		Id:          "client1",
		Name:        "Updated Client",
		Description: "Updated Description",
		ProjectId:   "project1",
	}

	mockStore.On("Update", ctx, client).Return(errors.New("database error"))

	err := service.Update(ctx, client)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error while updating client")
	assert.Contains(t, err.Error(), "database error")
	mockStore.AssertExpectations(t)
}

func TestEvent_Name(t *testing.T) {
	ctx := context.Background()
	client := sdk.Client{Id: "client1", Name: "Test Client"}
	metadata := sdk.Metadata{User: &sdk.User{Id: "user1"}, ProjectIds: []string{"project1"}}

	t.Run("client_created_event", func(t *testing.T) {
		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Name()

		assert.Equal(t, goiamuniverse.EventClientCreated, result)
	})

	t.Run("client_updated_event", func(t *testing.T) {
		event := newEvent(ctx, goiamuniverse.EventClientUpdated, client, metadata)

		result := event.Name()

		assert.Equal(t, goiamuniverse.EventClientUpdated, result)
	})
}

func TestEvent_Payload(t *testing.T) {
	ctx := context.Background()
	metadata := sdk.Metadata{User: &sdk.User{Id: "user1"}, ProjectIds: []string{"project1"}}

	t.Run("simple_client", func(t *testing.T) {
		client := sdk.Client{
			Id:        "client1",
			Name:      "Test Client",
			ProjectId: "project1",
		}

		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Payload()

		assert.Equal(t, client, result)
		assert.Equal(t, "client1", result.Id)
		assert.Equal(t, "Test Client", result.Name)
		assert.Equal(t, "project1", result.ProjectId)
	})

	t.Run("complex_client_with_all_fields", func(t *testing.T) {
		client := sdk.Client{
			Id:                    "client2",
			Name:                  "Complex Client",
			Description:           "A complex test client",
			Secret:                "secret123",
			Tags:                  []string{"tag1", "tag2"},
			RedirectURLs:          []string{"https://example.com/callback"},
			Scopes:                []string{"read", "write"},
			ProjectId:             "project2",
			DefaultAuthProviderId: "provider1",
			GoIamClient:           true,
			Enabled:               true,
		}

		event := newEvent(ctx, goiamuniverse.EventClientUpdated, client, metadata)

		result := event.Payload()

		assert.Equal(t, client, result)
		assert.Equal(t, "client2", result.Id)
		assert.Equal(t, "Complex Client", result.Name)
		assert.Equal(t, "A complex test client", result.Description)
		assert.Equal(t, "secret123", result.Secret)
		assert.Equal(t, []string{"tag1", "tag2"}, result.Tags)
		assert.Equal(t, []string{"https://example.com/callback"}, result.RedirectURLs)
		assert.Equal(t, []string{"read", "write"}, result.Scopes)
		assert.Equal(t, "project2", result.ProjectId)
		assert.Equal(t, "provider1", result.DefaultAuthProviderId)
		assert.True(t, result.GoIamClient)
		assert.True(t, result.Enabled)
	})

	t.Run("empty_client", func(t *testing.T) {
		client := sdk.Client{}

		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Payload()

		assert.Equal(t, client, result)
		assert.Empty(t, result.Id)
		assert.Empty(t, result.Name)
		assert.Empty(t, result.ProjectId)
		assert.False(t, result.GoIamClient)
		assert.False(t, result.Enabled)
	})
}

func TestEvent_Metadata(t *testing.T) {
	ctx := context.Background()
	client := sdk.Client{Id: "client1", Name: "Test Client"}

	t.Run("metadata_with_user_and_projects", func(t *testing.T) {
		user := &sdk.User{
			Id:    "user1",
			Name:  "Test User",
			Email: "test@example.com",
		}
		metadata := sdk.Metadata{
			User:       user,
			ProjectIds: []string{"project1", "project2", "project3"},
		}

		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Metadata()

		assert.Equal(t, metadata, result)
		assert.Equal(t, user, result.User)
		assert.Equal(t, []string{"project1", "project2", "project3"}, result.ProjectIds)
		assert.Equal(t, "user1", result.User.Id)
		assert.Equal(t, "Test User", result.User.Name)
		assert.Equal(t, "test@example.com", result.User.Email)
	})

	t.Run("metadata_with_nil_user", func(t *testing.T) {
		metadata := sdk.Metadata{
			User:       nil,
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventClientUpdated, client, metadata)

		result := event.Metadata()

		assert.Equal(t, metadata, result)
		assert.Nil(t, result.User)
		assert.Equal(t, []string{"project1"}, result.ProjectIds)
	})

	t.Run("metadata_with_empty_projects", func(t *testing.T) {
		user := &sdk.User{Id: "user2", Name: "Another User"}
		metadata := sdk.Metadata{
			User:       user,
			ProjectIds: []string{},
		}

		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Metadata()

		assert.Equal(t, metadata, result)
		assert.Equal(t, user, result.User)
		assert.Empty(t, result.ProjectIds)
		assert.Equal(t, "user2", result.User.Id)
		assert.Equal(t, "Another User", result.User.Name)
	})

	t.Run("empty_metadata", func(t *testing.T) {
		metadata := sdk.Metadata{}

		event := newEvent(ctx, goiamuniverse.EventClientUpdated, client, metadata)

		result := event.Metadata()

		assert.Equal(t, metadata, result)
		assert.Nil(t, result.User)
		assert.Nil(t, result.ProjectIds)
	})
}

type ctxKeyType struct{}

func TestEvent_Context(t *testing.T) {
	client := sdk.Client{Id: "client1", Name: "Test Client"}
	metadata := sdk.Metadata{User: &sdk.User{Id: "user1"}, ProjectIds: []string{"project1"}}

	t.Run("background_context", func(t *testing.T) {
		ctx := context.Background()

		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Context()

		assert.Equal(t, ctx, result)
		assert.NotNil(t, result)
	})

	t.Run("context_with_values", func(t *testing.T) {
		baseCtx := context.Background()
		val := ctxKeyType{}
		ctx := context.WithValue(baseCtx, val, "testValue")

		event := newEvent(ctx, goiamuniverse.EventClientUpdated, client, metadata)

		result := event.Context()

		assert.Equal(t, ctx, result)
		assert.Equal(t, "testValue", result.Value(val))
	})

	t.Run("context_with_metadata", func(t *testing.T) {
		testUser := &sdk.User{Id: "user1", Name: "Test User"}
		contextMetadata := sdk.Metadata{
			User:       testUser,
			ProjectIds: []string{"project1", "project2"},
		}
		ctx := middlewares.AddMetadata(context.Background(), contextMetadata)

		event := newEvent(ctx, goiamuniverse.EventClientCreated, client, metadata)

		result := event.Context()

		assert.Equal(t, ctx, result)
		// Verify the context contains the middleware data
		retrievedProjects := middlewares.GetProjects(result)
		retrievedUser := middlewares.GetUser(result)
		assert.Equal(t, []string{"project1", "project2"}, retrievedProjects)
		assert.Equal(t, testUser, retrievedUser)
	})

	t.Run("context_with_cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		event := newEvent(ctx, goiamuniverse.EventClientUpdated, client, metadata)

		result := event.Context()

		assert.Equal(t, ctx, result)
		// Verify context is still active
		select {
		case <-result.Done():
			t.Error("Context should not be cancelled yet")
		default:
			// Expected - context is still active
		}

		// Cancel the context and verify
		cancel()
		select {
		case <-result.Done():
			// Expected - context is now cancelled
		default:
			t.Error("Context should be cancelled")
		}
	})
}

func TestNewEvent(t *testing.T) {
	t.Run("create_event_with_all_parameters", func(t *testing.T) {
		ctx := context.Background()
		eventName := goiamuniverse.EventClientCreated
		client := sdk.Client{
			Id:        "client1",
			Name:      "Test Client",
			ProjectId: "project1",
		}
		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{
			User:       user,
			ProjectIds: []string{"project1", "project2"},
		}

		event := newEvent(ctx, eventName, client, metadata)

		assert.NotNil(t, event)
		assert.Equal(t, eventName, event.Name())
		assert.Equal(t, client, event.Payload())
		assert.Equal(t, metadata, event.Metadata())
		assert.Equal(t, ctx, event.Context())
	})

	t.Run("create_event_with_minimal_parameters", func(t *testing.T) {
		ctx := context.Background()
		eventName := goiamuniverse.EventClientUpdated
		client := sdk.Client{}
		metadata := sdk.Metadata{}

		event := newEvent(ctx, eventName, client, metadata)

		assert.NotNil(t, event)
		assert.Equal(t, eventName, event.Name())
		assert.Equal(t, client, event.Payload())
		assert.Equal(t, metadata, event.Metadata())
		assert.Equal(t, ctx, event.Context())
	})
}
