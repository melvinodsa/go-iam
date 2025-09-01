package resource

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test helper to create context with metadata
func createTestContext() context.Context {
	metadata := sdk.Metadata{
		User: &sdk.User{
			Id:   "test-user-id",
			Name: "Test User",
		},
		ProjectIds: []string{"test-project-id"},
	}
	return middlewares.AddMetadata(context.Background(), metadata)
}

// MockStore implements Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Search(ctx context.Context, query sdk.ResourceQuery) (*sdk.ResourceList, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.ResourceList), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*sdk.Resource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Resource), args.Error(1)
}

func (m *MockStore) Create(ctx context.Context, resource *sdk.Resource) (string, error) {
	args := m.Called(ctx, resource)
	return args.String(0), args.Error(1)
}

func (m *MockStore) Update(ctx context.Context, resource *sdk.Resource) error {
	args := m.Called(ctx, resource)
	return args.Error(0)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSubscriber implements Subscriber interface for testing
type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) HandleEvent(event utils.Event[sdk.Resource]) {
	m.Called(event)
}

func TestNewService(t *testing.T) {
	mockStore := &MockStore{}

	service := NewService(mockStore)

	assert.NotNil(t, service)
	assert.Implements(t, (*Service)(nil), service)
}

func TestService_Search(t *testing.T) {
	t.Run("successful_search", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := middlewares.AddMetadata(context.Background(), sdk.Metadata{
			ProjectIds: []string{"project1", "project2"},
		})
		query := sdk.ResourceQuery{
			Name:  "test",
			Skip:  0,
			Limit: 10,
		}

		expectedQuery := sdk.ResourceQuery{
			Name:       "test",
			Skip:       0,
			Limit:      10,
			ProjectIds: []string{"project1", "project2"},
		}

		expectedResult := &sdk.ResourceList{
			Resources: []sdk.Resource{
				{
					ID:   "resource1",
					Key:  "users",
					Name: "Users Resource",
				},
			},
			Total: 1,
			Skip:  0,
			Limit: 10,
		}

		mockStore.On("Search", ctx, expectedQuery).Return(expectedResult, nil)

		result, err := service.Search(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := middlewares.AddMetadata(context.Background(), sdk.Metadata{
			ProjectIds: []string{"project1"},
		})
		query := sdk.ResourceQuery{
			Name: "test",
		}

		mockStore.On("Search", ctx, mock.AnythingOfType("sdk.ResourceQuery")).Return(nil, errors.New("store error"))

		result, err := service.Search(ctx, query)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "store error")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resourceId := "resource1"

		expectedResource := &sdk.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		}

		mockStore.On("Get", ctx, resourceId).Return(expectedResource, nil)

		result, err := service.Get(ctx, resourceId)

		assert.NoError(t, err)
		assert.Equal(t, expectedResource, result)
		mockStore.AssertExpectations(t)
	})

	t.Run("resource_not_found", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resourceId := "nonexistent"

		mockStore.On("Get", ctx, resourceId).Return(nil, errors.New("resource not found"))

		result, err := service.Get(ctx, resourceId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "resource not found")
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resourceId := "resource1"

		mockStore.On("Get", ctx, resourceId).Return(nil, errors.New("database error"))

		result, err := service.Get(ctx, resourceId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("successful_create", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := middlewares.AddMetadata(context.Background(), sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		})

		resource := &sdk.Resource{
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		}

		mockStore.On("Create", ctx, resource).Return("resource1", nil)

		err := service.Create(ctx, resource)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resource := &sdk.Resource{
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		}

		mockStore.On("Create", ctx, resource).Return("", errors.New("creation failed"))

		err := service.Create(ctx, resource)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "creation failed")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resource := &sdk.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Updated Users Resource",
			Description: "Updated resource for user management",
		}

		mockStore.On("Update", ctx, resource).Return(nil)

		err := service.Update(ctx, resource)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resource := &sdk.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Updated Users Resource",
			Description: "Updated resource for user management",
		}

		mockStore.On("Update", ctx, resource).Return(errors.New("update failed"))

		err := service.Update(ctx, resource)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resourceId := "resource1"

		mockStore.On("Delete", ctx, resourceId).Return(nil)

		err := service.Delete(ctx, resourceId)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resourceId := "resource1"

		mockStore.On("Delete", ctx, resourceId).Return(errors.New("delete failed"))

		err := service.Delete(ctx, resourceId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Emit(t *testing.T) {
	t.Run("emit_valid_event", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		ctx := createTestContext()
		resource := sdk.Resource{
			ID:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		// This should not panic or cause issues
		service.Emit(event)

		// No assertions needed as Emit doesn't return anything
		// The test passes if no panic occurs
	})

	t.Run("emit_nil_event", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		// This should not panic
		service.Emit(nil)

		// No assertions needed as Emit should handle nil gracefully
	})
}

