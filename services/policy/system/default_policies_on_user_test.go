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

// mockUserEvent implements utils.Event[sdk.User] for testing
type mockUserEvent struct {
	name     goiamuniverse.Event
	payload  sdk.User
	metadata sdk.Metadata
	ctx      context.Context
}

func (e mockUserEvent) Name() goiamuniverse.Event {
	return e.name
}

func (e mockUserEvent) Payload() sdk.User {
	return e.payload
}

func (e mockUserEvent) Metadata() sdk.Metadata {
	return e.metadata
}

func (e mockUserEvent) Context() context.Context {
	return e.ctx
}

func newMockUserEvent(ctx context.Context, name goiamuniverse.Event, payload sdk.User, metadata sdk.Metadata) mockUserEvent {
	return mockUserEvent{ctx: ctx, name: name, payload: payload, metadata: metadata}
}

func TestNewDefaultPoliciesOnUser(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewDefaultPoliciesOnUser(userSvc)

	assert.Equal(t, "@policy/system/default_policies_on_user", policy.ID())
	assert.NotNil(t, policy.userSvc)
}

func TestDefaultPoliciesOnUser_HandleEvent_Success(t *testing.T) {
	userSvc := &services.MockUserService{}

	ctx := context.Background()
	userId := "user123"
	user := sdk.User{
		Id:    userId,
		Email: "test@example.com",
		Name:  "Test User",
	}

	event := newMockUserEvent(
		ctx,
		goiamuniverse.EventUserCreated,
		user,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Mock the user service
	expectedPolicies := map[string]sdk.UserPolicy{
		"@policy/system/access_to_created_resource": {
			Name: "User get access to the resource created by the user",
		},
	}
	userSvc.On("AddPolicyToUser", ctx, userId, expectedPolicies).Return(nil)

	// Execute
	policy := NewDefaultPoliciesOnUser(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}

func TestDefaultPoliciesOnUser_HandleEvent_AddPolicyError(t *testing.T) {
	userSvc := &services.MockUserService{}

	ctx := context.Background()
	userId := "user123"
	user := sdk.User{
		Id:    userId,
		Email: "test@example.com",
		Name:  "Test User",
	}

	event := newMockUserEvent(
		ctx,
		goiamuniverse.EventUserCreated,
		user,
		sdk.Metadata{User: &sdk.User{Id: userId}},
	)

	// Mock the user service to return error
	expectedPolicies := map[string]sdk.UserPolicy{
		"@policy/system/access_to_created_resource": {
			Name: "User get access to the resource created by the user",
		},
	}
	userSvc.On("AddPolicyToUser", ctx, userId, expectedPolicies).Return(errors.New("add policy error"))

	// Execute
	policy := NewDefaultPoliciesOnUser(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}

func TestDefaultPoliciesOnUser_HandleEvent_DifferentUser(t *testing.T) {
	userSvc := &services.MockUserService{}

	ctx := context.Background()
	originalUserId := "original123"
	targetUserId := "target456"
	user := sdk.User{
		Id:    targetUserId,
		Email: "target@example.com",
		Name:  "Target User",
	}

	event := newMockUserEvent(
		ctx,
		goiamuniverse.EventUserCreated,
		user,
		sdk.Metadata{User: &sdk.User{Id: originalUserId}},
	)

	// Mock the user service - should use payload user ID, not metadata user ID
	expectedPolicies := map[string]sdk.UserPolicy{
		"@policy/system/access_to_created_resource": {
			Name: "User get access to the resource created by the user",
		},
	}
	userSvc.On("AddPolicyToUser", ctx, targetUserId, expectedPolicies).Return(nil)

	// Execute
	policy := NewDefaultPoliciesOnUser(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}
