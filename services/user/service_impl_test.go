package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/melvinodsa/go-iam/middlewares"
	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test/services"
)

// MockStore is a mock implementation of the Store interface
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Create(ctx context.Context, user *sdk.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, user *sdk.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStore) GetByEmail(ctx context.Context, email string, projectId string) (*sdk.User, error) {
	args := m.Called(ctx, email, projectId)
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockStore) GetById(ctx context.Context, id string) (*sdk.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockStore) GetByPhone(ctx context.Context, phone string, projectId string) (*sdk.User, error) {
	args := m.Called(ctx, phone, projectId)
	return args.Get(0).(*sdk.User), args.Error(1)
}

func (m *MockStore) GetAll(ctx context.Context, query sdk.UserQuery) (*sdk.UserList, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sdk.UserList), args.Error(1)
}

func (m *MockStore) RemoveResourceFromAll(ctx context.Context, resourceKey string) error {
	args := m.Called(ctx, resourceKey)
	return args.Error(0)
}

// Test helper to create a test user
func createTestUser() *sdk.User {
	now := time.Now()
	return &sdk.User{
		Id:        "user-123",
		Email:     "test@example.com",
		Phone:     "+1234567890",
		Name:      "Test User",
		ProjectId: "project-123",
		Enabled:   true,
		Expiry:    nil,
		Roles:     make(map[string]sdk.UserRole),
		Resources: make(map[string]sdk.UserResource),
		Policies:  make(map[string]sdk.UserPolicy),
		CreatedAt: &now,
		CreatedBy: "admin",
		UpdatedAt: &now,
		UpdatedBy: "admin",
	}
}

// Test helper to create a test role
func createTestRole() *sdk.Role {
	now := time.Now()
	return &sdk.Role{
		Id:        "role-123",
		Name:      "Test Role",
		ProjectId: "project-123",
		Resources: make(map[string]sdk.Resources),
		CreatedAt: &now,
		CreatedBy: "admin",
		UpdatedAt: &now,
		UpdatedBy: "admin",
	}
}

func setupUserService() (*service, *MockStore, *services.MockRoleService) {
	mockStore := &MockStore{}
	mockRoleService := &services.MockRoleService{}

	svc := &service{
		store:   mockStore,
		roleSvc: mockRoleService,
		e:       utils.NewEmitter[utils.Event[sdk.User]](),
	}

	return svc, mockStore, mockRoleService
}

// Helper function to create a context with metadata
func createContextWithMetadata() context.Context {
	metadata := sdk.Metadata{
		User:       createTestUser(),
		ProjectIds: []string{"project-123"},
	}
	return middlewares.AddMetadata(context.Background(), metadata)
}

// TestNewService tests the NewService constructor
func TestNewService(t *testing.T) {
	mockStore := &MockStore{}
	mockRoleService := &services.MockRoleService{}

	svc := NewService(mockStore, mockRoleService)

	require.NotNil(t, svc)
	assert.IsType(t, &service{}, svc)
}

