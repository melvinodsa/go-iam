package system

import (
	"context"
	"errors"
	"testing"

	"github.com/melvinodsa/go-iam/sdk"
	"github.com/melvinodsa/go-iam/utils"
	"github.com/melvinodsa/go-iam/utils/goiamuniverse"
	"github.com/melvinodsa/go-iam/utils/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPolicyCheck for testing
type MockPolicyCheck struct {
	mock.Mock
}

func (m *MockPolicyCheck) RunCheck(ctx context.Context, id, userId string) (*sdk.User, bool, error) {
	args := m.Called(ctx, id, userId)
	if args.Get(0) == nil {
		return nil, args.Bool(1), args.Error(2)
	}
	return args.Get(0).(*sdk.User), args.Bool(1), args.Error(2)
}

// mockEvent implements utils.Event[sdk.Resource] for testing
type mockEvent struct {
	name     goiamuniverse.Event
	payload  sdk.Resource
	metadata sdk.Metadata
	ctx      context.Context
}

func (e mockEvent) Name() goiamuniverse.Event {
	return e.name
}

func (e mockEvent) Payload() sdk.Resource {
	return e.payload
}

func (e mockEvent) Metadata() sdk.Metadata {
	return e.metadata
}

func (e mockEvent) Context() context.Context {
	return e.ctx
}

func newMockEvent(ctx context.Context, name goiamuniverse.Event, payload sdk.Resource, metadata sdk.Metadata) utils.Event[sdk.Resource] {
	return mockEvent{ctx: ctx, name: name, payload: payload, metadata: metadata}
}

func TestNewAccessToCreatedResource(t *testing.T) {
	userSvc := &services.MockUserService{}
	policy := NewAccessToCreatedResource(userSvc)

	assert.Equal(t, "@policy/system/access_to_created_resource", policy.ID())
	assert.Equal(t, "User get access to the resource created by the user", policy.Name())
	assert.NotNil(t, policy.userSvc)
	assert.NotNil(t, policy.pc)
}

func TestAccessToCreatedResource_HandleEvent_Success(t *testing.T) {
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

	// Mock the user service for policy check
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/access_to_created_resource": {Name: "Access Policy"},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)
	userSvc.On("AddResourceToUser", ctx, userId, sdk.AddUserResourceRequest{
		PolicyId: "@policy/system/access_to_created_resource",
		Key:      resource.Key,
		Name:     resource.Name,
	}).Return(nil)

	// Execute
	policy := NewAccessToCreatedResource(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}

func TestAccessToCreatedResource_HandleEvent_PolicyCheckError(t *testing.T) {
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
	policy := NewAccessToCreatedResource(userSvc)
	policy.HandleEvent(event)

	// Verify - no AddResourceToUser calls should be made due to error
	userSvc.AssertExpectations(t)
	userSvc.AssertNotCalled(t, "AddResourceToUser")
}

func TestAccessToCreatedResource_HandleEvent_PolicyNotExists(t *testing.T) {
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
	policy := NewAccessToCreatedResource(userSvc)
	policy.HandleEvent(event)

	// Verify - no AddResourceToUser calls should be made when policy doesn't exist
	userSvc.AssertExpectations(t)
	userSvc.AssertNotCalled(t, "AddResourceToUser")
}

func TestAccessToCreatedResource_HandleEvent_AddResourceError(t *testing.T) {
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

	// Mock the user service for policy check
	testUser := &sdk.User{
		Id: userId,
		Policies: map[string]sdk.UserPolicy{
			"@policy/system/access_to_created_resource": {Name: "Access Policy"},
		},
	}
	userSvc.On("GetById", ctx, userId).Return(testUser, nil)
	userSvc.On("AddResourceToUser", ctx, userId, sdk.AddUserResourceRequest{
		PolicyId: "@policy/system/access_to_created_resource",
		Key:      resource.Key,
		Name:     resource.Name,
	}).Return(errors.New("add resource error"))

	// Execute
	policy := NewAccessToCreatedResource(userSvc)
	policy.HandleEvent(event)

	// Verify
	userSvc.AssertExpectations(t)
}
