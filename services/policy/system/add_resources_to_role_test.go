package system

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRoleService implements role.Service interface for testing
type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) Create(ctx context.Context, role *sdk.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleService) Update(ctx context.Context, role *sdk.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleService) GetById(ctx context.Context, id string) (*sdk.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.Role), args.Error(1)
}

func (m *MockRoleService) GetAll(ctx context.Context, query sdk.RoleQuery) (*sdk.RoleList, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sdk.RoleList), args.Error(1)
}

func (m *MockRoleService) AddResource(ctx context.Context, roleId string, resource sdk.Resources) error {
	args := m.Called(ctx, roleId, resource)
	return args.Error(0)
}

func (m *MockRoleService) RemoveResource(ctx context.Context, roleId string, resourceId string) error {
	args := m.Called(ctx, roleId, resourceId)
	return args.Error(0)
}

func (m *MockRoleService) Emit(event utils.Event[sdk.Role]) {
	m.Called(event)
}

func (m *MockRoleService) Subscribe(eventName goiamuniverse.Event, subscriber utils.Subscriber[utils.Event[sdk.Role], sdk.Role]) {
	m.Called(eventName, subscriber)
}

func TestNewAddResourcesToRole(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}
	policy := NewAddResourcesToRole(userSvc, roleSvc)

	assert.Equal(t, "@policy/system/add_resources_to_role", policy.ID())
	assert.Equal(t, "Add resources to role specified in user policy", policy.Name())
	assert.NotNil(t, policy.userSvc)
	assert.NotNil(t, policy.roleSvc)
	assert.NotNil(t, policy.pc)
}

func TestAddResourcesToRole_HandleEvent_Success(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}

	ctx := context.Background()
	userId := "user123"
	roleId := "role456"
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceCreated,
		resource,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Mock the user service for policy check
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/add_resources_to_role": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{
						"@roleId": {Static: roleId},
					},
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)
	roleSvc.On("AddResource", ctx, roleId, sdk.Resources{
		Id:   resource.ID,
		Key:  resource.Key,
		Name: resource.Name,
	}).Return(nil)

	// Execute
	policy := NewAddResourcesToRole(userSvc, roleSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
	roleSvc.AssertExpectations(t)
}

func TestAddResourcesToRole_HandleEvent_PolicyCheckError(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}

	ctx := context.Background()
	userId := "user123"
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceCreated,
		resource,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Setup mocks - policy check returns error
	userSvc.On("GetById", ctx, userId).Return(nil, errors.New("policy check error"))

	// Execute
	policy := NewAddResourcesToRole(userSvc, roleSvc)
	policy.HandleEvent(event)

	// Verify - no role service calls should be made due to error
	userSvc.AssertExpectations(t)
	roleSvc.AssertNotCalled(t, "AddResource")
}

func TestAddResourcesToRole_HandleEvent_PolicyNotExists(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}

	ctx := context.Background()
	userId := "user123"
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceCreated,
		resource,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Setup mocks - user exists but policy doesn't exist for user
	testUser := &sdk.User{
		Id:       userId,
		Policies: map[string]sdk.UserPolicy{}, // Empty policies
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	policy := NewAddResourcesToRole(userSvc, roleSvc)
	policy.HandleEvent(event)

	// Verify - no role service calls should be made when policy doesn't exist
	userSvc.AssertExpectations(t)
	roleSvc.AssertNotCalled(t, "AddResource")
}

func TestAddResourcesToRole_HandleEvent_NoRoleIdInPolicy(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}

	ctx := context.Background()
	userId := "user123"
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceCreated,
		resource,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Mock the user service for policy check - policy exists but no @roleId argument
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/add_resources_to_role": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{}, // No @roleId
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	policy := NewAddResourcesToRole(userSvc, roleSvc)
	policy.HandleEvent(event)

	// Verify - no role service calls should be made when roleId is missing
	userSvc.AssertExpectations(t)
	roleSvc.AssertNotCalled(t, "AddResource")
}