func TestService_Subscribe(t *testing.T) {
	t.Run("subscribe_to_event", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		mockSubscriber := &MockSubscriber{}
		eventName := goiamuniverse.EventResourceCreated

		// This should not panic
		service.Subscribe(eventName, mockSubscriber)

		// No assertions needed as Subscribe doesn't return anything
		// The test passes if no panic occurs
	})
}

func TestEvent_Name(t *testing.T) {
	t.Run("resource_created_event", func(t *testing.T) {
		ctx := createTestContext()
		resource := sdk.Resource{
			ID:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		assert.Equal(t, goiamuniverse.EventResourceCreated, event.Name())
	})
}

func TestEvent_Payload(t *testing.T) {
	t.Run("simple_resource", func(t *testing.T) {
		ctx := createTestContext()
		resource := sdk.Resource{
			ID:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		payload := event.Payload()
		assert.Equal(t, "resource1", payload.ID)
		assert.Equal(t, "users", payload.Key)
		assert.Equal(t, "Users Resource", payload.Name)
	})

	t.Run("resource_with_description", func(t *testing.T) {
		ctx := createTestContext()
		resource := sdk.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		payload := event.Payload()
		assert.Equal(t, "resource1", payload.ID)
		assert.Equal(t, "users", payload.Key)
		assert.Equal(t, "Users Resource", payload.Name)
		assert.Equal(t, "Resource for user management", payload.Description)
	})
}

func TestEvent_Metadata(t *testing.T) {
	t.Run("metadata_with_user_and_projects", func(t *testing.T) {
		ctx := createTestContext()
		resource := sdk.Resource{
			ID:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1", "project2"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		eventMetadata := event.Metadata()
		assert.Equal(t, "user1", eventMetadata.User.Id)
		assert.Equal(t, []string{"project1", "project2"}, eventMetadata.ProjectIds)
	})
}

func TestEvent_Context(t *testing.T) {
	t.Run("background_context", func(t *testing.T) {
		ctx := createTestContext()
		resource := sdk.Resource{
			ID:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		assert.Equal(t, ctx, event.Context())
	})

	t.Run("context_with_metadata", func(t *testing.T) {
		baseCtx := createTestContext()
		ctx := middlewares.AddMetadata(baseCtx, sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		})

		resource := sdk.Resource{
			ID:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1"},
		}

		event := newEvent(ctx, goiamuniverse.EventResourceCreated, resource, metadata)

		assert.Equal(t, ctx, event.Context())

		// Verify we can retrieve metadata from context
		retrievedMetadata := middlewares.GetMetadata(event.Context())
		assert.Equal(t, "user1", retrievedMetadata.User.Id)
		assert.Equal(t, []string{"project1"}, retrievedMetadata.ProjectIds)
	})
}

func TestNewEvent(t *testing.T) {
	t.Run("create_event_with_all_parameters", func(t *testing.T) {
		ctx := createTestContext()
		eventName := goiamuniverse.EventResourceCreated
		resource := sdk.Resource{
			ID:          "resource1",
			Key:         "users",
			Name:        "Users Resource",
			Description: "Resource for user management",
		}
		metadata := sdk.Metadata{
			User:       &sdk.User{Id: "user1"},
			ProjectIds: []string{"project1", "project2"},
		}

		event := newEvent(ctx, eventName, resource, metadata)

		assert.NotNil(t, event)
		assert.Equal(t, eventName, event.Name())
		assert.Equal(t, resource, event.Payload())
		assert.Equal(t, metadata, event.Metadata())
		assert.Equal(t, ctx, event.Context())
	})
}
