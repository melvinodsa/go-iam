package role

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStore implements Store interface for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Create(ctx context.Context, role *sdk.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, role *sdk.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockStore) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Role), args.Error(1)
}

func (m *MockStore) GetAll(ctx context.Context, query sdk.RoleQuery) (*sdk.RoleList, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.RoleList), args.Error(1)
}

func (m *MockStore) RemoveResourceFromAll(ctx context.Context, resourceKey string) error {
	args := m.Called(ctx, resourceKey)
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	mockStore := &MockStore{}

	service := NewService(mockStore)

	assert.NotNil(t, service)
	assert.Implements(t, (*Service)(nil), service)
}

func TestService_Create(t *testing.T) {
	t.Run("successful_create", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		role := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockStore.On("Create", ctx, role).Return(nil)

		err := service.Create(ctx, role)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		role := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockStore.On("Create", ctx, role).Return(errors.New("database error"))

		err := service.Create(ctx, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		role := &sdk.Role{
			Id:          "role1",
			Name:        "Updated Role",
			Description: "An updated role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockStore.On("Update", ctx, role).Return(nil)

		err := service.Update(ctx, role)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		role := &sdk.Role{
			Id:          "role1",
			Name:        "Updated Role",
			Description: "An updated role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockStore.On("Update", ctx, role).Return(errors.New("database error"))

		err := service.Update(ctx, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update role")
		assert.Contains(t, err.Error(), "database error")
		mockStore.AssertExpectations(t)
	})
}

func TestService_GetById(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		expectedRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		mockStore.On("GetById", ctx, "role1").Return(expectedRole, nil)

		result, err := service.GetById(ctx, "role1")

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, result)
		mockStore.AssertExpectations(t)
	})

	t.Run("role_not_found", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		mockStore.On("GetById", ctx, "nonexistent").Return(nil, sdk.ErrRoleNotFound)

		result, err := service.GetById(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, sdk.ErrRoleNotFound, err)
		assert.Nil(t, result)
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		mockStore.On("GetById", ctx, "role1").Return(nil, errors.New("database error"))

		result, err := service.GetById(ctx, "role1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		mockStore.AssertExpectations(t)
	})
}

func TestService_GetAll(t *testing.T) {
	t.Run("successful_get_all", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1", "project2"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		query := sdk.RoleQuery{
			SearchQuery: "test",
			Skip:        0,
			Limit:       10,
		}

		expectedRoles := &sdk.RoleList{
			Roles: []sdk.Role{
				{
					Id:          "role1",
					Name:        "Test Role 1",
					Description: "A test role",
					ProjectId:   "project1",
					Enabled:     true,
				},
				{
					Id:          "role2",
					Name:        "Test Role 2",
					Description: "Another test role",
					ProjectId:   "project2",
					Enabled:     false,
				},
			},
			Total: 2,
			Skip:  0,
			Limit: 10,
		}

		expectedQuery := query
		expectedQuery.ProjectIds = []string{"project1", "project2"}

		mockStore.On("GetAll", ctx, expectedQuery).Return(expectedRoles, nil)

		result, err := service.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, result)
		assert.Equal(t, 2, len(result.Roles))
		mockStore.AssertExpectations(t)
	})

	t.Run("empty_result", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		query := sdk.RoleQuery{
			SearchQuery: "nonexistent",
			Skip:        0,
			Limit:       10,
		}

		expectedRoles := &sdk.RoleList{
			Roles: []sdk.Role{},
			Total: 0,
			Skip:  0,
			Limit: 10,
		}

		expectedQuery := query
		expectedQuery.ProjectIds = []string{"project1"}

		mockStore.On("GetAll", ctx, expectedQuery).Return(expectedRoles, nil)

		result, err := service.GetAll(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, result)
		assert.Equal(t, 0, len(result.Roles))
		mockStore.AssertExpectations(t)
	})

	t.Run("store_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		query := sdk.RoleQuery{
			SearchQuery: "test",
			Skip:        0,
			Limit:       10,
		}

		expectedQuery := query
		expectedQuery.ProjectIds = []string{"project1"}

		mockStore.On("GetAll", ctx, expectedQuery).Return(nil, errors.New("database error"))

		result, err := service.GetAll(ctx, query)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		mockStore.AssertExpectations(t)
	})
}