func TestAddResourcesToRole_HandleEvent_EmptyRoleId(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}

	ctx := context.Background()
	userId := "user123"
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceCreated,
		resource,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Mock the user service for policy check - policy exists but empty @roleId
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/add_resources_to_role": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{
						"@roleId": {Static: ""}, // Empty role ID
					},
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	policy := NewAddResourcesToRole(userSvc, roleSvc)
	policy.HandleEvent(event)

	// Verify - no role service calls should be made when roleId is empty
	userSvc.AssertExpectations(t)
	roleSvc.AssertNotCalled(t, "AddResource")
}

func TestAddResourcesToRole_HandleEvent_AddResourceError(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}

	ctx := context.Background()
	userId := "user123"
	roleId := "role456"
	resource := sdk.Resource{
		ID:   "resource123",
		Key:  "test-resource",
		Name: "Test Resource",
	}

	event := newMockEvent(
		ctx,
		goiamuniverse.EventResourceCreated,
		resource,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Mock the user service for policy check
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/add_resources_to_role": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{
						"@roleId": {Static: roleId},
					},
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)
	roleSvc.On("AddResource", ctx, roleId, sdk.Resources{
		Id:   resource.ID,
		Key:  resource.Key,
		Name: resource.Name,
	}).Return(errors.New("add resource error"))

	// Execute
	policy := NewAddResourcesToRole(userSvc, roleSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
	roleSvc.AssertExpectations(t)
}

func TestAddResourcesToRole_getTargetRoleId(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}
	policy := NewAddResourcesToRole(userSvc, roleSvc)

	tests := []struct {
		name           string
		user           *sdk.User
		expectedRoleId string
		expectedOk     bool
	}{
		{
			name: "valid_role_id",
			user: &sdk.User{
				Id: "user123",
				Policies: map[string]sdk.UserPolicy{
					"@policy/system/add_resources_to_role": {
						Name: "Add Resources Policy",
						Mapping: sdk.UserPolicyMapping{
							Arguments: map[string]sdk.UserPolicyMappingValue{
								"@roleId": {Static: "role456"},
							},
						},
					},
				},
			},
			expectedRoleId: "role456",
			expectedOk:     true,
		},
		{
			name: "policy_not_found",
			user: &sdk.User{
				Id:       "user123",
				Policies: map[string]sdk.UserPolicy{},
			},
			expectedRoleId: "",
			expectedOk:     false,
		},
		{
			name: "role_id_argument_not_found",
			user: &sdk.User{
				Id: "user123",
				Policies: map[string]sdk.UserPolicy{
					"@policy/system/add_resources_to_role": {
						Name: "Add Resources Policy",
						Mapping: sdk.UserPolicyMapping{
							Arguments: map[string]sdk.UserPolicyMappingValue{},
						},
					},
				},
			},
			expectedRoleId: "",
			expectedOk:     false,
		},
		{
			name: "empty_role_id",
			user: &sdk.User{
				Id: "user123",
				Policies: map[string]sdk.UserPolicy{
					"@policy/system/add_resources_to_role": {
						Name: "Add Resources Policy",
						Mapping: sdk.UserPolicyMapping{
							Arguments: map[string]sdk.UserPolicyMappingValue{
								"@roleId": {Static: ""},
							},
						},
					},
				},
			},
			expectedRoleId: "",
			expectedOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleId, ok := policy.getTargetRoleId(tt.user)
			assert.Equal(t, tt.expectedRoleId, roleId)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}

func TestAddResourcesToRole_PolicyDef(t *testing.T) {
	userSvc := &MockUserService{}
	roleSvc := &MockRoleService{}
	policy := NewAddResourcesToRole(userSvc, roleSvc)

	policyDef := policy.PolicyDef()

	assert.Equal(t, "@policy/system/add_resources_to_role", policyDef.Id)
	assert.Equal(t, "Add resources to role specified in user policy", policyDef.Name)
	assert.Equal(t, "This policy adds the created resource to the role specified in the user policy.", policyDef.Description)
	assert.Len(t, policyDef.Definition.Arguments, 1)
	assert.Equal(t, "@roleId", policyDef.Definition.Arguments[0].Name)
	assert.Equal(t, "The role to which the resource access is granted.", policyDef.Definition.Arguments[0].Description)
	assert.Equal(t, goiamuniverse.Role, policyDef.Definition.Arguments[0].DataType)
}
