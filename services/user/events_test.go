package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
)

// TestHandleEvent tests the HandleEvent method
func TestHandleEvent(t *testing.T) {
	ctx := createContextWithMetadata()

	testRole := createTestRole()
	testMetadata := middlewares.GetMetadata(ctx)

	t.Run("success - handle role updated event", func(t *testing.T) {
		svc, mockStore, _ := setupUserService()

		// Setup mocks for role update handling
		mockStore.On("GetAll", ctx, mock.MatchedBy(func(query sdk.UserQuery) bool {
			return query.RoleId == testRole.Id
		})).Return(&sdk.UserList{Users: []sdk.User{}, Total: 0}, nil)

		// Create a mock event for role updated
		event := &mockEvent{
			name:     goiamuniverse.EventRoleUpdated,
			payload:  *testRole,
			metadata: testMetadata,
			ctx:      ctx,
		}

		// This should not panic and should handle the event
		assert.NotPanics(t, func() {
			svc.HandleEvent(event)
		})

		mockStore.AssertExpectations(t)
	})

	t.Run("success - handle unknown event", func(t *testing.T) {
		svc, _, _ := setupUserService()

		// Create a mock event for unknown event type
		event := &mockEvent{
			name:     goiamuniverse.EventUserCreated, // Different event type
			payload:  *testRole,
			metadata: testMetadata,
			ctx:      ctx,
		}

		// This should not panic and should return early
		assert.NotPanics(t, func() {
			svc.HandleEvent(event)
		})
	})

	t.Run("success - handle client created event", func(t *testing.T) {
		svc, _, _ := setupUserService()

		// Create a mock event for client created (should be ignored)
		event := &mockEvent{
			name:     goiamuniverse.EventClientCreated,
			payload:  *testRole,
			metadata: testMetadata,
			ctx:      ctx,
		}

		// This should not panic and should return early
		assert.NotPanics(t, func() {
			svc.HandleEvent(event)
		})
	})
}

// TestHandleRoleUpdate tests the handleRoleUpdate method
func TestHandleRoleUpdate(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testRole := createTestRole()
	testMetadata := middlewares.GetMetadata(ctx)

	t.Run("success - handle role update with users", func(t *testing.T) {
		// Setup mock for fetching users with role
		users := []sdk.User{
			*createTestUser(),
			{
				Id:        "user-2",
				Email:     "user2@example.com",
				ProjectId: "project-123",
				Roles: map[string]sdk.UserRole{
					testRole.Id: {Id: testRole.Id, Name: testRole.Name},
				},
				Resources: map[string]sdk.UserResource{},
			},
		}

		userList := &sdk.UserList{
			Users: users,
			Total: int64(len(users)),
		}

		// Mock the GetAll call for fetching users with role
		mockStore.On("GetAll", ctx, mock.MatchedBy(func(query sdk.UserQuery) bool {
			return query.RoleId == testRole.Id && query.Skip == 0 && query.Limit == 10
		})).Return(userList, nil).Once()

		// Mock the GetAll call for second page (empty)
		mockStore.On("GetAll", ctx, mock.MatchedBy(func(query sdk.UserQuery) bool {
			return query.RoleId == testRole.Id && query.Skip == 10 && query.Limit == 10
		})).Return(&sdk.UserList{Users: []sdk.User{}, Total: 0}, nil).Once()

		// Mock Update calls for each user
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil).Times(len(users))

		event := &mockEvent{
			name:     goiamuniverse.EventRoleUpdated,
			payload:  *testRole,
			metadata: testMetadata,
			ctx:      ctx,
		}

		// Should handle the role update without panicking
		assert.NotPanics(t, func() {
			svc.handleRoleUpdate(event)
		})

		mockStore.AssertExpectations(t)
	})

	t.Run("error - fetch users fails", func(t *testing.T) {
		// Mock GetAll to return error
		mockStore.ExpectedCalls = nil
		mockStore.On("GetAll", ctx, mock.AnythingOfType("sdk.UserQuery")).Return((*sdk.UserList)(nil), errors.New("database error"))

		event := &mockEvent{
			name:     goiamuniverse.EventRoleUpdated,
			payload:  *testRole,
			metadata: testMetadata,
			ctx:      ctx,
		}

		// Should handle the error gracefully (logs error but doesn't panic)
		assert.NotPanics(t, func() {
			svc.handleRoleUpdate(event)
		})

		mockStore.AssertExpectations(t)
	})
}

