package system

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
)

func TestNewAddResourcesToUser(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewAddResourcesToUser(userSvc)

	assert.Equal(t, "@policy/system/add_resources_to_user", policy.ID())
	assert.Equal(t, "Add resources to user specified in user policy", policy.Name())
	assert.NotNil(t, policy.userSvc)
	assert.NotNil(t, policy.pc)
}

func TestAddResourcesToUser_HandleEvent_Success(t *testing.T) {
	userSvc := &services.MockUserService{}

	ctx := context.Background()
	userId := "user123"
	targetUserId := "target456"
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
			"@policy/system/add_resources_to_user": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{
						"@userId": {Static: targetUserId},
					},
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)
	userSvc.On("AddResourceToUser", ctx, targetUserId, sdk.AddUserResourceRequest{
		PolicyId: "@policy/system/add_resources_to_user",
		Key:      resource.Key,
		Name:     resource.Name,
	}).Return(nil)

	// Execute
	policy := NewAddResourcesToUser(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}

func TestAddResourcesToUser_HandleEvent_PolicyCheckError(t *testing.T) {
	userSvc := &services.MockUserService{}

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
	userSvc.On("GetById", ctx, userId).Return(&sdk.User{}, errors.New("policy check error"))

	// Execute
	policy := NewAddResourcesToUser(userSvc)
	policy.HandleEvent(event)

	// Verify - no AddResourceToUser calls should be made due to error
	userSvc.AssertExpectations(t)
	userSvc.AssertNotCalled(t, "AddResourceToUser")
}

func TestAddResourcesToUser_HandleEvent_PolicyNotExists(t *testing.T) {
	userSvc := &services.MockUserService{}

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
	policy := NewAddResourcesToUser(userSvc)
	policy.HandleEvent(event)

	// Verify - no AddResourceToUser calls should be made when policy doesn't exist
	userSvc.AssertExpectations(t)
	userSvc.AssertNotCalled(t, "AddResourceToUser")
}

func TestAddResourcesToUser_HandleEvent_NoUserIdInPolicy(t *testing.T) {
	userSvc := &services.MockUserService{}

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

	// Mock the user service for policy check - policy exists but no @userId argument
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/add_resources_to_user": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{}, // No @userId
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	policy := NewAddResourcesToUser(userSvc)
	policy.HandleEvent(event)

	// Verify - no AddResourceToUser calls should be made when userId is missing
	userSvc.AssertExpectations(t)
	userSvc.AssertNotCalled(t, "AddResourceToUser")
}

func TestAddResourcesToUser_HandleEvent_EmptyUserId(t *testing.T) {
	userSvc := &services.MockUserService{}

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

	// Mock the user service for policy check - policy exists but empty @userId
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/add_resources_to_user": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{
						"@userId": {Static: ""}, // Empty user ID
					},
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)

	// Execute
	policy := NewAddResourcesToUser(userSvc)
	policy.HandleEvent(event)

	// Verify - no AddResourceToUser calls should be made when userId is empty
	userSvc.AssertExpectations(t)
	userSvc.AssertNotCalled(t, "AddResourceToUser")
}

func TestAddResourcesToUser_HandleEvent_AddResourceError(t *testing.T) {
	userSvc := &services.MockUserService{}

	ctx := context.Background()
	userId := "user123"
	targetUserId := "target456"
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
			"@policy/system/add_resources_to_user": {
				Name: "Add Resources Policy",
				Mapping: sdk.UserPolicyMapping{
					Arguments: map[string]sdk.UserPolicyMappingValue{
						"@userId": {Static: targetUserId},
					},
				},
			},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)
	userSvc.On("AddResourceToUser", ctx, targetUserId, sdk.AddUserResourceRequest{
		PolicyId: "@policy/system/add_resources_to_user",
		Key:      resource.Key,
		Name:     resource.Name,
	}).Return(errors.New("add resource error"))

	// Execute
	policy := NewAddResourcesToUser(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}

func TestAddResourcesToUser_getTargetUserId(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewAddResourcesToUser(userSvc)

	tests := []struct {
		name           string
		user           *sdk.User
		expectedUserId string
		expectedOk     bool
	}{
		{
			name: "valid_user_id",
			user: &sdk.User{
				Id: "user123",
				Policies: map[string]sdk.UserPolicy{
					"@policy/system/add_resources_to_user": {
						Name: "Add Resources Policy",
						Mapping: sdk.UserPolicyMapping{
							Arguments: map[string]sdk.UserPolicyMappingValue{
								"@userId": {Static: "target456"},
							},
						},
					},
				},
			},
			expectedUserId: "target456",
			expectedOk:     true,
		},
		{
			name: "policy_not_found",
			user: &sdk.User{
				Id:       "user123",
				Policies: map[string]sdk.UserPolicy{},
			},
			expectedUserId: "",
			expectedOk:     false,
		},
		{
			name: "user_id_argument_not_found",
			user: &sdk.User{
				Id: "user123",
				Policies: map[string]sdk.UserPolicy{
					"@policy/system/add_resources_to_user": {
						Name: "Add Resources Policy",
						Mapping: sdk.UserPolicyMapping{
							Arguments: map[string]sdk.UserPolicyMappingValue{},
						},
					},
				},
			},
			expectedUserId: "",
			expectedOk:     false,
		},
		{
			name: "empty_user_id",
			user: &sdk.User{
				Id: "user123",
				Policies: map[string]sdk.UserPolicy{
					"@policy/system/add_resources_to_user": {
						Name: "Add Resources Policy",
						Mapping: sdk.UserPolicyMapping{
							Arguments: map[string]sdk.UserPolicyMappingValue{
								"@userId": {Static: ""},
							},
						},
					},
				},
			},
			expectedUserId: "",
			expectedOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId, ok := policy.getTargetUserId(tt.user)
			assert.Equal(t, tt.expectedUserId, userId)
			assert.Equal(t, tt.expectedOk, ok)
		})
	}
}

func TestAddResourcesToUser_PolicyDef(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewAddResourcesToUser(userSvc)

	policyDef := policy.PolicyDef()

	assert.Equal(t, "@policy/system/add_resources_to_user", policyDef.Id)
	assert.Equal(t, "Add resources to user specified in user policy", policyDef.Name)
	assert.Equal(t, "This policy adds the created resource to the user specified in the user policy.", policyDef.Description)
	assert.Len(t, policyDef.Definition.Arguments, 1)
	assert.Equal(t, "@userId", policyDef.Definition.Arguments[0].Name)
	assert.Equal(t, "The user to whom the resource access is granted.", policyDef.Definition.Arguments[0].Description)
	assert.Equal(t, goiamuniverse.User, policyDef.Definition.Arguments[0].DataType)
}