// TestCreate tests the Create method
func TestCreate(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	tests := []struct {
		name          string
		user          *sdk.User
		setupMocks    func()
		expectedError string
	}{
		{
			name: "success - create user",
			user: createTestUser(),
			setupMocks: func() {
				mockStore.On("Create", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name: "error - store create fails",
			user: createTestUser(),
			setupMocks: func() {
				mockStore.On("Create", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.Create(ctx, tt.user)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	tests := []struct {
		name          string
		user          *sdk.User
		setupMocks    func()
		expectedError string
	}{
		{
			name: "success - update user",
			user: createTestUser(),
			setupMocks: func() {
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name: "error - store update fails",
			user: createTestUser(),
			setupMocks: func() {
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("update failed"))
			},
			expectedError: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.Update(ctx, tt.user)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestGetByEmail tests the GetByEmail method
func TestGetByEmail(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUser := createTestUser()

	tests := []struct {
		name          string
		email         string
		projectId     string
		setupMocks    func()
		expectedUser  *sdk.User
		expectedError string
	}{
		{
			name:      "success - user found",
			email:     "test@example.com",
			projectId: "project-123",
			setupMocks: func() {
				mockStore.On("GetByEmail", ctx, "test@example.com", "project-123").Return(testUser, nil)
			},
			expectedUser: testUser,
		},
		{
			name:      "error - user not found",
			email:     "notfound@example.com",
			projectId: "project-123",
			setupMocks: func() {
				mockStore.On("GetByEmail", ctx, "notfound@example.com", "project-123").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			user, err := svc.GetByEmail(ctx, tt.email, tt.projectId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestGetById tests the GetById method
func TestGetById(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUser := createTestUser()

	tests := []struct {
		name          string
		userId        string
		setupMocks    func()
		expectedUser  *sdk.User
		expectedError string
	}{
		{
			name:   "success - user found",
			userId: "user-123",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
			},
			expectedUser: testUser,
		},
		{
			name:   "error - user not found",
			userId: "user-999",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			user, err := svc.GetById(ctx, tt.userId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestGetByPhone tests the GetByPhone method
func TestGetByPhone(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUser := createTestUser()

	tests := []struct {
		name          string
		phone         string
		projectId     string
		setupMocks    func()
		expectedUser  *sdk.User
		expectedError string
	}{
		{
			name:      "success - user found",
			phone:     "+1234567890",
			projectId: "project-123",
			setupMocks: func() {
				mockStore.On("GetByPhone", ctx, "+1234567890", "project-123").Return(testUser, nil)
			},
			expectedUser: testUser,
		},
		{
			name:      "error - user not found",
			phone:     "+9999999999",
			projectId: "project-123",
			setupMocks: func() {
				mockStore.On("GetByPhone", ctx, "+9999999999", "project-123").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			user, err := svc.GetByPhone(ctx, tt.phone, tt.projectId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestGetAll tests the GetAll method
func TestGetAll(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUsers := &sdk.UserList{
		Users: []sdk.User{*createTestUser()},
		Total: 1,
	}

	query := sdk.UserQuery{
		ProjectIds: []string{"project-123"},
		Limit:      10,
		Skip:       0,
	}

	tests := []struct {
		name          string
		query         sdk.UserQuery
		setupMocks    func()
		expectedUsers *sdk.UserList
		expectedError string
	}{
		{
			name:  "success - users found",
			query: query,
			setupMocks: func() {
				mockStore.On("GetAll", ctx, query).Return(testUsers, nil)
			},
			expectedUsers: testUsers,
		},
		{
			name:  "error - store query fails",
			query: query,
			setupMocks: func() {
				mockStore.On("GetAll", ctx, query).Return((*sdk.UserList)(nil), errors.New("query failed"))
			},
			expectedError: "query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			users, err := svc.GetAll(ctx, tt.query)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, users)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUsers, users)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestAddRoleToUser tests the AddRoleToUser method
func TestAddRoleToUser(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, mockRoleService := setupUserService()

	testUser := createTestUser()
	testRole := createTestRole()

	tests := []struct {
		name          string
		userId        string
		roleId        string
		setupMocks    func()
		expectedError string
	}{
		{
			name:   "success - add role to user",
			userId: "user-123",
			roleId: "role-123",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockRoleService.On("GetById", ctx, "role-123").Return(testRole, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:   "error - user not found",
			userId: "user-999",
			roleId: "role-123",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:   "error - role not found",
			userId: "user-123",
			roleId: "role-999",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockRoleService.On("GetById", ctx, "role-999").Return((*sdk.Role)(nil), errors.New("role not found"))
			},
			expectedError: "role not found",
		},
		{
			name:   "error - store update fails",
			userId: "user-123",
			roleId: "role-123",
			setupMocks: func() {
				// Use a fresh user without the role so the update will be called
				freshUser := createTestUser()
				freshUser.Roles = make(map[string]sdk.UserRole) // Ensure no existing roles
				mockStore.On("GetById", ctx, "user-123").Return(freshUser, nil)
				mockRoleService.On("GetById", ctx, "role-123").Return(testRole, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database connection failed"))
			},
			expectedError: "failed to add role to user",
		},
		{
			name:   "error - empty user ID",
			userId: "",
			roleId: "role-123",
			setupMocks: func() {
				// No mocks needed as validation happens before store calls
			},
			expectedError: "user ID and role ID are required",
		},
		{
			name:   "error - empty role ID",
			userId: "user-123",
			roleId: "",
			setupMocks: func() {
				// No mocks needed as validation happens before store calls
			},
			expectedError: "user ID and role ID are required",
		},
		{
			name:   "success - role already exists (no-op)",
			userId: "user-123",
			roleId: "role-123",
			setupMocks: func() {
				// Create a user that already has this role
				userWithRole := createTestUser()
				userWithRole.Roles["role-123"] = sdk.UserRole{Id: "role-123", Name: "Test Role"}
				mockStore.On("GetById", ctx, "user-123").Return(userWithRole, nil)
				mockRoleService.On("GetById", ctx, "role-123").Return(testRole, nil)
				// No Update call expected since role already exists
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil
			mockRoleService.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.AddRoleToUser(ctx, tt.userId, tt.roleId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
			mockRoleService.AssertExpectations(t)
		})
	}
}

// TestRemoveRoleFromUser tests the RemoveRoleFromUser method
func TestRemoveRoleFromUser(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, mockRoleService := setupUserService()

	testUser := createTestUser()
	testUser.Roles["role-123"] = sdk.UserRole{Id: "role-123", Name: "Test Role"}
	testRole := createTestRole()

	tests := []struct {
		name          string
		userId        string
		roleId        string
		setupMocks    func()
		expectedError string
	}{
		{
			name:   "success - remove role from user",
			userId: "user-123",
			roleId: "role-123",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockRoleService.On("GetById", ctx, "role-123").Return(testRole, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:   "error - user not found",
			userId: "user-999",
			roleId: "role-123",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:   "success - role not in user (no-op)",
			userId: "user-123",
			roleId: "role-999",
			setupMocks: func() {
				// User doesn't have this role, so should return without doing anything
				mockStore.On("GetById", ctx, "user-123").Return(createTestUser(), nil)
			},
		},
		{
			name:          "error - empty user ID",
			userId:        "",
			roleId:        "role-123",
			setupMocks:    func() {},
			expectedError: "user ID and role ID are required",
		},
		{
			name:          "error - empty role ID",
			userId:        "user-123",
			roleId:        "",
			setupMocks:    func() {},
			expectedError: "user ID and role ID are required",
		},
		{
			name:   "error - role service fails",
			userId: "user-123",
			roleId: "role-123",
			setupMocks: func() {
				userWithRole := createTestUser()
				userWithRole.Roles["role-123"] = sdk.UserRole{Id: "role-123", Name: "Test Role"}
				mockStore.On("GetById", ctx, "user-123").Return(userWithRole, nil)
				mockRoleService.On("GetById", ctx, "role-123").Return((*sdk.Role)(nil), errors.New("role service error"))
			},
			expectedError: "role service error",
		},
		{
			name:   "error - store update fails",
			userId: "user-123",
			roleId: "role-123",
			setupMocks: func() {
				userWithRole := createTestUser()
				userWithRole.Roles["role-123"] = sdk.UserRole{Id: "role-123", Name: "Test Role"}
				mockStore.On("GetById", ctx, "user-123").Return(userWithRole, nil)
				mockRoleService.On("GetById", ctx, "role-123").Return(testRole, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error"))
			},
			expectedError: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil
			mockRoleService.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.RemoveRoleFromUser(ctx, tt.userId, tt.roleId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
			if len(mockRoleService.ExpectedCalls) > 0 {
				mockRoleService.AssertExpectations(t)
			}
		})
	}
}

// TestAddResourceToUser tests the AddResourceToUser method
func TestAddResourceToUser(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUser := createTestUser()
	resourceRequest := sdk.AddUserResourceRequest{
		Key:  "resource-key",
		Name: "Resource Name",
	}

	tests := []struct {
		name          string
		userId        string
		request       sdk.AddUserResourceRequest
		setupMocks    func()
		expectedError string
	}{
		{
			name:    "success - add resource to user",
			userId:  "user-123",
			request: resourceRequest,
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:    "error - user not found",
			userId:  "user-999",
			request: resourceRequest,
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:    "error - store update fails",
			userId:  "user-123",
			request: resourceRequest,
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error"))
			},
			expectedError: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.AddResourceToUser(ctx, tt.userId, tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestAddPolicyToUser tests the AddPolicyToUser method
func TestAddPolicyToUser(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUser := createTestUser()
	policies := map[string]sdk.UserPolicy{
		"test-policy": {
			Name: "test-policy",
			Mapping: sdk.UserPolicyMapping{
				Arguments: map[string]sdk.UserPolicyMappingValue{
					"key": {Static: "value"},
				},
			},
		},
	}

	tests := []struct {
		name          string
		userId        string
		policies      map[string]sdk.UserPolicy
		setupMocks    func()
		expectedError string
	}{
		{
			name:     "success - add policies to user",
			userId:   "user-123",
			policies: policies,
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:     "error - user not found",
			userId:   "user-999",
			policies: policies,
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:     "error - store update fails",
			userId:   "user-123",
			policies: policies,
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error"))
			},
			expectedError: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.AddPolicyToUser(ctx, tt.userId, tt.policies)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestRemovePolicyFromUser tests the RemovePolicyFromUser method
func TestRemovePolicyFromUser(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	testUser := createTestUser()
	testUser.Policies["test-policy"] = sdk.UserPolicy{Name: "test-policy"}

	tests := []struct {
		name          string
		userId        string
		policyIds     []string
		setupMocks    func()
		expectedError string
	}{
		{
			name:      "success - remove policies from user",
			userId:    "user-123",
			policyIds: []string{"test-policy"},
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:      "error - user not found",
			userId:    "user-999",
			policyIds: []string{"test-policy"},
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:      "error - store update fails",
			userId:    "user-123",
			policyIds: []string{"test-policy"},
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(testUser, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error"))
			},
			expectedError: "failed to update user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.RemovePolicyFromUser(ctx, tt.userId, tt.policyIds)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestRemoveResourceFromAllUsers(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, _ := setupUserService()

	tests := []struct {
		name          string
		resourceKey   string
		setupMocks    func()
		expectedError string
	}{
		{
			name:        "success - remove resource from all users",
			resourceKey: "resource-key",
			setupMocks: func() {
				mockStore.On("RemoveResourceFromAll", ctx, "resource-key").Return(nil)
			},
		},
		{
			name:        "error - store removal fails",
			resourceKey: "resource-key",
			setupMocks: func() {
				mockStore.On("RemoveResourceFromAll", ctx, "resource-key").Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.RemoveResourceFromAll(ctx, tt.resourceKey)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

func TestTransferOwnership(t *testing.T) {
	ctx := createContextWithMetadata()
	svc, mockStore, roleSvc := setupUserService()

	oldOwner := createTestUser()
	newOwner := createTestUser()
	newOwner.Id = "user-456"
	newOwner.Email = "new-owner@example.com"
	newOwner.Phone = "+1987654321"

	tests := []struct {
		name          string
		oldOwnerId    string
		newOwnerId    string
		setupMocks    func()
		expectedError string
	}{
		{
			name:       "success - transfer ownership",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(oldOwner, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newOwner, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:       "success - transfer ownership with policies",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				usr := createTestUser()
				usr.Policies = map[string]sdk.UserPolicy{
					"policy-1": {
						Name: "policy-1",
					},
				}
				mockStore.On("GetById", ctx, "user-123").Return(usr, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newOwner, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:       "success - transfer ownership with roles",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				usr := createTestUser()
				usr.Roles = map[string]sdk.UserRole{
					"role-1": {
						Name: "role-1",
					},
					"role-2": {
						Name: "role-2",
					},
				}
				newUsr := createTestUser()
				newUsr.Roles = map[string]sdk.UserRole{
					"role-2": {
						Name: "role-2",
					},
				}
				roleSvc.On("GetById", ctx, "role-1").Return(createTestRole(), nil)
				mockStore.On("GetById", ctx, "user-123").Return(usr, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newUsr, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:       "success - transfer ownership with resources",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				usr := createTestUser()
				usr.Resources = map[string]sdk.UserResource{
					"resource-1": {
						Name: "resource-1",
					},
					"resource-2": {
						Name:      "resource-2",
						PolicyIds: map[string]bool{"policy-1": true},
					},
				}
				newUsr := createTestUser()
				newUsr.Resources = map[string]sdk.UserResource{
					"resource-2": {
						Name:      "resource-2",
						PolicyIds: map[string]bool{"policy-2": true},
					},
				}
				mockStore.On("GetById", ctx, "user-123").Return(usr, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newUsr, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:       "success - transfer ownership with roles not found",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				usr := createTestUser()
				usr.Roles = map[string]sdk.UserRole{
					"role-2": {
						Name: "role-2",
					},
				}
				roleSvc.On("GetById", ctx, "role-2").Return(createTestRole(), errors.New("role not found"))
				mockStore.On("GetById", ctx, "user-123").Return(usr, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newOwner, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil)
			},
		},
		{
			name:       "error - old owner not found",
			oldOwnerId: "user-999",
			newOwnerId: "user-456",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:       "error - new owner not found",
			oldOwnerId: "user-123",
			newOwnerId: "user-999",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(oldOwner, nil)
				mockStore.On("GetById", ctx, "user-999").Return((*sdk.User)(nil), ErrorUserNotFound)
			},
			expectedError: "user not found",
		},
		{
			name:       "error - store update fails",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(oldOwner, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newOwner, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error"))
			},
			expectedError: "failed to update new owner",
		},
		{
			name:       "error - old user store update fails",
			oldOwnerId: "user-123",
			newOwnerId: "user-456",
			setupMocks: func() {
				mockStore.On("GetById", ctx, "user-123").Return(oldOwner, nil)
				mockStore.On("GetById", ctx, "user-456").Return(newOwner, nil)
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(nil).Once()                          // First call succeeds
				mockStore.On("Update", ctx, mock.AnythingOfType("*sdk.User")).Return(errors.New("database error")).Once() // Second call fails
			},
			expectedError: "failed to update old user",
		},
		{
			name:          "error - empty old owner ID",
			oldOwnerId:    "",
			newOwnerId:    "user-456",
			setupMocks:    func() {},
			expectedError: "user ID and new owner ID are required",
		},
		{
			name:          "error - empty new owner ID",
			oldOwnerId:    "user-123",
			newOwnerId:    "",
			setupMocks:    func() {},
			expectedError: "user ID and new owner ID are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockStore.ExpectedCalls = nil

			tt.setupMocks()

			err := svc.TransferOwnership(ctx, tt.oldOwnerId, tt.newOwnerId)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockStore.AssertExpectations(t)
		})
	}
}

// TestEventHandling tests the event handling functionality
func TestEventHandling(t *testing.T) {
	svc, _, _ := setupUserService()

	// Test that the service implements the Emitter interface
	var _ utils.Emitter[utils.Event[sdk.User], sdk.User] = svc

	// We can't easily test the actual emitter functionality without
	// creating complex event implementations, but we can verify the
	// service has the required interface methods
	assert.NotNil(t, svc)
}

// TestEventMethods tests the event struct methods: Metadata, Payload, Context, Name
func TestEventMethods(t *testing.T) {
	ctx := createContextWithMetadata()
	testUser := createTestUser()
	testMetadata := middlewares.GetMetadata(ctx)
	eventName := goiamuniverse.EventUserCreated

	// Create a new event using the newEvent function
	event := newEvent(ctx, eventName, *testUser, testMetadata)

	t.Run("Name method", func(t *testing.T) {
		assert.Equal(t, eventName, event.Name())
	})

	t.Run("Payload method", func(t *testing.T) {
		payload := event.Payload()
		assert.Equal(t, testUser.Id, payload.Id)
		assert.Equal(t, testUser.Email, payload.Email)
		assert.Equal(t, testUser.Phone, payload.Phone)
		assert.Equal(t, testUser.ProjectId, payload.ProjectId)
	})

	t.Run("Metadata method", func(t *testing.T) {
		metadata := event.Metadata()
		assert.Equal(t, testMetadata, metadata)
	})

	t.Run("Context method", func(t *testing.T) {
		eventCtx := event.Context()
		assert.Equal(t, ctx, eventCtx)
		// Verify the context contains the expected metadata
		assert.Equal(t, testMetadata, middlewares.GetMetadata(eventCtx))
	})

	t.Run("Event with empty metadata", func(t *testing.T) {
		emptyCtx := context.Background()
		emptyMetadata := sdk.Metadata{}
		emptyEvent := newEvent(emptyCtx, eventName, *testUser, emptyMetadata)

		assert.Equal(t, eventName, emptyEvent.Name())
		assert.Equal(t, *testUser, emptyEvent.Payload())
		assert.Equal(t, emptyMetadata, emptyEvent.Metadata())
		assert.Equal(t, emptyCtx, emptyEvent.Context())
	})

	t.Run("Event with different event types", func(t *testing.T) {
		updateEvent := newEvent(ctx, goiamuniverse.EventRoleUpdated, *testUser, testMetadata)
		assert.Equal(t, goiamuniverse.EventRoleUpdated, updateEvent.Name())
		assert.Equal(t, *testUser, updateEvent.Payload())
		assert.Equal(t, testMetadata, updateEvent.Metadata())
		assert.Equal(t, ctx, updateEvent.Context())
	})
}

// TestEmit tests the Emit method
func TestEmit(t *testing.T) {
	svc, _, _ := setupUserService()
	ctx := createContextWithMetadata()
	testUser := createTestUser()
	testMetadata := middlewares.GetMetadata(ctx)

	t.Run("success - emit valid event", func(t *testing.T) {
		event := newEvent(ctx, goiamuniverse.EventUserCreated, *testUser, testMetadata)

		// This should not panic or return error
		assert.NotPanics(t, func() {
			svc.Emit(event)
		})
	})

	t.Run("success - emit nil event", func(t *testing.T) {
		// This should handle nil gracefully and not panic
		assert.NotPanics(t, func() {
			svc.Emit(nil)
		})
	})

	t.Run("success - emit event with different types", func(t *testing.T) {
		// Test with different event types
		events := []goiamuniverse.Event{
			goiamuniverse.EventUserCreated,
			goiamuniverse.EventRoleUpdated,
			goiamuniverse.EventResourceCreated,
			goiamuniverse.EventClientCreated,
		}

		for _, eventType := range events {
			t.Run(string(eventType), func(t *testing.T) {
				event := newEvent(ctx, eventType, *testUser, testMetadata)
				assert.NotPanics(t, func() {
					svc.Emit(event)
				})
			})
		}
	})
}

// MockSubscriber implements the Subscriber interface for testing
type MockSubscriber struct {
	events []utils.Event[sdk.User]
}

func (m *MockSubscriber) HandleEvent(event utils.Event[sdk.User]) {
	m.events = append(m.events, event)
}

// TestSubscribe tests the Subscribe method
func TestSubscribe(t *testing.T) {
	svc, _, _ := setupUserService()

	t.Run("success - subscribe to user created event", func(t *testing.T) {
		mockSubscriber := &MockSubscriber{}
		assert.NotPanics(t, func() {
			svc.Subscribe(goiamuniverse.EventUserCreated, mockSubscriber)
		})
	})

	t.Run("success - subscribe to different events", func(t *testing.T) {
		events := []goiamuniverse.Event{
			goiamuniverse.EventUserCreated,
			goiamuniverse.EventRoleUpdated,
			goiamuniverse.EventResourceCreated,
			goiamuniverse.EventClientCreated,
		}

		for _, eventType := range events {
			t.Run(string(eventType), func(t *testing.T) {
				mockSubscriber := &MockSubscriber{}
				assert.NotPanics(t, func() {
					svc.Subscribe(eventType, mockSubscriber)
				})
			})
		}
	})
}

// TestEmitAndSubscribeIntegration tests the integration between Emit and Subscribe
func TestEmitAndSubscribeIntegration(t *testing.T) {
	svc, _, _ := setupUserService()
	ctx := createContextWithMetadata()
	testUser := createTestUser()
	testMetadata := middlewares.GetMetadata(ctx)

	t.Run("integration - subscribe and emit", func(t *testing.T) {
		mockSubscriber := &MockSubscriber{}

		// Subscribe to events
		svc.Subscribe(goiamuniverse.EventUserCreated, mockSubscriber)

		// Emit an event
		event := newEvent(ctx, goiamuniverse.EventUserCreated, *testUser, testMetadata)
		svc.Emit(event)

		// Verify the subscription and emission worked without panicking
		// Note: The actual delivery depends on the emitter implementation
		// We mainly test that the methods work correctly
		assert.NotNil(t, mockSubscriber)
		assert.NotNil(t, event)
	})
}