func TestService_AddResource(t *testing.T) {
	t.Run("successful_add_resource_to_new_role", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		existingRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
			Resources:   nil, // No existing resources
		}

		resource := sdk.Resources{
			Id:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}

		expectedUpdatedRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
			Resources: map[string]sdk.Resources{
				"users": resource,
			},
		}

		mockStore.On("GetById", ctx, "role1").Return(existingRole, nil)
		mockStore.On("Update", ctx, expectedUpdatedRole).Return(nil)

		err := service.AddResource(ctx, "role1", resource)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("successful_add_resource_to_existing_resources", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		existingRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
			Resources: map[string]sdk.Resources{
				"projects": {
					Id:   "resource0",
					Key:  "projects",
					Name: "Projects Resource",
				},
			},
		}

		resource := sdk.Resources{
			Id:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}

		expectedUpdatedRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
			Resources: map[string]sdk.Resources{
				"projects": {
					Id:   "resource0",
					Key:  "projects",
					Name: "Projects Resource",
				},
				"users": resource,
			},
		}

		mockStore.On("GetById", ctx, "role1").Return(existingRole, nil)
		mockStore.On("Update", ctx, expectedUpdatedRole).Return(nil)

		err := service.AddResource(ctx, "role1", resource)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("role_not_found", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		resource := sdk.Resources{
			Id:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}

		mockStore.On("GetById", ctx, "nonexistent").Return(nil, sdk.ErrRoleNotFound)

		err := service.AddResource(ctx, "nonexistent", resource)

		assert.Error(t, err)
		assert.Equal(t, sdk.ErrRoleNotFound, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("get_role_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		resource := sdk.Resources{
			Id:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}

		mockStore.On("GetById", ctx, "role1").Return(nil, errors.New("database error"))

		err := service.AddResource(ctx, "role1", resource)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockStore.AssertExpectations(t)
	})

	t.Run("update_role_error", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := middlewares.AddMetadata(context.Background(), metadata)

		existingRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
			Resources:   nil,
		}

		resource := sdk.Resources{
			Id:   "resource1",
			Key:  "users",
			Name: "Users Resource",
		}

		expectedUpdatedRole := &sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
			Resources: map[string]sdk.Resources{
				"users": resource,
			},
		}

		mockStore.On("GetById", ctx, "role1").Return(existingRole, nil)
		mockStore.On("Update", ctx, expectedUpdatedRole).Return(errors.New("update failed"))

		err := service.AddResource(ctx, "role1", resource)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update role")
		assert.Contains(t, err.Error(), "update failed")
		mockStore.AssertExpectations(t)
	})
}

func TestService_Emit(t *testing.T) {
	t.Run("emit_valid_event", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		role := sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{User: user, ProjectIds: []string{"project1"}}
		ctx := context.Background()

		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		// This should not panic
		assert.NotPanics(t, func() {
			service.Emit(event)
		})
	})

	t.Run("emit_nil_event", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		// This should not panic
		assert.NotPanics(t, func() {
			service.Emit(nil)
		})
	})
}

func TestService_Subscribe(t *testing.T) {
	t.Run("subscribe_to_event", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)

		// Create a mock subscriber that implements the Subscriber interface
		mockSubscriber := &MockSubscriber{}

		// This should not panic
		assert.NotPanics(t, func() {
			service.Subscribe(goiamuniverse.EventRoleUpdated, mockSubscriber)
		})
	})
}

func TestService_RemoveResourceFromAll(t *testing.T) {
	t.Run("successful_remove_resource_from_all_roles", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		resourceKey := "users"

		mockStore.On("RemoveResourceFromAll", ctx, resourceKey).Return(nil)

		err := service.RemoveResourceFromAll(ctx, resourceKey)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("failed_remove_resource_from_all_roles", func(t *testing.T) {
		mockStore := &MockStore{}
		service := NewService(mockStore)
		ctx := context.Background()

		resourceKey := "users"

		mockStore.On("RemoveResourceFromAll", ctx, resourceKey).Return(errors.New("database error"))

		err := service.RemoveResourceFromAll(ctx, resourceKey)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockStore.AssertExpectations(t)
	})
}

// MockSubscriber implements Subscriber interface for testing
type MockSubscriber struct {
	mock.Mock
}

func (m *MockSubscriber) HandleEvent(event utils.Event[sdk.Role]) {
	m.Called(event)
}

func TestEvent_Name(t *testing.T) {
	ctx := context.Background()
	role := sdk.Role{Id: "role1", Name: "Test Role"}
	metadata := sdk.Metadata{User: &sdk.User{Id: "user1"}, ProjectIds: []string{"project1"}}

	t.Run("role_updated_event", func(t *testing.T) {
		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		result := event.Name()

		assert.Equal(t, goiamuniverse.EventRoleUpdated, result)
	})
}

func TestEvent_Payload(t *testing.T) {
	ctx := context.Background()
	metadata := sdk.Metadata{User: &sdk.User{Id: "user1"}, ProjectIds: []string{"project1"}}

	t.Run("simple_role", func(t *testing.T) {
		role := sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}

		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		result := event.Payload()

		assert.Equal(t, role, result)
		assert.Equal(t, "role1", result.Id)
		assert.Equal(t, "Test Role", result.Name)
		assert.Equal(t, "A test role", result.Description)
		assert.Equal(t, "project1", result.ProjectId)
		assert.True(t, result.Enabled)
	})

	t.Run("role_with_resources", func(t *testing.T) {
		createdAt := time.Now()
		updatedAt := time.Now().Add(time.Hour)

		role := sdk.Role{
			Id:          "role1",
			Name:        "Complex Role",
			Description: "A complex test role",
			ProjectId:   "project1",
			Resources: map[string]sdk.Resources{
				"users": {
					Id:   "resource1",
					Key:  "users",
					Name: "Users Resource",
				},
				"projects": {
					Id:   "resource2",
					Key:  "projects",
					Name: "Projects Resource",
				},
			},
			Enabled:   true,
			CreatedAt: &createdAt,
			CreatedBy: "user1",
			UpdatedAt: &updatedAt,
			UpdatedBy: "user2",
		}

		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		result := event.Payload()

		assert.Equal(t, role, result)
		assert.Equal(t, "role1", result.Id)
		assert.Equal(t, "Complex Role", result.Name)
		assert.Equal(t, "A complex test role", result.Description)
		assert.Equal(t, "project1", result.ProjectId)
		assert.True(t, result.Enabled)
		assert.Equal(t, 2, len(result.Resources))
		assert.Contains(t, result.Resources, "users")
		assert.Contains(t, result.Resources, "projects")
		assert.Equal(t, "Users Resource", result.Resources["users"].Name)
		assert.Equal(t, "Projects Resource", result.Resources["projects"].Name)
		assert.Equal(t, &createdAt, result.CreatedAt)
		assert.Equal(t, "user1", result.CreatedBy)
		assert.Equal(t, &updatedAt, result.UpdatedAt)
		assert.Equal(t, "user2", result.UpdatedBy)
	})
}