// TestFetchAndUpdateUsersWithRole tests the fetchAndUpdateUsersWithRole method
func TestFetchAndUpdateUsersWithRole(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testRole := createTestRole()

	t.Run("success - fetch and update users with pagination", func(t *testing.T) {
		// Create test users for multiple pages
		users1 := make([]sdk.User, 10) // Full page
		for i := range users1 {
			users1[i] = sdk.User{
				Id:        "user-" + string(rune(i+1)),
				Email:     "user" + string(rune(i+1)) + "@example.com",
				ProjectId: "project-123",
				Roles: map[string]sdk.UserRole{
					testRole.Id: {Id: testRole.Id, Name: testRole.Name},
				},
				Resources: map[string]sdk.UserResource{},
			}
		}

		users2 := []sdk.User{ // Partial page
			{
				Id:        "user-11",
				Email:     "user11@example.com",
				ProjectId: "project-123",
				Roles: map[string]sdk.UserRole{
					testRole.Id: {Id: testRole.Id, Name: testRole.Name},
				},
				Resources: map[string]sdk.UserResource{},
			},
		}

		// Mock first page
		mockStore.On("GetAll", ctx, mock.MatchedBy(func(query sdk.UserQuery) bool {
			return query.RoleId == testRole.Id && query.Skip == 0 && query.Limit == 10
		})).Return(&sdk.UserList{Users: users1, Total: 11}, nil).Once()

		// Mock second page
		mockStore.On("GetAll", ctx, mock.MatchedBy(func(query sdk.UserQuery) bool {
			return query.RoleId == testRole.Id && query.Skip == 10 && query.Limit == 10
		})).Return(&sdk.UserList{Users: users2, Total: 11}, nil).Once()

		// Mock third page (empty)
		mockStore.On("GetAll", ctx, mock.MatchedBy(func(query sdk.UserQuery) bool {
			return query.RoleId == testRole.Id && query.Skip == 20 && query.Limit == 10
		})).Return(&sdk.UserList{Users: []sdk.User{}, Total: 11}, nil).Once()

		// Mock Update calls for all users
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil).Times(11)

		err := svc.fetchAndUpdateUsersWithRole(ctx, *testRole)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("success - no users with role", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		// Mock GetAll to return empty list
		mockStore.On("GetAll", ctx, mock.AnythingOfType("sdk.UserQuery")).Return(&sdk.UserList{Users: []sdk.User{}, Total: 0}, nil)

		err := svc.fetchAndUpdateUsersWithRole(ctx, *testRole)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("error - GetAll fails", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		// Mock GetAll to return error
		mockStore.On("GetAll", ctx, mock.AnythingOfType("sdk.UserQuery")).Return((*sdk.UserList)(nil), errors.New("database error"))

		err := svc.fetchAndUpdateUsersWithRole(ctx, *testRole)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockStore.AssertExpectations(t)
	})

	t.Run("error - updateUsersWithRole fails", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		users := []sdk.User{*createTestUser()}

		// Mock successful GetAll
		mockStore.On("GetAll", ctx, mock.AnythingOfType("sdk.UserQuery")).Return(&sdk.UserList{Users: users, Total: 1}, nil)

		// Mock Update to fail
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("update error"))

		err := svc.fetchAndUpdateUsersWithRole(ctx, *testRole)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
		mockStore.AssertExpectations(t)
	})
}

// TestUpdateUsersWithRole tests the updateUsersWithRole method
func TestUpdateUsersWithRole(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testRole := createTestRole()

	t.Run("success - update multiple users", func(t *testing.T) {
		users := []sdk.User{
			*createTestUser(),
			{
				Id:        "user-2",
				Email:     "user2@example.com",
				ProjectId: "project-123",
				Roles: map[string]sdk.UserRole{
					testRole.Id: {Id: testRole.Id, Name: testRole.Name},
				},
				Resources: map[string]sdk.UserResource{},
			},
		}

		// Mock Update calls for each user
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil).Times(len(users))

		err := svc.updateUsersWithRole(ctx, *testRole, users)

		assert.NoError(t, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("success - update empty users list", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		users := []sdk.User{}

		err := svc.updateUsersWithRole(ctx, *testRole, users)

		assert.NoError(t, err)
		// No store calls should be made
		mockStore.AssertExpectations(t)
	})

	t.Run("error - user update fails", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		users := []sdk.User{*createTestUser()}

		// Mock Update to fail
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("update failed"))

		err := svc.updateUsersWithRole(ctx, *testRole, users)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
		mockStore.AssertExpectations(t)
	})
}

// TestUpdateUser tests the updateUser method
func TestUpdateUser(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testRole := createTestRole()

	t.Run("success - update user with role", func(t *testing.T) {
		user := createTestUser()
		// Add the role to user initially
		user.Roles[testRole.Id] = sdk.UserRole{Id: testRole.Id, Name: testRole.Name}

		// Mock successful Update
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)

		err := svc.updateUser(ctx, *testRole, user)

		assert.NoError(t, err)

		// Verify the role is still in the user (removed and re-added)
		assert.Contains(t, user.Roles, testRole.Id)
		assert.Equal(t, testRole.Id, user.Roles[testRole.Id].Id)
		assert.Equal(t, testRole.Name, user.Roles[testRole.Id].Name)

		mockStore.AssertExpectations(t)
	})

	t.Run("error - store update fails", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		user := createTestUser()

		// Mock Update to fail
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("store error"))

		err := svc.updateUser(ctx, *testRole, user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store error")
		mockStore.AssertExpectations(t)
	})

	t.Run("success - update user without existing role", func(t *testing.T) {
		mockStore.ExpectedCalls = nil
		user := createTestUser()
		// Ensure user doesn't have the role initially
		delete(user.Roles, testRole.Id)

		// Mock successful Update
		mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)

		err := svc.updateUser(ctx, *testRole, user)

		assert.NoError(t, err)

		// Verify the role is added to the user
		assert.Contains(t, user.Roles, testRole.Id)
		assert.Equal(t, testRole.Id, user.Roles[testRole.Id].Id)
		assert.Equal(t, testRole.Name, user.Roles[testRole.Id].Name)

		mockStore.AssertExpectations(t)
	})
}

// mockEvent implements utils.Event[sdk.Role] for testing
type mockEvent struct {
	name     goiamuniverse.Event
	payload  sdk.Role
	metadata sdk.Metadata
	ctx      context.Context
}

func (e *mockEvent) Name() goiamuniverse.Event {
	return e.name
}

func (e *mockEvent) Payload() sdk.Role {
	return e.payload
}

func (e *mockEvent) Metadata() sdk.Metadata {
	return e.metadata
}

func (e *mockEvent) Context() context.Context {
	return e.ctx
}