func TestEvent_Metadata(t *testing.T) {
	ctx := context.Background()
	role := sdk.Role{Id: "role1", Name: "Test Role"}

	t.Run("metadata_with_user_and_projects", func(t *testing.T) {
		user := &sdk.User{
			Id:    "user1",
			Name:  "Test User",
			Email: "test@example.com",
		}
		metadata := sdk.Metadata{
			User:       user,
			ProjectIds: []string{"project1", "project2"},
		}

		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		result := event.Metadata()

		assert.Equal(t, metadata, result)
		assert.Equal(t, user, result.User)
		assert.Equal(t, []string{"project1", "project2"}, result.ProjectIds)
		assert.Equal(t, "user1", result.User.Id)
		assert.Equal(t, "Test User", result.User.Name)
		assert.Equal(t, "test@example.com", result.User.Email)
	})
}

func TestEvent_Context(t *testing.T) {
	role := sdk.Role{Id: "role1", Name: "Test Role"}
	metadata := sdk.Metadata{User: &sdk.User{Id: "user1"}, ProjectIds: []string{"project1"}}

	t.Run("background_context", func(t *testing.T) {
		ctx := context.Background()

		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		result := event.Context()

		assert.Equal(t, ctx, result)
		assert.NotNil(t, result)
	})

	t.Run("context_with_metadata", func(t *testing.T) {
		testUser := &sdk.User{Id: "user1", Name: "Test User"}
		contextMetadata := sdk.Metadata{
			User:       testUser,
			ProjectIds: []string{"project1", "project2"},
		}
		ctx := middlewares.AddMetadata(context.Background(), contextMetadata)

		event := newEvent(ctx, goiamuniverse.EventRoleUpdated, role, metadata)

		result := event.Context()

		assert.Equal(t, ctx, result)
		// Verify the context contains the middleware data
		retrievedProjects := middlewares.GetProjects(result)
		retrievedUser := middlewares.GetUser(result)
		assert.Equal(t, []string{"project1", "project2"}, retrievedProjects)
		assert.Equal(t, testUser, retrievedUser)
	})
}

func TestNewEvent(t *testing.T) {
	t.Run("create_event_with_all_parameters", func(t *testing.T) {
		ctx := context.Background()
		eventName := goiamuniverse.EventRoleUpdated
		role := sdk.Role{
			Id:          "role1",
			Name:        "Test Role",
			Description: "A test role",
			ProjectId:   "project1",
			Enabled:     true,
		}
		user := &sdk.User{Id: "user1", Name: "Test User"}
		metadata := sdk.Metadata{
			User:       user,
			ProjectIds: []string{"project1", "project2"},
		}

		event := newEvent(ctx, eventName, role, metadata)

		assert.NotNil(t, event)
		assert.Equal(t, eventName, event.Name())
		assert.Equal(t, role, event.Payload())
		assert.Equal(t, metadata, event.Metadata())
		assert.Equal(t, ctx, event.Context())
	})
}
